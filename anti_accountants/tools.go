package anti_accountants

import (
	"log"
	"math"
	"reflect"
	"sort"
)

func IS_IN[t any](element t, elements []t) bool {
	for _, a := range elements {
		if reflect.DeepEqual(a, element) {
			return true
		}
	}
	return false
}

func IS_SHORTER_THAN[t any](slice1, slice2 []t) bool {
	if len(slice1) < len(slice2) {
		return true
	}
	return false
}

func SMALLEST[t Number](a, b t) t {
	if a < b {
		return a
	}
	return b
}

func TRANSPOSE[t any](slice [][]t) [][]t {
	xl := len(slice[0])
	yl := len(slice)
	result := make([][]t, xl)
	for i := range result {
		result[i] = make([]t, yl)
	}
	for i := 0; i < xl; i++ {
		for j := 0; j < yl; j++ {
			result[i][j] = slice[j][i]
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

// func parse_date(string_date string, date_layouts []string) time.Time {
// 	for _, i := range date_layouts {
// 		date, err := time.Parse(i, string_date)
// 		if err == nil {
// 			return date
// 		}
// 	}
// 	error_date_layout(string_date)
// 	return time.Time{}
// }

func RETURN_SET_AND_DUPLICATES_SLICES[t any](accounts_numbers []t) ([]t, []t) {
	var set_of_elems, duplicated_element []t
big_loop:
	for _, element := range accounts_numbers {
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
	n := reflect.ValueOf(s).Len()
	swap := reflect.Swapper(s)
	for i, j := 0, n-1; i < j; i, j = i+1, j-1 {
		swap(i, j)
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

func TEST_FUNCTION[t any](actual, expected t) {
	if !reflect.DeepEqual(actual, expected) {
		log.Panic("actual : ", actual, " expected : ", expected)
	}
}

func SORT_BY_TIME_INVENTORY(slice []INVENTORY_TAG, is_ascending bool) {
	sort.Slice(slice, func(i, j int) bool {
		return slice[i].DATE_START.After(slice[j].DATE_START) == is_ascending
	})
}

func SORT_BY_TIME_JOURNAL(slice []JOURNAL_TAG, is_ascending bool) {
	sort.Slice(slice, func(i, j int) bool {
		return slice[i].DATE_START.After(slice[j].DATE_START) == is_ascending
	})
}

func POPUP[t any, type_index int | uint](slice []t, index_to_popup type_index) []t {
	return append(slice[:index_to_popup], slice[index_to_popup+1:]...)
}

func CUT_THE_SLICE[t any](slice []t, end int) []t {
	return slice[:len(slice)-end]
}
