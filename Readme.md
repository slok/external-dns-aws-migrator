# external-dns-aws-migrator

Utility that helps migrating old (not managed by [external-dns][external-dns]) route53 resource record sets from a hosted zone to being managed by external-dns and let update based on the kubernetes ingress address (set by other ingress controllers)

## Motivation

Sometimes you add external-dns to your clusters and the new flow doesn't have your old host entries from route53 but you want these entries being tracked and managed by external-dns, to do this you would need to create a txt entry for each one, this tool helps doing this migration all at once, filtered by hosts, with multiple external-dns...

## When should I use this?

You are using Kubernetes, AWS, you already have dns entries in route53 (manually or managed by another tool) and you want to start using external-dns to automate ingress loadbalancer addresses in route53.

## Usage

external-dns-aws-migrator reads hosts from the stdin (one per line) and tries to adopt the entries  so the external-dns starts managing the entries.

Example:

Get all ingress hosts from a cluster.

```bash
kubectl get ingress \
    --all-namespaces  \
    -o jsonpath='{.items[*].spec.rules[*].host}' \
    | sed "s/ /\n/g" > /tmp/ingresses.txt
```

Now adopt in dry run mode(only print the ones that will be applied) all `slok.xyz` hosts with the external-dns instance identifier `heritage=external-dns,external-dns/owner=slok-xyz`

```bash
external-dns-aws-migrator \
    -filter ".*\.slok\.xyz$" \
    --txt-owner-id "slok-xyz" \
    --dry-run < /tmp/ingresses.txt 
```

[external-dns]: https://github.com/kubernetes-incubator/external-dns