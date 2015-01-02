package validate

import (
	"bytes"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strings"
	"testing"

	"github.com/casualjim/go-swagger"
	"github.com/casualjim/go-swagger/spec"
	"github.com/stretchr/testify/assert"
)

type friend struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type jsonRequestParams struct {
	ID        int64    // path
	Name      string   // query
	Friend    friend   // body
	RequestID int64    // header
	Tags      []string // csv
}

type jsonRequestPtr struct {
	ID        int64    // path
	Name      string   // query
	RequestID int64    // header
	Tags      []string // csv
	Friend    *friend
}

type jsonRequestSlice struct {
	ID        int64    // path
	Name      string   // query
	RequestID int64    // header
	Tags      []string // csv
	Friend    []friend
}

type jsonRequestAllTypes struct {
	Confirmed bool
	Planned   swagger.Date
	Delivered swagger.DateTime
	Age       int32
	ID        int64
	Score     float32
	Factor    float64
	Friend    friend
	Name      string
	Tags      []string
	Picture   []byte
	RequestID int64
}

func parametersForAllTypes(fmt string) map[string]spec.Parameter {
	if fmt == "" {
		fmt = "csv"
	}
	nameParam := spec.QueryParam()
	nameParam.Name = "name"
	nameParam.Type = "string"

	idParam := spec.PathParam()
	idParam.Name = "id"
	idParam.Type = "integer"
	idParam.Format = "int64"

	ageParam := spec.QueryParam()
	ageParam.Name = "age"
	ageParam.Type = "integer"
	ageParam.Format = "int32"

	scoreParam := spec.QueryParam()
	scoreParam.Name = "score"
	scoreParam.Type = "number"
	scoreParam.Format = "float"

	factorParam := spec.QueryParam()
	factorParam.Name = "factor"
	factorParam.Type = "number"
	factorParam.Format = "double"

	friendParam := spec.BodyParam()
	friendParam.Name = "friend"
	friendParam.Type = "object"

	requestIDParam := spec.HeaderParam()
	requestIDParam.Name = "X-Request-Id"
	requestIDParam.Type = "integer"
	requestIDParam.Format = "int64"
	requestIDParam.Extensions = spec.Extensions(map[string]interface{}{})
	requestIDParam.Extensions.Add("go-name", "RequestID")

	tagsParam := spec.QueryParam()
	tagsParam.Name = "tags"
	tagsParam.Type = "array"
	tagsParam.Items = new(spec.Items)
	tagsParam.Items.Type = "string"
	tagsParam.CollectionFormat = fmt

	confirmedParam := spec.QueryParam()
	confirmedParam.Name = "confirmed"
	confirmedParam.Type = "boolean"

	plannedParam := spec.QueryParam()
	plannedParam.Name = "planned"
	plannedParam.Type = "string"
	plannedParam.Format = "date"

	deliveredParam := spec.QueryParam()
	deliveredParam.Name = "delivered"
	deliveredParam.Type = "string"
	deliveredParam.Format = "date-time"

	pictureParam := spec.QueryParam() // base64 encoded during transport
	pictureParam.Name = "picture"
	pictureParam.Type = "array"
	pictureParam.Items = new(spec.Items)
	pictureParam.Items.Type = "string"
	pictureParam.Items.Format = "byte"

	return map[string]spec.Parameter{
		"ID":        *idParam,
		"Name":      *nameParam,
		"RequestID": *requestIDParam,
		"Friend":    *friendParam,
		"Tags":      *tagsParam,
		"Age":       *ageParam,
		"Score":     *scoreParam,
		"Factor":    *factorParam,
		"Confirmed": *confirmedParam,
		"Planned":   *plannedParam,
		"Delivered": *deliveredParam,
		// "Picture":   *pictureParam,
	}
}

func parametersForJSONRequestParams(fmt string) map[string]spec.Parameter {
	if fmt == "" {
		fmt = "csv"
	}
	nameParam := spec.QueryParam()
	nameParam.Name = "name"
	nameParam.Type = "string"

	idParam := spec.PathParam()
	idParam.Name = "id"
	idParam.Type = "integer"
	idParam.Format = "int64"

	friendParam := spec.BodyParam()
	friendParam.Name = "friend"
	friendParam.Type = "object"

	requestIDParam := spec.HeaderParam()
	requestIDParam.Name = "X-Request-Id"
	requestIDParam.Type = "integer"
	requestIDParam.Format = "int64"
	requestIDParam.Extensions = spec.Extensions(map[string]interface{}{})
	requestIDParam.Extensions.Add("go-name", "RequestID")

	tagsParam := spec.QueryParam()
	tagsParam.Name = "tags"
	tagsParam.Type = "array"
	tagsParam.Items = new(spec.Items)
	tagsParam.Items.Type = "string"
	tagsParam.CollectionFormat = fmt

	return map[string]spec.Parameter{
		"ID":        *idParam,
		"Name":      *nameParam,
		"RequestID": *requestIDParam,
		"Friend":    *friendParam,
		"Tags":      *tagsParam,
	}
}

func TestRequestBindingForInvalid(t *testing.T) {

	invalidParam := spec.QueryParam()
	invalidParam.Name = "some"

	op1 := map[string]spec.Parameter{"Some": *invalidParam}

	binder := &operationBinder{Parameters: op1, Consumers: map[string]swagger.Consumer{}}
	req, _ := http.NewRequest("GET", "http://localhost:8002/hello?name=the-name", nil)

	err := binder.Bind(req, nil, new(jsonRequestParams))
	assert.Error(t, err)

	op2 := parametersForJSONRequestParams("")
	binder = &operationBinder{Parameters: op2, Consumers: map[string]swagger.Consumer{"application/json": swagger.JSONConsumer()}}

	req, _ = http.NewRequest("POST", "http://localhost:8002/hello/1?name=the-name", bytes.NewBuffer([]byte(`{"name":"toby","age":32}`)))
	req.Header.Set("X-Request-Id", "1325959595")

	data := jsonRequestParams{}
	err = binder.Bind(req, swagger.RouteParams([]swagger.RouteParam{{"id", "1"}}), &data)
	assert.Error(t, err)

	req, _ = http.NewRequest("POST", "http://localhost:8002/hello/1?name=the-name", bytes.NewBuffer([]byte(`{"name":"toby","age":32}`)))
	req.Header.Set("Content-Type", "application(")
	data = jsonRequestParams{}
	err = binder.Bind(req, swagger.RouteParams([]swagger.RouteParam{{"id", "1"}}), &data)
	assert.Error(t, err)

	req, _ = http.NewRequest("POST", "http://localhost:8002/hello/1?name=the-name", bytes.NewBuffer([]byte(`{]`)))
	req.Header.Set("Content-Type", "application/json")
	data = jsonRequestParams{}
	err = binder.Bind(req, swagger.RouteParams([]swagger.RouteParam{{"id", "1"}}), &data)
	assert.Error(t, err)

	invalidMultiParam := spec.HeaderParam()
	invalidMultiParam.Name = "tags"
	invalidMultiParam.Type = "array"
	invalidMultiParam.Items = new(spec.Items)
	invalidMultiParam.CollectionFormat = "multi"

	op3 := map[string]spec.Parameter{"Tags": *invalidMultiParam}
	binder = &operationBinder{Parameters: op3, Consumers: map[string]swagger.Consumer{"application/json": swagger.JSONConsumer()}}

	req, _ = http.NewRequest("POST", "http://localhost:8002/hello/1?name=the-name", bytes.NewBuffer([]byte(`{}`)))
	req.Header.Set("Content-Type", "application/json")
	data = jsonRequestParams{}
	err = binder.Bind(req, swagger.RouteParams([]swagger.RouteParam{{"id", "1"}}), &data)
	assert.Error(t, err)

	invalidMultiParam = spec.PathParam()
	invalidMultiParam.Name = "tags"
	invalidMultiParam.Type = "array"
	invalidMultiParam.Items = new(spec.Items)
	invalidMultiParam.CollectionFormat = "multi"

	op4 := map[string]spec.Parameter{"Tags": *invalidMultiParam}
	binder = &operationBinder{Parameters: op4, Consumers: map[string]swagger.Consumer{"application/json": swagger.JSONConsumer()}}

	req, _ = http.NewRequest("POST", "http://localhost:8002/hello/1?name=the-name", bytes.NewBuffer([]byte(`{}`)))
	req.Header.Set("Content-Type", "application/json")
	data = jsonRequestParams{}
	err = binder.Bind(req, swagger.RouteParams([]swagger.RouteParam{{"id", "1"}}), &data)
	assert.Error(t, err)

	invalidInParam := spec.HeaderParam()
	invalidInParam.Name = "tags"
	invalidInParam.Type = "string"
	invalidInParam.In = "invalid"
	op5 := map[string]spec.Parameter{"Tags": *invalidInParam}
	binder = &operationBinder{Parameters: op5, Consumers: map[string]swagger.Consumer{"application/json": swagger.JSONConsumer()}}

	req, _ = http.NewRequest("POST", "http://localhost:8002/hello/1?name=the-name", bytes.NewBuffer([]byte(`{}`)))
	req.Header.Set("Content-Type", "application/json")
	data = jsonRequestParams{}
	err = binder.Bind(req, swagger.RouteParams([]swagger.RouteParam{{"id", "1"}}), &data)
	assert.Error(t, err)
}

func TestRequestBindingForValid(t *testing.T) {

	for _, fmt := range []string{"csv", "pipes", "tsv", "ssv", "multi"} {
		op1 := parametersForJSONRequestParams(fmt)

		binder := &operationBinder{Parameters: op1, Consumers: map[string]swagger.Consumer{"application/json": swagger.JSONConsumer()}}

		lval := []string{"one", "two", "three"}
		queryString := ""
		switch fmt {
		case "multi":
			queryString = strings.Join(lval, "&tags=")
		case "ssv":
			queryString = strings.Join(lval, " ")
		case "pipes":
			queryString = strings.Join(lval, "|")
		case "tsv":
			queryString = strings.Join(lval, "\t")
		default:
			queryString = strings.Join(lval, ",")
		}

		urlStr := "http://localhost:8002/hello/1?name=the-name&tags=" + queryString

		req, _ := http.NewRequest("POST", urlStr, bytes.NewBuffer([]byte(`{"name":"toby","age":32}`)))
		req.Header.Set("Content-Type", "application/json;charset=utf-8")
		req.Header.Set("X-Request-Id", "1325959595")

		data := jsonRequestParams{}
		err := binder.Bind(req, swagger.RouteParams([]swagger.RouteParam{{"id", "1"}}), &data)

		expected := jsonRequestParams{
			ID:        1,
			Name:      "the-name",
			Friend:    friend{"toby", 32},
			RequestID: 1325959595,
			Tags:      []string{"one", "two", "three"},
		}
		assert.NoError(t, err)
		assert.Equal(t, expected, data)
	}

	op1 := parametersForJSONRequestParams("")

	binder := &operationBinder{Parameters: op1, Consumers: map[string]swagger.Consumer{"application/json": swagger.JSONConsumer()}}
	urlStr := "http://localhost:8002/hello/1?name=the-name&tags=one,two,three"
	req, _ := http.NewRequest("POST", urlStr, bytes.NewBuffer([]byte(`{"name":"toby","age":32}`)))
	req.Header.Set("Content-Type", "application/json;charset=utf-8")

	data2 := jsonRequestPtr{}
	err := binder.Bind(req, nil, &data2)

	expected2 := jsonRequestPtr{
		Friend: &friend{"toby", 32},
		Tags:   []string{"one", "two", "three"},
	}
	assert.NoError(t, err)
	assert.Equal(t, *expected2.Friend, *data2.Friend)
	assert.Equal(t, expected2.Tags, data2.Tags)

	req, _ = http.NewRequest("POST", urlStr, bytes.NewBuffer([]byte(`[{"name":"toby","age":32}]`)))
	req.Header.Set("Content-Type", "application/json;charset=utf-8")

	data3 := jsonRequestSlice{}
	err = binder.Bind(req, nil, &data3)

	expected3 := jsonRequestSlice{
		Friend: []friend{{"toby", 32}},
		Tags:   []string{"one", "two", "three"},
	}
	assert.NoError(t, err)
	assert.Equal(t, expected3.Friend, data3.Friend)
	assert.Equal(t, expected3.Tags, data3.Tags)
}

type formRequest struct {
	Name string
	Age  int
}

func parametersForFormUpload() map[string]spec.Parameter {
	nameParam := spec.FormDataParam()
	nameParam.Name = "name"
	nameParam.Type = "string"

	ageParam := spec.FormDataParam()
	ageParam.Name = "age"
	ageParam.Type = "integer"
	ageParam.Format = "int32"

	return map[string]spec.Parameter{"Name": *nameParam, "Age": *ageParam}
}

func TestFormUpload(t *testing.T) {
	params := parametersForFormUpload()
	binder := &operationBinder{Parameters: params, Consumers: map[string]swagger.Consumer{"application/json": swagger.JSONConsumer()}}

	urlStr := "http://localhost:8002/hello"
	req, _ := http.NewRequest("POST", urlStr, bytes.NewBufferString(`name=the-name&age=32`))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	data := formRequest{}
	assert.NoError(t, binder.Bind(req, nil, &data))
	assert.Equal(t, "the-name", data.Name)
	assert.Equal(t, 32, data.Age)

	req, _ = http.NewRequest("POST", urlStr, bytes.NewBufferString(`name=%3&age=32`))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	data = formRequest{}
	assert.Error(t, binder.Bind(req, nil, &data))
}

type fileRequest struct {
	Name string       // body
	File swagger.File // upload
}

func paramsForFileUpload() *operationBinder {
	nameParam := spec.FormDataParam()
	nameParam.Name = "name"
	nameParam.Type = "string"

	fileParam := spec.FormDataParam()
	fileParam.Name = "file"
	fileParam.Type = "file"

	params := map[string]spec.Parameter{"Name": *nameParam, "File": *fileParam}
	return &operationBinder{
		Parameters: params,
		Consumers:  map[string]swagger.Consumer{"application/json": swagger.JSONConsumer()},
	}
}

func TestBindingFileUpload(t *testing.T) {
	binder := paramsForFileUpload()

	body := bytes.NewBuffer(nil)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "plain-jane.txt")
	assert.NoError(t, err)

	part.Write([]byte("the file contents"))
	writer.WriteField("name", "the-name")
	assert.NoError(t, writer.Close())

	urlStr := "http://localhost:8002/hello"
	req, _ := http.NewRequest("POST", urlStr, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	data := fileRequest{}
	assert.NoError(t, binder.Bind(req, nil, &data))
	assert.Equal(t, "the-name", data.Name)
	assert.NotNil(t, data.File)
	assert.NotNil(t, data.File.Header)
	assert.Equal(t, "plain-jane.txt", data.File.Header.Filename)

	bb, err := ioutil.ReadAll(data.File.Data)
	assert.NoError(t, err)
	assert.Equal(t, []byte("the file contents"), bb)

	req, _ = http.NewRequest("POST", urlStr, body)
	req.Header.Set("Content-Type", "application/json")
	data = fileRequest{}
	assert.Error(t, binder.Bind(req, nil, &data))

	req, _ = http.NewRequest("POST", urlStr, body)
	req.Header.Set("Content-Type", "application(")
	data = fileRequest{}
	assert.Error(t, binder.Bind(req, nil, &data))

	body = bytes.NewBuffer(nil)
	writer = multipart.NewWriter(body)
	part, err = writer.CreateFormFile("bad-name", "plain-jane.txt")
	assert.NoError(t, err)

	part.Write([]byte("the file contents"))
	writer.WriteField("name", "the-name")
	assert.NoError(t, writer.Close())
	req, _ = http.NewRequest("POST", urlStr, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	data = fileRequest{}
	assert.Error(t, binder.Bind(req, nil, &data))

	req, _ = http.NewRequest("POST", urlStr, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.MultipartReader()

	data = fileRequest{}
	assert.Error(t, binder.Bind(req, nil, &data))

}
