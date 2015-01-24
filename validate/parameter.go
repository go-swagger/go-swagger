package validate

import (
	"encoding"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/casualjim/go-swagger"
	"github.com/casualjim/go-swagger/errors"
	"github.com/casualjim/go-swagger/spec"
	"github.com/casualjim/go-swagger/util"
)

func newParamBinder(param spec.Parameter, spec *spec.Swagger, formats formats, formatValidators map[string]FormatValidator) *paramBinder {
	binder := new(paramBinder)
	binder.name = param.Name
	binder.parameter = &param
	binder.formats = formats
	if param.In != "body" {
		binder.validator = newParamValidator(&param, formatValidators)
	} else {
		binder.validator = newSchemaValidator(param.Schema, spec, param.Name, formatValidators)
	}

	return binder
}

type paramBinder struct {
	parameter *spec.Parameter
	formats   formats
	name      string
	validator entityValidator
}

func (p *paramBinder) Type() reflect.Type {
	return p.typeForSchema(p.parameter.Type, p.parameter.Format, p.parameter.Items)
}

func (p *paramBinder) typeForSchema(tpe, format string, items *spec.Items) reflect.Type {
	for _, fmt := range p.formats {
		if tpe == "string" && format == fmt.Name() {
			return fmt.Type()
		}
	}

	switch tpe {

	case "boolean":
		return reflect.TypeOf(true)

	case "string":
		switch strings.Replace(format, "-", "", -1) {
		case "byte":
			return reflect.TypeOf([]byte{})
		case "date":
			return reflect.TypeOf(swagger.Date{})
		case "datetime":
			return reflect.TypeOf(swagger.DateTime{})
		case "uri":
			return reflect.TypeOf(swagger.URI(""))
		case "email":
			return reflect.TypeOf(swagger.Email(""))
		case "hostname":
			return reflect.TypeOf(swagger.Hostname(""))
		case "ipv4":
			return reflect.TypeOf(swagger.IPv4(""))
		case "ipv6":
			return reflect.TypeOf(swagger.IPv6(""))
		case "uuid":
			return reflect.TypeOf(swagger.UUID(""))
		case "uuid3":
			return reflect.TypeOf(swagger.UUID3(""))
		case "uuid4":
			return reflect.TypeOf(swagger.UUID4(""))
		case "uuid5":
			return reflect.TypeOf(swagger.UUID5(""))
		case "isbn":
			return reflect.TypeOf(swagger.ISBN(""))
		case "isbn10":
			return reflect.TypeOf(swagger.ISBN10(""))
		case "isbn13":
			return reflect.TypeOf(swagger.ISBN13(""))
		case "creditcard":
			return reflect.TypeOf(swagger.CreditCard(""))
		case "ssn":
			return reflect.TypeOf(swagger.SSN(""))
		case "hexcolor":
			return reflect.TypeOf(swagger.HexColor(""))
		case "rgbcolor":
			return reflect.TypeOf(swagger.RGBColor(""))
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

func (p *paramBinder) readValue(values interface{}, target reflect.Value) ([]string, bool, error) {
	name, in, cf, tpe := p.parameter.Name, p.parameter.In, p.parameter.CollectionFormat, p.parameter.Type
	if tpe == "array" {
		if cf == "multi" {
			if !p.allowsMulti() {
				return nil, false, errors.InvalidCollectionFormat(name, in, cf)
			}
			return values.(url.Values)[name], false, nil
		}

		v := readSingle(values.(getValue), name)
		return p.readFormattedSliceFieldValue(v, target)
	}

	v := readSingle(values.(getValue), name)
	if v == "" {
		return nil, false, nil
	}
	return []string{v}, false, nil
}

func (p *paramBinder) Bind(request *http.Request, routeParams swagger.RouteParams, consumer swagger.Consumer, target reflect.Value) error {
	// fmt.Println("binding", p.name, "as", p.Type())
	switch p.parameter.In {
	case "query":
		data, custom, err := p.readValue(request.URL.Query(), target)
		if err != nil {
			return err
		}
		if custom {
			return nil
		}

		return p.bindValue(data, target)

	case "header":
		data, custom, err := p.readValue(request.Header, target)
		if err != nil {
			return err
		}
		if custom {
			return nil
		}
		return p.bindValue(data, target)

	case "path":
		data, custom, err := p.readValue(routeParams, target)
		if err != nil {
			return err
		}
		if custom {
			return nil
		}
		return p.bindValue(data, target)

	case "formData":
		mt, err := contentType(request)
		if err != nil {
			return errors.InvalidContentType("", []string{"multipart/form-data", "application/x-www-form-urlencoded"})
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

		if p.parameter.Type == "file" {
			file, header, err := request.FormFile(p.parameter.Name)
			if err != nil {
				return err
			}
			target.Set(reflect.ValueOf(swagger.File{Data: file, Header: header}))
			return nil
		}

		if request.MultipartForm != nil {
			data, custom, err := p.readValue(url.Values(request.MultipartForm.Value), target)
			if err != nil {
				return err
			}
			if custom {
				return nil
			}
			return p.bindValue(data, target)
		}
		data, custom, err := p.readValue(url.Values(request.PostForm), target)
		if err != nil {
			return err
		}
		if custom {
			return nil
		}
		return p.bindValue(data, target)

	case "body":
		newValue := reflect.New(target.Type())
		if err := consumer.Consume(request.Body, newValue.Interface()); err != nil {
			if err == io.EOF && p.parameter.Default != nil {
				target.Set(reflect.ValueOf(p.parameter.Default))
				return nil
			}
			tpe := p.parameter.Type
			if p.parameter.Format != "" {
				tpe = p.parameter.Format
			}
			return errors.InvalidType(p.name, p.parameter.In, tpe, nil)
		}
		target.Set(reflect.Indirect(newValue))
		return nil
	default:
		return errors.New(500, fmt.Sprintf("invalid parameter location %q", p.parameter.In))
	}
}

func (p *paramBinder) bindValue(data []string, target reflect.Value) error {
	if p.parameter.Type == "array" {
		return p.setSliceFieldValue(target, p.parameter.Default, data)
	}
	var d string
	if len(data) > 0 {
		d = data[0]
	}
	return p.setFieldValue(target, p.parameter.Default, d)
}

func (p *paramBinder) setFieldValue(target reflect.Value, defaultValue interface{}, data string) error {
	tpe := p.parameter.Type
	if p.parameter.Format != "" {
		tpe = p.parameter.Format
	}

	if data == "" && p.parameter.Required && p.parameter.Default == nil {
		return errors.Required(p.name, p.parameter.In)
	}

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
		value := data
		if value == "" {
			value = defVal.String()
		}
		// validate string
		target.SetString(value)

	case reflect.Ptr:
		if data == "" && defVal.Kind() == reflect.Ptr {
			target.Set(defVal)
			return nil
		}
		newVal := reflect.New(target.Type().Elem())
		if err := p.setFieldValue(reflect.Indirect(newVal), defVal, data); err != nil {
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

func (p *paramBinder) readFormattedSliceFieldValue(data string, target reflect.Value) ([]string, bool, error) {
	ok, err := p.tryUnmarshaler(target, p.parameter.Default, data)
	if err != nil {
		return nil, true, err
	}
	if ok {
		return nil, true, nil
	}

	return split(data, p.parameter.CollectionFormat), false, nil
}

func (p *paramBinder) setSliceFieldValue(target reflect.Value, defaultValue interface{}, data []string) error {
	if len(data) == 0 && p.parameter.Required && p.parameter.Default == nil {
		return errors.Required(p.name, p.parameter.In)
	}
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
