package main

import (
	"fmt"
	"os"
	"reflect"
	"text/tabwriter"
)

func FFindAccountFromNameOrBarcode(nameOrBarcode string) (SAccount, int, error) {
	for k1, v1 := range VAccounts {
		if nameOrBarcode == v1.TAccountName {
			return v1, k1, nil
		}
		_, isIn := FFind(nameOrBarcode, v1.TAccountBarcode)
		if isIn {
			return v1, k1, nil
		}
	}
	return SAccount{}, 0, VErrorNotListed
}

func FFindAccountFromName(accountName string) (SAccount, int, error) {
	for k1, v1 := range VAccounts {
		if accountName == v1.TAccountName {
			return v1, k1, nil
		}
	}
	return SAccount{}, 0, VErrorNotListed
}

func FFindAutoCompletionFromName(accountName string) (SAutoCompletion, int, error) {
	for k1, v1 := range VAutoCompletionEntries {
		if accountName == v1.TAccountName {
			return v1, k1, nil
		}
	}
	return SAutoCompletion{}, 0, VErrorNotListed
}

func FAddAccount(account SAccount) error {
	_, _, err := FFindAccountFromName(account.TAccountName)
	if err == nil {
		return VErrorAccountNameIsUsed
	}
	if FIsBarcodesUsed(account.TAccountBarcode) {
		return VErrorBarcodeIsUsed
	}

	VAccounts = append(VAccounts, account)
	FSetTheAccounts()
	FDbInsertIntoAccounts()
	return nil
}

func FCheckIfAccountNumberDuplicated() []error {
	var errors []error
	maxLen := FMaxLenForAccountNumber()
	for k1 := 0; k1 < maxLen; k1++ {
	big_loop:
		for k2, v2 := range VAccounts {
			if len(v2.TAccountNumber[k1]) > 0 {
				for k3, v3 := range VAccounts {
					if k2 != k3 && reflect.DeepEqual(v2.TAccountNumber[k1], v3.TAccountNumber[k1]) {
						errors = append(errors, fmt.Errorf("the account number %v for %v is duplicated", v2.TAccountNumber[k1], v2))
						continue big_loop
					}
				}
			}
		}
	}
	errors, _ = FReturnSetAndDuplicatesSlices(errors)
	return errors
}

func FCheckIfLowLevelAccountForAll() []error {
	var errors []error
	maxLen := FMaxLenForAccountNumber()
	for k1 := 1; k1 < maxLen; k1++ {
	big_loop:
		for k2, v2 := range VAccounts {
			if len(v2.TAccountNumber[k1]) > 0 {
				for _, v3 := range VAccounts {
					if len(v3.TAccountNumber[k1]) > 0 {
						if FIsItSubAccountUsingNumber(v2.TAccountNumber[k1], v3.TAccountNumber[k1]) {
							continue big_loop
						}
					}
				}
				if VAccounts[k2].TCostFlowType != CHighLevelAccount {
					errors = append(errors, fmt.Errorf("should be low level account in all account numbers %v", VAccounts[k2]))
				}
			}
		}
	}
	return errors
}

func FCheckIfTheTreeConnected() []error {
	var errors []error
	maxLen := FMaxLenForAccountNumber()
	for k1 := 0; k1 < maxLen; k1++ {
	big_loop:
		for _, v2 := range VAccounts {
			if len(v2.TAccountNumber[k1]) > 1 {
				for _, v3 := range VAccounts {
					if FIsItTheFather(v3.TAccountNumber[k1], v2.TAccountNumber[k1]) {
						continue big_loop
					}
				}
				errors = append(errors, fmt.Errorf("the account number %v for %v not conected to the tree", v2.TAccountNumber[k1], v2))
			}
		}
	}
	return errors
}

func FCheckTheTree() []error {
	var errorsMessages []error
	errorsMessages = append(errorsMessages, FCheckIfLowLevelAccountForAll()...)
	errorsMessages = append(errorsMessages, FCheckIfAccountNumberDuplicated()...)
	errorsMessages = append(errorsMessages, FCheckIfTheTreeConnected()...)
	return errorsMessages
}

func FEditAccount(isDelete bool, index int, account SAccount) {
	newAccountName := account.TAccountName
	oldAccountName := VAccounts[index].TAccountName

	if !FIsUsedInJournal(oldAccountName) {
		if isDelete {
			VAccounts = FRemove(VAccounts, index)
			FSetTheAccounts()
			FDbInsertIntoAccounts()
			return
		}

		VAccounts[index].TCostFlowType = account.TCostFlowType
		VAccounts[index].TIsCredit = account.TIsCredit
	}

	if oldAccountName != newAccountName && newAccountName != "" {
		_, _, err := FFindAccountFromName(newAccountName)
		if err != nil {
			FChangeAccountName(oldAccountName, newAccountName)
			VAccounts[index].TAccountName = newAccountName
		}
	}

	if !FIsBarcodesUsed(account.TAccountBarcode) {
		VAccounts[index].TAccountBarcode = account.TAccountBarcode
	}

	_, isIn := FFind(account.TCostFlowType, VCostFlowType)
	if isIn && VAccounts[index].TCostFlowType != CHighLevelAccount {
		VAccounts[index].TCostFlowType = account.TCostFlowType
	}
	VAccounts[index].TAccountNotes = account.TAccountNotes
	VAccounts[index].TAccountImage = account.TAccountImage
	VAccounts[index].TAccountNumber = account.TAccountNumber

	FSetTheAccounts()
	FDbInsertIntoAccounts()
}

func FFormatSliceOfSliceOfStringToString(a [][]string) string {
	var str string
	for _, v1 := range a {
		str += "{"
		for _, v2 := range v1 {
			str += "\"" + v2 + "\","
		}
		str += "}\t,"
	}
	return "[][]string{" + str + "}"
}

func FFormatSliceOfSliceOfUintToString(a [][]uint) string {
	var str string
	for _, v1 := range a {
		str += "{"
		for _, v2 := range v1 {
			str += fmt.Sprint(v2) + ","
		}
		str += "}\t,"
	}
	return "[][]uint{" + str + "}"
}

func FFormatSliceOfUintToString(a []uint) string {
	var str string
	for _, v1 := range a {
		str += fmt.Sprint(v1) + ","
	}
	return "[]uint{" + str + "}"
}

func FFormatStringSliceToString(a []string) string {
	var str string
	for _, v1 := range a {
		str += "\"" + v1 + "\","
	}
	return "[]string{" + str + "}"
}

func FIsBarcodesUsed(barcode []string) bool {
	for _, v1 := range VAccounts {
		for _, v2 := range barcode {
			_, isIn := FFind(v2, v1.TAccountBarcode)
			if isIn {
				return true
			}
		}
	}
	return false
}

func FIsItHighThanByOrder(accountNumber1, accountNumber2 []uint) bool {
	l1 := len(accountNumber1)
	l2 := len(accountNumber2)
	for k1 := 0; k1 < FSmallest(l1, l2); k1++ {
		if accountNumber1[k1] < accountNumber2[k1] {
			return true
		} else if accountNumber1[k1] > accountNumber2[k1] {
			return false
		}
	}
	return l2 > l1
}

func FIsItPossibleToBeSubAccount(higherLevelAccountNumber, lowerLevelAccountNumber []uint) bool {
	lenHigherLevelAccountNumber := len(higherLevelAccountNumber)
	lenLowerLevelAccountNumber := len(lowerLevelAccountNumber)
	if lenHigherLevelAccountNumber == 0 || lenLowerLevelAccountNumber == 0 {
		return false
	}
	if lenHigherLevelAccountNumber >= lenLowerLevelAccountNumber {
		return false
	}
	if reflect.DeepEqual(higherLevelAccountNumber, lowerLevelAccountNumber) {
		return false
	}
	return true
}

func FIsItSubAccountUsingName(higherLevelAccount, lowerLevelAccount string) bool {
	a1, _, _ := FFindAccountFromName(higherLevelAccount)
	a2, _, _ := FFindAccountFromName(lowerLevelAccount)
	return FIsItSubAccountUsingNumber(a1.TAccountNumber[VIndexOfAccountNumber], a2.TAccountNumber[VIndexOfAccountNumber])
}

func FIsItSubAccountUsingNumber(higherLevelAccountNumber, lowerLevelAccountNumber []uint) bool {
	if !FIsItPossibleToBeSubAccount(higherLevelAccountNumber, lowerLevelAccountNumber) {
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
	if !FIsItPossibleToBeSubAccount(higherLevelAccountNumber, lowerLevelAccountNumber) {
		return false
	}
	return reflect.DeepEqual(higherLevelAccountNumber, FCutTheSlice(lowerLevelAccountNumber, 1))
}

func FIsUsedInJournal(accountName string) bool {
	_, journal := FDbRead[SJournal](VDbJournal)
	for _, v1 := range journal {
		if accountName == v1.AccountCredit || accountName == v1.AccountDebit {
			return true
		}
	}
	return false
}

func FMaxLenForAccountNumber() int {
	var maxLen int
	for _, v1 := range VAccounts {
		var length int
		for _, v2 := range v1.TAccountNumber {
			if len(v2) > 0 {
				length++
			}
		}
		if length > maxLen {
			maxLen = length
		}
	}
	return maxLen
}

func FPrintFormatedAccounts() {
	p := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	for _, v1 := range VAccounts {
		isCredit := fmt.Sprint(v1.TIsCredit)
		costFlowType := "\t,\"" + v1.TCostFlowType + "\""
		accountName := "\t,\"" + v1.TAccountName + "\""
		notes := "\t,\"" + v1.TAccountNotes + "\""
		image := "\t," + FFormatStringSliceToString(v1.TAccountImage)
		barcodes := "\t," + FFormatStringSliceToString(v1.TAccountBarcode)
		accountNumber := "\t," + FFormatSliceOfSliceOfUintToString(v1.TAccountNumber)
		accountLevels := "\t," + FFormatSliceOfUintToString(v1.TAccountLevels)
		fathersAccountsName := "\t," + FFormatSliceOfSliceOfStringToString(v1.TAccountFathersName)
		fmt.Fprintln(p, "{", isCredit, costFlowType, accountName, notes,
			image, barcodes, accountNumber, accountLevels, fathersAccountsName, "},")
	}
	p.Flush()
}

func FSetTheAccounts() {
	maxLen := FMaxLenForAccountNumber()

	for k1, v1 := range VAccounts {
		VAccounts[k1].TAccountFathersName = make([][]string, maxLen)
		VAccounts[k1].TAccountNumber = make([][]uint, maxLen)
		VAccounts[k1].TAccountLevels = make([]uint, maxLen)
		for k2, v2 := range v1.TAccountNumber {
			if k2 < maxLen {
				VAccounts[k1].TAccountNumber[k2] = v2
				VAccounts[k1].TAccountLevels[k2] = uint(len(v2))
			}
		}

		_, isIn := FFind(v1.TCostFlowType, VCostFlowType)
		if !isIn {
			VAccounts[k1].TCostFlowType = CHighLevelAccount
		}
	}

	for k1 := 0; k1 < maxLen; k1++ {
		for k2, v2 := range VAccounts {
			if len(v2.TAccountNumber[k1]) > 1 {
				for _, v3 := range VAccounts {
					if len(v3.TAccountNumber[k1]) > 0 {
						if FIsItSubAccountUsingNumber(v3.TAccountNumber[k1], v2.TAccountNumber[k1]) {
							VAccounts[k2].TAccountFathersName[k1] = append(VAccounts[k2].TAccountFathersName[k1], v3.TAccountName)
						}
					}
				}
			}
		}
	}

	for k1 := range VAccounts {
		for k2 := range VAccounts {
			if k1 < k2 && !FIsItHighThanByOrder(VAccounts[k1].TAccountNumber[VIndexOfAccountNumber], VAccounts[k2].TAccountNumber[VIndexOfAccountNumber]) {
				FSwap(VAccounts, k1, k2)
			}
		}
	}
}

func FAddAutoCompletion(a SAutoCompletion) error {
	account, isExist, err := FAccountTerms(a.TAccountName, true, false)
	if isExist && err != nil {
		return err
	}
	accountCost, isExistCost, err := FAccountTerms(CPrefixCost+a.TAccountName, true, false)
	if isExistCost && err != nil {
		return err
	}
	accountDiscount, isExistDiscount, err := FAccountTerms(CPrefixDiscount+a.TAccountName, true, false)
	if isExistDiscount && err != nil {
		return err
	}
	accountTaxExpenses, isExistTaxExpenses, err := FAccountTerms(CPrefixTaxExpenses+a.TAccountName, true, false)
	if isExistTaxExpenses && err != nil {
		return err
	}
	accountTaxLiability, isExistTaxLiability, err := FAccountTerms(CPrefixTaxLiability+a.TAccountName, true, true)
	if isExistTaxLiability && err != nil {
		return err
	}
	accountRevenue, isExistRevenue, err := FAccountTerms(CPrefixRevenue+a.TAccountName, true, true)
	if isExistRevenue && err != nil {
		return err
	}

	if !isExist {
		FAddAccount(account)
	}
	if !isExistCost {
		FAddAccount(accountCost)
	}
	if !isExistDiscount {
		FAddAccount(accountDiscount)
	}
	if !isExistTaxExpenses {
		FAddAccount(accountTaxExpenses)
	}
	if !isExistTaxLiability {
		FAddAccount(accountTaxLiability)
	}
	if !isExistRevenue {
		FAddAccount(accountRevenue)
	}

	a.PriceRevenue = FAbs(a.PriceRevenue)
	a.PriceTax = FAbs(a.PriceTax)
	for k1, v1 := range a.PriceDiscount {
		a.PriceDiscount[k1].TPrice = FAbs(v1.TPrice)
		a.PriceDiscount[k1].TQuantity = FAbs(v1.TQuantity)
	}

	FDbUpdate(VDbAutoCompletionEntries, []byte(a.TAccountName), a)
	_, VAutoCompletionEntries = FDbRead[SAutoCompletion](VDbAutoCompletionEntries)

	return nil
}

func FAccountTerms(accountName string, isLowLevel, isCredit bool) (SAccount, bool, error) {
	account, _, err := FFindAccountFromName(accountName)
	newAccount := SAccount{TCostFlowType: CFifo, TIsCredit: isCredit, TAccountName: accountName}

	if err != nil {
		return newAccount, false, err
	}

	if isLowLevel {
		if account.TCostFlowType == "" {
			return newAccount, true, fmt.Errorf("(%v) should be Low Level account", account.TAccountName)
		}
	} else {
		if account.TCostFlowType != "" {
			return newAccount, true, fmt.Errorf("(%v) should not be Low Level account", account.TAccountName)
		}
	}

	if isCredit {
		if !account.TIsCredit {
			return newAccount, true, fmt.Errorf("(%v) should be credit account", account.TAccountName)
		}
	} else {
		if account.TIsCredit {
			return newAccount, true, fmt.Errorf("(%v) should not be credit account", account.TAccountName)
		}
	}

	return newAccount, true, nil
}
