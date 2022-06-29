package anti_accountants

import (
	"fmt"
	"log"
	"os"
	"testing"
	"text/tabwriter"
)

func TestAddAccount(t *testing.T) {
	var a1 SAccount3

	a1 = FAddAccount(true, SAccount1{
		IsCredit:     false,
		CostFlowType: CHighLevelAccount,
		Name:         "assets",
		Number:       [][]uint{{1}},
	})
	FTest(true, a1, SAccount3{})

	a1 = FAddAccount(true, SAccount1{
		IsCredit:     false,
		CostFlowType: CHighLevelAccount,
		Name:         "current assets",
		Number:       [][]uint{{1, 1}},
	})
	FTest(true, a1, SAccount3{})

	a1 = FAddAccount(true, SAccount1{
		IsCredit:     false,
		CostFlowType: CWma,
		Name:         "cash",
		Number:       [][]uint{{1, 1, 1}},
	})
	FTest(true, a1, SAccount3{})

	a1 = FAddAccount(true, SAccount1{
		IsCredit:     false,
		CostFlowType: CHighLevelAccount,
		Name:         "invnetory",
		Number:       [][]uint{{1, 1, 3}},
	})
	FTest(true, a1, SAccount3{})

	a1 = FAddAccount(true, SAccount1{
		IsCredit:     false,
		CostFlowType: CFifo,
		Name:         "book",
		Number:       [][]uint{{1, 1, 3, 1}},
	})
	FTest(true, a1, SAccount3{})

	a1 = FAddAccount(true, SAccount1{
		IsCredit:     true,
		CostFlowType: CHighLevelAccount,
		Name:         "liabilities",
		Number:       [][]uint{{2}},
	})
	FTest(true, a1, SAccount3{})

	a1 = FAddAccount(true, SAccount1{
		IsCredit:     true,
		CostFlowType: CHighLevelAccount,
		Name:         "owner's equity",
		Number:       [][]uint{{3}},
	})
	FTest(true, a1, SAccount3{})

	a1 = FAddAccount(true, SAccount1{
		IsCredit:     true,
		CostFlowType: CHighLevelAccount,
		Name:         "retained earnings",
		Number:       [][]uint{{3, 1}},
	})
	FTest(true, a1, SAccount3{})

	a1 = FAddAccount(true, SAccount1{
		IsCredit:     true,
		CostFlowType: CHighLevelAccount,
		Name:         "income",
		Number:       [][]uint{{3, 2}},
	})
	FTest(true, a1, SAccount3{})

	a1 = FAddAccount(true, SAccount1{
		IsCredit:     true,
		CostFlowType: CHighLevelAccount,
		Name:         "revenue",
		Number:       [][]uint{{3, 2, 1}},
	})
	FTest(true, a1, SAccount3{})

	a1 = FAddAccount(true, SAccount1{
		IsCredit:     true,
		CostFlowType: CWma,
		Name:         "rent",
		Number:       [][]uint{{3, 2, 1, 1}},
	})
	FTest(true, a1, SAccount3{})

	a1 = FAddAccount(true, SAccount1{
		IsCredit:     false,
		CostFlowType: CHighLevelAccount,
		Name:         "expense",
		Number:       [][]uint{{3, 2, 2}},
	})
	FTest(true, a1, SAccount3{})

	a1 = FAddAccount(true, SAccount1{
		IsCredit:     false,
		CostFlowType: CHighLevelAccount,
		Name:         "discounts",
		Number:       [][]uint{{3, 2, 2, 1}},
	})
	FTest(true, a1, SAccount3{})

	a1 = FAddAccount(true, SAccount1{
		IsCredit:     false,
		CostFlowType: CWma,
		Name:         "Invoice PQ",
		Number:       [][]uint{{3, 2, 2, 1, 1}},
	})
	FTest(true, a1, SAccount3{})

	a1 = FAddAccount(true, SAccount1{
		IsCredit:     false,
		CostFlowType: CWma,
		Name:         CPrefixDiscount + "book",
		Number:       [][]uint{{3, 2, 2, 1, 2}},
	})
	FTest(true, a1, SAccount3{})

	a1 = FAddAccount(true, SAccount1{
		IsCredit:     false,
		CostFlowType: CHighLevelAccount,
		Name:         "cost of goods sold",
		Number:       [][]uint{{3, 2, 2, 2}},
	})
	FTest(true, a1, SAccount3{})

	a1 = FAddAccount(true, SAccount1{
		IsCredit:     false,
		CostFlowType: CWma,
		Name:         CPrefixCost + "book",
		Number:       [][]uint{{3, 2, 2, 2, 1}},
	})
	FTest(true, a1, SAccount3{})
}

func TestEditAccount(t *testing.T) {
	account, index, err := FFindAccountFromName("cash")
	fmt.Println(err)
	if err == nil {
		account.IsCredit = false
		FEditAccount(false, index, account)
	}
}

func TestIsItHighThanByOrder(t *testing.T) {
	FTest(true, FIsItHighThanByOrder([]uint{1}, []uint{1, 2}), true)
	FTest(true, FIsItHighThanByOrder([]uint{1}, []uint{2}), true)
	FTest(true, FIsItHighThanByOrder([]uint{}, []uint{2}), true)
	FTest(true, FIsItHighThanByOrder([]uint{1, 1}, []uint{1, 1, 2}), true)
	FTest(true, FIsItHighThanByOrder([]uint{1}, []uint{2, 1}), true)
	FTest(true, FIsItHighThanByOrder([]uint{2, 1}, []uint{1, 1}), false)
	FTest(true, FIsItHighThanByOrder([]uint{1, 1}, []uint{1}), false)
}

func TestIsItSubAccountUsingNumber(t *testing.T) {
	FTest(true, FIsItSubAccountUsingNumber([]uint{1}, []uint{1, 2}), true)
	FTest(true, FIsItSubAccountUsingNumber([]uint{1}, []uint{2}), false)
	FTest(true, FIsItSubAccountUsingNumber([]uint{}, []uint{2}), false)
	FTest(true, FIsItSubAccountUsingNumber([]uint{1}, []uint{1, 1, 2}), true)
	FTest(true, FIsItSubAccountUsingNumber([]uint{1}, []uint{2, 1}), false)
	FTest(true, FIsItSubAccountUsingNumber([]uint{1}, []uint{1, 1}), true)
}

func TestIsItTheFather(t *testing.T) {
	FTest(true, FIsItTheFather([]uint{1}, []uint{1, 2}), true)
	FTest(true, FIsItTheFather([]uint{1}, []uint{2}), false)
	FTest(true, FIsItTheFather([]uint{}, []uint{2}), false)
	FTest(true, FIsItTheFather([]uint{1}, []uint{1, 1, 2}), false)
	FTest(true, FIsItTheFather([]uint{1}, []uint{2, 1}), false)
	FTest(true, FIsItTheFather([]uint{1}, []uint{1, 1}), true)
}

func TestIsUsedInJournal(t *testing.T) {
	a := FIsUsedInJournal("book")
	FTest(true, a, false)
}

func TestMaxLenForAccountNumber(t *testing.T) {
	FTest(true, FMaxLenForAccountNumber(), 2)
}

func TestSetTheAccounts(t *testing.T) {
	FSetTheAccounts()
}

func TestAddAutoCompletion(t *testing.T) {
	a1 := FAddAutoCompletion(SAutoCompletion{
		AccountName:          "book",
		PriceRevenue:         1250,
		PriceTax:             250,
		DiscountWay:          CDiscountPerOne,
		DiscountPerOne:       250,
		DiscountTotal:        0,
		DiscountPerQuantity:  SPQ{},
		DiscountDecisionTree: []SPQ{},
	})
	fmt.Println(a1)
	FPrintStructSlice(false, VAutoCompletionEntries)
}

func TestFAddAccount(t *testing.T) {
	a1 := FAddAccount(true, SAccount1{
		IsCredit:     false,
		CostFlowType: "",
		Inventory:    "home",
		Name:         "book2",
		Notes:        "",
		Image:        []string{},
		Barcode:      []string{"1", "2"},
		Number:       [][]uint{{1}, {1}},
		Levels:       []uint{},
		FathersName:  [][]string{},
	})
	log.Println(a1)
}

func TestPrint(t *testing.T) {
	p := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	for _, v1 := range VAccounts {
		fmt.Fprintf(p, "%#v\n", v1)
	}
	p.Flush()
}

func TestFCheckTheTree(t *testing.T) {
	a1 := FSetTheAccounts()
	FPrintStructSlice(true, a1)
}
