package swagger

import (
	"encoding/json"
	"io"
)

// JSONConsumer creates a new JSON consumer
func JSONConsumer() Consumer {
	return FuncConsumer(func(reader io.Reader, data interface{}) error {
		dec := json.NewDecoder(reader)
		return dec.Decode(data)
	})
}

// JSONProducer creates a new JSON producer
func JSONProducer() Producer {
	return FuncProducer(func(writer io.Writer, data interface{}) error {
		enc := json.NewEncoder(writer)
		return enc.Encode(data)
	})
}
