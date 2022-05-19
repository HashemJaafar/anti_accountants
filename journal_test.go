package main

import (
	"fmt"
	"testing"
)

func TestCheckDebitEqualCredit(t *testing.T) {
	i1 := []SAPQA{
		{"book", 1, 10, SAccount{IsCredit: false}},
		{"cash", 1, 10, SAccount{IsCredit: false}},
		{"rent", 1, 10, SAccount{IsCredit: true}},
		{"rent", 1, 10, SAccount{IsCredit: true}},
	}
	a1, a2, a3 := FCheckDebitEqualCredit(i1)
	FPrintSlice(a1)
	FPrintSlice(a2)
	fmt.Println(a3)
}

func TestSetPriceAndQuantity(t *testing.T) {
	_, inventory := FDbRead[SAPQ](VDbInventory)
	FPrintSlice(inventory)
	i1 := SAPQA{"rent", 0, -1, SAccount{IsCredit: false}}
	a1 := FSetPriceAndQuantity(i1, true)
	fmt.Println(a1)
	_, inventory = FDbRead[SAPQ](VDbInventory)
	FPrintSlice(inventory)
	FDbClose()
}

func TestGroupByAccount(t *testing.T) {
	i1 := []SAPQA{
		{"book", 1, 10, SAccount{IsCredit: false, CostFlowType: CLifo}},
		{"book", 5, 10, SAccount{IsCredit: false, CostFlowType: CLifo}},
		{"book", 3, 10, SAccount{IsCredit: false, CostFlowType: CLifo}},
		{"rent", 1, 10, SAccount{IsCredit: true, CostFlowType: CWma}},
		{"cash", 1, 10, SAccount{IsCredit: false, CostFlowType: CWma}},
	}
	a1 := FGroupByAccount(i1)
	e1 := []SAPQA{
		{"book", 3, 30, SAccount{IsCredit: false, CostFlowType: CLifo}},
		{"rent", 1, 10, SAccount{IsCredit: true, CostFlowType: CWma}},
		{"cash", 1, 10, SAccount{IsCredit: false, CostFlowType: CWma}},
	}
	FTest(true, a1, e1)

}
func TestSimpleJournalEntry(t *testing.T) {
	var i1 []SAPQB
	var a1 []SAPQB
	var a2 error

	i1 = []SAPQB{
		{"cash", 1, 1000, ""},
		{"rent", 1, 1000, ""},
	}
	a1, a2 = FSimpleJournalEntry(i1, SEntryInfo{"ksdfjpaodka", "yasa", "hashem", "invoice"}, true)
	FPrintSlice(a1)
	FTest(true, a2, nil)

	i1 = []SAPQB{
		{"cash", 1, 1000, ""},
		{"rent", 1, 1000, ""},
	}
	a1, a2 = FSimpleJournalEntry(i1, SEntryInfo{"ksdfjpaodka", "yasa", "hashem", "invoice"}, true)
	FPrintSlice(a1)
	FTest(true, a2, nil)

	i1 = []SAPQB{
		{"cash", 1, -400, ""},
		{"book", 2, 200, ""},
	}
	a1, a2 = FSimpleJournalEntry(i1, SEntryInfo{"ksdfjpaodka", "yasa", "hashem", "payment"}, true)
	FPrintSlice(a1)
	FTest(true, a2, nil)

	i1 = []SAPQB{
		{"cash", 1, -350, ""},
		{"book", 1.4, 250, ""},
	}
	a1, a2 = FSimpleJournalEntry(i1, SEntryInfo{"ksdfjpaodka", "yasa", "hashem", "payment"}, true)
	FPrintSlice(a1)
	FTest(true, a2, nil)

	i1 = []SAPQB{
		{"cash", 1, 20, ""},
		{"book", 1, -10, ""},
	}
	a1, a2 = FSimpleJournalEntry(i1, SEntryInfo{"ksdfjpaodka", "yasa", "hashem", "invoice"}, true)
	FPrintSlice(a1)
	FTest(true, a2, nil)

	i1 = []SAPQB{
		{"cash", 1, 36, ""},
		{"book", 1, -18, ""},
	}
	a1, a2 = FSimpleJournalEntry(i1, SEntryInfo{"ksdfjpaodka", "zizi", "hashem", "invoice"}, true)
	FPrintSlice(a1)
	FTest(true, a2, nil)

	i1 = []SAPQB{
		{"cash", 1, 20, ""},
		{"book", 1, -10, ""},
	}
	a1, a2 = FSimpleJournalEntry(i1, SEntryInfo{"ksdfjpaodka", "yasa", "hashem", "invoice"}, true)
	FPrintSlice(a1)
	FTest(true, a2, nil)

	FDbClose()
	FPrintFormatedAccounts()
}

func TestStage1(t *testing.T) {
	FPrintFormatedAccounts()
	i1 := []SAPQB{
		{"cash", 1, 10, "2"},
		{"book", 1, 10, "1"},
		{"cash", 1, 10, ""},
		{"cash", 0, 10, ""},
		{"cash", 10, 0, ""},
		{"ca", 10, 10, ""},
	}
	a1 := FStage1(i1, false)
	e1 := []SAPQA{
		{"book", 1, 10, SAccount{IsCredit: false, CostFlowType: CLifo}},
		{"rent", 1, 10, SAccount{IsCredit: true, CostFlowType: CWma}},
		{"cash", 1, 10, SAccount{IsCredit: false, CostFlowType: CWma}},
	}
	FTest(true, a1, e1)
}

func TestReverseEntries(t *testing.T) {
	keys, journal := FFindEntryFromNumber(8, 0)
	FReverseEntries(keys, journal, "hashem")
	FDbClose()
}

func TestConvertPriceQuantityAccountToPriceQuantityAccountBarcode(t *testing.T) {
	a1 := FConvertAPQICToAPQB([]SAPQA{{
		Name:     "cash",
		Price:    5,
		Quantity: 8,
		Account:  SAccount{},
	}})
	e1 := []SAPQB{{"cash", 5, 8, ""}}
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

func TestJournalFilter(t *testing.T) {
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

func TestValueAfterAdjustUsingAdjustingMethods(t *testing.T) {
	a1 := FValueAfterAdjustUsingAdjustingMethods("", 2, 100, 10, 100)
	fmt.Println(a1)
}

func TestInvoiceJournalEntry(t *testing.T) {
	VAutoCompletionEntries = []SAutoCompletion{{
		AccountInvnetory: "book",
		PriceRevenue:     12,
		PriceTax:         0,
		PriceDiscount: []SPQ{{
			Price:    5,
			Quantity: 2,
		}},
	}}

	a1, a2 := FInvoiceJournalEntry("cash", 1, 4, []SAPQB{{"book", 5, -10, ""}}, SEntryInfo{}, true)
	_, inventory := FDbRead[SAPQ](VDbInventory)
	_, journal := FDbRead[SJournal](VDbJournal)
	FDbClose()
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
		Name:     "book",
		Price:    0,
		Quantity: -1,
		Account:  SAccount{},
	}}, "cash", 1)
	e1 := [][]SAPQA{
		{{"book", 0, -1, SAccount{IsCredit: false}}, {CPrefixCost + "book", 0, 1, SAccount{IsCredit: false}}},
		{{CPrefixRevenue + "book", 12, 1, SAccount{IsCredit: true}}, {CPrefixCost + "book", 5, 1, SAccount{IsCredit: false}}, {"cash", 1, 7, SAccount{IsCredit: false}}},
	}
	FTest(true, a1, e1)

	FDbClose()
	FPrintSlice(a1)
}

func TestFilterJournalFromReverseEntry(t *testing.T) {
	keys, journal := FDbRead[SJournal](VDbJournal)
	a1, a2 := FFilterJournalFromReverseEntry(keys, journal)
	FDbClose()
	FPrintSlice(a1)
	FPrintSlice(a2)
}

func TestConvertJournalToAPQA(t *testing.T) {
	_, journal := FFindEntryFromNumber(8, 0)
	FDbClose()
	a1 := FConvertJournalToAPQA(journal)
	FPrintSlice(a1)
	FPrintJournal(journal)
}
