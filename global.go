package anti_accountants

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
)

var (
	COMPANY_NAME            = "anti_accountants"
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
	PRINT_TABLE          = tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	STANDARD_DAYS        = []string{SATURDAY, SUNDAY, MONDAY, TUESDAY, WEDNESDAY, THURSDAY, FRIDAY}
	DEPRECIATION_METHODS = []string{LINEAR, EXPONENTIAL, LOGARITHMIC}
	COST_FLOW_TYPE       = []string{FIFO, LIFO, WMA}
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
	INVENTORY_ACCOUNT         string
	COST_OF_GOOD_SOLD_ACCOUNT string
	REVENUE_ACCOUNT           string
	DESCOUNT_ACCOUNT          string
	SELLING_PRICE             float64
	DESCOUNT_PRICE            float64
}
type FILTERED_STATEMENT struct {
	KEY_ACCOUNT_FLOW string
	KEY_ACCOUNT      string
	KEY_NAME         string
	KEY_VPQ          string
	KEY_NUMBER       string
	NUMBER           float64
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
	REVERSE               bool
	ENTRY_NUMBER_COMPOUND int
	ENTRY_NUMBER_SIMPLE   int
	VALUE                 float64
	PRICE_DEBIT           float64
	PRICE_CREDIT          float64
	QUANTITY_DEBIT        float64
	QUANTITY_CREDIT       float64
	ACCOUNT_DEBIT         string
	ACCOUNT_CREDIT        string
	NOTES                 string
	NAME                  string
	NAME_EMPLOYEE         string
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
