#!/bin/bash
set -e -o pipefail

FIXTURES=fixtures/canary

if [[ "${SWAGGER_BIN}" ]]; then
  cp "${SWAGGER_BIN}" /go/bin/
fi

if [ ! -f `which swagger` ]; then
  echo "can't find swagger in the PATH"
  exit 1
fi

for dir in $(ls "${FIXTURES}")
do
  if [ $dir != "bitbucket.org" ]; then
    pushd "${FIXTURES}/$dir"
    rm -rf client models restapi cmd
    swagger generate client --skip-validation
    go test ./...
    if [ $dir != 'kubernetes' ] && [ $dir != 'ms-cog-sci' ] ; then
      swagger generate server --skip-validation
      go test ./...
    fi
    popd
  else
    echo "$dir is disabled for now"
  fi
done
