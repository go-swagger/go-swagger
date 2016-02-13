#!/bin/bash

version=${1-"0.2.5"}
mkdir -p /drone/dist/binaries
gox -os="linux darwin windows" -output="/drone/dist/binaries/{{.Dir}}_{{.OS}}_{{.Arch}}" ./cmd/swagger

cd /drone/dist/binaries
gzip -f swagger_*

ls

if [ -z "${GITHUB_TOKEN}" ]; then
  echo 'Please set GITHUB_TOKEN...'
  exit 1
fi

export GITHUB_USER="${GITHUB_USER:-"go-swagger"}"
export GITHUB_REPO="${GITHUB_REPO:-"go-swagger"}"

# Generate description
description=$(
if [[ "${version}" == "prerelease-"* ]]; then
  echo '**This is a PRERELEASE version.**'
fi
echo '
The binaries below are provided without warranty, following the [Apache license](LICENSE).
'
echo '
Instructions:
* Download the file relevant to your operating system
* Decompress (i.e. `gzip -d swagger_linux_amd64.gz`)
* Set the executable bit (i.e. `chmod +x swagger_linux_amd64`)
* Move the file to a directory in your `$PATH` (i.e. `mv swagger_linux_amd64 /usr/local/bin`)
'
echo '```'
echo '$ sha1sum swagger_*.gz'
sha1sum swagger_*.gz
echo '```'
)


echo "Creating release..."
github-release release --tag "${version}" --name "${name}" --description "${description}"

# Upload build artifacts
for f in /drone/dist/binaries/swagger_*; do
  b=$(basename ${f})
  echo "Uploading $f..."
  github-release upload --tag "${version}" --name "${b}" --file "${f}"
done

# cd /usr/share/dist
# mkdir -p /usr/share/dist/linux/amd64/usr/bin
# cp /usr/share/dist/binaries/swagger_linux_amd64 /usr/share/dist/linux/amd64/usr/bin/swagger
# fpm -t deb -s dir -C /usr/share/dist/linux/amd64 -v "$version" -n swagger --license "ASL 2.0" -a x86_64 -m "ivan@flanders.co.nz" --url "https://goswagger.io" usr
# fpm -t rpm -s dir -C /usr/share/dist/linux/amd64 -v "$version" -n swagger --license "ASL 2.0" -a x86_64 -m "ivan@flanders.co.nz" --url "https://goswagger.io" usr
