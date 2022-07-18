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
	e1 := SAccount3{IsCredit: "", CostFlowType: "", Name: "", Notes: "", Image: "", Number: []TErr{""}, Levels: "", FathersName: ""}

	a1 = FAddAccount(true, SAccount1{
		IsCredit:     false,
		CostFlowType: CHighLevelAccount,
		Name:         "assets",
		Number:       [][]uint{{1}},
	})
	FTest(true, a1, e1)

	a1 = FAddAccount(true, SAccount1{
		IsCredit:     false,
		CostFlowType: CHighLevelAccount,
		Name:         "current assets",
		Number:       [][]uint{{1, 1}},
	})
	FTest(true, a1, e1)

	a1 = FAddAccount(true, SAccount1{
		IsCredit:     false,
		CostFlowType: CWma,
		Name:         "cash",
		Number:       [][]uint{{1, 1, 1}},
	})
	FTest(true, a1, e1)

	a1 = FAddAccount(true, SAccount1{
		IsCredit:     false,
		CostFlowType: CHighLevelAccount,
		Name:         "invnetory",
		Number:       [][]uint{{1, 1, 3}},
	})
	FTest(true, a1, e1)

	a1 = FAddAccount(true, SAccount1{
		IsCredit:     false,
		CostFlowType: CFifo,
		Name:         "book",
		Number:       [][]uint{{1, 1, 3, 1}},
	})
	FTest(true, a1, e1)

	a1 = FAddAccount(true, SAccount1{
		IsCredit:     true,
		CostFlowType: CHighLevelAccount,
		Name:         "liabilities",
		Number:       [][]uint{{2}},
	})
	FTest(true, a1, e1)

	a1 = FAddAccount(true, SAccount1{
		IsCredit:     true,
		CostFlowType: CHighLevelAccount,
		Name:         "owner's equity",
		Number:       [][]uint{{3}},
	})
	FTest(true, a1, e1)

	a1 = FAddAccount(true, SAccount1{
		IsCredit:     true,
		CostFlowType: CHighLevelAccount,
		Name:         "retained earnings",
		Number:       [][]uint{{3, 1}},
	})
	FTest(true, a1, e1)

	a1 = FAddAccount(true, SAccount1{
		IsCredit:     true,
		CostFlowType: CHighLevelAccount,
		Name:         "income",
		Number:       [][]uint{{3, 2}},
	})
	FTest(true, a1, e1)

	a1 = FAddAccount(true, SAccount1{
		IsCredit:     true,
		CostFlowType: CHighLevelAccount,
		Name:         "revenue",
		Number:       [][]uint{{3, 2, 1}},
	})
	FTest(true, a1, e1)

	a1 = FAddAccount(true, SAccount1{
		IsCredit:     true,
		CostFlowType: CWma,
		Name:         "rent",
		Number:       [][]uint{{3, 2, 1, 1}},
	})
	FTest(true, a1, e1)

	a1 = FAddAccount(true, SAccount1{
		IsCredit:     false,
		CostFlowType: CHighLevelAccount,
		Name:         "expense",
		Number:       [][]uint{{3, 2, 2}},
	})
	FTest(true, a1, e1)

	a1 = FAddAccount(true, SAccount1{
		IsCredit:     false,
		CostFlowType: CHighLevelAccount,
		Name:         "discounts",
		Number:       [][]uint{{3, 2, 2, 1}},
	})
	FTest(true, a1, e1)

	a1 = FAddAccount(true, SAccount1{
		IsCredit:     false,
		CostFlowType: CWma,
		Name:         "Invoice discount",
		Number:       [][]uint{{3, 2, 2, 1, 1}},
	})
	FTest(true, a1, e1)

	a1 = FAddAccount(true, SAccount1{
		IsCredit:     false,
		CostFlowType: CWma,
		Name:         "Discount of book",
		Number:       [][]uint{{3, 2, 2, 1, 2}},
	})
	FTest(true, a1, e1)

	a1 = FAddAccount(true, SAccount1{
		IsCredit:     false,
		CostFlowType: CHighLevelAccount,
		Name:         "cost of goods sold",
		Number:       [][]uint{{3, 2, 2, 2}},
	})
	FTest(true, a1, e1)

	a1 = FAddAccount(true, SAccount1{
		IsCredit:     false,
		CostFlowType: CWma,
		Name:         "Cost of book",
		Number:       [][]uint{{3, 2, 2, 2, 1}},
	})
	FTest(true, a1, e1)
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
	FTest(true, FIsItHighThanByOrder([]uint{1, 1}, []uint{1, 1}), false)
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
	VDbAccounts.DropAll()
	VDbAutoCompletionEntries.DropAll()

	TestAddAccount(t)
	var a1 error

	a1 = FAddAutoCompletion(SAutoCompletion1{
		Group:                "1",
		Barcode:              []string{},
		Inventory:            "Inventory item 1",
		CostOfGoodsSold:      "CostOfGoodsSold item 1",
		TaxExpenses:          "TaxExpenses item 1",
		TaxLiability:         "TaxLiability item 1",
		Revenue:              "Revenue item 1",
		Discount:             "Discount item 1",
		PriceTax:             0,
		PriceRevenue:         1250,
		DiscountWay:          CDiscountPrice,
		DiscountPrice:        0,
		DiscountPercent:      0,
		DiscountTotal:        0,
		DiscountPerQuantity:  SPQ{},
		DiscountDecisionTree: []SPQ{},
	})
	FTest(true, a1, nil)

	a1 = FAddAutoCompletion(SAutoCompletion1{
		Group:                "1",
		Barcode:              []string{},
		Inventory:            "Inventory item 2",
		CostOfGoodsSold:      "CostOfGoodsSold item 2",
		TaxExpenses:          "TaxExpenses item 2",
		TaxLiability:         "TaxLiability item 2",
		Revenue:              "Revenue item 2",
		Discount:             "Discount item 2",
		PriceTax:             200,
		PriceRevenue:         1000,
		DiscountWay:          CDiscountPrice,
		DiscountPrice:        100,
		DiscountPercent:      0,
		DiscountTotal:        0,
		DiscountPerQuantity:  SPQ{},
		DiscountDecisionTree: []SPQ{},
	})
	FTest(true, a1, nil)

	a1 = FAddAutoCompletion(SAutoCompletion1{
		Group:                "1",
		Barcode:              []string{},
		Inventory:            "",
		CostOfGoodsSold:      "",
		TaxExpenses:          "",
		TaxLiability:         "TaxLiability item 3",
		Revenue:              "Revenue item 3",
		Discount:             "Discount item 3",
		PriceTax:             200,
		PriceRevenue:         1000,
		DiscountWay:          CDiscountPrice,
		DiscountPrice:        100,
		DiscountPercent:      0,
		DiscountTotal:        0,
		DiscountPerQuantity:  SPQ{},
		DiscountDecisionTree: []SPQ{},
	})
	FTest(true, a1, nil)
}

func TestFAddAccount(t *testing.T) {
	a1 := FAddAccount(true, SAccount1{
		IsCredit:     false,
		CostFlowType: "",
		Name:         "book2",
		Notes:        "",
		Image:        []string{},
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
