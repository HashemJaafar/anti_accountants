package main

import (
	"testing"
)

func TestMixCostVolumeProfit(t *testing.T) {
	a1 := map[string]map[string]float64{
		"book":  {Units: 10, Sales: 500, FixedCost: 250, VariableCost: 3},
		"book1": {Units: 10, Sales: 500, FixedCost: 250, VariableCost: 3},
		"book2": {Units: 10, Sales: 500, FixedCost: 250, VariableCost: 3},
		"book3": {Units: 10, Sales: 500, FixedCost: 250, VariableCost: 3},
	}
	MixCostVolumeProfit(true, true, a1)
	PrintMap2(a1)
}

func TestFCostVolumeProfit(t *testing.T) {
	a1, a2, a3, a4 := FCostVolumeProfit(1200, 12000, []APQ{
		{"a", 250, 1},
		{"b", 500, 1},
		{"c", 300, 1},
		{"d", 8000, 100},
	}, []APQ{
		{"e", 500000, 1000},
		{"f", 700000, 600},
	})

	PrintCvp(a1)
	PrintCvp(a2)
	PrintSlice(a3)
	PrintSlice(a4)
}
