package httpkit

import (
	"encoding/json"
	"io"
)

// JSONConsumer creates a new JSON consumer
func JSONConsumer() Consumer {
	return ConsumerFunc(func(reader io.Reader, data interface{}) error {
		dec := json.NewDecoder(reader)
		return dec.Decode(data)
	})
}

// JSONProducer creates a new JSON producer
func JSONProducer() Producer {
	return ProducerFunc(func(writer io.Writer, data interface{}) error {
		enc := json.NewEncoder(writer)
		return enc.Encode(data)
	})
}
