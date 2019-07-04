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
func TestMapDiff(t *testing.T) {
	mapA := map[string]interface{}{"abc": 1, "def": 2, "ghi": 3, "jkl": 4}
	added, deleted, common := FromStringMap(mapA).
		DiffsTo(mapA)
	assertThat(t, added, isMapWithItems(map[string]interface{}{}))
	assertThat(t, deleted, isMapWithItems(map[string]interface{}{}))
	commonDiffs := map[string]interface{}{"abc": Pair{1, 1}, "def": Pair{2, 2}, "ghi": Pair{3, 3}, "jkl": Pair{4, 4}}

	assertThat(t, common, isMapWithItems(commonDiffs))

	mapB := map[string]interface{}{"abc": 2, "ghi": 3, "jkl": 4, "xyz": 5, "fgh": 6}
	added, deleted, common = FromStringMap(mapA).
		DiffsTo(mapB)
	assertThat(t, added, isMapWithItems(map[string]interface{}{"xyz": 5, "fgh": 6}))
	assertThat(t, deleted, isMapWithItems(map[string]interface{}{"def": 2}))
	commonDiffs = map[string]interface{}{"abc": Pair{1, 2}, "ghi": Pair{3, 3}, "jkl": Pair{4, 4}}
	assertThat(t, common, isMapWithItems(commonDiffs))

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

func isMapWithItems(other map[string]interface{}) *gocrest.Matcher {
	matcher := new(gocrest.Matcher)
	matcher.Describe = fmt.Sprintf("map with items:%v", other)
	matcher.Matches = func(actual interface{}) bool {
		if actual == nil {
			return other == nil
		}
		if actualValue, ok := actual.(map[string]interface{}); ok {
			if len(actualValue) == 0 {
				return len(other) == 0
			}
			leftToMatch := len(actualValue)
			for keyActual, actualItem := range actualValue {
				for keyOther, otherItem := range other {
					if actualItem == otherItem && keyActual == keyOther {
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
