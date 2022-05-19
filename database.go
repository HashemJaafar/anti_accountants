package main

import (
	"encoding/json"

	badger "github.com/dgraph-io/badger/v3"
)

func DbClose() {
	DbAccounts.Close()
	DbJournal.Close()
	DbInventory.Close()
	DbAutoCompletionEntries.Close()
}

func DbInsertIntoAccounts() {
	DbAccounts.DropAll()
	for _, v1 := range Accounts {
		DbUpdate(DbAccounts, []byte(v1.Name), v1)
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
	keys, inventory := DbRead[APQ](DbInventory)
	for k1, v1 := range inventory {
		if v1.Name == old {
			v1.Name = new
			DbUpdate(DbJournal, keys[k1], v1)
		}
	}
}

func WeightedAverage(account string) {
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

	keys, inventory := DbRead[APQ](DbInventory)
	for k1, v1 := range inventory {
		if v1.Name == account {
			DbDelete(DbInventory, keys[k1])
		}
	}

	DbUpdate(DbInventory, Now(), APQ{account, Abs(totalValue / totalQuantity), Abs(totalQuantity)})
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

func ChangeName(old, new string) {
	keys, journal := DbRead[Journal](DbJournal)
	for k1, v1 := range journal {
		if v1.Name == old {
			v1.Name = new
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

func ChangeEntryInfoByEntryNumberCompund(entryNumberCompund int, new EntryInfo) {
	keys, journal := DbRead[Journal](DbJournal)
	for k1, v1 := range journal {
		if v1.EntryNumberCompound == entryNumberCompund {
			v1.Notes = new.Notes
			v1.Name = new.Name
			v1.TypeOfCompoundEntry = new.TypeOfCompoundEntry
			DbUpdate(DbJournal, keys[k1], v1)
		}
	}
}
