package httpkit

import (
	"io"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// YAMLConsumer creates a consumer for yaml data
func YAMLConsumer() Consumer {
	return ConsumerFunc(func(r io.Reader, v interface{}) error {
		buf, err := ioutil.ReadAll(r)
		if err != nil {
			return err
		}
		return yaml.Unmarshal(buf, v)
	})
}

// YAMLProducer creates a producer for yaml data
func YAMLProducer() Producer {
	return ProducerFunc(func(w io.Writer, v interface{}) error {
		b, _ := yaml.Marshal(v) // can't make this error come up
		_, err := w.Write(b)
		return err
	})
}
