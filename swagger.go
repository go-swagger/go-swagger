package swagger

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/casualjim/go-swagger/reflection"
)

// Definitions contains the models explicitly defined in this spec
// An object to hold data types that can be consumed and produced by operations.
// These data types can be primitives, arrays or models.
//
// For more information: http://goo.gl/8us55a#definitionsObject
type Definitions map[string]Schema

// SecurityDefinitions a declaration of the security schemes available to be used in the specification.
// This does not enforce the security schemes on the operations and only serves to provide
// the relevant details for each scheme.
//
// For more information: http://goo.gl/8us55a#securityDefinitionsObject
type SecurityDefinitions map[string]*SecurityScheme

// Extensions vendor specific extensions
type Extensions map[string]interface{}

// Add adds a value to these extensions
func (e Extensions) Add(key string, value interface{}) {
	realKey := strings.ToLower(key)
	e[realKey] = value
}

// GetString gets a string value from the extensions
func (e Extensions) GetString(key string) (string, bool) {
	if v, ok := e[strings.ToLower(key)]; ok {
		str, ok := v.(string)
		return str, ok
	}
	return "", false
}

func addExtensions(res map[string]interface{}, extensions map[string]interface{}) {
	for k, v := range extensions {
		key := k
		if key != "" {
			if !strings.HasPrefix(key, "x-") {
				key = "x-" + key
			}
		}
		valueType := reflect.TypeOf(v)
		if !reflection.IsZero(reflect.ValueOf(v)) {
			switch valueType.Kind() {
			case reflect.Struct, reflect.Map:
				res[key] = reflection.MarshalMap(v)
			default:
				if reflection.IsZero(reflect.ValueOf(v)) {
					res[key] = make(map[string]interface{})
				} else {
					res[key] = v
				}
			}
		}
	}
}

func readExtensions(res map[string]interface{}) map[string]interface{} {
	var extensions map[string]interface{}
	for k, v := range res {
		if strings.HasPrefix(k, "x-") {
			if extensions == nil {
				extensions = make(map[string]interface{})
			}
			if reflection.IsZero(reflect.ValueOf(v)) {
				extensions[k] = make(map[string]interface{})
			} else {
				extensions[k] = v
			}
			delete(res, k)
		}
	}
	return extensions
}

// StringOrArray represents a value that can either be a string
// or an array of strings. Mainly here for serialization purposes
type StringOrArray struct {
	Single string
	Multi  []string
}

// UnmarshalMap hydrates this string or string array with data from the map
func (s *StringOrArray) UnmarshalMap(data interface{}) error {
	return s.unmarshalInterface(data, []byte{})
}

// UnmarshalJSON unmarshals this string or array object from a JSON array or JSON string
func (s *StringOrArray) UnmarshalJSON(data []byte) error {
	var parsed interface{}
	err := json.Unmarshal(data, &parsed)
	if err != nil {
		return err
	}
	return s.unmarshalInterface(parsed, data)
}

// UnmarshalYAML unmarshals this string or array object from a YAML array or YAML string
func (s *StringOrArray) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var parsed interface{}
	err := unmarshal(&parsed)
	if err != nil {
		return err
	}
	return s.unmarshalInterface(parsed, []byte{})
}

func (s *StringOrArray) unmarshalInterface(parsed interface{}, data []byte) error {
	switch parsed.(type) {
	case string:
		s.Single = parsed.(string)
		return nil
	case []interface{}:
		arr := parsed.([]interface{})
		var multi []string
		for _, v := range arr {
			if v == nil {
				continue
			}
			str, ok := v.(string)
			if !ok {
				return fmt.Errorf("only string array are allowed for string or array")
			}
			multi = append(multi, str)
		}
		s.Multi = multi
		return nil
	case nil:
		return nil
	default:
		return fmt.Errorf("could not unmarshal string or array from: `%s`", data)
	}
}

// MarshalYAML converts this string or array to a YAML array or YAML string
func (s StringOrArray) MarshalYAML() (interface{}, error) {
	if s.Single != "" {
		return s.Single, nil
	}
	return s.Multi, nil
}

// MarshalJSON converts this string or array to a JSON array or JSON string
func (s StringOrArray) MarshalJSON() ([]byte, error) {
	if s.Single != "" {
		return json.Marshal(s.Single)
	}
	return json.Marshal(s.Multi)
}

// SchemaOrArray represents a value that can either be a Schema
// or an array of Schema. Mainly here for serialization purposes
type SchemaOrArray struct {
	Single *Schema
	Multi  []Schema
}

// MarshalYAML converts this schema object or array into YAML structure
func (s SchemaOrArray) MarshalYAML() (interface{}, error) {
	if s.Single != nil {
		return s.Single, nil
	}
	return s.Multi, nil
}

// MarshalJSON converts this schema object or array into JSON structure
func (s SchemaOrArray) MarshalJSON() ([]byte, error) {
	if s.Single != nil {
		return json.Marshal(s.Single)
	}
	return json.Marshal(s.Multi)
}

// UnmarshalMap hydrates this schema or array with data from a map
func (s *SchemaOrArray) UnmarshalMap(data interface{}) error {
	return s.unmarshalInterface(data, []byte{})
}

// UnmarshalJSON converts this schema object or array from a JSON structure
func (s *SchemaOrArray) UnmarshalJSON(data []byte) error {
	var parsed interface{}
	err := json.Unmarshal(data, &parsed)
	if err != nil {
		return err
	}
	return s.unmarshalInterface(parsed, data)
}

// UnmarshalYAML converts this schema object or array from a YAML structure
func (s *SchemaOrArray) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var parsed interface{}
	err := unmarshal(&parsed)
	if err != nil {
		return err
	}
	return s.unmarshalInterface(parsed, []byte{})
}

func (s *SchemaOrArray) unmarshalInterface(parsed interface{}, data []byte) error {
	switch parsed.(type) {
	case map[string]interface{}:
		val := &Schema{}
		err := reflection.UnmarshalMap(parsed.(map[string]interface{}), val)
		if err != nil {
			return err
		}
		s.Single = val
		return nil
	case map[interface{}]interface{}:
		val := &Schema{}
		stringMap := make(map[string]interface{})
		for k, v := range parsed.(map[interface{}]interface{}) {
			stringMap[k.(string)] = v
		}
		err := reflection.UnmarshalMap(stringMap, val)
		if err != nil {
			return err
		}
		s.Single = val
		return nil
	case []interface{}:
		val := []Schema{}
		for _, v := range parsed.([]interface{}) {
			if dict, ok := v.(map[string]interface{}); ok {
				elem := Schema{}
				err := reflection.UnmarshalMap(dict, &elem)
				if err != nil {
					return err
				}
				val = append(val, elem)
			} else if dict, ok := v.(map[interface{}]interface{}); ok {
				elem := Schema{}
				stringMap := make(map[string]interface{})
				for k, vv := range dict {
					stringMap[k.(string)] = vv
				}

				err := reflection.UnmarshalMap(stringMap, &elem)
				if err != nil {
					return err
				}
				val = append(val, elem)
			}
		}
		s.Multi = val
		return nil
	case nil:
		return nil
	default:
		return fmt.Errorf("could not unmarshal string or array from: `%s`", data)
	}
}
