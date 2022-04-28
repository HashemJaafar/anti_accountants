package main

import "testing"

func TestFUNCTION_FILTER_DUPLICATE(t *testing.T) {
	a1 := FUNCTION_FILTER_DUPLICATE("lkjds", "ojdi", true)
	TEST(true, a1, false)
	a1 = FUNCTION_FILTER_DUPLICATE(4496, 546, true)
	TEST(true, a1, false)
	a1 = FUNCTION_FILTER_DUPLICATE(true, true, true)
	TEST(true, a1, true)
	a1 = FUNCTION_FILTER_DUPLICATE("lkjds", "ojdi", false)
	TEST(true, a1, true)
}
