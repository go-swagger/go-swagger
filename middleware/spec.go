package middleware

import "net/http"

func specMiddleware(ctx *Context, next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/swagger.json" {
			rw.Header().Set("Content-Type", "application/json")
			rw.WriteHeader(http.StatusOK)
			rw.Write(ctx.spec.Raw())
			return
		}
		if next == nil {
			http.NotFound(rw, r)
			return
		}
		next.ServeHTTP(rw, r)
	})
}
