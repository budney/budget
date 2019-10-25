// Copyright 2017 Len Budney. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package budget defines constants and provides functions for
// working with a budget spreadsheet. This is where knowledge
// is found about where transactions go in the sheet, as well
// as any other information needed to update the sheet with
// data downloaded from the bank.
package budget

import (
	"github.com/budney/budget/index"
	"google.golang.org/api/sheets/v4"
	"log"
	"sort"
	"sync"
)

// HeaderRange gives the location of the transaction header
const HeaderRange = "A1:H1"
// DataRange gives the location of the transactions
const DataRange = "A2:H"

// Spreadsheet has the same structure as a Record, and holds
// high-level information about a spreadsheet.
type Spreadsheet struct {
	index.Record   // Location, date range covered, etc.
	sheets.Service // Adds the Google Sheets API to this struct
}

// AppendFromChannel runs a goroutine that listens to a channel for
// transactions, filters out the ones that don't apply, and appends
// the rest to the budget spreadsheet for the specified account. It
// does the append when the channel is closed by the writer.
func (spreadsheet *Spreadsheet) AppendFromChannel(input <-chan Transaction, wait *sync.WaitGroup, worksheet string, category string) {
	go func() {
		transactions := make([]Transaction, 0, 2)

		for transaction := range input {
			transactions = append(transactions, transaction)
		}

		spreadsheet.AppendArray(transactions, worksheet, category)
		wait.Done()
	}()
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
	sort.Sort(byDate(transactions))

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
