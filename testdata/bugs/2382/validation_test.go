// SPDX-FileCopyrightText: Copyright 2026 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

//go:build ignore
// +build ignore

package operations

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
)

func TestNullArrayItems2382(t *testing.T) {
	route := &middleware.MatchedRoute{
		Consumer: runtime.JSONConsumer(),
		Formats:  strfmt.Default,
	}

	defaultParams := NewUpdateFooParams()
	defaultRequest := httptest.NewRequest("PUT", "/foo", strings.NewReader("[null]"))
	defaultRequest.Header.Set("Content-Type", "application/json")
	if err := defaultParams.BindRequest(defaultRequest, route); err == nil {
		t.Fatal("expected null array item to fail validation")
	}

	nullableParams := NewUpdateNullableFooParams()
	nullableRequest := httptest.NewRequest("PUT", "/nullable-foo", strings.NewReader("[null]"))
	nullableRequest.Header.Set("Content-Type", "application/json")
	if err := nullableParams.BindRequest(nullableRequest, route); err != nil {
		t.Fatalf("expected explicitly nullable array item to pass validation: %v", err)
	}
}
