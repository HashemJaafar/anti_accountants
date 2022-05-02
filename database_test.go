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

func TestDbInsert(t *testing.T) {
	DbInventory.DropAll()
	DbInsert(DbInventory, []InventoryTag{
		{1, 10, "book"},
		{2, 10, "book"},
		{3, 10, "book"},
		{4, 10, "book"},
		{1, 10, "cash"},
		{1, 10, "cash"},
		{2, 10, "rent"},
		{9, 10, "rent"},
	})
	_, inventory := DbRead[InventoryTag](DbInventory)
	for _, v1 := range inventory {
		fmt.Println(v1)
	}
	DbClose()
}

func TestDbInsertIntoAccounts(t *testing.T) {
}

func TestDbLastLine(t *testing.T) {
	a1 := DbLastLine[JournalTag](DbJournal)
	fmt.Println(a1)
}

func TestDbRead(t *testing.T) {
	_, inventory := DbRead[InventoryTag](DbInventory)
	_, journal := DbRead[JournalTag](DbJournal)
	DbClose()
	PrintSlice(inventory)
	PrintSlice(journal)
}

func TestDbUpdate(t *testing.T) {
	DbUpdate(DbInventory, Now(), InventoryTag{1, 10, "book1"})
}

func TestWeightedAverage(t *testing.T) {
	WeightedAverage("rent")
	_, inventory := DbRead[InventoryTag](DbInventory)
	DbClose()
	PrintSlice(inventory)
	//e1:=
	//Test(true,a1,e1)
}

func TestChangeName(t *testing.T) {
	ChangeName("zizi", "zaid")
}

func TestChangeNameByEntryNumberCompund(t *testing.T) {
	ChangeNameByEntryNumberCompund(11, "zaid")
}
