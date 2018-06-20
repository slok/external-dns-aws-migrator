package adopt

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/route53iface"

	"github.com/slok/external-dns-aws-migrator/pkg/log"
	"github.com/slok/external-dns-aws-migrator/pkg/model"
)

// RSAdopter is the Route53 AWS record set adopter, it will get a txt Entry and it will addopt the entry in the required route53 hosted zone.
type RSAdopter interface {
	Adopt(*model.Entry) error
}

type adopter struct {
	r53Svc route53iface.Route53API
	dryRun bool
	logger log.Logger
}

// NewRSAdopter is the implementation of the RSAdopter
func NewRSAdopter(dryRun bool, r53Svc route53iface.Route53API, logger log.Logger) RSAdopter {
	return &adopter{
		r53Svc: r53Svc,
		dryRun: dryRun,
		logger: logger,
	}
}

func (a *adopter) Adopt(entry *model.Entry) error {
	// Get the right hosted zone.
	hzid, err := a.findHostedZone(entry.Host)
	if err != nil {
		return err
	}
	// Can create the txt?
	err = a.canCreateTXTEntry(hzid, entry.Host)
	if err != nil {
		return err
	}

	// Create the txt.
	err = a.createTXTEntry(hzid, entry)
	if err != nil {
		return err
	}
	return nil
}

// findHostedZone will find the correct hosted zone for the adopting host.
func (a *adopter) findHostedZone(domain string) (string, error) {
	// Get all the available HZ.
	// TODO: If paginated then get all (more than 100 HZ).
	req := a.r53Svc.ListHostedZonesRequest(&route53.ListHostedZonesInput{})
	res, err := req.Send()
	if err != nil {
		return "", err
	}
	zones := map[string]route53.HostedZone{}
	for _, zone := range res.HostedZones {
		zone := zone
		name := strings.TrimSuffix(*zone.Name, ".")
		zones[name] = zone
	}

	// Sanitize domain and get the different subdomain levels.
	domain = strings.TrimSuffix(domain, ".")
	splDomain := strings.Split(domain, ".")[1:] // We get rid of the first one (wildcard or direct one)

	// Get the correct hosted zone. On each iteration it will remove a subdomain level
	// until it finds the zone.
	for i := 0; i < len(splDomain)-1; i++ {
		// Generate domain.
		domain := strings.Join(splDomain[i:], ".")

		// If HZ found then finish.
		if zone, ok := zones[domain]; ok {
			return *zone.Id, nil
		}
	}

	return "", fmt.Errorf("no hosted zones available for domain %s", domain)
}

func (a *adopter) canCreateTXTEntry(hzID, domain string) error {
	rrs, err := a.getRecordSets(hzID, domain)
	if err != nil {
		return err
	}

	// Check the host exists.
	ts := []route53.RRType{
		route53.RRTypeA,
		route53.RRTypeAaaa,
		route53.RRTypeCname,
	}
	found := a.findRecordSetType(ts, rrs)
	if !found {
		return fmt.Errorf("not present record set for A, AAAA or CNAME types with host %s", domain)
	}

	// Check the host txt exists.
	ts = []route53.RRType{route53.RRTypeTxt}
	found = a.findRecordSetType(ts, rrs)
	if found {
		return fmt.Errorf("txt record set already present for domain: %s", domain)
	}

	return nil
}

func (a *adopter) findRecordSetType(types []route53.RRType, rrs []route53.ResourceRecordSet) bool {
	for _, rr := range rrs {
		for _, t := range types {
			if rr.Type == t {
				return true
			}
		}
	}

	return false
}

func (a *adopter) getRecordSets(hzID, domain string) ([]route53.ResourceRecordSet, error) {
	rrs := []route53.ResourceRecordSet{}
	domain = strings.TrimRight(domain, ".") + "." // Set always the dot at the end.

	params := &route53.ListResourceRecordSetsInput{
		HostedZoneId: aws.String(hzID),
	}
	for {
		req := a.r53Svc.ListResourceRecordSetsRequest(params)
		resp, err := req.Send()
		if err != nil {
			return nil, err
		}

		// Save current records.
		for _, rs := range resp.ResourceRecordSets {
			if aws.StringValue(rs.Name) == domain {
				rrs = append(rrs, rs)
			}
		}

		// No more? then exit loop.
		if !aws.BoolValue(resp.IsTruncated) {
			break
		}

		// prepare the call to grab the next ones.
		params.StartRecordName = resp.NextRecordName
	}

	return rrs, nil
}

func (a *adopter) createTXTEntry(hzID string, entry *model.Entry) error {
	// Ensure string is between quotes.
	txt := fmt.Sprintf(`"%s"`, strings.Trim(entry.TXT, `"`))

	logger := a.logger.With("hz", hzID).
		With("host", entry.Host).
		With("txt", entry.TXT)
	if a.dryRun {
		logger.Infof("not creating txt record set because of dry-run")
		return nil
	}

	input := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &route53.ChangeBatch{
			Changes: []route53.Change{
				{
					Action: route53.ChangeActionCreate,
					ResourceRecordSet: &route53.ResourceRecordSet{
						Name: aws.String(entry.Host),
						Type: route53.RRTypeTxt,
						TTL:  aws.Int64(300),
						ResourceRecords: []route53.ResourceRecord{
							route53.ResourceRecord{
								Value: aws.String(txt),
							},
						},
					},
				},
			},
			Comment: aws.String("Add txt entry"),
		},
		HostedZoneId: aws.String(hzID),
	}

	req := a.r53Svc.ChangeResourceRecordSetsRequest(input)
	_, err := req.Send()
	if err != nil {
		return err
	}

	logger.Infof("txt record set created")
	return nil
}
