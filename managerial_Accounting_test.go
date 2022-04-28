package main

import "testing"

func TestMIX_COST_VOLUME_PROFIT(t *testing.T) {
	a1 := map[string]map[string]float64{
		"book":  {units: 10, sales: 500, fixed_cost: 250, variable_cost: 3},
		"book1": {units: 10, sales: 500, fixed_cost: 250, variable_cost: 3},
		"book2": {units: 10, sales: 500, fixed_cost: 250, variable_cost: 3},
		"book3": {units: 10, sales: 500, fixed_cost: 250, variable_cost: 3},
	}
	MIX_COST_VOLUME_PROFIT(true, true, a1)
	print_map_2(a1)
}
