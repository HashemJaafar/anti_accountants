package anti_accountants

import (
	"encoding/json"
	"time"

	badger "github.com/dgraph-io/badger/v3"
)

func DB_OPEN(path string) *badger.DB {
	for {
		db, err := badger.Open(badger.DefaultOptions(path))
		if err == nil {
			return db
		}
	}
}

func DB_CLOSE() {
	DB_ACCOUNTS.Close()
	DB_JOURNAL.Close()
	DB_JOURNAL_DRAFT.Close()
	DB_INVENTORY.Close()
}

func DB_UPDATE[t any](db *badger.DB, key []byte, value t) {
	db.Update(func(txn *badger.Txn) error {
		json_value, _ := json.Marshal(value)
		txn.Set(key, json_value)
		return nil
	})
}

func DB_INSERT_INTO_ACCOUNTS() {
	DB_ACCOUNTS.DropAll()
	for _, a := range ACCOUNTS {
		DB_UPDATE(DB_ACCOUNTS, []byte(a.ACCOUNT_NAME), a)
	}
}

func DB_INSERT[t any](db *badger.DB, slice []t) {
	for _, a := range slice {
		DB_UPDATE(db, []byte(time.Now().String()), a)
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

func ACCOUNT_BALANCE(account string) float64 {
	_, journal := DB_READ[JOURNAL_TAG](DB_JOURNAL)
	var value_debit, value_credit float64
	for _, entry := range journal {
		if account == entry.ACCOUNT_CREDIT {
			value_credit += entry.VALUE
		}
		if account == entry.ACCOUNT_DEBIT {
			value_debit += entry.VALUE
		}
	}
	account_struct, _, _ := ACCOUNT_STRUCT_FROM_NAME(account)
	if account_struct.IS_CREDIT {
		return value_credit - value_debit
	}
	return value_debit - value_credit
}

func DB_LAST_LINE[t any](db *badger.DB) t {
	var tag t
	db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := txn.NewIterator(opts)
		defer it.Close()
		var str []byte
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			item.Value(func(val []byte) error {
				str = val
				return nil
			})
		}
		var tag t
		json.Unmarshal(str, &tag)
		return nil
	})
	return tag
}

func WEIGHTED_AVERAGE(account string) {
	_, inventory := DB_READ[INVENTORY_TAG](DB_INVENTORY)
	var total_value, total_quantity float64
	for _, entry := range inventory {
		total_value += entry.PRICE * entry.QUANTITY
		total_quantity += entry.QUANTITY
	}
	price := total_value / total_quantity

	DB_INVENTORY.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			item.Value(func(val []byte) error {
				var tag INVENTORY_TAG
				json.Unmarshal(val, &tag)
				tag.PRICE = price
				json_entry, _ := json.Marshal(tag)
				txn.Set([]byte(item.Key()), []byte(json_entry))
				return nil
			})
		}
		return nil
	})
}
