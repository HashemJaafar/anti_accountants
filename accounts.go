package anti_accountants

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"text/tabwriter"
)

var cost_flow_type = []string{"fifo", "lifo", "wma", "barcode"}

type ACCOUNT struct {
	IS_CREDIT                    bool
	COST_FLOW_TYPE, ACCOUNT_NAME string
	ACCOUNT_NUMBER               []uint
}

func (s FINANCIAL_ACCOUNTING) is_credit(account_name string) bool {
	for _, a := range s.ACCOUNTS {
		if a.ACCOUNT_NAME == account_name {
			return a.IS_CREDIT
		}
	}
	error_account_name_is_not_listed(account_name)
	return false
}

func (s FINANCIAL_ACCOUNTING) return_cost_flow_type(account_name string) string {
	for _, a := range s.ACCOUNTS {
		if a.ACCOUNT_NAME == account_name {
			return a.COST_FLOW_TYPE
		}
	}
	error_account_name_is_not_listed(account_name)
	return ""
}

func (s FINANCIAL_ACCOUNTING) account_number(account_name string) []uint {
	for _, a := range s.ACCOUNTS {
		if a.ACCOUNT_NAME == account_name {
			return a.ACCOUNT_NUMBER
		}
	}
	error_account_name_is_not_listed(account_name)
	return []uint{}
}

func (s FINANCIAL_ACCOUNTING) is_it_sub_account_using_name(higher_level_account, lower_level_account string) bool {
	return is_it_sub_account_using_number(s.account_number(higher_level_account), s.account_number(lower_level_account))
}

func is_it_sub_account_using_number(higher_level_account_number, lower_level_account_number []uint) bool {
	if reflect.DeepEqual(higher_level_account_number, lower_level_account_number) {
		return false
	}
	if !is_shorter_than(higher_level_account_number, lower_level_account_number) {
		return false
	}
	for i, h := range higher_level_account_number {
		if h != lower_level_account_number[i] {
			return false
		}
	}
	return true
}

func (s FINANCIAL_ACCOUNTING) is_it_high_by_level(account_number []uint) bool {
	for _, a := range s.ACCOUNTS {
		if is_it_sub_account_using_number(account_number, a.ACCOUNT_NUMBER) {
			return true
		}
	}
	return false
}

func (s FINANCIAL_ACCOUNTING) find_all_higher_level_accounts(account_name string) []string {
	account_number := s.account_number(account_name)
	var higher_level_accounts []string
	for _, a := range s.ACCOUNTS {
		if is_it_sub_account_using_number(a.ACCOUNT_NUMBER, account_number) {
			higher_level_accounts = append(higher_level_accounts, a.ACCOUNT_NAME)
		}
	}
	return higher_level_accounts
}

func is_it_first_sub_level_account_using_number(higher_level_account_number, lower_level_account_number []uint) bool {
	if len(higher_level_account_number)+1 != len(lower_level_account_number) {
		return false
	}
	return is_it_sub_account_using_number(higher_level_account_number, lower_level_account_number)
}

func (s FINANCIAL_ACCOUNTING) check_if_the_tree_connected() {
big_loop:
	for _, a := range s.ACCOUNTS {
		if len(a.ACCOUNT_NUMBER) > 1 {
			for _, b := range s.ACCOUNTS {
				if is_it_first_sub_level_account_using_number(b.ACCOUNT_NUMBER, a.ACCOUNT_NUMBER) {
					continue big_loop
				}
			}
			error_not_connected_tree(a)
		}
	}
}

func (s FINANCIAL_ACCOUNTING) check_cost_flow_type() {
	retained_earnings := s.account_number(s.RETAINED_EARNINGS)
	receivables := s.account_number(s.RECEIVABLES)
	liabilities := s.account_number(s.LIABILITIES)
	for _, a := range s.ACCOUNTS {
		is_in_cost_flow_type := IS_IN(a.COST_FLOW_TYPE, cost_flow_type)
		is_it_sub_from_retained_earnings := is_it_sub_account_using_number(retained_earnings, a.ACCOUNT_NUMBER) || reflect.DeepEqual(retained_earnings, a.ACCOUNT_NUMBER)
		is_it_sub_from_receivables := is_it_sub_account_using_number(receivables, a.ACCOUNT_NUMBER) || reflect.DeepEqual(receivables, a.ACCOUNT_NUMBER)
		is_it_sub_from_liabilities := is_it_sub_account_using_number(liabilities, a.ACCOUNT_NUMBER) || reflect.DeepEqual(liabilities, a.ACCOUNT_NUMBER)
		is_it_high_by_level := s.is_it_high_by_level(a.ACCOUNT_NUMBER) && len(a.ACCOUNT_NUMBER) != 0
		if is_in_cost_flow_type {
			if a.IS_CREDIT {
				error_cost_flow_type_used_with___account(a, "credit")
			}
			if is_it_sub_from_retained_earnings {
				error_cost_flow_type_used_with___account(a, "temporary")
			}
			if is_it_sub_from_receivables {
				error_cost_flow_type_used_with___account(a, "receivables")
			}
			if is_it_sub_from_liabilities {
				error_cost_flow_type_used_with___account(a, "liabilities")
			}
			if is_it_high_by_level {
				error_cost_flow_type_used_with___account(a, "high level")
			}
		} else if a.COST_FLOW_TYPE != "" {
			error_element_is_not_in_elements(a.COST_FLOW_TYPE, cost_flow_type)
		}
		if !is_in_cost_flow_type && !is_it_sub_from_retained_earnings && !is_it_sub_from_receivables && !is_it_sub_from_liabilities && !is_it_high_by_level && !a.IS_CREDIT {
			error_you_should_use_cost_flow_type(a.ACCOUNT_NAME)
		}
	}
}

func (s FINANCIAL_ACCOUNTING) check_if_duplicated() {
	for indexa, a := range s.ACCOUNTS {
		not_empty_account_number := len(a.ACCOUNT_NUMBER) != 0
		for indexb, b := range s.ACCOUNTS {
			if indexa != indexb {
				if reflect.DeepEqual(a, b) {
					error_duplicate_value(a)
				}
				if not_empty_account_number && reflect.DeepEqual(a.ACCOUNT_NUMBER, b.ACCOUNT_NUMBER) {
					error_duplicate_value(a.ACCOUNT_NUMBER)
				}
				if a.ACCOUNT_NAME == b.ACCOUNT_NAME {
					error_duplicate_value(a.ACCOUNT_NAME)
				}
			}
		}
	}
}

func (s FINANCIAL_ACCOUNTING) check_if_the_tree_ordered() {
	switch {
	case !s.is_it_sub_account_using_name(s.ASSETS, s.CURRENT_ASSETS):
		error_should_be_one_of_the_fathers(s.ASSETS, s.CURRENT_ASSETS)
	case !s.is_it_sub_account_using_name(s.CURRENT_ASSETS, s.CASH_AND_CASH_EQUIVALENTS):
		error_should_be_one_of_the_fathers(s.CURRENT_ASSETS, s.CASH_AND_CASH_EQUIVALENTS)
	case !s.is_it_sub_account_using_name(s.CURRENT_ASSETS, s.SHORT_TERM_INVESTMENTS):
		error_should_be_one_of_the_fathers(s.CURRENT_ASSETS, s.SHORT_TERM_INVESTMENTS)
	case !s.is_it_sub_account_using_name(s.CURRENT_ASSETS, s.RECEIVABLES):
		error_should_be_one_of_the_fathers(s.CURRENT_ASSETS, s.RECEIVABLES)
	case !s.is_it_sub_account_using_name(s.CURRENT_ASSETS, s.INVENTORY):
		error_should_be_one_of_the_fathers(s.CURRENT_ASSETS, s.INVENTORY)
	case !s.is_it_sub_account_using_name(s.LIABILITIES, s.CURRENT_LIABILITIES):
		error_should_be_one_of_the_fathers(s.LIABILITIES, s.CURRENT_LIABILITIES)
	case !s.is_it_sub_account_using_name(s.EQUITY, s.RETAINED_EARNINGS):
		error_should_be_one_of_the_fathers(s.EQUITY, s.RETAINED_EARNINGS)
	case !s.is_it_sub_account_using_name(s.RETAINED_EARNINGS, s.DIVIDENDS):
		error_should_be_one_of_the_fathers(s.RETAINED_EARNINGS, s.DIVIDENDS)
	case !s.is_it_sub_account_using_name(s.RETAINED_EARNINGS, s.INCOME_STATEMENT):
		error_should_be_one_of_the_fathers(s.RETAINED_EARNINGS, s.INCOME_STATEMENT)
	case !s.is_it_sub_account_using_name(s.INCOME_STATEMENT, s.EBITDA):
		error_should_be_one_of_the_fathers(s.INCOME_STATEMENT, s.EBITDA)
	case !s.is_it_sub_account_using_name(s.INCOME_STATEMENT, s.INTEREST_EXPENSE):
		error_should_be_one_of_the_fathers(s.INCOME_STATEMENT, s.INTEREST_EXPENSE)
	case !s.is_it_sub_account_using_name(s.EBITDA, s.SALES):
		error_should_be_one_of_the_fathers(s.EBITDA, s.SALES)
	case !s.is_it_sub_account_using_name(s.EBITDA, s.COST_OF_GOODS_SOLD):
		error_should_be_one_of_the_fathers(s.EBITDA, s.COST_OF_GOODS_SOLD)
	case !s.is_it_sub_account_using_name(s.EBITDA, s.DISCOUNTS):
		error_should_be_one_of_the_fathers(s.EBITDA, s.DISCOUNTS)
	case !s.is_it_sub_account_using_name(s.DISCOUNTS, s.INVOICE_DISCOUNT):
		error_should_be_one_of_the_fathers(s.DISCOUNTS, s.INVOICE_DISCOUNT)
	}
}

func (s FINANCIAL_ACCOUNTING) inventory_accounts() []string {
	var inventory []string
	for _, a := range s.ACCOUNTS {
		if IS_IN(a.COST_FLOW_TYPE, cost_flow_type) {
			inventory = append(inventory, a.ACCOUNT_NAME)
		}
	}
	return inventory
}

func (s FINANCIAL_ACCOUNTING) SORT_THE_ACCOUNT_BY_ACCOUNT_NUMBER() {
	for index := range s.ACCOUNTS {
		for indexb := range s.ACCOUNTS {
			if index < indexb {
				if !is_it_high_than_by_order(s.ACCOUNTS[index].ACCOUNT_NUMBER, s.ACCOUNTS[indexb].ACCOUNT_NUMBER) {
					s.ACCOUNTS[index], s.ACCOUNTS[indexb] = s.ACCOUNTS[indexb], s.ACCOUNTS[index]
				}
			}
		}
	}
	s.print_formated_accounts()
}

func (s FINANCIAL_ACCOUNTING) print_formated_accounts() {
	p := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	for _, a := range s.ACCOUNTS {
		var account_number string
		for _, b := range a.ACCOUNT_NUMBER {
			account_number += strconv.Itoa(int(b)) + ","
		}
		fmt.Fprintln(p, "{", a.IS_CREDIT, "\t", ",\""+a.COST_FLOW_TYPE+"\"", "\t", ",\""+a.ACCOUNT_NAME+"\"", "\t", ",[]uint{", account_number, "}", "\t", "},")
	}
	p.Flush()
}

func is_it_high_than_by_order(account_number1, account_number2 []uint) bool {
	var short_number int
	if is_shorter_than(account_number1, account_number2) {
		short_number = len(account_number1)
	} else {
		short_number = len(account_number2)
	}
	for index := 0; index < short_number; index++ {
		if account_number1[index] < account_number2[index] {
			return true
		} else if account_number1[index] > account_number2[index] {
			return false
		}
	}
	return is_shorter_than(account_number1, account_number2)
}

func is_shorter_than(slice1, slice2 []uint) bool {
	if len(slice1) < len(slice2) {
		return true
	} else {
		return false
	}
}
