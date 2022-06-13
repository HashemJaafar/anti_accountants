package main

import (
	"encoding/json"
	"log"
	"sync"

	badger "github.com/dgraph-io/badger/v3"
)

func FDbOpenAll() {
	var wait sync.WaitGroup

	open := func(db **badger.DB, path string) {
		*db = FDbOpen(*db, CPathDataBase+VCompanyName+path)
		wait.Done()
	}

	wait.Add(7)
	go open(&VDbAccounts, CPathAccounts)
	go open(&VDbJournal, CPathJournal)
	go open(&VDbInventory, CPathInventory)
	go open(&VDbAutoCompletionEntries, CPathAutoCompletionEntries)
	go open(&VDbEmployees, CPathEmployees)
	go open(&VDbJournalDrafts, CPathJournalDrafts)
	go open(&VDbInvoiceDrafts, CPathInvoiceDrafts)
	wait.Wait()

	wait.Add(2)
	go func() {
		_, VAccounts = FDbRead[SAccount1](VDbAccounts)
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

	close := func(db *badger.DB) {
		FDbClose(db)
		wait.Done()
	}

	wait.Add(7)
	go close(VDbAccounts)
	go close(VDbJournal)
	go close(VDbInventory)
	go close(VDbAutoCompletionEntries)
	go close(VDbEmployees)
	go close(VDbJournalDrafts)
	go close(VDbInvoiceDrafts)
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
		FDbUpdate(VDbAccounts, []byte(v1.Name), v1)
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

func FDbOpen(oldDB *badger.DB, path string) *badger.DB {
	FDbClose(oldDB)
	for {
		newDB, err := badger.Open(badger.DefaultOptions(path))
		if err == nil {
			return newDB
		} else {
			log.Println(err)
		}
		// if err != nil &&
		// 	(strings.Contains(err.Error(), "Cannot acquire directory lock on") ||
		// 		strings.Contains(err.Error(), "Another process is using this Badger database.")) {
		// 	return oldDB
		// }
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
	keys, journal := FDbRead[SJournal1](VDbJournal)
	for k1, v1 := range journal {
		if v1.CreditAccountName == old {
			v1.CreditAccountName = new
			FDbUpdate(VDbJournal, keys[k1], v1)
		}
		if v1.DebitAccountName == old {
			v1.DebitAccountName = new
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
	_, price, quantity := FTotalValuePriceQuantity(account)

	keys, inventory := FDbRead[SAPQ](VDbInventory)
	for k1, v1 := range inventory {
		if v1.TAccountName == account {
			FDbDelete(VDbInventory, keys[k1])
		}
	}

	FDbUpdate(VDbInventory, FNow(), SAPQ{account, price, FAbs(quantity)})
}

func FChangeNotes(old, new string) {
	keys, journal := FDbRead[SJournal1](VDbJournal)
	for k1, v1 := range journal {
		if v1.Notes == old {
			v1.Notes = new
			FDbUpdate(VDbJournal, keys[k1], v1)
		}
	}
}

func FChangeName(old, new string) {
	keys, journal := FDbRead[SJournal1](VDbJournal)
	for k1, v1 := range journal {
		if v1.Name == old {
			v1.Name = new
			FDbUpdate(VDbJournal, keys[k1], v1)
		}
	}
}

func FChangeEmployeeName(old, new string) {
	keys, journal := FDbRead[SJournal1](VDbJournal)
	for k1, v1 := range journal {
		if v1.Employee == old {
			v1.Employee = new
			FDbUpdate(VDbJournal, keys[k1], v1)
		}
	}
}

func FChangeTypeOfCompoundEntry(old, new string) {
	keys, journal := FDbRead[SJournal1](VDbJournal)
	for k1, v1 := range journal {
		if v1.TypeOfCompoundEntry == old {
			v1.TypeOfCompoundEntry = new
			FDbUpdate(VDbJournal, keys[k1], v1)
		}
	}
}

func FChangeEntryInfoByEntryNumberCompund(entryNumberCompund uint, new SEntry) {
	keys, journal := FDbRead[SJournal1](VDbJournal)
	for k1, v1 := range journal {
		if v1.EntryNumberCompound == entryNumberCompund {
			v1.Notes = new.Notes
			v1.Name = new.Name
			v1.TypeOfCompoundEntry = new.TypeOfCompoundEntry
			FDbUpdate(VDbJournal, keys[k1], v1)
		}
	}
}
