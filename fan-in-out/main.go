package main

import (
	"fmt"
	"sync"
)

// Generator produces input values
func generator(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, n := range nums {
			out <- n
		}
	}()
	return out
}

// Workers ready from in, do processing and write to out
func worker(id int, in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			// processing step
			out <- n * n
		}
	}()
	return out
}

// Fan-in merges channels into one
func fanIn(chans ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	wg.Add(len(chans))
	for _, ch := range chans {
		c := ch
		go func() {
			defer wg.Done()
			for v := range c {
				out <- v
			}
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func main() {
	in := generator(1, 2, 3, 4, 5, 6, 7, 8)

	// Fan-out: 3 workers consume same input channel
	w1 := worker(1, in)
	w2 := worker(2, in)
	w3 := worker(3, in)

	// Merge all worker outputs
	out := fanIn(w1, w2, w3)
	for v := range out {
		fmt.Println(v)
	}
}
