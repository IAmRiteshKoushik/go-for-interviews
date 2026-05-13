package main

type IsolationLevel uint8

// Loosest to strictest (top to bottom)
const (
	ReadUncommittedIsolation IsolationLevel = iota
	ReadCommittedIsolation
	RepeatableReadIsolation
	SnapshotIsolation
	SerializableIsolation
)
