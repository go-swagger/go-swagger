package validate

import (
	"math"
	"reflect"
	"regexp"

	"github.com/casualjim/go-swagger/errors"
	"github.com/casualjim/go-swagger/spec"
)

// a param has very limited subset of validations to apply
type paramValidator struct {
	param *spec.Parameter
	name  string
}

func (p *paramValidator) Validate(data interface{}) *errors.Validation {
	val := reflect.Indirect(reflect.ValueOf(data))
	tpe := val.Type()

	switch tpe.Kind() {
	case reflect.String:
		if err := p.validateString(data.(string)); err != nil {
			return err
		}
		if err := p.validateCommon(data); err != nil {
			return err
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if err := p.validateNumber(float64(val.Int())); err != nil {
			return err
		}
		if err := p.validateCommon(data); err != nil {
			return err
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if err := p.validateNumber(float64(val.Uint())); err != nil {
			return err
		}
		if err := p.validateCommon(data); err != nil {
			return err
		}
	case reflect.Float32, reflect.Float64:
		if err := p.validateNumber(val.Float()); err != nil {
			return err
		}
		if err := p.validateCommon(data); err != nil {
			return err
		}
	}
	return nil
}

func (p *paramValidator) validateCommon(data interface{}) *errors.Validation {
	if len(p.param.Enum) > 0 {
		for i := 0; i < len(p.param.Enum); i++ {
			if reflect.DeepEqual(p.param.Enum[i], data) {
				return nil
			}
		}
		return errors.EnumFail(p.name, p.param.In, data, p.param.Enum)
	}
	return nil
}

func (p *paramValidator) validateNumber(data float64) *errors.Validation {
	return (&numberValidator{
		Name:             p.param.Name,
		In:               p.param.In,
		Default:          p.param.Default,
		MultipleOf:       p.param.MultipleOf,
		Maximum:          p.param.Maximum,
		ExclusiveMaximum: p.param.ExclusiveMaximum,
		Minimum:          p.param.Minimum,
		ExclusiveMinimum: p.param.ExclusiveMinimum,
	}).Validate(data)
}
func (p *paramValidator) validateString(data string) *errors.Validation {
	return (&stringValidator{
		Name:      p.param.Name,
		In:        p.param.In,
		Default:   p.param.Default,
		Required:  p.param.Required,
		MaxLength: p.param.MaxLength,
		MinLength: p.param.MinLength,
		Pattern:   p.param.Pattern,
	}).Validate(data)
}

type numberValidator struct {
	Name             string
	In               string
	Default          interface{}
	MultipleOf       *float64
	Maximum          *float64
	ExclusiveMaximum bool
	Minimum          *float64
	ExclusiveMinimum bool
}

func (n *numberValidator) Validate(data float64) *errors.Validation {
	if n.MultipleOf != nil && math.Mod(data, *n.MultipleOf) != 0 {
		return errors.NotMultipleOf(n.Name, n.In, *n.MultipleOf)
	}
	if n.Maximum != nil {
		max := *n.Maximum
		if n.ExclusiveMaximum {
			max--
		}
		if max < data {
			return errors.ExceedsMaximum(n.Name, n.In, *n.Maximum, n.ExclusiveMaximum)
		}
	}
	if n.Minimum != nil {
		min := *n.Minimum
		if n.ExclusiveMinimum {
			min++
		}
		if min > data {
			return errors.ExceedsMinimum(n.Name, n.In, *n.Minimum, n.ExclusiveMinimum)
		}
	}
	return nil
}

type stringValidator struct {
	Default   interface{}
	Required  bool
	MaxLength *int64
	MinLength *int64
	Pattern   string
	Name      string
	In        string
}

func (s *stringValidator) Validate(data string) *errors.Validation {
	if s.Required && s.Default == nil && data == "" {
		return errors.Required(s.Name, s.In)
	}

	if s.MaxLength != nil && int64(len(data)) > *s.MaxLength {
		return errors.TooLong(s.Name, s.In, *s.MaxLength)
	}

	if s.MinLength != nil && int64(len(data)) < *s.MinLength {
		return errors.TooShort(s.Name, s.In, *s.MinLength)
	}

	if s.Pattern != "" {
		// TODO: translate back and forth from javascript syntax, peferrably cached?
		//       perhaps allow the option of using a different regex engine like pcre?
		re := regexp.MustCompile(s.Pattern)
		if !re.MatchString(data) {
			return errors.FailedPattern(s.Name, s.In, s.Pattern)
		}
	}
	return nil
}

// type objectPath struct {
// 	head string
// 	tail *objectPath
// }

// func consPathNode(head string, tail *objectPath) *objectPath {
// 	return &objectPath{head, tail}
// }

// // String displays the context in reverse.
// // This plays well with the data structure's persistent nature with
// // Cons and a json document's tree structure.
// func (c *objectPath) String() string {
// 	byteArr := make([]byte, 0, c.stringLen())
// 	buf := bytes.NewBuffer(byteArr)
// 	c.writeStringToBuffer(buf)

// 	return buf.String()
// }

// func (c *objectPath) stringLen() int {
// 	length := 0
// 	if c.tail != nil {
// 		length = c.tail.stringLen() + 1 // add 1 for "."
// 	}

// 	length += len(c.head)
// 	return length
// }

// func (c *objectPath) writeStringToBuffer(buf *bytes.Buffer) {
// 	if c.tail != nil {
// 		c.tail.writeStringToBuffer(buf)
// 		buf.WriteString(".")
// 	}

// 	buf.WriteString(c.head)
// }

// // validates things after the type has been verified, defaults are set and the value was converted
// // knows how to validate schemas as the other things are just subsets of a schema
// type schemaValidator struct {
// 	path     objectPath
// 	schema   spec.Schema // copy of schema, or adapter for parameter, items and header
// 	required bool
// }

// func (s *schemaValidator) Validate(data interface{}) errors.Error {
// 	val := reflect.ValueOf(data)
// 	tpe := val.Type()
// 	if tpe.Kind() == reflect.Ptr {
// 		val = val.Elem()
// 		tpe = val.Type()
// 	}

// 	switch val.Kind() {
// 	case reflect.String:
// 		// validate required
// 		// validate max length
// 		// validate min length
// 		// validate pattern
// 	default:
// 		return errors.NotImplemented("validating other kinds than string is not yet implemented")
// 	}
// 	return nil
// }
