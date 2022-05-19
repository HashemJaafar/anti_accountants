package main

import (
	"errors"
	"fmt"
	"testing"
)

func TestAccountStructFromBarcode(t *testing.T) {
	VIndexOfAccountNumber = 0
	account, index, err := FFindAccountFromBarcode("kaslajs")
	FTest(true, err, nil)
	FTest(true, index, 1)
	FTest(true, account, SAccount{false, false, "", "CURRENT_ASSETS", "", []string{}, []string{"sijadpodjpao", "kaslajs"}, [][]uint{{1, 1}, {}}, []uint{2, 0}, [][]string{{"ASSETS"}, {}}})
}

func TestAccountStructFromName(t *testing.T) {
	VIndexOfAccountNumber = 0
	account, index, err := FFindAccountFromName("ASSETS")
	FTest(true, err, nil)
	FTest(true, index, 0)
	FTest(true, account, SAccount{false, false, "", "ASSETS", "", []string{}, []string{"nojdsjdpq"}, [][]uint{{1}, {}}, []uint{1, 0}, [][]string{{}, {}}})
}

func TestAddAccount(t *testing.T) {
	var a1 error

	a1 = FAddAccount(SAccount{
		IsLowLevel:   false,
		IsCredit:     false,
		CostFlowType: CWma,
		Name:         "assets",
		Number:       [][]uint{{1}},
	})
	FTest(true, a1, nil)

	a1 = FAddAccount(SAccount{
		IsLowLevel:   false,
		IsCredit:     false,
		CostFlowType: CWma,
		Name:         "current assets",
		Number:       [][]uint{{1, 1}},
	})
	FTest(true, a1, nil)

	a1 = FAddAccount(SAccount{
		IsLowLevel:   true,
		IsCredit:     false,
		CostFlowType: CWma,
		Name:         "cash",
		Number:       [][]uint{{1, 1, 1}},
	})
	FTest(true, a1, nil)

	a1 = FAddAccount(SAccount{
		IsLowLevel:   false,
		IsCredit:     false,
		CostFlowType: CWma,
		Name:         "invnetory",
		Number:       [][]uint{{1, 1, 3}},
	})
	FTest(true, a1, nil)

	a1 = FAddAccount(SAccount{
		IsLowLevel:   true,
		IsCredit:     false,
		CostFlowType: CFifo,
		Name:         "book",
		Number:       [][]uint{{1, 1, 3, 1}},
	})
	FTest(true, a1, nil)

	a1 = FAddAccount(SAccount{
		IsLowLevel:   false,
		IsCredit:     true,
		CostFlowType: CWma,
		Name:         "liabilities",
		Number:       [][]uint{{2}},
	})
	FTest(true, a1, nil)

	a1 = FAddAccount(SAccount{
		IsLowLevel:   false,
		IsCredit:     true,
		CostFlowType: CWma,
		Name:         "owner's equity",
		Number:       [][]uint{{3}},
	})
	FTest(true, a1, nil)

	a1 = FAddAccount(SAccount{
		IsLowLevel:   true,
		IsCredit:     true,
		CostFlowType: CWma,
		Name:         "retained earnings",
		Number:       [][]uint{{3, 1}},
	})
	FTest(true, a1, nil)

	a1 = FAddAccount(SAccount{
		IsLowLevel:   false,
		IsCredit:     true,
		CostFlowType: CWma,
		Name:         "income",
		Number:       [][]uint{{3, 2}},
	})
	FTest(true, a1, nil)

	a1 = FAddAccount(SAccount{
		IsLowLevel:   false,
		IsCredit:     true,
		CostFlowType: CWma,
		Name:         "revenue",
		Number:       [][]uint{{3, 2, 1}},
	})
	FTest(true, a1, nil)

	a1 = FAddAccount(SAccount{
		IsLowLevel:   true,
		IsCredit:     true,
		CostFlowType: CWma,
		Name:         "rent",
		Number:       [][]uint{{3, 2, 1, 1}},
	})
	FTest(true, a1, nil)

	a1 = FAddAccount(SAccount{
		IsLowLevel:   false,
		IsCredit:     false,
		CostFlowType: CWma,
		Name:         "expense",
		Number:       [][]uint{{3, 2, 2}},
	})
	FTest(true, a1, nil)

	a1 = FAddAccount(SAccount{
		IsLowLevel:   false,
		IsCredit:     false,
		CostFlowType: CWma,
		Name:         "discounts",
		Number:       [][]uint{{3, 2, 2, 1}},
	})
	FTest(true, a1, nil)

	a1 = FAddAccount(SAccount{
		IsLowLevel:   true,
		IsCredit:     false,
		CostFlowType: CWma,
		Name:         VInvoiceDiscount,
		Number:       [][]uint{{3, 2, 2, 1, 1}},
	})
	FTest(true, a1, nil)

	a1 = FAddAccount(SAccount{
		IsLowLevel:   true,
		IsCredit:     false,
		CostFlowType: CWma,
		Name:         CPrefixDiscount + "book",
		Number:       [][]uint{{3, 2, 2, 1, 2}},
	})
	FTest(true, a1, nil)

	a1 = FAddAccount(SAccount{
		IsLowLevel:   false,
		IsCredit:     false,
		CostFlowType: CWma,
		Name:         "cost of goods sold",
		Number:       [][]uint{{3, 2, 2, 2}},
	})
	FTest(true, a1, nil)

	a1 = FAddAccount(SAccount{
		IsLowLevel:   true,
		IsCredit:     false,
		CostFlowType: CWma,
		Name:         CPrefixCost + "book",
		Number:       [][]uint{{3, 2, 2, 2, 1}},
	})
	FTest(true, a1, nil)

	FDbClose()
	FPrintFormatedAccounts()
}

func TestCheckIfAccountNumberDuplicated(t *testing.T) {
	a := FCheckIfAccountNumberDuplicated()
	fmt.Println(a)
	FTest(true, a, []error{errors.New("the account number [2] for {false false true Fifo e  [] [] [[4] [2]] [] [] 0 0 0 false} duplicated")})
}

func TestCheckIfLowLevelAccountForAll(t *testing.T) {
	a := FCheckIfLowLevelAccountForAll()
	FTest(true, a, []error{errors.New("should be low level account in all account numbers {false false false  b  [] [] [[1 1] [1 2]] [] [] 0 0 0 false}")})
}

func TestCheckIfTheTreeConnected(t *testing.T) {
	a := FCheckIfTheTreeConnected()
	fmt.Println(a)
	FTest(true, a, []error{errors.New("the account number [2 1 8] for {true false true Fifo f  [] [] [[4 1] [2 1 8]] [] [] 0 0 0 false} not conected to the tree")})
}

func TestCheckTheTree(t *testing.T) {
	a1 := FCheckTheTree()
	FTest(true, a1, []error{
		errors.New("should be low level account in all account numbers {false false false  b  [] [] [[1 1] [1 2]] [] [] 0 0 0 false}"),
		errors.New("the account number [2 1 8] for {true false true Fifo f  [] [] [[4 1] [2 1 8]] [] [] 0 0 0 false} not conected to the tree"),
	})
}

func TestEditAccount(t *testing.T) {
	account, index, err := FFindAccountFromName(CPrefixRevenue + "book")
	fmt.Println(err)
	if err == nil {
		account.Number = [][]uint{{3, 2, 1, 2}}
		FEditAccount(false, index, account)
	}
	FDbClose()
	FPrintFormatedAccounts()
}

func TestIsBarcodesUsed(t *testing.T) {
	a := FIsBarcodesUsed([]string{"a", "b"})
	FTest(true, a, true)
	a = FIsBarcodesUsed([]string{"c", "b"})
	FTest(true, a, false)
}

func TestIsItHighThanByOrder(t *testing.T) {
	a := FIsItHighThanByOrder([]uint{1}, []uint{1, 2})
	FTest(true, a, true)
	a = FIsItHighThanByOrder([]uint{1}, []uint{1})
	FTest(true, a, false)
	a = FIsItHighThanByOrder([]uint{1, 2}, []uint{1})
	FTest(true, a, false)
	a = FIsItHighThanByOrder([]uint{3}, []uint{1, 1})
	FTest(true, a, false)
	a = FIsItHighThanByOrder([]uint{1, 5}, []uint{3})
	FTest(true, a, true)
	a = FIsItHighThanByOrder([]uint{1, 1}, []uint{1, 2})
	FTest(true, a, true)
	a = FIsItHighThanByOrder([]uint{4}, []uint{1, 2})
	FTest(true, a, false)
}

func TestIsItPossibleToBeSubAccount(t *testing.T) {
	a := FIsItPossibleToBeSubAccount([]uint{1}, []uint{1, 2})
	FTest(true, a, true)
	a = FIsItPossibleToBeSubAccount([]uint{1}, []uint{2})
	FTest(true, a, false)
	a = FIsItPossibleToBeSubAccount([]uint{}, []uint{2})
	FTest(true, a, false)
	a = FIsItPossibleToBeSubAccount([]uint{1}, []uint{1, 1, 2})
	FTest(true, a, true)
}

func TestIsItSubAccountUsingName(t *testing.T) {
	VIndexOfAccountNumber = 0
	a := FIsItSubAccountUsingName("ASSETS", "CASH_AND_CASH_EQUIVALENTS")
	FTest(true, a, true)
	a = FIsItSubAccountUsingName("CASH_AND_CASH_EQUIVALENTS", "ASSETS")
	FTest(true, a, false)
	VIndexOfAccountNumber = 1
	a = FIsItSubAccountUsingName("ASSETS", "CASH_AND_CASH_EQUIVALENTS")
	FTest(true, a, false)
}

func TestIsItSubAccountUsingNumber(t *testing.T) {
	a := FIsItSubAccountUsingNumber([]uint{1}, []uint{1, 2})
	FTest(true, a, true)
	a = FIsItSubAccountUsingNumber([]uint{1}, []uint{2})
	FTest(true, a, false)
	a = FIsItSubAccountUsingNumber([]uint{}, []uint{2})
	FTest(true, a, false)
}

func TestIsItTheFather(t *testing.T) {
	a := FIsItTheFather([]uint{1}, []uint{1, 2})
	FTest(true, a, true)
	a = FIsItTheFather([]uint{1}, []uint{2})
	FTest(true, a, false)
	a = FIsItTheFather([]uint{}, []uint{2})
	FTest(true, a, false)
	a = FIsItTheFather([]uint{1}, []uint{1, 1, 2})
	FTest(true, a, false)
}

func TestIsUsedInJournal(t *testing.T) {
	a := FIsUsedInJournal("book")
	FTest(true, a, false)
}

func TestMaxLenForAccountNumber(t *testing.T) {
	a := FMaxLenForAccountNumber()
	FTest(true, a, 2)
}

func TestPrintFormatedAccounts(t *testing.T) {
	FPrintFormatedAccounts()
}

func TestSetTheAccounts(t *testing.T) {
	FSetTheAccounts()
	FPrintFormatedAccounts()
}

func TestAddAutoCompletion(t *testing.T) {
	a1 := FAddAutoCompletion(SAutoCompletion{"book2", 5, 0, []SPQ{}})
	FDbClose()
	FPrintSlice(VAutoCompletionEntries)
	FTest(true, a1, nil)
}
