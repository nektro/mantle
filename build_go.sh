#!/usr/bin/env bash

set -e

export CGO_ENABLED=1
export GOOS=$1
export GOARCH=$2
export GOARM=7
tag="r$(./release_num.sh)"
rev=$(git log --format=%h -1)
echo "$tag.$rev $GOOS $GOARCH"
go build -ldflags="-s -w -X main.Version=$tag" -o ./bin/mantle-$GOOS-$GOARCH$ext
