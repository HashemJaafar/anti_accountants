package main

import (
	"errors"
	"fmt"
	"math"
	"sync"
	"time"
)

//TODO
// func FValueAfterAdjustUsingAdjustingMethods(adjustingMethod string, minutesCurrent, minutesTotal, minutesPast, valueTotal float64) float64 {
// 	percent := FRoot(valueTotal, minutesTotal)
// 	switch adjustingMethod {
// 	case CExponential:
// 		return math.Pow(percent, minutesPast+minutesCurrent) - math.Pow(percent, minutesPast)
// 	case CLogarithmic:
// 		return (valueTotal / math.Pow(percent, minutesPast)) - (valueTotal / math.Pow(percent, minutesPast+minutesCurrent))
// 	default:
// 		return minutesCurrent * (valueTotal / minutesTotal)
// 	}
// }

func FCheckDebitEqualCredit(entries []SAPQ12SAccount1) ([]SAPQ12SAccount1, []SAPQ12SAccount1, error) {
	var debitEntries, creditEntries []SAPQ12SAccount1
	var zero float64
	for _, v1 := range entries {
		value := v1.SAPQ1.Price * v1.SAPQ1.Quantity
		if value == 0 {
			continue
		}

		switch v1.SAccount1.IsCredit {
		case false:
			zero += value
			if value > 0 {
				debitEntries = append(debitEntries, v1)
			} else {
				creditEntries = append(creditEntries, v1)
			}
		case true:
			zero -= value
			if value < 0 {
				debitEntries = append(debitEntries, v1)
			} else {
				creditEntries = append(creditEntries, v1)
			}
		}
	}
	s := fmt.Sprintf(", the debitEntries is %v and the creditEntries is %v", debitEntries, creditEntries)
	if len(debitEntries) != 1 && len(creditEntries) != 1 {
		return debitEntries, creditEntries, errors.New("should be one debit or one credit in the entry" + s)
	}
	if zero != 0 {
		return debitEntries, creditEntries, fmt.Errorf("the debit and credit should be equal. and the debit is more than credit by %v"+s, zero)
	}
	return debitEntries, creditEntries, nil
}

func FSetPriceAndQuantity(account SAPQ12SAccount1, insert bool) SAPQ12SAccount1 {
	if account.SAPQ1.Quantity > 0 {
		return account
	}

	var keys [][]byte
	var inventory []SAPQ1
	switch account.SAccount1.CostFlowType {
	case CFifo:
		keys, inventory = FDbRead[SAPQ1](VDbInventory)
	case CLifo:
		keys, inventory = FDbRead[SAPQ1](VDbInventory)
		FReverseSlice(keys)
		FReverseSlice(inventory)
	case CWma:
		FWeightedAverage(account.SAPQ1.AccountName)
		keys, inventory = FDbRead[SAPQ1](VDbInventory)
	default:
		return account
	}

	QuantityCount := FAbs(account.SAPQ1.Quantity)
	var costs float64
	for k1, v1 := range inventory {
		if v1.AccountName == account.SAPQ1.AccountName {
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

	if QuantityCount > 0 {
		account.SAPQ2.Quantity = TErr("the quantity is too much")
	}
	account.SAPQ1.Quantity += QuantityCount
	account.SAPQ1.Price = costs / account.SAPQ1.Quantity
	return account
}

func FInsertToDatabaseJournal(journal []SJournal1) {
	last := FDbLastLine[SJournal1](VDbJournal)
	for k1, v1 := range journal {
		v1.Date = time.Now()
		v1.EntryNumberCompound = last.EntryNumberCompound + 1
		v1.EntryNumberSimple = uint(k1) + 1

		calculateBalance := func(accountName string, isCreditInTheEntry bool, inputQuantity float64) (float64, float64, float64) {

			account, _, _ := FFindAccountFromName(accountName)
			sign := 1.0
			if account.IsCredit != isCreditInTheEntry {
				sign = -1
			}

			value, _, quantity := FTotalValuePriceQuantity(accountName)
			value += v1.Value * sign
			quantity += inputQuantity * sign

			return value, value / quantity, quantity
		}

		v1.DebitBalanceValue, v1.DebitBalancePrice, v1.DebitBalanceQuantity = calculateBalance(v1.DebitAccountName, false, v1.DebitQuantity)
		v1.CreditBalanceValue, v1.CreditBalancePrice, v1.CreditBalanceQuantity = calculateBalance(v1.CreditAccountName, true, v1.CreditQuantity)

		FDbUpdate(VDbJournal, []byte(v1.Date.Format(CTimeLayout)), v1)
	}
}

func FInsertToJournal(debitEntries, creditEntries []SAPQ12SAccount1, entryInfo SEntry) []SJournal1 {
	var simpleEntries []SJournal1
	for _, debitEntry := range debitEntries {
		for _, creditEntry := range creditEntries {
			value := FSmallest(FAbs(debitEntry.SAPQ1.Price*debitEntry.SAPQ1.Quantity), FAbs(creditEntry.SAPQ1.Price*creditEntry.SAPQ1.Quantity))
			simpleEntries = append(simpleEntries, SJournal1{
				Date:                       time.Time{},
				IsReverse:                  false,
				IsReversed:                 false,
				ReverseEntryNumberCompound: 0,
				ReverseEntryNumberSimple:   0,
				EntryNumberCompound:        0,
				EntryNumberSimple:          0,
				Value:                      value,
				DebitAccountName:           debitEntry.SAPQ1.AccountName,
				DebitPrice:                 debitEntry.SAPQ1.Price,
				DebitQuantity:              value / debitEntry.SAPQ1.Price,
				DebitBalanceValue:          0,
				DebitBalancePrice:          0,
				DebitBalanceQuantity:       0,
				CreditAccountName:          creditEntry.SAPQ1.AccountName,
				CreditPrice:                creditEntry.SAPQ1.Price,
				CreditQuantity:             value / creditEntry.SAPQ1.Price,
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

func FSetEntries(entries []SAPQ1) []SAPQ12SAccount1 {
	var newEntries []SAPQ12SAccount1
	for _, v1 := range entries {
		account, _, errAccountName := FFindAccountFromNameOrBarcode(v1.AccountName)

		var err SAPQ2
		if errAccountName != nil {
			err.AccountName = TErr(errAccountName.Error())
		}
		if account.CostFlowType == CHighLevelAccount {
			err.AccountName = TErr("is " + CHighLevelAccount)
		}
		if v1.Price == 0 {
			err.Price = TErr("is zero")
		}
		if v1.Quantity == 0 {
			err.Quantity = TErr("is zero")
		}

		newEntries = append(newEntries, SAPQ12SAccount1{
			SAPQ1:     SAPQ1{v1.AccountName, FAbs(v1.Price), v1.Quantity},
			SAPQ2:     err,
			SAccount1: account,
		})
	}
	return newEntries
}

func FGroupByAccount(entries []SAPQ12SAccount1) []SAPQ12SAccount1 {
	m := map[string]*SAPQ12SAccount1{}
	for _, v1 := range entries {
		key := v1.SAPQ1.AccountName
		sums := m[key]
		if sums == nil {
			sums = &SAPQ12SAccount1{}
			m[key] = sums
		}
		sums.SAPQ1.AccountName = v1.SAPQ1.AccountName
		sums.SAPQ1.Price += v1.SAPQ1.Price * v1.SAPQ1.Quantity //here i store the VALUE in Price field
		sums.SAPQ1.Quantity += v1.SAPQ1.Quantity
		sums.SAPQ2 = v1.SAPQ2
		sums.SAccount1 = v1.SAccount1
	}
	entries = []SAPQ12SAccount1{}
	for _, v1 := range m {
		v1.SAPQ1.Price = v1.SAPQ1.Price / v1.SAPQ1.Quantity
		entries = append(entries, SAPQ12SAccount1{
			SAPQ1:     v1.SAPQ1,
			SAPQ2:     v1.SAPQ2,
			SAccount1: v1.SAccount1,
		})
	}
	return entries
}

func FSimpleJournalEntry(entries []SAPQ1, entryInfo SEntry, insert bool) ([]SAPQ12SAccount1, error) {
	newEntries1 := FSetEntries(entries)
	newEntries1 = FGroupByAccount(newEntries1)

	for k1, v1 := range newEntries1 {
		newEntries1[k1] = FSetPriceAndQuantity(v1, false)
	}

	debitEntries, creditEntries, err := FCheckDebitEqualCredit(newEntries1)
	newEntries1 = append(debitEntries, creditEntries...)

	if err != nil {
		return newEntries1, err
	}

	for _, v1 := range newEntries1 {
		if v1.SAPQ2.AccountName != "" || v1.SAPQ2.Price != "" || v1.SAPQ2.Quantity != "" {
			return newEntries1, nil
		}
	}

	if insert {
		FInsertToDatabase(FInsertToJournal(debitEntries, creditEntries, entryInfo), newEntries1)
	}

	return newEntries1, nil
}

func FInvoiceJournalEntry(AccountNamePay, AccountNameInvoiceDiscount string, PricePay, TotalInvoiceDiscount float64, entries []SAPQ1, entryInfo SEntry, insert bool) ([]SAPQ12SAccount1, error) {

	newEntries1 := FSetEntries(entries)
	newEntries1 = FGroupByAccount(newEntries1)

	type entry struct {
		Inventory       SAPQ12SAccount1
		Cost            SAPQ12SAccount1
		TaxExpenses     SAPQ12SAccount1
		TaxLiability    SAPQ12SAccount1
		Revenue         SAPQ12SAccount1
		Discount        SAPQ12SAccount1
		Pay             SAPQ12SAccount1
		InvoiceDiscount SAPQ12SAccount1
	}
	var totalvalueOfInvoiceDiscountAndPay float64
	var newEntries2 []entry
	for k1, v1 := range newEntries1 {

		v1.SAPQ2.Price = ""
		v1.SAPQ1.Quantity = -FAbs(v1.SAPQ1.Quantity)
		v1 = FSetPriceAndQuantity(v1, false)
		autoCompletion, _, err := FFindAutoCompletionFromName(v1.SAPQ1.AccountName)
		if err != nil {
			v1.SAPQ2.AccountName = TErr(err.Error())
		}

		quantity := FAbs(v1.SAPQ1.Quantity)
		addEntry := func(entry *SAPQ12SAccount1, prefix string, price float64, isCredit bool) {
			if price > 0 {
				q := quantity
				if prefix == "" {
					q = -quantity
				}
				*entry = SAPQ12SAccount1{
					SAPQ1:     SAPQ1{prefix + autoCompletion.AccountName, price, q},
					SAPQ2:     SAPQ2{},
					SAccount1: SAccount1{IsCredit: isCredit},
				}
			}
		}

		var entries entry
		addEntry(&entries.Inventory, "", v1.SAPQ1.Price, false)
		addEntry(&entries.Cost, CPrefixCost, v1.SAPQ1.Price, false)
		addEntry(&entries.TaxExpenses, CPrefixTaxExpenses, autoCompletion.PriceTax, false)
		addEntry(&entries.TaxLiability, CPrefixTaxLiability, autoCompletion.PriceTax, true)
		addEntry(&entries.Revenue, CPrefixRevenue, autoCompletion.PriceRevenue, true)

		setDiscount := func(discountPrice float64) {
			addEntry(&entries.Discount, CPrefixDiscount, discountPrice, false)
			totalvalueOfInvoiceDiscountAndPay += (autoCompletion.PriceRevenue - discountPrice) * quantity
		}

		switch autoCompletion.DiscountWay {
		case CDiscountPerOne:
			setDiscount(autoCompletion.DiscountPerOne)
		case CDiscountTotal:
			discountPrice := FAbs(autoCompletion.DiscountTotal / quantity)
			if discountPrice > autoCompletion.PriceRevenue {
				discountPrice = autoCompletion.PriceRevenue
			}
			setDiscount(discountPrice)
		case CDiscountPerQuantity:
			setDiscount(math.Floor(quantity/autoCompletion.DiscountPerQuantity.TQuantity) * autoCompletion.DiscountPerQuantity.TPrice)
		case CDiscountDecisionTree:
			var discountPrice float64
			for _, v2 := range autoCompletion.DiscountDecisionTree {
				if v2.TQuantity > quantity {
					discountPrice = v2.TPrice
				}
			}
			setDiscount(discountPrice)
		}

		newEntries2 = append(newEntries2, entries)
		newEntries1[k1] = v1
	}

	if TotalInvoiceDiscount > totalvalueOfInvoiceDiscountAndPay {
		TotalInvoiceDiscount = totalvalueOfInvoiceDiscountAndPay
	}

	var newEntries3 [][]SAPQ12SAccount1
	for _, v1 := range newEntries2 {
		valueRevenue := v1.Revenue.SAPQ1.Price * v1.Revenue.SAPQ1.Quantity
		valueDiscount := v1.Discount.SAPQ1.Price * v1.Discount.SAPQ1.Quantity
		valueOfInvoiceDiscountAndPay := valueRevenue - valueDiscount

		valueInvoiceDiscount := (TotalInvoiceDiscount / totalvalueOfInvoiceDiscountAndPay) * valueOfInvoiceDiscountAndPay
		v1.InvoiceDiscount = SAPQ12SAccount1{
			SAPQ1:     SAPQ1{AccountNameInvoiceDiscount, valueInvoiceDiscount, 1},
			SAPQ2:     SAPQ2{},
			SAccount1: SAccount1{IsCredit: false},
		}

		valuePay := valueOfInvoiceDiscountAndPay - valueInvoiceDiscount
		v1.Pay = SAPQ12SAccount1{
			SAPQ1:     SAPQ1{AccountNamePay, PricePay, valuePay / PricePay},
			SAPQ2:     SAPQ2{},
			SAccount1: SAccount1{IsCredit: false},
		}

		newEntries3 = append(newEntries3, []SAPQ12SAccount1{v1.Inventory, v1.Cost}, []SAPQ12SAccount1{v1.TaxExpenses, v1.TaxLiability}, []SAPQ12SAccount1{v1.Revenue, v1.Discount, v1.Pay, v1.InvoiceDiscount})
	}

	var newEntries4 []SJournal1
	for _, v1 := range newEntries3 {
		debitEntries, creditEntries, err := FCheckDebitEqualCredit(v1)
		if err != nil {
			fmt.Println(err)
			continue
		}
		newEntries4 = append(newEntries4, FInsertToJournal(debitEntries, creditEntries, entryInfo)...)
	}

	if _, isExist, err := FAccountTerms(AccountNamePay, true, false); !isExist || err != nil {
		return newEntries1, err
	}
	if _, isExist, err := FAccountTerms(AccountNameInvoiceDiscount, true, false); !isExist || err != nil {
		return newEntries1, err
	}

	for _, v1 := range newEntries1 {
		if v1.SAPQ2.AccountName != "" || v1.SAPQ2.Price != "" || v1.SAPQ2.Quantity != "" {
			return newEntries1, nil
		}
	}

	if insert {
		FInsertToDatabase(newEntries4, newEntries1)
	}

	return newEntries1, nil
}

func FInsertToDatabase(entriesJournal []SJournal1, entriesInventory []SAPQ12SAccount1) {
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

func FInsertToDatabaseInventory(entries []SAPQ12SAccount1) {
	for _, v1 := range entries {
		if v1.SAPQ1.Quantity > 0 {
			FDbUpdate(VDbInventory, FNow(), SAPQ1{v1.SAPQ1.AccountName, v1.SAPQ1.Price, v1.SAPQ1.Quantity})
		} else {
			FSetPriceAndQuantity(v1, true)
		}
	}
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

		entryCredit := FSetPriceAndQuantity(SAPQ12SAccount1{
			SAPQ1:     SAPQ1{v1.CreditAccountName, v1.CreditPrice, v1.CreditQuantity},
			SAPQ2:     SAPQ2{},
			SAccount1: account,
		}, false)
		entryDebit := FSetPriceAndQuantity(SAPQ12SAccount1{
			SAPQ1:     SAPQ1{v1.DebitAccountName, v1.DebitPrice, v1.DebitQuantity},
			SAPQ2:     SAPQ2{},
			SAccount1: account,
		}, false)

		if entryCredit.SAPQ1.Quantity == v1.CreditQuantity && entryDebit.SAPQ1.Quantity == v1.DebitQuantity {

			entryCredit.SAccount1.CostFlowType = CWma
			entryDebit.SAccount1.CostFlowType = CWma

			FInsertToDatabaseInventory([]SAPQ12SAccount1{entryCredit, entryDebit})

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
	_, journal := FDbRead[SJournal1](VDbJournal)
	for k1 := len(journal) - 1; k1 >= 0; k1-- {
		if journal[k1].DebitAccountName == account {
			return journal[k1].DebitBalanceValue, journal[k1].DebitBalancePrice, journal[k1].DebitBalanceQuantity
		}
		if journal[k1].CreditAccountName == account {
			return journal[k1].CreditBalanceValue, journal[k1].CreditBalancePrice, journal[k1].CreditBalanceQuantity
		}
	}
	return 0, 0, 0
}
