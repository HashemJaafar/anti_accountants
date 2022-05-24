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

func FCutTheSlice[t any](a []t, b int) []t { return a[:len(a)-b] }
func FRemove[t any](a []t, b int) []t      { return append(a[:b], a[b+1:]...) }
func FSwap[t any](a []t, b, c int)         { a[b], a[c] = a[c], a[b] }

// FNow this function to get the current time in the format of TimeLayout to make the error less likely
func FNow() []byte {
	return []byte(time.Now().Format(CTimeLayout))
}

func FAssignNumberIfNumber(m map[string]float64, str string) {
	number, err := strconv.ParseFloat(str, 64)
	if err == nil {
		m[str] = number
	}
}

func FConvertNanToZero(VALUE float64) float64 {
	if math.IsNaN(VALUE) {
		return 0
	}
	return VALUE
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

func FReturnSameSignOfNumberSign(numberSign, number float64) float64 {
	if numberSign < 0 {
		return -FAbs(number)
	}
	return FAbs(number)
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
		return slice[k1].SStatement.TNumber > slice[k2].SStatement.TNumber == isAscending
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

func (s SFilterDate) FFilter(input time.Time) bool {
	if !s.IsFilter {
		return true
	}

	if s.Big.Before(s.Small) {
		s.Big, s.Small = s.Small, s.Big
	}

	isSmallerThanSmall := input.Before(s.Small)
	isBiggerThanBig := input.After(s.Big)

	switch s.Way {
	case CBetween:
		return !isSmallerThanSmall && !isBiggerThanBig
	case CNotBetween:
		return isSmallerThanSmall || isBiggerThanBig
	case CBigger:
		return input.After(s.Big)
	case CSmaller:
		return input.Before(s.Small)
	case CEqualToOneOfThem:
		return input == s.Small || input == s.Big
	}

	return false
}

func (s SFilterNumber) FFilter(input float64) bool {
	if !s.IsFilter {
		return true
	}

	if s.Big < s.Small {
		s.Big, s.Small = s.Small, s.Big
	}

	isSmallerThanSmall := input < s.Small
	isBiggerThanBig := input > s.Big

	switch s.Way {
	case CBetween:
		return !isSmallerThanSmall && !isBiggerThanBig
	case CNotBetween:
		return isSmallerThanSmall || isBiggerThanBig
	case CBigger:
		return input > s.Big
	case CSmaller:
		return input < s.Small
	case CEqualToOneOfThem:
		return input == s.Small || input == s.Big
	}

	return false
}

func (s SFilterString) FFilter(input string) bool {
	if !s.IsFilter {
		return true
	}

	switch s.Way {
	case CInSlice:
		_, isIn := FFind(input, s.Slice)
		return isIn
	case CNotInSlice:
		_, isIn := FFind(input, s.Slice)
		return isIn == false
	case CElementsInElement:
		return FElementsInElement(input, s.Slice)
	case CElementsNotInElement:
		return FElementsInElement(input, s.Slice) == false
	}
	return false
}

func (s SFilterAccount) FFilter(account SAccount, err error) bool {
	if !s.IsFilter {
		return true
	}

	return (err != nil) || // here if the account is not listed in the account list like AllAccounts it will show in statment
		(s.IsCredit.FFilter(account.TIsCredit) &&
			s.FathersName.FFilter(account.TAccountName, account.TAccountFathersName[VIndexOfAccountNumber]) &&
			s.Levels.FFilter(account.TAccountLevels[VIndexOfAccountNumber]))
}

func (s SFilterFathersAccountsName) FFilter(accountName string, fathersAccountsNameForAccount []string) bool {
	if !s.IsFilter {
		return true
	}
	for _, v1 := range s.FathersName {
		if v1 == accountName { // if accountName is in the slice
			return s.InAccountName
		}
		for _, v2 := range fathersAccountsNameForAccount {
			if v1 == v2 {
				return s.InFathersName
			}
		}
	}
	return !s.InFathersName
}

func (s SFilterSliceUint) FFilter(input uint) bool {
	if !s.IsFilter {
		return true
	}
	_, isIn := FFind(input, s.Slice)
	return isIn == s.InSlice
}

func (s SFilterBool) FFilter(input bool) bool {
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

func FFilterDuplicate[t comparable](input1, input2 t, f bool) bool {
	if !f {
		return true
	}
	return input1 == input2
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

func FPrintStatement(slice []SStatement) {
	fmt.Fprintln(VPrintTable, "Account1", "\t", "Account2", "\t", "Name", "\t", "Vpq", "\t", "TypeOfVpq", "\t", "ChangeOrRatioOrBalance", "\t", "Number")
	for _, v1 := range slice {
		fmt.Fprintln(VPrintTable, v1.TAccount1Name, "\t", v1.TAccount2Name, "\t", v1.TPersonName, "\t", v1.TVpq, "\t", v1.TTypeOfVpq, "\t", v1.TChangeOrRatioOrBalance, "\t", v1.TNumber)
	}
	VPrintTable.Flush()
}

func FPrintJournal(slice []SJournal) {
	for _, v1 := range slice {
		fmt.Fprintln(VPrintTable,
			"\t", v1.IsReverse,
			"\t", v1.IsReversed,
			"\t", v1.ReverseEntryNumberCompound,
			"\t", v1.ReverseEntryNumberSimple,
			"\t", v1.EntryNumberCompound,
			"\t", v1.EntryNumberSimple,
			"\t", v1.Value,
			"\t", v1.PriceDebit,
			"\t", v1.PriceCredit,
			"\t", v1.QuantityDebit,
			"\t", v1.QuantityCredit,
			"\t", v1.AccountDebit,
			"\t", v1.AccountCredit,
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
