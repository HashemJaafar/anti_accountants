package anti_accountants

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"unicode/utf8"
)

func Test_lower(t *testing.T) {
	fmt.Println(strings.ToLower("ALERT_FOR_MINIMUM_QUANTITY_BY_TURNOVER_IN_DAYS"))
	fmt.Println(strings.ToLower("ALERT_FOR_MINIMUM_QUANTITY_BY_QUINTITY"))
	fmt.Println(strings.ToLower("TARGET_BALANCE"))
	fmt.Println(strings.ToLower("IF_THE_TARGET_BALANCE_IS_LESS_IS_GOOD"))
}

func Test_a(t *testing.T) {
	add_under_score("remove the accounts not in accounts list")
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

func add_under_score(str string) {
	byt := []byte(str)
	for indexa, a := range byt {
		if a == ' ' {
			byt[indexa] = '_'
		}
	}
	fmt.Println(string(byt))
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
