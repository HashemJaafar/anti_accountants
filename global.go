package main

import (
	"errors"
	"os"
	"text/tabwriter"
	"time"
)

const (
	FIFO        = "fifo"
	LIFO        = "lifo"
	WMA         = "wma"
	LINEAR      = "linear"
	EXPONENTIAL = "exponential"
	LOGARITHMIC = "logarithmic"
	SATURDAY    = "Saturday"
	SUNDAY      = "Sunday"
	MONDAY      = "Monday"
	TUESDAY     = "Tuesday"
	WEDNESDAY   = "Wednesday"
	THURSDAY    = "Thursday"
	FRIDAY      = "Friday"
	TIME_LAYOUT = "2006-01-02 15:04:05.999999999 -0700 MST"
	//constants for financial_statements
	// vpq
	VALUE    = "value"
	PRICE    = "price"
	QUANTITY = "quantity"
	// type_of_vpq
	beginning_balance           = "beginning_balance"
	ending_balance              = "ending_balance"
	inflow                      = "inflow"
	outflow                     = "outflow"
	flow                        = "flow"
	average                     = "average"
	turnover                    = "turnover"
	turnover_days               = "turnover_days"
	growth_ratio                = "growth_ratio"
	name_percent                = "name_percent"
	change_since_base_period    = "change_since_base_period"
	growth_ratio_to_base_period = "growth_ratio_to_base_period"
	// key words for statment columns in financial statement
	all_names    = "all_names"
	names        = "names"
	all_accounts = "all_accounts"

	// all cvp keyword
	variable_cost                 = "variable_cost"
	variable_cost_per_units       = "variable_cost_per_units"
	units                         = "units"
	fixed_cost                    = "fixed_cost"
	fixed_cost_per_units          = "fixed_cost_per_units"
	mixed_cost                    = "mixed_cost"
	mixed_cost_per_units          = "mixed_cost_per_units"
	sales                         = "sales"
	sales_per_units               = "sales_per_units"
	profit                        = "profit"
	profit_per_units              = "profit_per_units"
	contribution_margin           = "contribution_margin"
	contribution_margin_per_units = "contribution_margin_per_units"
	break_even_in_sales           = "break_even_in_sales"
	break_even_in_units           = "break_even_in_units"
	contribution_margin_ratio     = "contribution_margin_ratio"
	degree_of_operating_leverage  = "degree_of_operating_leverage"
	units_gap                     = "units_gap"
	actual_units                  = "actual_units"

	// filter key words for numbers and dates
	between              = "between"
	not_between          = "not_between"
	bigger               = "bigger"
	smaller              = "smaller"
	equal_to_one_of_them = "equal_to_one_of_them"

	// filter key words for string
	in_slice                 = "in_slice"
	not_in_slice             = "not_in_slice"
	contain_one_in_slice     = "contain_one_in_slice"
	not_contain_one_in_slice = "not_contain_one_in_slice"
)

var (
	COMPANY_NAME            = "anti_accountants"
	EMPLOYEE_NAME           = "hashem"
	INDEX_OF_ACCOUNT_NUMBER = 0
	INVOICE_DISCOUNTS_LIST  [][2]float64
	AUTO_COMPLETE_ENTRIES   []AUTO_COMPLETE_ENTRIE
	ERRORS_MESSAGES         = CHECK_THE_TREE()
	PRIMARY_ACCOUNTS_NAMES  = FINANCIAL_ACCOUNTING{
		ASSETS:                    []string{"ASSETS"},
		CURRENT_ASSETS:            []string{"CURRENT_ASSETS"},
		CASH_AND_CASH_EQUIVALENTS: []string{"CASH_AND_CASH_EQUIVALENTS"},
		SHORT_TERM_INVESTMENTS:    []string{"SHORT_TERM_INVESTMENTS"},
		RECEIVABLES:               []string{"RECEIVABLES"},
		INVENTORY:                 []string{"INVENTORY"},
		LIABILITIES:               []string{"LIABILITIES"},
		CURRENT_LIABILITIES:       []string{"CURRENT_LIABILITIES"},
		EQUITY:                    []string{"EQUITY"},
		RETAINED_EARNINGS:         []string{"RETAINED_EARNINGS"},
		DIVIDENDS:                 []string{"DIVIDENDS"},
		INCOME_STATEMENT:          []string{"INCOME_STATEMENT"},
		EBITDA:                    []string{"EBITDA"},
		SALES:                     []string{"SALES"},
		COST_OF_GOODS_SOLD:        []string{"COST_OF_GOODS_SOLD"},
		DISCOUNTS:                 []string{"DISCOUNTS"},
		INVOICE_DISCOUNT:          []string{"INVOICE_DISCOUNT"},
		INTEREST_EXPENSE:          []string{"INTEREST_EXPENSE"},
	}
	// data base
	DB_ACCOUNTS  = DB_OPEN(DB_PATH_ACCOUNTS)
	DB_JOURNAL   = DB_OPEN(DB_PATH_JOURNAL)
	DB_INVENTORY = DB_OPEN(DB_PATH_INVENTORY)
	_, ACCOUNTS  = DB_READ[ACCOUNT](DB_ACCOUNTS)
	// all the below is final
	// pathes
	DB_PATH_ACCOUNTS  = "./db/" + COMPANY_NAME + "/accounts"
	DB_PATH_JOURNAL   = "./db/" + COMPANY_NAME + "/journal"
	DB_PATH_INVENTORY = "./db/" + COMPANY_NAME + "/inventory"
	// standards
	PRINT_TABLE = tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	// STANDARD_DAYS        = []string{SATURDAY, SUNDAY, MONDAY, TUESDAY, WEDNESDAY, THURSDAY, FRIDAY}
	// DEPRECIATION_METHODS = []string{LINEAR, EXPONENTIAL, LOGARITHMIC}
	COST_FLOW_TYPE = []string{FIFO, LIFO, WMA}
	//errors
	ERROR_NOT_LISTED             = errors.New("is not listed")
	ERROR_NOT_INVENTORY_ACCOUNT  = errors.New("not inventory account")
	ERROR_SHOULD_BE_NEGATIVE     = errors.New("the quantity should be negative")
	ERROR_ACCOUNT_NAME_IS_USED   = errors.New("account name is used")
	ERROR_BARCODE_IS_USED        = errors.New("barcode is used")
	ERROR_ACCOUNT_NUMBER_IS_USED = errors.New("account number is used")
	ERROR_ACCOUNT_NAME_IS_EMPTY  = errors.New("account name is empty")

	//this vaiable for TEST function
	fail_test_number int
)

type NUMBER interface{ INTEGER | float64 | float32 }
type INTEGER interface{ int | int64 | uint }
type ACCOUNT struct { // 									   configer		|change				|correct	|necessary	|is unique
	IS_LOW_LEVEL_ACCOUNT                           bool       // manual		|if not in journal	|cant		|yes		|no
	IS_CREDIT                                      bool       // manual		|if not in journal	|cant		|yes		|no
	IS_TEMPORARY                                   bool       // manual		|if not in journal	|auto		|yes		|no
	COST_FLOW_TYPE                                 string     // manual		|manual				|auto		|yes		|no
	ACCOUNT_NAME                                   string     // manual		|if not used		|manual		|yes		|yes
	NOTES                                          string     // manual		|manual				|manual		|no			|no
	IMAGE                                          []string   // manual		|manual				|manual		|no			|no
	BARCODE                                        []string   // manual		|if not used		|manual		|yes		|yes
	ACCOUNT_NUMBER                                 [][]uint   // manual		|manual				|manual		|yes		|yes it should be but we don't inforce you
	ACCOUNT_LEVELS                                 []uint     // auto		|auto				|auto		|yes		|no
	FATHER_AND_GRANDPA_ACCOUNTS_NAME               [][]string // auto		|auto				|auto		|yes		|no
	ALERT_FOR_MINIMUM_QUANTITY_BY_TURNOVER_IN_DAYS uint       // manual		|manual				|manual		|no			|no
	ALERT_FOR_MINIMUM_QUANTITY_BY_QUINTITY         float64    // manual		|manual				|manual		|no			|no
	TARGET_BALANCE                                 float64    // manual		|manual				|manual		|no			|no
	IF_THE_TARGET_BALANCE_IS_LESS_IS_GOOD          bool       // manual		|manual				|manual		|no			|no
}
type DAY_START_END struct {
	DAY          string
	START_HOUR   int
	START_MINUTE int
	END_HOUR     int
	END_MINUTE   int
}
type start_end_minutes struct {
	date_start time.Time
	date_end   time.Time
	minutes    float64
}
type AUTO_COMPLETE_ENTRIE struct {
	ACCOUNT_NAME    string
	ACCOUNT_CREDIT  string
	ACCOUNT_DEBIT   string
	PRICE_DEBIT     float64
	PRICE_CREDIT    float64
	QUANTITY        float64
	QUANTITY_DEBIT  float64
	QUANTITY_CREDIT float64
}
type FILTERED_STATEMENT struct {
	ACCOUNT1    string
	ACCOUNT2    string
	NAME        string
	VPQ         string
	TYPE_OF_VPQ string
	NUMBER      float64
}
type THE_JOURNAL_DUPLICATE_FILTER struct {
	IS_REVERSE                    bool
	IS_REVERSED                   bool
	REVERSE_ENTRY_NUMBER_COMPOUND bool
	REVERSE_ENTRY_NUMBER_SIMPLE   bool
	VALUE                         bool
	PRICE_DEBIT                   bool
	PRICE_CREDIT                  bool
	QUANTITY_DEBIT                bool
	QUANTITY_CREDIT               bool
	ACCOUNT_DEBIT                 bool
	ACCOUNT_CREDIT                bool
	NOTES                         bool
	NAME                          bool
	NAME_EMPLOYEE                 bool
}
type THE_FILTER_OF_THE_STATEMENT struct {
	old_statement         map[string]map[string]map[string]map[string]map[string]float64
	account1              []string
	account2              []string
	name                  []string
	vpq                   []string
	type_of_vpq           []string
	account1_levels       []uint
	account2_levels       []uint
	is_in_account1        bool
	is_in_account2        bool
	is_in_name            bool
	is_in_vpq             bool
	is_in_type_of_vpq     bool
	is_in_account1_levels bool
	is_in_account2_levels bool
}
type THE_JOURNAL_FILTER struct {
	DATE                          FILTER_DATE
	IS_REVERSE                    FILTER_BOOL
	IS_REVERSED                   FILTER_BOOL
	REVERSE_ENTRY_NUMBER_COMPOUND FILTER_NUMBER
	REVERSE_ENTRY_NUMBER_SIMPLE   FILTER_NUMBER
	ENTRY_NUMBER_COMPOUND         FILTER_NUMBER
	ENTRY_NUMBER_SIMPLE           FILTER_NUMBER
	VALUE                         FILTER_NUMBER
	PRICE_DEBIT                   FILTER_NUMBER
	PRICE_CREDIT                  FILTER_NUMBER
	QUANTITY_DEBIT                FILTER_NUMBER
	QUANTITY_CREDIT               FILTER_NUMBER
	ACCOUNT_DEBIT                 FILTER_STRING
	ACCOUNT_CREDIT                FILTER_STRING
	NOTES                         FILTER_STRING
	NAME                          FILTER_STRING
	NAME_EMPLOYEE                 FILTER_STRING
}
type FILTER_DATE struct {
	FILTER bool
	WAY    string // here you have some method : between, not_between, bigger, smaller, equal_to_one_of_them
	BIG    time.Time
	SMALL  time.Time
}
type FILTER_NUMBER struct {
	FILTER bool
	WAY    string // here you have some method : between, not_between, bigger, smaller, equal_to_one_of_them
	BIG    float64
	SMALL  float64
}
type FILTER_STRING struct {
	FILTER bool
	WAY    string // here you have some method : in_slice, not_in_slice, contain_one_in_slice , not_contain_one_in_slice
	SLICE  []string
}
type FILTER_BOOL struct {
	FILTER     bool
	BOOL_VALUE bool
}
type FINANCIAL_ACCOUNTING struct {
	ASSETS                    []string
	CURRENT_ASSETS            []string
	CASH_AND_CASH_EQUIVALENTS []string
	SHORT_TERM_INVESTMENTS    []string
	RECEIVABLES               []string
	INVENTORY                 []string
	LIABILITIES               []string
	CURRENT_LIABILITIES       []string
	EQUITY                    []string
	RETAINED_EARNINGS         []string
	DIVIDENDS                 []string
	INCOME_STATEMENT          []string
	EBITDA                    []string
	SALES                     []string
	COST_OF_GOODS_SOLD        []string
	DISCOUNTS                 []string
	INVOICE_DISCOUNT          []string
	INTEREST_EXPENSE          []string
}
type JOURNAL_TAG struct {
	IS_REVERSE                    bool    // this is true to the new entry when you enter reverse old entry
	IS_REVERSED                   bool    // this is true to the old entry when you enter reverse old entry
	REVERSE_ENTRY_NUMBER_COMPOUND int     // that mean if this is reverse entry what the entry compound was reversed
	REVERSE_ENTRY_NUMBER_SIMPLE   int     // that mean if this is reverse entry what the entry simple was reversed
	ENTRY_NUMBER_COMPOUND         int     // that mean number entry you made
	ENTRY_NUMBER_SIMPLE           int     // that mean the index of the simple entry in the you made
	VALUE                         float64 // this sould be positive
	PRICE_DEBIT                   float64 // this sould be positive
	PRICE_CREDIT                  float64 // this sould be positive
	QUANTITY_DEBIT                float64 // this sould be positive
	QUANTITY_CREDIT               float64 // this sould be positive
	ACCOUNT_DEBIT                 string  // the account name in the debit side
	ACCOUNT_CREDIT                string  // the account name in the credit side
	NOTES                         string  // your nots on the entry
	NAME                          string  // the name of the dealer or customer
	NAME_EMPLOYEE                 string  // the name of the employee that made the entry
}
type INVENTORY_TAG struct {
	PRICE        float64
	QUANTITY     float64
	ACCOUNT_NAME string
}
type INVOICE_STRUCT struct {
	VALUE        float64
	PRICE        float64
	QUANTITY     float64
	ACCOUNT_NAME string
}
type FINANCIAL_ANALYSIS struct {
	CURRENT_ASSETS                             float64
	CURRENT_LIABILITIES                        float64
	CASH                                       float64
	SHORT_TERM_INVESTMENTS                     float64
	NET_RECEIVABLES                            float64
	NET_CREDIT_SALES                           float64
	AVERAGE_NET_RECEIVABLES                    float64
	COST_OF_GOODS_SOLD                         float64
	AVERAGE_INVENTORY                          float64
	NET_INCOME                                 float64
	NET_SALES                                  float64
	AVERAGE_ASSETS                             float64
	AVERAGE_EQUITY                             float64
	PREFERRED_DIVIDENDS                        float64
	AVERAGE_COMMON_STOCKHOLDERS_EQUITY         float64
	MARKET_PRICE_PER_SHARES_OUTSTANDING        float64
	CASH_DIVIDENDS                             float64
	TOTAL_DEBT                                 float64
	TOTAL_ASSETS                               float64
	EBITDA                                     float64
	INTEREST_EXPENSE                           float64
	WEIGHTED_AVERAGE_COMMON_SHARES_OUTSTANDING float64
}
type FINANCIAL_ANALYSIS_STATEMENT struct {
	CURRENT_RATIO                        float64 // CURRENT_ASSETS / CURRENT_LIABILITIES
	ACID_TEST                            float64 // (CASH + SHORT_TERM_INVESTMENTS + NET_RECEIVABLES) / CURRENT_LIABILITIES
	RECEIVABLES_TURNOVER                 float64 // NET_CREDIT_SALES / AVERAGE_NET_RECEIVABLES
	INVENTORY_TURNOVER                   float64 // COST_OF_GOODS_SOLD / AVERAGE_INVENTORY
	PROFIT_MARGIN                        float64 // NET_INCOME / NET_SALES
	ASSET_TURNOVER                       float64 // NET_SALES / AVERAGE_ASSETS
	RETURN_ON_ASSETS                     float64 // NET_INCOME / AVERAGE_ASSETS
	RETURN_ON_EQUITY                     float64 // NET_INCOME / AVERAGE_EQUITY
	PAYOUT_RATIO                         float64 // CASH_DIVIDENDS / NET_INCOME
	DEBT_TO_TOTAL_ASSETS_RATIO           float64 // TOTAL_DEBT / TOTAL_ASSETS
	TIMES_INTEREST_EARNED                float64 // EBITDA / INTEREST_EXPENSE
	RETURN_ON_COMMON_STOCKHOLDERS_EQUITY float64 // (NET_INCOME - PREFERRED_DIVIDENDS) / AVERAGE_COMMON_STOCKHOLDERS_EQUITY
	EARNINGS_PER_SHARE                   float64 // (NET_INCOME - PREFERRED_DIVIDENDS) / WEIGHTED_AVERAGE_COMMON_SHARES_OUTSTANDING
	PRICE_EARNINGS_RATIO                 float64 // MARKET_PRICE_PER_SHARES_OUTSTANDING / EARNINGS_PER_SHARE
}
type PRICE_QUANTITY_ACCOUNT_BARCODE struct {
	PRICE        float64
	QUANTITY     float64
	ACCOUNT_NAME string
	BARCODE      string
}
type PRICE_QUANTITY_ACCOUNT struct {
	IS_CREDIT      bool
	COST_FLOW_TYPE string
	ACCOUNT_NAME   string
	PRICE          float64
	QUANTITY       float64
}
type ONE_STEP_DISTRIBUTION struct {
	SALES_OR_VARIABLE_OR_FIXED string
	DISTRIBUTION_METHOD        string
	AMOUNT                     float64
	FROM                       map[string]float64
	TO                         map[string]float64
}
