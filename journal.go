package main

import (
	"errors"
	"fmt"
	"math"
	"sync"
	"time"
)

//TODO
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

func FCheckDebitEqualCredit(entries []SAPQAE) ([]SAPQAE, []SAPQAE, error) {
	var debitEntries, creditEntries []SAPQAE
	var zero float64
	for _, v1 := range entries {
		VALUE := v1.TPrice * v1.TQuantity
		switch v1.SAccount1.IsCredit {
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
			fmt.Errorf("the debit and credit should be equal. and the debit is more than credit by %v, the debitEntries is %v and the creditEntries is %v", zero, debitEntries, creditEntries)
	}
	if len(debitEntries) != 1 && len(creditEntries) != 1 {
		return debitEntries, creditEntries, errors.New("should be one debit or one credit in the entry")
	}
	return debitEntries, creditEntries, nil
}

func FSetPriceAndQuantity(account SAPQAE, insert bool) SAPQAE {
	if account.TQuantity > 0 {
		return account
	}

	var keys [][]byte
	var inventory []SAPQ
	switch account.SAccount1.CostFlowType {
	case CFifo:
		keys, inventory = FDbRead[SAPQ](VDbInventory)
	case CLifo:
		keys, inventory = FDbRead[SAPQ](VDbInventory)
		FReverseSlice(keys)
		FReverseSlice(inventory)
	case CWma:
		FWeightedAverage(account.TAccountName)
		keys, inventory = FDbRead[SAPQ](VDbInventory)
	default:
		return account
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
	if QuantityCount > 0 {
		account.error = fmt.Errorf("the quantity of %s is too much", account.TAccountName)
	}
	account.TQuantity += QuantityCount
	account.TPrice = costs / account.TQuantity
	return account
}

func FGroupByAccount(entries []SAPQAE) []SAPQAE {
	m := map[string]*SAPQAE{}
	for _, v1 := range entries {
		key := v1.TAccountName
		sums := m[key]
		if sums == nil {
			sums = &SAPQAE{}
			m[key] = sums
		}
		sums.TAccountName = v1.TAccountName
		sums.TPrice += v1.TPrice * v1.TQuantity //here i store the VALUE in Price field
		sums.TQuantity += v1.TQuantity
		sums.SAccount1 = v1.SAccount1
		sums.error = v1.error
	}
	entries = []SAPQAE{}
	for _, v1 := range m {
		entries = append(entries, SAPQAE{
			TAccountName: v1.TAccountName,
			TPrice:       v1.TPrice / v1.TQuantity,
			TQuantity:    v1.TQuantity,
			SAccount1:    v1.SAccount1,
			error:        v1.error,
		})
	}
	return entries
}

func FInsertToDatabaseJournal(journal []SJournal1) {
	last := FDbLastLine[SJournal1](VDbJournal)
	for k1, v1 := range journal {
		v1.EntryNumberCompound = last.EntryNumberCompound + 1
		v1.EntryNumberSimple = uint(k1) + 1

		v1.DebitBalanceValue, v1.DebitBalancePrice, v1.DebitBalanceQuantity = FTotalValuePriceQuantity(v1.DebitAccountName)
		v1.CreditBalanceValue, v1.CreditBalancePrice, v1.CreditBalanceQuantity = FTotalValuePriceQuantity(v1.CreditAccountName)

		FDbUpdate(VDbJournal, FNow(), v1)
	}
}

func FInsertToJournal(debitEntries, creditEntries []SAPQAE, entryInfo SEntry) []SJournal1 {
	var simpleEntries []SJournal1
	for _, debitEntry := range debitEntries {
		for _, creditEntry := range creditEntries {
			value := FSmallest(FAbs(debitEntry.TPrice*debitEntry.TQuantity), FAbs(creditEntry.TPrice*creditEntry.TQuantity))
			simpleEntries = append(simpleEntries, SJournal1{
				Date:                       time.Now(),
				IsReverse:                  false,
				IsReversed:                 false,
				ReverseEntryNumberCompound: 0,
				ReverseEntryNumberSimple:   0,
				EntryNumberCompound:        0,
				EntryNumberSimple:          0,
				Value:                      value,
				DebitAccountName:           debitEntry.TAccountName,
				DebitPrice:                 debitEntry.TPrice,
				DebitQuantity:              value / debitEntry.TPrice,
				DebitBalanceValue:          0,
				DebitBalancePrice:          0,
				DebitBalanceQuantity:       0,
				CreditAccountName:          creditEntry.TAccountName,
				CreditPrice:                creditEntry.TPrice,
				CreditQuantity:             value / creditEntry.TPrice,
				CreditBalanceValue:         0,
				CreditBalancePrice:         0,
				CreditBalanceQuantity:      0,
				Notes:                      entryInfo.Notes,
				Name:                       entryInfo.Name,
				Employee:                   entryInfo.Employee,
				TypeOfCompoundEntry:        entryInfo.TypeOfCompoundEntry,
			})
		}
	}
	return simpleEntries
}

func FSimpleJournalEntry(entries []SAPQ, entryInfo SEntry, insert bool) ([]SAPQAE, error) {
	newEntries1 := FSetEntries(entries, false)
	newEntries1 = FGroupByAccount(newEntries1)

	for k1, v1 := range newEntries1 {
		if v1.error == nil {
			newEntries1[k1] = FSetPriceAndQuantity(v1, false)
		}
	}

	debitEntries, creditEntries, err := FCheckDebitEqualCredit(newEntries1)
	if err != nil {
		return append(debitEntries, creditEntries...), err
	}

	if insert {
		simpleEntries := FInsertToJournal(debitEntries, creditEntries, entryInfo)
		FInsertToDatabase(simpleEntries, newEntries1)
	}

	return append(debitEntries, creditEntries...), nil
}

func FInvoiceJournalEntry(payAccountName string, payAccountPrice, invoiceDiscountPrice float64, inventoryAccounts []SAPQ, entryInfo SEntry, insert bool) ([]SAPQ, error) {

	if _, isExist, err := FAccountTerms(payAccountName, true, false); !isExist || err != nil {
		return inventoryAccounts, err
	}
	if _, isExist, err := FAccountTerms(VInvoiceDiscount, true, false); !isExist || err != nil {
		return inventoryAccounts, err
	}

	newEntries1 := FSetEntries(inventoryAccounts, true)
	newEntries1 = FGroupByAccount(newEntries1)

	for k1, v1 := range newEntries1 {
		newEntries1[k1] = FSetPriceAndQuantity(v1, false)
	}

	newEntries2 := FAutoComplete(newEntries1, payAccountName, payAccountPrice)
	inventoryAccounts = FConvertAPQICToAPQB(newEntries1)

	var newEntries3 []SJournal1
	newEntries1 = []SAPQAE{}
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
		FInsertToDatabase(newEntries3, newEntries1)
	}

	return inventoryAccounts, nil
}

func FInsertToDatabase(entriesJournal []SJournal1, entriesInventory []SAPQAE) {
	var wait sync.WaitGroup
	wait.Add(2)

	go func() {
		FInsertToDatabaseJournal(entriesJournal)
		wait.Done()
	}()
	go func() {
		FInsertToDatabaseInventory(entriesInventory)
		wait.Done()
	}()

	wait.Wait()
}

func FAutoComplete(inventoryAccounts []SAPQAE, payAccountName string, payAccountPrice float64) [][]SAPQAE {
	var simpleInfoAPQ [][]SAPQAE
	for _, v1 := range inventoryAccounts {

		autoCompletion, _, err := FFindAutoCompletionFromName(v1.TAccountName)
		if err != nil {
			continue
		}

		var inv []SAPQAE
		var tax []SAPQAE
		var pay []SAPQAE

		quantity := FAbs(v1.TQuantity)
		var valueRevenue float64
		var valueDiscount float64

		inv = append(inv, SAPQAE{autoCompletion.TAccountName, v1.TPrice, v1.TQuantity, SAccount1{IsCredit: false}, nil})
		inv = append(inv, SAPQAE{CPrefixCost + autoCompletion.TAccountName, v1.TPrice, quantity, SAccount1{IsCredit: false}, nil})

		if autoCompletion.PriceTax > 0 {
			tax = append(tax, SAPQAE{CPrefixTaxExpenses + autoCompletion.TAccountName, autoCompletion.PriceTax, quantity, SAccount1{IsCredit: false}, nil})
			tax = append(tax, SAPQAE{CPrefixTaxLiability + autoCompletion.TAccountName, autoCompletion.PriceTax, quantity, SAccount1{IsCredit: true}, nil})
		}

		if autoCompletion.PriceRevenue > 0 {
			pay = append(pay, SAPQAE{CPrefixRevenue + autoCompletion.TAccountName, autoCompletion.PriceRevenue, quantity, SAccount1{IsCredit: true}, nil})
			valueRevenue = quantity * autoCompletion.PriceRevenue
		}

		var priceDiscount float64
		for _, v1 := range autoCompletion.PriceDiscount {
			if v1.TQuantity > quantity {
				priceDiscount = v1.TPrice
			}
		}
		priceDiscount = FAbs(priceDiscount)

		if priceDiscount > 0 {
			pay = append(pay, SAPQAE{CPrefixDiscount + autoCompletion.TAccountName, priceDiscount, quantity, SAccount1{IsCredit: false}, nil})
			valueDiscount = quantity * priceDiscount
		}

		payValue := valueRevenue - valueDiscount
		pay = append(pay, SAPQAE{payAccountName, payAccountPrice, payValue / payAccountPrice, SAccount1{IsCredit: false}, nil})

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

func FConvertAPQICToAPQB(entries []SAPQAE) []SAPQ {
	var newEntries []SAPQ
	for _, v1 := range entries {
		newEntries = append(newEntries, SAPQ{v1.TAccountName, v1.TPrice, v1.TQuantity})
	}
	return newEntries
}

func FInsertToDatabaseInventory(entries []SAPQAE) {
	for _, v1 := range entries {
		if v1.TQuantity > 0 {
			FDbUpdate(VDbInventory, FNow(), SAPQ{v1.TAccountName, v1.TPrice, v1.TQuantity})
		} else {
			FSetPriceAndQuantity(v1, true)
		}
	}
}

func FSetEntries(entries []SAPQ, isInvoice bool) []SAPQAE {
	var newEntries []SAPQAE
	for _, v1 := range entries {
		account, _, err1 := FFindAccountFromNameOrBarcode(v1.TAccountName)
		if isInvoice {
			if account.IsCredit || v1.TQuantity >= 0 {
				continue
			}
			v1.TPrice = 1
		}

		var err2 error
		switch {
		case err1 != nil:
			err2 = err1
		case account.CostFlowType == CHighLevelAccount:
			err2 = errors.New("you can't use high level account in journal")
		case v1.TQuantity == 0:
			err2 = errors.New("quantity can't be 0")
		case v1.TPrice == 0:
			err2 = errors.New("price can't be 0")
		}

		newEntries = append(newEntries, SAPQAE{
			TAccountName: v1.TAccountName,
			TPrice:       FAbs(v1.TPrice),
			TQuantity:    v1.TQuantity,
			SAccount1:    account,
			error:        err2,
		})
	}
	return newEntries
}

func FReverseEntries(entriesKeys [][]byte, entries []SJournal1, nameEmployee string) {
	var entryToReverse []SJournal1
	for k1, v1 := range entries {
		if v1.IsReversed {
			continue
		}
		account, _, _ := FFindAccountFromName(v1.CreditAccountName)
		v1.CreditQuantity = FConvertTheSignOfDoubleEntryToSingleEntry(account.IsCredit, true, v1.CreditQuantity)
		account, _, _ = FFindAccountFromName(v1.DebitAccountName)
		v1.DebitQuantity = FConvertTheSignOfDoubleEntryToSingleEntry(account.IsCredit, false, v1.DebitQuantity)

		entryCredit := FSetPriceAndQuantity(SAPQAE{v1.CreditAccountName, v1.CreditPrice, v1.CreditQuantity, SAccount1{IsCredit: false, CostFlowType: CFifo}, nil}, false)
		entryDebit := FSetPriceAndQuantity(SAPQAE{v1.DebitAccountName, v1.DebitPrice, v1.DebitQuantity, SAccount1{IsCredit: false, CostFlowType: CFifo}, nil}, false)

		if entryCredit.TQuantity == v1.CreditQuantity && entryDebit.TQuantity == v1.DebitQuantity {

			entryCredit.SAccount1.CostFlowType = CWma
			entryDebit.SAccount1.CostFlowType = CWma

			FInsertToDatabaseInventory([]SAPQAE{entryCredit, entryDebit})

			v1.CreditPrice, v1.DebitPrice = v1.DebitPrice, v1.CreditPrice
			v1.CreditQuantity, v1.DebitQuantity = FAbs(v1.DebitQuantity), FAbs(v1.CreditQuantity)
			v1.CreditAccountName, v1.DebitAccountName = v1.DebitAccountName, v1.CreditAccountName

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

func FConvertJournalToAPQA(entries []SJournal1) []SAPQAE {
	var newEntries []SAPQAE
	for _, v1 := range entries {
		account, _, _ := FFindAccountFromName(v1.DebitAccountName)
		newEntries = append(newEntries, SAPQAE{
			TAccountName: v1.DebitAccountName,
			TPrice:       v1.DebitPrice,
			TQuantity:    FConvertTheSignOfDoubleEntryToSingleEntry(account.IsCredit, false, v1.DebitQuantity),
			SAccount1:    account,
		})

		account, _, _ = FFindAccountFromName(v1.CreditAccountName)
		newEntries = append(newEntries, SAPQAE{
			TAccountName: v1.CreditAccountName,
			TPrice:       v1.CreditPrice,
			TQuantity:    FConvertTheSignOfDoubleEntryToSingleEntry(account.IsCredit, true, v1.CreditQuantity),
			SAccount1:    account,
		})
	}
	return FGroupByAccount(newEntries)
}

func FConvertAPQAToAPQB(entries []SAPQAE) []SAPQ {
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

func FJournalFilter(journal []SJournal1, f SJournal2, isDebitAndCredit bool) []SJournal1 {
	if f.Date.Way == CDontFilter &&
		!f.IsReverse.IsFilter &&
		!f.IsReversed.IsFilter &&
		f.ReverseEntryNumberCompound.Way == CDontFilter &&
		f.ReverseEntryNumberSimple.Way == CDontFilter &&
		f.EntryNumberCompound.Way == CDontFilter &&
		f.EntryNumberSimple.Way == CDontFilter &&
		f.Value.Way == CDontFilter &&
		f.DebitAccountName.Way == CDontFilter &&
		f.DebitPrice.Way == CDontFilter &&
		f.DebitQuantity.Way == CDontFilter &&
		f.DebitBalanceValue.Way == CDontFilter &&
		f.DebitBalancePrice.Way == CDontFilter &&
		f.DebitBalanceQuantity.Way == CDontFilter &&
		f.CreditAccountName.Way == CDontFilter &&
		f.CreditPrice.Way == CDontFilter &&
		f.CreditQuantity.Way == CDontFilter &&
		f.CreditBalanceValue.Way == CDontFilter &&
		f.CreditBalancePrice.Way == CDontFilter &&
		f.CreditBalanceQuantity.Way == CDontFilter &&
		f.Notes.Way == CDontFilter &&
		f.Name.Way == CDontFilter &&
		f.Employee.Way == CDontFilter &&
		f.TypeOfCompoundEntry.Way == CDontFilter {
		return journal
	}

	var newJournal []SJournal1
	for _, v1 := range journal {

		var isTheAccounts bool
		if isDebitAndCredit {
			isTheAccounts = FFilterString(v1.DebitAccountName, f.DebitAccountName) && FFilterString(v1.CreditAccountName, f.CreditAccountName)
		} else {
			isTheAccounts = FFilterString(v1.DebitAccountName, f.DebitAccountName) || FFilterString(v1.CreditAccountName, f.CreditAccountName)
		}

		if isTheAccounts &&
			FFilterTime(v1.Date, f.Date) &&
			FFilterBool(v1.IsReverse, f.IsReverse) &&
			FFilterBool(v1.IsReversed, f.IsReversed) &&
			FFilterNumber(v1.ReverseEntryNumberCompound, f.ReverseEntryNumberCompound) &&
			FFilterNumber(v1.ReverseEntryNumberSimple, f.ReverseEntryNumberSimple) &&
			FFilterNumber(v1.EntryNumberCompound, f.EntryNumberCompound) &&
			FFilterNumber(v1.EntryNumberSimple, f.EntryNumberSimple) &&
			FFilterNumber(v1.Value, f.Value) &&
			FFilterString(v1.DebitAccountName, f.DebitAccountName) &&
			FFilterNumber(v1.DebitPrice, f.DebitPrice) &&
			FFilterNumber(v1.DebitQuantity, f.DebitQuantity) &&
			FFilterNumber(v1.DebitBalanceValue, f.DebitBalanceValue) &&
			FFilterNumber(v1.DebitBalancePrice, f.DebitBalancePrice) &&
			FFilterNumber(v1.DebitBalanceQuantity, f.DebitBalanceQuantity) &&
			FFilterString(v1.CreditAccountName, f.CreditAccountName) &&
			FFilterNumber(v1.CreditPrice, f.CreditPrice) &&
			FFilterNumber(v1.CreditQuantity, f.CreditQuantity) &&
			FFilterNumber(v1.CreditBalanceValue, f.CreditBalanceValue) &&
			FFilterNumber(v1.CreditBalancePrice, f.CreditBalancePrice) &&
			FFilterNumber(v1.CreditBalanceQuantity, f.CreditBalanceQuantity) &&
			FFilterString(v1.Notes, f.Notes) &&
			FFilterString(v1.Name, f.Name) &&
			FFilterString(v1.Employee, f.Employee) &&
			FFilterString(v1.TypeOfCompoundEntry, f.TypeOfCompoundEntry) {
			newJournal = append(newJournal, v1)
		}
	}
	return newJournal
}

func FFindDuplicateElement(dates []time.Time, journal []SJournal1, f SJournal3) ([]time.Time, []SJournal1) {
	if !f.IsReverse &&
		!f.IsReversed &&
		!f.ReverseEntryNumberCompound &&
		!f.ReverseEntryNumberSimple &&
		!f.EntryNumberCompound &&
		!f.EntryNumberSimple &&
		!f.Value &&
		!f.DebitAccountName &&
		!f.DebitPrice &&
		!f.DebitQuantity &&
		!f.DebitBalanceValue &&
		!f.DebitBalancePrice &&
		!f.DebitBalanceQuantity &&
		!f.CreditAccountName &&
		!f.CreditPrice &&
		!f.CreditQuantity &&
		!f.CreditBalanceValue &&
		!f.CreditBalancePrice &&
		!f.CreditBalanceQuantity &&
		!f.Notes &&
		!f.Name &&
		!f.Employee &&
		!f.TypeOfCompoundEntry {
		return dates, journal
	}

	var newJournal []SJournal1
	var newDates []time.Time
	for k1, v1 := range journal {
		for k2, v2 := range journal {
			if k1 != k2 &&
				FFilterDuplicate(v1.IsReverse, v2.IsReverse, f.IsReverse) &&
				FFilterDuplicate(v1.IsReversed, v2.IsReversed, f.IsReversed) &&
				FFilterDuplicate(v1.ReverseEntryNumberCompound, v2.ReverseEntryNumberCompound, f.ReverseEntryNumberCompound) &&
				FFilterDuplicate(v1.ReverseEntryNumberSimple, v2.ReverseEntryNumberSimple, f.ReverseEntryNumberSimple) &&
				FFilterDuplicate(v1.EntryNumberCompound, v2.EntryNumberCompound, f.EntryNumberCompound) &&
				FFilterDuplicate(v1.EntryNumberSimple, v2.EntryNumberSimple, f.EntryNumberSimple) &&
				FFilterDuplicate(v1.Value, v2.Value, f.Value) &&
				FFilterDuplicate(v1.DebitAccountName, v2.DebitAccountName, f.DebitAccountName) &&
				FFilterDuplicate(v1.DebitPrice, v2.DebitPrice, f.DebitPrice) &&
				FFilterDuplicate(v1.DebitQuantity, v2.DebitQuantity, f.DebitQuantity) &&
				FFilterDuplicate(v1.DebitBalanceValue, v2.DebitBalanceValue, f.DebitBalanceValue) &&
				FFilterDuplicate(v1.DebitBalancePrice, v2.DebitBalancePrice, f.DebitBalancePrice) &&
				FFilterDuplicate(v1.DebitBalanceQuantity, v2.DebitBalanceQuantity, f.DebitBalanceQuantity) &&
				FFilterDuplicate(v1.CreditAccountName, v2.CreditAccountName, f.CreditAccountName) &&
				FFilterDuplicate(v1.CreditPrice, v2.CreditPrice, f.CreditPrice) &&
				FFilterDuplicate(v1.CreditQuantity, v2.CreditQuantity, f.CreditQuantity) &&
				FFilterDuplicate(v1.CreditBalanceValue, v2.CreditBalanceValue, f.CreditBalanceValue) &&
				FFilterDuplicate(v1.CreditBalancePrice, v2.CreditBalancePrice, f.CreditBalancePrice) &&
				FFilterDuplicate(v1.CreditBalanceQuantity, v2.CreditBalanceQuantity, f.CreditBalanceQuantity) &&
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

func FConvertTheSignOfDoubleEntryToSingleEntry(isCredit, isCreditInTheEntry bool, number float64) float64 {
	number = math.Abs(number)
	if isCredit != isCreditInTheEntry {
		return -number
	}
	return number
}

func FTotalValuePriceQuantity(account string) (float64, float64, float64) {
	var totalValue, totalQuantity float64
	_, journal := FDbRead[SJournal1](VDbJournal)
	for _, v1 := range journal {
		if v1.CreditAccountName == account {
			totalValue += v1.Value
			totalQuantity += v1.CreditQuantity
		}
		if v1.DebitAccountName == account {
			totalValue -= v1.Value
			totalQuantity -= v1.DebitQuantity
		}
	}
	return FAbs(totalValue), FAbs(totalValue / totalQuantity), FAbs(totalQuantity)
}
