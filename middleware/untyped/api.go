package untyped

import (
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/casualjim/go-swagger"
	"github.com/casualjim/go-swagger/errors"
	"github.com/casualjim/go-swagger/spec"
	"github.com/casualjim/go-swagger/strfmt"
)

// NewAPI creates the default untyped API
func NewAPI(spec *spec.Document) *API {
	return &API{
		spec:            spec,
		DefaultProduces: "application/json",
		DefaultConsumes: "application/json",
		consumers: map[string]swagger.Consumer{
			"application/json": swagger.JSONConsumer(),
		},
		producers: map[string]swagger.Producer{
			"application/json": swagger.JSONProducer(),
		},
		authenticators: make(map[string]swagger.Authenticator),
		operations:     make(map[string]swagger.OperationHandler),
		ServeError:     errors.ServeError,
		Models:         make(map[string]func() interface{}),
		formats:        strfmt.NewFormats(),
	}
}

// API represents an untyped mux for a swagger spec
type API struct {
	spec            *spec.Document
	DefaultProduces string
	DefaultConsumes string
	consumers       map[string]swagger.Consumer
	producers       map[string]swagger.Producer
	authenticators  map[string]swagger.Authenticator
	operations      map[string]swagger.OperationHandler
	ServeError      func(http.ResponseWriter, *http.Request, error)
	Models          map[string]func() interface{}
	formats         strfmt.Registry
}

// Formats returns the registered string formats
func (d *API) Formats() strfmt.Registry {
	return d.formats
}

// RegisterFormat registers a custom format validator
func (d *API) RegisterFormat(name string, format strfmt.Format, validator strfmt.Validator) {
	d.formats.Add(name, format, validator)
}

// RegisterAuth registers an auth handler in this api
func (d *API) RegisterAuth(scheme string, handler swagger.Authenticator) {
	d.authenticators[scheme] = handler
}

// RegisterConsumer registers a consumer for a media type.
func (d *API) RegisterConsumer(mediaType string, handler swagger.Consumer) {
	d.consumers[strings.ToLower(mediaType)] = handler
}

// RegisterProducer registers a producer for a media type
func (d *API) RegisterProducer(mediaType string, handler swagger.Producer) {
	d.producers[strings.ToLower(mediaType)] = handler
}

// RegisterOperation registers an operation handler for an operation name
func (d *API) RegisterOperation(operationID string, handler swagger.OperationHandler) {
	d.operations[operationID] = handler
}

// OperationHandlerFor returns the operation handler for the specified id if it can be found
func (d *API) OperationHandlerFor(operationID string) (swagger.OperationHandler, bool) {
	h, ok := d.operations[operationID]
	return h, ok
}

// ConsumersFor gets the consumers for the specified media types
func (d *API) ConsumersFor(mediaTypes []string) map[string]swagger.Consumer {
	result := make(map[string]swagger.Consumer)
	for _, mt := range mediaTypes {
		if consumer, ok := d.consumers[mt]; ok {
			result[mt] = consumer
		}
	}
	return result
}

// ProducersFor gets the producers for the specified media types
func (d *API) ProducersFor(mediaTypes []string) map[string]swagger.Producer {
	result := make(map[string]swagger.Producer)
	for _, mt := range mediaTypes {
		if producer, ok := d.producers[mt]; ok {
			result[mt] = producer
		}
	}
	return result
}

// AuthenticatorsFor gets the authenticators for the specified security schemes
func (d *API) AuthenticatorsFor(schemes map[string]spec.SecurityScheme) map[string]swagger.Authenticator {
	result := make(map[string]swagger.Authenticator)
	for k := range schemes {
		if a, ok := d.authenticators[k]; ok {
			result[k] = a
		}
	}
	return result
}

// Validate validates this API for any missing items
func (d *API) Validate() error {
	return d.validate()
}

// validateWith validates the registrations in this API against the provided spec analyzer
func (d *API) validate() error {
	var consumes []string
	for k := range d.consumers {
		consumes = append(consumes, k)
	}

	var produces []string
	for k := range d.producers {
		produces = append(produces, k)
	}

	var authenticators []string
	for k := range d.authenticators {
		authenticators = append(authenticators, k)
	}

	var operations []string
	for k := range d.operations {
		operations = append(operations, k)
	}

	var definedAuths []string
	for k := range d.spec.Spec().SecurityDefinitions {
		definedAuths = append(definedAuths, k)
	}

	if err := d.verify("consumes", consumes, d.spec.RequiredConsumes()); err != nil {
		return err
	}
	if err := d.verify("produces", produces, d.spec.RequiredProduces()); err != nil {
		return err
	}
	if err := d.verify("operation", operations, d.spec.OperationIDs()); err != nil {
		return err
	}

	requiredAuths := d.spec.RequiredSchemes()
	if err := d.verify("auth scheme", authenticators, requiredAuths); err != nil {
		return err
	}
	fmt.Printf("comparing %s with %s\n", strings.Join(definedAuths, ","), strings.Join(requiredAuths, ","))
	if err := d.verify("security definitions", definedAuths, requiredAuths); err != nil {
		return err
	}
	return nil
}

func (d *API) verify(name string, registrations []string, expectations []string) error {

	sort.Sort(sort.StringSlice(registrations))
	sort.Sort(sort.StringSlice(expectations))

	expected := map[string]struct{}{}
	seen := map[string]struct{}{}

	for _, v := range expectations {
		expected[v] = struct{}{}
	}

	var unspecified []string
	for _, v := range registrations {
		seen[v] = struct{}{}
		if _, ok := expected[v]; !ok {
			unspecified = append(unspecified, v)
		}
	}

	for k := range seen {
		delete(expected, k)
	}

	var unregistered []string
	for k := range expected {
		unregistered = append(unregistered, k)
	}
	sort.Sort(sort.StringSlice(unspecified))
	sort.Sort(sort.StringSlice(unregistered))

	if len(unregistered) > 0 || len(unspecified) > 0 {
		return &errors.APIVerificationFailed{
			Section:              name,
			MissingSpecification: unspecified,
			MissingRegistration:  unregistered,
		}
	}

	return nil
}
