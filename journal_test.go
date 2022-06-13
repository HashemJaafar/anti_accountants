package main

import (
	"fmt"
	"testing"
)

func TestCheckDebitEqualCredit(t *testing.T) {
	i1 := []SAPQAE{
		{"book", 1, 10, SAccount1{IsCredit: false}, nil},
		{"cash", 1, 10, SAccount1{IsCredit: false}, nil},
		{"rent", 1, 10, SAccount1{IsCredit: true}, nil},
		{"rent", 1, 10, SAccount1{IsCredit: true}, nil},
	}
	a1, a2, a3 := FCheckDebitEqualCredit(i1)
	FPrintSlice(a1)
	FPrintSlice(a2)
	fmt.Println(a3)
}

func TestSetPriceAndQuantity(t *testing.T) {
	_, inventory := FDbRead[SAPQ](VDbInventory)
	FPrintSlice(inventory)
	i1 := SAPQAE{"rent", 0, -1, SAccount1{IsCredit: false}, nil}
	a1 := FSetPriceAndQuantity(i1, true)
	fmt.Println(a1)
	_, inventory = FDbRead[SAPQ](VDbInventory)
	FPrintSlice(inventory)
	FDbCloseAll()
}

func TestGroupByAccount(t *testing.T) {
	i1 := []SAPQAE{
		{"book", 1, 10, SAccount1{IsCredit: false, CostFlowType: CLifo}, nil},
		{"book", 5, 10, SAccount1{IsCredit: false, CostFlowType: CLifo}, nil},
		{"book", 3, 10, SAccount1{IsCredit: false, CostFlowType: CLifo}, nil},
		{"rent", 1, 10, SAccount1{IsCredit: true, CostFlowType: CWma}, nil},
		{"cash", 1, 10, SAccount1{IsCredit: false, CostFlowType: CWma}, nil},
	}
	a1 := FGroupByAccount(i1)
	e1 := []SAPQAE{
		{"book", 3, 30, SAccount1{IsCredit: false, CostFlowType: CLifo}, nil},
		{"rent", 1, 10, SAccount1{IsCredit: true, CostFlowType: CWma}, nil},
		{"cash", 1, 10, SAccount1{IsCredit: false, CostFlowType: CWma}, nil},
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
		{"book", 1, 10, SAccount1{IsCredit: false, CostFlowType: CLifo}, nil},
		{"rent", 1, 10, SAccount1{IsCredit: true, CostFlowType: CWma}, nil},
		{"cash", 1, 10, SAccount1{IsCredit: false, CostFlowType: CWma}, nil},
	}
	FTest(true, a1, e1)
}

func TestConvertPriceQuantityAccountToPriceQuantityAccountBarcode(t *testing.T) {
	a1 := FConvertAPQICToAPQB([]SAPQAE{{
		TAccountName: "cash",
		TPrice:       5,
		TQuantity:    8,
		SAccount1:    SAccount1{},
	}})
	e1 := []SAPQ{{"cash", 5, 8}}
	FTest(true, a1, e1)
}

func TestFindDuplicateElement(t *testing.T) {
	keys, journal := FDbRead[SJournal1](VDbJournal)
	dates := FConvertByteSliceToTime(keys)
	a1, a2 := FFindDuplicateElement(dates, journal, SJournal3{})
	FPrintSlice(a1)
	FPrintSlice(a2)
}

func TestFJournalFilter(t *testing.T) {
	_, journal := FDbRead[SJournal1](VDbJournal)
	i1 := SJournal2{}
	a1 := FJournalFilter(journal, i1, true)
	FPrintSlice(a1)
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
	_, journal := FDbRead[SJournal1](VDbJournal)
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
		SAccount1:    SAccount1{},
	}}, "cash", 1)
	e1 := [][]SAPQAE{
		{{"book", 0, -1, SAccount1{IsCredit: false}, nil}, {CPrefixCost + "book", 0, 1, SAccount1{IsCredit: false}, nil}},
		{{CPrefixRevenue + "book", 12, 1, SAccount1{IsCredit: true}, nil}, {CPrefixCost + "book", 5, 1, SAccount1{IsCredit: false}, nil}, {"cash", 1, 7, SAccount1{IsCredit: false}, nil}},
	}
	FTest(true, a1, e1)

	FDbCloseAll()
	FPrintSlice(a1)
}

func TestFReverseEntries(t *testing.T) {
	FDbOpenAll()
	keys, journal := FDbRead[SJournal1](VDbJournal)
	FReverseEntries(keys, journal, "hashem")
	keys, journal = FDbRead[SJournal1](VDbJournal)
	FPrintJournal(journal)
	FDbCloseAll()
}
