package main

import (
	"fmt"
	"log"
)

func main() {
	vendingMachine := newVendingMachine(1, 10)

	err := vendingMachine.requestItem()
	if err != nil {
		log.Fatal(err.Error())
	}

	err = vendingMachine.insertMoney(10)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = vendingMachine.dispenseItem()
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println()

	err = vendingMachine.addItem(2)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println()

	err = vendingMachine.requestItem()
	if err != nil {
		log.Fatal(err.Error())
	}

	err = vendingMachine.insertMoney(10)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = vendingMachine.dispenseItem()
	if err != nil {
		log.Fatal(err.Error())
	}

}
