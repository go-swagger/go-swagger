package swagger

import (
	"net/http"
	"reflect"
)

var requestBinderType = reflect.TypeOf(new(RequestBinder)).Elem()

// RequestBinder is an interface for types that want to take charge of customizing the binding process
// or want to sidestep the reflective binding of values.
type RequestBinder interface {
	BindRequest(*http.Request, RouteParams) error
}
