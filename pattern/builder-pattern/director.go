package main

// The director is not a strict necessity but if you are adding it, then it is
// a good place to put various construction routines so you can re-use them
// across your program

type Director struct {
	builder IBuilder
}

func newDirector(b IBuilder) *Director {
	return &Director{
		builder: b,
	}
}

func (d *Director) setBuilder(b IBuilder) {
	d.builder = b
}

func (d *Director) buildHouse() House {
	d.builder.setWindowType()
	d.builder.setDoorType()
	d.builder.setNumFloor()

	return d.builder.getHouse()
}
