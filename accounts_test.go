package anti_accountants

import (
	"log"
	"reflect"
	"testing"
)

func test_function(a, e interface{}) {
	if !reflect.DeepEqual(a, e) {
		log.Fatal("expected : ", e, " actual : ", a)
	}
}

func Test_account_struct_from_name(t *testing.T) {
	a := account_struct_from_name("ASSETS")
	e := ACCOUNT{
		is_low_level_account: false,
		IS_CREDIT:            false,
		IS_TEMPORARY:         false,
		COST_FLOW_TYPE:       "",
		ACCOUNT_NAME:         "ASSETS",
		DESCRIPTION:          "",
		ACCOUNT_NUMBER:       [][]uint{{1}},
		IMAGE:                "",
	}
	test_function(a, e)
}

func Test_account_struct_from_number(t *testing.T) {
	a := account_struct_from_number([]uint{1})
	e := ACCOUNT{
		is_low_level_account: false,
		IS_CREDIT:            false,
		IS_TEMPORARY:         false,
		COST_FLOW_TYPE:       "",
		ACCOUNT_NAME:         "ASSETS",
		DESCRIPTION:          "",
		ACCOUNT_NUMBER:       [][]uint{{1}},
		IMAGE:                "",
	}
	test_function(a, e)
}

func Test_is_it_sub_account_using_name(t *testing.T) {
	a := is_it_sub_account_using_name("ASSETS", "CURRENT_ASSETS")
	e := true
	test_function(a, e)
	a = is_it_sub_account_using_name("CURRENT_ASSETS", "ASSETS")
	e = false
	test_function(a, e)
}

func Test_is_it_sub_account_using_number(t *testing.T) {
	a := is_it_sub_account_using_number([]uint{1}, []uint{1, 1})
	e := true
	test_function(a, e)
	a = is_it_sub_account_using_number([]uint{1, 1}, []uint{1})
	e = false
	test_function(a, e)
}

func Test_is_it_high_by_level(t *testing.T) {
	// a := is_it_high_by_level("ASSETS")
	// e := true
	// test_function(a, e)
	// a = is_it_high_by_level("CASH_AND_CASH_EQUIVALENTS")
	// e = false
	// test_function(a, e)
}

func Test_find_all_higher_level_accounts(t *testing.T) {
}

func Test_is_it_first_sub_level_account_using_number(t *testing.T) {
}

func Test_check_if_the_tree_connected(t *testing.T) {
}

func Test_check_cost_flow_type(t *testing.T) {
}

func Test_check_if_duplicated(t *testing.T) {
}

func Test_inventory_accounts(t *testing.T) {
}

func Test_SORT_THE_ACCOUNT_BY_ACCOUNT_NUMBER(t *testing.T) {
	SORT_THE_ACCOUNT_BY_ACCOUNT_NUMBER()
}

func Test_print_formated_accounts(t *testing.T) {
}

func Test_is_it_high_than_by_order(t *testing.T) {
}

func Test_is_shorter_than(t *testing.T) {
}
