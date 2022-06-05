package main

import (
	"testing"
	"time"
)

func TestStatementStep1(t *testing.T) {
	keys, journal := FDbRead[SJournal](VDbJournal)
	FDbCloseAll()
	journalTimes := FConvertByteSliceToTime(keys)
	i1 := FStatementStep1(journalTimes, journal, time.Time{}, time.Now())
	FPrintMap6(i1)
}

func TestStatementStep2(t *testing.T) {
	keys, journal := FDbRead[SJournal](VDbJournal)
	FDbCloseAll()
	journalTimes := FConvertByteSliceToTime(keys)
	i1 := FStatementStep1(journalTimes, journal, time.Time{}, time.Now())
	i1 = FStatementStep2(i1)
	FPrintMap6(i1)
}

func TestStatementStep3(t *testing.T) {
	keys, journal := FDbRead[SJournal](VDbJournal)
	FDbCloseAll()
	journalTimes := FConvertByteSliceToTime(keys)
	i1 := FStatementStep1(journalTimes, journal, time.Time{}, time.Now())
	i1 = FStatementStep2(i1)
	a1 := FStatementStep3(i1)
	FPrintMap6(a1)
}

func TestStatementStep4(t *testing.T) {
	keys, journal := FDbRead[SJournal](VDbJournal)
	FDbCloseAll()
	journalTimes := FConvertByteSliceToTime(keys)
	i1 := FStatementStep1(journalTimes, journal, time.Time{}, time.Now())
	i1 = FStatementStep2(i1)
	i2 := FStatementStep3(i1)
	i2 = FStatementStep4(i2, true, []string{"yasa"})
	FPrintMap6(i2)
}

func TestStatementStep5(t *testing.T) {
	keys, journal := FDbRead[SJournal](VDbJournal)
	FDbCloseAll()
	journalTimes := FConvertByteSliceToTime(keys)
	i1 := FStatementStep1(journalTimes, journal, time.Time{}, time.Now())
	i1 = FStatementStep2(i1)
	i2 := FStatementStep3(i1)
	i2 = FStatementStep4(i2, true, []string{"yasa"})
	i3 := FStatementStep5(i2)
	FPrintMap5(i3)
}

func TestStatementStep6(t *testing.T) {
	keys, journal := FDbRead[SJournal](VDbJournal)
	FDbCloseAll()
	journalTimes := FConvertByteSliceToTime(keys)
	i1 := FStatementStep1(journalTimes, journal, time.Time{}, time.Now())
	i1 = FStatementStep2(i1)
	i2 := FStatementStep3(i1)
	i2 = FStatementStep4(i2, true, []string{"yasa"})
	i3 := FStatementStep5(i2)
	FStatementStep6(365, i3)
	FPrintMap5(i3)
}

func TestCalculatePrice(t *testing.T) {
	keys, journal := FDbRead[SJournal](VDbJournal)
	FDbCloseAll()
	journalTimes := FConvertByteSliceToTime(keys)
	i1 := FStatementStep1(journalTimes, journal, time.Time{}, time.Now())
	i1 = FStatementStep2(i1)
	i2 := FStatementStep3(i1)
	i2 = FStatementStep4(i2, true, []string{"yasa"})
	i3 := FStatementStep5(i2)
	FStatementStep6(365, i3)
	FPrintMap5(i3)
}

func TestStatementStep7(t *testing.T) {
	keys, journal := FDbRead[SJournal](VDbJournal)
	FDbCloseAll()
	journalTimes := FConvertByteSliceToTime(keys)
	i1 := FStatementStep1(journalTimes, journal, time.Time{}, time.Now())
	i1 = FStatementStep2(i1)
	i2 := FStatementStep3(i1)
	i2 = FStatementStep4(i2, true, []string{"yasa"})
	i3 := FStatementStep5(i2)
	FStatementStep6(365, i3)
	FStatementStep7(i3)
	FPrintMap5(i3)
}

func TestFinancialStatements(t *testing.T) {
	VCompanyName = "anti_accountants"
	FDbOpenAll()
	keys, journal := FDbRead[SJournal](VDbJournal)
	FDbCloseAll()
	FPrintJournal(journal)
	dates := FConvertByteSliceToTime(keys)
	dates, journal = FJournalFilter(dates, journal, SFilterJournal{}, false)
	a1, a2 := FStatement(dates, journal, []time.Time{time.Now()}, 1, []string{"yasa"}, true)
	a3 := FStatementFilter(a1[0], SFilterStatement{
		Account1: SFilterAccount{
			IsFilter:    false,
			IsCredit:    SFilterBool{IsFilter: false, BoolValue: false},
			Account:     SFilterString{IsFilter: false, Way: "", Slice: []string{}},
			FathersName: SFilterFathersAccountsName{IsFilter: false, InAccountName: false, InFathersName: false, FathersName: []string{"assets"}},
			Levels:      SFilterSliceUint{IsFilter: false, InSlice: false, Slice: []uint{}},
		},
		Account2: SFilterAccount{
			IsFilter:    false,
			IsCredit:    SFilterBool{IsFilter: false, BoolValue: false},
			Account:     SFilterString{IsFilter: true, Way: CInSlice, Slice: []string{CAllAccounts}},
			FathersName: SFilterFathersAccountsName{IsFilter: false, InAccountName: false, InFathersName: false, FathersName: []string{}},
			Levels:      SFilterSliceUint{IsFilter: false, InSlice: false, Slice: []uint{}},
		},
		Name:                   SFilterString{IsFilter: true, Way: CInSlice, Slice: []string{"zizi"}},
		Vpq:                    SFilterString{IsFilter: true, Way: CInSlice, Slice: []string{CValue}},
		TypeOfVpq:              SFilterString{IsFilter: true, Way: CInSlice, Slice: []string{CFlowEnding}},
		ChangeOrRatioOrBalance: SFilterString{IsFilter: true, Way: CInSlice, Slice: []string{CBalance}},
		Number:                 SFilterNumber{IsFilter: false, Way: "", Slice: []float64{}},
	})

	FTest(true, a2, nil)
	a4 := FConvertStatmentWithAccountToFilteredStatement(a3)
	FPrintStatement(a4)
}

func TestStatementFilterByGreedyAlgorithm(t *testing.T) {
	keys, journal := FDbRead[SJournal](VDbJournal)
	FDbCloseAll()
	dates := FConvertByteSliceToTime(keys)
	i1, _ := FStatement(dates, journal, []time.Time{time.Now()}, 1, []string{"yasa"}, true)
	a1 := FStatementFilter(i1[0], SFilterStatement{
		Account1: SFilterAccount{
			IsFilter:    false,
			IsCredit:    SFilterBool{IsFilter: false, BoolValue: false},
			Account:     SFilterString{IsFilter: false, Way: "", Slice: []string{}},
			FathersName: SFilterFathersAccountsName{IsFilter: false, InAccountName: false, InFathersName: false, FathersName: []string{"assets"}},
			Levels:      SFilterSliceUint{IsFilter: false, InSlice: false, Slice: []uint{}},
		},
		Account2: SFilterAccount{
			IsFilter:    false,
			IsCredit:    SFilterBool{IsFilter: false, BoolValue: false},
			Account:     SFilterString{IsFilter: true, Way: CInSlice, Slice: []string{CAllAccounts}},
			FathersName: SFilterFathersAccountsName{IsFilter: false, InAccountName: false, InFathersName: false, FathersName: []string{}},
			Levels:      SFilterSliceUint{IsFilter: false, InSlice: false, Slice: []uint{}},
		},
		Name:                   SFilterString{IsFilter: true, Way: CInSlice, Slice: []string{CAllNames, CNames}},
		Vpq:                    SFilterString{IsFilter: true, Way: CInSlice, Slice: []string{CValue}},
		TypeOfVpq:              SFilterString{IsFilter: true, Way: CInSlice, Slice: []string{CFlowEnding}},
		ChangeOrRatioOrBalance: SFilterString{},
		Number:                 SFilterNumber{IsFilter: false, Way: "", Slice: []float64{}},
	})

	a1 = FStatementFilterByGreedyAlgorithm(a1, true, 0.7)
	a2 := FConvertStatmentWithAccountToFilteredStatement(a1)
	FPrintStatement(a2)
}

func TestSortByLevel(t *testing.T) {
	keys, journal := FDbRead[SJournal](VDbJournal)
	FDbCloseAll()
	dates := FConvertByteSliceToTime(keys)
	i1, _ := FStatement(dates, journal, []time.Time{time.Now()}, 1, []string{"yasa"}, true)
	a1 := FStatementFilter(i1[0], SFilterStatement{
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

	a1 = FSortByLevel(a1)
	a2 := FConvertStatmentWithAccountToFilteredStatement(a1)
	FPrintStatement(a2)
}

func TestMakeSpaceBeforeAccountInStatementStruct(t *testing.T) {
	keys, journal := FDbRead[SJournal](VDbJournal)
	FDbCloseAll()
	dates := FConvertByteSliceToTime(keys)
	i1, _ := FStatement(dates, journal, []time.Time{time.Now()}, 1, []string{"yasa"}, true)
	a1 := FStatementFilter(i1[0], SFilterStatement{
		Account1: SFilterAccount{
			IsFilter:    true,
			IsCredit:    SFilterBool{IsFilter: false, BoolValue: false},
			Account:     SFilterString{IsFilter: false, Way: "", Slice: []string{}},
			FathersName: SFilterFathersAccountsName{IsFilter: true, InAccountName: true, InFathersName: true, FathersName: []string{"owner's equity"}},
			Levels:      SFilterSliceUint{IsFilter: false, InSlice: false, Slice: []uint{}},
		},
		Account2: SFilterAccount{
			IsFilter:    false,
			IsCredit:    SFilterBool{IsFilter: false, BoolValue: false},
			Account:     SFilterString{IsFilter: true, Way: CInSlice, Slice: []string{CAllAccounts}},
			FathersName: SFilterFathersAccountsName{IsFilter: false, InAccountName: false, InFathersName: false, FathersName: []string{}},
			Levels:      SFilterSliceUint{IsFilter: false, InSlice: false, Slice: []uint{}},
		},
		Name:                   SFilterString{IsFilter: true, Way: CInSlice, Slice: []string{CAllNames, CNames}},
		Vpq:                    SFilterString{IsFilter: true, Way: CInSlice, Slice: []string{CValue}},
		TypeOfVpq:              SFilterString{IsFilter: true, Way: CInSlice, Slice: []string{CFlowEnding}},
		ChangeOrRatioOrBalance: SFilterString{IsFilter: true, Way: CInSlice, Slice: []string{CBalance}},
		Number:                 SFilterNumber{IsFilter: false, Way: "", Slice: []float64{}},
	})

	a1 = FSortByLevel(a1)
	FMakeSpaceBeforeAccountInStatementStruct(a1)
	a2 := FConvertStatmentWithAccountToFilteredStatement(a1)
	FPrintStatement(a2)
}
