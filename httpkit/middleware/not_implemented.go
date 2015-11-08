package middleware

import (
	"net/http"

	"github.com/go-swagger/go-swagger/httpkit"
)

type errorResp struct {
	code     int
	response interface{}
	headers  http.Header
}

func (e *errorResp) WriteResponse(rw http.ResponseWriter, producer httpkit.Producer) {
	for k, v := range e.headers {
		for _, val := range v {
			rw.Header().Add(k, val)
		}
	}
	if e.code > 0 {
		rw.WriteHeader(e.code)
	} else {
		rw.WriteHeader(http.StatusInternalServerError)
	}
	if err := producer.Produce(rw, e.response); err != nil {
		panic(err)
	}
}

// NotImplemented the error response when the response is not implemented
func NotImplemented(message string) Responder {
	return &errorResp{http.StatusNotImplemented, message, make(http.Header)}
}
