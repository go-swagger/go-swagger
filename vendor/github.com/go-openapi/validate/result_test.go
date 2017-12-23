package validate

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test AddError() uniqueness
func TestResult_AddError(t *testing.T) {
	r := Result{}
	r.AddErrors(fmt.Errorf("One error"))
	r.AddErrors(fmt.Errorf("Another error"))
	r.AddErrors(fmt.Errorf("One error"))
	r.AddErrors(fmt.Errorf("One error"))
	r.AddErrors(fmt.Errorf("One error"))
	r.AddErrors(fmt.Errorf("One error"), fmt.Errorf("Another error"))

	assert.Len(t, r.Errors, 2)
	assert.Contains(t, r.Errors, fmt.Errorf("One error"))
	assert.Contains(t, r.Errors, fmt.Errorf("Another error"))
}

func TestResult_AddNilError(t *testing.T) {
	r := Result{}
	r.AddErrors(nil)
	assert.Len(t, r.Errors, 0)

	errArray := []error{fmt.Errorf("One Error"), nil, fmt.Errorf("Another error")}
	r.AddErrors(errArray...)
	assert.Len(t, r.Errors, 2)
}

func TestResult_AddWarnings(t *testing.T) {
	r := Result{}
	r.AddErrors(fmt.Errorf("One Error"))
	assert.Len(t, r.Errors, 1)
	assert.Len(t, r.Warnings, 0)

	r.AddWarnings(fmt.Errorf("One Warning"))
	assert.Len(t, r.Errors, 1)
	assert.Len(t, r.Warnings, 1)
}

func TestResult_Merge(t *testing.T) {
	r := Result{}
	r.AddErrors(fmt.Errorf("One Error"))
	r.AddWarnings(fmt.Errorf("One Warning"))
	r.Inc()
	assert.Len(t, r.Errors, 1)
	assert.Len(t, r.Warnings, 1)
	assert.Equal(t, r.MatchCount, 1)

	// Merge with same
	r2 := Result{}
	r2.AddErrors(fmt.Errorf("One Error"))
	r2.AddWarnings(fmt.Errorf("One Warning"))
	r2.Inc()

	r.Merge(&r2)

	assert.Len(t, r.Errors, 1)
	assert.Len(t, r.Warnings, 1)
	assert.Equal(t, r.MatchCount, 2)

	// Merge with new
	r3 := Result{}
	r3.AddErrors(fmt.Errorf("New Error"))
	r3.AddWarnings(fmt.Errorf("New Warning"))
	r3.Inc()

	r.Merge(&r3)

	assert.Len(t, r.Errors, 2)
	assert.Len(t, r.Warnings, 2)
	assert.Equal(t, r.MatchCount, 3)
}

func TestResult_MergeAsErrors(t *testing.T) {
	r := Result{}
	r.AddErrors(fmt.Errorf("One Error"))
	r.AddWarnings(fmt.Errorf("One Warning"))
	r.Inc()
	assert.Len(t, r.Errors, 1)
	assert.Len(t, r.Warnings, 1)
	assert.Equal(t, r.MatchCount, 1)

	// same
	r2 := Result{}
	r2.AddErrors(fmt.Errorf("One Error"))
	r2.AddWarnings(fmt.Errorf("One Warning"))
	r2.Inc()

	// new
	r3 := Result{}
	r3.AddErrors(fmt.Errorf("New Error"))
	r3.AddWarnings(fmt.Errorf("New Warning"))
	r3.Inc()

	r.MergeAsErrors(&r2, &r3)

	assert.Len(t, r.Errors, 4) // One Warning added to Errors
	assert.Len(t, r.Warnings, 1)
	assert.Equal(t, r.MatchCount, 3)
}

func TestResult_MergeAsWarnings(t *testing.T) {
	r := Result{}
	r.AddErrors(fmt.Errorf("One Error"))
	r.AddWarnings(fmt.Errorf("One Warning"))
	r.Inc()
	assert.Len(t, r.Errors, 1)
	assert.Len(t, r.Warnings, 1)
	assert.Equal(t, r.MatchCount, 1)

	// same
	r2 := Result{}
	r2.AddErrors(fmt.Errorf("One Error"))
	r2.AddWarnings(fmt.Errorf("One Warning"))
	r2.Inc()

	// new
	r3 := Result{}
	r3.AddErrors(fmt.Errorf("New Error"))
	r3.AddWarnings(fmt.Errorf("New Warning"))
	r3.Inc()

	r.MergeAsWarnings(&r2, &r3)

	assert.Len(t, r.Errors, 1) // One Warning added to Errors
	assert.Len(t, r.Warnings, 4)
	assert.Equal(t, r.MatchCount, 3)
}

func TestResult_IsValid(t *testing.T) {
	r := Result{}

	assert.True(t, r.IsValid())
	assert.False(t, r.HasErrors())

	r.AddWarnings(fmt.Errorf("One Warning"))
	assert.True(t, r.IsValid())
	assert.False(t, r.HasErrors())

	r.AddErrors(fmt.Errorf("One Error"))
	assert.False(t, r.IsValid())
	assert.True(t, r.HasErrors())
}

func TestResult_HasWarnings(t *testing.T) {
	r := Result{}

	assert.False(t, r.HasWarnings())

	r.AddErrors(fmt.Errorf("One Error"))
	assert.False(t, r.HasWarnings())

	r.AddWarnings(fmt.Errorf("One Warning"))
	assert.True(t, r.HasWarnings())
}

func TestResult_KeepRelevantErrors(t *testing.T) {
	r := Result{}
	r.AddErrors(fmt.Errorf("One Error"))
	r.AddErrors(fmt.Errorf("IMPORTANT!Another Error"))
	r.AddWarnings(fmt.Errorf("One warning"))
	r.AddWarnings(fmt.Errorf("IMPORTANT!Another warning"))
	assert.Len(t, r.KeepRelevantErrors().Errors, 1)
	assert.Len(t, r.KeepRelevantErrors().Warnings, 1)
}

func TestResult_AsError(t *testing.T) {
	r := Result{}
	assert.Nil(t, r.AsError())
	r.AddErrors(fmt.Errorf("One Error"))
	r.AddErrors(fmt.Errorf("Additional Error"))
	res := r.AsError()
	if assert.NotNil(t, res) {
		assert.Contains(t, res.Error(), "validation failure list:") // Expected from pkg errors
		assert.Contains(t, res.Error(), "One Error")                // Expected from pkg errors
		assert.Contains(t, res.Error(), "Additional Error")         // Expected from pkg errors
	}
}
