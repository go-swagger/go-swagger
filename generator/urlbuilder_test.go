package generator

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestURLBuilder_SimplePathParams(t *testing.T) {
	gen, err := opBuilder("simplePathParams", "../fixtures/codegen/todolist.url.simple.yml")
	require.NoError(t, err)

	op, err := gen.MakeOperation()
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	opts := opts()
	err = templates.MustGet("serverUrlbuilder").Execute(buf, op)
	require.NoError(t, err)

	ff, err := opts.LanguageOpts.FormatContent("simple_path_params.go", buf.Bytes())
	require.NoErrorf(t, err, buf.String())

	res := string(ff)
	assertInCode(t, "var _path = \"/singleValuePath/{siString}/{siInt}/{siInt32}/{siInt64}/{siFloat}/{siFloat32}/{siFloat64}/{siBool}\"", res)
	assertInCode(t, "siBool := swag.FormatBool(o.SiBool)", res)
	assertInCode(t, `if siBool != ""`, res)
	assertInCode(t, `_path = strings.Replace(_path, "{siBool}", siBool, -1)`, res)
	assertInCode(t, `return nil, errors.New("siBool is required on SimplePathParamsURL")`, res)
	assertInCode(t, `siFloat := swag.FormatFloat64(o.SiFloat)`, res)
	assertInCode(t, `if siFloat != ""`, res)
	assertInCode(t, `_path = strings.Replace(_path, "{siFloat}", siFloat, -1)`, res)
	assertInCode(t, `siFloat32 := swag.FormatFloat32(o.SiFloat32)`, res)
	assertInCode(t, `if siFloat32 != "" `, res)
	assertInCode(t, `_path = strings.Replace(_path, "{siFloat32}", siFloat32, -1)`, res)
	assertInCode(t, `return nil, errors.New("siFloat32 is required on SimplePathParamsURL")`, res)
	assertInCode(t, `siFloat64 := swag.FormatFloat64(o.SiFloat64)`, res)
	assertInCode(t, `if siFloat64 != ""`, res)
	assertInCode(t, `_path = strings.Replace(_path, "{siFloat64}", siFloat64, -1)`, res)
	assertInCode(t, `return nil, errors.New("siFloat64 is required on SimplePathParamsURL")`, res)
	assertInCode(t, `siInt := swag.FormatInt64(o.SiInt)`, res)
	assertInCode(t, `if siInt != ""`, res)
	assertInCode(t, `_path = strings.Replace(_path, "{siInt}", siInt, -1)`, res)
	assertInCode(t, `return nil, errors.New("siInt is required on SimplePathParamsURL")`, res)
	assertInCode(t, `siInt32 := swag.FormatInt32(o.SiInt32)`, res)
	assertInCode(t, `if siInt32 != ""`, res)
	assertInCode(t, `_path = strings.Replace(_path, "{siInt32}", siInt32, -1)`, res)
	assertInCode(t, `return nil, errors.New("siInt32 is required on SimplePathParamsURL")`, res)
	assertInCode(t, `siInt64 := swag.FormatInt64(o.SiInt64)`, res)
	assertInCode(t, `if siInt64 != ""`, res)
	assertInCode(t, `_path = strings.Replace(_path, "{siInt64}", siInt64, -1)`, res)
	assertInCode(t, `return nil, errors.New("siInt64 is required on SimplePathParamsURL")`, res)
	assertInCode(t, `siString := o.SiString`, res)
	assertInCode(t, `if siString != ""`, res)
	assertInCode(t, `_path = strings.Replace(_path, "{siString}", siString, -1)`, res)
	assertInCode(t, `return nil, errors.New("siString is required on SimplePathParamsURL")`, res)
	assertInCode(t, `result.Path = golangswaggerpaths.Join(_basePath, _path)`, res)
	assertNotInCode(t, `result.RawQuery = qs.Encode()`, res)
}

func TestURLBuilder_SimpleQueryParams(t *testing.T) {
	gen, err := opBuilder("simpleQueryParams", "../fixtures/codegen/todolist.url.simple.yml")
	require.NoError(t, err)

	op, err := gen.MakeOperation()
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	opts := opts()
	err = templates.MustGet("serverUrlbuilder").Execute(buf, op)
	require.NoError(t, err)

	ff, err := opts.LanguageOpts.FormatContent("simple_query_params.go", buf.Bytes())
	require.NoErrorf(t, err, buf.String())

	res := string(ff)
	assertInCode(t, "var _path = \"/singleValueQuery/{id}\"", res)
	assertInCode(t, `id := swag.FormatInt32(o.ID)`, res)
	assertInCode(t, `if id != ""`, res)
	assertInCode(t, `_path = strings.Replace(_path, "{id}", id, -1)`, res)
	assertInCode(t, `return nil, errors.New("id is required on SimpleQueryParamsURL")`, res)
	assertInCode(t, `result.Path = golangswaggerpaths.Join(_basePath, _path)`, res)
	assertInCode(t, `qs := make(url.Values)`, res)
	assertInCode(t, `siBoolQ := swag.FormatBool(o.SiBool)`, res)
	assertInCode(t, `if siBoolQ != ""`, res)
	assertInCode(t, `qs.Set("siBool", siBoolQ)`, res)
	assertInCode(t, `var siFloatQ string`, res)
	assertInCode(t, `if o.SiFloat != nil`, res)
	assertInCode(t, `siFloatQ = swag.FormatFloat64(*o.SiFloat)`, res)
	assertInCode(t, `if siFloatQ != ""`, res)
	assertInCode(t, `qs.Set("siFloat", siFloatQ)`, res)
	assertInCode(t, `var siFloat32Q string`, res)
	assertInCode(t, `if o.SiFloat32 != nil {`, res)
	assertInCode(t, `siFloat32Q = swag.FormatFloat32(*o.SiFloat32)`, res)
	assertInCode(t, `if siFloat32Q != ""`, res)
	assertInCode(t, `qs.Set("siFloat32", siFloat32Q)`, res)
	assertInCode(t, `siFloat64Q := swag.FormatFloat64(o.SiFloat64)`, res)
	assertInCode(t, `if siFloat64Q != ""`, res)
	assertInCode(t, `qs.Set("siFloat64", siFloat64Q)`, res)
	assertInCode(t, `var siIntQ string`, res)
	assertInCode(t, `if o.SiInt != nil`, res)
	assertInCode(t, `siIntQ = swag.FormatInt64(*o.SiInt)`, res)
	assertInCode(t, `if siIntQ != ""`, res)
	assertInCode(t, `qs.Set("siInt", siIntQ)`, res)
	assertInCode(t, `var siInt32Q string`, res)
	assertInCode(t, `if o.SiInt32 != nil`, res)
	assertInCode(t, `siInt32Q = swag.FormatInt32(*o.SiInt32)`, res)
	assertInCode(t, `if siInt32Q != ""`, res)
	assertInCode(t, `qs.Set("siInt32", siInt32Q)`, res)
	assertInCode(t, `siInt64Q := swag.FormatInt64(o.SiInt64)`, res)
	assertInCode(t, `if siInt64Q != ""`, res)
	assertInCode(t, `qs.Set("siInt64", siInt64Q)`, res)
	assertInCode(t, `siStringQ := o.SiString`, res)
	assertInCode(t, `if siStringQ != ""`, res)
	assertInCode(t, `qs.Set("siString", siStringQ)`, res)
	assertInCode(t, `result.RawQuery = qs.Encode()`, res)
}

func TestURLBuilder_ArrayQueryParams(t *testing.T) {
	testArrayQueryParams(t, "../fixtures/codegen/todolist.url.simple.yml", "")
}

func TestURLBuilder_ArrayQueryParams_BasePath(t *testing.T) {
	testArrayQueryParams(t, "../fixtures/codegen/todolist.url.basepath.yml", "/v1/api")
}

func testArrayQueryParams(t testing.TB, filePath, basePath string) {
	gen, err := opBuilder("arrayQueryParams", filePath)
	require.NoError(t, err)

	op, err := gen.MakeOperation()
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	opts := opts()
	err = templates.MustGet("serverUrlbuilder").Execute(buf, op)
	require.NoError(t, err)

	ff, err := opts.LanguageOpts.FormatContent("array_query_params.go", buf.Bytes())
	require.NoErrorf(t, err, buf.String())

	res := string(ff)

	assertInCode(t, "var _path = \"/arrayValueQuery/{id}\"", res)
	assertInCode(t, `id := swag.FormatInt32(o.ID)`, res)
	assertInCode(t, `if id != ""`, res)
	assertInCode(t, `_path = strings.Replace(_path, "{id}", id, -1)`, res)
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
	assertInCode(t, `siBoolIS := swag.FormatBool(siBoolI)`, res)
	assertInCode(t, `if siBoolIS != ""`, res)
	assertInCode(t, `siBoolIR = append(siBoolIR, siBoolIS)`, res)
	assertInCode(t, `siBool := swag.JoinByFormat(siBoolIR, "ssv")`, res)
	assertInCode(t, `if len(siBool) > 0`, res)
	assertInCode(t, `qsv := siBool[0]`, res)
	assertInCode(t, `qs.Set("siBool", qsv)`, res)

	assertInCode(t, `var siFloatIR []string`, res)
	assertInCode(t, `for _, siFloatI := range o.SiFloat {`, res)
	assertInCode(t, `siFloatIS := swag.FormatFloat64(siFloatI)`, res)
	assertInCode(t, `if siFloatIS != ""`, res)
	assertInCode(t, `siFloatIR = append(siFloatIR, siFloatIS)`, res)
	assertInCode(t, `siFloat := swag.JoinByFormat(siFloatIR, "multi")`, res)
	assertInCode(t, `for _, qsv := range siFloat`, res)
	assertInCode(t, `qs.Add("siFloat", qsv)`, res)

	assertInCode(t, `var siFloat32IR []string`, res)
	assertInCode(t, `for _, siFloat32I := range o.SiFloat32 {`, res)
	assertInCode(t, `siFloat32IS := swag.FormatFloat32(siFloat32I)`, res)
	assertInCode(t, `if siFloat32IS != ""`, res)
	assertInCode(t, `siFloat32IR = append(siFloat32IR, siFloat32IS)`, res)
	assertInCode(t, `siFloat32 := swag.JoinByFormat(siFloat32IR, "")`, res)
	assertInCode(t, `if len(siFloat32) > 0`, res)
	assertInCode(t, `qsv := siFloat32[0]`, res)
	assertInCode(t, `qs.Set("siFloat32", qsv)`, res)

	assertInCode(t, `var siFloat64IR []string`, res)
	assertInCode(t, `for _, siFloat64I := range o.SiFloat64 {`, res)
	assertInCode(t, `siFloat64IS := swag.FormatFloat64(siFloat64I)`, res)
	assertInCode(t, `if siFloat64IS != ""`, res)
	assertInCode(t, `siFloat64IR = append(siFloat64IR, siFloat64IS)`, res)
	assertInCode(t, `siFloat64 := swag.JoinByFormat(siFloat64IR, "pipes")`, res)
	assertInCode(t, `if len(siFloat64) > 0`, res)
	assertInCode(t, `qsv := siFloat64[0]`, res)
	assertInCode(t, `qs.Set("siFloat64", qsv)`, res)

	assertInCode(t, `var siIntIR []string`, res)
	assertInCode(t, `for _, siIntI := range o.SiInt {`, res)
	assertInCode(t, `siIntIS := swag.FormatInt64(siIntI)`, res)
	assertInCode(t, `if siIntIS != ""`, res)
	assertInCode(t, `siIntIR = append(siIntIR, siIntIS)`, res)
	assertInCode(t, `siInt := swag.JoinByFormat(siIntIR, "pipes")`, res)
	assertInCode(t, `if len(siInt) > 0`, res)
	assertInCode(t, `qsv := siInt[0]`, res)
	assertInCode(t, `qs.Set("siInt", qsv)`, res)

	assertInCode(t, `var siInt32IR []string`, res)
	assertInCode(t, `for _, siInt32I := range o.SiInt32 {`, res)
	assertInCode(t, `siInt32IS := swag.FormatInt32(siInt32I)`, res)
	assertInCode(t, `if siInt32IS != ""`, res)
	assertInCode(t, `siInt32IR = append(siInt32IR, siInt32IS)`, res)
	assertInCode(t, `siInt32 := swag.JoinByFormat(siInt32IR, "tsv")`, res)
	assertInCode(t, `if len(siInt32) > 0`, res)
	assertInCode(t, `qsv := siInt32[0]`, res)
	assertInCode(t, `qs.Set("siInt32", qsv)`, res)

	assertInCode(t, `var siInt64IR []string`, res)
	assertInCode(t, `for _, siInt64I := range o.SiInt64 {`, res)
	assertInCode(t, `siInt64IS := swag.FormatInt64(siInt64I)`, res)
	assertInCode(t, `if siInt64IS != ""`, res)
	assertInCode(t, `siInt64IR = append(siInt64IR, siInt64IS)`, res)
	assertInCode(t, `siInt64 := swag.JoinByFormat(siInt64IR, "ssv")`, res)
	assertInCode(t, `if len(siInt64) > 0`, res)
	assertInCode(t, `qsv := siInt64[0]`, res)
	assertInCode(t, `qs.Set("siInt64", qsv)`, res)

	assertInCode(t, `var siStringIR []string`, res)
	assertInCode(t, `for _, siStringI := range o.SiString {`, res)
	assertInCode(t, `siStringIS := siStringI`, res)
	assertInCode(t, `if siStringIS != ""`, res)
	assertInCode(t, `siStringIR = append(siStringIR, siStringIS)`, res)
	assertInCode(t, `siString := swag.JoinByFormat(siStringIR, "csv")`, res)
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
	assertInCode(t, `siNestedIIS := swag.JoinByFormat(siNestedIIIR, "csv")`, res)
	assertInCode(t, `if len(siNestedIIS) > 0`, res)
	assertInCode(t, `siNestedIISs := siNestedIIS[0]`, res)
	assertInCode(t, `if siNestedIISs != ""`, res)
	assertInCode(t, `siNestedIIR = append(siNestedIIR, siNestedIISs)`, res)
	assertInCode(t, `siNestedIS := swag.JoinByFormat(siNestedIIR, "pipes")`, res)
	assertInCode(t, `if len(siNestedIS) > 0`, res)
	assertInCode(t, `siNestedISs := siNestedIS[0]`, res)
	assertInCode(t, `if siNestedISs != ""`, res)
	assertInCode(t, `siNestedIR = append(siNestedIR, siNestedISs)`, res)
	assertInCode(t, `siNested := swag.JoinByFormat(siNestedIR, "multi")`, res)
	assertInCode(t, `for _, qsv := range siNested`, res)
	assertInCode(t, `qs.Add("siNested", qsv)`, res)

	assertInCode(t, `result.RawQuery = qs.Encode()`, res)
}

func TestURLBuilder_Issue2167(t *testing.T) {
	gen, err := opBuilder("xGoNameInParams", "../fixtures/enhancements/2167/swagger.yml")
	require.NoError(t, err)

	op, err := gen.MakeOperation()
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	opts := opts()
	err = templates.MustGet("serverUrlbuilder").Execute(buf, op)
	require.NoError(t, err)

	ff, err := opts.LanguageOpts.FormatContent("get_test_test_name_urlbuilder.go", buf.Bytes())
	require.NoErrorf(t, err, buf.String())

	res := string(ff)
	assertRegexpInCode(t, `(?m)^\tMyPathName\s+string$`, res)
	assertRegexpInCode(t, `(?m)^\tTestRegion\s+string$`, res)
	assertRegexpInCode(t, `(?m)^\tMyQueryCount\s+\*int64$`, res)
	assertRegexpInCode(t, `(?m)^\tTestLimit\s+\*int64$`, res)
}

func TestURLBuilder_Issue2167_Error(t *testing.T) {
	gen, err := opBuilder("xGoNameInParams", "../fixtures/enhancements/2167/swagger-error.yml")
	require.NoError(t, err)

	_, err = gen.MakeOperation()
	require.Error(t, err)
	assert.Contains(t, `GET /test/{test_name}, parameter "test_name": "x-go-name" field must be a string, not a []interface {}`, err.Error())
}
