package anti_accountants

import (
	"fmt"
	"testing"
	"time"
)

func Test_open_db(t *testing.T) {
	db := db_open(db_accounts)
	db.DropAll()
	db.Close()
	db = db_open(db_inventory)
	db.DropAll()
	db.Close()
	db = db_open(db_journal)
	db.DropAll()
	db.Close()
}

func Test_db_insert_into_accounts(t *testing.T) {
	db_insert_into_accounts()
}

func Test_db_insert_into_journal(t *testing.T) {
	var j []JOURNAL_TAG
	_, _, es := entry_number()
	for i := es; i < es+10000; i++ {
		j = append(j, JOURNAL_TAG{
			REVERSE:               false,
			ENTRY_NUMBER:          uint(i),
			ENTRY_NUMBER_COMPOUND: uint(i),
			ENTRY_NUMBER_SIMPLE:   uint(i),
			VALUE:                 float64(i),
			PRICE_DEBIT:           float64(i),
			PRICE_CREDIT:          float64(i),
			QUANTITY_DEBIT:        float64(i),
			QUANTITY_CREDIT:       float64(i),
			ACCOUNT_DEBIT:         "ACCOUNT_DEBIT",
			ACCOUNT_CREDIT:        "ACCOUNT_CREDIT",
			NOTES:                 "",
			NAME:                  "NAME",
			NAME_EMPLOYEE:         "",
			DATE_START:            time.Time{},
			DATE_END:              time.Time{},
			DATE_ENTRY:            time.Time{},
		})
	}
	db_insert_into_journal(j)
}

func Test_db_insert_into_inventory(t *testing.T) {
	var j []INVENTORY_TAG
	_, _, es := entry_number()
	for i := es; i < es+10000; i++ {
		j = append(j, INVENTORY_TAG{
			PRICE:        float64(i),
			QUANTITY:     float64(i),
			ACCOUNT_NAME: "ACCOUNT",
			DATE_START:   time.Time{},
			DATE_END:     time.Time{},
		})
	}
	db_insert_into_inventory(j)
}

func Test_db_read_accounts(t *testing.T) {
	a := db_read_accounts()
	e := ACCOUNTS
	test_function(a, e)
}

func Test_db_read_journal(t *testing.T) {
	fmt.Println(len(db_read_journal()))
}

func Test_db_read_inventory(t *testing.T) {
	fmt.Println(len(db_read_inventory()))
}

func Test_last_line_in_db(t *testing.T) {
	fmt.Println(last_line_in_db())
}

func Test_entry_number(t *testing.T) {
	fmt.Println(entry_number())
}
