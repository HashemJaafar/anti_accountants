package anti_accountants

import (
	"fmt"
	"math"
	"time"
)

func FINANCIAL_STATEMENTS(date_start, date_end time.Time, periods int, names_you_want []string, in_names bool) []map[string]map[string]map[string]map[string]map[string]float64 {

	if date_start.After(date_end) {
		date_start = date_end
	}

	keys, journal := DB_READ[JOURNAL_TAG](DB_JOURNAL)
	journal_times := CONVERT_BYTE_SLICE_TO_TIME(keys)

	days := int(date_end.Sub(date_start).Hours() / 24)
	var statements []map[string]map[string]map[string]map[string]map[string]float64

	for k1 := 0; k1 < periods; k1++ {
		trailing_balance_sheet := STATEMENT_STEP_1(journal_times, journal, date_start.AddDate(0, 0, -days*k1), date_end.AddDate(0, 0, -days*k1))
		trailing_balance_sheet = STATEMENT_STEP_2(trailing_balance_sheet)
		trailing_balance_sheet = STATEMENT_STEP_3(trailing_balance_sheet)
		trailing_balance_sheet = STATEMENT_STEP_4(in_names, names_you_want, trailing_balance_sheet)
		statement := STATEMENT_STEP_5(trailing_balance_sheet)
		VERTICAL_ANALYSIS(days, statement)
		// 		statement = sum_1st_column(statement)
		// 		statement = sum_2nd_column(statement)
		// 		sum_3rd_column(statement, []string{}, []string{}, "all", false)
		// 		sum_3rd_column(statement, names, []string{"all"}, "names", in_names)
		// 		vertical_analysis(statement, float64(days))
		// 		statements = append(statements, statement)
		// 	}
		// 	for _, statement_current := range statements {
		// 		horizontal_analysis(statement_current, statements[periods-1])
		// 		prepare_statement(statement_current)
		// 		calculate_price(statement_current)
	}
	return statements
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

func STATEMENT_STEP_2(old_statement map[string]map[string]map[string]map[string]map[bool]map[bool]float64) map[string]map[string]map[string]map[string]map[bool]map[bool]float64 {
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

func STATEMENT_STEP_3(old_statement map[string]map[string]map[string]map[string]map[bool]map[bool]float64) map[string]map[string]map[string]map[string]map[bool]map[bool]float64 {
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

func STATEMENT_STEP_4(in_names bool, names_you_want []string, old_statement map[string]map[string]map[string]map[string]map[bool]map[bool]float64) map[string]map[string]map[string]map[string]map[bool]map[bool]float64 {
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

func STATEMENT_STEP_5(old_statement map[string]map[string]map[string]map[string]map[bool]map[bool]float64) map[string]map[string]map[string]map[string]map[string]float64 {
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

func VERTICAL_ANALYSIS(days int, old_statement map[string]map[string]map[string]map[string]map[string]float64) {
	// in this function we make vertical analysis of the statement

	// the sequanse of the columns is:account1,account2,name,vpq,type_of_vpq,number
	for _, v1 := range old_statement { //account1
		for _, v2 := range v1 { //account2
			for k3, v3 := range v2 { //name
				for k4, v4 := range v3 { //vpq
					v4[flow] = v4[inflow] - v4[outflow]
					// here i calculate the ending balance of the account by the total flow
					// because the total inflow is the same of the total increase of the account
					// and total outflow is the same of the total decrease of the account
					x := v1[all_accounts][k3][k4]
					v4[ending_balance] = v4[beginning_balance] + x[inflow] - x[outflow]
					v4[average] = (v4[ending_balance] + v4[beginning_balance]) / 2
					v4[turnover] = v4[outflow] / v4[average]
					v4[turnover_days] = float64(days) / v4[turnover]
					v4[growth_ratio] = v4[ending_balance] / v4[beginning_balance]
					// v4[percent] = v4[ending_balance] / ending_balance(statement, key_account_flow, key_account_flow, key_name, key_vpq)
					// v4[name_percent] = v4[ending_balance] / ending_balance(statement, key_account_flow, key_account, "all", key_vpq)
				}
			}
		}
	}
}

func fill_new_statement(new_statement map[string]map[string]map[string]map[string]map[bool]map[bool]float64, v1 JOURNAL_TAG, is_before_date_start bool) {
	m := INITIALIZE_MAP_6(new_statement, v1.ACCOUNT_CREDIT, v1.ACCOUNT_DEBIT, v1.NAME, VALUE, is_before_date_start)
	m[true] += v1.VALUE
	m = INITIALIZE_MAP_6(new_statement, v1.ACCOUNT_CREDIT, v1.ACCOUNT_DEBIT, v1.NAME, QUANTITY, is_before_date_start)
	m[true] += math.Abs(v1.QUANTITY_CREDIT)
	m = INITIALIZE_MAP_6(new_statement, v1.ACCOUNT_DEBIT, v1.ACCOUNT_CREDIT, v1.NAME, VALUE, is_before_date_start)
	m[false] += v1.VALUE
	m = INITIALIZE_MAP_6(new_statement, v1.ACCOUNT_DEBIT, v1.ACCOUNT_CREDIT, v1.NAME, QUANTITY, is_before_date_start)
	m[false] += math.Abs(v1.QUANTITY_DEBIT)
}
func print_map_6(columns map[string]map[string]map[string]map[string]map[bool]map[bool]float64) {
	for k1, v1 := range columns {
		for k2, v2 := range v1 {
			for k3, v3 := range v2 {
				for k4, v4 := range v3 {
					for k5, v5 := range v4 {
						for k6, v6 := range v5 {
							fmt.Fprintln(PRINT_TABLE, k1, "\t", k2, "\t", k3, "\t", k4, "\t", k5, "\t", k6, "\t", v6)
						}
					}
				}
			}
		}
	}
	fmt.Println("//////////////////////////////////////////")
	PRINT_TABLE.Flush()
}
func print_map_5[ta, tb, tc, td, te comparable, tr any](m map[ta]map[tb]map[tc]map[td]map[te]tr) {
	for k1, v1 := range m {
		for k2, v2 := range v1 {
			for k3, v3 := range v2 {
				for k4, v4 := range v3 {
					for k5, v5 := range v4 {
						fmt.Fprintln(PRINT_TABLE, k1, "\t", k2, "\t", k3, "\t", k4, "\t", k5, "\t", v5)
					}
				}
			}
		}
	}
	fmt.Println("//////////////////////////////////////////")
	PRINT_TABLE.Flush()
}
func print_map_4[ta, tb, tc, td comparable, tr any](m map[ta]map[tb]map[tc]map[td]tr) {
	for k1, v1 := range m {
		for k2, v2 := range v1 {
			for k3, v3 := range v2 {
				for k4, v4 := range v3 {
					fmt.Fprintln(PRINT_TABLE, k1, "\t", k2, "\t", k3, "\t", k4, "\t", v4)
				}
			}
		}
	}
	fmt.Println("//////////////////////////////////////////")
	PRINT_TABLE.Flush()
}
