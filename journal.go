package main

import (
	"errors"
	"fmt"
	"log"
	"math"
	"reflect"
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

func FConvertTheSignOfDoubleEntryToSingleEntry(isCredit, isCreditInTheEntry bool, number float64) float64 {
	number = math.Abs(number)
	if isCredit != isCreditInTheEntry {
		return -number
	}
	return number
}

func FSetAccount(entries []SAPQ1) []SAPQ12SAccount1 {
	var newEntries []SAPQ12SAccount1
	for _, v1 := range entries {
		account, _, errAccountName := FFindAccountFromName(v1.AccountName)

		var err SAPQ2
		if errAccountName != nil {
			err.AccountName = TErr(errAccountName.Error())
		}
		if account.CostFlowType == CHighLevelAccount {
			err.AccountName = TErr("is " + CHighLevelAccount)
		}

		newEntries = append(newEntries, SAPQ12SAccount1{
			SAPQ1:     SAPQ1{v1.AccountName, math.Abs(v1.Price), v1.Quantity},
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

func FSetErr(entries []SAPQ12SAccount1) {
	for k1 := range entries {
		if entries[k1].SAPQ1.Price == 0 {
			entries[k1].SAPQ2.Price = TErr("is zero")
		}
		if entries[k1].SAPQ1.Quantity == 0 {
			entries[k1].SAPQ2.Quantity = TErr("is zero")
		}
	}
}

func FSetPriceAndQuantityByQuantity(account SAPQ12SAccount1, insert bool) SAPQ12SAccount1 {
	if account.SAPQ1.Quantity >= 0 {
		return account
	}

	keys, inventory := FSortInventoryByCostFlow(account)

	QuantityCount := math.Abs(account.SAPQ1.Quantity)
	var costs float64
	for k1, v1 := range inventory {
		if v1.AccountName == account.SAPQ1.AccountName {
			if QuantityCount < v1.Quantity {
				costs -= v1.Price * QuantityCount
				if insert {
					inventory[k1].Quantity -= QuantityCount
					FDbUpdate(VDbInventory, keys[k1], inventory[k1])
				}
				QuantityCount = 0
				break
			}
			if QuantityCount >= v1.Quantity {
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

func FSetPriceAndQuantityByValue(account SAPQ12SAccount1) SAPQ12SAccount1 {
	value := account.SAPQ1.Price * account.SAPQ1.Quantity
	if value >= 0 {
		return account
	}

	_, inventory := FSortInventoryByCostFlow(account)

	ValueCount := math.Abs(value)
	var costs, quantity float64
	for _, v1 := range inventory {
		if v1.AccountName == account.SAPQ1.AccountName {
			value := v1.Price * v1.Quantity
			if ValueCount < value {
				quantity -= ValueCount / v1.Price
				costs -= ValueCount
				ValueCount = 0
				break
			}
			if ValueCount >= value {
				quantity -= value / v1.Price
				costs -= value
				ValueCount -= value
			}
		}
	}

	if ValueCount > 0 {
		account.SAPQ2.Quantity = TErr("you don't have enough quantity")
	}
	account.SAPQ1.Quantity = quantity
	account.SAPQ1.Price = costs / quantity
	return account
}

func FSortInventoryByCostFlow(account SAPQ12SAccount1) ([][]byte, []SAPQ1) {
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
		log.Println("the cost flow type ", account.SAccount1.CostFlowType, " is not in ", account.SAccount1.CostFlowType, VLowCostFlowType)
	}
	return keys, inventory
}

func FSeperateDebitFromCredit(entries []SAPQ12SAccount1) ([]SAPQ12SAccount1, []SAPQ12SAccount1, float64, float64) {
	var debitEntries, creditEntries []SAPQ12SAccount1
	var debit, credit float64
	for _, v1 := range entries {
		value := v1.SAPQ1.Price * v1.SAPQ1.Quantity
		if value == 0 {
			continue
		}

		switch v1.SAccount1.IsCredit {
		case false:
			if value > 0 {
				debitEntries = append(debitEntries, v1)
				debit += math.Abs(value)
			} else {
				creditEntries = append(creditEntries, v1)
				credit += math.Abs(value)
			}
		case true:
			if value < 0 {
				debitEntries = append(debitEntries, v1)
				debit += math.Abs(value)
			} else {
				creditEntries = append(creditEntries, v1)
				credit += math.Abs(value)
			}
		}
	}
	return debitEntries, creditEntries, debit, credit
}

func FCheckDebitEqualCredit(debit, credit float64) error {
	if debit != credit {
		return fmt.Errorf("the debit and credit should be equal. and the debit is %v and credit is %v", debit, credit)
	}
	return nil
}

func FCheckOneDebitOrOneCredit(debit, credit []SAPQ12SAccount1) error {
	if len(debit) != 1 && len(credit) != 1 {
		return fmt.Errorf("should be one debit or one credit in the entry \n the debit entries is:\n\t %v \n the credit entries is:\n\t %v", debit, credit)
	}
	return nil
}

func FInsertToJournal(debitEntries, creditEntries []SAPQ12SAccount1, entryInfo SEntry) []SJournal1 {
	var simpleEntries []SJournal1
	for _, debitEntry := range debitEntries {
		for _, creditEntry := range creditEntries {
			value := FSmallest(math.Abs(debitEntry.SAPQ1.Price*debitEntry.SAPQ1.Quantity), math.Abs(creditEntry.SAPQ1.Price*creditEntry.SAPQ1.Quantity))
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
				TypeOfCompoundEntry:        entryInfo.Label,
			})
		}
	}
	return simpleEntries
}

func FTotalValuePriceQuantity(account string) (float64, float64, float64) {
	_, journal := FDbRead[SJournal1](VDbJournal)
	for k1 := len(journal) - 1; k1 >= 0; k1-- {
		a := journal[k1]
		if account == a.DebitAccountName {
			return a.DebitBalanceValue, a.DebitBalancePrice, a.DebitBalanceQuantity
		}
		if account == a.CreditAccountName {
			return a.CreditBalanceValue, a.CreditBalancePrice, a.CreditBalanceQuantity
		}
	}
	return 0, 0, 0
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

func FInsertToDatabaseJournal(journal []SJournal1) {
	last := FDbLastLine[SJournal1](VDbJournal)
	for k1, v1 := range journal {
		v1.Date = time.Now()
		v1.EntryNumberCompound = last.EntryNumberCompound + 1
		v1.EntryNumberSimple = uint(k1) + 1

		if v1.DebitAccountName == v1.CreditAccountName {
			account, _, _ := FFindAccountFromName(v1.DebitAccountName)
			value, _, quantity := FTotalValuePriceQuantity(v1.DebitAccountName)

			if account.IsCredit {
				quantity -= v1.DebitQuantity
				quantity += v1.CreditQuantity
			} else {
				quantity += v1.DebitQuantity
				quantity -= v1.CreditQuantity
			}

			quantity = math.Abs(quantity)
			price := FConvertNanToZero(value / quantity)

			v1.DebitBalanceValue, v1.DebitBalancePrice, v1.DebitBalanceQuantity = value, price, quantity
			v1.CreditBalanceValue, v1.CreditBalancePrice, v1.CreditBalanceQuantity = value, price, quantity
		} else {
			calculateBalance := func(accountName string, isCreditInTheEntry bool, inputQuantity float64) (float64, float64, float64) {
				account, _, _ := FFindAccountFromName(accountName)
				value, _, quantity := FTotalValuePriceQuantity(accountName)

				if account.IsCredit == isCreditInTheEntry {
					value += v1.Value
					quantity += inputQuantity
				} else {
					value -= v1.Value
					quantity -= inputQuantity
				}

				return value, FConvertNanToZero(value / quantity), quantity
			}

			v1.DebitBalanceValue, v1.DebitBalancePrice, v1.DebitBalanceQuantity = calculateBalance(v1.DebitAccountName, false, v1.DebitQuantity)
			v1.CreditBalanceValue, v1.CreditBalancePrice, v1.CreditBalanceQuantity = calculateBalance(v1.CreditAccountName, true, v1.CreditQuantity)
		}

		FDbUpdate(VDbJournal, []byte(v1.Date.Format(CTimeLayout)), v1)
	}
}

func FInsertToDatabaseInventory(entries []SAPQ12SAccount1) {
	for _, v1 := range entries {
		if v1.SAPQ1.Quantity > 0 {
			FDbUpdate(VDbInventory, FNow(), SAPQ1{v1.SAPQ1.AccountName, v1.SAPQ1.Price, v1.SAPQ1.Quantity})
		} else {
			FSetPriceAndQuantityByQuantity(v1, true)
		}
	}
}

func FEntryAutoComplete(entries []SAPQ1, entryInfo SEntry, insert bool, accountToComplete string) ([]SAPQ12SAccount1, error) {
	newEntries1 := FSetAccount(entries)
	newEntries1 = FGroupByAccount(newEntries1)
	FSetErr(newEntries1)

	var oneEntry SAPQ12SAccount1
	for k1 := 0; k1 < len(newEntries1); k1++ {
		if newEntries1[k1].SAPQ1.AccountName == accountToComplete {
			oneEntry = newEntries1[k1]
			newEntries1 = FRemove(newEntries1, k1)
			continue
		}
		newEntries1[k1] = FSetPriceAndQuantityByQuantity(newEntries1[k1], false)
	}

	debitEntries, creditEntries, debit, credit := FSeperateDebitFromCredit(newEntries1)

	if debit != credit && !reflect.DeepEqual(oneEntry, SAPQ12SAccount1{}) {
		value := math.Abs(debit - credit)
		if oneEntry.IsCredit {
			if debit > credit {
				oneEntry.SAPQ1.Quantity = value / oneEntry.SAPQ1.Price
				creditEntries = append(creditEntries, oneEntry)
				credit += value
			} else {
				oneEntry.SAPQ1.Price = 1
				oneEntry.SAPQ1.Quantity = -value
				oneEntry = FSetPriceAndQuantityByValue(oneEntry)
				debitEntries = append(debitEntries, oneEntry)
				debit += value
			}
		} else {
			if debit > credit {
				oneEntry.SAPQ1.Price = 1
				oneEntry.SAPQ1.Quantity = -value
				oneEntry = FSetPriceAndQuantityByValue(oneEntry)
				creditEntries = append(creditEntries, oneEntry)
				credit += value
			} else {
				oneEntry.SAPQ1.Quantity = value / oneEntry.SAPQ1.Price
				debitEntries = append(debitEntries, oneEntry)
				debit += value
			}
		}
	}

	newEntries1 = append(debitEntries, creditEntries...)

	if err := FCheckDebitEqualCredit(debit, credit); err != nil {
		return newEntries1, err
	}

	if err := FCheckOneDebitOrOneCredit(debitEntries, creditEntries); err != nil {
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

func FEntryClose(accountToCloseName, accountToCloseWithName string, accountToCloseWithPrice float64, entryInfo SEntry, insert bool) error {
	accountToClose, _, accountToCloseErr := FFindAccountFromName(accountToCloseName)
	if accountToCloseErr != nil {
		return fmt.Errorf("%v is not found", accountToCloseName)
	}

	accountToCloseWith, _, accountToCloseWithErr := FFindAccountFromName(accountToCloseWithName)
	if accountToCloseWithErr != nil {
		return fmt.Errorf("%v is not found", accountToCloseWithName)
	}

	value, accountToClosePrice, accountToCloseQuantity := FTotalValuePriceQuantity(accountToCloseName)
	if value == 0 {
		return errors.New("you can't close the account because it has zero value")
	}

	accountToCloseWithPrice = math.Abs(accountToCloseWithPrice)
	accountQuantityToCloseWith := value / accountToCloseWithPrice
	if accountToCloseWith.IsCredit != accountToClose.IsCredit {
		accountQuantityToCloseWith *= -1
	}

	accountToCloseWith1 := FSetPriceAndQuantityByValue(SAPQ12SAccount1{SAPQ1{accountToCloseWithName, accountToCloseWithPrice, accountQuantityToCloseWith}, SAPQ2{}, accountToCloseWith})
	if accountToCloseWith1.SAPQ2.Quantity != "" {
		return fmt.Errorf("you don't have enough quantity for %v", accountToCloseWith1.SAPQ1.AccountName)
	}

	newEntries := []SAPQ12SAccount1{
		{SAPQ1{accountToCloseName, accountToClosePrice, -accountToCloseQuantity}, SAPQ2{}, accountToClose},
		accountToCloseWith1,
	}

	debitEntries, creditEntries, debit, credit := FSeperateDebitFromCredit(newEntries)
	FPanicIfErr(FCheckDebitEqualCredit(debit, credit))
	FPanicIfErr(FCheckOneDebitOrOneCredit(debitEntries, creditEntries))

	if insert {
		FInsertToDatabase(FInsertToJournal(debitEntries, creditEntries, entryInfo), newEntries)
	}

	return nil
}

func FEntryReconciliation(account1Name string, account1Quantity float64, entryInfo SEntry, insert bool) error {
	account2 := SAPQ1{AccountName: account1Name}
	_, account2.Price, account2.Quantity = FTotalValuePriceQuantity(account2.AccountName)
	account2.Quantity = -account2.Quantity

	account1And2Value := account2.Price * account2.Quantity

	account1 := SAPQ1{AccountName: account1Name}
	account1.Quantity = math.Abs(account1Quantity)
	account1.Price = account1And2Value / account1.Quantity

	newEntries1 := FSetAccount([]SAPQ1{account1, account2})
	FSetErr(newEntries1)

	for k1, v1 := range newEntries1 {
		newEntries1[k1] = FSetPriceAndQuantityByQuantity(v1, false)
	}

	for _, v1 := range newEntries1 {
		switch {
		case v1.SAPQ2.AccountName != "":
			return errors.New(v1.SAPQ1.AccountName + " " + string(v1.SAPQ2.AccountName))
		case v1.SAPQ2.Price != "":
			return errors.New(v1.SAPQ1.AccountName + " " + string(v1.SAPQ2.Price))
		case v1.SAPQ2.Quantity != "":
			return errors.New(v1.SAPQ1.AccountName + " " + string(v1.SAPQ2.Quantity))
		}
	}

	debitEntries, creditEntries, debit, credit := FSeperateDebitFromCredit(newEntries1)
	newEntries1 = append(debitEntries, creditEntries...)
	FPanicIfErr(FCheckDebitEqualCredit(debit, credit))
	FPanicIfErr(FCheckOneDebitOrOneCredit(debitEntries, creditEntries))

	if insert {
		FInsertToDatabase(FInsertToJournal(debitEntries, creditEntries, entryInfo), newEntries1)
	}

	return nil
}

func FEntryReconciliationWithAccount(account1 SAPQ1, account2Name string, account2Price float64, entryInfo SEntry, insert bool) error {
	account1.Price = math.Abs(account1.Price)
	account1.Quantity = math.Abs(account1.Quantity)

	account3 := SAPQ1{AccountName: account1.AccountName}
	_, account3.Price, account3.Quantity = FTotalValuePriceQuantity(account3.AccountName)
	account3.Quantity = -account3.Quantity

	account1Value := account1.Price * account1.Quantity
	account3Value := account3.Price * account3.Quantity
	account2Value := account1Value + account3Value

	newEntries1 := FSetAccount([]SAPQ1{account1, account3, {account2Name, account2Price, math.Abs(account2Value / account2Price)}})
	FSetErr(newEntries1)

	max, min := func() (float64, float64) {
		sign := func(number float64) float64 {
			if number > 0 {
				return 1
			}
			return -1
		}
		if math.Abs(account1Value) > math.Abs(account3Value) {
			return sign(account1Value), sign(account3Value)
		}
		return sign(account3Value), sign(account1Value)
	}()

	if newEntries1[2].IsCredit == newEntries1[0].IsCredit {
		newEntries1[2].SAPQ1.Quantity *= min
	} else {
		newEntries1[2].SAPQ1.Quantity *= max
	}

	for k1, v1 := range newEntries1 {
		newEntries1[k1] = FSetPriceAndQuantityByQuantity(v1, false)
	}

	debitEntries, creditEntries, debit, credit := FSeperateDebitFromCredit(newEntries1)
	newEntries1 = append(debitEntries, creditEntries...)

	err := FCheckDebitEqualCredit(debit, credit)
	if err != nil {
		return fmt.Errorf("you don't have enough %v to reconcile", account2Name)
	}

	err = FCheckOneDebitOrOneCredit(debitEntries, creditEntries)
	if err != nil {
		return fmt.Errorf("you don't have enough %v to reconcile", account2Name)
	}

	for _, v1 := range newEntries1 {
		switch {
		case v1.SAPQ2.AccountName != "":
			return errors.New(v1.SAPQ1.AccountName + " " + string(v1.SAPQ2.AccountName))
		case v1.SAPQ2.Price != "":
			return errors.New(v1.SAPQ1.AccountName + " " + string(v1.SAPQ2.Price))
		case v1.SAPQ2.Quantity != "":
			return errors.New(v1.SAPQ1.AccountName + " " + string(v1.SAPQ2.Quantity))
		}
	}

	if insert {
		FInsertToDatabase(FInsertToJournal(debitEntries, creditEntries, entryInfo), newEntries1)
	}

	return nil
}

func FEntryInvoice(AccountNamePay, AccountNameInvoiceDiscount string, PricePay, TotalInvoiceDiscount float64, entries []SAPQ1, entryInfo SEntry, insert bool) ([]SAPQ12SAccount1, error, error) {

	newEntries1 := FSetAccount(entries)
	newEntries1 = FGroupByAccount(newEntries1)
	FSetErr(newEntries1)

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
		v1.SAPQ1.Quantity = -math.Abs(v1.SAPQ1.Quantity)
		v1 = FSetPriceAndQuantityByQuantity(v1, false)
		autoCompletion, _, err := FFindAutoCompletionFromName(v1.SAPQ1.AccountName)
		if err != nil {
			v1.SAPQ2.AccountName = TErr(err.Error())
		}

		quantity := math.Abs(v1.SAPQ1.Quantity)
		addEntry := func(prefix string, price float64, isCredit bool) SAPQ12SAccount1 {
			if price > 0 {
				q := quantity
				if prefix == "" {
					q = -quantity
				}
				return SAPQ12SAccount1{
					SAPQ1:     SAPQ1{prefix + autoCompletion.AccountName, price, q},
					SAPQ2:     SAPQ2{},
					SAccount1: SAccount1{IsCredit: isCredit, CostFlowType: v1.CostFlowType},
				}
			}
			return SAPQ12SAccount1{}
		}

		var entries entry
		entries.Inventory = addEntry("", v1.SAPQ1.Price, false)
		entries.Cost = addEntry(CPrefixCost, v1.SAPQ1.Price, false)
		entries.TaxExpenses = addEntry(CPrefixTaxExpenses, autoCompletion.PriceTax, false)
		entries.TaxLiability = addEntry(CPrefixTaxLiability, autoCompletion.PriceTax, true)
		entries.Revenue = addEntry(CPrefixRevenue, autoCompletion.PriceRevenue, true)

		setDiscount := func(discountPrice float64) {
			entries.Discount = addEntry(CPrefixDiscount, discountPrice, false)
			totalvalueOfInvoiceDiscountAndPay += (autoCompletion.PriceRevenue - discountPrice) * quantity
		}

		switch autoCompletion.DiscountWay {
		case CDiscountPerOne:
			setDiscount(autoCompletion.DiscountPerOne)
		case CDiscountTotal:
			discountPrice := math.Abs(autoCompletion.DiscountTotal / quantity)
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
	var newEntries4 []SAPQ12SAccount1
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

		if v1.Inventory.SAPQ1.Price > 0 {
			newEntries3 = append(newEntries3, []SAPQ12SAccount1{v1.Inventory, v1.Cost})
		}
		if v1.TaxExpenses.SAPQ1.Price > 0 {
			newEntries3 = append(newEntries3, []SAPQ12SAccount1{v1.TaxExpenses, v1.TaxLiability})
		}
		if v1.Revenue.SAPQ1.Price > 0 {
			newEntries3 = append(newEntries3, []SAPQ12SAccount1{v1.Revenue, v1.Discount, v1.Pay, v1.InvoiceDiscount})
		}

		newEntries4 = append(newEntries4, v1.Inventory, v1.Cost, v1.TaxExpenses, v1.TaxLiability, v1.Revenue, v1.Discount, v1.Pay, v1.InvoiceDiscount)
	}

	_, _, errAccountNamePay := FAccountTerms(AccountNamePay, true, false)
	_, _, errAccountNameInvoiceDiscount := FAccountTerms(AccountNameInvoiceDiscount, true, false)

	if errAccountNamePay != nil || errAccountNameInvoiceDiscount != nil {
		return newEntries1, errAccountNamePay, errAccountNameInvoiceDiscount
	}

	for _, v1 := range newEntries1 {
		if v1.SAPQ2.AccountName != "" || v1.SAPQ2.Price != "" || v1.SAPQ2.Quantity != "" {
			return newEntries1, errAccountNamePay, errAccountNameInvoiceDiscount
		}
	}

	var newEntries5 []SJournal1
	for _, v1 := range newEntries3 {
		debitEntries, creditEntries, debit, credit := FSeperateDebitFromCredit(v1)
		FPanicIfErr(FCheckDebitEqualCredit(debit, credit))
		FPanicIfErr(FCheckOneDebitOrOneCredit(debitEntries, creditEntries))
		newEntries5 = append(newEntries5, FInsertToJournal(debitEntries, creditEntries, entryInfo)...)
	}

	if insert {
		FInsertToDatabase(newEntries5, newEntries4)
	}

	return newEntries1, errAccountNamePay, errAccountNameInvoiceDiscount
}

func FEntryReverse(entriesKeys [][]byte, entries []SJournal1, nameEmployee string) {
	var entryToReverse []SJournal1
	for k1, v1 := range entries {
		if v1.IsReversed {
			continue
		}
		account, _, _ := FFindAccountFromName(v1.CreditAccountName)
		v1.CreditQuantity = FConvertTheSignOfDoubleEntryToSingleEntry(account.IsCredit, true, v1.CreditQuantity)
		account, _, _ = FFindAccountFromName(v1.DebitAccountName)
		v1.DebitQuantity = FConvertTheSignOfDoubleEntryToSingleEntry(account.IsCredit, false, v1.DebitQuantity)

		entryCredit := FSetPriceAndQuantityByQuantity(SAPQ12SAccount1{
			SAPQ1:     SAPQ1{v1.CreditAccountName, v1.CreditPrice, v1.CreditQuantity},
			SAPQ2:     SAPQ2{},
			SAccount1: account,
		}, false)
		entryDebit := FSetPriceAndQuantityByQuantity(SAPQ12SAccount1{
			SAPQ1:     SAPQ1{v1.DebitAccountName, v1.DebitPrice, v1.DebitQuantity},
			SAPQ2:     SAPQ2{},
			SAccount1: account,
		}, false)

		if entryCredit.SAPQ1.Quantity == v1.CreditQuantity && entryDebit.SAPQ1.Quantity == v1.DebitQuantity {

			entryCredit.SAccount1.CostFlowType = CWma
			entryDebit.SAccount1.CostFlowType = CWma

			FInsertToDatabaseInventory([]SAPQ12SAccount1{entryCredit, entryDebit})

			v1.CreditPrice, v1.DebitPrice = v1.DebitPrice, v1.CreditPrice
			v1.CreditQuantity, v1.DebitQuantity = math.Abs(v1.DebitQuantity), math.Abs(v1.CreditQuantity)
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

func FFindDuplicateElement(journal []SJournal1, f SJournal3) []SJournal1 {
	if !f.Date &&
		!f.IsReverse &&
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
		return journal
	}

	var newJournal []SJournal1
	for k1, v1 := range journal {
		for k2, v2 := range journal {
			if k1 != k2 &&
				FFilterDuplicate(v1.Date, v2.Date, f.Date) &&
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
				break
			}
		}
	}
	return newJournal
}
