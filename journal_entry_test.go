package anti_accountants

import (
	"testing"
	"time"
)

func Test_SET_DATE_END_TO_ZERO_IF_SMALLER_THAN_DATE_START(t *testing.T) {
	a := SET_DATE_END_TO_ZERO_IF_SMALLER_THAN_DATE_START(time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC), time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC))
	TEST(true, a, time.Time{})
	x := time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)
	a = SET_DATE_END_TO_ZERO_IF_SMALLER_THAN_DATE_START(time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC), x)
	TEST(true, a, x)
}
func Test_SET_ADJUSTING_METHOD(t *testing.T) {
	a := SET_ADJUSTING_METHOD(time.Time{}, "liea", []PRICE_QUANTITY_ACCOUNT_BARCODE{
		{1, 10, "book1", "1"},
	})
	TEST(true, a, "")
	a = SET_ADJUSTING_METHOD(time.Now(), "lmd", []PRICE_QUANTITY_ACCOUNT_BARCODE{
		{1, 10, "assets", "1"},
	})
	TEST(true, a, LINEAR)
	a = SET_ADJUSTING_METHOD(time.Now(), LINEAR, []PRICE_QUANTITY_ACCOUNT_BARCODE{
		{1, 10, "book1", "1"},
	})
	PRINT_FORMATED_ACCOUNTS()
	TEST(true, a, EXPIRE)
}
func Test_GROUP_BY_ACCOUNT_AND_BARCODE(t *testing.T) {
	a := GROUP_BY_ACCOUNT([]PRICE_QUANTITY_ACCOUNT_BARCODE{{1, 10, "book1", "1"}, {1, 10, "book1", "2"}})
	TEST(true, a, []PRICE_QUANTITY_ACCOUNT_BARCODE{{1, 20, "book1", ""}})
}
func Test_FIND_COST(t *testing.T) {
	// a := FIND_COST()
	// TEST(true,a,)
}
func Test_INSERT_TO_JOURNAL_TAG(t *testing.T) {
	// a := INSERT_TO_JOURNAL_TAG()
	// TEST(true,a,)
}
func Test_INCREASE_THE_VALUE_TO_MAKE_THE_NEW_BALANCE_FOR_THE_ACCOUNT_POSITIVE(t *testing.T) {
	// a := INCREASE_THE_VALUE_TO_MAKE_THE_NEW_BALANCE_FOR_THE_ACCOUNT_POSITIVE()
	// TEST(true,a,)
}
func Test_FIND_ACCOUNT_FROM_BARCODE(t *testing.T) {
	// a := FIND_ACCOUNT_FROM_BARCODE()
	// TEST(true,a,)
}
func Test_SET_THE_SIGN_OF_THE_VALUE_SAME_SIGN_OF_QUANTITY(t *testing.T) {
	// a := SET_THE_SIGN_OF_THE_VALUE_SAME_SIGN_OF_QUANTITY()
	// TEST(true,a,)
}
func Test_COST_FLOW(t *testing.T) {
	a1, _, a2 := COST_FLOW("book1", 10, false)
	TEST(true, a1, 0)
	TEST(true, a2, ERROR_SHOULD_BE_NEGATIVE)
	a1, _, a2 = COST_FLOW("assets", -10, false)
	TEST(true, a1, 0)
	TEST(true, a2, ERROR_NOT_INVENTORY_ACCOUNT)
	a1, _, a2 = COST_FLOW("book1", -40, true)
	TEST(true, a1, 40)
	TEST(true, a2, nil)
}
func Test_FLOW_TYPE(t *testing.T) {
	a1, a2 := FLOW_TYPE("book1")
	TEST(true, a1, FIFO)
	TEST(true, a2, true)
	a1, a2 = FLOW_TYPE("book")
	TEST(true, a1, "")
	TEST(true, a2, true)
}
func Test_INSERT_TO_DATABASE(t *testing.T) {
	// a := INSERT_TO_DATABASE()
	// TEST(true,a,)
}
func Test_INSERT_ENTRY_NUMBER(t *testing.T) {
	// a := INSERT_ENTRY_NUMBER()
	// TEST(true,a,)
}
func Test_CALCULATE_AND_INSERT_VALUE_PRICE_QUANTITY(t *testing.T) {
	// a := CALCULATE_AND_INSERT_VALUE_PRICE_QUANTITY()
	// TEST(true,a,)
}
func Test_INSERT_IF_NOT_ZERO(t *testing.T) {
	// a := INSERT_IF_NOT_ZERO()
	// TEST(true,a,)
}
func Test_REMOVE_THE_ACCOUNTS(t *testing.T) {
	// a := REMOVE_THE_ACCOUNTS()
	// TEST(true,a,)
}
func Test_CHECK_DEBIT_EQUAL_CREDIT(t *testing.T) {
	// a := CHECK_DEBIT_EQUAL_CREDIT()
	// TEST(true,a,)
}
func Test_JOURNAL_ENTRY(t *testing.T) {
	// a := JOURNAL_ENTRY()
	// TEST(true,a,)
}
