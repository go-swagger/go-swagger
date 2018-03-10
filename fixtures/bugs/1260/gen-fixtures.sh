#! /bin/bash
if [[ ${1} == "--clean" ]] ; then
    # remove generation after successul pass
    clean=1
fi

. ./build-script.sh

# A small utility to build fixture servers
testcases="fixture-realiased-types.yaml"
testcases="${testcases} test3-swagger.yaml test3-bis-swagger.yaml test3-ter-swagger.yaml test3-ter-swagger-flat.json"
build_fixtures "server" ${testcases}
