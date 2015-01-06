package errors

import "fmt"

const (
	invalidType       = `%s is an invalid type name`
	typeFail          = `%s in %s must be of type %s`
	typeFailWithData  = `%s in %s must be of type %s: %q`
	typeFailWithError = `%s in %s must be of type %s, because: %s`
	requiredFail      = `%s in %s is required`
	tooLongMessage    = `%s in %s should be at most %d chars long`
	tooShortMessage   = `%s in %s should be at least %d chars long`
	patternFail       = `%s in %s should match '%s'`
	enumFail          = `%s in %s should be one of %v`
	mulitpleOfFail    = `%s in %s should be a multiple of %v`
	maxIncFail        = `%s in %s should be less than or equal to %v`
	maxExcFail        = `%s in %s should be less than %v`
	minIncFail        = `%s in %s should be greater than or equal to %v`
	minExcFail        = `%s in %s should be greater than %v`
	uniqueFail        = `%s in %s should shouldn't contain duplicates`
	maxItemsFail      = `%s in %s should at most have %d items`
	minItemsFail      = `%s in %s should at most have %d items`
)

// CompositeValidationError an error to wrap a bunch of other errors
func CompositeValidationError(errors ...Error) *Validation {
	return &Validation{
		code:    422,
		Value:   append([]Error{}, errors...),
		message: "validation failure list",
	}
}

// InvalidCollectionFormat another flavor of invalid type error
func InvalidCollectionFormat(name, in, format string) *Validation {
	return &Validation{
		code:    422,
		Name:    name,
		In:      in,
		Value:   format,
		message: fmt.Sprintf("the collection format %q is not supported for a %s param", format, in),
	}
}

// InvalidTypeName an error for when the type is invalid
func InvalidTypeName(typeName string) *Validation {
	return &Validation{
		code:    422,
		message: fmt.Sprintf(invalidType, typeName),
	}
}

// InvalidType creates an error for when the type is invalid
func InvalidType(name, in, typeName string, value interface{}) *Validation {
	var message string
	switch value.(type) {
	case string:
		message = fmt.Sprintf(typeFailWithData, name, in, typeName, value)
	case error:
		message = fmt.Sprintf(typeFailWithError, name, in, typeName, value)
	default:
		message = fmt.Sprintf(typeFail, name, in, typeName)
	}
	return &Validation{
		code:    422,
		Name:    name,
		In:      in,
		Value:   value,
		message: message,
	}
}

// DuplicateItems error for when an array contains duplicates
func DuplicateItems(name, in string) *Validation {
	return &Validation{
		code:    422,
		Name:    name,
		In:      in,
		message: fmt.Sprintf(uniqueFail, name, in),
	}
}

// TooManyItems error for when an array contains too many items
func TooManyItems(name, in string, max int64) *Validation {
	return &Validation{
		code:    422,
		Name:    name,
		In:      in,
		message: fmt.Sprintf(maxItemsFail, name, in, max),
	}
}

// TooFewItems error for when an array contains too few items
func TooFewItems(name, in string, min int64) *Validation {
	return &Validation{
		code:    422,
		Name:    name,
		In:      in,
		message: fmt.Sprintf(minItemsFail, name, in, min),
	}
}

// ExceedsMaximum error for when maxinum validation fails
func ExceedsMaximum(name, in string, max float64, exclusive bool) *Validation {
	message := maxIncFail
	if exclusive {
		message = maxExcFail
	}
	return &Validation{
		code:    422,
		Name:    name,
		In:      in,
		Value:   max,
		message: fmt.Sprintf(message, name, in, max),
	}
}

// ExceedsMinimum error for when maxinum validation fails
func ExceedsMinimum(name, in string, min float64, exclusive bool) *Validation {
	message := minIncFail
	if exclusive {
		message = minExcFail
	}
	return &Validation{
		code:    422,
		Name:    name,
		In:      in,
		Value:   min,
		message: fmt.Sprintf(message, name, in, min),
	}
}

// NotMultipleOf error for when multiple of validation fails
func NotMultipleOf(name, in string, multiple float64) *Validation {
	return &Validation{
		code:    422,
		Name:    name,
		In:      in,
		Value:   multiple,
		message: fmt.Sprintf(mulitpleOfFail, name, in, multiple),
	}
}

// EnumFail error for when an enum validation fails
func EnumFail(name, in string, value interface{}, values []interface{}) *Validation {
	return &Validation{
		code:    422,
		Name:    name,
		In:      in,
		Value:   value,
		Values:  values,
		message: fmt.Sprintf(enumFail, name, in, values),
	}
}

// Required error for when a value is missing
func Required(name, in string) *Validation {
	return &Validation{
		code:    422,
		Name:    name,
		In:      in,
		message: fmt.Sprintf(requiredFail, name, in),
	}
}

// TooLong error for when a string is too long
func TooLong(name, in string, max int64) *Validation {
	return &Validation{
		code:    422,
		Name:    name,
		In:      in,
		message: fmt.Sprintf(tooLongMessage, name, in, max),
	}
}

// TooShort error for when a string is too short
func TooShort(name, in string, min int64) *Validation {
	return &Validation{
		code:    422,
		Name:    name,
		In:      in,
		message: fmt.Sprintf(tooShortMessage, name, in, min),
	}
}

// FailedPattern error for when a string fails a regex pattern match
// the pattern that is returned is the ECMA syntax version of the pattern not the golang version.
func FailedPattern(name, in, pattern string) *Validation {
	return &Validation{
		code:    422,
		Name:    name,
		In:      in,
		message: fmt.Sprintf(patternFail, name, in, pattern),
	}
}
