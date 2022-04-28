package main

import (
	"fmt"
	"testing"
)

func TestDB_CLOSE(t *testing.T) {
	TEST(true, DB_ACCOUNTS.IsClosed(), false)
	TEST(true, DB_JOURNAL.IsClosed(), false)
	TEST(true, DB_INVENTORY.IsClosed(), false)
	DB_CLOSE()
	TEST(true, DB_ACCOUNTS.IsClosed(), true)
	TEST(true, DB_JOURNAL.IsClosed(), true)
	TEST(true, DB_INVENTORY.IsClosed(), true)
}

func TestDB_INSERT(t *testing.T) {
	DB_INVENTORY.DropAll()
	DB_INSERT(DB_INVENTORY, []INVENTORY_TAG{
		{1, 10, "book"},
		{2, 10, "book"},
		{3, 10, "book"},
		{4, 10, "book"},
		{1, 10, "cash"},
		{1, 10, "cash"},
		{2, 10, "rent"},
		{9, 10, "rent"},
	})
	_, inventory := DB_READ[INVENTORY_TAG](DB_INVENTORY)
	for _, v1 := range inventory {
		fmt.Println(v1)
	}
	DB_CLOSE()
}

func TestDB_INSERT_INTO_ACCOUNTS(t *testing.T) {
}

func TestDB_LAST_LINE(t *testing.T) {
	a1 := DB_LAST_LINE[JOURNAL_TAG](DB_JOURNAL)
	fmt.Println(a1)
	// e1 := JOURNAL_TAG{false, 0, 0, 0, 0, 0, 0, 0, 0, "", "", "", "", "", time.Time{}, time.Time{}}
	// TEST(true, a1, e1)
}

func TestDB_READ(t *testing.T) {
	_, inventory := DB_READ[INVENTORY_TAG](DB_INVENTORY)
	_, journal := DB_READ[JOURNAL_TAG](DB_JOURNAL)
	DB_CLOSE()
	PRINT_SLICE(inventory)
	PRINT_SLICE(journal)
}

func TestDB_UPDATE(t *testing.T) {
	DB_UPDATE(DB_INVENTORY, NOW(), INVENTORY_TAG{1, 10, "book1"})
}

func TestWEIGHTED_AVERAGE(t *testing.T) {
	WEIGHTED_AVERAGE("rent")
	_, inventory := DB_READ[INVENTORY_TAG](DB_INVENTORY)
	DB_CLOSE()
	PRINT_SLICE(inventory)
	//e1:=
	//TEST(true,a1,e1)
}
