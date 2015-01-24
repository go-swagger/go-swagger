package middleware

import "net/http"

func newOperationExecutor(ctx *Context) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		// use context to lookup routes
		route, _ := ctx.RouteInfo(r)
		bound, _ := ctx.BindAndValidate(r, route)
		result, err := route.Handler.Handle(bound)

		if err != nil {
			ctx.Respond(rw, r, route.Produces, err)
			return
		}

		ctx.Respond(rw, r, route.Produces, result)
	})
}
