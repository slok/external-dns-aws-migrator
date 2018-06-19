package mocks

import (
	"github.com/aws/aws-sdk-go-v2/service/route53/route53iface"
)

// These are simple wrappers around third party libs so mockery tool can create the mocks.
// https://github.com/vektra/mockery/issues/181

// Route53API is a route53iface.Route53API wrapper.
type Route53API interface{ route53iface.Route53API }
