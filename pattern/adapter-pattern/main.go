package main

import "fmt"

func main() {
	client := &Client{}

	mac := &Mac{}
	windows := &Windows{}
	windowsAdapter := &WindowsAdapter{
		windowMachine: windows,
	}

	client.InsertLightningConnectorIntoComputer(mac)
	fmt.Println("")
	client.InsertLightningConnectorIntoComputer(windowsAdapter)
}
