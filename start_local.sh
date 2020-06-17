#!/usr/bin/env bash

set -e
set -x

go build
./mantle \
--skip-translation-fetch \
--config './data/config.json' \
