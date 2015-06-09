// Package assets contains the embedded assets like the json schema json doc and the swagger 2.0 schema doc
package assets

//go:generate go-bindata -pkg=assets -prefix=../schemas -ignore=.*\.md ../schemas/...
//go:generate perl -pi -e s,Json,JSON,g bindata.go
