package anti_accountants

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type ACCOUNT_VALUE_QUANTITY_BARCODE struct {
	ACCOUNT  string
	VALUE    float64
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

func check_the_params(entry_expair time.Time, adjusting_method string, date time.Time, array_of_entry []ACCOUNT_VALUE_QUANTITY_BARCODE, array_day_start_end []DAY_START_END) []DAY_START_END {
	if entry_expair.IsZero() == IS_IN(adjusting_method, adjusting_methods[:]) {
		log.Panic("check entry_expair => ", entry_expair, " and adjusting_method => ", adjusting_method, " should be in ", adjusting_methods)
	}
	if !entry_expair.IsZero() {
		check_dates(date, entry_expair)
	}
	for _, entry := range array_of_entry {
		if IS_IN(entry.ACCOUNT, inventory) && !IS_IN(adjusting_method, []string{"expire", ""}) {
			log.Panic(entry.ACCOUNT + " is in inventory you just can use expire or make it empty")
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
				log.Panic("error ", element.DAY, " for ", element, " is not in ", standard_days)
			case element.START_HOUR < 0:
				log.Panic("error ", element.START_HOUR, " for ", element, " is < 0")
			case element.START_HOUR > 23:
				log.Panic("error ", element.START_HOUR, " for ", element, " is > 23")
			case element.START_MINUTE < 0:
				log.Panic("error ", element.START_MINUTE, " for ", element, " is < 0")
			case element.START_MINUTE > 59:
				log.Panic("error ", element.START_MINUTE, " for ", element, " is > 59")
			case element.END_HOUR < 0:
				log.Panic("error ", element.END_HOUR, " for ", element, " is < 0")
			case element.END_HOUR > 23:
				log.Panic("error ", element.END_HOUR, " for ", element, " is > 23")
			case element.END_MINUTE < 0:
				log.Panic("error ", element.END_MINUTE, " for ", element, " is < 0")
			case element.END_MINUTE > 59:
				log.Panic("error ", element.END_MINUTE, " for ", element, " is > 59")
			}
		}
	}
	return array_day_start_end
}

func group_by_account_and_barcode(array_of_entry []ACCOUNT_VALUE_QUANTITY_BARCODE) []ACCOUNT_VALUE_QUANTITY_BARCODE {
	type account_barcode struct {
		account, barcode string
	}
	g := map[account_barcode]*ACCOUNT_VALUE_QUANTITY_BARCODE{}
	for _, v := range array_of_entry {
		key := account_barcode{v.ACCOUNT, v.BARCODE}
		sums := g[key]
		if sums == nil {
			sums = &ACCOUNT_VALUE_QUANTITY_BARCODE{}
			g[key] = sums
		}
		sums.VALUE += v.VALUE
		sums.QUANTITY += v.QUANTITY
	}
	array_of_entry = []ACCOUNT_VALUE_QUANTITY_BARCODE{}
	for key, v := range g {
		array_of_entry = append(array_of_entry, ACCOUNT_VALUE_QUANTITY_BARCODE{key.account, v.VALUE, v.QUANTITY, key.barcode})
	}
	return array_of_entry
}

func remove_zero_values(array_of_entry []ACCOUNT_VALUE_QUANTITY_BARCODE) []ACCOUNT_VALUE_QUANTITY_BARCODE {
	var index int
	for index < len(array_of_entry) {
		if array_of_entry[index].VALUE == 0 || array_of_entry[index].QUANTITY == 0 {
			// fmt.Println(array_of_entry[index], " is removed because one of the values is 0")
			array_of_entry = append(array_of_entry[:index], array_of_entry[index+1:]...)
		} else {
			index++
		}
	}
	return array_of_entry
}

func find_account_from_barcode(array_of_entry []ACCOUNT_VALUE_QUANTITY_BARCODE) {
	for index, entry := range array_of_entry {
		if entry.ACCOUNT == "" && entry.BARCODE == "" {
			log.Panic("can't find the account name if the barcode is empty in ", entry)
		}
		if entry.ACCOUNT == "" {
			err := db.QueryRow("select account from journal where barcode=? limit 1", entry.BARCODE).Scan(&array_of_entry[index].ACCOUNT)
			if err != nil {
				log.Panic("the barcode is wrong for ", entry)
			}
		}
	}
}

func (s FINANCIAL_ACCOUNTING) auto_completion_the_entry(array_of_entry []ACCOUNT_VALUE_QUANTITY_BARCODE, auto_completion bool) []ACCOUNT_VALUE_QUANTITY_BARCODE {
	var new_array_of_entry [][]ACCOUNT_VALUE_QUANTITY_BARCODE
	if auto_completion {
		for _, entry := range array_of_entry {
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
					// 	new_array_of_entry = append(new_array_of_entry, []ACCOUNT_VALUE_QUANTITY_BARCODE{
					// 		{complement.ACCOUNT_1, complement.NUMBER * entry.VALUE, complement.NUMBER * entry.QUANTITY, barcode_1},
					// 		{complement.ACCOUNT_2, complement.NUMBER * entry.VALUE, complement.NUMBER * entry.QUANTITY, barcode_2},
					// 	})
					// case false:
					// 	new_array_of_entry = append(new_array_of_entry, []ACCOUNT_VALUE_QUANTITY_BARCODE{
					// 		{complement.ACCOUNT_1, complement.NUMBER, complement.NUMBER / complement.PRICE_1, barcode_1},
					// 		{complement.ACCOUNT_2, complement.NUMBER, complement.NUMBER / complement.PRICE_2, barcode_2},
					// 	})
					// }
				}
			}
		}
	}
	fmt.Println(new_array_of_entry)
	return array_of_entry
}

func (s FINANCIAL_ACCOUNTING) find_cost(array_of_entry []ACCOUNT_VALUE_QUANTITY_BARCODE) {
	for index, entry := range array_of_entry {
		costs := s.cost_flow(entry.ACCOUNT, entry.QUANTITY, entry.BARCODE, false)
		if costs != 0 {
			array_of_entry[index].VALUE = -costs
		}
	}
}

func (s FINANCIAL_ACCOUNTING) auto_completion_the_invoice_discount(auto_completion bool, array_of_entry []ACCOUNT_VALUE_QUANTITY_BARCODE) []ACCOUNT_VALUE_QUANTITY_BARCODE {
	if auto_completion {
		total_invoice_before_invoice_discount := s.total_invoice_before_invoice_discount(array_of_entry)
		_, discount := X_UNDER_X(s.INVOICE_DISCOUNTS_LIST, total_invoice_before_invoice_discount)
		invoice_discount := discount_tax_calculator(total_invoice_before_invoice_discount, discount)
		array_of_entry = append(array_of_entry, ACCOUNT_VALUE_QUANTITY_BARCODE{s.INVOICE_DISCOUNT, invoice_discount, 1, ""})
	}
	return array_of_entry
}

func (s FINANCIAL_ACCOUNTING) total_invoice_before_invoice_discount(array_of_entry []ACCOUNT_VALUE_QUANTITY_BARCODE) float64 {
	var total_invoice_before_invoice_discount float64
	for _, entry := range array_of_entry {
		if s.is_father(s.INCOME_STATEMENT, entry.ACCOUNT) && s.is_credit(entry.ACCOUNT) {
			total_invoice_before_invoice_discount += entry.VALUE
		} else if s.is_father(s.DISCOUNTS, entry.ACCOUNT) && !s.is_credit(entry.ACCOUNT) {
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

func (s FINANCIAL_ACCOUNTING) can_the_account_be_negative(array_of_entry []ACCOUNT_VALUE_QUANTITY_BARCODE) {
	for _, entry := range array_of_entry {
		if !(s.is_father(s.EQUITY, entry.ACCOUNT) && s.is_credit(entry.ACCOUNT)) {
			var account_balance float64
			db.QueryRow("select sum(value) from journal where account=? and date<?", entry.ACCOUNT, NOW.String()).Scan(&account_balance)
			if account_balance+entry.VALUE < 0 {
				log.Panic("you cant enter ", entry, " because you have ", account_balance, " and that will make the balance of ", entry.ACCOUNT, " negative ", account_balance+entry.VALUE, " and that you just can do it in equity_normal accounts not other accounts")
			}
		}
	}
}

func (s FINANCIAL_ACCOUNTING) convert_to_simple_entry(debit_entries, credit_entries []ACCOUNT_VALUE_QUANTITY_BARCODE) [][]ACCOUNT_VALUE_QUANTITY_BARCODE {
	simple_entries := [][]ACCOUNT_VALUE_QUANTITY_BARCODE{}
	for _, debit_entry := range debit_entries {
		for _, credit_entry := range credit_entries {
			simple_entries = append(simple_entries, []ACCOUNT_VALUE_QUANTITY_BARCODE{debit_entry, credit_entry})
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

func insert_to_journal_tag(array_of_entry []ACCOUNT_VALUE_QUANTITY_BARCODE, date time.Time, entry_expair time.Time, description string, name string, employee_name string) []journal_tag {
	var array_to_insert []journal_tag
	for _, entry := range array_of_entry {
		price := entry.VALUE / entry.QUANTITY
		if price < 0 {
			log.Panic("the ", entry.VALUE, " and ", entry.QUANTITY, " for ", entry, " should be positive both or negative both")
		}
		array_to_insert = append(array_to_insert, journal_tag{
			DATE:          date.String(),
			ENTRY_NUMBER:  0,
			ACCOUNT:       entry.ACCOUNT,
			VALUE:         entry.VALUE,
			PRICE:         price,
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

func adjuste_the_array(entry_expair time.Time, date time.Time, array_day_start_end []DAY_START_END, array_to_insert []journal_tag, adjusting_method string, description string, name string, employee_name string) [][]journal_tag {
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
	var adjusted_array_to_insert [][]journal_tag
	for _, entry := range array_to_insert {
		var value, value_counter, second_counter float64
		var one_account_adjusted_list []journal_tag
		total_value := math.Abs(entry.VALUE)
		deprecation := ROOT(total_value, total_minutes)
		value_per_second := entry.VALUE / total_minutes
		for index, element := range day_start_end_date_minutes_array {
			switch adjusting_method {
			case "linear":
				value = element.minutes * value_per_second
			case "exponential":
				value = math.Pow(deprecation, second_counter+element.minutes) - math.Pow(deprecation, second_counter)
			case "logarithmic":
				value = (total_value / math.Pow(deprecation, second_counter)) - (total_value / math.Pow(deprecation, second_counter+element.minutes))
			}
			second_counter += element.minutes

			quantity := value / entry.PRICE
			if index >= delta_days-1 {
				value = math.Abs(total_value - value_counter)
				quantity = value / entry.PRICE
			}
			value_counter += math.Abs(value)
			if entry.VALUE < 0 {
				value = -math.Abs(value)
			}
			if entry.QUANTITY < 0 {
				quantity = -math.Abs(quantity)
			}

			one_account_adjusted_list = append(one_account_adjusted_list, journal_tag{
				DATE:          element.start_date.String(),
				ENTRY_NUMBER:  0,
				ACCOUNT:       entry.ACCOUNT,
				VALUE:         value,
				PRICE:         entry.PRICE,
				QUANTITY:      quantity,
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

func (s FINANCIAL_ACCOUNTING) cost_flow(account string, quantity float64, barcode string, insert bool) float64 {
	var order_by_date_asc_or_desc string
	switch {
	case quantity > 0:
		return 0
	case s.return_cost_flow_type(account) == "fifo":
		order_by_date_asc_or_desc = "asc"
	case s.return_cost_flow_type(account) == "lifo":
		order_by_date_asc_or_desc = "desc"
	case s.return_cost_flow_type(account) == "wma":
		weighted_average([]string{account})
		order_by_date_asc_or_desc = "asc"
	default:
		return 0
	}
	rows, _ := db.Query("select price,quantity from inventory where quantity>0 and account=? and barcode=? order by date "+order_by_date_asc_or_desc, account, barcode)
	var inventory []journal_tag
	for rows.Next() {
		var tag journal_tag
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
				db.Exec("update inventory set quantity=quantity-? where account=? and price=? and quantity=? and barcode=? order by date "+order_by_date_asc_or_desc+" limit 1", quantity_count, account, item.PRICE, item.QUANTITY, barcode)
			}
			quantity_count = 0
			break
		}
		if item.QUANTITY <= quantity_count {
			costs += item.PRICE * item.QUANTITY
			if insert {
				db.Exec("delete from inventory where account=? and price=? and quantity=? and barcode=? order by date "+order_by_date_asc_or_desc+" limit 1", account, item.PRICE, item.QUANTITY, barcode)
			}
			quantity_count -= item.QUANTITY
		}
	}
	if quantity_count != 0 {
		log.Panic("you order ", quantity, " but you have ", quantity-quantity_count, " ", account, " with barcode ", barcode)
	}
	return costs
}

func (s FINANCIAL_ACCOUNTING) insert_to_database(array_of_journal_tag []journal_tag, insert_into_journal, insert_into_inventory, inventory_flow bool) {
	entry_number := float64(entry_number())
	for indexa, entry := range array_of_journal_tag {
		entry.ENTRY_NUMBER = int(entry_number)
		array_of_journal_tag[indexa].ENTRY_NUMBER = int(entry_number)
		entry_number += 0.5
		if insert_into_journal {
			db.Exec("insert into journal(date,entry_number,account,value,price,quantity,barcode,entry_expair,description,name,employee_name,entry_date,reverse) values (?,?,?,?,?,?,?,?,?,?,?,?,?)",
				&entry.DATE, &entry.ENTRY_NUMBER, &entry.ACCOUNT, &entry.VALUE, &entry.PRICE, &entry.QUANTITY, &entry.BARCODE,
				&entry.ENTRY_EXPAIR, &entry.DESCRIPTION, &entry.NAME, &entry.EMPLOYEE_NAME, &entry.ENTRY_DATE, &entry.REVERSE)
		}
		if IS_IN(entry.ACCOUNT, inventory) {
			costs := s.cost_flow(entry.ACCOUNT, entry.QUANTITY, entry.BARCODE, inventory_flow)
			if insert_into_inventory && costs == 0 {
				db.Exec("insert into inventory(date,account,price,quantity,barcode,entry_expair,name,employee_name,entry_date)values (?,?,?,?,?,?,?,?,?)",
					&entry.DATE, &entry.ACCOUNT, &entry.PRICE, &entry.QUANTITY, &entry.BARCODE, &entry.ENTRY_EXPAIR, &entry.NAME, &entry.EMPLOYEE_NAME, &entry.ENTRY_DATE)

			}
		}
	}
}

func entry_number() int {
	var tag int
	err := db.QueryRow("select max(entry_number) from journal").Scan(&tag)
	if err != nil {
		tag = 0
	}
	return tag + 1
}

func weighted_average(array_of_accounts []string) {
	for _, account := range array_of_accounts {
		db.Exec("update inventory set price=(select sum(value)/sum(quantity) from journal where account=?) where account=?", account, account)
	}
}

func (s FINANCIAL_ACCOUNTING) REVERSE_ENTRY(entry_number uint, employee_name string) {
	var array_of_entry_to_reverse []journal_tag
	rows, _ := db.Query("select * from journal where entry_number=? order by date", entry_number)
	array_of_journal_tag := select_from_journal(rows)
	if len(array_of_journal_tag) == 0 {
		log.Panic("this entry not exist")
	}
	for _, entry := range array_of_journal_tag {
		if !entry.REVERSE {
			if PARSE_DATE(entry.DATE, s.DATE_LAYOUT).Before(NOW) {
				db.Exec("update journal set reverse=True where date=? and entry_number=? and account=? and value=? and price=? and quantity=? and barcode=? and entry_expair=? and description=? and name=? and employee_name=? and entry_date=? and reverse=?",
					entry.DATE, entry.ENTRY_NUMBER, entry.ACCOUNT, entry.VALUE, entry.PRICE, entry.QUANTITY, entry.BARCODE, entry.ENTRY_EXPAIR, entry.DESCRIPTION, entry.NAME, entry.EMPLOYEE_NAME, entry.ENTRY_DATE, entry.REVERSE)
				entry.DESCRIPTION = "(reverse entry for entry number " + strconv.Itoa(entry.ENTRY_NUMBER) + " entered by " + entry.EMPLOYEE_NAME + " and revised by " + employee_name + ")"
				entry.DATE = NOW.String()
				entry.VALUE *= -1
				entry.QUANTITY *= -1
				entry.ENTRY_EXPAIR = time.Time{}.String()
				entry.EMPLOYEE_NAME = employee_name
				entry.ENTRY_DATE = NOW.String()
				array_of_entry_to_reverse = append(array_of_entry_to_reverse, entry)
				weighted_average([]string{entry.ACCOUNT})
			} else {
				db.Exec("delete from journal where date=? and entry_number=? and account=? and value=? and price=? and quantity=? and barcode=? and entry_expair=? and description=? and name=? and employee_name=? and entry_date=? and reverse=?",
					entry.DATE, entry.ENTRY_NUMBER, entry.ACCOUNT, entry.VALUE, entry.PRICE, entry.QUANTITY, entry.BARCODE, entry.ENTRY_EXPAIR, entry.DESCRIPTION, entry.NAME, entry.EMPLOYEE_NAME, entry.ENTRY_DATE, entry.REVERSE)
			}
		}
	}
	s.insert_to_database(array_of_entry_to_reverse, true, true, true)
}

func (s FINANCIAL_ACCOUNTING) JOURNAL_ENTRY(array_of_entry []ACCOUNT_VALUE_QUANTITY_BARCODE, insert, auto_completion bool, date time.Time, entry_expair time.Time, adjusting_method string,
	description string, name string, employee_name string, array_day_start_end []DAY_START_END) []journal_tag {
	array_day_start_end = check_the_params(entry_expair, adjusting_method, date, array_of_entry, array_day_start_end)
	array_of_entry = group_by_account_and_barcode(array_of_entry)
	array_of_entry = remove_zero_values(array_of_entry)
	find_account_from_barcode(array_of_entry)
	s.find_cost(array_of_entry)
	array_of_entry = s.auto_completion_the_entry(array_of_entry, auto_completion)
	// array_of_entry = s.auto_completion_the_invoice_discount(auto_completion, array_of_entry)
	array_of_entry = group_by_account_and_barcode(array_of_entry)
	array_of_entry = remove_zero_values(array_of_entry)
	s.can_the_account_be_negative(array_of_entry)
	debit_entries, credit_entries := s.check_debit_equal_credit(array_of_entry, false)
	simple_entries := s.convert_to_simple_entry(debit_entries, credit_entries)
	var all_array_to_insert []journal_tag
	for _, simple_entry := range simple_entries {
		array_to_insert := insert_to_journal_tag(simple_entry, date, entry_expair, description, name, employee_name)
		if IS_IN(adjusting_method, depreciation_methods[:]) {
			adjusted_array_to_insert := adjuste_the_array(entry_expair, date, array_day_start_end, array_to_insert, adjusting_method, description, name, employee_name)
			adjusted_array_to_insert = transpose(adjusted_array_to_insert)
			array_to_insert = unpack_the_array(adjusted_array_to_insert)
		}
		all_array_to_insert = append(all_array_to_insert, array_to_insert...)
	}
	s.insert_to_database(all_array_to_insert, insert, insert, insert)
	return all_array_to_insert
}
