package is

import (
	"fmt"
	"github.com/corbym/gocrest"
)

//AllOf takes some matchers and checks if all the matchers return true.
//Returns a matcher that performs the the test on the input matchers.
func AllOf(allMatchers ...*gocrest.Matcher) *gocrest.Matcher {
	matcher := new(gocrest.Matcher)
	matcher.Describe = fmt.Sprintf("all of (%s)", describe(allMatchers, "and"))
	matcher.Matches = matchAll(allMatchers, matcher)
	return matcher
}

func matchAll(allMatchers []*gocrest.Matcher, allOf *gocrest.Matcher) func(actual interface{}) bool {
	return func(actual interface{}) bool {
		matches := true
		allOf.AppendActual(fmt.Sprintf("actual <%v>", actual))
		for x := 0; x < len(allMatchers); x++ {
			if !allMatchers[x].Matches(actual) {
				matches = false
			}
			allOf.AppendActual(allMatchers[x].Actual)
		}
		return matches
	}
}
