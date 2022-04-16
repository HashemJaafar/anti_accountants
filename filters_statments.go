
package anti_accountants
// import (
// 	"sort"
// 	"strings"
// )
// func STATEMENT_FILTER(
// 	all_financial_statements []map[string]map[string]map[string]map[string]map[string]float64,
// 	account_flow_slice, account_slice, name_slice, vpq_slice, number_slice []string,
// 	account_flow_levels, account_levels []uint,
// 	in_account_flow_slice, in_account_slice, in_name_slice, in_vpq_slice, in_number_slice,
// 	in_account_flow_levels, in_account_levels bool) [][]FILTERED_STATEMENT {
// 	var all_statements_struct [][]FILTERED_STATEMENT
// 	for _, statement := range all_financial_statements {
// 		var one_statement_struct []FILTERED_STATEMENT
// 		for key_account_flow, map_account_flow := range statement {
// 			if IS_IN(key_account_flow, account_flow_slice) == in_account_flow_slice {
// 				if accept_level_using_name(key_account_flow, account_flow_levels, in_account_flow_levels) {
// 					for key_account, map_account := range map_account_flow {
// 						if IS_IN(key_account, account_slice) == in_account_slice {
// 							if accept_level_using_name(key_account, account_levels, in_account_levels) {
// 								for key_name, map_name := range map_account {
// 									if IS_IN(key_name, name_slice) == in_name_slice {
// 										for key_vpq, map_vpq := range map_name {
// 											if IS_IN(key_vpq, vpq_slice) == in_vpq_slice {
// 												for key_number, number := range map_vpq {
// 													if IS_IN(key_number, number_slice) == in_number_slice {
// 														one_statement_struct = append(one_statement_struct, FILTERED_STATEMENT{key_account_flow, key_account, key_name, key_vpq, key_number, number})
// 													}
// 												}
// 											}
// 										}
// 									}
// 								}
// 							}
// 						}
// 					}
// 				}
// 			}
// 		}
// 		all_statements_struct = append(all_statements_struct, one_statement_struct)
// 	}
// 	return all_statements_struct
// }
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
// 		default:
// 			error_element_is_not_in_elements(sort_by, []string{"account_number", "number"})
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
