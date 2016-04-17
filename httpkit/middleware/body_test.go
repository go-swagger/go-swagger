package middleware

import (
	"io"
	"net/http"
	"testing"

	"github.com/go-swagger/go-swagger/httpkit"
	"github.com/go-swagger/go-swagger/internal/testing/petstore"
	"github.com/stretchr/testify/assert"
)

type eofReader struct {
}

func (r *eofReader) Read(b []byte) (int, error) {
	return 0, io.EOF
}

func (r *eofReader) Close() error {
	return nil
}

type rbn func(*http.Request, *MatchedRoute) error

func (b rbn) BindRequest(r *http.Request, rr *MatchedRoute) error {
	return b(r, rr)
}

func TestBindRequest_BodyValidation(t *testing.T) {
	spec, api := petstore.NewAPI(t)
	ctx := NewContext(spec, api, nil)
	api.DefaultConsumes = httpkit.JSONMime
	ctx.router = DefaultRouter(spec, ctx.api)

	req, err := http.NewRequest("GET", "/pets", new(eofReader))
	if assert.NoError(t, err) {
		req.Header.Set("Content-Type", httpkit.JSONMime)

		ri, ok := ctx.RouteInfo(req)
		if assert.True(t, ok) {

			err := ctx.BindValidRequest(req, ri, rbn(func(r *http.Request, rr *MatchedRoute) error {
				defer r.Body.Close()
				var data interface{}
				err := httpkit.JSONConsumer().Consume(r.Body, &data)
				_ = data
				return err
			}))

			assert.Error(t, err)
			assert.Equal(t, io.EOF, err)
		}
	}
}

func TestBindRequest_DeleteNoBody(t *testing.T) {
	spec, api := petstore.NewAPI(t)
	ctx := NewContext(spec, api, nil)
	api.DefaultConsumes = httpkit.JSONMime
	ctx.router = DefaultRouter(spec, ctx.api)

	req, err := http.NewRequest("DELETE", "/pets/123", new(eofReader))
	if assert.NoError(t, err) {
		req.Header.Set("Accept", "*/*")

		ri, ok := ctx.RouteInfo(req)
		if assert.True(t, ok) {

			err := ctx.BindValidRequest(req, ri, rbn(func(r *http.Request, rr *MatchedRoute) error {
				return nil
			}))

			assert.NoError(t, err)
			//assert.Equal(t, io.EOF, err)
		}
	}

	req, err = http.NewRequest("DELETE", "/pets/123", new(eofReader))
	if assert.NoError(t, err) {
		req.Header.Set("Accept", "*/*")
		req.Header.Set("Content-Type", httpkit.JSONMime)
		req.ContentLength = 1

		ri, ok := ctx.RouteInfo(req)
		if assert.True(t, ok) {

			err := ctx.BindValidRequest(req, ri, rbn(func(r *http.Request, rr *MatchedRoute) error {
				defer r.Body.Close()
				var data interface{}
				err := httpkit.JSONConsumer().Consume(r.Body, &data)
				_ = data
				return err
			}))

			assert.Error(t, err)
			assert.Equal(t, io.EOF, err)
		}
	}
}
