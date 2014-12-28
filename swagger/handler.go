package swagger

import (
	"net/http"

	"github.com/casualjim/go-swagger"
)

type handler struct {
	Consumes   []string
	Consumers  map[string]Consumer
	Produces   []string
	Producers  map[string]Producer
	Operation  *swagger.Operation
	Handler    OperationHandler
	Analyzer   *specAnalyzer
	ParamNames map[string]string
	Binder     *operationBinder
}

func createHandler(api *API, operation *swagger.Operation, h OperationHandler, analyzer *specAnalyzer) *handler {
	consumes := analyzer.ConsumesFor(operation)
	consumers := api.ConsumersFor(consumes)
	produces := analyzer.ProducesFor(operation)
	producers := api.ProducersFor(produces)
	parameters := analyzer.ParametersFor(operation)

	return &handler{
		Consumes:  consumes,
		Consumers: consumers,
		Produces:  produces,
		Producers: producers,
		Operation: operation,
		Handler:   h,
		Analyzer:  analyzer,
		Binder:    &operationBinder{parameters, consumers},
	}
}

func (h *handler) Handle(rw http.ResponseWriter, req *http.Request, routeParams RouteParams) {
	// authenticate

	// perform non-body validations

	// create new instance
	parameters := h.Handler.ParameterModel()
	if err := h.Binder.Bind(req, routeParams, parameters); err != nil {
		// use renderer to render an error with the appropriate status code
		//h.Renderer(rw, err)
		return
	}

	// validate body

	// execute

	// render
	rw.WriteHeader(http.StatusNotImplemented)
}
