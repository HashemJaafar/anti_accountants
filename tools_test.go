package anti_accountants

import (
	"fmt"
	"testing"
)

func Test_REVERSE_SLICE(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5, 6}
	REVERSE_SLICE(slice)
	fmt.Println(slice)
}

func Test_TRANSPOSE(t *testing.T) {
	slice := [][]int{{1, 2, 3}, {4, 5, 6}}
	slice = TRANSPOSE(slice)
	fmt.Println(slice)
}

func Test_TEST(t *testing.T) {
	TEST(false, true, false)
	TEST(true, true, true)
}
