#!/usr/bin/env bash

set -o errexit
set -o nounset

go test `go list ./... | grep -v vendor` -v -tags='integration'