// SPDX-FileCopyrightText: Copyright 2026 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

//go:build ignore
// +build ignore

package models

import (
	"encoding/json"
	"testing"

	"github.com/go-openapi/strfmt"
)

func TestNullArrayItemsModel2382(t *testing.T) {
	var values FooArray
	if err := json.Unmarshal([]byte("[null]"), &values); err != nil {
		t.Fatal(err)
	}
	if err := values.Validate(strfmt.Default); err == nil {
		t.Fatal("expected null array item to fail model validation")
	}

	var nullableValues NullableFooArray
	if err := json.Unmarshal([]byte("[null]"), &nullableValues); err != nil {
		t.Fatal(err)
	}
	if err := nullableValues.Validate(strfmt.Default); err != nil {
		t.Fatalf("expected explicitly nullable array item to pass model validation: %v", err)
	}
}
