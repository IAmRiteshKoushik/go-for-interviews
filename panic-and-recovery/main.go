package main

import "fmt"

func divide(a, b int) {

	// Has to run in the same go-routine
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("panic recovered:", err)
		}
	}()

	if b == 0 {
		panic(nil)
	}
	fmt.Println("result:", a/b)
}

func main() {
	divide(10, 2)
	divide(10, 0)
	fmt.Println("End of main")
}
