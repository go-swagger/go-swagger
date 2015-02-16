package validate

import (
	"reflect"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/casualjim/go-swagger/errors"
)

// Enum validates if the data is a member of the enum
func Enum(path, in string, data interface{}, enum interface{}) *errors.Validation {
	val := reflect.ValueOf(data)
	if val.Kind() != reflect.Slice {
		return nil
	}

	var values []interface{}
	for i := 0; i < val.Len(); i++ {
		ele := val.Index(i)
		enumValue := ele.Interface()
		if data != nil && reflect.DeepEqual(data, enumValue) {
			return nil
		}
		values = append(values, enumValue)
	}
	return errors.EnumFail(path, in, data, values)
}

// MinItems validates that there are at least n items in a slice
func MinItems(path, in string, size, min int64) *errors.Validation {
	if size < min {
		return errors.TooFewItems(path, in, min)
	}
	return nil
}

// MaxItems validates that there are at most n items in a slice
func MaxItems(path, in string, size, max int64) *errors.Validation {
	if size > max {
		return errors.TooManyItems(path, in, max)
	}
	return nil
}

// UniqueItems validates that the provided slice has unique elements
func UniqueItems(path, in string, data interface{}) *errors.Validation {
	val := reflect.ValueOf(data)
	if val.Kind() != reflect.Slice {
		return nil
	}

	dict := make(map[interface{}]struct{})
	for i := 0; i < val.Len(); i++ {
		ele := val.Index(i)
		if _, ok := dict[ele.Interface()]; ok {
			return errors.DuplicateItems(path, in)
		}
		dict[ele.Interface()] = struct{}{}
	}
	return nil
}

// MinLength validates a string for minimum length
func MinLength(path, in, data string, minLength int64) *errors.Validation {
	strLen := int64(utf8.RuneCount([]byte(data)))
	if strLen < minLength {
		return errors.TooShort(path, in, minLength)
	}
	return nil
}

// MaxLength validates a string for maximum length
func MaxLength(path, in, data string, maxLength int64) *errors.Validation {
	strLen := int64(utf8.RuneCount([]byte(data)))
	if strLen > maxLength {
		return errors.TooLong(path, in, maxLength)
	}
	return nil
}

// Required validates an interface for requiredness
func Required(path, in string, data interface{}) *errors.Validation {
	val := reflect.ValueOf(data)
	if reflect.DeepEqual(reflect.Zero(val.Type()), val) {
		return errors.Required(path, in)
	}
	return nil
}

// RequiredString validates a string for requiredness
func RequiredString(path, in, data string) *errors.Validation {
	if data == "" {
		return errors.Required(path, in)
	}
	return nil
}

// RequiredNumber validates a number for requiredness
func RequiredNumber(path, in string, data float64) *errors.Validation {
	if data == 0 {
		return errors.Required(path, in)
	}
	return nil
}

// Pattern validates a string against a regular expression
func Pattern(path, in, data, pattern string) *errors.Validation {
	re := regexp.MustCompile(pattern)
	if !re.MatchString(data) {
		return errors.FailedPattern(path, in, pattern)
	}
	return nil
}

// Maximum validates if a number is smaller than a given maximum
func Maximum(path, in string, data, max float64, exclusive bool) *errors.Validation {
	if (!exclusive && data > max) || (exclusive && data >= max) {
		return errors.ExceedsMaximum(path, in, max, exclusive)
	}
	return nil
}

// Minimum validates if a number is smaller than a given minimum
func Minimum(path, in string, data, min float64, exclusive bool) *errors.Validation {
	if (!exclusive && data < min) || (exclusive && data <= min) {
		return errors.ExceedsMinimum(path, in, min, exclusive)
	}
	return nil
}

// MultipleOf validates if the provided number is a multiple of the factor
func MultipleOf(path, in string, data, factor float64) *errors.Validation {
	if !isFloat64AnInteger(data / factor) {
		return errors.NotMultipleOf(path, in, factor)
	}
	return nil
}

// FormatOf validates if a string matches a format in the format registry
func FormatOf(path, in, format, data string) *errors.Validation {
	validate, ok := formatCheckers[strings.Replace(format, "-", "", -1)]
	if !ok {
		return errors.InvalidTypeName(format)
	}
	if validate(data) {
		return nil
	}
	return errors.InvalidType(path, in, format, data)
}
