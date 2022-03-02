package anti_accountants

import (
	"database/sql"
	"os"
	"text/tabwriter"
	"time"
)

var (
	// exportable
	NOW                    = time.Now()
	DB                     *sql.DB
	INVOICE_DISCOUNTS_LIST [][2]float64
	AUTO_COMPLETE_ENTRIES  []AUTO_COMPLETE_ENTRIE
	DATE_LAYOUT            = []string{"2006-01-02 15:04:05.999999999 -0700 +03 m=+0.999999999", "2006-01-02 15:04:05.999999999 -0700 +03"}
	PRIMARY_ACCOUNTS_NAMES = FINANCIAL_ACCOUNTING{
		ASSETS:                    "ASSETS",
		CURRENT_ASSETS:            "CURRENT_ASSETS",
		CASH_AND_CASH_EQUIVALENTS: "CASH_AND_CASH_EQUIVALENTS",
		SHORT_TERM_INVESTMENTS:    "SHORT_TERM_INVESTMENTS",
		RECEIVABLES:               "RECEIVABLES",
		INVENTORY:                 "INVENTORY",
		LIABILITIES:               "LIABILITIES",
		CURRENT_LIABILITIES:       "CURRENT_LIABILITIES",
		EQUITY:                    "EQUITY",
		RETAINED_EARNINGS:         "RETAINED_EARNINGS",
		DIVIDENDS:                 "DIVIDENDS",
		INCOME_STATEMENT:          "INCOME_STATEMENT",
		EBITDA:                    "EBITDA",
		SALES:                     "SALES",
		COST_OF_GOODS_SOLD:        "COST_OF_GOODS_SOLD",
		DISCOUNTS:                 "DISCOUNTS",
		INVOICE_DISCOUNT:          "INVOICE_DISCOUNT",
		INTEREST_EXPENSE:          "INTEREST_EXPENSE",
	}
	ACCOUNTS = []ACCOUNT{
		{false, "", "ASSETS", []uint{1}},
		{false, "", "CURRENT_ASSETS", []uint{1, 1}},
		{false, "", "CASH_AND_CASH_EQUIVALENTS", []uint{1, 1, 1}},
		{false, "", "SHORT_TERM_INVESTMENTS", []uint{1, 2}},
		{false, "", "RECEIVABLES", []uint{1, 3}},
		{false, "", "INVENTORY", []uint{1, 4}},
		{false, "", "LIABILITIES", []uint{2}},
		{false, "", "CURRENT_LIABILITIES", []uint{2, 1}},
		{false, "", "EQUITY", []uint{3}},
		{false, "", "RETAINED_EARNINGS", []uint{3, 1}},
		{false, "", "DIVIDENDS", []uint{3, 1, 1}},
		{false, "", "INCOME_STATEMENT", []uint{3, 1, 2}},
		{false, "", "EBITDA", []uint{3, 1, 2, 1}},
		{false, "", "SALES", []uint{3, 1, 2, 1, 1}},
		{false, "", "COST_OF_GOODS_SOLD", []uint{3, 1, 2, 1, 2}},
		{false, "", "DISCOUNTS", []uint{3, 1, 2, 1, 3}},
		{false, "", "INVOICE_DISCOUNT", []uint{3, 1, 2, 1, 3, 1}},
		{false, "", "INTEREST_EXPENSE", []uint{3, 1, 2, 1, 4}},
	}

	// var
	inventory []string

	// const
	print_table          = tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	standard_days        = []string{"Saturday", "Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday"}
	adjusting_methods    = []string{"linear", "exponential", "logarithmic", "expire"}
	depreciation_methods = []string{"linear", "exponential", "logarithmic"}
	cost_flow_type       = []string{"fifo", "lifo", "wma", "barcode"}
)

type ACCOUNT struct {
	IS_CREDIT                    bool
	COST_FLOW_TYPE, ACCOUNT_NAME string
	ACCOUNT_NUMBER               []uint
}

type DAY_START_END struct {
	DAY          string
	START_HOUR   int
	START_MINUTE int
	END_HOUR     int
	END_MINUTE   int
}

type start_end_minutes struct {
	start_date time.Time
	end_date   time.Time
	minutes    float64
}

type AUTO_COMPLETE_ENTRIE struct {
	INVENTORY_ACCOUNT, COST_OF_GOOD_SOLD_ACCOUNT, REVENUE_ACCOUNT, DESCOUNT_ACCOUNT string
	SELLING_PRICE, DESCOUNT_PRICE                                                   float64
}

type FILTERED_STATEMENT struct {
	KEY_ACCOUNT_FLOW, KEY_ACCOUNT, KEY_NAME, KEY_VPQ, KEY_NUMBER string
	NUMBER                                                       float64
}

type FINANCIAL_ACCOUNTING struct {
	ASSETS                    string
	CURRENT_ASSETS            string
	CASH_AND_CASH_EQUIVALENTS string
	SHORT_TERM_INVESTMENTS    string
	RECEIVABLES               string
	INVENTORY                 string
	LIABILITIES               string
	CURRENT_LIABILITIES       string
	EQUITY                    string
	RETAINED_EARNINGS         string
	DIVIDENDS                 string
	INCOME_STATEMENT          string
	EBITDA                    string
	SALES                     string
	COST_OF_GOODS_SOLD        string
	DISCOUNTS                 string
	INVOICE_DISCOUNT          string
	INTEREST_EXPENSE          string
}

type JOURNAL_TAG struct {
	DATE          string
	ENTRY_NUMBER  int
	ACCOUNT       string
	VALUE         float64
	PRICE         float64
	QUANTITY      float64
	BARCODE       string
	ENTRY_EXPAIR  string
	DESCRIPTION   string
	NAME          string
	EMPLOYEE_NAME string
	ENTRY_DATE    string
	REVERSE       bool
}

type INVOICE_STRUCT struct {
	ACCOUNT                string
	VALUE, PRICE, QUANTITY float64
}

type FINANCIAL_ANALYSIS struct {
	CURRENT_ASSETS,
	CURRENT_LIABILITIES,
	CASH,
	SHORT_TERM_INVESTMENTS,
	NET_RECEIVABLES,
	NET_CREDIT_SALES,
	AVERAGE_NET_RECEIVABLES,
	COST_OF_GOODS_SOLD,
	AVERAGE_INVENTORY,
	NET_INCOME,
	NET_SALES,
	AVERAGE_ASSETS,
	AVERAGE_EQUITY,
	PREFERRED_DIVIDENDS,
	AVERAGE_COMMON_STOCKHOLDERS_EQUITY,
	MARKET_PRICE_PER_SHARES_OUTSTANDING,
	CASH_DIVIDENDS,
	TOTAL_DEBT,
	TOTAL_ASSETS,
	EBITDA,
	INTEREST_EXPENSE,
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

type ACCOUNT_VALUE_PRICE_QUANTITY_BARCODE struct {
	ACCOUNT  string
	VALUE    float64
	PRICE    float64
	QUANTITY float64
	BARCODE  string
}

type ONE_STEP_DISTRIBUTION struct {
	SALES_OR_VARIABLE_OR_FIXED, DISTRIBUTION_METHOD string
	AMOUNT                                          float64
	FROM, TO                                        map[string]float64
}
