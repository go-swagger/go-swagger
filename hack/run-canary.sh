#!/bin/bash
#
# A small utility to generate clients and servers on
# well known specifications.

# Bails on any command failure
set -e -o pipefail

basedir=`cd ${0%/*};pwd`
FIXTURES=${basedir}/../fixtures/canary

if [[ "${SWAGGER_BIN}" ]]; then
  cp "${SWAGGER_BIN}" /go/bin/
fi

if [ ! -f `which swagger` ]; then
  echo "can't find swagger in the PATH"
  exit 1
fi

for dir in $(ls "${FIXTURES}")
do
    echo $dir
    pushd "${FIXTURES}/$dir"

    case $dir in
    "bitbucket.org")          
        # bitbucket.org generates wrong model 
        client=false
        server=false
        echo "$dir is disabled for now"
        ;;
    kubernetes|docker)
        # docker has an invalid spec with duplicate operationIds. Generates on docker-fixed
        # kubernetes uses unsupported media type options (issue#1377)
        # ms-cog-sci
        client=true
        server=false
        echo "$dir is disabled for server generation now"
        ;;
    *)
        client=true
        server=true
        ;;
    esac
    if [[ ${client} == "true" ]] ; then
        rm -rf client models restapi cmd
        swagger generate client --skip-validation --quiet
        go test ./...
    fi

    if [[ ${server} == "true" ]] ; then
        swagger generate server --skip-validation --quiet
        go test ./...
    fi
    popd
done
