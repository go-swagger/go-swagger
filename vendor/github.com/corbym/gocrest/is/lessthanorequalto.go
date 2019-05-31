package is

import "github.com/corbym/gocrest"

//LessThanOrEqualTo is a short hand matcher for anyOf(lessThan(x), equalTo(x))
//Returns a matcher matching if actual <= expected (using deepEquals).
func LessThanOrEqualTo(expected interface{}) *gocrest.Matcher {
	matcher := new(gocrest.Matcher)
	matcher.Matches = func(actual interface{}) bool {
		anyOf := AnyOf(LessThan(expected), EqualTo(expected))
		anyOfMatches := anyOf.Matches(actual)
		matcher.Describe = anyOf.Describe
		return anyOfMatches
	}
	return matcher
}
