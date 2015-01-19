package middleware

import (
	"net/http"

	"github.com/casualjim/go-swagger/errors"
	"github.com/gorilla/context"
)

func newRouter(ctx *Context, next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		defer context.Clear(r)
		// use context to lookup routes
		if _, ok := ctx.RouteInfo(r); ok {
			next.ServeHTTP(rw, r)
			return
		}

		// Not found, check if it exists in the other methods first
		if others := ctx.AllowedMethods(r); len(others) > 0 {
			ctx.Respond(rw, r, ctx.spec.RequiredProduces(), errors.MethodNotAllowed(r.Method, others))
			return
		}
		ctx.Respond(rw, r, ctx.spec.RequiredProduces(), errors.NotFound("path %s was not found", r.URL.Path))
	})
}
