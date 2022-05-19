package main

import (
	"fmt"
	"testing"
)

func TestDbClose(t *testing.T) {
	FTest(true, VDbAccounts.IsClosed(), false)
	FTest(true, VDbJournal.IsClosed(), false)
	FTest(true, VDbInventory.IsClosed(), false)
	FDbClose()
	FTest(true, VDbAccounts.IsClosed(), true)
	FTest(true, VDbJournal.IsClosed(), true)
	FTest(true, VDbInventory.IsClosed(), true)
}

func TestDbLastLine(t *testing.T) {
	a1 := FDbLastLine[SJournal](VDbJournal)
	fmt.Println(a1)
}

func TestDbRead(t *testing.T) {
	_, inventory := FDbRead[SAPQ](VDbInventory)
	_, journal := FDbRead[SJournal](VDbJournal)
	_, VAutoCompletionEntries = FDbRead[SAutoCompletion](VDbAutoCompletionEntries)
	FDbClose()
	FPrintSlice(inventory)
	FPrintJournal(journal)
	FPrintSlice(VAutoCompletionEntries)
}

func TestDbUpdate(t *testing.T) {
	FDbUpdate(VDbInventory, FNow(), SAPQ{"book1", 1, 10})
}

func TestWeightedAverage(t *testing.T) {
	FWeightedAverage("rent")
	_, inventory := FDbRead[SAPQ](VDbInventory)
	FDbClose()
	FPrintSlice(inventory)
}

func TestChangeName(t *testing.T) {
	FChangeName("zizi", "zaid")
}
