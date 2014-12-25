package swagger

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var consProdYAML = "name: Somebody\nid: 1\n"

func TestYAMLConsumer(t *testing.T) {
	cons := YAMLConsumer()
	var data struct {
		Name string
		ID   int
	}
	err := cons.Consume(bytes.NewBuffer([]byte(consProdYAML)), &data)
	assert.NoError(t, err)
	assert.Equal(t, "Somebody", data.Name)
	assert.Equal(t, 1, data.ID)
}

func TestYAMLProducer(t *testing.T) {
	prod := YAMLProducer()
	data := struct {
		Name string `yaml:"name"`
		ID   int    `yaml:"id"`
	}{Name: "Somebody", ID: 1}

	rw := httptest.NewRecorder()
	err := prod.Produce(rw, data)
	assert.NoError(t, err)
	assert.Equal(t, consProdYAML, rw.Body.String())
}
