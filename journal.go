package main

import (
	"fmt"
	"math"
	"time"
)

func FValueAfterAdjustUsingAdjustingMethods(adjustingMethod string, minutesCurrent, minutesTotal, minutesPast, valueTotal float64) float64 {
	percent := FRoot(valueTotal, minutesTotal)
	switch adjustingMethod {
	case CExponential:
		return math.Pow(percent, minutesPast+minutesCurrent) - math.Pow(percent, minutesPast)
	case CLogarithmic:
		return (valueTotal / math.Pow(percent, minutesPast)) - (valueTotal / math.Pow(percent, minutesPast+minutesCurrent))
	default:
		return minutesCurrent * (valueTotal / minutesTotal)
	}
}

func FCheckDebitEqualCredit(entries []SAPQA) ([]SAPQA, []SAPQA, error) {
	var debitEntries, creditEntries []SAPQA
	var zero float64
	for _, v1 := range entries {
		VALUE := v1.Price * v1.Quantity
		switch v1.Account.IsCredit {
		case false:
			zero += VALUE
			if VALUE > 0 {
				debitEntries = append(debitEntries, v1)
			} else if VALUE < 0 {
				creditEntries = append(creditEntries, v1)
			}
		case true:
			zero -= VALUE
			if VALUE < 0 {
				debitEntries = append(debitEntries, v1)
			} else if VALUE > 0 {
				creditEntries = append(creditEntries, v1)
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

func FSetPriceAndQuantity(account SAPQA, insert bool) SAPQA {
	if account.Quantity > 0 {
		return account
	}

	var keys [][]byte
	var inventory []SAPQ
	switch account.Account.CostFlowType {
	case CFifo:
		keys, inventory = FDbRead[SAPQ](VDbInventory)
	case CLifo:
		keys, inventory = FDbRead[SAPQ](VDbInventory)
		FReverseSlice(keys)
		FReverseSlice(inventory)
	case CWma:
		FWeightedAverage(account.Name)
		keys, inventory = FDbRead[SAPQ](VDbInventory)
	}

	QuantityCount := FAbs(account.Quantity)
	var costs float64
	for k1, v1 := range inventory {
		if v1.Name == account.Name {
			if QuantityCount <= v1.Quantity {
				costs -= v1.Price * QuantityCount
				if insert {
					inventory[k1].Quantity -= QuantityCount
					FDbUpdate(VDbInventory, keys[k1], inventory[k1])
				}
				QuantityCount = 0
				break
			}
			if QuantityCount > v1.Quantity {
				costs -= v1.Price * v1.Quantity
				if insert {
					FDbDelete(VDbInventory, keys[k1])
				}
				QuantityCount -= v1.Quantity
			}
		}
	}
	account.Quantity += QuantityCount
	account.Price = costs / account.Quantity
	return account
}

func FGroupByAccount(entries []SAPQA) []SAPQA {
	m := map[string]*SAPQA{}
	for _, v1 := range entries {
		key := v1.Name
		sums := m[key]
		if sums == nil {
			sums = &SAPQA{}
			m[key] = sums
		}
		sums.Name = v1.Name
		sums.Price += v1.Price * v1.Quantity //here i store the VALUE in Price field
		sums.Quantity += v1.Quantity
		sums.Account = v1.Account
	}
	entries = []SAPQA{}
	for _, v1 := range m {
		entries = append(entries, SAPQA{
			Name:     v1.Name,
			Price:    v1.Price / v1.Quantity,
			Quantity: v1.Quantity,
			Account:  v1.Account,
		})
	}
	return entries
}

func FInsertToDatabaseJournal(journal []SJournal) {
	last := FDbLastLine[SJournal](VDbJournal)
	for k1, v1 := range journal {
		v1.EntryNumberCompound = last.EntryNumberCompound + 1
		v1.EntryNumberSimple = k1 + 1

		FDbUpdate(VDbJournal, FNow(), v1)
	}
}

func FInsertToJournal(debitEntries, creditEntries []SAPQA, entryInfo SEntryInfo) []SJournal {
	var simpleEntries []SJournal
	for _, debitEntry := range debitEntries {
		for _, creditEntry := range creditEntries {
			VALUE := FSmallest(FAbs(debitEntry.Price*debitEntry.Quantity), FAbs(creditEntry.Price*creditEntry.Quantity))
			simpleEntries = append(simpleEntries, SJournal{
				IsReverse:                  false,
				IsReversed:                 false,
				ReverseEntryNumberCompound: 0,
				ReverseEntryNumberSimple:   0,
				EntryNumberCompound:        0,
				EntryNumberSimple:          0,
				Value:                      VALUE,
				PriceDebit:                 debitEntry.Price,
				PriceCredit:                creditEntry.Price,
				QuantityDebit:              VALUE / debitEntry.Price,
				QuantityCredit:             VALUE / creditEntry.Price,
				AccountDebit:               debitEntry.Name,
				AccountCredit:              creditEntry.Name,
				Notes:                      entryInfo.Notes,
				Name:                       entryInfo.Name,
				Employee:                   entryInfo.Employee,
				TypeOfCompoundEntry:        entryInfo.TypeOfCompoundEntry,
			})
		}
	}
	return simpleEntries
}

func FSimpleJournalEntry(entries []SAPQB, entryInfo SEntryInfo, insert bool) ([]SAPQB, error) {
	newEntries1 := FStage1(entries, false)
	newEntries1 = FGroupByAccount(newEntries1)

	for k1, v1 := range newEntries1 {
		newEntries1[k1] = FSetPriceAndQuantity(v1, false)
	}

	debitEntries, creditEntries, err := FCheckDebitEqualCredit(newEntries1)
	newEntries2 := FConvertAPQICToAPQB(append(debitEntries, creditEntries...))
	if err != nil {
		return newEntries2, err
	}

	if insert {
		simpleEntries := FInsertToJournal(debitEntries, creditEntries, entryInfo)
		FInsertToDatabaseJournal(simpleEntries)
		FInsertToDatabaseInventory(newEntries1)
	}

	return newEntries2, nil
}

func FInvoiceJournalEntry(payAccountName string, payAccountPrice, invoiceDiscountPrice float64, inventoryAccounts []SAPQB, entryInfo SEntryInfo, insert bool) ([]SAPQB, error) {

	if _, isExist, err := FAccountTerms(payAccountName, true, false); !isExist || err != nil {
		return inventoryAccounts, err
	}
	if _, isExist, err := FAccountTerms(VInvoiceDiscount, true, false); !isExist || err != nil {
		return inventoryAccounts, err
	}

	newEntries1 := FStage1(inventoryAccounts, true)
	newEntries1 = FGroupByAccount(newEntries1)

	for k1, v1 := range newEntries1 {
		newEntries1[k1] = FSetPriceAndQuantity(v1, false)
	}

	newEntries2 := FAutoComplete(newEntries1, payAccountName, payAccountPrice)
	inventoryAccounts = FConvertAPQICToAPQB(newEntries1)

	var newEntries3 []SJournal
	newEntries1 = []SAPQA{}
	for _, v1 := range newEntries2 {
		debitEntries, creditEntries, err := FCheckDebitEqualCredit(v1)
		if err != nil {
			return inventoryAccounts, err
		}
		newEntries1 = append(newEntries1, debitEntries...)
		newEntries1 = append(newEntries1, creditEntries...)
		newEntries3 = append(newEntries3, FInsertToJournal(debitEntries, creditEntries, entryInfo)...)
	}

	if insert {
		FInsertToDatabaseJournal(newEntries3)
		FInsertToDatabaseInventory(newEntries1)
	}

	return inventoryAccounts, nil
}

func FAutoComplete(inventoryAccounts []SAPQA, payAccountName string, payAccountPrice float64) [][]SAPQA {
	var simpleInfoAPQ [][]SAPQA
	for _, v1 := range inventoryAccounts {

		autoCompletion, _, err := FFindAutoCompletionFromName(v1.Name)
		if err != nil {
			continue
		}

		var inv []SAPQA
		var tax []SAPQA
		var pay []SAPQA

		quantity := FAbs(v1.Quantity)
		var valueRevenue float64
		var valueDiscount float64

		inv = append(inv, SAPQA{autoCompletion.AccountInvnetory, v1.Price, v1.Quantity, SAccount{IsCredit: false}})
		inv = append(inv, SAPQA{CPrefixCost + autoCompletion.AccountInvnetory, v1.Price, quantity, SAccount{IsCredit: false}})

		if autoCompletion.PriceTax > 0 {
			tax = append(tax, SAPQA{CPrefixTaxExpenses + autoCompletion.AccountInvnetory, autoCompletion.PriceTax, quantity, SAccount{IsCredit: false}})
			tax = append(tax, SAPQA{CPrefixTaxLiability + autoCompletion.AccountInvnetory, autoCompletion.PriceTax, quantity, SAccount{IsCredit: true}})
		}

		if autoCompletion.PriceRevenue > 0 {
			pay = append(pay, SAPQA{CPrefixRevenue + autoCompletion.AccountInvnetory, autoCompletion.PriceRevenue, quantity, SAccount{IsCredit: true}})
			valueRevenue = quantity * autoCompletion.PriceRevenue
		}

		if priceDiscount := FMaxDiscount(autoCompletion.PriceDiscount, quantity); priceDiscount > 0 {
			pay = append(pay, SAPQA{CPrefixDiscount + autoCompletion.AccountInvnetory, priceDiscount, quantity, SAccount{IsCredit: false}})
			valueDiscount = quantity * priceDiscount
		}

		payValue := valueRevenue - valueDiscount
		pay = append(pay, SAPQA{payAccountName, payAccountPrice, payValue / payAccountPrice, SAccount{IsCredit: false}})

		simpleInfoAPQ = append(simpleInfoAPQ, inv)
		if len(tax) > 1 {
			simpleInfoAPQ = append(simpleInfoAPQ, tax)
		}
		if len(pay) > 1 {
			simpleInfoAPQ = append(simpleInfoAPQ, pay)
		}

	}
	return simpleInfoAPQ
}

func FConvertAPQICToAPQB(entries []SAPQA) []SAPQB {
	var newEntries []SAPQB
	for _, v1 := range entries {
		newEntries = append(newEntries, SAPQB{v1.Name, v1.Price, v1.Quantity, ""})
	}
	return newEntries
}

func FInsertToDatabaseInventory(entries []SAPQA) {
	for _, v1 := range entries {
		if v1.Quantity > 0 {
			FDbUpdate(VDbInventory, FNow(), SAPQ{v1.Name, v1.Price, v1.Quantity})
		} else {
			FSetPriceAndQuantity(v1, true)
		}
	}
}

func FStage1(entries []SAPQB, isInvoice bool) []SAPQA {
	var newEntries []SAPQA
	for _, v1 := range entries {
		account, _, err := FFindAccountFromBarcode(v1.Barcode)
		if err != nil {
			account, _, err = FFindAccountFromName(FFormatTheString(v1.Name))
		}
		if isInvoice {
			if account.IsCredit || v1.Quantity >= 0 {
				continue
			}
			v1.Price = 1
		}
		if err == nil && account.IsLowLevel && v1.Quantity != 0 && v1.Price != 0 {
			newEntries = append(newEntries, SAPQA{
				Name:     account.Name,
				Price:    FAbs(v1.Price),
				Quantity: v1.Quantity,
				Account:  account,
			})
		}
	}
	return newEntries
}

func FReverseEntries(entriesKeys [][]byte, entries []SJournal, nameEmployee string) {
	var entryToReverse []SJournal
	for k1, v1 := range entries {
		if v1.IsReversed {
			continue
		}
		account, _, _ := FFindAccountFromName(v1.AccountCredit)
		v1.QuantityCredit = FConvertTheSignOfDoubleEntryToSingleEntry(account.IsCredit, true, v1.QuantityCredit)
		account, _, _ = FFindAccountFromName(v1.AccountDebit)
		v1.QuantityDebit = FConvertTheSignOfDoubleEntryToSingleEntry(account.IsCredit, false, v1.QuantityDebit)

		entryCredit := FSetPriceAndQuantity(SAPQA{v1.AccountCredit, v1.PriceCredit, v1.QuantityCredit, SAccount{IsCredit: false, CostFlowType: CFifo}}, false)
		entryDebit := FSetPriceAndQuantity(SAPQA{v1.AccountDebit, v1.PriceDebit, v1.QuantityDebit, SAccount{IsCredit: false, CostFlowType: CFifo}}, false)

		if entryCredit.Quantity == v1.QuantityCredit && entryDebit.Quantity == v1.QuantityDebit {

			entryCredit.Account.CostFlowType = CWma
			entryDebit.Account.CostFlowType = CWma

			FInsertToDatabaseInventory([]SAPQA{entryCredit, entryDebit})

			v1.PriceCredit, v1.PriceDebit = v1.PriceDebit, v1.PriceCredit
			v1.QuantityCredit, v1.QuantityDebit = FAbs(v1.QuantityDebit), FAbs(v1.QuantityCredit)
			v1.AccountCredit, v1.AccountDebit = v1.AccountDebit, v1.AccountCredit

			v1.IsReverse = true
			v1.ReverseEntryNumberCompound = v1.EntryNumberCompound
			v1.ReverseEntryNumberSimple = v1.EntryNumberSimple
			v1.Notes = "revese entry for entry was entered by " + v1.Employee
			v1.Employee = nameEmployee

			entryToReverse = append(entryToReverse, v1)

			entries[k1].IsReversed = true
			FDbUpdate(VDbJournal, entriesKeys[k1], entries[k1])
		}
	}

	FInsertToDatabaseJournal(entryToReverse)
}

func FFindEntryFromNumber(entryNumberCompound int, entryNumberSimple int) ([][]byte, []SJournal) {
	var entries []SJournal
	var entriesKeys [][]byte
	keys, journal := FDbRead[SJournal](VDbJournal)
	for k1, v1 := range journal {
		if v1.EntryNumberCompound == entryNumberCompound && (entryNumberSimple == 0 || v1.EntryNumberSimple == entryNumberSimple) {
			entries = append(entries, v1)
			entriesKeys = append(entriesKeys, keys[k1])
		}
	}
	return entriesKeys, entries
}

func FConvertJournalToAPQA(entries []SJournal) []SAPQA {
	var newEntries []SAPQA
	for _, v1 := range entries {
		account, _, _ := FFindAccountFromName(v1.AccountDebit)
		newEntries = append(newEntries, SAPQA{
			Name:     v1.AccountDebit,
			Price:    v1.PriceDebit,
			Quantity: FConvertTheSignOfDoubleEntryToSingleEntry(account.IsCredit, false, v1.QuantityDebit),
			Account:  account,
		})

		account, _, _ = FFindAccountFromName(v1.AccountCredit)
		newEntries = append(newEntries, SAPQA{
			Name:     v1.AccountCredit,
			Price:    v1.PriceCredit,
			Quantity: FConvertTheSignOfDoubleEntryToSingleEntry(account.IsCredit, true, v1.QuantityCredit),
			Account:  account,
		})
	}
	return FGroupByAccount(newEntries)
}

func FConvertAPQAToAPQB(entries []SAPQA) []SAPQB {
	var newEntries []SAPQB
	for _, v1 := range entries {
		newEntries = append(newEntries, SAPQB{
			Name:     v1.Name,
			Price:    v1.Price,
			Quantity: v1.Quantity,
			Barcode:  v1.Account.Barcode[0],
		})
	}
	return newEntries
}

func FExtractEntryInfoFromJournal(entry SJournal) SEntryInfo {
	return SEntryInfo{
		Notes:               entry.Notes,
		Name:                entry.Name,
		Employee:            entry.Employee,
		TypeOfCompoundEntry: entry.TypeOfCompoundEntry,
	}
}

func FJournalFilter(dates []time.Time, journal []SJournal, f SFilterJournal, isDebitAndCredit bool) ([]time.Time, []SJournal) {
	if !f.Date.IsFilter &&
		!f.IsReverse.IsFilter &&
		!f.IsReversed.IsFilter &&
		!f.ReverseEntryNumberCompound.IsFilter &&
		!f.ReverseEntryNumberSimple.IsFilter &&
		!f.EntryNumberCompound.IsFilter &&
		!f.EntryNumberSimple.IsFilter &&
		!f.Value.IsFilter &&
		!f.PriceDebit.IsFilter &&
		!f.PriceCredit.IsFilter &&
		!f.QuantityDebit.IsFilter &&
		!f.QuantityCredit.IsFilter &&
		!f.AccountDebit.IsFilter &&
		!f.AccountCredit.IsFilter &&
		!f.Notes.IsFilter &&
		!f.Name.IsFilter &&
		!f.Employee.IsFilter &&
		!f.TypeOfCompoundEntry.IsFilter {
		return dates, journal
	}

	var newJournal []SJournal
	var newDates []time.Time
	for k1, v1 := range journal {

		var isTheAccounts bool
		if isDebitAndCredit {
			isTheAccounts = f.AccountDebit.FFilter(v1.AccountDebit) && f.AccountCredit.FFilter(v1.AccountCredit)
		} else {
			isTheAccounts = f.AccountDebit.FFilter(v1.AccountDebit) || f.AccountCredit.FFilter(v1.AccountCredit)
		}

		if isTheAccounts &&
			f.Date.FFilter(dates[k1]) &&
			f.IsReverse.FFilter(v1.IsReverse) &&
			f.IsReversed.FFilter(v1.IsReversed) &&
			f.ReverseEntryNumberCompound.FFilter(float64(v1.ReverseEntryNumberCompound)) &&
			f.ReverseEntryNumberSimple.FFilter(float64(v1.ReverseEntryNumberSimple)) &&
			f.EntryNumberCompound.FFilter(float64(v1.EntryNumberCompound)) &&
			f.EntryNumberSimple.FFilter(float64(v1.EntryNumberSimple)) &&
			f.Value.FFilter(v1.Value) &&
			f.PriceDebit.FFilter(v1.PriceDebit) &&
			f.PriceCredit.FFilter(v1.PriceCredit) &&
			f.QuantityDebit.FFilter(v1.QuantityDebit) &&
			f.QuantityCredit.FFilter(v1.QuantityCredit) &&
			f.Notes.FFilter(v1.Notes) &&
			f.Name.FFilter(v1.Name) &&
			f.Employee.FFilter(v1.Employee) &&
			f.TypeOfCompoundEntry.FFilter(v1.TypeOfCompoundEntry) {
			newJournal = append(newJournal, v1)
			newDates = append(newDates, dates[k1])
		}
	}
	return newDates, newJournal
}

func FFindDuplicateElement(dates []time.Time, journal []SJournal, f SFilterJournalDuplicate) ([]time.Time, []SJournal) {
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
		!f.Employee &&
		!f.TypeOfCompoundEntry {
		return dates, journal
	}

	var newJournal []SJournal
	var newDates []time.Time
	for k1, v1 := range journal {
		for k2, v2 := range journal {
			if k1 != k2 &&
				FFilterDuplicate(v1.IsReverse, v2.IsReverse, f.IsReverse) &&
				FFilterDuplicate(v1.IsReversed, v2.IsReversed, f.IsReversed) &&
				FFilterDuplicate(v1.ReverseEntryNumberCompound, v2.ReverseEntryNumberCompound, f.ReverseEntryNumberCompound) &&
				FFilterDuplicate(v1.ReverseEntryNumberSimple, v2.ReverseEntryNumberSimple, f.ReverseEntryNumberSimple) &&
				FFilterDuplicate(v1.Value, v2.Value, f.Value) &&
				FFilterDuplicate(v1.PriceDebit, v2.PriceDebit, f.PriceDebit) &&
				FFilterDuplicate(v1.PriceCredit, v2.PriceCredit, f.PriceCredit) &&
				FFilterDuplicate(v1.QuantityDebit, v2.QuantityDebit, f.QuantityDebit) &&
				FFilterDuplicate(v1.QuantityCredit, v2.QuantityCredit, f.QuantityCredit) &&
				FFilterDuplicate(v1.AccountDebit, v2.AccountDebit, f.AccountDebit) &&
				FFilterDuplicate(v1.AccountCredit, v2.AccountCredit, f.AccountCredit) &&
				FFilterDuplicate(v1.Notes, v2.Notes, f.Notes) &&
				FFilterDuplicate(v1.Name, v2.Name, f.Name) &&
				FFilterDuplicate(v1.Employee, v2.Employee, f.Employee) &&
				FFilterDuplicate(v1.TypeOfCompoundEntry, v2.TypeOfCompoundEntry, f.TypeOfCompoundEntry) {
				newJournal = append(newJournal, v1)
				newDates = append(newDates, dates[k1])
				break
			}
		}
	}
	return newDates, newJournal
}

func FMaxDiscount(discounts []SPQ, quantity float64) float64 {
	var price float64
	for _, v1 := range discounts {
		if v1.Quantity > quantity {
			price = v1.Price
		}
	}
	return FAbs(price)
}

func FConvertTheSignOfDoubleEntryToSingleEntry(isCredit, isCreditInTheEntry bool, number float64) float64 {
	number = math.Abs(number)
	if isCredit != isCreditInTheEntry {
		return -number
	}
	return number
}
