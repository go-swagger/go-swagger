#!/bin/bash

set -e -o pipefail

# Run test coverage on each subdirectories and merge the coverage profile.
echo "mode: ${GOCOVMODE-atomic}" > coverage.txt
repo_pref="github.com/${CIRCLE_PROJECT_USERNAME-"$(basename `pwd`)"}/${CIRCLE_PROJECT_REPONAME-"$(basename `pwd`)"}/"
# Standard go tooling behavior is to ignore dirs with leading underscores
for dir in $(go list ./... | grep -v -E 'vendor|fixtures|examples')
do
  pth="${dir//*$repo_pref}"
  go test -tags netgo -installsuffix netgo -covermode=${GOCOVMODE-atomic} -coverprofile=${pth}/profile.tmp $dir
  if [ -f $pth/profile.tmp ]
  then
      cat $pth/profile.tmp | tail -n +2 >> coverage.txt
      rm -f $pth/profile.tmp
  fi
done

LDFLAGS="-s -w -linkmode external -extldflags \"-static\""
LDFLAGS="$LDFLAGS -X github.com/$CIRCLE_PROJECT_USERNAME/$CIRCLE_PROJECT_REPONAME/cmd/swagger/commands.Commit=${CIRCLE_SHA1} -X github.com/$CIRCLE_PROJECT_USERNAME/$CIRCLE_PROJECT_REPONAME/cmd/swagger/commands.Version=${CIRCLE_TAG-dev}"

go tool cover -func coverage.txt
gocov convert coverage.txt | gocov report
gocov convert coverage.txt | gocov-html > /usr/share/coverage/coverage-${CIRCLE_BUILD_NUM-"0"}.html
go build -o /usr/share/dist/swagger --ldflags "$LDFLAGS" ./cmd/swagger

go install ./cmd/swagger

./hack/run-canary.sh
