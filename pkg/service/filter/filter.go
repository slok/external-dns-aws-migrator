package filter

import (
	"fmt"
	"regexp"

	"github.com/slok/external-dns-aws-migrator/pkg/model"
)

const (
	txtEntryFMT = "heritage=external-dns,external-dns/owner=%s"
)

// EntryValidator will validate an entry.
type EntryValidator interface {
	Validate(host string) (*model.Entry, error)
}

type validator struct {
	filter *regexp.Regexp
	txt    string
}

// NewEntryValidator returns a new entry validator.
func NewEntryValidator(filter, ownerID string) (EntryValidator, error) {
	r, err := regexp.Compile(filter)
	if err != nil {
		return nil, err
	}
	return &validator{
		filter: r,
		txt:    fmt.Sprintf(txtEntryFMT, ownerID),
	}, nil
}

func (v *validator) Validate(host string) (*model.Entry, error) {
	// Check the regexp filters
	if !v.filter.MatchString(host) {
		return nil, fmt.Errorf("%s not a valid host for the loaded filter", host)
	}

	return &model.Entry{
		Host: host,
		TXT:  v.txt,
	}, nil
}
