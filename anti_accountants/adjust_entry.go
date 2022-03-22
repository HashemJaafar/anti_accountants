package anti_accountants

import (
	"math"
	"strings"
	"time"
)

func SET_SLICE_DAY_START_END(array_day_start_end []DAY_START_END) {
	if len(array_day_start_end) == 0 {
		array_day_start_end = []DAY_START_END{
			{SATURDAY, 0, 0, 23, 59},
			{SUNDAY, 0, 0, 23, 59},
			{MONDAY, 0, 0, 23, 59},
			{TUESDAY, 0, 0, 23, 59},
			{WEDNESDAY, 0, 0, 23, 59},
			{THURSDAY, 0, 0, 23, 59},
			{FRIDAY, 0, 0, 23, 59}}
	}
	for index := range array_day_start_end {
		array_day_start_end[index].DAY = strings.Title(array_day_start_end[index].DAY)

		if !IS_IN(array_day_start_end[index].DAY, STANDARD_DAYS) {
			array_day_start_end[index].DAY = SUNDAY
		}

		if array_day_start_end[index].START_HOUR < 0 {
			array_day_start_end[index].START_HOUR = 0
		}
		if array_day_start_end[index].START_HOUR > 23 {
			array_day_start_end[index].START_HOUR = 23
		}
		if array_day_start_end[index].START_MINUTE < 0 {
			array_day_start_end[index].START_MINUTE = 0
		}
		if array_day_start_end[index].START_MINUTE > 59 {
			array_day_start_end[index].START_MINUTE = 59
		}
		if array_day_start_end[index].END_HOUR < 0 {
			array_day_start_end[index].END_HOUR = 0
		}
		if array_day_start_end[index].END_HOUR > 23 {
			array_day_start_end[index].END_HOUR = 23
		}
		if array_day_start_end[index].END_MINUTE < 0 {
			array_day_start_end[index].END_MINUTE = 0
		}
		if array_day_start_end[index].END_MINUTE > 59 {
			array_day_start_end[index].END_MINUTE = 59
		}

		if array_day_start_end[index].START_HOUR > array_day_start_end[index].END_HOUR {
			array_day_start_end[index].START_HOUR = 0
		}
		if array_day_start_end[index].START_HOUR == array_day_start_end[index].END_HOUR && array_day_start_end[index].START_MINUTE > array_day_start_end[index].END_MINUTE {
			array_day_start_end[index].START_MINUTE = 0
		}
	}
}

func CREATE_ARRAY_START_END_MINUTES(date_end, date_start time.Time, array_day_start_end []DAY_START_END) []start_end_minutes {
	var array_start_end_minutes []start_end_minutes
	var previous_date_end time.Time
	delta_days := int(date_end.Sub(date_start).Hours()/24 + 1)
	year, month, day := date_start.Date()
	for day_counter := 0; day_counter < delta_days; day_counter++ {
		for _, element := range array_day_start_end {
			start := time.Date(year, month, day+day_counter, element.START_HOUR, element.START_MINUTE, 0, 0, time.Local)
			if start.Weekday().String() == element.DAY {
				end := time.Date(year, month, day+day_counter, element.END_HOUR, element.END_MINUTE, 0, 0, time.Local)
				start, end = SHIFT_AND_ARRANGE_THE_TIME_SERIES(previous_date_end, start, end)
				array_start_end_minutes = append(array_start_end_minutes, start_end_minutes{start, end, end.Sub(start).Minutes()})
				previous_date_end = end
			}
		}
	}
	return array_start_end_minutes
}

func SHIFT_AND_ARRANGE_THE_TIME_SERIES(previous_date_end, date_start, date_end time.Time) (time.Time, time.Time) {
	if previous_date_end.After(date_start) {
		date_start = previous_date_end
	}
	if date_start.After(date_end) {
		date_end = date_start
	}
	return date_start, date_end
}

func TOTAL_MINUTES(array_start_end_minutes []start_end_minutes) float64 {
	var TOTAL_MINUTES float64
	for _, element := range array_start_end_minutes {
		TOTAL_MINUTES += element.minutes
	}
	return TOTAL_MINUTES
}

func ADJUST_THE_ARRAY(array_to_insert []JOURNAL_TAG, array_start_end_minutes []start_end_minutes, adjusting_method string) [][]JOURNAL_TAG {
	var adjusted_array_to_insert [][]JOURNAL_TAG
	TOTAL_MINUTES := TOTAL_MINUTES(array_start_end_minutes)
	array_len_start_end_minutes := len(array_start_end_minutes) - 1
	for _, entry := range array_to_insert {
		var value_counter, time_unit_counter float64
		var one_account_adjusted_list []JOURNAL_TAG
		for index, element := range array_start_end_minutes {
			value := VALUE_AFTER_ADJUST_USING_ADJUSTING_METHODS(adjusting_method, element.minutes, TOTAL_MINUTES, time_unit_counter, entry.VALUE)

			if index == array_len_start_end_minutes {
				value = entry.VALUE - value_counter
			}

			time_unit_counter += element.minutes
			value_counter += value
			one_account_adjusted_list = append(one_account_adjusted_list, JOURNAL_TAG{
				REVERSE:               false,
				ENTRY_NUMBER:          0,
				ENTRY_NUMBER_COMPOUND: index,
				ENTRY_NUMBER_SIMPLE:   0,
				VALUE:                 value,
				PRICE_DEBIT:           entry.PRICE_DEBIT,
				PRICE_CREDIT:          entry.PRICE_CREDIT,
				QUANTITY_DEBIT:        RETURN_SAME_SIGN_OF_NUMBER_SIGN(value/entry.PRICE_DEBIT, entry.QUANTITY_DEBIT),
				QUANTITY_CREDIT:       RETURN_SAME_SIGN_OF_NUMBER_SIGN(value/entry.PRICE_CREDIT, entry.QUANTITY_CREDIT),
				ACCOUNT_DEBIT:         entry.ACCOUNT_DEBIT,
				ACCOUNT_CREDIT:        entry.ACCOUNT_CREDIT,
				NOTES:                 entry.NOTES,
				NAME:                  entry.NAME,
				NAME_EMPLOYEE:         entry.NAME_EMPLOYEE,
				DATE_START:            element.date_start,
				DATE_END:              element.date_end,
			})
		}
		adjusted_array_to_insert = append(adjusted_array_to_insert, one_account_adjusted_list)
	}
	return adjusted_array_to_insert
}

func VALUE_AFTER_ADJUST_USING_ADJUSTING_METHODS(adjusting_method string, minutes, TOTAL_MINUTES, time_unit_counter, total_value float64) float64 {
	percent := ROOT(total_value, TOTAL_MINUTES)
	switch adjusting_method {
	case EXPONENTIAL:
		return math.Pow(percent, time_unit_counter+minutes) - math.Pow(percent, time_unit_counter)
	case LOGARITHMIC:
		return (total_value / math.Pow(percent, time_unit_counter)) - (total_value / math.Pow(percent, time_unit_counter+minutes))
	default:
		return minutes * (total_value / TOTAL_MINUTES)
	}
}
