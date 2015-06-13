package client

import (
	"bytes"
	"io"
	"io/ioutil"
	"testing"

	"github.com/casualjim/go-swagger/httpkit"
	"github.com/stretchr/testify/assert"
)

type response struct {
}

func (r response) Code() int {
	return 490
}
func (r response) Message() string {
	return "the message"
}
func (r response) GetHeader(_ string) string {
	return "the header"
}
func (r response) Body() io.ReadCloser {
	return ioutil.NopCloser(bytes.NewBufferString("the content"))
}

func TestResponseReaderFunc(t *testing.T) {
	var actual struct {
		Header, Message, Body string
		Code                  int
	}
	reader := ResponseReaderFunc(func(r Response, _ httpkit.Consumer) (interface{}, error) {
		b, _ := ioutil.ReadAll(r.Body())
		actual.Body = string(b)
		actual.Code = r.Code()
		actual.Message = r.Message()
		actual.Header = r.GetHeader("blah")
		return actual, nil
	})
	reader.ReadResponse(response{}, nil)
	assert.Equal(t, "the content", actual.Body)
	assert.Equal(t, "the message", actual.Message)
	assert.Equal(t, "the header", actual.Header)
	assert.Equal(t, 490, actual.Code)
}
