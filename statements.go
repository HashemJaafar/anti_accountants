package main

import (
	"fmt"
	"time"
)

func FinancialStatements(allEndDates []time.Time, periodInDaysBeforeEndDate uint, namesYouWant []string, inNames bool) ([]map[string]map[string]map[string]map[string]map[string]float64, error) {
	var statements []map[string]map[string]map[string]map[string]map[string]float64

	SortTime(allEndDates, true)

	keys, journal := DbRead[JournalTag](DbJournal)
	journalTimes := ConvertByteSliceToTime(keys)

	for _, v1 := range allEndDates {
		trailingBalanceSheet := StatementStep1(journalTimes, journal, v1.AddDate(0, 0, -int(periodInDaysBeforeEndDate)), v1)
		trailingBalanceSheet = StatementStep2(trailingBalanceSheet)
		trailingBalanceSheet = StatementStep3(trailingBalanceSheet)
		trailingBalanceSheet = StatementStep4(trailingBalanceSheet)
		trailingBalanceSheet = StatementStep5(inNames, namesYouWant, trailingBalanceSheet)
		statement := StatementStep6(trailingBalanceSheet)
		StatementStep7(periodInDaysBeforeEndDate, statement)
		StatementStep8(statement)
		statements = append(statements, statement)
	}

	for _, v1 := range statements {
		HorizontalAnalysis(v1, statements[0])
		// 		prepare_statement(statement_current)
		CalculatePrice(v1)
	}

	return statements, nil
}

func StatementStep1(journalTimes []time.Time, journal []JournalTag, dateStart, dateEnd time.Time) map[string]map[string]map[string]map[string]map[bool]map[bool]float64 {
	// in this function we create the statement map

	// the sequanse of the columns is:account1,account2,name,vpq,isBeforeDateStart,is_credit,number
	newStatement := map[string]map[string]map[string]map[string]map[bool]map[bool]float64{}

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

func StatementStep2(oldStatement map[string]map[string]map[string]map[string]map[bool]map[bool]float64) map[string]map[string]map[string]map[string]map[bool]map[bool]float64 {
	// in this function i insert retained earnings account in column account 1

	newStatement := map[string]map[string]map[string]map[string]map[bool]map[bool]float64{}

	for k1, v1 := range oldStatement { //account1
		accountStruct, _, _ := AccountStructFromName(k1)
		for k2, v2 := range v1 { //account2
			for k3, v3 := range v2 { //name
				for k4, v4 := range v3 { //vpq
					for k5, v5 := range v4 { //isBeforeDateStart
						for k6, v6 := range v5 { //is_credit
							// if the account is temporary account and the entry is before the date start
							// then i insert retained earnings account in column account 1
							// else i dont do anything
							if accountStruct.IsTemporary && k5 {
								m := InitializeMap6(newStatement, RetinedEarnings, k2, k3, k4, k5)
								m[k6] += v6
							} else {
								// here i copy the map
								m := InitializeMap6(newStatement, k1, k2, k3, k4, k5)
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

func StatementStep3(oldStatement map[string]map[string]map[string]map[string]map[bool]map[bool]float64) map[string]map[string]map[string]map[string]map[bool]map[bool]float64 {
	// in this function i insert the father accounts in column account1
	// i sum the credit to credit and debit to debit .Like:
	// if there is three accounts like this:
	// assets 			debit
	// 	|_equipment 	debit
	// 	|_depreciation 	credit
	// i will sum the debit side of the equipment and depreciation to debit side of assets
	// and the credit side of the equipment and depreciation to credit side of assets

	newStatement := map[string]map[string]map[string]map[string]map[bool]map[bool]float64{}

	for k1, v1 := range oldStatement { //account1
		accountStruct, _, err := AccountStructFromName(k1)
		for k2, v2 := range v1 { //account2
			for k3, v3 := range v2 { //name
				for k4, v4 := range v3 { //vpq
					for k5, v5 := range v4 { //isBeforeDateStart
						for k6, v6 := range v5 { //is_credit
							// here i copy the map
							m := InitializeMap6(newStatement, k1, k2, k3, k4, k5)
							m[k6] += v6

							if err == nil {
								for _, v7 := range accountStruct.FathersAccountsName[IndexOfAccountNumber] {
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

func StatementStep4(oldStatement map[string]map[string]map[string]map[string]map[bool]map[bool]float64) map[string]map[string]map[string]map[string]map[bool]map[bool]float64 {
	// in this function i insert the father accounts and 'AllAccounts' key word in column account2
	// i sum the credit to credit and debit to debit .Like:
	// if there is three accounts like this:
	// assets debit,equipment debit,depreciation credit
	// i will sum the debit side of the equipment and depreciation to debit side of assets
	// and the credit side of the equipment and depreciation to credit side of assets

	newStatement := map[string]map[string]map[string]map[string]map[bool]map[bool]float64{}

	for k1, v1 := range oldStatement { //account1
		for k2, v2 := range v1 { //account2
			accountStruct, _, _ := AccountStructFromName(k2)
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

							for _, v7 := range accountStruct.FathersAccountsName[IndexOfAccountNumber] {
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

func StatementStep5(inNames bool, namesYouWant []string, oldStatement map[string]map[string]map[string]map[string]map[bool]map[bool]float64) map[string]map[string]map[string]map[string]map[bool]map[bool]float64 {
	// in this function i insert the key word 'AllNames' and 'Names' in column name

	newStatement := map[string]map[string]map[string]map[string]map[bool]map[bool]float64{}

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

func StatementStep6(oldStatement map[string]map[string]map[string]map[string]map[bool]map[bool]float64) map[string]map[string]map[string]map[string]map[string]float64 {
	// in this function i insert the type_of_vpq and remove column is before_date_start and is_credit

	// the sequanse of the columns is:account1,account2,name,vpq,type_of_vpq,number
	newStatement := map[string]map[string]map[string]map[string]map[string]float64{}

	for k1, v1 := range oldStatement { //account1
		accountStruct1, _, _ := AccountStructFromName(k1)
		for k2, v2 := range v1 { //account2
			for k3, v3 := range v2 { //name
				for k4, v4 := range v3 { //vpq
					for k5, v5 := range v4 { //isBeforeDateStart
						for k6, v6 := range v5 { //is_credit
							if k5 {
								// here i insert BeginningBalance
								if accountStruct1.IsCredit == k6 {
									m := InitializeMap5(newStatement, k1, k2, k3, k4)
									m[BeginningBalance] += v6
								} else {
									m := InitializeMap5(newStatement, k1, k2, k3, k4)
									m[BeginningBalance] -= v6
								}
							} else {
								// here i insert Inflow and Outflow
								if accountStruct1.IsCredit == k6 {
									m := InitializeMap5(newStatement, k1, k2, k3, k4)
									m[Inflow] += v6
								} else {
									m := InitializeMap5(newStatement, k1, k2, k3, k4)
									m[Outflow] += v6
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

func StatementStep7(days uint, oldStatement map[string]map[string]map[string]map[string]map[string]float64) {
	// in this function we make vertical analysis of the statement

	// the sequanse of the columns is:account1,account2,name,vpq,type_of_vpq,number
	for _, v1 := range oldStatement { //account1
		for _, v2 := range v1 { //account2
			for _, v3 := range v2 { //name
				for _, v4 := range v3 { //vpq
					v4[Flow] = v4[Inflow] - v4[Outflow]
					// here i calculate the ending balance of the account by sum the Flow with BeginningBalance
					// because the Inflow is the same of the increase of the account
					// and Outflow is the same of the decrease of the account
					v4[EndingBalance] = v4[BeginningBalance] + v4[Flow]
					v4[Average] = (v4[EndingBalance] + v4[BeginningBalance]) / 2
					v4[Turnover] = v4[Outflow] / v4[Average]
					v4[TurnoverDays] = float64(days) / v4[Turnover]
					v4[GrowthRatio] = v4[EndingBalance] / v4[BeginningBalance]
				}
			}
		}
	}
}

func StatementStep8(oldStatement map[string]map[string]map[string]map[string]map[string]float64) {
	// in this function i complete vertical analysis of the statement
	// but here i calculate the percentage of the account from account father

	// the sequanse of the columns is:account1,account2,name,vpq,type_of_vpq,number
	for _, v1 := range oldStatement { //account1
		for _, v2 := range v1 { //account2
			for _, v3 := range v2 { //name
				for k4, v4 := range v3 { //vpq
					v4[NamePercent] = v4[EndingBalance] / v2[AllNames][k4][EndingBalance]
				}
			}
		}
	}
}

func HorizontalAnalysis(oldStatement, baseStatement map[string]map[string]map[string]map[string]map[string]float64) {
	// the sequanse of the columns is:account1,account2,name,vpq,type_of_vpq,number
	for k1, v1 := range oldStatement { //account1
		for k2, v2 := range v1 { //account2
			for k3, v3 := range v2 { //name
				for k4, v4 := range v3 { //vpq
					v4[ChangeSinceBasePeriod] = v4[EndingBalance] - baseStatement[k1][k2][k3][k4][EndingBalance]
					v4[GrowthRatioToBasePeriod] = v4[EndingBalance] / baseStatement[k1][k2][k3][k4][EndingBalance]
				}
			}
		}
	}
}

func CalculatePrice(oldStatement map[string]map[string]map[string]map[string]map[string]float64) {
	// in this function we calculate the Price

	// the sequanse of the columns is:account1,account2,name,vpq,type_of_vpq,number
	for _, v1 := range oldStatement { //account1
		for _, v2 := range v1 { //account2
			for _, v3 := range v2 { //name
				if v3[Price] == nil {
					v3[Price] = map[string]float64{}
				}
				for _, v4 := range v3 { //vpq
					for k5 := range v4 { //type_of_vpq
						v3[Price][k5] = v3[Value][k5] / v3[Quantity][k5]
					}
				}
			}
		}
	}
}

func FillNewStatement(newStatement map[string]map[string]map[string]map[string]map[bool]map[bool]float64, v1 JournalTag, isBeforeDateStart bool) {
	m := InitializeMap6(newStatement, v1.AccountCredit, v1.AccountDebit, v1.Name, Value, isBeforeDateStart)
	m[true] += v1.Value
	m = InitializeMap6(newStatement, v1.AccountCredit, v1.AccountDebit, v1.Name, Quantity, isBeforeDateStart)
	m[true] += Abs(v1.QuantityCredit)
	m = InitializeMap6(newStatement, v1.AccountDebit, v1.AccountCredit, v1.Name, Value, isBeforeDateStart)
	m[false] += v1.Value
	m = InitializeMap6(newStatement, v1.AccountDebit, v1.AccountCredit, v1.Name, Quantity, isBeforeDateStart)
	m[false] += Abs(v1.QuantityDebit)
}

func StatementFilter(oldStatement map[string]map[string]map[string]map[string]map[string]float64, f FilterStatement) []FilteredStatement {
	var newStatement []FilteredStatement

	// the sequanse of the columns is:account1,account2,name,vpq,type_of_vpq,number
	for k1, v1 := range oldStatement { //account1
		if f.Account1.Filter(k1) {
			for k2, v2 := range v1 { //account2
				if f.Account2.Filter(k2) {
					for k3, v3 := range v2 { //name
						if f.Name.Filter(k3) {
							for k4, v4 := range v3 { //vpq
								if f.Vpq.Filter(k4) {
									for k5, v5 := range v4 { //type_of_vpq
										if f.TypeOfVpq.Filter(k5) && f.Number.Filter(v5) {
											newStatement = append(newStatement, FilteredStatement{k1, k2, k3, k4, k5, v5})
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

func StatementFilterAccounts(oldStatement []FilteredStatement, f1 FilterAccount, f2 FilterAccount) []FilteredStatement {
	var newStatement []FilteredStatement
	for _, v1 := range oldStatement {
		account1, _, _ := AccountStructFromName(v1.Account1)
		if f1.IsLowLevelAccount.Filter(account1.IsLowLevelAccount) &&
			f1.IsCredit.Filter(account1.IsCredit) &&
			f1.IsTemporary.Filter(account1.IsTemporary) &&
			f1.FathersAccountsName.Filter(account1.FathersAccountsName[IndexOfAccountNumber]) &&
			f1.AccountLevels.Filter(account1.AccountLevels[IndexOfAccountNumber]) {

			if v1.Account2 == AllAccounts {
				newStatement = append(newStatement, v1)
				continue
			}

			account2, _, _ := AccountStructFromName(v1.Account2)
			if f2.IsLowLevelAccount.Filter(account2.IsLowLevelAccount) &&
				f2.IsCredit.Filter(account2.IsCredit) &&
				f2.IsTemporary.Filter(account2.IsTemporary) &&
				f2.FathersAccountsName.Filter(account1.FathersAccountsName[IndexOfAccountNumber]) &&
				f2.AccountLevels.Filter(account2.AccountLevels[IndexOfAccountNumber]) {

				newStatement = append(newStatement, v1)
			}
		}
	}
	return newStatement
}

func StatementFilterByGreedyAlgorithm(oldStatement []FilteredStatement, isPercent bool, targetUnits float64) []FilteredStatement {
	if targetUnits == 0 {
		return oldStatement
	}

	var totalUnits float64
	if isPercent {
		// here i sum the numbers to find the total amount to calculate the percentage in units
		for _, v1 := range oldStatement {
			totalUnits += v1.Number
		}
		// here i convert the percent to a units
		targetUnits = totalUnits * targetUnits
	}

	SortStatementNumber(oldStatement, true)
	var newStatement []FilteredStatement
	var currentUnits float64
	for _, v1 := range oldStatement {
		if currentUnits < targetUnits {
			newStatement = append(newStatement, v1)
			currentUnits += v1.Number
		}
	}
	return newStatement
}

func StatementSort(oldStatement []FilteredStatement, way string) []FilteredStatement {
	switch way {
	case "ascending":
		SortStatementNumber(oldStatement, true)
		return oldStatement
	case "descending":
		SortStatementNumber(oldStatement, false)
		return oldStatement
	default:
		SortByLevel(oldStatement)
		return oldStatement
	}
}

func SortByLevel(oldStatement []FilteredStatement) []FilteredStatement {
	type statmentWithAccount struct {
		Account1 Account
		Account2 Account
		Statment FilteredStatement
	}

	var swa []statmentWithAccount
	for _, v1 := range oldStatement {
		account1, _, _ := AccountStructFromName(v1.Account1)
		account2, _, _ := AccountStructFromName(v1.Account2)
		swa = append(swa, statmentWithAccount{account1, account2, v1})
	}

	for k1 := range swa {
		for k2 := range swa {
			if k1 < k2 &&
				!IsItHighThanByOrder(swa[k1].Account1.AccountNumber[IndexOfAccountNumber], swa[k2].Account1.AccountNumber[IndexOfAccountNumber]) {
				fmt.Println(swa[k1].Account1.AccountNumber[IndexOfAccountNumber], swa[k2].Account1.AccountNumber[IndexOfAccountNumber])
				Swap(swa, k1, k2)
			}
		}
	}

	var newStatement []FilteredStatement
	for _, v1 := range swa {
		newStatement = append(newStatement, v1.Statment)
	}

	return newStatement
}

func StatementAnalysis(s FinancialAnalysis) FinancialAnalysisStatement {
	currentRatio := s.CurrentAssets / s.CurrentLiabilities
	acidTest := (s.Cash + s.ShortTermInvestments + s.NetReceivables) / s.CurrentLiabilities
	receivablesTurnover := s.NetCreditSales / s.AverageNetReceivables
	inventoryTurnover := s.CostOfGoodsSold / s.AverageInventory
	profitMargin := s.NetIncome / s.NetSales
	assetTurnover := s.NetSales / s.AverageAssets
	returnOnAssets := s.NetIncome / s.AverageAssets
	returnOnEquity := s.NetIncome / s.AverageEquity
	payoutRatio := s.CashDividends / s.NetIncome
	debtToTotalAssetsRatio := s.TotalDebt / s.TotalAssets
	timesInterestEarned := s.Ebitda / s.InterestExpense
	returnOnCommonStockholdersEquity := (s.NetIncome - s.PreferredDividends) / s.AverageCommonStockholdersEquity
	earningsPerShare := (s.NetIncome - s.PreferredDividends) / s.WeightedAverageCommonSharesOutstanding
	priceEarningsRatio := s.MarketPricePerSharesOutstanding / earningsPerShare
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

// func ANALYSIS(statements []map[string]map[string]map[string]map[string]map[string]float64) []FinancialAnalysisStatement {
// 	var all_analysis []FinancialAnalysisStatement
// 	for _, statement := range statements {
// 		analysis := FINANCIAL_ANALYSIS_STATEMENT_func(FinancialAnalysis{
// 			CURRENT_ASSETS:                      statement[PRIMARY_ACCOUNTS_NAMES.CASH_AND_CASH_EQUIVALENTS][PRIMARY_ACCOUNTS_NAMES.CURRENT_ASSETS]["names"]["VALUE"]["ending_balance"],
// 			CURRENT_LIABILITIES:                 statement[PRIMARY_ACCOUNTS_NAMES.CASH_AND_CASH_EQUIVALENTS][PRIMARY_ACCOUNTS_NAMES.CURRENT_LIABILITIES]["names"]["VALUE"]["ending_balance"],
// 			CASH:                                statement[PRIMARY_ACCOUNTS_NAMES.CASH_AND_CASH_EQUIVALENTS][PRIMARY_ACCOUNTS_NAMES.CASH_AND_CASH_EQUIVALENTS]["names"]["VALUE"]["ending_balance"],
// 			SHORT_TERM_INVESTMENTS:              statement[PRIMARY_ACCOUNTS_NAMES.CASH_AND_CASH_EQUIVALENTS][PRIMARY_ACCOUNTS_NAMES.SHORT_TERM_INVESTMENTS]["names"]["VALUE"]["ending_balance"],
// 			NET_RECEIVABLES:                     statement[PRIMARY_ACCOUNTS_NAMES.CASH_AND_CASH_EQUIVALENTS][PRIMARY_ACCOUNTS_NAMES.RECEIVABLES]["names"]["VALUE"]["ending_balance"],
// 			NET_CREDIT_SALES:                    statement[PRIMARY_ACCOUNTS_NAMES.SALES][PRIMARY_ACCOUNTS_NAMES.RECEIVABLES]["names"]["VALUE"]["flow"],
// 			AVERAGE_NET_RECEIVABLES:             statement[PRIMARY_ACCOUNTS_NAMES.CASH_AND_CASH_EQUIVALENTS][PRIMARY_ACCOUNTS_NAMES.RECEIVABLES]["names"]["VALUE"]["average"],
// 			COST_OF_GOODS_SOLD:                  statement[PRIMARY_ACCOUNTS_NAMES.CASH_AND_CASH_EQUIVALENTS][PRIMARY_ACCOUNTS_NAMES.COST_OF_GOODS_SOLD]["names"]["VALUE"]["ending_balance"],
// 			AVERAGE_INVENTORY:                   statement[PRIMARY_ACCOUNTS_NAMES.CASH_AND_CASH_EQUIVALENTS][PRIMARY_ACCOUNTS_NAMES.INVENTORY]["names"]["VALUE"]["average"],
// 			NET_INCOME:                          statement[PRIMARY_ACCOUNTS_NAMES.CASH_AND_CASH_EQUIVALENTS][PRIMARY_ACCOUNTS_NAMES.INCOME_STATEMENT]["names"]["VALUE"]["ending_balance"],
// 			NET_SALES:                           statement[PRIMARY_ACCOUNTS_NAMES.CASH_AND_CASH_EQUIVALENTS][PRIMARY_ACCOUNTS_NAMES.SALES]["names"]["VALUE"]["ending_balance"],
// 			AVERAGE_ASSETS:                      statement[PRIMARY_ACCOUNTS_NAMES.CASH_AND_CASH_EQUIVALENTS][PRIMARY_ACCOUNTS_NAMES.ASSETS]["names"]["VALUE"]["average"],
// 			AVERAGE_EQUITY:                      statement[PRIMARY_ACCOUNTS_NAMES.CASH_AND_CASH_EQUIVALENTS][PRIMARY_ACCOUNTS_NAMES.EQUITY]["names"]["VALUE"]["average"],
// 			PREFERRED_DIVIDENDS:                 0,
// 			AVERAGE_COMMON_STOCKHOLDERS_EQUITY:  0,
// 			MARKET_PRICE_PER_SHARES_OUTSTANDING: 0,
// 			CASH_DIVIDENDS:                      statement[PRIMARY_ACCOUNTS_NAMES.CASH_AND_CASH_EQUIVALENTS][PRIMARY_ACCOUNTS_NAMES.DIVIDENDS]["names"]["VALUE"]["flow"],
// 			TOTAL_DEBT:                          statement[PRIMARY_ACCOUNTS_NAMES.CASH_AND_CASH_EQUIVALENTS][PRIMARY_ACCOUNTS_NAMES.LIABILITIES]["names"]["VALUE"]["ending_balance"],
// 			TOTAL_ASSETS:                        statement[PRIMARY_ACCOUNTS_NAMES.CASH_AND_CASH_EQUIVALENTS][PRIMARY_ACCOUNTS_NAMES.ASSETS]["names"]["VALUE"]["ending_balance"],
// 			EBITDA:                              statement[PRIMARY_ACCOUNTS_NAMES.CASH_AND_CASH_EQUIVALENTS][PRIMARY_ACCOUNTS_NAMES.EBITDA]["names"]["VALUE"]["ending_balance"],
// 			INTEREST_EXPENSE:                    statement[PRIMARY_ACCOUNTS_NAMES.CASH_AND_CASH_EQUIVALENTS][PRIMARY_ACCOUNTS_NAMES.INTEREST_EXPENSE]["names"]["VALUE"]["ending_balance"],
// 			WEIGHTED_AVERAGE_COMMON_SHARES_OUTSTANDING: 0,
// 		})
// 		all_analysis = append(all_analysis, analysis)
// 	}
// 	return all_analysis
// }
