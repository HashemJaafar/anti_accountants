package main

import (
	"fmt"
	"time"
)

// func AdjustTheArray(array_to_insert []JournalTag, array_start_end_minutes []start_end_minutes, adjusting_method string) [][]JournalTag {
// 	var adjusted_array_to_insert [][]JournalTag
// 	total_minutes := TOTAL_MINUTES(array_start_end_minutes)
// 	array_len_start_end_minutes := len(array_start_end_minutes) - 1
// 	for _, entry := range array_to_insert {
// 		var VALUE_counter, time_unit_counter float64
// 		var one_account_adjusted_list []JournalTag
// 		for index, element := range array_start_end_minutes {
// 			VALUE := VALUE_AFTER_ADJUST_USING_ADJUSTING_METHODS(adjusting_method, element.minutes, total_minutes, time_unit_counter, entry.VALUE)

// 			if index == array_len_start_end_minutes {
// 				VALUE = entry.VALUE - VALUE_counter
// 			}

// 			time_unit_counter += element.minutes
// 			VALUE_counter += VALUE
// 			one_account_adjusted_list = append(one_account_adjusted_list, JournalTag{
// 				IsReversed:           false,
// 				EntryNumberCompound: index,
// 				EntryNumberSimple:   0,
// 				VALUE:                 VALUE,
// 				PriceDebit:           entry.PriceDebit,
// 				PriceCredit:          entry.PriceCredit,
// 				QuantityDebit:        RETURN_SAME_SIGN_OF_NUMBER_SIGN(VALUE/entry.PriceDebit, entry.QuantityDebit),
// 				QuantityCredit:       RETURN_SAME_SIGN_OF_NUMBER_SIGN(VALUE/entry.PriceCredit, entry.QuantityCredit),
// 				AccountDebit:         entry.AccountDebit,
// 				AccountCredit:        entry.AccountCredit,
// 				NOTES:                 entry.NOTES,
// 				NAME:                  entry.NAME,
// 				NAME_EMPLOYEE:         entry.NAME_EMPLOYEE,
// 			})
// 		}
// 		adjusted_array_to_insert = append(adjusted_array_to_insert, one_account_adjusted_list)
// 	}
// 	return adjusted_array_to_insert
// }

// func CreateArrayStartEndMinutes(date_start, date_end time.Time, array_day_start_end []DAY_START_END) []start_end_minutes {
// 	var array_start_end_minutes []start_end_minutes
// 	var previous_date_end time.Time
// 	delta_days := int(date_end.Sub(date_start).Hours()/24 + 1)
// 	year, month, day := date_start.Date()
// 	for day_counter := 0; day_counter < delta_days; day_counter++ {
// 		for _, element := range array_day_start_end {
// 			start := time.Date(year, month, day+day_counter, element.START_HOUR, element.START_MINUTE, 0, 0, time.Local)
// 			if start.Weekday().String() == element.DAY {
// 				end := time.Date(year, month, day+day_counter, element.END_HOUR, element.END_MINUTE, 0, 0, time.Local)
// 				start, end = SHIFT_AND_ARRANGE_THE_TIME_SERIES(previous_date_end, start, end)
// 				array_start_end_minutes = append(array_start_end_minutes, start_end_minutes{start, end, end.Sub(start).Minutes()})
// 				previous_date_end = end
// 			}
// 		}
// 	}
// 	return array_start_end_minutes
// }

// func SetAdjustingMethod(entry_expair time.Time, adjusting_method string) string {
// 	if entry_expair.IsZero() {
// 		return ""
// 	}
// 	if !IS_IN(adjusting_method, DEPRECIATION_METHODS) {
// 		return Linear
// 	}
// 	return adjusting_method
// }

// func SetDateEndToZeroIfSmallerThanDateStart(date_start, date_end time.Time) time.Time {
// 	if !date_end.IsZero() && date_start.Before(date_end) {
// 		return time.Time{}
// 	}
// 	return date_end
// }

// func SetSliceDayStartEnd(array_day_start_end []DAY_START_END) []DAY_START_END {
// 	if len(array_day_start_end) == 0 {
// 		array_day_start_end = []DAY_START_END{
// 			{Saturday, 0, 0, 23, 59},
// 			{Sunday, 0, 0, 23, 59},
// 			{Monday, 0, 0, 23, 59},
// 			{Tuesday, 0, 0, 23, 59},
// 			{Wednesday, 0, 0, 23, 59},
// 			{Thursday, 0, 0, 23, 59},
// 			{Friday, 0, 0, 23, 59}}
// 	}
// 	for index := range array_day_start_end {
// 		array_day_start_end[index].DAY = strings.Title(array_day_start_end[index].DAY)

// 		if !IS_IN(array_day_start_end[index].DAY, STANDARD_DAYS) {
// 			array_day_start_end[index].DAY = Sunday
// 		}

// 		if array_day_start_end[index].START_HOUR < 0 {
// 			array_day_start_end[index].START_HOUR = 0
// 		}
// 		if array_day_start_end[index].START_HOUR > 23 {
// 			array_day_start_end[index].START_HOUR = 23
// 		}
// 		if array_day_start_end[index].START_MINUTE < 0 {
// 			array_day_start_end[index].START_MINUTE = 0
// 		}
// 		if array_day_start_end[index].START_MINUTE > 59 {
// 			array_day_start_end[index].START_MINUTE = 59
// 		}
// 		if array_day_start_end[index].END_HOUR < 0 {
// 			array_day_start_end[index].END_HOUR = 0
// 		}
// 		if array_day_start_end[index].END_HOUR > 23 {
// 			array_day_start_end[index].END_HOUR = 23
// 		}
// 		if array_day_start_end[index].END_MINUTE < 0 {
// 			array_day_start_end[index].END_MINUTE = 0
// 		}
// 		if array_day_start_end[index].END_MINUTE > 59 {
// 			array_day_start_end[index].END_MINUTE = 59
// 		}

// 		if array_day_start_end[index].START_HOUR > array_day_start_end[index].END_HOUR {
// 			array_day_start_end[index].START_HOUR = 0
// 		}
// 		if array_day_start_end[index].START_HOUR == array_day_start_end[index].END_HOUR && array_day_start_end[index].START_MINUTE > array_day_start_end[index].END_MINUTE {
// 			array_day_start_end[index].START_MINUTE = 0
// 		}
// 	}
// 	return array_day_start_end
// }

// func ShiftAndArrangeTheTimeSeries(previous_date_end, date_start, date_end time.Time) (time.Time, time.Time) {
// 	if previous_date_end.After(date_start) {
// 		date_start = previous_date_end
// 	}
// 	if date_start.After(date_end) {
// 		date_end = date_start
// 	}
// 	return date_start, date_end
// }

// func TotalMinutes(array_start_end_minutes []start_end_minutes) float64 {
// 	var total_minutes float64
// 	for _, element := range array_start_end_minutes {
// 		total_minutes += element.minutes
// 	}
// 	return total_minutes
// }

// func ValueAfterAdjustUsingAdjustingMethods(adjusting_method string, minutes, TOTAL_MINUTES, time_unit_counter, total_VALUE float64) float64 {
// 	percent := ROOT(total_VALUE, TOTAL_MINUTES)
// 	switch adjusting_method {
// 	case Exponential:
// 		return math.Pow(percent, time_unit_counter+minutes) - math.Pow(percent, time_unit_counter)
// 	case Logarithmic:
// 		return (total_VALUE / math.Pow(percent, time_unit_counter)) - (total_VALUE / math.Pow(percent, time_unit_counter+minutes))
// 	default:
// 		return minutes * (total_VALUE / TOTAL_MINUTES)
// 	}
// }

func CheckDebitEqualCredit(entries []PriceQuantityAccount) ([]PriceQuantityAccount, []PriceQuantityAccount, error) {
	var debitEntries, creditEntries []PriceQuantityAccount
	var zero float64
	for _, entry := range entries {
		VALUE := entry.Price * entry.Quantity
		switch entry.IsCredit {
		case false:
			zero += VALUE
			if VALUE > 0 {
				debitEntries = append(debitEntries, entry)
			} else if VALUE < 0 {
				creditEntries = append(creditEntries, entry)
			}
		case true:
			zero -= VALUE
			if VALUE < 0 {
				debitEntries = append(debitEntries, entry)
			} else if VALUE > 0 {
				creditEntries = append(creditEntries, entry)
			}
		}
	}
	if zero != 0 {
		return debitEntries, creditEntries,
			fmt.Errorf("the debit and credit should be equal. and the debit is more than credit by %f, the debitEntries is %v and the creditEntries is %v", zero, debitEntries, creditEntries)
	}
	if len(debitEntries) != 1 && len(creditEntries) != 1 {
		return debitEntries, creditEntries,
			fmt.Errorf("should be one debit or one credit in the entry, the debitEntries is %v and the creditEntries is %v", debitEntries, creditEntries)
	}
	return debitEntries, creditEntries, nil
}

func SetPriceAndQuantity(account PriceQuantityAccount, isUpdate bool) PriceQuantityAccount {
	if account.Quantity > 0 {
		return account
	}

	// i make it this way just to make it faster when using Wma case
	var keys [][]byte
	var inventory []InventoryTag
	switch account.CostFlowType {
	case Fifo:
		keys, inventory = DbRead[InventoryTag](DbInventory)
	case Lifo:
		keys, inventory = DbRead[InventoryTag](DbInventory)
		ReverseSlice(keys)
		ReverseSlice(inventory)
	case Wma:
		WeightedAverage(account.AccountName)
		keys, inventory = DbRead[InventoryTag](DbInventory)
	}

	QuantityCount := Abs(account.Quantity)
	var costs float64
	for k1, v1 := range inventory {
		if v1.AccountName == account.AccountName {
			if QuantityCount <= v1.Quantity {
				costs -= v1.Price * QuantityCount
				if isUpdate {
					inventory[k1].Quantity -= QuantityCount
					DbUpdate(DbInventory, keys[k1], inventory[k1])
				}
				QuantityCount = 0
				break
			}
			if QuantityCount > v1.Quantity {
				costs -= v1.Price * v1.Quantity
				if isUpdate {
					DbDelete(DbInventory, keys[k1])
				}
				QuantityCount -= v1.Quantity
			}
		}
	}
	account.Quantity += QuantityCount
	account.Price = costs / account.Quantity
	return account
}

func GroupByAccount(entries []PriceQuantityAccount) []PriceQuantityAccount {
	m := map[string]*PriceQuantityAccount{}
	for _, v1 := range entries {
		key := v1.AccountName
		sums := m[key]
		if sums == nil {
			sums = &PriceQuantityAccount{}
			m[key] = sums
		}
		// i make this to store the VALUE and then devide it by the Quantity to get the Price
		sums.IsCredit = v1.IsCredit
		sums.CostFlowType = v1.CostFlowType
		sums.AccountName = v1.AccountName
		sums.Price += v1.Price * v1.Quantity //here i store the VALUE in Price field
		sums.Quantity += v1.Quantity
	}
	entries = []PriceQuantityAccount{}
	for _, v1 := range m {
		entries = append(entries, PriceQuantityAccount{
			IsCredit:     v1.IsCredit,
			CostFlowType: v1.CostFlowType,
			AccountName:  v1.AccountName,
			Price:        v1.Price / v1.Quantity,
			Quantity:     v1.Quantity,
		})
	}
	return entries
}

func InsertEntryNumber(arrayOfJournalTag []JournalTag) {
	journalTag := DbLastLine[JournalTag](DbJournal)
	var lastEntryNumberCompound int
	var entryNumberSimple int
	for k1, v1 := range arrayOfJournalTag {
		arrayOfJournalTag[k1].EntryNumberCompound = journalTag.EntryNumberCompound + 1
		if v1.EntryNumberCompound != lastEntryNumberCompound {
			entryNumberSimple = 0
			lastEntryNumberCompound = v1.EntryNumberCompound
		}
		entryNumberSimple++
		arrayOfJournalTag[k1].EntryNumberSimple = entryNumberSimple
	}
}

func InsertToDatabaseJournal(entries []JournalTag) {
	InsertEntryNumber(entries)
	for _, v1 := range entries {
		DbUpdate(DbJournal, Now(), v1)
	}
}

func InsertToJournalTag(debitEntries, creditEntries []PriceQuantityAccount, notes, name, nameEmployee string) []JournalTag {
	var simpleEntries []JournalTag
	for _, debitEntry := range debitEntries {
		for _, creditEntry := range creditEntries {
			VALUE := Smallest(Abs(debitEntry.Price*debitEntry.Quantity), Abs(creditEntry.Price*creditEntry.Quantity))
			simpleEntries = append(simpleEntries, JournalTag{
				IsReversed:          false,
				EntryNumberCompound: 0,
				EntryNumberSimple:   0,
				Value:               VALUE,
				PriceDebit:          debitEntry.Price,
				PriceCredit:         creditEntry.Price,
				QuantityDebit:       VALUE / debitEntry.Price,
				QuantityCredit:      VALUE / creditEntry.Price,
				AccountDebit:        debitEntry.AccountName,
				AccountCredit:       creditEntry.AccountName,
				Notes:               notes,
				Name:                name,
				NameEmployee:        nameEmployee,
			})
		}
	}
	return simpleEntries
}

func SimpleJournalEntry(
	entries []PriceQuantityAccountBarcode,
	insert, autoCompletion, invoiceDiscount bool,
	notes, name, nameEmployee string) ([]PriceQuantityAccountBarcode, error) {

	sliceOfPriceQuantityAccount := Stage1(entries)
	sliceOfPriceQuantityAccount = GroupByAccount(sliceOfPriceQuantityAccount)

	for k1, v1 := range sliceOfPriceQuantityAccount {
		sliceOfPriceQuantityAccount[k1] = SetPriceAndQuantity(v1, false)
	}

	// if autoCompletion {
	// 	AUTO_COMPLETION_THE_ENTRY(sliceOfPriceQuantityAccount)
	// }
	// if invoiceDiscount {
	// 	entries = auto_completion_the_invoice_discount(entries)
	// }

	sliceOfPriceQuantityAccount = GroupByAccount(sliceOfPriceQuantityAccount)
	debitEntries, creditEntries, err := CheckDebitEqualCredit(sliceOfPriceQuantityAccount)
	newEntries := ConvertPriceQuantityAccountToPriceQuantityAccountBarcode(append(debitEntries, creditEntries...))
	if err != nil {
		return newEntries, err
	}

	if insert {
		simpleEntries := InsertToJournalTag(debitEntries, creditEntries, notes, name, nameEmployee)
		InsertToDatabaseJournal(simpleEntries)
		InsertToDatabaseInventory(sliceOfPriceQuantityAccount)
	}

	return newEntries, nil
}

func ConvertPriceQuantityAccountToPriceQuantityAccountBarcode(entries []PriceQuantityAccount) []PriceQuantityAccountBarcode {
	var newEntries []PriceQuantityAccountBarcode
	for _, v1 := range entries {
		newEntries = append(newEntries, PriceQuantityAccountBarcode{
			Price:       v1.Price,
			Quantity:    v1.Quantity,
			AccountName: v1.AccountName,
			Barcode:     "",
		})
	}
	return newEntries
}

func InsertToDatabaseInventory(entries []PriceQuantityAccount) {
	for _, v1 := range entries {
		if v1.Quantity > 0 {
			DbUpdate(DbInventory, Now(), InventoryTag{v1.Price, v1.Quantity, v1.AccountName})
		} else {
			SetPriceAndQuantity(v1, true)
		}
	}
}

func Stage1(entries []PriceQuantityAccountBarcode) []PriceQuantityAccount {
	var arrayPriceQuantityAccount []PriceQuantityAccount
	for _, v1 := range entries {
		accountStruct, _, err := AccountStructFromBarcode(v1.Barcode)
		if err != nil {
			accountStruct, _, err = AccountStructFromName(FormatTheString(v1.AccountName))
		}
		if err == nil && accountStruct.IsLowLevelAccount && v1.Quantity != 0 && v1.Price != 0 {
			arrayPriceQuantityAccount = append(arrayPriceQuantityAccount, PriceQuantityAccount{
				IsCredit:     accountStruct.IsCredit,
				CostFlowType: accountStruct.CostFlowType,
				AccountName:  accountStruct.AccountName,
				Price:        Abs(v1.Price),
				Quantity:     v1.Quantity,
			})
		}
	}
	return arrayPriceQuantityAccount
}

func ReverseEntries(entryNumberCompound, entryNumberSimple int, nameEmployee string) {
	var entries []JournalTag
	var entriesKeys [][]byte
	keys, journal := DbRead[JournalTag](DbJournal)
	for k1, v1 := range journal {
		if v1.EntryNumberCompound == entryNumberCompound && (entryNumberSimple == 0 || v1.EntryNumberSimple == entryNumberSimple) && v1.IsReversed == false {
			entries = append(entries, v1)
			entriesKeys = append(entriesKeys, keys[k1])
		}
	}

	var entryToReverse []JournalTag
	for k1, v1 := range entries {
		// here i check if the credit side credit nature then it will be negative Quantity and vice versa
		accountStructCredit, _, _ := AccountStructFromName(v1.AccountCredit)
		if accountStructCredit.IsCredit {
			v1.QuantityCredit *= -1
		}
		// here i check if the debit side debit nature then it will be negative Quantity and vice versa
		accountStructDebit, _, _ := AccountStructFromName(v1.AccountDebit)
		if !accountStructDebit.IsCredit {
			v1.QuantityDebit *= -1
		}

		// here i check if the account can be negative by seeing the difference in Quantity after the find the cost in inventory.
		// because i dont want to make the account negative balance
		entryCredit := SetPriceAndQuantity(PriceQuantityAccount{false, Fifo, v1.AccountCredit, v1.PriceCredit, v1.QuantityCredit}, false)
		entryDebit := SetPriceAndQuantity(PriceQuantityAccount{false, Fifo, v1.AccountDebit, v1.PriceDebit, v1.QuantityDebit}, false)

		// here i compare the Quantity if it is the same i will reverse the entry
		if entryCredit.Quantity == v1.QuantityCredit && entryDebit.Quantity == v1.QuantityDebit {

			// here i change the cost flow to Wma just to make outflow from the inventory without error
			entryCredit.CostFlowType = Wma
			entryDebit.CostFlowType = Wma

			// here i insert to the inventory
			InsertToDatabaseInventory([]PriceQuantityAccount{entryCredit, entryDebit})

			// i swap the debit and credit with each other but the Quantity after i swap it will be positive
			v1.PriceCredit, v1.PriceDebit = v1.PriceDebit, v1.PriceCredit
			v1.QuantityCredit, v1.QuantityDebit = Abs(v1.QuantityDebit), Abs(v1.QuantityCredit)
			v1.AccountCredit, v1.AccountDebit = v1.AccountDebit, v1.AccountCredit

			v1.Notes = "revese entry for entry was entered by " + v1.NameEmployee
			v1.NameEmployee = nameEmployee
			v1.IsReverse = true
			v1.ReverseEntryNumberCompound = v1.EntryNumberCompound
			v1.ReverseEntryNumberSimple = v1.EntryNumberSimple

			// here i append the entry to the journal to reverse all in one entry number compound
			entryToReverse = append(entryToReverse, v1)

			// i make the reverse field in the entry true just to not reverse it again
			entries[k1].IsReversed = true
			DbUpdate(DbJournal, entriesKeys[k1], entries[k1])
		}
	}

	// and then i insert to database
	InsertToDatabaseJournal(entryToReverse)
}

func JournalFilter(dates []time.Time, journal []JournalTag, f FilterJournal, isDebitAndCredit bool) ([]time.Time, []JournalTag) {
	var filteredJournal []JournalTag
	var filteredDates []time.Time
	for k1, v1 := range journal {

		// here icheck if the user whant the debit and credit or one of them each time
		var isTheAccounts bool
		if isDebitAndCredit {
			isTheAccounts = f.AccountDebit.Filter(v1.AccountDebit) && f.AccountCredit.Filter(v1.AccountCredit)
		} else {
			isTheAccounts = f.AccountDebit.Filter(v1.AccountDebit) || f.AccountCredit.Filter(v1.AccountCredit)
		}

		if isTheAccounts &&
			f.Date.Filter(dates[k1]) &&
			f.IsReverse.Filter(v1.IsReverse) &&
			f.IsReversed.Filter(v1.IsReversed) &&
			f.ReverseEntryNumberCompound.Filter(float64(v1.ReverseEntryNumberCompound)) &&
			f.ReverseEntryNumberSimple.Filter(float64(v1.ReverseEntryNumberSimple)) &&
			f.EntryNumberCompound.Filter(float64(v1.EntryNumberCompound)) &&
			f.EntryNumberSimple.Filter(float64(v1.EntryNumberSimple)) &&
			f.Value.Filter(v1.Value) &&
			f.PriceDebit.Filter(v1.PriceDebit) &&
			f.PriceCredit.Filter(v1.PriceCredit) &&
			f.QuantityDebit.Filter(v1.QuantityDebit) &&
			f.QuantityCredit.Filter(v1.QuantityCredit) &&
			f.Notes.Filter(v1.Notes) &&
			f.Name.Filter(v1.Name) &&
			f.NameEmployee.Filter(v1.NameEmployee) {
			filteredJournal = append(filteredJournal, v1)
			filteredDates = append(filteredDates, dates[k1])
		}
	}
	return filteredDates, filteredJournal
}

func FindDuplicateElement(dates []time.Time, journal []JournalTag, f TheJournalDuplicateFilter) ([]time.Time, []JournalTag) {
	// here i make it return the same input if the filter is is all false because that mean i dont want to filter
	if !f.IsReverse &&
		!f.IsReversed &&
		!f.ReverseEntryNumberCompound &&
		!f.ReverseEntryNumberSimple &&
		!f.Value &&
		!f.PriceDebit &&
		!f.PriceCredit &&
		!f.QuantityDebit &&
		!f.QuantityCredit &&
		!f.AccountDebit &&
		!f.AccountCredit &&
		!f.Notes &&
		!f.Name &&
		!f.NameEmployee {
		return dates, journal
	}

	var filteredJournal []JournalTag
	var filteredDates []time.Time
	for k1, v1 := range journal {
		for k2, v2 := range journal {
			if k1 != k2 &&
				FunctionFilterDuplicate(v1.IsReverse, v2.IsReverse, f.IsReverse) &&
				FunctionFilterDuplicate(v1.IsReversed, v2.IsReversed, f.IsReversed) &&
				FunctionFilterDuplicate(v1.ReverseEntryNumberCompound, v2.ReverseEntryNumberCompound, f.ReverseEntryNumberCompound) &&
				FunctionFilterDuplicate(v1.ReverseEntryNumberSimple, v2.ReverseEntryNumberSimple, f.ReverseEntryNumberSimple) &&
				FunctionFilterDuplicate(v1.Value, v2.Value, f.Value) &&
				FunctionFilterDuplicate(v1.PriceDebit, v2.PriceDebit, f.PriceDebit) &&
				FunctionFilterDuplicate(v1.PriceCredit, v2.PriceCredit, f.PriceCredit) &&
				FunctionFilterDuplicate(v1.QuantityDebit, v2.QuantityDebit, f.QuantityDebit) &&
				FunctionFilterDuplicate(v1.QuantityCredit, v2.QuantityCredit, f.QuantityCredit) &&
				FunctionFilterDuplicate(v1.AccountDebit, v2.AccountDebit, f.AccountDebit) &&
				FunctionFilterDuplicate(v1.AccountCredit, v2.AccountCredit, f.AccountCredit) &&
				FunctionFilterDuplicate(v1.Notes, v2.Notes, f.Notes) &&
				FunctionFilterDuplicate(v1.Name, v2.Name, f.Name) &&
				FunctionFilterDuplicate(v1.NameEmployee, v2.NameEmployee, f.NameEmployee) {
				filteredJournal = append(filteredJournal, v1)
				filteredDates = append(filteredDates, dates[k1])
				break
			}
		}
	}
	return filteredDates, filteredJournal
}
