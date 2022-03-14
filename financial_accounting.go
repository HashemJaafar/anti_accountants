package anti_accountants

// import "log"

// func check_debit_equal_credit(entries []ACCOUNT_VALUE_PRICE_QUANTITY_BARCODE) {
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

// func separate_debit_from_credit(entries []ACCOUNT_VALUE_PRICE_QUANTITY_BARCODE) ([]ACCOUNT_VALUE_PRICE_QUANTITY_BARCODE, []ACCOUNT_VALUE_PRICE_QUANTITY_BARCODE) {
// 	var debit_entries, credit_entries []ACCOUNT_VALUE_PRICE_QUANTITY_BARCODE
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

// func check_one_debit_or_one_credit(debit_entries, credit_entries []ACCOUNT_VALUE_PRICE_QUANTITY_BARCODE) {
// 	if (len(debit_entries) != 1) && (len(credit_entries) != 1) {
// 		error_one_credit___one_debit("or", debit_entries, credit_entries)
// 	}
// }

// func check_one_debit_and_one_credit(debit_entries, credit_entries []ACCOUNT_VALUE_PRICE_QUANTITY_BARCODE) {
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

func INITIALIZE(driverName, dataSourceName, database_name string) {
	// open_and_create_database(driverName, dataSourceName, database_name)
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
	inventory_accounts()
	// check_accounts("account", "inventory", " is not have fifo lifo wma on cost_flow_type field", inventory)
	sort_the_accounts_by_account_number()

	// entry_number := entry_number()
	// var array_to_insert []JOURNAL_TAG
	// expair_expenses := JOURNAL_TAG{NOW.String(), entry_number, expair_expenses, 0, 0, 0, "", time.Time{}.String(), "to record the expiry of the goods automatically", "", "", NOW.String(), false}
	// expair_goods, _ := DB.Query("select account,price*quantity*-1,price,quantity*-1,barcode from inventory where entry_expair<? and entry_expair!='0001-01-01 00:00:00 +0000 UTC'", NOW.String())
	// for expair_goods.Next() {
	// 	tag := expair_expenses
	// 	expair_goods.Scan(&tag.ACCOUNT, &tag.value, &tag.price, &tag.quantity, &tag.barcode)
	// 	expair_expenses.value -= tag.value
	// 	expair_expenses.quantity -= tag.quantity
	// 	array_to_insert = append(array_to_insert, tag)
	// }
	// expair_expenses.price = expair_expenses.value / expair_expenses.quantity
	// array_to_insert = append(array_to_insert, expair_expenses)
	// insert_to_database(array_to_insert, true, false)
	// DB.Exec("delete from inventory where entry_expair<? and entry_expair!='0001-01-01 00:00:00 +0000 UTC'", NOW.String())
	// DB.Exec("delete from inventory where quantity=0")

	// check_debit_equal_credit_and_check_one_debit_and_one_credit_in_the_journal(JOURNAL_ORDERED_BY_DATE_ENTRY_NUMBER())
}

// func check_debit_equal_credit_and_check_one_debit_and_one_credit_in_the_journal(JOURNAL_ORDERED_BY_DATE_ENTRY_NUMBER []JOURNAL_TAG) {
// 	var double_entry []ACCOUNT_VALUE_PRICE_QUANTITY_BARCODE
// 	previous_entry_number := 1
// 	for _, entry := range JOURNAL_ORDERED_BY_DATE_ENTRY_NUMBER {
// 		if previous_entry_number != entry.ENTRY_NUMBER {
// 			delete_not_double_entry(double_entry, previous_entry_number)
// 			if len(double_entry) == 2 {
// 				check_debit_equal_credit(double_entry)
// 				debit_entries, credit_entries := separate_debit_from_credit(double_entry)
// 				check_one_debit_and_one_credit(debit_entries, credit_entries)
// 			}
// 			double_entry = []ACCOUNT_VALUE_PRICE_QUANTITY_BARCODE{}
// 		}
// 		double_entry = append(double_entry, ACCOUNT_VALUE_PRICE_QUANTITY_BARCODE{
// 			ACCOUNT:  entry.ACCOUNT,
// 			VALUE:    entry.VALUE,
// 			PRICE:    entry.PRICE,
// 			QUANTITY: entry.QUANTITY,
// 			BARCODE:  entry.BARCODE,
// 		})
// 		previous_entry_number = entry.ENTRY_NUMBER
// 	}
// 	delete_not_double_entry(double_entry, previous_entry_number)
// 	check_debit_equal_credit(double_entry)
// 	debit_entries, credit_entries := separate_debit_from_credit(double_entry)
// 	check_one_debit_and_one_credit(debit_entries, credit_entries)
// }
