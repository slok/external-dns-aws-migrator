// Code generated by mockery v1.0.0
package adopt

import mock "github.com/stretchr/testify/mock"
import model "github.com/slok/external-dns-aws-adopter/pkg/model"

// EntryValidator is an autogenerated mock type for the EntryValidator type
type EntryValidator struct {
	mock.Mock
}

// Validate provides a mock function with given fields: host
func (_m *EntryValidator) Validate(host string) (*model.Entry, error) {
	ret := _m.Called(host)

	var r0 *model.Entry
	if rf, ok := ret.Get(0).(func(string) *model.Entry); ok {
		r0 = rf(host)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Entry)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(host)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
