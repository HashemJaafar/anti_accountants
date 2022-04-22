package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"reflect"
	"runtime/debug"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"
)

func NOW() []byte {
	// i use this function to get the current time in the format of TIME_LAYOUT to make the error less likely
	return []byte(time.Now().Format(TIME_LAYOUT))
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

func CUT_THE_SLICE[t any](s []t, a int) []t { return s[:len(s)-a] }

func FORMAT_THE_STRING(str string) string {
	return strings.ToLower(strings.Join(strings.Fields(str), " "))
}

func INITIALIZE_MAP_6[t1, t2, t3, t4, t5, t6 comparable, tr any](m map[t1]map[t2]map[t3]map[t4]map[t5]map[t6]tr, i1 t1, i2 t2, i3 t3, i4 t4, i5 t5) map[t6]tr {
	if m == nil {
		m = map[t1]map[t2]map[t3]map[t4]map[t5]map[t6]tr{}
	}
	if m[i1] == nil {
		m[i1] = map[t2]map[t3]map[t4]map[t5]map[t6]tr{}
	}
	if m[i1][i2] == nil {
		m[i1][i2] = map[t3]map[t4]map[t5]map[t6]tr{}
	}
	if m[i1][i2][i3] == nil {
		m[i1][i2][i3] = map[t4]map[t5]map[t6]tr{}
	}
	if m[i1][i2][i3][i4] == nil {
		m[i1][i2][i3][i4] = map[t5]map[t6]tr{}
	}
	if m[i1][i2][i3][i4][i5] == nil {
		m[i1][i2][i3][i4][i5] = map[t6]tr{}
	}
	return m[i1][i2][i3][i4][i5]
}

func INITIALIZE_MAP_5[t1, t2, t3, t4, t5 comparable, tr any](m map[t1]map[t2]map[t3]map[t4]map[t5]tr, i1 t1, i2 t2, i3 t3, i4 t4) map[t5]tr {
	if m == nil {
		m = map[t1]map[t2]map[t3]map[t4]map[t5]tr{}
	}
	if m[i1] == nil {
		m[i1] = map[t2]map[t3]map[t4]map[t5]tr{}
	}
	if m[i1][i2] == nil {
		m[i1][i2] = map[t3]map[t4]map[t5]tr{}
	}
	if m[i1][i2][i3] == nil {
		m[i1][i2][i3] = map[t4]map[t5]tr{}
	}
	if m[i1][i2][i3][i4] == nil {
		m[i1][i2][i3][i4] = map[t5]tr{}
	}
	return m[i1][i2][i3][i4]
}

func IS_IN[t any](element t, elements []t) bool {
	for _, a := range elements {
		if reflect.DeepEqual(a, element) {
			return true
		}
	}
	return false
}

func IS_INF_IN(numbers ...float64) bool {
	for _, a := range numbers {
		if math.IsInf(a, 0) {
			return true
		}
	}
	return false
}

func PACK[t any](len_new_slice int, slice []t) []t {
	new_slice := make([]t, len_new_slice)
	for indexa, a := range slice {
		new_slice[indexa] = a
	}
	return new_slice
}

func REMOVE[t any](s []t, a int) []t { return append(s[:a], s[a+1:]...) }

func RETURN_SAME_SIGN_OF_NUMBER_SIGN(number_sign, number float64) float64 {
	if number_sign < 0 {
		return -ABS(number)
	}
	return ABS(number)
}

func ABS[t NUMBER](n t) t {
	// this is alternative of math.Abs
	if n < 0 {
		return -n
	}
	return n
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

func SMALLEST[t NUMBER](a, b t) t {
	if a < b {
		return a
	}
	return b
}

func SORT_TIME(slice1 []time.Time, is_ascending bool) {
	for k1 := range slice1 {
		for k2 := range slice1 {
			if k1 < k2 && (slice1[k1]).After((slice1[k2])) == is_ascending {
				SWAP(slice1, k1, k2)
			}
		}
	}
}
func CONVERT_BYTE_SLICE_TO_TIME(slice [][]byte) []time.Time {
	var slice_of_time []time.Time
	for _, v1 := range slice {
		date, _ := time.Parse(TIME_LAYOUT, string(v1))
		slice_of_time = append(slice_of_time, date)
	}
	return slice_of_time
}

func SWAP[t any](s []t, a, b int) { s[a], s[b] = s[b], s[a] }

func TEST[t any](should_equal bool, actual, expected t) {
	if reflect.DeepEqual(actual, expected) != should_equal {
		fail_test_number++
		p := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
		fmt.Fprintln(p, "\033[32m", "fail_test_number\t:", fail_test_number) //green
		fmt.Fprintln(p, "\033[35m", "should_equal\t:", should_equal)         //purple
		fmt.Fprintln(p, "\033[34m", "actual\t:", actual)                     //blue
		fmt.Fprintln(p, "\033[33m", "expected\t:", expected)                 //yellow
		p.Flush()

		// fmt.Println("\033[34m") //blue
		// spew.Dump(actual)
		// fmt.Println("\033[33m") //yellow
		// spew.Dump(expected)

		fmt.Println("\033[31m") //red
		defer func() {
			if r := recover(); r != nil {
				fmt.Println(string(debug.Stack()), "\033[0m") //reset
			}
		}()
		log.Panic()
	}
}
func PRINT_SLICE[t any](a1 []t) {
	for _, v1 := range a1 {
		fmt.Println(v1)
	}
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

func FILTER_NUMBER[t NUMBER](number, bellow, above t, is_just_between bool) bool {
	is_bellow := number <= bellow
	is_above := number >= above
	if is_just_between {
		return is_bellow && is_above
	}
	return is_bellow || is_above
}

func print_map_6[t1, t2, t3, t4, t5, t6 comparable, tr any](m map[t1]map[t2]map[t3]map[t4]map[t5]map[t6]tr) {
	for k1, v1 := range m {
		for k2, v2 := range v1 {
			for k3, v3 := range v2 {
				for k4, v4 := range v3 {
					for k5, v5 := range v4 {
						for k6, v6 := range v5 {
							fmt.Fprintln(PRINT_TABLE, k1, "\t", k2, "\t", k3, "\t", k4, "\t", k5, "\t", k6, "\t", v6)
						}
					}
				}
			}
		}
	}
	fmt.Println("//////////////////////////////////////////")
	PRINT_TABLE.Flush()
}

func print_map_5[t1, t2, t3, t4, t5 comparable, tr any](m map[t1]map[t2]map[t3]map[t4]map[t5]tr) {
	for k1, v1 := range m {
		for k2, v2 := range v1 {
			for k3, v3 := range v2 {
				for k4, v4 := range v3 {
					for k5, v5 := range v4 {
						fmt.Fprintln(PRINT_TABLE, k1, "\t", k2, "\t", k3, "\t", k4, "\t", k5, "\t", v5)
					}
				}
			}
		}
	}
	fmt.Println("//////////////////////////////////////////")
	PRINT_TABLE.Flush()
}

func print_map_2[t1, t2 comparable, tr any](m map[t1]map[t2]tr) {
	for k1, v1 := range m {
		for k2, v2 := range v1 {
			fmt.Fprintln(PRINT_TABLE, k1, "\t", k2, "\t", v2)
		}
	}
	fmt.Println("//////////////////////////////////////////")
	PRINT_TABLE.Flush()
}
