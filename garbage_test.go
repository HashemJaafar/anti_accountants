//this file created automatically
package anti_accountants

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"
	"unicode/utf8"
)

func FuzzReverse(f *testing.F) {
	testcases := []string{"Hello, world", " ", "!12345"}
	for _, tc := range testcases {
		f.Add(tc) // Use f.Add to provide a seed corpus
	}
	f.Fuzz(func(t *testing.T, orig string) {
		rev, err1 := Reverse(orig)
		if err1 != nil {
			return
		}
		doubleRev, err2 := Reverse(rev)
		if err2 != nil {
			return
		}
		if orig != doubleRev {
			t.Errorf("Before: %q, after: %q", orig, doubleRev)
		}
		if utf8.ValidString(orig) && !utf8.ValidString(rev) {
			t.Errorf("Reverse produced invalid UTF-8 string %q", rev)
		}
	})
}

func Reverse(s string) (string, error) {
	if !utf8.ValidString(s) {
		return s, errors.New("input is not valid UTF-8")
	}
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r), nil
}

func Test2(t *testing.T) {
	fmt.Println(len(FORMAT_THE_STRING("	   ")))
}

func TestDropAll_database(t *testing.T) {
	DB_ACCOUNTS.DropAll()
	DB_JOURNAL.DropAll()
	DB_INVENTORY.DropAll()
	DB_CLOSE()
	PRINT_FORMATED_ACCOUNTS()
}

func Testa(t *testing.T) {
	add_under_score("remove the accounts not in accounts list")
}

func Testlower(t *testing.T) {
	fmt.Println(strings.ToLower("ALERT_FOR_MINIMUM_QUANTITY_BY_TURNOVER_IN_DAYS"))
	fmt.Println(strings.ToLower("ALERT_FOR_MINIMUM_QUANTITY_BY_QUINTITY"))
	fmt.Println(strings.ToLower("TARGET_BALANCE"))
	fmt.Println(strings.ToLower("IF_THE_TARGET_BALANCE_IS_LESS_IS_GOOD"))
}

func Testm(t *testing.T) {
	b, err := json.Marshal(10)
	fmt.Println(err, b)
}

func Testupper(t *testing.T) {
	fmt.Println(strings.ToLower("ENTRY_NUMBER_COMPOUND int, ENTRY_NUMBER_SIMPLE int"))
}

func add_under_score(str string) {
	byt := []byte(str)
	for indexa, a := range byt {
		if a == ' ' {
			byt[indexa] = '_'
		}
	}
	fmt.Println(string(byt))
}

func remove_under_score(str string) {
	byt := []byte(str)
	for indexa, a := range byt {
		if a == '_' {
			byt[indexa] = ' '
		}
	}
	fmt.Println(string(byt))
}
