#!/usr/bin/env bash

set -e

tag=r$(./release_num.sh)

GITHUB_TOKEN="$1"
PROJECT_USERNAME=$(echo $GITHUB_REPOSITORY | cut -d'/' -f1)
PROJECT_REPONAME=$(echo $GITHUB_REPOSITORY | cut -d'/' -f2)

~/.zigmod/bin/ghr \
    -t ${GITHUB_TOKEN} \
    -u ${PROJECT_USERNAME} \
    -r ${PROJECT_REPONAME} \
    -b "$(./changelog.sh)" \
    "$tag" \
    "./bin/"
