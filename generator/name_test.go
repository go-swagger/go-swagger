// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package generator

import (
	"testing"

	"github.com/go-openapi/spec"
	"github.com/go-openapi/swag/mangling"
	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"

	golangfuncs "github.com/go-swagger/go-swagger/generator/internal/funcmaps/golang"
)

func testNameMangler() mangling.NameMangler {
	return mangling.NewNameMangler(
		mangling.WithGoNamePrefixFunc(golangfuncs.PrefixForName),
	)
}

func TestExportGoName_Explicit(t *testing.T) {
	t.Parallel()

	mangler := testNameMangler()

	cases := []struct {
		name     string
		raw      string
		explicit bool
		want     string
	}{
		{
			name:     "preserves intentional casing from x-go-name",
			raw:      "NoTls",
			explicit: true,
			want:     "NoTls",
		},
		{
			name:     "exports lowercase explicit name without initialism rewrite",
			raw:      "noTls",
			explicit: true,
			want:     "NoTls",
		},
		{
			name:     "preserves VmId",
			raw:      "VmId",
			explicit: true,
			want:     "VmId",
		},
		{
			name:     "preserves all-caps explicit acronym",
			raw:      "ID",
			explicit: true,
			want:     "ID",
		},
		{
			name:     "empty raw name becomes Empty",
			raw:      "",
			explicit: true,
			want:     "Empty",
		},
		{
			name:     "special prefix uses PrefixForName",
			raw:      "-myField",
			explicit: true,
			want:     golangfuncs.PrefixForName("-myField"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := exportGoName(tc.raw, tc.explicit, mangler)
			assert.EqualT(t, tc.want, got)
		})
	}
}

func TestExportGoName_Implicit(t *testing.T) {
	t.Parallel()

	mangler := testNameMangler()

	cases := []struct {
		name string
		raw  string
	}{
		{name: "swagger property key", raw: "backend_tls_skip_verify"},
		{name: "camelCase property key", raw: "vmId"},
		{name: "snake_case property key", raw: "record_id"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			want := mangler.ToGoName(mangler.ToGoName(tc.raw))
			got := exportGoName(tc.raw, false, mangler)
			assert.EqualT(t, want, got)
		})
	}
}

func TestSchemaGoName(t *testing.T) {
	t.Parallel()

	mangler := testNameMangler()

	cases := []struct {
		name     string
		fallback string
		schema   func() *spec.Schema
		want     string
	}{
		{
			name:     "x-go-name overrides fallback and skips initialisms",
			fallback: "backendTlsSkipVerify",
			schema: func() *spec.Schema {
				sch := spec.Schema{}
				sch.AddExtension(xGoName, "NoTls")
				return &sch
			},
			want: "NoTls",
		},
		{
			name:     "x-go-name exports lowercase value",
			fallback: "vmId",
			schema: func() *spec.Schema {
				sch := spec.Schema{}
				sch.AddExtension(xGoName, "VmId")
				return &sch
			},
			want: "VmId",
		},
		{
			name:     "x-go-name preserves acronym casing",
			fallback: "recordId",
			schema: func() *spec.Schema {
				sch := spec.Schema{}
				sch.AddExtension(xGoName, "ID")
				return &sch
			},
			want: "ID",
		},
		{
			name:     "without x-go-name uses normal mangling",
			fallback: "backend_tls_skip_verify",
			schema: func() *spec.Schema {
				return &spec.Schema{}
			},
			want: mangler.ToGoName(mangler.ToGoName("backend_tls_skip_verify")),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			sch := tc.schema()
			require.NotNil(t, sch)

			got := schemaGoName(sch, tc.fallback, mangler)
			assert.EqualT(t, tc.want, got)
		})
	}
}

func TestExtensionGoNameOrError_InvalidType(t *testing.T) {
	t.Parallel()

	param := spec.Parameter{}
	param.Extensions = spec.Extensions{xGoName: []any{"not", "a", "string"}}

	_, err := extensionGoNameOrError(param.Extensions, "fallback", testNameMangler())
	require.Error(t, err)
	assert.StringContainsT(t, err.Error(), `must be a string`)
}

func TestExtensionGoName_InvalidTypeFallsBack(t *testing.T) {
	t.Parallel()

	mangler := testNameMangler()
	ext := spec.Extensions{xGoName: []any{"not", "a", "string"}}
	want := mangler.ToGoName(mangler.ToGoName("backend_tls_skip_verify"))

	got := extensionGoName(ext, "backend_tls_skip_verify", mangler)
	assert.EqualT(t, want, got)
}

func TestSchemaGoName_InvalidTypeFallsBack(t *testing.T) {
	t.Parallel()

	mangler := testNameMangler()
	sch := spec.Schema{}
	sch.AddExtension(xGoName, 42)
	want := mangler.ToGoName(mangler.ToGoName("plain_tls_flag"))

	got := schemaGoName(&sch, "plain_tls_flag", mangler)
	assert.EqualT(t, want, got)
}
