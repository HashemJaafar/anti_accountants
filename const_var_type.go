package anti_accountants

import (
	"database/sql"
	"os"
	"text/tabwriter"
	"time"
)

var (
	// exportable
	NOW                     = time.Now()
	DB                      *sql.DB
	INDEX_OF_ACCOUNT_NUMBER = 0
	INVOICE_DISCOUNTS_LIST  [][2]float64
	AUTO_COMPLETE_ENTRIES   []AUTO_COMPLETE_ENTRIE
	DATE_LAYOUT             = []string{"2006-01-02 15:04:05.999999999 -0700 +03 m=+0.999999999", "2006-01-02 15:04:05.999999999 -0700 +03"}
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
	ACCOUNTS = []ACCOUNT{
		{false, false, false, "", "ASSETS", "", [][]uint{{1}, {}}, []uint{1, 0}, [][]string{{}, {}}, ""},
		{false, false, false, "", "CURRENT_ASSETS", "", [][]uint{{1, 1}, {2}}, []uint{2, 1}, [][]string{{"ASSETS"}, {}}, ""},
		{true, false, false, "fifo", "CASH_AND_CASH_EQUIVALENTS", "", [][]uint{{1, 1, 1}, {2, 1}}, []uint{3, 2}, [][]string{{"ASSETS", "CURRENT_ASSETS"}, {"CURRENT_ASSETS"}}, ""},
		{true, false, false, "fifo", "SHORT_TERM_INVESTMENTS", "", [][]uint{{1, 2}, {2, 2}}, []uint{2, 2}, [][]string{{"ASSETS"}, {"CURRENT_ASSETS"}}, ""},
		{true, false, false, "", "RECEIVABLES", "", [][]uint{{1, 3}, {2, 3}}, []uint{2, 2}, [][]string{{"ASSETS"}, {"CURRENT_ASSETS"}}, ""},
		{true, false, false, "wma", "INVENTORY", "", [][]uint{{1, 4}, {2, 4}}, []uint{2, 2}, [][]string{{"ASSETS"}, {"CURRENT_ASSETS"}}, ""},
		{false, true, false, "", "LIABILITIES", "", [][]uint{{2}, {}}, []uint{1, 0}, [][]string{{}, {}}, ""},
		{true, true, false, "", "CURRENT_LIABILITIES", "", [][]uint{{2, 1}, {4}}, []uint{2, 1}, [][]string{{"LIABILITIES"}, {}}, ""},
		{false, true, false, "", "EQUITY", "", [][]uint{{3}, {5}}, []uint{1, 1}, [][]string{{}, {}}, ""},
		{false, true, false, "", "RETAINED_EARNINGS", "", [][]uint{{3, 1}, {}}, []uint{2, 0}, [][]string{{"EQUITY"}, {}}, ""},
		{true, false, true, "", "DIVIDENDS", "", [][]uint{{3, 1, 1}, {5, 2}}, []uint{3, 2}, [][]string{{"EQUITY", "RETAINED_EARNINGS"}, {"EQUITY"}}, ""},
		{false, true, false, "", "INCOME_STATEMENT", "", [][]uint{{3, 1, 2}, {5, 3}}, []uint{3, 2}, [][]string{{"EQUITY", "RETAINED_EARNINGS"}, {"EQUITY"}}, ""},
		{false, true, false, "", "EBITDA", "", [][]uint{{3, 1, 2, 1}, {}}, []uint{4, 0}, [][]string{{"EQUITY", "RETAINED_EARNINGS", "INCOME_STATEMENT"}, {}}, ""},
		{true, true, true, "", "SALES", "", [][]uint{{3, 1, 2, 1, 1}, {5, 3, 2}}, []uint{5, 3}, [][]string{{"EQUITY", "RETAINED_EARNINGS", "INCOME_STATEMENT", "EBITDA"}, {"EQUITY", "INCOME_STATEMENT"}}, ""},
		{true, false, true, "", "COST_OF_GOODS_SOLD", "", [][]uint{{3, 1, 2, 1, 2}, {5, 3, 6}}, []uint{5, 3}, [][]string{{"EQUITY", "RETAINED_EARNINGS", "INCOME_STATEMENT", "EBITDA"}, {"EQUITY", "INCOME_STATEMENT"}}, ""},
		{false, false, false, "", "DISCOUNTS", "", [][]uint{{3, 1, 2, 1, 3}, {}}, []uint{5, 0}, [][]string{{"EQUITY", "RETAINED_EARNINGS", "INCOME_STATEMENT", "EBITDA"}, {}}, ""},
		{true, false, true, "", "INVOICE_DISCOUNT", "", [][]uint{{3, 1, 2, 1, 3, 1}, {6}}, []uint{6, 1}, [][]string{{"EQUITY", "RETAINED_EARNINGS", "INCOME_STATEMENT", "EBITDA", "DISCOUNTS"}, {}}, ""}}
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
	is_low_level_account, IS_CREDIT, IS_TEMPORARY bool
	COST_FLOW_TYPE, ACCOUNT_NAME, DESCRIPTION     string
	ACCOUNT_NUMBER                                [][]uint
	account_levels                                []uint
	father_and_grandpa_accounts_name              [][]string
	IMAGE                                         string
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
	ASSETS,
	CURRENT_ASSETS,
	CASH_AND_CASH_EQUIVALENTS,
	SHORT_TERM_INVESTMENTS,
	RECEIVABLES,
	INVENTORY,
	LIABILITIES,
	CURRENT_LIABILITIES,
	EQUITY,
	RETAINED_EARNINGS,
	DIVIDENDS,
	INCOME_STATEMENT,
	EBITDA,
	SALES,
	COST_OF_GOODS_SOLD,
	DISCOUNTS,
	INVOICE_DISCOUNT,
	INTEREST_EXPENSE []string
}

type JOURNAL_TAG struct {
	REVERSE         bool
	ENTRY_NUMBER    int
	LINE_NUMBER     int
	VALUE           float64
	PRICE_DEBIT     float64
	PRICE_CREDIT    float64
	QUANTITY_DEBIT  float64
	QUANTITY_CREDIT float64
	ACCOUNT_DEBIT   string
	ACCOUNT_CREDIT  string
	BARCODE_DEBIT   string
	BARCODE_CREDIT  string
	DESCRIPTION     string
	NAME            string
	EMPLOYEE_NAME   string
	DATE            string
	ENTRY_EXPAIR    string
	ENTRY_DATE      string
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
