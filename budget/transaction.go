// Copyright 2017 Len Budney. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package budget

import (
    "time"
)

// A Transaction contains information about a single transaction.
type Transaction struct {
	Index          int       // A counter for sorting transactions on the same Date
	Date           time.Time // The date of the transaction
	Type           string    // A type description, such as POS, Check, ATM, etc.
	Description    string    // Usually the payor / payee of the transaction
	DebitPennies   int64     // The debit amount, in pennies
	CreditPennies  int64     // The credit amount, in pennies
	BalancePennies int64     // The balance, in pennies, after the transaction
}
