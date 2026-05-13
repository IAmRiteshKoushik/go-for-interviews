package main

import (
	"github.com/tidwall/btree"
)

type Value struct {
	txStartId uint64
	txEndId   uint64
	value     string
}

type Transaction struct {
	isolation IsolationLevel
	id        uint64
	state     TransactionSate

	// Used only by repeatable read and stricter
	inprogress btree.Set[uint64]

	// Used only by snapshot isolation and stricter
	writeset btree.Set[string]
	readset  btree.Set[string]
}
