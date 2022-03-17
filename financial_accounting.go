package anti_accountants

// import "log"

// func check_debit_equal_credit(entries []VALUE_PRICE_QUANTITY_ACCOUNT_BARCODE) {
// 	var zero float64
// 	for _, entry := range entries {
// 		switch account_struct_from_name(entry.ACCOUNT).IS_CREDIT {
// 		case false:
// 			zero += entry.VALUE
// 		case true:
// 			zero -= entry.VALUE
// 		}
// 	}
// 	if zero != 0 {
// 		error_debit_not_equal_credit(zero, entries)
// 	}
// }

// func separate_debit_from_credit(entries []VALUE_PRICE_QUANTITY_ACCOUNT_BARCODE) ([]VALUE_PRICE_QUANTITY_ACCOUNT_BARCODE, []VALUE_PRICE_QUANTITY_ACCOUNT_BARCODE) {
// 	var debit_entries, credit_entries []VALUE_PRICE_QUANTITY_ACCOUNT_BARCODE
// 	for _, entry := range entries {
// 		switch account_struct_from_name(entry.ACCOUNT).IS_CREDIT {
// 		case false:
// 			if entry.VALUE > 0 {
// 				debit_entries = append(debit_entries, entry)
// 			} else if entry.VALUE < 0 {
// 				credit_entries = append(credit_entries, entry)
// 			} else {
// 				log.Panic("value is zero for entry: ", entry)
// 			}
// 		case true:
// 			if entry.VALUE < 0 {
// 				debit_entries = append(debit_entries, entry)
// 			} else if entry.VALUE > 0 {
// 				credit_entries = append(credit_entries, entry)
// 			} else {
// 				log.Panic("value is zero for entry: ", entry)
// 			}
// 		}
// 	}
// 	return debit_entries, credit_entries
// }

// func check_one_debit_or_one_credit(debit_entries, credit_entries []VALUE_PRICE_QUANTITY_ACCOUNT_BARCODE) {
// 	if (len(debit_entries) != 1) && (len(credit_entries) != 1) {
// 		error_one_credit___one_debit("or", debit_entries, credit_entries)
// 	}
// }

// func check_one_debit_and_one_credit(debit_entries, credit_entries []VALUE_PRICE_QUANTITY_ACCOUNT_BARCODE) {
// 	if !((len(debit_entries) == 1) && (len(credit_entries) == 1)) {
// 		error_one_credit___one_debit("and", debit_entries, credit_entries)
// 	}
// }

// func INVOICE(array_of_journal_tag []JOURNAL_TAG) []INVOICE_STRUCT {
// 	m := map[string]*INVOICE_STRUCT{}
// 	for _, entry := range array_of_journal_tag {
// 		var key string
// 		switch {
// 		case is_it_sub_account_using_name(PRIMARY_ACCOUNTS_NAMES.ASSETS, entry.ACCOUNT) && !is_credit(entry.ACCOUNT) && !is_in(entry.ACCOUNT, inventory) && entry.VALUE > 0:
// 			key = "total"
// 		case is_it_sub_account_using_name(PRIMARY_ACCOUNTS_NAMES.DISCOUNTS, entry.ACCOUNT) && !is_credit(entry.ACCOUNT):
// 			key = "total discounts"
// 		case is_it_sub_account_using_name(PRIMARY_ACCOUNTS_NAMES.SALES, entry.ACCOUNT) && is_credit(entry.ACCOUNT):
// 			key = entry.ACCOUNT
// 		default:
// 			continue
// 		}
// 		sums := m[key]
// 		if sums == nil {
// 			sums = &INVOICE_STRUCT{}
// 			m[key] = sums
// 		}
// 		sums.VALUE += entry.VALUE
// 		sums.QUANTITY += entry.QUANTITY
// 	}
// 	invoice := []INVOICE_STRUCT{}
// 	for k, v := range m {
// 		invoice = append(invoice, INVOICE_STRUCT{k, v.VALUE, v.VALUE / v.QUANTITY, v.QUANTITY})
// 	}
// 	return invoice
// }

func INITIALIZE() {
	ACCOUNTS = db_read_accounts()
	init_account_numbers_and_father_and_grandpa_accounts_name()
	remove_duplicate_accounts_name()
	remove_duplicate_accounts_barcode()
	remove_duplicate_accounts_number()
	check_if_the_tree_connected()
	set_low_level_accounts()
	check_if_low_level_account_for_all()
	set_high_level_account_to_permanent()
	set_account_levels()
	set_father_and_grandpa_accounts_name()
	set_cost_flow_type()
	sort_the_accounts_by_account_number()
	db_insert_into_accounts()
	print_formated_accounts()
	check_if_journal_accounts_is_listed_in_accounts()
	check_if_inventory_accounts_still_used_as_inventory()
}
