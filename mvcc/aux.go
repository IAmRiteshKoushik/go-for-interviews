package main

import (
	"fmt"
	"os"
	"slices"
)

func assert(b bool, msg string) {
	if !b {
		panic(msg)
	}
}

func assertEQ[C comparable](a C, b C, prefix string) {
	if a != b {
		panic(fmt.Sprintf("%s %v != %v", prefix, a, b))
	}
}

var DEBUG = slices.Contains(os.Args, "--debug")

func debug(a ...any) {
	if !DEBUG {
		return
	}

	args := append([]any{"[DEBUG]"}, a...)
	fmt.Println(args...)
}
