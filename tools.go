package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"reflect"
	"runtime/debug"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"
)

func CutTheSlice[t any](s []t, a int) []t { return s[:len(s)-a] }
func Remove[t any](s []t, a int) []t      { return append(s[:a], s[a+1:]...) }
func Swap[t any](s []t, a, b int)         { s[a], s[b] = s[b], s[a] }

func Now() []byte {
	// i use this function to get the current time in the format of TimeLayout to make the error less likely
	return []byte(time.Now().Format(TimeLayout))
}

func AssignNumberIfNumber(m map[string]float64, str string) {
	number, err := strconv.ParseFloat(str, 64)
	if err == nil {
		m[str] = number
	}
}

func ConvertNanToZero(VALUE float64) float64 {
	if math.IsNaN(VALUE) {
		return 0
	}
	return VALUE
}

func FormatTheString(str string) string {
	return strings.ToLower(strings.Join(strings.Fields(str), " "))
}

func InitializeMap6[t1, t2, t3, t4, t5, t6 comparable, tr any](m map[t1]map[t2]map[t3]map[t4]map[t5]map[t6]tr, i1 t1, i2 t2, i3 t3, i4 t4, i5 t5) map[t6]tr {
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

func InitializeMap5[t1, t2, t3, t4, t5 comparable, tr any](m map[t1]map[t2]map[t3]map[t4]map[t5]tr, i1 t1, i2 t2, i3 t3, i4 t4) map[t5]tr {
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

func IsIn[t comparable](element t, elements []t) bool {
	for _, v1 := range elements {
		if v1 == element {
			return true
		}
	}
	return false
}

func IsInfIn(numbers ...float64) bool {
	for _, a := range numbers {
		if math.IsInf(a, 0) {
			return true
		}
	}
	return false
}

func Pack[t any](lenNewSlice int, slice []t) []t {
	newSlice := make([]t, lenNewSlice)
	for indexa, a := range slice {
		newSlice[indexa] = a
	}
	return newSlice
}

func ReturnSameSignOfNumberSign(numberSign, number float64) float64 {
	if numberSign < 0 {
		return -Abs(number)
	}
	return Abs(number)
}

func Abs[t Number](n t) t {
	// this is alternative of math.Abs
	if n < 0 {
		return -n
	}
	return n
}

func ReturnSetAndDuplicatesSlices[t any](slice []t) ([]t, []t) {
	var setOfElems, duplicatedElement []t
big_loop:
	for _, element := range slice {
		for _, b := range setOfElems {
			if reflect.DeepEqual(b, element) {
				duplicatedElement = append(duplicatedElement, element)
				continue big_loop
			}
		}
		setOfElems = append(setOfElems, element)
	}
	return setOfElems, duplicatedElement
}

func ReverseSlice[t any](s []t) {
	for a, b := 0, len(s)-1; a < b; a, b = a+1, b-1 {
		Swap(s, a, b)
	}
}

func Smallest[t Number](a, b t) t {
	if a < b {
		return a
	}
	return b
}

func SortTime(slice1 []time.Time, isAscending bool) {
	for k1 := range slice1 {
		for k2 := range slice1 {
			if k1 < k2 && (slice1[k1]).After((slice1[k2])) == isAscending {
				Swap(slice1, k1, k2)
			}
		}
	}
}

func SortStatementNumber(slice1 []FilteredStatement, isAscending bool) {
	sort.Slice(slice1, func(k1, k2 int) bool {
		return slice1[k1].Number > slice1[k2].Number == isAscending
	})
}

func ConvertByteSliceToTime(slice [][]byte) []time.Time {
	var sliceOfTime []time.Time
	for _, v1 := range slice {
		date, _ := time.Parse(TimeLayout, string(v1))
		sliceOfTime = append(sliceOfTime, date)
	}
	return sliceOfTime
}

func Test[t any](shouldEqual bool, actual, expected t) {
	if reflect.DeepEqual(actual, expected) != shouldEqual {
		FailTestNumber++
		p := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
		fmt.Fprintln(p, "\033[32m", "fail_test_number\t:", FailTestNumber) //green
		fmt.Fprintln(p, "\033[35m", "should_equal\t:", shouldEqual)        //purple
		fmt.Fprintln(p, "\033[34m", "actual\t:", actual)                   //blue
		fmt.Fprintln(p, "\033[33m", "expected\t:", expected)               //yellow
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

func PrintSlice[t any](a1 []t) {
	for _, v1 := range a1 {
		fmt.Println(v1)
	}
}

func Transpose[t any](slice [][]t) [][]t {
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

func Unpack[t any](slice [][]t) []t {
	var result []t
	for _, element := range slice {
		result = append(result, element...)
	}
	return result
}

func (s FilterDate) Filter(input time.Time) bool {
	if !s.IsFilter {
		return true
	}

	if s.Big.Before(s.Small) {
		s.Big, s.Small = s.Small, s.Big
	}

	isSmallerThanSmall := input.Before(s.Small)
	isBiggerThanBig := input.After(s.Big)

	switch s.Way {
	case Between:
		return !isSmallerThanSmall && !isBiggerThanBig
	case NotBetween:
		return isSmallerThanSmall || isBiggerThanBig
	case Bigger:
		return input.After(s.Big)
	case Smaller:
		return input.Before(s.Small)
	case EqualToOneOfThem:
		return input == s.Small || input == s.Big
	}

	return false
}

func (s FilterNumber) Filter(input float64) bool {
	if !s.IsFilter {
		return true
	}

	if s.Big < s.Small {
		s.Big, s.Small = s.Small, s.Big
	}

	isSmallerThanSmall := input < s.Small
	isBiggerThanBig := input > s.Big

	switch s.Way {
	case Between:
		return !isSmallerThanSmall && !isBiggerThanBig
	case NotBetween:
		return isSmallerThanSmall || isBiggerThanBig
	case Bigger:
		return input > s.Big
	case Smaller:
		return input < s.Small
	case EqualToOneOfThem:
		return input == s.Small || input == s.Big
	}

	return false
}

func (s FilterString) Filter(input string) bool {
	if !s.IsFilter {
		return true
	}

	switch s.Way {
	case InSlice:
		return IsIn(input, s.Slice)
	case NotInSlice:
		return IsIn(input, s.Slice) == false
	case ElementsInElement:
		return FElementsInElement(input, s.Slice)
	case ElementsNotInElement:
		return FElementsInElement(input, s.Slice) == false
	}
	return false
}

func (s FilterSliceString) Filter(input []string) bool {
	if !s.IsFilter {
		return true
	}
	for _, v1 := range input {
		for _, v2 := range s.Slice {
			if v1 == v2 {
				return true == s.InSlice
			}
		}
	}
	return false == s.InSlice
}

func (s FilterSliceUint) Filter(input uint) bool {
	if !s.IsFilter {
		return true
	}
	return IsIn(input, s.Slice) == s.InSlice
}

func (s FilterBool) Filter(input bool) bool {
	if !s.IsFilter {
		return true
	}
	return input == s.BoolValue
}

func FElementsInElement(element string, elements []string) bool {
	for _, v1 := range elements {
		if strings.Contains(element, v1) {
			return true
		}
	}
	return false
}

func FunctionFilterDuplicate[t comparable](input1, input2 t, f bool) bool {
	if !f {
		return true
	}
	return input1 == input2
}

func PrintMap6[t1, t2, t3, t4, t5, t6 comparable, tr any](m map[t1]map[t2]map[t3]map[t4]map[t5]map[t6]tr) {
	for k1, v1 := range m {
		for k2, v2 := range v1 {
			for k3, v3 := range v2 {
				for k4, v4 := range v3 {
					for k5, v5 := range v4 {
						for k6, v6 := range v5 {
							fmt.Fprintln(PrintTable, k1, "\t", k2, "\t", k3, "\t", k4, "\t", k5, "\t", k6, "\t", v6)
						}
					}
				}
			}
		}
	}
	fmt.Println("//////////////////////////////////////////")
	PrintTable.Flush()
}

func PrintMap5[t1, t2, t3, t4, t5 comparable, tr any](m map[t1]map[t2]map[t3]map[t4]map[t5]tr) {
	for k1, v1 := range m {
		for k2, v2 := range v1 {
			for k3, v3 := range v2 {
				for k4, v4 := range v3 {
					for k5, v5 := range v4 {
						fmt.Fprintln(PrintTable, k1, "\t", k2, "\t", k3, "\t", k4, "\t", k5, "\t", v5)
					}
				}
			}
		}
	}
	fmt.Println("//////////////////////////////////////////")
	PrintTable.Flush()
}

func PrintMap2[t1, t2 comparable, tr any](m map[t1]map[t2]tr) {
	for k1, v1 := range m {
		for k2, v2 := range v1 {
			fmt.Fprintln(PrintTable, k1, "\t", k2, "\t", v2)
		}
	}
	fmt.Println("//////////////////////////////////////////")
	PrintTable.Flush()
}

func PrintFilteredStatement(slice []FilteredStatement) {
	fmt.Fprintln(PrintTable, "Account1", "\t", "Account2", "\t", "Name", "\t", "Vpq", "\t", "TypeOfVpq", "\t", "Number")
	for _, v1 := range slice {
		fmt.Fprintln(PrintTable, v1.Account1, "\t", v1.Account2, "\t", v1.Name, "\t", v1.Vpq, "\t", v1.TypeOfVpq, "\t", v1.Number)
	}
	PrintTable.Flush()
}
