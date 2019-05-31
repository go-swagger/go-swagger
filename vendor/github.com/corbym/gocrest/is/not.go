package is

import (
	"github.com/corbym/gocrest"
)

//Not negates the given matcher.
//Returns a matcher that returns logical not of the matcher given.
func Not(matcher *gocrest.Matcher) *gocrest.Matcher {
	match := new(gocrest.Matcher)
	match.Describe = "not(" + matcher.Describe + ")"
	match.Matches = func(actual interface{}) bool {
		matches := !matcher.Matches(actual)
		match.Actual = matcher.Actual
		return matches
	}
	return match
}
