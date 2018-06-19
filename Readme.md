# external-dns-aws-adopter

Utility to let [external-dns][external-dns] adopt route53 record sets not managed by [external-dns][external-dns]

## Motivation

Sometimes you add external-dns to your clusters and the new flow doesn't have your old host entries from route53 but you want these entries being tracked and managed by external-dns, to do this you would need to create a txt entry for each one, this tool helps doing this migration all at once, filtered by hosts, with multiple external-dns...

[external-dns]: https://github.com/kubernetes-incubator/external-dns