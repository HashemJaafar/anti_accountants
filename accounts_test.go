package main

import (
	"errors"
	"fmt"
	"testing"
)

func TestAccountStructFromBarcode(t *testing.T) {
	IndexOfAccountNumber = 0
	account, index, err := FindAccountFromBarcode("kaslajs")
	Test(true, err, nil)
	Test(true, index, 1)
	Test(true, account, Account{false, false, "", "CURRENT_ASSETS", "", []string{}, []string{"sijadpodjpao", "kaslajs"}, [][]uint{{1, 1}, {}}, []uint{2, 0}, [][]string{{"ASSETS"}, {}}})
}

func TestAccountStructFromName(t *testing.T) {
	IndexOfAccountNumber = 0
	account, index, err := FindAccountFromName("ASSETS")
	Test(true, err, nil)
	Test(true, index, 0)
	Test(true, account, Account{false, false, "", "ASSETS", "", []string{}, []string{"nojdsjdpq"}, [][]uint{{1}, {}}, []uint{1, 0}, [][]string{{}, {}}})
}

func TestAddAccount(t *testing.T) {
	var a1 error

	a1 = AddAccount(Account{
		IsLowLevel:   false,
		IsCredit:     false,
		CostFlowType: Wma,
		Name:         "assets",
		Number:       [][]uint{{1}},
	})
	Test(true, a1, nil)

	a1 = AddAccount(Account{
		IsLowLevel:   false,
		IsCredit:     false,
		CostFlowType: Wma,
		Name:         "current assets",
		Number:       [][]uint{{1, 1}},
	})
	Test(true, a1, nil)

	a1 = AddAccount(Account{
		IsLowLevel:   true,
		IsCredit:     false,
		CostFlowType: Wma,
		Name:         "cash",
		Number:       [][]uint{{1, 1, 1}},
	})
	Test(true, a1, nil)

	a1 = AddAccount(Account{
		IsLowLevel:   false,
		IsCredit:     false,
		CostFlowType: Wma,
		Name:         "invnetory",
		Number:       [][]uint{{1, 1, 3}},
	})
	Test(true, a1, nil)

	a1 = AddAccount(Account{
		IsLowLevel:   true,
		IsCredit:     false,
		CostFlowType: Fifo,
		Name:         "book",
		Number:       [][]uint{{1, 1, 3, 1}},
	})
	Test(true, a1, nil)

	a1 = AddAccount(Account{
		IsLowLevel:   false,
		IsCredit:     true,
		CostFlowType: Wma,
		Name:         "liabilities",
		Number:       [][]uint{{2}},
	})
	Test(true, a1, nil)

	a1 = AddAccount(Account{
		IsLowLevel:   false,
		IsCredit:     true,
		CostFlowType: Wma,
		Name:         "owner's equity",
		Number:       [][]uint{{3}},
	})
	Test(true, a1, nil)

	a1 = AddAccount(Account{
		IsLowLevel:   true,
		IsCredit:     true,
		CostFlowType: Wma,
		Name:         "retained earnings",
		Number:       [][]uint{{3, 1}},
	})
	Test(true, a1, nil)

	a1 = AddAccount(Account{
		IsLowLevel:   false,
		IsCredit:     true,
		CostFlowType: Wma,
		Name:         "income",
		Number:       [][]uint{{3, 2}},
	})
	Test(true, a1, nil)

	a1 = AddAccount(Account{
		IsLowLevel:   false,
		IsCredit:     true,
		CostFlowType: Wma,
		Name:         "revenue",
		Number:       [][]uint{{3, 2, 1}},
	})
	Test(true, a1, nil)

	a1 = AddAccount(Account{
		IsLowLevel:   true,
		IsCredit:     true,
		CostFlowType: Wma,
		Name:         "rent",
		Number:       [][]uint{{3, 2, 1, 1}},
	})
	Test(true, a1, nil)

	a1 = AddAccount(Account{
		IsLowLevel:   false,
		IsCredit:     false,
		CostFlowType: Wma,
		Name:         "expense",
		Number:       [][]uint{{3, 2, 2}},
	})
	Test(true, a1, nil)

	a1 = AddAccount(Account{
		IsLowLevel:   false,
		IsCredit:     false,
		CostFlowType: Wma,
		Name:         "discounts",
		Number:       [][]uint{{3, 2, 2, 1}},
	})
	Test(true, a1, nil)

	a1 = AddAccount(Account{
		IsLowLevel:   true,
		IsCredit:     false,
		CostFlowType: Wma,
		Name:         InvoiceDiscount,
		Number:       [][]uint{{3, 2, 2, 1, 1}},
	})
	Test(true, a1, nil)

	a1 = AddAccount(Account{
		IsLowLevel:   true,
		IsCredit:     false,
		CostFlowType: Wma,
		Name:         PrefixDiscount + "book",
		Number:       [][]uint{{3, 2, 2, 1, 2}},
	})
	Test(true, a1, nil)

	a1 = AddAccount(Account{
		IsLowLevel:   false,
		IsCredit:     false,
		CostFlowType: Wma,
		Name:         "cost of goods sold",
		Number:       [][]uint{{3, 2, 2, 2}},
	})
	Test(true, a1, nil)

	a1 = AddAccount(Account{
		IsLowLevel:   true,
		IsCredit:     false,
		CostFlowType: Wma,
		Name:         PrefixCost + "book",
		Number:       [][]uint{{3, 2, 2, 2, 1}},
	})
	Test(true, a1, nil)

	DbClose()
	PrintFormatedAccounts()
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
	a1 := CheckTheTree()
	Test(true, a1, []error{
		errors.New("should be low level account in all account numbers {false false false  b  [] [] [[1 1] [1 2]] [] [] 0 0 0 false}"),
		errors.New("the account number [2 1 8] for {true false true Fifo f  [] [] [[4 1] [2 1 8]] [] [] 0 0 0 false} not conected to the tree"),
	})
}

func TestEditAccount(t *testing.T) {
	account, index, err := FindAccountFromName(PrefixRevenue + "book")
	fmt.Println(err)
	if err == nil {
		account.Number = [][]uint{{3, 2, 1, 2}}
		EditAccount(false, index, account)
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
}

func TestAddAutoCompletion(t *testing.T) {
	a1 := AddAutoCompletion(AutoCompletion{"book2", 5, 0, []Discount{}})
	DbClose()
	PrintSlice(AutoCompletionEntries)
	Test(true, a1, nil)
}
