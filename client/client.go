package client

// A Transport implementor knows how to submit Request objects to some destination
type Transport interface {
	Submit(string, RequestWriter, ResponseReader) (interface{}, error)
}
