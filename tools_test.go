package anti_accountants

import (
	"fmt"
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

func TestFTest(t *testing.T) {
	FTest(false, SAccount1{}, SAccount1{})
}

func Test1(t *testing.T) {
	fmt.Println(len(`goroutine 141 [running]:
runtime/debug.Stack()
        /usr/local/go/src/runtime/debug/stack.go:24 +0x65
anti_accountants.FTest[...](_, {0x0, {0x0, 0x0}, {0x0, 0x0}, {0x0, 0x0}, {0x0, 0x0, ...}, ...}, ...)
        /home/hashem/anti_accountants/tools.go:160 +0x495
anti_accountants.TestFTest(0x0?)
        /home/hashem/anti_accountants/tools_test.go:20 +0x89`))

	fmt.Println(len(`goroutine 141 [running]:
		runtime/debug.Stack()
				/usr/local/go/src/runtime/debug/stack.go:24 +0x65
		anti_accountants.FTest[...](_, {0x0, {0x0, 0x0}, {0x0, 0x0}, {0x0, 0x0}, {0x0, 0x0, ...}, ...}, ...)
				/home/hashem/anti_accountants/tools.go:160 +0x495
		anti_accountants.TestFTest(0x0?)`))
}

func Test2(t *testing.T) {
	fmt.Println('\n')
}
