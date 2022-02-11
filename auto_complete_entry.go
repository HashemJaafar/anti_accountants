package anti_accountants

import (
	"fmt"
	"math"
)

type AUTO_COMPLETE_ENTRIE struct {
	INVENTORY_ACCOUNT, COST_OF_GOOD_SOLD_ACCOUNT, REVENUE_ACCOUNT, DESCOUNT_ACCOUNT string
	SELLING_PRICE, DESCOUNT_PRICE                                                   float64
}

func (s FINANCIAL_ACCOUNTING) auto_completion_the_entry(entries []ACCOUNT_VALUE_PRICE_QUANTITY_BARCODE) []ACCOUNT_VALUE_PRICE_QUANTITY_BARCODE {
	var new_entries []ACCOUNT_VALUE_PRICE_QUANTITY_BARCODE
	for _, entry := range entries {
		for _, complement := range s.AUTO_COMPLETE_ENTRIES {
			if complement.INVENTORY_ACCOUNT == entry.ACCOUNT {
				fmt.Println(complement)
				new_entries = append(new_entries, ACCOUNT_VALUE_PRICE_QUANTITY_BARCODE{
					ACCOUNT:  complement.COST_OF_GOOD_SOLD_ACCOUNT,
					VALUE:    entry.VALUE,
					PRICE:    entry.PRICE,
					QUANTITY: entry.QUANTITY,
					BARCODE:  entry.BARCODE,
				})
				new_entries = append(new_entries, ACCOUNT_VALUE_PRICE_QUANTITY_BARCODE{
					ACCOUNT:  entry.ACCOUNT,
					VALUE:    entry.VALUE,
					PRICE:    entry.PRICE,
					QUANTITY: entry.QUANTITY,
					BARCODE:  entry.BARCODE,
				})
				new_entries = append(new_entries, ACCOUNT_VALUE_PRICE_QUANTITY_BARCODE{
					ACCOUNT:  entry.ACCOUNT,
					VALUE:    entry.VALUE,
					PRICE:    entry.PRICE,
					QUANTITY: entry.QUANTITY,
					BARCODE:  entry.BARCODE,
				})
				new_entries = append(new_entries, ACCOUNT_VALUE_PRICE_QUANTITY_BARCODE{
					ACCOUNT:  complement.REVENUE_ACCOUNT,
					VALUE:    entry.QUANTITY * complement.SELLING_PRICE,
					PRICE:    complement.SELLING_PRICE,
					QUANTITY: entry.QUANTITY,
					BARCODE:  entry.BARCODE,
				})
			}
		}
	}
	fmt.Println(new_entries)
	return entries
}

func (s FINANCIAL_ACCOUNTING) auto_completion_the_invoice_discount(entries []ACCOUNT_VALUE_PRICE_QUANTITY_BARCODE) []ACCOUNT_VALUE_PRICE_QUANTITY_BARCODE {
	total_invoice_before_invoice_discount := s.total_invoice_before_invoice_discount(entries)
	_, discount := X_UNDER_X(s.INVOICE_DISCOUNTS_LIST, total_invoice_before_invoice_discount)
	invoice_discount := discount_tax_calculator(total_invoice_before_invoice_discount, discount)
	entries = append(entries, ACCOUNT_VALUE_PRICE_QUANTITY_BARCODE{
		ACCOUNT:  s.INVOICE_DISCOUNT,
		VALUE:    invoice_discount,
		PRICE:    invoice_discount,
		QUANTITY: 1,
	})
	return entries
}

func (s FINANCIAL_ACCOUNTING) total_invoice_before_invoice_discount(entries []ACCOUNT_VALUE_PRICE_QUANTITY_BARCODE) float64 {
	var total_invoice_before_invoice_discount float64
	for _, entry := range entries {
		if s.is_it_sub_account_using_name(s.INCOME_STATEMENT, entry.ACCOUNT) && s.is_credit(entry.ACCOUNT) {
			total_invoice_before_invoice_discount += entry.VALUE
		} else if s.is_it_sub_account_using_name(s.DISCOUNTS, entry.ACCOUNT) && !s.is_credit(entry.ACCOUNT) {
			total_invoice_before_invoice_discount -= entry.VALUE
		}
	}
	return total_invoice_before_invoice_discount
}

func discount_tax_calculator(price, discount_tax float64) float64 {
	if discount_tax < 0 {
		discount_tax = math.Abs(discount_tax)
	} else if discount_tax > 0 {
		discount_tax = price * discount_tax
	}
	return discount_tax
}
