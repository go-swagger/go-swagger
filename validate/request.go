package validate

import (
	"encoding"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/casualjim/go-swagger"
	"github.com/casualjim/go-swagger/errors"
	"github.com/casualjim/go-swagger/httputils"
	"github.com/casualjim/go-swagger/spec"
)

var textUnmarshalType = reflect.TypeOf(new(encoding.TextUnmarshaler)).Elem()

type formats map[string]map[string]reflect.Type

// RequestBinder binds and validates the data from a http request
type RequestBinder struct {
	Parameters map[string]spec.Parameter
	Consumer   swagger.Consumer
	Formats    formats
}

// Bind perform the databinding and validation
func (o *RequestBinder) Bind(request *http.Request, routeParams swagger.RouteParams, data interface{}) errors.Error {
	val := reflect.Indirect(reflect.ValueOf(data))
	isMap := val.Kind() == reflect.Map

	for fieldName, param := range o.Parameters {
		binder := new(paramBinder)
		binder.name = fieldName
		binder.parameter = &param
		binder.formats = o.Formats

		var target reflect.Value
		if !isMap {
			target = val.FieldByName(fieldName)
		}

		if isMap {
			binder.name = param.Name
			tpe := binder.Type()
			if tpe == nil {
				continue
			}
			target = reflect.Indirect(reflect.New(tpe))
		}

		if !target.IsValid() {
			return errors.New(500, fmt.Sprintf("parameter name %q is an unknown field", binder.name))
		}

		if err := binder.Bind(request, routeParams, o.Consumer, target); err != nil {
			switch err.(type) {
			case *errors.Validation:
				return err.(*errors.Validation)
			case errors.Error:
				return err.(errors.Error)
			default:
				return errors.New(500, err.Error())
			}
		}

		if isMap {
			val.SetMapIndex(reflect.ValueOf(param.Name), target)
		}
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
