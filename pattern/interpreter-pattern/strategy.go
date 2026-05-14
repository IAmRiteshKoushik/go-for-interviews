package main

type StrategyParameters struct {
	EntryConditions []Condition `json:"entry_conditions,omitempty"`
	ExitConditions  []Condition `json:"exit_conditions,omitempty"`
}

func NewStrategy(entryConditions []Condition, exitConditions []Condition) *StrategyParameters {
	return &StrategyParameters{
		EntryConditions: entryConditions,
		ExitConditions:  exitConditions,
	}
}

func (x *StrategyParameters) GetEntryExpression() Expression {
	return x.getExpression(x.EntryConditions)
}

func (x *StrategyParameters) GetExitExpression() Expression {
	return x.getExpression(x.ExitConditions)
}

// Helper function
func (x *StrategyParameters) getExpression(conditions []Condition) Expression {
	if len(conditions) == 0 {
		return &NilExpression{}
	}
	if len(conditions) == 1 {
		return &conditions[0]
	}
	return &AndExpression{
		Left:  &conditions[0],
		Right: x.getExpression(conditions[1:]),
	}
}
