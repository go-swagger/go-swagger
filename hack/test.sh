#!/bin/bash

# Bails on any command failure
set -e -o pipefail

cd ${0%/*}/..
echo "Running tests in $(pwd)..."
# List of packages to test
# Currently no packaged tests are available in fixtures or examples
packages=$(go list ./... | grep -v -E 'vendor|fixtures|examples')
repo_pref="github.com/${CIRCLE_PROJECT_USERNAME-"$(basename `pwd`)"}/${CIRCLE_PROJECT_REPONAME-"$(basename `pwd`)"}/"

if [[ ${1} == "--nocover" ]] ; then
    # Run simple tests without coverage computations, but with race detector turned on
    echo "Running unit tests with race detector"
    go test -race -v ${packages}
else
    # Run test coverage on each subdirectories and merge the coverage profile.
    echo "Running CI unit tests with coverage calculation"
    echo "mode: ${GOCOVMODE-atomic}" > coverage.txt
    # Standard go tooling behavior is to ignore dirs with leading underscores
    for dir in ${packages} ; do
        pth="${dir//*$repo_pref}"
        # -tags netgo: test as statically linked
        # -installsuffix netgo: produce suffixed object for this statically linked build
        go test -tags netgo -installsuffix netgo -covermode=${GOCOVMODE-atomic} -coverprofile=${pth}/profile.tmp $dir
        if [[ -f $pth/profile.tmp ]] ; then
            cat $pth/profile.tmp | tail -n +2 >> coverage.txt
            rm -f $pth/profile.tmp
        fi
    done
    go tool cover -func coverage.txt
    # print out coverage report
    gocov convert coverage.txt | gocov report
    outputdir="/usr/share/coverage"
    if [[ ! -d ${outputdir} ]] ; then
        mkdir -p ${outputdir}
    fi
    gocov convert coverage.txt | gocov-html > ${outputdir}/coverage-${CIRCLE_BUILD_NUM-"0"}.html
fi
