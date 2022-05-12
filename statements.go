package main

import (
	"strings"
	"time"
)

// the columns is:account1,account2,name,vpq,isBeforeDateStart,is_credit,number
type statement1 map[string]map[string]map[string]map[string]map[bool]map[bool]float64

// the columns is:account1,account2,name,vpq,type_of_vpq,number
type statement2 map[string]map[string]map[string]map[string]map[string]float64

// the columns is:account1,account2,name,vpq,type_of_vpq,Change_or_Ratio_or_balance,number
type statement3 map[string]map[string]map[string]map[string]map[string]map[string]float64

func FinancialStatements(allEndDates []time.Time, periodInDaysBeforeEndDate uint, namesYouWant []string, inNames, withoutReverseEntry bool) ([]statement3, error) {

	keys, journal := DbRead[Journal](DbJournal)

	if withoutReverseEntry {
		keys, journal = FilterJournalFromReverseEntry(keys, journal)
	}

	journalTimes := ConvertByteSliceToTime(keys)

	var statements1 []statement2
	SortTime(allEndDates, true)
	for _, v1 := range allEndDates {
		trailingBalanceSheet := StatementStep1(journalTimes, journal, v1.AddDate(0, 0, -int(periodInDaysBeforeEndDate)), v1)
		trailingBalanceSheet = StatementStep2(trailingBalanceSheet)
		trailingBalanceSheet = StatementStep3(trailingBalanceSheet)
		trailingBalanceSheet = StatementStep4(inNames, namesYouWant, trailingBalanceSheet)
		statement := StatementStep5(trailingBalanceSheet)
		StatementStep6(periodInDaysBeforeEndDate, statement)
		StatementStep7(statement)
		statements1 = append(statements1, statement)
	}

	var statements2 []statement3
	for _, v1 := range statements1 {
		statement := StatementStep8(v1, statements1[0])
		StatementStep9(statement)
		statements2 = append(statements2, statement)
	}

	return statements2, nil
}

func StatementStep1(journalTimes []time.Time, journal []Journal, dateStart, dateEnd time.Time) statement1 {
	// in this function we create the statement map

	newStatement := statement1{}
	for k1, v1 := range journal {
		switch {
		case journalTimes[k1].Before(dateStart):
			FillNewStatement(newStatement, v1, true)
		case journalTimes[k1].Before(dateEnd):
			FillNewStatement(newStatement, v1, false)
		default:
			break
		}
	}

	return newStatement
}

func StatementStep2(oldStatement statement1) statement1 {
	// in this function i insert the father accounts in column account1
	// i sum the credit to credit and debit to debit .Like:
	// if there is three accounts like this:
	// assets 			debit
	// 	|_equipment 	debit
	// 	|_depreciation 	credit
	// i will sum the debit side of the equipment and depreciation to debit side of assets
	// and the credit side of the equipment and depreciation to credit side of assets

	newStatement := statement1{}

	for k1, v1 := range oldStatement { //account1
		account, _, err := FindAccountFromName(k1)
		for k2, v2 := range v1 { //account2
			for k3, v3 := range v2 { //name
				for k4, v4 := range v3 { //vpq
					for k5, v5 := range v4 { //isBeforeDateStart
						for k6, v6 := range v5 { //is_credit
							// here i copy the map
							m := InitializeMap6(newStatement, k1, k2, k3, k4, k5)
							m[k6] += v6

							if err == nil {
								for _, v7 := range account.FathersName[IndexOfAccountNumber] {
									m = InitializeMap6(newStatement, v7, k2, k3, k4, k5)
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

func StatementStep3(oldStatement statement1) statement1 {
	// in this function i insert the father accounts and 'AllAccounts' key word in column account2
	// i sum the credit to credit and debit to debit .Like:
	// if there is three accounts like this:
	// assets debit,equipment debit,depreciation credit
	// i will sum the debit side of the equipment and depreciation to debit side of assets
	// and the credit side of the equipment and depreciation to credit side of assets

	newStatement := statement1{}

	for k1, v1 := range oldStatement { //account1
		for k2, v2 := range v1 { //account2
			account, _, _ := FindAccountFromName(k2)
			for k3, v3 := range v2 { //name
				for k4, v4 := range v3 { //vpq
					for k5, v5 := range v4 { //isBeforeDateStart
						for k6, v6 := range v5 { //is_credit
							// here i copy the map
							m := InitializeMap6(newStatement, k1, k2, k3, k4, k5)
							m[k6] += v6

							// here i insert the key word 'AllAccounts' in column account2
							m = InitializeMap6(newStatement, k1, AllAccounts, k3, k4, k5)
							m[k6] += v6

							for _, v7 := range account.FathersName[IndexOfAccountNumber] {
								m = InitializeMap6(newStatement, k1, v7, k3, k4, k5)
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

func StatementStep4(inNames bool, namesYouWant []string, oldStatement statement1) statement1 {
	// in this function i insert the key word 'AllNames' and 'Names' in column name

	newStatement := statement1{}

	for k1, v1 := range oldStatement { //account1
		for k2, v2 := range v1 { //account2
			for k3, v3 := range v2 { //name
				for k4, v4 := range v3 { //vpq
					for k5, v5 := range v4 { //isBeforeDateStart
						for k6, v6 := range v5 { //is_credit
							// here i copy the map
							m := InitializeMap6(newStatement, k1, k2, k3, k4, k5)
							m[k6] += v6

							// here i insert the key word 'AllNames' in column name
							m = InitializeMap6(newStatement, k1, k2, AllNames, k4, k5)
							m[k6] += v6

							if IsIn(k2, namesYouWant) == inNames {
								// here i insert the key word 'Names' in column name
								m = InitializeMap6(newStatement, k1, k2, Names, k4, k5)
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

func StatementStep5(oldStatement statement1) statement2 {
	// in this function i insert the type_of_vpq and remove column isBeforeDateStart and is_credit

	newStatement := statement2{}
	for k1, v1 := range oldStatement { //account1
		accountStruct1, _, _ := FindAccountFromName(k1)
		for k2, v2 := range v1 { //account2
			for k3, v3 := range v2 { //name
				for k4, v4 := range v3 { //vpq
					for k5, v5 := range v4 { //isBeforeDateStart
						for k6, v6 := range v5 { //is_credit
							if k5 {
								// here i insert FlowInBeginning and FlowOutBeginning
								if accountStruct1.IsCredit == k6 {
									m := InitializeMap5(newStatement, k1, k2, k3, k4)
									m[FlowInBeginning] += v6
								} else {
									m := InitializeMap5(newStatement, k1, k2, k3, k4)
									m[FlowOutBeginning] += v6
								}
							} else {
								// here i insert FlowInPeriod and FlowOutPeriod
								if accountStruct1.IsCredit == k6 {
									m := InitializeMap5(newStatement, k1, k2, k3, k4)
									m[FlowInPeriod] += v6
								} else {
									m := InitializeMap5(newStatement, k1, k2, k3, k4)
									m[FlowOutPeriod] += v6
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

func StatementStep6(days uint, oldStatement statement2) {
	// in this function we make vertical analysis of the statement

	for _, v1 := range oldStatement { //account1
		for _, v2 := range v1 { //account2
			for _, v3 := range v2 { //name
				for _, v4 := range v3 { //vpq
					v4[FlowInEnding] = v4[FlowInPeriod] + v4[FlowInBeginning]
					v4[FlowOutEnding] = v4[FlowOutPeriod] + v4[FlowOutBeginning]

					v4[FlowBeginning] = v4[FlowInBeginning] - v4[FlowOutBeginning]
					v4[FlowPeriod] = v4[FlowInPeriod] - v4[FlowOutPeriod]
					v4[FlowEnding] = v4[FlowInEnding] - v4[FlowOutEnding]

					v4[Average] = (v4[FlowEnding] + v4[FlowBeginning]) / 2
					v4[Turnover] = v4[FlowOutPeriod] / v4[Average]
					v4[TurnoverDays] = float64(days) / v4[Turnover]
					v4[GrowthRatio] = v4[FlowEnding] / v4[FlowBeginning]
				}
			}
		}
	}
}

func StatementStep7(oldStatement statement2) {
	// in this function i complete vertical analysis of the statement
	// but here i calculate the percentage of the name from the account

	for _, v1 := range oldStatement { //account1
		for _, v2 := range v1 { //account2
			for _, v3 := range v2 { //name
				for k4, v4 := range v3 { //vpq
					v4[NamePercent] = v4[FlowEnding] / v2[AllNames][k4][FlowEnding]
				}
			}
		}
	}
}

func StatementStep8(oldStatement, baseStatement statement2) statement3 {
	newStatement := statement3{}

	for k1, v1 := range oldStatement { //account1
		for k2, v2 := range v1 { //account2
			for k3, v3 := range v2 { //name
				for k4, v4 := range v3 { //vpq
					for k5, v5 := range v4 { //type_of_vpq
						m := InitializeMap6(newStatement, k1, k2, k3, k4, k5)
						m[Balance] = v5
						m[ChangeSinceBasePeriod] = v5 - baseStatement[k1][k2][k3][k4][k5]
						m[GrowthRatioToBasePeriod] = v5 / baseStatement[k1][k2][k3][k4][k5]
					}
				}
			}
		}
	}
	return newStatement
}

func StatementStep9(oldStatement statement3) {
	// in this function we calculate the Price

	for _, v1 := range oldStatement { //account1
		for _, v2 := range v1 { //account2
			for _, v3 := range v2 { //name
				if v3[Price] == nil {
					v3[Price] = map[string]map[string]float64{}
				}
				for _, v4 := range v3 { //vpq
					for k5, v5 := range v4 { //type_of_vpq
						if v3[Price][k5] == nil {
							v3[Price][k5] = map[string]float64{}
						}
						for k6 := range v5 { //Change_or_Ratio_or_balance
							v3[Price][k5][k6] = v3[Value][k5][k6] / v3[Quantity][k5][k6]
						}
					}
				}
			}
		}
	}
}

func FillNewStatement(newStatement statement1, v1 Journal, isBeforeDateStart bool) {
	m := InitializeMap6(newStatement, v1.AccountCredit, v1.AccountDebit, v1.Name, Value, isBeforeDateStart)
	m[true] += v1.Value
	m = InitializeMap6(newStatement, v1.AccountCredit, v1.AccountDebit, v1.Name, Quantity, isBeforeDateStart)
	m[true] += Abs(v1.QuantityCredit)
	m = InitializeMap6(newStatement, v1.AccountDebit, v1.AccountCredit, v1.Name, Value, isBeforeDateStart)
	m[false] += v1.Value
	m = InitializeMap6(newStatement, v1.AccountDebit, v1.AccountCredit, v1.Name, Quantity, isBeforeDateStart)
	m[false] += Abs(v1.QuantityDebit)
}

func StatementFilter(oldStatement statement3, f FilterStatement) []StatmentWithAccount {
	var newStatement []StatmentWithAccount

	for k1, v1 := range oldStatement { //account1
		if f.Account1.Account.Filter(k1) {
			account1, _, err := FindAccountFromName(k1)
			if f.Account1.Filter(account1, err) {
				for k2, v2 := range v1 { //account2
					if f.Account2.Account.Filter(k2) {
						account2, _, err := FindAccountFromName(k2)
						if f.Account2.Filter(account1, err) {
							for k3, v3 := range v2 { //name
								if f.Name.Filter(k3) {
									for k4, v4 := range v3 { //vpq
										if f.Vpq.Filter(k4) {
											for k5, v5 := range v4 { //type_of_vpq
												if f.TypeOfVpq.Filter(k5) {
													for k6, v6 := range v5 { //Change_or_Ratio_or_balance
														if f.ChangeOrRatioOrBalance.Filter(k6) {
															if f.Number.Filter(v6) {

																// here i prefer to show the account struct in the statment to use it later in sorting the account
																newStatement = append(newStatement, StatmentWithAccount{
																	Account1: account1,
																	Account2: account2,
																	Statment: Statement{k1, k2, k3, k4, k5, k6, v6},
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

func StatementFilterByGreedyAlgorithm(oldStatement []StatmentWithAccount, isPercent bool, numberTarget float64) []StatmentWithAccount {
	if numberTarget == 0 {
		return oldStatement
	}

	var numberTotal float64
	if isPercent {
		// here i sum the numbers to find the total amount to calculate the percentage in units
		for _, v1 := range oldStatement {
			numberTotal += v1.Statment.Number
		}
		// here i convert the percent to a units
		numberTarget = numberTotal * numberTarget
	}

	SortStatementNumber(oldStatement, true)
	var newStatement []StatmentWithAccount
	var numberCurrent float64
	for _, v1 := range oldStatement {
		if numberCurrent < numberTarget {
			newStatement = append(newStatement, v1)
			numberCurrent += v1.Statment.Number
		}
	}
	return newStatement
}

func StatementSort(statement []StatmentWithAccount, way string) []StatmentWithAccount {
	switch way {
	case Ascending:
		SortStatementNumber(statement, true)
		return statement
	case Descending:
		SortStatementNumber(statement, false)
		return statement
	default:
		SortByLevel(statement)
		return statement
	}
}

func SortByLevel(s []StatmentWithAccount) []StatmentWithAccount {
	for k1 := range s {
		for k2 := range s {
			if k1 < k2 &&
				!IsItHighThanByOrder(s[k1].Account1.Number[IndexOfAccountNumber], s[k2].Account1.Number[IndexOfAccountNumber]) {
				Swap(s, k1, k2) // account1

				if s[k1].Statment.Account1 == s[k2].Statment.Account1 &&
					(s[k1].Statment.Account2 == AllAccounts || s[k2].Statment.Account2 == AllAccounts ||
						!IsItHighThanByOrder(s[k1].Account2.Number[IndexOfAccountNumber], s[k2].Account2.Number[IndexOfAccountNumber])) {
					Swap(s, k1, k2) // account2

					if s[k1].Statment.Account2 == s[k2].Statment.Account2 &&
						s[k1].Statment.Name > s[k2].Statment.Name {
						Swap(s, k1, k2) // name

						if s[k1].Statment.Name == s[k2].Statment.Name &&
							s[k1].Statment.Vpq > s[k2].Statment.Vpq {
							Swap(s, k1, k2) // vpq

							if s[k1].Statment.Vpq == s[k2].Statment.Vpq &&
								s[k1].Statment.TypeOfVpq > s[k2].Statment.TypeOfVpq {
								Swap(s, k1, k2) // typeOfVpq

								if s[k1].Statment.TypeOfVpq == s[k2].Statment.TypeOfVpq &&
									s[k1].Statment.ChangeOrRatioOrBalance > s[k2].Statment.ChangeOrRatioOrBalance {
									Swap(s, k1, k2) // ChangeOrRatioOrBalance
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

func MakeSpaceBeforeAccountInStatementStruct(oldStatement []StatmentWithAccount) {
	for k1, v1 := range oldStatement {
		oldStatement[k1].Statment.Account1 = strings.Repeat("  ", int(v1.Account1.Levels[IndexOfAccountNumber])) + v1.Statment.Account1
		if v1.Statment.Account2 != AllAccounts {
			oldStatement[k1].Statment.Account2 = strings.Repeat("  ", int(v1.Account2.Levels[IndexOfAccountNumber])) + v1.Statment.Account2
		}
	}
}

func ConvertStatmentWithAccountToFilteredStatement(oldStatement []StatmentWithAccount) []Statement {
	var newStatement []Statement
	for _, v1 := range oldStatement {
		newStatement = append(newStatement, v1.Statment)
	}
	return newStatement
}

func StatementAnalysis(i FinancialAnalysis) FinancialAnalysisStatement {
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
	return FinancialAnalysisStatement{
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
