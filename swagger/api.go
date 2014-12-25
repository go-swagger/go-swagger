package swagger

import (
	"io"
	"net/http"
	"strings"

	"github.com/casualjim/go-swagger"
)

// NewAPI creates the default untyped mux
func NewAPI(spec *swagger.Spec) *API {
	return &API{
		spec: spec,
		consumers: map[string]Consumer{
			"application/json": JSONConsumer(),
		},
		producers: map[string]Producer{
			"application/json": JSONProducer(),
		},
		authHandlers:         make(map[string]AuthHandler),
		registeredOperations: make(map[string]OperationHandler),
	}
}

// API represents an untyped mux for a swagger spec
type API struct {
	spec                 *swagger.Spec
	consumers            map[string]Consumer
	producers            map[string]Producer
	authHandlers         map[string]AuthHandler
	registeredOperations map[string]OperationHandler
}

// Spec returns the swagger spec this untyped mux will serve
func (d *API) Spec() *swagger.Spec {
	return d.spec
}

// RegisterAuth registers an auth handler in this api
func (d *API) RegisterAuth(scheme string, handler AuthHandler) {
	d.authHandlers[strings.ToUpper(scheme)] = handler
}

// RegisterConsumer registers a consumer for a media type.
func (d *API) RegisterConsumer(mediaType string, handler Consumer) {
	d.consumers[strings.ToLower(mediaType)] = handler
}

// RegisterProducer registers a producer for a media type
func (d *API) RegisterProducer(mediaType string, handler Producer) {
	d.producers[strings.ToLower(mediaType)] = handler
}

// RegisterOperation registers an operation handler for an operation name
func (d *API) RegisterOperation(operationID string, handler OperationHandler) {
	d.registeredOperations[operationID] = handler
}

// ValidateWith validates the registrations in this API against the provided spec analyzer
func (d *API) ValidateWith(analyzer *SpecAnalyzer) error {
	var consumes []string
	for k := range d.consumers {
		consumes = append(consumes, k)
	}
	var produces []string
	for k := range d.producers {
		produces = append(produces, k)
	}
	// TODO: implement auth handlers later
	var authHandlers []string
	// for k := range d.authHandlers {
	// 	authHandlers = append(authHandlers, k)
	// }

	var operations []string
	for k := range d.registeredOperations {
		operations = append(operations, k)
	}

	return analyzer.ValidateRegistrations(consumes, produces, authHandlers, operations)
}

// HandlerFunc represents a swagger enabled handler func
type HandlerFunc func(http.ResponseWriter, *http.Request, RouteParams)

// OperationHandler a handler for a swagger operation
type OperationHandler func(interface{}) (interface{}, error)

// ConsumerFunc represents a function that can be used as a consumer
type ConsumerFunc func(io.Reader, interface{}) error

type funcConsumer struct {
	fn ConsumerFunc
}

// Consume consumes the reader into the data parameter
func (f *funcConsumer) Consume(reader io.Reader, data interface{}) error {
	return f.fn(reader, data)
}

// FuncConsumer creates a consumer from a function
func FuncConsumer(fn ConsumerFunc) Consumer {
	return &funcConsumer{fn: fn}
}

// Consumer implementations know how to bind the values on the provided interface to
// data provided by the request body
type Consumer interface {
	// Consume performs the binding of request values
	Consume(io.Reader, interface{}) error
}

// ProducerFunc represents a function that can be used as a producer
type ProducerFunc func(io.Writer, interface{}) error

type funcProducer struct {
	fn ProducerFunc
}

func (f *funcProducer) Produce(writer io.Writer, data interface{}) error {
	return f.fn(writer, data)
}

// FuncProducer creates a producer implemenation from a function
func FuncProducer(fn ProducerFunc) Producer {
	return &funcProducer{fn: fn}
}

// Producer implementations know how to turn the provided interface into a valid
// HTTP response
type Producer interface {
	// Produce writes to the http response
	Produce(io.Writer, interface{}) error
}

// AuthHandler handles authentication for an API
type AuthHandler interface {
	// Authenticate peforms tha authentication
	Authenticate(*http.Request, RouteParams) interface{}
}
