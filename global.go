package main

import (
	"errors"
	"os"
	"text/tabwriter"
	"time"
)

const (
	Fifo        = "Fifo"
	Lifo        = "Lifo"
	Wma         = "Wma"
	Linear      = "Linear"
	Exponential = "Exponential"
	Logarithmic = "Logarithmic"
	Saturday    = "Saturday"
	Sunday      = "Sunday"
	Monday      = "Monday"
	Tuesday     = "Tuesday"
	Wednesday   = "Wednesday"
	Thursday    = "Thursday"
	Friday      = "Friday"
	TimeLayout  = "2006-01-02 15:04:05.999999999 -0700 MST"
	//constants for financial_statements
	// vpq
	Value    = "Value"
	Price    = "Price"
	Quantity = "Quantity"
	// type_of_vpq
	BeginningBalance        = "BeginningBalance"
	EndingBalance           = "EndingBalance"
	Inflow                  = "Inflow"
	Outflow                 = "Outflow"
	Flow                    = "Flow"
	Average                 = "Average"
	Turnover                = "Turnover"
	TurnoverDays            = "TurnoverDays"
	GrowthRatio             = "GrowthRatio"
	NamePercent             = "NamePercent"
	ChangeSinceBasePeriod   = "ChangeSinceBasePeriod"
	GrowthRatioToBasePeriod = "GrowthRatioToBasePeriod"
	// key words for statment columns in financial statement
	Names       = "Names"
	AllNames    = "AllNames"
	AllAccounts = "AllAccounts"

	// all cvp keyword
	VariableCost                  = "VariableCost"
	VariableCostPerUnits          = "VariableCostPerUnits"
	Units                         = "Units"
	FixedCost                     = "FixedCost"
	FixedCostPerUnits             = "FixedCostPerUnits"
	MixedCost                     = "MixedCost"
	MixedCostPerUnits             = "MixedCostPerUnits"
	Sales                         = "Sales"
	SalesPerUnits                 = "SalesPerUnits"
	Profit                        = "Profit"
	ProfitPerUnits                = "ProfitPerUnits"
	ContributionMargin            = "ContributionMargin"
	ContributionMarginPerUnits    = "ContributionMarginPerUnits"
	BreakEvenInSales              = "BreakEvenInSales"
	BreakEvenInUnits              = "BreakEvenInUnits"
	ContributionMarginRatio       = "ContributionMarginRatio"
	DegreeOfOperatingLeverage     = "DegreeOfOperatingLeverage"
	UnitsGap                      = "UnitsGap"
	ActualUnits                   = "ActualUnits"
	Total                         = "Total"
	PercentFromVariableCost       = "PercentFromVariableCost"
	PercentFromFixedCost          = "PercentFromFixedCost"
	PercentFromMixedCost          = "PercentFromMixedCost"
	PercentFromSales              = "PercentFromSales"
	PercentFromProfit             = "PercentFromProfit"
	PercentFromContributionMargin = "PercentFromContributionMargin"
	Portions                      = "Portions"

	// filter key words for numbers and dates
	Between          = "Between"    // between big and small
	NotBetween       = "NotBetween" // not between big and small
	Bigger           = "Bigger"     // bigger than big
	Smaller          = "Smaller"    // smaller than small
	EqualToOneOfThem = "EqualToOneOfThem"

	// filter key words for string
	InSlice              = "InSlice"
	NotInSlice           = "NotInSlice"
	ElementsInElement    = "ElementsInElement"
	ElementsNotInElement = "ElementsNotInElement"
)

var (
	CompanyName          = "anti_accountants"
	EmployeeName         = "hashem"
	IndexOfAccountNumber = 0
	RetinedEarnings      = FormatTheString("Retained Earnings")
	InvoiceDiscountsList [][2]float64
	AutoCompleteEntries  []AutoCompleteEntrie
	ErrorsMessages       = CheckTheTree()
	// all the below is final
	// pathes
	DbPathAccounts  = "./db/" + CompanyName + "/accounts"
	DbPathJournal   = "./db/" + CompanyName + "/journal"
	DbPathInventory = "./db/" + CompanyName + "/inventory"
	// data base
	DbAccounts  = DbOpen(DbPathAccounts)
	DbJournal   = DbOpen(DbPathJournal)
	DbInventory = DbOpen(DbPathInventory)
	_, Accounts = DbRead[Account](DbAccounts)
	// standards
	PrintTable = tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	// STANDARD_DAYS        = []string{SATURDAY, SUNDAY, MONDAY, TUESDAY, WEDNESDAY, THURSDAY, FRIDAY}
	// DEPRECIATION_METHODS = []string{LINEAR, EXPONENTIAL, LOGARITHMIC}
	CostFlowType = []string{Fifo, Lifo, Wma}
	//errors
	ErrorNotListed           = errors.New("is not listed")
	ErrorNotInventoryAccount = errors.New("not inventory account")
	ErrorShouldBeNegative    = errors.New("the QUANTITY should be negative")
	ErrorAccountNameIsUsed   = errors.New("account name is used")
	ErrorBarcodeIsUsed       = errors.New("barcode is used")
	ErrorAccountNumberIsUsed = errors.New("account number is used")
	ErrorAccountNameIsEmpty  = errors.New("account name is empty")

	//this vaiable for TEST function
	FailTestNumber int
)

type Number interface{ Integer | float64 | float32 }
type Integer interface{ int | int64 | uint }
type Account struct { // 									   configer		|change				|correct	|necessary	|is unique
	IsLowLevelAccount                       bool       // manual		|if not in journal	|cant		|yes		|no
	IsCredit                                bool       // manual		|if not in journal	|cant		|yes		|no
	IsTemporary                             bool       // manual		|if not in journal	|auto		|yes		|no
	CostFlowType                            string     // manual		|manual				|auto		|yes		|no
	AccountName                             string     // manual		|if not used		|manual		|yes		|yes
	Notes                                   string     // manual		|manual				|manual		|no			|no
	Image                                   []string   // manual		|manual				|manual		|no			|no
	Barcode                                 []string   // manual		|if not used		|manual		|yes		|yes
	AccountNumber                           [][]uint   // manual		|manual				|manual		|yes		|yes it should be but we don't inforce you
	AccountLevels                           []uint     // auto		|auto				|auto		|yes		|no
	FathersAccountsName                     [][]string // auto		|auto				|auto		|yes		|no
	AlertForMinimumQuantityByTurnoverInDays uint       // manual		|manual				|manual		|no			|no
	AlertForMinimumQuantityByQuintity       float64    // manual		|manual				|manual		|no			|no
	TargetBalance                           float64    // manual		|manual				|manual		|no			|no
	IfTheTargetBalanceIsLessIsGood          bool       // manual		|manual				|manual		|no			|no
}
type DayStartEnd struct {
	Day         string
	StartHour   int
	StartMinute int
	EndHour     int
	EndMinute   int
}
type StartEndMinutes struct {
	DateStart time.Time
	DateEnd   time.Time
	Minutes   float64
}
type AutoCompleteEntrie struct {
	AccountName    string
	AccountCredit  string
	AccountDebit   string
	PriceDebit     float64
	PriceCredit    float64
	Quantity       float64
	QuantityDebit  float64
	QuantityCredit float64
}
type FilteredStatement struct {
	Account1  string
	Account2  string
	Name      string
	Vpq       string
	TypeOfVpq string
	Number    float64
}
type TheJournalDuplicateFilter struct {
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
	NameEmployee               bool
}
type FilterStatement struct {
	Account1  FilterString
	Account2  FilterString
	Name      FilterString
	Vpq       FilterString
	TypeOfVpq FilterString
	Number    FilterNumber
}
type FilterJournal struct {
	Date                       FilterDate
	IsReverse                  FilterBool
	IsReversed                 FilterBool
	ReverseEntryNumberCompound FilterNumber
	ReverseEntryNumberSimple   FilterNumber
	EntryNumberCompound        FilterNumber
	EntryNumberSimple          FilterNumber
	Value                      FilterNumber
	PriceDebit                 FilterNumber
	PriceCredit                FilterNumber
	QuantityDebit              FilterNumber
	QuantityCredit             FilterNumber
	AccountDebit               FilterString
	AccountCredit              FilterString
	Notes                      FilterString
	Name                       FilterString
	NameEmployee               FilterString
}
type FilterAccount struct {
	IsLowLevelAccount   FilterBool
	IsCredit            FilterBool
	IsTemporary         FilterBool
	FathersAccountsName FilterSliceString
	AccountLevels       FilterSliceUint
}
type FilterSliceString struct {
	IsFilter bool
	InSlice  bool
	Slice    []string
}
type FilterSliceUint struct {
	IsFilter bool
	InSlice  bool
	Slice    []uint
}
type FilterDate struct {
	IsFilter bool
	Way      string // here you have some method : between, not_between, bigger, smaller, equal_to_one_of_them
	Big      time.Time
	Small    time.Time
}
type FilterNumber struct {
	IsFilter bool
	Way      string // here you have some method : between, not_between, bigger, smaller, equal_to_one_of_them
	Big      float64
	Small    float64
}
type FilterString struct {
	IsFilter bool
	Way      string // here you have some method : in_slice, not_in_slice, elements_in_element , elements_not_in_element
	Slice    []string
}
type FilterBool struct {
	IsFilter  bool
	BoolValue bool
}
type FinancialAccounting struct {
	Assets                 []string
	CurrentAssets          []string
	CashAndCashEquivalents []string
	ShortTermInvestments   []string
	Receivables            []string
	Inventory              []string
	Liabilities            []string
	CurrentLiabilities     []string
	Equity                 []string
	RetainedEarnings       []string
	Dividends              []string
	IncomeStatement        []string
	Ebitda                 []string
	Sales                  []string
	CostOfGoodsSold        []string
	Discounts              []string
	InvoiceDiscount        []string
	InterestExpense        []string
}
type JournalTag struct {
	IsReverse                  bool    // this is true to the new entry when you enter reverse old entry
	IsReversed                 bool    // this is true to the old entry when you enter reverse old entry
	ReverseEntryNumberCompound int     // that mean if this is reverse entry what the entry compound was reversed
	ReverseEntryNumberSimple   int     // that mean if this is reverse entry what the entry simple was reversed
	EntryNumberCompound        int     // that mean number entry you made
	EntryNumberSimple          int     // that mean the index of the simple entry in the you made
	Value                      float64 // this sould be positive
	PriceDebit                 float64 // this sould be positive
	PriceCredit                float64 // this sould be positive
	QuantityDebit              float64 // this sould be positive
	QuantityCredit             float64 // this sould be positive
	AccountDebit               string  // the account name in the debit side
	AccountCredit              string  // the account name in the credit side
	Notes                      string  // your nots on the entry
	Name                       string  // the name of the dealer or customer
	NameEmployee               string  // the name of the employee that made the entry
}
type InventoryTag struct {
	Price       float64
	Quantity    float64
	AccountName string
}
type InvoiceStruct struct {
	Value       float64
	Price       float64
	Quantity    float64
	AccountName string
}
type FinancialAnalysis struct {
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
type FinancialAnalysisStatement struct {
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
type PriceQuantityAccountBarcode struct {
	Price       float64
	Quantity    float64
	AccountName string
	Barcode     string
}
type PriceQuantityAccount struct {
	IsCredit     bool
	CostFlowType string
	AccountName  string
	Price        float64
	Quantity     float64
}
type OneStepDistribution struct {
	SalesOrVariableOrFixed string
	DistributionMethod     string
	Amount                 float64
	From                   map[string]float64
	To                     map[string]float64
}
