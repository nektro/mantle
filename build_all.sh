#!/usr/bin/env bash

set -e

GOPATH=$(go env GOPATH)
go install github.com/rakyll/statik@latest
$GOPATH/bin/statik -src="./www/" -f

./build_go.sh linux amd64
