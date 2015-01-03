package validate

import (
	"reflect"
	"testing"

	"github.com/casualjim/go-swagger/errors"
	"github.com/casualjim/go-swagger/spec"
	"github.com/stretchr/testify/assert"
)

// common validations: enum, allOf, anyOf, oneOf, not, definitions

func maxError(param *spec.Parameter) *errors.Validation {
	return errors.ExceedsMaximum(param.Name, param.In, *param.Maximum, param.ExclusiveMaximum)
}

func minError(param *spec.Parameter) *errors.Validation {
	return errors.ExceedsMinimum(param.Name, param.In, *param.Minimum, param.ExclusiveMinimum)
}

func multipleOfError(param *spec.Parameter) *errors.Validation {
	return errors.NotMultipleOf(param.Name, param.In, *param.MultipleOf)
}

func makeFloat(data interface{}) float64 {
	val := reflect.ValueOf(data)
	knd := val.Kind()
	switch {
	case knd >= reflect.Int && knd <= reflect.Int64:
		return float64(val.Int())
	case knd >= reflect.Uint && knd <= reflect.Uint64:
		return float64(val.Uint())
	default:
		return val.Float()
	}
}

func TestNumberParameterValidation(t *testing.T) {

	values := [][]interface{}{
		[]interface{}{23, 49, 56, 21, 14, 35, 28, 7, 42},
		[]interface{}{uint(23), uint(49), uint(56), uint(21), uint(14), uint(35), uint(28), uint(7), uint(42)},
		[]interface{}{float64(23), float64(49), float64(56), float64(21), float64(14), float64(35), float64(28), float64(7), float64(42)},
	}

	for _, v := range values {
		factorParam := spec.QueryParam("factor")
		factorParam.WithMaximum(makeFloat(v[1]), false)
		factorParam.WithMinimum(makeFloat(v[3]), false)
		factorParam.WithMultipleOf(makeFloat(v[7]))
		factorParam.WithEnum(v[3], v[6], v[8], v[1])
		validator := &paramValidator{factorParam, factorParam.Name}

		// MultipleOf
		err := validator.Validate(v[0])
		assert.Error(t, err)
		assert.EqualError(t, multipleOfError(factorParam), err.Error())

		// Maximum
		err = validator.Validate(v[1])
		assert.NoError(t, err)
		err = validator.Validate(v[2])
		assert.Error(t, err)
		assert.EqualError(t, maxError(factorParam), err.Error())
		// ExclusiveMaximum
		factorParam.ExclusiveMaximum = true
		err = validator.Validate(v[1])
		assert.Error(t, err)
		assert.EqualError(t, maxError(factorParam), err.Error())

		// Minimum
		err = validator.Validate(v[3])
		assert.NoError(t, err)
		err = validator.Validate(v[4])
		assert.Error(t, err)
		assert.EqualError(t, minError(factorParam), err.Error())
		// ExclusiveMinimum
		factorParam.ExclusiveMinimum = true
		err = validator.Validate(v[3])
		assert.Error(t, err)
		assert.EqualError(t, minError(factorParam), err.Error())

		// Enum
		err = validator.Validate(v[5])
		assert.Error(t, err)
		assert.EqualError(t, enumFail(factorParam, v[5]), err.Error())

		assert.NoError(t, validator.Validate(v[6]))
	}

	// AllOf
	// AnyOf
	// OneOf
	// Not
	// Definitions
}

func maxLengthError(param *spec.Parameter) *errors.Validation {
	return errors.TooLong(param.Name, param.In, *param.MaxLength)
}

func minLengthError(param *spec.Parameter) *errors.Validation {
	return errors.TooShort(param.Name, param.In, *param.MinLength)
}

func patternFail(param *spec.Parameter) *errors.Validation {
	return errors.FailedPattern(param.Name, param.In, param.Pattern)
}

func enumFail(param *spec.Parameter, data interface{}) *errors.Validation {
	return errors.EnumFail(param.Name, param.In, data, param.Enum)
}

func TestStringParameterValidation(t *testing.T) {
	nameParam := spec.QueryParam("name").AsRequired().WithMinLength(3).WithMaxLength(5).WithPattern(`^[a-z]+$`)
	nameParam.WithEnum("aaa", "bbb", "ccc")
	validator := &paramValidator{nameParam, nameParam.Name}

	// required
	err := validator.Validate("")
	assert.Error(t, err)
	assert.EqualError(t, requiredError(nameParam), err.Error())
	// MaxLength
	err = validator.Validate("abcdef")
	assert.Error(t, err)
	assert.EqualError(t, maxLengthError(nameParam), err.Error())
	// MinLength
	err = validator.Validate("a")
	assert.Error(t, err)
	assert.EqualError(t, minLengthError(nameParam), err.Error())
	// Pattern
	err = validator.Validate("a394")
	assert.Error(t, err)
	assert.EqualError(t, patternFail(nameParam), err.Error())

	// Enum
	err = validator.Validate("abcde")
	assert.Error(t, err)
	assert.EqualError(t, enumFail(nameParam, "abcde"), err.Error())

	// Valid passes
	err = validator.Validate("bbb")
	assert.NoError(t, err)

	// Not required in a parameter
	// AllOf
	// AnyOf
	// OneOf
	// Not
	// Definitions
}

func TestArrayParameterValidation(t *testing.T) {
	// Additional items
	// Items
	// MaxItems
	// MinItems
	// UniqueItems

	// Enum
	// AllOf
	// AnyOf
	// OneOf
	// Not
	// Definitions
}

func TestObjectParameterValidation(t *testing.T) {
	// MaxProperties
	// MinProperties
	// Required
	// AdditionalProperties
	// Properties
	// PatternProperties
	// Dependencies

	// Enum
	// AllOf
	// AnyOf
	// OneOf
	// Not
	// Definitions
}
