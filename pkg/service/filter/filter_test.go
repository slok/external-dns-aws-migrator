package filter_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/slok/external-dns-aws-adopter/pkg/model"
	"github.com/slok/external-dns-aws-adopter/pkg/service/filter"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name     string
		filter   string
		txt      string
		host     string
		expEntry *model.Entry
		expErr   bool
	}{
		{
			name:   "A valid hsot with a matching regex should return that is valid",
			filter: `.*batman\.com$`,
			txt:    "test-owner-id",
			host:   "bruce-wayne.is.batman.com",
			expEntry: &model.Entry{
				Host: "bruce-wayne.is.batman.com",
				TXT:  "heritage=external-dns,external-dns/owner=test-owner-id",
			},
		},
		{
			name:   "A invalid hsot with a matching regex should return that is invalid",
			filter: `.*batman\.com$`,
			txt:    "test-owner-id",
			host:   "bruce-wayne.is.spiderman.com",
			expErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require := require.New(t)
			assert := assert.New(t)

			ev, err := filter.NewEntryValidator(test.filter, test.txt)
			require.NoError(err)
			gotEntry, err := ev.Validate(test.host)

			if test.expErr {
				assert.Error(err)
			} else if assert.NoError(err) {
				assert.Equal(test.expEntry, gotEntry)
			}
		})
	}
}
