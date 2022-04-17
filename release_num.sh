#!/usr/bin/env bash

set -e

tagcount=$(git tag | wc -l)

echo $tagcount
