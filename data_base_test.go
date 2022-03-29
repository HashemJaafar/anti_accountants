package anti_accountants

import (
	"fmt"
	"testing"
	"time"
)

func Test_DB_OPEN(t *testing.T) {
	// db := DB_OPEN(DB_PATH_ACCOUNTS)
	// a := db.IsClosed()
	// TEST(true,a, false)
	// db.Close()
	// a = db.IsClosed()
	// TEST(true,a, true)
}
func Test_DB_CLOSE(t *testing.T) {
	TEST(true, DB_ACCOUNTS.IsClosed(), false)
	TEST(true, DB_JOURNAL.IsClosed(), false)
	TEST(true, DB_JOURNAL_DRAFT.IsClosed(), false)
	TEST(true, DB_INVENTORY.IsClosed(), false)
	DB_CLOSE()
	TEST(true, DB_ACCOUNTS.IsClosed(), true)
	TEST(true, DB_JOURNAL.IsClosed(), true)
	TEST(true, DB_JOURNAL_DRAFT.IsClosed(), true)
	TEST(true, DB_INVENTORY.IsClosed(), true)
}
func Test_DB_INSERT_INTO_ACCOUNTS(t *testing.T) {
}

func Test_DB_UPDATE(t *testing.T) {
	DB_UPDATE(DB_INVENTORY, []byte(time.Now().String()), INVENTORY_TAG{1, 10, "book1", time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC), time.Time{}})
}
func Test_DB_INSERT(t *testing.T) {
	// DB_INVENTORY.DropAll()
	DB_INSERT(DB_INVENTORY, []INVENTORY_TAG{
		{1, 10, "book1", time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC), time.Time{}},
		{2, 10, "book1", time.Date(2022, time.January, 2, 0, 0, 0, 0, time.UTC), time.Time{}},
		{3, 10, "book1", time.Date(2022, time.January, 3, 0, 0, 0, 0, time.UTC), time.Time{}},
		{4, 10, "book1", time.Date(2022, time.January, 4, 0, 0, 0, 0, time.UTC), time.Time{}},
	})
	DB_CLOSE()
}
func Test_DB_READ(t *testing.T) {
	_, inventory := DB_READ[INVENTORY_TAG](DB_INVENTORY)
	for _, a := range inventory {
		fmt.Println(a)
	}
}
func Test_DB_UPDATE_ACCOUNT_NAME_IN_JOURNAL(t *testing.T) {
}
func Test_ACCOUNT_BALANCE(t *testing.T) {
}
func Test_DB_LAST_LINE_IN_JOURNAL(t *testing.T) {
}
func Test_WEIGHTED_AVERAGE(t *testing.T) {
}
