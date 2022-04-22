package main

import (
	"testing"
	"time"
)

func TestSTATEMENT_STEP_1(t *testing.T) {
	keys, journal := DB_READ[JOURNAL_TAG](DB_JOURNAL)
	DB_CLOSE()
	journal_times := CONVERT_BYTE_SLICE_TO_TIME(keys)
	i1 := STATEMENT_STEP_1(journal_times, journal, time.Time{}, time.Now())
	print_map_6(i1)
}

func TestSTATEMENT_STEP_3(t *testing.T) {
	keys, journal := DB_READ[JOURNAL_TAG](DB_JOURNAL)
	DB_CLOSE()
	journal_times := CONVERT_BYTE_SLICE_TO_TIME(keys)
	i1 := STATEMENT_STEP_1(journal_times, journal, time.Time{}, time.Now())
	i1 = STATEMENT_STEP_3(i1)
	print_map_6(i1)
}

func TestSTATEMENT_STEP_4(t *testing.T) {
	keys, journal := DB_READ[JOURNAL_TAG](DB_JOURNAL)
	DB_CLOSE()
	journal_times := CONVERT_BYTE_SLICE_TO_TIME(keys)
	i1 := STATEMENT_STEP_1(journal_times, journal, time.Time{}, time.Now())
	i1 = STATEMENT_STEP_3(i1)
	a1 := STATEMENT_STEP_4(i1)
	print_map_6(a1)
}

func TestSTATEMENT_STEP_5(t *testing.T) {
	keys, journal := DB_READ[JOURNAL_TAG](DB_JOURNAL)
	DB_CLOSE()
	journal_times := CONVERT_BYTE_SLICE_TO_TIME(keys)
	i1 := STATEMENT_STEP_1(journal_times, journal, time.Time{}, time.Now())
	i1 = STATEMENT_STEP_3(i1)
	i2 := STATEMENT_STEP_4(i1)
	i2 = STATEMENT_STEP_5(true, []string{"yasa"}, i2)
	print_map_6(i2)
}

func TestSTATEMENT_STEP_6(t *testing.T) {
	keys, journal := DB_READ[JOURNAL_TAG](DB_JOURNAL)
	DB_CLOSE()
	journal_times := CONVERT_BYTE_SLICE_TO_TIME(keys)
	i1 := STATEMENT_STEP_1(journal_times, journal, time.Time{}, time.Now())
	i1 = STATEMENT_STEP_3(i1)
	i2 := STATEMENT_STEP_4(i1)
	i2 = STATEMENT_STEP_5(true, []string{"yasa"}, i2)
	i3 := STATEMENT_STEP_6(i2)
	print_map_5(i3)
}

func TestSTATEMENT_STEP_7(t *testing.T) {
	keys, journal := DB_READ[JOURNAL_TAG](DB_JOURNAL)
	DB_CLOSE()
	journal_times := CONVERT_BYTE_SLICE_TO_TIME(keys)
	i1 := STATEMENT_STEP_1(journal_times, journal, time.Time{}, time.Now())
	i1 = STATEMENT_STEP_3(i1)
	i2 := STATEMENT_STEP_4(i1)
	i2 = STATEMENT_STEP_5(true, []string{"yasa"}, i2)
	i3 := STATEMENT_STEP_6(i2)
	STATEMENT_STEP_7(365, i3)
	print_map_5(i3)
}

func TestCALCULATE_PRICE(t *testing.T) {
	keys, journal := DB_READ[JOURNAL_TAG](DB_JOURNAL)
	DB_CLOSE()
	journal_times := CONVERT_BYTE_SLICE_TO_TIME(keys)
	i1 := STATEMENT_STEP_1(journal_times, journal, time.Time{}, time.Now())
	i1 = STATEMENT_STEP_3(i1)
	i2 := STATEMENT_STEP_4(i1)
	i2 = STATEMENT_STEP_5(true, []string{"yasa"}, i2)
	i3 := STATEMENT_STEP_6(i2)
	STATEMENT_STEP_7(365, i3)
	CALCULATE_PRICE(i3)
	print_map_5(i3)
}

func TestSTATEMENT_STEP_2(t *testing.T) {
	keys, journal := DB_READ[JOURNAL_TAG](DB_JOURNAL)
	DB_CLOSE()
	journal_times := CONVERT_BYTE_SLICE_TO_TIME(keys)
	i1 := STATEMENT_STEP_1(journal_times, journal, time.Now(), time.Now())
	i1 = STATEMENT_STEP_2(i1, SET_RETAINED_EARNINGS_ACCOUNT(ACCOUNT{ACCOUNT_NAME: "retined_earnings hh"}).ACCOUNT_NAME)
	print_map_6(i1)
}

func TestSTATEMENT_STEP_8(t *testing.T) {
	keys, journal := DB_READ[JOURNAL_TAG](DB_JOURNAL)
	DB_CLOSE()
	journal_times := CONVERT_BYTE_SLICE_TO_TIME(keys)
	i1 := STATEMENT_STEP_1(journal_times, journal, time.Time{}, time.Now())
	i1 = STATEMENT_STEP_3(i1)
	i2 := STATEMENT_STEP_4(i1)
	i2 = STATEMENT_STEP_5(true, []string{"yasa"}, i2)
	i3 := STATEMENT_STEP_6(i2)
	STATEMENT_STEP_7(365, i3)
	STATEMENT_STEP_8(i3)
	print_map_5(i3)
}

func TestFINANCIAL_STATEMENTS(t *testing.T) {
	a1, a2 := FINANCIAL_STATEMENTS([]time.Time{time.Now()}, 1, []string{"yasa"}, true, ACCOUNT{ACCOUNT_NAME: "retined_earnings"})
	TEST(true, a2, nil)
	for _, v1 := range a1 {
		print_map_5(v1)
	}
}

func TestSTATEMENT_FILTER(t *testing.T) {
	retained_earnings_account := ACCOUNT{ACCOUNT_NAME: "retined_earnings", ACCOUNT_LEVELS: []uint{1, 5}}
	i1, _ := FINANCIAL_STATEMENTS([]time.Time{time.Now()}, 10, []string{"yasa"}, true, retained_earnings_account)
	DB_CLOSE()
	// retained_earnings_account = SET_RETAINED_EARNINGS_ACCOUNT(retained_earnings_account)
	// ACCOUNTS = append(ACCOUNTS, retained_earnings_account)
	// SET_THE_ACCOUNTS()
	PRINT_FORMATED_ACCOUNTS()

	a1 := STATEMENT_FILTER(THE_FILTER_OF_THE_STATEMENT{
		old_statement:         i1[0],
		account1:              []string{},
		account2:              []string{},
		name:                  []string{},
		vpq:                   []string{},
		type_of_vpq:           []string{},
		account1_levels:       []uint{},
		account2_levels:       []uint{},
		is_in_account1:        false,
		is_in_account2:        false,
		is_in_name:            false,
		is_in_vpq:             false,
		is_in_type_of_vpq:     false,
		is_in_account1_levels: false,
		is_in_account2_levels: false,
	})
	PRINT_SLICE(a1)
}
