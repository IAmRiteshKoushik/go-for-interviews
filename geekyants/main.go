package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	oddChan := make(chan struct{}, 1)
	evenChan := make(chan struct{}, 1)

	wg.Add(1)
	go func() {
		defer wg.Done()

		for i := range 10 {
			if i%2 == 0 {
				<-evenChan
				fmt.Println(i)
				oddChan <- struct{}{}
			} else {
				<-oddChan
				fmt.Println(i)
				// There is no receiver after i == 9, so it causes a deadlock
				if i < 9 {
					evenChan <- struct{}{}
				}
			}
		}
	}()

	// init trigger
	evenChan <- struct{}{} // start with even = 0
	wg.Wait()
	close(oddChan)
	close(evenChan)
}
