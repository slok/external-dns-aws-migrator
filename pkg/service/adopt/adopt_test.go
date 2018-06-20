package adopt_test

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/slok/external-dns-aws-migrator/pkg/log"
	mroute53iface "github.com/slok/external-dns-aws-migrator/pkg/mocks/github.com/aws/aws-sdk-go-v2/service/route53/route53iface"
	"github.com/slok/external-dns-aws-migrator/pkg/model"
	"github.com/slok/external-dns-aws-migrator/pkg/service/adopt"
)

func mockDefaultListHostedZones() route53.ListHostedZonesRequest {
	// precreated HostedZones.
	var (
		hz0 = route53.HostedZone{
			Name: aws.String("peter.parker.spiderman.marvel.superheroes.comics."),
			Id:   aws.String("peter.parker.spiderman.marvel.superheroes.comics."),
		}
		hz1 = route53.HostedZone{
			Name: aws.String("batman.dc.superheroes.comics."),
			Id:   aws.String("batman.dc.superheroes.comics."),
		}
		hz2 = route53.HostedZone{
			Name: aws.String("dc.superheroes.comics."),
			Id:   aws.String("dc.superheroes.comics."),
		}
	)

	return mockListHostedZones(&route53.ListHostedZonesOutput{
		HostedZones: []route53.HostedZone{hz0, hz1, hz2},
	})
}

func mockListHostedZones(v *route53.ListHostedZonesOutput) route53.ListHostedZonesRequest {
	return route53.ListHostedZonesRequest{
		Request: &aws.Request{
			Data: v,
		},
	}
}

func mockDefaultListResourceRecordSetsRequest() route53.ListResourceRecordSetsRequest {
	var (
		rrs0 = route53.ResourceRecordSet{
			Name: aws.String("valid.with.txt.batman.dc.superheroes.comics."),
			Type: route53.RRTypeA,
		}
		rrs1 = route53.ResourceRecordSet{
			Name: aws.String("valid.with.txt.batman.dc.superheroes.comics."),
			Type: route53.RRTypeTxt,
		}
		rrs2 = route53.ResourceRecordSet{
			Name: aws.String("valid.without.txt.batman.dc.superheroes.comics."),
			Type: route53.RRTypeA,
		}
	)
	return mockListResourceRecordSetsRequest(&route53.ListResourceRecordSetsOutput{
		ResourceRecordSets: []route53.ResourceRecordSet{rrs0, rrs1, rrs2},
	})
}

func mockListResourceRecordSetsRequest(v *route53.ListResourceRecordSetsOutput) route53.ListResourceRecordSetsRequest {
	return route53.ListResourceRecordSetsRequest{
		Request: &aws.Request{
			Data: v,
		},
	}
}

func mockChangeResourceRecordSetsRequest(v *route53.ChangeResourceRecordSetsOutput) route53.ChangeResourceRecordSetsRequest {
	return route53.ChangeResourceRecordSetsRequest{
		Request: &aws.Request{
			Data: v,
		},
	}
}

func getTXTResourceRecordSetMatchedByFunc(expTXT, expHZID, expHost string) func(*route53.ChangeResourceRecordSetsInput) bool {
	return func(input *route53.ChangeResourceRecordSetsInput) bool {
		if aws.StringValue(input.HostedZoneId) != expHZID {
			return false
		}

		// Only one change.
		if len(input.ChangeBatch.Changes) > 1 {
			return false
		}

		ch := input.ChangeBatch.Changes[0]

		// Check types
		if ch.Action != route53.ChangeActionCreate {
			return false
		}
		if ch.ResourceRecordSet.Type != route53.RRTypeTxt {
			return false
		}

		// Check data.
		if aws.StringValue(ch.ResourceRecordSet.Name) != expHost {
			return false
		}
		// For now only allowed entry creation, so this should be one value.
		if aws.StringValue(ch.ResourceRecordSet.ResourceRecords[0].Value) != expTXT {
			return false
		}

		return true
	}
}

func TestAdopterAdopt(t *testing.T) {
	tests := []struct {
		name         string
		dryRun       bool
		entry        *model.Entry
		expEntryHZ   string
		expEntryTXT  string
		expEntryHost string
		expErr       bool
	}{
		{
			name:   "If there is not valid HZ for the domain it should error.",
			dryRun: false,
			entry: &model.Entry{
				Host: "domain.with.no.hosted-zone.com",
				TXT:  "heritage=external-dns,external-dns/owner=default",
			},
			expErr: true,
		},
		{
			name:   "If there is not a A, AAAA or CNAME entry with the host on the HZ it should fail.",
			dryRun: false,
			entry: &model.Entry{
				Host: "no.a.aaaa.or.cname.entry.batman.dc.superheroes.comics",
				TXT:  "heritage=external-dns,external-dns/owner=default",
			},
			expErr: true,
		},
		{
			name:   "If there is a A, AAAA or CNAME already with the host and also a TXT it should fail.",
			dryRun: false,
			entry: &model.Entry{
				Host: "valid.with.txt.batman.dc.superheroes.comics",
				TXT:  "heritage=external-dns,external-dns/owner=default",
			},
			expErr: true,
		},
		{
			name:   "If there is a A, AAAA or CNAME already with the host and not a TXT it should create the entry.",
			dryRun: false,
			entry: &model.Entry{
				Host: "valid.without.txt.batman.dc.superheroes.comics",
				TXT:  "heritage=external-dns,external-dns/owner=default",
			},
			expEntryHZ:   "batman.dc.superheroes.comics.",
			expEntryTXT:  `"heritage=external-dns,external-dns/owner=default"`,
			expEntryHost: "valid.without.txt.batman.dc.superheroes.comics",
			expErr:       false,
		},
		{
			name:   "If there is a A, AAAA or CNAME already with the host and not a TXT in dry run mode it shouldn't create the entry.",
			dryRun: true,
			entry: &model.Entry{
				Host: "valid.without.txt.batman.dc.superheroes.comics",
				TXT:  "heritage=external-dns,external-dns/owner=default",
			},
			expEntryHZ:   "batman.dc.superheroes.comics.",
			expEntryTXT:  `"heritage=external-dns,external-dns/owner=default"`,
			expEntryHost: "valid.without.txt.batman.dc.superheroes.comics",
			expErr:       false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)

			// Mocks.
			mr53 := &mroute53iface.Route53API{}

			// Mock hosted zones.
			mr53.On("ListHostedZonesRequest", mock.Anything).Return(mockDefaultListHostedZones())
			mr53.On("ListResourceRecordSetsRequest", mock.Anything).Return(mockDefaultListResourceRecordSetsRequest())
			// Only for no dry run.
			if !test.dryRun {
				mbf := getTXTResourceRecordSetMatchedByFunc(test.expEntryTXT, test.expEntryHZ, test.expEntryHost)
				mr53.On("ChangeResourceRecordSetsRequest", mock.MatchedBy(mbf)).Return(mockChangeResourceRecordSetsRequest(nil))
			}

			ad := adopt.NewRSAdopter(test.dryRun, mr53, log.Dummy)

			err := ad.Adopt(test.entry)
			if test.expErr {
				assert.Error(err)
			} else if assert.NoError(err) {
				// In the happy path assert the expectations.
				mr53.AssertExpectations(t)
			}
		})
	}

}
