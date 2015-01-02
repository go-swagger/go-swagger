package reflection

// import (
// 	"math"
// 	"reflect"
// )

// // IsInteger returns true when the value is an integer number
// // or convertible to an integer number without precision loss
// func IsInteger(value reflect.Value) bool {
// 	kind := value.Kind()
// 	switch {
// 	case kind >= reflect.Int && kind <= reflect.Uint64:
// 		return true
// 	case kind >= reflect.Float32 && kind <= reflect.Float64:
// 		f := value.Float()
// 		if math.IsNaN(f) || math.IsInf(f, 0) {
// 			return false
// 		}
// 		return f == float64(int64(f)) || f == float64(uint64(f))
// 	case kind == reflect.Ptr:
// 		return IsInteger(reflect.Indirect(value))
// 	default:
// 		return false
// 	}
// }

// // IsNumeric returns true when the value is a numeric value
// func IsNumeric(value reflect.Value) bool {
// 	kind := value.Kind()
// 	switch {
// 	case kind >= reflect.Int && kind <= reflect.Uint64:
// 		return true
// 	case kind >= reflect.Float32 && kind <= reflect.Float64:
// 		return true
// 	case kind == reflect.Ptr:
// 		return IsNumeric(reflect.Indirect(value))
// 	default:
// 		return false
// 	}
// }

// var simpleTypes = map[string]struct{}{
// 	"bool":       struct{}{},
// 	"uint":       struct{}{},
// 	"uint8":      struct{}{},
// 	"uint16":     struct{}{},
// 	"uint32":     struct{}{},
// 	"uint64":     struct{}{},
// 	"int":        struct{}{},
// 	"int8":       struct{}{},
// 	"int16":      struct{}{},
// 	"int32":      struct{}{},
// 	"int64":      struct{}{},
// 	"float32":    struct{}{},
// 	"float64":    struct{}{},
// 	"string":     struct{}{},
// 	"complex64":  struct{}{},
// 	"complex128": struct{}{},
// 	"byte":       struct{}{},
// 	"rune":       struct{}{},
// 	"uintptr":    struct{}{},
// 	"error":      struct{}{},
// 	"Time":       struct{}{},
// }

// // IsSimpleType returns true when the value is a simple type.
// // simple types are bools, uints, ints, floats, strings, complex numbers, bytes, runes, errors and time.Time
// func IsSimpleType(value reflect.Value) (ok bool) {
// 	if reflect.Ptr == value.Type().Kind() {
// 		return IsSimpleType(reflect.Indirect(value))
// 	}
// 	_, ok = simpleTypes[value.Type().Name()]
// 	return
// }
