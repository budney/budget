// Copyright 2017 Len Budney. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
	Command budget-update fetches transactions from one or
	more sources, and forwards them one-by-one to one or
	more destinations. Usually this means fetching from
	a bank web site and writing to a Google spreadsheet.
*/
package main

import (
	"github.com/budney/budget/app"
	"github.com/budney/budget/budget"
	"github.com/budney/budget/index"
	"github.com/budney/google/sheets"
	"github.com/budney/tdbank"

	"log"
	"sync"
	"time"
)

func main() {
	flags := app.ParseFlags()

	index := getBudgetIndex(flags)
	srv, _ := sheets.GetService(flags.Sheets.AppSecretFile, flags.Sheets.UserAuthFile)
	spreadsheet := budget.Spreadsheet{index[len(index)-1], *srv}

	wait := new(sync.WaitGroup)
	wait.Add(1)

	channel := make(chan budget.Transaction)
	spreadsheet.AppendFromChannel(channel, wait, "Joint Checking", "Uncategorized")

	transactions := getTransactions(flags)
	for _, v := range transactions {
		channel <- budget.Transaction(v)
	}

	close(channel)
	wait.Wait()
}

func getTransactions(flags app.Flags) []tdbank.HistoryRecord {
	var client tdbank.Client
	client.Start()
	defer client.Stop()

	auth := tdbank.Auth{
		LoginUrl:          flags.Bank.LoginUrl,
		Username:          flags.Bank.Username,
		Password:          flags.Bank.Password,
		SecurityQuestions: flags.Bank.SecurityQuestions,
	}
	client.Login(auth)

	start := time.Date(2017, time.December, 31, 0, 0, 0, 0, time.Local)
	end := time.Date(2018, time.January, 10, 0, 0, 0, 0, time.Local)
	log.Printf("Date range: %s - %s", start.Format("01/02/2006"), end.Format("01/02/2006"))

	client.DownloadAccountHistory("Joint Checking", start, end)
	history, err := client.ParseAccountHistory()
	if err != nil {
		log.Fatalf("Failed to read history: %s", err)
	}

	return history
}

func getBudgetIndex(flags app.Flags) []index.Record {
	srv, err := sheets.GetService(flags.Sheets.AppSecretFile, flags.Sheets.UserAuthFile)
	if err != nil {
		log.Fatalf("Couldn't initialize sheets service: %s", err)
	}

	index, err := index.FromGoogleSheet(srv, flags.Sheets.IndexSheetId)
	if err != nil {
		log.Fatalf("Couldn't read budget index: %s", err)
	}

	/*
		values := [][]interface{}{{"A", "B", "C", "3.14", "E"}}
		valuerange := &google.ValueRange{Range: "A1:E", MajorDimension: "ROWS", Values: values}
		_, err = srv.Spreadsheets.Values.Append(flags.Sheets.IndexSheetId, "A1:E", valuerange).ValueInputOption("USER_ENTERED").Do()
		if err != nil {
			log.Fatalf("Couldn't append stuff: %s", err)
		}
	*/

	return index
}
