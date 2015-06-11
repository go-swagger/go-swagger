package client

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"

	"github.com/casualjim/go-swagger/client"
	"github.com/casualjim/go-swagger/httpkit"
	"github.com/casualjim/go-swagger/spec"
	"github.com/casualjim/go-swagger/strfmt"
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

	client          *http.Client
	Formats         strfmt.Registry
	methodsAndPaths map[string]methodAndPath
}

// A ResultFactory creates a new instance of a result for a given status code.
// when this is a response without a body the bool will be true
// This needs to produce the result or it loses the type information.
// That's the explanation for the somewhat many args to this function
type ResultFactory func(int, io.Reader, httpkit.Consumer) (interface{}, error)

// A ResponseConsumer is responsible for turning a response into a proper return
// for a swagger operation.
type ResponseConsumer interface {
	Consume(*http.Response) (interface{}, error)
}

type responseConsumer struct {
	runtime *Runtime
	create  ResultFactory
}

func (hr *responseConsumer) Consume(response *http.Response) (interface{}, error) {
	// work out the consumer
	// TODO: be smarter about this
	consumer := hr.runtime.Consumers[hr.runtime.DefaultMediaType]
	return hr.create(response.StatusCode, response.Body, consumer)
}

// New creates a new default runtime for a swagger api client.
func New(swaggerSpec *spec.Document) *Runtime {
	var rt Runtime
	rt.DefaultMediaType = httpkit.JSONMime
	rt.Consumers = map[string]httpkit.Consumer{
		httpkit.JSONMime: httpkit.JSONConsumer(),
	}
	rt.Producers = map[string]httpkit.Producer{
		httpkit.JSONMime: httpkit.JSONProducer(),
	}
	rt.Spec = swaggerSpec
	rt.Transport = http.DefaultTransport
	rt.client = http.DefaultClient
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
func (r *Runtime) Submit(operationID string, params client.RequestWriter, responses ResultFactory) (interface{}, error) {
	mthPth, ok := r.methodsAndPaths[operationID]
	if !ok {
		return nil, fmt.Errorf("unknown operation: %q", operationID)
	}
	request, err := NewRequest(mthPth.Method, mthPth.PathPattern, params)
	if err != nil {
		return nil, err
	}

	// TODO: Something smarter for the content type
	request.SetHeaderParam(httpkit.HeaderContentType, r.DefaultMediaType)
	// TODO: Something smarter for the Accept headers, like say multiple sorted desc by q, preferring json as default
	request.SetHeaderParam(httpkit.HeaderAccept, r.DefaultMediaType)

	req, err := request.BuildHTTP(r.Producers[r.DefaultMediaType], r.Formats)
	req.URL.Scheme = "http"
	req.URL.Host = r.Spec.Host()
	req.URL.Path = filepath.Join(r.Spec.BasePath(), req.URL.Path)
	if err != nil {
		return nil, err
	}

	res, err := r.client.Do(req) // make requests, by default follows 10 redirects before failing
	if err != nil {
		return nil, err
	}

	return (&responseConsumer{r, responses}).Consume(res)
}
