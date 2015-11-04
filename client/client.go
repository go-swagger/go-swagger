package client

// Operation represents the context for a swagger operation to be submitted to the transport
type Operation struct {
	ID       string
	AuthInfo AuthInfoWriter
	Params   RequestWriter
	Reader   ResponseReader
}

// A Transport implementor knows how to submit Request objects to some destination
type Transport interface {
	//Submit(string, RequestWriter, ResponseReader, AuthInfoWriter) (interface{}, error)
	Submit(*Operation) (interface{}, error)
}
