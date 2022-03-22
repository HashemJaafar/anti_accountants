package anti_accountants

import (
	"encoding/json"
	"time"

	badger "github.com/dgraph-io/badger/v3"
)

func DB_OPEN(path string) *badger.DB {
	var db *badger.DB
	for {
		var err error
		db, err = badger.Open(badger.DefaultOptions(path))
		if err == nil {
			return db
		}
	}
}

func DB_INSERT_INTO_ACCOUNTS() {
	db := DB_OPEN(DB_ACCOUNTS)
	db.DropAll()
	defer db.Close()
	for _, entry := range ACCOUNTS {
		db.Update(func(txn *badger.Txn) error {
			json_entry, _ := json.Marshal(entry)
			txn.Set([]byte(entry.ACCOUNT_NAME), []byte(json_entry))
			return nil
		})
	}
}

func DB_INSERT_INTO_JOURNAL_OR_INVENTORY[t JOURNAL_TAG | INVENTORY_TAG](path string, array []t) {
	db := DB_OPEN(path)
	defer db.Close()
	for _, entry := range array {
		db.Update(func(txn *badger.Txn) error {
			json_entry, _ := json.Marshal(entry)
			txn.Set([]byte(time.Now().String()), []byte(json_entry))
			return nil
		})
	}
}

func DB_READ_ACCOUNTS() []ACCOUNT {
	db := DB_OPEN(DB_ACCOUNTS)
	defer db.Close()
	var array []ACCOUNT
	db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			item.Value(func(val []byte) error {
				var tag ACCOUNT
				json.Unmarshal(val, &tag)
				array = append(array, tag)
				return nil
			})
		}
		return nil
	})
	return array
}

func DB_READ_JOURNAL(path string) []JOURNAL_TAG {
	db := DB_OPEN(path)
	defer db.Close()
	var array []JOURNAL_TAG
	db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			item.Value(func(val []byte) error {
				var tag JOURNAL_TAG
				json.Unmarshal(val, &tag)
				array = append(array, tag)
				return nil
			})
		}
		return nil
	})
	return array
}

func DB_READ_INVENTORY(account_name string) []INVENTORY_TAG {
	db := DB_OPEN(DB_INVENTORY + account_name)
	defer db.Close()
	var array []INVENTORY_TAG
	db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			item.Value(func(val []byte) error {
				var tag INVENTORY_TAG
				json.Unmarshal(val, &tag)
				array = append(array, tag)
				return nil
			})
		}
		return nil
	})
	return array
}

func DB_UPDATE_ACCOUNT_NAME_IN_JOURNAL(account_name, new_name string) {
	db := DB_OPEN(DB_JOURNAL)
	defer db.Close()
	db.Update(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			item.Value(func(val []byte) error {
				var tag JOURNAL_TAG
				json.Unmarshal(val, &tag)
				if tag.ACCOUNT_CREDIT == account_name {
					tag.ACCOUNT_CREDIT = new_name
				}
				if tag.ACCOUNT_DEBIT == account_name {
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
	journal := DB_READ_JOURNAL(DB_JOURNAL)
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

func DB_LAST_LINE_IN_JOURNAL() JOURNAL_TAG {
	db := DB_OPEN(DB_JOURNAL)
	defer db.Close()
	var journal_tag JOURNAL_TAG
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
		var journal_tag JOURNAL_TAG
		json.Unmarshal(str, &journal_tag)
		return nil
	})
	return journal_tag
}

func WEIGHTED_AVERAGE(account string) {
	inventory := DB_READ_INVENTORY(account)
	var total_value, total_quantity float64
	for _, entry := range inventory {
		total_value += entry.PRICE * entry.QUANTITY
		total_quantity += entry.QUANTITY
	}
	price := total_value / total_quantity

	db := DB_OPEN(DB_INVENTORY + account)
	defer db.Close()
	db.View(func(txn *badger.Txn) error {
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
