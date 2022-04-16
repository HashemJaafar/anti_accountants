package anti_accountants

import (
	"errors"
	"fmt"
	"testing"
)

func Benchmark_ACCOUNT_STRUCT_FROM_BARCODE(b *testing.B) {
	//TODO
	//for i := 0; i < b.N; i++ {
	//	ACCOUNT_STRUCT_FROM_BARCODE()
	//}
}

func Example_ACCOUNT_STRUCT_FROM_BARCODE() {
	//TODO
	//a1:=ACCOUNT_STRUCT_FROM_BARCODE()
	//fmt.Println(a1)
}

func Fuzz_ACCOUNT_STRUCT_FROM_BARCODE(f *testing.F) {
	//TODO
	//f.Add()
	//f.Fuzz(func(t *testing.T){
	//	a1:=ACCOUNT_STRUCT_FROM_BARCODE()
	//	e1:=
	//	TEST(true,a1,e1)
	//})
}

func Test_ACCOUNT_STRUCT_FROM_BARCODE(t *testing.T) {
	INDEX_OF_ACCOUNT_NUMBER = 0
	account_struct, index, err := ACCOUNT_STRUCT_FROM_BARCODE("kaslajs")
	TEST(true, err, nil)
	TEST(true, index, 1)
	TEST(true, account_struct, ACCOUNT{false, false, false, "", "CURRENT_ASSETS", "", []string{}, []string{"sijadpodjpao", "kaslajs"}, [][]uint{{1, 1}, {}}, []uint{2, 0}, [][]string{{"ASSETS"}, {}}, 0, 0, 0, false})
}

func Benchmark_ACCOUNT_STRUCT_FROM_NAME(b *testing.B) {
	//TODO
	//for i := 0; i < b.N; i++ {
	//	ACCOUNT_STRUCT_FROM_NAME()
	//}
}

func Example_ACCOUNT_STRUCT_FROM_NAME() {
	//TODO
	//a1:=ACCOUNT_STRUCT_FROM_NAME()
	//fmt.Println(a1)
}

func Fuzz_ACCOUNT_STRUCT_FROM_NAME(f *testing.F) {
	//TODO
	//f.Add()
	//f.Fuzz(func(t *testing.T){
	//	a1:=ACCOUNT_STRUCT_FROM_NAME()
	//	e1:=
	//	TEST(true,a1,e1)
	//})
}

func Test_ACCOUNT_STRUCT_FROM_NAME(t *testing.T) {
	INDEX_OF_ACCOUNT_NUMBER = 0
	account_struct, index, err := ACCOUNT_STRUCT_FROM_NAME("ASSETS")
	TEST(true, err, nil)
	TEST(true, index, 0)
	TEST(true, account_struct, ACCOUNT{false, false, false, "", "ASSETS", "", []string{}, []string{"nojdsjdpq"}, [][]uint{{1}, {}}, []uint{1, 0}, [][]string{{}, {}}, 0, 0, 0, false})
}

func Benchmark_ADD_ACCOUNT(b *testing.B) {
	//TODO
	//for i := 0; i < b.N; i++ {
	//	ADD_ACCOUNT()
	//}
}

func Example_ADD_ACCOUNT() {
	//TODO
	//a1:=ADD_ACCOUNT()
	//fmt.Println(a1)
}

func Fuzz_ADD_ACCOUNT(f *testing.F) {
	//TODO
	//f.Add()
	//f.Fuzz(func(t *testing.T){
	//	a1:=ADD_ACCOUNT()
	//	e1:=
	//	TEST(true,a1,e1)
	//})
}

func Test_ADD_ACCOUNT(t *testing.T) {
	// DB_ACCOUNTS.DropAll()
	a := ADD_ACCOUNT(ACCOUNT{
		IS_LOW_LEVEL_ACCOUNT:             true,
		IS_CREDIT:                        false,
		IS_TEMPORARY:                     false,
		COST_FLOW_TYPE:                   WMA,
		ACCOUNT_NAME:                     "BOOK",
		NOTES:                            "",
		IMAGE:                            []string{},
		BARCODE:                          []string{"2"},
		ACCOUNT_NUMBER:                   [][]uint{{1, 3}},
		ACCOUNT_LEVELS:                   []uint{},
		FATHER_AND_GRANDPA_ACCOUNTS_NAME: [][]string{},
		ALERT_FOR_MINIMUM_QUANTITY_BY_TURNOVER_IN_DAYS: 0,
		ALERT_FOR_MINIMUM_QUANTITY_BY_QUINTITY:         0,
		TARGET_BALANCE:                                 0,
		IF_THE_TARGET_BALANCE_IS_LESS_IS_GOOD:          false,
	})
	PRINT_FORMATED_ACCOUNTS()
	DB_CLOSE()
	TEST(true, a, nil)
}

func Benchmark_CHECK_IF_ACCOUNT_NUMBER_DUPLICATED(b *testing.B) {
	//TODO
	//for i := 0; i < b.N; i++ {
	//	CHECK_IF_ACCOUNT_NUMBER_DUPLICATED()
	//}
}

func Example_CHECK_IF_ACCOUNT_NUMBER_DUPLICATED() {
	//TODO
	//a1:=CHECK_IF_ACCOUNT_NUMBER_DUPLICATED()
	//fmt.Println(a1)
}

func Fuzz_CHECK_IF_ACCOUNT_NUMBER_DUPLICATED(f *testing.F) {
	//TODO
	//f.Add()
	//f.Fuzz(func(t *testing.T){
	//	a1:=CHECK_IF_ACCOUNT_NUMBER_DUPLICATED()
	//	e1:=
	//	TEST(true,a1,e1)
	//})
}

func Test_CHECK_IF_ACCOUNT_NUMBER_DUPLICATED(t *testing.T) {
	a := CHECK_IF_ACCOUNT_NUMBER_DUPLICATED()
	fmt.Println(a)
	TEST(true, a, []error{errors.New("the account number [2] for {false false true fifo e  [] [] [[4] [2]] [] [] 0 0 0 false} duplicated")})
}

func Benchmark_CHECK_IF_LOW_LEVEL_ACCOUNT_FOR_ALL(b *testing.B) {
	//TODO
	//for i := 0; i < b.N; i++ {
	//	CHECK_IF_LOW_LEVEL_ACCOUNT_FOR_ALL()
	//}
}

func Example_CHECK_IF_LOW_LEVEL_ACCOUNT_FOR_ALL() {
	//TODO
	//a1:=CHECK_IF_LOW_LEVEL_ACCOUNT_FOR_ALL()
	//fmt.Println(a1)
}

func Fuzz_CHECK_IF_LOW_LEVEL_ACCOUNT_FOR_ALL(f *testing.F) {
	//TODO
	//f.Add()
	//f.Fuzz(func(t *testing.T){
	//	a1:=CHECK_IF_LOW_LEVEL_ACCOUNT_FOR_ALL()
	//	e1:=
	//	TEST(true,a1,e1)
	//})
}

func Test_CHECK_IF_LOW_LEVEL_ACCOUNT_FOR_ALL(t *testing.T) {
	a := CHECK_IF_LOW_LEVEL_ACCOUNT_FOR_ALL()
	TEST(true, a, []error{errors.New("should be low level account in all account numbers {false false false  b  [] [] [[1 1] [1 2]] [] [] 0 0 0 false}")})
}

func Benchmark_CHECK_IF_THE_TREE_CONNECTED(b *testing.B) {
	//TODO
	//for i := 0; i < b.N; i++ {
	//	CHECK_IF_THE_TREE_CONNECTED()
	//}
}

func Example_CHECK_IF_THE_TREE_CONNECTED() {
	//TODO
	//a1:=CHECK_IF_THE_TREE_CONNECTED()
	//fmt.Println(a1)
}

func Fuzz_CHECK_IF_THE_TREE_CONNECTED(f *testing.F) {
	//TODO
	//f.Add()
	//f.Fuzz(func(t *testing.T){
	//	a1:=CHECK_IF_THE_TREE_CONNECTED()
	//	e1:=
	//	TEST(true,a1,e1)
	//})
}

func Test_CHECK_IF_THE_TREE_CONNECTED(t *testing.T) {
	a := CHECK_IF_THE_TREE_CONNECTED()
	fmt.Println(a)
	TEST(true, a, []error{errors.New("the account number [2 1 8] for {true false true fifo f  [] [] [[4 1] [2 1 8]] [] [] 0 0 0 false} not conected to the tree")})
}

func Benchmark_CHECK_THE_TREE(b *testing.B) {
	//TODO
	//for i := 0; i < b.N; i++ {
	//	CHECK_THE_TREE()
	//}
}

func Example_CHECK_THE_TREE() {
	//TODO
	//a1:=CHECK_THE_TREE()
	//fmt.Println(a1)
}

func Fuzz_CHECK_THE_TREE(f *testing.F) {
	//TODO
	//f.Add()
	//f.Fuzz(func(t *testing.T){
	//	a1:=CHECK_THE_TREE()
	//	e1:=
	//	TEST(true,a1,e1)
	//})
}

func Test_CHECK_THE_TREE(t *testing.T) {
	CHECK_THE_TREE()
	TEST(true, ERRORS_MESSAGES, []error{
		errors.New("should be low level account in all account numbers {false false false  b  [] [] [[1 1] [1 2]] [] [] 0 0 0 false}"),
		errors.New("the account number [2 1 8] for {true false true fifo f  [] [] [[4 1] [2 1 8]] [] [] 0 0 0 false} not conected to the tree"),
	})
}

func Benchmark_EDIT_ACCOUNT(b *testing.B) {
	//TODO
	//for i := 0; i < b.N; i++ {
	//	EDIT_ACCOUNT()
	//}
}

func Example_EDIT_ACCOUNT() {
	//TODO
	//a1:=EDIT_ACCOUNT()
	//fmt.Println(a1)
}

func Fuzz_EDIT_ACCOUNT(f *testing.F) {
	//TODO
	//f.Add()
	//f.Fuzz(func(t *testing.T){
	//	a1:=EDIT_ACCOUNT()
	//	e1:=
	//	TEST(true,a1,e1)
	//})
}

func Test_EDIT_ACCOUNT(t *testing.T) {
	account_struct, index, err := ACCOUNT_STRUCT_FROM_NAME("rent")
	fmt.Println(err)
	if err == nil {
		account_struct.IS_CREDIT = true
		account_struct.IS_TEMPORARY = true
		account_struct.BARCODE = []string{"1"}
		account_struct.ACCOUNT_NUMBER = [][]uint{{4, 1}}
		account_struct.COST_FLOW_TYPE = WMA
		EDIT_ACCOUNT(false, index, account_struct)
	}
	DB_CLOSE()
	// TEST(true,)
	PRINT_FORMATED_ACCOUNTS()
}

func Benchmark_FORMAT_SLICE_OF_SLICE_OF_STRING_TO_STRING(b *testing.B) {
	//TODO
	//for i := 0; i < b.N; i++ {
	//	FORMAT_SLICE_OF_SLICE_OF_STRING_TO_STRING()
	//}
}

func Example_FORMAT_SLICE_OF_SLICE_OF_STRING_TO_STRING() {
	//TODO
	//a1:=FORMAT_SLICE_OF_SLICE_OF_STRING_TO_STRING()
	//fmt.Println(a1)
}

func Fuzz_FORMAT_SLICE_OF_SLICE_OF_STRING_TO_STRING(f *testing.F) {
	//TODO
	//f.Add()
	//f.Fuzz(func(t *testing.T){
	//	a1:=FORMAT_SLICE_OF_SLICE_OF_STRING_TO_STRING()
	//	e1:=
	//	TEST(true,a1,e1)
	//})
}

func Test_FORMAT_SLICE_OF_SLICE_OF_STRING_TO_STRING(t *testing.T) {
	//TODO
	//a1:=FORMAT_SLICE_OF_SLICE_OF_STRING_TO_STRING()
	//e1:=
	//TEST(true,a1,e1)
}

func Benchmark_FORMAT_SLICE_OF_SLICE_OF_UINT_TO_STRING(b *testing.B) {
	//TODO
	//for i := 0; i < b.N; i++ {
	//	FORMAT_SLICE_OF_SLICE_OF_UINT_TO_STRING()
	//}
}

func Example_FORMAT_SLICE_OF_SLICE_OF_UINT_TO_STRING() {
	//TODO
	//a1:=FORMAT_SLICE_OF_SLICE_OF_UINT_TO_STRING()
	//fmt.Println(a1)
}

func Fuzz_FORMAT_SLICE_OF_SLICE_OF_UINT_TO_STRING(f *testing.F) {
	//TODO
	//f.Add()
	//f.Fuzz(func(t *testing.T){
	//	a1:=FORMAT_SLICE_OF_SLICE_OF_UINT_TO_STRING()
	//	e1:=
	//	TEST(true,a1,e1)
	//})
}

func Test_FORMAT_SLICE_OF_SLICE_OF_UINT_TO_STRING(t *testing.T) {
	//TODO
	//a1:=FORMAT_SLICE_OF_SLICE_OF_UINT_TO_STRING()
	//e1:=
	//TEST(true,a1,e1)
}

func Benchmark_FORMAT_SLICE_OF_UINT_TO_STRING(b *testing.B) {
	//TODO
	//for i := 0; i < b.N; i++ {
	//	FORMAT_SLICE_OF_UINT_TO_STRING()
	//}
}

func Example_FORMAT_SLICE_OF_UINT_TO_STRING() {
	//TODO
	//a1:=FORMAT_SLICE_OF_UINT_TO_STRING()
	//fmt.Println(a1)
}

func Fuzz_FORMAT_SLICE_OF_UINT_TO_STRING(f *testing.F) {
	//TODO
	//f.Add()
	//f.Fuzz(func(t *testing.T){
	//	a1:=FORMAT_SLICE_OF_UINT_TO_STRING()
	//	e1:=
	//	TEST(true,a1,e1)
	//})
}

func Test_FORMAT_SLICE_OF_UINT_TO_STRING(t *testing.T) {
	//TODO
	//a1:=FORMAT_SLICE_OF_UINT_TO_STRING()
	//e1:=
	//TEST(true,a1,e1)
}

func Benchmark_FORMAT_STRING_SLICE_TO_STRING(b *testing.B) {
	//TODO
	//for i := 0; i < b.N; i++ {
	//	FORMAT_STRING_SLICE_TO_STRING()
	//}
}

func Example_FORMAT_STRING_SLICE_TO_STRING() {
	//TODO
	//a1:=FORMAT_STRING_SLICE_TO_STRING()
	//fmt.Println(a1)
}

func Fuzz_FORMAT_STRING_SLICE_TO_STRING(f *testing.F) {
	//TODO
	//f.Add()
	//f.Fuzz(func(t *testing.T){
	//	a1:=FORMAT_STRING_SLICE_TO_STRING()
	//	e1:=
	//	TEST(true,a1,e1)
	//})
}

func Test_FORMAT_STRING_SLICE_TO_STRING(t *testing.T) {
	//TODO
	//a1:=FORMAT_STRING_SLICE_TO_STRING()
	//e1:=
	//TEST(true,a1,e1)
}

func Benchmark_IS_BARCODES_USED(b *testing.B) {
	//TODO
	//for i := 0; i < b.N; i++ {
	//	IS_BARCODES_USED()
	//}
}

func Example_IS_BARCODES_USED() {
	//TODO
	//a1:=IS_BARCODES_USED()
	//fmt.Println(a1)
}

func Fuzz_IS_BARCODES_USED(f *testing.F) {
	//TODO
	//f.Add()
	//f.Fuzz(func(t *testing.T){
	//	a1:=IS_BARCODES_USED()
	//	e1:=
	//	TEST(true,a1,e1)
	//})
}

func Test_IS_BARCODES_USED(t *testing.T) {
	a := IS_BARCODES_USED([]string{"a", "b"})
	TEST(true, a, true)
	a = IS_BARCODES_USED([]string{"c", "b"})
	TEST(true, a, false)
}

func Benchmark_IS_IT_HIGH_THAN_BY_ORDER(b *testing.B) {
	//TODO
	//for i := 0; i < b.N; i++ {
	//	IS_IT_HIGH_THAN_BY_ORDER()
	//}
}

func Example_IS_IT_HIGH_THAN_BY_ORDER() {
	//TODO
	//a1:=IS_IT_HIGH_THAN_BY_ORDER()
	//fmt.Println(a1)
}

func Fuzz_IS_IT_HIGH_THAN_BY_ORDER(f *testing.F) {
	//TODO
	//f.Add()
	//f.Fuzz(func(t *testing.T){
	//	a1:=IS_IT_HIGH_THAN_BY_ORDER()
	//	e1:=
	//	TEST(true,a1,e1)
	//})
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

func Benchmark_IS_IT_POSSIBLE_TO_BE_SUB_ACCOUNT(b *testing.B) {
	//TODO
	//for i := 0; i < b.N; i++ {
	//	IS_IT_POSSIBLE_TO_BE_SUB_ACCOUNT()
	//}
}

func Example_IS_IT_POSSIBLE_TO_BE_SUB_ACCOUNT() {
	//TODO
	//a1:=IS_IT_POSSIBLE_TO_BE_SUB_ACCOUNT()
	//fmt.Println(a1)
}

func Fuzz_IS_IT_POSSIBLE_TO_BE_SUB_ACCOUNT(f *testing.F) {
	//TODO
	//f.Add()
	//f.Fuzz(func(t *testing.T){
	//	a1:=IS_IT_POSSIBLE_TO_BE_SUB_ACCOUNT()
	//	e1:=
	//	TEST(true,a1,e1)
	//})
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

func Benchmark_IS_IT_SUB_ACCOUNT_USING_NAME(b *testing.B) {
	//TODO
	//for i := 0; i < b.N; i++ {
	//	IS_IT_SUB_ACCOUNT_USING_NAME()
	//}
}

func Example_IS_IT_SUB_ACCOUNT_USING_NAME() {
	//TODO
	//a1:=IS_IT_SUB_ACCOUNT_USING_NAME()
	//fmt.Println(a1)
}

func Fuzz_IS_IT_SUB_ACCOUNT_USING_NAME(f *testing.F) {
	//TODO
	//f.Add()
	//f.Fuzz(func(t *testing.T){
	//	a1:=IS_IT_SUB_ACCOUNT_USING_NAME()
	//	e1:=
	//	TEST(true,a1,e1)
	//})
}

func Test_IS_IT_SUB_ACCOUNT_USING_NAME(t *testing.T) {
	INDEX_OF_ACCOUNT_NUMBER = 0
	a := IS_IT_SUB_ACCOUNT_USING_NAME("ASSETS", "CASH_AND_CASH_EQUIVALENTS")
	TEST(true, a, true)
	a = IS_IT_SUB_ACCOUNT_USING_NAME("CASH_AND_CASH_EQUIVALENTS", "ASSETS")
	TEST(true, a, false)
	INDEX_OF_ACCOUNT_NUMBER = 1
	a = IS_IT_SUB_ACCOUNT_USING_NAME("ASSETS", "CASH_AND_CASH_EQUIVALENTS")
	TEST(true, a, false)
}

func Benchmark_IS_IT_SUB_ACCOUNT_USING_NUMBER(b *testing.B) {
	//TODO
	//for i := 0; i < b.N; i++ {
	//	IS_IT_SUB_ACCOUNT_USING_NUMBER()
	//}
}

func Example_IS_IT_SUB_ACCOUNT_USING_NUMBER() {
	//TODO
	//a1:=IS_IT_SUB_ACCOUNT_USING_NUMBER()
	//fmt.Println(a1)
}

func Fuzz_IS_IT_SUB_ACCOUNT_USING_NUMBER(f *testing.F) {
	//TODO
	//f.Add()
	//f.Fuzz(func(t *testing.T){
	//	a1:=IS_IT_SUB_ACCOUNT_USING_NUMBER()
	//	e1:=
	//	TEST(true,a1,e1)
	//})
}

func Test_IS_IT_SUB_ACCOUNT_USING_NUMBER(t *testing.T) {
	a := IS_IT_SUB_ACCOUNT_USING_NUMBER([]uint{1}, []uint{1, 2})
	TEST(true, a, true)
	a = IS_IT_SUB_ACCOUNT_USING_NUMBER([]uint{1}, []uint{2})
	TEST(true, a, false)
	a = IS_IT_SUB_ACCOUNT_USING_NUMBER([]uint{}, []uint{2})
	TEST(true, a, false)
}

func Benchmark_IS_IT_THE_FATHER(b *testing.B) {
	//TODO
	//for i := 0; i < b.N; i++ {
	//	IS_IT_THE_FATHER()
	//}
}

func Example_IS_IT_THE_FATHER() {
	//TODO
	//a1:=IS_IT_THE_FATHER()
	//fmt.Println(a1)
}

func Fuzz_IS_IT_THE_FATHER(f *testing.F) {
	//TODO
	//f.Add()
	//f.Fuzz(func(t *testing.T){
	//	a1:=IS_IT_THE_FATHER()
	//	e1:=
	//	TEST(true,a1,e1)
	//})
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

func Benchmark_IS_USED_IN_JOURNAL(b *testing.B) {
	//TODO
	//for i := 0; i < b.N; i++ {
	//	IS_USED_IN_JOURNAL()
	//}
}

func Example_IS_USED_IN_JOURNAL() {
	//TODO
	//a1:=IS_USED_IN_JOURNAL()
	//fmt.Println(a1)
}

func Fuzz_IS_USED_IN_JOURNAL(f *testing.F) {
	//TODO
	//f.Add()
	//f.Fuzz(func(t *testing.T){
	//	a1:=IS_USED_IN_JOURNAL()
	//	e1:=
	//	TEST(true,a1,e1)
	//})
}

func Test_IS_USED_IN_JOURNAL(t *testing.T) {
	a := IS_USED_IN_JOURNAL("book")
	TEST(true, a, false)
}

func Benchmark_MAX_LEN_FOR_ACCOUNT_NUMBER(b *testing.B) {
	//TODO
	//for i := 0; i < b.N; i++ {
	//	MAX_LEN_FOR_ACCOUNT_NUMBER()
	//}
}

func Example_MAX_LEN_FOR_ACCOUNT_NUMBER() {
	//TODO
	//a1:=MAX_LEN_FOR_ACCOUNT_NUMBER()
	//fmt.Println(a1)
}

func Fuzz_MAX_LEN_FOR_ACCOUNT_NUMBER(f *testing.F) {
	//TODO
	//f.Add()
	//f.Fuzz(func(t *testing.T){
	//	a1:=MAX_LEN_FOR_ACCOUNT_NUMBER()
	//	e1:=
	//	TEST(true,a1,e1)
	//})
}

func Test_MAX_LEN_FOR_ACCOUNT_NUMBER(t *testing.T) {
	a := MAX_LEN_FOR_ACCOUNT_NUMBER()
	TEST(true, a, 2)
}

func Benchmark_PRINT_FORMATED_ACCOUNTS(b *testing.B) {
	//TODO
	//for i := 0; i < b.N; i++ {
	//	PRINT_FORMATED_ACCOUNTS()
	//}
}

func Example_PRINT_FORMATED_ACCOUNTS() {
	//TODO
	//a1:=PRINT_FORMATED_ACCOUNTS()
	//fmt.Println(a1)
}

func Fuzz_PRINT_FORMATED_ACCOUNTS(f *testing.F) {
	//TODO
	//f.Add()
	//f.Fuzz(func(t *testing.T){
	//	a1:=PRINT_FORMATED_ACCOUNTS()
	//	e1:=
	//	TEST(true,a1,e1)
	//})
}

func Test_PRINT_FORMATED_ACCOUNTS(t *testing.T) {
	PRINT_FORMATED_ACCOUNTS()
	//e1:=
	//TEST(true,a1,e1)
}

func Benchmark_SET_FATHER_AND_GRANDPA_ACCOUNTS_NAME(b *testing.B) {
	//TODO
	//for i := 0; i < b.N; i++ {
	//	SET_FATHER_AND_GRANDPA_ACCOUNTS_NAME()
	//}
}

func Example_SET_FATHER_AND_GRANDPA_ACCOUNTS_NAME() {
	//TODO
	//a1:=SET_FATHER_AND_GRANDPA_ACCOUNTS_NAME()
	//fmt.Println(a1)
}

func Fuzz_SET_FATHER_AND_GRANDPA_ACCOUNTS_NAME(f *testing.F) {
	//TODO
	//f.Add()
	//f.Fuzz(func(t *testing.T){
	//	a1:=SET_FATHER_AND_GRANDPA_ACCOUNTS_NAME()
	//	e1:=
	//	TEST(true,a1,e1)
	//})
}

func Test_SET_FATHER_AND_GRANDPA_ACCOUNTS_NAME(t *testing.T) {
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

func Benchmark_SET_THE_ACCOUNTS(b *testing.B) {
	//TODO
	//for i := 0; i < b.N; i++ {
	//	SET_THE_ACCOUNTS()
	//}
}

func Example_SET_THE_ACCOUNTS() {
	//TODO
	//a1:=SET_THE_ACCOUNTS()
	//fmt.Println(a1)
}

func Fuzz_SET_THE_ACCOUNTS(f *testing.F) {
	//TODO
	//f.Add()
	//f.Fuzz(func(t *testing.T){
	//	a1:=SET_THE_ACCOUNTS()
	//	e1:=
	//	TEST(true,a1,e1)
	//})
}

func Test_SET_THE_ACCOUNTS(t *testing.T) {
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

func Benchmark_SORT_THE_ACCOUNTS_BY_ACCOUNT_NUMBER(b *testing.B) {
	//TODO
	//for i := 0; i < b.N; i++ {
	//	SORT_THE_ACCOUNTS_BY_ACCOUNT_NUMBER()
	//}
}

func Example_SORT_THE_ACCOUNTS_BY_ACCOUNT_NUMBER() {
	//TODO
	//a1:=SORT_THE_ACCOUNTS_BY_ACCOUNT_NUMBER()
	//fmt.Println(a1)
}

func Fuzz_SORT_THE_ACCOUNTS_BY_ACCOUNT_NUMBER(f *testing.F) {
	//TODO
	//f.Add()
	//f.Fuzz(func(t *testing.T){
	//	a1:=SORT_THE_ACCOUNTS_BY_ACCOUNT_NUMBER()
	//	e1:=
	//	TEST(true,a1,e1)
	//})
}

func Test_SORT_THE_ACCOUNTS_BY_ACCOUNT_NUMBER(t *testing.T) {
	//TODO
	//a1:=SORT_THE_ACCOUNTS_BY_ACCOUNT_NUMBER()
	//e1:=
	//TEST(true,a1,e1)
}
