package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

type RateLimiter struct {
	tokens     float64
	capacity   float64
	refillRate float64
	lastRefill time.Time

	mu sync.Mutex
}

func NewRateLimiter(capacity int, refillPerSec float64) *RateLimiter {
	return &RateLimiter{
		tokens:     float64(capacity),
		capacity:   float64(capacity),
		refillRate: refillPerSec,
		lastRefill: time.Now(),
	}
}

// Token bucket
func (rl *RateLimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(rl.lastRefill).Seconds()

	rl.tokens += elapsed * rl.refillRate
	if rl.tokens > rl.capacity {
		rl.tokens = rl.capacity // avoiding overflows
	}
	rl.lastRefill = now

	// Consume a token
	if rl.tokens >= 1 {
		rl.tokens--
		return true
	}
	return false
}

func main() {
	limiter := NewRateLimiter(10, 5)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		fmt.Fprintln(w, "ok")
	})

	fmt.Println("Listening on :8080")
	_ = http.ListenAndServe(":8080", nil)
}
