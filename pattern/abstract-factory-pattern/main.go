package main

import (
	"fmt"
)

// Define an abstract factory
type GUIFactory interface {
	CreateButton() Button
	CreateCheckbox() Checkbox
}

// Define interfaces for products
type (
	Button interface {
		Press() string
	}

	Checkbox interface {
		Check() string
	}

	// Implement concrete products for Windows
	WindowsButton   struct{}
	WindowsCheckbox struct{}

	// Implement concrete products for Mac
	MacButton   struct{}
	MacCheckbox struct{}

	// Implement factories for each platform
	WindowsFactory struct{}
	MacFactory     struct{}
)

func (w *WindowsButton) Press() string { return "Windows Button Pressed" }

func (w *WindowsCheckbox) Check() string { return "Windows Checkbox Checked" }

func (m *MacButton) Press() string { return "Mac Button Pressed" }

func (m *MacCheckbox) Check() string { return "Mac Checkbox Checked" }

func (w *WindowsFactory) CreateButton() Button     { return &WindowsButton{} }
func (w *WindowsFactory) CreateCheckbox() Checkbox { return &WindowsCheckbox{} }

func (m *MacFactory) CreateButton() Button     { return &MacButton{} }
func (m *MacFactory) CreateCheckbox() Checkbox { return &MacCheckbox{} }

func main() {
	// Get a Windows factory
	var wf GUIFactory = &WindowsFactory{}
	button := wf.CreateButton()
	checkbox := wf.CreateCheckbox()

	fmt.Println(button.Press())   // Output: Windows Button Pressed
	fmt.Println(checkbox.Check()) // Output: Windows Checkbox Checked

	var mf GUIFactory = &MacFactory{}
	button = mf.CreateButton()
	checkbox = mf.CreateCheckbox()

	fmt.Println(button.Press())   // Output: Mac Button Pressed
	fmt.Println(checkbox.Check()) // Output: Mac Checkbox Checked
}
