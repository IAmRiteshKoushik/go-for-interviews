package main

import "sync/atomic"

// 70% of performance bottleneck in highly concurrent Go applications trace
// back to lock contention, not algorithmic complexity. The lcks we add
// to protect data become the very choking point.

// Compare and Swap operations (CAS) - Persistent Optimism

type node[T any] struct {
	value T
	next  *node[T]
}

type Stack[T any] struct {
	head atomic.Pointer[node[T]]
}

func (s *Stack[T]) Push(v T) {

}

func (s *Stack[T]) Pop() (T, bool) {
	var zero T

	return zero, true
}
