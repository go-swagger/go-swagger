package middleware

import (
	"mime"
	"net/http"

	"github.com/casualjim/go-swagger"
	"github.com/casualjim/go-swagger/errors"
	"github.com/casualjim/go-swagger/httputils"
	"github.com/casualjim/go-swagger/util"
	"github.com/casualjim/go-swagger/validate"
)

// NewValidation starts a new validation middleware
func newValidation(ctx *Context, next http.Handler) http.Handler {

	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		matched, _ := ctx.RouteInfo(r)
		_, result := ctx.BindAndValidate(r, matched)

		if result.HasErrors() {
			ctx.Respond(rw, r, matched.Produces, matched, result.Errors[0])
			return
		}

		next.ServeHTTP(rw, r)
	})
}

type validation struct {
	context *Context
	result  *validate.Result
	request *http.Request
	route   *MatchedRoute
	bound   map[string]interface{}
}

type untypedBinder map[string]interface{}

func (ub untypedBinder) BindRequest(r *http.Request, route *MatchedRoute, consumer swagger.Consumer) error {
	if res := route.Binder.Bind(r, route.Params, consumer, ub); res != nil && res.HasErrors() {
		return errors.CompositeValidationError(res.Errors...)
	}
	return nil
}

// ContentType validates the content type of a request
func validateContentType(allowed []string, actual string) *errors.Validation {
	mt, _, err := mime.ParseMediaType(actual)
	if err != nil {
		return errors.InvalidContentType(actual, allowed)
	}
	if util.ContainsStringsCI(allowed, mt) {
		return nil
	}
	return errors.InvalidContentType(actual, allowed)
}

func validateRequest(ctx *Context, request *http.Request, route *MatchedRoute) *validation {
	validate := &validation{
		context: ctx,
		result:  new(validate.Result),
		request: request,
		route:   route,
		bound:   make(map[string]interface{}),
	}

	validate.contentType()
	validate.responseFormat()
	if validate.result.IsValid() {
		validate.parameters()
	}

	return validate
}

func (v *validation) parameters() {
	result := v.route.Binder.Bind(v.request, v.route.Params, v.route.Consumer, v.bound)
	v.result.Merge(result)
}

func (v *validation) contentType() {
	if httputils.CanHaveBody(v.request.Method) {
		ct, _, err := v.context.ContentType(v.request)
		if err != nil {
			v.result.AddErrors(err)
		} else {
			if err := validateContentType(v.route.Consumes, ct); err != nil {
				v.result.AddErrors(err)
			}
			v.route.Consumer = v.route.Consumers[ct]
		}
	}
}

func (v *validation) responseFormat() {
	if str := v.context.ResponseFormat(v.request, v.route.Produces); str == "" {
		v.result.AddErrors(errors.InvalidResponseFormat(v.request.Header.Get(httputils.HeaderAccept), v.route.Produces))
	}
}
