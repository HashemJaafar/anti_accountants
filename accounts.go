package main

import (
	"fmt"
	"os"
	"reflect"
	"text/tabwriter"
)

func AccountStructFromBarcode(barcode string) (Account, int, error) {
	for k1, v1 := range Accounts {
		if IsIn(barcode, v1.Barcode) {
			return v1, k1, nil
		}
	}
	return Account{}, 0, ErrorNotListed
}

func AccountStructFromName(accountName string) (Account, int, error) {
	for k1, v1 := range Accounts {
		if v1.AccountName == accountName {
			return v1, k1, nil
		}
	}
	return Account{}, 0, ErrorNotListed
}

func AddAccount(account Account) error {
	account.AccountName = FormatTheString(account.AccountName)
	if account.AccountName == "" {
		return ErrorAccountNameIsEmpty
	}
	_, _, err := AccountStructFromName(account.AccountName)
	if err == nil {
		return ErrorAccountNameIsUsed
	}
	if IsBarcodesUsed(account.Barcode) {
		return ErrorBarcodeIsUsed
	}

	Accounts = append(Accounts, account)
	SetTheAccounts()
	DbInsertIntoAccounts()
	return nil
}

func CheckIfAccountNumberDuplicated() []error {
	var errors []error
	maxLen := MaxLenForAccountNumber()
	for k1 := 0; k1 < maxLen; k1++ {
	big_loop:
		for k2, v2 := range Accounts {
			if len(v2.AccountNumber[k1]) > 0 {
				for indexb, b := range Accounts {
					if k2 != indexb && reflect.DeepEqual(v2.AccountNumber[k1], b.AccountNumber[k1]) {
						errors = append(errors, fmt.Errorf("the account number %v for %v is duplicated", v2.AccountNumber[k1], v2))
						continue big_loop
					}
				}
			}
		}
	}
	errors, _ = ReturnSetAndDuplicatesSlices(errors)
	return errors
}

func CheckIfLowLevelAccountForAll() []error {
	var errors []error
	maxLen := MaxLenForAccountNumber()
	for k1 := 1; k1 < maxLen; k1++ {
	big_loop:
		for k2, v2 := range Accounts {
			if len(v2.AccountNumber[k1]) > 0 {
				for _, v3 := range Accounts {
					if len(v3.AccountNumber[k1]) > 0 {
						if IsItSubAccountUsingNumber(v2.AccountNumber[k1], v3.AccountNumber[k1]) {
							continue big_loop
						}
					}
				}
				if !Accounts[k2].IsLowLevelAccount {
					errors = append(errors, fmt.Errorf("should be low level account in all account numbers %v", Accounts[k2]))
				}
			}
		}
	}
	return errors
}

func CheckIfTheTreeConnected() []error {
	var errors []error
	maxLen := MaxLenForAccountNumber()
	for k1 := 0; k1 < maxLen; k1++ {
	big_loop:
		for _, v2 := range Accounts {
			if len(v2.AccountNumber[k1]) > 1 {
				for _, v3 := range Accounts {
					if IsItTheFather(v3.AccountNumber[k1], v2.AccountNumber[k1]) {
						continue big_loop
					}
				}
				errors = append(errors, fmt.Errorf("the account number %v for %v not conected to the tree", v2.AccountNumber[k1], v2))
			}
		}
	}
	return errors
}

func CheckTheTree() []error {
	var errorsMessages []error
	errorsMessages = append(errorsMessages, CheckIfLowLevelAccountForAll()...)
	errorsMessages = append(errorsMessages, CheckIfAccountNumberDuplicated()...)
	errorsMessages = append(errorsMessages, CheckIfTheTreeConnected()...)
	return errorsMessages
}

func EditAccount(isDelete bool, index int, account Account) {
	newAccountName := FormatTheString(account.AccountName)
	oldAccountName := Accounts[index].AccountName

	// here i will search for oldAccountName in journal if not used i can delete it or chenge it
	if !IsUsedInJournal(oldAccountName) {
		if isDelete {
			Accounts = Remove(Accounts, index)
			SetTheAccounts()
			DbInsertIntoAccounts()
			return
		}

		Accounts[index].IsLowLevelAccount = account.IsLowLevelAccount
		Accounts[index].IsCredit = account.IsCredit
	}

	if oldAccountName != newAccountName && newAccountName != "" {
		// if the account not used in journal then the account is not used in inventory then
		// i will search for the account newAccountName in accounts database if it is not used then i can chenge the name
		_, _, err := AccountStructFromName(newAccountName)
		if err != nil {
			ChangeAccountName(oldAccountName, newAccountName)
			Accounts[index].AccountName = newAccountName
		}
	}

	if !IsBarcodesUsed(account.Barcode) {
		Accounts[index].Barcode = account.Barcode
	}

	Accounts[index].IsTemporary = account.IsTemporary
	Accounts[index].CostFlowType = account.CostFlowType
	Accounts[index].Notes = account.Notes
	Accounts[index].Image = account.Image
	Accounts[index].AccountNumber = account.AccountNumber
	Accounts[index].AlertForMinimumQuantityByTurnoverInDays = account.AlertForMinimumQuantityByTurnoverInDays
	Accounts[index].AlertForMinimumQuantityByQuintity = account.AlertForMinimumQuantityByQuintity
	Accounts[index].TargetBalance = account.TargetBalance
	Accounts[index].IfTheTargetBalanceIsLessIsGood = account.IfTheTargetBalanceIsLessIsGood

	SetTheAccounts()
	DbInsertIntoAccounts()
}

func FormatSliceOfSliceOfStringToString(a [][]string) string {
	var str string
	for _, b := range a {
		str += "{"
		for _, c := range b {
			str += "\"" + c + "\","
		}
		str += "}\t,"
	}
	return "[][]string{" + str + "}"
}

func FormatSliceOfSliceOfUintToString(a [][]uint) string {
	var str string
	for _, b := range a {
		str += "{"
		for _, c := range b {
			str += fmt.Sprint(c) + ","
		}
		str += "}\t,"
	}
	return "[][]uint{" + str + "}"
}

func FormatSliceOfUintToString(a []uint) string {
	var str string
	for _, b := range a {
		str += fmt.Sprint(b) + ","
	}
	return "[]uint{" + str + "}"
}

func FormatStringSliceToString(a []string) string {
	var str string
	for _, b := range a {
		str += "\"" + b + "\","
	}
	return "[]string{" + str + "}"
}

func IsBarcodesUsed(barcode []string) bool {
	for _, v1 := range Accounts {
		for _, v2 := range barcode {
			if IsIn(v2, v1.Barcode) {
				return true
			}
		}
	}
	return false
}

func IsItHighThanByOrder(accountNumber1, accountNumber2 []uint) bool {
	l1 := len(accountNumber1)
	l2 := len(accountNumber2)
	for index := 0; index < Smallest(l1, l2); index++ {
		if accountNumber1[index] < accountNumber2[index] {
			return true
		} else if accountNumber1[index] > accountNumber2[index] {
			return false
		}
	}
	return l2 > l1
}

func IsItPossibleToBeSubAccount(higherLevelAccountNumber, lowerLevelAccountNumber []uint) bool {
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

func IsItSubAccountUsingName(higherLevelAccount, lowerLevelAccount string) bool {
	a1, _, _ := AccountStructFromName(higherLevelAccount)
	a2, _, _ := AccountStructFromName(lowerLevelAccount)
	return IsItSubAccountUsingNumber(a1.AccountNumber[IndexOfAccountNumber], a2.AccountNumber[IndexOfAccountNumber])
}

func IsItSubAccountUsingNumber(higherLevelAccountNumber, lowerLevelAccountNumber []uint) bool {
	if !IsItPossibleToBeSubAccount(higherLevelAccountNumber, lowerLevelAccountNumber) {
		return false
	}
	for i, h := range higherLevelAccountNumber {
		if h != lowerLevelAccountNumber[i] {
			return false
		}
	}
	return true
}

func IsItTheFather(higherLevelAccountNumber, lowerLevelAccountNumber []uint) bool {
	if !IsItPossibleToBeSubAccount(higherLevelAccountNumber, lowerLevelAccountNumber) {
		return false
	}
	return reflect.DeepEqual(higherLevelAccountNumber, CutTheSlice(lowerLevelAccountNumber, 1))
}

func IsUsedInJournal(accountName string) bool {
	_, journal := DbRead[Journal](DbJournal)
	for _, i := range journal {
		if accountName == i.AccountCredit || accountName == i.AccountDebit {
			return true
		}
	}
	return false
}

func MaxLenForAccountNumber() int {
	var maxLen int
	for _, a := range Accounts {
		var length int
		for _, b := range a.AccountNumber {
			if len(b) > 0 {
				length++
			}
		}
		if length > maxLen {
			maxLen = length
		}
	}
	return maxLen
}

func PrintFormatedAccounts() {
	p := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	for _, a := range Accounts {
		isLowLevelAccount := a.IsLowLevelAccount
		isCredit := "\t," + fmt.Sprint(a.IsCredit)
		isTemporary := "\t," + fmt.Sprint(a.IsTemporary)
		costFlowType := "\t,\"" + a.CostFlowType + "\""
		accountName := "\t,\"" + a.AccountName + "\""
		notes := "\t,\"" + a.Notes + "\""
		image := "\t," + FormatStringSliceToString(a.Image)
		barcodes := "\t," + FormatStringSliceToString(a.Barcode)
		accountNumber := "\t," + FormatSliceOfSliceOfUintToString(a.AccountNumber)
		accountLevels := "\t," + FormatSliceOfUintToString(a.AccountLevels)
		fathersAccountsName := "\t," + FormatSliceOfSliceOfStringToString(a.FathersAccountsName)
		alertForMinimumQuantityByTurnoverInDays := "\t," + fmt.Sprint(a.AlertForMinimumQuantityByTurnoverInDays)
		alertForMinimumQuantityByQuintity := "\t," + fmt.Sprint(a.AlertForMinimumQuantityByQuintity)
		targetBalance := "\t," + fmt.Sprint(a.TargetBalance)
		ifTheTargetBalanceIsLessIsGood := "\t," + fmt.Sprint(a.IfTheTargetBalanceIsLessIsGood)
		fmt.Fprintln(p, "{", isLowLevelAccount, isCredit, isTemporary, costFlowType, accountName, notes,
			image, barcodes, accountNumber, accountLevels, fathersAccountsName,
			alertForMinimumQuantityByTurnoverInDays, alertForMinimumQuantityByQuintity, targetBalance, ifTheTargetBalanceIsLessIsGood, "},")
	}
	p.Flush()
}

func SetTheAccounts() {
	maxLen := MaxLenForAccountNumber()

	for k1, v1 := range Accounts {
		// init the slices
		Accounts[k1].FathersAccountsName = make([][]string, maxLen)
		Accounts[k1].AccountNumber = make([][]uint, maxLen)
		Accounts[k1].AccountLevels = make([]uint, maxLen)
		for k2, v2 := range v1.AccountNumber {
			if k2 < maxLen {
				Accounts[k1].AccountNumber[k2] = v2
				Accounts[k1].AccountLevels[k2] = uint(len(v2))
			}
		}

		// set high level account to permanent
		// set cost flow type . the cost flow should be used for every low level account
		if !v1.IsLowLevelAccount {
			Accounts[k1].IsTemporary = false
			Accounts[k1].CostFlowType = ""
		} else if !IsIn(v1.CostFlowType, CostFlowType) {
			Accounts[k1].CostFlowType = Fifo
		}
	}

	// here i set the father and grandpa accounts name
	for k1 := 0; k1 < maxLen; k1++ {
		for k2, v2 := range Accounts { // here i loop over account
			if len(v2.AccountNumber[k1]) > 1 {
				for _, v3 := range Accounts { // but here i loop over account to find the father or grandpa account
					if len(v3.AccountNumber[k1]) > 0 {
						if IsItSubAccountUsingNumber(v3.AccountNumber[k1], v2.AccountNumber[k1]) {
							Accounts[k2].FathersAccountsName[k1] = append(Accounts[k2].FathersAccountsName[k1], v3.AccountName)
						}
					}
				}
			}
		}
	}

	// here i sort the accounts by there account number
	for k1 := range Accounts {
		for k2 := range Accounts {
			if k1 < k2 && !IsItHighThanByOrder(Accounts[k1].AccountNumber[IndexOfAccountNumber], Accounts[k2].AccountNumber[IndexOfAccountNumber]) {
				Swap(Accounts, k1, k2)
			}
		}
	}
}

func SetRetainedEarningsAccount(account Account) Account {
	// in this function i fix the account field to the retained earnings account
	// just to know the RETAINED_EARNINGS is low level account but i dont want to use it in journal
	account.IsLowLevelAccount = true
	account.IsCredit = true
	account.IsTemporary = false
	return account
}
