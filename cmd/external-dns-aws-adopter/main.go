package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/route53iface"

	"github.com/slok/external-dns-aws-adopter/pkg/log"
	"github.com/slok/external-dns-aws-adopter/pkg/service/adopt"
	"github.com/slok/external-dns-aws-adopter/pkg/service/filter"
	"github.com/slok/external-dns-aws-adopter/pkg/service/process"
)

// Main is the Main program.
type Main struct {
	flags  *Flags
	logger log.Logger
}

// Main is the main function that will be executed.
func (m *Main) Main() error {
	r53cli := m.createAWSCli(defAWSRegion)

	if m.flags.Debug {
		m.logger.Set("debug")
	}

	// Create services.
	fsvc, err := filter.NewEntryValidator(m.flags.Filter, m.flags.TXTOwnerID)
	if err != nil {
		return err
	}
	adsvc := adopt.NewRSAdopter(m.flags.DryRun, r53cli, m.logger)
	spsvc := process.NewStreamAdopter(adsvc, fsvc, m.logger)

	// Start adopting.
	err = spsvc.AdoptStream(os.Stdin)
	if err != nil {
		return err
	}

	return nil
}

func (m *Main) createAWSCli(awsRegion string) route53iface.Route53API {
	// and credentials values from the environment variables, shared
	// credentials, and shared configuration files
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		panic("unable to load SDK config, " + err.Error())
	}

	// Set the AWS Region that the service clients should use
	cfg.Region = awsRegion

	return route53.New(cfg)
}

func main() {
	flags := NewFlags()
	m := &Main{
		flags:  flags,
		logger: log.Base(),
	}
	err := m.Main()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing program: %s\n", err)
		os.Exit(1)
	}

	os.Exit(0)
}
