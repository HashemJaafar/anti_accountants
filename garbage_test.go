//this file created automatically
package main

import (
	"fmt"
	"strings"
	"testing"
)

func Test2(t *testing.T) {
	fmt.Println(len(FormatTheString("	   ")))
}

func Test1(t *testing.T) {
	AddUnderScore("remove the accounts not in accounts list")
}

func TestLower(t *testing.T) {
	fmt.Println(strings.ToLower("ALERT_FOR_MINIMUM_QUANTITY_BY_TURNOVER_IN_DAYS"))
	fmt.Println(strings.ToLower("ALERT_FOR_MINIMUM_QUANTITY_BY_QUINTITY"))
	fmt.Println(strings.ToLower("TARGET_BALANCE"))
	fmt.Println(strings.ToLower("IF_THE_TARGET_BALANCE_IS_LESS_IS_GOOD"))
}

func TestUpper(t *testing.T) {
	fmt.Println(strings.ToUpper(`anti accountants`))
}

func AddUnderScore(str string) {
	byt := []byte(str)
	for indexa, a := range byt {
		if a == ' ' {
			byt[indexa] = '_'
		}
	}
	fmt.Println(string(byt))
}

func RemoveUnderScore(str string) {
	byt := []byte(str)
	for indexa, a := range byt {
		if a == '_' {
			byt[indexa] = ' '
		}
	}
	fmt.Println(string(byt))
}
