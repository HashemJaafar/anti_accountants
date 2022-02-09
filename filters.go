package anti_accountants

import (
	"sort"
	"strings"
	"time"
)

type FILTERED_STATEMENT struct {
	KEY_ACCOUNT_FLOW, KEY_ACCOUNT, KEY_NAME, KEY_VPQ, KEY_NUMBER string
	NUMBER                                                       float64
}

func (s FINANCIAL_ACCOUNTING) JOURNAL_FILTER(
	JOURNAL_ORDERED_BY_DATE_ENTRY_NUMBER []JOURNAL_TAG,

	filter_date, in_date bool,
	min_date, max_date time.Time,

	filter_entry_number, in_entry_number bool,
	min_entry_number, max_entry_number int,

	filter_account, in_account bool,
	account []string,

	filter_value, in_value bool,
	min_value, max_value float64,

	filter_price, in_price bool,
	min_price, max_price float64,

	filter_quantity, in_quantity bool,
	min_quantity, max_quantity float64,

	filter_barcode, in_barcode bool,
	barcode []string,

	filter_entry_expair, in_entry_expair bool,
	min_entry_expair, max_entry_expair time.Time,

	filter_description, in_description bool,
	description []string,

	filter_name, in_name bool,
	name []string,

	filter_employee_name, in_employee_name bool,
	employee_name []string,

	filter_entry_date, in_entry_date bool,
	min_entry_date, max_entry_date time.Time,

	filter_reverse,
	reverse bool,

) []JOURNAL_TAG {

	check_dates(min_date, max_date)
	check_dates(min_entry_expair, max_entry_expair)
	check_dates(min_entry_date, max_entry_date)

	var filtered_journal []JOURNAL_TAG

	for _, entry := range JOURNAL_ORDERED_BY_DATE_ENTRY_NUMBER {

		date := PARSE_DATE(entry.DATE, s.DATE_LAYOUT)
		entry_expair := PARSE_DATE(entry.ENTRY_EXPAIR, s.DATE_LAYOUT)
		entry_date := PARSE_DATE(entry.ENTRY_DATE, s.DATE_LAYOUT)

		if filter_date || (in_date == (date.After(min_date) && date.Before(max_date))) {
			if filter_entry_number || (in_entry_number == (entry.ENTRY_NUMBER >= min_entry_number && entry.ENTRY_NUMBER <= max_entry_number)) {
				if filter_account || (in_account == IS_IN(entry.ACCOUNT, account)) {
					if filter_value || (in_value == (entry.VALUE >= min_value && entry.VALUE <= max_value)) {
						if filter_price || (in_price == (entry.PRICE >= min_price && entry.PRICE <= max_price)) {
							if filter_quantity || (in_quantity == (entry.QUANTITY >= min_quantity && entry.QUANTITY <= max_quantity)) {
								if filter_barcode || (in_barcode == IS_IN(entry.BARCODE, barcode)) {
									if filter_entry_expair || (in_entry_expair == (entry_expair.After(min_entry_expair) && entry_expair.Before(max_entry_expair))) {
										if filter_description || (in_description == IS_IN(entry.DESCRIPTION, description)) {
											if filter_name || (in_name == IS_IN(entry.NAME, name)) {
												if filter_employee_name || (in_employee_name == IS_IN(entry.EMPLOYEE_NAME, employee_name)) {
													if filter_entry_date || (in_entry_date == (entry_date.After(min_entry_date) && entry_date.Before(max_entry_date))) {
														if filter_reverse || reverse == entry.REVERSE {
															filtered_journal = append(filtered_journal, entry)
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
	return filtered_journal
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
		case "account_number":
			s.sort_statement_by_account_number(one_statement_struct)
		case "multiple_alphabet_column":
			s.sort_by_multiple_alphabet_column(one_statement_struct)
		case "number":
			s.sort_by_number(one_statement_struct)
		default:
			error_element_is_not_in_elements(sort_by, []string{"pre_order", "account_number", "multiple_alphabet_column", "number"})
		}
		if is_reverse {
			REVERSE_SLICE(one_statement_struct)
		}
		s.make_space_before_account_in_statement_struct(one_statement_struct)
	}
}

func (s FINANCIAL_ACCOUNTING) sort_statement_by_account_number(one_statement_struct []FILTERED_STATEMENT) {
	for indexa := range one_statement_struct {
		for indexb := range one_statement_struct {
			if indexa < indexb && !is_it_high_than_by_order(s.account_number(one_statement_struct[indexa].KEY_ACCOUNT), s.account_number(one_statement_struct[indexb].KEY_ACCOUNT)) {
				one_statement_struct[indexa], one_statement_struct[indexb] = one_statement_struct[indexb], one_statement_struct[indexa]
			}
		}
	}
}

func (s FINANCIAL_ACCOUNTING) sort_by_multiple_alphabet_column(one_statement_struct []FILTERED_STATEMENT) { // later to complete
}

func (s FINANCIAL_ACCOUNTING) sort_by_number(one_statement_struct []FILTERED_STATEMENT) {
	sort.Slice(one_statement_struct, func(p, q int) bool { return one_statement_struct[p].NUMBER < one_statement_struct[q].NUMBER })
}

func (s FINANCIAL_ACCOUNTING) sort_statement_by_pre_order_in_insertion_sort(one_statement_struct []FILTERED_STATEMENT) {
	var indexa int
	for _, a := range s.ACCOUNTS {
		for indexb, b := range one_statement_struct {
			if a.ACCOUNT_NAME == b.KEY_ACCOUNT {
				one_statement_struct[indexa], one_statement_struct[indexb] = one_statement_struct[indexb], one_statement_struct[indexa]
				indexa++
				break
			}
		}
	}
}

func (s FINANCIAL_ACCOUNTING) make_space_before_account_in_statement_struct(one_statement_struct []FILTERED_STATEMENT) {
	for indexa, a := range one_statement_struct {
		lenght_of_account_number := len(s.account_number(a.KEY_ACCOUNT))
		one_statement_struct[indexa].KEY_ACCOUNT = strings.Repeat("  ", lenght_of_account_number) + a.KEY_ACCOUNT
	}
}
