
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

package assertions

import "fmt"

// ShouldPanic receives a void, niladic function and expects to recover a panic.
func ShouldPanic(actual interface{}, expected ...interface{}) (message string) {
	if fail := need(0, expected); fail != success {
		return fail
	}

	action, _ := actual.(func())

	if action == nil {
		message = shouldUseVoidNiladicFunction
		return
	}

	defer func() {
		recovered := recover()
		if recovered == nil {
			message = shouldHavePanicked
		} else {
			message = success
		}
	}()
	action()

	return
}

// ShouldNotPanic receives a void, niladic function and expects to execute the function without any panic.
func ShouldNotPanic(actual interface{}, expected ...interface{}) (message string) {
	if fail := need(0, expected); fail != success {
		return fail
	}

	action, _ := actual.(func())

	if action == nil {
		message = shouldUseVoidNiladicFunction
		return
	}

	defer func() {
		recovered := recover()
		if recovered != nil {
			message = fmt.Sprintf(shouldNotHavePanicked, recovered)
		} else {
			message = success
		}
	}()
	action()

	return
}

// ShouldPanicWith receives a void, niladic function and expects to recover a panic with the second argument as the content.
func ShouldPanicWith(actual interface{}, expected ...interface{}) (message string) {
	if fail := need(1, expected); fail != success {
		return fail
	}

	action, _ := actual.(func())

	if action == nil {
		message = shouldUseVoidNiladicFunction
		return
	}

	defer func() {
		recovered := recover()
		if recovered == nil {
			message = shouldHavePanicked
		} else {
			if equal := ShouldEqual(recovered, expected[0]); equal != success {
				message = serializer.serialize(expected[0], recovered, fmt.Sprintf(shouldHavePanickedWith, expected[0], recovered))
			} else {
				message = success
			}
		}
	}()
	action()

	return
}

// ShouldNotPanicWith receives a void, niladic function and expects to recover a panic whose content differs from the second argument.
func ShouldNotPanicWith(actual interface{}, expected ...interface{}) (message string) {
	if fail := need(1, expected); fail != success {
		return fail
	}

	action, _ := actual.(func())

	if action == nil {
		message = shouldUseVoidNiladicFunction
		return
	}

	defer func() {
		recovered := recover()
		if recovered == nil {
			message = success
		} else {
			if equal := ShouldEqual(recovered, expected[0]); equal == success {
				message = fmt.Sprintf(shouldNotHavePanickedWith, expected[0])
			} else {
				message = success
			}
		}
	}()
	action()

	return
}
