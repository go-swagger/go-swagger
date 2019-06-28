#!/bin/bash

set -e -o pipefail

workdir=cover
profile="$workdir/cover.out"
mode=${GOCOVMODE-atomic}

# Run test coverage on each subdirectories and merge the coverage profile.
echo "mode: $mode" > "$profile"

# Standard go tooling behavior is to ignore dirs with leading underscores
for dir in $(go list ./... | grep -v -E 'vendor|fixtures|examples')
do
  f="$workdir/$(echo "$dir" | tr / -).cover"
  go test -tags netgo -installsuffix netgo -covermode="$mode" -coverprofile="$f" "$dir"
  if [ -f "$f" ]
  then
    cat "$f" | tail -n +2 >> "$profile"
    rm "$f"
  fi
done

go tool cover -func "$profile"
gocov convert "$profile" | gocov report
gocov convert "$profile" | gocov-html > "$workdir/coverage.html"
