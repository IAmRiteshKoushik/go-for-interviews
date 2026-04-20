package main

import "sync"

type Semaphore struct {
	value   int
	mutex   *sync.Mutex
	waiters []chan struct{}
}

// Create a new semaphore with an initial value
func NewSemaphore(value int) *Semaphore {
	return &Semaphore{
		value: value,
		mutex: &sync.Mutex{},
	}
}

func (sem *Semaphore) Wait() {
	sem.mutex.Lock()

	sem.value--
	if sem.value >= 0 {
	}
}

func (sem *Semaphore) Signal() {}
