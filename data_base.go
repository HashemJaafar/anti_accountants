package main

import (
	"encoding/json"

	badger "github.com/dgraph-io/badger/v3"
)

func ACCOUNT_BALANCE(account string) float64 {
	_, journal := DB_READ[JOURNAL_TAG](DB_JOURNAL)
	var value_credit float64
	for _, v1 := range journal {
		if account == v1.ACCOUNT_CREDIT {
			value_credit += v1.VALUE
		}
		if account == v1.ACCOUNT_DEBIT {
			value_credit -= v1.VALUE
		}
	}
	account_struct, _, _ := ACCOUNT_STRUCT_FROM_NAME(account)
	if account_struct.IS_CREDIT {
		return value_credit
	}
	return -value_credit
}

func DB_CLOSE() {
	DB_ACCOUNTS.Close()
	DB_JOURNAL.Close()
	DB_INVENTORY.Close()
}

func DB_INSERT[t any](db *badger.DB, slice []t) {
	for _, a := range slice {
		DB_UPDATE(db, NOW(), a)
	}
}

func DB_INSERT_INTO_ACCOUNTS() {
	DB_ACCOUNTS.DropAll()
	for _, a := range ACCOUNTS {
		DB_UPDATE(DB_ACCOUNTS, []byte(a.ACCOUNT_NAME), a)
	}
}

func DB_LAST_LINE[t any](db *badger.DB) t {
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

func DB_OPEN(path string) *badger.DB {
	for {
		db, err := badger.Open(badger.DefaultOptions(path))
		if err == nil {
			return db
		}
	}
}

func DB_READ[t any](db *badger.DB) ([][]byte, []t) {
	var values []t
	var keys [][]byte
	db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			item.Value(func(val []byte) error {
				var value t
				json.Unmarshal(val, &value)
				values = append(values, value)
				keys = append(keys, item.Key())
				return nil
			})
		}
		return nil
	})
	return keys, values
}

func DB_UPDATE[t any](db *badger.DB, key []byte, value t) {
	// db.Update(func(txn *badger.Txn) error {
	// 	json_value, _ := json.Marshal(value)
	// 	txn.Set(key, json_value)
	// 	return nil
	// })
	txn := db.NewTransaction(true)
	defer txn.Commit()
	json_value, _ := json.Marshal(value)
	txn.Set(key, json_value)
}

func DB_UPDATE_ACCOUNT_NAME_IN_INVENTORY(old_name, new_name string) {
	DB_INVENTORY.Update(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			item.Value(func(val []byte) error {
				var tag INVENTORY_TAG
				json.Unmarshal(val, &tag)
				if tag.ACCOUNT_NAME == old_name {
					tag.ACCOUNT_NAME = new_name
				}
				json_entry, _ := json.Marshal(tag)
				txn.Set([]byte(item.Key()), []byte(json_entry))
				return nil
			})
		}
		return nil
	})
}

func DB_UPDATE_ACCOUNT_NAME_IN_JOURNAL(old_name, new_name string) {
	DB_JOURNAL.Update(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			item.Value(func(val []byte) error {
				var tag JOURNAL_TAG
				json.Unmarshal(val, &tag)
				if tag.ACCOUNT_CREDIT == old_name {
					tag.ACCOUNT_CREDIT = new_name
				}
				if tag.ACCOUNT_DEBIT == old_name {
					tag.ACCOUNT_DEBIT = new_name
				}
				json_entry, _ := json.Marshal(tag)
				txn.Set([]byte(item.Key()), []byte(json_entry))
				return nil
			})
		}
		return nil
	})
}

func WEIGHTED_AVERAGE(account string) {
	keys, inventory := DB_READ[INVENTORY_TAG](DB_INVENTORY)
	var total_value, total_quantity float64
	for k1, v1 := range inventory {
		if v1.ACCOUNT_NAME == account {
			total_value += v1.PRICE * v1.QUANTITY
			total_quantity += v1.QUANTITY
			DB_INVENTORY.Update(func(txn *badger.Txn) error {
				return txn.Delete(keys[k1])
			})
		}
	}
	DB_UPDATE(DB_INVENTORY, NOW(), INVENTORY_TAG{total_value / total_quantity, total_quantity, account})
}
