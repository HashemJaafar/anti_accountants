package anti_accountants

import (
	"errors"
	"fmt"
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

	e1 = []SInvoiceEntry{{nil, nil, "1", "Revenue item 1", 5, 1, CDiscountTotal, 2, 4}}
	a1, a2, a3 = FEntryInvoice("cash", "Invoice discount", 1, 2500, e1, SEntry1{Labels: []string{"Invoice"}}, true)
	FTest(true, a1, e1)
	FTest(true, a2, nil)
	FTest(true, a3, nil)

	e1 = []SInvoiceEntry{{nil, nil, "1", "Revenue item 1", 0, 0, "", 0, 4}}
	a1, a2, a3 = FEntryInvoice("cash", "Invoice discount", 1, 2500, e1, SEntry1{Labels: []string{"Invoice"}}, true)
	FTest(true, a1, e1)
	FTest(true, a2, nil)
	FTest(true, a3, nil)

	e1 = []SInvoiceEntry{{nil, nil, "1", "Revenue item 2", 0, 0, "", 0, 4}}
	a1, a2, a3 = FEntryInvoice("cash", "Invoice discount", 1, 2500, e1, SEntry1{Labels: []string{"Invoice"}}, true)
	FTest(true, a1, e1)
	FTest(true, a2, nil)
	FTest(true, a3, nil)

	e1 = []SInvoiceEntry{{nil, errors.New("you order 30 and you have 16"), "1", "Revenue item 2", 0, 0, "", 0, 30}}
	a1, a2, a3 = FEntryInvoice("cash", "Invoice discount", 1, 2500, e1, SEntry1{Labels: []string{"Invoice"}}, true)
	fmt.Println(a1[0].QuantityError)
	FTest(true, a1, e1)
	FTest(true, a2, nil)
	FTest(true, a3, nil)
}

func TestFEntryReconciliationWithAccount(t *testing.T) {
	var err error

	err = FEntryReconciliationWithAccount(SAPQ1{
		AccountName: "rent",
		Price:       250,
		Quantity:    6,
	}, "cash", 1, SEntry1{Labels: []string{"ReconciliationWithAccount"}}, true)
	FTest(true, err, nil)

	err = FEntryReconciliationWithAccount(SAPQ1{
		AccountName: "rent",
		Price:       250,
		Quantity:    3,
	}, "cash", 1, SEntry1{Labels: []string{"ReconciliationWithAccount"}}, true)
	FTest(true, err, nil)

	err = FEntryReconciliationWithAccount(SAPQ1{
		AccountName: "book",
		Price:       1,
		Quantity:    100,
	}, "cash", 1, SEntry1{Labels: []string{"ReconciliationWithAccount"}}, true)
	FTest(true, err, nil)

	err = FEntryReconciliationWithAccount(SAPQ1{
		AccountName: "book",
		Price:       1,
		Quantity:    200,
	}, "cash", 1, SEntry1{Labels: []string{"ReconciliationWithAccount"}}, true)
	FTest(true, err, nil)

	err = FEntryReconciliationWithAccount(SAPQ1{
		AccountName: "book",
		Price:       1,
		Quantity:    150,
	}, "cash", 1, SEntry1{Labels: []string{"ReconciliationWithAccount"}}, true)
	FTest(true, err, nil)

	err = FEntryReconciliationWithAccount(SAPQ1{
		AccountName: "book",
		Price:       2,
		Quantity:    300000000,
	}, "cash", 1, SEntry1{Labels: []string{"ReconciliationWithAccount"}}, true)
	FTest(true, err, errors.New("you don't have enough cash to reconcile"))

	err = FEntryReconciliationWithAccount(SAPQ1{
		AccountName: "book",
		Price:       0,
		Quantity:    0,
	}, "cash", 1, SEntry1{Labels: []string{"ReconciliationWithAccount"}}, true)
	FTest(true, err, errors.New("you don't have enough cash to reconcile"))
}

func TestFEntryReconciliation(t *testing.T) {
	FTest(true, FEntryReconciliation("book", 200, SEntry1{Labels: []string{"Reconciliation"}}, true), nil)
	FTest(false, FEntryReconciliation("b", 100, SEntry1{Labels: []string{"Reconciliation"}}, true), nil)
	FTest(true, FEntryReconciliation("book", 300, SEntry1{Labels: []string{"Reconciliation"}}, true), nil)
	FTest(true, FEntryReconciliation("book", 50, SEntry1{Labels: []string{"Reconciliation"}}, true), nil)
	FTest(true, FEntryReconciliation("book", 2000, SEntry1{Labels: []string{"Reconciliation"}}, true), nil)
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
}

func TestMakeJournal(t *testing.T) {
	VDbJournal.DropAll()
	VDbInventory.DropAll()

	TestFEntryAutoComplete(t)
	TestFEntryInvoice(t)
	TestFEntryReconciliationWithAccount(t)
	TestFEntryReconciliation(t)
	TestFEntryClose(t)
}
