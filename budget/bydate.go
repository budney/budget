// Copyright 2017 Len Budney. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package budget defines constants and provides functions for
// working with a budget spreadsheet. This is where knowledge
// is found about where transactions go in the sheet, as well
// as any other information needed to update the sheet with
// data downloaded from the bank.
package budget

// A byDate is an array of Transaction structs, which implements
// sort.Interface for sorting transactions by date and index.
type byDate []Transaction

// Len returns the length of a byDate array of transactions
func (t byDate) Len() int { return len(t) }

// Swap interchanges two elements of a byDate array of transactions
func (t byDate) Swap(i, j int) { t[i], t[j] = t[j], t[i] }

// Less compares two transactions
func (t byDate) Less(i, j int) bool {
	if t[i].Date.Before(t[j].Date) {
		return true
	}
	if t[i].Date.After(t[j].Date) {
		return false
	}
	if t[i].Index < t[j].Index {
		return true
	}

	return false
}
