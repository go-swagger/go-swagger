package middleware

import (
	"net/http"

	"github.com/casualjim/go-swagger"
	"github.com/casualjim/go-swagger/errors"
	"github.com/casualjim/go-swagger/httputils"
	"github.com/casualjim/go-swagger/spec"
	"github.com/casualjim/go-swagger/swagger-ui"
	"github.com/casualjim/go-swagger/validate"
	"github.com/golang/gddo/httputil"
	"github.com/gorilla/context"
)

// RequestBinder is an interface for types to implement
// when they want to be able to bind from a request
type RequestBinder interface {
	BindRequest(*http.Request, *MatchedRoute, swagger.Consumer) *validate.Result
}

// Context is a type safe wrapper around an untyped request context
// used throughout to store request context with the gorilla context module
type Context struct {
	spec   *spec.Document
	api    *swagger.API
	router Router
}

type routableUntypedAPI struct {
	api      *swagger.API
	handlers map[string]http.Handler
}

func newRoutableUntypedAPI(spec *spec.Document, api *swagger.API, context *Context) *routableUntypedAPI {
	var handlers map[string]http.Handler
	if spec == nil || api == nil {
		return nil
	}
	for _, hls := range spec.Operations() {
		for _, op := range hls {
			schemes := spec.SecurityDefinitionsFor(op)

			if oh, ok := api.OperationHandlerFor(op.ID); ok {
				if handlers == nil {
					handlers = make(map[string]http.Handler)
				}

				handlers[op.ID] = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					// lookup route info in the context
					route, _ := context.RouteInfo(r)

					// bind and validate the request using reflection
					bound, validation := context.BindAndValidate(r, route)
					if validation.HasErrors() {
						context.Respond(w, r, route.Produces, route, validation.Errors[0])
						return
					}

					// actually handle the request
					result, err := oh.Handle(bound)
					if err != nil {
						// respond with failure
						context.Respond(w, r, route.Produces, route, err)
						return
					}

					// respond with success
					context.Respond(w, r, route.Produces, route, result)
				})

				if len(schemes) > 0 {
					handlers[op.ID] = newSecureAPI(context, handlers[op.ID])
				}
			}
		}
	}

	return &routableUntypedAPI{api: api, handlers: handlers}
}

func (r *routableUntypedAPI) HandlerFor(operationID string) (http.Handler, bool) {
	handler, ok := r.handlers[operationID]
	return handler, ok
}
func (r *routableUntypedAPI) ServeErrorFor(operationID string) func(http.ResponseWriter, *http.Request, error) {
	return r.api.ServeError
}
func (r *routableUntypedAPI) ConsumersFor(mediaTypes []string) map[string]swagger.Consumer {
	return r.api.ConsumersFor(mediaTypes)
}
func (r *routableUntypedAPI) ProducersFor(mediaTypes []string) map[string]swagger.Producer {
	return r.api.ProducersFor(mediaTypes)
}
func (r *routableUntypedAPI) AuthenticatorsFor(schemes map[string]spec.SecurityScheme) map[string]swagger.Authenticator {
	return r.api.AuthenticatorsFor(schemes)
}

// NewContext creates a new context wrapper
func NewContext(spec *spec.Document, api *swagger.API, routes Router) *Context {
	ctx := &Context{spec: spec, api: api}
	if routes == nil {
		routableAPI := newRoutableUntypedAPI(spec, api, ctx)
		routes = DefaultRouter(spec, routableAPI)
	}
	ctx.router = routes
	return ctx
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
func (c *Context) BindValidRequest(request *http.Request, route *MatchedRoute, binder RequestBinder) error {
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
func (c *Context) LookupRoute(request *http.Request) (*MatchedRoute, bool) {
	if route, ok := c.router.Lookup(request.Method, request.URL.Path); ok {
		return route, ok
	}
	return nil, false
}

// RouteInfo tries to match a route for this request
func (c *Context) RouteInfo(request *http.Request) (*MatchedRoute, bool) {
	if v, ok := context.GetOk(request, ctxMatchedRoute); ok {
		if val, ok := v.(*MatchedRoute); ok {
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
func (c *Context) Authorize(request *http.Request, route *MatchedRoute) (interface{}, error) {
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
func (c *Context) BindAndValidate(request *http.Request, matched *MatchedRoute) (interface{}, *validate.Result) {
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
func (c *Context) Respond(rw http.ResponseWriter, r *http.Request, produces []string, route *MatchedRoute, data interface{}) {
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

// APIHandler returns a handler to serve
func (c *Context) APIHandler() http.Handler {
	return specMiddleware(c, newRouter(c, newOperationExecutor(c)))
}

// UIMiddleware creates a new swagger UI middleware for this context
func (c *Context) UIMiddleware(handler http.Handler) http.Handler {
	return swaggerui.Middleware("", handler)
}
