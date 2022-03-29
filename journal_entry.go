package anti_accountants

import (
	"fmt"
	"math"
	"strings"
	"time"
)

func SET_DATE_END_TO_ZERO_IF_SMALLER_THAN_DATE_START(date_start, date_end time.Time) time.Time {
	if !date_end.IsZero() {
		if !date_start.Before(date_end) {
			return time.Time{}
		}
	}
	return date_end
}

func SET_ADJUSTING_METHOD(entry_expair time.Time, adjusting_method string, entries []PRICE_QUANTITY_ACCOUNT_BARCODE) string {
	if entry_expair.IsZero() {
		return ""
	}
	if !IS_IN(adjusting_method, ADJUSTING_METHODS) {
		adjusting_method = LINEAR
	}
	if IS_IN(adjusting_method, DEPRECIATION_METHODS) {
		for _, entry := range entries {
			account_struct, _, _ := ACCOUNT_STRUCT_FROM_NAME(entry.ACCOUNT_NAME)
			if account_struct.COST_FLOW_TYPE != "" {
				adjusting_method = EXPIRE
			}
		}
	}
	return adjusting_method
}

func GROUP_BY_ACCOUNT(entries []PRICE_QUANTITY_ACCOUNT_BARCODE) []PRICE_QUANTITY_ACCOUNT_BARCODE {
	m := map[string]*PRICE_QUANTITY_ACCOUNT_BARCODE{}
	for _, v := range entries {
		key := v.ACCOUNT_NAME
		sums := m[key]
		if sums == nil {
			sums = &PRICE_QUANTITY_ACCOUNT_BARCODE{}
			m[key] = sums
		}
		sums.PRICE += v.PRICE * v.QUANTITY // i make this to store the value and then devide it by the quantity to get the price
		sums.QUANTITY += v.QUANTITY
	}
	entries = []PRICE_QUANTITY_ACCOUNT_BARCODE{}
	for key, v := range m {
		entries = append(entries, PRICE_QUANTITY_ACCOUNT_BARCODE{
			PRICE:        v.PRICE / v.QUANTITY,
			QUANTITY:     v.QUANTITY,
			ACCOUNT_NAME: key,
		})
	}
	return entries
}

func FIND_COST(entries []PRICE_QUANTITY_ACCOUNT_BARCODE) {
	for index, entry := range entries {
		costs, _, err := COST_FLOW(entry.ACCOUNT_NAME, entry.QUANTITY, false)
		if err == nil {
			entries[index].PRICE = -costs / entry.QUANTITY
		}
	}
}

func INSERT_TO_JOURNAL_TAG(debit_entries, credit_entries []PRICE_QUANTITY_ACCOUNT_BARCODE, date_start, date_end time.Time, notes, name, name_employee string) []JOURNAL_TAG {
	var simple_entries []JOURNAL_TAG
	for _, debit_entry := range debit_entries {
		for _, credit_entry := range credit_entries {
			simple_entries = append(simple_entries, JOURNAL_TAG{
				VALUE:           math.Abs(SMALLEST(debit_entry.PRICE/debit_entry.QUANTITY, credit_entry.PRICE/credit_entry.QUANTITY)),
				PRICE_DEBIT:     debit_entry.PRICE,
				PRICE_CREDIT:    credit_entry.PRICE,
				QUANTITY_DEBIT:  debit_entry.QUANTITY,
				QUANTITY_CREDIT: credit_entry.QUANTITY,
				ACCOUNT_DEBIT:   debit_entry.ACCOUNT_NAME,
				ACCOUNT_CREDIT:  credit_entry.ACCOUNT_NAME,
				NOTES:           notes,
				NAME:            name,
				NAME_EMPLOYEE:   name_employee,
				DATE_START:      date_start,
				DATE_END:        date_end,
			})
		}
	}
	return simple_entries
}

func INCREASE_THE_VALUE_TO_MAKE_THE_NEW_BALANCE_FOR_THE_ACCOUNT_POSITIVE(entries []PRICE_QUANTITY_ACCOUNT_BARCODE) {
	for indexa, a := range entries {
		new_balance := ACCOUNT_BALANCE(a.ACCOUNT_NAME) + a.PRICE*a.QUANTITY
		if new_balance < 0 {
			entries[indexa].QUANTITY -= new_balance / entries[indexa].PRICE
		}
	}
}

func COST_FLOW(account string, quantity float64, is_update bool) (float64, float64, error) {
	if quantity > 0 {
		return 0, quantity, ERROR_SHOULD_BE_NEGATIVE
	}
	flow_type, is_ascending := FLOW_TYPE(account)
	if flow_type == "" {
		return 0, quantity, ERROR_NOT_INVENTORY_ACCOUNT
	}
	keys, inventory := DB_READ[INVENTORY_TAG](DB_INVENTORY)
	SORT_BY_TIME_INVENTORY(inventory, keys, is_ascending)
	quantity_count := math.Abs(quantity)
	var costs float64
	for indexa, a := range inventory {
		if a.ACCOUNT_NAME == account {
			if quantity_count <= a.QUANTITY {
				costs += a.PRICE * quantity_count
				if is_update {
					a.QUANTITY -= quantity_count
					DB_UPDATE(DB_INVENTORY, keys[indexa], a)
				}
				quantity_count -= quantity_count
				break
			}
			if quantity_count > a.QUANTITY {
				costs += a.PRICE * a.QUANTITY
				if is_update {
					a.QUANTITY -= a.QUANTITY
					DB_UPDATE(DB_INVENTORY, keys[indexa], a)
				}
				quantity_count -= a.QUANTITY
			}
		}
	}
	quantity += quantity_count
	return costs, quantity, nil
}

func FLOW_TYPE(account string) (string, bool) {
	account_struct, _, _ := ACCOUNT_STRUCT_FROM_NAME(account)
	switch account_struct.COST_FLOW_TYPE {
	case SPECIFIC_IDENTIFICATION:
		return SPECIFIC_IDENTIFICATION, true
	case LIFO:
		return LIFO, false
	case FIFO:
		return FIFO, true
	case WMA:
		WEIGHTED_AVERAGE(account)
		return FIFO, true
	}
	return "", true
}

func INSERT_TO_DATABASE(array_of_journal_tag []JOURNAL_TAG) {
	INSERT_ENTRY_NUMBER(array_of_journal_tag)
	DB_INSERT(DB_JOURNAL, array_of_journal_tag)
	// for _, entry := range array_of_journal_tag {
	// 	cost,err:=COST_FLOW(entry.ACCOUNT_DEBIT, entry.QUANTITY_DEBIT, true)
	// 	if err!=nil{ {
	// 		DB_INSERT_INTO_JOURNAL_OR_INVENTORY(DB_PATH_INVENTORY+entry.ACCOUNT_DEBIT, array_of_journal_tag)
	// 	}
	// }
}

func INSERT_ENTRY_NUMBER(array_of_journal_tag []JOURNAL_TAG) {
	journal_tag := DB_LAST_LINE[JOURNAL_TAG](DB_JOURNAL)
	for indexa := range array_of_journal_tag {
		array_of_journal_tag[indexa].ENTRY_NUMBER = journal_tag.ENTRY_NUMBER + 1
		array_of_journal_tag[indexa].ENTRY_NUMBER_COMPOUND += journal_tag.ENTRY_NUMBER_COMPOUND + 1
		array_of_journal_tag[indexa].ENTRY_NUMBER_SIMPLE = journal_tag.ENTRY_NUMBER_SIMPLE + indexa + 1
	}
}

func STAGE_1(entries []PRICE_QUANTITY_ACCOUNT_BARCODE) {
	var indexa int
	for indexa < len(entries) {
		// set the price to abs
		entries[indexa].PRICE = math.Abs(entries[indexa].PRICE)
		// find account name from barcode
		account_struct, _, err := ACCOUNT_STRUCT_FROM_BARCODE(entries[indexa].BARCODE)
		if err == nil {
			entries[indexa].ACCOUNT_NAME = account_struct.ACCOUNT_NAME
		}
		// search for the account name
		account_struct, _, err = ACCOUNT_STRUCT_FROM_NAME(entries[indexa].ACCOUNT_NAME)
		// delete the entry
		if err != nil || !account_struct.IS_LOW_LEVEL_ACCOUNT || entries[indexa].QUANTITY == 0 {
			entries = POPUP(entries, indexa)
		} else {
			indexa++
		}
	}
}

func CHECK_DEBIT_EQUAL_CREDIT(entries []PRICE_QUANTITY_ACCOUNT_BARCODE) ([]PRICE_QUANTITY_ACCOUNT_BARCODE, []PRICE_QUANTITY_ACCOUNT_BARCODE, error) {
	var debit_entries, credit_entries []PRICE_QUANTITY_ACCOUNT_BARCODE
	var zero float64
	for _, entry := range entries {
		account_struct, _, _ := ACCOUNT_STRUCT_FROM_NAME(entry.ACCOUNT_NAME)
		value := entry.PRICE * entry.QUANTITY
		switch account_struct.IS_CREDIT {
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
		return []PRICE_QUANTITY_ACCOUNT_BARCODE{}, []PRICE_QUANTITY_ACCOUNT_BARCODE{},
			fmt.Errorf("the debit and credit should be equal and the deffrence is %f", zero)
	}
	if len(debit_entries) != 1 && len(credit_entries) != 1 {
		return []PRICE_QUANTITY_ACCOUNT_BARCODE{}, []PRICE_QUANTITY_ACCOUNT_BARCODE{}, ERROR_SHOULD_BE_ONE_DEBIT_OR_ONE_CREDIT
	}
	return debit_entries, credit_entries, nil
}

func SET_SLICE_DAY_START_END(array_day_start_end []DAY_START_END) {
	if len(array_day_start_end) == 0 {
		array_day_start_end = []DAY_START_END{
			{SATURDAY, 0, 0, 23, 59},
			{SUNDAY, 0, 0, 23, 59},
			{MONDAY, 0, 0, 23, 59},
			{TUESDAY, 0, 0, 23, 59},
			{WEDNESDAY, 0, 0, 23, 59},
			{THURSDAY, 0, 0, 23, 59},
			{FRIDAY, 0, 0, 23, 59}}
	}
	for index := range array_day_start_end {
		array_day_start_end[index].DAY = strings.Title(array_day_start_end[index].DAY)

		if !IS_IN(array_day_start_end[index].DAY, STANDARD_DAYS) {
			array_day_start_end[index].DAY = SUNDAY
		}

		if array_day_start_end[index].START_HOUR < 0 {
			array_day_start_end[index].START_HOUR = 0
		}
		if array_day_start_end[index].START_HOUR > 23 {
			array_day_start_end[index].START_HOUR = 23
		}
		if array_day_start_end[index].START_MINUTE < 0 {
			array_day_start_end[index].START_MINUTE = 0
		}
		if array_day_start_end[index].START_MINUTE > 59 {
			array_day_start_end[index].START_MINUTE = 59
		}
		if array_day_start_end[index].END_HOUR < 0 {
			array_day_start_end[index].END_HOUR = 0
		}
		if array_day_start_end[index].END_HOUR > 23 {
			array_day_start_end[index].END_HOUR = 23
		}
		if array_day_start_end[index].END_MINUTE < 0 {
			array_day_start_end[index].END_MINUTE = 0
		}
		if array_day_start_end[index].END_MINUTE > 59 {
			array_day_start_end[index].END_MINUTE = 59
		}

		if array_day_start_end[index].START_HOUR > array_day_start_end[index].END_HOUR {
			array_day_start_end[index].START_HOUR = 0
		}
		if array_day_start_end[index].START_HOUR == array_day_start_end[index].END_HOUR && array_day_start_end[index].START_MINUTE > array_day_start_end[index].END_MINUTE {
			array_day_start_end[index].START_MINUTE = 0
		}
	}
}

func CREATE_ARRAY_START_END_MINUTES(date_end, date_start time.Time, array_day_start_end []DAY_START_END) []start_end_minutes {
	var array_start_end_minutes []start_end_minutes
	var previous_date_end time.Time
	delta_days := int(date_end.Sub(date_start).Hours()/24 + 1)
	year, month, day := date_start.Date()
	for day_counter := 0; day_counter < delta_days; day_counter++ {
		for _, element := range array_day_start_end {
			start := time.Date(year, month, day+day_counter, element.START_HOUR, element.START_MINUTE, 0, 0, time.Local)
			if start.Weekday().String() == element.DAY {
				end := time.Date(year, month, day+day_counter, element.END_HOUR, element.END_MINUTE, 0, 0, time.Local)
				start, end = SHIFT_AND_ARRANGE_THE_TIME_SERIES(previous_date_end, start, end)
				array_start_end_minutes = append(array_start_end_minutes, start_end_minutes{start, end, end.Sub(start).Minutes()})
				previous_date_end = end
			}
		}
	}
	return array_start_end_minutes
}

func SHIFT_AND_ARRANGE_THE_TIME_SERIES(previous_date_end, date_start, date_end time.Time) (time.Time, time.Time) {
	if previous_date_end.After(date_start) {
		date_start = previous_date_end
	}
	if date_start.After(date_end) {
		date_end = date_start
	}
	return date_start, date_end
}

func TOTAL_MINUTES(array_start_end_minutes []start_end_minutes) float64 {
	var total_minutes float64
	for _, element := range array_start_end_minutes {
		total_minutes += element.minutes
	}
	return total_minutes
}

func ADJUST_THE_ARRAY(array_to_insert []JOURNAL_TAG, array_start_end_minutes []start_end_minutes, adjusting_method string) [][]JOURNAL_TAG {
	var adjusted_array_to_insert [][]JOURNAL_TAG
	total_minutes := TOTAL_MINUTES(array_start_end_minutes)
	array_len_start_end_minutes := len(array_start_end_minutes) - 1
	for _, entry := range array_to_insert {
		var value_counter, time_unit_counter float64
		var one_account_adjusted_list []JOURNAL_TAG
		for index, element := range array_start_end_minutes {
			value := VALUE_AFTER_ADJUST_USING_ADJUSTING_METHODS(adjusting_method, element.minutes, total_minutes, time_unit_counter, entry.VALUE)

			if index == array_len_start_end_minutes {
				value = entry.VALUE - value_counter
			}

			time_unit_counter += element.minutes
			value_counter += value
			one_account_adjusted_list = append(one_account_adjusted_list, JOURNAL_TAG{
				REVERSE:               false,
				ENTRY_NUMBER:          0,
				ENTRY_NUMBER_COMPOUND: index,
				ENTRY_NUMBER_SIMPLE:   0,
				VALUE:                 value,
				PRICE_DEBIT:           entry.PRICE_DEBIT,
				PRICE_CREDIT:          entry.PRICE_CREDIT,
				QUANTITY_DEBIT:        RETURN_SAME_SIGN_OF_NUMBER_SIGN(value/entry.PRICE_DEBIT, entry.QUANTITY_DEBIT),
				QUANTITY_CREDIT:       RETURN_SAME_SIGN_OF_NUMBER_SIGN(value/entry.PRICE_CREDIT, entry.QUANTITY_CREDIT),
				ACCOUNT_DEBIT:         entry.ACCOUNT_DEBIT,
				ACCOUNT_CREDIT:        entry.ACCOUNT_CREDIT,
				NOTES:                 entry.NOTES,
				NAME:                  entry.NAME,
				NAME_EMPLOYEE:         entry.NAME_EMPLOYEE,
				DATE_START:            element.date_start,
				DATE_END:              element.date_end,
			})
		}
		adjusted_array_to_insert = append(adjusted_array_to_insert, one_account_adjusted_list)
	}
	return adjusted_array_to_insert
}

func VALUE_AFTER_ADJUST_USING_ADJUSTING_METHODS(adjusting_method string, minutes, TOTAL_MINUTES, time_unit_counter, total_value float64) float64 {
	percent := ROOT(total_value, TOTAL_MINUTES)
	switch adjusting_method {
	case EXPONENTIAL:
		return math.Pow(percent, time_unit_counter+minutes) - math.Pow(percent, time_unit_counter)
	case LOGARITHMIC:
		return (total_value / math.Pow(percent, time_unit_counter)) - (total_value / math.Pow(percent, time_unit_counter+minutes))
	default:
		return minutes * (total_value / TOTAL_MINUTES)
	}
}

func JOURNAL_ENTRY(
	entries []PRICE_QUANTITY_ACCOUNT_BARCODE,
	insert, auto_completion, invoice_discount bool,
	date_start, date_end time.Time,
	adjusting_method, notes, name, name_employee string,
	array_day_start_end []DAY_START_END) []JOURNAL_TAG {

	STAGE_1(entries)
	SET_SLICE_DAY_START_END(array_day_start_end)
	date_end = SET_DATE_END_TO_ZERO_IF_SMALLER_THAN_DATE_START(date_start, date_end)
	adjusting_method = SET_ADJUSTING_METHOD(date_end, adjusting_method, entries)
	entries = GROUP_BY_ACCOUNT(entries)
	FIND_COST(entries)

	// if auto_completion {
	// 	entries = auto_completion_the_entry(entries)
	// }
	// if invoice_discount {
	// 	entries = auto_completion_the_invoice_discount(entries)
	// }

	entries = GROUP_BY_ACCOUNT(entries)

	INCREASE_THE_VALUE_TO_MAKE_THE_NEW_BALANCE_FOR_THE_ACCOUNT_POSITIVE(entries)

	debit_entries, credit_entries, _ := CHECK_DEBIT_EQUAL_CREDIT(entries)
	simple_entries := INSERT_TO_JOURNAL_TAG(debit_entries, credit_entries, date_start, date_end, notes, name, name_employee)

	if IS_IN(adjusting_method, DEPRECIATION_METHODS) {
		array_start_end_minutes := CREATE_ARRAY_START_END_MINUTES(date_end, date_start, array_day_start_end)
		adjusted_array_to_insert := ADJUST_THE_ARRAY(simple_entries, array_start_end_minutes, adjusting_method)
		adjusted_array_to_insert = TRANSPOSE(adjusted_array_to_insert)
		simple_entries = UNPACK(adjusted_array_to_insert)
	}

	INSERT_TO_DATABASE(simple_entries)
	return simple_entries
}
