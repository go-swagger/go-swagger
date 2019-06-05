package diff

import (
	"fmt"
	"testing"

	"github.com/corbym/gocrest"
	"github.com/corbym/gocrest/then"
)

var assertThat = then.AssertThat

func TestArrayDiff(t *testing.T) {
	listA := []string{"abc", "def", "ghi", "jkl"}
	added, deleted, common := FromStringArray(listA).
		DiffsTo(listA)
	assertThat(t, added, isListWithItems([]string{}))
	assertThat(t, deleted, isListWithItems([]string{}))
	assertThat(t, common, isListWithItems(listA))

	listB := []string{"abc", "ghi", "jkl", "xyz", "fgh"}
	added, deleted, common = FromStringArray(listA).
		DiffsTo(listB)
	assertThat(t, added, isListWithItems([]string{"xyz", "fgh"}))
	assertThat(t, deleted, isListWithItems([]string{"def"}))
	assertThat(t, common, isListWithItems([]string{"abc", "ghi", "jkl"}))

}

func isListWithItems(other []string) *gocrest.Matcher {
	matcher := new(gocrest.Matcher)
	matcher.Describe = fmt.Sprintf("list with items:%v", other)
	matcher.Matches = func(actual interface{}) bool {
		if actual == nil {
			return other == nil
		}
		if actualValue, ok := actual.([]string); ok {
			if len(actualValue) == 0 {
				return len(other) == 0
			}
			leftToMatch := len(actualValue)
			for _, actualItem := range actualValue {
				for _, otherItem := range other {
					if actualItem == otherItem {
						leftToMatch--
						break
					}
				}
			}
			return leftToMatch == 0
		}
		return false
	}
	return matcher
}
