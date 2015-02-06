package errors

import "fmt"

const (
	invalidType           = `%s is an invalid type name`
	typeFail              = `%s in %s must be of type %s`
	typeFailWithData      = `%s in %s must be of type %s: %q`
	typeFailWithError     = `%s in %s must be of type %s, because: %s`
	requiredFail          = `%s in %s is required`
	tooLongMessage        = `%s in %s should be at most %d chars long`
	tooShortMessage       = `%s in %s should be at least %d chars long`
	patternFail           = `%s in %s should match '%s'`
	enumFail              = `%s in %s should be one of %v`
	mulitpleOfFail        = `%s in %s should be a multiple of %v`
	maxIncFail            = `%s in %s should be less than or equal to %v`
	maxExcFail            = `%s in %s should be less than %v`
	minIncFail            = `%s in %s should be greater than or equal to %v`
	minExcFail            = `%s in %s should be greater than %v`
	uniqueFail            = `%s in %s shouldn't contain duplicates`
	maxItemsFail          = `%s in %s should have at most %d items`
	minItemsFail          = `%s in %s should have at least %d items`
	typeFailNoIn          = `%s must be of type %s`
	typeFailWithDataNoIn  = `%s must be of type %s: %q`
	typeFailWithErrorNoIn = `%s must be of type %s, because: %s`
	requiredFailNoIn      = `%s is required`
	tooLongMessageNoIn    = `%s should be at most %d chars long`
	tooShortMessageNoIn   = `%s should be at least %d chars long`
	patternFailNoIn       = `%s should match '%s'`
	enumFailNoIn          = `%s should be one of %v`
	mulitpleOfFailNoIn    = `%s should be a multiple of %v`
	maxIncFailNoIn        = `%s should be less than or equal to %v`
	maxExcFailNoIn        = `%s should be less than %v`
	minIncFailNoIn        = `%s should be greater than or equal to %v`
	minExcFailNoIn        = `%s should be greater than %v`
	uniqueFailNoIn        = `%s shouldn't contain duplicates`
	maxItemsFailNoIn      = `%s should have at most %d items`
	minItemsFailNoIn      = `%s should have at least %d items`
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
		message: fmt.Sprintf("the collection format %q is not supported for the %s param %q", format, in, name),
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

	if in != "" {
		switch value.(type) {
		case string:
			message = fmt.Sprintf(typeFailWithData, name, in, typeName, value)
		case error:
			message = fmt.Sprintf(typeFailWithError, name, in, typeName, value)
		default:
			message = fmt.Sprintf(typeFail, name, in, typeName)
		}
	} else {
		switch value.(type) {
		case string:
			message = fmt.Sprintf(typeFailWithDataNoIn, name, typeName, value)
		case error:
			message = fmt.Sprintf(typeFailWithErrorNoIn, name, typeName, value)
		default:
			message = fmt.Sprintf(typeFailNoIn, name, typeName)
		}
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
	msg := fmt.Sprintf(uniqueFail, name, in)
	if in == "" {
		msg = fmt.Sprintf(uniqueFailNoIn, name)
	}
	return &Validation{
		code:    422,
		Name:    name,
		In:      in,
		message: msg,
	}
}

// TooManyItems error for when an array contains too many items
func TooManyItems(name, in string, max int64) *Validation {
	msg := fmt.Sprintf(maxItemsFail, name, in, max)
	if in == "" {
		msg = fmt.Sprintf(maxItemsFailNoIn, name, max)
	}

	return &Validation{
		code:    422,
		Name:    name,
		In:      in,
		message: msg,
	}
}

// TooFewItems error for when an array contains too few items
func TooFewItems(name, in string, min int64) *Validation {
	msg := fmt.Sprintf(minItemsFail, name, in, min)
	if in == "" {
		msg = fmt.Sprintf(minItemsFailNoIn, name, min)
	}
	return &Validation{
		code:    422,
		Name:    name,
		In:      in,
		message: msg,
	}
}

// ExceedsMaximum error for when maxinum validation fails
func ExceedsMaximum(name, in string, max float64, exclusive bool) *Validation {
	var message string
	if in == "" {
		m := maxIncFailNoIn
		if exclusive {
			m = maxExcFailNoIn
		}
		message = fmt.Sprintf(m, name, max)
	} else {
		m := maxIncFail
		if exclusive {
			m = maxExcFail
		}
		message = fmt.Sprintf(m, name, in, max)
	}
	return &Validation{
		code:    422,
		Name:    name,
		In:      in,
		Value:   max,
		message: message,
	}
}

// ExceedsMinimum error for when maxinum validation fails
func ExceedsMinimum(name, in string, min float64, exclusive bool) *Validation {
	var message string
	if in == "" {
		m := minIncFailNoIn
		if exclusive {
			m = minExcFailNoIn
		}
		message = fmt.Sprintf(m, name, min)
	} else {
		m := minIncFail
		if exclusive {
			m = minExcFail
		}
		message = fmt.Sprintf(m, name, in, min)
	}
	return &Validation{
		code:    422,
		Name:    name,
		In:      in,
		Value:   min,
		message: message,
	}
}

// NotMultipleOf error for when multiple of validation fails
func NotMultipleOf(name, in string, multiple float64) *Validation {
	var msg string
	if in == "" {
		msg = fmt.Sprintf(mulitpleOfFailNoIn, name, multiple)
	} else {
		msg = fmt.Sprintf(mulitpleOfFail, name, in, multiple)
	}
	return &Validation{
		code:    422,
		Name:    name,
		In:      in,
		Value:   multiple,
		message: msg,
	}
}

// EnumFail error for when an enum validation fails
func EnumFail(name, in string, value interface{}, values []interface{}) *Validation {
	var msg string
	if in == "" {
		msg = fmt.Sprintf(enumFailNoIn, name, values)
	} else {
		msg = fmt.Sprintf(enumFail, name, in, values)
	}

	return &Validation{
		code:    422,
		Name:    name,
		In:      in,
		Value:   value,
		Values:  values,
		message: msg,
	}
}

// Required error for when a value is missing
func Required(name, in string) *Validation {
	var msg string
	if in == "" {
		msg = fmt.Sprintf(requiredFailNoIn, name)
	} else {
		msg = fmt.Sprintf(requiredFail, name, in)
	}
	return &Validation{
		code:    422,
		Name:    name,
		In:      in,
		message: msg,
	}
}

// TooLong error for when a string is too long
func TooLong(name, in string, max int64) *Validation {
	var msg string
	if in == "" {
		msg = fmt.Sprintf(tooLongMessageNoIn, name, max)
	} else {
		msg = fmt.Sprintf(tooLongMessage, name, in, max)
	}
	return &Validation{
		code:    422,
		Name:    name,
		In:      in,
		message: msg,
	}
}

// TooShort error for when a string is too short
func TooShort(name, in string, min int64) *Validation {
	var msg string
	if in == "" {
		msg = fmt.Sprintf(tooShortMessageNoIn, name, min)
	} else {
		msg = fmt.Sprintf(tooShortMessage, name, in, min)
	}

	return &Validation{
		code:    422,
		Name:    name,
		In:      in,
		message: msg,
	}
}

// FailedPattern error for when a string fails a regex pattern match
// the pattern that is returned is the ECMA syntax version of the pattern not the golang version.
func FailedPattern(name, in, pattern string) *Validation {
	var msg string
	if in == "" {
		msg = fmt.Sprintf(patternFailNoIn, name, pattern)
	} else {
		msg = fmt.Sprintf(patternFail, name, in, pattern)
	}

	return &Validation{
		code:    422,
		Name:    name,
		In:      in,
		message: msg,
	}
}
