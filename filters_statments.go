package main

// import (
// 	"sort"
// 	"strings"
// )
// func SORT_THE_STATMENT(all_statements_struct [][]FILTERED_STATEMENT, sort_by string, is_reverse bool) {
// 	for _, one_statement_struct := range all_statements_struct {
// 		switch sort_by {
// 		case "account_number":
// 			column1_sort_statement_by_account_number(one_statement_struct)
// 			column2_sort_statement_by_account_number(one_statement_struct)
// 			column3_sort_statement_by_alphabet(one_statement_struct)
// 			column4_sort_statement_by_alphabet(one_statement_struct)
// 			column5_sort_statement_by_alphabet(one_statement_struct)
// 		case "number":
// 			sort_by_number(one_statement_struct)
// 		}
// 		if is_reverse {
// 			REVERSE_SLICE(one_statement_struct)
// 		}
// 		make_space_before_account_in_statement_struct(one_statement_struct)
// 	}
// }
// func column1_sort_statement_by_account_number(one_statement_struct []FILTERED_STATEMENT) {
// 	for indexa := range one_statement_struct {
// 		for indexb := range one_statement_struct {
// 			if indexa < indexb {
// 				if one_statement_struct[indexa].KEY_ACCOUNT_FLOW == "financial_statement" && one_statement_struct[indexb].KEY_ACCOUNT_FLOW != "financial_statement" {
// 					compare_the_numbers_and_swap(one_statement_struct, indexa, indexb, []uint{}, account_number(one_statement_struct[indexb].KEY_ACCOUNT_FLOW))
// 				} else if one_statement_struct[indexa].KEY_ACCOUNT_FLOW != "financial_statement" && one_statement_struct[indexb].KEY_ACCOUNT_FLOW == "financial_statement" {
// 					compare_the_numbers_and_swap(one_statement_struct, indexa, indexb, account_number(one_statement_struct[indexa].KEY_ACCOUNT_FLOW), []uint{})
// 				} else if one_statement_struct[indexa].KEY_ACCOUNT_FLOW == "financial_statement" && one_statement_struct[indexb].KEY_ACCOUNT_FLOW == "financial_statement" {
// 					compare_the_numbers_and_swap(one_statement_struct, indexa, indexb, []uint{}, []uint{})
// 				} else {
// 					compare_the_numbers_and_swap(one_statement_struct, indexa, indexb, account_number(one_statement_struct[indexa].KEY_ACCOUNT_FLOW), account_number(one_statement_struct[indexb].KEY_ACCOUNT_FLOW))
// 				}
// 			}
// 		}
// 	}
// }
// func column2_sort_statement_by_account_number(one_statement_struct []FILTERED_STATEMENT) {
// 	for indexa := range one_statement_struct {
// 		for indexb := range one_statement_struct {
// 			if indexa < indexb {
// 				if one_statement_struct[indexa].KEY_ACCOUNT_FLOW == one_statement_struct[indexb].KEY_ACCOUNT_FLOW {
// 					compare_the_numbers_and_swap(one_statement_struct, indexa, indexb, account_number(one_statement_struct[indexa].KEY_ACCOUNT), account_number(one_statement_struct[indexb].KEY_ACCOUNT))
// 				}
// 			}
// 		}
// 	}
// }
// func column3_sort_statement_by_alphabet(one_statement_struct []FILTERED_STATEMENT) {
// 	for indexa := range one_statement_struct {
// 		for indexb := range one_statement_struct {
// 			if indexa < indexb {
// 				if one_statement_struct[indexa].KEY_ACCOUNT_FLOW == one_statement_struct[indexb].KEY_ACCOUNT_FLOW {
// 					if one_statement_struct[indexa].KEY_ACCOUNT == one_statement_struct[indexb].KEY_ACCOUNT {
// 						if one_statement_struct[indexa].KEY_NAME > one_statement_struct[indexb].KEY_NAME {
// 							one_statement_struct[indexa], one_statement_struct[indexb] = one_statement_struct[indexb], one_statement_struct[indexa]
// 						}
// 					}
// 				}
// 			}
// 		}
// 	}
// }
// func column4_sort_statement_by_alphabet(one_statement_struct []FILTERED_STATEMENT) {
// 	for indexa := range one_statement_struct {
// 		for indexb := range one_statement_struct {
// 			if indexa < indexb {
// 				if one_statement_struct[indexa].KEY_ACCOUNT_FLOW == one_statement_struct[indexb].KEY_ACCOUNT_FLOW {
// 					if one_statement_struct[indexa].KEY_ACCOUNT == one_statement_struct[indexb].KEY_ACCOUNT {
// 						if one_statement_struct[indexa].KEY_NAME == one_statement_struct[indexb].KEY_NAME {
// 							if one_statement_struct[indexa].KEY_VPQ > one_statement_struct[indexb].KEY_VPQ {
// 								one_statement_struct[indexa], one_statement_struct[indexb] = one_statement_struct[indexb], one_statement_struct[indexa]
// 							}
// 						}
// 					}
// 				}
// 			}
// 		}
// 	}
// }
// func column5_sort_statement_by_alphabet(one_statement_struct []FILTERED_STATEMENT) {
// 	for indexa := range one_statement_struct {
// 		for indexb := range one_statement_struct {
// 			if indexa < indexb {
// 				if one_statement_struct[indexa].KEY_ACCOUNT_FLOW == one_statement_struct[indexb].KEY_ACCOUNT_FLOW {
// 					if one_statement_struct[indexa].KEY_ACCOUNT == one_statement_struct[indexb].KEY_ACCOUNT {
// 						if one_statement_struct[indexa].KEY_NAME == one_statement_struct[indexb].KEY_NAME {
// 							if one_statement_struct[indexa].KEY_VPQ == one_statement_struct[indexb].KEY_VPQ {
// 								if one_statement_struct[indexa].KEY_NUMBER > one_statement_struct[indexb].KEY_NUMBER {
// 									one_statement_struct[indexa], one_statement_struct[indexb] = one_statement_struct[indexb], one_statement_struct[indexa]
// 								}
// 							}
// 						}
// 					}
// 				}
// 			}
// 		}
// 	}
// }
// func compare_the_numbers_and_swap(one_statement_struct []FILTERED_STATEMENT, indexa, indexb int, account_number1, account_number2 []uint) {
// 	if !IS_IT_HIGH_THAN_BY_ORDER(account_number1, account_number2) {
// 		one_statement_struct[indexa], one_statement_struct[indexb] = one_statement_struct[indexb], one_statement_struct[indexa]
// 	}
// }
// func sort_by_number(one_statement_struct []FILTERED_STATEMENT) {
// 	sort.Slice(one_statement_struct, func(p, q int) bool { return one_statement_struct[p].NUMBER < one_statement_struct[q].NUMBER })
// }
// func make_space_before_account_in_statement_struct(one_statement_struct []FILTERED_STATEMENT) {
// 	for indexa, a := range one_statement_struct {
// 		if a.KEY_ACCOUNT_FLOW != "financial_statement" {
// 			one_statement_struct[indexa].KEY_ACCOUNT_FLOW = make_space_before_account_name(a.KEY_ACCOUNT_FLOW)
// 		}
// 		one_statement_struct[indexa].KEY_ACCOUNT = make_space_before_account_name(a.KEY_ACCOUNT)
// 	}
// }
// func make_space_before_account_name(account_name string) string {
// 	return strings.Repeat("  ", len(account_number(account_name))) + account_name
// }
// func accept_level_using_name(account_name string, levels []uint, in_levels bool) bool {
// 	var len_account_number int
// 	if account_name != "financial_statement" {
// 		len_account_number = len(account_number(account_name))
// 	}
// 	for _, a := range levels {
// 		if (len_account_number == int(a)) == in_levels {
// 			return true
// 		}
// 	}
// 	return false
// }
// func sort_statement_by_pre_order_in_insertion_sort(one_statement_struct []FILTERED_STATEMENT) {
// 	var indexa int
// 	for _, a := range ACCOUNTS {
// 		for indexb, b := range one_statement_struct {
// 			if a.ACCOUNT_NAME == b.KEY_ACCOUNT {
// 				one_statement_struct[indexa], one_statement_struct[indexb] = one_statement_struct[indexb], one_statement_struct[indexa]
// 				indexa++
// 				break
// 			}
// 		}
// 	}
// }
