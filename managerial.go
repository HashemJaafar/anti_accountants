package main

import "log"

func CalculateCvpMap(cvp map[string]map[string]float64, print, checkIfKeysInTheEquations bool) {
	for _, i := range cvp {
		CostVolumeProfit(print, checkIfKeysInTheEquations, i)
		_, okVariableCostPerUnits := i[VariableCostPerUnits]
		if !okVariableCostPerUnits {
			i[VariableCostPerUnits] = 0
			CostVolumeProfit(print, false, i)
		}
		_, okFixedCost := i[FixedCost]
		if !okFixedCost {
			i[FixedCost] = 0
			CostVolumeProfit(print, false, i)
		}
		_, okSalesPerUnits := i[SalesPerUnits]
		if !okSalesPerUnits {
			i[SalesPerUnits] = 0
			CostVolumeProfit(print, false, i)
		}
		_, okUnits := i[Units]
		if !okUnits {
			i[Units] = 0
			CostVolumeProfit(print, false, i)
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
	for _, step := range distributionSteps {
		var totalMixedCost, totalPortionsTo, totalColumnToDistribute float64
		if simple {
			totalMixedCost = step.Amount
		} else {
			totalMixedCost = TotalMixedCostInComplicatedAndMultiLevelStep(cvp, step, totalMixedCost)
		}
		for keyPortionsTo, portionsTo := range step.To {
			totalPortionsTo += portionsTo
			totalColumnToDistribute += cvp[keyPortionsTo][step.DistributionMethod]
		}
		for keyPortionsTo, portionsTo := range step.To {
			var totalOverheadCostToSum float64
			switch step.DistributionMethod {
			case UnitsGap:
				totalOverheadCostToSum = cvp[keyPortionsTo][UnitsGap] * cvp[keyPortionsTo][VariableCostPerUnits]
				cvp[keyPortionsTo][Units] -= cvp[keyPortionsTo][UnitsGap]
				cvp[keyPortionsTo][UnitsGap] = 0
			case "1":
				totalOverheadCostToSum = totalMixedCost
			case Portions:
				totalOverheadCostToSum = portionsTo / totalPortionsTo * totalMixedCost
			case Units:
				totalOverheadCostToSum = cvp[keyPortionsTo][Units] / totalColumnToDistribute * totalMixedCost
			case VariableCost:
				totalOverheadCostToSum = cvp[keyPortionsTo][VariableCost] / totalColumnToDistribute * totalMixedCost
			case FixedCost:
				totalOverheadCostToSum = cvp[keyPortionsTo][FixedCost] / totalColumnToDistribute * totalMixedCost
			case MixedCost:
				totalOverheadCostToSum = cvp[keyPortionsTo][MixedCost] / totalColumnToDistribute * totalMixedCost
			case Sales:
				totalOverheadCostToSum = cvp[keyPortionsTo][Sales] / totalColumnToDistribute * totalMixedCost
			case Profit:
				totalOverheadCostToSum = cvp[keyPortionsTo][Profit] / totalColumnToDistribute * totalMixedCost
			case ContributionMargin:
				totalOverheadCostToSum = cvp[keyPortionsTo][ContributionMargin] / totalColumnToDistribute * totalMixedCost
			case PercentFromVariableCost:
				totalOverheadCostToSum = cvp[keyPortionsTo][VariableCost] * portionsTo
			case PercentFromFixedCost:
				totalOverheadCostToSum = cvp[keyPortionsTo][FixedCost] * portionsTo
			case PercentFromMixedCost:
				totalOverheadCostToSum = cvp[keyPortionsTo][MixedCost] * portionsTo
			case PercentFromSales:
				totalOverheadCostToSum = cvp[keyPortionsTo][Sales] * portionsTo
			case PercentFromProfit:
				totalOverheadCostToSum = cvp[keyPortionsTo][Profit] * portionsTo
			case PercentFromContributionMargin:
				totalOverheadCostToSum = cvp[keyPortionsTo][ContributionMargin] * portionsTo
			default:
				totalOverheadCostToSum = float64(len(step.To)) * totalMixedCost
			}
			switch step.SalesOrVariableOrFixed {
			case VariableCost:
				cvp[keyPortionsTo][VariableCostPerUnits] = ((cvp[keyPortionsTo][VariableCostPerUnits] * cvp[keyPortionsTo][Units]) + totalOverheadCostToSum) / cvp[keyPortionsTo][Units]
			case FixedCost:
				cvp[keyPortionsTo][FixedCost] += totalOverheadCostToSum
			default:
				cvp[keyPortionsTo][SalesPerUnits] = ((cvp[keyPortionsTo][SalesPerUnits] * cvp[keyPortionsTo][Units]) - totalOverheadCostToSum) / cvp[keyPortionsTo][Units]
			}
			for k1, v1 := range cvp {
				cvp[k1] = map[string]float64{"units_gap": v1[UnitsGap], "units": v1[Units],
					"sales_per_units": v1[SalesPerUnits], VariableCostPerUnits: v1[VariableCostPerUnits], "fixedCost": v1[FixedCost]}
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
	for keyPortionsFrom, portions := range step.From {
		if cvp[keyPortionsFrom][Units] < portions {
			log.Panic(portions, " for ", keyPortionsFrom, " in ", step.From, " is more than ", cvp[keyPortionsFrom][Units])
		}
		totalMixedCost += portions * cvp[keyPortionsFrom][MixedCostPerUnits]
		cvp[keyPortionsFrom][FixedCost] -= (cvp[keyPortionsFrom][FixedCost] / cvp[keyPortionsFrom][Units]) * portions
		cvp[keyPortionsFrom][Units] -= portions
		if cvp[keyPortionsFrom][Units] == 0 {
			cvp[keyPortionsFrom][VariableCostPerUnits] = 0
		}
	}
	return totalMixedCost
}
