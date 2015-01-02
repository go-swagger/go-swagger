package validate

import (
	"encoding"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/casualjim/go-swagger"
	"github.com/casualjim/go-swagger/errors"
	"github.com/casualjim/go-swagger/httputils"
	"github.com/casualjim/go-swagger/spec"
	"github.com/casualjim/go-swagger/util"
)

var textUnmarshalType = reflect.TypeOf(new(encoding.TextUnmarshaler)).Elem()

type operationBinder struct {
	Parameters map[string]spec.Parameter
	Consumers  map[string]swagger.Consumer
	Formats    map[string]map[string]reflect.Type
}

func (o *operationBinder) Bind(request *http.Request, routeParams swagger.RouteParams, data interface{}) error {
	val := reflect.Indirect(reflect.ValueOf(data))
	isMap := val.Kind() == reflect.Map
	for fieldName, param := range o.Parameters {
		name := fieldName
		var target reflect.Value
		if !isMap {
			target = val.FieldByName(name)
		}

		if isMap {
			name = param.Name
			tpe := o.typeForSchema(param.Type, param.Format, param.Items)
			if tpe == nil {
				continue
			}
			target = reflect.Indirect(reflect.New(tpe))
		}

		if !target.IsValid() {
			return fmt.Errorf("parameter name %q is an unknown field", name)
		}

		if err := o.setParamValue(target, &param, request, routeParams); err != nil {
			return err
		}
		if isMap {
			val.SetMapIndex(reflect.ValueOf(param.Name), target)
		}
	}

	return nil
}

func (o *operationBinder) typeForSchema(tpe, format string, items *spec.Items) reflect.Type {
	switch tpe {
	case "boolean":
		return reflect.TypeOf(true)
	case "string":
		switch format {
		case "byte":
			return reflect.TypeOf(byte(1))
		case "date":
			return reflect.TypeOf(swagger.Date{})
		case "date-time":
			return reflect.TypeOf(swagger.DateTime{})
		default:
			if fmts, ok := o.Formats[tpe]; ok {
				if tp, ok := fmts[format]; ok {
					return tp
				}
			}
			return reflect.TypeOf("")
		}
	case "integer":
		switch format {
		case "int32":
			return reflect.TypeOf(int32(0))
		case "int64":
			return reflect.TypeOf(int64(0))
		}
	case "number":
		switch format {
		case "float":
			return reflect.TypeOf(float32(0))
		case "double":
			return reflect.TypeOf(float64(0))
		}
	case "array":
		itemsType := o.typeForSchema(items.Type, items.Format, nil)
		if itemsType == nil {
			return nil
		}
		return reflect.MakeSlice(reflect.SliceOf(itemsType), 0, 0).Type()
	case "file":
		return reflect.TypeOf(swagger.File{})
	case "object":
		return reflect.TypeOf(map[string]interface{}{})
	}
	return nil
}

const defaultMaxMemory = 32 << 20

func contentType(req *http.Request) (string, error) {
	mt, _, err := httputils.ContentType(req.Header)
	if err != nil {
		return "", err
	}
	return mt, nil
}

func (o *operationBinder) setParamValue(target reflect.Value, param *spec.Parameter, request *http.Request, routeParams swagger.RouteParams) error {
	tgt := target

	switch param.In {
	case "query":
		return o.setValue(tgt, request.URL.Query(), param, true)

	case "header":
		return o.setValue(tgt, request.Header, param, false)

	case "path":
		return o.setValue(tgt, routeParams, param, false)

	case "formData":
		mt, err := contentType(request)
		if err != nil {
			return err
		}
		if mt != "multipart/form-data" && mt != "application/x-www-form-urlencoded" {
			return errors.InvalidContentType(mt, []string{"multipart/form-data", "application/x-www-form-urlencoded"})
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
			tgt.Set(reflect.ValueOf(swagger.File{Data: file, Header: header}))
			return nil
		}

		if request.MultipartForm != nil {
			return o.setValue(tgt, url.Values(request.MultipartForm.Value), param, true)
		}
		return o.setValue(tgt, request.PostForm, param, true)

	case "body":
		mt, err := contentType(request)
		if err != nil {
			return err
		}
		if consumer, ok := o.Consumers[mt]; ok {
			newValue := reflect.New(tgt.Type())
			if err := consumer.Consume(request.Body, newValue.Interface()); err != nil {
				return err
			}
			tgt.Set(reflect.Indirect(newValue))
			return nil
		}
		return fmt.Errorf("no consumer for the body as %q", mt)
	default:
		return fmt.Errorf("invalid parameter location %q", param.In)
	}
}

func (o *operationBinder) setValue(target reflect.Value, values interface{}, param *spec.Parameter, multi bool) error {
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
		target.SetBool(util.ContainsStringsCI(evaluatesAsTrue, data))

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
