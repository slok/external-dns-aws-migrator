/*
Package mocks will have all the mocks of the library.
*/
package mocks // import "github.com/slok/external-dns-aws-migrator/pkg/mocks"

// AWS mocks.
//go:generate mockery -output ./github.com/aws/aws-sdk-go-v2/service/route53/route53iface -outpkg route53iface -dir ./ -name Route53API

// Service mocks.
//go:generate mockery -output ./service/adopt -outpkg adopt -dir ../service/adopt -name RSAdopter
//go:generate mockery -output ./service/filter -outpkg adopt -dir ../service/filter -name EntryValidator
