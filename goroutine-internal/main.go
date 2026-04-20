package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

// It is trivial to spin up 100_000 goroutines
func demoGorutineCheapness() {
	const n = 100_000
	var wg sync.WaitGroup
	wg.Add(n)

	start := time.Now()
	for range n {
		go func() {
			defer wg.Done()
			// Some small work
			_ = compute(10)
		}()
	}
	wg.Wait()
	fmt.Printf("Spawned and joined %d goroutines in %v\n", n, time.Since(start))
}

// Work: Sum of n natural numbers. Stack growth using recursive func call
func compute(depth int) int {
	if depth == 0 {
		return 1
	}
	return depth + compute(depth-1)
}

// shows the P (logical processor) model - GOMAXPROCS controls parallelism
func demoWorkStealing() {
	// With GOMAXPROCS = 1, all goroutines share one P - purely concurrent
	// With OGMAXPROCS = N, the scheduler work-steals across N Ps - paralley
	procs := runtime.GOMAXPROCS(0) // 0 = query current value
	fmt.Printf("Running in GOMAXPROCS=%d\n", procs)

	var wg sync.WaitGroup
	results := make([]int, procs)

	for i := range procs {
		wg.Add(i)
		go func(id int) {
			defer wg.Done()
			sum := 0
			for j := range 1_000_000 {
				sum += j
			}
			results[id] = sum
		}(i)
	}

	wg.Wait()
	fmt.Println("All Ps finished their work-stolen queues:", results[0])
}

func main() {
	demoGorutineCheapness()
	demoWorkStealing()
}
