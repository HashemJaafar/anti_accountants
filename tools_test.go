package main

import (
	"testing"
)

func TestFilter(t *testing.T) {
	a1 := FilterDuplicate("lkjds", "ojdi", true)
	Test(true, a1, false)
	a1 = FilterDuplicate(4496, 546, true)
	Test(true, a1, false)
	a1 = FilterDuplicate(true, true, true)
	Test(true, a1, true)
	a1 = FilterDuplicate("lkjds", "ojdi", false)
	Test(true, a1, true)
}

func TestFilterNumber_Filter(t *testing.T) {
	a1 := FilterNumber{IsFilter: true, Way: Between, Big: 900, Small: 0}.Filter(1000)
	Test(true, a1, false)
	a1 = FilterNumber{IsFilter: true, Way: NotBetween, Big: 900, Small: 0}.Filter(1000)
	Test(true, a1, true)
	a1 = FilterNumber{IsFilter: true, Way: Bigger, Big: 900, Small: 0}.Filter(1000)
	Test(true, a1, true)
	a1 = FilterNumber{IsFilter: true, Way: Smaller, Big: 900, Small: 0}.Filter(1000)
	Test(true, a1, false)
	a1 = FilterNumber{IsFilter: true, Way: EqualToOneOfThem, Big: 900, Small: 0}.Filter(1000)
	Test(true, a1, false)
}

func TestFilterSliceString_Filter(t *testing.T) {
	i1 := []string{"1", "2", "3"}
	slice := []string{"4", "5", "6"}

	a1 := FilterFathersAccountsName{true, false, true, i1}.Filter("1", i1)
	Test(true, a1, false)
	a1 = FilterFathersAccountsName{true, false, false, i1}.Filter("1", i1)
	Test(true, a1, false)
	a1 = FilterFathersAccountsName{true, false, true, slice}.Filter("1", i1)
	Test(true, a1, false)
	a1 = FilterFathersAccountsName{true, false, false, slice}.Filter("1", i1)
	Test(true, a1, true)
	a1 = FilterFathersAccountsName{true, true, true, i1}.Filter("1", []string{})
	Test(true, a1, true)
	a1 = FilterFathersAccountsName{true, false, true, i1}.Filter("1", []string{})
	Test(true, a1, false)
}
