package anti_accountants

import (
	"log"
	"math"
	"sort"
	"time"
)

type FINANCIAL_ANALYSIS struct {
	CURRENT_ASSETS,
	CURRENT_LIABILITIES,
	CASH,
	SHORT_TERM_INVESTMENTS,
	NET_RECEIVABLES,
	NET_CREDIT_SALES,
	AVERAGE_NET_RECEIVABLES,
	COST_OF_GOODS_SOLD,
	AVERAGE_INVENTORY,
	NET_INCOME,
	NET_SALES,
	AVERAGE_ASSETS,
	AVERAGE_EQUITY,
	PREFERRED_DIVIDENDS,
	AVERAGE_COMMON_STOCKHOLDERS_EQUITY,
	MARKET_PRICE_PER_SHARES_OUTSTANDING,
	CASH_DIVIDENDS,
	TOTAL_DEBT,
	TOTAL_ASSETS,
	EBITDA,
	INTEREST_EXPENSE,
	WEIGHTED_AVERAGE_COMMON_SHARES_OUTSTANDING float64
}

type FINANCIAL_ANALYSIS_STATEMENT struct {
	CURRENT_RATIO                        float64 // current_assets / current_liabilities
	ACID_TEST                            float64 // (cash + short_term_investments + net_receivables) / current_liabilities
	RECEIVABLES_TURNOVER                 float64 // net_credit_sales / average_net_receivables
	INVENTORY_TURNOVER                   float64 // cost_of_goods_sold / average_inventory
	ASSET_TURNOVER                       float64 // net_sales / average_assets
	PROFIT_MARGIN                        float64 // net_income / net_sales
	RETURN_ON_ASSETS                     float64 // net_income / average_assets
	RETURN_ON_EQUITY                     float64 // net_income / average_equity
	PAYOUT_RATIO                         float64 // cash_dividends / net_income
	DEBT_TO_TOTAL_ASSETS_RATIO           float64 // total_debt / total_assets
	TIMES_INTEREST_EARNED                float64 // ebitda / interest_expense
	RETURN_ON_COMMON_STOCKHOLDERS_EQUITY float64 // (net_income - preferred_dividends) / average_common_stockholders_equity
	EARNINGS_PER_SHARE                   float64 // (net_income - preferred_dividends) / weighted_average_common_shares_outstanding
	PRICE_EARNINGS_RATIO                 float64 // market_price_per_shares_outstanding / earnings_per_share
}

type FILTERED_STATEMENT struct {
	KEY_ACCOUNT_FLOW, KEY_ACCOUNT, KEY_NAME, KEY_VPQ, KEY_NUMBER string
	NUMBER                                                       float64
}

func ending_balance(statement map[string]map[string]map[string]map[string]map[string]float64, key_account_flow, key_account, key_name, key_vpq string) float64 {
	return statement[key_account_flow][key_account][key_name][key_vpq]["beginning_balance"] + statement[key_account][key_account][key_name][key_vpq]["increase"] - statement[key_account][key_account][key_name][key_vpq]["decrease"]
}

func sum_flows(b JOURNAL_TAG, x float64, map_v, map_q map[string]float64) {
	if b.VALUE*x < 0 {
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
		case s.is_father(s.RETAINED_EARNINGS, entry.ACCOUNT) && s.is_credit(entry.ACCOUNT):
			map_v2["beginning_balance"] += entry.VALUE
			map_q2["beginning_balance"] += entry.QUANTITY
		case s.is_father(s.RETAINED_EARNINGS, entry.ACCOUNT) && !s.is_credit(entry.ACCOUNT):
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
				if b.ACCOUNT == a.ACCOUNT || s.is_credit(b.ACCOUNT) != s.is_credit(a.ACCOUNT) {
					sum_flows(b, 1, map_v, map_q)
				} else {
					sum_flows(b, -1, map_v, map_q)
				}
			}
		}
	}
}

func (s FINANCIAL_ACCOUNTING) analysis(statement map[string]map[string]map[string]map[string]map[string]float64) FINANCIAL_ANALYSIS_STATEMENT {
	return FINANCIAL_ANALYSIS{
		CURRENT_ASSETS:                      statement[s.CASH_AND_CASH_EQUIVALENTS][s.CURRENT_ASSETS]["names"]["value"]["ending_balance"],
		CURRENT_LIABILITIES:                 statement[s.CASH_AND_CASH_EQUIVALENTS][s.CURRENT_LIABILITIES]["names"]["value"]["ending_balance"],
		CASH:                                statement[s.CASH_AND_CASH_EQUIVALENTS][s.CASH_AND_CASH_EQUIVALENTS]["names"]["value"]["ending_balance"],
		SHORT_TERM_INVESTMENTS:              statement[s.CASH_AND_CASH_EQUIVALENTS][s.SHORT_TERM_INVESTMENTS]["names"]["value"]["ending_balance"],
		NET_RECEIVABLES:                     statement[s.CASH_AND_CASH_EQUIVALENTS][s.RECEIVABLES]["names"]["value"]["ending_balance"],
		NET_CREDIT_SALES:                    statement[s.SALES][s.RECEIVABLES]["names"]["value"]["flow"],
		AVERAGE_NET_RECEIVABLES:             statement[s.CASH_AND_CASH_EQUIVALENTS][s.RECEIVABLES]["names"]["value"]["average"],
		COST_OF_GOODS_SOLD:                  statement[s.CASH_AND_CASH_EQUIVALENTS][s.COST_OF_GOODS_SOLD]["names"]["value"]["ending_balance"],
		AVERAGE_INVENTORY:                   statement[s.CASH_AND_CASH_EQUIVALENTS][s.INVENTORY]["names"]["value"]["average"],
		NET_INCOME:                          statement[s.CASH_AND_CASH_EQUIVALENTS][s.INCOME_STATEMENT]["names"]["value"]["ending_balance"],
		NET_SALES:                           statement[s.CASH_AND_CASH_EQUIVALENTS][s.SALES]["names"]["value"]["ending_balance"],
		AVERAGE_ASSETS:                      statement[s.CASH_AND_CASH_EQUIVALENTS][s.ASSETS]["names"]["value"]["average"],
		AVERAGE_EQUITY:                      statement[s.CASH_AND_CASH_EQUIVALENTS][s.EQUITY]["names"]["value"]["average"],
		PREFERRED_DIVIDENDS:                 0,
		AVERAGE_COMMON_STOCKHOLDERS_EQUITY:  0,
		MARKET_PRICE_PER_SHARES_OUTSTANDING: 0,
		CASH_DIVIDENDS:                      statement[s.CASH_AND_CASH_EQUIVALENTS][s.DIVIDENDS]["names"]["value"]["flow"],
		TOTAL_DEBT:                          statement[s.CASH_AND_CASH_EQUIVALENTS][s.LIABILITIES]["names"]["value"]["ending_balance"],
		TOTAL_ASSETS:                        statement[s.CASH_AND_CASH_EQUIVALENTS][s.ASSETS]["names"]["value"]["ending_balance"],
		EBITDA:                              statement[s.CASH_AND_CASH_EQUIVALENTS][s.EBITDA]["names"]["value"]["ending_balance"],
		INTEREST_EXPENSE:                    statement[s.CASH_AND_CASH_EQUIVALENTS][s.INTEREST_EXPENSE]["names"]["value"]["ending_balance"],
		WEIGHTED_AVERAGE_COMMON_SHARES_OUTSTANDING: 0,
	}.FINANCIAL_ANALYSIS_STATEMENT()
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
							if !s.is_father(s.INCOME_STATEMENT, key_account) {
								map_vpq1["percent"] = statement[s.INCOME_STATEMENT][key_account][key_name][key_vpq]["percent"]
							} else {
								map_vpq1["percent"] = statement[s.ASSETS][key_account][key_name][key_vpq]["percent"]
							}
							switch {
							case s.is_father(s.INVENTORY, key_account):
								map_vpq1["turnover"] = statement[s.COST_OF_GOODS_SOLD][key_account][key_name][key_vpq]["turnover"]
								map_vpq1["turnover_days"] = statement[s.COST_OF_GOODS_SOLD][key_account][key_name][key_vpq]["turnover_days"]
							case s.is_father(s.ASSETS, key_account):
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
			var last_name string
			key1 := key_account
			for {
				for _, ss := range s.ACCOUNTS {
					if ss.NAME == key_account {
						key_account = ss.FATHER
						for key_name, map_name := range map_account {
							for key_vpq, map_vpq := range map_name {
								map_vpq1 := initialize_map_4(new_statement, key_account_flow, ss.NAME, key_name, key_vpq)
								for key_number, number := range map_vpq {
									switch {
									case !IS_IN(key_number, []string{"inflow", "outflow"}):
										if s.is_credit(key1) == s.is_credit(ss.NAME) {
											map_vpq1[key_number] += number
										} else {
											map_vpq1[key_number] -= number
										}
									case key_account_flow != key1:
										map_vpq1[key_number] += number
									case key_account_flow == ss.NAME:
										new_statement[key_account_flow][key1][key_name][key_vpq][key_number] += number
									}
								}
							}
						}
					}
				}
				if last_name == key_account {
					break
				}
				last_name = key_account
			}
		}
	}
	return new_statement
}

func (s FINANCIAL_ACCOUNTING) sum_1st_column(statement map[string]map[string]map[string]map[string]map[string]float64) map[string]map[string]map[string]map[string]map[string]float64 {
	new_statement := map[string]map[string]map[string]map[string]map[string]float64{}
	var flow_accounts []string
	for _, a := range s.ACCOUNTS {
		for _, b := range s.ACCOUNTS {
			if s.is_father(a.NAME, b.NAME) {
				flow_accounts = append(flow_accounts, b.NAME)
			}
		}
		for key_account_flow, map_account_flow := range statement {
			if IS_IN(key_account_flow, flow_accounts) {
				for key_account, map_account := range map_account_flow {
					for key_name, map_name := range map_account {
						for key_vpq, map_vpq := range map_name {
							map_vpq1 := initialize_map_4(new_statement, a.NAME, key_account, key_name, key_vpq)
							for key_number, number := range map_vpq {
								switch {
								case IS_IN(key_number, []string{"inflow", "outflow"}):
									if s.is_credit(a.NAME) == s.is_credit(key_account_flow) {
										map_vpq1[key_number] += number
									} else {
										map_vpq1[key_number] -= number
									}
								default:
									map_vpq1[key_number] = number
								}
							}
						}
					}
				}
			}
		}
		flow_accounts = []string{}
	}
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

func (s FINANCIAL_ANALYSIS) FINANCIAL_ANALYSIS_STATEMENT() FINANCIAL_ANALYSIS_STATEMENT {
	CURRENT_RATIO := s.CURRENT_ASSETS / s.CURRENT_LIABILITIES
	ACID_TEST := (s.CASH + s.SHORT_TERM_INVESTMENTS + s.NET_RECEIVABLES) / s.CURRENT_LIABILITIES
	RECEIVABLES_TURNOVER := s.NET_CREDIT_SALES / s.AVERAGE_NET_RECEIVABLES
	INVENTORY_TURNOVER := s.COST_OF_GOODS_SOLD / s.AVERAGE_INVENTORY
	PROFIT_MARGIN := s.NET_INCOME / s.NET_SALES
	ASSET_TURNOVER := s.NET_SALES / s.AVERAGE_ASSETS
	RETURN_ON_ASSETS := s.NET_INCOME / s.AVERAGE_ASSETS
	RETURN_ON_EQUITY := s.NET_INCOME / s.AVERAGE_EQUITY
	PAYOUT_RATIO := s.CASH_DIVIDENDS / s.NET_INCOME
	DEBT_TO_TOTAL_ASSETS_RATIO := s.TOTAL_DEBT / s.TOTAL_ASSETS
	TIMES_INTEREST_EARNED := s.EBITDA / s.INTEREST_EXPENSE
	RETURN_ON_COMMON_STOCKHOLDERS_EQUITY := (s.NET_INCOME - s.PREFERRED_DIVIDENDS) / s.AVERAGE_COMMON_STOCKHOLDERS_EQUITY
	EARNINGS_PER_SHARE := (s.NET_INCOME - s.PREFERRED_DIVIDENDS) / s.WEIGHTED_AVERAGE_COMMON_SHARES_OUTSTANDING
	PRICE_EARNINGS_RATIO := s.MARKET_PRICE_PER_SHARES_OUTSTANDING / EARNINGS_PER_SHARE
	return FINANCIAL_ANALYSIS_STATEMENT{
		CURRENT_RATIO:                        CURRENT_RATIO,
		ACID_TEST:                            ACID_TEST,
		RECEIVABLES_TURNOVER:                 RECEIVABLES_TURNOVER,
		INVENTORY_TURNOVER:                   INVENTORY_TURNOVER,
		ASSET_TURNOVER:                       ASSET_TURNOVER,
		PROFIT_MARGIN:                        PROFIT_MARGIN,
		RETURN_ON_ASSETS:                     RETURN_ON_ASSETS,
		RETURN_ON_EQUITY:                     RETURN_ON_EQUITY,
		PAYOUT_RATIO:                         PAYOUT_RATIO,
		DEBT_TO_TOTAL_ASSETS_RATIO:           DEBT_TO_TOTAL_ASSETS_RATIO,
		TIMES_INTEREST_EARNED:                TIMES_INTEREST_EARNED,
		RETURN_ON_COMMON_STOCKHOLDERS_EQUITY: RETURN_ON_COMMON_STOCKHOLDERS_EQUITY,
		EARNINGS_PER_SHARE:                   EARNINGS_PER_SHARE,
		PRICE_EARNINGS_RATIO:                 PRICE_EARNINGS_RATIO}
}

func (s FINANCIAL_ACCOUNTING) FINANCIAL_STATEMENTS(start_date, end_date time.Time, periods int, names []string, in_names bool) ([]map[string]map[string]map[string]map[string]map[string]float64, []FINANCIAL_ANALYSIS_STATEMENT, []JOURNAL_TAG) {
	check_dates(start_date, end_date)
	days := int(end_date.Sub(start_date).Hours() / 24)
	rows, _ := DB.Query("select * from journal order by date,entry_number")
	journal := select_from_journal(rows)
	statements := []map[string]map[string]map[string]map[string]map[string]float64{}
	for a := 0; a < periods; a++ {
		flow_statement, nan_flow_statement := s.statement(journal, start_date.AddDate(0, 0, -days*a), end_date.AddDate(0, 0, -days*a))
		statement := combine_statements(flow_statement, nan_flow_statement)
		statement = s.sum_1st_column(statement)
		statement = s.sum_2nd_column(statement)
		sum_3rd_column(statement, []string{}, []string{}, "all", false)
		sum_3rd_column(statement, names, []string{"all"}, "names", in_names)
		vertical_analysis(statement, float64(days))
		statements = append(statements, statement)
	}
	var all_analysis []FINANCIAL_ANALYSIS_STATEMENT
	for _, statement_current := range statements {
		horizontal_analysis(statement_current, statements[periods-1])
		s.prepare_statement(statement_current)
		calculate_price(statement_current)
		analysis := s.analysis(statement_current)
		all_analysis = append(all_analysis, analysis)
	}
	return statements, all_analysis, journal
}

func (s FINANCIAL_ACCOUNTING) STATEMENT_FILTER(all_financial_statements []map[string]map[string]map[string]map[string]map[string]float64, account_flow_slice, account_slice, name_slice, vpq_slice, number_slice []string,
	in_account_flow_slice, in_account_slice, in_name_slice, in_vpq_slice, in_number_slice bool) [][]FILTERED_STATEMENT {
	var all_statements_struct [][]FILTERED_STATEMENT
	for _, statement := range all_financial_statements {
		var one_statement_struct []FILTERED_STATEMENT
		for key_account_flow, map_account_flow := range statement {
			if IS_IN(key_account_flow, account_flow_slice) == in_account_flow_slice {
				for key_account, map_account := range map_account_flow {
					if IS_IN(key_account, account_slice) == in_account_slice {
						for key_name, map_name := range map_account {
							if IS_IN(key_name, name_slice) == in_name_slice {
								for key_vpq, map_vpq := range map_name {
									if IS_IN(key_vpq, vpq_slice) == in_vpq_slice {
										for key_number, number := range map_vpq {
											if IS_IN(key_number, number_slice) == in_number_slice {
												one_statement_struct = append(one_statement_struct, FILTERED_STATEMENT{key_account_flow, key_account, key_name, key_vpq, key_number, number})
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
		all_statements_struct = append(all_statements_struct, one_statement_struct)
	}
	return all_statements_struct
}

func (s FINANCIAL_ACCOUNTING) SORT_THE_STATMENT(all_statements_struct [][]FILTERED_STATEMENT, sort_by string, is_reverse bool) {
	for _, one_statement_struct := range all_statements_struct {
		switch sort_by {
		case "pre_order":
			s.sort_statement_by_pre_order_in_insertion_sort(one_statement_struct)
		case "father_name":
			s.sort_statement_by_father_name(one_statement_struct)
		case "multiple_alphabet_column":
			s.sort_by_multiple_alphabet_column(one_statement_struct)
		case "number":
			s.sort_by_number(one_statement_struct)
		case "no_order":
		default:
			log.Panic(sort_by, " is not in [pre_order,father_name,multiple_alphabet_column,number,no_order]")
		}
		if is_reverse {
			REVERSE_SLICE(one_statement_struct)
		}
	}
}

func (s FINANCIAL_ACCOUNTING) sort_statement_by_father_name(one_statement_struct []FILTERED_STATEMENT) {
	// later to complete
}

func (s FINANCIAL_ACCOUNTING) sort_by_multiple_alphabet_column(one_statement_struct []FILTERED_STATEMENT) {
	// later to complete
}

func (s FINANCIAL_ACCOUNTING) sort_by_number(one_statement_struct []FILTERED_STATEMENT) {
	sort.Slice(one_statement_struct, func(p, q int) bool { return one_statement_struct[p].NUMBER < one_statement_struct[q].NUMBER })
}

func (s FINANCIAL_ACCOUNTING) sort_statement_by_pre_order_in_insertion_sort(one_statement_struct []FILTERED_STATEMENT) {
	var indexa int
	for _, a := range s.ACCOUNTS {
		for indexb, b := range one_statement_struct {
			if a.NAME == b.KEY_ACCOUNT {
				one_statement_struct[indexa], one_statement_struct[indexb] = one_statement_struct[indexb], one_statement_struct[indexa]
				indexa++
				break
			}
		}
	}
}
