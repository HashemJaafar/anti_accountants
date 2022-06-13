package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"reflect"
	"runtime/debug"
	"sort"
	"strings"
	"text/tabwriter"
	"time"
)

func FCutTheSlice[t any](a []t, b int) []t { return a[:len(a)-b] }
func FRemove[t any](a []t, b int) []t      { return append(a[:b], a[b+1:]...) }
func FSwap[t any](a []t, b, c int)         { a[b], a[c] = a[c], a[b] }

// FNow this function to get the current time in the format of TimeLayout to make the error less likely
func FNow() []byte {
	return []byte(time.Now().Format(CTimeLayout))
}

func FInitializeMap6[t1, t2, t3, t4, t5, t6 comparable, tr any](m map[t1]map[t2]map[t3]map[t4]map[t5]map[t6]tr, i1 t1, i2 t2, i3 t3, i4 t4, i5 t5) map[t6]tr {
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

func FInitializeMap5[t1, t2, t3, t4, t5 comparable, tr any](m map[t1]map[t2]map[t3]map[t4]map[t5]tr, i1 t1, i2 t2, i3 t3, i4 t4) map[t5]tr {
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

func FFind[t comparable](element t, elements []t) (int, bool) {
	for k1, v1 := range elements {
		if v1 == element {
			return k1, true
		}
	}
	return 0, false
}

func FIsInfIn(numbers ...float64) bool {
	for _, v1 := range numbers {
		if math.IsInf(v1, 0) {
			return true
		}
	}
	return false
}

func FAbs[t INumber](n t) t {
	// this is alternative of math.Abs
	if n < 0 {
		return -n
	}
	return n
}

func FReturnSetAndDuplicatesSlices[t any](slice []t) ([]t, []t) {
	var setOfElems, duplicatedElement []t
big_loop:
	for _, v1 := range slice {
		for _, v2 := range setOfElems {
			if reflect.DeepEqual(v2, v1) {
				duplicatedElement = append(duplicatedElement, v1)
				continue big_loop
			}
		}
		setOfElems = append(setOfElems, v1)
	}
	return setOfElems, duplicatedElement
}

func FReverseSlice[t any](a []t) {
	for k1, k2 := 0, len(a)-1; k1 < k2; k1, k2 = k1+1, k2-1 {
		FSwap(a, k1, k2)
	}
}

func FSmallest[t INumber](a, b t) t {
	if a < b {
		return a
	}
	return b
}

func FSortTime(slice []time.Time, isAscending bool) {
	for k1 := range slice {
		for k2 := range slice {
			if k1 < k2 && (slice[k1]).After((slice[k2])) == isAscending {
				FSwap(slice, k1, k2)
			}
		}
	}
}

func FSortStatementNumber(slice []SStatmentWithAccount, isAscending bool) {
	sort.Slice(slice, func(k1, k2 int) bool {
		return slice[k1].SStatement1.Number > slice[k2].SStatement1.Number == isAscending
	})
}

func FConvertByteSliceToTime(slice [][]byte) []time.Time {
	var sliceOfTime []time.Time
	for _, v1 := range slice {
		date, _ := time.Parse(CTimeLayout, string(v1))
		sliceOfTime = append(sliceOfTime, date)
	}
	return sliceOfTime
}

func FConvertFromByteSliceToString(slice [][]byte) []string {
	var values []string
	for _, v1 := range slice {
		values = append(values, string(v1))
	}
	return values
}

func FTest[t any](shouldEqual bool, actual, expected t) {
	if reflect.DeepEqual(actual, expected) != shouldEqual {
		VFailTestNumber++
		p := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
		fmt.Fprintln(p, "\033[32m", "fail_test_number\t:", VFailTestNumber) //green
		fmt.Fprintln(p, "\033[35m", "should_equal\t:", shouldEqual)         //purple
		fmt.Fprintln(p, "\033[34m", "actual\t:", actual)                    //blue
		fmt.Fprintln(p, "\033[33m", "expected\t:", expected)                //yellow
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

func FPrintSlice[t any](a1 []t) {
	for _, v1 := range a1 {
		fmt.Println(v1)
	}
}

func FTranspose[t any](slice [][]t) [][]t {
	xl := len(slice[0])
	yl := len(slice)
	result := make([][]t, xl)
	for k1 := range result {
		result[k1] = make([]t, yl)
	}
	for k1 := 0; k1 < xl; k1++ {
		for k2 := 0; k2 < yl; k2++ {
			result[k1][k2] = slice[k2][k1]
		}
	}
	return result
}

func FUnpack[t any](slice [][]t) []t {
	var result []t
	for _, v1 := range slice {
		result = append(result, v1...)
	}
	return result
}

func FElementsInElement(element string, elements []string) bool {
	for _, v1 := range elements {
		if strings.Contains(element, v1) {
			return true
		}
	}
	return false
}

func FPrintMap6[t1, t2, t3, t4, t5, t6 comparable, tr any](m map[t1]map[t2]map[t3]map[t4]map[t5]map[t6]tr) {
	for k1, v1 := range m {
		for k2, v2 := range v1 {
			for k3, v3 := range v2 {
				for k4, v4 := range v3 {
					for k5, v5 := range v4 {
						for k6, v6 := range v5 {
							fmt.Fprintln(VPrintTable, k1, "\t", k2, "\t", k3, "\t", k4, "\t", k5, "\t", k6, "\t", v6)
						}
					}
				}
			}
		}
	}
	fmt.Println("//////////////////////////////////////////")
	VPrintTable.Flush()
}

func FPrintMap5[t1, t2, t3, t4, t5 comparable, tr any](m map[t1]map[t2]map[t3]map[t4]map[t5]tr) {
	for k1, v1 := range m {
		for k2, v2 := range v1 {
			for k3, v3 := range v2 {
				for k4, v4 := range v3 {
					for k5, v5 := range v4 {
						fmt.Fprintln(VPrintTable, k1, "\t", k2, "\t", k3, "\t", k4, "\t", k5, "\t", v5)
					}
				}
			}
		}
	}
	fmt.Println("//////////////////////////////////////////")
	VPrintTable.Flush()
}

func FPrintMap2[t1, t2 comparable, tr any](m map[t1]map[t2]tr) {
	for k1, v1 := range m {
		for k2, v2 := range v1 {
			fmt.Fprintln(VPrintTable, k1, "\t", k2, "\t", v2)
		}
	}
	fmt.Println("//////////////////////////////////////////")
	VPrintTable.Flush()
}

func FPrintStatement(slice []SStatement1) {
	fmt.Fprintln(VPrintTable, "Account1", "\t", "Account2", "\t", "Name", "\t", "Vpq", "\t", "TypeOfVpq", "\t", "ChangeOrRatioOrBalance", "\t", "Number")
	for _, v1 := range slice {
		fmt.Fprintln(VPrintTable, v1.Account1Name, "\t", v1.Account2Name, "\t", v1.PersonName, "\t", v1.Vpq, "\t", v1.TypeOfVpq, "\t", v1.ChangeOrRatioOrBalance, "\t", v1.Number)
	}
	VPrintTable.Flush()
}

func FPrintJournal(slice []SJournal1) {
	for _, v1 := range slice {
		fmt.Fprintln(VPrintTable,
			"\t", v1.IsReverse,
			"\t", v1.IsReversed,
			"\t", v1.ReverseEntryNumberCompound,
			"\t", v1.ReverseEntryNumberSimple,
			"\t", v1.EntryNumberCompound,
			"\t", v1.EntryNumberSimple,
			"\t", v1.Value,
			"\t", v1.DebitPrice,
			"\t", v1.CreditPrice,
			"\t", v1.DebitQuantity,
			"\t", v1.CreditQuantity,
			"\t", v1.DebitAccountName,
			"\t", v1.CreditAccountName,
			"\t", v1.Notes,
			"\t", v1.Name,
			"\t", v1.Employee,
			"\t", v1.TypeOfCompoundEntry)
	}
	VPrintTable.Flush()
}

func FPrintCvp(a SCvp) {
	fmt.Fprintln(VPrintTable, "VariableCost", "\t", a.VariableCost)
	fmt.Fprintln(VPrintTable, "FixedCost", "\t", a.FixedCost)
	fmt.Fprintln(VPrintTable, "MixedCost", "\t", a.MixedCost)
	fmt.Fprintln(VPrintTable, "Sales", "\t", a.Sales)
	fmt.Fprintln(VPrintTable, "Profit", "\t", a.Profit)
	fmt.Fprintln(VPrintTable, "ContributionMargin", "\t", a.ContributionMargin)
	VPrintTable.Flush()
}

func FFilesName(dir string) ([]string, error) {
	var filesName []string
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return filesName, err
	}
	for _, v1 := range files {
		filesName = append(filesName, v1.Name())
	}
	return filesName, nil
}

func FMaxMinDate(slice []time.Time) (time.Time, time.Time) {
	var max, min time.Time
	for _, v1 := range slice {
		if v1.After(max) {
			max = v1
		}
		if v1.Before(min) {
			min = v1
		}
	}
	return max, min
}

func FFilterDuplicate[t comparable](input1, input2 t, isFilter bool) bool {
	return !isFilter || input1 == input2
}

func FFilterBool(input bool, f SFilterBool) bool {
	return !f.IsFilter || input == f.BoolValue
}

func FFilterNumber[t uint | float64](input t, f SFilter[t]) bool {
	switch f.Way {
	case CDontFilter:
		return true
	case CBetween:
		max, min := FMaxMin(f.Slice)
		return input >= min && input <= max
	case CNotBetween:
		max, min := FMaxMin(f.Slice)
		return !(input >= min && input <= max)
	case CBigger:
		max, _ := FMaxMin(f.Slice)
		return input > max
	case CSmaller:
		_, min := FMaxMin(f.Slice)
		return input < min
	case CInSlice:
		_, isIn := FFind(input, f.Slice)
		return isIn
	case CNotInSlice:
		_, isIn := FFind(input, f.Slice)
		return !isIn
	}

	return true
}

func FFilterTime(input time.Time, f SFilter[time.Time]) bool {
	switch f.Way {
	case CDontFilter:
		return true
	case CBetween:
		max, min := FMaxMinDate(f.Slice)
		return input.After(min) && input.Before(max)
	case CNotBetween:
		max, min := FMaxMinDate(f.Slice)
		return !(input.After(min) && input.Before(max))
	case CBigger:
		max, _ := FMaxMinDate(f.Slice)
		return input.After(max)
	case CSmaller:
		_, min := FMaxMinDate(f.Slice)
		return input.Before(min)
	case CInSlice:
		_, isIn := FFind(input, f.Slice)
		return isIn
	case CNotInSlice:
		_, isIn := FFind(input, f.Slice)
		return !isIn
	}

	return true
}

func FFilterString(input string, f SFilter[string]) bool {
	switch f.Way {
	case CDontFilter:
		return true
	case CInSlice:
		_, isIn := FFind(input, f.Slice)
		return isIn
	case CNotInSlice:
		_, isIn := FFind(input, f.Slice)
		return !isIn
	case CElementsInElement:
		return FElementsInElement(input, f.Slice)
	case CElementsNotInElement:
		return !FElementsInElement(input, f.Slice)
	}

	return true
}

func FFilterSlice[t uint | string](input []t, f SFilter[t]) bool {
	switch f.Way {
	case CDontFilter:
		return true
	case CInSlice:
		for _, v1 := range input {
			_, isIn := FFind(v1, f.Slice)
			return isIn
		}
	case CNotInSlice:
		for _, v1 := range input {
			_, isIn := FFind(v1, f.Slice)
			return !isIn
		}
	}

	return true
}

func FFilterAccount(input SAccount1, f SAccount2) bool {
	if FFilterBool(input.IsCredit, f.IsCredit) &&
		FFilterString(input.CostFlowType, f.CostFlowType) &&
		FFilterString(input.Inventory, f.Inventory) &&
		FFilterString(input.Name, f.Name) &&
		FFilterString(input.Notes, f.Notes) &&
		FFilterSlice(input.Image, f.Image) &&
		FFilterSlice(input.Barcode, f.Barcode) &&
		FFilterSlice(input.Number[VIndexOfAccountNumber], f.Number) &&
		FFilterNumber(input.Levels[VIndexOfAccountNumber], f.Levels) &&
		FFilterSlice(input.FathersName[VIndexOfAccountNumber], f.FathersName) {
		return true
	}
	return false
}

func FSetSliceOfTErr(err []error) []TErr {
	var slice []TErr
	for _, v1 := range err {
		if v1 != nil {
			slice = append(slice, TErr(v1.Error()))
		} else {
			slice = append(slice, TErr(""))
		}
	}
	return slice
}
