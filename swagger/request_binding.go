package swagger

import (
	"encoding"
	"errors"
	"fmt"
	"mime"
	"mime/multipart"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/casualjim/go-swagger"
	"github.com/casualjim/go-swagger/swagger/util"
)

var requestBinderType = reflect.TypeOf(new(RequestBinder)).Elem()
var textUnmarshalType = reflect.TypeOf(new(encoding.TextUnmarshaler)).Elem()

// File represents an uploaded file.
type File struct {
	Data   multipart.File
	Header *multipart.FileHeader
}

// RequestBinder is an interface for types that want to take charge of customizing the binding process
// or want to sidestep the reflective binding of values.
type RequestBinder interface {
	BindRequest(*http.Request, RouteParams) error
}

type operationBinder struct {
	Parameters []swagger.Parameter
	Consumers  map[string]Consumer
}

func (o *operationBinder) Bind(request *http.Request, routeParams RouteParams, data interface{}) error {
	val := reflect.Indirect(reflect.ValueOf(data))

	for _, param := range o.Parameters {

		fieldName := fieldNameFromParam(&param)
		target := val.FieldByName(fieldName)

		if !target.IsValid() {
			return fmt.Errorf("parameter name %q is an unknown field", fieldName)
		}

		if err := o.setParamValue(target, &param, request, routeParams); err != nil {
			return err
		}
	}

	return nil
}

const defaultMaxMemory = 32 << 20

func contentType(req *http.Request) (string, error) {
	ct := req.Header.Get("Content-Type")
	if ct == "" {
		ct = "application/octet-stream"
	}
	mt, _, err := mime.ParseMediaType(ct)
	if err != nil {
		return "", err
	}
	return mt, nil
}

func (o *operationBinder) setParamValue(target reflect.Value, param *swagger.Parameter, request *http.Request, routeParams RouteParams) error {

	switch param.In {
	case "query":
		return o.setValue(target, request.URL.Query(), param, true)

	case "header":
		return o.setValue(target, request.Header, param, false)

	case "path":
		return o.setValue(target, routeParams, param, false)

	case "formData":
		mt, err := contentType(request)
		if err != nil {
			return err
		}
		if mt != "multipart/form-data" && mt != "application/x-www-form-urlencoded" {
			return errors.New("invalid content type, should be \"multipart/form-data\" or \"application/x-www-form-urlencoded\"")
		}
		if mt == "multipart/form-data" {
			if err := request.ParseMultipartForm(defaultMaxMemory); err != nil {
				return err
			}
		}
		if err := request.ParseForm(); err != nil {
			return err
		}

		if param.Type == "file" {
			file, header, err := request.FormFile(param.Name)
			if err != nil {
				return err
			}
			target.Set(reflect.ValueOf(&File{file, header}))
			return nil
		}

		if request.MultipartForm != nil {
			return o.setValue(target, url.Values(request.MultipartForm.Value), param, true)
		}
		return o.setValue(target, request.PostForm, param, true)

	case "body":
		mt, err := contentType(request)
		if err != nil {
			return err
		}
		if consumer, ok := o.Consumers[mt]; ok {
			newValue := reflect.New(target.Type())
			if err := consumer.Consume(request.Body, newValue.Interface()); err != nil {
				return err
			}
			target.Set(reflect.Indirect(newValue))
			return nil
		}
		return fmt.Errorf("no consumer for the body as %q", mt)
	default:
		return fmt.Errorf("invalid parameter location %q", param.In)
	}
}

func (o *operationBinder) setValue(target reflect.Value, values interface{}, param *swagger.Parameter, multi bool) error {
	if param.Type == "array" {
		if param.CollectionFormat == "multi" {
			if !multi {
				return fmt.Errorf("the collection format %q is not supported for a %s param", param.CollectionFormat, param.In)
			}
			return setSliceFieldValue(target, values.(url.Values)[param.Name], param.Default)
		}
		return setFormattedSliceFieldValue(target, readSingle(values.(interface {
			Get(string) string
		}), param.Name), param.CollectionFormat, param.Default)
	}
	return setFieldValue(target, readSingle(values.(interface {
		Get(string) string
	}), param.Name), param.Default)
}

func readSingle(from interface {
	Get(string) string
}, name string) string {
	return from.Get(name)
}

func fieldNameFromParam(param *swagger.Parameter) string {
	if nm, ok := param.Extensions.GetString("GO-NAME"); ok {
		return nm
	}
	return util.ToGoName(param.Name)
}

func setFieldValue(target reflect.Value, data string, defaultValue interface{}) error {
	ok, err := tryUnmarshaler(target, data, defaultValue)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}

	defVal := reflect.Zero(target.Type())
	if defaultValue != nil {
		defVal = reflect.ValueOf(defaultValue)
	}

	switch target.Kind() {
	case reflect.Bool:
		if data == "" {
			target.SetBool(defVal.Bool())
			return nil
		}
		target.SetBool(contains(evaluatesAsTrue, data))

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if data == "" {
			target.SetInt(defVal.Int())
			return nil
		}
		i, err := strconv.ParseInt(data, 10, 64)
		if err != nil {
			return err
		}
		target.SetInt(i)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if data == "" {
			target.SetUint(defVal.Uint())
			return nil
		}
		u, err := strconv.ParseUint(data, 10, 64)
		if err != nil {
			return err
		}
		target.SetUint(u)

	case reflect.Float32, reflect.Float64:
		if data == "" {
			target.SetFloat(defVal.Float())
			return nil
		}
		f, err := strconv.ParseFloat(data, 64)
		if err != nil {
			return err
		}
		target.SetFloat(f)

	case reflect.String:
		if data == "" {
			target.SetString(defVal.String())
			return nil
		}
		target.SetString(data)

	case reflect.Ptr:
		if data == "" && defVal.Kind() == reflect.Ptr {
			target.Set(defVal)
			return nil
		}
		newVal := reflect.New(target.Type().Elem())
		err := setFieldValue(reflect.Indirect(newVal), data, defaultValue)
		if err != nil {
			return err
		}
		target.Set(newVal)
	default:
		return fmt.Errorf("Don't know how to convert %q to a %s (kind %s)", data, target.Type(), target.Kind())
	}
	return nil
}

func tryUnmarshaler(target reflect.Value, data string, defaultValue interface{}) (bool, error) {
	// When a type implements encoding.TextUnmarshaler we'll use that instead of reflecting some more
	if reflect.PtrTo(target.Type()).Implements(textUnmarshalType) {
		if defaultValue != nil && len(data) == 0 {
			target.Set(reflect.ValueOf(defaultValue))
			return true, nil
		}
		value := reflect.New(target.Type())
		if err := value.Interface().(encoding.TextUnmarshaler).UnmarshalText([]byte(data)); err != nil {
			return true, err
		}
		target.Set(reflect.Indirect(value))
		return true, nil
	}
	return false, nil
}

func setFormattedSliceFieldValue(target reflect.Value, data, format string, defaultValue interface{}) error {
	ok, err := tryUnmarshaler(target, data, defaultValue)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}

	return setSliceFieldValue(target, split(data, format), defaultValue)
}

func setSliceFieldValue(target reflect.Value, data []string, defaultValue interface{}) error {
	defVal := reflect.Zero(target.Type())
	if defaultValue != nil {
		defVal = reflect.ValueOf(defaultValue)
	}
	if len(data) == 0 {
		target.Set(defVal)
		return nil
	}

	sz := len(data)
	value := reflect.MakeSlice(reflect.SliceOf(target.Type().Elem()), sz, sz)

	for i := 0; i < sz; i++ {
		if err := setFieldValue(value.Index(i), data[i], nil); err != nil {
			return err
		}
	}

	target.Set(value)

	return nil
}

func split(data, format string) []string {
	if data == "" {
		return nil
	}
	var sep string
	switch format {
	case "ssv":
		sep = " "
	case "tsv":
		sep = "\t"
	case "pipes":
		sep = "|"
	case "multi":
		return nil
	default:
		sep = ","
	}
	var result []string
	for _, s := range strings.Split(data, sep) {
		if ts := strings.TrimSpace(s); ts != "" {
			result = append(result, ts)
		}
	}
	return result
}

var evaluatesAsTrue = []string{"true", "1", "yes", "ok", "y", "on", "selected", "checked", "t", "enabled"}

func contains(s []string, e string) bool {
	for _, a := range s {
		if strings.EqualFold(a, e) {
			return true
		}
	}
	return false
}
