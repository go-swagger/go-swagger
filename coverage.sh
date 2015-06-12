#!/bin/bash

set -e

# Run test coverage on each subdirectories and merge the coverage profile.
echo "mode: count" > profile.cov

# Standard go tooling behavior is to ignore dirs with leading underscores
for dir in $(find . -maxdepth 10 -not -path './.git*' -not -path '*/_*' -type d -not -path 'Godeps' -not -path './generator' -not -path './examples*' -not -path './fixtures*' -not -path './swagger-ui')
do
if ls $dir/*.go &> /dev/null; then
    godep go test -covermode=count -coverprofile=$dir/profile.tmp $dir
    if [ -f $dir/profile.tmp ]
    then
        cat $dir/profile.tmp | tail -n +2 >> profile.cov
        rm $dir/profile.tmp
    fi
fi
done

godep go tool cover -func profile.cov
if [ "$TRAVIS_SECURE_ENV_VARS" = "true" ]; then
  goveralls -coverprofile=profile.cov -service=travis-ci -repotoken=$COVERALLS_TOKEN;
fi
if [ "$CIRCLECI" = "true" ]; then
  goveralls -coverprofile=profile.cov -service=circleci -repotoken=$COVERALLS_REPO_TOKEN;
fi
