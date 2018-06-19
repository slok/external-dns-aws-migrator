#!/usr/bin/env sh

set -o errexit
set -o nounset


DIR="$( cd "$( dirname "${0}" )" && pwd )"
ROOT_DIR=${DIR}/../..
SRC=github.com/slok/external-dns-aws-adopter/cmd/external-dns-aws-adopter
OUT=${ROOT_DIR}/bin/external-dns-aws-adopter
LDF_CMP="-w -extldflags '-static'"

echo "Building binary at ${OUT}"
CGO_ENABLED=0 go build -o ${OUT} --ldflags "${LDF_CMP}"  ${SRC}