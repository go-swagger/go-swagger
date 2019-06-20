package then

import (
	"fmt"
	"github.com/corbym/gocrest"
)

//AssertThat calls a given matcher and fails the test with a message if the matcher doesn't match.
func AssertThat(t gocrest.TestingT, actual interface{}, m *gocrest.Matcher) {
	t.Helper()
	matches := m.Matches(actual)
	if !matches {
		t.Errorf("%s\nExpected: %s"+
			"\n     but: <%s>\n",
			m.ReasonString,
			m.Describe,
			actualAsString(m, actual),
		)
	}
}

func actualAsString(matcher *gocrest.Matcher, actual interface{}) string {
	if matcher.Actual != "" {
		return matcher.Actual
	}
	return fmt.Sprintf("%v", actual)
}
