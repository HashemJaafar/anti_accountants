package anti_accountants

import (
	"errors"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Write code here to run before tests
	VCompanyName = "anti_accountants"
	FDbOpenAll()

	// Run tests
	exitVal := m.Run()

	// Write code here to run after tests
	_, inventory := FDbRead[SAPQ1](VDbInventory)
	_, journal := FDbRead[SJournal1](VDbJournal)
	FDbCloseAll()
	FPrintStructSlice(true, VAccounts)
	FPrintStructSlice(true, VAutoCompletionEntries)
	FPrintStructSlice(true, inventory)
	FPrintStructSlice(false, journal)

	// Exit with exit value from tests
	os.Exit(exitVal)
}

func TestFEntryInvoice(t *testing.T) {
	var a1 []SInvoiceEntry
	var a2, a3 error
	var e1 []SInvoiceEntry

	e1 = []SInvoiceEntry{{nil, nil, "1", "Revenue item 1", 5, 1, 2, CDiscountTotal, 4}}
	a1, a2, a3 = FEntryInvoice("cash", "Invoice discount", 1, 2500, e1, SEntry1{Labels: []string{"Invoice"}}, true)
	FTest(true, a1, e1)
	FTest(true, a2, nil)
	FTest(true, a3, nil)

	e1 = []SInvoiceEntry{{nil, nil, "1", "Revenue item 1", 0, 0, 0, CDiscountTotal, 4}}
	a1, a2, a3 = FEntryInvoice("cash", "Invoice discount", 1, 2500, e1, SEntry1{Labels: []string{"Invoice"}}, true)
	FTest(true, a1, e1)
	FTest(true, a2, nil)
	FTest(true, a3, nil)

	e1 = []SInvoiceEntry{{nil, nil, "1", "Revenue item 2", 0, 0, 0, CDiscountTotal, 4}}
	a1, a2, a3 = FEntryInvoice("cash", "Invoice discount", 1, 2500, e1, SEntry1{Labels: []string{"Invoice"}}, true)
	FTest(true, a1, e1)
	FTest(true, a2, nil)
	FTest(true, a3, nil)

	e1 = []SInvoiceEntry{{nil, errors.New("you order 30 and you have 16"), "1", "Revenue item 2", 0, 0, 0, CDiscountTotal, 30}}
	a1, a2, a3 = FEntryInvoice("cash", "Invoice discount", 1, 2500, e1, SEntry1{Labels: []string{"Invoice"}}, true)
	FTest(true, a1, e1)
	FTest(true, a2, nil)
	FTest(true, a3, nil)
}

func TestFEntryChangeQuantity(t *testing.T) {
	FTest(true, FEntryChangeQuantity("book", 200, SEntry1{Labels: []string{"ChangeQuantity"}}, true), nil)
	FTest(false, FEntryChangeQuantity("b", 100, SEntry1{Labels: []string{"ChangeQuantity"}}, true), nil)
	FTest(true, FEntryChangeQuantity("book", 300, SEntry1{Labels: []string{"ChangeQuantity"}}, true), nil)
	FTest(true, FEntryChangeQuantity("book", 50, SEntry1{Labels: []string{"ChangeQuantity"}}, true), nil)
	FTest(true, FEntryChangeQuantity("book", 2000, SEntry1{Labels: []string{"ChangeQuantity"}}, true), nil)
}

func TestFEntryClose(t *testing.T) {
	FTest(true, FEntryClose("book", "cash", 1, SEntry1{Labels: []string{"Close"}}, true), nil)
	FTest(true, FEntryClose("rent", "cash", 1, SEntry1{Labels: []string{"Close"}}, true), nil)
}

func TestFSetPriceAndQuantityByValue(t *testing.T) {
	a1 := FSetPriceAndQuantityByValue(SAPQ12SAccount1{
		SAPQ1: SAPQ1{
			AccountName: "book",
			Price:       1,
			Quantity:    -800,
		},
		SAPQ2:     SAPQ2{},
		SAccount1: SAccount1{CostFlowType: CFifo},
	})

	FTest(true, a1, SAPQ12SAccount1{
		SAPQ1:     SAPQ1{"book", 1.9875776397515528, -402.5},
		SAPQ2:     SAPQ2{"", "", ""},
		SAccount1: SAccount1{false, "Fifo", "", "", []string(nil), [][]uint(nil), []uint(nil), [][]string(nil)},
	})
}

func TestFEntryAutoComplete(t *testing.T) {
	var a2 error

	_, a2 = FEntryAutoComplete([]SAPQ1{
		{"cash", 1, 1000},
		{"rent", 1000, 1},
	}, SEntry1{Labels: []string{"AutoComplete"}}, true, "")
	FTest(true, a2, nil)

	_, a2 = FEntryAutoComplete([]SAPQ1{
		{"cash", 1, -400},
		{"book", 2, 200},
	}, SEntry1{Labels: []string{"AutoComplete"}}, true, "")
	FTest(true, a2, nil)

	_, a2 = FEntryAutoComplete([]SAPQ1{
		{"cash", 1, 1500},
		{"rent", 1500, 1500},
	}, SEntry1{Labels: []string{"AutoComplete"}}, true, "")
	FTest(false, a2, nil)

	_, a2 = FEntryAutoComplete([]SAPQ1{
		{"book", 1, -1},
		{"cash", 1, 2},
	}, SEntry1{Labels: []string{"AutoComplete"}}, true, "")
	FTest(true, a2, nil)

	_, a2 = FEntryAutoComplete([]SAPQ1{
		{"book", 1, 800},
		{"cash", 1, -5},
	}, SEntry1{Labels: []string{"AutoComplete"}}, true, "")
	FTest(false, a2, nil)

	_, a2 = FEntryAutoComplete([]SAPQ1{
		{"cash", 1, -5},
		{"book", 1, 800},
	}, SEntry1{Labels: []string{"AutoComplete"}}, true, "")
	FTest(false, a2, nil)

	_, a2 = FEntryAutoComplete([]SAPQ1{
		{"cash", 1, -285},
		{"revenue of book", 1, 285},
		{"cost of book", 1.5, 190},
		{"book", 1.5, 190},
	}, SEntry1{Labels: []string{"AutoComplete"}}, true, "")
	FTest(false, a2, nil)

	_, a2 = FEntryAutoComplete([]SAPQ1{
		{"cash", 1, 1},
		{"rent", 1000, 1},
	}, SEntry1{Labels: []string{"AutoComplete"}}, true, "cash")
	FTest(true, a2, nil)

	_, a2 = FEntryAutoComplete([]SAPQ1{
		{"cash", 1, -1},
		{"book", 2, 200},
	}, SEntry1{Labels: []string{"AutoComplete"}}, true, "cash")
	FTest(true, a2, nil)

	_, a2 = FEntryAutoComplete([]SAPQ1{
		{"cash", 1, 1500},
		{"rent", 1500, 0},
	}, SEntry1{Labels: []string{"AutoComplete"}}, true, "rent")
	FTest(true, a2, nil)

	_, a2 = FEntryAutoComplete([]SAPQ1{
		{"book", 1, -2},
		{"cash", 1, 2},
	}, SEntry1{Labels: []string{"AutoComplete"}}, true, "book")
	FTest(true, a2, nil)

	_, a2 = FEntryAutoComplete([]SAPQ1{
		{"book", 1, 800},
		{"cash", 1, -5},
	}, SEntry1{Labels: []string{"AutoComplete"}}, true, "book")
	FTest(true, a2, nil)

	_, a2 = FEntryAutoComplete([]SAPQ1{
		{"cash", 1, -5},
		{"book", 1, 800},
	}, SEntry1{Labels: []string{"AutoComplete"}}, true, "")
	FTest(false, a2, nil)

	_, a2 = FEntryAutoComplete([]SAPQ1{
		{"cash", 1, -5},
		{"Inventory item 1", 1, 20},
	}, SEntry1{Labels: []string{"AutoComplete"}}, true, "cash")
	FTest(true, a2, nil)

	_, a2 = FEntryAutoComplete([]SAPQ1{
		{"cash", 1, -5},
		{"Inventory item 2", 1, 20},
	}, SEntry1{Labels: []string{"AutoComplete"}}, true, "cash")
	FTest(true, a2, nil)

	_, a2 = FEntryAutoComplete([]SAPQ1{
		{"cash", 1, -5},
		{"cars", 1, 20},
	}, SEntry1{Labels: []string{"AutoComplete"}}, true, "cash")
	FTest(true, a2, nil)

	_, a2 = FEntryAutoComplete([]SAPQ1{
		{"cash", 1, -5},
		{"color", 1, 20},
	}, SEntry1{Labels: []string{"AutoComplete"}}, true, "cash")
	FTest(true, a2, nil)
}

func TestMakeJournal(t *testing.T) {
	VDbJournal.DropAll()
	VDbInventory.DropAll()

	TestFEntryAutoComplete(t)
	TestFEntryInvoice(t)
	TestFEntryChangeQuantity(t)
	TestFEntryClose(t)
	TestFEntryAddValueWithoutChangeQuantity(t)
}

func TestFEntryAddValueWithoutChangeQuantity(t *testing.T) {
	var a1 error

	a1 = FEntryAddValueWithoutChangeQuantity(SAPQ1{"color", 2, -4}, "cars", SEntry1{Labels: []string{"AddValueWithoutChangeQuantity"}}, false)
	FTest(true, a1, nil)
	a1 = FEntryAddValueWithoutChangeQuantity(SAPQ1{"color", 2, 4}, "cars", SEntry1{Labels: []string{"AddValueWithoutChangeQuantity"}}, false)
	FTest(true, a1, nil)
	a1 = FEntryAddValueWithoutChangeQuantity(SAPQ1{"color", 2, 50}, "cars", SEntry1{Labels: []string{"AddValueWithoutChangeQuantity"}}, false)
	FTest(true, a1, nil)
	a1 = FEntryAddValueWithoutChangeQuantity(SAPQ1{"color", 2, -50}, "cars", SEntry1{Labels: []string{"AddValueWithoutChangeQuantity"}}, false)
	FTest(true, a1, nil)
}

func TestFCompleteTheEntry(t *testing.T) {
	{
		i1 := []SAPQ12SAccount1{{
			SAPQ1:     SAPQ1{"cars", 100, 1},
			SAPQ2:     SAPQ2{},
			SAccount1: SAccount1{CostFlowType: CWma},
		}}
		i2 := []SAPQ12SAccount1{{
			SAPQ1:     SAPQ1{"color", 1, -20},
			SAPQ2:     SAPQ2{},
			SAccount1: SAccount1{CostFlowType: CWma},
		}}
		i3 := 100.0
		i4 := 20.0

		e1 := i1
		e2 := i2
		e3 := i3
		e4 := i4

		a1 := FCompleteTheEntry(&i1, &i2, &i3, &i4,
			SAPQ12SAccount1{
				SAPQ1:     SAPQ1{"cars", 100, -4},
				SAPQ2:     SAPQ2{},
				SAccount1: SAccount1{CostFlowType: CWma},
			},
		)

		FTest(true, a1, errors.New("you don't have enough quantity"))
		FTest(true, i1, e1)
		FTest(true, i2, e2)
		FTest(true, i3, e3)
		FTest(true, i4, e4)
	}

	{
		i1 := []SAPQ12SAccount1{{
			SAPQ1:     SAPQ1{"cars", 10, 1},
			SAPQ2:     SAPQ2{},
			SAccount1: SAccount1{CostFlowType: CWma},
		}}
		i2 := []SAPQ12SAccount1{{
			SAPQ1:     SAPQ1{"color", 1, -6},
			SAPQ2:     SAPQ2{},
			SAccount1: SAccount1{CostFlowType: CWma},
		}}
		i3 := 10.0
		i4 := 6.0

		e1 := i1
		e2 := []SAPQ12SAccount1{
			{SAPQ1{"color", 1, -6}, SAPQ2{"", "", ""}, SAccount1{false, "Wma", "", "", []string(nil), [][]uint(nil), []uint(nil), [][]string(nil)}},
			{SAPQ1{"cars", 1, -4}, SAPQ2{"", "", ""}, SAccount1{false, "Wma", "", "", []string(nil), [][]uint(nil), []uint(nil), [][]string(nil)}},
		}
		e3 := i3
		e4 := 10.0

		a1 := FCompleteTheEntry(&i1, &i2, &i3, &i4,
			SAPQ12SAccount1{
				SAPQ1:     SAPQ1{"cars", 10, -4},
				SAPQ2:     SAPQ2{},
				SAccount1: SAccount1{CostFlowType: CWma},
			},
		)

		FTest(true, a1, nil)
		FTest(true, i1, e1)
		FTest(true, i2, e2)
		FTest(true, i3, e3)
		FTest(true, i4, e4)
	}
}
