package anti_accountants

import (
	"log"
	"math"
	"strings"
	"time"
)

func initialize_array_day_start_end(array_day_start_end []DAY_START_END) []DAY_START_END {
	if len(array_day_start_end) == 0 {
		array_day_start_end = []DAY_START_END{
			{"saturday", 0, 0, 23, 59},
			{"sunday", 0, 0, 23, 59},
			{"monday", 0, 0, 23, 59},
			{"tuesday", 0, 0, 23, 59},
			{"wednesday", 0, 0, 23, 59},
			{"thursday", 0, 0, 23, 59},
			{"friday", 0, 0, 23, 59}}
	}
	return array_day_start_end
}

func check_array_day_start_end(array_day_start_end []DAY_START_END) {
	for index, element := range array_day_start_end {
		array_day_start_end[index].DAY = strings.Title(element.DAY)
		switch {
		case !is_in(array_day_start_end[index].DAY, standard_days):
			error_element_is_not_in_elements(element.DAY, standard_days)
		case element.START_HOUR < 0:
			error_the_time_is_not_in_range(element, 0)
		case element.START_HOUR > 23:
			error_the_time_is_not_in_range(element, 23)
		case element.START_MINUTE < 0:
			error_the_time_is_not_in_range(element, 0)
		case element.START_MINUTE > 59:
			error_the_time_is_not_in_range(element, 59)
		case element.END_HOUR < 0:
			error_the_time_is_not_in_range(element, 0)
		case element.END_HOUR > 23:
			error_the_time_is_not_in_range(element, 23)
		case element.END_MINUTE < 0:
			error_the_time_is_not_in_range(element, 0)
		case element.END_MINUTE > 59:
			error_the_time_is_not_in_range(element, 59)
		}
	}
}

func create_array_start_end_minutes(entry_expair, date time.Time, array_day_start_end []DAY_START_END) []start_end_minutes {
	var array_start_end_minutes []start_end_minutes
	var previous_end_date, end time.Time
	delta_days := int(entry_expair.Sub(date).Hours()/24 + 1)
	year, month_string, day := date.Date()
	for day_counter := 0; day_counter < delta_days; day_counter++ {
		for index, element := range array_day_start_end {
			start := time.Date(year, month_string, day+day_counter, element.START_HOUR, element.START_MINUTE, 0, 0, time.Local)
			if start.Weekday().String() == element.DAY {
				end = time.Date(year, month_string, day+day_counter, element.END_HOUR, element.END_MINUTE, 0, 0, time.Local)

				if previous_end_date.After(start) {
					log.Panic("the end_hour and end_minute for ", array_day_start_end[index-1], " should be smaller than start_hour and start_minute for the second ", array_day_start_end[index])
				}
				if start.After(end) {
					log.Panic("the start_hour and start_minute should be smaller than end_hour and end_minute for ", array_day_start_end[index])
				}

				// fmt.Println(element.DAY, start, end, end.Sub(start).Minutes())
				array_start_end_minutes = append(array_start_end_minutes, start_end_minutes{start, end, end.Sub(start).Minutes()})
				previous_end_date = end
			}
		}
	}
	return array_start_end_minutes
}

func total_minutes(array_start_end_minutes []start_end_minutes) float64 {
	var total_minutes float64
	for _, element := range array_start_end_minutes {
		total_minutes += element.minutes
	}
	return total_minutes
}

func adjust_the_array(array_to_insert []JOURNAL_TAG, array_start_end_minutes []start_end_minutes, total_minutes float64, adjusting_method, description, name, employee_name string) [][]JOURNAL_TAG {
	var adjusted_array_to_insert [][]JOURNAL_TAG
	array_len_start_end_minutes := len(array_start_end_minutes) - 1
	for _, entry := range array_to_insert {
		var value_counter, time_unit_counter float64
		var one_account_adjusted_list []JOURNAL_TAG
		total_value := math.Abs(entry.VALUE)
		for index, element := range array_start_end_minutes {
			value := value_after_adjust_using_adjusting_methods(adjusting_method, element, total_minutes, time_unit_counter, total_value)

			if index == array_len_start_end_minutes {
				value = math.Abs(total_value - value_counter)
			}

			time_unit_counter += element.minutes
			value_counter += math.Abs(value)
			value = return_same_sign_of_number_sign(entry.VALUE, value)
			one_account_adjusted_list = append(one_account_adjusted_list, JOURNAL_TAG{
				REVERSE:               false,
				ENTRY_NUMBER:          0,
				ENTRY_NUMBER_COMPOUND: 0,
				ENTRY_NUMBER_SIMPLE:   0,
				VALUE:                 value,
				PRICE_DEBIT:           entry.PRICE_DEBIT,
				PRICE_CREDIT:          entry.PRICE_CREDIT,
				QUANTITY_DEBIT:        value / entry.PRICE_DEBIT,
				QUANTITY_CREDIT:       value / entry.PRICE_CREDIT,
				ACCOUNT_DEBIT:         entry.ACCOUNT_DEBIT,
				ACCOUNT_CREDIT:        entry.ACCOUNT_CREDIT,
				NOTES:                 "",
				NAME:                  name,
				NAME_EMPLOYEE:         name,
				DATE_START:            time.Time{},
				DATE_END:              time.Time{},
				DATE_ENTRY:            time.Time{},
			})
		}
		adjusted_array_to_insert = append(adjusted_array_to_insert, one_account_adjusted_list)
	}
	return adjusted_array_to_insert
}

func value_after_adjust_using_adjusting_methods(adjusting_method string, element start_end_minutes, total_minutes, time_unit_counter, total_value float64) float64 {
	percent := root(total_value, total_minutes)
	switch adjusting_method {
	case "linear":
		return element.minutes * (total_value / total_minutes)
	case "exponential":
		return math.Pow(percent, time_unit_counter+element.minutes) - math.Pow(percent, time_unit_counter)
	case "logarithmic":
		return (total_value / math.Pow(percent, time_unit_counter)) - (total_value / math.Pow(percent, time_unit_counter+element.minutes))
	}
	return 0
}
