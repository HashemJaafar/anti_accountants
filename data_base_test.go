package anti_accountants

import (
	"fmt"
	"testing"
)

func Benchmark_ACCOUNT_BALANCE(b *testing.B) {
	//TODO
	//for i := 0; i < b.N; i++ {
	//	ACCOUNT_BALANCE()
	//}
}

func Example_ACCOUNT_BALANCE() {
	//TODO
	//a1:=ACCOUNT_BALANCE()
	//fmt.Println(a1)
}

func Fuzz_ACCOUNT_BALANCE(f *testing.F) {
	//TODO
	//f.Add()
	//f.Fuzz(func(t *testing.T){
	//	a1:=ACCOUNT_BALANCE()
	//	e1:=
	//	TEST(true,a1,e1)
	//})
}

func Test_ACCOUNT_BALANCE(t *testing.T) {
	a1 := ACCOUNT_BALANCE("cash")
	e1 := 0.0
	TEST(true, a1, e1)
}

func Benchmark_DB_CLOSE(b *testing.B) {
	//TODO
	//for i := 0; i < b.N; i++ {
	//	DB_CLOSE()
	//}
}

func Example_DB_CLOSE() {
	//TODO
	//a1:=DB_CLOSE()
	//fmt.Println(a1)
}

func Fuzz_DB_CLOSE(f *testing.F) {
	//TODO
	//f.Add()
	//f.Fuzz(func(t *testing.T){
	//	a1:=DB_CLOSE()
	//	e1:=
	//	TEST(true,a1,e1)
	//})
}

func Test_DB_CLOSE(t *testing.T) {
	TEST(true, DB_ACCOUNTS.IsClosed(), false)
	TEST(true, DB_JOURNAL.IsClosed(), false)
	TEST(true, DB_INVENTORY.IsClosed(), false)
	DB_CLOSE()
	TEST(true, DB_ACCOUNTS.IsClosed(), true)
	TEST(true, DB_JOURNAL.IsClosed(), true)
	TEST(true, DB_INVENTORY.IsClosed(), true)
}

func Benchmark_DB_INSERT(b *testing.B) {
	//TODO
	//for i := 0; i < b.N; i++ {
	//	DB_INSERT()
	//}
}

func Example_DB_INSERT() {
	//TODO
	//a1:=DB_INSERT()
	//fmt.Println(a1)
}

func Fuzz_DB_INSERT(f *testing.F) {
	//TODO
	//f.Add()
	//f.Fuzz(func(t *testing.T){
	//	a1:=DB_INSERT()
	//	e1:=
	//	TEST(true,a1,e1)
	//})
}

func Test_DB_INSERT(t *testing.T) {
	DB_INVENTORY.DropAll()
	DB_INSERT(DB_INVENTORY, []INVENTORY_TAG{
		{1, 10, "book"},
		{2, 10, "book"},
		{3, 10, "book"},
		{4, 10, "book"},
		{1, 10, "cash"},
		{1, 10, "cash"},
		{2, 10, "rent"},
		{9, 10, "rent"},
	})
	_, inventory := DB_READ[INVENTORY_TAG](DB_INVENTORY)
	for _, v1 := range inventory {
		fmt.Println(v1)
	}
	DB_CLOSE()
}

func Benchmark_DB_INSERT_INTO_ACCOUNTS(b *testing.B) {
	//TODO
	//for i := 0; i < b.N; i++ {
	//	DB_INSERT_INTO_ACCOUNTS()
	//}
}

func Example_DB_INSERT_INTO_ACCOUNTS() {
	//TODO
	//a1:=DB_INSERT_INTO_ACCOUNTS()
	//fmt.Println(a1)
}

func Fuzz_DB_INSERT_INTO_ACCOUNTS(f *testing.F) {
	//TODO
	//f.Add()
	//f.Fuzz(func(t *testing.T){
	//	a1:=DB_INSERT_INTO_ACCOUNTS()
	//	e1:=
	//	TEST(true,a1,e1)
	//})
}

func Test_DB_INSERT_INTO_ACCOUNTS(t *testing.T) {
}

func Benchmark_DB_LAST_LINE(b *testing.B) {
	//TODO
	//for i := 0; i < b.N; i++ {
	//	DB_LAST_LINE()
	//}
}

func Example_DB_LAST_LINE() {
	//TODO
	//a1:=DB_LAST_LINE()
	//fmt.Println(a1)
}

func Fuzz_DB_LAST_LINE(f *testing.F) {
	//TODO
	//f.Add()
	//f.Fuzz(func(t *testing.T){
	//	a1:=DB_LAST_LINE()
	//	e1:=
	//	TEST(true,a1,e1)
	//})
}

func Test_DB_LAST_LINE(t *testing.T) {
	a1 := DB_LAST_LINE[JOURNAL_TAG](DB_JOURNAL)
	fmt.Println(a1)
	// e1 := JOURNAL_TAG{false, 0, 0, 0, 0, 0, 0, 0, 0, "", "", "", "", "", time.Time{}, time.Time{}}
	// TEST(true, a1, e1)
}

func Benchmark_DB_OPEN(b *testing.B) {
	//TODO
	//for i := 0; i < b.N; i++ {
	//	DB_OPEN()
	//}
}

func Example_DB_OPEN() {
	//TODO
	//a1:=DB_OPEN()
	//fmt.Println(a1)
}

func Fuzz_DB_OPEN(f *testing.F) {
	//TODO
	//f.Add()
	//f.Fuzz(func(t *testing.T){
	//	a1:=DB_OPEN()
	//	e1:=
	//	TEST(true,a1,e1)
	//})
}

func Test_DB_OPEN(t *testing.T) {
	//TODO
	//a1:=DB_OPEN()
	//e1:=
	//TEST(true,a1,e1)
}

func Benchmark_DB_READ(b *testing.B) {
	//TODO
	//for i := 0; i < b.N; i++ {
	//	DB_READ()
	//}
}

func Example_DB_READ() {
	//TODO
	//a1:=DB_READ()
	//fmt.Println(a1)
}

func Fuzz_DB_READ(f *testing.F) {
	//TODO
	//f.Add()
	//f.Fuzz(func(t *testing.T){
	//	a1:=DB_READ()
	//	e1:=
	//	TEST(true,a1,e1)
	//})
}

func Test_DB_READ(t *testing.T) {
	_, inventory := DB_READ[INVENTORY_TAG](DB_INVENTORY)
	_, journal := DB_READ[JOURNAL_TAG](DB_JOURNAL)
	DB_CLOSE()
	PRINT_SLICE(inventory)
	PRINT_SLICE(journal)
}

func Benchmark_DB_UPDATE(b *testing.B) {
	// db := DB_OPEN("./test")
	// for i := 0; i < 1000; i++ {
	// 	DB_UPDATE(db, NOW(), i)
	// }
}

func Example_DB_UPDATE() {
	//TODO
	//a1:=DB_UPDATE()
	//fmt.Println(a1)
}

func Fuzz_DB_UPDATE(f *testing.F) {
	//TODO
	//f.Add()
	//f.Fuzz(func(t *testing.T){
	//	a1:=DB_UPDATE()
	//	e1:=
	//	TEST(true,a1,e1)
	//})
}

func Test_DB_UPDATE(t *testing.T) {
	DB_UPDATE(DB_INVENTORY, NOW(), INVENTORY_TAG{1, 10, "book1"})
}

func Benchmark_DB_UPDATE_ACCOUNT_NAME_IN_INVENTORY(b *testing.B) {
	//TODO
	//for i := 0; i < b.N; i++ {
	//	DB_UPDATE_ACCOUNT_NAME_IN_INVENTORY()
	//}
}

func Example_DB_UPDATE_ACCOUNT_NAME_IN_INVENTORY() {
	//TODO
	//a1:=DB_UPDATE_ACCOUNT_NAME_IN_INVENTORY()
	//fmt.Println(a1)
}

func Fuzz_DB_UPDATE_ACCOUNT_NAME_IN_INVENTORY(f *testing.F) {
	//TODO
	//f.Add()
	//f.Fuzz(func(t *testing.T){
	//	a1:=DB_UPDATE_ACCOUNT_NAME_IN_INVENTORY()
	//	e1:=
	//	TEST(true,a1,e1)
	//})
}

func Test_DB_UPDATE_ACCOUNT_NAME_IN_INVENTORY(t *testing.T) {
	//TODO
	//a1:=DB_UPDATE_ACCOUNT_NAME_IN_INVENTORY()
	//e1:=
	//TEST(true,a1,e1)
}

func Benchmark_DB_UPDATE_ACCOUNT_NAME_IN_JOURNAL(b *testing.B) {
	//TODO
	//for i := 0; i < b.N; i++ {
	//	DB_UPDATE_ACCOUNT_NAME_IN_JOURNAL()
	//}
}

func Example_DB_UPDATE_ACCOUNT_NAME_IN_JOURNAL() {
	//TODO
	//a1:=DB_UPDATE_ACCOUNT_NAME_IN_JOURNAL()
	//fmt.Println(a1)
}

func Fuzz_DB_UPDATE_ACCOUNT_NAME_IN_JOURNAL(f *testing.F) {
	//TODO
	//f.Add()
	//f.Fuzz(func(t *testing.T){
	//	a1:=DB_UPDATE_ACCOUNT_NAME_IN_JOURNAL()
	//	e1:=
	//	TEST(true,a1,e1)
	//})
}

func Test_DB_UPDATE_ACCOUNT_NAME_IN_JOURNAL(t *testing.T) {
	//TODO
	//a1:=DB_UPDATE_ACCOUNT_NAME_IN_JOURNAL()
	//e1:=
	//TEST(true,a1,e1)
}

func Benchmark_WEIGHTED_AVERAGE(b *testing.B) {
	//TODO
	//for i := 0; i < b.N; i++ {
	//	WEIGHTED_AVERAGE()
	//}
}

func Example_WEIGHTED_AVERAGE() {
	//TODO
	//a1:=WEIGHTED_AVERAGE()
	//fmt.Println(a1)
}

func Fuzz_WEIGHTED_AVERAGE(f *testing.F) {
	//TODO
	//f.Add()
	//f.Fuzz(func(t *testing.T){
	//	a1:=WEIGHTED_AVERAGE()
	//	e1:=
	//	TEST(true,a1,e1)
	//})
}

func Test_WEIGHTED_AVERAGE(t *testing.T) {
	WEIGHTED_AVERAGE("rent")
	_, inventory := DB_READ[INVENTORY_TAG](DB_INVENTORY)
	DB_CLOSE()
	PRINT_SLICE(inventory)
	//e1:=
	//TEST(true,a1,e1)
}
