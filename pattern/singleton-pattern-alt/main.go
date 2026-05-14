package main

import "fmt"

func main() {
	for range 30 {
		go getInstance()
	}

	fmt.Scanln()
}
