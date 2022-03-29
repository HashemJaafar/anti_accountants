package anti_accountants

// import (
// 	"fmt"
// 	"math"
// )

// func auto_completion_the_entry(entries []PRICE_QUANTITY_ACCOUNT_BARCODE) []PRICE_QUANTITY_ACCOUNT_BARCODE {
// 	var new_entries []PRICE_QUANTITY_ACCOUNT_BARCODE
// 	for _, entry := range entries {
// 		for _, complement := range AUTO_COMPLETE_ENTRIES {
// 			if complement.INVENTORY_ACCOUNT == entry.ACCOUNT {
// 				fmt.Println(complement)
// 				new_entries = append(new_entries, PRICE_QUANTITY_ACCOUNT_BARCODE{
// 					ACCOUNT:  complement.COST_OF_GOOD_SOLD_ACCOUNT,
// 					VALUE:    entry.VALUE,
// 					PRICE:    entry.PRICE,
// 					QUANTITY: entry.QUANTITY,
// 					BARCODE:  entry.BARCODE,
// 				})
// 				new_entries = append(new_entries, PRICE_QUANTITY_ACCOUNT_BARCODE{
// 					ACCOUNT:  entry.ACCOUNT,
// 					VALUE:    entry.VALUE,
// 					PRICE:    entry.PRICE,
// 					QUANTITY: entry.QUANTITY,
// 					BARCODE:  entry.BARCODE,
// 				})
// 				new_entries = append(new_entries, PRICE_QUANTITY_ACCOUNT_BARCODE{
// 					ACCOUNT:  entry.ACCOUNT,
// 					VALUE:    entry.VALUE,
// 					PRICE:    entry.PRICE,
// 					QUANTITY: entry.QUANTITY,
// 					BARCODE:  entry.BARCODE,
// 				})
// 				new_entries = append(new_entries, PRICE_QUANTITY_ACCOUNT_BARCODE{
// 					ACCOUNT:  complement.REVENUE_ACCOUNT,
// 					VALUE:    entry.QUANTITY * complement.SELLING_PRICE,
// 					PRICE:    complement.SELLING_PRICE,
// 					QUANTITY: entry.QUANTITY,
// 					BARCODE:  entry.BARCODE,
// 				})
// 			}
// 		}
// 	}
// 	fmt.Println(new_entries)
// 	return entries
// }

// func auto_completion_the_invoice_discount(entries []PRICE_QUANTITY_ACCOUNT_BARCODE) []PRICE_QUANTITY_ACCOUNT_BARCODE {
// 	total_invoice_before_invoice_discount := total_invoice_before_invoice_discount(entries)
// 	_, discount := X_UNDER_X(INVOICE_DISCOUNTS_LIST, total_invoice_before_invoice_discount)
// 	invoice_discount := discount_tax_calculator(total_invoice_before_invoice_discount, discount)
// 	entries = append(entries, PRICE_QUANTITY_ACCOUNT_BARCODE{
// 		ACCOUNT:  PRIMARY_ACCOUNTS_NAMES.INVOICE_DISCOUNT,
// 		VALUE:    invoice_discount,
// 		PRICE:    invoice_discount,
// 		QUANTITY: 1,
// 	})
// 	return entries
// }

// func total_invoice_before_invoice_discount(entries []PRICE_QUANTITY_ACCOUNT_BARCODE) float64 {
// 	var total_invoice_before_invoice_discount float64
// 	for _, entry := range entries {
// 		if IS_IT_SUB_ACCOUNT_USING_NAME(PRIMARY_ACCOUNTS_NAMES.INCOME_STATEMENT, entry.ACCOUNT) && is_credit(entry.ACCOUNT) {
// 			total_invoice_before_invoice_discount += entry.VALUE
// 		} else if IS_IT_SUB_ACCOUNT_USING_NAME(PRIMARY_ACCOUNTS_NAMES.DISCOUNTS, entry.ACCOUNT) && !is_credit(entry.ACCOUNT) {
// 			total_invoice_before_invoice_discount -= entry.VALUE
// 		}
// 	}
// 	return total_invoice_before_invoice_discount
// }

// func discount_tax_calculator(price, discount_tax float64) float64 {
// 	if discount_tax < 0 {
// 		discount_tax = math.Abs(discount_tax)
// 	} else if discount_tax > 0 {
// 		discount_tax = price * discount_tax
// 	}
// 	return discount_tax
// }
