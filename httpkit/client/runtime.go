package client

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/go-swagger/go-swagger/client"
	"github.com/go-swagger/go-swagger/httpkit"
	"github.com/go-swagger/go-swagger/spec"
	"github.com/go-swagger/go-swagger/strfmt"
)

// Runtime represents an API client that uses the transport
// to make http requests based on a swagger specification.
type Runtime struct {
	DefaultMediaType string
	Consumers        map[string]httpkit.Consumer
	Producers        map[string]httpkit.Producer
	Transport        http.RoundTripper
	Spec             *spec.Document
	Host             string
	BasePath         string
	Formats          strfmt.Registry

	client          *http.Client
	methodsAndPaths map[string]methodAndPath
}

// New creates a new default runtime for a swagger api client.
func New(swaggerSpec *spec.Document) *Runtime {
	var rt Runtime
	rt.DefaultMediaType = httpkit.JSONMime
	// TODO: actually infer this stuff from the spec
	rt.Consumers = map[string]httpkit.Consumer{
		httpkit.JSONMime: httpkit.JSONConsumer(),
	}
	rt.Producers = map[string]httpkit.Producer{
		httpkit.JSONMime: httpkit.JSONProducer(),
	}
	rt.Spec = swaggerSpec
	rt.Transport = http.DefaultTransport
	rt.client = http.DefaultClient
	rt.client.Transport = rt.Transport
	rt.Host = swaggerSpec.Host()
	rt.BasePath = swaggerSpec.BasePath()
	rt.methodsAndPaths = make(map[string]methodAndPath)
	for mth, pathItem := range rt.Spec.Operations() {
		for pth, op := range pathItem {
			rt.methodsAndPaths[op.ID] = methodAndPath{mth, pth}
		}
	}
	return &rt
}

// Submit a request and when there is a body on success it will turn that into the result
// all other things are turned into an api error for swagger which retains the status code
func (r *Runtime) Submit(operationID string, params client.RequestWriter, readResponse client.ResponseReader) (interface{}, error) {
	mthPth, ok := r.methodsAndPaths[operationID]
	if !ok {
		return nil, fmt.Errorf("unknown operation: %q", operationID)
	}
	request, err := newRequest(mthPth.Method, mthPth.PathPattern, params)
	if err != nil {
		return nil, err
	}

	request.SetHeaderParam(httpkit.HeaderContentType, r.DefaultMediaType)
	var accept []string
	for k := range r.Consumers {
		accept = append(accept, k)
	}
	request.SetHeaderParam(httpkit.HeaderAccept, accept...)

	req, err := request.BuildHTTP(r.Producers[r.DefaultMediaType], r.Formats)
	// TODO: work out scheme based on the operations and the default scheme
	req.URL.Scheme = "http"
	req.URL.Host = r.Host
	req.URL.Path = filepath.Join(r.BasePath, req.URL.Host)
	if err != nil {
		return nil, err
	}

	res, err := r.client.Do(req) // make requests, by default follows 10 redirects before failing
	if err != nil {
		return nil, err
	}
	ct := res.Header.Get(httpkit.HeaderContentType)
	if ct == "" { // this should really really never occur
		ct = r.DefaultMediaType
	}

	// TODO: normalize this (ct) and only match on media type,
	// skip the params like charset unless a tie breaker is needed
	cons, ok := r.Consumers[ct]
	if !ok {
		// scream about not knowing what to do
		return nil, fmt.Errorf("no consumer: %q", ct)
	}
	return readResponse.ReadResponse(response{res}, cons)
}
