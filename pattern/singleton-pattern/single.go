package main

import (
	"fmt"
	"sync"
)

var lock = &sync.Mutex{}

var singleInstance *single

type single struct {
}

func getInstance() *single {
	if singleInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		// While u wait for the lock, another thread creates the instance so
		// you double check
		if singleInstance == nil {
			fmt.Println("Creating single instance now.")
			singleInstance = &single{}
		} else {
			fmt.Println("Single instance already created while waiting for lock")
		}

	} else {
		fmt.Println("Single instance already exists")
	}

	return singleInstance
}
