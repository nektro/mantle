#!/usr/bin/env bash

set -e
set -x

docker build --rm --build-arg=BUILD_NUM=4 . -t mantle_local
docker run --rm -p 80:8000 -v ~/dev/golang/mantle/data/:/data/ mantle_local
