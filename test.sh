#!/usr/bin/env bash

set -e
set -x
shopt -s globstar


########################
#    Golang Testing    #
########################

# https://golang.org/pkg/testing/
go test

# https://golang.org/cmd/vet/
go vet -composites=false


########################
#  Javascript Testing  #
########################

# https://jshint.com/
./node_modules/.bin/jshint ./www/chat/js/**/*.js ./tests/**/*.js

# https://eslint.org/
./node_modules/.bin/eslint ./www/chat/js/**/*.js ./tests/**/*.js
