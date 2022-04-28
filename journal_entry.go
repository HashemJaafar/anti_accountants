package main

import (
	"fmt"
	"time"
)

// func ADJUST_THE_ARRAY(array_to_insert []JOURNAL_TAG, array_start_end_minutes []start_end_minutes, adjusting_method string) [][]JOURNAL_TAG {
// 	var adjusted_array_to_insert [][]JOURNAL_TAG
// 	total_minutes := TOTAL_MINUTES(array_start_end_minutes)
// 	array_len_start_end_minutes := len(array_start_end_minutes) - 1
// 	for _, entry := range array_to_insert {
// 		var value_counter, time_unit_counter float64
// 		var one_account_adjusted_list []JOURNAL_TAG
// 		for index, element := range array_start_end_minutes {
// 			value := VALUE_AFTER_ADJUST_USING_ADJUSTING_METHODS(adjusting_method, element.minutes, total_minutes, time_unit_counter, entry.VALUE)

// 			if index == array_len_start_end_minutes {
// 				value = entry.VALUE - value_counter
// 			}

// 			time_unit_counter += element.minutes
// 			value_counter += value
// 			one_account_adjusted_list = append(one_account_adjusted_list, JOURNAL_TAG{
// 				IS_REVERSED:           false,
// 				ENTRY_NUMBER_COMPOUND: index,
// 				ENTRY_NUMBER_SIMPLE:   0,
// 				VALUE:                 value,
// 				PRICE_DEBIT:           entry.PRICE_DEBIT,
// 				PRICE_CREDIT:          entry.PRICE_CREDIT,
// 				QUANTITY_DEBIT:        RETURN_SAME_SIGN_OF_NUMBER_SIGN(value/entry.PRICE_DEBIT, entry.QUANTITY_DEBIT),
// 				QUANTITY_CREDIT:       RETURN_SAME_SIGN_OF_NUMBER_SIGN(value/entry.PRICE_CREDIT, entry.QUANTITY_CREDIT),
// 				ACCOUNT_DEBIT:         entry.ACCOUNT_DEBIT,
// 				ACCOUNT_CREDIT:        entry.ACCOUNT_CREDIT,
// 				NOTES:                 entry.NOTES,
// 				NAME:                  entry.NAME,
// 				NAME_EMPLOYEE:         entry.NAME_EMPLOYEE,
// 			})
// 		}
// 		adjusted_array_to_insert = append(adjusted_array_to_insert, one_account_adjusted_list)
// 	}
// 	return adjusted_array_to_insert
// }

// func CREATE_ARRAY_START_END_MINUTES(date_start, date_end time.Time, array_day_start_end []DAY_START_END) []start_end_minutes {
// 	var array_start_end_minutes []start_end_minutes
// 	var previous_date_end time.Time
// 	delta_days := int(date_end.Sub(date_start).Hours()/24 + 1)
// 	year, month, day := date_start.Date()
// 	for day_counter := 0; day_counter < delta_days; day_counter++ {
// 		for _, element := range array_day_start_end {
// 			start := time.Date(year, month, day+day_counter, element.START_HOUR, element.START_MINUTE, 0, 0, time.Local)
// 			if start.Weekday().String() == element.DAY {
// 				end := time.Date(year, month, day+day_counter, element.END_HOUR, element.END_MINUTE, 0, 0, time.Local)
// 				start, end = SHIFT_AND_ARRANGE_THE_TIME_SERIES(previous_date_end, start, end)
// 				array_start_end_minutes = append(array_start_end_minutes, start_end_minutes{start, end, end.Sub(start).Minutes()})
// 				previous_date_end = end
// 			}
// 		}
// 	}
// 	return array_start_end_minutes
// }

// func SET_ADJUSTING_METHOD(entry_expair time.Time, adjusting_method string) string {
// 	if entry_expair.IsZero() {
// 		return ""
// 	}
// 	if !IS_IN(adjusting_method, DEPRECIATION_METHODS) {
// 		return LINEAR
// 	}
// 	return adjusting_method
// }

// func SET_DATE_END_TO_ZERO_IF_SMALLER_THAN_DATE_START(date_start, date_end time.Time) time.Time {
// 	if !date_end.IsZero() && date_start.Before(date_end) {
// 		return time.Time{}
// 	}
// 	return date_end
// }

// func SET_SLICE_DAY_START_END(array_day_start_end []DAY_START_END) []DAY_START_END {
// 	if len(array_day_start_end) == 0 {
// 		array_day_start_end = []DAY_START_END{
// 			{SATURDAY, 0, 0, 23, 59},
// 			{SUNDAY, 0, 0, 23, 59},
// 			{MONDAY, 0, 0, 23, 59},
// 			{TUESDAY, 0, 0, 23, 59},
// 			{WEDNESDAY, 0, 0, 23, 59},
// 			{THURSDAY, 0, 0, 23, 59},
// 			{FRIDAY, 0, 0, 23, 59}}
// 	}
// 	for index := range array_day_start_end {
// 		array_day_start_end[index].DAY = strings.Title(array_day_start_end[index].DAY)

// 		if !IS_IN(array_day_start_end[index].DAY, STANDARD_DAYS) {
// 			array_day_start_end[index].DAY = SUNDAY
// 		}

// 		if array_day_start_end[index].START_HOUR < 0 {
// 			array_day_start_end[index].START_HOUR = 0
// 		}
// 		if array_day_start_end[index].START_HOUR > 23 {
// 			array_day_start_end[index].START_HOUR = 23
// 		}
// 		if array_day_start_end[index].START_MINUTE < 0 {
// 			array_day_start_end[index].START_MINUTE = 0
// 		}
// 		if array_day_start_end[index].START_MINUTE > 59 {
// 			array_day_start_end[index].START_MINUTE = 59
// 		}
// 		if array_day_start_end[index].END_HOUR < 0 {
// 			array_day_start_end[index].END_HOUR = 0
// 		}
// 		if array_day_start_end[index].END_HOUR > 23 {
// 			array_day_start_end[index].END_HOUR = 23
// 		}
// 		if array_day_start_end[index].END_MINUTE < 0 {
// 			array_day_start_end[index].END_MINUTE = 0
// 		}
// 		if array_day_start_end[index].END_MINUTE > 59 {
// 			array_day_start_end[index].END_MINUTE = 59
// 		}

// 		if array_day_start_end[index].START_HOUR > array_day_start_end[index].END_HOUR {
// 			array_day_start_end[index].START_HOUR = 0
// 		}
// 		if array_day_start_end[index].START_HOUR == array_day_start_end[index].END_HOUR && array_day_start_end[index].START_MINUTE > array_day_start_end[index].END_MINUTE {
// 			array_day_start_end[index].START_MINUTE = 0
// 		}
// 	}
// 	return array_day_start_end
// }

// func SHIFT_AND_ARRANGE_THE_TIME_SERIES(previous_date_end, date_start, date_end time.Time) (time.Time, time.Time) {
// 	if previous_date_end.After(date_start) {
// 		date_start = previous_date_end
// 	}
// 	if date_start.After(date_end) {
// 		date_end = date_start
// 	}
// 	return date_start, date_end
// }

// func TOTAL_MINUTES(array_start_end_minutes []start_end_minutes) float64 {
// 	var total_minutes float64
// 	for _, element := range array_start_end_minutes {
// 		total_minutes += element.minutes
// 	}
// 	return total_minutes
// }

// func VALUE_AFTER_ADJUST_USING_ADJUSTING_METHODS(adjusting_method string, minutes, TOTAL_MINUTES, time_unit_counter, total_value float64) float64 {
// 	percent := ROOT(total_value, TOTAL_MINUTES)
// 	switch adjusting_method {
// 	case EXPONENTIAL:
// 		return math.Pow(percent, time_unit_counter+minutes) - math.Pow(percent, time_unit_counter)
// 	case LOGARITHMIC:
// 		return (total_value / math.Pow(percent, time_unit_counter)) - (total_value / math.Pow(percent, time_unit_counter+minutes))
// 	default:
// 		return minutes * (total_value / TOTAL_MINUTES)
// 	}
// }

func CHECK_DEBIT_EQUAL_CREDIT(entries []PRICE_QUANTITY_ACCOUNT) ([]PRICE_QUANTITY_ACCOUNT, []PRICE_QUANTITY_ACCOUNT, error) {
	var debit_entries, credit_entries []PRICE_QUANTITY_ACCOUNT
	var zero float64
	for _, entry := range entries {
		value := entry.PRICE * entry.QUANTITY
		switch entry.IS_CREDIT {
		case false:
			zero += value
			if value > 0 {
				debit_entries = append(debit_entries, entry)
			} else if value < 0 {
				credit_entries = append(credit_entries, entry)
			}
		case true:
			zero -= value
			if value < 0 {
				debit_entries = append(debit_entries, entry)
			} else if value > 0 {
				credit_entries = append(credit_entries, entry)
			}
		}
	}
	if zero != 0 {
		return []PRICE_QUANTITY_ACCOUNT{}, []PRICE_QUANTITY_ACCOUNT{},
			fmt.Errorf("the debit and credit should be equal. and the debit is more than credit by %f, the debit_entries is %v and the credit_entries is %v", zero, debit_entries, credit_entries)
	}
	if len(debit_entries) != 1 && len(credit_entries) != 1 {
		return []PRICE_QUANTITY_ACCOUNT{}, []PRICE_QUANTITY_ACCOUNT{}, fmt.Errorf("should be one debit or one credit in the entry, the debit_entries is %v and the credit_entries is %v", debit_entries, credit_entries)
	}
	return debit_entries, credit_entries, nil
}

func SET_PRICE_AND_QUANTITY(account PRICE_QUANTITY_ACCOUNT, is_update bool) PRICE_QUANTITY_ACCOUNT {
	if account.QUANTITY > 0 {
		return account
	}

	// i make it this way just to make it faster when using WMA case
	var keys [][]byte
	var inventory []INVENTORY_TAG
	switch account.COST_FLOW_TYPE {
	case FIFO:
		keys, inventory = DB_READ[INVENTORY_TAG](DB_INVENTORY)
	case LIFO:
		keys, inventory = DB_READ[INVENTORY_TAG](DB_INVENTORY)
		REVERSE_SLICE(keys)
		REVERSE_SLICE(inventory)
	case WMA:
		WEIGHTED_AVERAGE(account.ACCOUNT_NAME)
		keys, inventory = DB_READ[INVENTORY_TAG](DB_INVENTORY)
	}

	quantity_count := ABS(account.QUANTITY)
	var costs float64
	for k1, v1 := range inventory {
		if v1.ACCOUNT_NAME == account.ACCOUNT_NAME {
			if quantity_count <= v1.QUANTITY {
				costs -= v1.PRICE * quantity_count
				if is_update {
					inventory[k1].QUANTITY -= quantity_count
					DB_UPDATE(DB_INVENTORY, keys[k1], inventory[k1])
				}
				quantity_count = 0
				break
			}
			if quantity_count > v1.QUANTITY {
				costs -= v1.PRICE * v1.QUANTITY
				if is_update {
					inventory[k1].QUANTITY = 0
					DB_UPDATE(DB_INVENTORY, keys[k1], inventory[k1])
				}
				quantity_count -= v1.QUANTITY
			}
		}
	}
	account.QUANTITY += quantity_count
	account.PRICE = costs / account.QUANTITY
	return account
}

func GROUP_BY_ACCOUNT(entries []PRICE_QUANTITY_ACCOUNT) []PRICE_QUANTITY_ACCOUNT {
	m := map[string]*PRICE_QUANTITY_ACCOUNT{}
	for _, v1 := range entries {
		key := v1.ACCOUNT_NAME
		sums := m[key]
		if sums == nil {
			sums = &PRICE_QUANTITY_ACCOUNT{}
			m[key] = sums
		}
		// i make this to store the value and then devide it by the quantity to get the price
		sums.IS_CREDIT = v1.IS_CREDIT
		sums.COST_FLOW_TYPE = v1.COST_FLOW_TYPE
		sums.ACCOUNT_NAME = v1.ACCOUNT_NAME
		sums.PRICE += v1.PRICE * v1.QUANTITY //here i store the value in price field
		sums.QUANTITY += v1.QUANTITY
	}
	entries = []PRICE_QUANTITY_ACCOUNT{}
	for _, v1 := range m {
		entries = append(entries, PRICE_QUANTITY_ACCOUNT{
			IS_CREDIT:      v1.IS_CREDIT,
			COST_FLOW_TYPE: v1.COST_FLOW_TYPE,
			ACCOUNT_NAME:   v1.ACCOUNT_NAME,
			PRICE:          v1.PRICE / v1.QUANTITY,
			QUANTITY:       v1.QUANTITY,
		})
	}
	return entries
}

func INSERT_ENTRY_NUMBER(array_of_journal_tag []JOURNAL_TAG) {
	journal_tag := DB_LAST_LINE[JOURNAL_TAG](DB_JOURNAL)
	var last_entry_number_compound int
	var entry_number_simple int
	for k1, v1 := range array_of_journal_tag {
		array_of_journal_tag[k1].ENTRY_NUMBER_COMPOUND = journal_tag.ENTRY_NUMBER_COMPOUND + 1
		if v1.ENTRY_NUMBER_COMPOUND != last_entry_number_compound {
			entry_number_simple = 0
			last_entry_number_compound = v1.ENTRY_NUMBER_COMPOUND
		}
		entry_number_simple++
		array_of_journal_tag[k1].ENTRY_NUMBER_SIMPLE = entry_number_simple
	}
}

func INSERT_TO_DATABASE_JOURNAL(entries []JOURNAL_TAG) {
	INSERT_ENTRY_NUMBER(entries)
	for _, v1 := range entries {
		DB_UPDATE(DB_JOURNAL, NOW(), v1)
	}
}

func INSERT_TO_JOURNAL_TAG(debit_entries, credit_entries []PRICE_QUANTITY_ACCOUNT, notes, name, name_employee string) []JOURNAL_TAG {
	var simple_entries []JOURNAL_TAG
	for _, debit_entry := range debit_entries {
		for _, credit_entry := range credit_entries {
			value := SMALLEST(ABS(debit_entry.PRICE*debit_entry.QUANTITY), ABS(credit_entry.PRICE*credit_entry.QUANTITY))
			simple_entries = append(simple_entries, JOURNAL_TAG{
				IS_REVERSED:           false,
				ENTRY_NUMBER_COMPOUND: 0,
				ENTRY_NUMBER_SIMPLE:   0,
				VALUE:                 value,
				PRICE_DEBIT:           debit_entry.PRICE,
				PRICE_CREDIT:          credit_entry.PRICE,
				QUANTITY_DEBIT:        value / debit_entry.PRICE,
				QUANTITY_CREDIT:       value / credit_entry.PRICE,
				ACCOUNT_DEBIT:         debit_entry.ACCOUNT_NAME,
				ACCOUNT_CREDIT:        credit_entry.ACCOUNT_NAME,
				NOTES:                 notes,
				NAME:                  name,
				NAME_EMPLOYEE:         name_employee,
			})
		}
	}
	return simple_entries
}

func SIMPLE_JOURNAL_ENTRY(
	entries []PRICE_QUANTITY_ACCOUNT_BARCODE,
	insert, auto_completion, invoice_discount bool,
	notes, name, name_employee string) ([]PRICE_QUANTITY_ACCOUNT_BARCODE, error) {

	slice_of_price_quantity_account := STAGE_1(entries)
	slice_of_price_quantity_account = GROUP_BY_ACCOUNT(slice_of_price_quantity_account)

	for k1, v1 := range slice_of_price_quantity_account {
		slice_of_price_quantity_account[k1] = SET_PRICE_AND_QUANTITY(v1, false)
	}

	// if auto_completion {
	// 	AUTO_COMPLETION_THE_ENTRY(slice_of_price_quantity_account)
	// }
	// if invoice_discount {
	// 	entries = auto_completion_the_invoice_discount(entries)
	// }

	slice_of_price_quantity_account = GROUP_BY_ACCOUNT(slice_of_price_quantity_account)
	debit_entries, credit_entries, err := CHECK_DEBIT_EQUAL_CREDIT(slice_of_price_quantity_account)
	if err != nil {
		return []PRICE_QUANTITY_ACCOUNT_BARCODE{}, err
	}
	simple_entries := INSERT_TO_JOURNAL_TAG(debit_entries, credit_entries, notes, name, name_employee)

	if insert {
		INSERT_TO_DATABASE_JOURNAL(simple_entries)
		INSERT_TO_DATABASE_INVENTORY(slice_of_price_quantity_account)
	}

	return CONVERT_PRICE_QUANTITY_ACCOUNT_TO_PRICE_QUANTITY_ACCOUNT_BARCODE(append(debit_entries, credit_entries...)), nil
}

func CONVERT_PRICE_QUANTITY_ACCOUNT_TO_PRICE_QUANTITY_ACCOUNT_BARCODE(entries []PRICE_QUANTITY_ACCOUNT) []PRICE_QUANTITY_ACCOUNT_BARCODE {
	var new_entries []PRICE_QUANTITY_ACCOUNT_BARCODE
	for _, v1 := range entries {
		new_entries = append(new_entries, PRICE_QUANTITY_ACCOUNT_BARCODE{
			PRICE:        v1.PRICE,
			QUANTITY:     v1.QUANTITY,
			ACCOUNT_NAME: v1.ACCOUNT_NAME,
			BARCODE:      "",
		})
	}
	return new_entries
}

func INSERT_TO_DATABASE_INVENTORY(entries []PRICE_QUANTITY_ACCOUNT) {
	for _, v1 := range entries {
		if v1.QUANTITY > 0 {
			DB_UPDATE(DB_INVENTORY, NOW(), INVENTORY_TAG{v1.PRICE, v1.QUANTITY, v1.ACCOUNT_NAME})
		} else {
			SET_PRICE_AND_QUANTITY(v1, true)
		}
	}
}

func STAGE_1(entries []PRICE_QUANTITY_ACCOUNT_BARCODE) []PRICE_QUANTITY_ACCOUNT {
	var array_price_quantity_account []PRICE_QUANTITY_ACCOUNT
	for _, v1 := range entries {
		account_struct, _, err := ACCOUNT_STRUCT_FROM_BARCODE(v1.BARCODE)
		if err != nil {
			account_struct, _, err = ACCOUNT_STRUCT_FROM_NAME(FORMAT_THE_STRING(v1.ACCOUNT_NAME))
		}
		if err == nil && account_struct.IS_LOW_LEVEL_ACCOUNT && v1.QUANTITY != 0 && v1.PRICE != 0 {
			array_price_quantity_account = append(array_price_quantity_account, PRICE_QUANTITY_ACCOUNT{
				IS_CREDIT:      account_struct.IS_CREDIT,
				COST_FLOW_TYPE: account_struct.COST_FLOW_TYPE,
				ACCOUNT_NAME:   account_struct.ACCOUNT_NAME,
				PRICE:          ABS(v1.PRICE),
				QUANTITY:       v1.QUANTITY,
			})
		}
	}
	return array_price_quantity_account
}

func REVERSE_ENTRIES(entry_number_compound, entry_number_simple int, name_employee string) {
	var entries []JOURNAL_TAG
	var entries_keys [][]byte
	keys, journal := DB_READ[JOURNAL_TAG](DB_JOURNAL)
	for k1, v1 := range journal {
		if v1.ENTRY_NUMBER_COMPOUND == entry_number_compound && (entry_number_simple == 0 || v1.ENTRY_NUMBER_SIMPLE == entry_number_simple) && v1.IS_REVERSED == false {
			entries = append(entries, v1)
			entries_keys = append(entries_keys, keys[k1])
		}
	}

	var entry_to_reverse []JOURNAL_TAG
	for k1, v1 := range entries {
		// here i check if the account in the entry is credit and it have debit nature then it will be negative quantity and vice versa
		account_struct_credit, _, _ := ACCOUNT_STRUCT_FROM_NAME(v1.ACCOUNT_CREDIT)
		if !account_struct_credit.IS_CREDIT {
			v1.QUANTITY_CREDIT *= -1
		}
		// here i check if the account in the entry is debit and it have credit nature then it will be negative quantity and vice versa
		account_struct_debit, _, _ := ACCOUNT_STRUCT_FROM_NAME(v1.ACCOUNT_DEBIT)
		if account_struct_debit.IS_CREDIT {
			v1.QUANTITY_DEBIT *= -1
		}

		// here i check if the account can be negative by seeing the difference in quantity after the find the cost in inventory.
		// because i dont want to make the account negative balance
		entry_credit := SET_PRICE_AND_QUANTITY(PRICE_QUANTITY_ACCOUNT{false, FIFO, v1.ACCOUNT_CREDIT, v1.PRICE_CREDIT, v1.QUANTITY_CREDIT}, false)
		entry_debit := SET_PRICE_AND_QUANTITY(PRICE_QUANTITY_ACCOUNT{false, FIFO, v1.ACCOUNT_DEBIT, v1.PRICE_DEBIT, v1.QUANTITY_DEBIT}, false)

		// here i compare the quantity if it is the same i will reverse the entry
		if entry_credit.QUANTITY == v1.QUANTITY_CREDIT && entry_debit.QUANTITY == v1.QUANTITY_DEBIT {

			// here i change the cost flow to wma just to make outflow from the inventory without error
			entry_credit.COST_FLOW_TYPE = WMA
			entry_debit.COST_FLOW_TYPE = WMA

			// here i insert to the inventory
			INSERT_TO_DATABASE_INVENTORY([]PRICE_QUANTITY_ACCOUNT{entry_credit, entry_debit})

			// i swap the debit and credit with each other but the quantity after i swap it will be positive
			v1.PRICE_CREDIT, v1.PRICE_DEBIT = v1.PRICE_DEBIT, v1.PRICE_CREDIT
			v1.QUANTITY_CREDIT, v1.QUANTITY_DEBIT = ABS(v1.QUANTITY_DEBIT), ABS(v1.QUANTITY_CREDIT)
			v1.ACCOUNT_CREDIT, v1.ACCOUNT_DEBIT = v1.ACCOUNT_DEBIT, v1.ACCOUNT_CREDIT

			v1.NOTES = "revese entry for entry was entered by " + v1.NAME_EMPLOYEE
			v1.NAME_EMPLOYEE = name_employee
			v1.IS_REVERSE = true
			v1.REVERSE_ENTRY_NUMBER_COMPOUND = v1.ENTRY_NUMBER_COMPOUND
			v1.REVERSE_ENTRY_NUMBER_SIMPLE = v1.ENTRY_NUMBER_SIMPLE

			// here i append the entry to the journal to reverse all in one entry number compound
			entry_to_reverse = append(entry_to_reverse, v1)

			// i make the reverse field in the entry true just to not reverse it again
			entries[k1].IS_REVERSED = true
			DB_UPDATE(DB_JOURNAL, entries_keys[k1], entries[k1])
		}
	}

	// and then i insert to database
	INSERT_TO_DATABASE_JOURNAL(entry_to_reverse)
}

func JOURNAL_FILTER(dates []time.Time, journal []JOURNAL_TAG, f THE_JOURNAL_FILTER, is_debit_and_credit bool) ([]time.Time, []JOURNAL_TAG) {
	var filtered_journal []JOURNAL_TAG
	var filtered_dates []time.Time
	for k1, v1 := range journal {

		// here icheck if the user whant the debit and credit or one of them each time
		var is_the_accounts bool
		if is_debit_and_credit {
			is_the_accounts = FUNCTION_FILTER_STRING(v1.ACCOUNT_DEBIT, f.ACCOUNT_DEBIT) && FUNCTION_FILTER_STRING(v1.ACCOUNT_CREDIT, f.ACCOUNT_CREDIT)
		} else {
			is_the_accounts = FUNCTION_FILTER_STRING(v1.ACCOUNT_DEBIT, f.ACCOUNT_DEBIT) || FUNCTION_FILTER_STRING(v1.ACCOUNT_CREDIT, f.ACCOUNT_CREDIT)
		}

		if is_the_accounts &&
			FUNCTION_FILTER_DATE(dates[k1], f.DATE) &&
			FUNCTION_FILTER_BOOL(v1.IS_REVERSE, f.IS_REVERSE) &&
			FUNCTION_FILTER_BOOL(v1.IS_REVERSED, f.IS_REVERSED) &&
			FUNCTION_FILTER_NUMBER(float64(v1.REVERSE_ENTRY_NUMBER_COMPOUND), f.REVERSE_ENTRY_NUMBER_COMPOUND) &&
			FUNCTION_FILTER_NUMBER(float64(v1.REVERSE_ENTRY_NUMBER_SIMPLE), f.REVERSE_ENTRY_NUMBER_SIMPLE) &&
			FUNCTION_FILTER_NUMBER(float64(v1.ENTRY_NUMBER_COMPOUND), f.ENTRY_NUMBER_COMPOUND) &&
			FUNCTION_FILTER_NUMBER(float64(v1.ENTRY_NUMBER_SIMPLE), f.ENTRY_NUMBER_SIMPLE) &&
			FUNCTION_FILTER_NUMBER(v1.VALUE, f.VALUE) &&
			FUNCTION_FILTER_NUMBER(v1.PRICE_DEBIT, f.PRICE_DEBIT) &&
			FUNCTION_FILTER_NUMBER(v1.PRICE_CREDIT, f.PRICE_CREDIT) &&
			FUNCTION_FILTER_NUMBER(v1.QUANTITY_DEBIT, f.QUANTITY_DEBIT) &&
			FUNCTION_FILTER_NUMBER(v1.QUANTITY_CREDIT, f.QUANTITY_CREDIT) &&
			FUNCTION_FILTER_STRING(v1.NOTES, f.NOTES) &&
			FUNCTION_FILTER_STRING(v1.NAME, f.NAME) &&
			FUNCTION_FILTER_STRING(v1.NAME_EMPLOYEE, f.NAME_EMPLOYEE) {
			filtered_journal = append(filtered_journal, v1)
			filtered_dates = append(filtered_dates, dates[k1])
		}
	}
	return filtered_dates, filtered_journal
}

func FIND_DUPLICATE_ELEMENT(journal []JOURNAL_TAG, f THE_JOURNAL_DUPLICATE_FILTER) []JOURNAL_TAG {
	var filtered_journal []JOURNAL_TAG
	for k1, v1 := range journal {
		for k2, v2 := range journal {
			if k1 < k2 &&
				FUNCTION_FILTER_DUPLICATE(v1.IS_REVERSE, v2.IS_REVERSE, f.IS_REVERSE) &&
				FUNCTION_FILTER_DUPLICATE(v1.IS_REVERSED, v2.IS_REVERSED, f.IS_REVERSED) &&
				FUNCTION_FILTER_DUPLICATE(v1.REVERSE_ENTRY_NUMBER_COMPOUND, v2.REVERSE_ENTRY_NUMBER_COMPOUND, f.REVERSE_ENTRY_NUMBER_COMPOUND) &&
				FUNCTION_FILTER_DUPLICATE(v1.REVERSE_ENTRY_NUMBER_SIMPLE, v2.REVERSE_ENTRY_NUMBER_SIMPLE, f.REVERSE_ENTRY_NUMBER_SIMPLE) &&
				FUNCTION_FILTER_DUPLICATE(v1.VALUE, v2.VALUE, f.VALUE) &&
				FUNCTION_FILTER_DUPLICATE(v1.PRICE_DEBIT, v2.PRICE_DEBIT, f.PRICE_DEBIT) &&
				FUNCTION_FILTER_DUPLICATE(v1.PRICE_CREDIT, v2.PRICE_CREDIT, f.PRICE_CREDIT) &&
				FUNCTION_FILTER_DUPLICATE(v1.QUANTITY_DEBIT, v2.QUANTITY_DEBIT, f.QUANTITY_DEBIT) &&
				FUNCTION_FILTER_DUPLICATE(v1.QUANTITY_CREDIT, v2.QUANTITY_CREDIT, f.QUANTITY_CREDIT) &&
				FUNCTION_FILTER_DUPLICATE(v1.ACCOUNT_DEBIT, v2.ACCOUNT_DEBIT, f.ACCOUNT_DEBIT) &&
				FUNCTION_FILTER_DUPLICATE(v1.ACCOUNT_CREDIT, v2.ACCOUNT_CREDIT, f.ACCOUNT_CREDIT) &&
				FUNCTION_FILTER_DUPLICATE(v1.NOTES, v2.NOTES, f.NOTES) &&
				FUNCTION_FILTER_DUPLICATE(v1.NAME, v2.NAME, f.NAME) &&
				FUNCTION_FILTER_DUPLICATE(v1.NAME_EMPLOYEE, v2.NAME_EMPLOYEE, f.NAME_EMPLOYEE) {
				filtered_journal = append(filtered_journal, v1)
			}
		}
	}
	return filtered_journal
}
