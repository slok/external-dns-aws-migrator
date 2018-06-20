package main

import (
	"flag"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws/endpoints"
)

// Defaults.
const (
	defTXTOwnerID  = "default"
	defFilter      = `^.+$`
	defAWSRegion   = endpoints.EuWest1RegionID
	defDryRun      = false
	defDebug       = false
	defShowVersion = false
)

// Flags are the flags of the program.
type Flags struct {
	AWSRegion   string
	Filter      string
	TXTOwnerID  string
	DryRun      bool
	Debug       bool
	ShowVersion bool
}

// NewFlags returns the flags of the commandline.
func NewFlags() *Flags {
	flags := &Flags{}
	fl := flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	fl.StringVar(&flags.AWSRegion, "aws-region", defAWSRegion, "AWS region to act on hosted zones")
	fl.StringVar(&flags.Filter, "filter", defFilter, "regex to filter domains to act on")
	fl.StringVar(&flags.TXTOwnerID, "txt-owner-id", defTXTOwnerID, "the txt owner id that will be set on the txt registry")
	fl.BoolVar(&flags.DryRun, "dry-run", defDryRun, "run in dry-run mode")
	fl.BoolVar(&flags.Debug, "debug", defDebug, "run in debug mode")
	fl.BoolVar(&flags.ShowVersion, "version", defShowVersion, "show version of the app")

	fl.Parse(os.Args[1:])

	return flags
}
