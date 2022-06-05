package main

import (
	"errors"
	"os"
	"text/tabwriter"
	"time"

	badger "github.com/dgraph-io/badger/v3"
)

const (
	CPathDataBase              = "./db/"
	CPathAccounts              = "/accounts"
	CPathJournal               = "/journal"
	CPathInventory             = "/inventory"
	CPathAutoCompletionEntries = "/auto_completion_entries"
	CPathEmployees             = "/employees"
	CPathDrafts                = "/drafts"

	CFifo             = "Fifo"
	CLifo             = "Lifo"
	CWma              = "Wma"
	CHighLevelAccount = "HighLevelAccount"

	CCredit = "Credit"
	CDebit  = "Debit"

	CLinear      = "Linear"
	CExponential = "Exponential"
	CLogarithmic = "Logarithmic"
	CSaturday    = "Saturday"
	CSunday      = "Sunday"
	CMonday      = "Monday"
	CTuesday     = "Tuesday"
	CWednesday   = "Wednesday"
	CThursday    = "Thursday"
	CFriday      = "Friday"
	CTimeLayout  = "2006-01-02 15:04:05.999999999 -0700 MST"
	//constants for financial_statements
	// vpq
	CValue    = "Value"
	CPrice    = "Price"
	CQuantity = "Quantity"
	// type_of_vpq
	CFlowInBeginning         = "FlowInBeginning"
	CFlowOutBeginning        = "FlowOutBeginning"
	CFlowBeginning           = "FlowBeginning"
	CFlowInPeriod            = "FlowInPeriod"
	CFlowOutPeriod           = "FlowOutPeriod"
	CFlowPeriod              = "FlowPeriod"
	CFlowInEnding            = "FlowInEnding"
	CFlowOutEnding           = "FlowOutEnding"
	CFlowEnding              = "FlowEnding"
	CAverage                 = "Average"
	CTurnover                = "Turnover"
	CTurnoverDays            = "TurnoverDays"
	CGrowthRatio             = "GrowthRatio"
	CNamePercent             = "NamePercent"
	CBalance                 = "Balance"
	CChangeSinceBasePeriod   = "ChangeSinceBasePeriod"
	CGrowthRatioToBasePeriod = "GrowthRatioToBasePeriod"
	// key words for statment columns in financial statement
	CNames       = "Names"
	CAllNames    = "AllNames"
	CAllAccounts = "AllAccounts"

	// all cvp keyword
	CVariableCost                  = "VariableCost"
	CVariableCostPerUnits          = "VariableCostPerUnits"
	CUnits                         = "Units"
	CFixedCost                     = "FixedCost"
	CFixedCostPerUnits             = "FixedCostPerUnits"
	CMixedCost                     = "MixedCost"
	CMixedCostPerUnits             = "MixedCostPerUnits"
	CSales                         = "Sales"
	CSalesPerUnits                 = "SalesPerUnits"
	CProfit                        = "Profit"
	CProfitPerUnits                = "ProfitPerUnits"
	CContributionMargin            = "ContributionMargin"
	CContributionMarginPerUnits    = "ContributionMarginPerUnits"
	CBreakEvenInSales              = "BreakEvenInSales"
	CBreakEvenInUnits              = "BreakEvenInUnits"
	CContributionMarginRatio       = "ContributionMarginRatio"
	CDegreeOfOperatingLeverage     = "DegreeOfOperatingLeverage"
	CUnitsGap                      = "UnitsGap"
	CActualUnits                   = "ActualUnits"
	CTotal                         = "Total"
	CPercentFromVariableCost       = "PercentFromVariableCost"
	CPercentFromFixedCost          = "PercentFromFixedCost"
	CPercentFromMixedCost          = "PercentFromMixedCost"
	CPercentFromSales              = "PercentFromSales"
	CPercentFromProfit             = "PercentFromProfit"
	CPercentFromContributionMargin = "PercentFromContributionMargin"
	CPortions                      = "Portions"

	// filter key words for numbers and dates
	CBetween    = "Between"    // between big and small
	CNotBetween = "NotBetween" // not between big and small
	CBigger     = "Bigger"     // bigger than big
	CSmaller    = "Smaller"    // smaller than small

	// filter key words for string
	CInSlice              = "InSlice"
	CNotInSlice           = "NotInSlice"
	CElementsInElement    = "ElementsInElement"
	CElementsNotInElement = "ElementsNotInElement"

	// way to sort statment
	CAscending  = "Ascending"
	CDescending = "Descending"

	// Prefixes of inventory account
	CPrefixCost         = "cost of "
	CPrefixDiscount     = "discount of "
	CPrefixTaxExpenses  = "tax expenses of "
	CPrefixTaxLiability = "tax liability of "
	CPrefixRevenue      = "revenue of "
)

var (
	VIndexOfAccountNumber = 0
	VInvoiceDiscount      = "Invoice PQ"

	VInvoiceDiscountsList    []SPQ
	VCompanyName             string
	VEmployeeName            string
	VDbAccounts              *badger.DB
	VDbJournal               *badger.DB
	VDbInventory             *badger.DB
	VDbAutoCompletionEntries *badger.DB
	VDbEmployees             *badger.DB
	VDbJournalDrafts         *badger.DB
	VAccounts                []SAccount
	VAutoCompletionEntries   []SAutoCompletion

	// standards
	VPrintTable   = tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	VCostFlowType = []string{CFifo, CLifo, CWma}

	//errors
	VErrorNotListed          = errors.New("is not listed")
	VErrorAccountNameIsUsed  = errors.New("account name is used")
	VErrorBarcodeIsUsed      = errors.New("barcode is used")
	VErrorAccountNameIsEmpty = errors.New("account name is empty")

	//this var for Test function
	VFailTestNumber int
)

type INumber interface{ IInteger | float64 | float32 }
type IInteger interface{ int | int64 | uint }

type (
	TAccountName            = string
	TAccount1Name           = string
	TAccount2Name           = string
	TPrice                  = float64
	TQuantity               = float64
	TPersonName             = string
	TVpq                    = string
	TTypeOfVpq              = string
	TChangeOrRatioOrBalance = string
	TIsBeforeDateStart      = bool
	TIsCredit               = bool
	TNumber                 = float64
	TStatement1             = map[TAccount1Name]map[TAccount2Name]map[TPersonName]map[TVpq]map[TIsBeforeDateStart]map[TIsCredit]TNumber
	TStatement2             = map[TAccount1Name]map[TAccount2Name]map[TPersonName]map[TVpq]map[TTypeOfVpq]TNumber
	TStatement3             = map[TAccount1Name]map[TAccount2Name]map[TPersonName]map[TVpq]map[TTypeOfVpq]map[TChangeOrRatioOrBalance]TNumber
)

type SEntry struct {
	Notes               string
	Name                string
	Employee            string
	TypeOfCompoundEntry string
}
type SPQ struct {
	TPrice
	TQuantity
}
type SAPQ struct {
	TAccountName
	TPrice
	TQuantity
}
type SAPQAE struct {
	TAccountName
	TPrice
	TQuantity
	SAccount
	error
}
type SAccount struct {
	TIsCredit           bool
	TCostFlowType       string
	TAccountName        string
	TAccountNotes       string
	TAccountImage       []string
	TAccountBarcode     []string
	TAccountNumber      [][]uint
	TAccountLevels      []uint
	TAccountFathersName [][]string
}
type SJournal struct {
	IsReverse                  bool
	IsReversed                 bool
	ReverseEntryNumberCompound int
	ReverseEntryNumberSimple   int
	EntryNumberCompound        int
	EntryNumberSimple          int
	Value                      float64
	DebitAccountName           string
	DebitPrice                 float64
	DebitQuantity              float64
	DebitBalanceValue          float64
	DebitBalancePrice          float64
	DebitBalanceQuantity       float64
	CreditAccountName          string
	CreditPrice                float64
	CreditQuantity             float64
	CreditBalanceValue         float64
	CreditBalancePrice         float64
	CreditBalanceQuantity      float64
	SEntry
}
type SStatement struct {
	TAccount1Name
	TAccount2Name
	TPersonName
	TVpq
	TTypeOfVpq
	TChangeOrRatioOrBalance
	TNumber
}
type SStatmentWithAccount struct {
	Account1 SAccount
	Account2 SAccount
	SStatement
}
type SFilterStatement struct {
	Account1               SFilterAccount
	Account2               SFilterAccount
	Name                   SFilterString
	Vpq                    SFilterString
	TypeOfVpq              SFilterString
	ChangeOrRatioOrBalance SFilterString
	Number                 SFilterNumber
}
type SFilterJournal struct {
	Date                       SFilterDate
	IsReverse                  SFilterBool
	IsReversed                 SFilterBool
	ReverseEntryNumberCompound SFilterNumber
	ReverseEntryNumberSimple   SFilterNumber
	EntryNumberCompound        SFilterNumber
	EntryNumberSimple          SFilterNumber
	Value                      SFilterNumber
	DebitAccountName           SFilterString
	DebitPrice                 SFilterNumber
	DebitQuantity              SFilterNumber
	DebitBalanceValue          SFilterNumber
	DebitBalancePrice          SFilterNumber
	DebitBalanceQuantity       SFilterNumber
	CreditAccountName          SFilterString
	CreditPrice                SFilterNumber
	CreditQuantity             SFilterNumber
	CreditBalanceValue         SFilterNumber
	CreditBalancePrice         SFilterNumber
	CreditBalanceQuantity      SFilterNumber
	Notes                      SFilterString
	Name                       SFilterString
	Employee                   SFilterString
	TypeOfCompoundEntry        SFilterString
}
type SFilterJournalDuplicate struct {
	IsReverse                  bool
	IsReversed                 bool
	ReverseEntryNumberCompound bool
	ReverseEntryNumberSimple   bool
	Value                      bool
	PriceDebit                 bool
	PriceCredit                bool
	QuantityDebit              bool
	QuantityCredit             bool
	AccountDebit               bool
	AccountCredit              bool
	Notes                      bool
	Name                       bool
	Employee                   bool
	TypeOfCompoundEntry        bool
}
type SFilterAccount struct {
	IsFilter    bool
	IsCredit    SFilterBool
	Account     SFilterString
	FathersName SFilterFathersAccountsName
	Levels      SFilterSliceUint
}
type SFilterFathersAccountsName struct {
	IsFilter      bool
	InAccountName bool
	InFathersName bool
	FathersName   []string
}
type SFilterSliceUint struct {
	IsFilter bool
	InSlice  bool
	Slice    []uint
}
type SFilterDate struct {
	IsFilter bool
	Way      string
	Slice    []time.Time
}
type SFilterNumber struct {
	IsFilter bool
	Way      string
	Slice    []float64
}
type SFilterString struct {
	IsFilter bool
	Way      string
	Slice    []string
}
type SFilterBool struct {
	IsFilter  bool
	BoolValue bool
}
type SFinancialAnalysis struct {
	CurrentAssets                          float64
	CurrentLiabilities                     float64
	Cash                                   float64
	ShortTermInvestments                   float64
	NetReceivables                         float64
	NetCreditSales                         float64
	AverageNetReceivables                  float64
	CostOfGoodsSold                        float64
	AverageInventory                       float64
	NetIncome                              float64
	NetSales                               float64
	AverageAssets                          float64
	AverageEquity                          float64
	PreferredDividends                     float64
	AverageCommonStockholdersEquity        float64
	MarketPricePerSharesOutstanding        float64
	CashDividends                          float64
	TotalDebt                              float64
	TotalAssets                            float64
	Ebitda                                 float64
	InterestExpense                        float64
	WeightedAverageCommonSharesOutstanding float64
}
type SFinancialAnalysisStatement struct {
	CurrentRatio                     float64 // CURRENT_ASSETS / CURRENT_LIABILITIES
	AcidTest                         float64 // (CASH + SHORT_TERM_INVESTMENTS + NET_RECEIVABLES) / CURRENT_LIABILITIES
	ReceivablesTurnover              float64 // NET_CREDIT_SALES / AVERAGE_NET_RECEIVABLES
	InventoryTurnover                float64 // COST_OF_GOODS_SOLD / AVERAGE_INVENTORY
	ProfitMargin                     float64 // NET_INCOME / NET_SALES
	AssetTurnover                    float64 // NET_SALES / AVERAGE_ASSETS
	ReturnOnAssets                   float64 // NET_INCOME / AVERAGE_ASSETS
	ReturnOnEquity                   float64 // NET_INCOME / AVERAGE_EQUITY
	PayoutRatio                      float64 // CASH_DIVIDENDS / NET_INCOME
	DebtToTotalAssetsRatio           float64 // TOTAL_DEBT / TOTAL_ASSETS
	TimesInterestEarned              float64 // EBITDA / INTEREST_EXPENSE
	ReturnOnCommonStockholdersEquity float64 // (NET_INCOME - PREFERRED_DIVIDENDS) / AVERAGE_COMMON_STOCKHOLDERS_EQUITY
	EarningsPerShare                 float64 // (NET_INCOME - PREFERRED_DIVIDENDS) / WEIGHTED_AVERAGE_COMMON_SHARES_OUTSTANDING
	PriceEarningsRatio               float64 // MARKET_PRICE_PER_SHARES_OUTSTANDING / EARNINGS_PER_SHARE
}
type SOneStepDistribution struct {
	SalesOrVariableOrFixed string
	DistributionMethod     string
	Amount                 float64
	From                   map[string]float64
	To                     map[string]float64
}
type SAutoCompletion struct {
	TAccountName  string
	PriceRevenue  float64
	PriceTax      float64
	PriceDiscount []SPQ
}
type SCvp struct {
	VariableCost       float64
	FixedCost          float64
	MixedCost          float64
	Sales              float64
	Profit             float64
	ContributionMargin float64
}
type SAVQ struct {
	TAccountName
	TValue float64
	TQuantity
}
