package main

import "fmt"

type MarketData struct {
	Symbol string  `json:"symbol"`
	Open   float64 `json:"open"`
	High   float64 `json:"high"`
	Close  float64 `json:"close"`
	Low    float64 `json:"low"`
	Volume float64 `json:"volume"`
}

func (x *MarketData) GetIndicator(indicator string) float64 {
	switch indicator {
	case "open":
		return x.Open
	case "high":
		return x.High
	case "close", "price":
		return x.Close
	case "low":
		return x.Low
	case "volume":
		return x.Volume
	}
	return 0
}

func main() {
	// Mock DATA
	data1 := &MarketData{
		Symbol: "AAPL",
		Open:   100.0,
		High:   110.0,
		Close:  105.0,
		Low:    95.0,
		Volume: 1000000.0,
	}
	data2 := &MarketData{
		Symbol: "GOOG",
		Open:   1000.0,
		High:   1100.0,
		Close:  1050.0,
		Low:    950.0,
		Volume: 10000000.0,
	}
	context := map[string]*MarketData{
		"AAPL": data1,
		"GOOG": data2,
	}

	strategy := NewStrategy([]Condition{
		{
			Symbol:    "AAPL",
			Indicator: "close",
			Operator:  ">",
			Value:     100.0,
		},
		{
			Symbol:    "GOOG",
			Indicator: "volume",
			Operator:  ">",
			Value:     5000000.0,
		},
	}, []Condition{})

	entryExpression := strategy.GetEntryExpression()
	if entryExpression.Interpret(context) {
		fmt.Println("enter trade")
	} else {
		fmt.Println("Do not enter trade")
	}

	exitExpression := strategy.GetExitExpression()
	if exitExpression.Interpret(context) {
		fmt.Println("Exit trade")
	} else {
		fmt.Println("Do not exit trade")
	}
}
