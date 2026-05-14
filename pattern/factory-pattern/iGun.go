package main

type IGun interface {
	setName(name string)
	getName() string
	setPower(power int)
	getPower() int
}
