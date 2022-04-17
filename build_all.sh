#!/usr/bin/env bash

set -e

GOPATH=$(go env GOPATH)
go get -v github.com/rakyll/statik
$GOPATH/bin/statik -src="./www/" -f

./build_go.sh linux amd64
