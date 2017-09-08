#!/bin/bash

set -o errexit
set -o pipefail

gometalinter \
  --exclude='error return value not checked.*(Close|Log|Print|Shutdown).*\(errcheck\)$' \
  --skip=fixtures \
  --skip=examples \
  --tests \
  --vendor \
  --disable=aligncheck \
  --disable=gotype \
  --disable=goconst \
  --cyclo-over=20 \
  --deadline=300s \
  ./cmd/...

gometalinter \
  --exclude='^generator/bindata\.go.*$' \
  --exclude='error return value not checked.*(Close|Log|Print|RemoveAll|Setenv|Shutdown).*\(errcheck\)$' \
  --exclude='^scan/schema\.go.*pkg can be fmt.Stringer \(interfacer\)$' \
  --skip=fixtures \
  --skip=examples \
  --skip=cmd \
  --tests \
  --vendor \
  --disable=aligncheck \
  --disable=gotype \
  --disable=goconst \
  --disable=gocyclo \
  --deadline=300s \
  ./...
