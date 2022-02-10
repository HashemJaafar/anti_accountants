package anti_accountants

import (
	"time"
)

func (s FINANCIAL_ACCOUNTING) JOURNAL_FILTER(
	JOURNAL_ORDERED_BY_DATE_ENTRY_NUMBER []JOURNAL_TAG,

	filter_date, in_date bool,
	min_date, max_date time.Time,

	filter_entry_number, in_entry_number bool,
	min_entry_number, max_entry_number int,

	filter_account, in_account bool,
	account []string,

	filter_value, in_value bool,
	min_value, max_value float64,

	filter_price, in_price bool,
	min_price, max_price float64,

	filter_quantity, in_quantity bool,
	min_quantity, max_quantity float64,

	filter_barcode, in_barcode bool,
	barcode []string,

	filter_entry_expair, in_entry_expair bool,
	min_entry_expair, max_entry_expair time.Time,

	filter_description, in_description bool,
	description []string,

	filter_name, in_name bool,
	name []string,

	filter_employee_name, in_employee_name bool,
	employee_name []string,

	filter_entry_date, in_entry_date bool,
	min_entry_date, max_entry_date time.Time,

	filter_reverse,
	reverse bool,

) []JOURNAL_TAG {

	check_dates(min_date, max_date)
	check_dates(min_entry_expair, max_entry_expair)
	check_dates(min_entry_date, max_entry_date)

	var filtered_journal []JOURNAL_TAG

	for _, entry := range JOURNAL_ORDERED_BY_DATE_ENTRY_NUMBER {

		date := PARSE_DATE(entry.DATE, s.DATE_LAYOUT)
		entry_expair := PARSE_DATE(entry.ENTRY_EXPAIR, s.DATE_LAYOUT)
		entry_date := PARSE_DATE(entry.ENTRY_DATE, s.DATE_LAYOUT)

		if filter_date || (in_date == (date.After(min_date) && date.Before(max_date))) {
			if filter_entry_number || (in_entry_number == (entry.ENTRY_NUMBER >= min_entry_number && entry.ENTRY_NUMBER <= max_entry_number)) {
				if filter_account || (in_account == IS_IN(entry.ACCOUNT, account)) {
					if filter_value || (in_value == (entry.VALUE >= min_value && entry.VALUE <= max_value)) {
						if filter_price || (in_price == (entry.PRICE >= min_price && entry.PRICE <= max_price)) {
							if filter_quantity || (in_quantity == (entry.QUANTITY >= min_quantity && entry.QUANTITY <= max_quantity)) {
								if filter_barcode || (in_barcode == IS_IN(entry.BARCODE, barcode)) {
									if filter_entry_expair || (in_entry_expair == (entry_expair.After(min_entry_expair) && entry_expair.Before(max_entry_expair))) {
										if filter_description || (in_description == IS_IN(entry.DESCRIPTION, description)) {
											if filter_name || (in_name == IS_IN(entry.NAME, name)) {
												if filter_employee_name || (in_employee_name == IS_IN(entry.EMPLOYEE_NAME, employee_name)) {
													if filter_entry_date || (in_entry_date == (entry_date.After(min_entry_date) && entry_date.Before(max_entry_date))) {
														if filter_reverse || reverse == entry.REVERSE {
															filtered_journal = append(filtered_journal, entry)
														}
													}
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return filtered_journal
}
