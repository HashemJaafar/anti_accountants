package main

import (
	"fmt"
	"testing"
)

func TestFInvoiceJournalEntry(t *testing.T) {
	VCompanyName = "anti_accountants"
	FDbOpenAll()
	a1, a2 := FInvoiceJournalEntry("cash", "Invoice PQ", 1, 2500, []SAPQ1{
		{"book", 0, 30},
	}, SEntry{}, true)

	FDbCloseAll()
	fmt.Println("a2:", a2)
	FPrintStructSlice(true, a1)
}

func TestFSimpleJournalEntry(t *testing.T) {
	VCompanyName = "anti_accountants"
	FDbOpenAll()
	var a1 []SAPQ12SAccount1
	var a2 error

	a1, a2 = FSimpleJournalEntry([]SAPQ1{
		{"cash", 1, 1000},
		{"rent", 1000, 1},
	}, SEntry{
		Notes:               "rent",
		Name:                "zaid",
		Employee:            "hashem",
		TypeOfCompoundEntry: "audited",
	}, true)
	FPrintStructSlice(false, a1)
	fmt.Println("a2:", a2)

	a1, a2 = FSimpleJournalEntry([]SAPQ1{
		{"cash", 1, -400},
		{"book", 2, 200},
	}, SEntry{
		Notes:               "rent",
		Name:                "zaid",
		Employee:            "hashem",
		TypeOfCompoundEntry: "audited",
	}, true)
	FPrintStructSlice(false, a1)
	fmt.Println("a2:", a2)

	FDbCloseAll()
}

func TestFSimpleJournalEntry1(t *testing.T) {
	VCompanyName = "anti_accountants"
	FDbOpenAll()
	a1, a2 := FSimpleJournalEntry([]SAPQ1{
		{"cash", 1, -400},
		{"book", 2, 200},
	}, SEntry{
		Notes:               "rent",
		Name:                "zaid",
		Employee:            "hashem",
		TypeOfCompoundEntry: "audited",
	}, true)
	FPrintStructSlice(false, a1)
	fmt.Println("a2:", a2)
	FDbCloseAll()
}
