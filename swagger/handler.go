package swagger

import (
	"net/http"

	"github.com/casualjim/go-swagger"
)

type handler struct {
	Consumes  []string
	Consumers map[string]Consumer
	Produces  []string
	Producers map[string]Producer
	Operation *swagger.Operation
	Handler   OperationHandler
	Analyzer  *specAnalyzer
	Binder    *requestBinder
}

func createHandler(api *API, operation *swagger.Operation, h OperationHandler, analyzer *specAnalyzer) *handler {
	consumes := analyzer.ConsumesFor(operation)
	consumers := api.ConsumersFor(consumes)

	return &handler{
		Consumes:  consumes,
		Consumers: consumers,
		Produces:  nil,
		Producers: nil,
		Operation: operation,
		Handler:   h,
		Analyzer:  analyzer,
		Binder:    newRequestBinder(operation, consumers),
	}
}

func (h *handler) Handle(rw http.ResponseWriter, req *http.Request, routeParams RouteParams) {
	rw.WriteHeader(http.StatusNotImplemented)
}
