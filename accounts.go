package anti_accountants

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"text/tabwriter"
)

func account_struct_from_name(account_name string) (ACCOUNT, int, error) {
	for indexa, a := range ACCOUNTS {
		if a.ACCOUNT_NAME == account_name {
			return a, indexa, nil
		}
	}
	return ACCOUNT{}, 0, error_not_listed
}

func account_struct_from_barcode(barcode string) (ACCOUNT, int, error) {
	for indexa, a := range ACCOUNTS {
		if is_in(barcode, a.BARCODE) {
			return a, indexa, nil
		}
	}
	return ACCOUNT{}, 0, error_not_listed
}

func account_struct_from_number(account_number []uint) (ACCOUNT, int, error) {
	for indexa, a := range ACCOUNTS {
		if reflect.DeepEqual(a.ACCOUNT_NUMBER[INDEX_OF_ACCOUNT_NUMBER], account_number) {
			return a, indexa, nil
		}
	}
	return ACCOUNT{}, 0, error_not_listed
}

func is_it_sub_account_using_name(higher_level_account, lower_level_account string) bool {
	a1, _, _ := account_struct_from_name(higher_level_account)
	a2, _, _ := account_struct_from_name(lower_level_account)
	return is_it_sub_account_using_number(a1.ACCOUNT_NUMBER[INDEX_OF_ACCOUNT_NUMBER], a2.ACCOUNT_NUMBER[INDEX_OF_ACCOUNT_NUMBER])
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

func check_if_the_tree_connected() error {
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
				return fmt.Errorf("the account number %f for %f not conected to the tree", a.ACCOUNT_NUMBER[i], a)
			}
		}
	}
	return nil
}

func set_cost_flow_type() {
	for indexa, a := range ACCOUNTS {
		is_in_cost_flow_type := is_in(a.COST_FLOW_TYPE, cost_flow_type)
		is_in_receivables := is_in(a.ACCOUNT_NAME, PRIMARY_ACCOUNTS_NAMES.RECEIVABLES)
		is_in_liabilities := is_in(a.ACCOUNT_NAME, PRIMARY_ACCOUNTS_NAMES.LIABILITIES)
		cost_flow_rules := !a.IS_CREDIT && !a.IS_TEMPORARY && a.IS_LOW_LEVEL_ACCOUNT && !is_in_receivables && !is_in_liabilities
		if !cost_flow_rules {
			ACCOUNTS[indexa].COST_FLOW_TYPE = ""
		} else if !is_in_cost_flow_type {
			ACCOUNTS[indexa].COST_FLOW_TYPE = FIFO
		}
	}
}

func sort_the_accounts_by_account_number() {
	for index := range ACCOUNTS {
		for indexb := range ACCOUNTS {
			if index < indexb && !is_it_high_than_by_order(ACCOUNTS[index].ACCOUNT_NUMBER[INDEX_OF_ACCOUNT_NUMBER], ACCOUNTS[indexb].ACCOUNT_NUMBER[INDEX_OF_ACCOUNT_NUMBER]) {
				ACCOUNTS[index], ACCOUNTS[indexb] = ACCOUNTS[indexb], ACCOUNTS[index]
			}
		}
	}
}

func print_formated_accounts() {
	p := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	for _, a := range ACCOUNTS {
		is_low_level_account := a.IS_LOW_LEVEL_ACCOUNT
		is_credit := "\t," + strconv.FormatBool(a.IS_CREDIT)
		is_temporary := "\t," + strconv.FormatBool(a.IS_TEMPORARY)
		cost_flow_type := "\t,\"" + a.COST_FLOW_TYPE + "\""
		account_name := "\t,\"" + a.ACCOUNT_NAME + "\""
		notes := "\t,\"" + a.NOTES + "\""
		image := "\t," + format_string_slice_to_string(a.IMAGE)
		barcodes := "\t," + format_string_slice_to_string(a.BARCODE)
		account_number := "\t," + format_slice_of_slice_of_uint_to_string(a.ACCOUNT_NUMBER)
		account_levels := "\t," + format_slice_of_uint_to_string(a.ACCOUNT_LEVELS)
		father_and_grandpa_accounts_name := "\t," + format_slice_of_slice_of_string_to_string(a.FATHER_AND_GRANDPA_ACCOUNTS_NAME)
		alert_for_minimum_quantity_by_turnover_in_days := "\t," + strconv.FormatUint(uint64(a.ALERT_FOR_MINIMUM_QUANTITY_BY_TURNOVER_IN_DAYS), 36)
		alert_for_minimum_quantity_by_quintity := "\t," + strconv.FormatFloat(a.ALERT_FOR_MINIMUM_QUANTITY_BY_QUINTITY, 'E', 64, 64)
		target_balance := "\t," + strconv.FormatFloat(a.TARGET_BALANCE, 'E', 64, 64)
		if_the_target_balance_is_less_is_good := "\t," + strconv.FormatBool(a.IF_THE_TARGET_BALANCE_IS_LESS_IS_GOOD)
		fmt.Fprintln(p, "{", is_low_level_account, is_credit, is_temporary, cost_flow_type, account_name, notes,
			image, barcodes, account_number, account_levels, father_and_grandpa_accounts_name,
			alert_for_minimum_quantity_by_turnover_in_days, alert_for_minimum_quantity_by_quintity, target_balance, if_the_target_balance_is_less_is_good, "},")
	}
	p.Flush()
}

func format_slice_of_slice_of_string_to_string(a [][]string) string {
	var str string
	for _, b := range a {
		str += "{"
		for _, c := range b {
			str += "\"" + c + "\","
		}
		str += "}\t,"
	}
	return "[][]string{" + str + "}"
}

func format_slice_of_uint_to_string(a []uint) string {
	var str string
	for _, b := range a {
		str += strconv.Itoa(int(b)) + ","
	}
	return "[]uint{" + str + "}"
}

func format_slice_of_slice_of_uint_to_string(a [][]uint) string {
	var str string
	for _, b := range a {
		str += "{"
		for _, c := range b {
			str += strconv.Itoa(int(c)) + ","
		}
		str += "}\t,"
	}
	return "[][]uint{" + str + "}"
}

func format_string_slice_to_string(a []string) string {
	var str string
	for _, b := range a {
		str += "\"" + b + "\","
	}
	return "[]string{" + str + "}"
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
		ACCOUNTS[indexa].FATHER_AND_GRANDPA_ACCOUNTS_NAME = make([][]string, max_len)
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

func remove_empty_and_duplicate_accounts_name() {
	var indexa, indexb int
	for indexa < len(ACCOUNTS) {
		for indexb < len(ACCOUNTS) {
			if indexa < indexb && ACCOUNTS[indexa].ACCOUNT_NAME == ACCOUNTS[indexb].ACCOUNT_NAME || ACCOUNTS[indexb].ACCOUNT_NAME == "" {
				popup(ACCOUNTS, indexb)
			} else {
				indexb++
			}
		}
		indexb = 0
		indexa++
	}
}

func remove_empty_and_duplicate_accounts_barcode() {
	var barcodes []string
	for indexa := range ACCOUNTS {
		var indexb int
		for indexb < len(ACCOUNTS[indexa].BARCODE) {
			if is_in(ACCOUNTS[indexa].BARCODE[indexb], barcodes) || ACCOUNTS[indexa].BARCODE[indexb] == "" {
				popup(ACCOUNTS[indexa].BARCODE, indexb)
			} else {
				barcodes = append(barcodes, ACCOUNTS[indexa].BARCODE[indexb])
				indexb++
			}
		}
	}
}

func empty_the_duplicate_accounts_number() {
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
		ACCOUNTS[indexa].IS_LOW_LEVEL_ACCOUNT = true
	}
}

func check_if_low_level_account_for_all() error {
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
				if !ACCOUNTS[indexa].IS_LOW_LEVEL_ACCOUNT {
					return fmt.Errorf("should be low level or high level account in all account numbers %f", ACCOUNTS[indexa])
				}
			}
		}
	}
	return nil
}

func set_account_levels() {
	for indexa, a := range ACCOUNTS {
		ACCOUNTS[indexa].ACCOUNT_LEVELS = []uint{}
		for _, b := range a.ACCOUNT_NUMBER {
			ACCOUNTS[indexa].ACCOUNT_LEVELS = append(ACCOUNTS[indexa].ACCOUNT_LEVELS, uint(len(b)))
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
							ACCOUNTS[indexa].FATHER_AND_GRANDPA_ACCOUNTS_NAME[i] = append(ACCOUNTS[indexa].FATHER_AND_GRANDPA_ACCOUNTS_NAME[i], b.ACCOUNT_NAME)
						}
					}
				}
			}
		}
	}
}

func set_high_level_account_to_permanent() {
	for indexa, a := range ACCOUNTS {
		if !a.IS_LOW_LEVEL_ACCOUNT {
			ACCOUNTS[indexa].IS_TEMPORARY = false
		}
	}
}

func find_all_inventory_files() []string {
	files, _ := ioutil.ReadDir(db_inventory)
	var inventory_files_name []string
	for _, f := range files {
		inventory_files_name = append(inventory_files_name, f.Name())
	}
	return inventory_files_name
}

func check_inventory() {
	for _, inventory_account := range find_all_inventory_files() {
		account, index, err := account_struct_from_name(inventory_account)
		if err != nil {
			add_new_account(inventory_account)
		}
		if account.COST_FLOW_TYPE == "" {
			ACCOUNTS[index] = ACCOUNT{
				IS_LOW_LEVEL_ACCOUNT:             true,
				IS_CREDIT:                        false,
				IS_TEMPORARY:                     false,
				COST_FLOW_TYPE:                   FIFO,
				ACCOUNT_NAME:                     inventory_account,
				NOTES:                            "",
				IMAGE:                            []string{},
				BARCODE:                          []string{},
				ACCOUNT_NUMBER:                   max_account_numbers(),
				ACCOUNT_LEVELS:                   []uint{},
				FATHER_AND_GRANDPA_ACCOUNTS_NAME: [][]string{},
				ALERT_FOR_MINIMUM_QUANTITY_BY_TURNOVER_IN_DAYS: 0,
				ALERT_FOR_MINIMUM_QUANTITY_BY_QUINTITY:         0,
				TARGET_BALANCE:                                 0,
				IF_THE_TARGET_BALANCE_IS_LESS_IS_GOOD:          false,
			}
		}
	}
}

func check_journal() {
	for _, a := range db_read_journal() {
		a1, index, err := account_struct_from_name(a.ACCOUNT_DEBIT)
		if err != nil {
			add_new_account(a1.ACCOUNT_NAME)
		} else if !a1.IS_LOW_LEVEL_ACCOUNT {
			ACCOUNTS[index].ACCOUNT_NUMBER = max_account_numbers()
		}
		a2, index, err := account_struct_from_name(a.ACCOUNT_CREDIT)
		if err != nil {
			add_new_account(a2.ACCOUNT_NAME)
		} else if !a2.IS_LOW_LEVEL_ACCOUNT {
			ACCOUNTS[index].ACCOUNT_NUMBER = max_account_numbers()
		}
	}
}

func add_new_account(account_name string) {
	ACCOUNTS = append(ACCOUNTS, ACCOUNT{
		IS_CREDIT:      false,
		COST_FLOW_TYPE: FIFO,
		ACCOUNT_NAME:   account_name,
		NOTES:          "it is auto added",
		ACCOUNT_NUMBER: max_account_numbers(),
	})
}

func max_account_numbers() [][]uint {
	var account_number [][]uint
	l := len(ACCOUNTS[0].ACCOUNT_NUMBER)
	for i := 0; i < l; i++ {
		var max_number uint
		for _, a := range ACCOUNTS {
			if len(a.ACCOUNT_NUMBER[i]) > 0 && max_number < a.ACCOUNT_NUMBER[i][0] {
				max_number = a.ACCOUNT_NUMBER[i][0]
			}
		}
		account_number = append(account_number, []uint{max_number + 1})
	}
	return account_number
}

func INITIALIZE() {
	ACCOUNTS = db_read_accounts()
	var a []ACCOUNT
	for reflect.DeepEqual(a, ACCOUNTS) {
		a = ACCOUNTS
		init_account_numbers_and_father_and_grandpa_accounts_name()
		remove_empty_and_duplicate_accounts_name()
		remove_empty_and_duplicate_accounts_barcode()
		empty_the_duplicate_accounts_number()
		_ = check_if_the_tree_connected()
		set_low_level_accounts()
		_ = check_if_low_level_account_for_all()
		set_high_level_account_to_permanent()
		set_account_levels()
		set_father_and_grandpa_accounts_name()
		set_cost_flow_type()
		sort_the_accounts_by_account_number()
		print_formated_accounts()
		check_journal()
		check_inventory()
	}
	db_insert_into_accounts()
}
