package anti_accountants

import (
	"database/sql"
	"log"
)

var (
	DB        *sql.DB
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

type JOURNAL_TAG struct {
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

type INVOICE_STRUCT struct {
	ACCOUNT                string
	VALUE, PRICE, QUANTITY float64
}

func (s FINANCIAL_ACCOUNTING) check_debit_equal_credit(entries []ACCOUNT_VALUE_PRICE_QUANTITY_BARCODE) {
	var zero float64
	for _, entry := range entries {
		switch s.is_credit(entry.ACCOUNT) {
		case false:
			zero += entry.VALUE
		case true:
			zero -= entry.VALUE
		}
	}
	if zero != 0 {
		error_debit_not_equal_credit(zero, entries)
	}
}

func (s FINANCIAL_ACCOUNTING) separate_debit_from_credit(entries []ACCOUNT_VALUE_PRICE_QUANTITY_BARCODE) ([]ACCOUNT_VALUE_PRICE_QUANTITY_BARCODE, []ACCOUNT_VALUE_PRICE_QUANTITY_BARCODE) {
	var debit_entries, credit_entries []ACCOUNT_VALUE_PRICE_QUANTITY_BARCODE
	for _, entry := range entries {
		switch s.is_credit(entry.ACCOUNT) {
		case false:
			if entry.VALUE > 0 {
				debit_entries = append(debit_entries, entry)
			} else if entry.VALUE < 0 {
				credit_entries = append(credit_entries, entry)
			} else {
				log.Panic("value is zero for entry: ", entry)
			}
		case true:
			if entry.VALUE < 0 {
				debit_entries = append(debit_entries, entry)
			} else if entry.VALUE > 0 {
				credit_entries = append(credit_entries, entry)
			} else {
				log.Panic("value is zero for entry: ", entry)
			}
		}
	}
	return debit_entries, credit_entries
}

func check_one_debit_or_one_credit(debit_entries, credit_entries []ACCOUNT_VALUE_PRICE_QUANTITY_BARCODE) {
	if (len(debit_entries) != 1) && (len(credit_entries) != 1) {
		error_one_credit___one_debit("or", debit_entries, credit_entries)
	}
}

func check_one_debit_and_one_credit(debit_entries, credit_entries []ACCOUNT_VALUE_PRICE_QUANTITY_BARCODE) {
	if !((len(debit_entries) == 1) && (len(credit_entries) == 1)) {
		error_one_credit___one_debit("and", debit_entries, credit_entries)
	}
}

func (s FINANCIAL_ACCOUNTING) INVOICE(array_of_journal_tag []JOURNAL_TAG) []INVOICE_STRUCT {
	m := map[string]*INVOICE_STRUCT{}
	for _, entry := range array_of_journal_tag {
		var key string
		switch {
		case s.is_it_sub_account_using_name(s.ASSETS, entry.ACCOUNT) && !s.is_credit(entry.ACCOUNT) && !IS_IN(entry.ACCOUNT, inventory) && entry.VALUE > 0:
			key = "total"
		case s.is_it_sub_account_using_name(s.DISCOUNTS, entry.ACCOUNT) && !s.is_credit(entry.ACCOUNT):
			key = "total discounts"
		case s.is_it_sub_account_using_name(s.SALES, entry.ACCOUNT) && s.is_credit(entry.ACCOUNT):
			key = entry.ACCOUNT
		default:
			continue
		}
		sums := m[key]
		if sums == nil {
			sums = &INVOICE_STRUCT{}
			m[key] = sums
		}
		sums.VALUE += entry.VALUE
		sums.QUANTITY += entry.QUANTITY
	}
	invoice := []INVOICE_STRUCT{}
	for k, v := range m {
		invoice = append(invoice, INVOICE_STRUCT{k, v.VALUE, v.VALUE / v.QUANTITY, v.QUANTITY})
	}
	return invoice
}

func (s FINANCIAL_ACCOUNTING) INITIALIZE() {
	s.open_and_create_database()
	s.check_if_the_tree_connected()
	s.check_cost_flow_type()
	s.check_if_duplicated()
	s.check_if_the_tree_ordered()
	inventory = s.inventory_accounts()
	check_accounts("account", "inventory", " is not have fifo lifo wma on cost_flow_type field", inventory)

	// entry_number := entry_number()
	// var array_to_insert []JOURNAL_TAG
	// expair_expenses := JOURNAL_TAG{NOW.String(), entry_number, s.expair_expenses, 0, 0, 0, "", time.Time{}.String(), "to record the expiry of the goods automatically", "", "", NOW.String(), false}
	// expair_goods, _ := DB.Query("select account,price*quantity*-1,price,quantity*-1,barcode from inventory where entry_expair<? and entry_expair!='0001-01-01 00:00:00 +0000 UTC'", NOW.String())
	// for expair_goods.Next() {
	// 	tag := expair_expenses
	// 	expair_goods.Scan(&tag.ACCOUNT, &tag.value, &tag.price, &tag.quantity, &tag.barcode)
	// 	expair_expenses.value -= tag.value
	// 	expair_expenses.quantity -= tag.quantity
	// 	array_to_insert = append(array_to_insert, tag)
	// }
	// expair_expenses.price = expair_expenses.value / expair_expenses.quantity
	// array_to_insert = append(array_to_insert, expair_expenses)
	// s.insert_to_database(array_to_insert, true, false)
	// DB.Exec("delete from inventory where entry_expair<? and entry_expair!='0001-01-01 00:00:00 +0000 UTC'", NOW.String())
	DB.Exec("delete from inventory where quantity=0")

	s.check_debit_equal_credit_and_check_one_debit_and_one_credit_in_the_journal(JOURNAL_ORDERED_BY_DATE_ENTRY_NUMBER())
}

func (s FINANCIAL_ACCOUNTING) check_debit_equal_credit_and_check_one_debit_and_one_credit_in_the_journal(JOURNAL_ORDERED_BY_DATE_ENTRY_NUMBER []JOURNAL_TAG) {
	var double_entry []ACCOUNT_VALUE_PRICE_QUANTITY_BARCODE
	previous_entry_number := 1
	for _, entry := range JOURNAL_ORDERED_BY_DATE_ENTRY_NUMBER {
		if previous_entry_number != entry.ENTRY_NUMBER {
			delete_not_double_entry(double_entry, previous_entry_number)
			if len(double_entry) == 2 {
				s.check_debit_equal_credit(double_entry)
				debit_entries, credit_entries := s.separate_debit_from_credit(double_entry)
				check_one_debit_and_one_credit(debit_entries, credit_entries)
			}
			double_entry = []ACCOUNT_VALUE_PRICE_QUANTITY_BARCODE{}
		}
		double_entry = append(double_entry, ACCOUNT_VALUE_PRICE_QUANTITY_BARCODE{
			ACCOUNT:  entry.ACCOUNT,
			VALUE:    entry.VALUE,
			PRICE:    entry.PRICE,
			QUANTITY: entry.QUANTITY,
			BARCODE:  entry.BARCODE,
		})
		previous_entry_number = entry.ENTRY_NUMBER
	}
	delete_not_double_entry(double_entry, previous_entry_number)
	s.check_debit_equal_credit(double_entry)
	debit_entries, credit_entries := s.separate_debit_from_credit(double_entry)
	check_one_debit_and_one_credit(debit_entries, credit_entries)
}
