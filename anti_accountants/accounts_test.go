package anti_accountants

import (
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
	TEST_FUNCTION(err, nil)
	TEST_FUNCTION(index, 0)
	TEST_FUNCTION(account_struct, ACCOUNT{false, false, false, "", "ASSETS", "", []string{}, []string{"nojdsjdpq"}, [][]uint{{1}, {}}, []uint{1, 0}, [][]string{{}, {}}, 0, 0, 0, false})
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
	TEST_FUNCTION(err, nil)
	TEST_FUNCTION(index, 1)
	TEST_FUNCTION(account_struct, ACCOUNT{false, false, false, "", "CURRENT_ASSETS", "", []string{}, []string{"sijadpodjpao", "kaslajs"}, [][]uint{{1, 1}, {}}, []uint{2, 0}, [][]string{{"ASSETS"}, {}}, 0, 0, 0, false})
}

func Test_ACCOUNT_STRUCT_FROM_NUMBER(t *testing.T) {
	ACCOUNTS = []ACCOUNT{
		{false, false, false, "", "ASSETS", "", []string{}, []string{"nojdsjdpq"}, [][]uint{{1}, {}}, []uint{1, 0}, [][]string{{}, {}}, 0, 0, 0, false},
		{false, false, false, "", "CURRENT_ASSETS", "", []string{}, []string{"sijadpodjpao", "kaslajs"}, [][]uint{{1, 1}, {}}, []uint{2, 0}, [][]string{{"ASSETS"}, {}}, 0, 0, 0, false},
		{true, false, false, "fifo", "CASH_AND_CASH_EQUIVALENTS", "", []string{}, []string{"888"}, [][]uint{{1, 1, 1}, {2}}, []uint{3, 1}, [][]string{{"ASSETS", "CURRENT_ASSETS"}, {}}, 0, 0, 0, false},
		{true, false, false, "fifo", "SHORT_TERM_INVESTMENTS", "", []string{}, []string{"SHORT_TERM_INVESTMENTS"}, [][]uint{{1, 2}, {5}}, []uint{2, 1}, [][]string{{"ASSETS"}, {}}, 0, 0, 0, false},
		{true, false, false, "", "RECEIVABLES", "", []string{}, []string{"RECEIVABLES"}, [][]uint{{1, 3}, {}}, []uint{2, 0}, [][]string{{"ASSETS"}, {}}, 0, 0, 0, false},
	}
	INDEX_OF_ACCOUNT_NUMBER = 1
	account_struct, index, err := ACCOUNT_STRUCT_FROM_NUMBER([]uint{2})
	TEST_FUNCTION(err, nil)
	TEST_FUNCTION(index, 2)
	TEST_FUNCTION(account_struct, ACCOUNT{true, false, false, "fifo", "CASH_AND_CASH_EQUIVALENTS", "", []string{}, []string{"888"}, [][]uint{{1, 1, 1}, {2}}, []uint{3, 1}, [][]string{{"ASSETS", "CURRENT_ASSETS"}, {}}, 0, 0, 0, false})
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
	TEST_FUNCTION(a, true)
	a = IS_IT_SUB_ACCOUNT_USING_NAME("CASH_AND_CASH_EQUIVALENTS", "ASSETS")
	TEST_FUNCTION(a, false)
	INDEX_OF_ACCOUNT_NUMBER = 1
	a = IS_IT_SUB_ACCOUNT_USING_NAME("ASSETS", "CASH_AND_CASH_EQUIVALENTS")
	TEST_FUNCTION(a, false)
}

func Test_IS_IT_SUB_ACCOUNT_USING_NUMBER(t *testing.T) {
	a := IS_IT_SUB_ACCOUNT_USING_NUMBER([]uint{1}, []uint{1, 2})
	TEST_FUNCTION(a, true)
	a = IS_IT_SUB_ACCOUNT_USING_NUMBER([]uint{1}, []uint{2})
	TEST_FUNCTION(a, false)
	a = IS_IT_SUB_ACCOUNT_USING_NUMBER([]uint{}, []uint{2})
	TEST_FUNCTION(a, false)
}

func Test_IS_IT_POSSIBLE_TO_BE_SUB_ACCOUNT(t *testing.T) {
	a := IS_IT_POSSIBLE_TO_BE_SUB_ACCOUNT([]uint{1}, []uint{1, 2})
	TEST_FUNCTION(a, true)
	a = IS_IT_POSSIBLE_TO_BE_SUB_ACCOUNT([]uint{1}, []uint{2})
	TEST_FUNCTION(a, false)
	a = IS_IT_POSSIBLE_TO_BE_SUB_ACCOUNT([]uint{}, []uint{2})
	TEST_FUNCTION(a, false)
	a = IS_IT_POSSIBLE_TO_BE_SUB_ACCOUNT([]uint{1}, []uint{1, 1, 2})
	TEST_FUNCTION(a, true)
}

func Test_IS_IT_THE_FATHER(t *testing.T) {
	a := IS_IT_THE_FATHER([]uint{1}, []uint{1, 2})
	TEST_FUNCTION(a, true)
	a = IS_IT_THE_FATHER([]uint{1}, []uint{2})
	TEST_FUNCTION(a, false)
	a = IS_IT_THE_FATHER([]uint{}, []uint{2})
	TEST_FUNCTION(a, false)
	a = IS_IT_THE_FATHER([]uint{1}, []uint{1, 1, 2})
	TEST_FUNCTION(a, false)
}

func Test_1(t *testing.T) {
	ACCOUNTS = []ACCOUNT{
		{false, false, false, "", "a", "", []string{}, []string{}, [][]uint{{1}, {1}}, []uint{}, [][]string{}, 0, 0, 0, false},
		{false, false, false, "", "b", "", []string{}, []string{}, [][]uint{{1, 1}, {1, 2}}, []uint{}, [][]string{}, 0, 0, 0, false},
		{false, false, false, "", "c", "", []string{}, []string{}, [][]uint{{1, 2}, {1, 3}}, []uint{}, [][]string{}, 0, 0, 0, false},
		{false, false, true, "fifo", "d", "", []string{}, []string{}, [][]uint{{2}, {2}}, []uint{}, [][]string{}, 0, 0, 0, false},
		{false, false, true, "fifo", "e", "", []string{}, []string{}, [][]uint{{4}, {2}}, []uint{}, [][]string{}, 0, 0, 0, false},
		{false, false, true, "fifo", "f", "", []string{}, []string{}, [][]uint{{4, 1}, {2, 1, 8}}, []uint{}, [][]string{}, 0, 0, 0, false},
	}
	INIT_ACCOUNT_NUMBERS_AND_FATHER_AND_GRANDPA_ACCOUNTS_NAME()

	SORT_THE_ACCOUNTS_BY_ACCOUNT_NUMBER()
	PRINT_FORMATED_ACCOUNTS()
	fmt.Println("//////////////////////SORT_THE_ACCOUNTS_BY_ACCOUNT_NUMBER")

	PRINT_FORMATED_ACCOUNTS()
	fmt.Println("//////////////////////SET_LOW_LEVEL_ACCOUNTS")

	SET_ACCOUNT_LEVELS()
	PRINT_FORMATED_ACCOUNTS()
	fmt.Println("//////////////////////SET_ACCOUNT_LEVELS")

	CHECK_THE_TREE()
	for _, a := range ERRORS_MESSAGES {
		fmt.Println(a)
	}
	PRINT_FORMATED_ACCOUNTS()
	fmt.Println("//////////////////////CHECK_THE_TREE")

	SET_HIGH_LEVEL_ACCOUNT_TO_PERMANENT()
	PRINT_FORMATED_ACCOUNTS()
	fmt.Println("//////////////////////SET_HIGH_LEVEL_ACCOUNT_TO_PERMANENT")

	SET_COST_FLOW_TYPE()
	PRINT_FORMATED_ACCOUNTS()
	fmt.Println("//////////////////////SET_COST_FLOW_TYPE")

	SET_FATHER_AND_GRANDPA_ACCOUNTS_NAME()
	PRINT_FORMATED_ACCOUNTS()
	fmt.Println("//////////////////////SET_FATHER_AND_GRANDPA_ACCOUNTS_NAME")
}

func Test_IS_IT_HIGH_THAN_BY_ORDER(t *testing.T) {
	a := IS_IT_HIGH_THAN_BY_ORDER([]uint{1}, []uint{1, 2})
	TEST_FUNCTION(a, true)
	a = IS_IT_HIGH_THAN_BY_ORDER([]uint{1}, []uint{1})
	TEST_FUNCTION(a, false)
	a = IS_IT_HIGH_THAN_BY_ORDER([]uint{1, 2}, []uint{1})
	TEST_FUNCTION(a, false)
	a = IS_IT_HIGH_THAN_BY_ORDER([]uint{3}, []uint{1, 1})
	TEST_FUNCTION(a, false)
	a = IS_IT_HIGH_THAN_BY_ORDER([]uint{1, 5}, []uint{3})
	TEST_FUNCTION(a, true)
}
