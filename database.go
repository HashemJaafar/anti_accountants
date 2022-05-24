package main

import (
	"encoding/json"

	badger "github.com/dgraph-io/badger/v3"
)

func FDbClose() {
	VDbAccounts.Close()
	VDbJournal.Close()
	VDbInventory.Close()
	VDbAutoCompletionEntries.Close()
}

func FDbInsertIntoAccounts() {
	VDbAccounts.DropAll()
	for _, v1 := range VAccounts {
		FDbUpdate(VDbAccounts, []byte(v1.TAccountName), v1)
	}
}

func FDbLastLine[t any](db *badger.DB) t {
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

func FDbOpen(path string) *badger.DB {
	for {
		db, err := badger.Open(badger.DefaultOptions(path))
		if err == nil {
			return db
		}
	}
}

func FDbRead[t any](db *badger.DB) ([][]byte, []t) {
	var Values []t
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
				Values = append(Values, Value)
				keys = append(keys, item.Key())
				return nil
			})
		}
		return nil
	})
	return keys, Values
}

func FDbDelete(db *badger.DB, key []byte) {
	db.Update(func(txn *badger.Txn) error {
		return txn.Delete(key)
	})
}

func FDbUpdate[t any](db *badger.DB, key []byte, Value t) {
	txn := db.NewTransaction(true)
	defer txn.Commit()
	jsonValue, _ := json.Marshal(Value)
	txn.Set(key, jsonValue)
}

func FChangeAccountName(old, new string) {
	keys, journal := FDbRead[SJournal](VDbJournal)
	for k1, v1 := range journal {
		if v1.AccountCredit == old {
			v1.AccountCredit = new
			FDbUpdate(VDbJournal, keys[k1], v1)
		}
		if v1.AccountDebit == old {
			v1.AccountDebit = new
			FDbUpdate(VDbJournal, keys[k1], v1)
		}
	}
	keys, inventory := FDbRead[SAPQ](VDbInventory)
	for k1, v1 := range inventory {
		if v1.TAccountName == old {
			v1.TAccountName = new
			FDbUpdate(VDbJournal, keys[k1], v1)
		}
	}
}

func FWeightedAverage(account string) {
	var totalValue, totalQuantity float64
	_, journal := FDbRead[SJournal](VDbJournal)
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

	keys, inventory := FDbRead[SAPQ](VDbInventory)
	for k1, v1 := range inventory {
		if v1.TAccountName == account {
			FDbDelete(VDbInventory, keys[k1])
		}
	}

	FDbUpdate(VDbInventory, FNow(), SAPQ{account, FAbs(totalValue / totalQuantity), FAbs(totalQuantity)})
}

func FChangeNotes(old, new string) {
	keys, journal := FDbRead[SJournal](VDbJournal)
	for k1, v1 := range journal {
		if v1.Notes == old {
			v1.Notes = new
			FDbUpdate(VDbJournal, keys[k1], v1)
		}
	}
}

func FChangeName(old, new string) {
	keys, journal := FDbRead[SJournal](VDbJournal)
	for k1, v1 := range journal {
		if v1.Name == old {
			v1.Name = new
			FDbUpdate(VDbJournal, keys[k1], v1)
		}
	}
}

func FChangeTypeOfCompoundEntry(old, new string) {
	keys, journal := FDbRead[SJournal](VDbJournal)
	for k1, v1 := range journal {
		if v1.TypeOfCompoundEntry == old {
			v1.TypeOfCompoundEntry = new
			FDbUpdate(VDbJournal, keys[k1], v1)
		}
	}
}

func FChangeEntryInfoByEntryNumberCompund(entryNumberCompund int, new SEntry) {
	keys, journal := FDbRead[SJournal](VDbJournal)
	for k1, v1 := range journal {
		if v1.EntryNumberCompound == entryNumberCompund {
			v1.Notes = new.TEntryNotes
			v1.Name = new.TPersonName
			v1.TypeOfCompoundEntry = new.TTypeOfCompoundEntry
			FDbUpdate(VDbJournal, keys[k1], v1)
		}
	}
}
