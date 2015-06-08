package httpkit

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var consProdJSON = `{"name":"Somebody","id":1}`

func TestJSONConsumer(t *testing.T) {
	cons := JSONConsumer()
	var data struct {
		Name string
		ID   int
	}
	err := cons.Consume(bytes.NewBuffer([]byte(consProdJSON)), &data)
	assert.NoError(t, err)
	assert.Equal(t, "Somebody", data.Name)
	assert.Equal(t, 1, data.ID)
}

func TestJSONProducer(t *testing.T) {
	prod := JSONProducer()
	data := struct {
		Name string `json:"name"`
		ID   int    `json:"id"`
	}{Name: "Somebody", ID: 1}

	rw := httptest.NewRecorder()
	err := prod.Produce(rw, data)
	assert.NoError(t, err)
	assert.Equal(t, consProdJSON+"\n", rw.Body.String())
}
