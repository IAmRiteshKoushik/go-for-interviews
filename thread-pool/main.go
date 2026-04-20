package main

import (
	"fmt"
	"sync"
	"time"
)

type Job func()

type Pool struct {
	workQueue chan Job
	wg        sync.WaitGroup
}

func NewPool(workerCount int) *Pool {
	pool := &Pool{
		workQueue: make(chan Job),
	}
	pool.wg.Add(workerCount)
	for range workerCount {
		go func() {
			defer pool.wg.Done()

			// Whenever there is a job added to the channel, one of the threads is
			// picking them up and executing. This is a blocking call as we are
			// dependent on the channel to give us new jobs
			for job := range pool.workQueue {
				job()
			}
		}()
	}
	return pool
}

func (p *Pool) AddJob(job Job) {
	// Introducing a new job in the pool
	p.workQueue <- job
}

func (p *Pool) Wait() {
	p.wg.Wait()
	close(p.workQueue) // channels are blocking bounded buffers
}

func main() {
	pool := NewPool(5)

	for range 30 {
		job := func() {
			time.Sleep(1 * time.Second)
			fmt.Println("Job: Completed")
		}
		pool.AddJob(job)
	}

	pool.Wait()
}
