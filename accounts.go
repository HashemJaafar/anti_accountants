package anti_accountants

import (
	"errors"
	"fmt"
	"math"
	"reflect"
)

func FFindAccountFromNameOrBarcode(inventory, nameOrBarcode string) (SAccount1, int, error) {
	for k1, v1 := range VAccounts {
		if v1.Inventory == inventory {
			if nameOrBarcode == v1.Name {
				return v1, k1, nil
			}
			_, isIn := FFind(nameOrBarcode, v1.Barcode)
			if isIn {
				return v1, k1, nil
			}
		}
	}
	return SAccount1{}, 0, VErrorNotListed
}

func FFindAccountFromName(accountName string) (SAccount1, int, error) {
	for k1, v1 := range VAccounts {
		if accountName == v1.Name {
			return v1, k1, nil
		}
	}
	return SAccount1{}, 0, VErrorNotListed
}

func FFindAutoCompletionFromName(accountName string) (SAutoCompletion, int, error) {
	for k1, v1 := range VAutoCompletionEntries {
		if accountName == v1.AccountName {
			return v1, k1, nil
		}
	}
	return SAutoCompletion{}, 0, VErrorNotListed
}

func FSetTheAccountFieldToStandard(account SAccount1) (SAccount1, error) {
	filter := func(slice []string) []string {
		for k1 := 0; k1 < len(slice); k1++ {
			if slice[k1] == "" {
				slice = FRemove(slice, k1)
			} else {
				k1++
			}
		}
		return slice
	}

	account.Image = filter(account.Image)
	account.Barcode = filter(account.Barcode)

	if account.Name == "" {
		return account, errors.New("the account name is empty")
	}

	return account, nil
}

func FAddAccount(isSave bool, account SAccount1) SAccount3 {
	var isErr bool
	var accountError SAccount3

	e := func(err error, ifTreu bool) {
		if ifTreu {
			accountError.Name = TErr(err.Error())
			isErr = true
		}
	}

	account, err := FSetTheAccountFieldToStandard(account)
	e(err, err != nil)
	_, _, err = FFindAccountFromName(account.Name)
	e(VErrorIsUsed, err == nil)

	accountError.Barcode = FSetSliceOfTErr(FErrorUsedBarcode(account))
	accountError.Number = FSetSliceOfTErr(FErrorUsedNumber(account))

	if isErr || !isSave {
		return accountError
	}

	VAccounts = append(VAccounts, account)
	FSetTheAccounts()
	FDbInsertIntoAccounts()
	return accountError
}

func FErrorUsedNumber(account SAccount1) []error {
	err := make([]error, len(account.Number))
	for k1, v1 := range account.Number {
		for _, v2 := range VAccounts {
			for k3, v3 := range v2.Number {
				if k1 == k3 {
					if reflect.DeepEqual(v1, v3) {
						err[k1] = VErrorIsUsed
					}
				}
			}
		}
	}
	return err
}

func FErrorUsedBarcode(account SAccount1) []error {
	err := make([]error, len(account.Barcode))
	for k1, v1 := range account.Barcode {
		for _, v2 := range VAccounts {
			if v2.Inventory == account.Inventory {
				for _, v3 := range v2.Barcode {
					if v1 == v3 {
						err[k1] = VErrorIsUsed
					}
				}
			}
		}
	}
	return err
}

func FEditAccount(isDelete bool, index int, account SAccount1) {
	account, err := FSetTheAccountFieldToStandard(account)
	if err != nil {
		return
	}

	newAccountName := account.Name
	oldAccountName := VAccounts[index].Name

	if !FIsUsedInJournal(oldAccountName) {
		if isDelete {
			VAccounts = FRemove(VAccounts, index)
			FSetTheAccounts()
			FDbInsertIntoAccounts()
			return
		}

		VAccounts[index].CostFlowType = account.CostFlowType
		VAccounts[index].IsCredit = account.IsCredit
	}

	if oldAccountName != newAccountName && newAccountName != "" {
		_, _, err := FFindAccountFromName(newAccountName)
		if err != nil {
			FChangeAccountName(oldAccountName, newAccountName)
			VAccounts[index].Name = newAccountName
		}
	}

	_, isIn1 := FFind(VAccounts[index].CostFlowType, VLowCostFlowType)
	_, isIn2 := FFind(account.CostFlowType, VLowCostFlowType)
	if VAccounts[index].CostFlowType == CHighLevelAccount || (isIn1 && !isIn2) {
		VAccounts[index].CostFlowType = account.CostFlowType
	}

	VAccounts[index].Inventory = account.Inventory
	VAccounts[index].Barcode = account.Barcode
	VAccounts[index].Notes = account.Notes
	VAccounts[index].Image = account.Image
	VAccounts[index].Number = account.Number

	FSetTheAccounts()
	FDbInsertIntoAccounts()
}

func FIsItHighThanByOrder(accountNumber1, accountNumber2 []uint) bool {
	l1 := len(accountNumber1)
	l2 := len(accountNumber2)
	for k1 := 0; k1 < FSmallest(l1, l2); k1++ {
		if accountNumber1[k1] < accountNumber2[k1] {
			return true
		}
		if accountNumber1[k1] > accountNumber2[k1] {
			return false
		}
	}
	return l2 > l1
}

func FIsItSubAccountUsingNumber(higherLevelAccountNumber, lowerLevelAccountNumber []uint) bool {
	lenHigherLevelAccountNumber := len(higherLevelAccountNumber)
	lenLowerLevelAccountNumber := len(lowerLevelAccountNumber)

	if lenHigherLevelAccountNumber == 0 || lenLowerLevelAccountNumber == 0 ||
		lenHigherLevelAccountNumber >= lenLowerLevelAccountNumber {
		return false
	}

	for k1, v1 := range higherLevelAccountNumber {
		if v1 != lowerLevelAccountNumber[k1] {
			return false
		}
	}
	return true
}

func FIsItTheFather(higherLevelAccountNumber, lowerLevelAccountNumber []uint) bool {
	lenHigherLevelAccountNumber := len(higherLevelAccountNumber)
	lenLowerLevelAccountNumber := len(lowerLevelAccountNumber)
	if lenHigherLevelAccountNumber == 0 || lenLowerLevelAccountNumber == 0 {
		return false
	}
	return reflect.DeepEqual(higherLevelAccountNumber, lowerLevelAccountNumber[:lenLowerLevelAccountNumber-1])
}

func FIsUsedInJournal(accountName string) bool {
	_, journal := FDbRead[SJournal1](VDbJournal)
	for _, v1 := range journal {
		if accountName == v1.CreditAccountName || accountName == v1.DebitAccountName {
			return true
		}
	}
	return false
}

func FMaxLenForAccountNumber() int {
	var maxLen int
	for _, v1 := range VAccounts {
		length := len(v1.Number)
		if length > maxLen {
			maxLen = length
		}
	}
	return maxLen
}

func FSetTheAccounts() []SAccount3 {
	maxLen := FMaxLenForAccountNumber()

	for k1, v1 := range VAccounts {
		VAccounts[k1].FathersName = make([][]string, maxLen)
		VAccounts[k1].Number = make([][]uint, maxLen)
		VAccounts[k1].Levels = make([]uint, maxLen)
		for k2, v2 := range v1.Number {
			if k2 < maxLen {
				VAccounts[k1].Number[k2] = v2
				VAccounts[k1].Levels[k2] = uint(len(v2))
			}
		}

		_, isIn := FFind(v1.CostFlowType, VAllCostFlowType)
		if !isIn {
			VAccounts[k1].CostFlowType = CHighLevelAccount
		}
	}

	for k1 := 0; k1 < maxLen; k1++ {
		for k2, v2 := range VAccounts {
			if len(v2.Number[k1]) > 1 {
				for _, v3 := range VAccounts {
					if len(v3.Number[k1]) > 0 {
						if FIsItSubAccountUsingNumber(v3.Number[k1], v2.Number[k1]) {
							VAccounts[k2].FathersName[k1] = append(VAccounts[k2].FathersName[k1], v3.Name)
						}
					}
				}
			}
		}
	}

	for k1 := range VAccounts {
		if len(VAccounts[k1].Number) > 0 {
			for k2 := range VAccounts {
				if len(VAccounts[k2].Number) > 0 {
					if k1 < k2 && !FIsItHighThanByOrder(VAccounts[k1].Number[VIndexOfAccountNumber], VAccounts[k2].Number[VIndexOfAccountNumber]) {
						FSwap(VAccounts, k1, k2)
					}
				}
			}
		}
	}

	accountErrors := make([]SAccount3, len(VAccounts))
	for k1, v1 := range VAccounts {
		accountErrors[k1].Barcode = make([]TErr, len(v1.Barcode))
		accountErrors[k1].Number = make([]TErr, maxLen)
	}

	for k1 := 0; k1 < maxLen; k1++ {
		for k2, v2 := range VAccounts {
			lenV2Number := len(v2.Number[k1])
			if lenV2Number > 0 {
				var haveFather bool
				for k3, v3 := range VAccounts {
					if len(v3.Number[k1]) > 0 {
						if VAccounts[k2].CostFlowType != CHighLevelAccount && FIsItSubAccountUsingNumber(v2.Number[k1], v3.Number[k1]) {
							accountErrors[k2].Number[k1] = TErr("should be low level account in all columns number")
						}
						if k2 != k3 && reflect.DeepEqual(v2.Number[k1], v3.Number[k1]) {
							accountErrors[k2].Number[k1] = TErr("duplicated")
						}
						if FIsItTheFather(v3.Number[k1], v2.Number[k1]) {
							haveFather = true
						}
					}
				}
				if lenV2Number > 1 && !haveFather {
					accountErrors[k2].Number[k1] = TErr("this account don't have father")
				}
			}
		}
	}

	for k1, v1 := range VAccounts {
		for k2, v2 := range VAccounts {
			if k1 != k2 && v1.Inventory == v2.Inventory {
				for k3, v3 := range v1.Barcode {
					for _, v4 := range v2.Barcode {
						if v3 == v4 {
							accountErrors[k1].Barcode[k3] = TErr("duplicated")
						}
					}
				}
			}
		}
	}

	return accountErrors
}

func FAddAutoCompletion(a SAutoCompletion) error {
	account, isExist, err := FAccountTerms(a.AccountName, true, false)
	if isExist && err != nil {
		return err
	}
	accountCost, isExistCost, err := FAccountTerms(CPrefixCost+a.AccountName, true, false)
	if isExistCost && err != nil {
		return err
	}
	accountDiscount, isExistDiscount, err := FAccountTerms(CPrefixDiscount+a.AccountName, true, false)
	if isExistDiscount && err != nil {
		return err
	}
	accountTaxExpenses, isExistTaxExpenses, err := FAccountTerms(CPrefixTaxExpenses+a.AccountName, true, false)
	if isExistTaxExpenses && err != nil {
		return err
	}
	accountTaxLiability, isExistTaxLiability, err := FAccountTerms(CPrefixTaxLiability+a.AccountName, true, true)
	if isExistTaxLiability && err != nil {
		return err
	}
	accountRevenue, isExistRevenue, err := FAccountTerms(CPrefixRevenue+a.AccountName, true, true)
	if isExistRevenue && err != nil {
		return err
	}

	add := func(isExist bool, account SAccount1) {
		if !isExist {
			FAddAccount(true, account)
		}
	}

	add(isExist, account)
	add(isExistCost, accountCost)
	add(isExistDiscount, accountDiscount)
	add(isExistTaxExpenses, accountTaxExpenses)
	add(isExistTaxLiability, accountTaxLiability)
	add(isExistRevenue, accountRevenue)

	a.PriceRevenue = math.Abs(a.PriceRevenue)
	a.PriceTax = math.Abs(a.PriceTax)

	setDiscont := func(discountPrice float64) float64 {
		discountPrice = math.Abs(discountPrice)
		if discountPrice > a.PriceRevenue {
			discountPrice = a.PriceRevenue
		}
		return discountPrice
	}

	a.DiscountPerOne = setDiscont(a.DiscountPerOne)
	a.DiscountTotal = math.Abs(a.DiscountTotal)
	a.DiscountPerQuantity.TPrice = setDiscont(a.DiscountPerQuantity.TPrice)
	a.DiscountPerQuantity.TQuantity = math.Abs(a.DiscountPerQuantity.TQuantity)
	for k1, v1 := range a.DiscountDecisionTree {
		a.DiscountDecisionTree[k1].TPrice = setDiscont(v1.TPrice)
		a.DiscountDecisionTree[k1].TQuantity = math.Abs(v1.TQuantity)
	}

	if _, isIn := FFind(a.DiscountWay, VDiscountWay); !isIn {
		a.DiscountWay = CDiscountPerOne
	}

	FDbUpdate(VDbAutoCompletionEntries, []byte(a.AccountName), a)
	_, VAutoCompletionEntries = FDbRead[SAutoCompletion](VDbAutoCompletionEntries)

	return nil
}

func FAccountTerms(accountName string, isLowLevel, isCredit bool) (SAccount1, bool, error) {
	account, _, err := FFindAccountFromName(accountName)
	newAccount := SAccount1{CostFlowType: CFifo, IsCredit: isCredit, Name: accountName}

	if err != nil {
		return newAccount, false, err
	}

	if isLowLevel {
		if account.CostFlowType == CHighLevelAccount {
			return newAccount, true, fmt.Errorf("(%v) should be Low Level account", accountName)
		}
	} else {
		if account.CostFlowType != CHighLevelAccount {
			return newAccount, true, fmt.Errorf("(%v) should not be Low Level account", accountName)
		}
	}

	if isCredit {
		if !account.IsCredit {
			return newAccount, true, fmt.Errorf("(%v) should be credit account", accountName)
		}
	} else {
		if account.IsCredit {
			return newAccount, true, fmt.Errorf("(%v) should not be credit account", accountName)
		}
	}

	return newAccount, true, nil
}
