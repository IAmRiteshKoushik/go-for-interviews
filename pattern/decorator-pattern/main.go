package main

import "fmt"

func main() {
	pizza := &VeggieMania{}

	// Add cheese toppings
	pizzaWithCheese := &CheeseTopping{
		pizza: pizza,
	}
	// Add tomato toppings
	pizzaWithCheeseAndTomato := &TomatoTopping{
		pizza: pizzaWithCheese,
	}

	fmt.Printf("Price of veggieMania with tomato and cheese: %d\n", pizzaWithCheeseAndTomato.getPrice())
}
