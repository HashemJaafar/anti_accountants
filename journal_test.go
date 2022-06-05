package main

import (
	"fmt"
	"testing"
)

func TestCheckDebitEqualCredit(t *testing.T) {
	i1 := []SAPQAE{
		{"book", 1, 10, SAccount{TIsCredit: false}, nil},
		{"cash", 1, 10, SAccount{TIsCredit: false}, nil},
		{"rent", 1, 10, SAccount{TIsCredit: true}, nil},
		{"rent", 1, 10, SAccount{TIsCredit: true}, nil},
	}
	a1, a2, a3 := FCheckDebitEqualCredit(i1)
	FPrintSlice(a1)
	FPrintSlice(a2)
	fmt.Println(a3)
}

func TestSetPriceAndQuantity(t *testing.T) {
	_, inventory := FDbRead[SAPQ](VDbInventory)
	FPrintSlice(inventory)
	i1 := SAPQAE{"rent", 0, -1, SAccount{TIsCredit: false}, nil}
	a1 := FSetPriceAndQuantity(i1, true)
	fmt.Println(a1)
	_, inventory = FDbRead[SAPQ](VDbInventory)
	FPrintSlice(inventory)
	FDbCloseAll()
}

func TestGroupByAccount(t *testing.T) {
	i1 := []SAPQAE{
		{"book", 1, 10, SAccount{TIsCredit: false, TCostFlowType: CLifo}, nil},
		{"book", 5, 10, SAccount{TIsCredit: false, TCostFlowType: CLifo}, nil},
		{"book", 3, 10, SAccount{TIsCredit: false, TCostFlowType: CLifo}, nil},
		{"rent", 1, 10, SAccount{TIsCredit: true, TCostFlowType: CWma}, nil},
		{"cash", 1, 10, SAccount{TIsCredit: false, TCostFlowType: CWma}, nil},
	}
	a1 := FGroupByAccount(i1)
	e1 := []SAPQAE{
		{"book", 3, 30, SAccount{TIsCredit: false, TCostFlowType: CLifo}, nil},
		{"rent", 1, 10, SAccount{TIsCredit: true, TCostFlowType: CWma}, nil},
		{"cash", 1, 10, SAccount{TIsCredit: false, TCostFlowType: CWma}, nil},
	}
	FTest(true, a1, e1)

}
func TestSimpleJournalEntry(t *testing.T) {
	VCompanyName = "anti_accountants"
	FDbOpenAll()
	var i1 []SAPQ
	var a1 []SAPQAE
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
	a1 := FSetEntries(i1, false)
	e1 := []SAPQAE{
		{"book", 1, 10, SAccount{TIsCredit: false, TCostFlowType: CLifo}, nil},
		{"rent", 1, 10, SAccount{TIsCredit: true, TCostFlowType: CWma}, nil},
		{"cash", 1, 10, SAccount{TIsCredit: false, TCostFlowType: CWma}, nil},
	}
	FTest(true, a1, e1)
}

func TestConvertPriceQuantityAccountToPriceQuantityAccountBarcode(t *testing.T) {
	a1 := FConvertAPQICToAPQB([]SAPQAE{{
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
	i1 := SFilterJournal{}
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

	a1 := FAutoComplete([]SAPQAE{{
		TAccountName: "book",
		TPrice:       0,
		TQuantity:    -1,
		SAccount:     SAccount{},
	}}, "cash", 1)
	e1 := [][]SAPQAE{
		{{"book", 0, -1, SAccount{TIsCredit: false}, nil}, {CPrefixCost + "book", 0, 1, SAccount{TIsCredit: false}, nil}},
		{{CPrefixRevenue + "book", 12, 1, SAccount{TIsCredit: true}, nil}, {CPrefixCost + "book", 5, 1, SAccount{TIsCredit: false}, nil}, {"cash", 1, 7, SAccount{TIsCredit: false}, nil}},
	}
	FTest(true, a1, e1)

	FDbCloseAll()
	FPrintSlice(a1)
}

func TestFReverseEntries(t *testing.T) {
	FDbOpenAll()
	keys, journal := FDbRead[SJournal](VDbJournal)
	FReverseEntries(keys, journal, "hashem")
	keys, journal = FDbRead[SJournal](VDbJournal)
	FPrintJournal(journal)
	FDbCloseAll()
}
