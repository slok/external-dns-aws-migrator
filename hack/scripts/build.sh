#!/usr/bin/env sh

set -o errexit
set -o nounset


DIR="$( cd "$( dirname "${0}" )" && pwd )"
ROOT_DIR=${DIR}/../..
SRC=github.com/slok/external-dns-aws-migrator/cmd/external-dns-aws-migrator
OUT=${ROOT_DIR}/bin/external-dns-aws-migrator
LDF_CMP="-w -extldflags '-static'"

echo "Building binary at ${OUT}"
CGO_ENABLED=0 go build -o ${OUT} --ldflags "${LDF_CMP}"  ${SRC}