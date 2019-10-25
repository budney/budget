// Copyright 2017 Len Budney. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package app provides the core app functionality. This file, flags.go,
// provides support for parsing command-line flags and reading a JSON file
// of runtime options, and merging the two so that options override the
// config file which overrides the defailts.
package app

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os/user"
	"path/filepath"
	"strings"
)

const defaultConfigDir = ".budget-update"    // The defaultConfigDir contains all config files
const defaultConfigFile = "options.json"     // The defaultConfigFile is read unless overridden
const defaultSecretFile = "client-auth.json" // The defaultSecretFile contains app authentication
const defaultAuthFile = "user-auth.json"     // The defaultAuthFile contains user authentication
const nullString = string(byte(0))           // A string with a null byte

// Sheets holds command-line flags related to spreadsheets
type Sheets struct {
	IndexSheetID   string
	ConfigFileName string
	AppSecretFile  string
	UserAuthFile   string
}

// Type for holding repeated arguments
type arrayFlags []string

func (i *arrayFlags) String() string {
	return strings.Join(*i, ", ")
}

func (i *arrayFlags) Set(value string) error {
	log.Printf("Called with value %s", value)
	*i = append(*i, value)
	return nil
}

// Bank holds command-line flags related to web banking
type Bank struct {
	LoginURL          string
	Username          string
	Password          string
	Accounts          arrayFlags
	SecurityQuestions map[string]string
}

// Flags holds all the command-line flags
type Flags struct {
	Sheets Sheets
	Bank   Bank
}

// readCommandLine uses the flag package to configure command-line options.
func readCommandLine(flags Flags) {
	// Configure command-line options
	flag.StringVar(&flags.Sheets.IndexSheetID, "index-sheet-id", nullString, "Google drive `sheet-id` of the budget index")
	flag.StringVar(&flags.Sheets.ConfigFileName, "config-file", nullString, "The `filename` of the config file to read at startup")
	flag.StringVar(&flags.Sheets.AppSecretFile, "app-secret-file", nullString, "The `filename` for the app to authenticate with Google Drive")
	flag.StringVar(&flags.Sheets.UserAuthFile, "user-auth-file", nullString, "The `filename` with cached user credentials for Google Drive")
	flag.StringVar(&flags.Bank.LoginURL, "bank-url", nullString, "The `URL` of the online banking web page")
	flag.StringVar(&flags.Bank.Username, "bank-username", nullString, "Your online banking `username`")
	flag.StringVar(&flags.Bank.Password, "bank-password", nullString, "Your online banking `password`")
	flag.Var(&flags.Bank.Accounts, "account", "Name(s) of account(s) to download transactions for")

	// Parse the command line
	flag.Parse()
}

// readConfigFile reads the config file specified in flags, or else the
// default.
func readConfigFile(flags Flags) Flags {
	var configFile string

	// Read the defaults file, if any
	if flags.Sheets.ConfigFileName != "" && flags.Sheets.ConfigFileName != nullString {
		configFile = flags.Sheets.ConfigFileName
	} else {
		configFile = defaultPath(defaultConfigFile)
	}

	// Read the configs from a file, and then overwrite with options
	// that were set on the command line
	return flagsFromFile(configFile)
}

func copyOptions(src Flags, dest Flags) {
	// Copy spreadsheet options
	if src.Sheets.IndexSheetID != nullString {
		dest.Sheets.IndexSheetID = src.Sheets.IndexSheetID
	}
	if src.Sheets.ConfigFileName != nullString {
		dest.Sheets.ConfigFileName = src.Sheets.ConfigFileName
	}
	if src.Sheets.AppSecretFile != nullString {
		dest.Sheets.AppSecretFile = src.Sheets.AppSecretFile
	}

	// Copy bank options
	if src.Bank.LoginURL != nullString {
		dest.Bank.LoginURL = src.Bank.LoginURL
	}
	if src.Bank.Username != nullString {
		dest.Bank.Username = src.Bank.Username
	}
	if src.Bank.Password != nullString {
		dest.Bank.Password = src.Bank.Password
	}

	// Copy the list of accounts
	if len(src.Bank.Accounts) > 0 {
		dest.Bank.Accounts = src.Bank.Accounts
	}
}

// ParseFlags parses command-line flags
func ParseFlags() Flags {
	// Read the command line flags. We have to do this
	// first, in case they specify a different config file.
	var flags Flags
	readCommandLine(flags)

	// Read the config file flags, and replace them with
	// any command-line flags we received.
	var options = readConfigFile(flags)
	copyOptions(flags, options)

	// Now set default authorization values, if they weren't already set
	if options.Sheets.AppSecretFile == "" {
		options.Sheets.AppSecretFile = defaultPath(defaultSecretFile)
	}
	if options.Sheets.UserAuthFile == "" {
		options.Sheets.UserAuthFile = defaultPath(defaultAuthFile)
	}

	return options
}

// defaultPath takes the specified filename, and converts it into
// an absolute path in the default subdirectory of the user's home
// directory.
func defaultPath(fileName string) string {
	// Get user info from the OS
	usr, err := user.Current()
	if err != nil {
		log.Fatalf("Can't get current user: %v", err)
	}

	// Return the default path to fileName
	return filepath.Join(usr.HomeDir, defaultConfigDir, filepath.Clean(fileName))
}

// flagsFromFile reads the specified file and converts it into a Flags record.
func flagsFromFile(fileName string) Flags {
	flags := Flags{}

	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalf("Unable to read config file: %v", err)
	}

	err = json.Unmarshal(b, &flags)
	if err != nil {
		log.Fatalf("Unable to process config file: %v", err)
	}

	return flags
}
