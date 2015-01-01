package jsonschema

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"strings"

	"github.com/casualjim/go-swagger/swagger/jsonreference"
)

// Loader represents a loader for a json schema
type Loader interface {
	URL() string
	Load() (interface{}, error)
}

type readerLoader struct {
	Reader   io.Reader
	Location string
}

// NewLoader creates a new loader that gets the schema from a io.Reader
func NewLoader(reader io.Reader, url string) Loader {
	return &readerLoader{reader, url}
}

func (r *readerLoader) URL() string {
	return r.Location
}

func (r *readerLoader) Load() (interface{}, error) {
	return readAll(r.Reader)
}

func readAll(reader io.Reader) (interface{}, error) {
	bodyBuff, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	var document interface{}
	err = json.Unmarshal(bodyBuff, &document)
	if err != nil {
		return nil, err
	}

	return document, nil
}

type referenceLoader struct {
	ref *jsonreference.Ref
}

// NewReferenceLoader creates a loader for a jsonreference.Ref
func NewReferenceLoader(reference *jsonreference.Ref) Loader {
	return &referenceLoader{reference}
}

func (r *referenceLoader) URL() string {
	return r.ref.String()
}

func (r *referenceLoader) Load() (interface{}, error) {
	refToURL := r.ref
	refToURL.GetURL().Fragment = ""
	if refToURL.HasFileScheme {
		// Load from file
		filename := strings.Replace(refToURL.String(), "file://", "", -1)
		return GetFileJSON(filename)
	}
	return GetHTTPJSON(refToURL.String())
}
