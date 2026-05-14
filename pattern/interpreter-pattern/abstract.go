package main

type Expression interface {
	Interpret(context map[string]*MarketData) bool
}

type NilExpression struct {
}

func (x *NilExpression) Interpret(context map[string]*MarketData) bool {
	return false
}

type AndExpression struct {
	Left  Expression
	Right Expression
}

func (x *AndExpression) Interpret(context map[string]*MarketData) bool {
	return x.Left.Interpret(context) && x.Right.Interpret(context)
}

type OrExpression struct {
	Left  Expression
	Right Expression
}

func (x *OrExpression) Interpret(context map[string]*MarketData) bool {
	return x.Left.Interpret(context) || x.Right.Interpret(context)
}
