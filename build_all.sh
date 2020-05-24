#!/usr/bin/env bash

go get -v github.com/rakyll/statik
$GOPATH/bin/statik -src="./www/" -f

date=$(date +'%Y.%m.%d')
version=${CIRCLE_BUILD_NUM-$date}
tag=v$version-$(git log --format=%h -1)

go get -v github.com/karalabe/xgo
$GOPATH/bin/xgo -ldflags="-s -w -X main.Version=$tag" -out="./bin/mantle-$tag" --targets=*/amd64,linux/* github.com/nektro/mantle

go mod tidy
