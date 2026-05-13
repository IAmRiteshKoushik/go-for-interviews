package main

type TransactionSate uint8

const (
	InProgressTransaction TransactionSate = iota
	AbortedTransaction
	CommittedTransaction
)
