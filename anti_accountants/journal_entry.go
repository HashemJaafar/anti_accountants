package anti_accountants

import (
	"errors"
	"fmt"
	"math"
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

func SET_ADJUSTING_METHOD(entry_expair time.Time, adjusting_method string, entries []VALUE_PRICE_QUANTITY_ACCOUNT_BARCODE) string {
	if !IS_IN(adjusting_method, ADJUSTING_METHODS) {
		return ""
	}
	if entry_expair.IsZero() {
		return ""
	}
	is_in_depreciation_methods := IS_IN(adjusting_method, DEPRECIATION_METHODS)
	for _, entry := range entries {
		account_struct, _, _ := ACCOUNT_STRUCT_FROM_NAME(entry.ACCOUNT_NAME)
		if account_struct.COST_FLOW_TYPE != "" && is_in_depreciation_methods {
			return ""
		}
	}
	return adjusting_method
}

func GROUP_BY_ACCOUNT_AND_BARCODE(entries []VALUE_PRICE_QUANTITY_ACCOUNT_BARCODE) []VALUE_PRICE_QUANTITY_ACCOUNT_BARCODE {
	g := map[string]*VALUE_PRICE_QUANTITY_ACCOUNT_BARCODE{}
	for _, v := range entries {
		key := v.ACCOUNT_NAME
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
			ACCOUNT_NAME: key,
			BARCODE:      v.BARCODE,
		})
	}
	return entries
}

func REMOVE_ZERO_VALUES(entries []VALUE_PRICE_QUANTITY_ACCOUNT_BARCODE) {
	var index int
	for index < len(entries) {
		if entries[index].VALUE == 0 || entries[index].QUANTITY == 0 {
			entries = POPUP(entries, index)
		} else {
			index++
		}
	}
}

func FIND_COST(entries []VALUE_PRICE_QUANTITY_ACCOUNT_BARCODE) {
	for index, entry := range entries {
		costs, err := COST_FLOW(entry.ACCOUNT_NAME, entry.QUANTITY, false)
		if err == nil {
			entries[index].VALUE = -costs
			entries[index].PRICE = -costs / entry.QUANTITY
		}
	}
}

func INSERT_TO_JOURNAL_TAG(debit_entries, credit_entries []VALUE_PRICE_QUANTITY_ACCOUNT_BARCODE, date_start, date_end time.Time, notes, name, name_employee string) []JOURNAL_TAG {
	var simple_entries []JOURNAL_TAG
	for _, debit_entry := range debit_entries {
		for _, credit_entry := range credit_entries {
			simple_entries = append(simple_entries, JOURNAL_TAG{
				VALUE:           math.Abs(SMALLEST(debit_entry.VALUE, credit_entry.VALUE)),
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

func INCREASE_THE_VALUE_TO_MAKE_THE_NEW_BALANCE_FOR_THE_ACCOUNT_POSITIVE(entries []VALUE_PRICE_QUANTITY_ACCOUNT_BARCODE) {
	for indexa, a := range entries {
		ACCOUNT_BALANCE := ACCOUNT_BALANCE(a.ACCOUNT_NAME)
		new_balance := ACCOUNT_BALANCE + a.VALUE
		if new_balance < 0 {
			entries[indexa].VALUE -= new_balance
			entries[indexa].QUANTITY = entries[indexa].VALUE / entries[indexa].PRICE
		}
	}
}

func FIND_ACCOUNT_FROM_BARCODE(entries []VALUE_PRICE_QUANTITY_ACCOUNT_BARCODE) {
	for indexa, a := range entries {
		account_struct, _, err := ACCOUNT_STRUCT_FROM_BARCODE(a.BARCODE)
		if err == nil {
			entries[indexa].ACCOUNT_NAME = account_struct.ACCOUNT_NAME
		}
	}
}

func SET_THE_SIGN_OF_THE_VALUE_SAME_SIGN_OF_QUANTITY(entries []VALUE_PRICE_QUANTITY_ACCOUNT_BARCODE) {
	for index := range entries {
		entries[index].VALUE = RETURN_SAME_SIGN_OF_NUMBER_SIGN(entries[index].QUANTITY, entries[index].VALUE)
	}
}

func COST_FLOW(account string, quantity float64, insert bool) (float64, error) {
	if quantity > 0 {
		return 0, ERROR_SHOULD_BE_NEGATIVE
	}
	IS_ASCENDING, err := IS_ASCENDING(account)
	if err != nil {
		return 0, ERROR_NOT_INVENTORY_ACCOUNT
	}
	inventory := DB_READ_INVENTORY(account)
	SORT_BY_TIME_INVENTORY(inventory, IS_ASCENDING)
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

func IS_ASCENDING(account string) (bool, error) {
	account_struct, _, _ := ACCOUNT_STRUCT_FROM_NAME(account)
	switch account_struct.COST_FLOW_TYPE {
	case LIFO:
		return false, nil
	case FIFO:
		return true, nil
	case WMA:
		WEIGHTED_AVERAGE(account)
		return true, nil
	}
	return false, errors.New("is not inventory account")
}

func INSERT_TO_DATABASE(array_of_journal_tag []JOURNAL_TAG) {
	INSERT_ENTRY_NUMBER(array_of_journal_tag)
	DB_INSERT_INTO_JOURNAL_OR_INVENTORY(DB_JOURNAL, array_of_journal_tag)
	// for _, entry := range array_of_journal_tag {
	// 	cost,err:=COST_FLOW(entry.ACCOUNT_DEBIT, entry.QUANTITY_DEBIT, true)
	// 	if err!=nil{ {
	// 		DB_INSERT_INTO_JOURNAL_OR_INVENTORY(DB_INVENTORY+entry.ACCOUNT_DEBIT, array_of_journal_tag)
	// 	}
	// }
}

func INSERT_ENTRY_NUMBER(array_of_journal_tag []JOURNAL_TAG) {
	journal_tag := DB_LAST_LINE_IN_JOURNAL()
	for indexa := range array_of_journal_tag {
		array_of_journal_tag[indexa].ENTRY_NUMBER = journal_tag.ENTRY_NUMBER + 1
		array_of_journal_tag[indexa].ENTRY_NUMBER_COMPOUND += journal_tag.ENTRY_NUMBER_COMPOUND + 1
		array_of_journal_tag[indexa].ENTRY_NUMBER_SIMPLE = journal_tag.ENTRY_NUMBER_SIMPLE + indexa + 1
	}
}

func CALCULATE_AND_INSERT_VALUE_PRICE_QUANTITY(entries []VALUE_PRICE_QUANTITY_ACCOUNT_BARCODE) {
	for index, entry := range entries {
		m := map[string]float64{}
		INSERT_IF_NOT_ZERO(m, "VALUE", entry.VALUE)
		INSERT_IF_NOT_ZERO(m, "PRICE", entry.PRICE)
		INSERT_IF_NOT_ZERO(m, "QUANTITY", entry.QUANTITY)
		EQUATION_SOLVER(false, m, "VALUE", "PRICE", "*", "QUANTITY")
		entries[index].VALUE = m["VALUE"]
		entries[index].PRICE = m["PRICE"]
		entries[index].QUANTITY = m["QUANTITY"]
	}
}

func INSERT_IF_NOT_ZERO(m map[string]float64, str string, number float64) {
	if number != 0 {
		m[str] = number
	}
}

func REMOVE_HIGH_LEVEL_ACCOUNT(entries []VALUE_PRICE_QUANTITY_ACCOUNT_BARCODE) {
	var indexa int
	for indexa < len(entries) {
		account_struct, _, _ := ACCOUNT_STRUCT_FROM_NAME(entries[indexa].ACCOUNT_NAME)
		if !account_struct.IS_LOW_LEVEL_ACCOUNT {
			entries = POPUP(entries, indexa)
		} else {
			indexa++
		}
	}
}

func REMOVE_THE_ACCOUNTS_NOT_IN_ACCOUNTS_LIST(entries []VALUE_PRICE_QUANTITY_ACCOUNT_BARCODE) {
	var indexa int
	for indexa < len(entries) {
		_, _, err := ACCOUNT_STRUCT_FROM_NAME(entries[indexa].ACCOUNT_NAME)
		if err != nil {
			entries = POPUP(entries, indexa)
		} else {
			indexa++
		}
	}
}

func CHECK_DEBIT_EQUAL_CREDIT(entries []VALUE_PRICE_QUANTITY_ACCOUNT_BARCODE) ([]VALUE_PRICE_QUANTITY_ACCOUNT_BARCODE, []VALUE_PRICE_QUANTITY_ACCOUNT_BARCODE, error) {
	var debit_entries, credit_entries []VALUE_PRICE_QUANTITY_ACCOUNT_BARCODE
	var zero float64
	for _, entry := range entries {
		account_struct, _, _ := ACCOUNT_STRUCT_FROM_NAME(entry.ACCOUNT_NAME)
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
		return []VALUE_PRICE_QUANTITY_ACCOUNT_BARCODE{}, []VALUE_PRICE_QUANTITY_ACCOUNT_BARCODE{}, ERROR_SHOULD_BE_ONE_DEBIT_OR_ONE_CREDIT
	}
	return debit_entries, credit_entries, nil
}

func JOURNAL_ENTRY(
	entries []VALUE_PRICE_QUANTITY_ACCOUNT_BARCODE,
	insert, auto_completion, invoice_discount bool,
	date_start, date_end time.Time,
	adjusting_method, notes, name, name_employee string,
	array_day_start_end []DAY_START_END) []JOURNAL_TAG {

	FIND_ACCOUNT_FROM_BARCODE(entries)
	REMOVE_THE_ACCOUNTS_NOT_IN_ACCOUNTS_LIST(entries)
	SET_SLICE_DAY_START_END(array_day_start_end)
	date_end = SET_DATE_END_TO_ZERO_IF_SMALLER_THAN_DATE_START(date_start, date_end)
	adjusting_method = SET_ADJUSTING_METHOD(date_end, adjusting_method, entries)
	entries = GROUP_BY_ACCOUNT_AND_BARCODE(entries)
	REMOVE_ZERO_VALUES(entries)
	SET_THE_SIGN_OF_THE_VALUE_SAME_SIGN_OF_QUANTITY(entries)
	CALCULATE_AND_INSERT_VALUE_PRICE_QUANTITY(entries)
	FIND_COST(entries)

	// if auto_completion {
	// 	entries = auto_completion_the_entry(entries)
	// }
	// if invoice_discount {
	// 	entries = auto_completion_the_invoice_discount(entries)
	// }

	entries = GROUP_BY_ACCOUNT_AND_BARCODE(entries)
	REMOVE_ZERO_VALUES(entries)

	REMOVE_HIGH_LEVEL_ACCOUNT(entries)
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
