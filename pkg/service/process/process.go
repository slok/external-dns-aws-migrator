package process

import (
	"bufio"
	"io"

	"github.com/slok/external-dns-aws-adopter/pkg/log"
	"github.com/slok/external-dns-aws-adopter/pkg/service/adopt"
	"github.com/slok/external-dns-aws-adopter/pkg/service/filter"
)

// StreamAdopter knows how to process a stream.
type StreamAdopter interface {
	AdoptStream(io.Reader) error
}

type streamAdopter struct {
	adSvc  adopt.RSAdopter
	flSvc  filter.EntryValidator
	logger log.Logger
}

// NewStreamAdopter returns a new stream adopter.
func NewStreamAdopter(adSvc adopt.RSAdopter, flSvc filter.EntryValidator, logger log.Logger) StreamAdopter {
	return &streamAdopter{
		adSvc:  adSvc,
		flSvc:  flSvc,
		logger: logger,
	}
}

func (s *streamAdopter) AdoptStream(r io.Reader) error {
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		domain := sc.Text()
		if domain != "" {
			err := s.adoptEntry(domain)
			if err != nil {
				s.logger.Warningf("error adopting entry: %s", err)
			}
		}
		if err := sc.Err(); err != nil {
			return nil
		}
	}
	return nil
}

func (s *streamAdopter) adoptEntry(domain string) error {
	entry, err := s.flSvc.Validate(domain)
	if err != nil {
		s.logger.Debugf("ignoring domain %s", domain)
		return nil
	}

	err = s.adSvc.Adopt(entry)
	if err != nil {
		return err
	}
	return nil
}
