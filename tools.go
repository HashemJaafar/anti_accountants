package anti_accountants

import (
	"fmt"
	"log"
	"math"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

func IS_IN[t any](element t, elements []t) bool {
	for _, a := range elements {
		if reflect.DeepEqual(a, element) {
			return true
		}
	}
	return false
}

func SMALLEST[t NUMBER](a, b t) t {
	if a < b {
		return a
	}
	return b
}

func TRANSPOSE[t any](slice [][]t) [][]t {
	xl := len(slice[0])
	yl := len(slice)
	result := make([][]t, xl)
	for a := range result {
		result[a] = make([]t, yl)
	}
	for a := 0; a < xl; a++ {
		for b := 0; b < yl; b++ {
			result[a][b] = slice[b][a]
		}
	}
	return result
}

func UNPACK[t any](slice [][]t) []t {
	var result []t
	for _, element := range slice {
		result = append(result, element...)
	}
	return result
}

func PACK[t any](len_new_slice int, slice []t) []t {
	new_slice := make([]t, len_new_slice)
	for indexa, a := range slice {
		new_slice[indexa] = a
	}
	return new_slice
}

func RETURN_SAME_SIGN_OF_NUMBER_SIGN(number_sign, number float64) float64 {
	if number_sign < 0 {
		return -math.Abs(number)
	}
	return math.Abs(number)
}

func RETURN_SET_AND_DUPLICATES_SLICES[t any](slice []t) ([]t, []t) {
	var set_of_elems, duplicated_element []t
big_loop:
	for _, element := range slice {
		for _, b := range set_of_elems {
			if reflect.DeepEqual(b, element) {
				duplicated_element = append(duplicated_element, element)
				continue big_loop
			}
		}
		set_of_elems = append(set_of_elems, element)
	}
	return set_of_elems, duplicated_element
}

func REVERSE_SLICE[t any](s []t) {
	for a, b := 0, len(s)-1; a < b; a, b = a+1, b-1 {
		SWAP(s, a, b)
	}
}

func INITIALIZE_MAP_4(m map[string]map[string]map[string]map[string]map[string]float64, a, b, c, d string) map[string]float64 {
	if m[a] == nil {
		m[a] = map[string]map[string]map[string]map[string]float64{}
	}
	if m[a][b] == nil {
		m[a][b] = map[string]map[string]map[string]float64{}
	}
	if m[a][b][c] == nil {
		m[a][b][c] = map[string]map[string]float64{}
	}
	if m[a][b][c][d] == nil {
		m[a][b][c][d] = map[string]float64{}
	}
	return m[a][b][c][d]
}

func INITIALIZE_MAP_3(m map[string]map[string]map[string]map[string]float64, a, b, c string) map[string]float64 {
	if m[a] == nil {
		m[a] = map[string]map[string]map[string]float64{}
	}
	if m[a][b] == nil {
		m[a][b] = map[string]map[string]float64{}
	}
	if m[a][b][c] == nil {
		m[a][b][c] = map[string]float64{}
	}
	return m[a][b][c]
}

func TEST[t any](should_equal bool, actual, expected t) {
	if !reflect.DeepEqual(actual, expected) == should_equal {
		fmt.Fprintln(PRINT_TABLE, "\033[35m", "should_equal\t:", should_equal) //purple
		fmt.Fprintln(PRINT_TABLE, "\033[34m", "actual\t:", actual)             //blue
		fmt.Fprintln(PRINT_TABLE, "\033[33m", "expected\t:", expected)         //yellow
		PRINT_TABLE.Flush()
		fmt.Println("\033[31m") //red
		log.Panic()
	}
}

func SORT_BY_TIME_INVENTORY(slice1 []INVENTORY_TAG, slice2 [][]byte, is_ascending bool) {
	for indexa := range slice1 {
		for indexb := range slice1 {
			if indexa < indexb && slice1[indexa].DATE_START.After(slice1[indexb].DATE_START) == is_ascending {
				SWAP(slice1, indexa, indexb)
				SWAP(slice2, indexa, indexb)
			}
		}
	}
}

func SORT_BY_TIME_JOURNAL(slice []JOURNAL_TAG, is_ascending bool) {
	sort.Slice(slice, func(i, j int) bool {
		return slice[i].DATE_START.Before(slice[j].DATE_START) == is_ascending
	})
}

func CUT_THE_SLICE[t any](s []t, a int) []t { return s[:len(s)-a] }
func POPUP[t any](s []t, a int) []t         { return append(s[:a], s[a+1:]...) }
func SWAP[t any](s []t, a, b int)           { s[a], s[b] = s[b], s[a] }

func FORMAT_THE_STRING(str string) string {
	return strings.ToLower(strings.Join(strings.Fields(str), " "))
}

func ASSIGN_NUMBER_IF_NUMBER(m map[string]float64, str string) {
	number, err := strconv.ParseFloat(str, 64)
	if err == nil {
		m[str] = number
	}
}

func CONVERT_NAN_TO_ZERO(value float64) float64 {
	if math.IsNaN(value) {
		return 0
	}
	return value
}

func IS_INF_IN(numbers ...float64) bool {
	for _, a := range numbers {
		if math.IsInf(a, 0) {
			return true
		}
	}
	return false
}
