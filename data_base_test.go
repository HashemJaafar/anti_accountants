package anti_accountants

import (
	"testing"
)

func Test_open_and_create_database(t *testing.T) {
	open_and_create_database("mysql", "dataSourceName", "acc")
}
