#!/bin/sh

go get -u -v github.com/golang/lint/golint
go get -u -v golang.org/x/tools/cmd/...
go get -u -v github.com/tools/godep
go get -u -v github.com/jteeuwen/go-bindata/...
go get -u -v github.com/elazarl/go-bindata-assetfs/...
go get -u -v github.com/redefiance/go-find-references
go get -u -v github.com/sqs/goreturns
go get -u -v code.google.com/p/gomock/gomock
go get -u -v code.google.com/p/gomock/mockgen
go get -u -v github.com/axw/gocov/gocov
go get -u -v gopkg.in/matm/v1/gocov-html
go get -u -v github.com/AlekSi/gocov-xml
go get -u -v github.com/nsf/gocode
go get -u -v github.com/kisielk/errcheck
go get -u -v github.com/jstemmer/gotags
go get -u -v github.com/smartystreets/goconvey
go get -u -v github.com/rogpeppe/godef
go get -u -v github.com/pquerna/ffjson
go get -u -v github.com/clipperhouse/gen
