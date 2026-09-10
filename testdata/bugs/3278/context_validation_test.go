// SPDX-FileCopyrightText: Copyright 2026 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

//go:build ignore
// +build ignore

package models

import (
	"context"
	"testing"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
)

func TestContextValidateContinuesPastZeroOrder(t *testing.T) {
	request := &BatchOrderRequest{
		Orders: []Order{
			{},
			{ID: "client-supplied value"},
		},
	}

	ctx := validate.WithOperationRequest(context.Background())
	if err := request.ContextValidate(ctx, strfmt.Default); err == nil {
		t.Fatal("expected read-only ID in the second order to fail context validation")
	}
}
