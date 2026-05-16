package main

import "fmt"

type MiddleCoordinates struct {
	x, y int
}

func (a *MiddleCoordinates) visitForSquare(s *Square) {
	fmt.Println("Calculating mid-pt for square")
}

func (a *MiddleCoordinates) visitForCircle(s *Circle) {
	fmt.Println("Calculating mid-pt for circle")
}

func (a *MiddleCoordinates) visitForRectangle(s *Rectangle) {
	fmt.Println("Calculating mid-pt for rectangle")
}
