package router

import (
	"regexp"
	"strings"

	"github.com/casualjim/go-swagger"
	swagger_api "github.com/casualjim/go-swagger/swagger"
	"github.com/casualjim/go-swagger/swagger/spec"
	"github.com/naoina/denco"
	"github.com/wsxiaoys/terminal/color"
)

// Router represents a swagger aware router
type Router interface {
	Lookup(method, path string) (*MatchedRoute, bool)
	OtherMethods(method, path string) []string
}

type defaultRouteBuilder struct {
	spec    *spec.Document
	api     *swagger_api.API
	records map[string][]denco.Record
}

type defaultRouter struct {
	spec    *spec.Document
	api     *swagger_api.API
	routers map[string]*denco.Router
}

func newDefaultRouteBuilder(spec *spec.Document, api *swagger_api.API) *defaultRouteBuilder {
	return &defaultRouteBuilder{
		spec:    spec,
		api:     api,
		records: make(map[string][]denco.Record),
	}
}

// Default dreates a default implemenation of the router
func Default(spec *spec.Document, api *swagger_api.API) Router {
	builder := newDefaultRouteBuilder(spec, api)
	if spec != nil {
		for method, paths := range spec.Operations() {
			for path, operation := range paths {
				builder.AddRoute(method, spec.BasePath()+path, operation)
			}
		}
	}
	return builder.Build()
}

type routeEntry struct {
	PathPattern string
	BasePath    string
	Operation   *swagger.Operation
	Consumes    []string
	Consumers   map[string]swagger_api.Consumer
	Produces    []string
	Producers   map[string]swagger_api.Producer
	Parameters  map[string]swagger.Parameter
	Handler     swagger_api.OperationHandler
}

// MatchedRoute represents the route that was matched in this request
type MatchedRoute struct {
	routeEntry
	Params swagger_api.RouteParams
}

func (d *defaultRouter) Lookup(method, path string) (*MatchedRoute, bool) {
	if router, ok := d.routers[strings.ToUpper(method)]; ok {
		if m, rp, ok := router.Lookup(path); ok && m != nil {
			if entry, ok := m.(*routeEntry); ok {
				var params swagger_api.RouteParams
				for _, p := range rp {
					params = append(params, swagger_api.RouteParam{Name: p.Name, Value: p.Value})
				}
				return &MatchedRoute{routeEntry: *entry, Params: params}, true
			}
		}
	}
	return nil, false
}

func (d *defaultRouter) OtherMethods(method, path string) []string {
	mn := strings.ToUpper(method)
	var methods []string
	for k, v := range d.routers {
		if k != mn {
			if _, _, ok := v.Lookup(path); ok {
				methods = append(methods, k)
				continue
			}
		}
	}
	return methods
}

var pathConverter = regexp.MustCompile(`{(\w+)}`)

func (d *defaultRouteBuilder) AddRoute(method, path string, operation *swagger.Operation) {
	mn := strings.ToUpper(method)

	if handler, ok := d.api.OperationHandlerFor(operation.ID); ok {
		consumes := d.spec.ConsumesFor(operation)
		produces := d.spec.ProducesFor(operation)

		record := denco.NewRecord(pathConverter.ReplaceAllString(path, ":$1"), &routeEntry{
			Operation:  operation,
			Handler:    handler,
			Consumes:   consumes,
			Produces:   produces,
			Consumers:  d.api.ConsumersFor(consumes),
			Producers:  d.api.ProducersFor(produces),
			Parameters: d.spec.ParamsFor(method, path),
		})
		color.Printf("registered route @{c}%s@{|}\t@{y}%q@{|}\t@{m}%s@{|}\n", mn, path, operation.ID)
		d.records[mn] = append(d.records[mn], record)
	}
}

func (d *defaultRouteBuilder) Build() *defaultRouter {
	routers := make(map[string]*denco.Router)
	for method, records := range d.records {
		router := denco.New()
		router.Build(records)
		routers[method] = router
	}
	return &defaultRouter{
		spec:    d.spec,
		routers: routers,
	}
}
