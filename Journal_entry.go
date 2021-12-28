package anti_accountants

import (
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Account_value_quantity_barcode struct {
	Account  string
	Value    float64
	Quantity float64
	Barcode  string
}

type Day_start_end struct {
	Day          string
	Start_hour   int
	Start_minute int
	End_hour     int
	End_minute   int
}

type Account_method_value_price struct {
	Account, Method         string
	Value_or_percent, Price float64
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
	Now                  = time.Now()
)

func check_the_params(entry_expair time.Time, adjusting_method string, date time.Time, array_of_entry []Account_value_quantity_barcode, array_day_start_end []Day_start_end) []Day_start_end {
	if entry_expair.IsZero() == IS_IN(adjusting_method, adjusting_methods[:]) {
		log.Panic("check entry_expair => ", entry_expair, " and adjusting_method => ", adjusting_method, " should be in ", adjusting_methods)
	}
	if !entry_expair.IsZero() {
		check_dates(date, entry_expair)
	}
	for _, entry := range array_of_entry {
		if IS_IN(entry.Account, inventory) && !IS_IN(adjusting_method, []string{"expire", ""}) {
			log.Panic(entry.Account + " is in inventory you just can use expire or make it empty")
		}
	}
	if IS_IN(adjusting_method, depreciation_methods[:]) {
		if len(array_day_start_end) == 0 {
			array_day_start_end = []Day_start_end{
				{"saturday", 0, 0, 23, 59},
				{"sunday", 0, 0, 23, 59},
				{"monday", 0, 0, 23, 59},
				{"tuesday", 0, 0, 23, 59},
				{"wednesday", 0, 0, 23, 59},
				{"thursday", 0, 0, 23, 59},
				{"friday", 0, 0, 23, 59}}
		}
		for index, element := range array_day_start_end {
			array_day_start_end[index].Day = strings.Title(element.Day)
			switch {
			case !IS_IN(array_day_start_end[index].Day, standard_days[:]):
				log.Panic("error ", element.Day, " for ", element, " is not in ", standard_days)
			case element.Start_hour < 0:
				log.Panic("error ", element.Start_hour, " for ", element, " is < 0")
			case element.Start_hour > 23:
				log.Panic("error ", element.Start_hour, " for ", element, " is > 23")
			case element.Start_minute < 0:
				log.Panic("error ", element.Start_minute, " for ", element, " is < 0")
			case element.Start_minute > 59:
				log.Panic("error ", element.Start_minute, " for ", element, " is > 59")
			case element.End_hour < 0:
				log.Panic("error ", element.End_hour, " for ", element, " is < 0")
			case element.End_hour > 23:
				log.Panic("error ", element.End_hour, " for ", element, " is > 23")
			case element.End_minute < 0:
				log.Panic("error ", element.End_minute, " for ", element, " is < 0")
			case element.End_minute > 59:
				log.Panic("error ", element.End_minute, " for ", element, " is > 59")
			}
		}
	}
	return array_day_start_end
}

func group_by_account_and_barcode(array_of_entry []Account_value_quantity_barcode) []Account_value_quantity_barcode {
	type Account_barcode struct {
		Account, barcode string
	}
	g := map[Account_barcode]*Account_value_quantity_barcode{}
	for _, v := range array_of_entry {
		key := Account_barcode{v.Account, v.Barcode}
		sums := g[key]
		if sums == nil {
			sums = &Account_value_quantity_barcode{}
			g[key] = sums
		}
		sums.Value += v.Value
		sums.Quantity += v.Quantity
	}
	array_of_entry = []Account_value_quantity_barcode{}
	for key, v := range g {
		array_of_entry = append(array_of_entry, Account_value_quantity_barcode{key.Account, v.Value, v.Quantity, key.barcode})
	}
	return array_of_entry
}

func remove_zero_values(array_of_entry []Account_value_quantity_barcode) []Account_value_quantity_barcode {
	var index int
	for index < len(array_of_entry) {
		if array_of_entry[index].Value == 0 || array_of_entry[index].Quantity == 0 {
			// fmt.Println(array_of_entry[index], " is removed because one of the values is 0")
			array_of_entry = append(array_of_entry[:index], array_of_entry[index+1:]...)
		} else {
			index++
		}
	}
	return array_of_entry
}

func find_barcode(array_of_entry []Account_value_quantity_barcode) {
	for index, entry := range array_of_entry {
		if entry.Account == "" && entry.Barcode == "" {
			log.Panic("can't find the account name if the barcode is empty in ", entry)
		}
		var tag string
		if entry.Account == "" {
			err := db.QueryRow("select account from journal where barcode=? limit 1", entry.Barcode).Scan(&tag)
			if err != nil {
				log.Panic("the barcode is wrong for ", entry)
			}
			array_of_entry[index].Account = tag
		}
	}
}

func (s Financial_accounting) auto_completion_the_entry(array_of_entry []Account_value_quantity_barcode, auto_completion bool) []Account_value_quantity_barcode {
	for index, entry := range array_of_entry {
		costs := s.cost_flow(entry.Account, entry.Quantity, entry.Barcode, false)
		if costs != 0 {
			array_of_entry[index] = Account_value_quantity_barcode{entry.Account, -costs, entry.Quantity, entry.Barcode}
		}
		if auto_completion {
			for _, complement := range s.Auto_complete_entries {
				if complement[0].Account == entry.Account && (entry.Quantity >= 0) == (complement[0].Value_or_percent >= 0) {
					if costs == 0 {
						array_of_entry[index] = Account_value_quantity_barcode{complement[0].Account, complement[0].Price * entry.Quantity, entry.Quantity, ""}
					}
					for _, i := range complement[1:] {
						switch i.Method {
						case "copy_abs":
							array_of_entry = append(array_of_entry, Account_value_quantity_barcode{i.Account, math.Abs(array_of_entry[index].Value), math.Abs(array_of_entry[index].Quantity), ""})
						case "copy":
							array_of_entry = append(array_of_entry, Account_value_quantity_barcode{i.Account, array_of_entry[index].Value, array_of_entry[index].Quantity, ""})
						case "quantity_ratio":
							array_of_entry = append(array_of_entry, Account_value_quantity_barcode{i.Account, math.Abs(array_of_entry[index].Quantity) * i.Price * i.Value_or_percent, math.Abs(array_of_entry[index].Quantity) * i.Value_or_percent, ""})
						case "value":
							array_of_entry = append(array_of_entry, Account_value_quantity_barcode{i.Account, i.Value_or_percent, i.Value_or_percent / i.Price, ""})
						default:
							log.Panic(i.Method, "in the method field for ", i, " dose not exist you just can use copy_abs or copy or quantity_ratio or value")
						}
					}
				}
			}
		}
	}
	return array_of_entry
}

func (s Financial_accounting) auto_completion_the_invoice_discount(auto_completion bool, array_of_entry []Account_value_quantity_barcode) []Account_value_quantity_barcode {
	if auto_completion {
		var total_invoice_before_invoice_discount, discount float64
		for _, entry := range array_of_entry {
			if s.is_father(s.Income_statement, entry.Account) && s.is_credit(entry.Account) {
				total_invoice_before_invoice_discount += entry.Value
			} else if s.is_father(s.Discounts, entry.Account) && !s.is_credit(entry.Account) {
				total_invoice_before_invoice_discount -= entry.Value
			}
		}
		for _, i := range s.Invoice_discounts_list {
			if total_invoice_before_invoice_discount >= i[0] {
				discount = i[1]
			}
		}
		invoice_discount := discount_tax_calculator(total_invoice_before_invoice_discount, discount)
		array_of_entry = append(array_of_entry, Account_value_quantity_barcode{s.Invoice_discount, invoice_discount, 1, ""})
	}
	return array_of_entry
}

func discount_tax_calculator(price, discount_tax float64) float64 {
	if discount_tax < 0 {
		discount_tax = math.Abs(discount_tax)
	} else if discount_tax > 0 {
		discount_tax = price * discount_tax
	}
	return discount_tax
}

func (s Financial_accounting) can_the_account_be_negative(array_of_entry []Account_value_quantity_barcode) {
	for _, entry := range array_of_entry {
		if !(s.is_father(s.Equity, entry.Account) && s.is_credit(entry.Account)) {
			var account_balance float64
			db.QueryRow("select sum(value) from journal where account=? and date<?", entry.Account, Now.String()).Scan(&account_balance)
			if account_balance+entry.Value < 0 {
				log.Panic("you cant enter ", entry, " because you have ", account_balance, " and that will make the balance of ", entry.Account, " negative ", account_balance+entry.Value, " and that you just can do it in equity_normal accounts not other accounts")
			}
		}
	}
}

func (s Financial_accounting) convert_to_simple_entry(debit_entries, credit_entries []Account_value_quantity_barcode) [][]Account_value_quantity_barcode {
	simple_entries := [][]Account_value_quantity_barcode{}
	for _, debit_entry := range debit_entries {
		for _, credit_entry := range credit_entries {
			simple_entries = append(simple_entries, []Account_value_quantity_barcode{debit_entry, credit_entry})
		}
	}
	for _, a := range simple_entries {
		switch math.Abs(a[0].Value) >= math.Abs(a[1].Value) {
		case true:
			sign := a[0].Value / a[1].Value
			price := a[0].Value / a[0].Quantity
			a[0].Value = a[1].Value * sign / math.Abs(sign)
			a[0].Quantity = a[0].Value / price
		case false:
			sign := a[0].Value / a[1].Value
			price := a[1].Value / a[1].Quantity
			a[1].Value = a[0].Value * sign / math.Abs(sign)
			a[1].Quantity = a[1].Value / price
		}
	}
	return simple_entries
}

func insert_to_journal_tag(array_of_entry []Account_value_quantity_barcode, date time.Time, entry_expair time.Time, description string, name string, employee_name string) []journal_tag {
	var array_to_insert []journal_tag
	for _, entry := range array_of_entry {
		price := entry.Value / entry.Quantity
		if price < 0 {
			log.Panic("the ", entry.Value, " and ", entry.Quantity, " for ", entry, " should be positive both or negative both")
		}
		array_to_insert = append(array_to_insert, journal_tag{
			Date:          date.String(),
			Entry_number:  0,
			Account:       entry.Account,
			Value:         entry.Value,
			Price:         price,
			Quantity:      entry.Quantity,
			Barcode:       entry.Barcode,
			Entry_expair:  entry_expair.String(),
			Description:   description,
			Name:          name,
			Employee_name: employee_name,
			Entry_date:    Now.String(),
			Reverse:       false,
		})
	}
	return array_to_insert
}

func adjuste_the_array(entry_expair time.Time, date time.Time, array_day_start_end []Day_start_end, array_to_insert []journal_tag, adjusting_method string, description string, name string, employee_name string) [][]journal_tag {
	var day_start_end_date_minutes_array []day_start_end_date_minutes
	var total_minutes float64
	var previous_end_date, end time.Time
	delta_days := int(entry_expair.Sub(date).Hours()/24 + 1)
	year, month_sting, day := date.Date()
	for day_counter := 0; day_counter < delta_days; day_counter++ {
		for _, element := range array_day_start_end {
			if start := time.Date(year, month_sting, day+day_counter, element.Start_hour, element.Start_minute, 0, 0, time.Local); start.Weekday().String() == element.Day {
				previous_end_date = end
				end = time.Date(year, month_sting, day+day_counter, element.End_hour, element.End_minute, 0, 0, time.Local)
				if start.After(end) {
					log.Panic("the start_hour and start_minute should be smaller than end_hour and end_minute for ", element)
				}
				if previous_end_date.After(start) {
					log.Panic("the end_hour and end_minute for ", element.Day, " should be smaller than start_hour and start_minute for the second ", element)
				}
				minutes := end.Sub(start).Minutes()
				total_minutes += minutes
				day_start_end_date_minutes_array = append(day_start_end_date_minutes_array, day_start_end_date_minutes{element.Day, start, end, minutes})
			}
		}
	}
	var adjusted_array_to_insert [][]journal_tag
	for _, entry := range array_to_insert {
		var value, value_counter, second_counter float64
		var one_account_adjusted_list []journal_tag
		total_value := math.Abs(entry.Value)
		deprecation := Root(total_value, total_minutes)
		value_per_second := entry.Value / total_minutes
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

			quantity := value / entry.Price
			if index >= delta_days-1 {
				value = math.Abs(total_value - value_counter)
				quantity = value / entry.Price
			}
			value_counter += math.Abs(value)
			if entry.Value < 0 {
				value = -math.Abs(value)
			}
			if entry.Quantity < 0 {
				quantity = -math.Abs(quantity)
			}

			one_account_adjusted_list = append(one_account_adjusted_list, journal_tag{
				Date:          element.start_date.String(),
				Entry_number:  0,
				Account:       entry.Account,
				Value:         value,
				Price:         entry.Price,
				Quantity:      quantity,
				Barcode:       entry.Barcode,
				Entry_expair:  element.end_date.String(),
				Description:   description,
				Name:          name,
				Employee_name: employee_name,
				Entry_date:    Now.String(),
				Reverse:       false,
			})
		}
		adjusted_array_to_insert = append(adjusted_array_to_insert, one_account_adjusted_list)
	}
	return adjusted_array_to_insert
}

func (s Financial_accounting) cost_flow(account string, quantity float64, barcode string, insert bool) float64 {
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
		rows.Scan(&tag.Price, &tag.Quantity)
		inventory = append(inventory, tag)
	}
	quantity = math.Abs(quantity)
	quantity_count := quantity
	var costs float64
	for _, item := range inventory {
		if item.Quantity > quantity_count {
			costs += item.Price * quantity_count
			if insert {
				db.Exec("update inventory set quantity=quantity-? where account=? and price=? and quantity=? and barcode=? order by date "+order_by_date_asc_or_desc+" limit 1", quantity_count, account, item.Price, item.Quantity, barcode)
			}
			quantity_count = 0
			break
		}
		if item.Quantity <= quantity_count {
			costs += item.Price * item.Quantity
			if insert {
				db.Exec("delete from inventory where account=? and price=? and quantity=? and barcode=? order by date "+order_by_date_asc_or_desc+" limit 1", account, item.Price, item.Quantity, barcode)
			}
			quantity_count -= item.Quantity
		}
	}
	if quantity_count != 0 {
		log.Panic("you order ", quantity, " but you have ", quantity-quantity_count, " ", account, " with barcode ", barcode)
	}
	return costs
}

func (s Financial_accounting) insert_to_database(array_of_journal_tag []journal_tag, insert_into_journal, insert_into_inventory, inventory_flow bool) {
	entry_number := float64(entry_number())
	for indexa, entry := range array_of_journal_tag {
		entry.Entry_number = int(entry_number)
		array_of_journal_tag[indexa].Entry_number = int(entry_number)
		entry_number += 0.5
		if insert_into_journal {
			db.Exec("insert into journal(date,entry_number,account,value,price,quantity,barcode,entry_expair,description,name,employee_name,entry_date,reverse) values (?,?,?,?,?,?,?,?,?,?,?,?,?)",
				&entry.Date, &entry.Entry_number, &entry.Account, &entry.Value, &entry.Price, &entry.Quantity, &entry.Barcode,
				&entry.Entry_expair, &entry.Description, &entry.Name, &entry.Employee_name, &entry.Entry_date, &entry.Reverse)
		}
		if IS_IN(entry.Account, inventory) {
			costs := s.cost_flow(entry.Account, entry.Quantity, entry.Barcode, inventory_flow)
			if insert_into_inventory && costs == 0 {
				db.Exec("insert into inventory(date,account,price,quantity,barcode,entry_expair,name,employee_name,entry_date)values (?,?,?,?,?,?,?,?,?)",
					&entry.Date, &entry.Account, &entry.Price, &entry.Quantity, &entry.Barcode, &entry.Entry_expair, &entry.Name, &entry.Employee_name, &entry.Entry_date)

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

func (s Financial_accounting) Reverse_entry(entry_number uint, employee_name string) {
	var array_of_entry_to_reverse, array_of_journal_tag []journal_tag
	rows, _ := db.Query("select * from journal where entry_number=? order by date", entry_number)
	for rows.Next() {
		var tag journal_tag
		rows.Scan(&tag.Date, &tag.Entry_number, &tag.Account, &tag.Value, &tag.Price, &tag.Quantity, &tag.Barcode, &tag.Entry_expair, &tag.Description, &tag.Name, &tag.Employee_name, &tag.Entry_date, &tag.Reverse)
		array_of_journal_tag = append(array_of_journal_tag, tag)
	}
	if len(array_of_journal_tag) == 0 {
		log.Panic("this entry not exist")
	}
	for _, entry := range array_of_journal_tag {
		if !entry.Reverse {
			if Parse_date(entry.Date, s.Date_layout).Before(Now) {
				db.Exec("update journal set reverse=True where date=? and entry_number=? and account=? and value=? and price=? and quantity=? and barcode=? and entry_expair=? and description=? and name=? and employee_name=? and entry_date=? and reverse=?",
					entry.Date, entry.Entry_number, entry.Account, entry.Value, entry.Price, entry.Quantity, entry.Barcode, entry.Entry_expair, entry.Description, entry.Name, entry.Employee_name, entry.Entry_date, entry.Reverse)
				entry.Description = "(reverse entry for entry number " + strconv.Itoa(entry.Entry_number) + " entered by " + entry.Employee_name + " and revised by " + employee_name + ")"
				entry.Date = Now.String()
				entry.Value *= -1
				entry.Quantity *= -1
				entry.Entry_expair = time.Time{}.String()
				entry.Employee_name = employee_name
				entry.Entry_date = Now.String()
				array_of_entry_to_reverse = append(array_of_entry_to_reverse, entry)
				weighted_average([]string{entry.Account})
			} else {
				db.Exec("delete from journal where date=? and entry_number=? and account=? and value=? and price=? and quantity=? and barcode=? and entry_expair=? and description=? and name=? and employee_name=? and entry_date=? and reverse=?",
					entry.Date, entry.Entry_number, entry.Account, entry.Value, entry.Price, entry.Quantity, entry.Barcode, entry.Entry_expair, entry.Description, entry.Name, entry.Employee_name, entry.Entry_date, entry.Reverse)
			}
		}
	}
	s.insert_to_database(array_of_entry_to_reverse, true, true, true)
}

func (s Financial_accounting) Journal_entry(array_of_entry []Account_value_quantity_barcode, insert, auto_completion bool, date time.Time, entry_expair time.Time, adjusting_method string,
	description string, name string, employee_name string, array_day_start_end []Day_start_end) []journal_tag {
	array_day_start_end = check_the_params(entry_expair, adjusting_method, date, array_of_entry, array_day_start_end)
	array_of_entry = group_by_account_and_barcode(array_of_entry)
	array_of_entry = remove_zero_values(array_of_entry)
	find_barcode(array_of_entry)
	array_of_entry = s.auto_completion_the_entry(array_of_entry, auto_completion)
	array_of_entry = s.auto_completion_the_invoice_discount(auto_completion, array_of_entry)
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
