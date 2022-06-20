package main

import (
	"testing"
)

func TestFilter(t *testing.T) {
	a1 := FFilterDuplicate("lkjds", "ojdi", true)
	FTest(true, a1, false)
	a1 = FFilterDuplicate(4496, 546, true)
	FTest(true, a1, false)
	a1 = FFilterDuplicate(true, true, true)
	FTest(true, a1, true)
	a1 = FFilterDuplicate("lkjds", "ojdi", false)
	FTest(true, a1, true)
}
