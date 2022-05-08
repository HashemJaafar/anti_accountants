package main

import (
	"fmt"
	"testing"
)

func TestDbClose(t *testing.T) {
	Test(true, DbAccounts.IsClosed(), false)
	Test(true, DbJournal.IsClosed(), false)
	Test(true, DbInventory.IsClosed(), false)
	DbClose()
	Test(true, DbAccounts.IsClosed(), true)
	Test(true, DbJournal.IsClosed(), true)
	Test(true, DbInventory.IsClosed(), true)
}

func TestDbLastLine(t *testing.T) {
	a1 := DbLastLine[Journal](DbJournal)
	fmt.Println(a1)
}

func TestDbRead(t *testing.T) {
	_, inventory := DbRead[APQ](DbInventory)
	_, journal := DbRead[Journal](DbJournal)
	_, AutoCompletionEntries = DbRead[AutoCompletion](DbAutoCompletionEntries)
	DbClose()
	PrintSlice(inventory)
	PrintJournal(journal)
	PrintSlice(AutoCompletionEntries)
}

func TestDbUpdate(t *testing.T) {
	DbUpdate(DbInventory, Now(), APQ{"book1", 1, 10})
}

func TestWeightedAverage(t *testing.T) {
	WeightedAverage("rent")
	_, inventory := DbRead[APQ](DbInventory)
	DbClose()
	PrintSlice(inventory)
}

func TestChangeName(t *testing.T) {
	ChangeName("zizi", "zaid")
}
