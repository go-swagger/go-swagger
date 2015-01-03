package validate

import (
	"bytes"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

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
	nameParam := spec.QueryParam("name").Typed("string", "")
	idParam := spec.PathParam("id").Typed("integer", "int64")
	ageParam := spec.QueryParam("age").Typed("integer", "int32")
	scoreParam := spec.QueryParam("score").Typed("number", "float")
	factorParam := spec.QueryParam("factor").Typed("number", "double")

	friendParam := spec.BodyParam("friend", nil)

	requestIDParam := spec.HeaderParam("X-Request-Id").Typed("integer", "int64")
	requestIDParam.Extensions = spec.Extensions(map[string]interface{}{})
	requestIDParam.Extensions.Add("go-name", "RequestID")

	items := new(spec.Items)
	items.Type = "string"
	tagsParam := spec.QueryParam("tags").CollectionOf(items, fmt)

	confirmedParam := spec.QueryParam("confirmed").Typed("boolean", "")
	plannedParam := spec.QueryParam("planned").Typed("string", "date")
	deliveredParam := spec.QueryParam("delivered").Typed("string", "date-time")
	pictureParam := spec.QueryParam("picture").Typed("string", "byte") // base64 encoded during transport

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
		"Picture":   *pictureParam,
	}
}

func parametersForJSONRequestParams(fmt string) map[string]spec.Parameter {
	if fmt == "" {
		fmt = "csv"
	}
	nameParam := spec.QueryParam("name").Typed("string", "")
	idParam := spec.PathParam("id").Typed("integer", "int64")
	friendParam := spec.BodyParam("friend", nil)

	requestIDParam := spec.HeaderParam("X-Request-Id").Typed("integer", "int64")
	requestIDParam.Extensions = spec.Extensions(map[string]interface{}{})
	requestIDParam.Extensions.Add("go-name", "RequestID")

	items := new(spec.Items)
	items.Type = "string"
	tagsParam := spec.QueryParam("tags").CollectionOf(items, fmt)

	return map[string]spec.Parameter{
		"ID":        *idParam,
		"Name":      *nameParam,
		"RequestID": *requestIDParam,
		"Friend":    *friendParam,
		"Tags":      *tagsParam,
	}
}

func TestRequestBindingDefaultValue(t *testing.T) {

	confirmed := true
	name := "thomas"
	friend := map[string]interface{}{"name": "toby", "age": float64(32)}
	id, age, score, factor := int64(7575), int32(348), float32(5.309), float64(37.403)
	requestID := 19394858
	tags := []string{"one", "two", "three"}
	dt1 := time.Date(2014, 8, 9, 0, 0, 0, 0, time.UTC)
	planned := swagger.Date{Time: dt1}
	dt2 := time.Date(2014, 10, 12, 8, 5, 5, 0, time.UTC)
	delivered := swagger.DateTime{Time: dt2}
	// picture := base64.StdEncoding.EncodeToString([]byte("hello"))
	uri, _ := url.Parse("http://localhost:8002/hello")
	defaults := map[string]interface{}{
		"id":           id,
		"age":          age,
		"score":        score,
		"factor":       factor,
		"name":         name,
		"friend":       friend,
		"X-Request-Id": requestID,
		"tags":         tags,
		"confirmed":    confirmed,
		"planned":      planned,
		"delivered":    delivered,
		"picture":      []byte("hello"),
	}
	op2 := parametersForAllTypes("")
	op3 := make(map[string]spec.Parameter)
	for k, p := range op2 {
		p.Default = defaults[p.Name]
		op3[k] = p
	}

	req, _ := http.NewRequest("POST", uri.String(), bytes.NewBuffer([]byte(`{}`)))
	req.Header.Set("Content-Type", "application/json")
	binder := &operationBinder{
		Parameters: op3,
		Consumers:  map[string]swagger.Consumer{"application/json": swagger.JSONConsumer()},
	}

	data := make(map[string]interface{})
	err := binder.Bind(req, swagger.RouteParams(nil), &data)
	assert.NoError(t, err)
	assert.Equal(t, defaults["id"], data["id"])
	assert.Equal(t, name, data["name"])
	// assert.Equal(t, friend, data["friend"])
	assert.Equal(t, requestID, data["X-Request-Id"])
	assert.Equal(t, tags, data["tags"])
	assert.Equal(t, planned, data["planned"])
	assert.Equal(t, delivered, data["delivered"])
	assert.Equal(t, confirmed, data["confirmed"])
	assert.Equal(t, age, data["age"])
	assert.Equal(t, factor, data["factor"])
	assert.Equal(t, score, data["score"])
	assert.Equal(t, "hello", string(data["picture"].([]byte)))
}

func TestRequestBindingForInvalid(t *testing.T) {

	invalidParam := spec.QueryParam("some")

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

	invalidMultiParam := spec.HeaderParam("tags").CollectionOf(new(spec.Items), "multi")
	op3 := map[string]spec.Parameter{"Tags": *invalidMultiParam}
	binder = &operationBinder{Parameters: op3, Consumers: map[string]swagger.Consumer{"application/json": swagger.JSONConsumer()}}

	req, _ = http.NewRequest("POST", "http://localhost:8002/hello/1?name=the-name", bytes.NewBuffer([]byte(`{}`)))
	req.Header.Set("Content-Type", "application/json")
	data = jsonRequestParams{}
	err = binder.Bind(req, swagger.RouteParams([]swagger.RouteParam{{"id", "1"}}), &data)
	assert.Error(t, err)

	invalidMultiParam = spec.PathParam("").CollectionOf(new(spec.Items), "multi")

	op4 := map[string]spec.Parameter{"Tags": *invalidMultiParam}
	binder = &operationBinder{Parameters: op4, Consumers: map[string]swagger.Consumer{"application/json": swagger.JSONConsumer()}}

	req, _ = http.NewRequest("POST", "http://localhost:8002/hello/1?name=the-name", bytes.NewBuffer([]byte(`{}`)))
	req.Header.Set("Content-Type", "application/json")
	data = jsonRequestParams{}
	err = binder.Bind(req, swagger.RouteParams([]swagger.RouteParam{{"id", "1"}}), &data)
	assert.Error(t, err)

	invalidInParam := spec.HeaderParam("tags").Typed("string", "")
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
	nameParam := spec.FormDataParam("name").Typed("string", "")

	ageParam := spec.FormDataParam("age").Typed("integer", "int32")

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
	nameParam := spec.FormDataParam("name").Typed("string", "")

	fileParam := spec.FileParam("file")

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
