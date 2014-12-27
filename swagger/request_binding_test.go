package swagger

import (
	"bytes"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strings"
	"testing"

	"github.com/casualjim/go-swagger"
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

func parametersForJSONRequestParams(fmt string) map[string]swagger.Parameter {
	if fmt == "" {
		fmt = "csv"
	}
	nameParam := swagger.QueryParam()
	nameParam.Name = "name"
	nameParam.Type = "string"

	idParam := swagger.PathParam()
	idParam.Name = "id"
	idParam.Type = "integer"
	idParam.Format = "int64"

	friendParam := swagger.BodyParam()
	friendParam.Name = "friend"
	friendParam.Type = "object"

	requestIDParam := swagger.HeaderParam()
	requestIDParam.Name = "X-Request-Id"
	requestIDParam.Type = "string"
	requestIDParam.Extensions = swagger.Extensions(map[string]interface{}{})
	requestIDParam.Extensions.Add("go-name", "RequestID")

	tagsParam := swagger.QueryParam()
	tagsParam.Name = "tags"
	tagsParam.Type = "array"
	tagsParam.Items = new(swagger.Items)
	tagsParam.Items.Type = "string"
	tagsParam.CollectionFormat = fmt

	return map[string]swagger.Parameter{"ID": *idParam, "Name": *nameParam, "RequestID": *requestIDParam, "Friend": *friendParam, "Tags": *tagsParam}
}

func TestRequestBindingForInvalid(t *testing.T) {

	invalidParam := swagger.QueryParam()
	invalidParam.Name = "some"

	op1 := map[string]swagger.Parameter{"Some": *invalidParam}

	binder := &operationBinder{op1, map[string]Consumer{}}
	req, _ := http.NewRequest("GET", "http://localhost:8002/hello?name=the-name", nil)

	err := binder.Bind(req, nil, new(jsonRequestParams))
	assert.Error(t, err)

	op2 := parametersForJSONRequestParams("")
	binder = &operationBinder{op2, map[string]Consumer{"application/json": JSONConsumer()}}

	req, _ = http.NewRequest("POST", "http://localhost:8002/hello/1?name=the-name", bytes.NewBuffer([]byte(`{"name":"toby","age":32}`)))
	req.Header.Set("X-Request-Id", "1325959595")

	data := jsonRequestParams{}
	err = binder.Bind(req, RouteParams([]RouteParam{{"id", "1"}}), &data)
	assert.Error(t, err)

	req, _ = http.NewRequest("POST", "http://localhost:8002/hello/1?name=the-name", bytes.NewBuffer([]byte(`{"name":"toby","age":32}`)))
	req.Header.Set("Content-Type", "application(")
	data = jsonRequestParams{}
	err = binder.Bind(req, RouteParams([]RouteParam{{"id", "1"}}), &data)
	assert.Error(t, err)

	req, _ = http.NewRequest("POST", "http://localhost:8002/hello/1?name=the-name", bytes.NewBuffer([]byte(`{]`)))
	req.Header.Set("Content-Type", "application/json")
	data = jsonRequestParams{}
	err = binder.Bind(req, RouteParams([]RouteParam{{"id", "1"}}), &data)
	assert.Error(t, err)

	invalidMultiParam := swagger.HeaderParam()
	invalidMultiParam.Name = "tags"
	invalidMultiParam.Type = "array"
	invalidMultiParam.Items = new(swagger.Items)
	invalidMultiParam.CollectionFormat = "multi"

	op3 := map[string]swagger.Parameter{"Tags": *invalidMultiParam}
	binder = &operationBinder{op3, map[string]Consumer{"application/json": JSONConsumer()}}

	req, _ = http.NewRequest("POST", "http://localhost:8002/hello/1?name=the-name", bytes.NewBuffer([]byte(`{}`)))
	req.Header.Set("Content-Type", "application/json")
	data = jsonRequestParams{}
	err = binder.Bind(req, RouteParams([]RouteParam{{"id", "1"}}), &data)
	assert.Error(t, err)

	invalidMultiParam = swagger.PathParam()
	invalidMultiParam.Name = "tags"
	invalidMultiParam.Type = "array"
	invalidMultiParam.Items = new(swagger.Items)
	invalidMultiParam.CollectionFormat = "multi"

	op4 := map[string]swagger.Parameter{"Tags": *invalidMultiParam}
	binder = &operationBinder{op4, map[string]Consumer{"application/json": JSONConsumer()}}

	req, _ = http.NewRequest("POST", "http://localhost:8002/hello/1?name=the-name", bytes.NewBuffer([]byte(`{}`)))
	req.Header.Set("Content-Type", "application/json")
	data = jsonRequestParams{}
	err = binder.Bind(req, RouteParams([]RouteParam{{"id", "1"}}), &data)
	assert.Error(t, err)

	invalidInParam := swagger.HeaderParam()
	invalidInParam.Name = "tags"
	invalidInParam.Type = "string"
	invalidInParam.In = "invalid"
	op5 := map[string]swagger.Parameter{"Tags": *invalidInParam}
	binder = &operationBinder{op5, map[string]Consumer{"application/json": JSONConsumer()}}

	req, _ = http.NewRequest("POST", "http://localhost:8002/hello/1?name=the-name", bytes.NewBuffer([]byte(`{}`)))
	req.Header.Set("Content-Type", "application/json")
	data = jsonRequestParams{}
	err = binder.Bind(req, RouteParams([]RouteParam{{"id", "1"}}), &data)
	assert.Error(t, err)
}

func TestRequestBindingForValid(t *testing.T) {

	for _, fmt := range []string{"csv", "pipes", "tsv", "ssv", "multi"} {
		op1 := parametersForJSONRequestParams(fmt)

		binder := &operationBinder{op1, map[string]Consumer{"application/json": JSONConsumer()}}

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
		err := binder.Bind(req, RouteParams([]RouteParam{{"id", "1"}}), &data)

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

	binder := &operationBinder{op1, map[string]Consumer{"application/json": JSONConsumer()}}
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

func TestFormUpload(t *testing.T) {
	nameParam := swagger.FormDataParam()
	nameParam.Name = "name"
	nameParam.Type = "string"

	ageParam := swagger.FormDataParam()
	ageParam.Name = "age"
	ageParam.Type = "integer"
	ageParam.Format = "int32"

	params := map[string]swagger.Parameter{"Name": *nameParam, "Age": *ageParam}
	binder := &operationBinder{params, map[string]Consumer{"application/json": JSONConsumer()}}

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
	Name string // body
	File *File  // upload
}

func TestBindingFileUpload(t *testing.T) {
	nameParam := swagger.FormDataParam()
	nameParam.Name = "name"
	nameParam.Type = "string"

	fileParam := swagger.FormDataParam()
	fileParam.Name = "file"
	fileParam.Type = "file"

	params := map[string]swagger.Parameter{"Name": *nameParam, "File": *fileParam}
	binder := &operationBinder{params, map[string]Consumer{"application/json": JSONConsumer()}}

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
