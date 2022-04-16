package anti_accountants

import (
	"fmt"
	"os"
	"reflect"
	"text/tabwriter"
)

func ACCOUNT_STRUCT_FROM_BARCODE(barcode string) (ACCOUNT, int, error) {
	for indexa, a := range ACCOUNTS {
		if IS_IN(barcode, a.BARCODE) {
			return a, indexa, nil
		}
	}
	return ACCOUNT{}, 0, ERROR_NOT_LISTED
}

func ACCOUNT_STRUCT_FROM_NAME(account_name string) (ACCOUNT, int, error) {
	for indexa, a := range ACCOUNTS {
		if a.ACCOUNT_NAME == account_name {
			return a, indexa, nil
		}
	}
	return ACCOUNT{}, 0, ERROR_NOT_LISTED
}

func ADD_ACCOUNT(account ACCOUNT) error {
	account.ACCOUNT_NAME = FORMAT_THE_STRING(account.ACCOUNT_NAME)
	if account.ACCOUNT_NAME == "" {
		return ERROR_ACCOUNT_NAME_IS_EMPTY
	}
	_, _, err := ACCOUNT_STRUCT_FROM_NAME(account.ACCOUNT_NAME)
	if err == nil {
		return ERROR_ACCOUNT_NAME_IS_USED
	}
	if IS_BARCODES_USED(account.BARCODE) {
		return ERROR_BARCODE_IS_USED
	}

	ACCOUNTS = append(ACCOUNTS, account)
	SET_THE_ACCOUNTS()
	DB_INSERT_INTO_ACCOUNTS()
	return nil
}

func CHECK_IF_ACCOUNT_NUMBER_DUPLICATED() []error {
	var errors []error
	max_len := MAX_LEN_FOR_ACCOUNT_NUMBER()
	for i := 0; i < max_len; i++ {
	big_loop:
		for indexa, a := range ACCOUNTS {
			if len(a.ACCOUNT_NUMBER[i]) > 0 {
				for indexb, b := range ACCOUNTS {
					if indexa != indexb && reflect.DeepEqual(a.ACCOUNT_NUMBER[i], b.ACCOUNT_NUMBER[i]) {
						errors = append(errors, fmt.Errorf("the account number %v for %v is duplicated", a.ACCOUNT_NUMBER[i], a))
						continue big_loop
					}
				}
			}
		}
	}
	errors, _ = RETURN_SET_AND_DUPLICATES_SLICES(errors)
	return errors
}

func CHECK_IF_LOW_LEVEL_ACCOUNT_FOR_ALL() []error {
	var errors []error
	max_len := MAX_LEN_FOR_ACCOUNT_NUMBER()
	for i := 1; i < max_len; i++ {
	big_loop:
		for indexa, a := range ACCOUNTS {
			if len(a.ACCOUNT_NUMBER[i]) > 0 {
				for _, b := range ACCOUNTS {
					if len(b.ACCOUNT_NUMBER[i]) > 0 {
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
	max_len := MAX_LEN_FOR_ACCOUNT_NUMBER()
	for i := 0; i < max_len; i++ {
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

func CHECK_THE_TREE() []error {
	var errors_messages []error
	errors_messages = append(errors_messages, CHECK_IF_LOW_LEVEL_ACCOUNT_FOR_ALL()...)
	errors_messages = append(errors_messages, CHECK_IF_ACCOUNT_NUMBER_DUPLICATED()...)
	errors_messages = append(errors_messages, CHECK_IF_THE_TREE_CONNECTED()...)
	return errors_messages
}

func EDIT_ACCOUNT(is_delete bool, index int, account ACCOUNT) {
	if !IS_USED_IN_JOURNAL(account.ACCOUNT_NAME) {
		if is_delete {
			ACCOUNTS = REMOVE(ACCOUNTS, index)
			SET_THE_ACCOUNTS()
			DB_INSERT_INTO_ACCOUNTS()
			return
		}

		ACCOUNTS[index].IS_LOW_LEVEL_ACCOUNT = account.IS_LOW_LEVEL_ACCOUNT
		ACCOUNTS[index].IS_CREDIT = account.IS_CREDIT

		// if the account not used in journal then the account is not used in inventory then
		// i will search for the account name in accounts database if it is not used then i can chenge the name
		new_name := FORMAT_THE_STRING(account.ACCOUNT_NAME)
		old_name := ACCOUNTS[index].ACCOUNT_NAME
		if old_name != new_name && new_name != "" {
			_, _, err := ACCOUNT_STRUCT_FROM_NAME(new_name)
			if err != nil {
				DB_UPDATE_ACCOUNT_NAME_IN_JOURNAL(old_name, new_name)
				DB_UPDATE_ACCOUNT_NAME_IN_INVENTORY(old_name, new_name)
				ACCOUNTS[index].ACCOUNT_NAME = new_name
			}
		}
	}

	if !IS_BARCODES_USED(account.BARCODE) {
		ACCOUNTS[index].BARCODE = account.BARCODE
	}

	ACCOUNTS[index].IS_TEMPORARY = account.IS_TEMPORARY
	ACCOUNTS[index].COST_FLOW_TYPE = account.COST_FLOW_TYPE
	ACCOUNTS[index].NOTES = account.NOTES
	ACCOUNTS[index].IMAGE = account.IMAGE
	ACCOUNTS[index].ACCOUNT_NUMBER = account.ACCOUNT_NUMBER
	ACCOUNTS[index].ALERT_FOR_MINIMUM_QUANTITY_BY_TURNOVER_IN_DAYS = account.ALERT_FOR_MINIMUM_QUANTITY_BY_TURNOVER_IN_DAYS
	ACCOUNTS[index].ALERT_FOR_MINIMUM_QUANTITY_BY_QUINTITY = account.ALERT_FOR_MINIMUM_QUANTITY_BY_QUINTITY
	ACCOUNTS[index].TARGET_BALANCE = account.TARGET_BALANCE
	ACCOUNTS[index].IF_THE_TARGET_BALANCE_IS_LESS_IS_GOOD = account.IF_THE_TARGET_BALANCE_IS_LESS_IS_GOOD

	SET_THE_ACCOUNTS()
	DB_INSERT_INTO_ACCOUNTS()
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

func FORMAT_SLICE_OF_UINT_TO_STRING(a []uint) string {
	var str string
	for _, b := range a {
		str += fmt.Sprint(b) + ","
	}
	return "[]uint{" + str + "}"
}

func FORMAT_STRING_SLICE_TO_STRING(a []string) string {
	var str string
	for _, b := range a {
		str += "\"" + b + "\","
	}
	return "[]string{" + str + "}"
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

func IS_IT_HIGH_THAN_BY_ORDER(account_number1, account_number2 []uint) bool {
	l1 := len(account_number1)
	l2 := len(account_number2)
	for index := 0; index < SMALLEST(l1, l2); index++ {
		if account_number1[index] < account_number2[index] {
			return true
		} else if account_number1[index] > account_number2[index] {
			return false
		}
	}
	return l2 > l1
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

func IS_IT_THE_FATHER(higher_level_account_number, lower_level_account_number []uint) bool {
	if !IS_IT_POSSIBLE_TO_BE_SUB_ACCOUNT(higher_level_account_number, lower_level_account_number) {
		return false
	}
	return reflect.DeepEqual(higher_level_account_number, CUT_THE_SLICE(lower_level_account_number, 1))
}

func IS_USED_IN_JOURNAL(account_name string) bool {
	_, journal := DB_READ[JOURNAL_TAG](DB_JOURNAL)
	for _, i := range journal {
		if account_name == i.ACCOUNT_CREDIT || account_name == i.ACCOUNT_DEBIT {
			return true
		}
	}
	return false
}

func MAX_LEN_FOR_ACCOUNT_NUMBER() int {
	var max_len int
	for _, a := range ACCOUNTS {
		var length int
		for _, b := range a.ACCOUNT_NUMBER {
			if len(b) > 0 {
				length++
			}
		}
		if length > max_len {
			max_len = length
		}
	}
	return max_len
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

func SET_FATHER_AND_GRANDPA_ACCOUNTS_NAME() {
	max_len := MAX_LEN_FOR_ACCOUNT_NUMBER()
	for i := 0; i < max_len; i++ {
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

func SET_THE_ACCOUNTS() {
	max_len := MAX_LEN_FOR_ACCOUNT_NUMBER()
	for k1, v1 := range ACCOUNTS {
		// init the slices
		ACCOUNTS[k1].FATHER_AND_GRANDPA_ACCOUNTS_NAME = make([][]string, max_len)
		ACCOUNTS[k1].ACCOUNT_NUMBER = make([][]uint, max_len)
		ACCOUNTS[k1].ACCOUNT_LEVELS = make([]uint, max_len)
		for k2, v2 := range v1.ACCOUNT_NUMBER {
			if k2 < max_len {
				ACCOUNTS[k1].ACCOUNT_NUMBER[k2] = v2
				ACCOUNTS[k1].ACCOUNT_LEVELS[k2] = uint(len(v2))
			}
		}

		// set high level account to permanent
		// set cost flow type . the cost flow should be used for every low level account
		if !v1.IS_LOW_LEVEL_ACCOUNT {
			ACCOUNTS[k1].IS_TEMPORARY = false
			ACCOUNTS[k1].COST_FLOW_TYPE = ""
		} else if !IS_IN(v1.COST_FLOW_TYPE, COST_FLOW_TYPE) {
			ACCOUNTS[k1].COST_FLOW_TYPE = FIFO
		}
	}
	SET_FATHER_AND_GRANDPA_ACCOUNTS_NAME()
	SORT_THE_ACCOUNTS_BY_ACCOUNT_NUMBER()
}

func SORT_THE_ACCOUNTS_BY_ACCOUNT_NUMBER() {
	for indexa := range ACCOUNTS {
		for indexb := range ACCOUNTS {
			if indexa < indexb && !IS_IT_HIGH_THAN_BY_ORDER(ACCOUNTS[indexa].ACCOUNT_NUMBER[INDEX_OF_ACCOUNT_NUMBER], ACCOUNTS[indexb].ACCOUNT_NUMBER[INDEX_OF_ACCOUNT_NUMBER]) {
				SWAP(ACCOUNTS, indexa, indexb)
			}
		}
	}
}
