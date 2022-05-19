package main

import (
	"fmt"
	"math"
	"time"
)

func ValueAfterAdjustUsingAdjustingMethods(adjustingMethod string, minutesCurrent, minutesTotal, minutesPast, valueTotal float64) float64 {
	percent := Root(valueTotal, minutesTotal)
	switch adjustingMethod {
	case Exponential:
		return math.Pow(percent, minutesPast+minutesCurrent) - math.Pow(percent, minutesPast)
	case Logarithmic:
		return (valueTotal / math.Pow(percent, minutesPast)) - (valueTotal / math.Pow(percent, minutesPast+minutesCurrent))
	default:
		return minutesCurrent * (valueTotal / minutesTotal)
	}
}

func CheckDebitEqualCredit(entries []APQA) ([]APQA, []APQA, error) {
	var debitEntries, creditEntries []APQA
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

func SetPriceAndQuantity(account APQA, insert bool) APQA {
	if account.Quantity > 0 {
		return account
	}

	var keys [][]byte
	var inventory []APQ
	switch account.Account.CostFlowType {
	case Fifo:
		keys, inventory = DbRead[APQ](DbInventory)
	case Lifo:
		keys, inventory = DbRead[APQ](DbInventory)
		ReverseSlice(keys)
		ReverseSlice(inventory)
	case Wma:
		WeightedAverage(account.Name)
		keys, inventory = DbRead[APQ](DbInventory)
	}

	QuantityCount := Abs(account.Quantity)
	var costs float64
	for k1, v1 := range inventory {
		if v1.Name == account.Name {
			if QuantityCount <= v1.Quantity {
				costs -= v1.Price * QuantityCount
				if insert {
					inventory[k1].Quantity -= QuantityCount
					DbUpdate(DbInventory, keys[k1], inventory[k1])
				}
				QuantityCount = 0
				break
			}
			if QuantityCount > v1.Quantity {
				costs -= v1.Price * v1.Quantity
				if insert {
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

func GroupByAccount(entries []APQA) []APQA {
	m := map[string]*APQA{}
	for _, v1 := range entries {
		key := v1.Name
		sums := m[key]
		if sums == nil {
			sums = &APQA{}
			m[key] = sums
		}
		sums.Name = v1.Name
		sums.Price += v1.Price * v1.Quantity //here i store the VALUE in Price field
		sums.Quantity += v1.Quantity
		sums.Account = v1.Account
	}
	entries = []APQA{}
	for _, v1 := range m {
		entries = append(entries, APQA{
			Name:     v1.Name,
			Price:    v1.Price / v1.Quantity,
			Quantity: v1.Quantity,
			Account:  v1.Account,
		})
	}
	return entries
}

func InsertToDatabaseJournal(journal []Journal) {
	last := DbLastLine[Journal](DbJournal)
	for k1, v1 := range journal {
		v1.EntryNumberCompound = last.EntryNumberCompound + 1
		v1.EntryNumberSimple = k1 + 1

		DbUpdate(DbJournal, Now(), v1)
	}
}

func InsertToJournal(debitEntries, creditEntries []APQA, entryInfo EntryInfo) []Journal {
	var simpleEntries []Journal
	for _, debitEntry := range debitEntries {
		for _, creditEntry := range creditEntries {
			VALUE := Smallest(Abs(debitEntry.Price*debitEntry.Quantity), Abs(creditEntry.Price*creditEntry.Quantity))
			simpleEntries = append(simpleEntries, Journal{
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

func SimpleJournalEntry(entries []APQB, entryInfo EntryInfo, insert bool) ([]APQB, error) {
	newEntries1 := Stage1(entries, false)
	newEntries1 = GroupByAccount(newEntries1)

	for k1, v1 := range newEntries1 {
		newEntries1[k1] = SetPriceAndQuantity(v1, false)
	}

	debitEntries, creditEntries, err := CheckDebitEqualCredit(newEntries1)
	newEntries2 := ConvertAPQICToAPQB(append(debitEntries, creditEntries...))
	if err != nil {
		return newEntries2, err
	}

	if insert {
		simpleEntries := InsertToJournal(debitEntries, creditEntries, entryInfo)
		InsertToDatabaseJournal(simpleEntries)
		InsertToDatabaseInventory(newEntries1)
	}

	return newEntries2, nil
}

func InvoiceJournalEntry(payAccountName string, payAccountPrice, invoiceDiscountPrice float64, inventoryAccounts []APQB, entryInfo EntryInfo, insert bool) ([]APQB, error) {

	if _, isExist, err := AccountTerms(payAccountName, true, false); !isExist || err != nil {
		return inventoryAccounts, err
	}
	if _, isExist, err := AccountTerms(InvoiceDiscount, true, false); !isExist || err != nil {
		return inventoryAccounts, err
	}

	newEntries1 := Stage1(inventoryAccounts, true)
	newEntries1 = GroupByAccount(newEntries1)

	for k1, v1 := range newEntries1 {
		newEntries1[k1] = SetPriceAndQuantity(v1, false)
	}

	newEntries2 := AutoComplete(newEntries1, payAccountName, payAccountPrice)
	inventoryAccounts = ConvertAPQICToAPQB(newEntries1)

	var newEntries3 []Journal
	newEntries1 = []APQA{}
	for _, v1 := range newEntries2 {
		debitEntries, creditEntries, err := CheckDebitEqualCredit(v1)
		if err != nil {
			return inventoryAccounts, err
		}
		newEntries1 = append(newEntries1, debitEntries...)
		newEntries1 = append(newEntries1, creditEntries...)
		newEntries3 = append(newEntries3, InsertToJournal(debitEntries, creditEntries, entryInfo)...)
	}

	if insert {
		InsertToDatabaseJournal(newEntries3)
		InsertToDatabaseInventory(newEntries1)
	}

	return inventoryAccounts, nil
}

func AutoComplete(inventoryAccounts []APQA, payAccountName string, payAccountPrice float64) [][]APQA {
	var simpleInfoAPQ [][]APQA
	for _, v1 := range inventoryAccounts {

		autoCompletion, _, err := FindAutoCompletionFromName(v1.Name)
		if err != nil {
			continue
		}

		var inv []APQA
		var tax []APQA
		var pay []APQA

		quantity := Abs(v1.Quantity)
		var valueRevenue float64
		var valueDiscount float64

		inv = append(inv, APQA{autoCompletion.AccountInvnetory, v1.Price, v1.Quantity, Account{IsCredit: false}})
		inv = append(inv, APQA{PrefixCost + autoCompletion.AccountInvnetory, v1.Price, quantity, Account{IsCredit: false}})

		if autoCompletion.PriceTax > 0 {
			tax = append(tax, APQA{PrefixTaxExpenses + autoCompletion.AccountInvnetory, autoCompletion.PriceTax, quantity, Account{IsCredit: false}})
			tax = append(tax, APQA{PrefixTaxLiability + autoCompletion.AccountInvnetory, autoCompletion.PriceTax, quantity, Account{IsCredit: true}})
		}

		if autoCompletion.PriceRevenue > 0 {
			pay = append(pay, APQA{PrefixRevenue + autoCompletion.AccountInvnetory, autoCompletion.PriceRevenue, quantity, Account{IsCredit: true}})
			valueRevenue = quantity * autoCompletion.PriceRevenue
		}

		if priceDiscount := MaxDiscount(autoCompletion.PriceDiscount, quantity); priceDiscount > 0 {
			pay = append(pay, APQA{PrefixDiscount + autoCompletion.AccountInvnetory, priceDiscount, quantity, Account{IsCredit: false}})
			valueDiscount = quantity * priceDiscount
		}

		payValue := valueRevenue - valueDiscount
		pay = append(pay, APQA{payAccountName, payAccountPrice, payValue / payAccountPrice, Account{IsCredit: false}})

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

func ConvertAPQICToAPQB(entries []APQA) []APQB {
	var newEntries []APQB
	for _, v1 := range entries {
		newEntries = append(newEntries, APQB{v1.Name, v1.Price, v1.Quantity, ""})
	}
	return newEntries
}

func InsertToDatabaseInventory(entries []APQA) {
	for _, v1 := range entries {
		if v1.Quantity > 0 {
			DbUpdate(DbInventory, Now(), APQ{v1.Name, v1.Price, v1.Quantity})
		} else {
			SetPriceAndQuantity(v1, true)
		}
	}
}

func Stage1(entries []APQB, isInvoice bool) []APQA {
	var newEntries []APQA
	for _, v1 := range entries {
		account, _, err := FindAccountFromBarcode(v1.Barcode)
		if err != nil {
			account, _, err = FindAccountFromName(FormatTheString(v1.Name))
		}
		if isInvoice {
			if account.IsCredit || v1.Quantity >= 0 {
				continue
			}
			v1.Price = 1
		}
		if err == nil && account.IsLowLevel && v1.Quantity != 0 && v1.Price != 0 {
			newEntries = append(newEntries, APQA{
				Name:     account.Name,
				Price:    Abs(v1.Price),
				Quantity: v1.Quantity,
				Account:  account,
			})
		}
	}
	return newEntries
}

func ReverseEntries(entriesKeys [][]byte, entries []Journal, nameEmployee string) {
	var entryToReverse []Journal
	for k1, v1 := range entries {
		if v1.IsReversed {
			continue
		}
		account, _, _ := FindAccountFromName(v1.AccountCredit)
		v1.QuantityCredit = ConvertTheSignOfDoubleEntryToSingleEntry(account.IsCredit, true, v1.QuantityCredit)
		account, _, _ = FindAccountFromName(v1.AccountDebit)
		v1.QuantityDebit = ConvertTheSignOfDoubleEntryToSingleEntry(account.IsCredit, false, v1.QuantityDebit)

		entryCredit := SetPriceAndQuantity(APQA{v1.AccountCredit, v1.PriceCredit, v1.QuantityCredit, Account{IsCredit: false, CostFlowType: Fifo}}, false)
		entryDebit := SetPriceAndQuantity(APQA{v1.AccountDebit, v1.PriceDebit, v1.QuantityDebit, Account{IsCredit: false, CostFlowType: Fifo}}, false)

		if entryCredit.Quantity == v1.QuantityCredit && entryDebit.Quantity == v1.QuantityDebit {

			entryCredit.Account.CostFlowType = Wma
			entryDebit.Account.CostFlowType = Wma

			InsertToDatabaseInventory([]APQA{entryCredit, entryDebit})

			v1.PriceCredit, v1.PriceDebit = v1.PriceDebit, v1.PriceCredit
			v1.QuantityCredit, v1.QuantityDebit = Abs(v1.QuantityDebit), Abs(v1.QuantityCredit)
			v1.AccountCredit, v1.AccountDebit = v1.AccountDebit, v1.AccountCredit

			v1.IsReverse = true
			v1.ReverseEntryNumberCompound = v1.EntryNumberCompound
			v1.ReverseEntryNumberSimple = v1.EntryNumberSimple
			v1.Notes = "revese entry for entry was entered by " + v1.Employee
			v1.Employee = nameEmployee

			entryToReverse = append(entryToReverse, v1)

			entries[k1].IsReversed = true
			DbUpdate(DbJournal, entriesKeys[k1], entries[k1])
		}
	}

	InsertToDatabaseJournal(entryToReverse)
}

func FindEntryFromNumber(entryNumberCompound int, entryNumberSimple int) ([][]byte, []Journal) {
	var entries []Journal
	var entriesKeys [][]byte
	keys, journal := DbRead[Journal](DbJournal)
	for k1, v1 := range journal {
		if v1.EntryNumberCompound == entryNumberCompound && (entryNumberSimple == 0 || v1.EntryNumberSimple == entryNumberSimple) {
			entries = append(entries, v1)
			entriesKeys = append(entriesKeys, keys[k1])
		}
	}
	return entriesKeys, entries
}

func ConvertJournalToAPQA(entries []Journal) []APQA {
	var newEntries []APQA
	for _, v1 := range entries {
		account, _, _ := FindAccountFromName(v1.AccountDebit)
		newEntries = append(newEntries, APQA{
			Name:     v1.AccountDebit,
			Price:    v1.PriceDebit,
			Quantity: ConvertTheSignOfDoubleEntryToSingleEntry(account.IsCredit, false, v1.QuantityDebit),
			Account:  account,
		})

		account, _, _ = FindAccountFromName(v1.AccountCredit)
		newEntries = append(newEntries, APQA{
			Name:     v1.AccountCredit,
			Price:    v1.PriceCredit,
			Quantity: ConvertTheSignOfDoubleEntryToSingleEntry(account.IsCredit, true, v1.QuantityCredit),
			Account:  account,
		})
	}
	return GroupByAccount(newEntries)
}

func ConvertAPQAToAPQB(entries []APQA) []APQB {
	var newEntries []APQB
	for _, v1 := range entries {
		newEntries = append(newEntries, APQB{
			Name:     v1.Name,
			Price:    v1.Price,
			Quantity: v1.Quantity,
			Barcode:  v1.Account.Barcode[0],
		})
	}
	return newEntries
}
func ExtractEntryInfoFromJournal(entry Journal) EntryInfo {
	return EntryInfo{
		Notes:               entry.Notes,
		Name:                entry.Name,
		Employee:            entry.Employee,
		TypeOfCompoundEntry: entry.TypeOfCompoundEntry,
	}
}

func JournalFilter(dates []time.Time, journal []Journal, f FilterJournal, isDebitAndCredit bool) ([]time.Time, []Journal) {
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

	var newJournal []Journal
	var newDates []time.Time
	for k1, v1 := range journal {

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
			f.Employee.Filter(v1.Employee) &&
			f.TypeOfCompoundEntry.Filter(v1.TypeOfCompoundEntry) {
			newJournal = append(newJournal, v1)
			newDates = append(newDates, dates[k1])
		}
	}
	return newDates, newJournal
}

func FindDuplicateElement(dates []time.Time, journal []Journal, f FilterJournalDuplicate) ([]time.Time, []Journal) {
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

	var newJournal []Journal
	var newDates []time.Time
	for k1, v1 := range journal {
		for k2, v2 := range journal {
			if k1 != k2 &&
				FilterDuplicate(v1.IsReverse, v2.IsReverse, f.IsReverse) &&
				FilterDuplicate(v1.IsReversed, v2.IsReversed, f.IsReversed) &&
				FilterDuplicate(v1.ReverseEntryNumberCompound, v2.ReverseEntryNumberCompound, f.ReverseEntryNumberCompound) &&
				FilterDuplicate(v1.ReverseEntryNumberSimple, v2.ReverseEntryNumberSimple, f.ReverseEntryNumberSimple) &&
				FilterDuplicate(v1.Value, v2.Value, f.Value) &&
				FilterDuplicate(v1.PriceDebit, v2.PriceDebit, f.PriceDebit) &&
				FilterDuplicate(v1.PriceCredit, v2.PriceCredit, f.PriceCredit) &&
				FilterDuplicate(v1.QuantityDebit, v2.QuantityDebit, f.QuantityDebit) &&
				FilterDuplicate(v1.QuantityCredit, v2.QuantityCredit, f.QuantityCredit) &&
				FilterDuplicate(v1.AccountDebit, v2.AccountDebit, f.AccountDebit) &&
				FilterDuplicate(v1.AccountCredit, v2.AccountCredit, f.AccountCredit) &&
				FilterDuplicate(v1.Notes, v2.Notes, f.Notes) &&
				FilterDuplicate(v1.Name, v2.Name, f.Name) &&
				FilterDuplicate(v1.Employee, v2.Employee, f.Employee) &&
				FilterDuplicate(v1.TypeOfCompoundEntry, v2.TypeOfCompoundEntry, f.TypeOfCompoundEntry) {
				newJournal = append(newJournal, v1)
				newDates = append(newDates, dates[k1])
				break
			}
		}
	}
	return newDates, newJournal
}

func MaxDiscount(discounts []PQ, quantity float64) float64 {
	var price float64
	for _, v1 := range discounts {
		if v1.Quantity > quantity {
			price = v1.Price
		}
	}
	return Abs(price)
}

func FilterJournalFromReverseEntry(keys [][]byte, journal []Journal) ([][]byte, []Journal) {
	var newKeys [][]byte
	var newJournal []Journal
	for k1, v1 := range journal {
		if !v1.IsReverse && !v1.IsReversed {
			newKeys = append(newKeys, keys[k1])
			newJournal = append(newJournal, v1)
		}
	}
	return newKeys, newJournal
}

func ConvertTheSignOfDoubleEntryToSingleEntry(isCredit, isCreditInTheEntry bool, number float64) float64 {
	number = math.Abs(number)
	if isCredit != isCreditInTheEntry {
		return -number
	}
	return number
}
