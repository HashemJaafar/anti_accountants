//this file created automatically
package main

import (
	"fmt"
	"testing"
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
	// i1 := []PRICE_QUANTITY_ACCOUNT_BARCODE{
	// 	{1, 1000, "cash", ""},
	// 	{1, 1000, "rent", ""},
	// }
	// a1, a2 := SIMPLE_JOURNAL_ENTRY(i1, true, false, false, "ksdfjpaodka", "yasa", "hashem")
	// i1 = []PRICE_QUANTITY_ACCOUNT_BARCODE{
	// 	{1, 1000, "cash", ""},
	// 	{1, 1000, "rent", ""},
	// }
	// a1, a2 = SIMPLE_JOURNAL_ENTRY(i1, true, false, false, "ksdfjpaodka", "yasa", "hashem")
	// i1 = []PRICE_QUANTITY_ACCOUNT_BARCODE{
	// 	{1, -400, "cash", ""},
	// 	{2, 200, "book", ""},
	// }
	// a1, a2 = SIMPLE_JOURNAL_ENTRY(i1, true, false, false, "ksdfjpaodka", "yasa", "hashem")
	// i1 = []PRICE_QUANTITY_ACCOUNT_BARCODE{
	// 	{1, -350, "cash", ""},
	// 	{1.4, 250, "book", ""},
	// }
	// a1, a2 = SIMPLE_JOURNAL_ENTRY(i1, true, false, false, "ksdfjpaodka", "yasa", "hashem")
	// i1 = []PRICE_QUANTITY_ACCOUNT_BARCODE{
	// 	{1, 10 * 1.6666666666666667, "cash", ""},
	// 	{1, -10, "book", ""},
	// }
	// a1, a2 = SIMPLE_JOURNAL_ENTRY(i1, true, false, false, "ksdfjpaodka", "yasa", "hashem")
	i1 := []PRICE_QUANTITY_ACCOUNT_BARCODE{
		{1, 30, "cash", ""},
		{1, -18, "book", ""},
	}
	a1, a2 := SIMPLE_JOURNAL_ENTRY(i1, true, false, false, "ksdfjpaodka", "zizi", "hashem")
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
	// i1 := THE_JOURNAL_FILTER{
	// }
	// a1 := JOURNAL_FILTER(i1)
	// PRINT_SLICE(a1)
}

func TestCONVERT_PRICE_QUANTITY_ACCOUNT_TO_PRICE_QUANTITY_ACCOUNT_BARCODE(t *testing.T) {
	a1 := CONVERT_PRICE_QUANTITY_ACCOUNT_TO_PRICE_QUANTITY_ACCOUNT_BARCODE([]PRICE_QUANTITY_ACCOUNT{{
		IS_CREDIT:      false,
		COST_FLOW_TYPE: "",
		ACCOUNT_NAME:   "cash",
		PRICE:          5,
		QUANTITY:       8,
	}})
	e1 := []PRICE_QUANTITY_ACCOUNT_BARCODE{{
		PRICE:        5,
		QUANTITY:     8,
		ACCOUNT_NAME: "cash",
		BARCODE:      "",
	}}
	TEST(true, a1, e1)
}
