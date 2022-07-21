package anti_accountants

import (
	"testing"
)

func TestMixCostVolumeProfit(t *testing.T) {
	a1 := map[string]map[string]float64{
		"book":  {CUnits: 10, CSales: 500, CFixedCost: 250, CVariableCost: 3},
		"book1": {CUnits: 10, CSales: 500, CFixedCost: 250, CVariableCost: 3},
		"book2": {CUnits: 10, CSales: 500, CFixedCost: 250, CVariableCost: 3},
		"book3": {CUnits: 10, CSales: 500, CFixedCost: 250, CVariableCost: 3},
	}
	FMixCostVolumeProfit(true, a1)
	FPrintMap2(a1)
}

func TestFCostVolumeProfit(t *testing.T) {
	a1, a2, a3, a4 := FCostVolumeProfit2(1200, 12000, []SAPQ1{
		{"a", 250, 1},
		{"b", 500, 1},
		{"c", 300, 1},
		{"d", 8000, 100},
	}, []SAPQ1{
		{"e", 500000, 1000},
		{"f", 700000, 600},
	})

	FPrintCvp(a1)
	FPrintCvp(a2)
	FPrintStructSlice(false, a3)
	FPrintStructSlice(false, a4)
}
