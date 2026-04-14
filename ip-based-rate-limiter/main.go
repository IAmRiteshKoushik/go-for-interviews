package main

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

type RateLimiter struct {
	mu         sync.Mutex
	tokens     float64
	capacity   float64
	refillRate float64
	lastRefill time.Time
}

func NewRateLimiter(capacity int, refillPerSec float64) *RateLimiter {
	return &RateLimiter{
		tokens:     float64(capacity),
		capacity:   float64(capacity),
		refillRate: refillPerSec,
		lastRefill: time.Now(),
	}
}

func (rl *RateLimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(rl.lastRefill).Seconds()

	rl.tokens += elapsed * rl.refillRate
	if rl.tokens > rl.capacity {
		rl.tokens = rl.capacity
	}
	rl.lastRefill = now

	if rl.tokens >= 1 {
		rl.tokens--
		return true
	}
	return false
}

type clientEntry struct {
	limiter  *RateLimiter
	lastSeen time.Time
}

type IPLimiter struct {
	mu      sync.Mutex
	clients map[string]*clientEntry

	capacity int
	refill   float64
	ttl      time.Duration
}

func NewIPLimiter(capacity int, refill float64, ttl time.Duration) *IPLimiter {
	l := &IPLimiter{
		clients:  make(map[string]*clientEntry),
		capacity: capacity,
		refill:   refill,
		ttl:      ttl,
	}

	go l.cleanLoop()
	return l
}

func (l *IPLimiter) Allow(ip string) bool {
	now := time.Now()

	l.mu.Lock()

	entry, ok := l.clients[ip]
	if !ok {
		entry = &clientEntry{
			limiter:  NewRateLimiter(l.capacity, l.refill),
			lastSeen: now,
		}
		l.clients[ip] = entry
	}
	// If entry exists
	entry.lastSeen = now
	l.mu.Unlock()

	return entry.limiter.Allow()
}

func (l *IPLimiter) cleanLoop() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		cutoff := time.Now().Add(-l.ttl)
		l.mu.Lock()
		for ip, entry := range l.clients {
			if entry.lastSeen.Before(cutoff) {
				delete(l.clients, ip)
			}
		}
		l.mu.Unlock()
	}
}

func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

func main() {
	// Handle a burst of 20, refill 10 req/s, evict idle entries
	ipLimiter := NewIPLimiter(20, 10, 10*time.Minute)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ip := clientIP(r)
		if !ipLimiter.Allow(ip) {
			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		fmt.Fprintln(w, "ok")
	})

	fmt.Println("listening on :8080")
	_ = http.ListenAndServe(":8080", nil)
}
