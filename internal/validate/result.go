package validate

import "github.com/go-swagger/go-swagger/errors"

// Result represents a validation result
type Result struct {
	Errors     []error
	MatchCount int
}

// Merge merges this result with the other one, preserving match counts etc
func (r *Result) Merge(other *Result) *Result {
	if other == nil {
		return r
	}
	r.AddErrors(other.Errors...)
	r.MatchCount += other.MatchCount
	return r
}

// AddErrors adds errors to this validation result
func (r *Result) AddErrors(errors ...error) {
	r.Errors = append(r.Errors, errors...)
}

// IsValid returns true when this result is valid
func (r *Result) IsValid() bool {
	return len(r.Errors) == 0
}

// HasErrors returns true when this result is invalid
func (r *Result) HasErrors() bool {
	return !r.IsValid()
}

// Inc increments the match count
func (r *Result) Inc() {
	r.MatchCount++
}

func (r *Result) AsError() error {
	if r.IsValid() {
		return nil
	}
	return errors.CompositeValidationError(r.Errors...)
}
