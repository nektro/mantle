#!/usr/bin/env bash

set -e

go build
./mantle \
--skip-translation-fetch \
--config './data/config.json' \
--port 80 \
