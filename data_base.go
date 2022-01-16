package anti_accountants

import (
	"database/sql"
	"fmt"
	"log"
)

func (s FINANCIAL_ACCOUNTING) open_and_create_database() {
	DB, _ = sql.Open(s.DRIVER_NAME, s.DATA_SOURCE_NAME)
	err := DB.Ping()
	error_fatal(err)
	DB.Exec("create database if not exists " + s.DATABASE_NAME)
	_, err = DB.Exec("USE " + s.DATABASE_NAME)
	error_fatal(err)
	DB.Exec("create table if not exists journal (date text,entry_number integer,account text,value real,price real,quantity real,barcode text,entry_expair text,description text,name text,employee_name text,entry_date text,reverse bool)")
	DB.Exec("create table if not exists inventory (date text,account text,price real,quantity real,barcode text,entry_expair text,name text,employee_name text,entry_date text)")
}

func delete_not_double_entry(double_entry []JOURNAL_TAG, previous_entry_number int) {
	if len(double_entry) != 2 {
		DB.Exec("delete from journal where entry_number=?", previous_entry_number)
		fmt.Println("this entry is deleted ", double_entry)
	}
}

func check_accounts(column, table, panic string, elements []string) {
	results, err := DB.Query("select " + column + " from " + table)
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
	err := DB.QueryRow("select account from journal where account=? limit 1", new_name).Scan(&tag)
	if err == nil {
		log.Panic("you can't change the name of [", name, "] to [", new_name, "] as new name because it used")
	} else {
		DB.Exec("update journal set account=? where account=?", new_name, name)
		DB.Exec("update inventory set account=? where account=?", new_name, name)
	}
}

func JOURNAL_ORDERED_BY_DATE_ENTRY_NUMBER() []JOURNAL_TAG {
	rows, _ := DB.Query("select * from journal order by date,entry_number")
	return select_from_journal(rows)
}

func account_balance(account string) float64 {
	var account_balance float64
	DB.QueryRow("select sum(value) from journal where account=?", account).Scan(&account_balance)
	return account_balance
}

func (s FINANCIAL_ACCOUNTING) insert_into_inventory(array_of_journal_tag []JOURNAL_TAG) {
	for _, entry := range array_of_journal_tag {
		costs := s.cost_flow(entry.ACCOUNT, entry.QUANTITY, entry.BARCODE, true)
		if s.asc_of_desc(entry.ACCOUNT, entry.BARCODE) != "" && costs == 0 {
			DB.Exec("insert into inventory(date,account,price,quantity,barcode,entry_expair,name,employee_name,entry_date)values (?,?,?,?,?,?,?,?,?)",
				&entry.DATE, &entry.ACCOUNT, &entry.PRICE, &entry.QUANTITY, &entry.BARCODE, &entry.ENTRY_EXPAIR, &entry.NAME, &entry.EMPLOYEE_NAME, &entry.ENTRY_DATE)
		}
	}
}

func insert_into_journal_func(array_of_journal_tag []JOURNAL_TAG) {
	for _, entry := range array_of_journal_tag {
		DB.Exec("insert into journal(date,entry_number,account,value,price,quantity,barcode,entry_expair,description,name,employee_name,entry_date,reverse) values (?,?,?,?,?,?,?,?,?,?,?,?,?)",
			&entry.DATE, &entry.ENTRY_NUMBER, &entry.ACCOUNT, &entry.VALUE, &entry.PRICE, &entry.QUANTITY, &entry.BARCODE,
			&entry.ENTRY_EXPAIR, &entry.DESCRIPTION, &entry.NAME, &entry.EMPLOYEE_NAME, &entry.ENTRY_DATE, &entry.REVERSE)
	}
}

func entry_number() int {
	var tag int
	err := DB.QueryRow("select max(entry_number) from journal").Scan(&tag)
	if err != nil {
		tag = 0
	}
	return tag + 1
}

func weighted_average(account string) {
	DB.Exec("update inventory set price=(select sum(value)/sum(quantity) from journal where account=?) where account=?", account, account)
}

func weighted_average_for_barcode(account string, barcode string) {
	DB.Exec("update inventory set price=(select sum(value)/sum(quantity) from journal where account=? and barcode=?) where account=? and barcode=?", account, barcode, account, barcode)
}
