package anti_accountants

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
	CPathAdjustingEntry        = "/adjusting_entry"
	CPathEmployees             = "/employees"
	CPathJournalDrafts         = "/journal_drafts"
	CPathInvoiceDrafts         = "/invoice_drafts"

	CFifo             = "Fifo"
	CLifo             = "Lifo"
	CWma              = "Wma"
	CHighLevelAccount = "HighLevelAccount"

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
	CDontFilter           = "DontFilter"
	CInSlice              = "InSlice"
	CNotInSlice           = "NotInSlice"
	CElementsInElement    = "ElementsInElement"
	CElementsNotInElement = "ElementsNotInElement"

	// way to sort statment
	CAscending  = "Ascending"
	CDescending = "Descending"

	// Discount Way
	CDiscountPrice        = "Price"
	CDiscountPercent      = "Percent"
	CDiscountTotal        = "Total"
	CDiscountPerQuantity  = "PerQuantity"
	CDiscountDecisionTree = "DecisionTree"
)

var (
	VIndexOfAccountNumber    uint
	VCompanyName             string
	VEmployeeName            string
	VDbAccounts              *badger.DB
	VDbJournal               *badger.DB
	VDbInventory             *badger.DB
	VDbAutoCompletionEntries *badger.DB
	VDbAdjustingEntry        *badger.DB
	VDbEmployees             *badger.DB
	VDbJournalDrafts         *badger.DB
	VDbInvoiceDrafts         *badger.DB
	VAccounts                []SAccount1
	VAutoCompletionEntries   []SAutoCompletion1

	VPrintTable      = tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	VAllCostFlowType = []string{CHighLevelAccount, CFifo, CLifo, CWma}
	VLowCostFlowType = []string{CFifo, CLifo, CWma}
	VDiscountWay     = []string{CDiscountPrice, CDiscountPercent, CDiscountTotal, CDiscountPerQuantity, CDiscountDecisionTree}

	VErrorNotListed = errors.New("is not listed")
	VErrorIsUsed    = errors.New("is used")

	//this var for Test function
	VFailTestNumber uint
)

type INumber interface{ int | uint | float64 | float32 }
type TErr string
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

type SEntry1 SEntry[string, []string]
type SEntry[
	t1 string | bool | SFilter[string],
	t2 []string | bool | SFilter[string],
] struct {
	Notes    t1
	Name     t1
	Employee t1
	Labels   t2
}

type SPQ struct {
	TPrice
	TQuantity
}

type SAPQ1 SAPQ[string, float64]
type SAPQ2 SAPQ[TErr, TErr]
type SAPQ[t1 string | TErr, t2 float64 | TErr] struct {
	AccountName t1
	Price       t2
	Quantity    t2
}

type SAPQ12SAccount1 struct {
	SAPQ1
	SAPQ2
	SAccount1
}

type SAccount1 SAccount[bool, string, string, string, []string, [][]uint, []uint, [][]string]
type SAccount2 SAccount[SFilterBool, SFilter[string], SFilter[string], SFilter[string], SFilter[string], SFilter[uint], SFilter[uint], SFilter[string]]
type SAccount3 SAccount[TErr, TErr, TErr, TErr, TErr, []TErr, TErr, TErr]
type SAccount[
	IsCredit bool | SFilterBool | TErr,
	CostFlowType string | SFilter[string] | TErr,
	Name string | SFilter[string] | TErr,
	Notes string | SFilter[string] | TErr,
	Image []string | SFilter[string] | TErr,
	Number [][]uint | SFilter[uint] | []string | []TErr,
	Levels []uint | SFilter[uint] | TErr,
	FathersName [][]string | SFilter[string] | TErr,
] struct {
	IsCredit     IsCredit
	CostFlowType CostFlowType
	Name         Name
	Notes        Notes
	Image        Image
	Number       Number
	Levels       Levels
	FathersName  FathersName
}

type SJournal1 SJournal[time.Time, bool, uint, float64, string, []string]
type SJournal2 SJournal[SFilter[time.Time], SFilterBool, SFilter[uint], SFilter[float64], SFilter[string], SFilter[string]]
type SJournal3 SJournal[bool, bool, bool, bool, bool, bool]
type SJournal[
	t1 time.Time | bool | SFilter[time.Time],
	t2 bool | SFilterBool,
	t3 uint | bool | SFilter[uint],
	t4 float64 | bool | SFilter[float64],
	t5 string | bool | SFilter[string],
	t6 []string | bool | SFilter[string],
] struct {
	Date                       t1
	IsReverse                  t2
	IsReversed                 t2
	ReverseEntryNumberCompound t3
	ReverseEntryNumberSimple   t3
	EntryNumberCompound        t3
	EntryNumberSimple          t3
	Value                      t4
	DebitAccountName           t5
	DebitPrice                 t4
	DebitQuantity              t4
	DebitBalanceValue          t4
	DebitBalancePrice          t4
	DebitBalanceQuantity       t4
	CreditAccountName          t5
	CreditPrice                t4
	CreditQuantity             t4
	CreditBalanceValue         t4
	CreditBalancePrice         t4
	CreditBalanceQuantity      t4
	SEntry[t5, t6]
}

type SStatement1 SStatement[string, string, float64]
type SStatement2 SStatement[SAccount2, SFilter[string], SFilter[float64]]
type SStatement[
	t1 string | SAccount2,
	t2 string | SFilter[string],
	t3 float64 | SFilter[float64],
] struct {
	Account1Name           t1
	Account2Name           t1
	PersonName             t2
	Vpq                    t2
	TypeOfVpq              t2
	ChangeOrRatioOrBalance t2
	Number                 t3
}

type SStatmentWithAccount struct {
	Account1 SAccount1
	Account2 SAccount1
	SStatement1
}

type SFilterBool struct {
	IsFilter  bool
	BoolValue bool
}

type SFilter[t uint | float64 | string | time.Time] struct {
	Way   string
	Slice []t
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

type SInvoiceEntry struct {
	RevenueNameError error
	QuantityError    error
	Group            string
	Revenue          string
	PriceRevenue     float64
	PriceTax         float64
	Discount         float64
	DiscountWay      string
	Quantity         float64
}

type SAutoCompletion1 SAutoCompletion[string, []string, string, string, string, string, string, string, float64, float64, string, float64, float64, float64, SPQ, []SPQ]
type SAutoCompletion2 SAutoCompletion[TErr, TErr, TErr, TErr, TErr, TErr, TErr, TErr, TErr, TErr, TErr, TErr, TErr, TErr, TErr, TErr]
type SAutoCompletion[
	Group string | SFilter[string] | TErr,
	Barcode []string | SFilter[string] | TErr,
	Inventory string | SFilter[string] | TErr,
	CostOfGoodsSold string | SFilter[string] | TErr,
	TaxExpenses string | SFilter[string] | TErr,
	TaxLiability string | SFilter[string] | TErr,
	Revenue string | SFilter[string] | TErr,
	Discount string | SFilter[string] | TErr,
	PriceTax float64 | SFilter[float64] | TErr,
	PriceRevenue float64 | SFilter[float64] | TErr,
	DiscountWay string | SFilter[string] | TErr,
	DiscountPrice float64 | SFilter[float64] | TErr,
	DiscountPercent float64 | SFilter[float64] | TErr,
	DiscountTotal float64 | SFilter[float64] | TErr,
	DiscountPerQuantity SPQ | TErr,
	DiscountDecisionTree []SPQ | TErr,
] struct {
	Group                Group
	Barcode              Barcode
	Inventory            Inventory
	CostOfGoodsSold      CostOfGoodsSold
	TaxExpenses          TaxExpenses
	TaxLiability         TaxLiability
	Revenue              Revenue
	Discount             Discount
	PriceTax             PriceTax
	PriceRevenue         PriceRevenue
	DiscountWay          DiscountWay
	DiscountPrice        DiscountPrice
	DiscountPercent      DiscountPercent
	DiscountTotal        DiscountTotal
	DiscountPerQuantity  DiscountPerQuantity
	DiscountDecisionTree DiscountDecisionTree
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

type SAdjustingEntry struct {
	AccountName1    string
	AccountName2    string
	Price           float64
	Quantity        float64
	DateStart       time.Time
	DateEnd         time.Time
	AdjustingMethod string
	EntryInfo       SEntry1
}
