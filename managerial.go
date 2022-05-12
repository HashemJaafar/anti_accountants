package main

import (
	"log"
	"math"
)

func CalculateCvpMap(cvp map[string]map[string]float64, print, checkIfKeysInTheEquations bool) {
	for _, v1 := range cvp {
		CostVolumeProfit(print, checkIfKeysInTheEquations, v1)
		_, okVariableCostPerUnits := v1[VariableCostPerUnits]
		if !okVariableCostPerUnits {
			v1[VariableCostPerUnits] = 0
			CostVolumeProfit(print, false, v1)
		}
		_, okFixedCost := v1[FixedCost]
		if !okFixedCost {
			v1[FixedCost] = 0
			CostVolumeProfit(print, false, v1)
		}
		_, okSalesPerUnits := v1[SalesPerUnits]
		if !okSalesPerUnits {
			v1[SalesPerUnits] = 0
			CostVolumeProfit(print, false, v1)
		}
		_, okUnits := v1[Units]
		if !okUnits {
			v1[Units] = 0
			CostVolumeProfit(print, false, v1)
		}
	}
}

func CostVolumeProfit(print, checkIfKeysInTheEquations bool, m map[string]float64) {
	equations := [][]string{
		{VariableCost, VariableCostPerUnits, "*", "units"},
		{FixedCost, FixedCostPerUnits, "*", Units},
		{MixedCost, MixedCostPerUnits, "*", Units},
		{Sales, SalesPerUnits, "*", Units},
		{Profit, ProfitPerUnits, "*", Units},
		{ContributionMargin, ContributionMarginPerUnits, "*", Units},
		{MixedCost, FixedCost, "+", VariableCost},
		{Sales, Profit, "+", MixedCost},
		{ContributionMargin, Sales, "-", VariableCost},
		{BreakEvenInSales, BreakEvenInUnits, "*", SalesPerUnits},
		{BreakEvenInUnits, ContributionMarginPerUnits, "/", FixedCost},
		{ContributionMarginPerUnits, ContributionMarginRatio, "*", SalesPerUnits},
		{ContributionMargin, DegreeOfOperatingLeverage, "*", Profit},
		{UnitsGap, Units, "-", ActualUnits},
	}
	EquationsSolver(print, checkIfKeysInTheEquations, m, equations)
}

func MixCostVolumeProfit(print, checkIfKeysInTheEquations bool, m map[string]map[string]float64) {
	var units, sales, variableCost, fixedCost float64
	for _, v1 := range m {
		units += v1[Units]
		sales += v1[Sales]
		variableCost += v1[VariableCost]
		fixedCost += v1[FixedCost]
	}
	m[Total] = map[string]float64{"units": units, "sales": sales, "variableCost": variableCost, "fixedCost": fixedCost}
	CostVolumeProfit(print, false, m[Total])
}

func CostVolumeProfitSlice(cvp map[string]map[string]float64, distributionSteps []OneStepDistribution, print, simple bool) {
	CalculateCvpMap(cvp, print, true)
	for _, v1 := range distributionSteps {
		var totalMixedCost, totalPortionsTo, totalColumnToDistribute float64
		if simple {
			totalMixedCost = v1.Amount
		} else {
			totalMixedCost = TotalMixedCostInComplicatedAndMultiLevelStep(cvp, v1, totalMixedCost)
		}
		for k2, v2 := range v1.To {
			totalPortionsTo += v2
			totalColumnToDistribute += cvp[k2][v1.DistributionMethod]
		}
		for k2, v2 := range v1.To {
			var totalOverheadCostToSum float64
			switch v1.DistributionMethod {
			case UnitsGap:
				totalOverheadCostToSum = cvp[k2][UnitsGap] * cvp[k2][VariableCostPerUnits]
				cvp[k2][Units] -= cvp[k2][UnitsGap]
				cvp[k2][UnitsGap] = 0
			case "1":
				totalOverheadCostToSum = totalMixedCost
			case Portions:
				totalOverheadCostToSum = v2 / totalPortionsTo * totalMixedCost
			case Units:
				totalOverheadCostToSum = cvp[k2][Units] / totalColumnToDistribute * totalMixedCost
			case VariableCost:
				totalOverheadCostToSum = cvp[k2][VariableCost] / totalColumnToDistribute * totalMixedCost
			case FixedCost:
				totalOverheadCostToSum = cvp[k2][FixedCost] / totalColumnToDistribute * totalMixedCost
			case MixedCost:
				totalOverheadCostToSum = cvp[k2][MixedCost] / totalColumnToDistribute * totalMixedCost
			case Sales:
				totalOverheadCostToSum = cvp[k2][Sales] / totalColumnToDistribute * totalMixedCost
			case Profit:
				totalOverheadCostToSum = cvp[k2][Profit] / totalColumnToDistribute * totalMixedCost
			case ContributionMargin:
				totalOverheadCostToSum = cvp[k2][ContributionMargin] / totalColumnToDistribute * totalMixedCost
			case PercentFromVariableCost:
				totalOverheadCostToSum = cvp[k2][VariableCost] * v2
			case PercentFromFixedCost:
				totalOverheadCostToSum = cvp[k2][FixedCost] * v2
			case PercentFromMixedCost:
				totalOverheadCostToSum = cvp[k2][MixedCost] * v2
			case PercentFromSales:
				totalOverheadCostToSum = cvp[k2][Sales] * v2
			case PercentFromProfit:
				totalOverheadCostToSum = cvp[k2][Profit] * v2
			case PercentFromContributionMargin:
				totalOverheadCostToSum = cvp[k2][ContributionMargin] * v2
			default:
				totalOverheadCostToSum = float64(len(v1.To)) * totalMixedCost
			}
			switch v1.SalesOrVariableOrFixed {
			case VariableCost:
				cvp[k2][VariableCostPerUnits] = ((cvp[k2][VariableCostPerUnits] * cvp[k2][Units]) + totalOverheadCostToSum) / cvp[k2][Units]
			case FixedCost:
				cvp[k2][FixedCost] += totalOverheadCostToSum
			default:
				cvp[k2][SalesPerUnits] = ((cvp[k2][SalesPerUnits] * cvp[k2][Units]) - totalOverheadCostToSum) / cvp[k2][Units]
			}
			for k3, v3 := range cvp {
				cvp[k3] = map[string]float64{"units_gap": v3[UnitsGap], "units": v3[Units],
					"sales_per_units": v3[SalesPerUnits], VariableCostPerUnits: v3[VariableCostPerUnits], "fixedCost": v3[FixedCost]}
			}
			CalculateCvpMap(cvp, print, false)
		}
	}
	TotalCostVolumeProfit(cvp, print)
}

func LaborCost(print, checkIfKeysInTheEquations bool, m map[string]float64) {
	equations := [][]string{
		{"overtime_wage_rate", "bonus_percentage", "*", "hourly_wage_rate"},

		{"overtime_wage", "overtime_hours", "*", "overtime_wage_rate"},
		{"work_time_wage", "work_time_hours", "*", "hourly_wage_rate"},
		{"holiday_wage", "holiday_hours", "*", "hourly_wage_rate"},
		{"vacations_wage", "vacations_hours", "*", "hourly_wage_rate"},
		{"normal_lost_time_wage", "normal_lost_time_hours", "*", "hourly_wage_rate"},
		{"abnormal_lost_time_wage", "abnormal_lost_time_hours", "*", "hourly_wage_rate"},
	}
	EquationsSolver(print, checkIfKeysInTheEquations, m, equations)
}

func ProcessCosting(print, checkIfKeysInTheEquations bool, m map[string]float64) {
	equations := [][]string{
		{"increase_or_decrease", "increase", "-", "decrease"},
		{"increase_or_decrease", "ending_balance", "-", "beginning_balance"},
		{"cost_of_goods_sold", "decrease", "-", "decreases_in_account_caused_by_not_sell"},
		{"equivalent_units", "number_of_partially_completed_units", "*", "percentage_completion"},
		{"equivalent_units_of_production_weighted_average_method", "units_transferred_to_the_next_department_or_to_finished_goods", "+", "equivalent_units_in_ending_work_in_process_inventory"},
		{"cost_of_ending_work_in_process_inventory", "cost_of_beginning_work_in_process_inventory", "+", "cost_added_during_the_period"},
		{"cost_per_equivalent_unit_weighted_average_method", "cost_of_ending_work_in_process_inventory", "/", "equivalent_units_of_production_weighted_average_method"},
		{"equivalent_units_of_production_Fifo_method", "equivalent_units_of_production_weighted_average_method", "-", "equivalent_units_in_beginning_work_in_process_inventory"},
		{"percentage_completion_minus_one", "1", "-", "percentage_completion"},
		{"equivalent_units_to_complete_beginning_work_in_process_inventory", "equivalent_units_in_beginning_work_in_process_inventory", "*", "percentage_completion_minus_one"},
		{"cost_per_equivalent_unit_Fifo_method", "cost_added_during_the_period", "/", "equivalent_units_of_production_Fifo_method"},
	}
	EquationsSolver(print, checkIfKeysInTheEquations, m, equations)
}

func TotalCostVolumeProfit(cvp map[string]map[string]float64, print bool) {
	var units, sales, variableCost, fixedCost float64
	for k1, v1 := range cvp {
		if k1 != "total" {
			units += v1[Units]
			sales += v1[Sales]
			variableCost += v1[VariableCost]
			fixedCost += v1[FixedCost]
		}
	}
	cvp[Total] = map[string]float64{Units: units, Sales: sales, VariableCost: variableCost, FixedCost: fixedCost}
	CostVolumeProfit(print, false, cvp[Total])
}

func TotalMixedCostInComplicatedAndMultiLevelStep(cvp map[string]map[string]float64, step OneStepDistribution, totalMixedCost float64) float64 {
	for k1, v1 := range step.From {
		if cvp[k1][Units] < v1 {
			log.Panic(v1, " for ", k1, " in ", step.From, " is more than ", cvp[k1][Units])
		}
		totalMixedCost += v1 * cvp[k1][MixedCostPerUnits]
		cvp[k1][FixedCost] -= (cvp[k1][FixedCost] / cvp[k1][Units]) * v1
		cvp[k1][Units] -= v1
		if cvp[k1][Units] == 0 {
			cvp[k1][VariableCostPerUnits] = 0
		}
	}
	return totalMixedCost
}

func FCostVolumeProfit(units, salesPerUnit float64, variableCosts, fixedCosts []APQ) (Cvp, Cvp, []AVQ, []AVQ) {
	var a []AVQ
	var variableCost float64
	variableCost, a = newFunction(variableCosts, units, variableCost, a)
	var b []AVQ
	var fixedCost float64
	fixedCost, b = newFunction(fixedCosts, units, fixedCost, b)

	sales := salesPerUnit * units
	mixedCost := fixedCost + variableCost
	profit := sales - mixedCost
	contributionMargin := sales - variableCost

	o1 := Cvp{
		VariableCost:       variableCost,
		FixedCost:          fixedCost,
		MixedCost:          mixedCost,
		Sales:              sales,
		Profit:             profit,
		ContributionMargin: contributionMargin,
	}

	perUint := Cvp{
		VariableCost:       variableCost / units,
		FixedCost:          fixedCost / units,
		MixedCost:          mixedCost / units,
		Sales:              sales / units,
		Profit:             profit / units,
		ContributionMargin: contributionMargin / units,
	}

	return o1, perUint, a, b
}

func newFunction(variableCosts []APQ, units float64, variableCost float64, a []AVQ) (float64, []AVQ) {
	for _, v1 := range variableCosts {
		if v1.Quantity == 0 {
			continue
		}
		x := math.Ceil(units / v1.Quantity)
		variableCost += x * v1.Price

		a = append(a, AVQ{
			Name:     v1.Name,
			Value:    x * v1.Price,
			Quantity: x,
		})
	}
	return variableCost, a
}
