#!/usr/bin/env bash

set -e
set -x


########################
#    Golang Testing    #
########################

# https://golang.org/pkg/testing/
go test

# https://golang.org/cmd/vet/
go vet -composites=false

# https://github.com/fzipp/gocyclo
$GOPATH/bin/gocyclo -over 10 .

# https://github.com/golang/lint
$GOPATH/bin/golint .

# https://github.com/gordonklaus/ineffassign
$GOPATH/bin/ineffassign .

# https://github.com/client9/misspell
$GOPATH/bin/misspell -error .


########################
#  Javascript Testing  #
########################

# https://jshint.com/
./node_modules/.bin/jshint ./www/chat/js/

# https://eslint.org/
./node_modules/.bin/eslint ./www/chat/js/
