package main

import (
	"testing"
)

func TestFilter(t *testing.T) {
	a1 := FFilterDuplicate("lkjds", "ojdi", true)
	FTest(true, a1, false)
	a1 = FFilterDuplicate(4496, 546, true)
	FTest(true, a1, false)
	a1 = FFilterDuplicate(true, true, true)
	FTest(true, a1, true)
	a1 = FFilterDuplicate("lkjds", "ojdi", false)
	FTest(true, a1, true)
}

func TestFilterNumber_Filter(t *testing.T) {
	// a1 := SFilter[float64]{IsFilter: true, Way: CBetween, Slice: []float64{100, 500}}.FFilter(1000)
	// FTest(true, a1, false)
	// a1 = SFilter[float64]{IsFilter: true, Way: CNotBetween, Slice: []float64{100, 500}}.FFilter(1000)
	// FTest(true, a1, true)
	// a1 = SFilter[float64]{IsFilter: true, Way: CBigger, Slice: []float64{100, 500}}.FFilter(1000)
	// FTest(true, a1, true)
	// a1 = SFilter[float64]{IsFilter: true, Way: CSmaller, Slice: []float64{100, 500}}.FFilter(1000)
	// FTest(true, a1, false)
	// a1 = SFilter[float64]{IsFilter: true, Way: CInSlice, Slice: []float64{100, 500}}.FFilter(1000)
	// FTest(true, a1, false)
	// a1 = SFilter[float64]{IsFilter: true, Way: CNotInSlice, Slice: []float64{100, 500}}.FFilter(1000)
	// FTest(true, a1, true)
}

func TestFilterSliceString_Filter(t *testing.T) {
	// i1 := []string{"1", "2", "3"}
	// slice := []string{"4", "5", "6"}

	// a1 := SFilterFathersAccountsName{true, false, true, i1}.FFilter("1", i1)
	// FTest(true, a1, false)
	// a1 = SFilterFathersAccountsName{true, false, false, i1}.FFilter("1", i1)
	// FTest(true, a1, false)
	// a1 = SFilterFathersAccountsName{true, false, true, slice}.FFilter("1", i1)
	// FTest(true, a1, false)
	// a1 = SFilterFathersAccountsName{true, false, false, slice}.FFilter("1", i1)
	// FTest(true, a1, true)
	// a1 = SFilterFathersAccountsName{true, true, true, i1}.FFilter("1", []string{})
	// FTest(true, a1, true)
	// a1 = SFilterFathersAccountsName{true, false, true, i1}.FFilter("1", []string{})
	// FTest(true, a1, false)
}
