package anti_accountants

// func INVOICE(array_of_journal_tag []JOURNAL_TAG) []INVOICE_STRUCT {
// 	m := map[string]*INVOICE_STRUCT{}
// 	for _, entry := range array_of_journal_tag {
// 		var key string
// 		switch {
// 		case is_it_sub_account_using_name(PRIMARY_ACCOUNTS_NAMES.ASSETS, entry.ACCOUNT) && !is_credit(entry.ACCOUNT) && !is_in(entry.ACCOUNT, inventory) && entry.VALUE > 0:
// 			key = "total"
// 		case is_it_sub_account_using_name(PRIMARY_ACCOUNTS_NAMES.DISCOUNTS, entry.ACCOUNT) && !is_credit(entry.ACCOUNT):
// 			key = "total discounts"
// 		case is_it_sub_account_using_name(PRIMARY_ACCOUNTS_NAMES.SALES, entry.ACCOUNT) && is_credit(entry.ACCOUNT):
// 			key = entry.ACCOUNT
// 		default:
// 			continue
// 		}
// 		sums := m[key]
// 		if sums == nil {
// 			sums = &INVOICE_STRUCT{}
// 			m[key] = sums
// 		}
// 		sums.VALUE += entry.VALUE
// 		sums.QUANTITY += entry.QUANTITY
// 	}
// 	invoice := []INVOICE_STRUCT{}
// 	for k, v := range m {
// 		invoice = append(invoice, INVOICE_STRUCT{k, v.VALUE, v.VALUE / v.QUANTITY, v.QUANTITY})
// 	}
// 	return invoice
// }
