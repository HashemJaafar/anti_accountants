package anti_accountants

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

	wait.Add(9)
	go open(&VDbJournal, CPathJournal)
	go open(&VDbInventory, CPathInventory)
	go open(&VDbAdjustingEntry, CPathAdjustingEntry)
	go open(&VDbEmployees, CPathEmployees)
	go open(&VDbJournalDrafts, CPathJournalDrafts)
	go open(&VDbInvoiceDrafts, CPathInvoiceDrafts)
	go func() {
		open(&VDbAccounts, CPathAccounts)
		_, VAccounts = FDbRead[SAccount1](VDbAccounts)
		wait.Done()
	}()
	go func() {
		open(&VDbAutoCompletionEntries, CPathAutoCompletionEntries)
		_, VAutoCompletionEntries = FDbRead[SAutoCompletion1](VDbAutoCompletionEntries)
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
	go close(VDbAdjustingEntry)
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

func FIsAnyClosed() bool {
	return VDbAccounts.IsClosed() ||
		VDbJournal.IsClosed() ||
		VDbInventory.IsClosed() ||
		VDbAutoCompletionEntries.IsClosed() ||
		VDbEmployees.IsClosed() ||
		VDbJournalDrafts.IsClosed() ||
		VDbInvoiceDrafts.IsClosed()
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
				err := json.Unmarshal(val, &Value)
				FPanicIfErr(err)
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

	keys, inventory := FDbRead[SAPQ1](VDbInventory)
	for k1, v1 := range inventory {
		if v1.AccountName == old {
			v1.AccountName = new
			FDbUpdate(VDbJournal, keys[k1], v1)
		}
	}

	keys, VAutoCompletionEntries = FDbRead[SAutoCompletion1](VDbAutoCompletionEntries)
	for k1 := range VAutoCompletionEntries {
		if VAutoCompletionEntries[k1].Inventory == old {
			VAutoCompletionEntries[k1].Inventory = new
		}
		if VAutoCompletionEntries[k1].CostOfGoodsSold == old {
			VAutoCompletionEntries[k1].CostOfGoodsSold = new
		}
		if VAutoCompletionEntries[k1].TaxExpenses == old {
			VAutoCompletionEntries[k1].TaxExpenses = new
		}
		if VAutoCompletionEntries[k1].TaxLiability == old {
			VAutoCompletionEntries[k1].TaxLiability = new
		}
		if VAutoCompletionEntries[k1].Revenue == old {
			VAutoCompletionEntries[k1].Revenue = new
		}
		if VAutoCompletionEntries[k1].Discount == old {
			VAutoCompletionEntries[k1].Discount = new
		}
		FDbUpdate(VDbJournal, keys[k1], VAutoCompletionEntries[k1])
	}

	keys, adjustingEntry := FDbRead[SAdjustingEntry](VDbAdjustingEntry)
	for k1, v1 := range adjustingEntry {
		if v1.AccountName1 == old {
			v1.AccountName1 = new
			FDbUpdate(VDbJournal, keys[k1], v1)
		}
		if v1.AccountName2 == old {
			v1.AccountName2 = new
			FDbUpdate(VDbJournal, keys[k1], v1)
		}
	}
}

func FWeightedAverage(account string) {
	_, price, quantity := FTotalValuePriceQuantity(account)

	keys, inventory := FDbRead[SAPQ1](VDbInventory)
	for k1, v1 := range inventory {
		if v1.AccountName == account {
			FDbDelete(VDbInventory, keys[k1])
		}
	}

	FDbUpdate(VDbInventory, FNow(), SAPQ1{account, price, quantity})
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
		if index, isIn := FFind(old, v1.Labels); isIn {
			v1.Labels[index] = new
			FDbUpdate(VDbJournal, keys[k1], v1)
		}
	}
}

func FChangeEntryInfoByEntryNumberCompund(entryNumberCompund uint, new SEntry1) {
	keys, journal := FDbRead[SJournal1](VDbJournal)
	for k1, v1 := range journal {
		if v1.EntryNumberCompound == entryNumberCompund {
			v1.Notes = new.Notes
			v1.Name = new.Name
			v1.Labels = new.Labels
			FDbUpdate(VDbJournal, keys[k1], v1)
		}
	}
}
