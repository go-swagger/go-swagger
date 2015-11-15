
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
	io "io"
	runtime "runtime"
	unsafe "unsafe"
)

type MockReader interface {
	io.Reader
	oglemock.MockObject
}

type mockReader struct {
	controller  oglemock.Controller
	description string
}

func NewMockReader(
	c oglemock.Controller,
	desc string) MockReader {
	return &mockReader{
		controller:  c,
		description: desc,
	}
}

func (m *mockReader) Oglemock_Id() uintptr {
	return uintptr(unsafe.Pointer(m))
}

func (m *mockReader) Oglemock_Description() string {
	return m.description
}

func (m *mockReader) Read(p0 []uint8) (o0 int, o1 error) {
	// Get a file name and line number for the caller.
	_, file, line, _ := runtime.Caller(1)

	// Hand the call off to the controller, which does most of the work.
	retVals := m.controller.HandleMethodCall(
		m,
		"Read",
		file,
		line,
		[]interface{}{p0})

	if len(retVals) != 2 {
		panic(fmt.Sprintf("mockReader.Read: invalid return values: %v", retVals))
	}

	// o0 int
	if retVals[0] != nil {
		o0 = retVals[0].(int)
	}

	// o1 error
	if retVals[1] != nil {
		o1 = retVals[1].(error)
	}

	return
}

type MockWriter interface {
	io.Writer
	oglemock.MockObject
}

type mockWriter struct {
	controller  oglemock.Controller
	description string
}

func NewMockWriter(
	c oglemock.Controller,
	desc string) MockWriter {
	return &mockWriter{
		controller:  c,
		description: desc,
	}
}

func (m *mockWriter) Oglemock_Id() uintptr {
	return uintptr(unsafe.Pointer(m))
}

func (m *mockWriter) Oglemock_Description() string {
	return m.description
}

func (m *mockWriter) Write(p0 []uint8) (o0 int, o1 error) {
	// Get a file name and line number for the caller.
	_, file, line, _ := runtime.Caller(1)

	// Hand the call off to the controller, which does most of the work.
	retVals := m.controller.HandleMethodCall(
		m,
		"Write",
		file,
		line,
		[]interface{}{p0})

	if len(retVals) != 2 {
		panic(fmt.Sprintf("mockWriter.Write: invalid return values: %v", retVals))
	}

	// o0 int
	if retVals[0] != nil {
		o0 = retVals[0].(int)
	}

	// o1 error
	if retVals[1] != nil {
		o1 = retVals[1].(error)
	}

	return
}
