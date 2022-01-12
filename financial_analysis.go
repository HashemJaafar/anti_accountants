package anti_accountants

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
	CURRENT_RATIO                        float64 // current_assets / current_liabilities
	ACID_TEST                            float64 // (cash + short_term_investments + net_receivables) / current_liabilities
	RECEIVABLES_TURNOVER                 float64 // net_credit_sales / average_net_receivables
	INVENTORY_TURNOVER                   float64 // cost_of_goods_sold / average_inventory
	ASSET_TURNOVER                       float64 // net_sales / average_assets
	PROFIT_MARGIN                        float64 // net_income / net_sales
	RETURN_ON_ASSETS                     float64 // net_income / average_assets
	RETURN_ON_EQUITY                     float64 // net_income / average_equity
	PAYOUT_RATIO                         float64 // cash_dividends / net_income
	DEBT_TO_TOTAL_ASSETS_RATIO           float64 // total_debt / total_assets
	TIMES_INTEREST_EARNED                float64 // ebitda / interest_expense
	RETURN_ON_COMMON_STOCKHOLDERS_EQUITY float64 // (net_income - preferred_dividends) / average_common_stockholders_equity
	EARNINGS_PER_SHARE                   float64 // (net_income - preferred_dividends) / weighted_average_common_shares_outstanding
	PRICE_EARNINGS_RATIO                 float64 // market_price_per_shares_outstanding / earnings_per_share
}

func (s FINANCIAL_ANALYSIS) FINANCIAL_ANALYSIS_STATEMENT() FINANCIAL_ANALYSIS_STATEMENT {
	CURRENT_RATIO := s.CURRENT_ASSETS / s.CURRENT_LIABILITIES
	ACID_TEST := (s.CASH + s.SHORT_TERM_INVESTMENTS + s.NET_RECEIVABLES) / s.CURRENT_LIABILITIES
	RECEIVABLES_TURNOVER := s.NET_CREDIT_SALES / s.AVERAGE_NET_RECEIVABLES
	INVENTORY_TURNOVER := s.COST_OF_GOODS_SOLD / s.AVERAGE_INVENTORY
	PROFIT_MARGIN := s.NET_INCOME / s.NET_SALES
	ASSET_TURNOVER := s.NET_SALES / s.AVERAGE_ASSETS
	RETURN_ON_ASSETS := s.NET_INCOME / s.AVERAGE_ASSETS
	RETURN_ON_EQUITY := s.NET_INCOME / s.AVERAGE_EQUITY
	PAYOUT_RATIO := s.CASH_DIVIDENDS / s.NET_INCOME
	DEBT_TO_TOTAL_ASSETS_RATIO := s.TOTAL_DEBT / s.TOTAL_ASSETS
	TIMES_INTEREST_EARNED := s.EBITDA / s.INTEREST_EXPENSE
	RETURN_ON_COMMON_STOCKHOLDERS_EQUITY := (s.NET_INCOME - s.PREFERRED_DIVIDENDS) / s.AVERAGE_COMMON_STOCKHOLDERS_EQUITY
	EARNINGS_PER_SHARE := (s.NET_INCOME - s.PREFERRED_DIVIDENDS) / s.WEIGHTED_AVERAGE_COMMON_SHARES_OUTSTANDING
	PRICE_EARNINGS_RATIO := s.MARKET_PRICE_PER_SHARES_OUTSTANDING / EARNINGS_PER_SHARE
	return FINANCIAL_ANALYSIS_STATEMENT{
		CURRENT_RATIO:                        CURRENT_RATIO,
		ACID_TEST:                            ACID_TEST,
		RECEIVABLES_TURNOVER:                 RECEIVABLES_TURNOVER,
		INVENTORY_TURNOVER:                   INVENTORY_TURNOVER,
		ASSET_TURNOVER:                       ASSET_TURNOVER,
		PROFIT_MARGIN:                        PROFIT_MARGIN,
		RETURN_ON_ASSETS:                     RETURN_ON_ASSETS,
		RETURN_ON_EQUITY:                     RETURN_ON_EQUITY,
		PAYOUT_RATIO:                         PAYOUT_RATIO,
		DEBT_TO_TOTAL_ASSETS_RATIO:           DEBT_TO_TOTAL_ASSETS_RATIO,
		TIMES_INTEREST_EARNED:                TIMES_INTEREST_EARNED,
		RETURN_ON_COMMON_STOCKHOLDERS_EQUITY: RETURN_ON_COMMON_STOCKHOLDERS_EQUITY,
		EARNINGS_PER_SHARE:                   EARNINGS_PER_SHARE,
		PRICE_EARNINGS_RATIO:                 PRICE_EARNINGS_RATIO}
}

func (s FINANCIAL_ACCOUNTING) ANALYSIS(statements []map[string]map[string]map[string]map[string]map[string]float64) []FINANCIAL_ANALYSIS_STATEMENT {
	var all_analysis []FINANCIAL_ANALYSIS_STATEMENT
	for _, statement := range statements {
		analysis := FINANCIAL_ANALYSIS{
			CURRENT_ASSETS:                      statement[s.CASH_AND_CASH_EQUIVALENTS][s.CURRENT_ASSETS]["names"]["value"]["ending_balance"],
			CURRENT_LIABILITIES:                 statement[s.CASH_AND_CASH_EQUIVALENTS][s.CURRENT_LIABILITIES]["names"]["value"]["ending_balance"],
			CASH:                                statement[s.CASH_AND_CASH_EQUIVALENTS][s.CASH_AND_CASH_EQUIVALENTS]["names"]["value"]["ending_balance"],
			SHORT_TERM_INVESTMENTS:              statement[s.CASH_AND_CASH_EQUIVALENTS][s.SHORT_TERM_INVESTMENTS]["names"]["value"]["ending_balance"],
			NET_RECEIVABLES:                     statement[s.CASH_AND_CASH_EQUIVALENTS][s.RECEIVABLES]["names"]["value"]["ending_balance"],
			NET_CREDIT_SALES:                    statement[s.SALES][s.RECEIVABLES]["names"]["value"]["flow"],
			AVERAGE_NET_RECEIVABLES:             statement[s.CASH_AND_CASH_EQUIVALENTS][s.RECEIVABLES]["names"]["value"]["average"],
			COST_OF_GOODS_SOLD:                  statement[s.CASH_AND_CASH_EQUIVALENTS][s.COST_OF_GOODS_SOLD]["names"]["value"]["ending_balance"],
			AVERAGE_INVENTORY:                   statement[s.CASH_AND_CASH_EQUIVALENTS][s.INVENTORY]["names"]["value"]["average"],
			NET_INCOME:                          statement[s.CASH_AND_CASH_EQUIVALENTS][s.INCOME_STATEMENT]["names"]["value"]["ending_balance"],
			NET_SALES:                           statement[s.CASH_AND_CASH_EQUIVALENTS][s.SALES]["names"]["value"]["ending_balance"],
			AVERAGE_ASSETS:                      statement[s.CASH_AND_CASH_EQUIVALENTS][s.ASSETS]["names"]["value"]["average"],
			AVERAGE_EQUITY:                      statement[s.CASH_AND_CASH_EQUIVALENTS][s.EQUITY]["names"]["value"]["average"],
			PREFERRED_DIVIDENDS:                 0,
			AVERAGE_COMMON_STOCKHOLDERS_EQUITY:  0,
			MARKET_PRICE_PER_SHARES_OUTSTANDING: 0,
			CASH_DIVIDENDS:                      statement[s.CASH_AND_CASH_EQUIVALENTS][s.DIVIDENDS]["names"]["value"]["flow"],
			TOTAL_DEBT:                          statement[s.CASH_AND_CASH_EQUIVALENTS][s.LIABILITIES]["names"]["value"]["ending_balance"],
			TOTAL_ASSETS:                        statement[s.CASH_AND_CASH_EQUIVALENTS][s.ASSETS]["names"]["value"]["ending_balance"],
			EBITDA:                              statement[s.CASH_AND_CASH_EQUIVALENTS][s.EBITDA]["names"]["value"]["ending_balance"],
			INTEREST_EXPENSE:                    statement[s.CASH_AND_CASH_EQUIVALENTS][s.INTEREST_EXPENSE]["names"]["value"]["ending_balance"],
			WEIGHTED_AVERAGE_COMMON_SHARES_OUTSTANDING: 0,
		}.FINANCIAL_ANALYSIS_STATEMENT()
		all_analysis = append(all_analysis, analysis)
	}
	return all_analysis
}
