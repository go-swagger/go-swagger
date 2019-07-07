#!/bin/bash

# Bails on any command failure
set -e -o pipefail -x

cd $(git rev-parse --show-toplevel)
echo "Building swagger from $(pwd)..."

if [[ ${1} == "--circleci" ]] ; then
    # CI build mode (for releases)
    username="${CIRCLE_PROJECT_USERNAME-"$(basename `pwd`)"}"
    project="${CIRCLE_PROJECT_REPONAME-"$(basename `pwd`)"}"
    commit_property="github.com/$username/$project/cmd/swagger/commands.Commit=${CIRCLE_SHA1}"
    tag_property="github.com/$username/$project/cmd/swagger/commands.Version=${CIRCLE_TAG-dev}"

    LDFLAGS="-s -w -X $commit_property -X $tag_property"
    go build -a -o /usr/share/dist/swagger --ldflags "$LDFLAGS" ./cmd/swagger
else
    # manual build mode
    go build -o /usr/share/dist/swagger ./cmd/swagger
fi

go install ./cmd/swagger
