
// Copyright 2015 go-swagger maintainers
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package govalidator

// Iterator is the function that accepts element of slice/array and its index
type Iterator func(interface{}, int)

// ResultIterator is the function that accepts element of slice/array and its index and returns any result
type ResultIterator func(interface{}, int) interface{}

// ConditionIterator is the function that accepts element of slice/array and its index and returns boolean
type ConditionIterator func(interface{}, int) bool

// Each iterates over the slice and apply Iterator to every item
func Each(array []interface{}, iterator Iterator) {
	for index, data := range array {
		iterator(data, index)
	}
}

// Map iterates over the slice and apply ResultIterator to every item. Returns new slice as a result.
func Map(array []interface{}, iterator ResultIterator) []interface{} {
	var result []interface{} = make([]interface{}, len(array))
	for index, data := range array {
		result[index] = iterator(data, index)
	}
	return result
}

// Find iterates over the slice and apply ConditionIterator to every item. Returns first item that meet ConditionIterator or nil otherwise.
func Find(array []interface{}, iterator ConditionIterator) interface{} {
	for index, data := range array {
		if iterator(data, index) {
			return data
		}
	}
	return nil
}

// Filter iterates over the slice and apply ConditionIterator to every item. Returns new slice.
func Filter(array []interface{}, iterator ConditionIterator) []interface{} {
	var result []interface{} = make([]interface{}, 0)
	for index, data := range array {
		if iterator(data, index) {
			result = append(result, data)
		}
	}
	return result
}

// Count iterates over the slice and apply ConditionIterator to every item. Returns count of items that meets ConditionIterator.
func Count(array []interface{}, iterator ConditionIterator) int {
	count := 0
	for index, data := range array {
		if iterator(data, index) {
			count = count + 1
		}
	}
	return count
}
