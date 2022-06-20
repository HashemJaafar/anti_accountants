package main

import (
	"errors"
	"fmt"
	"log"
	"testing"
)

func TestAddAccount(t *testing.T) {
	VCompanyName = "anti_accountants"
	FDbOpenAll()
	var a1 SAccount3

	a1 = FAddAccount(true, SAccount1{
		IsCredit:     false,
		CostFlowType: CHighLevelAccount,
		Name:         "assets",
		Number:       [][]uint{{1}},
	})
	FTest(true, a1, SAccount3{})

	a1 = FAddAccount(true, SAccount1{
		IsCredit:     false,
		CostFlowType: CHighLevelAccount,
		Name:         "current assets",
		Number:       [][]uint{{1, 1}},
	})
	FTest(true, a1, SAccount3{})

	a1 = FAddAccount(true, SAccount1{
		IsCredit:     false,
		CostFlowType: CWma,
		Name:         "cash",
		Number:       [][]uint{{1, 1, 1}},
	})
	FTest(true, a1, SAccount3{})

	a1 = FAddAccount(true, SAccount1{
		IsCredit:     false,
		CostFlowType: CHighLevelAccount,
		Name:         "invnetory",
		Number:       [][]uint{{1, 1, 3}},
	})
	FTest(true, a1, SAccount3{})

	a1 = FAddAccount(true, SAccount1{
		IsCredit:     false,
		CostFlowType: CFifo,
		Name:         "book",
		Number:       [][]uint{{1, 1, 3, 1}},
	})
	FTest(true, a1, SAccount3{})

	a1 = FAddAccount(true, SAccount1{
		IsCredit:     true,
		CostFlowType: CHighLevelAccount,
		Name:         "liabilities",
		Number:       [][]uint{{2}},
	})
	FTest(true, a1, SAccount3{})

	a1 = FAddAccount(true, SAccount1{
		IsCredit:     true,
		CostFlowType: CHighLevelAccount,
		Name:         "owner's equity",
		Number:       [][]uint{{3}},
	})
	FTest(true, a1, SAccount3{})

	a1 = FAddAccount(true, SAccount1{
		IsCredit:     true,
		CostFlowType: CHighLevelAccount,
		Name:         "retained earnings",
		Number:       [][]uint{{3, 1}},
	})
	FTest(true, a1, SAccount3{})

	a1 = FAddAccount(true, SAccount1{
		IsCredit:     true,
		CostFlowType: CHighLevelAccount,
		Name:         "income",
		Number:       [][]uint{{3, 2}},
	})
	FTest(true, a1, SAccount3{})

	a1 = FAddAccount(true, SAccount1{
		IsCredit:     true,
		CostFlowType: CHighLevelAccount,
		Name:         "revenue",
		Number:       [][]uint{{3, 2, 1}},
	})
	FTest(true, a1, SAccount3{})

	a1 = FAddAccount(true, SAccount1{
		IsCredit:     true,
		CostFlowType: CWma,
		Name:         "rent",
		Number:       [][]uint{{3, 2, 1, 1}},
	})
	FTest(true, a1, SAccount3{})

	a1 = FAddAccount(true, SAccount1{
		IsCredit:     false,
		CostFlowType: CHighLevelAccount,
		Name:         "expense",
		Number:       [][]uint{{3, 2, 2}},
	})
	FTest(true, a1, SAccount3{})

	a1 = FAddAccount(true, SAccount1{
		IsCredit:     false,
		CostFlowType: CHighLevelAccount,
		Name:         "discounts",
		Number:       [][]uint{{3, 2, 2, 1}},
	})
	FTest(true, a1, SAccount3{})

	a1 = FAddAccount(true, SAccount1{
		IsCredit:     false,
		CostFlowType: CWma,
		Name:         "Invoice PQ",
		Number:       [][]uint{{3, 2, 2, 1, 1}},
	})
	FTest(true, a1, SAccount3{})

	a1 = FAddAccount(true, SAccount1{
		IsCredit:     false,
		CostFlowType: CWma,
		Name:         CPrefixDiscount + "book",
		Number:       [][]uint{{3, 2, 2, 1, 2}},
	})
	FTest(true, a1, SAccount3{})

	a1 = FAddAccount(true, SAccount1{
		IsCredit:     false,
		CostFlowType: CHighLevelAccount,
		Name:         "cost of goods sold",
		Number:       [][]uint{{3, 2, 2, 2}},
	})
	FTest(true, a1, SAccount3{})

	a1 = FAddAccount(true, SAccount1{
		IsCredit:     false,
		CostFlowType: CWma,
		Name:         CPrefixCost + "book",
		Number:       [][]uint{{3, 2, 2, 2, 1}},
	})
	FTest(true, a1, SAccount3{})

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
	VCompanyName = "anti_accountants"
	FDbOpenAll()
	account, index, err := FFindAccountFromName("cash")
	fmt.Println(err)
	if err == nil {
		account.IsCredit = false
		FEditAccount(false, index, account)
	}
	FDbCloseAll()
	FPrintFormatedAccounts()
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
	VCompanyName = "anti_accountants"
	FDbOpenAll()
	FDbCloseAll()
	FPrintFormatedAccounts()
}

func TestSetTheAccounts(t *testing.T) {
	VCompanyName = "anti_accountants"
	FDbOpenAll()
	FSetTheAccounts()
	FDbCloseAll()
	FPrintFormatedAccounts()
}

func TestAddAutoCompletion(t *testing.T) {
	VCompanyName = "anti_accountants"
	FDbOpenAll()
	a1 := FAddAutoCompletion(SAutoCompletion{
		AccountName:          "book",
		PriceRevenue:         1250,
		PriceTax:             250,
		DiscountWay:          CDiscountPerOne,
		DiscountPerOne:       250,
		DiscountTotal:        0,
		DiscountPerQuantity:  SPQ{},
		DiscountDecisionTree: []SPQ{},
	})
	FDbCloseAll()
	fmt.Println(a1)
	FPrintStructSlice(false, VAutoCompletionEntries)
}

func TestFAddAccount(t *testing.T) {
	VCompanyName = "a"
	FDbOpenAll()
	a1 := FAddAccount(true, SAccount1{
		IsCredit:     false,
		CostFlowType: "",
		Inventory:    "home",
		Name:         "book",
		Notes:        "",
		Image:        []string{},
		Barcode:      []string{"1", "2"},
		Number:       [][]uint{{1}, {1}},
		Levels:       []uint{},
		FathersName:  [][]string{},
	})
	FPrintFormatedAccounts()
	FDbCloseAll()
	log.Println(a1)
}
