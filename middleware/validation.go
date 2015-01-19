package middleware

import (
	"net/http"

	"github.com/casualjim/go-swagger"
	"github.com/casualjim/go-swagger/errors"
	"github.com/casualjim/go-swagger/httputils"
	"github.com/casualjim/go-swagger/router"
	"github.com/casualjim/go-swagger/validate"
)

// NewValidation starts a new validation middleware
func newValidation(context *Context, next http.Handler) http.Handler {

	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		matched, _ := context.RouteInfo(r)

		result := validateRequest(context, r, matched)
		if result.HasErrors() {
			context.Respond(rw, r, matched.Produces, result.Errors[0])
			return
		}

		next.ServeHTTP(rw, r)
	})
}

type validation struct {
	context  *Context
	result   *result
	request  *http.Request
	route    *router.MatchedRoute
	bound    map[string]interface{}
	consumer swagger.Consumer
}

func validateRequest(context *Context, request *http.Request, route *router.MatchedRoute) *result {
	validate := &validation{context, &result{}, request, route, make(map[string]interface{}), nil}

	validate.contentType()
	validate.responseFormat()
	if validate.result.IsValid() {
		validate.parameters()
	}

	return validate.result
}

func (v *validation) parameters() {
	result := v.route.Binder.Bind(v.request, v.route.Params, v.consumer, v.bound)
	v.result.AddErrors(result.Errors...)
}

func (v *validation) contentType() {
	if httputils.CanHaveBody(v.request.Method) {
		ct, _, err := v.context.ContentType(v.request)
		if err != nil {
			v.result.AddErrors(err)
		} else {
			if err := validate.ContentType(v.route.Consumes, ct); err != nil {
				v.result.AddErrors(err)
			}
			v.consumer = v.route.Consumers[ct]
		}
	}
}

func (v *validation) responseFormat() {
	if str := v.context.ResponseFormat(v.request, v.route.Produces); str == "" {
		v.result.AddErrors(errors.InvalidResponseFormat(v.request.Header.Get(httputils.HeaderAccept), v.route.Produces))
	}
}

type result struct {
	Errors []errors.Error
}

func (r *result) AddErrors(errors ...errors.Error) {
	r.Errors = append(r.Errors, errors...)
}

func (r *result) IsValid() bool {
	return len(r.Errors) == 0
}

func (r *result) HasErrors() bool {
	return !r.IsValid()
}
