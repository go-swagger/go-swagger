#!/bin/bash

set -e

cd "${GOPATH/:*/}src/github.com/casualjim/go-swagger"

godep go test -race ./...
GOPATH=`godep path:$GOPATH` goveralls -service=circle-ci -repotoken=$COVERALLS_TOKEN
