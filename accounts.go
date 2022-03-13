package anti_accountants

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"text/tabwriter"
)

func account_struct_from_name(account_name string) ACCOUNT {
	for _, a := range ACCOUNTS {
		if a.ACCOUNT_NAME == account_name {
			return a
		}
	}
	error_account_is_not_listed(account_name)
	return ACCOUNT{}
}

func account_struct_from_number(account_number []uint) ACCOUNT {
	for _, a := range ACCOUNTS {
		if reflect.DeepEqual(a.ACCOUNT_NUMBER[INDEX_OF_ACCOUNT_NUMBER], account_number) {
			return a
		}
	}
	error_account_is_not_listed(account_number)
	return ACCOUNT{}
}

func is_it_sub_account_using_name(higher_level_account, lower_level_account string) bool {
	a1 := account_struct_from_name(higher_level_account).ACCOUNT_NUMBER[INDEX_OF_ACCOUNT_NUMBER]
	a2 := account_struct_from_name(lower_level_account).ACCOUNT_NUMBER[INDEX_OF_ACCOUNT_NUMBER]
	return is_it_sub_account_using_number(a1, a2)
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

func is_it_first_higher_level_account_using_number(higher_level_account_number, lower_level_account_number []uint) bool {
	return reflect.DeepEqual(higher_level_account_number, lower_level_account_number[:len(lower_level_account_number)-1])
}

func check_if_the_tree_connected() {
	l := len(ACCOUNTS[0].ACCOUNT_NUMBER)
	for i := 0; i < l; i++ {
	big_loop:
		for _, a := range ACCOUNTS {
			if len(a.ACCOUNT_NUMBER[i]) > 1 {
				for _, b := range ACCOUNTS {
					if is_it_first_higher_level_account_using_number(b.ACCOUNT_NUMBER[i], a.ACCOUNT_NUMBER[i]) {
						continue big_loop
					}
				}
				error_not_connected_tree(a)
			}
		}
	}
}

func set_cost_flow_type() {
	for indexa, a := range ACCOUNTS {
		is_in_cost_flow_type := is_in(a.COST_FLOW_TYPE, cost_flow_type)
		is_in_receivables := is_in(a.ACCOUNT_NAME, PRIMARY_ACCOUNTS_NAMES.RECEIVABLES)
		is_in_liabilities := is_in(a.ACCOUNT_NAME, PRIMARY_ACCOUNTS_NAMES.LIABILITIES)
		cost_flow_rules := !a.IS_CREDIT && !a.IS_TEMPORARY && a.is_low_level_account && !is_in_receivables && !is_in_liabilities
		if !cost_flow_rules {
			ACCOUNTS[indexa].COST_FLOW_TYPE = ""
		} else if !is_in_cost_flow_type {
			ACCOUNTS[indexa].COST_FLOW_TYPE = "fifo"
		}
	}
}

func inventory_accounts() {
	for _, a := range ACCOUNTS {
		if is_in(a.COST_FLOW_TYPE, cost_flow_type) {
			inventory = append(inventory, a.ACCOUNT_NAME)
		}
	}
}

func sort_the_accounts_by_account_number() {
	for index := range ACCOUNTS {
		for indexb := range ACCOUNTS {
			if index < indexb {
				if !is_it_high_than_by_order(ACCOUNTS[index].ACCOUNT_NUMBER[INDEX_OF_ACCOUNT_NUMBER], ACCOUNTS[indexb].ACCOUNT_NUMBER[INDEX_OF_ACCOUNT_NUMBER]) {
					ACCOUNTS[index], ACCOUNTS[indexb] = ACCOUNTS[indexb], ACCOUNTS[index]
				}
			}
		}
	}
	print_formated_accounts()
}

func print_formated_accounts() {
	p := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	for _, a := range ACCOUNTS {
		is_low_level_account := a.is_low_level_account
		is_credit := "\t," + strconv.FormatBool(a.IS_CREDIT)
		is_temporary := "\t," + strconv.FormatBool(a.IS_TEMPORARY)
		cost_flow_type := "\t,\"" + a.COST_FLOW_TYPE + "\""
		account_name := "\t,\"" + a.ACCOUNT_NAME + "\""
		description := "\t,\"" + a.DESCRIPTION + "\""
		account_number := "\t,[][]uint{" + format_account_number_to_string(a) + "}"
		account_levels := "\t,[]uint{" + format_account_levels_to_string(a) + "}"
		father_and_grandpa_accounts_name := "\t,[][]string{" + format_father_and_grandpa_accounts_name_to_string(a) + "}"
		image := "\t,\"" + a.IMAGE + "\""
		fmt.Fprintln(p, "{", is_low_level_account, is_credit, is_temporary, cost_flow_type, account_name, description, account_number, account_levels, father_and_grandpa_accounts_name, image, "},")
	}
	p.Flush()
}

func format_father_and_grandpa_accounts_name_to_string(a ACCOUNT) string {
	var father_and_grandpa_accounts_name string
	for _, b := range a.father_and_grandpa_accounts_name {
		father_and_grandpa_accounts_name += "{"
		for _, c := range b {
			father_and_grandpa_accounts_name += "\"" + c + "\","
		}
		father_and_grandpa_accounts_name += "}\t,"
	}
	return father_and_grandpa_accounts_name
}

func format_account_levels_to_string(a ACCOUNT) string {
	var account_levels string
	for _, b := range a.ACCOUNT_NUMBER {
		account_levels += strconv.Itoa(len(b)) + ","
	}
	return account_levels
}

func format_account_number_to_string(a ACCOUNT) string {
	var account_number string
	for _, b := range a.ACCOUNT_NUMBER {
		account_number += "{"
		for _, c := range b {
			account_number += strconv.Itoa(int(c)) + ","
		}
		account_number += "}\t,"
	}
	return account_number
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

func init_account_numbers_and_father_and_grandpa_accounts_name() {
	max_len := max_len_for_account_number()
	for indexa, a := range ACCOUNTS {
		ACCOUNTS[indexa].father_and_grandpa_accounts_name = make([][]string, max_len)
		ACCOUNTS[indexa].ACCOUNT_NUMBER = make([][]uint, max_len)
		for indexb, b := range a.ACCOUNT_NUMBER {
			ACCOUNTS[indexa].ACCOUNT_NUMBER[indexb] = b
		}
	}
}

func max_len_for_account_number() int {
	var max_len int
	for _, a := range ACCOUNTS {
		length := len(a.ACCOUNT_NUMBER)
		if length > max_len {
			max_len = length
		}
	}
	return max_len
}

func remove_duplicate_accounts_name() {
	var indexa, indexb int
	for indexa < len(ACCOUNTS) {
		for indexb < len(ACCOUNTS) {
			if indexa < indexb && ACCOUNTS[indexa].ACCOUNT_NAME == ACCOUNTS[indexb].ACCOUNT_NAME {
				ACCOUNTS = append(ACCOUNTS[:indexb], ACCOUNTS[indexb+1:]...)
			} else {
				indexb++
			}
		}
		indexb = 0
		indexa++
	}
}

func remove_duplicate_accounts_number() {
	l := len(ACCOUNTS[0].ACCOUNT_NUMBER)
	for i := 0; i < l; i++ {
		for indexa, a := range ACCOUNTS {
			if len(a.ACCOUNT_NUMBER[i]) != 0 {
				for indexb, b := range ACCOUNTS {
					if len(b.ACCOUNT_NUMBER[i]) != 0 {
						if indexa < indexb && reflect.DeepEqual(a.ACCOUNT_NUMBER[i], b.ACCOUNT_NUMBER[i]) {
							ACCOUNTS[indexb].ACCOUNT_NUMBER[i] = []uint{}
						}
					}
				}
			}
		}
	}
}

func set_low_level_accounts() {
big_loop:
	for indexa, a := range ACCOUNTS {
		for _, b := range ACCOUNTS {
			if is_it_sub_account_using_number(a.ACCOUNT_NUMBER[0], b.ACCOUNT_NUMBER[0]) {
				continue big_loop
			}
		}
		ACCOUNTS[indexa].is_low_level_account = true
	}
}

func check_if_low_level_account_for_all() {
	l := len(ACCOUNTS[0].ACCOUNT_NUMBER)
	for i := 1; i < l; i++ {
	big_loop:
		for indexa, a := range ACCOUNTS {
			if len(a.ACCOUNT_NUMBER[i]) != 0 {
				for _, b := range ACCOUNTS {
					if len(b.ACCOUNT_NUMBER[i]) != 0 {
						if is_it_sub_account_using_number(a.ACCOUNT_NUMBER[i], b.ACCOUNT_NUMBER[i]) {
							continue big_loop
						}
					}
				}
				if !ACCOUNTS[indexa].is_low_level_account {
					log.Fatal("should be low level or high level account in all account numbers ", ACCOUNTS[indexa])
				}
			}
		}
	}
}

func set_account_levels() {
	for indexa, a := range ACCOUNTS {
		ACCOUNTS[indexa].account_levels = []uint{}
		for _, b := range a.ACCOUNT_NUMBER {
			ACCOUNTS[indexa].account_levels = append(ACCOUNTS[indexa].account_levels, uint(len(b)))
		}
	}
}

func set_father_and_grandpa_accounts_name() {
	l := len(ACCOUNTS[0].ACCOUNT_NUMBER)
	for i := 0; i < l; i++ {
		for indexa, a := range ACCOUNTS {
			if len(a.ACCOUNT_NUMBER[i]) > 1 {
				for _, b := range ACCOUNTS {
					if len(b.ACCOUNT_NUMBER[i]) > 0 {
						if is_it_sub_account_using_number(b.ACCOUNT_NUMBER[i], a.ACCOUNT_NUMBER[i]) {
							ACCOUNTS[indexa].father_and_grandpa_accounts_name[i] = append(ACCOUNTS[indexa].father_and_grandpa_accounts_name[i], b.ACCOUNT_NAME)
						}
					}
				}
			}
		}
	}
}

func set_high_level_account_to_permanent() {
	for indexa, a := range ACCOUNTS {
		if !a.is_low_level_account {
			ACCOUNTS[indexa].IS_TEMPORARY = false
		}
	}
}
