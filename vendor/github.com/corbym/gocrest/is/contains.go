package is

import (
	"fmt"
	"github.com/corbym/gocrest"
	"reflect"
	"strings"
)

//ValueContaining finds if x is contained in y.
// Acts like "ContainsAll", all elements given must be present in actual.
// If "expected" is an array or slice, we assume that actual is the same type.
// assertThat([]T, has.ValueContaining(a,b,c)) is also valid if variadic a,b,c are all type T.
// For maps, the expected must also be a map or a variadic and matches if both maps contain all key,values in expected
// or all variadic values respectively.
// For string, behaves like strings.Contains.
// Will panic if types cannot be converted correctly.
//Returns the Matcher that returns true if found.
func ValueContaining(expected ...interface{}) *gocrest.Matcher {
	match := new(gocrest.Matcher)
	correctVariadicExpected := correctExpectedValue(expected...)
	match.Describe = fmt.Sprintf("something that contains %v", correctVariadicExpected)
	match.Matches = func(actual interface{}) bool {
		expectedAsStr, expectedOk := expected[0].(string)
		actualAsStr, actualOk := actual.(string)
		if expectedOk && actualOk {
			return strings.Contains(actualAsStr, expectedAsStr)
		}
		actualValue := reflect.ValueOf(actual)
		expectedValue := reflect.ValueOf(correctVariadicExpected)
		switch actualValue.Kind() {
		case reflect.Array, reflect.Slice:
			return listContains(expectedValue, actualValue)
		case reflect.Map:
			if expectedValue.Kind() == reflect.Array || expectedValue.Kind() == reflect.Slice {
				return mapContainsList(expectedValue, actualValue)
			}
			return mapContains(expectedValue, actualValue)
		default:
			panic("cannot determine type of variadic actual, " + actualValue.String())
		}
	}
	return match
}

func mapContainsList(expected reflect.Value, mapValue reflect.Value) bool {
	contains := make(map[interface{}]bool)
	for i := 0; i < expected.Len(); i++ {
		for _, key := range mapValue.MapKeys() {
			itemValue := expected.Index(i).Interface()
			if mapValue.MapIndex(key).Interface() == itemValue {
				contains[itemValue] = true
			}
		}
	}
	return len(contains) == expected.Len()
}

func mapContains(expected reflect.Value, actual reflect.Value) bool {
	expectedKeys := expected.MapKeys()

	contains := make(map[interface{}]bool)
	for i := 0; i < len(expectedKeys); i++ {
		val := actual.MapIndex(expectedKeys[i])
		if val.IsValid() {
			if val.Interface() == expected.MapIndex(expectedKeys[i]).Interface() {
				contains[val] = true
			}
		}
	}
	return len(contains) == len(expected.MapKeys())
}

func listContains(expectedValue reflect.Value, actualValue reflect.Value) bool {
	contains := make(map[interface{}]bool)
	for i := 0; i < expectedValue.Len(); i++ {
		for y := 0; y < actualValue.Len(); y++ {
			exp := expectedValue.Index(i).Interface()
			act := actualValue.Index(y).Interface()
			if exp == act {
				contains[act] = true
			}
		}
	}
	return len(contains) == expectedValue.Len()
}

func correctExpectedValue(expected ...interface{}) interface{} {
	kind := reflect.ValueOf(expected[0]).Kind()
	if kind == reflect.Slice || kind == reflect.Map {
		return expected[0]
	}
	return expected
}
