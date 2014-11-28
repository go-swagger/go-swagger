/*
A POWERFUL INTERFACE TO YOUR API

Swagger is a simple yet powerful representation of your RESTful API. With the largest ecosystem of API tooling on the planet,
thousands of developers are supporting Swagger in almost every modern programming language and deployment environment.
With a Swagger-enabled API, you get interactive documentation, client SDK generation and discoverability.
We created Swagger to help fulfill the promise of APIs.
Swagger helps companies like Apigee, Getty Images, Intuit, LivingSocial, McKesson, Microsoft, Morningstar, and PayPal build the best possible services with RESTful APIs.
Now in version 2.0, Swagger is more enabling than ever.
And it's 100% open source software.
*/
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

// ResponsesMap contains the responses by key
type ResponsesMap map[string]Response

// SecurityDefinitions a declaration of the security schemes available to be used in the specification.
// This does not enforce the security schemes on the operations and only serves to provide
// the relevant details for each scheme.
//
// For more information: http://goo.gl/8us55a#securityDefinitionsObject
type SecurityDefinitions map[string]*SecurityScheme

func addExtensions(res map[string]interface{}, extensions map[string]interface{}) {
	for k, v := range extensions {
		key := k
		if key != "" {
			if !strings.HasPrefix(key, "x-") {
				key = "x-" + key
			}
			zero := reflect.Zero(reflect.TypeOf(v)).Interface()
			if !reflect.DeepEqual(v, zero) {
				res[key] = reflection.MarshalMap(v)
			}
		}
	}
}

// StringOrArray represents a value that can either be a string
// or an array of strings. Mainly here for serialization purposes
type StringOrArray struct {
	Single string
	Multi  []string
}

func (s *StringOrArray) UnmarshalJSON(data []byte) error {
	var parsed interface{}
	err := json.Unmarshal(data, &parsed)
	if err != nil {
		return err
	}
	return s.unmarshalInterface(parsed, data)
}

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
		return fmt.Errorf("Could not unmarshal string or array from: `%s`", data)
	}
}

func (s StringOrArray) MarshalYAML() (interface{}, error) {
	if s.Single != "" {
		return s.Single, nil
	}
	return s.Multi, nil
}

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

func (s SchemaOrArray) MarshalYAML() (interface{}, error) {
	if s.Single != nil {
		return s.Single, nil
	}
	return s.Multi, nil
}

func (s SchemaOrArray) MarshalJSON() ([]byte, error) {
	if s.Single != nil {
		return json.Marshal(s.Single)
	}
	return json.Marshal(s.Multi)
}

func (s *SchemaOrArray) UnmarshalJSON(data []byte) error {
	var parsed interface{}
	err := json.Unmarshal(data, &parsed)
	if err != nil {
		return err
	}
	return s.unmarshalInterface(parsed, data)
}

func (s *SchemaOrArray) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var parsed interface{}
	err := unmarshal(&parsed)
	if err != nil {
		return err
	}
	return s.unmarshalInterface(parsed, []byte{})
}

func (s *SchemaOrArray) unmarshalInterface(parsed interface{}, data []byte) error {
	fmt.Printf("\nschema or interface: %T %#v\n", parsed, parsed)

	switch parsed.(type) {
	case map[string]interface{}:
		schema, err := SchemaFromMap(parsed.(map[string]interface{}))
		if err != nil {
			return err
		}
		s.Single = schema
		return nil
	case []interface{}:
		return nil
	case nil:
		return nil
	default:
		return fmt.Errorf("Could not unmarshal string or array from: `%s`", data)
	}
}
