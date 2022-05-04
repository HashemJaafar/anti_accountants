package main

import (
	"errors"
	"fmt"
	"testing"
)

func TestAccountStructFromBarcode(t *testing.T) {
	IndexOfAccountNumber = 0
	accountStruct, index, err := AccountStructFromBarcode("kaslajs")
	Test(true, err, nil)
	Test(true, index, 1)
	Test(true, accountStruct, Account{false, false, false, "", "CURRENT_ASSETS", "", []string{}, []string{"sijadpodjpao", "kaslajs"}, [][]uint{{1, 1}, {}}, []uint{2, 0}, [][]string{{"ASSETS"}, {}}, 0, 0, 0, false})
}

func TestAccountStructFromName(t *testing.T) {
	IndexOfAccountNumber = 0
	accountStruct, index, err := AccountStructFromName("ASSETS")
	Test(true, err, nil)
	Test(true, index, 0)
	Test(true, accountStruct, Account{false, false, false, "", "ASSETS", "", []string{}, []string{"nojdsjdpq"}, [][]uint{{1}, {}}, []uint{1, 0}, [][]string{{}, {}}, 0, 0, 0, false})
}

func TestAddAccount(t *testing.T) {
	a := AddAccount(Account{
		IsLowLevelAccount:                       true,
		IsCredit:                                true,
		IsTemporary:                             true,
		CostFlowType:                            Wma,
		AccountName:                             RetinedEarnings,
		Notes:                                   "",
		Image:                                   []string{},
		Barcode:                                 []string{},
		AccountNumber:                           [][]uint{{3, 1}},
		AccountLevels:                           []uint{},
		FathersAccountsName:                     [][]string{},
		AlertForMinimumQuantityByTurnoverInDays: 0,
		AlertForMinimumQuantityByQuintity:       0,
		TargetBalance:                           0,
		IfTheTargetBalanceIsLessIsGood:          false,
	})
	DbClose()
	PrintFormatedAccounts()
	Test(true, a, nil)
}

func TestCheckIfAccountNumberDuplicated(t *testing.T) {
	a := CheckIfAccountNumberDuplicated()
	fmt.Println(a)
	Test(true, a, []error{errors.New("the account number [2] for {false false true Fifo e  [] [] [[4] [2]] [] [] 0 0 0 false} duplicated")})
}

func TestCheckIfLowLevelAccountForAll(t *testing.T) {
	a := CheckIfLowLevelAccountForAll()
	Test(true, a, []error{errors.New("should be low level account in all account numbers {false false false  b  [] [] [[1 1] [1 2]] [] [] 0 0 0 false}")})
}

func TestCheckIfTheTreeConnected(t *testing.T) {
	a := CheckIfTheTreeConnected()
	fmt.Println(a)
	Test(true, a, []error{errors.New("the account number [2 1 8] for {true false true Fifo f  [] [] [[4 1] [2 1 8]] [] [] 0 0 0 false} not conected to the tree")})
}

func TestCheckTheTree(t *testing.T) {
	CheckTheTree()
	Test(true, ErrorsMessages, []error{
		errors.New("should be low level account in all account numbers {false false false  b  [] [] [[1 1] [1 2]] [] [] 0 0 0 false}"),
		errors.New("the account number [2 1 8] for {true false true Fifo f  [] [] [[4 1] [2 1 8]] [] [] 0 0 0 false} not conected to the tree"),
	})
}

func TestEditAccount(t *testing.T) {
	accountStruct, index, err := AccountStructFromName("cash")
	fmt.Println(err)
	if err == nil {
		accountStruct.CostFlowType = Wma
		EditAccount(false, index, accountStruct)
	}
	DbClose()
	PrintFormatedAccounts()
}

func TestIsBarcodesUsed(t *testing.T) {
	a := IsBarcodesUsed([]string{"a", "b"})
	Test(true, a, true)
	a = IsBarcodesUsed([]string{"c", "b"})
	Test(true, a, false)
}

func TestIsItHighThanByOrder(t *testing.T) {
	a := IsItHighThanByOrder([]uint{1}, []uint{1, 2})
	Test(true, a, true)
	a = IsItHighThanByOrder([]uint{1}, []uint{1})
	Test(true, a, false)
	a = IsItHighThanByOrder([]uint{1, 2}, []uint{1})
	Test(true, a, false)
	a = IsItHighThanByOrder([]uint{3}, []uint{1, 1})
	Test(true, a, false)
	a = IsItHighThanByOrder([]uint{1, 5}, []uint{3})
	Test(true, a, true)
	a = IsItHighThanByOrder([]uint{1, 1}, []uint{1, 2})
	Test(true, a, true)
	a = IsItHighThanByOrder([]uint{4}, []uint{1, 2})
	Test(true, a, false)
}

func TestIsItPossibleToBeSubAccount(t *testing.T) {
	a := IsItPossibleToBeSubAccount([]uint{1}, []uint{1, 2})
	Test(true, a, true)
	a = IsItPossibleToBeSubAccount([]uint{1}, []uint{2})
	Test(true, a, false)
	a = IsItPossibleToBeSubAccount([]uint{}, []uint{2})
	Test(true, a, false)
	a = IsItPossibleToBeSubAccount([]uint{1}, []uint{1, 1, 2})
	Test(true, a, true)
}

func TestIsItSubAccountUsingName(t *testing.T) {
	IndexOfAccountNumber = 0
	a := IsItSubAccountUsingName("ASSETS", "CASH_AND_CASH_EQUIVALENTS")
	Test(true, a, true)
	a = IsItSubAccountUsingName("CASH_AND_CASH_EQUIVALENTS", "ASSETS")
	Test(true, a, false)
	IndexOfAccountNumber = 1
	a = IsItSubAccountUsingName("ASSETS", "CASH_AND_CASH_EQUIVALENTS")
	Test(true, a, false)
}

func TestIsItSubAccountUsingNumber(t *testing.T) {
	a := IsItSubAccountUsingNumber([]uint{1}, []uint{1, 2})
	Test(true, a, true)
	a = IsItSubAccountUsingNumber([]uint{1}, []uint{2})
	Test(true, a, false)
	a = IsItSubAccountUsingNumber([]uint{}, []uint{2})
	Test(true, a, false)
}

func TestIsItTheFather(t *testing.T) {
	a := IsItTheFather([]uint{1}, []uint{1, 2})
	Test(true, a, true)
	a = IsItTheFather([]uint{1}, []uint{2})
	Test(true, a, false)
	a = IsItTheFather([]uint{}, []uint{2})
	Test(true, a, false)
	a = IsItTheFather([]uint{1}, []uint{1, 1, 2})
	Test(true, a, false)
}

func TestIsUsedInJournal(t *testing.T) {
	a := IsUsedInJournal("book")
	Test(true, a, false)
}

func TestMaxLenForAccountNumber(t *testing.T) {
	a := MaxLenForAccountNumber()
	Test(true, a, 2)
}

func TestPrintFormatedAccounts(t *testing.T) {
	PrintFormatedAccounts()
}

func TestSetTheAccounts(t *testing.T) {
	SetTheAccounts()
	PrintFormatedAccounts()
	Test(true, Accounts, []Account{
		// {false, false, false, "", "a", "", []string{}, []string{"a"}, [][]uint{{1}, {1}}, []uint{1, 1}, [][]string{{}, {}}, 0, 0, 0, false},
		// {false, false, false, "", "b", "", []string{}, []string{}, [][]uint{{1, 1}, {1, 2}}, []uint{2, 2}, [][]string{{"a"}, {"a"}}, 0, 0, 0, false},
		// {true, false, false, "Fifo", "c", "", []string{}, []string{}, [][]uint{{1, 2}, {1, 3}}, []uint{2, 2}, [][]string{{"a"}, {"a"}}, 0, 0, 0, false},
		// {false, false, false, "", "d", "", []string{}, []string{}, [][]uint{{2}, {2}}, []uint{1, 1}, [][]string{{}, {}}, 0, 0, 0, false},
		// {false, false, false, "", "e", "", []string{}, []string{}, [][]uint{{4}, {2}}, []uint{1, 1}, [][]string{{}, {}}, 0, 0, 0, false},
		// {true, false, true, "", "f", "", []string{}, []string{}, [][]uint{{4, 1}, {2, 1, 8}}, []uint{2, 3}, [][]string{{"e"}, {"d", "e"}}, 0, 0, 0, false},
	})
}
