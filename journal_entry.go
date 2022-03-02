package anti_accountants

import (
	"math"
	"time"
)

func check_the_adjusting_method_and_date(entry_expair time.Time, date time.Time, adjusting_method string, entries []ACCOUNT_VALUE_PRICE_QUANTITY_BARCODE) {

	is_in_adjusting_methods := is_in(adjusting_method, adjusting_methods)
	is_entry_expair_zero := entry_expair.IsZero()

	if !is_entry_expair_zero {
		check_dates(date, entry_expair)
	}

	if !is_in_adjusting_methods && adjusting_method != "" {
		error_element_is_not_in_elements(adjusting_method, adjusting_methods)
	}
	if is_entry_expair_zero == is_in_adjusting_methods {
		error_you_cant_use_entry_expire()
	}
	for _, entry := range entries {
		if is_in(entry.ACCOUNT, inventory) && is_in(adjusting_method, depreciation_methods) {
			error_you_cant_use_depreciation_methods_with_inventory(entry.ACCOUNT)
		}
	}
}

func group_by_account_and_barcode(entries []ACCOUNT_VALUE_PRICE_QUANTITY_BARCODE) []ACCOUNT_VALUE_PRICE_QUANTITY_BARCODE {
	type account_barcode struct {
		account, barcode string
	}
	g := map[account_barcode]*ACCOUNT_VALUE_PRICE_QUANTITY_BARCODE{}
	for _, v := range entries {
		key := account_barcode{v.ACCOUNT, v.BARCODE}
		sums := g[key]
		if sums == nil {
			sums = &ACCOUNT_VALUE_PRICE_QUANTITY_BARCODE{}
			g[key] = sums
		}
		sums.VALUE += v.VALUE
		sums.QUANTITY += v.QUANTITY
	}
	entries = []ACCOUNT_VALUE_PRICE_QUANTITY_BARCODE{}
	for key, v := range g {
		entries = append(entries, ACCOUNT_VALUE_PRICE_QUANTITY_BARCODE{
			ACCOUNT:  key.account,
			VALUE:    v.VALUE,
			PRICE:    v.VALUE / v.QUANTITY,
			QUANTITY: v.QUANTITY,
			BARCODE:  key.barcode,
		})
	}
	return entries
}

func remove_zero_values(entries []ACCOUNT_VALUE_PRICE_QUANTITY_BARCODE) {
	var index int
	for index < len(entries) {
		if entries[index].VALUE == 0 || entries[index].QUANTITY == 0 {
			entries = append(entries[:index], entries[index+1:]...)
		} else {
			index++
		}
	}
}

func find_cost(entries []ACCOUNT_VALUE_PRICE_QUANTITY_BARCODE) {
	for index, entry := range entries {
		costs := cost_flow(entry.ACCOUNT, entry.QUANTITY, entry.BARCODE, false)
		if costs != 0 {
			entries[index].VALUE = -costs
			entries[index].PRICE = -costs / entry.QUANTITY
		}
	}
}

func convert_to_simple_entry(debit_entries, credit_entries []ACCOUNT_VALUE_PRICE_QUANTITY_BARCODE) [][]ACCOUNT_VALUE_PRICE_QUANTITY_BARCODE {
	simple_entries := [][]ACCOUNT_VALUE_PRICE_QUANTITY_BARCODE{}
	for _, debit_entry := range debit_entries {
		for _, credit_entry := range credit_entries {
			simple_entries = append(simple_entries, []ACCOUNT_VALUE_PRICE_QUANTITY_BARCODE{debit_entry, credit_entry})
		}
	}
	for _, a := range simple_entries {
		switch math.Abs(a[0].VALUE) >= math.Abs(a[1].VALUE) {
		case true:
			sign := a[0].VALUE / a[1].VALUE
			price := a[0].VALUE / a[0].QUANTITY
			a[0].VALUE = a[1].VALUE * sign / math.Abs(sign)
			a[0].QUANTITY = a[0].VALUE / price
		case false:
			sign := a[0].VALUE / a[1].VALUE
			price := a[1].VALUE / a[1].QUANTITY
			a[1].VALUE = a[0].VALUE * sign / math.Abs(sign)
			a[1].QUANTITY = a[1].VALUE / price
		}
	}
	return simple_entries
}

func can_the_account_be_negative(entries []ACCOUNT_VALUE_PRICE_QUANTITY_BARCODE) {
	for _, entry := range entries {
		if !(is_it_sub_account_using_name(PRIMARY_ACCOUNTS_NAMES.EQUITY, entry.ACCOUNT) && is_credit(entry.ACCOUNT)) {
			account_balance := account_balance(entry.ACCOUNT)
			if account_balance+entry.VALUE < 0 {
				error_make_nagtive_balance(entry, account_balance)
			}
		}
	}
}

func find_account_from_barcode(entries []ACCOUNT_VALUE_PRICE_QUANTITY_BARCODE) {
	for index, entry := range entries {
		if entry.QUANTITY < 0 && entry.BARCODE != "" {
			err := DB.QueryRow("select account from journal where barcode=? limit 1", entry.BARCODE).Scan(&entries[index].ACCOUNT)
			if err != nil {
				error_the_barcode_is_wrong(entry)
			}
		}
	}
}

func insert_to_JOURNAL_TAG(entries []ACCOUNT_VALUE_PRICE_QUANTITY_BARCODE, date time.Time, entry_expair time.Time, description string, name string, employee_name string) []JOURNAL_TAG {
	var array_to_insert []JOURNAL_TAG
	for _, entry := range entries {
		array_to_insert = append(array_to_insert, JOURNAL_TAG{
			DATE:          date.String(),
			ENTRY_NUMBER:  0,
			ACCOUNT:       entry.ACCOUNT,
			VALUE:         entry.VALUE,
			PRICE:         entry.VALUE / entry.QUANTITY,
			QUANTITY:      entry.QUANTITY,
			BARCODE:       entry.BARCODE,
			ENTRY_EXPAIR:  entry_expair.String(),
			DESCRIPTION:   description,
			NAME:          name,
			EMPLOYEE_NAME: employee_name,
			ENTRY_DATE:    NOW.String(),
			REVERSE:       false,
		})
	}
	return array_to_insert
}

func check_if_the_price_is_negative(entries []ACCOUNT_VALUE_PRICE_QUANTITY_BARCODE) {
	for index, entry := range entries {
		entries[index].PRICE = entry.VALUE / entry.QUANTITY
		if entries[index].PRICE < 0 {
			error_the_price_should_be_positive(entries[index])
		}
	}
}

func cost_flow(account string, quantity float64, barcode string, insert bool) float64 {
	if quantity > 0 {
		return 0
	}
	order_by_date_asc_or_desc := asc_or_desc(account)
	if order_by_date_asc_or_desc == "" {
		return 0
	}
	rows, _ := DB.Query("select price,quantity from inventory where quantity>0 and account=? and barcode=? order by date "+order_by_date_asc_or_desc, account, barcode)
	var inventory []JOURNAL_TAG
	for rows.Next() {
		var tag JOURNAL_TAG
		rows.Scan(&tag.PRICE, &tag.QUANTITY)
		inventory = append(inventory, tag)
	}
	quantity = math.Abs(quantity)
	quantity_count := quantity
	var costs float64
	for _, item := range inventory {
		if item.QUANTITY > quantity_count {
			costs += item.PRICE * quantity_count
			if insert {
				DB.Exec("update inventory set quantity=quantity-? where account=? and price=? and quantity=? and barcode=? order by date "+order_by_date_asc_or_desc+" limit 1", quantity_count, account, item.PRICE, item.QUANTITY, barcode)
			}
			quantity_count = 0
			break
		}
		if item.QUANTITY <= quantity_count {
			costs += item.PRICE * item.QUANTITY
			if insert {
				DB.Exec("delete from inventory where account=? and price=? and quantity=? and barcode=? order by date "+order_by_date_asc_or_desc+" limit 1", account, item.PRICE, item.QUANTITY, barcode)
			}
			quantity_count -= item.QUANTITY
		}
	}
	if quantity_count != 0 {
		error_the_order_out_of_stock(quantity, quantity_count, account, barcode)
	}
	return costs
}

func asc_or_desc(account string) string {
	switch return_cost_flow_type(account) {
	case "lifo":
		return "desc"
	case "fifo":
		return "asc"
	case "wma":
		weighted_average(account)
		return "asc"
	case "barcode":
		weighted_average_for_barcode(account)
		return "asc"
	}
	return ""
}

func insert_to_database(array_of_journal_tag []JOURNAL_TAG, insert_into_journal, insert_into_inventory bool) {
	insert_entry_number(array_of_journal_tag)
	if insert_into_journal {
		insert_into_journal_func(array_of_journal_tag)
	}
	if insert_into_inventory {
		insert_into_inventory_func(array_of_journal_tag)
	}
}

func insert_entry_number(array_of_journal_tag []JOURNAL_TAG) {
	entry_number := float64(entry_number())
	for indexa := range array_of_journal_tag {
		array_of_journal_tag[indexa].ENTRY_NUMBER = int(entry_number)
		entry_number += 0.5
	}
}

func calculate_and_insert_value_price_quantity(entries []ACCOUNT_VALUE_PRICE_QUANTITY_BARCODE) {
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

func check_if_the_account_is_high_by_level(entries []ACCOUNT_VALUE_PRICE_QUANTITY_BARCODE) {
	for _, entry := range entries {
		if is_it_high_by_level(account_number(entry.ACCOUNT)) {
			error_is_high_level_account(entry.ACCOUNT)
		}
	}
}

func JOURNAL_ENTRY(
	entries []ACCOUNT_VALUE_PRICE_QUANTITY_BARCODE,
	insert, auto_completion, invoice_discount bool,
	date time.Time, entry_expair time.Time,
	adjusting_method, description, name, employee_name string,
	array_day_start_end []DAY_START_END) []JOURNAL_TAG {

	find_account_from_barcode(entries)
	check_the_adjusting_method_and_date(entry_expair, date, adjusting_method, entries)
	calculate_and_insert_value_price_quantity(entries)

	entries = group_by_account_and_barcode(entries)
	remove_zero_values(entries)

	find_cost(entries)

	if auto_completion {
		entries = auto_completion_the_entry(entries)
	}
	if invoice_discount {
		entries = auto_completion_the_invoice_discount(entries)
	}

	entries = group_by_account_and_barcode(entries)
	remove_zero_values(entries)

	check_if_the_account_is_high_by_level(entries)
	can_the_account_be_negative(entries)

	check_debit_equal_credit(entries)
	debit_entries, credit_entries := separate_debit_from_credit(entries)
	check_one_debit_or_one_credit(debit_entries, credit_entries)
	simple_entries := convert_to_simple_entry(debit_entries, credit_entries)

	var all_array_to_insert []JOURNAL_TAG
	for _, simple_entry := range simple_entries {

		check_if_the_price_is_negative(simple_entry)
		array_to_insert := insert_to_JOURNAL_TAG(simple_entry, date, entry_expair, description, name, employee_name)

		if is_in(adjusting_method, depreciation_methods) {
			array_day_start_end = initialize_array_day_start_end(array_day_start_end)
			check_array_day_start_end(array_day_start_end)
			array_start_end_minutes := create_array_start_end_minutes(entry_expair, date, array_day_start_end)
			total_minutes := total_minutes(array_start_end_minutes)
			adjusted_array_to_insert := adjust_the_array(array_to_insert, array_start_end_minutes, total_minutes, adjusting_method, description, name, employee_name)
			adjusted_array_to_insert = transpose(adjusted_array_to_insert)
			array_to_insert = unpack_the_array(adjusted_array_to_insert)
		}

		all_array_to_insert = append(all_array_to_insert, array_to_insert...)
	}

	insert_to_database(all_array_to_insert, insert, insert)
	return all_array_to_insert
}
