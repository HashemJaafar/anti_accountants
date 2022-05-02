package main

import (
	"testing"
	"time"
)

func TestStatementStep1(t *testing.T) {
	keys, journal := DbRead[JournalTag](DbJournal)
	DbClose()
	journalTimes := ConvertByteSliceToTime(keys)
	i1 := StatementStep1(journalTimes, journal, time.Time{}, time.Now())
	PrintMap6(i1)
}

func TestStatementStep3(t *testing.T) {
	keys, journal := DbRead[JournalTag](DbJournal)
	DbClose()
	journalTimes := ConvertByteSliceToTime(keys)
	i1 := StatementStep1(journalTimes, journal, time.Time{}, time.Now())
	i1 = StatementStep3(i1)
	PrintMap6(i1)
}

func TestStatementStep4(t *testing.T) {
	keys, journal := DbRead[JournalTag](DbJournal)
	DbClose()
	journalTimes := ConvertByteSliceToTime(keys)
	i1 := StatementStep1(journalTimes, journal, time.Time{}, time.Now())
	i1 = StatementStep3(i1)
	a1 := StatementStep4(i1)
	PrintMap6(a1)
}

func TestStatementStep5(t *testing.T) {
	keys, journal := DbRead[JournalTag](DbJournal)
	DbClose()
	journalTimes := ConvertByteSliceToTime(keys)
	i1 := StatementStep1(journalTimes, journal, time.Time{}, time.Now())
	i1 = StatementStep3(i1)
	i2 := StatementStep4(i1)
	i2 = StatementStep5(true, []string{"yasa"}, i2)
	PrintMap6(i2)
}

func TestStatementStep6(t *testing.T) {
	keys, journal := DbRead[JournalTag](DbJournal)
	DbClose()
	journalTimes := ConvertByteSliceToTime(keys)
	i1 := StatementStep1(journalTimes, journal, time.Time{}, time.Now())
	i1 = StatementStep3(i1)
	i2 := StatementStep4(i1)
	i2 = StatementStep5(true, []string{"yasa"}, i2)
	i3 := StatementStep6(i2)
	PrintMap5(i3)
}

func TestStatementStep7(t *testing.T) {
	keys, journal := DbRead[JournalTag](DbJournal)
	DbClose()
	journalTimes := ConvertByteSliceToTime(keys)
	i1 := StatementStep1(journalTimes, journal, time.Time{}, time.Now())
	i1 = StatementStep3(i1)
	i2 := StatementStep4(i1)
	i2 = StatementStep5(true, []string{"yasa"}, i2)
	i3 := StatementStep6(i2)
	StatementStep7(365, i3)
	PrintMap5(i3)
}

func TestCalculatePrice(t *testing.T) {
	keys, journal := DbRead[JournalTag](DbJournal)
	DbClose()
	journalTimes := ConvertByteSliceToTime(keys)
	i1 := StatementStep1(journalTimes, journal, time.Time{}, time.Now())
	i1 = StatementStep3(i1)
	i2 := StatementStep4(i1)
	i2 = StatementStep5(true, []string{"yasa"}, i2)
	i3 := StatementStep6(i2)
	StatementStep7(365, i3)
	CalculatePrice(i3)
	PrintMap5(i3)
}

func TestStatementStep2(t *testing.T) {
	keys, journal := DbRead[JournalTag](DbJournal)
	DbClose()
	journalTimes := ConvertByteSliceToTime(keys)
	i1 := StatementStep1(journalTimes, journal, time.Now(), time.Now())
	i1 = StatementStep2(i1)
	PrintMap6(i1)
}

func TestStatementStep8(t *testing.T) {
	keys, journal := DbRead[JournalTag](DbJournal)
	DbClose()
	journalTimes := ConvertByteSliceToTime(keys)
	i1 := StatementStep1(journalTimes, journal, time.Time{}, time.Now())
	i1 = StatementStep3(i1)
	i2 := StatementStep4(i1)
	i2 = StatementStep5(true, []string{"yasa"}, i2)
	i3 := StatementStep6(i2)
	StatementStep7(365, i3)
	StatementStep8(i3)
	PrintMap5(i3)
}

func TestFinancialStatements(t *testing.T) {
	a1, a2 := FinancialStatements([]time.Time{time.Now()}, 1, []string{"yasa"}, true)
	Test(true, a2, nil)
	for _, v1 := range a1 {
		PrintMap5(v1)
	}
}

func TestStatementFilter(t *testing.T) {
	i1, _ := FinancialStatements([]time.Time{time.Now()}, 1, []string{"yasa"}, true)
	DbClose()

	a1 := StatementFilter(i1[0], FilterStatement{
		Account1:  FilterString{IsFilter: false, Way: "", Slice: []string{}},
		Account2:  FilterString{IsFilter: true, Way: InSlice, Slice: []string{AllAccounts}},
		Name:      FilterString{IsFilter: true, Way: InSlice, Slice: []string{AllNames, Names}},
		Vpq:       FilterString{IsFilter: true, Way: InSlice, Slice: []string{Value}},
		TypeOfVpq: FilterString{IsFilter: true, Way: InSlice, Slice: []string{EndingBalance}},
		Number:    FilterNumber{IsFilter: false, Way: "", Big: 0, Small: 0},
	})
	PrintFilteredStatement(a1)
}

func TestStatementFilterByGreedyAlgorithm(t *testing.T) {
	i1 := []FilteredStatement{
		{"cash", "", "", "", "", 200},
		{"book", "", "", "", "", 50},
		{"rent", "", "", "", "", 10},
		{"book1", "", "", "", "", 288},
		{"book2", "", "", "", "", 30},
	}
	a1 := StatementFilterByGreedyAlgorithm(i1, true, 0.7)
	PrintSlice(a1)
}

func TestStatementFilterAccounts(t *testing.T) {
	i1, _ := FinancialStatements([]time.Time{time.Now()}, 10, []string{"yasa"}, true)
	DbClose()

	a1 := StatementFilter(i1[0], FilterStatement{
		Account1:  FilterString{IsFilter: false, Way: "", Slice: []string{}},
		Account2:  FilterString{IsFilter: true, Way: InSlice, Slice: []string{AllAccounts}},
		Name:      FilterString{IsFilter: true, Way: InSlice, Slice: []string{AllNames, Names}},
		Vpq:       FilterString{IsFilter: true, Way: InSlice, Slice: []string{Value}},
		TypeOfVpq: FilterString{IsFilter: true, Way: InSlice, Slice: []string{EndingBalance}},
		Number:    FilterNumber{IsFilter: false, Way: "", Big: 0, Small: 0},
	})

	i2 := FilterAccount{
		IsLowLevelAccount:   FilterBool{IsFilter: true, BoolValue: true},
		IsCredit:            FilterBool{IsFilter: true, BoolValue: true},
		IsTemporary:         FilterBool{IsFilter: true, BoolValue: true},
		FathersAccountsName: FilterSliceString{},
		AccountLevels:       FilterSliceUint{},
	}

	i3 := FilterAccount{
		IsLowLevelAccount:   FilterBool{},
		IsCredit:            FilterBool{},
		IsTemporary:         FilterBool{},
		FathersAccountsName: FilterSliceString{},
		AccountLevels:       FilterSliceUint{},
	}
	a1 = StatementFilterAccounts(a1, i2, i3)
	PrintFilteredStatement(a1)
	a1 = StatementFilterByGreedyAlgorithm(a1, true, 0.7)
	PrintFilteredStatement(a1)
}

func TestSortByLevel(t *testing.T) {
	i1, _ := FinancialStatements([]time.Time{time.Now()}, 10, []string{"yasa"}, true)
	DbClose()
	a1 := StatementFilter(i1[0], FilterStatement{
		Account1:  FilterString{IsFilter: false, Way: "", Slice: []string{}},
		Account2:  FilterString{IsFilter: true, Way: InSlice, Slice: []string{AllAccounts}},
		Name:      FilterString{IsFilter: true, Way: InSlice, Slice: []string{AllNames, Names}},
		Vpq:       FilterString{IsFilter: true, Way: InSlice, Slice: []string{Value}},
		TypeOfVpq: FilterString{IsFilter: true, Way: InSlice, Slice: []string{EndingBalance}},
		Number:    FilterNumber{IsFilter: false, Way: "", Big: 0, Small: 0},
	})
	SortByLevel(a1)
	PrintFilteredStatement(a1)
}
