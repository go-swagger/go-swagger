package middleware

import (
	"net/http"

	"github.com/casualjim/go-swagger/errors"
	"github.com/gorilla/context"
)

func newCompleteMiddleware(ctx *Context) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		defer context.Clear(r)

		// use context to lookup routes
		if matched, ok := ctx.RouteInfo(r); ok {
			// TODO: add security checks here

			bound, validation := ctx.BindAndValidate(r, matched)
			if validation.HasErrors() {
				ctx.Respond(rw, r, matched.Produces, validation.Errors[0])
				return
			}

			result, err := matched.Handler.Handle(bound)
			if err != nil {
				ctx.Respond(rw, r, matched.Produces, err)
				return
			}

			ctx.Respond(rw, r, matched.Produces, result)
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
