package validate

import (
	"bytes"
	"net/http"
	"testing"

	swagger_api "github.com/casualjim/go-swagger"
	"github.com/stretchr/testify/assert"
)

func TestUntypedBindingTypesForValid(t *testing.T) {
	op1 := parametersForJSONRequestParams("")

	binder := &operationBinder{Parameters: op1, Consumers: map[string]swagger_api.Consumer{"application/json": swagger_api.JSONConsumer()}}

	qsList := "one,two,three"
	urlStr := "http://localhost:8002/hello/1?name=the-name&tags=" + qsList

	req, _ := http.NewRequest("POST", urlStr, bytes.NewBuffer([]byte(`{"name":"toby","age":32}`)))
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	req.Header.Set("X-Request-Id", "1325959595")

	data := make(map[string]interface{})
	err := binder.Bind(req, swagger_api.RouteParams([]swagger_api.RouteParam{{"id", "1"}}), &data)

	expected := map[string]interface{}{
		"id":           1,
		"name":         "the-name",
		"friend":       map[string]interface{}{"name": "toby", "age": float64(32)}, // json request
		"X-Request-Id": "1325959595",
		"tags":         []string{"one", "two", "three"},
	}
	assert.NoError(t, err)
	assert.Equal(t, expected["id"], data["id"])
	assert.Equal(t, expected["name"], data["name"])
	assert.Equal(t, expected["friend"], data["friend"])
	assert.Equal(t, expected["X-Request-Id"], data["X-Request-Id"])
	assert.Equal(t, expected["tags"], data["tags"])
}
