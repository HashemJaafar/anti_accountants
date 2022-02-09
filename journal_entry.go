package anti_accountants

import (
	"fmt"
	"log"
	"math"
	"strings"
	"time"
)

type ACCOUNT_VALUE_PRICE_QUANTITY_BARCODE struct {
	ACCOUNT  string
	VALUE    float64
	PRICE    float64
	QUANTITY float64
	BARCODE  string
}

type DAY_START_END struct {
	DAY          string
	START_HOUR   int
	START_MINUTE int
	END_HOUR     int
	END_MINUTE   int
}

type AUTO_COMPLETE_ENTRIE struct {
	ACCOUNT_0  string
	IS_PERCENT bool
	VALUE      float64
	ACCOUNT_1  string
	QUANTITY_1 float64
	ACCOUNT_2  string
	QUANTITY_2 float64
}

type day_start_end_date_minutes struct {
	day        string
	start_date time.Time
	end_date   time.Time
	minutes    float64
}

var (
	standard_days        = [7]string{"Saturday", "Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday"}
	adjusting_methods    = [4]string{"linear", "exponential", "logarithmic", "expire"}
	depreciation_methods = [3]string{"linear", "exponential", "logarithmic"}
	NOW                  = time.Now()
)

func check_the_params(entry_expair time.Time, adjusting_method string, date time.Time, entries []JOURNAL_TAG, array_day_start_end []DAY_START_END) []DAY_START_END {
	if !(IS_IN(adjusting_method, adjusting_methods[:]) || adjusting_method == "") {
		error_element_is_not_in_elements(adjusting_method, adjusting_methods[:])
	}
	if entry_expair.IsZero() == IS_IN(adjusting_method, adjusting_methods[:]) {
		error_you_cant_use_entry_expire()
	}
	if !entry_expair.IsZero() {
		check_dates(date, entry_expair)
	}
	for _, entry := range entries {
		if IS_IN(entry.ACCOUNT, inventory) && !IS_IN(adjusting_method, []string{"expire", ""}) {
			error_you_cant_use_depreciation_methods_with_inventory(entry.ACCOUNT)
		}
	}
	if IS_IN(adjusting_method, depreciation_methods[:]) {
		if len(array_day_start_end) == 0 {
			array_day_start_end = []DAY_START_END{
				{"saturday", 0, 0, 23, 59},
				{"sunday", 0, 0, 23, 59},
				{"monday", 0, 0, 23, 59},
				{"tuesday", 0, 0, 23, 59},
				{"wednesday", 0, 0, 23, 59},
				{"thursday", 0, 0, 23, 59},
				{"friday", 0, 0, 23, 59}}
		}
		for index, element := range array_day_start_end {
			array_day_start_end[index].DAY = strings.Title(element.DAY)
			switch {
			case !IS_IN(array_day_start_end[index].DAY, standard_days[:]):
				error_element_is_not_in_elements(element.DAY, standard_days[:])
			case element.START_HOUR < 0:
				error_the_time_is_not_in_range(element, 0)
			case element.START_HOUR > 23:
				error_the_time_is_not_in_range(element, 23)
			case element.START_MINUTE < 0:
				error_the_time_is_not_in_range(element, 0)
			case element.START_MINUTE > 59:
				error_the_time_is_not_in_range(element, 59)
			case element.END_HOUR < 0:
				error_the_time_is_not_in_range(element, 0)
			case element.END_HOUR > 23:
				error_the_time_is_not_in_range(element, 23)
			case element.END_MINUTE < 0:
				error_the_time_is_not_in_range(element, 0)
			case element.END_MINUTE > 59:
				error_the_time_is_not_in_range(element, 59)
			}
		}
	}
	return array_day_start_end
}

func group_by_account_and_barcode(entries []JOURNAL_TAG) []JOURNAL_TAG {
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
	entries = []JOURNAL_TAG{}
	for key, v := range g {
		entries = append(entries, JOURNAL_TAG{
			ACCOUNT:  key.account,
			VALUE:    v.VALUE,
			PRICE:    v.VALUE / v.QUANTITY,
			QUANTITY: v.QUANTITY,
			BARCODE:  key.barcode,
		}) //ACCOUNT_VALUE_PRICE_QUANTITY_BARCODE{key.account, v.VALUE, v.VALUE / v.QUANTITY, v.QUANTITY, key.barcode})
	}
	return entries
}

func remove_zero_values(entries []JOURNAL_TAG) {
	var index int
	for index < len(entries) {
		if entries[index].VALUE == 0 || entries[index].QUANTITY == 0 {
			entries = append(entries[:index], entries[index+1:]...)
		} else {
			index++
		}
	}
}

func (s FINANCIAL_ACCOUNTING) auto_completion_the_entry(entries []JOURNAL_TAG, auto_completion bool) []JOURNAL_TAG {
	var new_entries [][]JOURNAL_TAG
	if auto_completion {
		for _, entry := range entries {
			for _, complement := range s.AUTO_COMPLETE_ENTRIES {
				if complement.ACCOUNT_0 == entry.ACCOUNT {
					// var barcode_1, barcode_2 string
					// if complement.ACCOUNT_0 == complement.ACCOUNT_1 {
					// 	barcode_1 = entry.BARCODE
					// }
					// if complement.ACCOUNT_0 == complement.ACCOUNT_2 {
					// 	barcode_2 = entry.BARCODE
					// }
					// switch complement.IS_PERCENT {
					// case true:
					// 	new_entries = append(new_entries, []JOURNAL_TAG{
					// 		{complement.ACCOUNT_1, complement.NUMBER * entry.VALUE, complement.NUMBER * entry.QUANTITY, barcode_1},
					// 		{complement.ACCOUNT_2, complement.NUMBER * entry.VALUE, complement.NUMBER * entry.QUANTITY, barcode_2},
					// 	})
					// case false:
					// 	new_entries = append(new_entries, []JOURNAL_TAG{
					// 		{complement.ACCOUNT_1, complement.NUMBER, complement.NUMBER / complement.PRICE_1, barcode_1},
					// 		{complement.ACCOUNT_2, complement.NUMBER, complement.NUMBER / complement.PRICE_2, barcode_2},
					// 	})
					// }
				}
			}
		}
		fmt.Println(new_entries)
	}
	return entries
}

func (s FINANCIAL_ACCOUNTING) find_cost(entries []JOURNAL_TAG) {
	for index, entry := range entries {
		costs := s.cost_flow(entry.ACCOUNT, entry.QUANTITY, entry.BARCODE, false)
		if costs != 0 {
			entries[index].VALUE = -costs
			entries[index].PRICE = -costs / entry.QUANTITY
		}
	}
}

func (s FINANCIAL_ACCOUNTING) auto_completion_the_invoice_discount(auto_completion bool, entries []JOURNAL_TAG) []JOURNAL_TAG {
	if auto_completion {
		total_invoice_before_invoice_discount := s.total_invoice_before_invoice_discount(entries)
		_, discount := X_UNDER_X(s.INVOICE_DISCOUNTS_LIST, total_invoice_before_invoice_discount)
		invoice_discount := discount_tax_calculator(total_invoice_before_invoice_discount, discount)
		entries = append(entries, JOURNAL_TAG{
			ACCOUNT:  s.INVOICE_DISCOUNT,
			VALUE:    invoice_discount,
			PRICE:    invoice_discount,
			QUANTITY: 1,
		})
	}
	return entries
}

func (s FINANCIAL_ACCOUNTING) total_invoice_before_invoice_discount(entries []JOURNAL_TAG) float64 {
	var total_invoice_before_invoice_discount float64
	for _, entry := range entries {
		if s.is_it_sub_account_using_name(s.INCOME_STATEMENT, entry.ACCOUNT) && s.is_credit(entry.ACCOUNT) {
			total_invoice_before_invoice_discount += entry.VALUE
		} else if s.is_it_sub_account_using_name(s.DISCOUNTS, entry.ACCOUNT) && !s.is_credit(entry.ACCOUNT) {
			total_invoice_before_invoice_discount -= entry.VALUE
		}
	}
	return total_invoice_before_invoice_discount
}

func discount_tax_calculator(price, discount_tax float64) float64 {
	if discount_tax < 0 {
		discount_tax = math.Abs(discount_tax)
	} else if discount_tax > 0 {
		discount_tax = price * discount_tax
	}
	return discount_tax
}

func (s FINANCIAL_ACCOUNTING) convert_to_simple_entry(debit_entries, credit_entries []JOURNAL_TAG) [][]JOURNAL_TAG {
	simple_entries := [][]JOURNAL_TAG{}
	for _, debit_entry := range debit_entries {
		for _, credit_entry := range credit_entries {
			simple_entries = append(simple_entries, []JOURNAL_TAG{debit_entry, credit_entry})
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

func (s FINANCIAL_ACCOUNTING) can_the_account_be_negative(entries []JOURNAL_TAG) {
	for _, entry := range entries {
		if !(s.is_it_sub_account_using_name(s.EQUITY, entry.ACCOUNT) && s.is_credit(entry.ACCOUNT)) {
			account_balance := account_balance(entry.ACCOUNT)
			if account_balance+entry.VALUE < 0 {
				error_make_nagtive_balance(entry, account_balance)
			}
		}
	}
}

func find_account_from_barcode(entries []JOURNAL_TAG) {
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
		if entry.PRICE < 0 {
			error_the_price_should_be_positive(entry)
		}
		array_to_insert = append(array_to_insert, JOURNAL_TAG{
			DATE:          date.String(),
			ENTRY_NUMBER:  0,
			ACCOUNT:       entry.ACCOUNT,
			VALUE:         entry.VALUE,
			PRICE:         entry.PRICE,
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

func adjuste_the_array(entry_expair time.Time, date time.Time, array_day_start_end []DAY_START_END, array_to_insert []JOURNAL_TAG, adjusting_method string, description string, name string, employee_name string) [][]JOURNAL_TAG {
	var day_start_end_date_minutes_array []day_start_end_date_minutes
	var total_minutes float64
	var previous_end_date, end time.Time
	delta_days := int(entry_expair.Sub(date).Hours()/24 + 1)
	year, month_sting, day := date.Date()
	for day_counter := 0; day_counter < delta_days; day_counter++ {
		for _, element := range array_day_start_end {
			if start := time.Date(year, month_sting, day+day_counter, element.START_HOUR, element.START_MINUTE, 0, 0, time.Local); start.Weekday().String() == element.DAY {
				previous_end_date = end
				end = time.Date(year, month_sting, day+day_counter, element.END_HOUR, element.END_MINUTE, 0, 0, time.Local)
				if start.After(end) {
					log.Panic("the start_hour and start_minute should be smaller than END_HOUR and end_minute for ", element)
				}
				if previous_end_date.After(start) {
					log.Panic("the END_HOUR and end_minute for ", element.DAY, " should be smaller than start_hour and start_minute for the second ", element)
				}
				minutes := end.Sub(start).Minutes()
				total_minutes += minutes
				day_start_end_date_minutes_array = append(day_start_end_date_minutes_array, day_start_end_date_minutes{element.DAY, start, end, minutes})
			}
		}
	}
	var adjusted_array_to_insert [][]JOURNAL_TAG
	for _, entry := range array_to_insert {
		var value_counter, time_unit_counter float64
		var one_account_adjusted_list []JOURNAL_TAG
		total_value := math.Abs(entry.VALUE)
		for index, element := range day_start_end_date_minutes_array {
			value := value_after_adjust_using_adjusting_methods(adjusting_method, element, total_minutes, time_unit_counter, total_value)

			if index >= delta_days-1 {
				value = math.Abs(total_value - value_counter)
			}

			time_unit_counter += element.minutes
			value_counter += math.Abs(value)
			value = RETURN_SAME_SIGN_OF_NUMBER_SIGN(entry.VALUE, value)
			one_account_adjusted_list = append(one_account_adjusted_list, JOURNAL_TAG{
				DATE:          element.start_date.String(),
				ENTRY_NUMBER:  0,
				ACCOUNT:       entry.ACCOUNT,
				VALUE:         value,
				PRICE:         entry.PRICE,
				QUANTITY:      value / entry.PRICE,
				BARCODE:       entry.BARCODE,
				ENTRY_EXPAIR:  element.end_date.String(),
				DESCRIPTION:   description,
				NAME:          name,
				EMPLOYEE_NAME: employee_name,
				ENTRY_DATE:    NOW.String(),
				REVERSE:       false,
			})
		}
		adjusted_array_to_insert = append(adjusted_array_to_insert, one_account_adjusted_list)
	}
	return adjusted_array_to_insert
}

func value_after_adjust_using_adjusting_methods(adjusting_method string, element day_start_end_date_minutes, total_minutes, time_unit_counter, total_value float64) float64 {
	percent := ROOT(total_value, total_minutes)
	switch adjusting_method {
	case "linear":
		return element.minutes * (total_value / total_minutes)
	case "exponential":
		return math.Pow(percent, time_unit_counter+element.minutes) - math.Pow(percent, time_unit_counter)
	case "logarithmic":
		return (total_value / math.Pow(percent, time_unit_counter)) - (total_value / math.Pow(percent, time_unit_counter+element.minutes))
	}
	return 0
}

func (s FINANCIAL_ACCOUNTING) cost_flow(account string, quantity float64, barcode string, insert bool) float64 {
	if quantity > 0 {
		return 0
	}
	order_by_date_asc_or_desc := s.asc_or_desc(account)
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

func (s FINANCIAL_ACCOUNTING) asc_or_desc(account string) string {
	switch s.return_cost_flow_type(account) {
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

func (s FINANCIAL_ACCOUNTING) insert_to_database(array_of_journal_tag []JOURNAL_TAG, insert_into_journal, insert_into_inventory bool) {
	insert_entry_number(array_of_journal_tag)
	if insert_into_journal {
		insert_into_journal_func(array_of_journal_tag)
	}
	if insert_into_inventory {
		s.insert_into_inventory(array_of_journal_tag)
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

func (s FINANCIAL_ACCOUNTING) is_the_account_is_it_sub_account_using_name(entries []JOURNAL_TAG) {
	for _, entry := range entries {
		if s.is_it_high_by_level(s.account_number(entry.ACCOUNT)) {
			error_is_high_level_account(entry.ACCOUNT)
		}
	}
}

func (s FINANCIAL_ACCOUNTING) JOURNAL_ENTRY(entries []ACCOUNT_VALUE_PRICE_QUANTITY_BARCODE, insert, auto_completion bool, date time.Time, entry_expair time.Time, adjusting_method string,
	description string, name string, employee_name string, array_day_start_end []DAY_START_END) []JOURNAL_TAG {
	calculate_and_insert_value_price_quantity(entries)
	JOURNAL_TAG_entries := insert_to_JOURNAL_TAG(entries, date, entry_expair, description, name, employee_name)
	array_day_start_end = check_the_params(entry_expair, adjusting_method, date, JOURNAL_TAG_entries, array_day_start_end)
	JOURNAL_TAG_entries = group_by_account_and_barcode(JOURNAL_TAG_entries)
	remove_zero_values(JOURNAL_TAG_entries)
	find_account_from_barcode(JOURNAL_TAG_entries)
	s.find_cost(JOURNAL_TAG_entries)
	JOURNAL_TAG_entries = s.auto_completion_the_entry(JOURNAL_TAG_entries, auto_completion)
	// JOURNAL_TAG_entries = s.auto_completion_the_invoice_discount(auto_completion, JOURNAL_TAG_entries)
	JOURNAL_TAG_entries = group_by_account_and_barcode(JOURNAL_TAG_entries)
	remove_zero_values(JOURNAL_TAG_entries)
	s.is_the_account_is_it_sub_account_using_name(JOURNAL_TAG_entries)
	s.can_the_account_be_negative(JOURNAL_TAG_entries)
	debit_entries, credit_entries := s.check_debit_equal_credit(JOURNAL_TAG_entries, false)
	simple_entries := s.convert_to_simple_entry(debit_entries, credit_entries)
	var all_array_to_insert []JOURNAL_TAG
	for _, simple_entry := range simple_entries {
		array_to_insert := insert_to_JOURNAL_TAG(make_ACCOUNT_VALUE_PRICE_QUANTITY_BARCODE_from_JOURNAL_TAG(simple_entry), date, entry_expair, description, name, employee_name)
		if IS_IN(adjusting_method, depreciation_methods[:]) {
			adjusted_array_to_insert := adjuste_the_array(entry_expair, date, array_day_start_end, array_to_insert, adjusting_method, description, name, employee_name)
			adjusted_array_to_insert = transpose(adjusted_array_to_insert)
			array_to_insert = unpack_the_array(adjusted_array_to_insert)
		}
		all_array_to_insert = append(all_array_to_insert, array_to_insert...)
	}
	s.insert_to_database(all_array_to_insert, insert, insert)
	return all_array_to_insert
}
