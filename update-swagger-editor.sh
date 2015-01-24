#!/bin/sh

version=2.8.2
vname=v$version

curdir=`pwd`
rm -rf /tmp/swaggereditor || true
cd /tmp
curl -OL'#' https://github.com/swagger-api/swagger-editor/archive/$vname.tar.gz | tar zx 
cd swagger-editor-${version}
npm install
bower install
grunt build

[[ -f $curdir/swagger-editor/middleware.go ]] && cp $curdir/swagger-editor/middleware.go ./swagger-editor-middleware.go

rm -rf $curdir/swagger-editor
mv dist $curdir/swagger-editor

cd /tmp
rm -rf swagger-editor-$version

cd $curdir/swagger-editor

go-bindata -pkg=swaggereditor ./...

cat > ./middleware.go <<TEMPLATE
package swaggereditor

import (
  "fmt"
  "io"
  "net/http"
  "os"

  "github.com/casualjim/go-swagger/util"
)

func serveWritable(path string, rw http.ResponseWriter, r *http.Request) {
  if r.Method == "PUT" {
    specFile, err := os.Create(path)
    defer specFile.Close()

    if err != nil {
      http.Error(rw, err.Error(), http.StatusInternalServerError)
      return
    }

    if _, err := io.Copy(specFile, r.Body); err != nil {
      http.Error(rw, err.Error(), http.StatusInternalServerError)
      return
    }
    if err := specFile.Close(); err != nil {
      http.Error(rw, err.Error(), http.StatusInternalServerError)
      return
    }
  }

  specFile, err := os.Open(path)
  defer specFile.Close()

  if err != nil {
    http.Error(rw, err.Error(), http.StatusInternalServerError)
    return
  }

  rw.Header().Set("Content-Type", "application/json")
  rw.WriteHeader(http.StatusOK)

  if _, err := io.Copy(rw, specFile); err != nil {
    http.Error(rw, err.Error(), http.StatusInternalServerError)
    return
  }
}

var defaultConfig = \`{
  "codegen": {
    "servers": "http://generator.wordnik.com/online/api/gen/servers",
    "clients": "http://generator.wordnik.com/online/api/gen/clients",
    "server": "http://generator.wordnik.com/online/api/gen/servers/{language}",
    "client": "http://generator.wordnik.com/online/api/gen/clients/{language}"
  },
  "disableCodeGen": true,
  "autocompleteExtension": {},
  "useBackendForStorage": true,
  "backendEndpoint": "/swagger.json",
  "backendHelathCheckTimeout": 5000,
  "useYamlBackend": false,
  "disableFileMenu": true,
  "headerBranding": false,
  "enableTryIt": true,
  "brandingCssClass": "",
  "schemaUrl": "/schema/swagger.json"
}\`

// Middleware serves this static site as middleware on /swagger-editor
// passing nil will make this a regular handler not a middleware
func Middleware(path string, next http.Handler) http.Handler {
  assetFS := func() *util.AssetFS {
    return &util.AssetFS{Asset: Asset, AssetDir: AssetDir, Prefix: "/"}
  }

  editor := util.MiddlewareAt("/swagger-editor", assetFS, next)

  writableAPIDocs := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
    fmt.Println("writeable api docs")
    if r.URL.Path == "/swagger.json" {
      serveWritable(path, rw, r)
      return
    }
    editor.ServeHTTP(rw, r)
  })

  return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
    if r.URL.Path == "/config/defaults.json" {
      fmt.Println("serving /config/defaults.json")
      rw.Header().Set("Content-Type", "application/json")
      rw.WriteHeader(http.StatusOK)
      rw.Write([]byte(defaultConfig))
      return
    }
    writableAPIDocs.ServeHTTP(rw, r)
  })
}
TEMPLATE