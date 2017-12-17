// Copyright 2017 Len Budney. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
	Package Worksheet defines constants and provides functions
	for working with a budget spreadsheet. This is where knowledge
	is found about where transactions go in the sheet, as well
	as any other information needed to update the sheet with
	data downloaded from the bank.
*/
package worksheet

import (
	"fmt"
	"github.com/araddon/dateparse"
	"google.golang.org/api/sheets/v4"
	"log"
	"time"
)

const HeaderRange = "B1:H1" // HeaderRange gives the location of the transaction header
const DataRange = "B2:H"    // DataRange gives the location of the transactions
