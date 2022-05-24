package main

import (
	"log"
	"math"
)

func FCalculateCvpMap(cvp map[string]map[string]float64, print, checkIfKeysInTheEquations bool) {
	for _, v1 := range cvp {
		FCostVolumeProfit1(print, checkIfKeysInTheEquations, v1)
		_, okVariableCostPerUnits := v1[CVariableCostPerUnits]
		if !okVariableCostPerUnits {
			v1[CVariableCostPerUnits] = 0
			FCostVolumeProfit1(print, false, v1)
		}
		_, okFixedCost := v1[CFixedCost]
		if !okFixedCost {
			v1[CFixedCost] = 0
			FCostVolumeProfit1(print, false, v1)
		}
		_, okSalesPerUnits := v1[CSalesPerUnits]
		if !okSalesPerUnits {
			v1[CSalesPerUnits] = 0
			FCostVolumeProfit1(print, false, v1)
		}
		_, okUnits := v1[CUnits]
		if !okUnits {
			v1[CUnits] = 0
			FCostVolumeProfit1(print, false, v1)
		}
	}
}

func FCostVolumeProfit1(print, checkIfKeysInTheEquations bool, m map[string]float64) {
	equations := [][]string{
		{CVariableCost, CVariableCostPerUnits, "*", "units"},
		{CFixedCost, CFixedCostPerUnits, "*", CUnits},
		{CMixedCost, CMixedCostPerUnits, "*", CUnits},
		{CSales, CSalesPerUnits, "*", CUnits},
		{CProfit, CProfitPerUnits, "*", CUnits},
		{CContributionMargin, CContributionMarginPerUnits, "*", CUnits},
		{CMixedCost, CFixedCost, "+", CVariableCost},
		{CSales, CProfit, "+", CMixedCost},
		{CContributionMargin, CSales, "-", CVariableCost},
		{CBreakEvenInSales, CBreakEvenInUnits, "*", CSalesPerUnits},
		{CBreakEvenInUnits, CContributionMarginPerUnits, "/", CFixedCost},
		{CContributionMarginPerUnits, CContributionMarginRatio, "*", CSalesPerUnits},
		{CContributionMargin, CDegreeOfOperatingLeverage, "*", CProfit},
		{CUnitsGap, CUnits, "-", CActualUnits},
	}
	FEquationsSolver(print, checkIfKeysInTheEquations, m, equations)
}

func FMixCostVolumeProfit(print, checkIfKeysInTheEquations bool, m map[string]map[string]float64) {
	var units, sales, variableCost, fixedCost float64
	for _, v1 := range m {
		units += v1[CUnits]
		sales += v1[CSales]
		variableCost += v1[CVariableCost]
		fixedCost += v1[CFixedCost]
	}
	m[CTotal] = map[string]float64{"units": units, "sales": sales, "variableCost": variableCost, "fixedCost": fixedCost}
	FCostVolumeProfit1(print, false, m[CTotal])
}

func FCostVolumeProfitSlice(cvp map[string]map[string]float64, distributionSteps []SOneStepDistribution, print, simple bool) {
	FCalculateCvpMap(cvp, print, true)
	for _, v1 := range distributionSteps {
		var totalMixedCost, totalPortionsTo, totalColumnToDistribute float64
		if simple {
			totalMixedCost = v1.Amount
		} else {
			totalMixedCost = FTotalMixedCostInComplicatedAndMultiLevelStep(cvp, v1, totalMixedCost)
		}
		for k2, v2 := range v1.To {
			totalPortionsTo += v2
			totalColumnToDistribute += cvp[k2][v1.DistributionMethod]
		}
		for k2, v2 := range v1.To {
			var totalOverheadCostToSum float64
			switch v1.DistributionMethod {
			case CUnitsGap:
				totalOverheadCostToSum = cvp[k2][CUnitsGap] * cvp[k2][CVariableCostPerUnits]
				cvp[k2][CUnits] -= cvp[k2][CUnitsGap]
				cvp[k2][CUnitsGap] = 0
			case "1":
				totalOverheadCostToSum = totalMixedCost
			case CPortions:
				totalOverheadCostToSum = v2 / totalPortionsTo * totalMixedCost
			case CUnits:
				totalOverheadCostToSum = cvp[k2][CUnits] / totalColumnToDistribute * totalMixedCost
			case CVariableCost:
				totalOverheadCostToSum = cvp[k2][CVariableCost] / totalColumnToDistribute * totalMixedCost
			case CFixedCost:
				totalOverheadCostToSum = cvp[k2][CFixedCost] / totalColumnToDistribute * totalMixedCost
			case CMixedCost:
				totalOverheadCostToSum = cvp[k2][CMixedCost] / totalColumnToDistribute * totalMixedCost
			case CSales:
				totalOverheadCostToSum = cvp[k2][CSales] / totalColumnToDistribute * totalMixedCost
			case CProfit:
				totalOverheadCostToSum = cvp[k2][CProfit] / totalColumnToDistribute * totalMixedCost
			case CContributionMargin:
				totalOverheadCostToSum = cvp[k2][CContributionMargin] / totalColumnToDistribute * totalMixedCost
			case CPercentFromVariableCost:
				totalOverheadCostToSum = cvp[k2][CVariableCost] * v2
			case CPercentFromFixedCost:
				totalOverheadCostToSum = cvp[k2][CFixedCost] * v2
			case CPercentFromMixedCost:
				totalOverheadCostToSum = cvp[k2][CMixedCost] * v2
			case CPercentFromSales:
				totalOverheadCostToSum = cvp[k2][CSales] * v2
			case CPercentFromProfit:
				totalOverheadCostToSum = cvp[k2][CProfit] * v2
			case CPercentFromContributionMargin:
				totalOverheadCostToSum = cvp[k2][CContributionMargin] * v2
			default:
				totalOverheadCostToSum = float64(len(v1.To)) * totalMixedCost
			}
			switch v1.SalesOrVariableOrFixed {
			case CVariableCost:
				cvp[k2][CVariableCostPerUnits] = ((cvp[k2][CVariableCostPerUnits] * cvp[k2][CUnits]) + totalOverheadCostToSum) / cvp[k2][CUnits]
			case CFixedCost:
				cvp[k2][CFixedCost] += totalOverheadCostToSum
			default:
				cvp[k2][CSalesPerUnits] = ((cvp[k2][CSalesPerUnits] * cvp[k2][CUnits]) - totalOverheadCostToSum) / cvp[k2][CUnits]
			}
			for k3, v3 := range cvp {
				cvp[k3] = map[string]float64{"units_gap": v3[CUnitsGap], "units": v3[CUnits],
					"sales_per_units": v3[CSalesPerUnits], CVariableCostPerUnits: v3[CVariableCostPerUnits], "fixedCost": v3[CFixedCost]}
			}
			FCalculateCvpMap(cvp, print, false)
		}
	}
	FTotalCostVolumeProfit(cvp, print)
}

func FLaborCost(print, checkIfKeysInTheEquations bool, m map[string]float64) {
	equations := [][]string{
		{"overtime_wage_rate", "bonus_percentage", "*", "hourly_wage_rate"},

		{"overtime_wage", "overtime_hours", "*", "overtime_wage_rate"},
		{"work_time_wage", "work_time_hours", "*", "hourly_wage_rate"},
		{"holiday_wage", "holiday_hours", "*", "hourly_wage_rate"},
		{"vacations_wage", "vacations_hours", "*", "hourly_wage_rate"},
		{"normal_lost_time_wage", "normal_lost_time_hours", "*", "hourly_wage_rate"},
		{"abnormal_lost_time_wage", "abnormal_lost_time_hours", "*", "hourly_wage_rate"},
	}
	FEquationsSolver(print, checkIfKeysInTheEquations, m, equations)
}

func FProcessCosting(print, checkIfKeysInTheEquations bool, m map[string]float64) {
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
	FEquationsSolver(print, checkIfKeysInTheEquations, m, equations)
}

func FTotalCostVolumeProfit(cvp map[string]map[string]float64, print bool) {
	var units, sales, variableCost, fixedCost float64
	for k1, v1 := range cvp {
		if k1 != "total" {
			units += v1[CUnits]
			sales += v1[CSales]
			variableCost += v1[CVariableCost]
			fixedCost += v1[CFixedCost]
		}
	}
	cvp[CTotal] = map[string]float64{CUnits: units, CSales: sales, CVariableCost: variableCost, CFixedCost: fixedCost}
	FCostVolumeProfit1(print, false, cvp[CTotal])
}

func FTotalMixedCostInComplicatedAndMultiLevelStep(cvp map[string]map[string]float64, step SOneStepDistribution, totalMixedCost float64) float64 {
	for k1, v1 := range step.From {
		if cvp[k1][CUnits] < v1 {
			log.Panic(v1, " for ", k1, " in ", step.From, " is more than ", cvp[k1][CUnits])
		}
		totalMixedCost += v1 * cvp[k1][CMixedCostPerUnits]
		cvp[k1][CFixedCost] -= (cvp[k1][CFixedCost] / cvp[k1][CUnits]) * v1
		cvp[k1][CUnits] -= v1
		if cvp[k1][CUnits] == 0 {
			cvp[k1][CVariableCostPerUnits] = 0
		}
	}
	return totalMixedCost
}

func FCostVolumeProfit2(units, salesPerUnit float64, variableCosts, fixedCosts []SAPQ) (SCvp, SCvp, []SAVQ, []SAVQ) {
	var a []SAVQ
	var variableCost float64
	variableCost, a = FSetCostAndAVQ(variableCosts, units, variableCost, a)
	var b []SAVQ
	var fixedCost float64
	fixedCost, b = FSetCostAndAVQ(fixedCosts, units, fixedCost, b)

	sales := salesPerUnit * units
	mixedCost := fixedCost + variableCost
	profit := sales - mixedCost
	contributionMargin := sales - variableCost

	o1 := SCvp{
		VariableCost:       variableCost,
		FixedCost:          fixedCost,
		MixedCost:          mixedCost,
		Sales:              sales,
		Profit:             profit,
		ContributionMargin: contributionMargin,
	}

	perUint := SCvp{
		VariableCost:       variableCost / units,
		FixedCost:          fixedCost / units,
		MixedCost:          mixedCost / units,
		Sales:              sales / units,
		Profit:             profit / units,
		ContributionMargin: contributionMargin / units,
	}

	return o1, perUint, a, b
}

func FSetCostAndAVQ(variableCosts []SAPQ, units float64, variableCost float64, a []SAVQ) (float64, []SAVQ) {
	for _, v1 := range variableCosts {
		if v1.TQuantity == 0 {
			continue
		}
		x := math.Ceil(units / v1.TQuantity)
		variableCost += x * v1.TPrice

		a = append(a, SAVQ{
			TAccountName: v1.TAccountName,
			TValue:       x * v1.TPrice,
			TQuantity:    x,
		})
	}
	return variableCost, a
}
