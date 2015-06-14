package client

// A Transport implementor knows how to submit Request objects to some destination
type Transport interface {
	Submit(operationID string, params RequestWriter, readResponse ResponseReader) (interface{}, error)
}
