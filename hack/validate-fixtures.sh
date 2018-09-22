#!/bin/bash
# Run a spec validation on fixtures

cd $(git rev-parse --show-toplevel)
FIXTURES="${basedir}/fixtures"

if ! command -v swagger >/dev/null 2>&1; then
  echo "can't find swagger in the PATH"
  exit 1
fi

find "${FIXTURES}" -type f \( \( -name \*.json -o -name \.yaml -o -name \*.yml \)  -a -not -name \*codegen\* \) |\
while read -r spec
do
    swagger validate "${spec}" --skip-warnings --stop-on-error
done
