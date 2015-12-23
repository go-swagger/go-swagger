#!/bin/sh

go get -u -v github.com/alecthomas/gometalinter
go get -u -v golang.org/x/tools/cmd/goimports
go get -u -v golang.org/x/tools/cmd/cover
go get -u -v golang.org/x/tools/cmd/vet
go get -u -v golang.org/x/tools/cmd/stringer
go get -u -v golang.org/x/tools/cmd/gotype
go get -u -v golang.org/x/tools/cmd/godoc
go get -u -v github.com/FiloSottile/gvt
go get -u -v github.com/jteeuwen/go-bindata/...
go get -u -v github.com/elazarl/go-bindata-assetfs/...
go get -u -v github.com/redefiance/go-find-references
#go get -u -v code.google.com/p/gomock/gomock
#go get -u -v code.google.com/p/gomock/mockgen
go get -u -v github.com/rafrombrc/gomock/...
go get -u -v github.com/axw/gocov/gocov
go get -u -v gopkg.in/matm/v1/gocov-html
go get -u -v github.com/AlekSi/gocov-xml
go get -u -v github.com/nsf/gocode
go get -u -v github.com/jstemmer/gotags
go get -u -v github.com/smartystreets/goconvey
go get -u -v github.com/rogpeppe/godef
go get -u -v github.com/pquerna/ffjson
go get -u -v github.com/nathany/looper
go get -u -v github.com/kylelemons/godebug/...
go get -u -v github.com/aktau/github-release
go get -u -v github.com/xoebus/anderson
go get -u -v github.com/spf13/hugo
#go get -u -v github.com/clipperhouse/gen


# install all the linters
gometalinter --install --update
