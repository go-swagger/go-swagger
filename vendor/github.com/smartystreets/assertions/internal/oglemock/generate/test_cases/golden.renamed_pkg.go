
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

// This file was auto-generated using createmock. See the following page for
// more information:
//
//     https://github.com/smartystreets/assertions/internal/oglemock
//

package some_pkg

import (
	fmt "fmt"
	oglemock "github.com/smartystreets/assertions/internal/oglemock"
	tony "github.com/smartystreets/assertions/internal/oglemock/generate/test_cases/renamed_pkg"
	runtime "runtime"
	unsafe "unsafe"
)

type MockSomeInterface interface {
	tony.SomeInterface
	oglemock.MockObject
}

type mockSomeInterface struct {
	controller  oglemock.Controller
	description string
}

func NewMockSomeInterface(
	c oglemock.Controller,
	desc string) MockSomeInterface {
	return &mockSomeInterface{
		controller:  c,
		description: desc,
	}
}

func (m *mockSomeInterface) Oglemock_Id() uintptr {
	return uintptr(unsafe.Pointer(m))
}

func (m *mockSomeInterface) Oglemock_Description() string {
	return m.description
}

func (m *mockSomeInterface) DoFoo(p0 int) (o0 int) {
	// Get a file name and line number for the caller.
	_, file, line, _ := runtime.Caller(1)

	// Hand the call off to the controller, which does most of the work.
	retVals := m.controller.HandleMethodCall(
		m,
		"DoFoo",
		file,
		line,
		[]interface{}{p0})

	if len(retVals) != 1 {
		panic(fmt.Sprintf("mockSomeInterface.DoFoo: invalid return values: %v", retVals))
	}

	// o0 int
	if retVals[0] != nil {
		o0 = retVals[0].(int)
	}

	return
}
