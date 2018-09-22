#!/bin/bash
#
# A small utility to generate clients and servers on
# well known specifications.

# Bails on any command failure
set -e -o pipefail

cd $(git rev-parse --show-toplevel)

if [[ "${SWAGGER_BIN}" ]]; then
  cp "${SWAGGER_BIN}" /go/bin/
fi

if ! command -v swagger >/dev/null 2>&1; then
  echo "can't find swagger in the PATH"
  exit 1
fi

for dir in fixtures/canary/*
do
    [ ! -d "$dir" ] && continue
    echo "validating '$dir'"
    pushd "$dir" > /dev/null

    case $dir in
    "fixtures/canary/bitbucket.org")
        # bitbucket.org generates wrong model
        client=false
        server=false
        echo "$dir is disabled for now"
        ;;
    "fixtures/canary/kubernetes"|"fixtures/canary/docker")
        # docker has an invalid spec with duplicate operationIds. Generates on docker-fixed
        # kubernetes uses unsupported media type options (issue#1377)
        client=true
        server=false
        echo "$dir is disabled for server generation now (only client is generated)"
        ;;
    *)
        client=true
        server=true
        ;;
    esac
    if [[ "${client}" == "true" ]] ; then
        rm -rf client models restapi cmd
        echo "generating client for $dir..."
        swagger generate client --skip-validation --quiet
        go test -vet off ./...
    fi

    if [[ "${server}" == "true" ]] ; then
        echo "generating server for $dir..."
        swagger generate server --skip-validation --quiet
        go test -vet off ./...
    fi
    set +e
    popd > /dev/null || true
    set -e
done
