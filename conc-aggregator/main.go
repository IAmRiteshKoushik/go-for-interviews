package main

import (
	"fmt"
	"sync"
	"time"
)

func fetchUser() string {
	time.Sleep(time.Millisecond * 100)
	return "BOB"
}

func fetchUserLikes(_ string, respch chan any, wg *sync.WaitGroup) {
	time.Sleep(time.Millisecond * 150)
	respch <- 200
	wg.Done()
}

func fetchUserMatch(_ string, respch chan any, wg *sync.WaitGroup) {
	time.Sleep(time.Millisecond * 100)
	respch <- "SALLY"
	wg.Done()
}

func main() {
	start := time.Now()
	userName := fetchUser() // 100ms
	respch := make(chan any, 2)
	wg := &sync.WaitGroup{}

	wg.Add(2)

	// likes := fetchUserLikes(userName)
	// match := fetchUserMatch(userName)

	go fetchUserLikes(userName, respch, wg) // 150ms
	go fetchUserMatch(userName, respch, wg) // 100ms

	// We need to close the channel otherwise we are
	// going to get a fatal error. But the problem is
	// that the "go" keyword offloads tasks to separate
	// goroutines (threads). And we are not sure if the
	// threads have completed their respective jobs and
	// returned the values which we are iterating over as
	// a range in the following for loop

	wg.Wait() // block until 2 wg.Done() calls.
	close(respch)

	// Adding the wg.Wait() statement, has blocked the program and does not
	// allow early closing of the channels before the job is completed. After
	// each job is completed, the wait group counter is brought down and only
	// then do close the channels.

	// fmt.Println("Likes: ", likes)
	// fmt.Println("Likes: ", match)

	// In Golang, you can range over channels
	for resp := range respch {
		fmt.Println("Resp: ", resp)
	}
	fmt.Println("Took us ", time.Since(start))
}
