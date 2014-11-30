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
	Tag        *parsedTag
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

	fieldNameMap := fieldTagNameMap(tpe, val)
	for k, v := range data {
		if targetDes, ok := fieldNameMap[k]; ok {
			if err := convertValue(targetDes.Descriptor.Name, k, reflect.ValueOf(v), val.FieldByName(targetDes.Descriptor.Name), targetDes.Tag); err != nil {
				return err
			}
		}
	}
	return nil
}

func convertValue(name, key string, source, target reflect.Value, tag *parsedTag) error {
	if !target.IsValid() {
		target.Set(reflect.Zero(target.Type()))
		return nil
	}

	if !target.CanSet() {
		return fmt.Errorf("%s (key %q) %#v must be addressable", name, key, target)
	}

	//fmt.Printf("converting %s (key %s) from %s into %s\n", name, key, source.Type(), target.Type())
	if reflect.PtrTo(target.Type()).Implements(unmarshallerType) {
		//fmt.Printf("converting custom converter %s\n", key)
		value := reflect.New(target.Type())
		if err := value.Interface().(MapUnmarshaller).UnmarshalMap(MarshalMap(source.Interface())); err != nil {
			return err
		}
		target.Set(reflect.Indirect(value))
		return nil
	}

	if target.Kind() == reflect.Interface {
		//fmt.Printf("setting interface %s\n", key)
		target.Set(source)
		return nil
	}

	if source.Kind() == reflect.Interface {
		//fmt.Printf("converting interface %s\n", key)
		return convertInterface(name, key, source, target, tag)
	}

	switch target.Kind() {
	case reflect.Bool:
		//fmt.Printf("converting bool %s\n", key)
		return convertBool(name, key, source, target, tag)
	case reflect.String:
		//fmt.Printf("converting string %s\n", key)
		return convertString(name, key, source, target, tag)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		//fmt.Printf("converting int %s\n", key)
		return convertInt(name, key, source, target, tag)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		//fmt.Printf("converting uint %s\n", key)
		return convertUint(name, key, source, target, tag)
	case reflect.Float32, reflect.Float64:
		//fmt.Printf("converting float %s\n", key)
		return convertFloat(name, key, source, target, tag)
	case reflect.Map:
		//fmt.Printf("converting map %s\n", key)
		return convertMap(name, key, source, target, tag)
	case reflect.Slice:
		//fmt.Printf("converting slice %s\n", key)
		return convertSlice(name, key, source, target, tag)
	case reflect.Struct:
		//fmt.Printf("converting struct %s\n", key)
		return convertStruct(name, key, source, target, tag)
	case reflect.Ptr:
		//fmt.Printf("converting pointer %s\n", key)
		return convertPtr(name, key, source, target, tag)
	default:
		return makeError(name, key, source, target)
	}

}

func convertStruct(name, key string, source, target reflect.Value, tag *parsedTag) (err error) {
	if !source.IsValid() {
		return
	}

	sourceValue := reflect.Indirect(source)
	sourceValueKind := sourceValue.Type().Kind()

	if sourceValueKind != reflect.Interface && sourceValueKind != reflect.Map && sourceValueKind != reflect.Struct {
		return fmt.Errorf("structs can only be read back in from maps at this moment but got %s", sourceValueKind)
	}
	if tag.ByValue && source.Type().AssignableTo(target.Type()) {
		target.Set(source)
	}

	value := reflect.New(target.Type())
	keyName := fmt.Sprintf("%s.%s", name, key)
	if err := UnmarshalMap(MarshalMap(sourceValue.Interface()), value.Interface()); err != nil {
		return wrapError(name, keyName, sourceValue, value, err)
	}

	target.Set(reflect.Indirect(value))

	return
}

func convertPtr(name, key string, source, target reflect.Value, tag *parsedTag) error {
	valueptr := reflect.New(target.Type().Elem())
	if valueptr.Type().Implements(unmarshallerType) {
		valueptr.Interface().(MapUnmarshaller).UnmarshalMap(source.Interface())
	} else {
		result := reflect.Indirect(valueptr)
		if err := convertValue(name, key, source, result, tag); err != nil {
			return wrapError(name, key, source, result, err)
		}
	}
	target.Set(valueptr)
	return nil
}

func convertSlice(name, key string, source, target reflect.Value, tag *parsedTag) error {
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
		if err := convertValue(name, keyName, sourceData, targetElement, tag); err != nil {
			return wrapError(name, keyName, sourceData, targetElement, err)
		}
	}

	target.Set(value)
	return nil
}

func convertMap(name, key string, source, target reflect.Value, tag *parsedTag) error {
	fieldType := target.Type()
	value := reflect.MakeMap(reflect.MapOf(fieldType.Key(), fieldType.Elem()))

	sourceValue := reflect.Indirect(source)
	if sourceValue.Kind() != reflect.Map {
		return makeError(name, key, sourceValue, target)
	}

	for _, k := range sourceValue.MapKeys() {
		keyName := fmt.Sprintf("%s[%s]", name, k)

		newKey := reflect.Indirect(reflect.New(fieldType.Key()))
		if err := convertValue(name, keyName, k, newKey, tag); err != nil {
			return wrapError(name, keyName, newKey, k, err)
		}

		currentValue := sourceValue.MapIndex(k)
		currentElem := reflect.Indirect(reflect.New(fieldType.Elem()))

		if err := convertValue(name, keyName, currentValue, currentElem, tag); err != nil {
			return wrapError(name, keyName, currentElem, currentValue, err)
		}

		value.SetMapIndex(newKey, currentElem)
	}

	target.Set(value)
	return nil
}

func convertInterface(name, key string, source, target reflect.Value, tag *parsedTag) error {
	switch target.Kind() {
	case reflect.Bool:
		target.SetBool(source.Interface().(bool))
		return nil
	case reflect.String:
		target.SetString(source.Interface().(string))
		return nil
	case reflect.Int:
		target.SetInt(int64(source.Interface().(int)))
		return nil
	case reflect.Int8:
		target.SetInt(int64(source.Interface().(int8)))
		return nil
	case reflect.Int16:
		target.SetInt(int64(source.Interface().(int16)))
		return nil
	case reflect.Int32:
		target.SetInt(int64(source.Interface().(int32)))
		return nil
	case reflect.Int64:
		target.SetInt(source.Interface().(int64))
		return nil
	case reflect.Uint:
		target.SetUint(uint64(source.Interface().(uint)))
		return nil
	case reflect.Uint8:
		target.SetUint(uint64(source.Interface().(uint8)))
		return nil
	case reflect.Uint16:
		target.SetUint(uint64(source.Interface().(uint16)))
		return nil
	case reflect.Uint32:
		target.SetUint(uint64(source.Interface().(uint32)))
		return nil
	case reflect.Uint64:
		target.SetUint(source.Interface().(uint64))
		return nil
	case reflect.Float32:
		target.SetFloat(float64(source.Interface().(float32)))
		return nil
	case reflect.Float64:
		target.SetFloat(source.Interface().(float64))
		return nil
	case reflect.Struct:
		//fmt.Printf("convertInterface: this is a slice target\n")
		return convertValue(name, key, reflect.ValueOf(source.Interface()), target, tag)
	case reflect.Interface:
		//fmt.Printf("convertInterface: this is an interface target\n")
		target.Set(source)
		return nil
	case reflect.Map:
		//fmt.Printf("convertInterface: this is a map target\n")
		return convertValue(name, key, reflect.ValueOf(source.Interface()), target, tag)
	case reflect.Slice:
		//fmt.Printf("convertInterface: this is a slice target\n")
		return convertValue(name, key, reflect.ValueOf(source.Interface()), target, tag)
	default:
		if source.Type().AssignableTo(target.Type()) {
			target.Set(source)
			return nil
		}
		return makeError(name, key, source, target)
	}
}

func makeError(name, key string, source, target reflect.Value) error {
	return fmt.Errorf("Couldn't convert %s (key %q) '%s' to a '%s'", name, key, source.Type(), target.Type())
}
func wrapError(name, key string, source, target reflect.Value, err error) error {
	return fmt.Errorf("Couldn't convert %s (key %q) '%s' to a '%s', because %v", name, key, source.Type(), target.Type(), err)
}

func convertBool(name, key string, source, target reflect.Value, tag *parsedTag) error {
	switch source.Kind() {
	case reflect.Interface:
		target.SetBool(source.Interface().(bool))
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

func convertString(name, key string, source, target reflect.Value, tag *parsedTag) error {
	switch source.Kind() {
	case reflect.Interface:
		target.SetString(source.Interface().(string))
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

func convertInt(name, key string, source, target reflect.Value, tag *parsedTag) error {
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

func convertUint(name, key string, source, target reflect.Value, tag *parsedTag) error {
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

func convertFloat(name, key string, source, target reflect.Value, tag *parsedTag) error {
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

func fieldTagNameMap(tpe reflect.Type, val reflect.Value) map[string]fieldInfo {
	result := map[string]fieldInfo{}
	for i := 0; i < tpe.NumField(); i++ {
		targetDes := tpe.Field(i)

		if targetDes.PkgPath != "" {
			continue
		}

		tag := parseTag(targetDes.Tag.Get(TagName), targetDes.Name)
		if !tag.ShouldSkip && !targetDes.Anonymous {
			result[tag.Name] = fieldInfo{targetDes, tag}
		}
	}
	return result
}
