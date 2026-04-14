package main

import (
	"fmt"
	"sync"
	"time"
)

const (
	numWorkers = 3
	numJobs    = 30
)

type Job struct {
	ID int
	N  int
}

type Result struct {
	JobID int
	Value int
}

func worker(id int, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()
	// Sender channel
	for job := range jobs {
		// simulate work done
		time.Sleep(1*time.Second + 500*time.Millisecond)
		results <- Result{
			JobID: job.ID,
			Value: job.N * job.N,
		}
		// Receiver channel is 'results'
	}
}

func main() {
	jobs := make(chan Job, numJobs)
	results := make(chan Result, numJobs)

	// Spawn N workers (fixed size)
	var wg sync.WaitGroup
	for w := range numWorkers {
		wg.Add(1)
		go worker(w, jobs, results, &wg)
	}

	// Submit jobs to the channel
	for i := range numJobs {
		jobs <- Job{ID: i, N: i}
	}
	close(jobs)

	// Async close after all workers finish publishing
	go func() {
		wg.Wait()
		close(results)
	}()

	// Consume results
	for r := range results {
		fmt.Printf("job=%d result=%d\n", r.JobID, r.Value)
	}
}
