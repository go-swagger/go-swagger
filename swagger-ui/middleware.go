package swaggerui

import (
	"io"
	"net/http"
	"os"

	"github.com/casualjim/go-swagger/util"
)

// Middleware serves this static site as middleware on /swagger-ui
// passing nil will make this a regular handler not a middleware
func Middleware(path string, next http.Handler) http.Handler {

	assetFS := func() *util.AssetFS {
		return &util.AssetFS{Asset: Asset, AssetDir: AssetDir, Prefix: "/"}
	}

	swaggerUI := util.MiddlewareAt("/swagger-ui", assetFS, next)
	if path == "" {
		return swaggerUI
	}

	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api-docs" {
			specFile, err := os.Open(path)
			defer specFile.Close()

			if err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				return
			}

			rw.Header().Set("Content-Type", "application/json")
			rw.WriteHeader(http.StatusOK)
			io.Copy(rw, specFile)
			return
		}
		swaggerUI.ServeHTTP(rw, r)
	})
}
