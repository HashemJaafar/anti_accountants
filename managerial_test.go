package main

import "testing"

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
