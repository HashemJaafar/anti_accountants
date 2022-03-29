package anti_accountants

import (
	"errors"
	"fmt"
	"testing"
)

func Test_ACCOUNT_STRUCT_FROM_NAME(t *testing.T) {
	ACCOUNTS = []ACCOUNT{
		{false, false, false, "", "ASSETS", "", []string{}, []string{"nojdsjdpq"}, [][]uint{{1}, {}}, []uint{1, 0}, [][]string{{}, {}}, 0, 0, 0, false},
		{false, false, false, "", "CURRENT_ASSETS", "", []string{}, []string{"sijadpodjpao", "kaslajs"}, [][]uint{{1, 1}, {}}, []uint{2, 0}, [][]string{{"ASSETS"}, {}}, 0, 0, 0, false},
		{true, false, false, "fifo", "CASH_AND_CASH_EQUIVALENTS", "", []string{}, []string{"888"}, [][]uint{{1, 1, 1}, {2}}, []uint{3, 1}, [][]string{{"ASSETS", "CURRENT_ASSETS"}, {}}, 0, 0, 0, false},
		{true, false, false, "fifo", "SHORT_TERM_INVESTMENTS", "", []string{}, []string{"SHORT_TERM_INVESTMENTS"}, [][]uint{{1, 2}, {5}}, []uint{2, 1}, [][]string{{"ASSETS"}, {}}, 0, 0, 0, false},
		{true, false, false, "", "RECEIVABLES", "", []string{}, []string{"RECEIVABLES"}, [][]uint{{1, 3}, {}}, []uint{2, 0}, [][]string{{"ASSETS"}, {}}, 0, 0, 0, false},
	}
	INDEX_OF_ACCOUNT_NUMBER = 0
	account_struct, index, err := ACCOUNT_STRUCT_FROM_NAME("ASSETS")
	TEST(true, err, nil)
	TEST(true, index, 0)
	TEST(true, account_struct, ACCOUNT{false, false, false, "", "ASSETS", "", []string{}, []string{"nojdsjdpq"}, [][]uint{{1}, {}}, []uint{1, 0}, [][]string{{}, {}}, 0, 0, 0, false})
}

func Test_ACCOUNT_STRUCT_FROM_BARCODE(t *testing.T) {
	ACCOUNTS = []ACCOUNT{
		{false, false, false, "", "ASSETS", "", []string{}, []string{"nojdsjdpq"}, [][]uint{{1}, {}}, []uint{1, 0}, [][]string{{}, {}}, 0, 0, 0, false},
		{false, false, false, "", "CURRENT_ASSETS", "", []string{}, []string{"sijadpodjpao", "kaslajs"}, [][]uint{{1, 1}, {}}, []uint{2, 0}, [][]string{{"ASSETS"}, {}}, 0, 0, 0, false},
		{true, false, false, "fifo", "CASH_AND_CASH_EQUIVALENTS", "", []string{}, []string{"888"}, [][]uint{{1, 1, 1}, {2}}, []uint{3, 1}, [][]string{{"ASSETS", "CURRENT_ASSETS"}, {}}, 0, 0, 0, false},
		{true, false, false, "fifo", "SHORT_TERM_INVESTMENTS", "", []string{}, []string{"SHORT_TERM_INVESTMENTS"}, [][]uint{{1, 2}, {5}}, []uint{2, 1}, [][]string{{"ASSETS"}, {}}, 0, 0, 0, false},
		{true, false, false, "", "RECEIVABLES", "", []string{}, []string{"RECEIVABLES"}, [][]uint{{1, 3}, {}}, []uint{2, 0}, [][]string{{"ASSETS"}, {}}, 0, 0, 0, false},
	}
	INDEX_OF_ACCOUNT_NUMBER = 0
	account_struct, index, err := ACCOUNT_STRUCT_FROM_BARCODE("kaslajs")
	TEST(true, err, nil)
	TEST(true, index, 1)
	TEST(true, account_struct, ACCOUNT{false, false, false, "", "CURRENT_ASSETS", "", []string{}, []string{"sijadpodjpao", "kaslajs"}, [][]uint{{1, 1}, {}}, []uint{2, 0}, [][]string{{"ASSETS"}, {}}, 0, 0, 0, false})
}

func Test_IS_IT_SUB_ACCOUNT_USING_NAME(t *testing.T) {
	ACCOUNTS = []ACCOUNT{
		{false, false, false, "", "ASSETS", "", []string{}, []string{"nojdsjdpq"}, [][]uint{{1}, {1}}, []uint{1, 0}, [][]string{{}, {}}, 0, 0, 0, false},
		{false, false, false, "", "CURRENT_ASSETS", "", []string{}, []string{"sijadpodjpao", "kaslajs"}, [][]uint{{1, 1}, {}}, []uint{2, 0}, [][]string{{"ASSETS"}, {}}, 0, 0, 0, false},
		{true, false, false, "fifo", "CASH_AND_CASH_EQUIVALENTS", "", []string{}, []string{"888"}, [][]uint{{1, 1, 1}, {2}}, []uint{3, 1}, [][]string{{"ASSETS", "CURRENT_ASSETS"}, {}}, 0, 0, 0, false},
		{true, false, false, "fifo", "SHORT_TERM_INVESTMENTS", "", []string{}, []string{"SHORT_TERM_INVESTMENTS"}, [][]uint{{1, 2}, {5}}, []uint{2, 1}, [][]string{{"ASSETS"}, {}}, 0, 0, 0, false},
		{true, false, false, "", "RECEIVABLES", "", []string{}, []string{"RECEIVABLES"}, [][]uint{{1, 3}, {}}, []uint{2, 0}, [][]string{{"ASSETS"}, {}}, 0, 0, 0, false},
	}
	INDEX_OF_ACCOUNT_NUMBER = 0
	a := IS_IT_SUB_ACCOUNT_USING_NAME("ASSETS", "CASH_AND_CASH_EQUIVALENTS")
	TEST(true, a, true)
	a = IS_IT_SUB_ACCOUNT_USING_NAME("CASH_AND_CASH_EQUIVALENTS", "ASSETS")
	TEST(true, a, false)
	INDEX_OF_ACCOUNT_NUMBER = 1
	a = IS_IT_SUB_ACCOUNT_USING_NAME("ASSETS", "CASH_AND_CASH_EQUIVALENTS")
	TEST(true, a, false)
}

func Test_IS_IT_SUB_ACCOUNT_USING_NUMBER(t *testing.T) {
	a := IS_IT_SUB_ACCOUNT_USING_NUMBER([]uint{1}, []uint{1, 2})
	TEST(true, a, true)
	a = IS_IT_SUB_ACCOUNT_USING_NUMBER([]uint{1}, []uint{2})
	TEST(true, a, false)
	a = IS_IT_SUB_ACCOUNT_USING_NUMBER([]uint{}, []uint{2})
	TEST(true, a, false)
}

func Test_IS_IT_POSSIBLE_TO_BE_SUB_ACCOUNT(t *testing.T) {
	a := IS_IT_POSSIBLE_TO_BE_SUB_ACCOUNT([]uint{1}, []uint{1, 2})
	TEST(true, a, true)
	a = IS_IT_POSSIBLE_TO_BE_SUB_ACCOUNT([]uint{1}, []uint{2})
	TEST(true, a, false)
	a = IS_IT_POSSIBLE_TO_BE_SUB_ACCOUNT([]uint{}, []uint{2})
	TEST(true, a, false)
	a = IS_IT_POSSIBLE_TO_BE_SUB_ACCOUNT([]uint{1}, []uint{1, 1, 2})
	TEST(true, a, true)
}

func Test_IS_IT_THE_FATHER(t *testing.T) {
	a := IS_IT_THE_FATHER([]uint{1}, []uint{1, 2})
	TEST(true, a, true)
	a = IS_IT_THE_FATHER([]uint{1}, []uint{2})
	TEST(true, a, false)
	a = IS_IT_THE_FATHER([]uint{}, []uint{2})
	TEST(true, a, false)
	a = IS_IT_THE_FATHER([]uint{1}, []uint{1, 1, 2})
	TEST(true, a, false)
}

func Test_IS_IT_HIGH_THAN_BY_ORDER(t *testing.T) {
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

func Test_MAX_LEN_FOR_ACCOUNT_NUMBER(t *testing.T) {
	ACCOUNTS = []ACCOUNT{
		{false, false, false, "", "a", "", []string{}, []string{}, [][]uint{{1}, {1}}, []uint{}, [][]string{}, 0, 0, 0, false},
		{false, false, false, "", "b", "", []string{}, []string{}, [][]uint{{1, 1}, {1, 2}}, []uint{}, [][]string{}, 0, 0, 0, false},
		{false, false, false, "", "c", "", []string{}, []string{}, [][]uint{{1, 2}, {1, 3}}, []uint{}, [][]string{}, 0, 0, 0, false},
		{false, false, true, "fifo", "d", "", []string{}, []string{}, [][]uint{{2}, {2}}, []uint{}, [][]string{}, 0, 0, 0, false},
		{false, false, true, "fifo", "e", "", []string{}, []string{}, [][]uint{{4}, {2}}, []uint{}, [][]string{}, 0, 0, 0, false},
		{false, false, true, "fifo", "f", "", []string{}, []string{}, [][]uint{{4, 1}, {}, {2, 1, 8}, {}}, []uint{}, [][]string{}, 0, 0, 0, false},
	}
	a := MAX_LEN_FOR_ACCOUNT_NUMBER()
	TEST(true, a, 2)
}
func Test_SET_FATHER_AND_GRANDPA_ACCOUNTS_NAME(t *testing.T) {
	ACCOUNTS = []ACCOUNT{
		{false, false, false, "", "a", "", []string{}, []string{}, [][]uint{{1}, {1}}, []uint{}, [][]string{}, 0, 0, 0, false},
		{false, false, false, "", "b", "", []string{}, []string{}, [][]uint{{1, 1}, {1, 2}}, []uint{}, [][]string{}, 0, 0, 0, false},
		{false, false, false, "", "c", "", []string{}, []string{}, [][]uint{{1, 2}, {1, 3}}, []uint{}, [][]string{}, 0, 0, 0, false},
		{false, false, true, "fifo", "d", "", []string{}, []string{}, [][]uint{{2}, {2}}, []uint{}, [][]string{}, 0, 0, 0, false},
		{false, false, true, "fifo", "e", "", []string{}, []string{}, [][]uint{{4}, {2}}, []uint{}, [][]string{}, 0, 0, 0, false},
		{false, false, true, "fifo", "f", "", []string{}, []string{}, [][]uint{{4, 1}, {2, 1, 8}}, []uint{}, [][]string{}, 0, 0, 0, false},
	}
	SET_THE_ACCOUNTS()
	SET_FATHER_AND_GRANDPA_ACCOUNTS_NAME()
	TEST(true, ACCOUNTS, []ACCOUNT{
		{false, false, false, "", "a", "", []string{}, []string{}, [][]uint{{1}, {1}}, []uint{}, [][]string{{}, {}}, 0, 0, 0, false},
		{false, false, false, "", "b", "", []string{}, []string{}, [][]uint{{1, 1}, {1, 2}}, []uint{}, [][]string{{"a"}, {"a"}}, 0, 0, 0, false},
		{false, false, false, "", "c", "", []string{}, []string{}, [][]uint{{1, 2}, {1, 3}}, []uint{}, [][]string{{"a"}, {"a"}}, 0, 0, 0, false},
		{false, false, true, "fifo", "d", "", []string{}, []string{}, [][]uint{{2}, {2}}, []uint{}, [][]string{{}, {}}, 0, 0, 0, false},
		{false, false, true, "fifo", "e", "", []string{}, []string{}, [][]uint{{4}, {2}}, []uint{}, [][]string{{}, {}}, 0, 0, 0, false},
		{false, false, true, "fifo", "f", "", []string{}, []string{}, [][]uint{{4, 1}, {2, 1, 8}}, []uint{}, [][]string{{"e"}, {"d", "e"}}, 0, 0, 0, false},
	})
}
func Test_FIND_ALL_INVENTORY_FILES(t *testing.T) {
	a := FIND_ALL_INVENTORY_FILES()
	TEST(true, a, []string{})
}

func Test_IS_BARCODES_USED(t *testing.T) {
	ACCOUNTS = []ACCOUNT{
		{false, false, false, "", "a", "", []string{}, []string{"a"}, [][]uint{{1}, {1}}, []uint{}, [][]string{}, 0, 0, 0, false},
		{false, false, false, "", "b", "", []string{}, []string{}, [][]uint{{1, 1}, {1, 2}}, []uint{}, [][]string{}, 0, 0, 0, false},
		{false, false, false, "", "c", "", []string{}, []string{}, [][]uint{{1, 2}, {1, 3}}, []uint{}, [][]string{}, 0, 0, 0, false},
		{false, false, true, "fifo", "d", "", []string{}, []string{}, [][]uint{{2}, {2}}, []uint{}, [][]string{}, 0, 0, 0, false},
		{false, false, true, "fifo", "e", "", []string{}, []string{}, [][]uint{{4}, {2}}, []uint{}, [][]string{}, 0, 0, 0, false},
		{false, false, true, "fifo", "f", "", []string{}, []string{}, [][]uint{{4, 1}, {2, 1, 8}}, []uint{}, [][]string{}, 0, 0, 0, false},
	}
	a := IS_BARCODES_USED([]string{"a", "b"})
	TEST(true, a, true)
	a = IS_BARCODES_USED([]string{"c", "b"})
	TEST(true, a, false)
}
func Test_UPDATE_INVENTORY_FILE_NAME(t *testing.T) {
	UPDATE_INVENTORY_FILE_NAME("book", "book1")
}
func Test_SET_THE_ACCOUNTS(t *testing.T) {
	ACCOUNTS = []ACCOUNT{
		{false, false, false, "", "a", "", []string{}, []string{"a"}, [][]uint{{1}, {1}}, []uint{}, [][]string{}, 0, 0, 0, false},
		{false, false, false, "", "b", "", []string{}, []string{}, [][]uint{{1, 1}, {1, 2}}, []uint{}, [][]string{}, 0, 0, 0, false},
		{true, false, false, "", "c", "", []string{}, []string{}, [][]uint{{1, 2}, {1, 3}}, []uint{}, [][]string{}, 0, 0, 0, false},
		{false, false, true, "fifo", "d", "", []string{}, []string{}, [][]uint{{2}, {2}}, []uint{}, [][]string{}, 0, 0, 0, false},
		{false, false, true, "fifo", "e", "", []string{}, []string{}, [][]uint{{4}, {2}}, []uint{}, [][]string{}, 0, 0, 0, false},
		{true, false, true, "fifo", "f", "", []string{}, []string{}, [][]uint{{4, 1}, {2, 1, 8}}, []uint{}, [][]string{}, 0, 0, 0, false},
	}
	SET_THE_ACCOUNTS()
	PRINT_FORMATED_ACCOUNTS()
	TEST(true, ACCOUNTS, []ACCOUNT{
		{false, false, false, "", "a", "", []string{}, []string{"a"}, [][]uint{{1}, {1}}, []uint{1, 1}, [][]string{{}, {}}, 0, 0, 0, false},
		{false, false, false, "", "b", "", []string{}, []string{}, [][]uint{{1, 1}, {1, 2}}, []uint{2, 2}, [][]string{{"a"}, {"a"}}, 0, 0, 0, false},
		{true, false, false, "fifo", "c", "", []string{}, []string{}, [][]uint{{1, 2}, {1, 3}}, []uint{2, 2}, [][]string{{"a"}, {"a"}}, 0, 0, 0, false},
		{false, false, false, "", "d", "", []string{}, []string{}, [][]uint{{2}, {2}}, []uint{1, 1}, [][]string{{}, {}}, 0, 0, 0, false},
		{false, false, false, "", "e", "", []string{}, []string{}, [][]uint{{4}, {2}}, []uint{1, 1}, [][]string{{}, {}}, 0, 0, 0, false},
		{true, false, true, "", "f", "", []string{}, []string{}, [][]uint{{4, 1}, {2, 1, 8}}, []uint{2, 3}, [][]string{{"e"}, {"d", "e"}}, 0, 0, 0, false},
	})
}
func Test_CHECK_IF_LOW_LEVEL_ACCOUNT_FOR_ALL(t *testing.T) {
	ACCOUNTS = []ACCOUNT{
		{false, false, false, "", "a", "", []string{}, []string{"a"}, [][]uint{{1}, {1}}, []uint{}, [][]string{}, 0, 0, 0, false},
		{false, false, false, "", "b", "", []string{}, []string{}, [][]uint{{1, 1}, {1, 2}}, []uint{}, [][]string{}, 0, 0, 0, false},
		{true, false, false, "", "c", "", []string{}, []string{}, [][]uint{{1, 2}, {1, 3}}, []uint{}, [][]string{}, 0, 0, 0, false},
		{false, false, true, "fifo", "d", "", []string{}, []string{}, [][]uint{{2}, {2}}, []uint{}, [][]string{}, 0, 0, 0, false},
		{false, false, true, "fifo", "e", "", []string{}, []string{}, [][]uint{{4}, {2}}, []uint{}, [][]string{}, 0, 0, 0, false},
		{true, false, true, "fifo", "f", "", []string{}, []string{}, [][]uint{{4, 1}, {2, 1, 8}}, []uint{}, [][]string{}, 0, 0, 0, false},
	}
	a := CHECK_IF_LOW_LEVEL_ACCOUNT_FOR_ALL()
	TEST(true, a, []error{errors.New("should be low level account in all account numbers {false false false  b  [] [] [[1 1] [1 2]] [] [] 0 0 0 false}")})
}
func Test_CHECK_IF_THE_TREE_CONNECTED(t *testing.T) {
	ACCOUNTS = []ACCOUNT{
		{false, false, false, "", "a", "", []string{}, []string{"a"}, [][]uint{{1}, {1}}, []uint{}, [][]string{}, 0, 0, 0, false},
		{false, false, false, "", "b", "", []string{}, []string{}, [][]uint{{1, 1}, {1, 2}}, []uint{}, [][]string{}, 0, 0, 0, false},
		{true, false, false, "", "c", "", []string{}, []string{}, [][]uint{{1, 2}, {1, 3}}, []uint{}, [][]string{}, 0, 0, 0, false},
		{false, false, true, "fifo", "d", "", []string{}, []string{}, [][]uint{{2}, {2}}, []uint{}, [][]string{}, 0, 0, 0, false},
		{false, false, true, "fifo", "e", "", []string{}, []string{}, [][]uint{{4}, {2}}, []uint{}, [][]string{}, 0, 0, 0, false},
		{true, false, true, "fifo", "f", "", []string{}, []string{}, [][]uint{{4, 1}, {2, 1, 8}}, []uint{}, [][]string{}, 0, 0, 0, false},
	}
	a := CHECK_IF_THE_TREE_CONNECTED()
	fmt.Println(a)
	TEST(true, a, []error{errors.New("the account number [2 1 8] for {true false true fifo f  [] [] [[4 1] [2 1 8]] [] [] 0 0 0 false} not conected to the tree")})
}
func Test_CHECK_IF_ACCOUNT_NUMBER_DUPLICATED(t *testing.T) {
	ACCOUNTS = []ACCOUNT{
		{false, false, false, "", "a", "", []string{}, []string{"a"}, [][]uint{{1}, {1}}, []uint{}, [][]string{}, 0, 0, 0, false},
		{false, false, false, "", "b", "", []string{}, []string{}, [][]uint{{1, 1}, {1, 2}}, []uint{}, [][]string{}, 0, 0, 0, false},
		{true, false, false, "", "c", "", []string{}, []string{}, [][]uint{{1, 2}, {1, 3}}, []uint{}, [][]string{}, 0, 0, 0, false},
		{false, false, true, "fifo", "d", "", []string{}, []string{}, [][]uint{{2}, {2}}, []uint{}, [][]string{}, 0, 0, 0, false},
		{false, false, true, "fifo", "e", "", []string{}, []string{}, [][]uint{{4}, {2}}, []uint{}, [][]string{}, 0, 0, 0, false},
		{true, false, true, "fifo", "f", "", []string{}, []string{}, [][]uint{{4, 1}, {2, 1, 8}}, []uint{}, [][]string{}, 0, 0, 0, false},
	}
	a := CHECK_IF_ACCOUNT_NUMBER_DUPLICATED()
	fmt.Println(a)
	TEST(true, a, []error{errors.New("the account number [2] for {false false true fifo e  [] [] [[4] [2]] [] [] 0 0 0 false} duplicated")})
}
func Test_CHECK_THE_TREE(t *testing.T) {
	ACCOUNTS = []ACCOUNT{
		{false, false, false, "", "a", "", []string{}, []string{"a"}, [][]uint{{1}, {1}}, []uint{}, [][]string{}, 0, 0, 0, false},
		{false, false, false, "", "b", "", []string{}, []string{}, [][]uint{{1, 1}, {1, 2}}, []uint{}, [][]string{}, 0, 0, 0, false},
		{true, false, false, "", "c", "", []string{}, []string{}, [][]uint{{1, 2}, {1, 3}}, []uint{}, [][]string{}, 0, 0, 0, false},
		{false, false, true, "fifo", "d", "", []string{}, []string{}, [][]uint{{2}, {2}}, []uint{}, [][]string{}, 0, 0, 0, false},
		{false, false, true, "fifo", "e", "", []string{}, []string{}, [][]uint{{4}, {2}}, []uint{}, [][]string{}, 0, 0, 0, false},
		{true, false, true, "fifo", "f", "", []string{}, []string{}, [][]uint{{4, 1}, {2, 1, 8}}, []uint{}, [][]string{}, 0, 0, 0, false},
	}
	CHECK_THE_TREE()
	TEST(true, ERRORS_MESSAGES, []error{
		errors.New("should be low level account in all account numbers {false false false  b  [] [] [[1 1] [1 2]] [] [] 0 0 0 false}"),
		errors.New("the account number [2 1 8] for {true false true fifo f  [] [] [[4 1] [2 1 8]] [] [] 0 0 0 false} not conected to the tree"),
	})
}
func Test_IS_USED_IN_JOURNAL(t *testing.T) {
	a := IS_USED_IN_JOURNAL("book")
	TEST(true, a, false)
}
func Test_ADD_ACCOUNT(t *testing.T) {
	a := ADD_ACCOUNT(ACCOUNT{
		IS_LOW_LEVEL_ACCOUNT:             false,
		IS_CREDIT:                        false,
		IS_TEMPORARY:                     false,
		COST_FLOW_TYPE:                   "",
		ACCOUNT_NAME:                     "assets",
		NOTES:                            "",
		IMAGE:                            []string{},
		BARCODE:                          []string{},
		ACCOUNT_NUMBER:                   [][]uint{{1}, {1}},
		ACCOUNT_LEVELS:                   []uint{},
		FATHER_AND_GRANDPA_ACCOUNTS_NAME: [][]string{},
		ALERT_FOR_MINIMUM_QUANTITY_BY_TURNOVER_IN_DAYS: 0,
		ALERT_FOR_MINIMUM_QUANTITY_BY_QUINTITY:         0,
		TARGET_BALANCE:                                 0,
		IF_THE_TARGET_BALANCE_IS_LESS_IS_GOOD:          false,
	})
	TEST(true, a, nil)
}
func Test_EDIT_ACCOUNT(t *testing.T) {
	// account_struct, index, err := ACCOUNT_STRUCT_FROM_NAME("book3")
	// fmt.Println(err)
	// if err == nil {
	// 	account_struct.ACCOUNT_NUMBER = [][]uint{{1, 3, 3}, {1, 3, 3}}
	// 	EDIT_ACCOUNT(false, index, account_struct)
	// }
	a := CHECK_IF_ACCOUNT_NUMBER_DUPLICATED()
	for _, a := range a {
		fmt.Println(a)
	}

	SET_THE_ACCOUNTS()
	// TEST(true,)
	PRINT_FORMATED_ACCOUNTS()
}
