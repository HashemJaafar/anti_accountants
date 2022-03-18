package anti_accountants

import (
	"errors"
	"math"
	"time"
)

func set_date_end_to_zero_if_smaller_than_date_start(date_start, date_end time.Time) time.Time {
	if !date_end.IsZero() {
		if !date_start.Before(date_end) {
			return time.Time{}
		}
	}
	return date_end
}

func set_adjusting_method(entry_expair time.Time, adjusting_method string, entries []VALUE_PRICE_QUANTITY_ACCOUNT_BARCODE) string {
	if !is_in(adjusting_method, adjusting_methods) {
		return ""
	}
	if entry_expair.IsZero() {
		return ""
	}
	is_in_depreciation_methods := is_in(adjusting_method, depreciation_methods)
	for _, entry := range entries {
		if account_struct_from_name(entry.ACCOUNT_NAME).COST_FLOW_TYPE != "" && is_in_depreciation_methods {
			return ""
		}
	}
	return adjusting_method
}

func group_by_account_and_barcode(entries []VALUE_PRICE_QUANTITY_ACCOUNT_BARCODE) []VALUE_PRICE_QUANTITY_ACCOUNT_BARCODE {
	type account_barcode struct {
		account, barcode string
	}
	g := map[account_barcode]*VALUE_PRICE_QUANTITY_ACCOUNT_BARCODE{}
	for _, v := range entries {
		key := account_barcode{v.ACCOUNT_NAME, v.BARCODE}
		sums := g[key]
		if sums == nil {
			sums = &VALUE_PRICE_QUANTITY_ACCOUNT_BARCODE{}
			g[key] = sums
		}
		sums.VALUE += v.VALUE
		sums.QUANTITY += v.QUANTITY
	}
	entries = []VALUE_PRICE_QUANTITY_ACCOUNT_BARCODE{}
	for key, v := range g {
		entries = append(entries, VALUE_PRICE_QUANTITY_ACCOUNT_BARCODE{
			VALUE:        v.VALUE,
			PRICE:        v.VALUE / v.QUANTITY,
			QUANTITY:     v.QUANTITY,
			ACCOUNT_NAME: key.account,
			BARCODE:      key.barcode,
		})
	}
	return entries
}

func remove_zero_values(entries []VALUE_PRICE_QUANTITY_ACCOUNT_BARCODE) {
	var index int
	for index < len(entries) {
		if entries[index].VALUE == 0 || entries[index].QUANTITY == 0 {
			entries = append(entries[:index], entries[index+1:]...)
		} else {
			index++
		}
	}
}

func find_cost(entries []VALUE_PRICE_QUANTITY_ACCOUNT_BARCODE) {
	for index, entry := range entries {
		costs := cost_flow(entry.ACCOUNT_NAME, entry.QUANTITY, false)
		if costs != 0 {
			entries[index].VALUE = -costs
			entries[index].PRICE = -costs / entry.QUANTITY
		}
	}
}

// func convert_to_simple_entry(debit_entries, credit_entries []VALUE_PRICE_QUANTITY_ACCOUNT_BARCODE) [][]VALUE_PRICE_QUANTITY_ACCOUNT_BARCODE {
// 	simple_entries := [][]VALUE_PRICE_QUANTITY_ACCOUNT_BARCODE{}
// 	for _, debit_entry := range debit_entries {
// 		for _, credit_entry := range credit_entries {
// 			simple_entries = append(simple_entries, []VALUE_PRICE_QUANTITY_ACCOUNT_BARCODE{debit_entry, credit_entry})
// 		}
// 	}
// 	for _, a := range simple_entries {
// 		switch math.Abs(a[0].VALUE) >= math.Abs(a[1].VALUE) {
// 		case true:
// 			sign := a[0].VALUE / a[1].VALUE
// 			price := a[0].VALUE / a[0].QUANTITY
// 			a[0].VALUE = a[1].VALUE * sign / math.Abs(sign)
// 			a[0].QUANTITY = a[0].VALUE / price
// 		case false:
// 			sign := a[0].VALUE / a[1].VALUE
// 			price := a[1].VALUE / a[1].QUANTITY
// 			a[1].VALUE = a[0].VALUE * sign / math.Abs(sign)
// 			a[1].QUANTITY = a[1].VALUE / price
// 		}
// 	}
// 	return simple_entries
// }

// func can_the_account_be_negative(entries []VALUE_PRICE_QUANTITY_ACCOUNT_BARCODE) {
// 	for _, entry := range entries {
// 		if !(is_it_sub_account_using_name(PRIMARY_ACCOUNTS_NAMES.EQUITY, entry.ACCOUNT) && is_credit(entry.ACCOUNT)) {
// 			account_balance := account_balance(entry.ACCOUNT)
// 			if account_balance+entry.VALUE < 0 {
// 				error_make_nagtive_balance(entry, account_balance)
// 			}
// 		}
// 	}
// }

func find_account_from_barcode(entries []VALUE_PRICE_QUANTITY_ACCOUNT_BARCODE) {
	for indexa, a := range entries {
		entries[indexa].ACCOUNT_NAME = account_struct_from_barcode(a.BARCODE).ACCOUNT_NAME
	}
}

// func insert_to_JOURNAL_TAG(entries []VALUE_PRICE_QUANTITY_ACCOUNT_BARCODE, date time.Time, entry_expair time.Time, description string, name string, employee_name string) []JOURNAL_TAG {
// 	var array_to_insert []JOURNAL_TAG
// 	for _, entry := range entries {
// 		array_to_insert = append(array_to_insert, JOURNAL_TAG{
// 			DATE:          date.String(),
// 			ENTRY_NUMBER:  0,
// 			ACCOUNT:       entry.ACCOUNT,
// 			VALUE:         entry.VALUE,
// 			PRICE:         entry.VALUE / entry.QUANTITY,
// 			QUANTITY:      entry.QUANTITY,
// 			BARCODE:       entry.BARCODE,
// 			ENTRY_EXPAIR:  entry_expair.String(),
// 			DESCRIPTION:   description,
// 			NAME:          name,
// 			EMPLOYEE_NAME: employee_name,
// 			ENTRY_DATE:    NOW.String(),
// 			REVERSE:       false,
// 		})
// 	}
// 	return array_to_insert
// }

// func check_if_the_price_is_negative(entries []VALUE_PRICE_QUANTITY_ACCOUNT_BARCODE) {
// 	for index, entry := range entries {
// 		entries[index].PRICE = entry.VALUE / entry.QUANTITY
// 		if entries[index].PRICE < 0 {
// 			error_the_price_should_be_positive(entries[index])
// 		}
// 	}
// }

func cost_flow(account string, quantity float64, insert bool) float64 {
	if quantity > 0 {
		return 0
	}
	is_ascending, err := is_ascending(account)
	if err != nil {
		return 0
	}
	inventory := db_read_inventory()
	sort_by_time(inventory, is_ascending)
	quantity = math.Abs(quantity)
	quantity_count := quantity
	var costs float64
	for _, item := range inventory {
		if item.QUANTITY > quantity_count {
			costs += item.PRICE * quantity_count
			if insert {
				// DB.Exec("update inventory set quantity=quantity-? where account=? and price=? and quantity=? and barcode=? order by date "+order_by_date_asc_or_desc+" limit 1", quantity_count, account, item.PRICE, item.QUANTITY, barcode)
			}
			quantity_count = 0
			break
		}
		if item.QUANTITY <= quantity_count {
			costs += item.PRICE * item.QUANTITY
			if insert {
				// DB.Exec("delete from inventory where account=? and price=? and quantity=? and barcode=? order by date "+order_by_date_asc_or_desc+" limit 1", account, item.PRICE, item.QUANTITY, barcode)
			}
			quantity_count -= item.QUANTITY
		}
	}
	if quantity_count != 0 {
		error_the_order_out_of_stock(quantity, quantity_count, account)
	}
	return costs
}

func is_ascending(account string) (bool, error) {
	switch account_struct_from_name(account).COST_FLOW_TYPE {
	case "lifo":
		return false, nil
	case "fifo":
		return true, nil
	case "wma":
		weighted_average(account)
		return true, nil
	}
	return false, errors.New("is not inventory account")
}

// func insert_to_database(array_of_journal_tag []JOURNAL_TAG, db_insert_into_journal, insert_into_inventory bool) {
// 	insert_entry_number(array_of_journal_tag)
// 	if db_insert_into_journal {
// 		db_insert_into_journal_func(array_of_journal_tag)
// 	}
// 	if insert_into_inventory {
// 		db_insert_into_inventory(array_of_journal_tag)
// 	}
// }

// func insert_entry_number(array_of_journal_tag []JOURNAL_TAG) {
// 	entry_number := float64(entry_number())
// 	for indexa := range array_of_journal_tag {
// 		array_of_journal_tag[indexa].ENTRY_NUMBER = int(entry_number)
// 		entry_number += 0.5
// 	}
// }

func calculate_and_insert_value_price_quantity(entries []VALUE_PRICE_QUANTITY_ACCOUNT_BARCODE) {
	for index, entry := range entries {
		m := map[string]float64{}
		insert_if_not_zero(m, "VALUE", entry.VALUE)
		insert_if_not_zero(m, "PRICE", entry.PRICE)
		insert_if_not_zero(m, "QUANTITY", entry.QUANTITY)
		EQUATIONS_SOLVER(false, false, m, [][]string{{"VALUE", "PRICE", "*", "QUANTITY"}})
		entries[index].VALUE = m["VALUE"]
		entries[index].PRICE = m["PRICE"]
		entries[index].QUANTITY = m["QUANTITY"]
	}
}

func insert_if_not_zero(m map[string]float64, str string, number float64) {
	if number != 0 {
		m[str] = number
	}
}

func remove_high_level_account(entries []VALUE_PRICE_QUANTITY_ACCOUNT_BARCODE) {
	var indexa int
	for indexa < len(entries) {
		if !account_struct_from_name(entries[indexa].ACCOUNT_NAME).IS_LOW_LEVEL_ACCOUNT {
			entries = append(entries[:indexa], entries[indexa+1:]...)
		} else {
			indexa++
		}
	}
}

func JOURNAL_ENTRY(
	entries []VALUE_PRICE_QUANTITY_ACCOUNT_BARCODE,
	insert, auto_completion, invoice_discount bool,
	date_start, date_end time.Time,
	adjusting_method, description, name, employee_name string,
	array_day_start_end []DAY_START_END) []JOURNAL_TAG {

	find_account_from_barcode(entries)
	date_end = set_date_end_to_zero_if_smaller_than_date_start(date_start, date_end)
	adjusting_method = set_adjusting_method(date_end, adjusting_method, entries)
	calculate_and_insert_value_price_quantity(entries)

	entries = group_by_account_and_barcode(entries)
	remove_zero_values(entries)

	find_cost(entries)

	// if auto_completion {
	// 	entries = auto_completion_the_entry(entries)
	// }
	// if invoice_discount {
	// 	entries = auto_completion_the_invoice_discount(entries)
	// }

	entries = group_by_account_and_barcode(entries)
	remove_zero_values(entries)

	remove_high_level_account(entries)
	// can_the_account_be_negative(entries)

	// check_debit_equal_credit(entries)
	// debit_entries, credit_entries := separate_debit_from_credit(entries)
	// check_one_debit_or_one_credit(debit_entries, credit_entries)
	// simple_entries := convert_to_simple_entry(debit_entries, credit_entries)

	// var all_array_to_insert []JOURNAL_TAG
	// for _, simple_entry := range simple_entries {

	// 	check_if_the_price_is_negative(simple_entry)
	// 	array_to_insert := insert_to_JOURNAL_TAG(simple_entry, date, entry_expair, description, name, employee_name)

	// 	if is_in(adjusting_method, depreciation_methods) {
	// 		array_day_start_end = initialize_array_day_start_end(array_day_start_end)
	// 		check_array_day_start_end(array_day_start_end)
	// 		array_start_end_minutes := create_array_start_end_minutes(entry_expair, date, array_day_start_end)
	// 		total_minutes := total_minutes(array_start_end_minutes)
	// 		adjusted_array_to_insert := adjust_the_array(array_to_insert, array_start_end_minutes, total_minutes, adjusting_method, description, name, employee_name)
	// 		adjusted_array_to_insert = transpose(adjusted_array_to_insert)
	// 		array_to_insert = unpack_the_array(adjusted_array_to_insert)
	// 	}

	// 	all_array_to_insert = append(all_array_to_insert, array_to_insert...)
	// }

	// insert_to_database(all_array_to_insert, insert, insert)
	return []JOURNAL_TAG{}
}
