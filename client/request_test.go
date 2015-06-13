package client

import (
	"net/http"
	"testing"

	"github.com/casualjim/go-swagger/strfmt"
	"github.com/stretchr/testify/assert"
)

type trw struct {
	Headers http.Header
	Body    interface{}
}

func (t *trw) SetHeaderParam(name string, values ...string) error {
	if t.Headers == nil {
		t.Headers = make(http.Header)
	}
	t.Headers.Set(name, values[0])
	return nil
}

func (t *trw) SetQueryParam(_ string, _ ...string) error { return nil }

func (t *trw) SetFormParam(_ string, _ ...string) error { return nil }

func (t *trw) SetPathParam(_ string, _ string) error { return nil }

func (t *trw) SetFileParam(_ string, _ string) error { return nil }

func (t *trw) SetBodyParam(body interface{}) error {
	t.Body = body
	return nil
}

func TestRequestWriterFunc(t *testing.T) {

	hand := RequestWriterFunc(func(r Request, reg strfmt.Registry) error {
		r.SetHeaderParam("blah", "blah blah")
		r.SetBodyParam(struct{ Name string }{"Adriana"})
		return nil
	})

	tr := new(trw)
	hand.WriteToRequest(tr, nil)
	assert.Equal(t, "blah blah", tr.Headers.Get("blah"))
	assert.Equal(t, "Adriana", tr.Body.(struct{ Name string }).Name)
}
