#!/bin/bash
#set -e -o pipefail

FIXTURES=${GOPATH}/src/github.com/go-swagger/go-swagger/fixtures

if [ ! -f `which swagger` ]; then
  echo "can't find swagger in the PATH"
  exit 1
fi

find ${FIXTURES} -type f \( \( -name \*.json -o -name \.yaml -o -name \*.yml \)  -a -not -name \*codegen\* \) |\
while read spec
do
    swagger validate ${spec}
done
