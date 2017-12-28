// Copyright 2017 Len Budney. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
	Package App provides the core app functionality. This file, flags.go,
	provides support for parsing command-line flags and reading a JSON file
	of runtime options, and merging the two so that options override the
	config file which overrides the defailts.
*/
package app

import (
	"flag"
	"os"
	"testing"
)

// Test the arrayFlags struct
func TestArrayFlags(t *testing.T) {
	var f arrayFlags

	// Empty array stringifies as empty string
	t.Run("Empty", func(t *testing.T) {
		if f.String() != "" {
			t.Fail()
		}
	})

	// Add three items and check the stringification each time
	t.Run("N=1", func(t *testing.T) {
		if f.Set("A"); f.String() != "A" {
			t.Fail()
		}
	})
	t.Run("N=2", func(t *testing.T) {
		if f.Set("B"); f.String() != "A, B" {
			t.Fail()
		}
	})
	t.Run("N=3", func(t *testing.T) {
		if f.Set("C"); f.String() != "A, B, C" {
			t.Fail()
		}
	})

	// Final count should equal 3
	t.Run("Count", func(t *testing.T) {
		if len(f) != 3 {
			t.Fail()
		}
	})
}

// Now test them in-situ in a Flags struct
func TestArrayFlagsInPlace(t *testing.T) {
	var flags Flags
	f := &flags.Bank.Accounts

	// Empty array stringifies as empty string
	t.Run("Empty", func(t *testing.T) {
		if f.String() != "" {
			t.Fail()
		}
	})

	// Add three items and check the stringification each time
	t.Run("N=1", func(t *testing.T) {
		if f.Set("A"); f.String() != "A" {
			t.Fail()
		}
	})
	t.Run("N=2", func(t *testing.T) {
		if f.Set("B"); f.String() != "A, B" {
			t.Fail()
		}
	})
	t.Run("N=3", func(t *testing.T) {
		if f.Set("C"); f.String() != "A, B, C" {
			t.Fail()
		}
	})

	// Final count should equal 3
	t.Run("Count", func(t *testing.T) {
		if len(*f) != 3 {
			t.Fail()
		}
	})
}

// Test some command-line args
func TestParseDefaults(t *testing.T) {
	var defaults Flags
	defaults.Sheets.ConfigFileName = "flags_test_01.json"
	defaults.Sheets.AppSecretFile = defaultPath(defaultSecretFile)
	defaults.Sheets.UserAuthFile = defaultPath(defaultAuthFile)

	os.Args = []string{os.Args[0], "--config-file", "flags_test_01.json"}
	flags, _ := ParseFlags()

	if !isSame(defaults, flags) {
		t.Fail()
	}
}

func TestParseAccounts(t *testing.T) {
	// Reset the command line to parse again
	os.Args = []string{os.Args[0], "--config-file", "flags_test_01.json", "--account", "Account 1", "--account", "Account 2"}
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	flags, _ := ParseFlags()
	accounts := &flags.Bank.Accounts

	if len(*accounts) != 2 {
		t.Errorf("Expected 2 accounts, found %d", len(*accounts))
	} else if flags.Bank.Accounts[0] != "Account 1" {
		t.Fail()
	} else if flags.Bank.Accounts[1] != "Account 2" {
		t.Fail()
	}
}

// Utility to compare two configs. Quick 'n dirty: only checks
// that the account lists and security questions are the same length
func isSame(a, b Flags) bool {
	if &a == &b {
		return true
	}
	if a.Sheets.IndexSheetId != b.Sheets.IndexSheetId {
		return false
	}
	if a.Sheets.ConfigFileName != b.Sheets.ConfigFileName {
		return false
	}
	if a.Sheets.AppSecretFile != b.Sheets.AppSecretFile {
		return false
	}
	if a.Sheets.UserAuthFile != b.Sheets.UserAuthFile {
		return false
	}
	if a.Bank.LoginUrl != b.Bank.LoginUrl {
		return false
	}
	if a.Bank.Username != b.Bank.Username {
		return false
	}
	if a.Bank.Password != b.Bank.Password {
		return false
	}
	if len(a.Bank.Accounts) != len(b.Bank.Accounts) {
		return false
	}
	if len(a.Bank.SecurityQuestions) != len(b.Bank.SecurityQuestions) {
		return false
	}

	return true
}
