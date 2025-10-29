#!/bin/bash

# Bails on any command failure
set -e -o pipefail -x

cd "$(git rev-parse --show-toplevel)"
echo "Building swagger from $(pwd)..."

target="${1}"
shift
extra="${@}"

if [[ "${target}" == "--circleci" ]] ; then
    # CI build mode (for releases)
    username="${CIRCLE_PROJECT_USERNAME-"$(basename "$PWD")"}"
    project="${CIRCLE_PROJECT_REPONAME-"$(basename "$PWD")"}"
    commit_property="github.com/$username/$project/cmd/swagger/commands.Commit=${CIRCLE_SHA1}"
    tag_property="github.com/$username/$project/cmd/swagger/commands.Version=${CIRCLE_TAG-dev}"

    LDFLAGS="-s -w -X $commit_property -X $tag_property"
    go build -a -o /usr/share/dist/swagger ${extra} --ldflags "$LDFLAGS" ./cmd/swagger
elif [[ "${target}" == "--github-action" ]] ; then
    # Github workflows build mode (for releases)
    commit_property="github.com/${GITHUB_REPOSITORY}/cmd/swagger/commands.Commit=${GITHUB_SHA}"
    tag_property="github.com/${GITHUB_REPOSITORY}/cmd/swagger/commands.Version=${GITHUB_REF_NAME-dev}"

    LDFLAGS="-s -w -X $commit_property -X $tag_property"
    go install -a --ldflags "$LDFLAGS" ${extra} ./cmd/swagger
    go build -o ./dist/swagger --ldflags "$LDFLAGS" ${extra} ./cmd/swagger
else
    # manual build mode
    go build -o ./dist/swagger ${extra} ./cmd/swagger
    go install ${extra} ./cmd/swagger
fi


