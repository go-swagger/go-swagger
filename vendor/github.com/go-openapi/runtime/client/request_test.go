// Copyright 2015 go-swagger maintainers
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package client

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/assert"
)

var testProducers = map[string]runtime.Producer{
	runtime.JSONMime: runtime.JSONProducer(),
	runtime.XMLMime:  runtime.XMLProducer(),
	runtime.TextMime: runtime.TextProducer(),
}

func TestBuildRequest_SetHeaders(t *testing.T) {
	r, _ := newRequest("GET", "/flats/{id}/", nil)
	// single value
	_ = r.SetHeaderParam("X-Rate-Limit", "500")
	assert.Equal(t, "500", r.header.Get("X-Rate-Limit"))
	_ = r.SetHeaderParam("X-Rate-Limit", "400")
	assert.Equal(t, "400", r.header.Get("X-Rate-Limit"))

	// multi value
	_ = r.SetHeaderParam("X-Accepts", "json", "xml", "yaml")
	assert.EqualValues(t, []string{"json", "xml", "yaml"}, r.header["X-Accepts"])
}

func TestBuildRequest_SetPath(t *testing.T) {
	r, _ := newRequest("GET", "/flats/{id}/?hello=world", nil)

	_ = r.SetPathParam("id", "1345")
	assert.Equal(t, "1345", r.pathParams["id"])
}

func TestBuildRequest_SetQuery(t *testing.T) {
	r, _ := newRequest("GET", "/flats/{id}/", nil)

	// single value
	_ = r.SetQueryParam("hello", "there")
	assert.Equal(t, "there", r.query.Get("hello"))

	// multi value
	_ = r.SetQueryParam("goodbye", "cruel", "world")
	assert.Equal(t, []string{"cruel", "world"}, r.query["goodbye"])
}

func TestBuildRequest_SetForm(t *testing.T) {
	// non-multipart
	r, _ := newRequest("POST", "/flats", nil)
	_ = r.SetFormParam("hello", "world")
	assert.Equal(t, "world", r.formFields.Get("hello"))
	_ = r.SetFormParam("goodbye", "cruel", "world")
	assert.Equal(t, []string{"cruel", "world"}, r.formFields["goodbye"])
}

func TestBuildRequest_SetFile(t *testing.T) {
	// needs to convert form to multipart
	r, _ := newRequest("POST", "/flats/{id}/image", nil)
	// error if it isn't there
	err := r.SetFileParam("not there", os.NewFile(0, "./i-dont-exist"))
	assert.Error(t, err)
	// error if it isn't a file
	err = r.SetFileParam("directory", os.NewFile(0, "../client"))
	assert.Error(t, err)
	// success adds it to the map
	err = r.SetFileParam("file", mustGetFile("./runtime.go"))
	if assert.NoError(t, err) {
		fl, ok := r.fileFields["file"]
		if assert.True(t, ok) {
			assert.Equal(t, "runtime.go", filepath.Base(fl[0].Name()))
		}
	}
	// success adds a file param with multiple files
	err = r.SetFileParam("otherfiles", mustGetFile("./runtime.go"), mustGetFile("./request.go"))
	if assert.NoError(t, err) {
		fl, ok := r.fileFields["otherfiles"]
		if assert.True(t, ok) {
			assert.Equal(t, "runtime.go", filepath.Base(fl[0].Name()))
			assert.Equal(t, "request.go", filepath.Base(fl[1].Name()))
		}
	}
}

func mustGetFile(path string) *os.File {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	return f
}

func TestBuildRequest_SetBody(t *testing.T) {
	r, _ := newRequest("GET", "/flats/{id}/?hello=world", nil)
	bd := []struct{ Name, Hobby string }{{"Tom", "Organ trail"}, {"John", "Bird watching"}}

	_ = r.SetBodyParam(bd)
	assert.Equal(t, bd, r.payload)
}

func TestBuildRequest_BuildHTTP_NoPayload(t *testing.T) {
	reqWrtr := runtime.ClientRequestWriterFunc(func(req runtime.ClientRequest, reg strfmt.Registry) error {
		_ = req.SetBodyParam(nil)
		_ = req.SetQueryParam("hello", "world")
		_ = req.SetPathParam("id", "1234")
		_ = req.SetHeaderParam("X-Rate-Limit", "200")
		return nil
	})
	r, _ := newRequest("POST", "/flats/{id}/", reqWrtr)

	req, err := r.BuildHTTP(runtime.JSONMime, "", testProducers, nil)
	if assert.NoError(t, err) && assert.NotNil(t, req) {
		assert.Equal(t, "200", req.Header.Get("x-rate-limit"))
		assert.Equal(t, "world", req.URL.Query().Get("hello"))
		assert.Equal(t, "/flats/1234/", req.URL.Path)
		assert.Equal(t, runtime.JSONMime, req.Header.Get(runtime.HeaderContentType))
	}
}

func TestBuildRequest_BuildHTTP_Payload(t *testing.T) {
	bd := []struct{ Name, Hobby string }{{"Tom", "Organ trail"}, {"John", "Bird watching"}}
	reqWrtr := runtime.ClientRequestWriterFunc(func(req runtime.ClientRequest, reg strfmt.Registry) error {
		_ = req.SetBodyParam(bd)
		_ = req.SetQueryParam("hello", "world")
		_ = req.SetPathParam("id", "1234")
		_ = req.SetHeaderParam("X-Rate-Limit", "200")
		return nil
	})
	r, _ := newRequest("GET", "/flats/{id}/", reqWrtr)
	_ = r.SetHeaderParam(runtime.HeaderContentType, runtime.JSONMime)

	req, err := r.BuildHTTP(runtime.JSONMime, "", testProducers, nil)
	if assert.NoError(t, err) && assert.NotNil(t, req) {
		assert.Equal(t, "200", req.Header.Get("x-rate-limit"))
		assert.Equal(t, "world", req.URL.Query().Get("hello"))
		assert.Equal(t, "/flats/1234/", req.URL.Path)
		expectedBody, _ := json.Marshal(bd)
		actualBody, _ := ioutil.ReadAll(req.Body)
		assert.Equal(t, append(expectedBody, '\n'), actualBody)
	}
}

func TestBuildRequest_BuildHTTP_SetsInAuth(t *testing.T) {
	bd := []struct{ Name, Hobby string }{{"Tom", "Organ trail"}, {"John", "Bird watching"}}
	reqWrtr := runtime.ClientRequestWriterFunc(func(req runtime.ClientRequest, reg strfmt.Registry) error {
		_ = req.SetBodyParam(bd)
		_ = req.SetQueryParam("hello", "wrong")
		_ = req.SetPathParam("id", "wrong")
		_ = req.SetHeaderParam("X-Rate-Limit", "wrong")
		return nil
	})

	auth := runtime.ClientAuthInfoWriterFunc(func(req runtime.ClientRequest, reg strfmt.Registry) error {
		_ = req.SetBodyParam(bd)
		_ = req.SetQueryParam("hello", "world")
		_ = req.SetPathParam("id", "1234")
		_ = req.SetHeaderParam("X-Rate-Limit", "200")
		return nil
	})

	r, _ := newRequest("GET", "/flats/{id}/", reqWrtr)
	_ = r.SetHeaderParam(runtime.HeaderContentType, runtime.JSONMime)

	req, err := r.buildHTTP(runtime.JSONMime, "", testProducers, nil, auth)
	if assert.NoError(t, err) && assert.NotNil(t, req) {
		assert.Equal(t, "200", req.Header.Get("x-rate-limit"))
		assert.Equal(t, "world", req.URL.Query().Get("hello"))
		assert.Equal(t, "/flats/1234/", req.URL.Path)
		expectedBody, _ := json.Marshal(bd)
		actualBody, _ := ioutil.ReadAll(req.Body)
		assert.Equal(t, append(expectedBody, '\n'), actualBody)
	}
}

func TestBuildRequest_BuildHTTP_XMLPayload(t *testing.T) {
	bd := []struct {
		XMLName xml.Name `xml:"person"`
		Name    string   `xml:"name"`
		Hobby   string   `xml:"hobby"`
	}{{xml.Name{}, "Tom", "Organ trail"}, {xml.Name{}, "John", "Bird watching"}}
	reqWrtr := runtime.ClientRequestWriterFunc(func(req runtime.ClientRequest, reg strfmt.Registry) error {
		_ = req.SetBodyParam(bd)
		_ = req.SetQueryParam("hello", "world")
		_ = req.SetPathParam("id", "1234")
		_ = req.SetHeaderParam("X-Rate-Limit", "200")
		return nil
	})
	r, _ := newRequest("GET", "/flats/{id}/", reqWrtr)
	_ = r.SetHeaderParam(runtime.HeaderContentType, runtime.XMLMime)

	req, err := r.BuildHTTP(runtime.XMLMime, "", testProducers, nil)
	if assert.NoError(t, err) && assert.NotNil(t, req) {
		assert.Equal(t, "200", req.Header.Get("x-rate-limit"))
		assert.Equal(t, "world", req.URL.Query().Get("hello"))
		assert.Equal(t, "/flats/1234/", req.URL.Path)
		expectedBody, _ := xml.Marshal(bd)
		actualBody, _ := ioutil.ReadAll(req.Body)
		assert.Equal(t, expectedBody, actualBody)
	}
}

func TestBuildRequest_BuildHTTP_TextPayload(t *testing.T) {
	bd := "Tom: Organ trail; John: Bird watching"
	reqWrtr := runtime.ClientRequestWriterFunc(func(req runtime.ClientRequest, reg strfmt.Registry) error {
		_ = req.SetBodyParam(bd)
		_ = req.SetQueryParam("hello", "world")
		_ = req.SetPathParam("id", "1234")
		_ = req.SetHeaderParam("X-Rate-Limit", "200")
		return nil
	})
	r, _ := newRequest("GET", "/flats/{id}/", reqWrtr)
	_ = r.SetHeaderParam(runtime.HeaderContentType, runtime.TextMime)

	req, err := r.BuildHTTP(runtime.TextMime, "", testProducers, nil)
	if assert.NoError(t, err) && assert.NotNil(t, req) {
		assert.Equal(t, "200", req.Header.Get("x-rate-limit"))
		assert.Equal(t, "world", req.URL.Query().Get("hello"))
		assert.Equal(t, "/flats/1234/", req.URL.Path)
		expectedBody := []byte(bd)
		actualBody, _ := ioutil.ReadAll(req.Body)
		assert.Equal(t, expectedBody, actualBody)
	}
}

func TestBuildRequest_BuildHTTP_Form(t *testing.T) {
	reqWrtr := runtime.ClientRequestWriterFunc(func(req runtime.ClientRequest, reg strfmt.Registry) error {
		_ = req.SetFormParam("something", "some value")
		_ = req.SetQueryParam("hello", "world")
		_ = req.SetPathParam("id", "1234")
		_ = req.SetHeaderParam("X-Rate-Limit", "200")
		return nil
	})
	r, _ := newRequest("GET", "/flats/{id}/", reqWrtr)
	_ = r.SetHeaderParam(runtime.HeaderContentType, runtime.JSONMime)

	req, err := r.BuildHTTP(runtime.JSONMime, "", testProducers, nil)
	if assert.NoError(t, err) && assert.NotNil(t, req) {
		assert.Equal(t, "200", req.Header.Get("x-rate-limit"))
		assert.Equal(t, "world", req.URL.Query().Get("hello"))
		assert.Equal(t, "/flats/1234/", req.URL.Path)
		expected := []byte("something=some+value")
		actual, _ := ioutil.ReadAll(req.Body)
		assert.Equal(t, expected, actual)
	}
}

func TestBuildRequest_BuildHTTP_Form_Content_Length(t *testing.T) {
	reqWrtr := runtime.ClientRequestWriterFunc(func(req runtime.ClientRequest, reg strfmt.Registry) error {
		_ = req.SetFormParam("something", "some value")
		_ = req.SetQueryParam("hello", "world")
		_ = req.SetPathParam("id", "1234")
		_ = req.SetHeaderParam("X-Rate-Limit", "200")
		return nil
	})
	r, _ := newRequest("GET", "/flats/{id}/", reqWrtr)
	_ = r.SetHeaderParam(runtime.HeaderContentType, runtime.MultipartFormMime)

	req, err := r.BuildHTTP(runtime.JSONMime, "", testProducers, nil)
	if assert.NoError(t, err) && assert.NotNil(t, req) {
		assert.Equal(t, "200", req.Header.Get("x-rate-limit"))
		assert.Equal(t, "world", req.URL.Query().Get("hello"))
		assert.Equal(t, "/flats/1234/", req.URL.Path)
		assert.Condition(t, func() bool { return req.ContentLength > 0 },
			"ContentLength must great than 0. got %d", req.ContentLength)
		expected := []byte("something=some+value")
		actual, _ := ioutil.ReadAll(req.Body)
		assert.Equal(t, expected, actual)
	}
}

func TestBuildRequest_BuildHTTP_FormMultipart(t *testing.T) {
	reqWrtr := runtime.ClientRequestWriterFunc(func(req runtime.ClientRequest, reg strfmt.Registry) error {
		_ = req.SetFormParam("something", "some value")
		_ = req.SetQueryParam("hello", "world")
		_ = req.SetPathParam("id", "1234")
		_ = req.SetHeaderParam("X-Rate-Limit", "200")
		return nil
	})
	r, _ := newRequest("GET", "/flats/{id}/", reqWrtr)
	_ = r.SetHeaderParam(runtime.HeaderContentType, runtime.MultipartFormMime)

	req, err := r.BuildHTTP(runtime.MultipartFormMime, "", testProducers, nil)
	if assert.NoError(t, err) && assert.NotNil(t, req) {
		assert.Equal(t, "200", req.Header.Get("x-rate-limit"))
		assert.Equal(t, "world", req.URL.Query().Get("hello"))
		assert.Equal(t, "/flats/1234/", req.URL.Path)
		expected1 := []byte("Content-Disposition: form-data; name=\"something\"")
		expected2 := []byte("some value")
		actual, _ := ioutil.ReadAll(req.Body)
		actuallines := bytes.Split(actual, []byte("\r\n"))
		assert.Equal(t, 6, len(actuallines))
		boundary := string(actuallines[0])
		lastboundary := string(actuallines[4])
		assert.True(t, strings.HasPrefix(boundary, "--"))
		assert.True(t, strings.HasPrefix(lastboundary, "--") && strings.HasSuffix(lastboundary, "--"))
		assert.Equal(t, lastboundary, boundary+"--")
		assert.Equal(t, expected1, actuallines[1])
		assert.Equal(t, expected2, actuallines[3])
	}
}

func TestBuildRequest_BuildHTTP_FormMultiples(t *testing.T) {
	reqWrtr := runtime.ClientRequestWriterFunc(func(req runtime.ClientRequest, reg strfmt.Registry) error {
		_ = req.SetFormParam("something", "some value", "another value")
		_ = req.SetQueryParam("hello", "world")
		_ = req.SetPathParam("id", "1234")
		_ = req.SetHeaderParam("X-Rate-Limit", "200")
		return nil
	})
	r, _ := newRequest("GET", "/flats/{id}/", reqWrtr)
	_ = r.SetHeaderParam(runtime.HeaderContentType, runtime.MultipartFormMime)

	req, err := r.BuildHTTP(runtime.MultipartFormMime, "", testProducers, nil)
	if assert.NoError(t, err) && assert.NotNil(t, req) {
		assert.Equal(t, "200", req.Header.Get("x-rate-limit"))
		assert.Equal(t, "world", req.URL.Query().Get("hello"))
		assert.Equal(t, "/flats/1234/", req.URL.Path)
		expected1 := []byte("Content-Disposition: form-data; name=\"something\"")
		expected2 := []byte("some value")
		expected3 := []byte("another value")
		actual, _ := ioutil.ReadAll(req.Body)
		actuallines := bytes.Split(actual, []byte("\r\n"))
		assert.Equal(t, 10, len(actuallines))
		boundary := string(actuallines[0])
		lastboundary := string(actuallines[8])
		assert.True(t, strings.HasPrefix(boundary, "--"))
		assert.True(t, strings.HasPrefix(lastboundary, "--") && strings.HasSuffix(lastboundary, "--"))
		assert.Equal(t, lastboundary, boundary+"--")
		assert.Equal(t, expected1, actuallines[1])
		assert.Equal(t, expected2, actuallines[3])
		assert.Equal(t, actuallines[0], actuallines[4])
		assert.Equal(t, expected1, actuallines[5])
		assert.Equal(t, expected3, actuallines[7])
	}
}

func TestBuildRequest_BuildHTTP_Files(t *testing.T) {
	cont, _ := ioutil.ReadFile("./runtime.go")
	cont2, _ := ioutil.ReadFile("./request.go")
	reqWrtr := runtime.ClientRequestWriterFunc(func(req runtime.ClientRequest, reg strfmt.Registry) error {
		_ = req.SetFormParam("something", "some value")
		_ = req.SetFileParam("file", mustGetFile("./runtime.go"))
		_ = req.SetFileParam("otherfiles", mustGetFile("./runtime.go"), mustGetFile("./request.go"))
		_ = req.SetQueryParam("hello", "world")
		_ = req.SetPathParam("id", "1234")
		_ = req.SetHeaderParam("X-Rate-Limit", "200")
		return nil
	})
	r, _ := newRequest("GET", "/flats/{id}/", reqWrtr)
	_ = r.SetHeaderParam(runtime.HeaderContentType, runtime.JSONMime)
	req, err := r.BuildHTTP(runtime.JSONMime, "", testProducers, nil)
	if assert.NoError(t, err) && assert.NotNil(t, req) {
		assert.Equal(t, "200", req.Header.Get("x-rate-limit"))
		assert.Equal(t, "world", req.URL.Query().Get("hello"))
		assert.Equal(t, "/flats/1234/", req.URL.Path)
		mediaType, params, err := mime.ParseMediaType(req.Header.Get(runtime.HeaderContentType))
		if assert.NoError(t, err) {
			assert.Equal(t, runtime.MultipartFormMime, mediaType)
			boundary := params["boundary"]
			mr := multipart.NewReader(req.Body, boundary)
			defer req.Body.Close()
			frm, err := mr.ReadForm(1 << 20)
			if assert.NoError(t, err) {
				assert.Equal(t, "some value", frm.Value["something"][0])
				fileverifier := func(name string, index int, filename string, content []byte) {
					mpff := frm.File[name][index]
					mpf, _ := mpff.Open()
					defer mpf.Close()
					assert.Equal(t, filename, mpff.Filename)
					actual, _ := ioutil.ReadAll(mpf)
					assert.Equal(t, content, actual)
				}
				fileverifier("file", 0, "runtime.go", cont)

				fileverifier("otherfiles", 0, "runtime.go", cont)
				fileverifier("otherfiles", 1, "request.go", cont2)
			}
		}
	}
}
func TestBuildRequest_BuildHTTP_Files_URLEncoded(t *testing.T) {
	cont, _ := ioutil.ReadFile("./runtime.go")
	cont2, _ := ioutil.ReadFile("./request.go")
	reqWrtr := runtime.ClientRequestWriterFunc(func(req runtime.ClientRequest, reg strfmt.Registry) error {
		_ = req.SetFormParam("something", "some value")
		_ = req.SetFileParam("file", mustGetFile("./runtime.go"))
		_ = req.SetFileParam("otherfiles", mustGetFile("./runtime.go"), mustGetFile("./request.go"))
		_ = req.SetQueryParam("hello", "world")
		_ = req.SetPathParam("id", "1234")
		_ = req.SetHeaderParam("X-Rate-Limit", "200")
		return nil
	})
	r, _ := newRequest("GET", "/flats/{id}/", reqWrtr)
	_ = r.SetHeaderParam(runtime.HeaderContentType, runtime.URLencodedFormMime)
	req, err := r.BuildHTTP(runtime.URLencodedFormMime, "", testProducers, nil)
	if assert.NoError(t, err) && assert.NotNil(t, req) {
		assert.Equal(t, "200", req.Header.Get("x-rate-limit"))
		assert.Equal(t, "world", req.URL.Query().Get("hello"))
		assert.Equal(t, "/flats/1234/", req.URL.Path)
		mediaType, params, err := mime.ParseMediaType(req.Header.Get(runtime.HeaderContentType))
		if assert.NoError(t, err) {
			assert.Equal(t, runtime.URLencodedFormMime, mediaType)
			boundary := params["boundary"]
			mr := multipart.NewReader(req.Body, boundary)
			defer req.Body.Close()
			frm, err := mr.ReadForm(1 << 20)
			if assert.NoError(t, err) {
				assert.Equal(t, "some value", frm.Value["something"][0])
				fileverifier := func(name string, index int, filename string, content []byte) {
					mpff := frm.File[name][index]
					mpf, _ := mpff.Open()
					defer mpf.Close()
					assert.Equal(t, filename, mpff.Filename)
					actual, _ := ioutil.ReadAll(mpf)
					assert.Equal(t, content, actual)
				}
				fileverifier("file", 0, "runtime.go", cont)

				fileverifier("otherfiles", 0, "runtime.go", cont)
				fileverifier("otherfiles", 1, "request.go", cont2)
			}
		}
	}
}

func TestBuildRequest_BuildHTTP_BasePath(t *testing.T) {
	reqWrtr := runtime.ClientRequestWriterFunc(func(req runtime.ClientRequest, reg strfmt.Registry) error {
		_ = req.SetBodyParam(nil)
		_ = req.SetQueryParam("hello", "world")
		_ = req.SetPathParam("id", "1234")
		_ = req.SetHeaderParam("X-Rate-Limit", "200")
		return nil
	})
	r, _ := newRequest("POST", "/flats/{id}/", reqWrtr)

	req, err := r.BuildHTTP(runtime.JSONMime, "/basepath", testProducers, nil)
	if assert.NoError(t, err) && assert.NotNil(t, req) {
		assert.Equal(t, "200", req.Header.Get("x-rate-limit"))
		assert.Equal(t, "world", req.URL.Query().Get("hello"))
		assert.Equal(t, "/basepath/flats/1234/", req.URL.Path)
		assert.Equal(t, runtime.JSONMime, req.Header.Get(runtime.HeaderContentType))
	}
}

func TestBuildRequest_BuildHTTP_EscapedPath(t *testing.T) {
	reqWrtr := runtime.ClientRequestWriterFunc(func(req runtime.ClientRequest, reg strfmt.Registry) error {
		_ = req.SetBodyParam(nil)
		_ = req.SetQueryParam("hello", "world")
		_ = req.SetPathParam("id", "1234/?*&^%")
		_ = req.SetHeaderParam("X-Rate-Limit", "200")
		return nil
	})
	r, _ := newRequest("POST", "/flats/{id}/", reqWrtr)

	req, err := r.BuildHTTP(runtime.JSONMime, "/basepath", testProducers, nil)
	if assert.NoError(t, err) && assert.NotNil(t, req) {
		assert.Equal(t, "200", req.Header.Get("x-rate-limit"))
		assert.Equal(t, "world", req.URL.Query().Get("hello"))
		assert.Equal(t, "/basepath/flats/1234/?*&^%/", req.URL.Path)
		assert.Equal(t, "/basepath/flats/1234%2F%3F%2A&%5E%25/", req.URL.RawPath)
		assert.Equal(t, req.URL.RawPath, req.URL.EscapedPath())
		assert.Equal(t, runtime.JSONMime, req.Header.Get(runtime.HeaderContentType))
	}
}
