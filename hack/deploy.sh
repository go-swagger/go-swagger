#!/bin/bash

set -eux -o pipefail

prjdir=`git rev-parse --show-toplevel`

build_binary() {
  docker run \
    --rm \
    -v $prjdir:/go/src/github.com/$CIRCLE_PROJECT_USERNAME/$CIRCLE_PROJECT_REPONAME \
    -w /go/src/github.com/$CIRCLE_PROJECT_USERNAME/$CIRCLE_PROJECT_REPONAME \
    -e CIRCLE_TAG \
    -e CIRCLE_SHA1 \
    -e CIRCLE_PROJECT_USERNAME \
    -e CIRCLE_PROJECT_REPONAME \
    -e API_EMAIL \
    -e API_USERNAME \
    casualjim/gox -ldflags "-X github.com/$CIRCLE_PROJECT_USERNAME/$CIRCLE_PROJECT_REPONAME/cmd/swagger/commands.Commit=${CIRCLE_SHA1} -X github.com/$CIRCLE_PROJECT_USERNAME/$CIRCLE_PROJECT_REPONAME/cmd/swagger/commands.Version=${CIRCLE_TAG}" -output "./dist/bin/{{.Dir}}_{{.OS}}_{{.Arch}}" "$@" ./cmd/swagger
}

prepare() {
  rm -rf $prjdir/dist
  mkdir -p $prjdir/dist/{bin,build}
  mkdir -p $prjdir/dist/linux/amd64/usr/bin
}

prepare_linuxpkg() {
  cp $prjdir/dist/bin/swagger_linux_amd64 $prjdir/dist/linux/amd64/usr/bin/swagger
}

build_linuxpkg() {
  docker run \
    --rm \
    -v $prjdir:/go/src/github.com/$CIRCLE_PROJECT_USERNAME/$CIRCLE_PROJECT_REPONAME \
    -w /go/src/github.com/$CIRCLE_PROJECT_USERNAME/$CIRCLE_PROJECT_REPONAME \
    -e CIRCLE_TAG \
    -e CIRCLE_SHA1 \
    -e CIRCLE_PROJECT_USERNAME \
    -e CIRCLE_PROJECT_REPONAME \
    -e API_EMAIL \
    -e API_USERNAME \
    casualjim/fpm -t $1 -p ./dist/build -s dir -C ./dist/linux/amd64 -v $CIRCLE_TAG -n swagger --license "ASL 2.0" -a x86_64 -m $API_EMAIL --url "https://goswagger.io" usr
}

upload_to_github() {
  echo "uploading to github"
  cd $prjdir/dist/bin

  sha1sum * > sha1sum.txt
  sha256sum * > sha256sum.txt
  github-release release -u $CIRCLE_PROJECT_USERNAME -r $CIRCLE_PROJECT_REPONAME -t $CIRCLE_TAG -c $CIRCLE_SHA1 -d "$(cat ./notes/v$CIRCLE_TAG.md)"
  github-release upload -u $CIRCLE_PROJECT_USERNAME -r $CIRCLE_PROJECT_REPONAME -t $CIRCLE_TAG
}

upload_to_bintray() {
  curl \
    --retry 10 \
    --retry-delay 5 \
    -T ./dist/build/swagger-${CIRCLE_TAG//-/_}-1.x86_64.rpm \
    -u${API_USERNAME}:${BINTRAY_TOKEN} \
    https://api.bintray.com/content/go-swagger/goswagger-rpm/swagger/${CIRCLE_TAG}/swagger-${CIRCLE_TAG//-/_}-1.x86_64.rpm

  curl \
    --retry 10 \
    --retry-delay 5 \
    -T ./dist/build/swagger_${CIRCLE_TAG}_amd64.deb \
    -u${API_USERNAME}:${BINTRAY_TOKEN} \
    "https://api.bintray.com/content/go-swagger/goswagger-debian/swagger/${CIRCLE_TAG}/swagger_${CIRCLE_TAG}_amd64.deb;deb_distribution=ubuntu;deb_component=main;deb_architecture=amd64"
}

deploy_docker() {
  docker build -t quay.io/goswagger/swagger:$CIRCLE_TAG .
  docker tag quay.io/goswagger/swagger:$CIRCLE_TAG quay.io/goswagger/swagger:latest
  docker login -u $API_USERNAME -p $QUAY_PASS https://quay.io
  docker push quay.io/goswagger/swagger
}

prepare

# build binaries
build_binary -os="linux darwin windows" -arch="amd64 386"
build_binary -os="linux" -arch="arm64 arm"

# # build linux packages
prepare_linuxpkg
build_linuxpkg deb
build_linuxpkg rpm

upload binary packages
upload_to_github
upload_to_bintray

deploy_docker
