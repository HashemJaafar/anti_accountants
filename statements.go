package main

import (
	"math"
	"strings"
	"time"
)

func FStatement(journal []SJournal1, allEndDates []time.Time, periodInDaysBeforeEndDate uint, namesYouWant []TPersonName, inNames bool) ([]TStatement3, error) {
	var statements1 []TStatement2
	FSortTime(allEndDates, true)
	for _, v1 := range allEndDates {
		trailingBalanceSheet := FStatementStep1(journal, v1.AddDate(0, 0, -int(periodInDaysBeforeEndDate)), v1)
		trailingBalanceSheet = FStatementStep2(trailingBalanceSheet)
		trailingBalanceSheet = FStatementStep3(trailingBalanceSheet)
		trailingBalanceSheet = FStatementStep4(trailingBalanceSheet, inNames, namesYouWant)
		statement := FStatementStep5(trailingBalanceSheet)
		FStatementStep6(periodInDaysBeforeEndDate, statement)
		FStatementStep7(statement)
		statements1 = append(statements1, statement)
	}

	var statements2 []TStatement3
	for _, v1 := range statements1 {
		statement := FStatementStep8(v1, statements1[0])
		FStatementStep9(statement)
		statements2 = append(statements2, statement)
	}

	return statements2, nil
}

func FStatementStep1(journal []SJournal1, dateStart, dateEnd time.Time) TStatement1 {
	newStatement := TStatement1{}

	addToNewStatement := func(v1 SJournal1, isBeforeDateStart bool) {
		m := FInitializeMap6(newStatement, v1.CreditAccountName, v1.DebitAccountName, v1.Name, CValue, isBeforeDateStart)
		m[true] += v1.Value
		m = FInitializeMap6(newStatement, v1.CreditAccountName, v1.DebitAccountName, v1.Name, CQuantity, isBeforeDateStart)
		m[true] += math.Abs(v1.CreditQuantity)
		m = FInitializeMap6(newStatement, v1.DebitAccountName, v1.CreditAccountName, v1.Name, CValue, isBeforeDateStart)
		m[false] += v1.Value
		m = FInitializeMap6(newStatement, v1.DebitAccountName, v1.CreditAccountName, v1.Name, CQuantity, isBeforeDateStart)
		m[false] += math.Abs(v1.DebitQuantity)
	}

	for _, v1 := range journal {
		switch {
		case v1.Date.Before(dateStart):
			addToNewStatement(v1, true)
		case v1.Date.Before(dateEnd):
			addToNewStatement(v1, false)
		default:
			break
		}
	}

	return newStatement
}

func FStatementStep2(oldStatement TStatement1) TStatement1 {
	// in this function i insert the father accounts in column account1
	// i sum the credit to credit and debit to debit .Like:
	// if there is three accounts like this:
	// assets 			debit
	// 	|_equipment 	debit
	// 	|_depreciation 	credit
	// i will sum the debit side of the equipment and depreciation to debit side of assets
	// and the credit side of the equipment and depreciation to credit side of assets

	newStatement := TStatement1{}

	for k1, v1 := range oldStatement { //account1
		account, _, err := FFindAccountFromName(k1)
		for k2, v2 := range v1 { //account2
			for k3, v3 := range v2 { //name
				for k4, v4 := range v3 { //vpq
					for k5, v5 := range v4 { //isBeforeDateStart
						for k6, v6 := range v5 { //is_credit
							// here i copy the map
							m := FInitializeMap6(newStatement, k1, k2, k3, k4, k5)
							m[k6] += v6

							if err == nil {
								for _, v7 := range account.FathersName[VIndexOfAccountNumber] {
									m = FInitializeMap6(newStatement, v7, k2, k3, k4, k5)
									m[k6] += v6
								}
							}
						}
					}
				}
			}
		}
	}

	return newStatement
}

func FStatementStep3(oldStatement TStatement1) TStatement1 {
	// in this function i insert the father accounts and 'AllAccounts' key word in column account2
	// i sum the credit to credit and debit to debit .Like:
	// if there is three accounts like this:
	// assets debit,equipment debit,depreciation credit
	// i will sum the debit side of the equipment and depreciation to debit side of assets
	// and the credit side of the equipment and depreciation to credit side of assets

	newStatement := TStatement1{}

	for k1, v1 := range oldStatement { //account1
		for k2, v2 := range v1 { //account2
			account, _, _ := FFindAccountFromName(k2)
			for k3, v3 := range v2 { //name
				for k4, v4 := range v3 { //vpq
					for k5, v5 := range v4 { //isBeforeDateStart
						for k6, v6 := range v5 { //is_credit
							// here i copy the map
							m := FInitializeMap6(newStatement, k1, k2, k3, k4, k5)
							m[k6] += v6

							// here i insert the key word 'AllAccounts' in column account2
							m = FInitializeMap6(newStatement, k1, CAllAccounts, k3, k4, k5)
							m[k6] += v6

							for _, v7 := range account.FathersName[VIndexOfAccountNumber] {
								m = FInitializeMap6(newStatement, k1, v7, k3, k4, k5)
								m[k6] += v6
							}
						}
					}
				}
			}
		}
	}

	return newStatement
}

func FStatementStep4(oldStatement TStatement1, inNames bool, namesYouWant []string) TStatement1 {
	// in this function i insert the key word 'AllNames' and 'Names' in column name

	newStatement := TStatement1{}

	for k1, v1 := range oldStatement { //account1
		for k2, v2 := range v1 { //account2
			for k3, v3 := range v2 { //name
				for k4, v4 := range v3 { //vpq
					for k5, v5 := range v4 { //isBeforeDateStart
						for k6, v6 := range v5 { //is_credit
							// here i copy the map
							m := FInitializeMap6(newStatement, k1, k2, k3, k4, k5)
							m[k6] += v6

							// here i insert the key word 'AllNames' in column name
							m = FInitializeMap6(newStatement, k1, k2, CAllNames, k4, k5)
							m[k6] += v6

							_, isIn := FFind(k2, namesYouWant)
							if isIn == inNames {
								// here i insert the key word 'Names' in column name
								m = FInitializeMap6(newStatement, k1, k2, CNames, k4, k5)
								m[k6] += v6
							}
						}
					}
				}
			}
		}
	}

	return newStatement
}

func FStatementStep5(oldStatement TStatement1) TStatement2 {
	// in this function i insert the type_of_vpq and remove column isBeforeDateStart and is_credit

	newStatement := TStatement2{}
	for k1, v1 := range oldStatement { //account1
		accountStruct1, _, _ := FFindAccountFromName(k1)
		for k2, v2 := range v1 { //account2
			for k3, v3 := range v2 { //name
				for k4, v4 := range v3 { //vpq
					for k5, v5 := range v4 { //isBeforeDateStart
						for k6, v6 := range v5 { //is_credit
							if k5 {
								// here i insert FlowInBeginning and FlowOutBeginning
								if accountStruct1.IsCredit == k6 {
									m := FInitializeMap5(newStatement, k1, k2, k3, k4)
									m[CFlowInBeginning] += v6
								} else {
									m := FInitializeMap5(newStatement, k1, k2, k3, k4)
									m[CFlowOutBeginning] += v6
								}
							} else {
								// here i insert FlowInPeriod and FlowOutPeriod
								if accountStruct1.IsCredit == k6 {
									m := FInitializeMap5(newStatement, k1, k2, k3, k4)
									m[CFlowInPeriod] += v6
								} else {
									m := FInitializeMap5(newStatement, k1, k2, k3, k4)
									m[CFlowOutPeriod] += v6
								}
							}
						}
					}
				}
			}
		}
	}

	return newStatement
}

func FStatementStep6(days uint, oldStatement TStatement2) {
	// in this function we make vertical analysis of the statement

	for _, v1 := range oldStatement { //account1
		for _, v2 := range v1 { //account2
			for _, v3 := range v2 { //name
				for _, v4 := range v3 { //vpq
					v4[CFlowInEnding] = v4[CFlowInPeriod] + v4[CFlowInBeginning]
					v4[CFlowOutEnding] = v4[CFlowOutPeriod] + v4[CFlowOutBeginning]

					v4[CFlowBeginning] = v4[CFlowInBeginning] - v4[CFlowOutBeginning]
					v4[CFlowPeriod] = v4[CFlowInPeriod] - v4[CFlowOutPeriod]
					v4[CFlowEnding] = v4[CFlowInEnding] - v4[CFlowOutEnding]

					v4[CAverage] = (v4[CFlowEnding] + v4[CFlowBeginning]) / 2
					v4[CTurnover] = v4[CFlowOutPeriod] / v4[CAverage]
					v4[CTurnoverDays] = float64(days) / v4[CTurnover]
					v4[CGrowthRatio] = v4[CFlowEnding] / v4[CFlowBeginning]
				}
			}
		}
	}
}

func FStatementStep7(oldStatement TStatement2) {
	// in this function i complete vertical analysis of the statement
	// but here i calculate the percentage of the name from the account

	for _, v1 := range oldStatement { //account1
		for _, v2 := range v1 { //account2
			for _, v3 := range v2 { //name
				for k4, v4 := range v3 { //vpq
					v4[CNamePercent] = v4[CFlowEnding] / v2[CAllNames][k4][CFlowEnding]
				}
			}
		}
	}
}

func FStatementStep8(oldStatement, baseStatement TStatement2) TStatement3 {
	newStatement := TStatement3{}

	for k1, v1 := range oldStatement { //account1
		for k2, v2 := range v1 { //account2
			for k3, v3 := range v2 { //name
				for k4, v4 := range v3 { //vpq
					for k5, v5 := range v4 { //type_of_vpq
						m := FInitializeMap6(newStatement, k1, k2, k3, k4, k5)
						m[CBalance] = v5
						m[CChangeSinceBasePeriod] = v5 - baseStatement[k1][k2][k3][k4][k5]
						m[CGrowthRatioToBasePeriod] = v5 / baseStatement[k1][k2][k3][k4][k5]
					}
				}
			}
		}
	}
	return newStatement
}

func FStatementStep9(oldStatement TStatement3) {
	// in this function we calculate the Price

	for _, v1 := range oldStatement { //account1
		for _, v2 := range v1 { //account2
			for _, v3 := range v2 { //name
				if v3[CPrice] == nil {
					v3[CPrice] = map[string]map[string]float64{}
				}
				for _, v4 := range v3 { //vpq
					for k5, v5 := range v4 { //type_of_vpq
						if v3[CPrice][k5] == nil {
							v3[CPrice][k5] = map[string]float64{}
						}
						for k6 := range v5 { //Change_or_Ratio_or_balance
							v3[CPrice][k5][k6] = v3[CValue][k5][k6] / v3[CQuantity][k5][k6]
						}
					}
				}
			}
		}
	}
}

func FStatementFilter(oldStatement TStatement3, f SStatement2) []SStatmentWithAccount {
	var newStatement []SStatmentWithAccount

	for k1, v1 := range oldStatement { //account1
		if FFilterString(k1, f.Account1Name.Name) {
			account1, _, err := FFindAccountFromName(k1)
			if err != nil || FFilterAccount(account1, f.Account1Name) {
				for k2, v2 := range v1 { //account2
					if FFilterString(k2, f.Account2Name.Name) {
						account2, _, err := FFindAccountFromName(k2)
						if err != nil || FFilterAccount(account1, f.Account2Name) {
							for k3, v3 := range v2 { //name
								if FFilterString(k3, f.PersonName) {
									for k4, v4 := range v3 { //vpq
										if FFilterString(k4, f.Vpq) {
											for k5, v5 := range v4 { //type_of_vpq
												if FFilterString(k5, f.TypeOfVpq) {
													for k6, v6 := range v5 { //Change_or_Ratio_or_balance
														if FFilterString(k6, f.ChangeOrRatioOrBalance) {
															if FFilterNumber(v6, f.Number) {

																// here i prefer to show the account struct in the statment to use it later in sorting the account
																newStatement = append(newStatement, SStatmentWithAccount{
																	account1,
																	account2,
																	SStatement1{k1, k2, k3, k4, k5, k6, v6},
																})
															}
														}
													}
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return newStatement
}

func FStatementFilterByGreedyAlgorithm(oldStatement []SStatmentWithAccount, isPercent bool, numberTarget float64) []SStatmentWithAccount {
	if numberTarget == 0 {
		return oldStatement
	}

	var numberTotal float64
	if isPercent {
		// here i sum the numbers to find the total amount to calculate the percentage in units
		for _, v1 := range oldStatement {
			numberTotal += v1.SStatement1.Number
		}
		// here i convert the percent to a units
		numberTarget = numberTotal * numberTarget
	}

	FSortStatementNumber(oldStatement, true)
	var newStatement []SStatmentWithAccount
	var numberCurrent float64
	for _, v1 := range oldStatement {
		if numberCurrent < numberTarget {
			newStatement = append(newStatement, v1)
			numberCurrent += v1.SStatement1.Number
		}
	}
	return newStatement
}

func FStatementSort(statement []SStatmentWithAccount, way string) []SStatmentWithAccount {
	switch way {
	case CAscending:
		FSortStatementNumber(statement, true)
		return statement
	case CDescending:
		FSortStatementNumber(statement, false)
		return statement
	default:
		FSortByLevel(statement)
		return statement
	}
}

func FSortByLevel(s []SStatmentWithAccount) []SStatmentWithAccount {
	for k1 := range s {
		for k2 := range s {
			if k1 < k2 &&
				!FIsItHighThanByOrder(s[k1].Account1.Number[VIndexOfAccountNumber], s[k2].Account1.Number[VIndexOfAccountNumber]) {
				FSwap(s, k1, k2) // account1

				if s[k1].SStatement1.Account1Name == s[k2].SStatement1.Account1Name &&
					(s[k1].SStatement1.Account2Name == CAllAccounts || s[k2].SStatement1.Account2Name == CAllAccounts ||
						!FIsItHighThanByOrder(s[k1].Account2.Number[VIndexOfAccountNumber], s[k2].Account2.Number[VIndexOfAccountNumber])) {
					FSwap(s, k1, k2) // account2

					if s[k1].SStatement1.Account2Name == s[k2].SStatement1.Account2Name &&
						s[k1].SStatement1.PersonName > s[k2].SStatement1.PersonName {
						FSwap(s, k1, k2) // name

						if s[k1].SStatement1.PersonName == s[k2].SStatement1.PersonName &&
							s[k1].SStatement1.Vpq > s[k2].SStatement1.Vpq {
							FSwap(s, k1, k2) // vpq

							if s[k1].SStatement1.Vpq == s[k2].SStatement1.Vpq &&
								s[k1].SStatement1.TypeOfVpq > s[k2].SStatement1.TypeOfVpq {
								FSwap(s, k1, k2) // typeOfVpq

								if s[k1].SStatement1.TypeOfVpq == s[k2].SStatement1.TypeOfVpq &&
									s[k1].SStatement1.ChangeOrRatioOrBalance > s[k2].SStatement1.ChangeOrRatioOrBalance {
									FSwap(s, k1, k2) // ChangeOrRatioOrBalance
								}
							}
						}
					}
				}
			}
		}
	}
	return s
}

func FMakeSpaceBeforeAccountInStatementStruct(oldStatement []SStatmentWithAccount) {
	for k1, v1 := range oldStatement {
		oldStatement[k1].SStatement1.Account1Name = strings.Repeat("  ", int(v1.Account1.Levels[VIndexOfAccountNumber])) + v1.SStatement1.Account1Name
		if v1.SStatement1.Account2Name != CAllAccounts {
			oldStatement[k1].SStatement1.Account2Name = strings.Repeat("  ", int(v1.Account2.Levels[VIndexOfAccountNumber])) + v1.SStatement1.Account2Name
		}
	}
}

func FConvertStatmentWithAccountToFilteredStatement(oldStatement []SStatmentWithAccount) []SStatement1 {
	var newStatement []SStatement1
	for _, v1 := range oldStatement {
		newStatement = append(newStatement, v1.SStatement1)
	}
	return newStatement
}

func FStatementAnalysis(i SFinancialAnalysis) SFinancialAnalysisStatement {
	currentRatio := i.CurrentAssets / i.CurrentLiabilities
	acidTest := (i.Cash + i.ShortTermInvestments + i.NetReceivables) / i.CurrentLiabilities
	receivablesTurnover := i.NetCreditSales / i.AverageNetReceivables
	inventoryTurnover := i.CostOfGoodsSold / i.AverageInventory
	profitMargin := i.NetIncome / i.NetSales
	assetTurnover := i.NetSales / i.AverageAssets
	returnOnAssets := i.NetIncome / i.AverageAssets
	returnOnEquity := i.NetIncome / i.AverageEquity
	payoutRatio := i.CashDividends / i.NetIncome
	debtToTotalAssetsRatio := i.TotalDebt / i.TotalAssets
	timesInterestEarned := i.Ebitda / i.InterestExpense
	returnOnCommonStockholdersEquity := (i.NetIncome - i.PreferredDividends) / i.AverageCommonStockholdersEquity
	earningsPerShare := (i.NetIncome - i.PreferredDividends) / i.WeightedAverageCommonSharesOutstanding
	priceEarningsRatio := i.MarketPricePerSharesOutstanding / earningsPerShare
	return SFinancialAnalysisStatement{
		CurrentRatio:                     currentRatio,
		AcidTest:                         acidTest,
		ReceivablesTurnover:              receivablesTurnover,
		InventoryTurnover:                inventoryTurnover,
		AssetTurnover:                    assetTurnover,
		ProfitMargin:                     profitMargin,
		ReturnOnAssets:                   returnOnAssets,
		ReturnOnEquity:                   returnOnEquity,
		PayoutRatio:                      payoutRatio,
		DebtToTotalAssetsRatio:           debtToTotalAssetsRatio,
		TimesInterestEarned:              timesInterestEarned,
		ReturnOnCommonStockholdersEquity: returnOnCommonStockholdersEquity,
		EarningsPerShare:                 earningsPerShare,
		PriceEarningsRatio:               priceEarningsRatio}
}
