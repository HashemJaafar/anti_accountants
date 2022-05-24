package main

import (
	"errors"
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
		VALUE := v1.TPrice * v1.TQuantity
		switch v1.SAccount.TIsCredit {
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
		return debitEntries, creditEntries, errors.New("should be one debit or one credit in the entry")
	}
	return debitEntries, creditEntries, nil
}

func FSetPriceAndQuantity(account SAPQA, insert bool) SAPQA {
	if account.TQuantity > 0 {
		return account
	}

	var keys [][]byte
	var inventory []SAPQ
	switch account.SAccount.TCostFlowType {
	case CFifo:
		keys, inventory = FDbRead[SAPQ](VDbInventory)
	case CLifo:
		keys, inventory = FDbRead[SAPQ](VDbInventory)
		FReverseSlice(keys)
		FReverseSlice(inventory)
	case CWma:
		FWeightedAverage(account.TAccountName)
		keys, inventory = FDbRead[SAPQ](VDbInventory)
	}

	QuantityCount := FAbs(account.TQuantity)
	var costs float64
	for k1, v1 := range inventory {
		if v1.TAccountName == account.TAccountName {
			if QuantityCount <= v1.TQuantity {
				costs -= v1.TPrice * QuantityCount
				if insert {
					inventory[k1].TQuantity -= QuantityCount
					FDbUpdate(VDbInventory, keys[k1], inventory[k1])
				}
				QuantityCount = 0
				break
			}
			if QuantityCount > v1.TQuantity {
				costs -= v1.TPrice * v1.TQuantity
				if insert {
					FDbDelete(VDbInventory, keys[k1])
				}
				QuantityCount -= v1.TQuantity
			}
		}
	}
	account.TQuantity += QuantityCount
	account.TPrice = costs / account.TQuantity
	return account
}

func FGroupByAccount(entries []SAPQA) []SAPQA {
	m := map[string]*SAPQA{}
	for _, v1 := range entries {
		key := v1.TAccountName
		sums := m[key]
		if sums == nil {
			sums = &SAPQA{}
			m[key] = sums
		}
		sums.TAccountName = v1.TAccountName
		sums.TPrice += v1.TPrice * v1.TQuantity //here i store the VALUE in Price field
		sums.TQuantity += v1.TQuantity
		sums.SAccount = v1.SAccount
	}
	entries = []SAPQA{}
	for _, v1 := range m {
		entries = append(entries, SAPQA{
			TAccountName: v1.TAccountName,
			TPrice:       v1.TPrice / v1.TQuantity,
			TQuantity:    v1.TQuantity,
			SAccount:     v1.SAccount,
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

func FInsertToJournal(debitEntries, creditEntries []SAPQA, entryInfo SEntry) []SJournal {
	var simpleEntries []SJournal
	for _, debitEntry := range debitEntries {
		for _, creditEntry := range creditEntries {
			VALUE := FSmallest(FAbs(debitEntry.TPrice*debitEntry.TQuantity), FAbs(creditEntry.TPrice*creditEntry.TQuantity))
			simpleEntries = append(simpleEntries, SJournal{
				IsReverse:                  false,
				IsReversed:                 false,
				ReverseEntryNumberCompound: 0,
				ReverseEntryNumberSimple:   0,
				EntryNumberCompound:        0,
				EntryNumberSimple:          0,
				Value:                      VALUE,
				PriceDebit:                 debitEntry.TPrice,
				PriceCredit:                creditEntry.TPrice,
				QuantityDebit:              VALUE / debitEntry.TPrice,
				QuantityCredit:             VALUE / creditEntry.TPrice,
				AccountDebit:               debitEntry.TAccountName,
				AccountCredit:              creditEntry.TAccountName,
				Notes:                      entryInfo.TEntryNotes,
				Name:                       entryInfo.TPersonName,
				Employee:                   entryInfo.TEmployeeName,
				TypeOfCompoundEntry:        entryInfo.TTypeOfCompoundEntry,
			})
		}
	}
	return simpleEntries
}

func FSimpleJournalEntry(entries []SAPQ, entryInfo SEntry, insert bool) ([]SAPQ, error) {
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

func FInvoiceJournalEntry(payAccountName string, payAccountPrice, invoiceDiscountPrice float64, inventoryAccounts []SAPQ, entryInfo SEntry, insert bool) ([]SAPQ, error) {

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

		autoCompletion, _, err := FFindAutoCompletionFromName(v1.TAccountName)
		if err != nil {
			continue
		}

		var inv []SAPQA
		var tax []SAPQA
		var pay []SAPQA

		quantity := FAbs(v1.TQuantity)
		var valueRevenue float64
		var valueDiscount float64

		inv = append(inv, SAPQA{autoCompletion.TAccountName, v1.TPrice, v1.TQuantity, SAccount{TIsCredit: false}})
		inv = append(inv, SAPQA{CPrefixCost + autoCompletion.TAccountName, v1.TPrice, quantity, SAccount{TIsCredit: false}})

		if autoCompletion.PriceTax > 0 {
			tax = append(tax, SAPQA{CPrefixTaxExpenses + autoCompletion.TAccountName, autoCompletion.PriceTax, quantity, SAccount{TIsCredit: false}})
			tax = append(tax, SAPQA{CPrefixTaxLiability + autoCompletion.TAccountName, autoCompletion.PriceTax, quantity, SAccount{TIsCredit: true}})
		}

		if autoCompletion.PriceRevenue > 0 {
			pay = append(pay, SAPQA{CPrefixRevenue + autoCompletion.TAccountName, autoCompletion.PriceRevenue, quantity, SAccount{TIsCredit: true}})
			valueRevenue = quantity * autoCompletion.PriceRevenue
		}

		if priceDiscount := FMaxDiscount(autoCompletion.PriceDiscount, quantity); priceDiscount > 0 {
			pay = append(pay, SAPQA{CPrefixDiscount + autoCompletion.TAccountName, priceDiscount, quantity, SAccount{TIsCredit: false}})
			valueDiscount = quantity * priceDiscount
		}

		payValue := valueRevenue - valueDiscount
		pay = append(pay, SAPQA{payAccountName, payAccountPrice, payValue / payAccountPrice, SAccount{TIsCredit: false}})

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

func FConvertAPQICToAPQB(entries []SAPQA) []SAPQ {
	var newEntries []SAPQ
	for _, v1 := range entries {
		newEntries = append(newEntries, SAPQ{v1.TAccountName, v1.TPrice, v1.TQuantity})
	}
	return newEntries
}

func FInsertToDatabaseInventory(entries []SAPQA) {
	for _, v1 := range entries {
		if v1.TQuantity > 0 {
			FDbUpdate(VDbInventory, FNow(), SAPQ{v1.TAccountName, v1.TPrice, v1.TQuantity})
		} else {
			FSetPriceAndQuantity(v1, true)
		}
	}
}

func FStage1(entries []SAPQ, isInvoice bool) []SAPQA {
	var newEntries []SAPQA
	for _, v1 := range entries {
		account, _, err := FFindAccountFromNameOrBarcode(v1.TAccountName)
		if isInvoice {
			if account.TIsCredit || v1.TQuantity >= 0 {
				continue
			}
			v1.TPrice = 1
		}
		if err == nil && account.TCostFlowType != "" && v1.TQuantity != 0 && v1.TPrice != 0 {
			newEntries = append(newEntries, SAPQA{
				TAccountName: account.TAccountName,
				TPrice:       FAbs(v1.TPrice),
				TQuantity:    v1.TQuantity,
				SAccount:     account,
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
		v1.QuantityCredit = FConvertTheSignOfDoubleEntryToSingleEntry(account.TIsCredit, true, v1.QuantityCredit)
		account, _, _ = FFindAccountFromName(v1.AccountDebit)
		v1.QuantityDebit = FConvertTheSignOfDoubleEntryToSingleEntry(account.TIsCredit, false, v1.QuantityDebit)

		entryCredit := FSetPriceAndQuantity(SAPQA{v1.AccountCredit, v1.PriceCredit, v1.QuantityCredit, SAccount{TIsCredit: false, TCostFlowType: CFifo}}, false)
		entryDebit := FSetPriceAndQuantity(SAPQA{v1.AccountDebit, v1.PriceDebit, v1.QuantityDebit, SAccount{TIsCredit: false, TCostFlowType: CFifo}}, false)

		if entryCredit.TQuantity == v1.QuantityCredit && entryDebit.TQuantity == v1.QuantityDebit {

			entryCredit.SAccount.TCostFlowType = CWma
			entryDebit.SAccount.TCostFlowType = CWma

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
			TAccountName: v1.AccountDebit,
			TPrice:       v1.PriceDebit,
			TQuantity:    FConvertTheSignOfDoubleEntryToSingleEntry(account.TIsCredit, false, v1.QuantityDebit),
			SAccount:     account,
		})

		account, _, _ = FFindAccountFromName(v1.AccountCredit)
		newEntries = append(newEntries, SAPQA{
			TAccountName: v1.AccountCredit,
			TPrice:       v1.PriceCredit,
			TQuantity:    FConvertTheSignOfDoubleEntryToSingleEntry(account.TIsCredit, true, v1.QuantityCredit),
			SAccount:     account,
		})
	}
	return FGroupByAccount(newEntries)
}

func FConvertAPQAToAPQB(entries []SAPQA) []SAPQ {
	var newEntries []SAPQ
	for _, v1 := range entries {
		newEntries = append(newEntries, SAPQ{
			TAccountName: v1.TAccountName,
			TPrice:       v1.TPrice,
			TQuantity:    v1.TQuantity,
		})
	}
	return newEntries
}

func FExtractEntryInfoFromJournal(entry SJournal) SEntry {
	return SEntry{
		TEntryNotes:          entry.Notes,
		TPersonName:          entry.Name,
		TEmployeeName:        entry.Employee,
		TTypeOfCompoundEntry: entry.TypeOfCompoundEntry,
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
		if v1.TQuantity > quantity {
			price = v1.TPrice
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
