package anti_accountants

import (
	"os"
	"text/tabwriter"
	"time"
)

const (
	db_accounts           = "db_accounts"
	db_journal            = "db_journal"
	db_inventory          = "db_inventory"
	db_last_key_accounts  = "db_last_key_accounts"
	db_last_key_journal   = "db_last_key_journal"
	db_last_key_inventory = "db_last_key_inventory"
)

var (
	// exportable
	NOW                     = time.Now()
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
		{false, false, false, "", "ASSETS", "", []string{}, []string{"nojdsjdpq"}, [][]uint{{1}, {}}, []uint{1, 0}, [][]string{{}, {}}},
		{false, false, false, "", "CURRENT_ASSETS", "", []string{}, []string{"sijadpodjpao", "kaslajs"}, [][]uint{{1, 1}, {}}, []uint{2, 0}, [][]string{{"ASSETS"}, {}}},
		{true, false, false, "fifo", "CASH_AND_CASH_EQUIVALENTS", "", []string{}, []string{"888"}, [][]uint{{1, 1, 1}, {2}}, []uint{3, 1}, [][]string{{"ASSETS", "CURRENT_ASSETS"}, {}}},
		{true, false, false, "fifo", "SHORT_TERM_INVESTMENTS", "", []string{}, []string{}, [][]uint{{1, 2}, {5}}, []uint{2, 1}, [][]string{{"ASSETS"}, {}}},
		{true, false, false, "", "RECEIVABLES", "", []string{}, []string{}, [][]uint{{1, 3}, {}}, []uint{2, 0}, [][]string{{"ASSETS"}, {}}},
		{true, false, false, "wma", "INVENTORY", "", []string{}, []string{}, [][]uint{{1, 4}, {2, 4}}, []uint{2, 2}, [][]string{{"ASSETS"}, {"CASH_AND_CASH_EQUIVALENTS"}}},
		{false, true, false, "", "LIABILITIES", "", []string{}, []string{}, [][]uint{{2}, {}}, []uint{1, 0}, [][]string{{}, {}}},
		{true, true, false, "", "CURRENT_LIABILITIES", "", []string{}, []string{}, [][]uint{{2, 1}, {4}}, []uint{2, 1}, [][]string{{"LIABILITIES"}, {}}},
		{false, true, false, "", "EQUITY", "", []string{}, []string{}, [][]uint{{3}, {}}, []uint{1, 0}, [][]string{{}, {}}},
		{false, true, false, "", "RETAINED_EARNINGS", "", []string{}, []string{}, [][]uint{{3, 1}, {}}, []uint{2, 0}, [][]string{{"EQUITY"}, {}}},
		{true, false, true, "", "DIVIDENDS", "", []string{}, []string{}, [][]uint{{3, 1, 1}, {5, 2}}, []uint{3, 2}, [][]string{{"EQUITY", "RETAINED_EARNINGS"}, {"SHORT_TERM_INVESTMENTS"}}},
		{false, true, false, "", "INCOME_STATEMENT", "", []string{}, []string{}, [][]uint{{3, 1, 2}, {5, 3}}, []uint{3, 2}, [][]string{{"EQUITY", "RETAINED_EARNINGS"}, {"SHORT_TERM_INVESTMENTS"}}},
		{false, true, false, "", "EBITDA", "", []string{}, []string{}, [][]uint{{3, 1, 2, 1}, {}}, []uint{4, 0}, [][]string{{"EQUITY", "RETAINED_EARNINGS", "INCOME_STATEMENT"}, {}}},
		{true, true, true, "", "SALES", "", []string{}, []string{}, [][]uint{{3, 1, 2, 1, 1}, {5, 3, 2}}, []uint{5, 3}, [][]string{{"EQUITY", "RETAINED_EARNINGS", "INCOME_STATEMENT", "EBITDA"}, {"SHORT_TERM_INVESTMENTS", "INCOME_STATEMENT"}}},
		{true, false, true, "", "COST_OF_GOODS_SOLD", "", []string{}, []string{}, [][]uint{{3, 1, 2, 1, 2}, {5, 3, 6}}, []uint{5, 3}, [][]string{{"EQUITY", "RETAINED_EARNINGS", "INCOME_STATEMENT", "EBITDA"}, {"SHORT_TERM_INVESTMENTS", "INCOME_STATEMENT"}}},
		{false, false, false, "", "DISCOUNTS", "", []string{}, []string{}, [][]uint{{3, 1, 2, 1, 3}, {}}, []uint{5, 0}, [][]string{{"EQUITY", "RETAINED_EARNINGS", "INCOME_STATEMENT", "EBITDA"}, {}}},
		{true, false, true, "", "INVOICE_DISCOUNT", "", []string{}, []string{}, [][]uint{{3, 1, 2, 1, 3, 1}, {6}}, []uint{6, 1}, [][]string{{"EQUITY", "RETAINED_EARNINGS", "INCOME_STATEMENT", "EBITDA", "DISCOUNTS"}, {}}},
	}
	// var

	// const
	print_table          = tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	standard_days        = []string{"Saturday", "Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday"}
	adjusting_methods    = []string{"linear", "exponential", "logarithmic", "expire"}
	depreciation_methods = []string{"linear", "exponential", "logarithmic"}
	cost_flow_type       = []string{"fifo", "lifo", "wma"}
)

type ACCOUNT struct {
	IS_LOW_LEVEL_ACCOUNT             bool
	IS_CREDIT                        bool
	IS_TEMPORARY                     bool
	COST_FLOW_TYPE                   string
	ACCOUNT_NAME                     string
	NOTES                            string
	IMAGE                            []string
	BARCODE                          []string
	ACCOUNT_NUMBER                   [][]uint
	ACCOUNT_LEVELS                   []uint
	FATHER_AND_GRANDPA_ACCOUNTS_NAME [][]string
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
	ENTRY_NUMBER          uint
	ENTRY_NUMBER_COMPOUND uint
	ENTRY_NUMBER_SIMPLE   uint
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
	DATE_START            time.Time
	DATE_END              time.Time
	DATE_ENTRY            time.Time
}

type INVENTORY_TAG struct {
	PRICE      float64
	QUANTITY   float64
	ACCOUNT    string
	DATE_START time.Time
	DATE_END   time.Time
}

type INVOICE_STRUCT struct {
	VALUE    float64
	PRICE    float64
	QUANTITY float64
	ACCOUNT  string
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

type VALUE_PRICE_QUANTITY_ACCOUNT_BARCODE struct {
	VALUE    float64
	PRICE    float64
	QUANTITY float64
	ACCOUNT  string
	BARCODE  string
}

type ONE_STEP_DISTRIBUTION struct {
	SALES_OR_VARIABLE_OR_FIXED string
	DISTRIBUTION_METHOD        string
	AMOUNT                     float64
	FROM                       map[string]float64
	TO                         map[string]float64
}
