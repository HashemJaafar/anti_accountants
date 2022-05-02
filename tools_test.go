package main

import (
	"testing"
)

func TestFunctionFilterDuplicate(t *testing.T) {
	a1 := FunctionFilterDuplicate("lkjds", "ojdi", true)
	Test(true, a1, false)
	a1 = FunctionFilterDuplicate(4496, 546, true)
	Test(true, a1, false)
	a1 = FunctionFilterDuplicate(true, true, true)
	Test(true, a1, true)
	a1 = FunctionFilterDuplicate("lkjds", "ojdi", false)
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

	a1 := FilterSliceString{true, true, i1}.Filter(i1)
	Test(true, a1, true)
	a1 = FilterSliceString{true, false, i1}.Filter(i1)
	Test(true, a1, false)
	a1 = FilterSliceString{true, true, slice}.Filter(i1)
	Test(true, a1, false)
	a1 = FilterSliceString{true, false, slice}.Filter(i1)
	Test(true, a1, true)
}
