//this file created automatically
package main

import (
	"fmt"
	"testing"
)

func TestCheckDebitEqualCredit(t *testing.T) {
	i1 := []PriceQuantityAccount{
		{false, "", "book", 1, 10},
		{false, "", "cash", 1, 10},
		{true, "", "rent", 1, 10},
		{true, "", "rent", 1, 10},
	}
	a1, a2, a3 := CheckDebitEqualCredit(i1)
	PrintSlice(a1)
	PrintSlice(a2)
	fmt.Println(a3)
}

func TestSetPriceAndQuantity(t *testing.T) {
	_, inventory := DbRead[Inventory](DbInventory)
	PrintSlice(inventory)
	i1 := PriceQuantityAccount{false, Wma, "rent", 0, -1}
	a1 := SetPriceAndQuantity(i1, true)
	fmt.Println(a1)
	_, inventory = DbRead[Inventory](DbInventory)
	PrintSlice(inventory)
	DbClose()
}

func TestGroupByAccount(t *testing.T) {
	i1 := []PriceQuantityAccount{
		{false, Lifo, "book", 1, 10},
		{false, Lifo, "book", 5, 10},
		{false, Lifo, "book", 3, 10},
		{true, Wma, "rent", 1, 10},
		{false, Wma, "cash", 1, 10},
	}
	a1 := GroupByAccount(i1)
	e1 := []PriceQuantityAccount{
		{false, Lifo, "book", 3, 30},
		{true, Wma, "rent", 1, 10},
		{false, Wma, "cash", 1, 10},
	}
	Test(true, a1, e1)

}
func TestSimpleJournalEntry(t *testing.T) {
	i1 := []PriceQuantityAccountBarcode{
		{1, 1000, "cash", ""},
		{1, 1000, "rent", ""},
	}
	a1, a2 := SimpleJournalEntry(i1, true, false, false, "ksdfjpaodka", "yasa", "hashem", "invoice")

	i1 = []PriceQuantityAccountBarcode{
		{1, 1000, "cash", ""},
		{1, 1000, "rent", ""},
	}
	a1, a2 = SimpleJournalEntry(i1, true, false, false, "ksdfjpaodka", "yasa", "hashem", "invoice")

	i1 = []PriceQuantityAccountBarcode{
		{1, -400, "cash", ""},
		{2, 200, "book", ""},
	}
	a1, a2 = SimpleJournalEntry(i1, true, false, false, "ksdfjpaodka", "yasa", "hashem", "payment")

	i1 = []PriceQuantityAccountBarcode{
		{1, -350, "cash", ""},
		{1.4, 250, "book", ""},
	}
	a1, a2 = SimpleJournalEntry(i1, true, false, false, "ksdfjpaodka", "yasa", "hashem", "payment")

	i1 = []PriceQuantityAccountBarcode{
		{1, 10 * 1.6666666666666667, "cash", ""},
		{1, -10, "book", ""},
	}
	a1, a2 = SimpleJournalEntry(i1, true, false, false, "ksdfjpaodka", "yasa", "hashem", "invoice")

	i1 = []PriceQuantityAccountBarcode{
		{1, 36, "cash", ""},
		{1, -18, "book", ""},
	}
	a1, a2 = SimpleJournalEntry(i1, true, false, false, "ksdfjpaodka", "zizi", "hashem", "invoice")
	DbClose()
	PrintFormatedAccounts()
	PrintSlice(a1)
	Test(true, a2, nil)
}

func TestStage1(t *testing.T) {
	PrintFormatedAccounts()
	i1 := []PriceQuantityAccountBarcode{
		{1, 10, "cash", "2"},
		{1, 10, "book", "1"},
		{1, 10, "cash", ""},
		{0, 10, "cash", ""},
		{10, 0, "cash", ""},
		{10, 10, "ca", ""},
	}
	a1 := Stage1(i1)
	e1 := []PriceQuantityAccount{
		{false, Lifo, "book", 1, 10},
		{true, Wma, "rent", 1, 10},
		{false, Wma, "cash", 1, 10},
	}
	Test(true, a1, e1)
}

func TestReverseEntries(t *testing.T) {
	ReverseEntries(5, 0, "hashem")
	DbClose()
}

func TestConvertPriceQuantityAccountToPriceQuantityAccountBarcode(t *testing.T) {
	a1 := ConvertPriceQuantityAccountToPriceQuantityAccountBarcode([]PriceQuantityAccount{{
		IsCredit:     false,
		CostFlowType: "",
		AccountName:  "cash",
		Price:        5,
		Quantity:     8,
	}})
	e1 := []PriceQuantityAccountBarcode{{
		Price:       5,
		Quantity:    8,
		AccountName: "cash",
		Barcode:     "",
	}}
	Test(true, a1, e1)
}

func TestFindDuplicateElement(t *testing.T) {
	keys, journal := DbRead[Journal](DbJournal)
	dates := ConvertByteSliceToTime(keys)
	a1, a2 := FindDuplicateElement(dates, journal, FilterJournalDuplicate{
		IsReverse:                  false,
		IsReversed:                 false,
		ReverseEntryNumberCompound: false,
		ReverseEntryNumberSimple:   false,
		Value:                      false,
		PriceDebit:                 false,
		PriceCredit:                false,
		QuantityDebit:              false,
		QuantityCredit:             false,
		AccountDebit:               false,
		AccountCredit:              false,
		Notes:                      false,
		Name:                       false,
		NameEmployee:               false,
	})
	PrintSlice(a1)
	PrintSlice(a2)
}

func TestJournalFilter(t *testing.T) {
	keys, journal := DbRead[Journal](DbJournal)
	dates := ConvertByteSliceToTime(keys)
	i1 := FilterJournal{
		Date:                       FilterDate{},
		IsReverse:                  FilterBool{},
		IsReversed:                 FilterBool{},
		ReverseEntryNumberCompound: FilterNumber{},
		ReverseEntryNumberSimple:   FilterNumber{},
		EntryNumberCompound:        FilterNumber{},
		EntryNumberSimple:          FilterNumber{},
		Value:                      FilterNumber{},
		PriceDebit:                 FilterNumber{},
		PriceCredit:                FilterNumber{},
		QuantityDebit:              FilterNumber{},
		QuantityCredit:             FilterNumber{IsFilter: false, Way: NotBetween, Big: 999, Small: 0},
		AccountDebit:               FilterString{},
		AccountCredit:              FilterString{},
		Notes:                      FilterString{},
		Name:                       FilterString{},
		NameEmployee:               FilterString{},
		TypeOfCompoundEntry:        FilterString{IsFilter: true, Way: InSlice, Slice: []string{"payment"}},
	}
	a1, a2 := JournalFilter(dates, journal, i1, true)
	PrintSlice(a1)
	PrintSlice(a2)
}
