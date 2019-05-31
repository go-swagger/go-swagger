package is

import (
	"fmt"
	"github.com/corbym/gocrest"
)

//AnyOf takes some matchers and checks if at least one of the matchers return true.
//Returns a matcher that performs the the test on the input matchers.
func AnyOf(allMatchers ...*gocrest.Matcher) *gocrest.Matcher {
	matcher := new(gocrest.Matcher)
	matcher.Describe = fmt.Sprintf("any of (%s)", describe(allMatchers, "or"))
	matcher.Matches = anyMatcherMatches(allMatchers, matcher)
	return matcher
}

func anyMatcherMatches(allMatchers []*gocrest.Matcher, anyOf *gocrest.Matcher) func(actual interface{}) bool {
	return func(actual interface{}) bool {
		matches := false
		anyOf.AppendActual(fmt.Sprintf("actual <%v>", actual))
		for x := 0; x < len(allMatchers); x++ {
			if allMatchers[x].Matches(actual) {
				matches = true
			}
			anyOf.AppendActual(allMatchers[x].Actual)
		}
		return matches
	}
}
