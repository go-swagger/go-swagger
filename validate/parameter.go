package validate

import (
	"encoding"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strconv"

	"github.com/casualjim/go-swagger"
	"github.com/casualjim/go-swagger/errors"
	"github.com/casualjim/go-swagger/spec"
	"github.com/casualjim/go-swagger/util"
)

type paramBinder struct {
	request     *http.Request
	routeParams swagger.RouteParams
	target      reflect.Value
	parameter   *spec.Parameter
	formats     formats
	name        string
	consumers   map[string]swagger.Consumer
}

func (p *paramBinder) Type() reflect.Type {
	return p.typeForSchema(p.parameter.Type, p.parameter.Format, p.parameter.Items)
}
func (p *paramBinder) typeForSchema(tpe, format string, items *spec.Items) reflect.Type {
	if fmts, ok := p.formats[tpe]; ok {
		if tp, ok := fmts[format]; ok {
			return tp
		}
	}

	switch tpe {

	case "boolean":
		return reflect.TypeOf(true)

	case "string":
		switch format {
		case "byte":
			return reflect.TypeOf([]byte{})
		case "date":
			return reflect.TypeOf(swagger.Date{})
		case "date-time":
			return reflect.TypeOf(swagger.DateTime{})
		default:
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
		if items == nil {
			return nil
		}
		itemsType := p.typeForSchema(items.Type, items.Format, items.Items)
		if itemsType == nil {
			return nil
		}
		return reflect.MakeSlice(reflect.SliceOf(itemsType), 0, 0).Type()

	case "file":
		return reflect.TypeOf(&swagger.File{}).Elem()

	case "object":
		return reflect.TypeOf(map[string]interface{}{})
	}
	return nil
}

func (p *paramBinder) allowsMulti() bool {
	return p.parameter.In == "query" || p.parameter.In == "formData"
}

type getValue interface {
	Get(string) string
}

func (p *paramBinder) readValue(values interface{}) ([]string, bool, error) {
	name, in, cf, tpe := p.parameter.Name, p.parameter.In, p.parameter.CollectionFormat, p.parameter.Type
	if tpe == "array" {
		if cf == "multi" {
			if !p.allowsMulti() {
				return nil, false, errors.InvalidCollectionFormat(name, in, cf)
			}
			return values.(url.Values)[name], false, nil
		}

		v := readSingle(values.(getValue), name)
		return p.readFormattedSliceFieldValue(v)
	}

	v := readSingle(values.(getValue), name)
	if v == "" {
		return nil, false, nil
	}
	return []string{v}, false, nil
}

func (p *paramBinder) Bind() error {

	switch p.parameter.In {
	case "query":
		data, custom, err := p.readValue(p.request.URL.Query())
		if err != nil {
			return err
		}
		if custom {
			return nil
		}

		return p.bindValue(data)

	case "header":
		data, custom, err := p.readValue(p.request.Header)
		if err != nil {
			return err
		}
		if custom {
			return nil
		}
		return p.bindValue(data)

	case "path":
		data, custom, err := p.readValue(p.routeParams)
		if err != nil {
			return err
		}
		if custom {
			return nil
		}
		return p.bindValue(data)

	case "formData":
		mt, err := contentType(p.request)
		if err != nil {
			return errors.InvalidContentType("", []string{"multipart/form-data", "application/x-www-form-urlencoded"})
		}
		if mt != "multipart/form-data" && mt != "application/x-www-form-urlencoded" {
			return errors.InvalidContentType(mt, []string{"multipart/form-data", "application/x-www-form-urlencoded"})
		}
		if mt == "multipart/form-data" {
			if err := p.request.ParseMultipartForm(defaultMaxMemory); err != nil {
				return err
			}
		}
		if err := p.request.ParseForm(); err != nil {
			return err
		}

		if p.parameter.Type == "file" {
			file, header, err := p.request.FormFile(p.parameter.Name)
			if err != nil {
				return err
			}
			p.target.Set(reflect.ValueOf(swagger.File{Data: file, Header: header}))
			return nil
		}

		if p.request.MultipartForm != nil {
			data, custom, err := p.readValue(url.Values(p.request.MultipartForm.Value))
			if err != nil {
				return err
			}
			if custom {
				return nil
			}
			return p.bindValue(data)
		}
		data, custom, err := p.readValue(url.Values(p.request.PostForm))
		if err != nil {
			return err
		}
		if custom {
			return nil
		}
		return p.bindValue(data)

	case "body":
		mt, err := contentType(p.request)
		if err != nil {
			return err
		}
		if consumer, ok := p.consumers[mt]; ok {
			newValue := reflect.New(p.target.Type())
			if err := consumer.Consume(p.request.Body, newValue.Interface()); err != nil {
				tpe := p.parameter.Type
				if p.parameter.Format != "" {
					tpe = p.parameter.Format
				}
				return errors.InvalidType(p.name, p.parameter.In, tpe, nil)
			}
			p.target.Set(reflect.Indirect(newValue))
			return nil
		}
		var names []string
		for k := range p.consumers {
			names = append(names, k)
		}
		return errors.InvalidContentType(mt, names)
	default:
		return errors.New(500, fmt.Sprintf("invalid parameter location %q", p.parameter.In))
	}
}

func (p *paramBinder) bindValue(data []string) error {
	if p.parameter.Type == "array" {
		return p.setSliceFieldValue(p.target, p.parameter.Default, data)
	}
	var d string
	if len(data) > 0 {
		d = data[0]
	}
	return p.setFieldValue(p.target, p.parameter.Default, d)
}

func (p *paramBinder) setFieldValue(target reflect.Value, defaultValue interface{}, data string) error {
	tpe := p.parameter.Type
	if p.parameter.Format != "" {
		tpe = p.parameter.Format
	}
	// fmt.Println("The target is of kind", target.Kind())

	ok, err := p.tryUnmarshaler(target, defaultValue, data)
	if err != nil {
		return errors.InvalidType(p.name, p.parameter.In, tpe, data)
	}
	if ok {
		return nil
	}

	defVal := reflect.Zero(target.Type())
	if defaultValue != nil {
		defVal = reflect.ValueOf(defaultValue)
	}

	if tpe == "byte" {
		if data == "" {
			target.SetBytes(defVal.Bytes())
			return nil
		}

		b, err := base64.StdEncoding.DecodeString(data)
		if err != nil {
			b, err = base64.URLEncoding.DecodeString(data)
			if err != nil {
				return errors.InvalidType(p.name, p.parameter.In, tpe, data)
			}
		}
		target.SetBytes(b)
		return nil
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
			return errors.InvalidType(p.name, p.parameter.In, tpe, data)
		}
		if target.OverflowInt(i) {
			return errors.InvalidType(p.name, p.parameter.In, tpe, data)
		}
		target.SetInt(i)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if data == "" {
			target.SetUint(defVal.Uint())
			return nil
		}
		u, err := strconv.ParseUint(data, 10, 64)
		if err != nil {
			return errors.InvalidType(p.name, p.parameter.In, tpe, data)
		}
		if target.OverflowUint(u) {
			return errors.InvalidType(p.name, p.parameter.In, tpe, data)
		}
		target.SetUint(u)

	case reflect.Float32, reflect.Float64:
		if data == "" {
			target.SetFloat(defVal.Float())
			return nil
		}
		f, err := strconv.ParseFloat(data, 64)
		if err != nil {
			return errors.InvalidType(p.name, p.parameter.In, tpe, data)
		}
		if target.OverflowFloat(f) {
			return errors.InvalidType(p.name, p.parameter.In, tpe, data)
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
		err := p.setFieldValue(reflect.Indirect(newVal), defVal, data)
		if err != nil {
			return err
		}
		target.Set(newVal)
	default:
		return errors.InvalidType(p.name, p.parameter.In, tpe, data)
	}
	return nil
}

func (p *paramBinder) tryUnmarshaler(target reflect.Value, defaultValue interface{}, data string) (bool, error) {
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

func (p *paramBinder) readFormattedSliceFieldValue(data string) ([]string, bool, error) {
	ok, err := p.tryUnmarshaler(p.target, p.parameter.Default, data)
	if err != nil {
		return nil, true, err
	}
	if ok {
		return nil, true, nil
	}

	return split(data, p.parameter.CollectionFormat), false, nil
}

func (p *paramBinder) setSliceFieldValue(target reflect.Value, defaultValue interface{}, data []string) error {
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
		if err := p.setFieldValue(value.Index(i), nil, data[i]); err != nil {
			return err
		}
	}

	target.Set(value)

	return nil
}
