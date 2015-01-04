package validate

import (
	"fmt"
	"testing"

	"github.com/casualjim/go-swagger/errors"
	"github.com/casualjim/go-swagger/spec"
	"github.com/stretchr/testify/assert"
)

func maxErrorItems(path, in string, items *spec.Items) *errors.Validation {
	return errors.ExceedsMaximum(path, in, *items.Maximum, items.ExclusiveMaximum)
}

func minErrorItems(path, in string, items *spec.Items) *errors.Validation {
	return errors.ExceedsMinimum(path, in, *items.Minimum, items.ExclusiveMinimum)
}

func multipleOfErrorItems(path, in string, items *spec.Items) *errors.Validation {
	return errors.NotMultipleOf(path, in, *items.MultipleOf)
}

func requiredErrorItems(path, in string) *errors.Validation {
	return errors.Required(path, in)
}

func maxLengthErrorItems(path, in string, items *spec.Items) *errors.Validation {
	return errors.TooLong(path, in, *items.MaxLength)
}

func minLengthErrorItems(path, in string, items *spec.Items) *errors.Validation {
	return errors.TooShort(path, in, *items.MinLength)
}

func patternFailItems(path, in string, items *spec.Items) *errors.Validation {
	return errors.FailedPattern(path, in, items.Pattern)
}

func enumFailItems(path, in string, items *spec.Items, data interface{}) *errors.Validation {
	return errors.EnumFail(path, in, data, items.Enum)
}

func minItemsErrorItems(path, in string, items *spec.Items) *errors.Validation {
	return errors.TooFewItems(path, in, *items.MinItems)
}

func maxItemsErrorItems(path, in string, items *spec.Items) *errors.Validation {
	return errors.TooManyItems(path, in, *items.MaxItems)
}

func duplicatesErrorItems(path, in string) *errors.Validation {
	return errors.DuplicateItems(path, in)
}

func TestNumberItemsValidation(t *testing.T) {

	values := [][]interface{}{
		[]interface{}{23, 49, 56, 21, 14, 35, 28, 7, 42},
		[]interface{}{uint(23), uint(49), uint(56), uint(21), uint(14), uint(35), uint(28), uint(7), uint(42)},
		[]interface{}{float64(23), float64(49), float64(56), float64(21), float64(14), float64(35), float64(28), float64(7), float64(42)},
	}

	for i, v := range values {
		items := spec.NewItems()
		items.WithMaximum(makeFloat(v[1]), false)
		items.WithMinimum(makeFloat(v[3]), false)
		items.WithMultipleOf(makeFloat(v[7]))
		items.WithEnum(v[3], v[6], v[8], v[1])
		parent := spec.QueryParam("factors").CollectionOf(items, "")
		path := fmt.Sprintf("factors.%d", i)
		validator := &itemsValidator{items, parent, parent.Name, parent.In}

		// MultipleOf
		err := validator.Validate(i, v[0])
		assert.Error(t, err)
		assert.EqualError(t, multipleOfErrorItems(path, validator.in, items), err.Error())

		// Maximum
		err = validator.Validate(i, v[1])
		assert.NoError(t, err)
		err = validator.Validate(i, v[2])
		assert.Error(t, err)
		assert.EqualError(t, maxErrorItems(path, validator.in, items), err.Error())

		// ExclusiveMaximum
		items.ExclusiveMaximum = true
		err = validator.Validate(i, v[1])
		assert.Error(t, err)
		assert.EqualError(t, maxErrorItems(path, validator.in, items), err.Error())

		// Minimum
		err = validator.Validate(i, v[3])
		assert.NoError(t, err)
		err = validator.Validate(i, v[4])
		assert.Error(t, err)
		assert.EqualError(t, minErrorItems(path, validator.in, items), err.Error())

		// ExclusiveMinimum
		items.ExclusiveMinimum = true
		err = validator.Validate(i, v[3])
		assert.Error(t, err)
		assert.EqualError(t, minErrorItems(path, validator.in, items), err.Error())

		// Enum
		err = validator.Validate(i, v[5])
		assert.Error(t, err)
		assert.EqualError(t, enumFailItems(path, validator.in, items, v[5]), err.Error())

		// Valid passes
		assert.NoError(t, validator.Validate(i, v[6]))
	}

}

func TestStringItemsValidation(t *testing.T) {
	items := spec.NewItems().WithMinLength(3).WithMaxLength(5).WithPattern(`^[a-z]+$`)
	items.WithEnum("aaa", "bbb", "ccc")
	parent := spec.QueryParam("tags").CollectionOf(items, "")
	path := parent.Name + ".1"
	validator := &itemsValidator{items, parent, parent.Name, parent.In}

	// required
	err := validator.Validate(1, "")
	assert.Error(t, err)
	assert.EqualError(t, minLengthErrorItems(path, validator.in, items), err.Error())

	// MaxLength
	err = validator.Validate(1, "abcdef")
	assert.Error(t, err)
	assert.EqualError(t, maxLengthErrorItems(path, validator.in, items), err.Error())

	// MinLength
	err = validator.Validate(1, "a")
	assert.Error(t, err)
	assert.EqualError(t, minLengthErrorItems(path, validator.in, items), err.Error())

	// Pattern
	err = validator.Validate(1, "a394")
	assert.Error(t, err)
	assert.EqualError(t, patternFailItems(path, validator.in, items), err.Error())

	// Enum
	err = validator.Validate(1, "abcde")
	assert.Error(t, err)
	assert.EqualError(t, enumFailItems(path, validator.in, items, "abcde"), err.Error())

	// Valid passes
	err = validator.Validate(1, "bbb")
	assert.NoError(t, err)
}

func TestArrayItemsValidation(t *testing.T) {
	items := spec.NewItems().CollectionOf(stringItems, "").WithMinItems(1).WithMaxItems(5).UniqueValues()
	items.WithEnum("aaa", "bbb", "ccc")
	parent := spec.QueryParam("tags").CollectionOf(items, "")
	path := parent.Name + ".1"
	validator := &itemsValidator{items, parent, parent.Name, parent.In}

	// MinItems
	err := validator.Validate(1, []string{})
	assert.Error(t, err)
	assert.Error(t, minItemsErrorItems(path, validator.in, items), err.Error())
	// MaxItems
	err = validator.Validate(1, []string{"a", "b", "c", "d", "e", "f"})
	assert.Error(t, err)
	assert.Error(t, maxItemsErrorItems(path, validator.in, items), err.Error())
	// UniqueItems
	err = validator.Validate(1, []string{"a", "a"})
	assert.Error(t, err)
	assert.Error(t, duplicatesErrorItems(path, validator.in), err.Error())

	// Enum
	err = validator.Validate(1, []string{"a", "b", "a"})
	assert.Error(t, err)
	assert.Error(t, enumFailItems(path, validator.in, items, []string{"a", "b", "c"}), err.Error())

	// Items
	strItems := spec.NewItems().WithMinLength(3).WithMaxLength(5).WithPattern(`^[a-z]+$`)
	items = spec.NewItems().CollectionOf(strItems, "").WithMinItems(1).WithMaxItems(5).UniqueValues()
	validator = &itemsValidator{items, parent, parent.Name, parent.In}

	err = validator.Validate(1, []string{"aa", "bbb", "ccc"})
	assert.Error(t, err)
	assert.EqualError(t, minLengthErrorItems(path+".0", parent.In, strItems), err.Error())
}
