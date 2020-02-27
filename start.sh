#!/usr/bin/env bash

set -e
set -x

jshint ./www/chat/js/*.js

go test

go build
./mantle \
