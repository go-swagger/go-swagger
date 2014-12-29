package swagger

// type handler struct {
// 	Consumes   []string
// 	Consumers  map[string]Consumer
// 	Produces   []string
// 	Producers  map[string]Producer
// 	Operation  *swagger.Operation
// 	Handler    OperationHandler
// 	Spec       *spec.Document
// 	ParamNames map[string]string
// 	Binder     *operationBinder
// }

// func createHandler(api *API, operation *swagger.Operation, h OperationHandler, spec *spec.Document) *handler {
// 	consumes := spec.ConsumesFor(operation)
// 	consumers := api.ConsumersFor(consumes)
// 	produces := spec.ProducesFor(operation)
// 	producers := api.ProducersFor(produces)
// 	parameters := spec.ParametersFor(operation)

// 	return &handler{
// 		Consumes:  consumes,
// 		Consumers: consumers,
// 		Produces:  produces,
// 		Producers: producers,
// 		Operation: operation,
// 		Handler:   h,
// 		Spec:      spec,
// 		Binder:    &operationBinder{parameters, consumers},
// 	}
// }

// func (h *handler) Handle(rw http.ResponseWriter, req *http.Request, routeParams RouteParams) {
// 	rw.WriteHeader(http.StatusNotImplemented)
// }
