#!/bin/bash

set -o errexit
set -o pipefail

gometalinter \
  --exclude='error return value not checked.*(Close|Log|Print|Shutdown).*\(errcheck\)$' \
  --exclude='declaration of "err" shadows declaration.*\(vetshadow\)$' \
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
  --exclude='error return value not checked.*(Close|Log|Print|RemoveAll|Shutdown).*\(errcheck\)$' \
  --exclude='declaration of "err" shadows declaration.*\(vetshadow\)$' \
  --exclude='^scan/schema\.go.*pkg can be fmt.Stringer \(interfacer\)$' \
  --exclude='^scan/.*unused struct field.*\(structcheck\)$' \
  --exclude='^scan/scanner\.go.*unused.*\((deadcode|varcheck)\)$' \
  --exclude='^scan/schema\.go.*newSchemaAnnotationParser is unused.*\(deadcode\)$'\
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
