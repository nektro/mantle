#!/usr/bin/env bash

set -e

export CGO_ENABLED=1
export GOOS=$1
export GOARCH=$2
export GOARM=7
ext=$3
date=$(date +'%Y%m%d')
version=${CIRCLE_BUILD_NUM-$date}
tag=v$version
echo $tag-$GOOS-$GOARCH
go build -ldflags="-s -w -X main.Version=$tag" -o ./bin/mantle-$tag-$GOOS-$GOARCH$ext
