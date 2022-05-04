package main

import (
	"encoding/json"

	badger "github.com/dgraph-io/badger/v3"
)

func DbClose() {
	DbAccounts.Close()
	DbJournal.Close()
	DbInventory.Close()
}

func DbInsert[t any](db *badger.DB, slice []t) {
	for _, a := range slice {
		DbUpdate(db, Now(), a)
	}
}

func DbInsertIntoAccounts() {
	DbAccounts.DropAll()
	for _, a := range Accounts {
		DbUpdate(DbAccounts, []byte(a.AccountName), a)
	}
}

func DbLastLine[t any](db *badger.DB) t {
	var tag t
	var str []byte
	db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			item.Value(func(val []byte) error {
				str = val
				return nil
			})
		}
		return nil
	})
	json.Unmarshal(str, &tag)
	return tag
}

func DbOpen(path string) *badger.DB {
	for {
		db, err := badger.Open(badger.DefaultOptions(path))
		if err == nil {
			return db
		}
	}
}

func DbRead[t any](db *badger.DB) ([][]byte, []t) {
	var VALUEs []t
	var keys [][]byte
	db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			item.Value(func(val []byte) error {
				var Value t
				json.Unmarshal(val, &Value)
				VALUEs = append(VALUEs, Value)
				keys = append(keys, item.Key())
				return nil
			})
		}
		return nil
	})
	return keys, VALUEs
}

func DbDelete(db *badger.DB, key []byte) {
	db.Update(func(txn *badger.Txn) error {
		return txn.Delete(key)
	})
}

func DbUpdate[t any](db *badger.DB, key []byte, Value t) {
	txn := db.NewTransaction(true)
	defer txn.Commit()
	jsonValue, _ := json.Marshal(Value)
	txn.Set(key, jsonValue)
}

func ChangeAccountName(old, new string) {
	keys, journal := DbRead[Journal](DbJournal)
	for k1, v1 := range journal {
		if v1.AccountCredit == old {
			v1.AccountCredit = new
			DbUpdate(DbJournal, keys[k1], v1)
		}
		if v1.AccountDebit == old {
			v1.AccountDebit = new
			DbUpdate(DbJournal, keys[k1], v1)
		}
	}
	keys, inventory := DbRead[Inventory](DbInventory)
	for k1, v1 := range inventory {
		if v1.AccountName == old {
			v1.AccountName = new
			DbUpdate(DbJournal, keys[k1], v1)
		}
	}
}

func WeightedAverage(account string) {
	// i find the WMA from journal because it is the most accurate way when you enter reverese entry
	// i store the Value in var total but just to know the total should allways be positive
	// if it is negative that mean the nature of the account is debit
	// because i subtruct the debit from the total and sum the credit to the total

	// here i find the sum of the Value and sum of QUANTITY from journal
	var totalValue, totalQuantity float64
	_, journal := DbRead[Journal](DbJournal)
	for _, v1 := range journal {
		if v1.AccountCredit == account {
			totalValue += v1.Value
			totalQuantity += v1.QuantityCredit
		}
		if v1.AccountDebit == account {
			totalValue -= v1.Value
			totalQuantity -= v1.QuantityDebit
		}
	}

	// here i delete the account from the inventory
	keys, inventory := DbRead[Inventory](DbInventory)
	for k1, v1 := range inventory {
		if v1.AccountName == account {
			DbDelete(DbInventory, keys[k1])
		}
	}

	// here i insert the new account with the new PRICE and sum of the QUANTITY
	DbUpdate(DbInventory, Now(), Inventory{Abs(totalValue / totalQuantity), Abs(totalQuantity), account})
}

func ChangeNotes(old, new string) {
	keys, journal := DbRead[Journal](DbJournal)
	for k1, v1 := range journal {
		if v1.Notes == old {
			v1.Notes = new
			DbUpdate(DbJournal, keys[k1], v1)
		}
	}
}

func ChangeNotesByEntryNumberCompund(entryNumberCompund int, new string) {
	keys, journal := DbRead[Journal](DbJournal)
	for k1, v1 := range journal {
		if v1.EntryNumberCompound == entryNumberCompund {
			v1.Notes = new
			DbUpdate(DbJournal, keys[k1], v1)
		}
	}
}

func ChangeNameByEntryNumberCompund(entryNumberCompund int, new string) {
	keys, journal := DbRead[Journal](DbJournal)
	for k1, v1 := range journal {
		if v1.EntryNumberCompound == entryNumberCompund {
			v1.Name = new
			DbUpdate(DbJournal, keys[k1], v1)
		}
	}
}

func ChangeName(old, new string) {
	keys, journal := DbRead[Journal](DbJournal)
	for k1, v1 := range journal {
		if v1.Name == old {
			v1.Name = new
			DbUpdate(DbJournal, keys[k1], v1)
		}
	}
}

func ChangeTypeOfCompoundEntryByEntryNumberCompund(entryNumberCompund int, new string) {
	keys, journal := DbRead[Journal](DbJournal)
	for k1, v1 := range journal {
		if v1.EntryNumberCompound == entryNumberCompund {
			v1.TypeOfCompoundEntry = new
			DbUpdate(DbJournal, keys[k1], v1)
		}
	}
}

func ChangeTypeOfCompoundEntry(old, new string) {
	keys, journal := DbRead[Journal](DbJournal)
	for k1, v1 := range journal {
		if v1.TypeOfCompoundEntry == old {
			v1.TypeOfCompoundEntry = new
			DbUpdate(DbJournal, keys[k1], v1)
		}
	}
}
