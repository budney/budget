package app

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os/user"
	"path/filepath"
)

const defaultConfigDir = ".budget-update"    // Collect configs in one place
const defaultConfigFile = "options.json"     // This is just the command options
const defaultSecretFile = "client-auth.json" // This is just the command options
const defaultAuthFile = "user-auth.json"
const nullString = string(byte(0)) // Hard (but not impossible) to supply on command line

// Struct for holding command-line flags related to sheets
type Sheets struct {
	IndexSheetId   string
	ConfigFileName string
	AppSecretFile  string
	UserAuthFile   string
}

// Struct for holding command-line flags related to web banking
type Bank struct {
	LoginUrl          string
	Username          string
	Password          string
	SecurityQuestions map[string]string
}

// Struct for holding all the command-line flags
type Flags struct {
	Sheets Sheets
	Bank   Bank
}

func ParseFlags() (Flags, error) {
	var flags Flags

	flag.StringVar(&flags.Sheets.IndexSheetId, "index-sheet-id", nullString, "Google drive `sheet-id` of the budget index")
	flag.StringVar(&flags.Sheets.ConfigFileName, "config-file", nullString, "The `filename` of the config file to read at startup")
	flag.StringVar(&flags.Sheets.AppSecretFile, "app-secret-file", nullString, "The `filename` for the app to authenticate with Google Drive")
	flag.StringVar(&flags.Sheets.UserAuthFile, "user-auth-file", nullString, "The `filename` with cached user credentials for Google Drive")
	flag.StringVar(&flags.Bank.LoginUrl, "bank-url", nullString, "The `URL` of the online banking web page")
	flag.StringVar(&flags.Bank.Username, "bank-username", nullString, "Your online banking `username`")
	flag.StringVar(&flags.Bank.Password, "bank-password", nullString, "Your online banking `password`")

	// Need user info from the OS
	usr, err := user.Current()
	if err != nil {
		log.Fatalf("Can't get current user: %v", err)
	}

	// Parse the command line
	flag.Parse()

	var configFile string

	// Read the defaults file, if any
	if flags.Sheets.ConfigFileName != "" && flags.Sheets.ConfigFileName != nullString {
		configFile = flags.Sheets.ConfigFileName
	} else {
		configFile = filepath.Join(usr.HomeDir, defaultConfigDir, defaultConfigFile)
	}

	// Read the configs from a file, and then overwrite with options
	// that were set on the command line
	options, _ := flagsFromFile(configFile)

	// Copy any options set on the command line
	if flags.Sheets.IndexSheetId != nullString {
		options.Sheets.IndexSheetId = flags.Sheets.IndexSheetId
	}
	if flags.Sheets.ConfigFileName != nullString {
		options.Sheets.ConfigFileName = flags.Sheets.ConfigFileName
	}
	if flags.Sheets.AppSecretFile != nullString {
		options.Sheets.AppSecretFile = flags.Sheets.AppSecretFile
	}

	if flags.Bank.LoginUrl != nullString {
		options.Bank.LoginUrl = flags.Bank.LoginUrl
	}
	if flags.Bank.Username != nullString {
		options.Bank.Username = flags.Bank.Username
	}
	if flags.Bank.Password != nullString {
		options.Bank.Password = flags.Bank.Password
	}

	// Very last, set default values, if they weren't already set
	if options.Sheets.AppSecretFile == "" {
		options.Sheets.AppSecretFile = filepath.Join(usr.HomeDir, defaultConfigDir, defaultSecretFile)
	}
	if options.Sheets.UserAuthFile == "" {
		options.Sheets.UserAuthFile = filepath.Join(usr.HomeDir, defaultConfigDir, defaultAuthFile)
	}

	return options, nil
}

func flagsFromFile(fileName string) (Flags, error) {
	flags := Flags{}

	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Printf("Unable to read config file: %v", err)
		return flags, err
	}

	err = json.Unmarshal(b, &flags)
	if err != nil {
		log.Printf("Unable to process config file: %v", err)
		return flags, err
	}

	return flags, nil
}
