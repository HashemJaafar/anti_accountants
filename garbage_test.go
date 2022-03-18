package anti_accountants

import (
	"fmt"
	"strings"
	"testing"
)

// SumIntsOrFloats sums the values of map m. It supports both int64 and float64
// as types for map values.
func SumIntsOrFloats[K comparable, V int64 | float64](m map[K]V) V {
    var s V
    for _, v := range m {
        s += v
    }
    return s
}

func Test_lower(t *testing.T) {
	fmt.Println(strings.ToLower("ALERT_FOR_MINIMUM_QUANTITY_BY_TURNOVER_IN_DAYS"))
	fmt.Println(strings.ToLower("ALERT_FOR_MINIMUM_QUANTITY_BY_QUINTITY"))
	fmt.Println(strings.ToLower("TARGET_BALANCE"))
	fmt.Println(strings.ToLower("IF_THE_TARGET_BALANCE_IS_LESS_IS_GOOD"))
}

func Test_a(t *testing.T){
	fmt.Printf("Generic Sums with Constraint: %v and %v\n",
    SumNumbers(ints),
    SumNumbers(floats))
}