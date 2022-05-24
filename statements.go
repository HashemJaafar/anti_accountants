package main

import (
	"strings"
	"time"
)

func FStatement(allEndDates []time.Time, periodInDaysBeforeEndDate uint, namesYouWant []TPersonName, inNames, withoutReverseEntry bool) ([]TStatement3, error) {

	keys, journal := FDbRead[SJournal](VDbJournal)

	if withoutReverseEntry {
		keys, journal = FFilterJournalFromReverseEntry(keys, journal)
	}

	dates := FConvertByteSliceToTime(keys)

	var statements1 []TStatement2
	FSortTime(allEndDates, true)
	for _, v1 := range allEndDates {
		trailingBalanceSheet := FStatementStep1(dates, journal, v1.AddDate(0, 0, -int(periodInDaysBeforeEndDate)), v1)
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

func FFilterJournalFromReverseEntry(keys [][]byte, journal []SJournal) ([][]byte, []SJournal) {
	var newKeys [][]byte
	var newJournal []SJournal
	for k1, v1 := range journal {
		if !v1.IsReverse && !v1.IsReversed {
			newKeys = append(newKeys, keys[k1])
			newJournal = append(newJournal, v1)
		}
	}
	return newKeys, newJournal
}

func FStatementStep1(dates []time.Time, journal []SJournal, dateStart, dateEnd time.Time) TStatement1 {
	newStatement := TStatement1{}
	for k1, v1 := range journal {
		switch {
		case dates[k1].Before(dateStart):
			FFillNewStatement(newStatement, v1, true)
		case dates[k1].Before(dateEnd):
			FFillNewStatement(newStatement, v1, false)
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
								for _, v7 := range account.TAccountFathersName[VIndexOfAccountNumber] {
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

							for _, v7 := range account.TAccountFathersName[VIndexOfAccountNumber] {
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
								if accountStruct1.TIsCredit == k6 {
									m := FInitializeMap5(newStatement, k1, k2, k3, k4)
									m[CFlowInBeginning] += v6
								} else {
									m := FInitializeMap5(newStatement, k1, k2, k3, k4)
									m[CFlowOutBeginning] += v6
								}
							} else {
								// here i insert FlowInPeriod and FlowOutPeriod
								if accountStruct1.TIsCredit == k6 {
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

func FFillNewStatement(newStatement TStatement1, v1 SJournal, isBeforeDateStart bool) {
	m := FInitializeMap6(newStatement, v1.AccountCredit, v1.AccountDebit, v1.Name, CValue, isBeforeDateStart)
	m[true] += v1.Value
	m = FInitializeMap6(newStatement, v1.AccountCredit, v1.AccountDebit, v1.Name, CQuantity, isBeforeDateStart)
	m[true] += FAbs(v1.QuantityCredit)
	m = FInitializeMap6(newStatement, v1.AccountDebit, v1.AccountCredit, v1.Name, CValue, isBeforeDateStart)
	m[false] += v1.Value
	m = FInitializeMap6(newStatement, v1.AccountDebit, v1.AccountCredit, v1.Name, CQuantity, isBeforeDateStart)
	m[false] += FAbs(v1.QuantityDebit)
}

func FStatementFilter(oldStatement TStatement3, f SFilterStatement) []SStatmentWithAccount {
	var newStatement []SStatmentWithAccount

	for k1, v1 := range oldStatement { //account1
		if f.Account1.Account.FFilter(k1) {
			account1, _, err := FFindAccountFromName(k1)
			if f.Account1.FFilter(account1, err) {
				for k2, v2 := range v1 { //account2
					if f.Account2.Account.FFilter(k2) {
						account2, _, err := FFindAccountFromName(k2)
						if f.Account2.FFilter(account1, err) {
							for k3, v3 := range v2 { //name
								if f.Name.FFilter(k3) {
									for k4, v4 := range v3 { //vpq
										if f.Vpq.FFilter(k4) {
											for k5, v5 := range v4 { //type_of_vpq
												if f.TypeOfVpq.FFilter(k5) {
													for k6, v6 := range v5 { //Change_or_Ratio_or_balance
														if f.ChangeOrRatioOrBalance.FFilter(k6) {
															if f.Number.FFilter(v6) {

																// here i prefer to show the account struct in the statment to use it later in sorting the account
																newStatement = append(newStatement, SStatmentWithAccount{
																	account1,
																	account2,
																	SStatement{k1, k2, k3, k4, k5, k6, v6},
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
			numberTotal += v1.SStatement.TNumber
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
			numberCurrent += v1.SStatement.TNumber
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
				!FIsItHighThanByOrder(s[k1].Account1.TAccountNumber[VIndexOfAccountNumber], s[k2].Account1.TAccountNumber[VIndexOfAccountNumber]) {
				FSwap(s, k1, k2) // account1

				if s[k1].SStatement.TAccount1Name == s[k2].SStatement.TAccount1Name &&
					(s[k1].SStatement.TAccount2Name == CAllAccounts || s[k2].SStatement.TAccount2Name == CAllAccounts ||
						!FIsItHighThanByOrder(s[k1].Account2.TAccountNumber[VIndexOfAccountNumber], s[k2].Account2.TAccountNumber[VIndexOfAccountNumber])) {
					FSwap(s, k1, k2) // account2

					if s[k1].SStatement.TAccount2Name == s[k2].SStatement.TAccount2Name &&
						s[k1].SStatement.TPersonName > s[k2].SStatement.TPersonName {
						FSwap(s, k1, k2) // name

						if s[k1].SStatement.TPersonName == s[k2].SStatement.TPersonName &&
							s[k1].SStatement.TVpq > s[k2].SStatement.TVpq {
							FSwap(s, k1, k2) // vpq

							if s[k1].SStatement.TVpq == s[k2].SStatement.TVpq &&
								s[k1].SStatement.TTypeOfVpq > s[k2].SStatement.TTypeOfVpq {
								FSwap(s, k1, k2) // typeOfVpq

								if s[k1].SStatement.TTypeOfVpq == s[k2].SStatement.TTypeOfVpq &&
									s[k1].SStatement.TChangeOrRatioOrBalance > s[k2].SStatement.TChangeOrRatioOrBalance {
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
		oldStatement[k1].SStatement.TAccount1Name = strings.Repeat("  ", int(v1.Account1.TAccountLevels[VIndexOfAccountNumber])) + v1.SStatement.TAccount1Name
		if v1.SStatement.TAccount2Name != CAllAccounts {
			oldStatement[k1].SStatement.TAccount2Name = strings.Repeat("  ", int(v1.Account2.TAccountLevels[VIndexOfAccountNumber])) + v1.SStatement.TAccount2Name
		}
	}
}

func FConvertStatmentWithAccountToFilteredStatement(oldStatement []SStatmentWithAccount) []SStatement {
	var newStatement []SStatement
	for _, v1 := range oldStatement {
		newStatement = append(newStatement, v1.SStatement)
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
