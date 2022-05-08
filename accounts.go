package main

import (
	"fmt"
	"os"
	"reflect"
	"text/tabwriter"
)

func FindAccountFromBarcode(barcode string) (Account, int, error) {
	for k1, v1 := range Accounts {
		if IsIn(barcode, v1.Barcode) {
			return v1, k1, nil
		}
	}
	return Account{}, 0, ErrorNotListed
}

func FindAccountFromName(accountName string) (Account, int, error) {
	for k1, v1 := range Accounts {
		if accountName == v1.Name {
			return v1, k1, nil
		}
	}
	return Account{}, 0, ErrorNotListed
}

func FindAutoCompletionFromName(accountName string) (AutoCompletion, int, error) {
	for k1, v1 := range AutoCompletionEntries {
		if accountName == v1.AccountInvnetory {
			return v1, k1, nil
		}
	}
	return AutoCompletion{}, 0, ErrorNotListed
}

func AddAccount(account Account) error {
	account.Name = FormatTheString(account.Name)
	if account.Name == "" {
		return ErrorAccountNameIsEmpty
	}
	_, _, err := FindAccountFromName(account.Name)
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
			if len(v2.Number[k1]) > 0 {
				for k3, v3 := range Accounts {
					if k2 != k3 && reflect.DeepEqual(v2.Number[k1], v3.Number[k1]) {
						errors = append(errors, fmt.Errorf("the account number %v for %v is duplicated", v2.Number[k1], v2))
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
			if len(v2.Number[k1]) > 0 {
				for _, v3 := range Accounts {
					if len(v3.Number[k1]) > 0 {
						if IsItSubAccountUsingNumber(v2.Number[k1], v3.Number[k1]) {
							continue big_loop
						}
					}
				}
				if !Accounts[k2].IsLowLevel {
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
			if len(v2.Number[k1]) > 1 {
				for _, v3 := range Accounts {
					if IsItTheFather(v3.Number[k1], v2.Number[k1]) {
						continue big_loop
					}
				}
				errors = append(errors, fmt.Errorf("the account number %v for %v not conected to the tree", v2.Number[k1], v2))
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
	newAccountName := FormatTheString(account.Name)
	oldAccountName := Accounts[index].Name

	// here i will search for oldAccountName in journal if not used i can delete it or chenge it
	if !IsUsedInJournal(oldAccountName) {
		if isDelete {
			Accounts = Remove(Accounts, index)
			SetTheAccounts()
			DbInsertIntoAccounts()
			return
		}

		Accounts[index].IsLowLevel = account.IsLowLevel
		Accounts[index].IsCredit = account.IsCredit
	}

	if oldAccountName != newAccountName && newAccountName != "" {
		// if the account not used in journal then the account is not used in inventory then
		// i will search for the account newAccountName in accounts database if it is not used then i can chenge the name
		_, _, err := FindAccountFromName(newAccountName)
		if err != nil {
			ChangeAccountName(oldAccountName, newAccountName)
			Accounts[index].Name = newAccountName
		}
	}

	if !IsBarcodesUsed(account.Barcode) {
		Accounts[index].Barcode = account.Barcode
	}

	Accounts[index].CostFlowType = account.CostFlowType
	Accounts[index].Notes = account.Notes
	Accounts[index].Image = account.Image
	Accounts[index].Number = account.Number

	SetTheAccounts()
	DbInsertIntoAccounts()
}

func FormatSliceOfSliceOfStringToString(a [][]string) string {
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

func FormatSliceOfSliceOfUintToString(a [][]uint) string {
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

func FormatSliceOfUintToString(a []uint) string {
	var str string
	for _, v1 := range a {
		str += fmt.Sprint(v1) + ","
	}
	return "[]uint{" + str + "}"
}

func FormatStringSliceToString(a []string) string {
	var str string
	for _, v1 := range a {
		str += "\"" + v1 + "\","
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
	for k1 := 0; k1 < Smallest(l1, l2); k1++ {
		if accountNumber1[k1] < accountNumber2[k1] {
			return true
		} else if accountNumber1[k1] > accountNumber2[k1] {
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
	a1, _, _ := FindAccountFromName(higherLevelAccount)
	a2, _, _ := FindAccountFromName(lowerLevelAccount)
	return IsItSubAccountUsingNumber(a1.Number[IndexOfAccountNumber], a2.Number[IndexOfAccountNumber])
}

func IsItSubAccountUsingNumber(higherLevelAccountNumber, lowerLevelAccountNumber []uint) bool {
	if !IsItPossibleToBeSubAccount(higherLevelAccountNumber, lowerLevelAccountNumber) {
		return false
	}
	for k1, v1 := range higherLevelAccountNumber {
		if v1 != lowerLevelAccountNumber[k1] {
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
	for _, v1 := range journal {
		if accountName == v1.AccountCredit || accountName == v1.AccountDebit {
			return true
		}
	}
	return false
}

func MaxLenForAccountNumber() int {
	var maxLen int
	for _, v1 := range Accounts {
		var length int
		for _, v2 := range v1.Number {
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

func PrintFormatedAccounts() {
	p := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	for _, v1 := range Accounts {
		isLowLevelAccount := v1.IsLowLevel
		isCredit := "\t," + fmt.Sprint(v1.IsCredit)
		costFlowType := "\t,\"" + v1.CostFlowType + "\""
		accountName := "\t,\"" + v1.Name + "\""
		notes := "\t,\"" + v1.Notes + "\""
		image := "\t," + FormatStringSliceToString(v1.Image)
		barcodes := "\t," + FormatStringSliceToString(v1.Barcode)
		accountNumber := "\t," + FormatSliceOfSliceOfUintToString(v1.Number)
		accountLevels := "\t," + FormatSliceOfUintToString(v1.Levels)
		fathersAccountsName := "\t," + FormatSliceOfSliceOfStringToString(v1.FathersName)
		fmt.Fprintln(p, "{", isLowLevelAccount, isCredit, costFlowType, accountName, notes,
			image, barcodes, accountNumber, accountLevels, fathersAccountsName, "},")
	}
	p.Flush()
}

func SetTheAccounts() {
	maxLen := MaxLenForAccountNumber()

	for k1, v1 := range Accounts {
		// init the slices
		Accounts[k1].FathersName = make([][]string, maxLen)
		Accounts[k1].Number = make([][]uint, maxLen)
		Accounts[k1].Levels = make([]uint, maxLen)
		for k2, v2 := range v1.Number {
			if k2 < maxLen {
				Accounts[k1].Number[k2] = v2
				Accounts[k1].Levels[k2] = uint(len(v2))
			}
		}

		// set high level account to permanent
		// set cost flow type . the cost flow should be used for every low level account
		if !v1.IsLowLevel {
			Accounts[k1].CostFlowType = ""
		} else if !IsIn(v1.CostFlowType, CostFlowType) {
			Accounts[k1].CostFlowType = Fifo
		}
	}

	// here i set the father and grandpa accounts name
	for k1 := 0; k1 < maxLen; k1++ {
		for k2, v2 := range Accounts { // here i loop over account
			if len(v2.Number[k1]) > 1 {
				for _, v3 := range Accounts { // but here i loop over account to find the father or grandpa account
					if len(v3.Number[k1]) > 0 {
						if IsItSubAccountUsingNumber(v3.Number[k1], v2.Number[k1]) {
							Accounts[k2].FathersName[k1] = append(Accounts[k2].FathersName[k1], v3.Name)
						}
					}
				}
			}
		}
	}

	// here i sort the accounts by there account number
	for k1 := range Accounts {
		for k2 := range Accounts {
			if k1 < k2 && !IsItHighThanByOrder(Accounts[k1].Number[IndexOfAccountNumber], Accounts[k2].Number[IndexOfAccountNumber]) {
				Swap(Accounts, k1, k2)
			}
		}
	}
}

func AddAutoCompletion(a AutoCompletion) error {
	account, isExist, err := AccountTerms(a.AccountInvnetory, true, false)
	if isExist && err != nil {
		return err
	}
	accountCost, isExistCost, err := AccountTerms(PrefixCost+a.AccountInvnetory, true, false)
	if isExistCost && err != nil {
		return err
	}
	accountDiscount, isExistDiscount, err := AccountTerms(PrefixDiscount+a.AccountInvnetory, true, false)
	if isExistDiscount && err != nil {
		return err
	}
	accountTaxExpenses, isExistTaxExpenses, err := AccountTerms(PrefixTaxExpenses+a.AccountInvnetory, true, false)
	if isExistTaxExpenses && err != nil {
		return err
	}
	accountTaxLiability, isExistTaxLiability, err := AccountTerms(PrefixTaxLiability+a.AccountInvnetory, true, true)
	if isExistTaxLiability && err != nil {
		return err
	}
	accountRevenue, isExistRevenue, err := AccountTerms(PrefixRevenue+a.AccountInvnetory, true, true)
	if isExistRevenue && err != nil {
		return err
	}

	if !isExist {
		AddAccount(account)
	}
	if !isExistCost {
		AddAccount(accountCost)
	}
	if !isExistDiscount {
		AddAccount(accountDiscount)
	}
	if !isExistTaxExpenses {
		AddAccount(accountTaxExpenses)
	}
	if !isExistTaxLiability {
		AddAccount(accountTaxLiability)
	}
	if !isExistRevenue {
		AddAccount(accountRevenue)
	}

	a.PriceRevenue = Abs(a.PriceRevenue)
	a.PriceTax = Abs(a.PriceTax)
	for k1, v1 := range a.PriceDiscount {
		a.PriceDiscount[k1].Price = Abs(v1.Price)
		a.PriceDiscount[k1].Quantity = Abs(v1.Quantity)
	}

	DbUpdate(DbAutoCompletionEntries, []byte(a.AccountInvnetory), a)
	_, AutoCompletionEntries = DbRead[AutoCompletion](DbAutoCompletionEntries)

	return nil
}

func AccountTerms(accountName string, isLowLevel, isCredit bool) (Account, bool, error) {
	accountName = FormatTheString(accountName)
	account, _, err := FindAccountFromName(accountName)
	newAccount := Account{IsLowLevel: isLowLevel, IsCredit: isCredit, Name: accountName}

	if err != nil {
		return newAccount, false, err
	}

	if isLowLevel {
		if !account.IsLowLevel {
			return newAccount, true, fmt.Errorf("(%v) should be Low Level account", account.Name)
		}
	} else {
		if account.IsLowLevel {
			return newAccount, true, fmt.Errorf("(%v) should not be Low Level account", account.Name)
		}
	}

	if isCredit {
		if !account.IsCredit {
			return newAccount, true, fmt.Errorf("(%v) should be credit account", account.Name)
		}
	} else {
		if account.IsCredit {
			return newAccount, true, fmt.Errorf("(%v) should not be credit account", account.Name)
		}
	}

	return newAccount, true, nil
}
