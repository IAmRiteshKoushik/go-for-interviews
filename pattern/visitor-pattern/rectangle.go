package main

type Rectangle struct {
	l, b int
}

func (s *Rectangle) accept(v Visitor) {
	v.visitForRectangle(s)
}

func (s *Rectangle) getType() string {
	return "Rectangle"
}
