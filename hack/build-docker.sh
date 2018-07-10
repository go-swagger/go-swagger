#!/bin/bash

# Bails on any command failure
set -e -o pipefail -x

cd `git rev-parse --show-toplevel`
echo "Building swagger from $(pwd)..."

if [[ ${1} == "--circleci" ]] ; then
    # CI build mode (for releases)
    LDFLAGS="-s -w -linkmode external -extldflags \"-static\""
    LDFLAGS="$LDFLAGS -X github.com/${CIRCLE_PROJECT_USERNAME-"$(basename `pwd`)"}/${CIRCLE_PROJECT_REPONAME-"$(basename `pwd`)"}/cmd/swagger/commands.Commit=${CIRCLE_SHA1} -X github.com/${CIRCLE_PROJECT_USERNAME-"$(basename `pwd`)"}/${CIRCLE_PROJECT_REPONAME-"$(basename `pwd`)"}/cmd/swagger/commands.Version=${CIRCLE_TAG-dev}"
    go build -a -o /usr/share/dist/swagger --ldflags "$LDFLAGS" ./cmd/swagger
else
    # manual build mode
    go build -o /usr/share/dist/swagger ./cmd/swagger
fi

go install ./cmd/swagger
