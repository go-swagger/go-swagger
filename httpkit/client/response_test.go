package client

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/casualjim/go-swagger/client"
	"github.com/stretchr/testify/assert"
)

func TestResponse(t *testing.T) {
	under := new(http.Response)
	under.Status = "the status message"
	under.StatusCode = 392
	under.Header = make(http.Header)
	under.Header.Set("Blah", "blah blah")
	under.Body = ioutil.NopCloser(bytes.NewBufferString("some content"))

	var resp client.Response = response{under}
	assert.EqualValues(t, under.StatusCode, resp.Code())
	assert.Equal(t, under.Status, resp.Message())
	assert.Equal(t, "blah blah", resp.GetHeader("blah"))
	assert.Equal(t, under.Body, resp.Body())
}
