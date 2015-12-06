#!/bin/bash

set -e -x

go test -v -race $(go list ./... | grep -v vendor) | go-junit-report -dir ./shippable/testresults/junit.xml

# Run test coverage on each subdirectories and merge the coverage profile.
echo "mode: count" > profile.cov
repo_pref="github.com/${REPO_NAME-"$(basename `pwd`)/$(basename `pwd`)"}/"
# Standard go tooling behavior is to ignore dirs with leading underscores
set -x
for dir in $(go list ./... | grep -v -E 'vendor|generator')
do
  pth="${dir//*$repo_pref}"
  godep go test -covermode=set -coverprofile=${pth}/profile.tmp $dir
  if [ -f $pth/profile.tmp ]
  then
      cat $pth/profile.tmp | tail -n +2 >> profile.cov
      rm $pth/profile.tmp
  fi
done

set +x
godep go tool cover -func profile.cov
gocov convert profile.cov | gocov report
gocov convert profile.cov | gocov-html > ./shippable/codecoverage/coverage-$BUILD_NUMBER.html
