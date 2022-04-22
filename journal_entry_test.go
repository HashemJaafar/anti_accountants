//this file created automatically
package main

import (
	"fmt"
	"testing"
	"time"
)

func TestCHECK_DEBIT_EQUAL_CREDIT(t *testing.T) {
	i1 := []PRICE_QUANTITY_ACCOUNT{
		{false, "", "book", 1, 10},
		{false, "", "cash", 1, 10},
		{true, "", "rent", 1, 10},
		{true, "", "rent", 1, 10},
	}
	a1, a2, a3 := CHECK_DEBIT_EQUAL_CREDIT(i1)
	PRINT_SLICE(a1)
	PRINT_SLICE(a2)
	fmt.Println(a3)
	// e1:=
	// e2:=
	// e3:=
	// TEST(true,a1,e1)
	// TEST(true,a2,e2)
	// TEST(true,a3,e3)
}

func TestSET_PRICE_AND_QUANTITY(t *testing.T) {
	_, inventory := DB_READ[INVENTORY_TAG](DB_INVENTORY)
	PRINT_SLICE(inventory)
	i1 := PRICE_QUANTITY_ACCOUNT{false, WMA, "rent", 0, -1}
	a1 := SET_PRICE_AND_QUANTITY(i1, true)
	fmt.Println(a1)
	_, inventory = DB_READ[INVENTORY_TAG](DB_INVENTORY)
	PRINT_SLICE(inventory)
	DB_CLOSE()
	//e1:=
	//TEST(true,a1,e1)
}

func TestGROUP_BY_ACCOUNT(t *testing.T) {
	i1 := []PRICE_QUANTITY_ACCOUNT{
		{false, LIFO, "book", 1, 10},
		{false, LIFO, "book", 5, 10},
		{false, LIFO, "book", 3, 10},
		{true, WMA, "rent", 1, 10},
		{false, WMA, "cash", 1, 10},
	}
	a1 := GROUP_BY_ACCOUNT(i1)
	e1 := []PRICE_QUANTITY_ACCOUNT{
		{false, LIFO, "book", 3, 30},
		{true, WMA, "rent", 1, 10},
		{false, WMA, "cash", 1, 10},
	}
	TEST(true, a1, e1)
}
func TestSIMPLE_JOURNAL_ENTRY(t *testing.T) {
	i1 := []PRICE_QUANTITY_ACCOUNT_BARCODE{
		{1, 1000, "cash", ""},
		{1, 1000, "rent", ""},
	}
	a1, a2 := SIMPLE_JOURNAL_ENTRY(i1, true, false, false, "ksdfjpaodka", "yasa", "hashem")
	i1 = []PRICE_QUANTITY_ACCOUNT_BARCODE{
		{1, 1000, "cash", ""},
		{1, 1000, "rent", ""},
	}
	a1, a2 = SIMPLE_JOURNAL_ENTRY(i1, true, false, false, "ksdfjpaodka", "yasa", "hashem")
	i1 = []PRICE_QUANTITY_ACCOUNT_BARCODE{
		{1, -400, "cash", ""},
		{2, 200, "book", ""},
	}
	a1, a2 = SIMPLE_JOURNAL_ENTRY(i1, true, false, false, "ksdfjpaodka", "yasa", "hashem")
	i1 = []PRICE_QUANTITY_ACCOUNT_BARCODE{
		{1, -350, "cash", ""},
		{1.4, 250, "book", ""},
	}
	a1, a2 = SIMPLE_JOURNAL_ENTRY(i1, true, false, false, "ksdfjpaodka", "yasa", "hashem")
	i1 = []PRICE_QUANTITY_ACCOUNT_BARCODE{
		{1, 10 * 1.6666666666666667, "cash", ""},
		{1, -10, "book", ""},
	}
	a1, a2 = SIMPLE_JOURNAL_ENTRY(i1, true, false, false, "ksdfjpaodka", "yasa", "hashem")
	i1 = []PRICE_QUANTITY_ACCOUNT_BARCODE{
		{1, 30, "cash", ""},
		{1, -18, "book", ""},
	}
	a1, a2 = SIMPLE_JOURNAL_ENTRY(i1, true, false, false, "ksdfjpaodka", "zizi", "hashem")
	DB_CLOSE()
	PRINT_FORMATED_ACCOUNTS()
	PRINT_SLICE(a1)
	fmt.Println(a2)
	//e1:=
	//TEST(true,a1,e1)
}

func TestSTAGE_1(t *testing.T) {
	PRINT_FORMATED_ACCOUNTS()
	i1 := []PRICE_QUANTITY_ACCOUNT_BARCODE{
		{1, 10, "cash", "2"},
		{1, 10, "book", "1"},
		{1, 10, "cash", ""},
		{0, 10, "cash", ""},
		{10, 0, "cash", ""},
		{10, 10, "ca", ""},
	}
	a1 := STAGE_1(i1)
	e1 := []PRICE_QUANTITY_ACCOUNT{
		{false, LIFO, "book", 1, 10},
		{true, WMA, "rent", 1, 10},
		{false, WMA, "cash", 1, 10},
	}
	TEST(true, a1, e1)
}

func TestREVERSE_ENTRIES(t *testing.T) {
	REVERSE_ENTRIES(2, 1, "hashem")
	DB_CLOSE()
}

func TestJOURNAL_FILTER(t *testing.T) {
	i1 := THE_JOURNAL_FILTER{
		JUST_BETWEEN_DATE:                          false,
		JUST_BETWEEN_REVERSE_ENTRY_NUMBER_COMPOUND: false,
		JUST_BETWEEN_REVERSE_ENTRY_NUMBER_SIMPLE:   false,
		JUST_BETWEEN_ENTRY_NUMBER_COMPOUND:         false,
		JUST_BETWEEN_ENTRY_NUMBER_SIMPLE:           false,
		JUST_BETWEEN_VALUE:                         false,
		JUST_BETWEEN_PRICE_DEBIT:                   false,
		JUST_BETWEEN_PRICE_CREDIT:                  false,
		JUST_BETWEEN_QUANTITY_DEBIT:                false,
		JUST_BETWEEN_QUANTITY_CREDIT:               false,
		IS_IN_SLICE_IS_REVERSE:                     false,
		IS_IN_SLICE_IS_REVERSED:                    false,
		IS_IN_SLICE_ACCOUNT_DEBIT:                  false,
		IS_IN_SLICE_ACCOUNT_CREDIT:                 false,
		IS_IN_SLICE_NOTES:                          false,
		IS_IN_SLICE_NAME:                           false,
		IS_IN_SLICE_NAME_EMPLOYEE:                  false,
		ABOVE_DATE:                                 time.Time{},
		ABOVE_REVERSE_ENTRY_NUMBER_COMPOUND:        0,
		ABOVE_REVERSE_ENTRY_NUMBER_SIMPLE:          0,
		ABOVE_ENTRY_NUMBER_COMPOUND:                0,
		ABOVE_ENTRY_NUMBER_SIMPLE:                  0,
		ABOVE_VALUE:                                0,
		ABOVE_PRICE_DEBIT:                          0,
		ABOVE_PRICE_CREDIT:                         0,
		ABOVE_QUANTITY_DEBIT:                       0,
		ABOVE_QUANTITY_CREDIT:                      0,
		BELLOW_DATE:                                time.Time{},
		BELLOW_REVERSE_ENTRY_NUMBER_COMPOUND:       0,
		BELLOW_REVERSE_ENTRY_NUMBER_SIMPLE:         0,
		BELLOW_ENTRY_NUMBER_COMPOUND:               0,
		BELLOW_ENTRY_NUMBER_SIMPLE:                 0,
		BELLOW_VALUE:                               0,
		BELLOW_PRICE_DEBIT:                         0,
		BELLOW_PRICE_CREDIT:                        0,
		BELLOW_QUANTITY_DEBIT:                      0,
		BELLOW_QUANTITY_CREDIT:                     0,
		SLICE_IS_REVERSE:                           []bool{},
		SLICE_IS_REVERSED:                          []bool{},
		SLICE_ACCOUNT_DEBIT:                        []string{},
		SLICE_ACCOUNT_CREDIT:                       []string{},
		SLICE_NOTES:                                []string{},
		SLICE_NAME:                                 []string{},
		SLICE_NAME_EMPLOYEE:                        []string{},
	}
	a1 := JOURNAL_FILTER(i1)
	PRINT_SLICE(a1)
}
