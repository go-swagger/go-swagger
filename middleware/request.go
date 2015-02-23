package middleware

import (
	"net/http"
	"reflect"

	"github.com/casualjim/go-swagger"
	"github.com/casualjim/go-swagger/errors"
	"github.com/casualjim/go-swagger/spec"
	"github.com/casualjim/go-swagger/strfmt"
	"github.com/casualjim/go-swagger/validate"
)

// RequestBinder binds and validates the data from a http request
type untypedRequestBinder struct {
	Spec         *spec.Swagger
	Parameters   map[string]spec.Parameter
	Formats      strfmt.Registry
	paramBinders map[string]*untypedParamBinder
}

// NewRequestBinder creates a new binder for reading a request.
func newUntypedRequestBinder(parameters map[string]spec.Parameter, spec *spec.Swagger, formats strfmt.Registry) *untypedRequestBinder {
	binders := make(map[string]*untypedParamBinder)
	for fieldName, param := range parameters {
		binders[fieldName] = newUntypedParamBinder(param, spec, formats)
	}
	return &untypedRequestBinder{
		Parameters:   parameters,
		paramBinders: binders,
		Spec:         spec,
		Formats:      formats,
	}
}

// Bind perform the databinding and validation
func (o *untypedRequestBinder) Bind(request *http.Request, routeParams RouteParams, consumer swagger.Consumer, data interface{}) *validate.Result {
	val := reflect.Indirect(reflect.ValueOf(data))
	isMap := val.Kind() == reflect.Map
	result := new(validate.Result)

	for fieldName, param := range o.Parameters {
		binder := o.paramBinders[fieldName]

		var target reflect.Value
		if !isMap {
			binder.Name = fieldName
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
			result.AddErrors(errors.New(500, "parameter name %q is an unknown field", binder.Name))
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
