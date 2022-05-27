package main

import (
	"encoding/json"
	"log"
	"strings"
	"sync"

	badger "github.com/dgraph-io/badger/v3"
)

func FDbOpenAll() {
	var wait sync.WaitGroup
	wait.Add(5)

	go func() {
		VDbAccounts = FDbOpen(VDbAccounts, CPathDataBase+VCompanyName+CPathAccounts)
		wait.Done()
	}()
	go func() {
		VDbJournal = FDbOpen(VDbJournal, CPathDataBase+VCompanyName+CPathJournal)
		wait.Done()
	}()
	go func() {
		VDbInventory = FDbOpen(VDbInventory, CPathDataBase+VCompanyName+CPathInventory)
		wait.Done()
	}()
	go func() {
		VDbAutoCompletionEntries = FDbOpen(VDbAutoCompletionEntries, CPathDataBase+VCompanyName+CPathAutoCompletionEntries)
		wait.Done()
	}()
	go func() {
		VDbEmployees = FDbOpen(VDbEmployees, CPathDataBase+VCompanyName+CPathEmployees)
		wait.Done()
	}()

	wait.Wait()

	wait.Add(2)
	go func() {
		_, VAccounts = FDbRead[SAccount](VDbAccounts)
		wait.Done()
	}()
	go func() {
		_, VAutoCompletionEntries = FDbRead[SAutoCompletion](VDbAutoCompletionEntries)
		wait.Done()
	}()
	wait.Wait()
}

func FDbCloseAll() {
	var wait sync.WaitGroup
	wait.Add(5)

	go func() {
		FDbClose(VDbAccounts)
		wait.Done()
	}()
	go func() {
		FDbClose(VDbJournal)
		wait.Done()
	}()
	go func() {
		FDbClose(VDbInventory)
		wait.Done()
	}()
	go func() {
		FDbClose(VDbAutoCompletionEntries)
		wait.Done()
	}()
	go func() {
		FDbClose(VDbEmployees)
		wait.Done()
	}()

	wait.Wait()
}

func FDbClose(db *badger.DB) {
	if db != nil {
		db.Close()
	}
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

func FDbOpen(database *badger.DB, path string) *badger.DB {
	for {
		db, err := badger.Open(badger.DefaultOptions(path))
		if err != nil {
			log.Println(err)
		}
		if err != nil &&
			(strings.Contains(err.Error(), "Cannot acquire directory lock on") ||
				strings.Contains(err.Error(), "Another process is using this Badger database.")) {
			return database
		}
		if err == nil {
			FDbClose(database)
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
