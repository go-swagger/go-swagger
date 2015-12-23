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

// Package client contains a client to send http requests
// to a swagger API. This implementation is untyped
package client

import "fmt"

type methodAndPath struct {
	Method      string
	PathPattern string
	Schemes     []string
}

// NewAPIError creates a new API error
func NewAPIError(opName string, payload []byte, code int) *APIError {
	return &APIError{
		OperationName: opName,
		Payload:       payload,
		Code:          code,
	}
}

// APIError wraps an error model and captures the status code
type APIError struct {
	Method        string
	Path          string
	OperationName string
	Payload       []byte
	Code          int
}

func (a *APIError) Error() string {
	return fmt.Sprintf("%s (status %d): %+v ", a.OperationName, a.Code, string(a.Payload))
}
