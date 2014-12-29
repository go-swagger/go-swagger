package middleware

import (
	"net/http"

	"github.com/casualjim/go-swagger/swagger"
	"github.com/casualjim/go-swagger/swagger/httputils"
	"github.com/casualjim/go-swagger/swagger/router"
	"github.com/casualjim/go-swagger/swagger/spec"
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
func NewContext(spec *spec.Document, api *swagger.API) *Context {
	return &Context{spec: spec, api: api, router: router.Default(spec, api)}
}

type contextKey int8

const (
	_ contextKey = iota
	ctxContentType
	ctxMatchedRoute
	ctxAllowedMethods
	ctxConsumer
)

type contentTypeValue struct {
	MediaType string
	Charset   string
}

// ContentType gets the parsed value of a content type
func (c *Context) ContentType(request *http.Request) (string, string, error) {
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

// AllowedMethods gets the allowed methods for the path of this request
func (c *Context) AllowedMethods(request *http.Request) []string {
	return c.router.OtherMethods(request.Method, request.URL.Path)
}

// Respond renders the response after doing some content negotiation
func (c *Context) Respond(rw http.ResponseWriter, r *http.Request, data interface{}) {
	if err, ok := data.(error); ok {
		c.api.ServeError(rw, r, err)
		return
	}
	// perform content negotiation
	// pick a producer
	// render content with producer
}
