package middleware

import (
	"net/http"

	"github.com/casualjim/go-swagger"
	"github.com/casualjim/go-swagger/errors"
	"github.com/casualjim/go-swagger/httputils"
	"github.com/casualjim/go-swagger/router"
	"github.com/casualjim/go-swagger/spec"
	"github.com/casualjim/go-swagger/swagger-ui"
	"github.com/casualjim/go-swagger/validate"
	"github.com/golang/gddo/httputil"
	"github.com/gorilla/context"
)

// RequestBinder is an interface for types to implement
// when they want to be able to bind from a request
type RequestBinder interface {
	BindRequest(*http.Request, *router.MatchedRoute, swagger.Consumer) *validate.Result
}

// Context is a type safe wrapper around an untyped request context
// used throughout to store request context with the gorilla context module
type Context struct {
	spec   *spec.Document
	api    *swagger.API
	router router.Router
}

// NewContext creates a new context wrapper
func NewContext(spec *spec.Document, api *swagger.API, routes router.Router) *Context {
	if routes == nil {
		routes = router.Default(spec, api)
	}
	return &Context{spec: spec, api: api, router: routes}
}

// Serve serves the specified spec with the specified api registrations as a http.Handler
func Serve(spec *spec.Document, api *swagger.API) http.Handler {
	context := NewContext(spec, api, nil)
	return context.APIHandler()
}

// ServeWithUI serves the specified spec with the specified api registrations as a http.Handler
// it also enables the swagger docs ui on /swagger-ui
func ServeWithUI(spec *spec.Document, api *swagger.API) http.Handler {
	context := NewContext(spec, api, nil)
	return context.UIMiddleware(context.APIHandler())
}

type contextKey int8

const (
	_ contextKey = iota
	ctxContentType
	ctxResponseFormat
	ctxMatchedRoute
	ctxAllowedMethods
	ctxBoundParams
	ctxSecurityPrincipal

	ctxConsumer
)

type contentTypeValue struct {
	MediaType string
	Charset   string
}

// BasePath returns the base path for this API
func (c *Context) BasePath() string {
	return c.spec.BasePath()
}

// RequiredProduces returns the accepted content types for responses
func (c *Context) RequiredProduces() []string {
	return c.spec.RequiredProduces()
}

// BindValidRequest binds a params object to a request but only when the request is valid
// if the request is not valid an error will be returned
func (c *Context) BindValidRequest(request *http.Request, route *router.MatchedRoute, binder RequestBinder) error {
	res := new(validate.Result)
	var consumer swagger.Consumer

	// check and validate content type, select consumer
	if httputils.CanHaveBody(request.Method) {
		ct, _, err := httputils.ContentType(request.Header)
		if err != nil {
			res.AddErrors(err)
		} else {
			if err := validate.ContentType(route.Consumes, ct); err != nil {
				res.AddErrors(err)
			}
			consumer = route.Consumers[ct]
		}
	}

	// check and validate the response format
	if res.IsValid() {
		if str := httputil.NegotiateContentType(request, route.Produces, ""); str == "" {
			res.AddErrors(errors.InvalidResponseFormat(request.Header.Get(httputils.HeaderAccept), route.Produces))
		}
	}

	// now bind the request with the provided binder
	// it's assumed the binder will also validate the request and return an error if the
	// request is invalid
	if res.IsValid() {
		res.Merge(binder.BindRequest(request, route, consumer))
	}

	if res.HasErrors() {
		return errors.CompositeValidationError(res.Errors...)
	}
	return nil
}

// ContentType gets the parsed value of a content type
func (c *Context) ContentType(request *http.Request) (string, string, *errors.ParseError) {
	if v, ok := context.GetOk(request, ctxContentType); ok {
		if val, ok := v.(*contentTypeValue); ok {
			return val.MediaType, val.Charset, nil
		}
	}

	mt, cs, err := httputils.ContentType(request.Header)
	if err != nil {
		return "", "", err
	}
	context.Set(request, ctxContentType, &contentTypeValue{mt, cs})
	return mt, cs, nil
}

// LookupRoute looks a route up and returns true when it is found
func (c *Context) LookupRoute(request *http.Request) (*router.MatchedRoute, bool) {
	if route, ok := c.router.Lookup(request.Method, request.URL.Path); ok {
		return route, ok
	}
	return nil, false
}

// RouteInfo tries to match a route for this request
func (c *Context) RouteInfo(request *http.Request) (*router.MatchedRoute, bool) {
	if v, ok := context.GetOk(request, ctxMatchedRoute); ok {
		if val, ok := v.(*router.MatchedRoute); ok {
			return val, ok
		}
	}
	if route, ok := c.LookupRoute(request); ok {
		context.Set(request, ctxMatchedRoute, route)
		return route, ok
	}

	return nil, false
}

// ResponseFormat negotiates the response content type
func (c *Context) ResponseFormat(r *http.Request, offers []string) string {
	if v, ok := context.GetOk(r, ctxResponseFormat); ok {
		if val, ok := v.(string); ok {
			return val
		}
	}

	format := httputil.NegotiateContentType(r, offers, "")
	if format != "" {
		context.Set(r, ctxResponseFormat, format)
	}
	return format
}

// AllowedMethods gets the allowed methods for the path of this request
func (c *Context) AllowedMethods(request *http.Request) []string {
	return c.router.OtherMethods(request.Method, request.URL.Path)
}

// Authorize authorizes the request
func (c *Context) Authorize(request *http.Request, route *router.MatchedRoute) (interface{}, error) {
	if len(route.Authenticators) == 0 {
		return nil, nil
	}
	if v, ok := context.GetOk(request, ctxSecurityPrincipal); ok {
		return v, nil
	}
	for _, authenticator := range route.Authenticators {
		applies, usr, err := authenticator.Authenticate(request)
		if !applies || err != nil || usr == nil {
			continue
		}
		context.Set(request, ctxSecurityPrincipal, usr)
		return usr, nil
	}
	return nil, errors.Unauthenticated("invalid credentials")
}

// BindAndValidate binds and validates the request
func (c *Context) BindAndValidate(request *http.Request, matched *router.MatchedRoute) (interface{}, *validate.Result) {
	if v, ok := context.GetOk(request, ctxBoundParams); ok {
		if val, ok := v.(*validation); ok {
			return val.bound, val.result
		}
	}
	result := validateRequest(c, request, matched)
	if result != nil {
		context.Set(request, ctxBoundParams, result)
	}
	return result.bound, result.result
}

// NotFound the default not found responder for when no route has been matched yet
func (c *Context) NotFound(rw http.ResponseWriter, r *http.Request) {
	c.Respond(rw, r, []string{httputils.JSONMime}, nil, errors.NotFound("not found"))
}

// Respond renders the response after doing some content negotiation
func (c *Context) Respond(rw http.ResponseWriter, r *http.Request, produces []string, route *router.MatchedRoute, data interface{}) {
	format := c.ResponseFormat(r, produces)
	rw.Header().Set(httputils.HeaderContentType, format)

	if err, ok := data.(error); ok {
		if format == "" {
			rw.Header().Set(httputils.HeaderContentType, httputils.JSONMime)
		}
		c.api.ServeError(rw, r, err)
		return
	}
	if route == nil || route.Operation == nil {
		rw.WriteHeader(200)
		producers := c.api.ProducersFor(produces)
		prod, ok := producers[format]
		if !ok {
			panic(errors.New(http.StatusInternalServerError, "can't find a producer for "+format))
		}
		if err := prod.Produce(rw, data); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
		return
	}
	if _, code, ok := route.Operation.SuccessResponse(); ok {
		if code == 201 || code == 204 {
			rw.WriteHeader(code)
			return
		}

		rw.WriteHeader(code)
		producers := route.Producers
		prod, ok := producers[format]
		if !ok {
			panic(errors.New(http.StatusInternalServerError, "can't find a producer for "+format))
		}
		if err := prod.Produce(rw, data); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
		return
	}
	c.api.ServeError(rw, r, errors.New(http.StatusInternalServerError, "can't produce response"))
}

// SpecMiddleware generates a middleware for serving the swagger spec document at /swagger.json
func (c *Context) SpecMiddleware(handler http.Handler) http.Handler {
	return specMiddleware(c, handler)
}

// RouterMiddleware creates a new router middleware for this context
func (c *Context) RouterMiddleware(handler http.Handler) http.Handler {
	return newRouter(c, handler)
}

// ValidationMiddleware creates a new validation middleware for this context
func (c *Context) ValidationMiddleware(handler http.Handler) http.Handler {
	return newValidation(c, handler)
}

// OperationHandlerMiddleware creates a terminating http handler
func (c *Context) OperationHandlerMiddleware() http.Handler {
	return newOperationExecutor(c)
}

// APIHandler returns a handler to serve
func (c *Context) APIHandler() http.Handler {
	return c.SpecMiddleware(c.DefaultMiddlewares())
}

// UIMiddleware creates a new swagger UI middleware for this context
func (c *Context) UIMiddleware(handler http.Handler) http.Handler {
	return swaggerui.Middleware("", handler)
}

// SecurityMiddleware creates the middleware to provide security for the API
func (c *Context) SecurityMiddleware(handler http.Handler) http.Handler {
	return newSecureAPI(c, handler)
}

// DefaultMiddlewares generates the default middleware handler stack
func (c *Context) DefaultMiddlewares() http.Handler {
	terminator := c.OperationHandlerMiddleware()
	validator := c.ValidationMiddleware(terminator)
	secured := c.SecurityMiddleware(validator)
	return c.RouterMiddleware(secured)
}
