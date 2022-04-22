package main

import (
	"fmt"
	"testing"
)

func TestFILTER_NUMBER(t *testing.T) {
	a1 := FILTER_NUMBER(4, 10, 5, false)
	fmt.Println(a1)
}
