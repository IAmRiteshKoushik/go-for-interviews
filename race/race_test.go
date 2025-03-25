package main

import "testing"

// Simulating race conditions
func TestDataRaceCondition(t *testing.T) {
	var state int32

	for i := range 10 {
		go func(i int) {
			state += int32(i)
		}(i)
	}
}
