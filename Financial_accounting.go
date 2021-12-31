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

type FINANCIAL_ACCOUNTING struct {
	DATE_LAYOUT               []string
	DRIVER_NAME               string
	DATA_SOURCE_NAME          string
	DATABASE_NAME             string
	ASSETS                    string
	CURRENT_ASSETS            string
	CASH_AND_CASH_EQUIVALENTS string
	SHORT_TERM_INVESTMENTS    string
	RECEIVABLES               string
	INVENTORY                 string
	LIABILITIES               string
	CURRENT_LIABILITIES       string
	EQUITY                    string
	RETAINED_EARNINGS         string
	DIVIDENDS                 string
	INCOME_STATEMENT          string
	EBITDA                    string
	SALES                     string
	COST_OF_GOODS_SOLD        string
	DISCOUNTS                 string
	INVOICE_DISCOUNT          string
	INTEREST_EXPENSE          string
	ACCOUNTS                  []ACCOUNT
	INVOICE_DISCOUNTS_LIST    [][2]float64
	AUTO_COMPLETE_ENTRIES     []AUTO_COMPLETE_ENTRIE
}

type ACCOUNT struct {
	IS_CREDIT                    bool
	COST_FLOW_TYPE, FATHER, NAME string
}

type journal_tag struct {
	DATE          string
	ENTRY_NUMBER  int
	ACCOUNT       string
	VALUE         float64
	PRICE         float64
	QUANTITY      float64
	BARCODE       string
	ENTRY_EXPAIR  string
	DESCRIPTION   string
	NAME          string
	EMPLOYEE_NAME string
	ENTRY_DATE    string
	REVERSE       bool
}

type invoice_struct struct {
	account                string
	value, price, quantity float64
}

func (s FINANCIAL_ACCOUNTING) is_credit(name string) bool {
	for _, a := range s.ACCOUNTS {
		if a.NAME == name {
			return a.IS_CREDIT
		}
	}
	log.Panic(name, " is not debit nor credit")
	return false
}

func (s FINANCIAL_ACCOUNTING) return_cost_flow_type(name string) string {
	for _, a := range s.ACCOUNTS {
		if a.NAME == name {
			return a.COST_FLOW_TYPE
		}
	}
	return ""
}

func (s FINANCIAL_ACCOUNTING) is_father(father, name string) bool {
	var last_name string
	for {
		for _, a := range s.ACCOUNTS {
			if a.NAME == name {
				name = a.FATHER
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

func (s FINANCIAL_ACCOUNTING) check_debit_equal_credit(array_of_entry []ACCOUNT_VALUE_QUANTITY_BARCODE, check_one_debit_and_one_credit bool) ([]ACCOUNT_VALUE_QUANTITY_BARCODE, []ACCOUNT_VALUE_QUANTITY_BARCODE) {
	var debit_entries, credit_entries []ACCOUNT_VALUE_QUANTITY_BARCODE
	var zero float64
	for _, entry := range array_of_entry {
		switch s.is_credit(entry.ACCOUNT) {
		case false:
			zero += entry.VALUE
			if entry.VALUE >= 0 {
				debit_entries = append(debit_entries, entry)
			} else {
				credit_entries = append(credit_entries, entry)
			}
		case true:
			zero -= entry.VALUE
			if entry.VALUE <= 0 {
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

func CHANGE_ACCOUNT_NAME(name, new_name string) {
	var tag string
	err := db.QueryRow("select account from journal where account=? limit 1", new_name).Scan(&tag)
	if err == nil {
		log.Panic("you can't change the name of [", name, "] to [", new_name, "] as new name because it used")
	} else {
		db.Exec("update journal set account=? where account=?", new_name, name)
		db.Exec("update inventory set account=? where account=?", new_name, name)
	}
}

func (s FINANCIAL_ACCOUNTING) INVOICE(array_of_journal_tag []journal_tag) []invoice_struct {
	m := map[string]*invoice_struct{}
	for _, entry := range array_of_journal_tag {
		var key string
		switch {
		case s.is_father(s.ASSETS, entry.ACCOUNT) && !s.is_credit(entry.ACCOUNT) && !IS_IN(entry.ACCOUNT, inventory) && entry.VALUE > 0:
			key = "total"
		case s.is_father(s.DISCOUNTS, entry.ACCOUNT) && !s.is_credit(entry.ACCOUNT):
			key = "total discounts"
		case s.is_father(s.SALES, entry.ACCOUNT) && s.is_credit(entry.ACCOUNT):
			key = entry.ACCOUNT
		default:
			continue
		}
		sums := m[key]
		if sums == nil {
			sums = &invoice_struct{}
			m[key] = sums
		}
		sums.value += entry.VALUE
		sums.quantity += entry.QUANTITY
	}
	invoice := []invoice_struct{}
	for k, v := range m {
		invoice = append(invoice, invoice_struct{k, v.value, v.value / v.quantity, v.quantity})
	}
	return invoice
}

func SELECT_JOURNAL(entry_number uint, account string, start_date, end_date time.Time) []journal_tag {
	var rows *sql.Rows
	switch {
	case entry_number != 0 && account == "":
		rows, _ = db.Query("select * from journal where date>? and date<? and entry_number=? order by date", start_date.String(), end_date.String(), entry_number)
	case entry_number == 0 && account != "":
		rows, _ = db.Query("select * from journal where date>? and date<? and account=? order by date", start_date.String(), end_date.String(), account)
	default:
		log.Panic("should be one of these entry_number != 0 && account == '' or entry_number == 0 && account != '' ")
	}
	journal := select_from_journal(rows)
	return journal
}

func (s FINANCIAL_ACCOUNTING) INITIALIZE() {
	db, _ = sql.Open(s.DRIVER_NAME, s.DATA_SOURCE_NAME)
	err := db.Ping()
	error_fatal(err)
	db.Exec("create database if not exists " + s.DATABASE_NAME)
	_, err = db.Exec("USE " + s.DATABASE_NAME)
	error_fatal(err)
	db.Exec("create table if not exists journal (date text,entry_number integer,account text,value real,price real,quantity real,barcode text,entry_expair text,description text,name text,employee_name text,entry_date text,reverse bool)")
	db.Exec("create table if not exists inventory (date text,account text,price real,quantity real,barcode text,entry_expair text,name text,employee_name text,entry_date text)")

	var all_accounts []string
	for _, i := range s.ACCOUNTS {
		if !s.is_father("", i.NAME) {
			log.Panic(i.NAME, " account does not ends in ''")
		}
		all_accounts = append(all_accounts, i.NAME)
		switch {
		case IS_IN(i.COST_FLOW_TYPE, []string{"fifo", "lifo", "wma"}) && !s.is_father(s.RETAINED_EARNINGS, i.NAME) && !i.IS_CREDIT:
			inventory = append(inventory, i.NAME)
		case i.COST_FLOW_TYPE == "":
		default:
			log.Panic(i.COST_FLOW_TYPE, " for ", i.NAME, " is not in [fifo,lifo,wma,''] or you can't use it with ", s.RETAINED_EARNINGS, " or is_credit==true")
		}
	}

	switch {
	case !s.is_father(s.ASSETS, s.CURRENT_ASSETS):
		log.Panic(s.ASSETS, " should be one of the fathers of ", s.CURRENT_ASSETS)
	case !s.is_father(s.CURRENT_ASSETS, s.CASH_AND_CASH_EQUIVALENTS):
		log.Panic(s.CURRENT_ASSETS, " should be one of the fathers of ", s.CASH_AND_CASH_EQUIVALENTS)
	case !s.is_father(s.CURRENT_ASSETS, s.SHORT_TERM_INVESTMENTS):
		log.Panic(s.CURRENT_ASSETS, " should be one of the fathers of ", s.SHORT_TERM_INVESTMENTS)
	case !s.is_father(s.CURRENT_ASSETS, s.RECEIVABLES):
		log.Panic(s.CURRENT_ASSETS, " should be one of the fathers of ", s.RECEIVABLES)
	case !s.is_father(s.CURRENT_ASSETS, s.INVENTORY):
		log.Panic(s.CURRENT_ASSETS, " should be one of the fathers of ", s.INVENTORY)
	case !s.is_father(s.LIABILITIES, s.CURRENT_LIABILITIES):
		log.Panic(s.LIABILITIES, " should be one of the fathers of ", s.CURRENT_LIABILITIES)
	case !s.is_father(s.EQUITY, s.RETAINED_EARNINGS):
		log.Panic(s.EQUITY, " should be one of the fathers of ", s.RETAINED_EARNINGS)
	case !s.is_father(s.RETAINED_EARNINGS, s.DIVIDENDS):
		log.Panic(s.RETAINED_EARNINGS, " should be one of the fathers of ", s.DIVIDENDS)
	case !s.is_father(s.RETAINED_EARNINGS, s.INCOME_STATEMENT):
		log.Panic(s.RETAINED_EARNINGS, " should be one of the fathers of ", s.INCOME_STATEMENT)
	case !s.is_father(s.INCOME_STATEMENT, s.EBITDA):
		log.Panic(s.INCOME_STATEMENT, " should be one of the fathers of ", s.EBITDA)
	case !s.is_father(s.INCOME_STATEMENT, s.INTEREST_EXPENSE):
		log.Panic(s.INCOME_STATEMENT, " should be one of the fathers of ", s.INTEREST_EXPENSE)
	case !s.is_father(s.EBITDA, s.SALES):
		log.Panic(s.EBITDA, " should be one of the fathers of ", s.SALES)
	case !s.is_father(s.EBITDA, s.COST_OF_GOODS_SOLD):
		log.Panic(s.EBITDA, " should be one of the fathers of ", s.COST_OF_GOODS_SOLD)
	case !s.is_father(s.EBITDA, s.DISCOUNTS):
		log.Panic(s.EBITDA, " should be one of the fathers of ", s.DISCOUNTS)
	case !s.is_father(s.DISCOUNTS, s.INVOICE_DISCOUNT):
		log.Panic(s.DISCOUNTS, " should be one of the fathers of ", s.INVOICE_DISCOUNT)
	}
	_, duplicated_element := CHECK_IF_DUPLICATES(all_accounts)
	if len(duplicated_element) != 0 {
		log.Panic(duplicated_element, " is duplicated values in the fields of FINANCIAL_ACCOUNTING and that make error. you should remove the duplicate")
	}
	check_accounts("account", "inventory", " is not have fifo lifo wma on cost_flow_type field", inventory)

	// entry_number := entry_number()
	// var array_to_insert []journal_tag
	// expair_expenses := journal_tag{NOW.String(), entry_number, s.expair_expenses, 0, 0, 0, "", time.Time{}.String(), "to record the expiry of the goods automatically", "", "", NOW.String(), false}
	// expair_goods, _ := db.Query("select account,price*quantity*-1,price,quantity*-1,barcode from inventory where entry_expair<? and entry_expair!='0001-01-01 00:00:00 +0000 UTC'", NOW.String())
	// for expair_goods.Next() {
	// 	tag := expair_expenses
	// 	expair_goods.Scan(&tag.ACCOUNT, &tag.value, &tag.price, &tag.quantity, &tag.barcode)
	// 	expair_expenses.value -= tag.value
	// 	expair_expenses.quantity -= tag.quantity
	// 	array_to_insert = append(array_to_insert, tag)
	// }
	// expair_expenses.price = expair_expenses.value / expair_expenses.quantity
	// array_to_insert = append(array_to_insert, expair_expenses)
	// s.insert_to_database(array_to_insert, true, false, false)
	// db.Exec("delete from inventory where entry_expair<? and entry_expair!='0001-01-01 00:00:00 +0000 UTC'", NOW.String())
	db.Exec("delete from inventory where quantity=0")

	var double_entry []ACCOUNT_VALUE_QUANTITY_BARCODE
	previous_entry_number := 1
	rows, _ := db.Query("select entry_number,account,value from journal order by date,entry_number")
	for rows.Next() {
		var entry_number int
		var tag ACCOUNT_VALUE_QUANTITY_BARCODE
		rows.Scan(&entry_number, &tag.ACCOUNT, &tag.VALUE)
		if previous_entry_number != entry_number {
			s.check_debit_equal_credit(double_entry, true)
			double_entry = []ACCOUNT_VALUE_QUANTITY_BARCODE{}
		}
		double_entry = append(double_entry, tag)
		previous_entry_number = entry_number
	}
}
