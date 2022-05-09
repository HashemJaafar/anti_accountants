package main

import (
	"testing"
	"time"
)

func TestStatementStep1(t *testing.T) {
	keys, journal := DbRead[Journal](DbJournal)
	DbClose()
	journalTimes := ConvertByteSliceToTime(keys)
	i1 := StatementStep1(journalTimes, journal, time.Time{}, time.Now())
	PrintMap6(i1)
}

func TestStatementStep2(t *testing.T) {
	keys, journal := DbRead[Journal](DbJournal)
	DbClose()
	journalTimes := ConvertByteSliceToTime(keys)
	i1 := StatementStep1(journalTimes, journal, time.Time{}, time.Now())
	i1 = StatementStep2(i1)
	PrintMap6(i1)
}

func TestStatementStep3(t *testing.T) {
	keys, journal := DbRead[Journal](DbJournal)
	DbClose()
	journalTimes := ConvertByteSliceToTime(keys)
	i1 := StatementStep1(journalTimes, journal, time.Time{}, time.Now())
	i1 = StatementStep2(i1)
	a1 := StatementStep3(i1)
	PrintMap6(a1)
}

func TestStatementStep4(t *testing.T) {
	keys, journal := DbRead[Journal](DbJournal)
	DbClose()
	journalTimes := ConvertByteSliceToTime(keys)
	i1 := StatementStep1(journalTimes, journal, time.Time{}, time.Now())
	i1 = StatementStep2(i1)
	i2 := StatementStep3(i1)
	i2 = StatementStep4(true, []string{"yasa"}, i2)
	PrintMap6(i2)
}

func TestStatementStep5(t *testing.T) {
	keys, journal := DbRead[Journal](DbJournal)
	DbClose()
	journalTimes := ConvertByteSliceToTime(keys)
	i1 := StatementStep1(journalTimes, journal, time.Time{}, time.Now())
	i1 = StatementStep2(i1)
	i2 := StatementStep3(i1)
	i2 = StatementStep4(true, []string{"yasa"}, i2)
	i3 := StatementStep5(i2)
	PrintMap5(i3)
}

func TestStatementStep6(t *testing.T) {
	keys, journal := DbRead[Journal](DbJournal)
	DbClose()
	journalTimes := ConvertByteSliceToTime(keys)
	i1 := StatementStep1(journalTimes, journal, time.Time{}, time.Now())
	i1 = StatementStep2(i1)
	i2 := StatementStep3(i1)
	i2 = StatementStep4(true, []string{"yasa"}, i2)
	i3 := StatementStep5(i2)
	StatementStep6(365, i3)
	PrintMap5(i3)
}

func TestCalculatePrice(t *testing.T) {
	keys, journal := DbRead[Journal](DbJournal)
	DbClose()
	journalTimes := ConvertByteSliceToTime(keys)
	i1 := StatementStep1(journalTimes, journal, time.Time{}, time.Now())
	i1 = StatementStep2(i1)
	i2 := StatementStep3(i1)
	i2 = StatementStep4(true, []string{"yasa"}, i2)
	i3 := StatementStep5(i2)
	StatementStep6(365, i3)
	PrintMap5(i3)
}

func TestStatementStep7(t *testing.T) {
	keys, journal := DbRead[Journal](DbJournal)
	DbClose()
	journalTimes := ConvertByteSliceToTime(keys)
	i1 := StatementStep1(journalTimes, journal, time.Time{}, time.Now())
	i1 = StatementStep2(i1)
	i2 := StatementStep3(i1)
	i2 = StatementStep4(true, []string{"yasa"}, i2)
	i3 := StatementStep5(i2)
	StatementStep6(365, i3)
	StatementStep7(i3)
	PrintMap5(i3)
}

func TestFinancialStatements(t *testing.T) {
	a1, a2 := FinancialStatements([]time.Time{time.Now()}, 1, []string{"yasa"}, true, false)
	Test(true, a2, nil)
	for _, v1 := range a1 {
		PrintMap5(v1)
	}
}

func TestStatementFilterByGreedyAlgorithm(t *testing.T) {
	i1, _ := FinancialStatements([]time.Time{time.Now()}, 10, []string{"yasa"}, true, false)
	DbClose()
	a1 := StatementFilter(i1[0], FilterStatement{
		Account1: FilterAccount{
			IsFilter:    false,
			IsLowLevel:  FilterBool{IsFilter: true, BoolValue: true},
			IsCredit:    FilterBool{IsFilter: false, BoolValue: false},
			Account:     FilterString{IsFilter: false, Way: "", Slice: []string{}},
			FathersName: FilterFathersAccountsName{IsFilter: false, InAccountName: false, InFathersName: false, FathersName: []string{"assets"}},
			Levels:      FilterSliceUint{IsFilter: false, InSlice: false, Slice: []uint{}},
		},
		Account2: FilterAccount{
			IsFilter:    false,
			IsLowLevel:  FilterBool{IsFilter: false, BoolValue: false},
			IsCredit:    FilterBool{IsFilter: false, BoolValue: false},
			Account:     FilterString{IsFilter: true, Way: InSlice, Slice: []string{AllAccounts}},
			FathersName: FilterFathersAccountsName{IsFilter: false, InAccountName: false, InFathersName: false, FathersName: []string{}},
			Levels:      FilterSliceUint{IsFilter: false, InSlice: false, Slice: []uint{}},
		},
		Name:                   FilterString{IsFilter: true, Way: InSlice, Slice: []string{AllNames, Names}},
		Vpq:                    FilterString{IsFilter: true, Way: InSlice, Slice: []string{Value}},
		TypeOfVpq:              FilterString{IsFilter: true, Way: InSlice, Slice: []string{FlowEnding}},
		ChangeOrRatioOrBalance: FilterString{},
		Number:                 FilterNumber{IsFilter: false, Way: "", Big: 0, Small: 0},
	})

	a1 = StatementFilterByGreedyAlgorithm(a1, true, 0.7)
	a2 := ConvertStatmentWithAccountToFilteredStatement(a1)
	PrintStatement(a2)
}

func TestSortByLevel(t *testing.T) {
	i1, _ := FinancialStatements([]time.Time{time.Now()}, 10, []string{"yasa"}, true, false)
	DbClose()
	a1 := StatementFilter(i1[0], FilterStatement{
		// Account1: FilterAccount{
		// 	IsLowLevel:   FilterBool{IsFilter: false, BoolValue: false},
		// 	IsCredit:            FilterBool{IsFilter: false, BoolValue: false},
		// 	IsTemporary:         FilterBool{IsFilter: false, BoolValue: false},
		// 	Account:             FilterString{IsFilter: false, Way: "", Slice: []string{}},
		// 	FathersName: FilterFathersAccountsName{IsFilter: false, InAccountName: false, InFathersName: false, FathersName: []string{"assets"}},
		// 	Levels:       FilterSliceUint{IsFilter: false, InSlice: false, Slice: []uint{}},
		// },
		// Account2: FilterAccount{
		// 	IsLowLevel:   FilterBool{IsFilter: false, BoolValue: false},
		// 	IsCredit:            FilterBool{IsFilter: false, BoolValue: false},
		// 	IsTemporary:         FilterBool{IsFilter: false, BoolValue: false},
		// 	Account:             FilterString{IsFilter: true, Way: NotInSlice, Slice: []string{AllAccounts}},
		// 	FathersName: FilterFathersAccountsName{IsFilter: false, InAccountName: false, InFathersName: false, FathersName: []string{}},
		// 	Levels:       FilterSliceUint{IsFilter: false, InSlice: false, Slice: []uint{}},
		// },
		// Name:      FilterString{IsFilter: false, Way: InSlice, Slice: []string{AllNames, Names}},
		// Vpq:       FilterString{IsFilter: false, Way: InSlice, Slice: []string{Value}},
		// TypeOfVpq: FilterString{IsFilter: false, Way: InSlice, Slice: []string{FlowEnding}},
		// Number:    FilterNumber{IsFilter: false, Way: "", Big: 0, Small: 0},
	})

	a1 = SortByLevel(a1)
	a2 := ConvertStatmentWithAccountToFilteredStatement(a1)
	PrintStatement(a2)
}

func TestMakeSpaceBeforeAccountInStatementStruct(t *testing.T) {
	i1, _ := FinancialStatements([]time.Time{time.Now()}, 1, []string{"yasa"}, true, false)
	DbClose()
	a1 := StatementFilter(i1[0], FilterStatement{
		Account1: FilterAccount{
			IsFilter:    true,
			IsLowLevel:  FilterBool{IsFilter: false, BoolValue: false},
			IsCredit:    FilterBool{IsFilter: false, BoolValue: false},
			Account:     FilterString{IsFilter: false, Way: "", Slice: []string{}},
			FathersName: FilterFathersAccountsName{IsFilter: true, InAccountName: true, InFathersName: true, FathersName: []string{"owner's equity"}},
			Levels:      FilterSliceUint{IsFilter: false, InSlice: false, Slice: []uint{}},
		},
		Account2: FilterAccount{
			IsFilter:    false,
			IsLowLevel:  FilterBool{IsFilter: false, BoolValue: false},
			IsCredit:    FilterBool{IsFilter: false, BoolValue: false},
			Account:     FilterString{IsFilter: true, Way: InSlice, Slice: []string{AllAccounts}},
			FathersName: FilterFathersAccountsName{IsFilter: false, InAccountName: false, InFathersName: false, FathersName: []string{}},
			Levels:      FilterSliceUint{IsFilter: false, InSlice: false, Slice: []uint{}},
		},
		Name:                   FilterString{IsFilter: true, Way: InSlice, Slice: []string{AllNames, Names}},
		Vpq:                    FilterString{IsFilter: true, Way: InSlice, Slice: []string{Value}},
		TypeOfVpq:              FilterString{IsFilter: true, Way: InSlice, Slice: []string{FlowEnding}},
		ChangeOrRatioOrBalance: FilterString{IsFilter: true, Way: InSlice, Slice: []string{Balance}},
		Number:                 FilterNumber{IsFilter: false, Way: "", Big: 0, Small: 0},
	})

	a1 = SortByLevel(a1)
	MakeSpaceBeforeAccountInStatementStruct(a1)
	a2 := ConvertStatmentWithAccountToFilteredStatement(a1)
	PrintStatement(a2)
}
