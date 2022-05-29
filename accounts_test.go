package main

import (
	"errors"
	"fmt"
	"testing"
)

func TestAccountStructFromBarcode(t *testing.T) {
	VIndexOfAccountNumber = 0
	account, index, err := FFindAccountFromNameOrBarcode("kaslajs")
	FTest(true, err, nil)
	FTest(true, index, 1)
	FTest(true, account, SAccount{false, "", "CURRENT_ASSETS", "", []string{}, []string{"sijadpodjpao", "kaslajs"}, [][]uint{{1, 1}, {}}, []uint{2, 0}, [][]string{{"ASSETS"}, {}}})
}

func TestAccountStructFromName(t *testing.T) {
	VIndexOfAccountNumber = 0
	account, index, err := FFindAccountFromName("ASSETS")
	FTest(true, err, nil)
	FTest(true, index, 0)
	FTest(true, account, SAccount{false, "", "ASSETS", "", []string{}, []string{"nojdsjdpq"}, [][]uint{{1}, {}}, []uint{1, 0}, [][]string{{}, {}}})
}

func TestAddAccount(t *testing.T) {
	FDbOpenAll()
	var a1 error

	a1 = FAddAccount(SAccount{
		TIsCredit:      false,
		TCostFlowType:  CHighLevelAccount,
		TAccountName:   "assets",
		TAccountNumber: [][]uint{{1}},
	})
	FTest(true, a1, nil)

	a1 = FAddAccount(SAccount{
		TIsCredit:      false,
		TCostFlowType:  CHighLevelAccount,
		TAccountName:   "current assets",
		TAccountNumber: [][]uint{{1, 1}},
	})
	FTest(true, a1, nil)

	a1 = FAddAccount(SAccount{
		TIsCredit:      false,
		TCostFlowType:  CWma,
		TAccountName:   "cash",
		TAccountNumber: [][]uint{{1, 1, 1}},
	})
	FTest(true, a1, nil)

	a1 = FAddAccount(SAccount{
		TIsCredit:      false,
		TCostFlowType:  CHighLevelAccount,
		TAccountName:   "invnetory",
		TAccountNumber: [][]uint{{1, 1, 3}},
	})
	FTest(true, a1, nil)

	a1 = FAddAccount(SAccount{
		TIsCredit:      false,
		TCostFlowType:  CFifo,
		TAccountName:   "book",
		TAccountNumber: [][]uint{{1, 1, 3, 1}},
	})
	FTest(true, a1, nil)

	a1 = FAddAccount(SAccount{
		TIsCredit:      true,
		TCostFlowType:  CHighLevelAccount,
		TAccountName:   "liabilities",
		TAccountNumber: [][]uint{{2}},
	})
	FTest(true, a1, nil)

	a1 = FAddAccount(SAccount{
		TIsCredit:      true,
		TCostFlowType:  CHighLevelAccount,
		TAccountName:   "owner's equity",
		TAccountNumber: [][]uint{{3}},
	})
	FTest(true, a1, nil)

	a1 = FAddAccount(SAccount{
		TIsCredit:      true,
		TCostFlowType:  CHighLevelAccount,
		TAccountName:   "retained earnings",
		TAccountNumber: [][]uint{{3, 1}},
	})
	FTest(true, a1, nil)

	a1 = FAddAccount(SAccount{
		TIsCredit:      true,
		TCostFlowType:  CHighLevelAccount,
		TAccountName:   "income",
		TAccountNumber: [][]uint{{3, 2}},
	})
	FTest(true, a1, nil)

	a1 = FAddAccount(SAccount{
		TIsCredit:      true,
		TCostFlowType:  CHighLevelAccount,
		TAccountName:   "revenue",
		TAccountNumber: [][]uint{{3, 2, 1}},
	})
	FTest(true, a1, nil)

	a1 = FAddAccount(SAccount{
		TIsCredit:      true,
		TCostFlowType:  CWma,
		TAccountName:   "rent",
		TAccountNumber: [][]uint{{3, 2, 1, 1}},
	})
	FTest(true, a1, nil)

	a1 = FAddAccount(SAccount{
		TIsCredit:      false,
		TCostFlowType:  CHighLevelAccount,
		TAccountName:   "expense",
		TAccountNumber: [][]uint{{3, 2, 2}},
	})
	FTest(true, a1, nil)

	a1 = FAddAccount(SAccount{
		TIsCredit:      false,
		TCostFlowType:  CHighLevelAccount,
		TAccountName:   "discounts",
		TAccountNumber: [][]uint{{3, 2, 2, 1}},
	})
	FTest(true, a1, nil)

	a1 = FAddAccount(SAccount{
		TIsCredit:      false,
		TCostFlowType:  CWma,
		TAccountName:   VInvoiceDiscount,
		TAccountNumber: [][]uint{{3, 2, 2, 1, 1}},
	})
	FTest(true, a1, nil)

	a1 = FAddAccount(SAccount{
		TIsCredit:      false,
		TCostFlowType:  CWma,
		TAccountName:   CPrefixDiscount + "book",
		TAccountNumber: [][]uint{{3, 2, 2, 1, 2}},
	})
	FTest(true, a1, nil)

	a1 = FAddAccount(SAccount{
		TIsCredit:      false,
		TCostFlowType:  CHighLevelAccount,
		TAccountName:   "cost of goods sold",
		TAccountNumber: [][]uint{{3, 2, 2, 2}},
	})
	FTest(true, a1, nil)

	a1 = FAddAccount(SAccount{
		TIsCredit:      false,
		TCostFlowType:  CWma,
		TAccountName:   CPrefixCost + "book",
		TAccountNumber: [][]uint{{3, 2, 2, 2, 1}},
	})
	FTest(true, a1, nil)

	FDbCloseAll()
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
		account.TAccountNumber = [][]uint{{3, 2, 1, 2}}
		FEditAccount(false, index, account)
	}
	FDbCloseAll()
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
	FDbCloseAll()
	FPrintSlice(VAutoCompletionEntries)
	FTest(true, a1, nil)
}
