package anti_accountants

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

func TestSTATEMENT_STEP_2(t *testing.T) {
	keys, journal := DB_READ[JOURNAL_TAG](DB_JOURNAL)
	DB_CLOSE()
	journal_times := CONVERT_BYTE_SLICE_TO_TIME(keys)
	i1 := STATEMENT_STEP_1(journal_times, journal, time.Time{}, time.Now())
	print_map_6(i1)
	i1 = STATEMENT_STEP_2(i1)
	print_map_6(i1)
}

func TestSTATEMENT_STEP_3(t *testing.T) {
	keys, journal := DB_READ[JOURNAL_TAG](DB_JOURNAL)
	DB_CLOSE()
	journal_times := CONVERT_BYTE_SLICE_TO_TIME(keys)
	i1 := STATEMENT_STEP_1(journal_times, journal, time.Time{}, time.Now())
	print_map_6(i1)
	i1 = STATEMENT_STEP_2(i1)
	print_map_6(i1)
	a1 := STATEMENT_STEP_3(i1)
	print_map_6(a1)
}

func TestSTATEMENT_STEP_4(t *testing.T) {
	keys, journal := DB_READ[JOURNAL_TAG](DB_JOURNAL)
	DB_CLOSE()
	journal_times := CONVERT_BYTE_SLICE_TO_TIME(keys)
	i1 := STATEMENT_STEP_1(journal_times, journal, time.Time{}, time.Now())
	print_map_6(i1)
	i1 = STATEMENT_STEP_2(i1)
	print_map_6(i1)
	i2 := STATEMENT_STEP_3(i1)
	print_map_6(i2)
	i2 = STATEMENT_STEP_4(true, []string{"yasa"}, i2)
	print_map_6(i2)
}

func TestSTATEMENT_STEP_5(t *testing.T) {
	keys, journal := DB_READ[JOURNAL_TAG](DB_JOURNAL)
	DB_CLOSE()
	journal_times := CONVERT_BYTE_SLICE_TO_TIME(keys)
	i1 := STATEMENT_STEP_1(journal_times, journal, time.Time{}, time.Now())
	print_map_6(i1)
	i1 = STATEMENT_STEP_2(i1)
	print_map_6(i1)
	i2 := STATEMENT_STEP_3(i1)
	print_map_6(i2)
	i2 = STATEMENT_STEP_4(true, []string{"yasa"}, i2)
	print_map_6(i2)
	i3 := STATEMENT_STEP_5(i2)
	print_map_5(i3)
}

func TestVERTICAL_ANALYSIS(t *testing.T) {
	keys, journal := DB_READ[JOURNAL_TAG](DB_JOURNAL)
	DB_CLOSE()
	journal_times := CONVERT_BYTE_SLICE_TO_TIME(keys)
	i1 := STATEMENT_STEP_1(journal_times, journal, time.Time{}, time.Now())
	print_map_6(i1)
	i1 = STATEMENT_STEP_2(i1)
	print_map_6(i1)
	i2 := STATEMENT_STEP_3(i1)
	print_map_6(i2)
	i2 = STATEMENT_STEP_4(true, []string{"yasa"}, i2)
	print_map_6(i2)
	i3 := STATEMENT_STEP_5(i2)
	print_map_5(i3)
	VERTICAL_ANALYSIS(365, i3)
	print_map_5(i3)
	// for k1, v1 := range i3 {
	// 	for k2, v2 := range v1 {
	// 		for k3, v3 := range v2 {
	// 			for k4, v4 := range v3 {
	// 				for k5, v5 := range v4 {
	// 					if k5 == beginning_balance {
	// 						fmt.Fprintln(PRINT_TABLE, k1, "\t", k2, "\t", k3, "\t", k4, "\t", k5, "\t", v5)
	// 					}
	// 				}
	// 			}
	// 		}
	// 	}
	// }
	// fmt.Println("//////////////////////////////////////////")
	// PRINT_TABLE.Flush()
}
