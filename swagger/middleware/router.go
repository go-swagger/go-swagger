package middleware

import (
	"net/http"
	"strings"
)

// NewRouter creates a new router middleware function
func NewRouter(context *Context) func(http.ResponseWriter, *http.Request, func(http.ResponseWriter, *http.Request)) {
	return func(rw http.ResponseWriter, r *http.Request, next func(http.ResponseWriter, *http.Request)) {
		// use context to lookup routes
		if _, ok := context.RouteInfo(r); ok {
			next(rw, r)
			return
		}

		// Not found, check if it exists in the other methods first
		if others := context.router.OtherMethods(r.Method, r.URL.Path); len(others) > 0 {
			rw.Header().Add("Allow", strings.Join(others, ","))
			rw.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		rw.WriteHeader(http.StatusNotFound)
	}
}
