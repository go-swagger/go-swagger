package router

import (
	"regexp"
	"strings"

	"github.com/casualjim/go-swagger"
	"github.com/casualjim/go-swagger/spec"
	"github.com/casualjim/go-swagger/validate"
	"github.com/naoina/denco"
)

// Router represents a swagger aware router
type Router interface {
	Lookup(method, path string) (*MatchedRoute, bool)
	OtherMethods(method, path string) []string
}

type defaultRouteBuilder struct {
	spec    *spec.Document
	api     *swagger.API
	records map[string][]denco.Record
}

type defaultRouter struct {
	spec    *spec.Document
	api     *swagger.API
	routers map[string]*denco.Router
}

func newDefaultRouteBuilder(spec *spec.Document, api *swagger.API) *defaultRouteBuilder {
	return &defaultRouteBuilder{
		spec:    spec,
		api:     api,
		records: make(map[string][]denco.Record),
	}
}

// Default creates a default implemenation of the router
func Default(spec *spec.Document, api *swagger.API) Router {
	builder := newDefaultRouteBuilder(spec, api)
	if spec != nil {
		for method, paths := range spec.Operations() {
			for path, operation := range paths {
				builder.AddRoute(method, path, operation)
			}
		}
	}
	return builder.Build()
}

type routeEntry struct {
	PathPattern string
	BasePath    string
	Operation   *spec.Operation
	Consumes    []string
	Consumers   map[string]swagger.Consumer
	Produces    []string
	Producers   map[string]swagger.Producer
	Parameters  map[string]spec.Parameter
	Handler     swagger.OperationHandler
	Binder      *validate.RequestBinder
}

// MatchedRoute represents the route that was matched in this request
type MatchedRoute struct {
	routeEntry
	Params swagger.RouteParams
}

func (d *defaultRouter) Lookup(method, path string) (*MatchedRoute, bool) {
	if router, ok := d.routers[strings.ToUpper(method)]; ok {
		if m, rp, ok := router.Lookup(path); ok && m != nil {
			if entry, ok := m.(*routeEntry); ok {
				var params swagger.RouteParams
				for _, p := range rp {
					params = append(params, swagger.RouteParam{Name: p.Name, Value: p.Value})
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

func (d *defaultRouteBuilder) AddRoute(method, path string, operation *spec.Operation) {
	mn := strings.ToUpper(method)

	if handler, ok := d.api.OperationHandlerFor(operation.ID); ok {
		consumes := d.spec.ConsumesFor(operation)
		produces := d.spec.ProducesFor(operation)
		parameters := d.spec.ParamsFor(method, path)

		record := denco.NewRecord(pathConverter.ReplaceAllString(path, ":$1"), &routeEntry{
			Operation:  operation,
			Handler:    handler,
			Consumes:   consumes,
			Produces:   produces,
			Consumers:  d.api.ConsumersFor(consumes),
			Producers:  d.api.ProducersFor(produces),
			Parameters: parameters,
			Binder:     validate.NewRequestBinder(parameters, d.spec.Spec()),
		})
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
