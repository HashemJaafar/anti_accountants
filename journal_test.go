package main

import (
	"fmt"
	"testing"
)

func TestCheckDebitEqualCredit(t *testing.T) {
	i1 := []SAPQA{
		{"book", 1, 10, SAccount{TIsCredit: false}},
		{"cash", 1, 10, SAccount{TIsCredit: false}},
		{"rent", 1, 10, SAccount{TIsCredit: true}},
		{"rent", 1, 10, SAccount{TIsCredit: true}},
	}
	a1, a2, a3 := FCheckDebitEqualCredit(i1)
	FPrintSlice(a1)
	FPrintSlice(a2)
	fmt.Println(a3)
}

func TestSetPriceAndQuantity(t *testing.T) {
	_, inventory := FDbRead[SAPQ](VDbInventory)
	FPrintSlice(inventory)
	i1 := SAPQA{"rent", 0, -1, SAccount{TIsCredit: false}}
	a1 := FSetPriceAndQuantity(i1, true)
	fmt.Println(a1)
	_, inventory = FDbRead[SAPQ](VDbInventory)
	FPrintSlice(inventory)
	FDbCloseAll()
}

func TestGroupByAccount(t *testing.T) {
	i1 := []SAPQA{
		{"book", 1, 10, SAccount{TIsCredit: false, TCostFlowType: CLifo}},
		{"book", 5, 10, SAccount{TIsCredit: false, TCostFlowType: CLifo}},
		{"book", 3, 10, SAccount{TIsCredit: false, TCostFlowType: CLifo}},
		{"rent", 1, 10, SAccount{TIsCredit: true, TCostFlowType: CWma}},
		{"cash", 1, 10, SAccount{TIsCredit: false, TCostFlowType: CWma}},
	}
	a1 := FGroupByAccount(i1)
	e1 := []SAPQA{
		{"book", 3, 30, SAccount{TIsCredit: false, TCostFlowType: CLifo}},
		{"rent", 1, 10, SAccount{TIsCredit: true, TCostFlowType: CWma}},
		{"cash", 1, 10, SAccount{TIsCredit: false, TCostFlowType: CWma}},
	}
	FTest(true, a1, e1)

}
func TestSimpleJournalEntry(t *testing.T) {
	var i1 []SAPQ
	var a1 []SAPQ
	var a2 error

	i1 = []SAPQ{
		{"cash", 1, 1000},
		{"rent", 1, 1000},
	}
	a1, a2 = FSimpleJournalEntry(i1, SEntry{"ksdfjpaodka", "yasa", "hashem", "invoice"}, true)
	FPrintSlice(a1)
	FTest(true, a2, nil)

	i1 = []SAPQ{
		{"cash", 1, 1000},
		{"rent", 1, 1000},
	}
	a1, a2 = FSimpleJournalEntry(i1, SEntry{"ksdfjpaodka", "yasa", "hashem", "invoice"}, true)
	FPrintSlice(a1)
	FTest(true, a2, nil)

	i1 = []SAPQ{
		{"cash", 1, -400},
		{"book", 2, 200},
	}
	a1, a2 = FSimpleJournalEntry(i1, SEntry{"ksdfjpaodka", "yasa", "hashem", "payment"}, true)
	FPrintSlice(a1)
	FTest(true, a2, nil)

	i1 = []SAPQ{
		{"cash", 1, -350},
		{"book", 1.4, 250},
	}
	a1, a2 = FSimpleJournalEntry(i1, SEntry{"ksdfjpaodka", "yasa", "hashem", "payment"}, true)
	FPrintSlice(a1)
	FTest(true, a2, nil)

	i1 = []SAPQ{
		{"cash", 1, 20},
		{"book", 1, -10},
	}
	a1, a2 = FSimpleJournalEntry(i1, SEntry{"ksdfjpaodka", "yasa", "hashem", "invoice"}, true)
	FPrintSlice(a1)
	FTest(true, a2, nil)

	i1 = []SAPQ{
		{"cash", 1, 36},
		{"book", 1, -18},
	}
	a1, a2 = FSimpleJournalEntry(i1, SEntry{"ksdfjpaodka", "zizi", "hashem", "invoice"}, true)
	FPrintSlice(a1)
	FTest(true, a2, nil)

	i1 = []SAPQ{
		{"cash", 1, 20},
		{"book", 1, -10},
	}
	a1, a2 = FSimpleJournalEntry(i1, SEntry{"ksdfjpaodka", "yasa", "hashem", "invoice"}, true)
	FPrintSlice(a1)
	FTest(true, a2, nil)

	FDbCloseAll()
	FPrintFormatedAccounts()
}

func TestStage1(t *testing.T) {
	FPrintFormatedAccounts()
	i1 := []SAPQ{
		{"cash", 1, 10},
		{"book", 1, 10},
		{"cash", 1, 10},
		{"cash", 0, 10},
		{"cash", 10, 0},
		{"ca", 10, 10},
	}
	a1 := FStage1(i1, false)
	e1 := []SAPQA{
		{"book", 1, 10, SAccount{TIsCredit: false, TCostFlowType: CLifo}},
		{"rent", 1, 10, SAccount{TIsCredit: true, TCostFlowType: CWma}},
		{"cash", 1, 10, SAccount{TIsCredit: false, TCostFlowType: CWma}},
	}
	FTest(true, a1, e1)
}

func TestReverseEntries(t *testing.T) {
	keys, journal := FFindEntryFromNumber(8, 0)
	FReverseEntries(keys, journal, "hashem")
	FDbCloseAll()
}

func TestConvertPriceQuantityAccountToPriceQuantityAccountBarcode(t *testing.T) {
	a1 := FConvertAPQICToAPQB([]SAPQA{{
		TAccountName: "cash",
		TPrice:       5,
		TQuantity:    8,
		SAccount:     SAccount{},
	}})
	e1 := []SAPQ{{"cash", 5, 8}}
	FTest(true, a1, e1)
}

func TestFindDuplicateElement(t *testing.T) {
	keys, journal := FDbRead[SJournal](VDbJournal)
	dates := FConvertByteSliceToTime(keys)
	a1, a2 := FFindDuplicateElement(dates, journal, SFilterJournalDuplicate{
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
	FPrintSlice(a1)
	FPrintSlice(a2)
}

func TestFJournalFilter(t *testing.T) {
	keys, journal := FDbRead[SJournal](VDbJournal)
	dates := FConvertByteSliceToTime(keys)
	i1 := SFilterJournal{
		Date:                       SFilterDate{},
		IsReverse:                  SFilterBool{},
		IsReversed:                 SFilterBool{},
		ReverseEntryNumberCompound: SFilterNumber{},
		ReverseEntryNumberSimple:   SFilterNumber{},
		EntryNumberCompound:        SFilterNumber{},
		EntryNumberSimple:          SFilterNumber{},
		Value:                      SFilterNumber{},
		PriceDebit:                 SFilterNumber{},
		PriceCredit:                SFilterNumber{},
		QuantityDebit:              SFilterNumber{},
		QuantityCredit:             SFilterNumber{IsFilter: false, Way: CNotBetween, Big: 999, Small: 0},
		AccountDebit:               SFilterString{},
		AccountCredit:              SFilterString{},
		Notes:                      SFilterString{},
		Name:                       SFilterString{},
		Employee:                   SFilterString{},
		TypeOfCompoundEntry:        SFilterString{IsFilter: true, Way: CInSlice, Slice: []string{"payment"}},
	}
	a1, a2 := FJournalFilter(dates, journal, i1, true)
	FPrintSlice(a1)
	FPrintSlice(a2)
}

func TestInvoiceJournalEntry(t *testing.T) {
	VAutoCompletionEntries = []SAutoCompletion{{
		TAccountName: "book",
		PriceRevenue: 12,
		PriceTax:     0,
		PriceDiscount: []SPQ{{
			TPrice:    5,
			TQuantity: 2,
		}},
	}}

	a1, a2 := FInvoiceJournalEntry("cash", 1, 4, []SAPQ{{"book", 5, -10}}, SEntry{}, true)
	_, inventory := FDbRead[SAPQ](VDbInventory)
	_, journal := FDbRead[SJournal](VDbJournal)
	FDbCloseAll()
	FPrintFormatedAccounts()
	FPrintSlice(inventory)
	FPrintSlice(journal)
	FPrintSlice(a1)
	FTest(true, a2, nil)
}

func TestAutoComplete(t *testing.T) {
	VAutoCompletionEntries = []SAutoCompletion{
		{"book", 12, 0, []SPQ{{5, 2}}},
	}

	a1 := FAutoComplete([]SAPQA{{
		TAccountName: "book",
		TPrice:       0,
		TQuantity:    -1,
		SAccount:     SAccount{},
	}}, "cash", 1)
	e1 := [][]SAPQA{
		{{"book", 0, -1, SAccount{TIsCredit: false}}, {CPrefixCost + "book", 0, 1, SAccount{TIsCredit: false}}},
		{{CPrefixRevenue + "book", 12, 1, SAccount{TIsCredit: true}}, {CPrefixCost + "book", 5, 1, SAccount{TIsCredit: false}}, {"cash", 1, 7, SAccount{TIsCredit: false}}},
	}
	FTest(true, a1, e1)

	FDbCloseAll()
	FPrintSlice(a1)
}

func TestConvertJournalToAPQA(t *testing.T) {
	_, journal := FFindEntryFromNumber(8, 0)
	FDbCloseAll()
	a1 := FConvertJournalToAPQA(journal)
	FPrintSlice(a1)
	FPrintJournal(journal)
}
