package anti_accountants

import (
	"fmt"
	"testing"
)

func Test_open_db(t *testing.T) {
	db_open_accounts()
	db_close_accounts()
	db_open_journal()
	db_close_journal()
	db_open_inventory()
	db_close_inventory()
}

func Test_insert_into_journal(t *testing.T) {
	var j []JOURNAL_TAG
	for i := 0; i < 100000; i++ {
		j = append(j, JOURNAL_TAG{
			REVERSE:         false,
			ENTRY_NUMBER:    i,
			LINE_NUMBER:     i,
			VALUE:           float64(i),
			PRICE_DEBIT:     float64(i),
			PRICE_CREDIT:    float64(i),
			QUANTITY_DEBIT:  float64(i),
			QUANTITY_CREDIT: float64(i),
			ACCOUNT_DEBIT:   "ACCOUNT_DEBIT",
			ACCOUNT_CREDIT:  "ACCOUNT_CREDIT",
			DESCRIPTION:     "DESCRIPTION",
			NAME:            "NAME",
			EMPLOYEE_NAME:   "EMPLOYEE_NAME",
			DATE:            "DATE",
			ENTRY_EXPAIR:    "ENTRY_EXPAIR",
			ENTRY_DATE:      "ENTRY_DATE",
		})
	}
	insert_into_journal(j)
}

func Test_JOURNAL_ORDERED_BY_DATE_ENTRY_NUMBER(t *testing.T) {
	fmt.Println(len(JOURNAL_ORDERED_BY_DATE_ENTRY_NUMBER()))
}
