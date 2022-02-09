package anti_accountants

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"
)

func set_to_reverse(entry JOURNAL_TAG, is_before_now bool) {
	var str string
	switch is_before_now {
	case true:
		str = "update journal set reverse=True"
	case false:
		str = "delete from journal"
	}
	DB.Exec(str+" where date=? and entry_number=? and account=? and value=? and price=? and quantity=? and barcode=? and entry_expair=? and description=? and name=? and employee_name=? and entry_date=? and reverse=?",
		entry.DATE, entry.ENTRY_NUMBER, entry.ACCOUNT, entry.VALUE, entry.PRICE, entry.QUANTITY, entry.BARCODE, entry.ENTRY_EXPAIR, entry.DESCRIPTION, entry.NAME, entry.EMPLOYEE_NAME, entry.ENTRY_DATE, entry.REVERSE)
}

func rows(reverse_method, date, entry_date string, entry_number int) *sql.Rows {
	var rows *sql.Rows
	switch reverse_method {
	case "entry_number":
		rows, _ = DB.Query("select * from journal where entry_number=?", entry_number)
	case "date":
		rows, _ = DB.Query("select * from journal where date=?", date)
	case "entry_date":
		rows, _ = DB.Query("select * from journal where entry_date=?", entry_date)
	case "date_and_entry_date":
		rows, _ = DB.Query("select * from journal where date=? and entry_date=?", date, entry_date)
	case "bigger_than_date_and_in_entry_date":
		rows, _ = DB.Query("select * from journal where date>=? and entry_date=?", date, entry_date)
	case "smaller_than_date_and_in_entry_date":
		rows, _ = DB.Query("select * from journal where date<=? and entry_date=?", date, entry_date)
	default:
		error_element_is_not_in_elements(reverse_method, []string{"entry_number", "date", "entry_date", "date_and_entry_date", "bigger_than_date_and_in_entry_date", "smaller_than_date_and_in_entry_date"})
	}
	return rows
}

func make_ACCOUNT_VALUE_PRICE_QUANTITY_BARCODE_from_JOURNAL_TAG(reverse_entry []JOURNAL_TAG) []ACCOUNT_VALUE_PRICE_QUANTITY_BARCODE {
	var entries_use_ACCOUNT_VALUE_PRICE_QUANTITY_BARCODE []ACCOUNT_VALUE_PRICE_QUANTITY_BARCODE
	for _, entry := range reverse_entry {
		entries_use_ACCOUNT_VALUE_PRICE_QUANTITY_BARCODE = append(entries_use_ACCOUNT_VALUE_PRICE_QUANTITY_BARCODE, ACCOUNT_VALUE_PRICE_QUANTITY_BARCODE{entry.ACCOUNT, entry.VALUE, entry.PRICE, entry.QUANTITY, entry.BARCODE})
	}
	return entries_use_ACCOUNT_VALUE_PRICE_QUANTITY_BARCODE
}

func (s FINANCIAL_ACCOUNTING) make_the_reverse_entries(entries []JOURNAL_TAG, reverse_using_current_date bool, employee_name string) []JOURNAL_TAG {
	var entries_to_reverse []JOURNAL_TAG
	for _, entry := range entries {
		if !entry.REVERSE {
			if PARSE_DATE(entry.DATE, s.DATE_LAYOUT).Before(NOW) {
				if reverse_using_current_date {
					entry.DATE = NOW.String()
				}
				entry.VALUE *= -1
				entry.QUANTITY *= -1
				entry.ENTRY_EXPAIR = time.Time{}.String()
				entry.DESCRIPTION = "(reverse entry for entry number " + strconv.Itoa(entry.ENTRY_NUMBER) + " entered by " + entry.EMPLOYEE_NAME + " and revised by " + employee_name + ")"
				entry.EMPLOYEE_NAME = employee_name
				entry.ENTRY_DATE = NOW.String()
				entries_to_reverse = append(entries_to_reverse, entry)
			}
		}
	}
	return entries_to_reverse
}

func (s FINANCIAL_ACCOUNTING) REVERSE_ENTRIES(reverse_using_current_date bool, employee_name, reverse_method, date, entry_date string, entry_number int) {
	rows := rows(reverse_method, date, entry_date, entry_number)
	entries := select_from_journal(rows)
	REVERSE_SLICE(entries)

	if len(entries) == 0 {
		error_this_entry_not_exist()
	}

	for _, entry := range entries {
		if entry.REVERSE {
			fmt.Println("entry number ", entry.ENTRY_NUMBER, " was reversed")
		}
	}

	reverse_entry := s.make_the_reverse_entries(entries, reverse_using_current_date, employee_name)
	s.can_the_account_be_negative(reverse_entry)

	for _, entry := range entries {
		if !entry.REVERSE {
			weighted_average(entry.ACCOUNT)
			set_to_reverse(entry, PARSE_DATE(entry.DATE, s.DATE_LAYOUT).Before(NOW))
		}
	}

	s.insert_to_database(reverse_entry, true, true)
}
