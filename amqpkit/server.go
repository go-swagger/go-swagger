package amqpkit

import (
	"bytes"
	"net/http"

	"github.com/streadway/amqp"
)

// Server contains all the configuration options for an AMQP server
type Server struct {
	Queue   *amqp.Queue
	Channel *amqp.Channel
}

// Serve serves an HTTP handler over amqp
func (s *Server) Serve(incoming <-chan amqp.Delivery, handler http.Handler) error {
	for d := range incoming {
		req, err := amqpDeliveryRequest(&d)
		if err != nil {
			return err
		}

		rw := s.amqpResponseWriter(&d, req)
		handler.ServeHTTP(rw, req)

		// TODO: nacking is probably just as important as acking
		d.Ack(false)

		// once we said we received and processed the message, we should reply too
		s.Channel.Publish(
			"",        // exchange
			d.ReplyTo, // routing key
			false,     // mandatory
			false,     // immediate
			rw.close(),
		)
	}
	return nil
}

// the client should add these headers for it to be swagger enabled
var reservedHeaders = map[string]struct{}{
	"Swagger-Path":       struct{}{},
	"Swagger-Method":     struct{}{},
	"Swagger-Context-Id": struct{}{},
}

func amqpDeliveryRequest(delivery *amqp.Delivery) (*http.Request, error) {
	method := delivery.Headers["Swagger-Method"].(string)
	path := delivery.Headers["Swagger-Path"].(string)
	r, err := http.NewRequest(method, path, bytes.NewBuffer(delivery.Body))
	if err != nil {
		return nil, err
	}

	for k := range delivery.Headers {
		ck := http.CanonicalHeaderKey(k)
		if _, reserved := reservedHeaders[ck]; !reserved {
			if slice, safe := delivery.Headers[k].([]interface{}); safe {
				var hv []string
				for _, v := range slice {
					if str, ok := v.(string); ok {
						hv = append(hv, str)
					}
				}

				if r.Header == nil {
					r.Header = make(http.Header)
				}
				r.Header.Del(ck)
				for _, vv := range hv {
					r.Header.Set(ck, vv)
				}
			}
		}
	}
	return r, nil
}

func (s *Server) amqpResponseWriter(delivery *amqp.Delivery, req *http.Request) *amqpServerResponse {
	sr := &amqpServerResponse{
		header:   make(http.Header),
		req:      req,
		delivery: delivery,
		buf:      bytes.NewBuffer(nil),
	}
	return sr
}

// amqpServerResponse pretends to be a http.ResponseWriter.
//
// It collects all the data in memory and then
type amqpServerResponse struct {
	header     http.Header
	req        *http.Request
	delivery   *amqp.Delivery
	statusCode int
	buf        *bytes.Buffer
}

func (a *amqpServerResponse) Header() http.Header {
	return a.header
}

func (a *amqpServerResponse) Write(data []byte) (int, error) {
	return a.buf.Write(data)
}

func (a *amqpServerResponse) WriteHeader(code int) { a.statusCode = code }

func (a *amqpServerResponse) close() (reply amqp.Publishing) {
	reply.CorrelationId = a.delivery.CorrelationId
	reply.Body = a.buf.Bytes()

	// TODO: this is extremely optimistic, parse and adapt!
	reply.ContentType = a.header.Get("Content-Type")

	ah := make(amqp.Table)
	for k, vv := range a.header {
		if len(vv) > 0 {
			ah[http.CanonicalHeaderKey(k)] = vv
		}
	}
	reply.Headers = ah
	return
}
