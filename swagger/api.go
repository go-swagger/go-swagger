package swagger

import (
	"io"
	"net/http"
	"strings"

	"github.com/casualjim/go-swagger"
)

// NewAPI creates the default untyped API
func NewAPI(spec *swagger.Spec) *API {
	return &API{
		analyzer: newAnalyzer(spec),
		consumers: map[string]Consumer{
			"application/json": JSONConsumer(),
		},
		producers: map[string]Producer{
			"application/json": JSONProducer(),
		},
		authHandlers: make(map[string]AuthHandler),
		operations:   make(map[string]OperationHandler),
	}
}

// API represents an untyped mux for a swagger spec
type API struct {
	analyzer     *specAnalyzer
	consumers    map[string]Consumer
	producers    map[string]Producer
	authHandlers map[string]AuthHandler
	operations   map[string]OperationHandler
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
	d.operations[operationID] = handler
}

// OperationHandlerFor returns the operation handler for the specified id if it can be found
func (d *API) OperationHandlerFor(operationID string) (OperationHandler, bool) {
	h, ok := d.operations[operationID]
	return h, ok
}

// ConsumersFor gets the consumers for the specified media types
func (d *API) ConsumersFor(mediaTypes []string) map[string]Consumer {
	result := make(map[string]Consumer)
	for _, mt := range mediaTypes {
		if consumer, ok := d.consumers[mt]; ok {
			result[mt] = consumer
		}
	}
	return result
}

// ProducersFor gets the producers for the specified media types
func (d *API) ProducersFor(mediaTypes []string) map[string]Producer {
	result := make(map[string]Producer)
	for _, mt := range mediaTypes {
		if producer, ok := d.producers[mt]; ok {
			result[mt] = producer
		}
	}
	return result
}

// Validate validates this API for any missing items
func (d *API) Validate() error {
	return d.validateWith(d.analyzer)
}

// Handler takes the untyped API and a router with those it will validate the registrations in the API.
// When everything is found to be valid it will build the http.Handler with the provided router
//
// If there are missing consumers for registered media types it will return an error
// If there are missing producers for registered media types it will return an error
// If there are missing auth handlers for registered security schemes it will return an error
// If there are missing operation handlers for operationIds it will return an error
func (d *API) Handler(router Router) (http.Handler, error) {
	if router == nil {
		router = DefaultRouter()
	}

	if err := d.Validate(); err != nil {
		return nil, err
	}

	initializer := &routerInitializer{
		API:      d,
		Router:   router,
		Analyzer: d.analyzer,
	}

	return initializer.Initialize()
}

// validateWith validates the registrations in this API against the provided spec analyzer
func (d *API) validateWith(analyzer *specAnalyzer) error {
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
	for k := range d.operations {
		operations = append(operations, k)
	}

	return analyzer.ValidateRegistrations(consumes, produces, authHandlers, operations)
}

// HandlerFunc represents a swagger enabled handler func
type HandlerFunc func(http.ResponseWriter, *http.Request, RouteParams)

// OperationHandler a handler for a swagger operation
type OperationHandler interface {
	ParameterModel() interface{}
	Handle(interface{}) (interface{}, error)
}

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
	// Authenticate peforms the authentication
	Authenticate(*http.Request) interface{}
}
