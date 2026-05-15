package main

func main() {
	cashier := &Cashier{}
	medical := &Medical{}
	doctor := &Doctor{}
	reception := &Reception{}
	patient := &Patient{name: "abc"}

	// Chain of Responsibility
	reception.setNext(doctor)
	doctor.setNext(medical)
	medical.setNext(cashier)

	// Execution start
	reception.execute(patient)
}
