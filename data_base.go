package anti_accountants

import (
	"encoding/json"
	"strconv"

	badger "github.com/dgraph-io/badger/v3"
)

func db_open_accounts() {
	var err error
	db_accounts, err = badger.Open(badger.DefaultOptions("./db_accounts"))
	error_fatal(1, err)
}

func db_open_journal() {
	var err error
	db_journal, err = badger.Open(badger.DefaultOptions("./db_journal"))
	error_fatal(1, err)
}

func db_open_inventory() {
	var err error
	db_inventory, err = badger.Open(badger.DefaultOptions("./db_inventory"))
	error_fatal(1, err)
}

func db_close_accounts() {
	db_accounts.Close()
}

func db_close_journal() {
	db_journal.Close()
}

func db_close_inventory() {
	db_inventory.Close()
}

// func check_accounts(column, table, panic string, elements []string) {
// 	results, err := DB.Query("select " + column + " from " + table)
// 	error_fatal(err)
// 	for results.Next() {
// 		var tag string
// 		results.Scan(&tag)
// 		if !is_in(tag, elements) {
// 			log.Panic(tag + panic)
// 		}
// 	}
// }

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

func JOURNAL_ORDERED_BY_DATE_ENTRY_NUMBER() []JOURNAL_TAG {
	db_open_journal()
	defer db_close_journal()
	var journal_tag_array []JOURNAL_TAG
	err := db_journal.View(func(txn *badger.Txn) error {
		// item, err := txn.Get([]byte("40982"))
		// error_fatal(5, err)
		// item.Value(func(val []byte) error {
		// 	err = json.Unmarshal(val, &journal_tag_array)
		// 	error_fatal(6, err)
		// 	return err
		// })
		// return err
		opts := badger.DefaultIteratorOptions
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			err := item.Value(func(val []byte) error {
				var journal_tag JOURNAL_TAG
				err := json.Unmarshal(val, &journal_tag)
				journal_tag_array = append(journal_tag_array, journal_tag)
				return err
			})
			error_fatal(7, err)
		}
		return nil
	})
	error_fatal(5, err)
	return journal_tag_array
}

// func account_balance(account string) float64 {
// 	var account_balance float64
// 	DB.QueryRow("select sum(value) from journal where account=?", account).Scan(&account_balance)
// 	return account_balance
// }

// func insert_into_inventory_func(array_of_journal_tag []JOURNAL_TAG) {
// 	for _, entry := range array_of_journal_tag {
// 		costs := cost_flow(entry.ACCOUNT_DEBIT, entry.QUANTITY_DEBIT, entry.BARCODE_DEBIT, true)
// 		if asc_or_desc(entry.ACCOUNT) != "" && costs == 0 {
// 			DB.Exec("insert into inventory(date,account,price,quantity,barcode,entry_expair,name,employee_name,entry_date)values (?,?,?,?,?,?,?,?,?)",
// 				&entry.DATE, &entry.ACCOUNT, &entry.PRICE, &entry.QUANTITY, &entry.BARCODE, &entry.ENTRY_EXPAIR, &entry.NAME, &entry.EMPLOYEE_NAME, &entry.ENTRY_DATE)
// 		}
// 	}
// }

func insert_into_journal(array_of_journal_tag []JOURNAL_TAG) {
	db_open_journal()
	defer db_close_journal()
	for _, entry := range array_of_journal_tag {
		err := db_journal.Update(func(txn *badger.Txn) error {
			json_entry, err := json.Marshal(entry)
			error_fatal(3, err)
			err = txn.Set([]byte(strconv.Itoa(entry.LINE_NUMBER)), []byte(json_entry))
			error_fatal(4, err)
			return nil
		})
		error_fatal(2, err)
	}
}

// func entry_number() int {
// 	var tag int
// 	err := DB.QueryRow("select max(entry_number) from journal").Scan(&tag)
// 	if err != nil {
// 		tag = 0
// 	}
// 	return tag + 1
// }

// func weighted_average(account string) {
// 	DB.Exec("update inventory set price=(select sum(price*quantity)/sum(quantity) from inventory where account=?) where account=?", account, account)
// }

// func weighted_average_for_barcode(account string) {
// 	// DB.Exec("update inventory set price=(select sum(price*quantity)/sum(quantity) from inventory where account=? and barcode=?) where account=? and barcode=?", account, barcode, account, barcode)
// }
