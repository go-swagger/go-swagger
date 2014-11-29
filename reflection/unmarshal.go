package reflection

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
)

var (
	unmarshallerType = reflect.TypeOf(new(MapUnmarshaller)).Elem()
)

// MapUnmarshaller is an interface for things that need to customize the unmarshalling process
type MapUnmarshaller interface {
	UnmarshalMap(interface{}) error
}

// UnmarshalMapRecursed converts the provided map to the target interface, but skips the
// initial interface check
func UnmarshalMapRecursed(data map[string]interface{}, target interface{}) error {
	return unmarshalMap(data, target, true)
}

type fieldInfo struct {
	Descriptor reflect.StructField
	ByValue    bool
}

// UnmarshalMap converts the provided map to the target interface
func UnmarshalMap(data map[string]interface{}, target interface{}) error {
	return unmarshalMap(data, target, false)
}
func unmarshalMap(data map[string]interface{}, target interface{}, skipCheck bool) error {
	targetType := reflect.TypeOf(target)
	if targetType.Kind() != reflect.Ptr {
		return fmt.Errorf("target must be a pointer")
	}

	if !skipCheck && targetType.Implements(unmarshallerType) {
		return target.(MapUnmarshaller).UnmarshalMap(data)
	}

	val := reflect.Indirect(reflect.ValueOf(target))
	tpe := val.Type()

	fieldNameMap := fieldTagNameMap(tpe)
	for k, v := range data {
		if targetDes, ok := fieldNameMap[k]; ok {
			if err := convertValue(targetDes.Name, k, reflect.ValueOf(v), val.FieldByName(targetDes.Name)); err != nil {
				return err
			}
		}
	}
	return nil
}

func convertValue(name, key string, source, target reflect.Value) error {
	if !target.IsValid() {
		target.Set(reflect.Zero(target.Type()))
		return nil
	}

	if !target.CanSet() {
		return fmt.Errorf("%s (key %q) %#v must be addressable", name, key, target)
	}

	ptr := reflect.PtrTo(target.Type())
	if ptr.AssignableTo(unmarshallerType) {
		value := reflect.New(target.Type())
		if err := value.Interface().(MapUnmarshaller).UnmarshalMap(MarshalMap(source.Interface())); err != nil {
			return err
		}
		target.Set(reflect.Indirect(value))
		return nil
	}

	switch target.Kind() {
	case reflect.Interface:
		return convertInterface(name, key, source, target)
	case reflect.Bool:
		return convertBool(name, key, source, target)
	case reflect.String:
		return convertString(name, key, source, target)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return convertInt(name, key, source, target)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return convertUint(name, key, source, target)
	case reflect.Float32, reflect.Float64:
		return convertFloat(name, key, source, target)
	case reflect.Map:
		return convertMap(name, key, source, target)
	case reflect.Slice:
		return convertSlice(name, key, source, target)
	case reflect.Struct:
		return convertStruct(name, key, source, target)
	case reflect.Ptr:
		return convertPtr(name, key, source, target)
	default:
		return makeError(name, key, source, target)
	}

}

func convertStruct(name, key string, source, target reflect.Value) (err error) {
	if !source.IsValid() {
		return
	}

	sourceValue := reflect.Indirect(source)
	sourceValueKind := sourceValue.Kind()
	if sourceValueKind != reflect.Map && sourceValueKind != reflect.Struct {
		return fmt.Errorf("structs can only be read back in from maps at this moment (field %q key %q) but got %s", name, key, sourceValueKind)
	}

	value := reflect.New(target.Type())
	keyName := fmt.Sprintf("%s.%s", name, key)
	if err := UnmarshalMap(MarshalMap(sourceValue.Interface()), value.Interface()); err != nil {
		return wrapError(name, keyName, sourceValue, value, err)
	}

	target.Set(reflect.Indirect(value))

	return
}

func convertPtr(name, key string, source, target reflect.Value) error {
	valueptr := reflect.New(target.Type().Elem())
	if valueptr.Type().Implements(unmarshallerType) {
		valueptr.Interface().(MapUnmarshaller).UnmarshalMap(source.Interface())
	} else {
		result := reflect.Indirect(valueptr)
		if err := convertValue(name, key, source, result); err != nil {
			return wrapError(name, key, source, result, err)
		}
	}
	target.Set(valueptr)
	return nil
}

func convertSlice(name, key string, source, target reflect.Value) error {
	sourceValue := reflect.Indirect(source)
	sourceValueKind := sourceValue.Kind()

	if sourceValueKind != reflect.Array && sourceValueKind != reflect.Slice {
		return fmt.Errorf("Expected %s (key %q) to be an array or a slice but was %s", name, key, sourceValueKind)
	}

	sz := source.Len()
	value := reflect.MakeSlice(reflect.SliceOf(target.Type().Elem()), sz, sz)
	for i := 0; i < sz; i++ {
		sourceData := source.Index(i)
		targetElement := value.Index(i)

		keyName := fmt.Sprintf("%s[%d]", name, i)
		if err := convertValue(name, keyName, sourceData, targetElement); err != nil {
			return wrapError(name, keyName, sourceData, targetElement, err)
		}
	}

	target.Set(value)
	return nil
}

func convertMap(name, key string, source, target reflect.Value) error {
	fieldType := target.Type()
	value := reflect.MakeMap(reflect.MapOf(fieldType.Key(), fieldType.Elem()))

	sourceValue := reflect.Indirect(source)
	if sourceValue.Kind() != reflect.Map {
		return makeError(name, key, sourceValue, target)
	}

	for _, k := range sourceValue.MapKeys() {
		keyName := fmt.Sprintf("%s[%s]", name, k)

		currentKey := reflect.Indirect(reflect.New(fieldType.Key()))
		if err := convertValue(name, keyName, k, currentKey); err != nil {
			return wrapError(name, keyName, currentKey, k, err)
		}

		currentValue := sourceValue.MapIndex(k)
		currentElem := reflect.Indirect(reflect.New(fieldType.Elem()))

		if err := convertValue(name, keyName, currentValue, currentElem); err != nil {
			return wrapError(name, keyName, currentElem, currentValue, err)
		}

		value.SetMapIndex(currentKey, currentElem)
	}

	target.Set(value)
	return nil
}

func convertInterface(name, key string, source, target reflect.Value) error {
	if source.Type().AssignableTo(target.Type()) {
		target.Set(source)
		return nil
	}
	return makeError(name, key, source, target)
}

func makeError(name, key string, source, target reflect.Value) error {
	return fmt.Errorf("Couldn't convert %s (key %q) '%s' to a '%s'", name, key, source.Type(), target.Type())
}
func wrapError(name, key string, source, target reflect.Value, err error) error {
	return fmt.Errorf("Couldn't convert %s (key %q) '%s' to a '%s', because %v", name, key, source.Type(), target.Type(), err)
}

func convertBool(name, key string, source, target reflect.Value) error {
	switch source.Kind() {
	case reflect.Bool:
		target.SetBool(source.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		target.SetBool(source.Int() != 0)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		target.SetBool(source.Uint() != 0)
	case reflect.Float32, reflect.Float64:
		target.SetBool(source.Float() != 0)
	case reflect.String:
		b, err := strconv.ParseBool(source.String())
		if err != nil {
			if source.String() == "" {
				target.SetBool(false)
				return nil
			}
			return wrapError(name, key, source, target, err)
		}
		target.SetBool(b)

	default:
		return makeError(name, key, source, target)
	}
	return nil
}

func convertString(name, key string, source, target reflect.Value) error {
	switch source.Kind() {
	case reflect.Bool:
		target.SetString(strconv.FormatBool(source.Bool()))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		target.SetString(strconv.FormatInt(source.Int(), 10))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		target.SetString(strconv.FormatUint(source.Uint(), 10))
	case reflect.Float32, reflect.Float64:
		target.SetString(strconv.FormatFloat(source.Float(), 'f', -1, 64))
	case reflect.String:
		target.SetString(source.String())
	default:
		return makeError(name, key, source, target)
	}
	return nil
}

func round(x float64, prec int) float64 {
	var rounder float64
	pow := math.Pow(10, float64(prec))
	intermed := x * pow
	_, frac := math.Modf(intermed)
	if frac >= 0.5 {
		rounder = math.Ceil(intermed)
	} else {
		rounder = math.Floor(intermed)
	}

	return rounder / pow
}

func convertInt(name, key string, source, target reflect.Value) error {
	switch source.Kind() {
	case reflect.Bool:
		if source.Bool() {
			target.SetInt(1)
		} else {
			target.SetInt(0)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		target.SetInt(source.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		target.SetInt(int64(source.Uint()))
	case reflect.Float32, reflect.Float64:
		target.SetInt(int64(round(source.Float(), 0)))
	case reflect.String:
		ev, err := strconv.ParseInt(source.String(), 10, 64)
		if err != nil {
			return err
		}
		if target.OverflowInt(ev) {
			return &reflect.ValueError{"reflect.Value.OverflowInt", target.Kind()}
		}
		target.SetInt(ev)
	default:
		return fmt.Errorf("Couldn't convert %v (type %T) to int", source.Interface(), source.Interface())
	}
	return nil
}

func convertUint(name, key string, source, target reflect.Value) error {
	switch source.Kind() {
	case reflect.Bool:
		if source.Bool() {
			target.SetUint(1)
		} else {
			target.SetUint(0)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		target.SetUint(uint64(source.Int()))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		target.SetUint(source.Uint())
	case reflect.Float32, reflect.Float64:
		target.SetUint(uint64(round(source.Float(), 0)))
	case reflect.String:
		ev, err := strconv.ParseUint(source.String(), 10, 64)
		if err != nil {
			return err
		}
		if target.OverflowUint(ev) {
			return &reflect.ValueError{"reflect.Value.OverflowUint", target.Kind()}
		}
		target.SetUint(ev)
	default:
		return fmt.Errorf("Couldn't convert %v (type %T) to int", source.Interface(), source.Interface())
	}
	return nil
}

func convertFloat(name, key string, source, target reflect.Value) error {
	switch source.Kind() {
	case reflect.Bool:
		if source.Bool() {
			target.SetFloat(1)
		} else {
			target.SetFloat(0)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		target.SetFloat(float64(source.Int()))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		target.SetFloat(float64(source.Uint()))
	case reflect.Float32, reflect.Float64:
		target.SetFloat(source.Float())
	case reflect.String:
		vv, err := strconv.ParseFloat(source.String(), 64)
		if err != nil {
			return err
		}
		if target.OverflowFloat(vv) {
			return &reflect.ValueError{"reflect.Value.OverflowFloat", target.Kind()}
		}
		target.SetFloat(vv)
	default:
		return fmt.Errorf("Couldn't convert %v (type %T) to int", source.Interface(), source.Interface())
	}
	return nil
}

func fieldTagNameMap(tpe reflect.Type) map[string]reflect.StructField {
	result := map[string]reflect.StructField{}
	for i := 0; i < tpe.NumField(); i++ {
		targetDes := tpe.Field(i)

		if targetDes.PkgPath != "" {
			continue
		}

		tag := parseTag(targetDes.Tag.Get(TagName), targetDes.Name)
		if !tag.ShouldSkip {
			result[tag.Name] = targetDes
		}
	}
	return result
}
