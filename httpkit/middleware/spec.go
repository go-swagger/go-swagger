package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/casualjim/go-swagger/spec"
)

func specMiddleware(ctx *Context, next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/swagger.json" {
			rw.Header().Set("Content-Type", "application/json")
			rw.WriteHeader(http.StatusOK)
			rw.Write(ctx.spec.Raw())
			return
		}
		if next == nil {
			ctx.NotFound(rw, r)
			return
		}
		next.ServeHTTP(rw, r)
	})
}

// Spec creates a middleware to serve a swagger spec.
// This allows for altering the spec before starting the http listener.
// This can be useful
//
func Spec(basePath string, swsp *spec.Swagger, next http.Handler) http.Handler {
	if basePath == "/" {
		basePath = ""
	}
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if r.URL.Path == basePath+"/swagger.json" {
			rw.Header().Set("Content-Type", "application/json")
			rw.WriteHeader(http.StatusOK)
			dec := json.NewEncoder(rw)
			dec.Encode(swsp)
			return
		}
		if next == nil {
			rw.Header().Set("Content-Type", "application/json")
			rw.WriteHeader(http.StatusNotFound)
			return
		}
		next.ServeHTTP(rw, r)
	})
}
