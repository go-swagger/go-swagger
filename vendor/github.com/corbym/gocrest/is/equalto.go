package is

import (
	"fmt"
	"github.com/corbym/gocrest"
	"reflect"
)

//EqualTo checks if two values are equal. Uses DeepEqual (could be slow).
//Returns a matcher that will return true if two values are equal.
func EqualTo(expected interface{}) *gocrest.Matcher {
	match := new(gocrest.Matcher)
	match.Describe = fmt.Sprintf("value equal to <%v>", expected)
	match.Matches = func(actual interface{}) bool {
		return reflect.DeepEqual(expected, actual)
	}

	return match
}
