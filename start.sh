#!/usr/bin/env bash

set -e
set -x

./node_modules/.bin/jshint ./www/chat/js/*.js
./node_modules/.bin/eslint ./www/chat/js/*.js

go test
go vet -composites=false

go build
./mantle \
