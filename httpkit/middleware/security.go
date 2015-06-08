package middleware

import "net/http"

func newSecureAPI(ctx *Context, next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		route, _ := ctx.RouteInfo(r)
		if len(route.Authenticators) == 0 {
			next.ServeHTTP(rw, r)
			return
		}

		if _, err := ctx.Authorize(r, route); err != nil {
			ctx.Respond(rw, r, route.Produces, route, err)
			return
		}

		next.ServeHTTP(rw, r)
	})
}
