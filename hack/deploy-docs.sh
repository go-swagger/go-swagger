#!/bin/bash

echo -e "\033[0;32mDeploying updates to GitHub...\033[0m"
wd=`git rev-parse --show-toplevel`
cd "${wd}/docs"

# Build the project.
hugo # if using a theme, replace by `hugo -t <yourtheme>`

cd "${wd}"
# Add changes to git.
git add -A docs

# Commit changes.
msg="rebuilding site `date`"
if [ $# -eq 1 ]
  then msg="$1"
fi
git commit -a -m "$msg"

# Push source and build repos.
git subtree push --prefix=docs/public https://github.com/go-swagger/go-swagger.github.io master
git push

# Come Back
cd "${wd}"
