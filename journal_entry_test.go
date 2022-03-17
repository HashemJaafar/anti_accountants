package anti_accountants

import (
	"testing"
	"time"
)

func Test_JOURNAL_ENTRY(t *testing.T) {
	JOURNAL_ENTRY([]VALUE_PRICE_QUANTITY_ACCOUNT_BARCODE{
		{100, 0, 1, "book", "1"},
		{-100, 0, -100, "cash", "1"},
	},
		true, true, true, time.Time{}, time.Time{}, "", "", "saba", "hashem", []DAY_START_END{})
}
