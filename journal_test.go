//this file created automatically
package main

import (
	"fmt"
	"testing"
)

func TestCheckDebitEqualCredit(t *testing.T) {
	i1 := []APQA{
		{"book", 1, 10, Account{IsCredit: false}},
		{"cash", 1, 10, Account{IsCredit: false}},
		{"rent", 1, 10, Account{IsCredit: true}},
		{"rent", 1, 10, Account{IsCredit: true}},
	}
	a1, a2, a3 := CheckDebitEqualCredit(i1)
	PrintSlice(a1)
	PrintSlice(a2)
	fmt.Println(a3)
}

func TestSetPriceAndQuantity(t *testing.T) {
	_, inventory := DbRead[APQ](DbInventory)
	PrintSlice(inventory)
	i1 := APQA{"rent", 0, -1, Account{IsCredit: false}}
	a1 := SetPriceAndQuantity(i1, true)
	fmt.Println(a1)
	_, inventory = DbRead[APQ](DbInventory)
	PrintSlice(inventory)
	DbClose()
}

func TestGroupByAccount(t *testing.T) {
	i1 := []APQA{
		{"book", 1, 10, Account{IsCredit: false, CostFlowType: Lifo}},
		{"book", 5, 10, Account{IsCredit: false, CostFlowType: Lifo}},
		{"book", 3, 10, Account{IsCredit: false, CostFlowType: Lifo}},
		{"rent", 1, 10, Account{IsCredit: true, CostFlowType: Wma}},
		{"cash", 1, 10, Account{IsCredit: false, CostFlowType: Wma}},
	}
	a1 := GroupByAccount(i1)
	e1 := []APQA{
		{"book", 3, 30, Account{IsCredit: false, CostFlowType: Lifo}},
		{"rent", 1, 10, Account{IsCredit: true, CostFlowType: Wma}},
		{"cash", 1, 10, Account{IsCredit: false, CostFlowType: Wma}},
	}
	Test(true, a1, e1)

}
func TestSimpleJournalEntry(t *testing.T) {
	var i1 []APQB
	var a1 []APQB
	var a2 error

	i1 = []APQB{
		{"cash", 1, 1000, ""},
		{"rent", 1, 1000, ""},
	}
	a1, a2 = SimpleJournalEntry(i1, EntryInfo{"ksdfjpaodka", "yasa", "hashem", "invoice"}, true)
	PrintSlice(a1)
	Test(true, a2, nil)

	i1 = []APQB{
		{"cash", 1, 1000, ""},
		{"rent", 1, 1000, ""},
	}
	a1, a2 = SimpleJournalEntry(i1, EntryInfo{"ksdfjpaodka", "yasa", "hashem", "invoice"}, true)
	PrintSlice(a1)
	Test(true, a2, nil)

	i1 = []APQB{
		{"cash", 1, -400, ""},
		{"book", 2, 200, ""},
	}
	a1, a2 = SimpleJournalEntry(i1, EntryInfo{"ksdfjpaodka", "yasa", "hashem", "payment"}, true)
	PrintSlice(a1)
	Test(true, a2, nil)

	i1 = []APQB{
		{"cash", 1, -350, ""},
		{"book", 1.4, 250, ""},
	}
	a1, a2 = SimpleJournalEntry(i1, EntryInfo{"ksdfjpaodka", "yasa", "hashem", "payment"}, true)
	PrintSlice(a1)
	Test(true, a2, nil)

	i1 = []APQB{
		{"cash", 1, 20, ""},
		{"book", 1, -10, ""},
	}
	a1, a2 = SimpleJournalEntry(i1, EntryInfo{"ksdfjpaodka", "yasa", "hashem", "invoice"}, true)
	PrintSlice(a1)
	Test(true, a2, nil)

	i1 = []APQB{
		{"cash", 1, 36, ""},
		{"book", 1, -18, ""},
	}
	a1, a2 = SimpleJournalEntry(i1, EntryInfo{"ksdfjpaodka", "zizi", "hashem", "invoice"}, true)
	PrintSlice(a1)
	Test(true, a2, nil)

	i1 = []APQB{
		{"cash", 1, 20, ""},
		{"book", 1, -10, ""},
	}
	a1, a2 = SimpleJournalEntry(i1, EntryInfo{"ksdfjpaodka", "yasa", "hashem", "invoice"}, true)
	PrintSlice(a1)
	Test(true, a2, nil)

	DbClose()
	PrintFormatedAccounts()
}

func TestStage1(t *testing.T) {
	PrintFormatedAccounts()
	i1 := []APQB{
		{"cash", 1, 10, "2"},
		{"book", 1, 10, "1"},
		{"cash", 1, 10, ""},
		{"cash", 0, 10, ""},
		{"cash", 10, 0, ""},
		{"ca", 10, 10, ""},
	}
	a1 := Stage1(i1, false)
	e1 := []APQA{
		{"book", 1, 10, Account{IsCredit: false, CostFlowType: Lifo}},
		{"rent", 1, 10, Account{IsCredit: true, CostFlowType: Wma}},
		{"cash", 1, 10, Account{IsCredit: false, CostFlowType: Wma}},
	}
	Test(true, a1, e1)
}

func TestReverseEntries(t *testing.T) {
	ReverseEntries(8, 0, "hashem")
	DbClose()
}

func TestConvertPriceQuantityAccountToPriceQuantityAccountBarcode(t *testing.T) {
	a1 := ConvertAPQICToAPQB([]APQA{{
		Name:     "cash",
		Price:    5,
		Quantity: 8,
		Account:  Account{},
	}})
	e1 := []APQB{{"cash", 5, 8, ""}}
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
		Employee:                   false,
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
		Employee:                   FilterString{},
		TypeOfCompoundEntry:        FilterString{IsFilter: true, Way: InSlice, Slice: []string{"payment"}},
	}
	a1, a2 := JournalFilter(dates, journal, i1, true)
	PrintSlice(a1)
	PrintSlice(a2)
}

func TestValueAfterAdjustUsingAdjustingMethods(t *testing.T) {
	a1 := ValueAfterAdjustUsingAdjustingMethods("", 2, 100, 10, 100)
	fmt.Println(a1)
}

func TestInvoiceJournalEntry(t *testing.T) {
	AutoCompletionEntries = []AutoCompletion{{
		AccountInvnetory: "book",
		PriceRevenue:     12,
		PriceTax:         0,
		PriceDiscount: []Discount{{
			Price:    5,
			Quantity: 2,
		}},
	}}

	a1, a2 := InvoiceJournalEntry("cash", 1, 4, []APQB{{"book", 5, -10, ""}}, EntryInfo{}, true)
	_, inventory := DbRead[APQ](DbInventory)
	_, journal := DbRead[Journal](DbJournal)
	DbClose()
	PrintFormatedAccounts()
	PrintSlice(inventory)
	PrintSlice(journal)
	PrintSlice(a1)
	Test(true, a2, nil)
}

func TestAutoComplete(t *testing.T) {
	AutoCompletionEntries = []AutoCompletion{
		{"book", 12, 0, []Discount{{5, 2}}},
	}

	a1 := AutoComplete([]APQA{{
		Name:     "book",
		Price:    0,
		Quantity: -1,
		Account:  Account{},
	}}, "cash", 1)
	e1 := [][]APQA{
		{{"book", 0, -1, Account{IsCredit: false}}, {PrefixCost + "book", 0, 1, Account{IsCredit: false}}},
		{{PrefixRevenue + "book", 12, 1, Account{IsCredit: true}}, {PrefixCost + "book", 5, 1, Account{IsCredit: false}}, {"cash", 1, 7, Account{IsCredit: false}}},
	}
	Test(true, a1, e1)

	DbClose()
	PrintSlice(a1)
}

func TestFilterJournalFromReverseEntry(t *testing.T) {
	keys, journal := DbRead[Journal](DbJournal)
	a1, a2 := FilterJournalFromReverseEntry(keys, journal)
	DbClose()
	PrintSlice(a1)
	PrintSlice(a2)
}
