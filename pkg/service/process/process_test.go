package process_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/slok/external-dns-aws-migrator/pkg/log"
	madopt "github.com/slok/external-dns-aws-migrator/pkg/mocks/service/adopt"
	mfilter "github.com/slok/external-dns-aws-migrator/pkg/mocks/service/filter"
	"github.com/slok/external-dns-aws-migrator/pkg/service/process"
)

func TestAdoptStream(t *testing.T) {
	tests := []struct {
		name         string
		entries      string
		expTimeCalls int
	}{
		{
			name: "multiple entries should adopt multiple times",
			entries: `
batman.dc.comic.io
superman.dc.comic.io
deadpool.marvel.comic.io
spiderman.marvel.comic.io
wolverine.marvel.comic.io
`,
			expTimeCalls: 5,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)

			// Mocks
			mf := &mfilter.EntryValidator{}
			ma := &madopt.RSAdopter{}

			mf.On("Validate", mock.Anything).Times(test.expTimeCalls).Return(nil, nil)
			ma.On("Adopt", mock.Anything).Times(test.expTimeCalls).Return(nil)

			sa := process.NewStreamAdopter(ma, mf, log.Dummy)
			bs := bytes.NewBufferString(test.entries)
			err := sa.AdoptStream(bs)
			if assert.NoError(err) {
				mf.AssertExpectations(t)
				ma.AssertExpectations(t)
			}
		})
	}
}
