package anti_accountants

import (
	"math"
	"time"
)

func ending_balance(statement map[string]map[string]map[string]map[string]map[string]float64, key_account_flow, key_account, key_name, key_vpq string) float64 {
	return statement[key_account_flow][key_account][key_name][key_vpq]["beginning_balance"] + statement[key_account][key_account][key_name][key_vpq]["increase"] - statement[key_account][key_account][key_name][key_vpq]["decrease"]
}

func sum_flows(b JOURNAL_TAG, both_credit_or_debit bool, map_v, map_q map[string]float64) {
	if (b.VALUE > 0) == both_credit_or_debit {
		map_v["outflow"] += math.Abs(b.VALUE)
		map_q["outflow"] += math.Abs(b.QUANTITY)
	} else {
		map_v["inflow"] += math.Abs(b.VALUE)
		map_q["inflow"] += math.Abs(b.QUANTITY)
	}
}

func (s FINANCIAL_ACCOUNTING) sum_values(date, start_date time.Time, entry JOURNAL_TAG, nan_flow_statement map[string]map[string]map[string]map[string]float64) {
	map_v1 := initialize_map_3(nan_flow_statement, entry.ACCOUNT, entry.NAME, "value")
	map_q1 := initialize_map_3(nan_flow_statement, entry.ACCOUNT, entry.NAME, "quantity")
	map_v2 := initialize_map_3(nan_flow_statement, s.RETAINED_EARNINGS, entry.NAME, "value")
	map_q2 := initialize_map_3(nan_flow_statement, s.RETAINED_EARNINGS, entry.NAME, "quantity")
	if date.Before(start_date) {
		switch {
		case s.is_it_sub_account_using_name(s.RETAINED_EARNINGS, entry.ACCOUNT) && s.is_credit(entry.ACCOUNT):
			map_v2["beginning_balance"] += entry.VALUE
			map_q2["beginning_balance"] += entry.QUANTITY
		case s.is_it_sub_account_using_name(s.RETAINED_EARNINGS, entry.ACCOUNT) && !s.is_credit(entry.ACCOUNT):
			map_v2["beginning_balance"] -= entry.VALUE
			map_q2["beginning_balance"] -= entry.QUANTITY
		default:
			map_v1["beginning_balance"] += entry.VALUE
			map_q1["beginning_balance"] += entry.QUANTITY
		}
	}
	if date.After(start_date) {
		if entry.VALUE >= 0 {
			map_v1["increase"] += math.Abs(entry.VALUE)
			map_q1["increase"] += math.Abs(entry.QUANTITY)
		} else {
			map_v1["decrease"] += math.Abs(entry.VALUE)
			map_q1["decrease"] += math.Abs(entry.QUANTITY)
		}
	}
}

func (s FINANCIAL_ACCOUNTING) sum_flow(date, start_date time.Time, one_simple_entry []JOURNAL_TAG, flow_statement map[string]map[string]map[string]map[string]map[string]float64) {
	for _, a := range one_simple_entry {
		for _, b := range one_simple_entry {
			map_v := initialize_map_4(flow_statement, a.ACCOUNT, b.ACCOUNT, b.NAME, "value")
			map_q := initialize_map_4(flow_statement, a.ACCOUNT, b.ACCOUNT, b.NAME, "quantity")
			if date.After(start_date) {
				both_credit_or_debit := s.is_credit(b.ACCOUNT) == s.is_credit(a.ACCOUNT)
				if b.ACCOUNT == a.ACCOUNT {
					sum_flows(b, false, map_v, map_q)
				} else {
					sum_flows(b, both_credit_or_debit, map_v, map_q)
				}
			}
		}
	}
}

func calculate_price(statement map[string]map[string]map[string]map[string]map[string]float64) {
	for _, map_account_flow := range statement {
		for _, map_account := range map_account_flow {
			for _, map_name := range map_account {
				if map_name["price"] == nil {
					map_name["price"] = map[string]float64{}
				}
				for _, map_vpq := range map_name {
					for key_number := range map_vpq {
						map_name["price"][key_number] = map_name["value"][key_number] / map_name["quantity"][key_number]
					}
				}
			}
		}
	}
}

func (s FINANCIAL_ACCOUNTING) prepare_statement(statement map[string]map[string]map[string]map[string]map[string]float64) {
	for key_account_flow, map_account_flow := range statement {
		if key_account_flow == s.CASH_AND_CASH_EQUIVALENTS {
			for key_account, map_account := range map_account_flow {
				for key_name, map_name := range map_account {
					for key_vpq, map_vpq := range map_name {
						map_vpq1 := initialize_map_4(statement, "financial_statement", key_account, key_name, key_vpq)
						for key_number, number := range map_vpq {
							map_vpq1[key_number] = number
							if !s.is_it_sub_account_using_name(s.INCOME_STATEMENT, key_account) {
								map_vpq1["percent"] = statement[s.INCOME_STATEMENT][key_account][key_name][key_vpq]["percent"]
							} else {
								map_vpq1["percent"] = statement[s.ASSETS][key_account][key_name][key_vpq]["percent"]
							}
							switch {
							case s.is_it_sub_account_using_name(s.INVENTORY, key_account):
								map_vpq1["turnover"] = statement[s.COST_OF_GOODS_SOLD][key_account][key_name][key_vpq]["turnover"]
								map_vpq1["turnover_days"] = statement[s.COST_OF_GOODS_SOLD][key_account][key_name][key_vpq]["turnover_days"]
							case s.is_it_sub_account_using_name(s.ASSETS, key_account):
								map_vpq1["turnover"] = statement[s.SALES][key_account][key_name][key_vpq]["turnover"]
								map_vpq1["turnover_days"] = statement[s.SALES][key_account][key_name][key_vpq]["turnover_days"]
							}
						}
					}
				}
			}
		}
	}
}

func horizontal_analysis(statement_current, statement_base map[string]map[string]map[string]map[string]map[string]float64) {
	for key_account_flow, map_account_flow := range statement_current {
		for key_account, map_account := range map_account_flow {
			for key_name, map_name := range map_account {
				for key_vpq, map_vpq := range map_name {
					map_vpq["change_since_base_period"] = map_vpq["ending_balance"] - statement_base[key_account_flow][key_account][key_name][key_vpq]["ending_balance"]
					map_vpq["growth_ratio_to_base_period"] = map_vpq["ending_balance"] / statement_base[key_account_flow][key_account][key_name][key_vpq]["ending_balance"]
				}
			}
		}
	}
}

func vertical_analysis(statement map[string]map[string]map[string]map[string]map[string]float64, days float64) {
	for key_account_flow, map_account_flow := range statement {
		for key_account, map_account := range map_account_flow {
			for key_name, map_name := range map_account {
				for key_vpq, map_vpq := range map_name {
					map_vpq["increase_or_decrease"] = map_vpq["increase"] - map_vpq["decrease"]
					map_vpq["ending_balance"] = map_vpq["beginning_balance"] + map_vpq["increase_or_decrease"]
					map_vpq["flow"] = map_vpq["inflow"] - map_vpq["outflow"]
					map_vpq["average"] = (map_vpq["ending_balance"] + map_vpq["beginning_balance"]) / 2
					map_vpq["turnover"] = map_vpq["inflow"] / map_vpq["average"]
					map_vpq["turnover_days"] = days / map_vpq["turnover"]
					map_vpq["growth_ratio"] = map_vpq["ending_balance"] / map_vpq["beginning_balance"]
					map_vpq["percent"] = map_vpq["ending_balance"] / ending_balance(statement, key_account_flow, key_account_flow, key_name, key_vpq)
					map_vpq["name_percent"] = map_vpq["ending_balance"] / ending_balance(statement, key_account_flow, key_account, "all", key_vpq)
				}
			}
		}
	}
}

func sum_3rd_column(statement map[string]map[string]map[string]map[string]map[string]float64, names, exempt_names []string, name string, in_names bool) {
	for _, map_account_flow := range statement {
		for _, map_account := range map_account_flow {
			if map_account[name] == nil {
				map_account[name] = map[string]map[string]float64{}
			}
			for key_name, map_name := range map_account {
				var ok bool
				if !IS_IN(key_name, append(exempt_names, name)) {
					if IS_IN(key_name, names) == in_names {
						ok = true
					}
					if ok {
						for key_vpq, map_vpq := range map_name {
							if map_account[name][key_vpq] == nil {
								map_account[name][key_vpq] = map[string]float64{}
							}
							for key_number, number := range map_vpq {
								map_account[name][key_vpq][key_number] += number
							}
						}
					}
				}
			}
		}
	}
}

func (s FINANCIAL_ACCOUNTING) sum_2nd_column(statement map[string]map[string]map[string]map[string]map[string]float64) map[string]map[string]map[string]map[string]map[string]float64 {
	new_statement := map[string]map[string]map[string]map[string]map[string]float64{}
	for key_account_flow, map_account_flow := range statement {
		for key_account, map_account := range map_account_flow {
			higher_level_accounts := append(s.find_all_higher_level_accounts(key_account), key_account)
			for key_name, map_name := range map_account {
				for key_vpq, map_vpq := range map_name {
					for _, accuont := range higher_level_accounts {
						new_map_vpq := initialize_map_4(new_statement, key_account_flow, accuont, key_name, key_vpq)
						for key_number, number := range map_vpq {
							switch {
							case !IS_IN(key_number, []string{"inflow", "outflow"}):
								if s.is_credit(key_account) == s.is_credit(accuont) {
									new_map_vpq[key_number] += number
								} else {
									new_map_vpq[key_number] -= number
								}
							case key_account_flow != accuont:
								new_map_vpq[key_number] += number
							case key_account_flow == accuont:
								new_statement[key_account_flow][accuont][key_name][key_vpq][key_number] += number
							}
						}
					}
				}
			}
		}
	}
	// for key_account_flow, map_account_flow := range statement {
	// 	for key_account, map_account := range map_account_flow {
	// 		var last_name string
	// 		key1 := key_account
	// 		for {
	// 			for _, ss := range s.ACCOUNTS {
	// 				if ss.NAME == key_account {
	// 					key_account = ss.FATHER
	// 					for key_name, map_name := range map_account {
	// 						for key_vpq, map_vpq := range map_name {
	// 							map_vpq1 := initialize_map_4(new_statement, key_account_flow, ss.NAME, key_name, key_vpq)
	// 							for key_number, number := range map_vpq {
	// 								switch {
	// 								case !IS_IN(key_number, []string{"inflow", "outflow"}):
	// 									if s.is_credit(key1) == s.is_credit(ss.NAME) {
	// 										map_vpq1[key_number] += number
	// 									} else {
	// 										map_vpq1[key_number] -= number
	// 									}
	// 								case key_account_flow != key1:
	// 									map_vpq1[key_number] += number
	// 								case key_account_flow == ss.NAME:
	// 									new_statement[key_account_flow][key1][key_name][key_vpq][key_number] += number
	// 								}
	// 							}
	// 						}
	// 					}
	// 				}
	// 			}
	// 			if last_name == key_account {
	// 				break
	// 			}
	// 			last_name = key_account
	// 		}
	// 	}
	// }
	return new_statement
}

func (s FINANCIAL_ACCOUNTING) sum_1st_column(statement map[string]map[string]map[string]map[string]map[string]float64) map[string]map[string]map[string]map[string]map[string]float64 {
	new_statement := map[string]map[string]map[string]map[string]map[string]float64{}
	for key_account_flow, map_account_flow := range statement {
		higher_level_accounts := append(s.find_all_higher_level_accounts(key_account_flow), key_account_flow)
		for key_account, map_account := range map_account_flow {
			for key_name, map_name := range map_account {
				for key_vpq, map_vpq := range map_name {
					for _, accuont := range higher_level_accounts {
						new_map_vpq := initialize_map_4(new_statement, accuont, key_account, key_name, key_vpq)
						for key_number, number := range map_vpq {
							switch {
							case !IS_IN(key_number, []string{"inflow", "outflow"}):
								// if s.is_credit(key_account) == s.is_credit(accuont) {
								new_map_vpq[key_number] += number
								// } else {
								// 	new_map_vpq[key_number] -= number
								// }
							default:
								new_map_vpq[key_number] = number
							}
						}
					}
				}
			}
		}
	}
	// var flow_accounts []string
	// for _, a := range s.ACCOUNTS {
	// 	for _, b := range s.ACCOUNTS {
	// 		if s.is_it_sub_account_using_name(a.NAME, b.NAME) {
	// 			flow_accounts = append(flow_accounts, b.NAME)
	// 		}
	// 	}
	// 	for key_account_flow, map_account_flow := range statement {
	// 		if IS_IN(key_account_flow, flow_accounts) {
	// 			for key_account, map_account := range map_account_flow {
	// 				for key_name, map_name := range map_account {
	// 					for key_vpq, map_vpq := range map_name {
	// 						map_vpq1 := initialize_map_4(new_statement, a.NAME, key_account, key_name, key_vpq)
	// 						for key_number, number := range map_vpq {
	// 							switch {
	// 							case IS_IN(key_number, []string{"inflow", "outflow"}):
	// 								if s.is_credit(a.NAME) == s.is_credit(key_account_flow) {
	// 									map_vpq1[key_number] += number
	// 								} else {
	// 									map_vpq1[key_number] -= number
	// 								}
	// 							default:
	// 								map_vpq1[key_number] = number
	// 							}
	// 						}
	// 					}
	// 				}
	// 			}
	// 		}
	// 	}
	// 	flow_accounts = []string{}
	// }
	return new_statement
}

func combine_statements(flow_statement map[string]map[string]map[string]map[string]map[string]float64, nan_flow_statement map[string]map[string]map[string]map[string]float64) map[string]map[string]map[string]map[string]map[string]float64 {
	for key_account_flow := range nan_flow_statement {
		for key_account, map_account := range nan_flow_statement {
			for key_name, map_name := range map_account {
				for key_vpq, map_vpq := range map_name {
					map_vpq1 := initialize_map_4(flow_statement, key_account_flow, key_account, key_name, key_vpq)
					for key_number := range map_vpq {
						map_vpq1[key_number] = map_vpq[key_number]
					}
				}
			}
		}
	}
	return flow_statement
}

func (s FINANCIAL_ACCOUNTING) statement(journal []JOURNAL_TAG, start_date, end_date time.Time) (map[string]map[string]map[string]map[string]map[string]float64, map[string]map[string]map[string]map[string]float64) {
	var one_simple_entry []JOURNAL_TAG
	var previous_entry_number int
	var date time.Time
	flow_statement := map[string]map[string]map[string]map[string]map[string]float64{}
	nan_flow_statement := map[string]map[string]map[string]map[string]float64{}
	for _, entry := range journal {
		date = PARSE_DATE(entry.DATE, s.DATE_LAYOUT)
		if previous_entry_number != entry.ENTRY_NUMBER {
			s.sum_flow(date, start_date, one_simple_entry, flow_statement)
			one_simple_entry = []JOURNAL_TAG{}
		}
		if date.Before(end_date) {
			s.sum_values(date, start_date, entry, nan_flow_statement)
			one_simple_entry = append(one_simple_entry, entry)
		}
		previous_entry_number = entry.ENTRY_NUMBER
	}
	s.sum_flow(date, start_date, one_simple_entry, flow_statement)
	return flow_statement, nan_flow_statement
}

func (s FINANCIAL_ACCOUNTING) FINANCIAL_STATEMENTS(JOURNAL_ORDERED_BY_DATE_ENTRY_NUMBER []JOURNAL_TAG, start_date, end_date time.Time, periods int, names []string, in_names bool) []map[string]map[string]map[string]map[string]map[string]float64 {
	check_dates(start_date, end_date)
	days := int(end_date.Sub(start_date).Hours() / 24)
	statements := []map[string]map[string]map[string]map[string]map[string]float64{}
	for a := 0; a < periods; a++ {
		flow_statement, nan_flow_statement := s.statement(JOURNAL_ORDERED_BY_DATE_ENTRY_NUMBER, start_date.AddDate(0, 0, -days*a), end_date.AddDate(0, 0, -days*a))
		statement := combine_statements(flow_statement, nan_flow_statement)
		statement = s.sum_1st_column(statement)
		statement = s.sum_2nd_column(statement)
		sum_3rd_column(statement, []string{}, []string{}, "all", false)
		sum_3rd_column(statement, names, []string{"all"}, "names", in_names)
		vertical_analysis(statement, float64(days))
		statements = append(statements, statement)
	}
	for _, statement_current := range statements {
		horizontal_analysis(statement_current, statements[periods-1])
		s.prepare_statement(statement_current)
		calculate_price(statement_current)
	}
	return statements
}
