package middleware

import "net/http"

// NewValidation starts a new validation middleware
func NewValidation(context *Context) func(http.ResponseWriter, *http.Request, http.HandlerFunc) {
	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		if _, ok := context.RouteInfo(r); ok {
			rw.WriteHeader(http.StatusOK)
			return
		}
		next(rw, r)
	}
}

type validation struct {
}
