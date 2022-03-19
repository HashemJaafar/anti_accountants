package anti_accountants

import (
	"errors"
	"fmt"
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
		account_struct, _, _ := account_struct_from_name(entry.ACCOUNT_NAME)
		if account_struct.COST_FLOW_TYPE != "" && is_in_depreciation_methods {
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
			popup(entries, index)
		} else {
			index++
		}
	}
}

func find_cost(entries []VALUE_PRICE_QUANTITY_ACCOUNT_BARCODE) {
	for index, entry := range entries {
		costs, err := cost_flow(entry.ACCOUNT_NAME, entry.QUANTITY, false)
		if err == nil {
			entries[index].VALUE = -costs
			entries[index].PRICE = -costs / entry.QUANTITY
		}
	}
}

func convert_to_simple_entry(debit_entries, credit_entries []VALUE_PRICE_QUANTITY_ACCOUNT_BARCODE) []JOURNAL_TAG {
	var simple_entries []JOURNAL_TAG
	for _, debit_entry := range debit_entries {
		for _, credit_entry := range credit_entries {
			simple_entries = append(simple_entries, JOURNAL_TAG{
				VALUE:           math.Abs(smallest(debit_entry.VALUE, credit_entry.VALUE)),
				PRICE_DEBIT:     debit_entry.PRICE,
				PRICE_CREDIT:    credit_entry.PRICE,
				QUANTITY_DEBIT:  debit_entry.QUANTITY,
				QUANTITY_CREDIT: credit_entry.QUANTITY,
				ACCOUNT_DEBIT:   debit_entry.ACCOUNT_NAME,
				ACCOUNT_CREDIT:  credit_entry.ACCOUNT_NAME,
			})
		}
	}
	return simple_entries
}

func increase_the_value_to_make_the_new_balance_for_the_account_positive(entries []VALUE_PRICE_QUANTITY_ACCOUNT_BARCODE) {
	for indexa, a := range entries {
		account_balance := account_balance(a.ACCOUNT_NAME)
		new_balance := account_balance + a.VALUE
		if new_balance < 0 {
			entries[indexa].VALUE -= new_balance
			entries[indexa].QUANTITY = entries[indexa].VALUE / entries[indexa].PRICE
		}
	}
}

func find_account_from_barcode(entries []VALUE_PRICE_QUANTITY_ACCOUNT_BARCODE) {
	for indexa, a := range entries {
		account_struct, _, err := account_struct_from_barcode(a.BARCODE)
		if err == nil {
			entries[indexa].ACCOUNT_NAME = account_struct.ACCOUNT_NAME
		}
	}
}

func insert_to_journal_tag(entries []JOURNAL_TAG, date_start, date_end time.Time, notes, name, name_employee string) {
	for indexa := range entries {
		entries[indexa].DATE_START = date_start
		entries[indexa].DATE_END = date_end
		entries[indexa].NOTES = notes
		entries[indexa].NAME = name
		entries[indexa].NAME_EMPLOYEE = name_employee
	}
}

func set_price_to_positive(entries []VALUE_PRICE_QUANTITY_ACCOUNT_BARCODE) {
	for index, entry := range entries {
		entries[index].PRICE = entry.VALUE / entry.QUANTITY
		if entries[index].PRICE < 0 {
			entries[index].VALUE = math.Abs(entries[index].VALUE)
			entries[index].QUANTITY = math.Abs(entries[index].QUANTITY)
		}
	}
}

func cost_flow(account string, quantity float64, insert bool) (float64, error) {
	if quantity > 0 {
		return 0, error_should_be_negative
	}
	is_ascending, err := is_ascending(account)
	if err != nil {
		return 0, error_not_inventory_account
	}
	inventory := db_read_inventory(account)
	sort_by_time_inventory(inventory, is_ascending)
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
	}
	return costs, nil
}

func is_ascending(account string) (bool, error) {
	account_struct, _, _ := account_struct_from_name(account)
	switch account_struct.COST_FLOW_TYPE {
	case LIFO:
		return false, nil
	case FIFO:
		return true, nil
	case WMA:
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
		account_struct, _, _ := account_struct_from_name(entries[indexa].ACCOUNT_NAME)
		if !account_struct.IS_LOW_LEVEL_ACCOUNT {
			popup(entries, indexa)
		} else {
			indexa++
		}
	}
}

func remove_the_accounts_not_in_accounts_list(entries []VALUE_PRICE_QUANTITY_ACCOUNT_BARCODE) {
	var indexa int
	for indexa < len(entries) {
		_, _, err := account_struct_from_name(entries[indexa].ACCOUNT_NAME)
		if err != nil {
			popup(entries, indexa)
		} else {
			indexa++
		}
	}
}

func check_debit_equal_credit(entries []VALUE_PRICE_QUANTITY_ACCOUNT_BARCODE) ([]VALUE_PRICE_QUANTITY_ACCOUNT_BARCODE, []VALUE_PRICE_QUANTITY_ACCOUNT_BARCODE, error) {
	var debit_entries, credit_entries []VALUE_PRICE_QUANTITY_ACCOUNT_BARCODE
	var zero float64
	for _, entry := range entries {
		account_struct, _, _ := account_struct_from_name(entry.ACCOUNT_NAME)
		switch account_struct.IS_CREDIT {
		case false:
			zero += entry.VALUE
			if entry.VALUE > 0 {
				debit_entries = append(debit_entries, entry)
			} else if entry.VALUE < 0 {
				credit_entries = append(credit_entries, entry)
			}
		case true:
			zero -= entry.VALUE
			if entry.VALUE < 0 {
				debit_entries = append(debit_entries, entry)
			} else if entry.VALUE > 0 {
				credit_entries = append(credit_entries, entry)
			}
		}
	}
	if zero != 0 {
		return []VALUE_PRICE_QUANTITY_ACCOUNT_BARCODE{}, []VALUE_PRICE_QUANTITY_ACCOUNT_BARCODE{},
			fmt.Errorf("the debit and credit should be equal and the deffrence is %f", zero)
	}
	if (len(debit_entries) != 1) && (len(credit_entries) != 1) {
		return []VALUE_PRICE_QUANTITY_ACCOUNT_BARCODE{}, []VALUE_PRICE_QUANTITY_ACCOUNT_BARCODE{}, error_should_be_one_debit_or_one_credit
	}
	return debit_entries, credit_entries, nil
}

func JOURNAL_ENTRY(
	entries []VALUE_PRICE_QUANTITY_ACCOUNT_BARCODE,
	insert, auto_completion, invoice_discount bool,
	date_start, date_end time.Time,
	adjusting_method, notes, name, name_employee string,
	array_day_start_end []DAY_START_END) []JOURNAL_TAG {

	find_account_from_barcode(entries)
	remove_the_accounts_not_in_accounts_list(entries)
	date_end = set_date_end_to_zero_if_smaller_than_date_start(date_start, date_end)
	adjusting_method = set_adjusting_method(date_end, adjusting_method, entries)
	calculate_and_insert_value_price_quantity(entries)
	set_price_to_positive(entries)

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
	increase_the_value_to_make_the_new_balance_for_the_account_positive(entries)

	debit_entries, credit_entries, _ := check_debit_equal_credit(entries)
	simple_entries := convert_to_simple_entry(debit_entries, credit_entries)
	insert_to_journal_tag(simple_entries, date_start, date_end, notes, name, name_employee)

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
	return simple_entries
}
