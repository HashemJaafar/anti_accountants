package anti_accountants

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
	a1 := FDbLastLine[SJournal1](VDbJournal)
	fmt.Println(a1)
}

func TestDbRead(t *testing.T) {
	_, inventory := FDbRead[SAPQ1](VDbInventory)
	_, journal := FDbRead[SJournal1](VDbJournal)
	_, VAutoCompletionEntries = FDbRead[SAutoCompletion](VDbAutoCompletionEntries)
	FPrintStructSlice(false, inventory)
	FPrintStructSlice(false, journal)
	FPrintStructSlice(false, VAutoCompletionEntries)
}

func TestDbUpdate(t *testing.T) {
	VDbEmployees = FDbOpen(VDbEmployees, CPathDataBase+VCompanyName+CPathEmployees)
	FDbUpdate(VDbEmployees, []byte("hashem"), "hajasa")
}

func TestWeightedAverage(t *testing.T) {
	FWeightedAverage("rent")
	_, inventory := FDbRead[SAPQ1](VDbInventory)
	FPrintStructSlice(false, inventory)
}

func TestChangeName(t *testing.T) {
	FChangeName("zizi", "zaid")
}

func TestFDbOpen(t *testing.T) {
	VDbAccounts = FDbOpen(VDbAccounts, CPathDataBase+VCompanyName+CPathAccounts)
	VDbAccounts = FDbOpen(VDbAccounts, CPathDataBase+VCompanyName+CPathAccounts)
	VDbAccounts = FDbOpen(VDbAccounts, CPathDataBase+VCompanyName+CPathAccounts)
}

func TestFDbRead(t *testing.T) {
	VDbEmployees = FDbOpen(VDbEmployees, CPathDataBase+VCompanyName+CPathEmployees)
	keys, passwords := FDbRead[string](VDbEmployees)
	employees := FConvertFromByteSliceToString(keys)
	fmt.Println(employees)
	fmt.Println(passwords)
}

func TestFDbClose(t *testing.T) {
	VDbAccounts = FDbOpen(VDbAccounts, CPathDataBase+VCompanyName+CPathAccounts)
	FDbClose(VDbAccounts)
	FDbClose(VDbAccounts)
}
