package middleware

import "net/http"

func newOperationExecutor(ctx *Context) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		// use context to lookup routes
		route, _ := ctx.RouteInfo(r)
		route.Handler.ServeHTTP(rw, r)
	})
}
