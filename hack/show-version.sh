#!/bin/sh

current_tag_sha1=`git rev-list --tags --max-count=1`
current_tag=`git describe --tags ${current_tag_sha1}`
since_tag=`git rev-list --count ${current_tag}..HEAD`
commit_hash=`git rev-parse --short HEAD`

if [ $since_tag -gt 0 ]; then
  echo "${current_tag}-${since_tag}-${commit_hash}"
else
  echo "${current_tag}"
fi
