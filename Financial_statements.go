package anti_accountants

import (
	"math"
	"time"
)

type financial_analysis struct {
	Current_assets,
	Current_liabilities,
	Cash,
	Short_term_investments,
	Net_receivables,
	Net_credit_sales,
	Average_net_receivables,
	Cost_of_goods_sold,
	Average_inventory,
	Net_income,
	Net_sales,
	Average_assets,
	Average_equity,
	Preferred_dividends,
	Average_common_stockholders_equity,
	Market_price_per_shares_outstanding,
	Cash_dividends,
	Total_debt,
	Total_assets,
	Ebitda,
	Interest_expense,
	Weighted_average_common_shares_outstanding float64
}

type financial_analysis_statement struct {
	Current_ratio                        float64 // current_assets / current_liabilities
	Acid_test                            float64 // (cash + short_term_investments + net_receivables) / current_liabilities
	Receivables_turnover                 float64 // net_credit_sales / average_net_receivables
	Inventory_turnover                   float64 // cost_of_goods_sold / average_inventory
	Asset_turnover                       float64 // net_sales / average_assets
	Profit_margin                        float64 // net_income / net_sales
	Return_on_assets                     float64 // net_income / average_assets
	Return_on_equity                     float64 // net_income / average_equity
	Payout_ratio                         float64 // cash_dividends / net_income
	Debt_to_total_assets_ratio           float64 // total_debt / total_assets
	Times_interest_earned                float64 // ebitda / interest_expense
	Return_on_common_stockholders_equity float64 // (net_income - preferred_dividends) / average_common_stockholders_equity
	Earnings_per_share                   float64 // (net_income - preferred_dividends) / weighted_average_common_shares_outstanding
	Price_earnings_ratio                 float64 // market_price_per_shares_outstanding / earnings_per_share
}

type filtered_statement struct {
	Key_account_flow, Key_account, Key_name, Key_vpq, Key_number string
	Number                                                       float64
}

func ending_balance(statement map[string]map[string]map[string]map[string]map[string]float64, key_account_flow, key_account, key_name, key_vpq string) float64 {
	return statement[key_account_flow][key_account][key_name][key_vpq]["beginning_balance"] + statement[key_account][key_account][key_name][key_vpq]["increase"] - statement[key_account][key_account][key_name][key_vpq]["decrease"]
}

func sum_flows(b journal_tag, x float64, map_v, map_q map[string]float64) {
	if b.Value*x < 0 {
		map_v["outflow"] += math.Abs(b.Value)
		map_q["outflow"] += math.Abs(b.Quantity)
	} else {
		map_v["inflow"] += math.Abs(b.Value)
		map_q["inflow"] += math.Abs(b.Quantity)
	}
}

func (s Financial_accounting) sum_values(date, start_date time.Time, one_simple_entry []journal_tag, nan_flow_statement map[string]map[string]map[string]map[string]float64) {
	for _, b := range one_simple_entry {
		map_v1 := initialize_map_3(nan_flow_statement, b.Account, b.Name, "value")
		map_q1 := initialize_map_3(nan_flow_statement, b.Account, b.Name, "quantity")
		map_v2 := initialize_map_3(nan_flow_statement, s.Retained_earnings, b.Name, "value")
		map_q2 := initialize_map_3(nan_flow_statement, s.Retained_earnings, b.Name, "quantity")
		if date.Before(start_date) {
			switch {
			case s.is_father(s.Retained_earnings, b.Account) && s.is_credit(b.Account):
				map_v2["beginning_balance"] += b.Value
				map_q2["beginning_balance"] += b.Quantity
			case s.is_father(s.Retained_earnings, b.Account) && !s.is_credit(b.Account):
				map_v2["beginning_balance"] -= b.Value
				map_q2["beginning_balance"] -= b.Quantity
			default:
				map_v1["beginning_balance"] += b.Value
				map_q1["beginning_balance"] += b.Quantity
			}
		}
		if date.After(start_date) {
			if b.Value >= 0 {
				map_v1["increase"] += math.Abs(b.Value)
				map_q1["increase"] += math.Abs(b.Quantity)
			} else {
				map_v1["decrease"] += math.Abs(b.Value)
				map_q1["decrease"] += math.Abs(b.Quantity)
			}
		}
	}
}

func (s Financial_accounting) sum_flow(date, start_date time.Time, one_simple_entry []journal_tag, flow_statement map[string]map[string]map[string]map[string]map[string]float64) {
	for _, a := range one_simple_entry {
		for _, b := range one_simple_entry {
			map_v := initialize_map_4(flow_statement, a.Account, b.Account, b.Name, "value")
			map_q := initialize_map_4(flow_statement, a.Account, b.Account, b.Name, "quantity")
			if date.After(start_date) {
				if b.Account == a.Account || s.is_credit(b.Account) != s.is_credit(a.Account) {
					sum_flows(b, 1, map_v, map_q)
				} else {
					sum_flows(b, -1, map_v, map_q)
				}
			}
		}
	}
}

func (s Financial_accounting) analysis(statement map[string]map[string]map[string]map[string]map[string]float64) financial_analysis_statement {
	return financial_analysis{
		Current_assets:                      statement[s.Cash_and_cash_equivalents][s.Current_assets]["names"]["value"]["ending_balance"],
		Current_liabilities:                 statement[s.Cash_and_cash_equivalents][s.Current_liabilities]["names"]["value"]["ending_balance"],
		Cash:                                statement[s.Cash_and_cash_equivalents][s.Cash_and_cash_equivalents]["names"]["value"]["ending_balance"],
		Short_term_investments:              statement[s.Cash_and_cash_equivalents][s.Short_term_investments]["names"]["value"]["ending_balance"],
		Net_receivables:                     statement[s.Cash_and_cash_equivalents][s.Receivables]["names"]["value"]["ending_balance"],
		Net_credit_sales:                    statement[s.Sales][s.Receivables]["names"]["value"]["flow"],
		Average_net_receivables:             statement[s.Cash_and_cash_equivalents][s.Receivables]["names"]["value"]["average"],
		Cost_of_goods_sold:                  statement[s.Cash_and_cash_equivalents][s.Cost_of_goods_sold]["names"]["value"]["ending_balance"],
		Average_inventory:                   statement[s.Cash_and_cash_equivalents][s.Inventory]["names"]["value"]["average"],
		Net_income:                          statement[s.Cash_and_cash_equivalents][s.Income_statement]["names"]["value"]["ending_balance"],
		Net_sales:                           statement[s.Cash_and_cash_equivalents][s.Sales]["names"]["value"]["ending_balance"],
		Average_assets:                      statement[s.Cash_and_cash_equivalents][s.Assets]["names"]["value"]["average"],
		Average_equity:                      statement[s.Cash_and_cash_equivalents][s.Equity]["names"]["value"]["average"],
		Preferred_dividends:                 0,
		Average_common_stockholders_equity:  0,
		Market_price_per_shares_outstanding: 0,
		Cash_dividends:                      statement[s.Cash_and_cash_equivalents][s.Dividends]["names"]["value"]["flow"],
		Total_debt:                          statement[s.Cash_and_cash_equivalents][s.Liabilities]["names"]["value"]["ending_balance"],
		Total_assets:                        statement[s.Cash_and_cash_equivalents][s.Assets]["names"]["value"]["ending_balance"],
		Ebitda:                              statement[s.Cash_and_cash_equivalents][s.Ebitda]["names"]["value"]["ending_balance"],
		Interest_expense:                    statement[s.Cash_and_cash_equivalents][s.Interest_expense]["names"]["value"]["ending_balance"],
		Weighted_average_common_shares_outstanding: 0,
	}.Financial_analysis_statement()
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

func (s Financial_accounting) prepare_statement(statement map[string]map[string]map[string]map[string]map[string]float64) {
	for key_account_flow, map_account_flow := range statement {
		if key_account_flow == s.Cash_and_cash_equivalents {
			for key_account, map_account := range map_account_flow {
				for key_name, map_name := range map_account {
					for key_vpq, map_vpq := range map_name {
						map_vpq1 := initialize_map_4(statement, "financial_statement", key_account, key_name, key_vpq)
						for key_number, number := range map_vpq {
							map_vpq1[key_number] = number
							if !s.is_father(s.Income_statement, key_account) {
								map_vpq1["percent"] = statement[s.Income_statement][key_account][key_name][key_vpq]["percent"]
							} else {
								map_vpq1["percent"] = statement[s.Assets][key_account][key_name][key_vpq]["percent"]
							}
							switch {
							case s.is_father(s.Inventory, key_account):
								map_vpq1["turnover"] = statement[s.Cost_of_goods_sold][key_account][key_name][key_vpq]["turnover"]
								map_vpq1["turnover_days"] = statement[s.Cost_of_goods_sold][key_account][key_name][key_vpq]["turnover_days"]
							case s.is_father(s.Assets, key_account):
								map_vpq1["turnover"] = statement[s.Sales][key_account][key_name][key_vpq]["turnover"]
								map_vpq1["turnover_days"] = statement[s.Sales][key_account][key_name][key_vpq]["turnover_days"]
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

func (s Financial_accounting) sum_2nd_column(statement map[string]map[string]map[string]map[string]map[string]float64) map[string]map[string]map[string]map[string]map[string]float64 {
	new_statement := map[string]map[string]map[string]map[string]map[string]float64{}
	for key_account_flow, map_account_flow := range statement {
		for key_account, map_account := range map_account_flow {
			var last_name string
			key1 := key_account
			for {
				for _, ss := range s.Accounts {
					if ss.Name == key_account {
						key_account = ss.Father
						for key_name, map_name := range map_account {
							for key_vpq, map_vpq := range map_name {
								map_vpq1 := initialize_map_4(new_statement, key_account_flow, ss.Name, key_name, key_vpq)
								for key_number, number := range map_vpq {
									switch {
									case !IS_IN(key_number, []string{"inflow", "outflow"}):
										if s.is_credit(key1) == s.is_credit(ss.Name) {
											map_vpq1[key_number] += number
										} else {
											map_vpq1[key_number] -= number
										}
									case key_account_flow != key1:
										map_vpq1[key_number] += number
									case key_account_flow == ss.Name:
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

func (s Financial_accounting) sum_1st_column(statement map[string]map[string]map[string]map[string]map[string]float64) map[string]map[string]map[string]map[string]map[string]float64 {
	new_statement := map[string]map[string]map[string]map[string]map[string]float64{}
	var flow_accounts []string
	for _, a := range s.Accounts {
		for _, b := range s.Accounts {
			if s.is_father(a.Name, b.Name) {
				flow_accounts = append(flow_accounts, b.Name)
			}
		}
		for key_account_flow, map_account_flow := range statement {
			if IS_IN(key_account_flow, flow_accounts) {
				for key_account, map_account := range map_account_flow {
					for key_name, map_name := range map_account {
						for key_vpq, map_vpq := range map_name {
							map_vpq1 := initialize_map_4(new_statement, a.Name, key_account, key_name, key_vpq)
							for key_number, number := range map_vpq {
								switch {
								case IS_IN(key_number, []string{"inflow", "outflow"}):
									if s.is_credit(a.Name) == s.is_credit(key_account_flow) {
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

func (s Financial_accounting) statement(journal []journal_tag, start_date, end_date time.Time) (map[string]map[string]map[string]map[string]map[string]float64, map[string]map[string]map[string]map[string]float64) {
	var one_simple_entry []journal_tag
	var previous_entry_number int
	var date time.Time
	flow_statement := map[string]map[string]map[string]map[string]map[string]float64{}
	nan_flow_statement := map[string]map[string]map[string]map[string]float64{}
	for _, entry := range journal {
		date = Parse_date(entry.Date, s.Date_layout)
		if previous_entry_number != entry.Entry_number {
			s.sum_flow(date, start_date, one_simple_entry, flow_statement)
			s.sum_values(date, start_date, one_simple_entry, nan_flow_statement)
			one_simple_entry = []journal_tag{}
		}
		if date.Before(end_date) {
			one_simple_entry = append(one_simple_entry, entry)
		}
		previous_entry_number = entry.Entry_number
	}
	s.sum_flow(date, start_date, one_simple_entry, flow_statement)
	s.sum_values(date, start_date, one_simple_entry, nan_flow_statement)
	return flow_statement, nan_flow_statement
}

func (s financial_analysis) Financial_analysis_statement() financial_analysis_statement {
	Current_ratio := s.Current_assets / s.Current_liabilities
	Acid_test := (s.Cash + s.Short_term_investments + s.Net_receivables) / s.Current_liabilities
	Receivables_turnover := s.Net_credit_sales / s.Average_net_receivables
	Inventory_turnover := s.Cost_of_goods_sold / s.Average_inventory
	Profit_margin := s.Net_income / s.Net_sales
	Asset_turnover := s.Net_sales / s.Average_assets
	Return_on_assets := s.Net_income / s.Average_assets
	Return_on_equity := s.Net_income / s.Average_equity
	Payout_ratio := s.Cash_dividends / s.Net_income
	Debt_to_total_assets_ratio := s.Total_debt / s.Total_assets
	Times_interest_earned := s.Ebitda / s.Interest_expense
	Return_on_common_stockholders_equity := (s.Net_income - s.Preferred_dividends) / s.Average_common_stockholders_equity
	Earnings_per_share := (s.Net_income - s.Preferred_dividends) / s.Weighted_average_common_shares_outstanding
	Price_earnings_ratio := s.Market_price_per_shares_outstanding / Earnings_per_share
	return financial_analysis_statement{
		Current_ratio:                        Current_ratio,
		Acid_test:                            Acid_test,
		Receivables_turnover:                 Receivables_turnover,
		Inventory_turnover:                   Inventory_turnover,
		Asset_turnover:                       Asset_turnover,
		Profit_margin:                        Profit_margin,
		Return_on_assets:                     Return_on_assets,
		Return_on_equity:                     Return_on_equity,
		Payout_ratio:                         Payout_ratio,
		Debt_to_total_assets_ratio:           Debt_to_total_assets_ratio,
		Times_interest_earned:                Times_interest_earned,
		Return_on_common_stockholders_equity: Return_on_common_stockholders_equity,
		Earnings_per_share:                   Earnings_per_share,
		Price_earnings_ratio:                 Price_earnings_ratio}
}

func (s Financial_accounting) Financial_statements(start_date, end_date time.Time, periods int, names []string, in_names bool) ([]map[string]map[string]map[string]map[string]map[string]float64, []financial_analysis_statement, []journal_tag) {
	check_dates(start_date, end_date)
	days := int(end_date.Sub(start_date).Hours() / 24)
	var journal []journal_tag
	rows, _ := db.Query("select * from journal order by date,entry_number")
	for rows.Next() {
		var entry journal_tag
		rows.Scan(&entry.Date, &entry.Entry_number, &entry.Account, &entry.Value, &entry.Price, &entry.Quantity, &entry.Barcode, &entry.Entry_expair, &entry.Description, &entry.Name, &entry.Employee_name, &entry.Entry_date, &entry.Reverse)
		journal = append(journal, entry)
	}
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
	var all_analysis []financial_analysis_statement
	for _, statement_current := range statements {
		horizontal_analysis(statement_current, statements[periods-1])
		s.prepare_statement(statement_current)
		calculate_price(statement_current)
		analysis := s.analysis(statement_current)
		all_analysis = append(all_analysis, analysis)
	}
	return statements, all_analysis, journal
}

func (s Financial_accounting) Statement_filter(all_financial_statements []map[string]map[string]map[string]map[string]map[string]float64, account_flow_slice, account_slice, name_slice, vpq_slice, number_slice []string,
	in_account_flow_slice, in_account_slice, in_name_slice, in_vpq_slice, in_number_slice bool) [][]filtered_statement {
	var all_statements_struct [][]filtered_statement
	for _, statement := range all_financial_statements {
		var statement_struct []filtered_statement
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
												statement_struct = append(statement_struct, filtered_statement{key_account_flow, key_account, key_name, key_vpq, key_number, number})
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
		var indexa int
		for _, a := range s.Accounts {
			for indexb, b := range statement_struct {
				if a.Name == b.Key_account {
					statement_struct[indexa], statement_struct[indexb] = statement_struct[indexb], statement_struct[indexa]
					indexa++
					break
				}
			}
		}
		all_statements_struct = append(all_statements_struct, statement_struct)
	}
	return all_statements_struct
}
