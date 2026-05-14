package main

import (
	"fmt"
	"sync"
)

var once sync.Once
var singleInstance *single

type single struct {
}

func getInstance() *single {
	if singleInstance == nil {
		once.Do(
			func() {
				fmt.Println("Create single instance now.")
				singleInstance = &single{}
			})
	} else {
		fmt.Println("Single instance already created")
	}

	return singleInstance
}
