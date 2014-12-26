package swagger

import (
	"net/http"

	"github.com/casualjim/go-swagger"
)

// requestBinder binds the values from the request to the operation parameter struct
type requestBinder struct {
	Consumers map[string]Consumer
	Operation *swagger.Operation
}

// NewRequestBinder creates a new instance of a request binder
func newRequestBinder(operation *swagger.Operation, consumers map[string]Consumer) *requestBinder {
	return &requestBinder{Consumers: consumers, Operation: operation}
}

// Bind binds the request values to the provided struct
func (o *requestBinder) Bind(req *http.Request, routeParams RouteParams, data interface{}) error {
	return nil
}
