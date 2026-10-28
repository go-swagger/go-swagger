// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package generator

import (
	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"
)

// makeGenDefinition is the per-model variant used by tests: it analyzes the spec afresh at each call.
func makeGenDefinition(name, pkg string, schema spec.Schema, specDoc *loads.Document, opts *GenOpts) (*GenDefinition, error) {
	return makeGenDefinitionWithAnalysis(name, pkg, schema, specDoc, opts, nil)
}
