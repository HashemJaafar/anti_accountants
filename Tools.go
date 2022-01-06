package anti_accountants

import (
	"database/sql"
	"log"
	"math"
	"reflect"
	"time"
)

func IS_IN(element string, elements []string) bool {
	for _, a := range elements {
		if a == element {
			return true
		}
	}
	return false
}

func transpose(slice [][]journal_tag) [][]journal_tag {
	xl := len(slice[0])
	yl := len(slice)
	result := make([][]journal_tag, xl)
	for i := range result {
		result[i] = make([]journal_tag, yl)
	}
	for i := 0; i < xl; i++ {
		for j := 0; j < yl; j++ {
			result[i][j] = slice[j][i]
		}
	}
	return result
}

func unpack_the_array(adjusted_array_to_insert [][]journal_tag) []journal_tag {
	array_to_insert := []journal_tag{}
	for _, element := range adjusted_array_to_insert {
		array_to_insert = append(array_to_insert, element...)
	}
	return array_to_insert
}

func RETURN_SAME_SIGN_OF_NUMBER_SIGN(number_sign, number float64) float64 {
	if number_sign < 0 {
		number = -math.Abs(number)
	} else {
		number = math.Abs(number)
	}
	return number
}

func PARSE_DATE(string_date string, date_layouts []string) time.Time {
	for _, i := range date_layouts {
		date, err := time.Parse(i, string_date)
		if err == nil {
			return date
		}
	}
	return time.Time{}
}

func RETURN_SET_AND_DUPLICATES_SLICES(slice_of_elements []string) ([]string, []string) {
	var set_of_elems, duplicated_element []string
big_loop:
	for _, element := range slice_of_elements {
		for _, b := range set_of_elems {
			if b == element {
				duplicated_element = append(duplicated_element, element)
				continue big_loop
			}
		}
		set_of_elems = append(set_of_elems, element)
	}
	return set_of_elems, duplicated_element
}

func CONCAT(args ...interface{}) interface{} {
	n := 0
	for _, arg := range args {
		n += reflect.ValueOf(arg).Len()
	}
	v := reflect.MakeSlice(reflect.TypeOf(args[0]), 0, n)
	for _, arg := range args {
		v = reflect.AppendSlice(v, reflect.ValueOf(arg))
	}
	return v.Interface()
}

func initialize_map_4(m map[string]map[string]map[string]map[string]map[string]float64, a, b, c, d string) map[string]float64 {
	if m[a] == nil {
		m[a] = map[string]map[string]map[string]map[string]float64{}
	}
	if m[a][b] == nil {
		m[a][b] = map[string]map[string]map[string]float64{}
	}
	if m[a][b][c] == nil {
		m[a][b][c] = map[string]map[string]float64{}
	}
	if m[a][b][c][d] == nil {
		m[a][b][c][d] = map[string]float64{}
	}
	return m[a][b][c][d]
}

func initialize_map_3(m map[string]map[string]map[string]map[string]float64, a, b, c string) map[string]float64 {
	if m[a] == nil {
		m[a] = map[string]map[string]map[string]float64{}
	}
	if m[a][b] == nil {
		m[a][b] = map[string]map[string]float64{}
	}
	if m[a][b][c] == nil {
		m[a][b][c] = map[string]float64{}
	}
	return m[a][b][c]
}

func error_fatal(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func check_dates(start_date, end_date time.Time) {
	if !start_date.Before(end_date) {
		log.Panic("please enter the start_date<=end_date")
	}
}

func select_from_journal(rows *sql.Rows) []journal_tag {
	var journal []journal_tag
	for rows.Next() {
		var tag journal_tag
		rows.Scan(&tag.DATE, &tag.ENTRY_NUMBER, &tag.ACCOUNT, &tag.VALUE, &tag.PRICE, &tag.QUANTITY, &tag.BARCODE, &tag.ENTRY_EXPAIR, &tag.DESCRIPTION, &tag.NAME, &tag.EMPLOYEE_NAME, &tag.ENTRY_DATE, &tag.REVERSE)
		journal = append(journal, tag)
	}
	return journal
}
