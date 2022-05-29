package main

import (
	"fmt"
	"testing"
)

func TestDbClose(t *testing.T) {
	FTest(true, VDbAccounts.IsClosed(), false)
	FTest(true, VDbJournal.IsClosed(), false)
	FTest(true, VDbInventory.IsClosed(), false)
	FDbCloseAll()
	FTest(true, VDbAccounts.IsClosed(), true)
	FTest(true, VDbJournal.IsClosed(), true)
	FTest(true, VDbInventory.IsClosed(), true)
}

func TestDbLastLine(t *testing.T) {
	a1 := FDbLastLine[SJournal](VDbJournal)
	fmt.Println(a1)
}

func TestDbRead(t *testing.T) {
	FDbOpenAll()
	_, inventory := FDbRead[SAPQ](VDbInventory)
	_, journal := FDbRead[SJournal](VDbJournal)
	_, VAutoCompletionEntries = FDbRead[SAutoCompletion](VDbAutoCompletionEntries)
	FDbCloseAll()
	FPrintSlice(inventory)
	FPrintJournal(journal)
	FPrintSlice(VAutoCompletionEntries)
}

func TestDbUpdate(t *testing.T) {
	VCompanyName = "anti_accountants"
	VDbEmployees = FDbOpen(VDbEmployees, CPathDataBase+VCompanyName+CPathEmployees)
	FDbUpdate(VDbEmployees, []byte("hashem"), "hajasa")
	FDbCloseAll()
}

func TestWeightedAverage(t *testing.T) {
	FWeightedAverage("rent")
	_, inventory := FDbRead[SAPQ](VDbInventory)
	FDbCloseAll()
	FPrintSlice(inventory)
}

func TestChangeName(t *testing.T) {
	FChangeName("zizi", "zaid")
}

func TestFDbOpenAll(t *testing.T) {
	VCompanyName = "aisdj"
	FDbOpenAll()
	FDbCloseAll()
}

func TestFDbOpen(t *testing.T) {
	VDbAccounts = FDbOpen(VDbAccounts, CPathDataBase+VCompanyName+CPathAccounts)
	VDbAccounts = FDbOpen(VDbAccounts, CPathDataBase+VCompanyName+CPathAccounts)
	VDbAccounts = FDbOpen(VDbAccounts, CPathDataBase+VCompanyName+CPathAccounts)
}

func TestFDbRead(t *testing.T) {
	VCompanyName = "anti_accountants"
	VDbEmployees = FDbOpen(VDbEmployees, CPathDataBase+VCompanyName+CPathEmployees)
	keys, passwords := FDbRead[string](VDbEmployees)
	employees := FConvertFromByteSliceToString(keys)
	fmt.Println(employees)
	fmt.Println(passwords)
	FDbCloseAll()
}

func TestFDbClose(t *testing.T) {
	VDbAccounts = FDbOpen(VDbAccounts, CPathDataBase+VCompanyName+CPathAccounts)
	FDbClose(VDbAccounts)
	FDbClose(VDbAccounts)
}
