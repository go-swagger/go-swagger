/*
	Modified version of assetfs file here:
	https://github.com/elazarl/go-bindata-assetfs/blob/master/assetfs.go
*/

package swag

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/elazarl/go-bindata-assetfs"
)

// MiddlewareAt creates a middleware to serve swagger ui at the specified basePath
func MiddlewareAt(basePath string, assetFS func() *assetfs.AssetFS, next http.Handler) http.Handler {
	fileServer := http.FileServer(assetFS())
	if basePath != "" && basePath != "/" {
		fileServer = http.StripPrefix(basePath, fileServer)
	}

	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, basePath) {
			fmt.Println("serving url", r.URL)
			fileServer.ServeHTTP(rw, r)
			return
		}

		if next == nil {
			http.NotFound(rw, r)
		} else {
			next.ServeHTTP(rw, r)
		}
	})
}
