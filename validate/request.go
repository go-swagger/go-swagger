package validate

import (
	"encoding"
	"mime"
	"net/http"
	"reflect"
	"strings"

	"github.com/casualjim/go-swagger"
	"github.com/casualjim/go-swagger/errors"
	"github.com/casualjim/go-swagger/spec"
	"github.com/casualjim/go-swagger/strfmt"
)

var textUnmarshalType = reflect.TypeOf(new(encoding.TextUnmarshaler)).Elem()

// RequestBinder binds and validates the data from a http request
type RequestBinder struct {
	Spec         *spec.Swagger
	Parameters   map[string]spec.Parameter
	Formats      strfmt.Registry
	paramBinders map[string]*paramBinder
}

// NewRequestBinder creates a new binder for reading a request.
func NewRequestBinder(parameters map[string]spec.Parameter, spec *spec.Swagger, formats strfmt.Registry) *RequestBinder {
	binders := make(map[string]*paramBinder)
	for fieldName, param := range parameters {
		binders[fieldName] = newParamBinder(param, spec, formats)
	}
	return &RequestBinder{
		Parameters:   parameters,
		paramBinders: binders,
		Spec:         spec,
		Formats:      formats,
	}
}

// Bind perform the databinding and validation
func (o *RequestBinder) Bind(request *http.Request, routeParams swagger.RouteParams, consumer swagger.Consumer, data interface{}) *Result {
	val := reflect.Indirect(reflect.ValueOf(data))
	isMap := val.Kind() == reflect.Map
	result := new(Result)

	for fieldName, param := range o.Parameters {
		binder := o.paramBinders[fieldName]
		// fmt.Println("binding", binder.name, "from", param.In, "as", binder.Type()) //, "with", binder.validator)

		var target reflect.Value
		if !isMap {
			binder.name = fieldName
			target = val.FieldByName(fieldName)
		}

		if isMap {
			tpe := binder.Type()
			if tpe == nil {
				if param.Schema.Type.Contains("array") {
					tpe = reflect.TypeOf([]interface{}{})
				} else {
					tpe = reflect.TypeOf(map[string]interface{}{})
				}
			}
			target = reflect.Indirect(reflect.New(tpe))

		}

		if !target.IsValid() {
			result.AddErrors(errors.New(500, "parameter name %q is an unknown field", binder.name))
			continue
		}

		if err := binder.Bind(request, routeParams, consumer, target); err != nil {
			switch err.(type) {
			case *errors.Validation:
				result.AddErrors(err.(*errors.Validation))
			case errors.Error:
				result.AddErrors(err.(errors.Error))
			default:
				result.AddErrors(errors.New(500, err.Error()))
			}
			continue
		}

		if binder.validator != nil {
			result.Merge(binder.validator.Validate(target.Interface()))
		}

		if isMap {
			val.SetMapIndex(reflect.ValueOf(param.Name), target)
		}
	}

	return result
}

const defaultMaxMemory = 32 << 20

func contentType(req *http.Request) (string, error) {
	ct := req.Header.Get("Content-Type")
	orig := ct
	if ct == "" {
		ct = "application/octect-stream"
	}

	mt, _, err := mime.ParseMediaType(ct)
	if err != nil {
		return "", errors.NewParseError("Content-Type", "header", orig, err)
	}

	return mt, nil
}

func readSingle(from getValue, name string) string {
	return from.Get(name)
}

var evaluatesAsTrue = []string{"true", "1", "yes", "ok", "y", "on", "selected", "checked", "t", "enabled"}

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
