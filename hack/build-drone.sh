#!/bin/bash

repo_pref="${CI_BUILD_DIR##${GOPATH/%:*/}/src/}/"
echo "repo ref: $repo_pref"

set -e
mkdir -p /usr/share/{testresults,coverage,dist/swagger}
go test -v -race $(go list ./... | grep -v vendor) | go-junit-report -dir /usr/share/testresults

# Run test coverage on each subdirectories and merge the coverage profile.
echo "mode: ${GOCOVMODE}" > profile.cov

# Standard go tooling behavior is to ignore dirs with leading underscores
for dir in $(go list ./... | grep -v -E 'vendor|generator')
do
  pth="${dir//*$repo_pref}"
  go test -covermode=${GOCOVMODE} -coverprofile=${pth}/profile.tmp $dir
  if [ -f $pth/profile.tmp ]
  then
      cat $pth/profile.tmp | tail -n +2 >> profile.cov
      rm $pth/profile.tmp
  fi
done

go tool cover -func profile.cov
gocov convert profile.cov | gocov report
gocov convert profile.cov | gocov-html > /usr/share/coverage/coverage-${CI_BUILD_NUM-"0"}.html
rm -rf /usr/share/dist/swagger
go build -o /usr/share/dist/swagger ./cmd/swagger
