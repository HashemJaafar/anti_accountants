package anti_accountants

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"text/tabwriter"
)

func ACCOUNT_STRUCT_FROM_NAME(account_name string) (ACCOUNT, int, error) {
	for indexa, a := range ACCOUNTS {
		if a.ACCOUNT_NAME == account_name {
			return a, indexa, nil
		}
	}
	return ACCOUNT{}, 0, ERROR_NOT_LISTED
}

func ACCOUNT_STRUCT_FROM_BARCODE(barcode string) (ACCOUNT, int, error) {
	for indexa, a := range ACCOUNTS {
		if IS_IN(barcode, a.BARCODE) {
			return a, indexa, nil
		}
	}
	return ACCOUNT{}, 0, ERROR_NOT_LISTED
}

func ACCOUNT_STRUCT_FROM_NUMBER(account_number []uint) (ACCOUNT, int, error) {
	for indexa, a := range ACCOUNTS {
		if reflect.DeepEqual(a.ACCOUNT_NUMBER[INDEX_OF_ACCOUNT_NUMBER], account_number) {
			return a, indexa, nil
		}
	}
	return ACCOUNT{}, 0, ERROR_NOT_LISTED
}

func IS_IT_SUB_ACCOUNT_USING_NAME(higher_level_account, lower_level_account string) bool {
	a1, _, _ := ACCOUNT_STRUCT_FROM_NAME(higher_level_account)
	a2, _, _ := ACCOUNT_STRUCT_FROM_NAME(lower_level_account)
	return IS_IT_SUB_ACCOUNT_USING_NUMBER(a1.ACCOUNT_NUMBER[INDEX_OF_ACCOUNT_NUMBER], a2.ACCOUNT_NUMBER[INDEX_OF_ACCOUNT_NUMBER])
}

func IS_IT_SUB_ACCOUNT_USING_NUMBER(higher_level_account_number, lower_level_account_number []uint) bool {
	if !IS_IT_POSSIBLE_TO_BE_SUB_ACCOUNT(higher_level_account_number, lower_level_account_number) {
		return false
	}
	for i, h := range higher_level_account_number {
		if h != lower_level_account_number[i] {
			return false
		}
	}
	return true
}

func IS_IT_POSSIBLE_TO_BE_SUB_ACCOUNT(higher_level_account_number, lower_level_account_number []uint) bool {
	len_higher_level_account_number := len(higher_level_account_number)
	len_lower_level_account_number := len(lower_level_account_number)
	if len_higher_level_account_number == 0 || len_lower_level_account_number == 0 {
		return false
	}
	if len_higher_level_account_number >= len_lower_level_account_number {
		return false
	}
	if reflect.DeepEqual(higher_level_account_number, lower_level_account_number) {
		return false
	}
	return true
}

func IS_IT_THE_FATHER(higher_level_account_number, lower_level_account_number []uint) bool {
	if !IS_IT_POSSIBLE_TO_BE_SUB_ACCOUNT(higher_level_account_number, lower_level_account_number) {
		return false
	}
	return reflect.DeepEqual(higher_level_account_number, CUT_THE_SLICE(lower_level_account_number, 1))
}

func SET_COST_FLOW_TYPE() {
	for indexa, a := range ACCOUNTS {
		is_in_cost_flow_type := IS_IN(a.COST_FLOW_TYPE, COST_FLOW_TYPE)
		is_in_receivables := IS_IN(a.ACCOUNT_NAME, PRIMARY_ACCOUNTS_NAMES.RECEIVABLES)
		is_in_liabilities := IS_IN(a.ACCOUNT_NAME, PRIMARY_ACCOUNTS_NAMES.LIABILITIES)
		cost_flow_rules := !a.IS_CREDIT && !a.IS_TEMPORARY && a.IS_LOW_LEVEL_ACCOUNT && !is_in_receivables && !is_in_liabilities
		if !cost_flow_rules {
			ACCOUNTS[indexa].COST_FLOW_TYPE = ""
		} else if !is_in_cost_flow_type {
			ACCOUNTS[indexa].COST_FLOW_TYPE = FIFO
		}
	}
}

func SORT_THE_ACCOUNTS_BY_ACCOUNT_NUMBER() {
	for index := range ACCOUNTS {
		for indexb := range ACCOUNTS {
			if index < indexb && !IS_IT_HIGH_THAN_BY_ORDER(ACCOUNTS[index].ACCOUNT_NUMBER[INDEX_OF_ACCOUNT_NUMBER], ACCOUNTS[indexb].ACCOUNT_NUMBER[INDEX_OF_ACCOUNT_NUMBER]) {
				ACCOUNTS[index], ACCOUNTS[indexb] = ACCOUNTS[indexb], ACCOUNTS[index]
			}
		}
	}
}

func PRINT_FORMATED_ACCOUNTS() {
	p := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	for _, a := range ACCOUNTS {
		is_low_level_account := a.IS_LOW_LEVEL_ACCOUNT
		is_credit := "\t," + fmt.Sprint(a.IS_CREDIT)
		is_temporary := "\t," + fmt.Sprint(a.IS_TEMPORARY)
		COST_FLOW_TYPE := "\t,\"" + a.COST_FLOW_TYPE + "\""
		account_name := "\t,\"" + a.ACCOUNT_NAME + "\""
		notes := "\t,\"" + a.NOTES + "\""
		image := "\t," + FORMAT_STRING_SLICE_TO_STRING(a.IMAGE)
		barcodes := "\t," + FORMAT_STRING_SLICE_TO_STRING(a.BARCODE)
		account_number := "\t," + FORMAT_SLICE_OF_SLICE_OF_UINT_TO_STRING(a.ACCOUNT_NUMBER)
		account_levels := "\t," + FORMAT_SLICE_OF_UINT_TO_STRING(a.ACCOUNT_LEVELS)
		father_and_grandpa_accounts_name := "\t," + FORMAT_SLICE_OF_SLICE_OF_STRING_TO_STRING(a.FATHER_AND_GRANDPA_ACCOUNTS_NAME)
		alert_for_minimum_quantity_by_turnover_in_days := "\t," + fmt.Sprint(a.ALERT_FOR_MINIMUM_QUANTITY_BY_TURNOVER_IN_DAYS)
		alert_for_minimum_quantity_by_quintity := "\t," + fmt.Sprint(a.ALERT_FOR_MINIMUM_QUANTITY_BY_QUINTITY)
		target_balance := "\t," + fmt.Sprint(a.TARGET_BALANCE)
		if_the_target_balance_is_less_is_good := "\t," + fmt.Sprint(a.IF_THE_TARGET_BALANCE_IS_LESS_IS_GOOD)
		fmt.Fprintln(p, "{", is_low_level_account, is_credit, is_temporary, COST_FLOW_TYPE, account_name, notes,
			image, barcodes, account_number, account_levels, father_and_grandpa_accounts_name,
			alert_for_minimum_quantity_by_turnover_in_days, alert_for_minimum_quantity_by_quintity, target_balance, if_the_target_balance_is_less_is_good, "},")
	}
	p.Flush()
}

func FORMAT_SLICE_OF_SLICE_OF_STRING_TO_STRING(a [][]string) string {
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

func FORMAT_SLICE_OF_UINT_TO_STRING(a []uint) string {
	var str string
	for _, b := range a {
		str += fmt.Sprint(b) + ","
	}
	return "[]uint{" + str + "}"
}

func FORMAT_SLICE_OF_SLICE_OF_UINT_TO_STRING(a [][]uint) string {
	var str string
	for _, b := range a {
		str += "{"
		for _, c := range b {
			str += fmt.Sprint(c) + ","
		}
		str += "}\t,"
	}
	return "[][]uint{" + str + "}"
}

func FORMAT_STRING_SLICE_TO_STRING(a []string) string {
	var str string
	for _, b := range a {
		str += "\"" + b + "\","
	}
	return "[]string{" + str + "}"
}

func IS_IT_HIGH_THAN_BY_ORDER(account_number1, account_number2 []uint) bool {
	len_account_number1 := len(account_number1)
	len_account_number2 := len(account_number2)
	if len_account_number1 == len_account_number2 {
		return false
	}
	for index := 0; index < SMALLEST(len_account_number1, len_account_number2); index++ {
		if account_number1[index] < account_number2[index] {
			return true
		} else if account_number1[index] > account_number2[index] {
			return false
		}
	}
	return IS_SHORTER_THAN(account_number1, account_number2)
}

func INIT_ACCOUNT_NUMBERS_AND_FATHER_AND_GRANDPA_ACCOUNTS_NAME() {
	max_len := MAX_LEN_FOR_ACCOUNT_NUMBER()
	for indexa, a := range ACCOUNTS {
		ACCOUNTS[indexa].FATHER_AND_GRANDPA_ACCOUNTS_NAME = make([][]string, max_len)
		ACCOUNTS[indexa].ACCOUNT_NUMBER = make([][]uint, max_len)
		for indexb, b := range a.ACCOUNT_NUMBER {
			ACCOUNTS[indexa].ACCOUNT_NUMBER[indexb] = b
		}
	}
}

func MAX_LEN_FOR_ACCOUNT_NUMBER() int {
	var max_len int
	for _, a := range ACCOUNTS {
		length := len(a.ACCOUNT_NUMBER)
		if length > max_len {
			max_len = length
		}
	}
	return max_len
}

func SET_ACCOUNT_LEVELS() {
	for indexa, a := range ACCOUNTS {
		ACCOUNTS[indexa].ACCOUNT_LEVELS = []uint{}
		for _, b := range a.ACCOUNT_NUMBER {
			ACCOUNTS[indexa].ACCOUNT_LEVELS = append(ACCOUNTS[indexa].ACCOUNT_LEVELS, uint(len(b)))
		}
	}
}

func SET_FATHER_AND_GRANDPA_ACCOUNTS_NAME() {
	l := len(ACCOUNTS[0].ACCOUNT_NUMBER)
	for i := 0; i < l; i++ {
		for indexa, a := range ACCOUNTS {
			if len(a.ACCOUNT_NUMBER[i]) > 1 {
				for _, b := range ACCOUNTS {
					if len(b.ACCOUNT_NUMBER[i]) > 0 {
						if IS_IT_SUB_ACCOUNT_USING_NUMBER(b.ACCOUNT_NUMBER[i], a.ACCOUNT_NUMBER[i]) {
							ACCOUNTS[indexa].FATHER_AND_GRANDPA_ACCOUNTS_NAME[i] = append(ACCOUNTS[indexa].FATHER_AND_GRANDPA_ACCOUNTS_NAME[i], b.ACCOUNT_NAME)
						}
					}
				}
			}
		}
	}
}

func SET_HIGH_LEVEL_ACCOUNT_TO_PERMANENT() {
	for indexa, a := range ACCOUNTS {
		if !a.IS_LOW_LEVEL_ACCOUNT {
			ACCOUNTS[indexa].IS_TEMPORARY = false
		}
	}
}

func FIND_ALL_INVENTORY_FILES() []string {
	files, _ := ioutil.ReadDir(DB_INVENTORY)
	var inventory_files_name []string
	for _, f := range files {
		inventory_files_name = append(inventory_files_name, f.Name())
	}
	return inventory_files_name
}

func IS_ACCOUNT_NUMBERS_USED(account_number [][]uint) bool {
	l := len(account_number)
	for i := 0; i < l; i++ {
		for _, a := range ACCOUNTS {
			if reflect.DeepEqual(a.ACCOUNT_NUMBER[i], account_number[i]) {
				return true
			}
		}
	}
	return false
}

func PACK_THE_ACCOUNT_NUMBER_TO_MAX_LEN(account_number [][]uint) [][]uint {
	l := MAX_LEN_FOR_ACCOUNT_NUMBER()
	len_account_number := len(account_number)
	if len_account_number > l {
		l = len_account_number
	}
	return PACK(l, account_number)
}

func IS_BARCODES_USED(barcode []string) bool {
	for _, a := range ACCOUNTS {
		for _, b := range barcode {
			if IS_IN(b, a.BARCODE) {
				return true
			}
		}
	}
	return false
}

func UPDATE_INVENTORY_FILE_NAME(account_name, new_name string) {
	for _, inventory_account := range FIND_ALL_INVENTORY_FILES() {
		if inventory_account == account_name {
			os.Rename(DB_INVENTORY+account_name, DB_INVENTORY+new_name)
		}
	}
}

func SET_THE_ACCOUNTS() {
	INIT_ACCOUNT_NUMBERS_AND_FATHER_AND_GRANDPA_ACCOUNTS_NAME()
	SET_FATHER_AND_GRANDPA_ACCOUNTS_NAME()
	SORT_THE_ACCOUNTS_BY_ACCOUNT_NUMBER()
	SET_ACCOUNT_LEVELS()
	SET_HIGH_LEVEL_ACCOUNT_TO_PERMANENT()
	SET_COST_FLOW_TYPE()
}

func CHECK_IF_LOW_LEVEL_ACCOUNT_FOR_ALL() []error {
	var errors []error
	l := len(ACCOUNTS[0].ACCOUNT_NUMBER)
	for i := 1; i < l; i++ {
	big_loop:
		for indexa, a := range ACCOUNTS {
			if len(a.ACCOUNT_NUMBER[i]) != 0 {
				for _, b := range ACCOUNTS {
					if len(b.ACCOUNT_NUMBER[i]) != 0 {
						if IS_IT_SUB_ACCOUNT_USING_NUMBER(a.ACCOUNT_NUMBER[i], b.ACCOUNT_NUMBER[i]) {
							continue big_loop
						}
					}
				}
				if !ACCOUNTS[indexa].IS_LOW_LEVEL_ACCOUNT {
					errors = append(errors, fmt.Errorf("should be low level account in all account numbers %v", ACCOUNTS[indexa]))
				}
			}
		}
	}
	return errors
}

func CHECK_IF_THE_TREE_CONNECTED() []error {
	var errors []error
	l := len(ACCOUNTS[0].ACCOUNT_NUMBER)
	for i := 0; i < l; i++ {
	big_loop:
		for _, a := range ACCOUNTS {
			if len(a.ACCOUNT_NUMBER[i]) > 1 {
				for _, b := range ACCOUNTS {
					if IS_IT_THE_FATHER(b.ACCOUNT_NUMBER[i], a.ACCOUNT_NUMBER[i]) {
						continue big_loop
					}
				}
				errors = append(errors, fmt.Errorf("the account number %v for %v not conected to the tree", a.ACCOUNT_NUMBER[i], a))
			}
		}
	}
	return errors
}

func CHECK_THE_TREE() {
	ERRORS_MESSAGES = append(ERRORS_MESSAGES, CHECK_IF_LOW_LEVEL_ACCOUNT_FOR_ALL()...)
	ERRORS_MESSAGES = append(ERRORS_MESSAGES, CHECK_IF_THE_TREE_CONNECTED()...)
}

func IS_USED_IN_JOURNAL(account_name string) bool {
	for _, i := range DB_READ_JOURNAL(DB_JOURNAL) {
		if account_name == i.ACCOUNT_CREDIT || account_name == i.ACCOUNT_DEBIT {
			return true
		}
	}
	return false
}

// Account functions interface
func ADD_ACCOUNT(account ACCOUNT) error {
	_, _, err := ACCOUNT_STRUCT_FROM_NAME(account.ACCOUNT_NAME)
	if err == nil {
		return ERROR_ACCOUNT_NAME_IS_USED
	}
	account.ACCOUNT_NUMBER = PACK_THE_ACCOUNT_NUMBER_TO_MAX_LEN(account.ACCOUNT_NUMBER)
	if IS_ACCOUNT_NUMBERS_USED(account.ACCOUNT_NUMBER) {
		return ERROR_ACCOUNT_NUMBER_IS_USED
	}
	if IS_BARCODES_USED(account.BARCODE) {
		return ERROR_BARCODE_IS_USED
	}
	ACCOUNTS = append(ACCOUNTS, account)

	SET_THE_ACCOUNTS()
	DB_INSERT_INTO_ACCOUNTS()
	return nil
}

func EDIT_ACCOUNT(index int, account ACCOUNT, is_delete bool) {
	if !IS_USED_IN_JOURNAL(account.ACCOUNT_NAME) {
		if is_delete {
			ACCOUNTS = POPUP(ACCOUNTS, index)
			SET_THE_ACCOUNTS()
			DB_INSERT_INTO_ACCOUNTS()
			return
		}

		ACCOUNTS[index].IS_LOW_LEVEL_ACCOUNT = account.IS_LOW_LEVEL_ACCOUNT
		ACCOUNTS[index].IS_CREDIT = account.IS_CREDIT
		ACCOUNTS[index].IS_TEMPORARY = account.IS_TEMPORARY

		old_name := ACCOUNTS[index].ACCOUNT_NAME
		if old_name != account.ACCOUNT_NAME {
			_, _, err := ACCOUNT_STRUCT_FROM_NAME(account.ACCOUNT_NAME)
			if err == nil {
				ACCOUNTS[index].ACCOUNT_NAME = account.ACCOUNT_NAME
				DB_UPDATE_ACCOUNT_NAME_IN_JOURNAL(old_name, account.ACCOUNT_NAME)
				UPDATE_INVENTORY_FILE_NAME(old_name, account.ACCOUNT_NAME)
			}
		}
	}

	if !IS_ACCOUNT_NUMBERS_USED(account.ACCOUNT_NUMBER) {
		ACCOUNTS[index].ACCOUNT_NUMBER = account.ACCOUNT_NUMBER
	}

	if !IS_BARCODES_USED(account.BARCODE) {
		ACCOUNTS[index].BARCODE = account.BARCODE
	}

	ACCOUNTS[index].COST_FLOW_TYPE = account.COST_FLOW_TYPE
	ACCOUNTS[index].NOTES = account.NOTES
	ACCOUNTS[index].IMAGE = account.IMAGE
	ACCOUNTS[index].ALERT_FOR_MINIMUM_QUANTITY_BY_TURNOVER_IN_DAYS = account.ALERT_FOR_MINIMUM_QUANTITY_BY_TURNOVER_IN_DAYS
	ACCOUNTS[index].ALERT_FOR_MINIMUM_QUANTITY_BY_QUINTITY = account.ALERT_FOR_MINIMUM_QUANTITY_BY_QUINTITY
	ACCOUNTS[index].TARGET_BALANCE = account.TARGET_BALANCE
	ACCOUNTS[index].IF_THE_TARGET_BALANCE_IS_LESS_IS_GOOD = account.IF_THE_TARGET_BALANCE_IS_LESS_IS_GOOD

	SET_THE_ACCOUNTS()
	DB_INSERT_INTO_ACCOUNTS()
}
