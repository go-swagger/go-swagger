#!/bin/bash
version=$1
mkdir -p /usr/share/dist/{binaries,packages}
gox -os="linux darwin windows" -output="/usr/share/dist/binaries/{{.Dir}}_{{.OS}}_{{.Arch}}" ./cmd/swagger

cd /usr/share/dist
mkdir -p /usr/share/dist/linux/amd64/usr/bin
cp /usr/share/dist/binaries/swagger_linux_amd64 /usr/share/dist/linux/amd64/usr/bin/swagger
fpm -t deb -s dir -C /usr/share/dist/linux/amd64 -v "$version" -n swagger --license "ASL 2.0" -a x86_64 -m "ivan@flanders.co.nz" --url "https://goswagger.io" usr
fpm -t rpm -s dir -C /usr/share/dist/linux/amd64 -v "$version" -n swagger --license "ASL 2.0" -a x86_64 -m "ivan@flanders.co.nz" --url "https://goswagger.io" usr
