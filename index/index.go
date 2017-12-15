// Copyright 2017 Len Budney. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
	Package Index provides constants and functions for reading
	a spreadsheet that lists other spreadsheets: each one a budget
	covering a particular date range. The app uses these functions
	to look up the budget spreadsheets and determine which one(s)
	a transaction should be added to.
*/
package index

import (
	"fmt"
	"github.com/araddon/dateparse"
	"google.golang.org/api/sheets/v4"
	"log"
	"time"
)

// Where in the spreadsheet are index records found?
const Range = "Index!A2:E"

// Struct for holding an index entry, representing a budget spreadsheet row
// identifying the file, its start/end dates, and the last updated date/time.

type Record struct {
	Index         int
	Filename      string
	Start         time.Time
	End           time.Time
	LastUpdated   time.Time
	SpreadsheetId string
	IndexId       string
}

func getDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
}

// An active record overlaps the specified date range, AND was not
// updated after the end date
func getActiveRecordTester(start time.Time, end time.Time) func(Record) bool {
	start = getDate(start)
	end = getDate(end).Add(24 * time.Hour)

	return func(record Record) bool {
		a := getDate(record.Start)
		b := getDate(record.End).Add(24 * time.Hour)

		// This record ends before the time interval starts
		if start.After(b) {
			return false
		}

		// This record starts after the time interval ends
		if end.Before(a) {
			return false
		}

		// This record has never been updated
		if record.LastUpdated.IsZero() {
			return true
		}

		// This record was last updated before its period ended
		x := getDate(record.LastUpdated)
		if x.Before(b) {
			return false
		}

		return true
	}
}

func Filter(history []Record, test func(Record) bool) (ret []Record) {
	for _, record := range history {
		if test(record) {
			ret = append(ret, record)
		}
	}

	return ret
}

func FilterActiveRecords(history []Record, start time.Time, end time.Time) []Record {
	test := getActiveRecordTester(start, end)
	return Filter(history, test)
}

func FromSpreadsheet(srv *sheets.Service, spreadsheetId string) ([]Record, error) {
	var history []Record

	// Open the spreadsheet
	response, err := srv.Spreadsheets.Values.Get(spreadsheetId, Range).Do()
	if err != nil {
		log.Printf("Unable to retrieve index from sheet ID %s: %v", spreadsheetId, err)
		return history, err
	}

	// It's technically OK for there to be no index data, but we go
	// ahead and log it
	if len(response.Values) == 0 {
		log.Printf("No index data found in sheet ID %s", spreadsheetId)
		return history, nil
	}

	// OK, parse it
	for i, row := range response.Values {
		record, err := FromSpreadsheetRow(i+1, row)
		if err != nil {
			return history, err
		}

		record.IndexId = spreadsheetId

		history = append(history, record)
	}

	return history, nil
}

func FromSpreadsheetRow(index int, row []interface{}) (Record, error) {
	record := Record{}
	var err error

	// Strings are easy
	record.Index = index
	record.Filename = fmt.Sprintf("%s", row[0])
	record.SpreadsheetId = fmt.Sprintf("%s", row[4])

	// Dates need parsed, and error-checked
	record.Start, err = dateparse.ParseLocal(fmt.Sprintf("%s", row[1]))
	if err != nil {
		log.Printf("Failed to parse date \"%s\": %v", row[1], err)
		return record, err
	}
	record.End, err = dateparse.ParseLocal(fmt.Sprintf("%s", row[2]))
	if err != nil {
		log.Printf("Failed to parse date \"%s\": %v", row[2], err)
		return record, err
	}

	// LastUpdated is optional
	if fmt.Sprintf("%s", row[3]) != "" {
		record.LastUpdated, err = dateparse.ParseLocal(fmt.Sprintf("%s", row[3]))
		if err != nil {
			log.Printf("Failed to parse date \"%s\": %v", row[3], err)
			return record, err
		}
	}

	return record, nil
}
