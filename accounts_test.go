package main

import (
	"errors"
	"fmt"
	"testing"
)

func TestACCOUNT_STRUCT_FROM_BARCODE(t *testing.T) {
	INDEX_OF_ACCOUNT_NUMBER = 0
	account_struct, index, err := ACCOUNT_STRUCT_FROM_BARCODE("kaslajs")
	TEST(true, err, nil)
	TEST(true, index, 1)
	TEST(true, account_struct, ACCOUNT{false, false, false, "", "CURRENT_ASSETS", "", []string{}, []string{"sijadpodjpao", "kaslajs"}, [][]uint{{1, 1}, {}}, []uint{2, 0}, [][]string{{"ASSETS"}, {}}, 0, 0, 0, false})
}

func TestACCOUNT_STRUCT_FROM_NAME(t *testing.T) {
	INDEX_OF_ACCOUNT_NUMBER = 0
	account_struct, index, err := ACCOUNT_STRUCT_FROM_NAME("ASSETS")
	TEST(true, err, nil)
	TEST(true, index, 0)
	TEST(true, account_struct, ACCOUNT{false, false, false, "", "ASSETS", "", []string{}, []string{"nojdsjdpq"}, [][]uint{{1}, {}}, []uint{1, 0}, [][]string{{}, {}}, 0, 0, 0, false})
}

func TestADD_ACCOUNT(t *testing.T) {
	// DB_ACCOUNTS.DropAll()
	a := ADD_ACCOUNT(ACCOUNT{
		IS_LOW_LEVEL_ACCOUNT:             true,
		IS_CREDIT:                        true,
		IS_TEMPORARY:                     false,
		COST_FLOW_TYPE:                   WMA,
		ACCOUNT_NAME:                     "rent",
		NOTES:                            "",
		IMAGE:                            []string{},
		BARCODE:                          []string{},
		ACCOUNT_NUMBER:                   [][]uint{{4, 1}},
		ACCOUNT_LEVELS:                   []uint{},
		FATHER_AND_GRANDPA_ACCOUNTS_NAME: [][]string{},
		ALERT_FOR_MINIMUM_QUANTITY_BY_TURNOVER_IN_DAYS: 0,
		ALERT_FOR_MINIMUM_QUANTITY_BY_QUINTITY:         0,
		TARGET_BALANCE:                                 0,
		IF_THE_TARGET_BALANCE_IS_LESS_IS_GOOD:          false,
	})
	DB_CLOSE()
	PRINT_FORMATED_ACCOUNTS()
	TEST(true, a, nil)
}

func TestCHECK_IF_ACCOUNT_NUMBER_DUPLICATED(t *testing.T) {
	a := CHECK_IF_ACCOUNT_NUMBER_DUPLICATED()
	fmt.Println(a)
	TEST(true, a, []error{errors.New("the account number [2] for {false false true fifo e  [] [] [[4] [2]] [] [] 0 0 0 false} duplicated")})
}

func TestCHECK_IF_LOW_LEVEL_ACCOUNT_FOR_ALL(t *testing.T) {
	a := CHECK_IF_LOW_LEVEL_ACCOUNT_FOR_ALL()
	TEST(true, a, []error{errors.New("should be low level account in all account numbers {false false false  b  [] [] [[1 1] [1 2]] [] [] 0 0 0 false}")})
}

func TestCHECK_IF_THE_TREE_CONNECTED(t *testing.T) {
	a := CHECK_IF_THE_TREE_CONNECTED()
	fmt.Println(a)
	TEST(true, a, []error{errors.New("the account number [2 1 8] for {true false true fifo f  [] [] [[4 1] [2 1 8]] [] [] 0 0 0 false} not conected to the tree")})
}

func TestCHECK_THE_TREE(t *testing.T) {
	CHECK_THE_TREE()
	TEST(true, ERRORS_MESSAGES, []error{
		errors.New("should be low level account in all account numbers {false false false  b  [] [] [[1 1] [1 2]] [] [] 0 0 0 false}"),
		errors.New("the account number [2 1 8] for {true false true fifo f  [] [] [[4 1] [2 1 8]] [] [] 0 0 0 false} not conected to the tree"),
	})
}

func TestEDIT_ACCOUNT(t *testing.T) {
	account_struct, index, err := ACCOUNT_STRUCT_FROM_NAME("retined_earnings")
	fmt.Println(err)
	if err == nil {
		account_struct.IS_CREDIT = true
		account_struct.IS_TEMPORARY = true
		account_struct.BARCODE = []string{}
		account_struct.ACCOUNT_NUMBER = [][]uint{{4, 1}}
		account_struct.COST_FLOW_TYPE = WMA
		EDIT_ACCOUNT(true, index, account_struct)
	}
	DB_CLOSE()
	// TEST(true,)
	PRINT_FORMATED_ACCOUNTS()
}

func TestIS_BARCODES_USED(t *testing.T) {
	a := IS_BARCODES_USED([]string{"a", "b"})
	TEST(true, a, true)
	a = IS_BARCODES_USED([]string{"c", "b"})
	TEST(true, a, false)
}

func TestIS_IT_HIGH_THAN_BY_ORDER(t *testing.T) {
	a := IS_IT_HIGH_THAN_BY_ORDER([]uint{1}, []uint{1, 2})
	TEST(true, a, true)
	a = IS_IT_HIGH_THAN_BY_ORDER([]uint{1}, []uint{1})
	TEST(true, a, false)
	a = IS_IT_HIGH_THAN_BY_ORDER([]uint{1, 2}, []uint{1})
	TEST(true, a, false)
	a = IS_IT_HIGH_THAN_BY_ORDER([]uint{3}, []uint{1, 1})
	TEST(true, a, false)
	a = IS_IT_HIGH_THAN_BY_ORDER([]uint{1, 5}, []uint{3})
	TEST(true, a, true)
}

func TestIS_IT_POSSIBLE_TO_BE_SUB_ACCOUNT(t *testing.T) {
	a := IS_IT_POSSIBLE_TO_BE_SUB_ACCOUNT([]uint{1}, []uint{1, 2})
	TEST(true, a, true)
	a = IS_IT_POSSIBLE_TO_BE_SUB_ACCOUNT([]uint{1}, []uint{2})
	TEST(true, a, false)
	a = IS_IT_POSSIBLE_TO_BE_SUB_ACCOUNT([]uint{}, []uint{2})
	TEST(true, a, false)
	a = IS_IT_POSSIBLE_TO_BE_SUB_ACCOUNT([]uint{1}, []uint{1, 1, 2})
	TEST(true, a, true)
}

func TestIS_IT_SUB_ACCOUNT_USING_NAME(t *testing.T) {
	INDEX_OF_ACCOUNT_NUMBER = 0
	a := IS_IT_SUB_ACCOUNT_USING_NAME("ASSETS", "CASH_AND_CASH_EQUIVALENTS")
	TEST(true, a, true)
	a = IS_IT_SUB_ACCOUNT_USING_NAME("CASH_AND_CASH_EQUIVALENTS", "ASSETS")
	TEST(true, a, false)
	INDEX_OF_ACCOUNT_NUMBER = 1
	a = IS_IT_SUB_ACCOUNT_USING_NAME("ASSETS", "CASH_AND_CASH_EQUIVALENTS")
	TEST(true, a, false)
}

func TestIS_IT_SUB_ACCOUNT_USING_NUMBER(t *testing.T) {
	a := IS_IT_SUB_ACCOUNT_USING_NUMBER([]uint{1}, []uint{1, 2})
	TEST(true, a, true)
	a = IS_IT_SUB_ACCOUNT_USING_NUMBER([]uint{1}, []uint{2})
	TEST(true, a, false)
	a = IS_IT_SUB_ACCOUNT_USING_NUMBER([]uint{}, []uint{2})
	TEST(true, a, false)
}

func TestIS_IT_THE_FATHER(t *testing.T) {
	a := IS_IT_THE_FATHER([]uint{1}, []uint{1, 2})
	TEST(true, a, true)
	a = IS_IT_THE_FATHER([]uint{1}, []uint{2})
	TEST(true, a, false)
	a = IS_IT_THE_FATHER([]uint{}, []uint{2})
	TEST(true, a, false)
	a = IS_IT_THE_FATHER([]uint{1}, []uint{1, 1, 2})
	TEST(true, a, false)
}

func TestIS_USED_IN_JOURNAL(t *testing.T) {
	a := IS_USED_IN_JOURNAL("book")
	TEST(true, a, false)
}

func TestMAX_LEN_FOR_ACCOUNT_NUMBER(t *testing.T) {
	a := MAX_LEN_FOR_ACCOUNT_NUMBER()
	TEST(true, a, 2)
}

func TestPRINT_FORMATED_ACCOUNTS(t *testing.T) {
	PRINT_FORMATED_ACCOUNTS()
	//e1:=
	//TEST(true,a1,e1)
}

func TestSET_FATHER_AND_GRANDPA_ACCOUNTS_NAME(t *testing.T) {
	// ACCOUNTS = []ACCOUNT{
	// 	{false, false, false, "", "a", "", []string{}, []string{}, [][]uint{{1}, {1}}, []uint{}, [][]string{}, 0, 0, 0, false},
	// 	{false, false, false, "", "b", "", []string{}, []string{}, [][]uint{{1, 1}, {1, 2}}, []uint{}, [][]string{}, 0, 0, 0, false},
	// 	{false, false, false, "", "c", "", []string{}, []string{}, [][]uint{{1, 2}, {1, 3}}, []uint{}, [][]string{}, 0, 0, 0, false},
	// 	{false, false, true, "fifo", "d", "", []string{}, []string{}, [][]uint{{2}, {2}}, []uint{}, [][]string{}, 0, 0, 0, false},
	// 	{false, false, true, "fifo", "e", "", []string{}, []string{}, [][]uint{{4}, {2}}, []uint{}, [][]string{}, 0, 0, 0, false},
	// 	{false, false, true, "fifo", "f", "", []string{}, []string{}, [][]uint{{4, 1}, {2, 1, 8}}, []uint{}, [][]string{}, 0, 0, 0, false},
	// }
	// SET_THE_ACCOUNTS()
	// SET_FATHER_AND_GRANDPA_ACCOUNTS_NAME()
	// e := []ACCOUNT{
	// 	{false, false, false, "", "a", "", []string{}, []string{}, [][]uint{{1}, {1}}, []uint{1, 1}, [][]string{{}, {}}, 0, 0, 0, false},
	// 	{false, false, false, "", "b", "", []string{}, []string{}, [][]uint{{1, 1}, {1, 2}}, []uint{2, 2}, [][]string{{"a", "a"}, {"a", "a"}}, 0, 0, 0, false},
	// 	{false, false, false, "", "c", "", []string{}, []string{}, [][]uint{{1, 2}, {1, 3}}, []uint{2, 2}, [][]string{{"a", "a"}, {"a", "a"}}, 0, 0, 0, false},
	// 	{false, false, false, "", "d", "", []string{}, []string{}, [][]uint{{2}, {2}}, []uint{1, 1}, [][]string{{}, {}}, 0, 0, 0, false},
	// 	{false, false, false, "", "e", "", []string{}, []string{}, [][]uint{{4}, {2}}, []uint{1, 1}, [][]string{{}, {}}, 0, 0, 0, false},
	// 	{false, false, false, "", "f", "", []string{}, []string{}, [][]uint{{4, 1}, {2, 1, 8}}, []uint{2, 3}, [][]string{{"e", "e"}, {"d", "e", "d", "e"}}, 0, 0, 0, false},
	// }
	// TEST(true, ACCOUNTS, e)
}

func TestSET_THE_ACCOUNTS(t *testing.T) {
	SET_THE_ACCOUNTS()
	PRINT_FORMATED_ACCOUNTS()
	TEST(true, ACCOUNTS, []ACCOUNT{
		// {false, false, false, "", "a", "", []string{}, []string{"a"}, [][]uint{{1}, {1}}, []uint{1, 1}, [][]string{{}, {}}, 0, 0, 0, false},
		// {false, false, false, "", "b", "", []string{}, []string{}, [][]uint{{1, 1}, {1, 2}}, []uint{2, 2}, [][]string{{"a"}, {"a"}}, 0, 0, 0, false},
		// {true, false, false, "fifo", "c", "", []string{}, []string{}, [][]uint{{1, 2}, {1, 3}}, []uint{2, 2}, [][]string{{"a"}, {"a"}}, 0, 0, 0, false},
		// {false, false, false, "", "d", "", []string{}, []string{}, [][]uint{{2}, {2}}, []uint{1, 1}, [][]string{{}, {}}, 0, 0, 0, false},
		// {false, false, false, "", "e", "", []string{}, []string{}, [][]uint{{4}, {2}}, []uint{1, 1}, [][]string{{}, {}}, 0, 0, 0, false},
		// {true, false, true, "", "f", "", []string{}, []string{}, [][]uint{{4, 1}, {2, 1, 8}}, []uint{2, 3}, [][]string{{"e"}, {"d", "e"}}, 0, 0, 0, false},
	})
}
