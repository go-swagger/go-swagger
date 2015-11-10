// Copyright 2013 The Go Authors. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file or at
// https://developers.google.com/open-source/licenses/bsd.

package httputil

import (
	"bytes"
	"net/http"
	"strconv"
)

type ResponseBuffer struct {
	buf    bytes.Buffer
	status int
	header http.Header
}

func (rb *ResponseBuffer) Write(p []byte) (int, error) {
	return rb.buf.Write(p)
}

func (rb *ResponseBuffer) WriteHeader(status int) {
	rb.status = status
}

func (rb *ResponseBuffer) Header() http.Header {
	if rb.header == nil {
		rb.header = make(http.Header)
	}
	return rb.header
}

func (rb *ResponseBuffer) WriteTo(w http.ResponseWriter) error {
	for k, v := range rb.header {
		w.Header()[k] = v
	}
	if rb.buf.Len() > 0 {
		w.Header().Set("Content-Length", strconv.Itoa(rb.buf.Len()))
	}
	if rb.status != 0 {
		w.WriteHeader(rb.status)
	}
	if rb.buf.Len() > 0 {
		if _, err := w.Write(rb.buf.Bytes()); err != nil {
			return err
		}
	}
	return nil
}
