package middleware

import (
	"errors"
	"net/http"

	"github.com/casualjim/go-swagger"
	"github.com/casualjim/go-swagger/httputils"
	"github.com/casualjim/go-swagger/router"
	"github.com/casualjim/go-swagger/spec"
	"github.com/casualjim/go-swagger/swagger-ui"
	"github.com/casualjim/go-swagger/validate"
	"github.com/golang/gddo/httputil"
	"github.com/gorilla/context"
)

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

func defaultStack(context *Context) http.Handler {
	terminator := context.OperationHandlerMiddleware()
	validator := context.ValidationMiddleware(terminator)
	return context.RouterMiddleware(validator)
}

func defaultStackWithUI(context *Context) http.Handler {
	return swaggerui.Middleware("", defaultStack(context))
}

// Serve serves the specified spec with the specified api registrations as a http.Handler
func Serve(spec *spec.Document, api *swagger.API) http.Handler {
	context := NewContext(spec, api, nil)
	return specMiddleware(context, defaultStackWithUI(context))
}

type contextKey int8

const (
	_ contextKey = iota
	ctxContentType
	ctxResponseFormat
	ctxMatchedRoute
	ctxAllowedMethods
	ctxBoundParams

	ctxConsumer
)

type contentTypeValue struct {
	MediaType string
	Charset   string
}

// ContentType gets the parsed value of a content type
func (c *Context) ContentType(request *http.Request) (string, string, *httputils.ParseError) {
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

// RouteInfo tries to match a route for this request
func (c *Context) RouteInfo(request *http.Request) (*router.MatchedRoute, bool) {
	if v, ok := context.GetOk(request, ctxMatchedRoute); ok {
		if val, ok := v.(*router.MatchedRoute); ok {
			return val, ok
		}
	}

	if route, ok := c.router.Lookup(request.Method, request.URL.Path); ok {
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
	// fmt.Fprintf(os.Stdout, "content type %q format %q", ct, format)
	if format != "" {
		context.Set(r, ctxResponseFormat, format)
	}
	return format
}

// AllowedMethods gets the allowed methods for the path of this request
func (c *Context) AllowedMethods(request *http.Request) []string {
	return c.router.OtherMethods(request.Method, request.URL.Path)
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

// Respond renders the response after doing some content negotiation
func (c *Context) Respond(rw http.ResponseWriter, r *http.Request, produces []string, data interface{}) {
	if err, ok := data.(error); ok {
		c.api.ServeError(rw, r, err)
		return
	}

	format := c.ResponseFormat(r, produces)
	producers := c.api.ProducersFor([]string{format})
	prod, ok := producers[format]
	if !ok {
		panic(errors.New("can't find a producer for " + format))
	}
	if err := prod.Produce(rw, data); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
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
