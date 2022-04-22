package main

import (
	"fmt"
	"time"
)

func FINANCIAL_STATEMENTS(all_end_dates []time.Time, period_in_days_before_end_date uint, names_you_want []string, in_names bool, retained_earnings_account ACCOUNT) ([]map[string]map[string]map[string]map[string]map[string]float64, error) {
	var statements []map[string]map[string]map[string]map[string]map[string]float64

	// if the account name is used then the function is abort
	_, _, err := ACCOUNT_STRUCT_FROM_NAME(retained_earnings_account.ACCOUNT_NAME)
	if err == nil {
		return statements, ERROR_ACCOUNT_NAME_IS_USED
	}

	// here i store the ACCOUNTS to old_accounts to restore after add the retained_earnings_account
	old_accounts := ACCOUNTS

	retained_earnings_account = SET_RETAINED_EARNINGS_ACCOUNT(retained_earnings_account)
	ACCOUNTS = append(ACCOUNTS, retained_earnings_account)
	SET_THE_ACCOUNTS()

	SORT_TIME(all_end_dates, true)

	keys, journal := DB_READ[JOURNAL_TAG](DB_JOURNAL)
	journal_times := CONVERT_BYTE_SLICE_TO_TIME(keys)

	for _, v1 := range all_end_dates {
		trailing_balance_sheet := STATEMENT_STEP_1(journal_times, journal, v1.AddDate(0, 0, -int(period_in_days_before_end_date)), v1)
		trailing_balance_sheet = STATEMENT_STEP_2(trailing_balance_sheet, retained_earnings_account.ACCOUNT_NAME)
		trailing_balance_sheet = STATEMENT_STEP_3(trailing_balance_sheet)
		trailing_balance_sheet = STATEMENT_STEP_4(trailing_balance_sheet)
		trailing_balance_sheet = STATEMENT_STEP_5(in_names, names_you_want, trailing_balance_sheet)
		statement := STATEMENT_STEP_6(trailing_balance_sheet)
		STATEMENT_STEP_7(period_in_days_before_end_date, statement)
		STATEMENT_STEP_8(statement)
		statements = append(statements, statement)
	}

	for _, v1 := range statements {
		HORIZONTAL_ANALYSIS(v1, statements[0])
		// 		prepare_statement(statement_current)
		CALCULATE_PRICE(v1)
	}

	ACCOUNTS = old_accounts
	return statements, nil
}

func STATEMENT_STEP_1(journal_times []time.Time, journal []JOURNAL_TAG, date_start, date_end time.Time) map[string]map[string]map[string]map[string]map[bool]map[bool]float64 {
	// in this function we create the statement map

	// the sequanse of the columns is:account1,account2,name,vpq,is_before_date_start,is_credit,number
	new_statement := map[string]map[string]map[string]map[string]map[bool]map[bool]float64{}

	for k1, v1 := range journal {
		switch {
		case journal_times[k1].Before(date_start):
			fill_new_statement(new_statement, v1, true)
		case journal_times[k1].Before(date_end):
			fill_new_statement(new_statement, v1, false)
		default:
			break
		}
	}

	return new_statement
}

func STATEMENT_STEP_2(old_statement map[string]map[string]map[string]map[string]map[bool]map[bool]float64, retained_earnings_account_name string) map[string]map[string]map[string]map[string]map[bool]map[bool]float64 {
	// in this function i insert retained earnings account in column account 1

	new_statement := map[string]map[string]map[string]map[string]map[bool]map[bool]float64{}

	for k1, v1 := range old_statement { //account1
		account_struct, _, _ := ACCOUNT_STRUCT_FROM_NAME(k1)
		for k2, v2 := range v1 { //account2
			for k3, v3 := range v2 { //name
				for k4, v4 := range v3 { //vpq
					for k5, v5 := range v4 { //is_before_date_start
						for k6, v6 := range v5 { //is_credit
							// if the account is temporary account and the entry is before the date start
							// then i insert retained earnings account in column account 1
							// else i dont do anything
							if account_struct.IS_TEMPORARY && k5 {
								m := INITIALIZE_MAP_6(new_statement, retained_earnings_account_name, k2, k3, k4, k5)
								m[k6] += v6
							} else {
								// here i copy the map
								m := INITIALIZE_MAP_6(new_statement, k1, k2, k3, k4, k5)
								m[k6] += v6
							}
						}
					}
				}
			}
		}
	}

	return new_statement
}

func STATEMENT_STEP_3(old_statement map[string]map[string]map[string]map[string]map[bool]map[bool]float64) map[string]map[string]map[string]map[string]map[bool]map[bool]float64 {
	// in this function i insert the father accounts in column account1
	// i sum the credit to credit and debit to debit . like:
	// if there is three accounts like this:
	// assets debit,equipment debit,depreciation credit
	// i will sum the debit side of the equipment and depreciation to debit side of assets
	// and the credit side of the equipment and depreciation to credit side of assets

	new_statement := map[string]map[string]map[string]map[string]map[bool]map[bool]float64{}

	for k1, v1 := range old_statement { //account1
		account_struct, _, _ := ACCOUNT_STRUCT_FROM_NAME(k1)
		for k2, v2 := range v1 { //account2
			for k3, v3 := range v2 { //name
				for k4, v4 := range v3 { //vpq
					for k5, v5 := range v4 { //is_before_date_start
						for k6, v6 := range v5 { //is_credit
							// here i copy the map
							m := INITIALIZE_MAP_6(new_statement, k1, k2, k3, k4, k5)
							m[k6] += v6

							for _, v7 := range account_struct.FATHER_AND_GRANDPA_ACCOUNTS_NAME[INDEX_OF_ACCOUNT_NUMBER] {
								m = INITIALIZE_MAP_6(new_statement, v7, k2, k3, k4, k5)
								m[k6] += v6
							}
						}
					}
				}
			}
		}
	}

	return new_statement
}

func STATEMENT_STEP_4(old_statement map[string]map[string]map[string]map[string]map[bool]map[bool]float64) map[string]map[string]map[string]map[string]map[bool]map[bool]float64 {
	// in this function i insert the father accounts and 'all_accounts' key word in column account2
	// i sum the credit to credit and debit to debit . like:
	// if there is three accounts like this:
	// assets debit,equipment debit,depreciation credit
	// i will sum the debit side of the equipment and depreciation to debit side of assets
	// and the credit side of the equipment and depreciation to credit side of assets

	new_statement := map[string]map[string]map[string]map[string]map[bool]map[bool]float64{}

	for k1, v1 := range old_statement { //account1
		for k2, v2 := range v1 { //account2
			account_struct, _, _ := ACCOUNT_STRUCT_FROM_NAME(k2)
			for k3, v3 := range v2 { //name
				for k4, v4 := range v3 { //vpq
					for k5, v5 := range v4 { //is_before_date_start
						for k6, v6 := range v5 { //is_credit
							// here i copy the map
							m := INITIALIZE_MAP_6(new_statement, k1, k2, k3, k4, k5)
							m[k6] += v6

							// here i insert the key word 'all_accounts' in column account2
							m = INITIALIZE_MAP_6(new_statement, k1, all_accounts, k3, k4, k5)
							m[k6] += v6

							for _, v7 := range account_struct.FATHER_AND_GRANDPA_ACCOUNTS_NAME[INDEX_OF_ACCOUNT_NUMBER] {
								m = INITIALIZE_MAP_6(new_statement, k1, v7, k3, k4, k5)
								m[k6] += v6
							}
						}
					}
				}
			}
		}
	}

	return new_statement
}

func STATEMENT_STEP_5(in_names bool, names_you_want []string, old_statement map[string]map[string]map[string]map[string]map[bool]map[bool]float64) map[string]map[string]map[string]map[string]map[bool]map[bool]float64 {
	// in this function i insert the key word 'all_names' and 'names' in column name

	new_statement := map[string]map[string]map[string]map[string]map[bool]map[bool]float64{}

	for k1, v1 := range old_statement { //account1
		for k2, v2 := range v1 { //account2
			for k3, v3 := range v2 { //name
				for k4, v4 := range v3 { //vpq
					for k5, v5 := range v4 { //is_before_date_start
						for k6, v6 := range v5 { //is_credit
							// here i copy the map
							m := INITIALIZE_MAP_6(new_statement, k1, k2, k3, k4, k5)
							m[k6] += v6

							// here i insert the key word 'all_names' in column name
							m = INITIALIZE_MAP_6(new_statement, k1, k2, all_names, k4, k5)
							m[k6] += v6

							if IS_IN(k2, names_you_want) == in_names {
								// here i insert the key word 'names' in column name
								m = INITIALIZE_MAP_6(new_statement, k1, k2, names, k4, k5)
								m[k6] += v6
							}
						}
					}
				}
			}
		}
	}

	return new_statement
}

func STATEMENT_STEP_6(old_statement map[string]map[string]map[string]map[string]map[bool]map[bool]float64) map[string]map[string]map[string]map[string]map[string]float64 {
	// in this function i insert the type_of_vpq and remove column is before_date_start and is_credit

	// the sequanse of the columns is:account1,account2,name,vpq,type_of_vpq,number
	new_statement := map[string]map[string]map[string]map[string]map[string]float64{}

	for k1, v1 := range old_statement { //account1
		account_struct1, _, _ := ACCOUNT_STRUCT_FROM_NAME(k1)
		for k2, v2 := range v1 { //account2
			for k3, v3 := range v2 { //name
				for k4, v4 := range v3 { //vpq
					for k5, v5 := range v4 { //is_before_date_start
						for k6, v6 := range v5 { //is_credit
							if k5 {
								// here i insert beginning_balance
								if account_struct1.IS_CREDIT == k6 {
									m := INITIALIZE_MAP_5(new_statement, k1, k2, k3, k4)
									m[beginning_balance] += v6
								} else {
									m := INITIALIZE_MAP_5(new_statement, k1, k2, k3, k4)
									m[beginning_balance] -= v6
								}
							} else {
								// here i insert inflow and outflow
								if account_struct1.IS_CREDIT == k6 {
									m := INITIALIZE_MAP_5(new_statement, k1, k2, k3, k4)
									m[inflow] += v6
								} else {
									m := INITIALIZE_MAP_5(new_statement, k1, k2, k3, k4)
									m[outflow] += v6
								}
							}
						}
					}
				}
			}
		}
	}

	return new_statement
}

func STATEMENT_STEP_7(days uint, old_statement map[string]map[string]map[string]map[string]map[string]float64) {
	// in this function we make vertical analysis of the statement

	// the sequanse of the columns is:account1,account2,name,vpq,type_of_vpq,number
	for _, v1 := range old_statement { //account1
		for _, v2 := range v1 { //account2
			for _, v3 := range v2 { //name
				for _, v4 := range v3 { //vpq
					v4[flow] = v4[inflow] - v4[outflow]
					// here i calculate the ending balance of the account by sum the flow with beginning_balance
					// because the inflow is the same of the increase of the account
					// and outflow is the same of the decrease of the account
					v4[ending_balance] = v4[beginning_balance] + v4[flow]
					v4[average] = (v4[ending_balance] + v4[beginning_balance]) / 2
					v4[turnover] = v4[outflow] / v4[average]
					v4[turnover_days] = float64(days) / v4[turnover]
					v4[growth_ratio] = v4[ending_balance] / v4[beginning_balance]
				}
			}
		}
	}
}

func STATEMENT_STEP_8(old_statement map[string]map[string]map[string]map[string]map[string]float64) {
	// in this function i complete vertical analysis of the statement
	// but here i calculate the percentage of the account from account father

	// the sequanse of the columns is:account1,account2,name,vpq,type_of_vpq,number
	for _, v1 := range old_statement { //account1
		for _, v2 := range v1 { //account2
			for _, v3 := range v2 { //name
				for k4, v4 := range v3 { //vpq
					v4[name_percent] = v4[ending_balance] / v2[all_names][k4][ending_balance]
				}
			}
		}
	}
}

func HORIZONTAL_ANALYSIS(old_statement, base_statement map[string]map[string]map[string]map[string]map[string]float64) {
	// the sequanse of the columns is:account1,account2,name,vpq,type_of_vpq,number
	for k1, v1 := range old_statement { //account1
		for k2, v2 := range v1 { //account2
			for k3, v3 := range v2 { //name
				for k4, v4 := range v3 { //vpq
					v4[change_since_base_period] = v4[ending_balance] - base_statement[k1][k2][k3][k4][ending_balance]
					v4[growth_ratio_to_base_period] = v4[ending_balance] / base_statement[k1][k2][k3][k4][ending_balance]
				}
			}
		}
	}
}

func CALCULATE_PRICE(old_statement map[string]map[string]map[string]map[string]map[string]float64) {
	// in this function we calculate the price

	// the sequanse of the columns is:account1,account2,name,vpq,type_of_vpq,number
	for _, v1 := range old_statement { //account1
		for _, v2 := range v1 { //account2
			for _, v3 := range v2 { //name
				if v3[PRICE] == nil {
					v3[PRICE] = map[string]float64{}
				}
				for _, v4 := range v3 { //vpq
					for k5 := range v4 { //type_of_vpq
						v3[PRICE][k5] = v3[VALUE][k5] / v3[QUANTITY][k5]
					}
				}
			}
		}
	}
}

func fill_new_statement(new_statement map[string]map[string]map[string]map[string]map[bool]map[bool]float64, v1 JOURNAL_TAG, is_before_date_start bool) {
	m := INITIALIZE_MAP_6(new_statement, v1.ACCOUNT_CREDIT, v1.ACCOUNT_DEBIT, v1.NAME, VALUE, is_before_date_start)
	m[true] += v1.VALUE
	m = INITIALIZE_MAP_6(new_statement, v1.ACCOUNT_CREDIT, v1.ACCOUNT_DEBIT, v1.NAME, QUANTITY, is_before_date_start)
	m[true] += ABS(v1.QUANTITY_CREDIT)
	m = INITIALIZE_MAP_6(new_statement, v1.ACCOUNT_DEBIT, v1.ACCOUNT_CREDIT, v1.NAME, VALUE, is_before_date_start)
	m[false] += v1.VALUE
	m = INITIALIZE_MAP_6(new_statement, v1.ACCOUNT_DEBIT, v1.ACCOUNT_CREDIT, v1.NAME, QUANTITY, is_before_date_start)
	m[false] += ABS(v1.QUANTITY_DEBIT)
}

func STATEMENT_FILTER(a THE_FILTER_OF_THE_STATEMENT) []FILTERED_STATEMENT {
	var new_statement []FILTERED_STATEMENT

	// the sequanse of the columns is:account1,account2,name,vpq,type_of_vpq,number
	for k1, v1 := range a.old_statement { //account1
		account_struct1, _, _ := ACCOUNT_STRUCT_FROM_NAME(k1)
		if IS_IN(k1, a.account1) == a.is_in_account1 && IS_IN(account_struct1.ACCOUNT_LEVELS[INDEX_OF_ACCOUNT_NUMBER], a.account1_levels) == a.is_in_account1_levels {
			for k2, v2 := range v1 { //account2
				account_struct2, _, _ := ACCOUNT_STRUCT_FROM_NAME(k2)
				fmt.Println(account_struct2)
				if IS_IN(k2, a.account2) == a.is_in_account2 && IS_IN(account_struct2.ACCOUNT_LEVELS[INDEX_OF_ACCOUNT_NUMBER], a.account2_levels) == a.is_in_account2_levels {
					for k3, v3 := range v2 { //name
						if IS_IN(k3, a.name) == a.is_in_name {
							for k4, v4 := range v3 { //vpq
								if IS_IN(k4, a.vpq) == a.is_in_vpq {
									for k5, v5 := range v4 { //type_of_vpq
										if IS_IN(k5, a.type_of_vpq) == a.is_in_type_of_vpq {
											new_statement = append(new_statement, FILTERED_STATEMENT{k1, k2, k3, k4, k5, v5})
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

	return new_statement
}
