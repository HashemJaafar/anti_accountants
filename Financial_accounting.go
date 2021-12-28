package anti_accountants

import (
	"database/sql"
	"log"
	"time"
)

var (
	db        *sql.DB
	inventory []string
)

type Financial_accounting struct {
	Date_layout               []string
	DriverName                string
	DataSourceName            string
	Database_name             string
	Assets                    string
	Current_assets            string
	Cash_and_cash_equivalents string
	Short_term_investments    string
	Receivables               string
	Inventory                 string
	Liabilities               string
	Current_liabilities       string
	Equity                    string
	Retained_earnings         string
	Dividends                 string
	Income_statement          string
	Ebitda                    string
	Sales                     string
	Cost_of_goods_sold        string
	Discounts                 string
	Invoice_discount          string
	Interest_expense          string
	Accounts                  []Account
	Invoice_discounts_list    [][2]float64
	Auto_complete_entries     [][]Account_method_value_price
}

type Account struct {
	Is_credit                    bool
	Cost_flow_type, Father, Name string
}

type journal_tag struct {
	Date          string
	Entry_number  int
	Account       string
	Value         float64
	Price         float64
	Quantity      float64
	Barcode       string
	Entry_expair  string
	Description   string
	Name          string
	Employee_name string
	Entry_date    string
	Reverse       bool
}

type invoice_struct struct {
	account                string
	value, price, quantity float64
}

func (s Financial_accounting) is_credit(name string) bool {
	for _, a := range s.Accounts {
		if a.Name == name {
			return a.Is_credit
		}
	}
	log.Panic(name, " is not debit nor credit")
	return false
}

func (s Financial_accounting) return_cost_flow_type(name string) string {
	for _, a := range s.Accounts {
		if a.Name == name {
			return a.Cost_flow_type
		}
	}
	return ""
}

func (s Financial_accounting) is_father(father, name string) bool {
	var last_name string
	for {
		for _, a := range s.Accounts {
			if a.Name == name {
				name = a.Father
			}
			if father == name {
				return true
			}
		}
		if last_name == name {
			break
		}
		last_name = name
	}
	return false
}

func (s Financial_accounting) check_debit_equal_credit(array_of_entry []Account_value_quantity_barcode, check_one_debit_and_one_credit bool) ([]Account_value_quantity_barcode, []Account_value_quantity_barcode) {
	var debit_entries, credit_entries []Account_value_quantity_barcode
	var zero float64
	for _, entry := range array_of_entry {
		switch s.is_credit(entry.Account) {
		case false:
			zero += entry.Value
			if entry.Value >= 0 {
				debit_entries = append(debit_entries, entry)
			} else {
				credit_entries = append(credit_entries, entry)
			}
		case true:
			zero -= entry.Value
			if entry.Value <= 0 {
				debit_entries = append(debit_entries, entry)
			} else {
				credit_entries = append(credit_entries, entry)
			}
		}
	}
	len_debit_entries := len(debit_entries)
	len_credit_entries := len(credit_entries)
	if (len_debit_entries != 1) && (len_credit_entries != 1) {
		log.Panic("should be one credit or one debit in the entry ", array_of_entry)
	}
	if !((len_debit_entries == 1) && (len_credit_entries == 1)) && check_one_debit_and_one_credit {
		log.Panic("should be one credit and one debit in the entry ", array_of_entry)
	}
	if zero != 0 {
		log.Panic(zero, " not equal 0 if the number>0 it means debit overstated else credit overstated debit-credit should equal zero ", array_of_entry)
	}
	return debit_entries, credit_entries
}

func check_accounts(column, table, panic string, elements []string) {
	results, err := db.Query("select " + column + " from " + table)
	error_fatal(err)
	for results.Next() {
		var tag string
		results.Scan(&tag)
		if !IS_IN(tag, elements) {
			log.Panic(tag + panic)
		}
	}
}

func Change_account_name(name, new_name string) {
	var tag string
	err := db.QueryRow("select account from journal where account=? limit 1", new_name).Scan(&tag)
	if err == nil {
		log.Panic("you can't change the name of [", name, "] to [", new_name, "] as new name because it used")
	} else {
		db.Exec("update journal set account=? where account=?", new_name, name)
		db.Exec("update inventory set account=? where account=?", new_name, name)
	}
}

func (s Financial_accounting) Invoice(array_of_journal_tag []journal_tag) []invoice_struct {
	m := map[string]*invoice_struct{}
	for _, entry := range array_of_journal_tag {
		var key string
		switch {
		case s.is_father(s.Assets, entry.Account) && !s.is_credit(entry.Account) && !IS_IN(entry.Account, inventory) && entry.Value > 0:
			key = "total"
		case s.is_father(s.Discounts, entry.Account) && !s.is_credit(entry.Account):
			key = "total discounts"
		case s.is_father(s.Sales, entry.Account) && s.is_credit(entry.Account):
			key = entry.Account
		default:
			continue
		}
		sums := m[key]
		if sums == nil {
			sums = &invoice_struct{}
			m[key] = sums
		}
		sums.value += entry.Value
		sums.quantity += entry.Quantity
	}
	invoice := []invoice_struct{}
	for k, v := range m {
		invoice = append(invoice, invoice_struct{k, v.value, v.value / v.quantity, v.quantity})
	}
	return invoice
}

func Select_journal(entry_number uint, account string, start_date, end_date time.Time) []journal_tag {
	var journal []journal_tag
	var rows *sql.Rows
	switch {
	case entry_number != 0 && account == "":
		rows, _ = db.Query("select * from journal where date>? and date<? and entry_number=? order by date", start_date.String(), end_date.String(), entry_number)
	case entry_number == 0 && account != "":
		rows, _ = db.Query("select * from journal where date>? and date<? and account=? order by date", start_date.String(), end_date.String(), account)
	default:
		log.Panic("should be one of these entry_number != 0 && account == '' or entry_number == 0 && account != '' ")
	}
	for rows.Next() {
		var tag journal_tag
		rows.Scan(&tag.Date, &tag.Entry_number, &tag.Account, &tag.Value, &tag.Price, &tag.Quantity, &tag.Barcode, &tag.Entry_expair, &tag.Description, &tag.Name, &tag.Employee_name, &tag.Entry_date, &tag.Reverse)
		journal = append(journal, tag)
	}
	return journal
}

func (s Financial_accounting) Initialize() {
	db, _ = sql.Open(s.DriverName, s.DataSourceName)
	err := db.Ping()
	error_fatal(err)
	db.Exec("create database if not exists " + s.Database_name)
	_, err = db.Exec("USE " + s.Database_name)
	error_fatal(err)
	db.Exec("create table if not exists journal (date text,entry_number integer,account text,value real,price real,quantity real,barcode text,entry_expair text,description text,name text,employee_name text,entry_date text,reverse bool)")
	db.Exec("create table if not exists inventory (date text,account text,price real,quantity real,barcode text,entry_expair text,name text,employee_name text,entry_date text)")

	var all_accounts []string
	for _, i := range s.Accounts {
		if !s.is_father("", i.Name) {
			log.Panic(i.Name, " account does not ends in ''")
		}
		all_accounts = append(all_accounts, i.Name)
		switch {
		case IS_IN(i.Cost_flow_type, []string{"fifo", "lifo", "wma"}) && !s.is_father(s.Retained_earnings, i.Name) && !i.Is_credit:
			inventory = append(inventory, i.Name)
		case i.Cost_flow_type == "":
		default:
			log.Panic(i.Cost_flow_type, " for ", i.Name, " is not in [fifo,lifo,wma,''] or you can't use it with ", s.Retained_earnings, " or is_credit==true")
		}
	}

	switch {
	case !s.is_father(s.Assets, s.Current_assets):
		log.Panic(s.Assets, " should be one of the fathers of ", s.Current_assets)
	case !s.is_father(s.Current_assets, s.Cash_and_cash_equivalents):
		log.Panic(s.Current_assets, " should be one of the fathers of ", s.Cash_and_cash_equivalents)
	case !s.is_father(s.Current_assets, s.Short_term_investments):
		log.Panic(s.Current_assets, " should be one of the fathers of ", s.Short_term_investments)
	case !s.is_father(s.Current_assets, s.Receivables):
		log.Panic(s.Current_assets, " should be one of the fathers of ", s.Receivables)
	case !s.is_father(s.Current_assets, s.Inventory):
		log.Panic(s.Current_assets, " should be one of the fathers of ", s.Inventory)
	case !s.is_father(s.Liabilities, s.Current_liabilities):
		log.Panic(s.Liabilities, " should be one of the fathers of ", s.Current_liabilities)
	case !s.is_father(s.Equity, s.Retained_earnings):
		log.Panic(s.Equity, " should be one of the fathers of ", s.Retained_earnings)
	case !s.is_father(s.Retained_earnings, s.Dividends):
		log.Panic(s.Retained_earnings, " should be one of the fathers of ", s.Dividends)
	case !s.is_father(s.Retained_earnings, s.Income_statement):
		log.Panic(s.Retained_earnings, " should be one of the fathers of ", s.Income_statement)
	case !s.is_father(s.Income_statement, s.Ebitda):
		log.Panic(s.Income_statement, " should be one of the fathers of ", s.Ebitda)
	case !s.is_father(s.Income_statement, s.Interest_expense):
		log.Panic(s.Income_statement, " should be one of the fathers of ", s.Interest_expense)
	case !s.is_father(s.Ebitda, s.Sales):
		log.Panic(s.Ebitda, " should be one of the fathers of ", s.Sales)
	case !s.is_father(s.Ebitda, s.Cost_of_goods_sold):
		log.Panic(s.Ebitda, " should be one of the fathers of ", s.Cost_of_goods_sold)
	case !s.is_father(s.Ebitda, s.Discounts):
		log.Panic(s.Ebitda, " should be one of the fathers of ", s.Discounts)
	case !s.is_father(s.Discounts, s.Invoice_discount):
		log.Panic(s.Discounts, " should be one of the fathers of ", s.Invoice_discount)
	}
	Check_if_duplicates(all_accounts)
	check_accounts("account", "inventory", " is not have fifo lifo wma on cost_flow_type field", inventory)

	// entry_number := entry_number()
	// var array_to_insert []journal_tag
	// expair_expenses := journal_tag{Now.String(), entry_number, s.expair_expenses, 0, 0, 0, "", time.Time{}.String(), "to record the expiry of the goods automatically", "", "", Now.String(), false}
	// expair_goods, _ := db.Query("select account,price*quantity*-1,price,quantity*-1,barcode from inventory where entry_expair<? and entry_expair!='0001-01-01 00:00:00 +0000 UTC'", Now.String())
	// for expair_goods.Next() {
	// 	tag := expair_expenses
	// 	expair_goods.Scan(&tag.Account, &tag.value, &tag.price, &tag.quantity, &tag.barcode)
	// 	expair_expenses.value -= tag.value
	// 	expair_expenses.quantity -= tag.quantity
	// 	array_to_insert = append(array_to_insert, tag)
	// }
	// expair_expenses.price = expair_expenses.value / expair_expenses.quantity
	// array_to_insert = append(array_to_insert, expair_expenses)
	// s.insert_to_database(array_to_insert, true, false, false)
	// db.Exec("delete from inventory where entry_expair<? and entry_expair!='0001-01-01 00:00:00 +0000 UTC'", Now.String())
	db.Exec("delete from inventory where quantity=0")

	var double_entry []Account_value_quantity_barcode
	previous_entry_number := 1
	rows, _ := db.Query("select entry_number,account,value from journal order by date,entry_number")
	for rows.Next() {
		var entry_number int
		var tag Account_value_quantity_barcode
		rows.Scan(&entry_number, &tag.Account, &tag.Value)
		if previous_entry_number != entry_number {
			s.check_debit_equal_credit(double_entry, true)
			double_entry = []Account_value_quantity_barcode{}
		}
		double_entry = append(double_entry, tag)
		previous_entry_number = entry_number
	}
}
