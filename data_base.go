package anti_accountants

import (
	"encoding/json"
	"time"

	badger "github.com/dgraph-io/badger/v3"
)

func db_open(path string) *badger.DB {
	var db *badger.DB
	for {
		var err error
		db, err = badger.Open(badger.DefaultOptions("./" + path))
		if err == nil {
			return db
		}
	}
}

func db_insert_into_accounts() {
	db := db_open(db_accounts)
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

func db_insert_into_journal(array []JOURNAL_TAG) {
	db := db_open(db_journal)
	defer db.Close()
	for _, entry := range array {
		db.Update(func(txn *badger.Txn) error {
			json_entry, _ := json.Marshal(entry)
			txn.Set([]byte(time.Now().String()), []byte(json_entry))
			return nil
		})
	}
}

func db_insert_into_inventory(array []INVENTORY_TAG) {
	db := db_open(db_inventory)
	defer db.Close()
	for _, entry := range array {
		db.Update(func(txn *badger.Txn) error {
			json_entry, _ := json.Marshal(entry)
			txn.Set([]byte(time.Now().String()), []byte(json_entry))
			return nil
		})
	}
}

func db_read_accounts() []ACCOUNT {
	db := db_open(db_accounts)
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

func db_read_journal() []JOURNAL_TAG {
	db := db_open(db_journal)
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

func db_read_inventory() []INVENTORY_TAG {
	db := db_open(db_inventory)
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

// func CHANGE_ACCOUNT_NAME(name, new_name string) {
// 	var tag string
// 	err := DB.QueryRow("select account from journal where account=? limit 1", new_name).Scan(&tag)
// 	if err == nil {
// 		error_you_cant_change_the_name(name, new_name)
// 	} else {
// 		DB.Exec("update journal set account=? where account=?", new_name, name)
// 		DB.Exec("update inventory set account=? where account=?", new_name, name)
// 	}
// }

// func account_balance(account string) float64 {
// 	var account_balance float64
// 	DB.QueryRow("select sum(value) from journal where account=?", account).Scan(&account_balance)
// 	return account_balance
// }

func last_line_in_db() JOURNAL_TAG {
	db := db_open(db_journal)
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

func entry_number() (uint, uint, uint) {
	journal := last_line_in_db()
	return journal.ENTRY_NUMBER + 1, journal.ENTRY_NUMBER_COMPOUND + 1, journal.ENTRY_NUMBER_SIMPLE + 1
}

func weighted_average(account string) {
	inventory := db_read_inventory()
	var total_value, total_quantity float64
	for _, entry := range inventory {
		if entry.ACCOUNT_NAME == account {
			total_value += entry.PRICE * entry.QUANTITY
			total_quantity += entry.QUANTITY
		}
	}
	price := total_value / total_quantity

	db := db_open(db_inventory)
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
				if tag.ACCOUNT_NAME == account {
					tag.PRICE = price
					json_entry, _ := json.Marshal(tag)
					txn.Set([]byte(item.Key()), []byte(json_entry))
				}
				return nil
			})
		}
		return nil
	})
}
