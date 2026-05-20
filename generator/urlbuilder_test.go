// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package generator

import (
	"bytes"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
)

func TestURLBuilder_SimplePathParams(t *testing.T) {
	t.Parallel()

	t.Run("should construct an operation builder", func(t *testing.T) {
		gen, err := opBuilder("simplePathParams", "../fixtures/codegen/todolist.url.simple.yml")
		require.NoError(t, err)

		t.Run("should make an operation", func(t *testing.T) {
			op, err := gen.MakeOperation()
			require.NoError(t, err)

			t.Run("should generate go code", func(t *testing.T) {
				buf := bytes.NewBuffer(nil)
				opts := opts()
				require.NoError(t, opts.templates.MustGet("serverUrlbuilder").Execute(buf, op))

				t.Run("should format go code", func(t *testing.T) {
					ff, err := opts.LanguageOpts.FormatContent("simple_path_params.go", buf.Bytes())
					require.NoErrorf(t, err, "generated operation is badly formatted:\n%s", buf.String())

					t.Run("generated code should verify a few patterns", func(t *testing.T) {
						res := string(ff)
						assertInCode(t, "var _path = \"/singleValuePath/{siString}/{siInt}/{siInt32}/{siInt64}/{siFloat}/{siFloat32}/{siFloat64}/{siBool}\"", res)
						assertInCode(t, "siBool := conv.FormatBool(o.SiBool)", res)
						assertInCode(t, `if siBool != ""`, res)
						assertInCode(t, `_path = strings.ReplaceAll(_path, "{siBool}", siBool)`, res)
						assertInCode(t, `return nil, errors.New("siBool is required on SimplePathParamsURL")`, res)
						assertInCode(t, `siFloat := conv.FormatFloat(o.SiFloat)`, res)
						assertInCode(t, `if siFloat != ""`, res)
						assertInCode(t, `_path = strings.ReplaceAll(_path, "{siFloat}", siFloat)`, res)
						assertInCode(t, `siFloat32 := conv.FormatFloat(o.SiFloat32)`, res)
						assertInCode(t, `if siFloat32 != "" `, res)
						assertInCode(t, `_path = strings.ReplaceAll(_path, "{siFloat32}", siFloat32)`, res)
						assertInCode(t, `return nil, errors.New("siFloat32 is required on SimplePathParamsURL")`, res)
						assertInCode(t, `siFloat64 := conv.FormatFloat(o.SiFloat64)`, res)
						assertInCode(t, `if siFloat64 != ""`, res)
						assertInCode(t, `_path = strings.ReplaceAll(_path, "{siFloat64}", siFloat64)`, res)
						assertInCode(t, `return nil, errors.New("siFloat64 is required on SimplePathParamsURL")`, res)
						assertInCode(t, `siInt := conv.FormatInteger(o.SiInt)`, res)
						assertInCode(t, `if siInt != ""`, res)
						assertInCode(t, `_path = strings.ReplaceAll(_path, "{siInt}", siInt)`, res)
						assertInCode(t, `return nil, errors.New("siInt is required on SimplePathParamsURL")`, res)
						assertInCode(t, `siInt32 := conv.FormatInteger(o.SiInt32)`, res)
						assertInCode(t, `if siInt32 != ""`, res)
						assertInCode(t, `_path = strings.ReplaceAll(_path, "{siInt32}", siInt32)`, res)
						assertInCode(t, `return nil, errors.New("siInt32 is required on SimplePathParamsURL")`, res)
						assertInCode(t, `siInt64 := conv.FormatInteger(o.SiInt64)`, res)
						assertInCode(t, `if siInt64 != ""`, res)
						assertInCode(t, `_path = strings.ReplaceAll(_path, "{siInt64}", siInt64)`, res)
						assertInCode(t, `return nil, errors.New("siInt64 is required on SimplePathParamsURL")`, res)
						assertInCode(t, `siString := o.SiString`, res)
						assertInCode(t, `if siString != ""`, res)
						assertInCode(t, `_path = strings.ReplaceAll(_path, "{siString}", siString)`, res)
						assertInCode(t, `return nil, errors.New("siString is required on SimplePathParamsURL")`, res)
						assertInCode(t, `result.Path = golangswaggerpaths.Join(_basePath, _path)`, res)

						assertNotInCode(t, `result.RawQuery = qs.Encode()`, res)
					})
				})
			})
		})
	})
}

func TestURLBuilder_SimpleQueryParams(t *testing.T) {
	t.Parallel()

	t.Run("should construct an operation builder", func(t *testing.T) {
		gen, err := opBuilder("simpleQueryParams", "../fixtures/codegen/todolist.url.simple.yml")
		require.NoError(t, err)

		t.Run("should make an operation", func(t *testing.T) {
			op, err := gen.MakeOperation()
			require.NoError(t, err)

			t.Run("should generate go code", func(t *testing.T) {
				buf := bytes.NewBuffer(nil)
				opts := opts()
				require.NoError(t, opts.templates.MustGet("serverUrlbuilder").Execute(buf, op))

				t.Run("should format go code", func(t *testing.T) {
					ff, err := opts.LanguageOpts.FormatContent("simple_query_params.go", buf.Bytes())
					require.NoErrorf(t, err, "generated operation is badly formatted:\n%s", buf.String())

					t.Run("generated code should verify a few patterns", func(t *testing.T) {
						res := string(ff)
						assertInCode(t, "var _path = \"/singleValueQuery/{id}\"", res)
						assertInCode(t, `id := conv.FormatInteger(o.ID)`, res)
						assertInCode(t, `if id != ""`, res)
						assertInCode(t, `_path = strings.ReplaceAll(_path, "{id}", id)`, res)
						assertInCode(t, `return nil, errors.New("id is required on SimpleQueryParamsURL")`, res)
						assertInCode(t, `result.Path = golangswaggerpaths.Join(_basePath, _path)`, res)
						assertInCode(t, `qs := make(url.Values)`, res)
						assertInCode(t, `siBoolQ := conv.FormatBool(o.SiBool)`, res)
						assertInCode(t, `if siBoolQ != ""`, res)
						assertInCode(t, `qs.Set("siBool", siBoolQ)`, res)
						assertInCode(t, `var siFloatQ string`, res)
						assertInCode(t, `if o.SiFloat != nil`, res)
						assertInCode(t, `siFloatQ = conv.FormatFloat(*o.SiFloat)`, res)
						assertInCode(t, `if siFloatQ != ""`, res)
						assertInCode(t, `qs.Set("siFloat", siFloatQ)`, res)
						assertInCode(t, `var siFloat32Q string`, res)
						assertInCode(t, `if o.SiFloat32 != nil {`, res)
						assertInCode(t, `siFloat32Q = conv.FormatFloat(*o.SiFloat32)`, res)
						assertInCode(t, `if siFloat32Q != ""`, res)
						assertInCode(t, `qs.Set("siFloat32", siFloat32Q)`, res)
						assertInCode(t, `siFloat64Q := conv.FormatFloat(o.SiFloat64)`, res)
						assertInCode(t, `if siFloat64Q != ""`, res)
						assertInCode(t, `qs.Set("siFloat64", siFloat64Q)`, res)
						assertInCode(t, `var siIntQ string`, res)
						assertInCode(t, `if o.SiInt != nil`, res)
						assertInCode(t, `siIntQ = conv.FormatInteger(*o.SiInt)`, res)
						assertInCode(t, `if siIntQ != ""`, res)
						assertInCode(t, `qs.Set("siInt", siIntQ)`, res)
						assertInCode(t, `var siInt32Q string`, res)
						assertInCode(t, `if o.SiInt32 != nil`, res)
						assertInCode(t, `siInt32Q = conv.FormatInteger(*o.SiInt32)`, res)
						assertInCode(t, `if siInt32Q != ""`, res)
						assertInCode(t, `qs.Set("siInt32", siInt32Q)`, res)
						assertInCode(t, `siInt64Q := conv.FormatInteger(o.SiInt64)`, res)
						assertInCode(t, `if siInt64Q != ""`, res)
						assertInCode(t, `qs.Set("siInt64", siInt64Q)`, res)
						assertInCode(t, `siStringQ := o.SiString`, res)
						assertInCode(t, `if siStringQ != ""`, res)
						assertInCode(t, `qs.Set("siString", siStringQ)`, res)
						assertInCode(t, `result.RawQuery = qs.Encode()`, res)
					})
				})
			})
		})
	})
}

func TestURLBuilder_ArrayQueryParams(t *testing.T) {
	t.Parallel()

	testArrayQueryParams(t, "../fixtures/codegen/todolist.url.simple.yml", "")
}

func TestURLBuilder_ArrayQueryParams_BasePath(t *testing.T) {
	t.Parallel()

	testArrayQueryParams(t, "../fixtures/codegen/todolist.url.basepath.yml", "/v1/api")
}

func testArrayQueryParams(t *testing.T, filePath, basePath string) {
	t.Run("should construct an operation builder", func(t *testing.T) {
		gen, err := opBuilder("arrayQueryParams", filePath)
		require.NoError(t, err)

		t.Run("should make an operation", func(t *testing.T) {
			op, err := gen.MakeOperation()
			require.NoError(t, err)

			t.Run("should generate go code", func(t *testing.T) {
				buf := bytes.NewBuffer(nil)
				opts := opts()
				require.NoError(t, opts.templates.MustGet("serverUrlbuilder").Execute(buf, op))

				t.Run("should format go code", func(t *testing.T) {
					ff, err := opts.LanguageOpts.FormatContent("array_query_params.go", buf.Bytes())
					require.NoErrorf(t, err, "generated operation is badly formatted:\n%s", buf.String())

					t.Run("generated code should verify a few patterns", func(t *testing.T) {
						res := string(ff)

						assertInCode(t, "var _path = \"/arrayValueQuery/{id}\"", res)
						assertInCode(t, `id := conv.FormatInteger(o.ID)`, res)
						assertInCode(t, `if id != ""`, res)
						assertInCode(t, `_path = strings.ReplaceAll(_path, "{id}", id)`, res)
						assertInCode(t, `return nil, errors.New("id is required on ArrayQueryParamsURL")`, res)
						assertInCode(t, "_basePath := o._basePath", res)

						if basePath != "" {
							assertInCode(t, `if _basePath == ""`, res)
							assertInCode(t, `_basePath = "`+basePath+`"`, res)
						} else {
							assertNotInCode(t, `_basePath = "`+basePath+`"`, res)
						}

						assertInCode(t, `result.Path = golangswaggerpaths.Join(_basePath, _path)`, res)
						assertInCode(t, `qs := make(url.Values)`, res)

						assertInCode(t, `var siBoolIR []string`, res)
						assertInCode(t, `for _, siBoolI := range o.SiBool {`, res)
						assertInCode(t, `siBoolIS := conv.FormatBool(siBoolI)`, res)
						assertInCode(t, `if siBoolIS != ""`, res)
						assertInCode(t, `siBoolIR = append(siBoolIR, siBoolIS)`, res)
						assertInCode(t, `siBool := stringutils.JoinByFormat(siBoolIR, "ssv")`, res)
						assertInCode(t, `if len(siBool) > 0`, res)
						assertInCode(t, `qsv := siBool[0]`, res)
						assertInCode(t, `qs.Set("siBool", qsv)`, res)

						assertInCode(t, `var siFloatIR []string`, res)
						assertInCode(t, `for _, siFloatI := range o.SiFloat {`, res)
						assertInCode(t, `siFloatIS := conv.FormatFloat(siFloatI)`, res)
						assertInCode(t, `if siFloatIS != ""`, res)
						assertInCode(t, `siFloatIR = append(siFloatIR, siFloatIS)`, res)
						assertInCode(t, `siFloat := stringutils.JoinByFormat(siFloatIR, "multi")`, res)
						assertInCode(t, `for _, qsv := range siFloat`, res)
						assertInCode(t, `qs.Add("siFloat", qsv)`, res)

						assertInCode(t, `var siFloat32IR []string`, res)
						assertInCode(t, `for _, siFloat32I := range o.SiFloat32 {`, res)
						assertInCode(t, `siFloat32IS := conv.FormatFloat(siFloat32I)`, res)
						assertInCode(t, `if siFloat32IS != ""`, res)
						assertInCode(t, `siFloat32IR = append(siFloat32IR, siFloat32IS)`, res)
						assertInCode(t, `siFloat32 := stringutils.JoinByFormat(siFloat32IR, "")`, res)
						assertInCode(t, `if len(siFloat32) > 0`, res)
						assertInCode(t, `qsv := siFloat32[0]`, res)
						assertInCode(t, `qs.Set("siFloat32", qsv)`, res)

						assertInCode(t, `var siFloat64IR []string`, res)
						assertInCode(t, `for _, siFloat64I := range o.SiFloat64 {`, res)
						assertInCode(t, `siFloat64IS := conv.FormatFloat(siFloat64I)`, res)
						assertInCode(t, `if siFloat64IS != ""`, res)
						assertInCode(t, `siFloat64IR = append(siFloat64IR, siFloat64IS)`, res)
						assertInCode(t, `siFloat64 := stringutils.JoinByFormat(siFloat64IR, "pipes")`, res)
						assertInCode(t, `if len(siFloat64) > 0`, res)
						assertInCode(t, `qsv := siFloat64[0]`, res)
						assertInCode(t, `qs.Set("siFloat64", qsv)`, res)

						assertInCode(t, `var siIntIR []string`, res)
						assertInCode(t, `for _, siIntI := range o.SiInt {`, res)
						assertInCode(t, `siIntIS := conv.FormatInteger(siIntI)`, res)
						assertInCode(t, `if siIntIS != ""`, res)
						assertInCode(t, `siIntIR = append(siIntIR, siIntIS)`, res)
						assertInCode(t, `siInt := stringutils.JoinByFormat(siIntIR, "pipes")`, res)
						assertInCode(t, `if len(siInt) > 0`, res)
						assertInCode(t, `qsv := siInt[0]`, res)
						assertInCode(t, `qs.Set("siInt", qsv)`, res)

						assertInCode(t, `var siInt32IR []string`, res)
						assertInCode(t, `for _, siInt32I := range o.SiInt32 {`, res)
						assertInCode(t, `siInt32IS := conv.FormatInteger(siInt32I)`, res)
						assertInCode(t, `if siInt32IS != ""`, res)
						assertInCode(t, `siInt32IR = append(siInt32IR, siInt32IS)`, res)
						assertInCode(t, `siInt32 := stringutils.JoinByFormat(siInt32IR, "tsv")`, res)
						assertInCode(t, `if len(siInt32) > 0`, res)
						assertInCode(t, `qsv := siInt32[0]`, res)
						assertInCode(t, `qs.Set("siInt32", qsv)`, res)

						assertInCode(t, `var siInt64IR []string`, res)
						assertInCode(t, `for _, siInt64I := range o.SiInt64 {`, res)
						assertInCode(t, `siInt64IS := conv.FormatInteger(siInt64I)`, res)
						assertInCode(t, `if siInt64IS != ""`, res)
						assertInCode(t, `siInt64IR = append(siInt64IR, siInt64IS)`, res)
						assertInCode(t, `siInt64 := stringutils.JoinByFormat(siInt64IR, "ssv")`, res)
						assertInCode(t, `if len(siInt64) > 0`, res)
						assertInCode(t, `qsv := siInt64[0]`, res)
						assertInCode(t, `qs.Set("siInt64", qsv)`, res)

						assertInCode(t, `var siStringIR []string`, res)
						assertInCode(t, `for _, siStringI := range o.SiString {`, res)
						assertInCode(t, `siStringIS := siStringI`, res)
						assertInCode(t, `if siStringIS != ""`, res)
						assertInCode(t, `siStringIR = append(siStringIR, siStringIS)`, res)
						assertInCode(t, `siString := stringutils.JoinByFormat(siStringIR, "csv")`, res)
						assertInCode(t, `if len(siString) > 0`, res)
						assertInCode(t, `qsv := siString[0]`, res)
						assertInCode(t, `qs.Set("siString", qsv)`, res)

						assertInCode(t, `var siNestedIR []string`, res)
						assertInCode(t, `for _, siNestedI := range o.SiNested`, res)
						assertInCode(t, `var siNestedIIR []string`, res)
						assertInCode(t, `for _, siNestedII := range siNestedI`, res)
						assertInCode(t, `var siNestedIIIR []string`, res)
						assertInCode(t, `for _, siNestedIII := range siNestedII`, res)
						assertInCode(t, `siNestedIIIS := siNestedIII`, res)
						assertInCode(t, `if siNestedIIIS != ""`, res)
						assertInCode(t, `siNestedIIIR = append(siNestedIIIR, siNestedIIIS)`, res)
						assertInCode(t, `siNestedIIS := stringutils.JoinByFormat(siNestedIIIR, "csv")`, res)
						assertInCode(t, `if len(siNestedIIS) > 0`, res)
						assertInCode(t, `siNestedIISs := siNestedIIS[0]`, res)
						assertInCode(t, `if siNestedIISs != ""`, res)
						assertInCode(t, `siNestedIIR = append(siNestedIIR, siNestedIISs)`, res)
						assertInCode(t, `siNestedIS := stringutils.JoinByFormat(siNestedIIR, "pipes")`, res)
						assertInCode(t, `if len(siNestedIS) > 0`, res)
						assertInCode(t, `siNestedISs := siNestedIS[0]`, res)
						assertInCode(t, `if siNestedISs != ""`, res)
						assertInCode(t, `siNestedIR = append(siNestedIR, siNestedISs)`, res)
						assertInCode(t, `siNested := stringutils.JoinByFormat(siNestedIR, "multi")`, res)
						assertInCode(t, `for _, qsv := range siNested`, res)
						assertInCode(t, `qs.Add("siNested", qsv)`, res)
						assertInCode(t, `result.RawQuery = qs.Encode()`, res)
					})
				})
			})
		})
	})
}

func TestURLBuilder_Issue2167(t *testing.T) {
	t.Parallel()

	t.Run("with happy path", func(t *testing.T) {
		t.Run("should construct an operation builder", func(t *testing.T) {
			gen, err := opBuilder("xGoNameInParams", "../fixtures/enhancements/2167/swagger.yml")
			require.NoError(t, err)

			t.Run("should make an operation", func(t *testing.T) {
				op, err := gen.MakeOperation()
				require.NoError(t, err)

				t.Run("should generate go code", func(t *testing.T) {
					buf := bytes.NewBuffer(nil)
					opts := opts()
					require.NoError(t, opts.templates.MustGet("serverUrlbuilder").Execute(buf, op))

					t.Run("should format go code", func(t *testing.T) {
						ff, err := opts.LanguageOpts.FormatContent("get_test_test_name_urlbuilder.go", buf.Bytes())
						require.NoErrorf(t, err, "generated operation is badly formatted:\n%s", buf.String())

						t.Run("generated code should verify a few patterns", func(t *testing.T) {
							res := string(ff)
							assertRegexpInCode(t, `(?m)^\tMyPathName\s+string$`, res)
							assertRegexpInCode(t, `(?m)^\tTestRegion\s+string$`, res)
							assertRegexpInCode(t, `(?m)^\tMyQueryCount\s+\*int64$`, res)
							assertRegexpInCode(t, `(?m)^\tTestLimit\s+\*int64$`, res)
						})
					})
				})
			})
		})

		t.Run("with error path", func(t *testing.T) {
			t.Run("should construct an operation builder", func(t *testing.T) {
				gen, err := opBuilder("xGoNameInParams", "../fixtures/enhancements/2167/swagger-error.yml")
				require.NoError(t, err)

				t.Run("should NOT make an operation", func(t *testing.T) {
					_, err = gen.MakeOperation()
					require.Error(t, err)
					assert.StringContainsT(t, err.Error(), `GET /test/{test_name}, parameter "test_name": "x-go-name" field must be a string, not a []interface {}`)
				})
			})
		})
	})
}
