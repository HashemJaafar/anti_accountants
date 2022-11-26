package anti_accountants

import (
	"log"
	"math"
)

func FCalculateCvpMap(cvp map[string]map[string]float64, print, checkIfKeysInTheEquations bool) {
	for _, v1 := range cvp {
		FCostVolumeProfit1(print, checkIfKeysInTheEquations, v1)

		zeroIfNot := func(key string) {
			_, isExist := v1[key]
			if !isExist {
				v1[key] = 0
				FCostVolumeProfit1(print, false, v1)
			}
		}
		zeroIfNot(CVariableCostPerUnits)
		zeroIfNot(CFixedCost)
		zeroIfNot(CSalesPerUnits)
		zeroIfNot(CUnits)
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
		{CBreakEvenInUnits, CFixedCost, "/", CContributionMarginPerUnits},
		{CContributionMarginPerUnits, CContributionMarginRatio, "*", CSalesPerUnits},
		{CContributionMargin, CDegreeOfOperatingLeverage, "*", CProfit},
		{CUnitsGap, CUnits, "-", CActualUnits},
	}
	FEquationsSolver(print, checkIfKeysInTheEquations, m, equations)
}

func FMixCostVolumeProfit(print bool, m map[string]map[string]float64) {
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

func FCostVolumeProfit2(units, salesPerUnit float64, variableCosts, fixedCosts []SAPQ1) (SCvp, SCvp, []SAVQ, []SAVQ) {
	FSetCostAndAVQ := func(variableCosts []SAPQ1, a []SAVQ) (float64, []SAVQ) {
		var Cost float64
		for _, v1 := range variableCosts {
			if v1.Quantity == 0 {
				continue
			}
			x := math.Ceil(units / v1.Quantity)
			Cost += x * v1.Price

			a = append(a, SAVQ{
				TAccountName: v1.AccountName,
				TValue:       x * v1.Price,
				TQuantity:    x,
			})
		}
		return Cost, a
	}

	var a []SAVQ
	variableCost, a := FSetCostAndAVQ(variableCosts, a)
	var b []SAVQ
	fixedCost, b := FSetCostAndAVQ(fixedCosts, b)

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
