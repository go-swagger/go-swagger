package swagger

import (
	"net/http"
	"regexp"

	"github.com/casualjim/go-swagger"
	"github.com/casualjim/go-swagger/swagger/spec"
	"github.com/naoina/denco"
)

// RouteParam is a object to capture route params in a framework agnostic way.
// implementations of the muxer should use these route params to communicate with the
// swagger framework
type RouteParam struct {
	Name  string
	Value string
}

// RouteParams the collection of route params
type RouteParams []RouteParam

// Get gets the value for the route param for the specified key
func (r RouteParams) Get(name string) string {
	for _, p := range r {
		if p.Name == name {
			return p.Value
		}
	}
	return ""
}

// Router implementations provide the integration with url routers
type Router interface {
	AddRoute(string, string, HandlerFunc)
	Build() (http.Handler, error)
}

// DefaultRouter creates a new denco route registrar
func DefaultRouter() Router {
	return &defaultRouter{}
}

type defaultRouter struct {
	handlers []denco.Handler
}

func (d *defaultRouter) Build() (http.Handler, error) {
	mux := denco.NewMux()

	return mux.Build(d.handlers)
}

var pathConverter = regexp.MustCompile(`{(\w+)}`)

func (d *defaultRouter) AddRoute(method, path string, handler HandlerFunc) {
	d.handlers = append(d.handlers, denco.Handler{
		Method: method,
		Path:   d.convertPathPattern(path),
		Func:   d.wrapHandlerFunc(handler),
	})
}

func (d *defaultRouter) convertPathPattern(path string) string {
	converted := pathConverter.ReplaceAllString(path, ":$1")
	ln := len(converted)
	if ln > 1 && converted[ln-1] == '/' {
		return converted[:ln-1]
	}
	return converted
}

func (d *defaultRouter) wrapHandlerFunc(handler HandlerFunc) denco.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request, p denco.Params) {
		routeParams := make([]RouteParam, 0, len(p))
		for _, param := range p {
			routeParams = append(routeParams, RouteParam{Name: param.Name, Value: param.Value})
		}
		handler(rw, r, routeParams)
	}
}

type routerInitializer struct {
	Router Router
	Spec   *spec.Document
	API    *API
}

func (r *routerInitializer) Initialize() (http.Handler, error) {
	for path, pathItem := range r.Spec.AllPaths() {
		r.registerRoute("GET", path, pathItem.Get)
		r.registerRoute("HEAD", path, pathItem.Head)
		r.registerRoute("OPTIONS", path, pathItem.Options)
		r.registerRoute("POST", path, pathItem.Post)
		r.registerRoute("PUT", path, pathItem.Put)
		r.registerRoute("PATCH", path, pathItem.Patch)
		r.registerRoute("DELETE", path, pathItem.Delete)
	}
	return r.Router.Build()
}

func (r *routerInitializer) registerRoute(method, path string, operation *swagger.Operation) {
	if operation != nil {
		if h, ok := r.API.OperationHandlerFor(operation.ID); ok {
			r.Router.AddRoute(method, path, createHandler(r.API, operation, h, r.Spec).Handle)
		}
	}
}
