#!/usr/bin/env bash

set -e
set -x

GOPATH=$(go env GOPATH)

########################
#    Golang Testing    #
########################

go install github.com/fzipp/gocyclo/cmd/gocyclo@latest
go install golang.org/x/lint/golint@latest
go install github.com/gordonklaus/ineffassign@latest
go install github.com/client9/misspell/cmd/misspell@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# https://golang.org/pkg/testing/
go test

# https://golang.org/cmd/vet/
go vet -composites=false

# https://github.com/fzipp/gocyclo
$GOPATH/bin/gocyclo -over 9 .

# https://golang.org/x/lint
$GOPATH/bin/golint .

# https://github.com/gordonklaus/ineffassign
$GOPATH/bin/ineffassign .

# https://github.com/client9/misspell
$GOPATH/bin/misspell -error .

# https://github.com/golangci/golangci-lint
$GOPATH/bin/golangci-lint run -D errcheck

########################
#  Javascript Testing  #
########################

# https://jshint.com/
./node_modules/.bin/jshint --reporter ./scripts/.jshintrc.fmt.js ./www/

# https://eslint.org/
./node_modules/.bin/eslint --format ./scripts/.eslintrc.fmt.js ./www/ --ext .js,.html
