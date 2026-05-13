package main

import "github.com/tidwall/btree"

type Database struct {
	defaultIsolation  IsolationLevel
	store             map[string][]Value
	transactions      btree.Map[uint64, Transaction]
	nextTransactionId uint64
}

func newDatabase() Database {
	return Database{
		defaultIsolation: ReadCommittedIsolation,
		store:            map[string][]Value{},
		// "0" transaction id will be used to mean that the id is not set.
		// So all valid transaction ids must start at 1
		nextTransactionId: 1,
	}
}

func (d *Database) inprogress() btree.Set[uint64] {
	var ids btree.Set[uint64]
	iter := d.transactions.Iter()
	for ok := iter.First(); ok; ok = iter.Next() {
		if iter.Value().state == InProgressTransaction {
			ids.Insert(iter.Key())
		}
	}
	return ids
}

func (d *Database) newTransaction() *Transaction {
	t := Transaction{}
	t.isolation = d.defaultIsolation
	t.state = InProgressTransaction

	// Assign and increment txn id
	t.id = d.nextTransactionId
	d.nextTransactionId++

	// Store all inprogress transaction ids
	t.inprogress = d.inprogress()

	// Add this transaction to history
	d.transactions.Set(t.id, t)

	debug("starting transaction", t.id)

	return &t
}

func (d *Database) completeTransaction(t *Transaction, state TransactionSate) error {
	debug("completing transaction", t.id)

	// Update txns
	t.state = state
	d.transactions.Set(t.id, *t)

	return nil
}

func (d *Database) transactionState(txnId uint64) Transaction {
	t, ok := d.transactions.Get(txnId)
	assert(ok, "valid transaction")
	return t
}

func (d *Database) assertValidTransaction(t *Transaction) {
	assert(t.id > 0, "valid id")
	assert(d.transactionState(t.id).state == InProgressTransaction, "in progress")
}
