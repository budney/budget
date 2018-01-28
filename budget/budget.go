// Copyright 2017 Len Budney. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
	Package Budget defines constants and provides functions for
	working with a budget spreadsheet. This is where knowledge
	is found about where transactions go in the sheet, as well
	as any other information needed to update the sheet with
	data downloaded from the bank.
*/
package budget

import (
	"github.com/budney/budget/index"
	"google.golang.org/api/sheets/v4"
	"log"
	"sort"
	"time"
)

const HeaderRange = "A1:H1" // HeaderRange gives the location of the transaction header
const DataRange = "A2:H"    // DataRange gives the location of the transactions

// A budget.Spreadsheet has the same structure as an index.Record, and holds
// high-level information about a spreadsheet.
type Spreadsheet struct {
	index.Record   // Location, date range covered, etc.
	sheets.Service // Adds the Google Sheets API to this struct
}

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

// A ByDate is an array of Transaction structs, which implements
// sort.Interface to sort transactions by date and index.
type ByDate []Transaction

// Len returns the length of a ByDate array of transactions
func (t ByDate) Len() int      { return len(t) }
func (t ByDate) Swap(i, j int) { t[i], t[j] = t[j], t[i] }
func (t ByDate) Less(i, j int) bool {
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

// AppendArray accepts an array of transaction records and appends them
// to the spreadsheet, sorted by Date and Index. It uses the worksheet
// whose name exactly matches the account, and it puts the provided
// category in the first spreadsheet column.
//
// NOTE! This method appends everything it's given. It doesn't filter
// the records based on date, or anything else. If you call this
// method directly, you should know what you're doing.
func (spreadsheet *Spreadsheet) AppendArray(transactions []Transaction, worksheet string, category string) error {
	// Sort the transactions in place by Date and Index
	sort.Sort(ByDate(transactions))

	// Extract the transaction records in column order
	rows := make([][]interface{}, 2)
	for _, transaction := range transactions {
		rows = append(rows, []interface{}{
			category,
			transaction.Index,
			transaction.Date.Format("1/2/2006"),
			transaction.Type,
			transaction.Description,
			transaction.DebitPennies / 100.0,
			transaction.CreditPennies / 100.0,
			transaction.BalancePennies / 100.0,
		})
	}

	area := worksheet + "!" + DataRange
	valueRange := &sheets.ValueRange{Range: area, MajorDimension: "ROWS", Values: rows}
	_, err := spreadsheet.Spreadsheets.Values.Append(spreadsheet.SpreadsheetId, area, valueRange).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		log.Printf("Couldn't append transactions: %s", err)
		return err
	}

	return nil
}
