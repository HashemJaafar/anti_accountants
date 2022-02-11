package anti_accountants

import (
	"log"
	"time"
)

func error_account_name_is_not_listed(account_name string) {
	log.Panic(account_name, " is not listed")
}

func error_not_connected_tree(account ACCOUNT) {
	log.Panic("you can't use these account number ", account, " because there is no account with account number", account.ACCOUNT_NUMBER[:len(account.ACCOUNT_NUMBER)-1])
}

func error_element_is_not_in_elements(element string, elements []string) {
	log.Panic(element, " is not in ", elements)
}

func error_not_equal(m map[string]float64, a, b, sign, c string) {
	log.Fatal(a, " ", m[a], " != ", b, " ", m[b], " ", sign, " ", c, " ", m[c])
}

func error_duplicate_value(duplicated_element interface{}) {
	log.Panic(duplicated_element, " is duplicated values in the fields of FINANCIAL_ACCOUNTING and that make error. you should remove the duplicate")
}

func error_you_cant_change_the_name(name, new_name string) {
	log.Panic("you can't change the name of [", name, "] to [", new_name, "] as new name because it used")
}

func error_should_be_one_of_the_fathers(account_name, sub_account_name string) {
	log.Panic(account_name, " should be one of the fathers of ", sub_account_name)
}

func error_one_credit___one_debit(or_and string, entries ...[]ACCOUNT_VALUE_PRICE_QUANTITY_BARCODE) {
	log.Panic("should be one credit ", or_and, " one debit in the entry ", entries)
}

func error_debit_not_equal_credit(difference float64, entries []ACCOUNT_VALUE_PRICE_QUANTITY_BARCODE) {
	log.Panic(difference, " not equal 0 if the number>0 it means debit overstated else credit overstated debit-credit should equal zero ", entries)
}

func error_the_time_is_not_in_range(element DAY_START_END, min_max_time_number int) {
	log.Panic("error ", element.START_HOUR, " for ", element, " is < ", min_max_time_number)
}

func error_the_barcode_is_wrong(entry ACCOUNT_VALUE_PRICE_QUANTITY_BARCODE) {
	log.Panic("the barcode is wrong for ", entry)
}

func error_is_high_level_account(account_name string) {
	log.Panic(account_name, " is high level account that mean you can't used in the entry")
}

func error_this_entry_not_exist() {
	log.Panic("this entry not exist")
}

func error_smaller_than_or_equal(start_date, end_date time.Time) {
	log.Panic("error please enter the ", start_date, " <= ", end_date)
}

func error_date_layout(string_date string) {
	log.Panic("you don't have layout for this date ", string_date)
}

func error_panic(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func error_cost_flow_type_used_with___account(account ACCOUNT, error_state string) {
	log.Panic("you can't use cost flow type ", account.COST_FLOW_TYPE, " with ", error_state, " accounts like ", account)
}

func error_you_cant_use_entry_expire() {
	log.Panic("you can't use entry expire equal to zero with adjusting methods ", adjusting_methods)
}

func error_you_cant_use_depreciation_methods_with_inventory(account_name string) {
	log.Panic("you just can use ", []string{"expire", ""}, " with ", account_name, " because it is inventory account")
}

func error_make_nagtive_balance(entry ACCOUNT_VALUE_PRICE_QUANTITY_BARCODE, account_balance float64) {
	log.Panic("you can't enter ", entry, " because you have ", account_balance, " and that will make the balance of ", entry.ACCOUNT, " negative ", account_balance+entry.VALUE, " and that you just can do it in equity_normal accounts not other accounts")
}

func error_the_price_should_be_positive(entry ACCOUNT_VALUE_PRICE_QUANTITY_BARCODE) {
	log.Panic("the price of ", entry, " should be positive")
}

func error_the_order_out_of_stock(quantity, quantity_count float64, account, barcode string) {
	log.Panic("you order ", quantity, " but you have ", quantity-quantity_count, " ", account, " with barcode ", barcode)
}

func error_you_should_use_cost_flow_type(account_name string) {
	log.Panic("you should use cost flow type for account ", account_name)
}
