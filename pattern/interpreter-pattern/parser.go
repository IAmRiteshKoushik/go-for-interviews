package main

type Condition struct {
	Symbol    string  `json:"symbol,omitempty"`
	Indicator string  `json:"indicator,omitempty"`
	Operator  string  `json:"operator,omitempty"`
	Value     float64 `json:"value,omitempty"`
}

func (c *Condition) Interpret(context map[string]*MarketData) bool {
	data, ok := context[c.Symbol]
	if !ok {
		return false
	}

	indicator := data.GetIndicator(c.Indicator)
	switch c.Operator {
	case ">":
		return indicator > c.Value
	case "<":
		return indicator < c.Value
	case "==":
		return indicator == c.Value
	}
	return false
}
