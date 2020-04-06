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

package generator

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-openapi/analysis"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUniqueOperationNameMangling(t *testing.T) {
	doc, err := loads.Spec("../fixtures/bugs/2213/fixture-2213.yaml")
	require.NoError(t, err)
	analyzed := analysis.New(doc.Spec())
	ops := gatherOperations(analyzed, nil)
	assert.Contains(t, ops, "GetFoo")
	assert.Contains(t, ops, "GetAFoo")
}

func TestUniqueOperationNames(t *testing.T) {
	doc, err := loads.Spec("../fixtures/codegen/todolist.simple.yml")
	require.NoError(t, err)

	sp := doc.Spec()
	sp.Paths.Paths["/tasks"].Post.ID = "saveTask"
	sp.Paths.Paths["/tasks"].Post.AddExtension("origName", "createTask")
	sp.Paths.Paths["/tasks/{id}"].Put.ID = "saveTask"
	sp.Paths.Paths["/tasks/{id}"].Put.AddExtension("origName", "updateTask")
	analyzed := analysis.New(sp)

	ops := gatherOperations(analyzed, nil)
	assert.Len(t, ops, 6)
	_, exists := ops["saveTask"]
	assert.True(t, exists)
	_, exists = ops["PutTasksID"]
	assert.True(t, exists)
}

func TestEmptyOperationNames(t *testing.T) {
	doc, err := loads.Spec("../fixtures/codegen/todolist.simple.yml")
	require.NoError(t, err)

	sp := doc.Spec()
	sp.Paths.Paths["/tasks"].Post.ID = ""
	sp.Paths.Paths["/tasks"].Post.AddExtension("origName", "createTask")
	sp.Paths.Paths["/tasks/{id}"].Put.ID = ""
	sp.Paths.Paths["/tasks/{id}"].Put.AddExtension("origName", "updateTask")
	analyzed := analysis.New(sp)

	ops := gatherOperations(analyzed, nil)
	assert.Len(t, ops, 6)
	_, exists := ops["PostTasks"]
	assert.True(t, exists)
	_, exists = ops["PutTasksID"]
	assert.True(t, exists)
}

func TestMakeResponseHeader(t *testing.T) {
	b, err := opBuilder("getTasks", "")
	require.NoError(t, err)

	hdr := findResponseHeader(&b.Operation, 200, "X-Rate-Limit")
	gh, er := b.MakeHeader("a", "X-Rate-Limit", *hdr)
	require.NoError(t, er)

	assert.True(t, gh.IsPrimitive)
	assert.Equal(t, "int32", gh.GoType)
	assert.Equal(t, "X-Rate-Limit", gh.Name)
}

func TestMakeResponseHeaderDefaultValues(t *testing.T) {
	b, err := opBuilder("getTasks", "")
	require.NoError(t, err)

	var testCases = []struct {
		name         string      // input
		typeStr      string      // expected type
		defaultValue interface{} // expected result
	}{
		{"Access-Control-Allow-Origin", "string", "*"},
		{"X-Rate-Limit", "int32", nil},
		{"X-Rate-Limit-Remaining", "int32", float64(42)},
		{"X-Rate-Limit-Reset", "int32", "1449875311"},
		{"X-Rate-Limit-Reset-Human", "string", "3 days"},
		{"X-Rate-Limit-Reset-Human-Number", "string", float64(3)},
	}

	for _, tc := range testCases {
		hdr := findResponseHeader(&b.Operation, 200, tc.name)
		require.NotNil(t, hdr)

		gh, er := b.MakeHeader("a", tc.name, *hdr)
		require.NoError(t, er)

		assert.True(t, gh.IsPrimitive)
		assert.Equal(t, tc.typeStr, gh.GoType)
		assert.Equal(t, tc.name, gh.Name)
		assert.Exactly(t, tc.defaultValue, gh.Default)
	}
}

func TestMakeResponse(t *testing.T) {
	b, err := opBuilder("getTasks", "")
	require.NoError(t, err)

	resolver := &typeResolver{ModelsPackage: b.ModelsPackage, Doc: b.Doc}
	resolver.KnownDefs = make(map[string]struct{})
	for k := range b.Doc.Spec().Definitions {
		resolver.KnownDefs[k] = struct{}{}
	}
	gO, err := b.MakeResponse("a", "getTasksSuccess", true, resolver, 200, b.Operation.Responses.StatusCodeResponses[200])
	require.NoError(t, err)

	assert.Len(t, gO.Headers, 6)
	assert.NotNil(t, gO.Schema)
	assert.True(t, gO.Schema.IsArray)
	assert.NotNil(t, gO.Schema.Items)
	assert.False(t, gO.Schema.IsAnonymous)
	assert.Equal(t, "[]*models.Task", gO.Schema.GoType)
}

func TestMakeResponse_WithAllOfSchema(t *testing.T) {
	b, err := methodPathOpBuilder("get", "/media/search", "../fixtures/codegen/instagram.yml")
	require.NoError(t, err)

	resolver := &typeResolver{ModelsPackage: b.ModelsPackage, Doc: b.Doc}
	resolver.KnownDefs = make(map[string]struct{})
	for k := range b.Doc.Spec().Definitions {
		resolver.KnownDefs[k] = struct{}{}
	}
	gO, err := b.MakeResponse("a", "get /media/search", true, resolver, 200, b.Operation.Responses.StatusCodeResponses[200])
	require.NoError(t, err)

	require.NotNil(t, gO.Schema)
	assert.Equal(t, "GetMediaSearchBody", gO.Schema.GoType)

	require.NotEmpty(t, b.ExtraSchemas)
	body := b.ExtraSchemas["GetMediaSearchBody"]
	require.NotEmpty(t, body.Properties)

	prop := body.Properties[0]
	assert.Equal(t, "data", prop.Name)
	// is in models only when definition is flattened: otherwise, ExtraSchema is rendered in operations package
	assert.Equal(t, "[]*GetMediaSearchBodyDataItems0", prop.GoType)

	items := b.ExtraSchemas["GetMediaSearchBodyDataItems0"]
	require.NotEmpty(t, items.AllOf)

	media := items.AllOf[0]
	// expect #definitions/media to be captured and reused by ExtraSchema
	assert.Equal(t, "models.Media", media.GoType)
}

func TestMakeOperationParam(t *testing.T) {
	b, err := opBuilder("getTasks", "")
	require.NoError(t, err)

	resolver := &typeResolver{ModelsPackage: b.ModelsPackage, Doc: b.Doc}
	gO, err := b.MakeParameter("a", resolver, b.Operation.Parameters[0], nil)
	require.NoError(t, err)

	assert.Equal(t, "size", gO.Name)
	assert.True(t, gO.IsPrimitive)
}

func TestMakeOperationParamItem(t *testing.T) {
	b, err := opBuilder("arrayQueryParams", "../fixtures/codegen/todolist.arrayquery.yml")
	require.NoError(t, err)
	resolver := &typeResolver{ModelsPackage: b.ModelsPackage, Doc: b.Doc}
	gO, err := b.MakeParameterItem("a", "siString", "ii", "siString", "a.SiString", "query", resolver, b.Operation.Parameters[1].Items, nil)
	require.NoError(t, err)
	assert.Nil(t, gO.Parent)
	assert.True(t, gO.IsPrimitive)
}

func TestMakeOperation(t *testing.T) {
	b, err := opBuilder("getTasks", "")
	require.NoError(t, err)
	gO, err := b.MakeOperation()
	require.NoError(t, err)
	assert.Equal(t, "getTasks", gO.Name)
	assert.Equal(t, "GET", gO.Method)
	assert.Equal(t, "/tasks", gO.Path)
	assert.Len(t, gO.Params, 2)
	assert.Len(t, gO.Responses, 1)
	assert.NotNil(t, gO.DefaultResponse)
	assert.NotNil(t, gO.SuccessResponse)
}

func TestRenderOperation_InstagramSearch(t *testing.T) {
	defer discardOutput()()

	b, err := methodPathOpBuilder("get", "/media/search", "../fixtures/codegen/instagram.yml")
	require.NoError(t, err)

	gO, err := b.MakeOperation()
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	opt := opts()

	require.NoError(t, opt.templates.MustGet("serverOperation").Execute(buf, gO))

	ff, err := opt.LanguageOpts.FormatContent("operation.go", buf.Bytes())
	require.NoErrorf(t, err, buf.String())

	res := string(ff)
	assertInCode(t, "type GetMediaSearchOKBody struct {", res)
	// codegen does not assumes objects are only in models
	// this is inlined
	assertInCode(t, "Data []*GetMediaSearchOKBodyDataItems0 `json:\"data\"`", res)
	assertInCode(t, "type GetMediaSearchOKBodyDataItems0 struct {", res)
	// this is a definition: expect this definition to be reused from the models pkg
	assertInCode(t, "models.Media", res)

	buf = bytes.NewBuffer(nil)
	require.NoError(t, opt.templates.MustGet("serverResponses").Execute(buf, gO))

	ff, err = opt.LanguageOpts.FormatContent("response.go", buf.Bytes())
	require.NoErrorf(t, err, buf.String())

	res = string(ff)
	// codegen does not assumes objects are only in models
	assertInCode(t, "type GetMediaSearchOK struct {", res)
	assertInCode(t, "GetMediaSearchOKBody", res)

	b, err = methodPathOpBuilderWithFlatten("get", "/media/search", "../fixtures/codegen/instagram.yml")
	require.NoError(t, err)

	gO, err = b.MakeOperation()
	require.NoError(t, err)

	buf = bytes.NewBuffer(nil)
	opt = opts()
	require.NoError(t, opt.templates.MustGet("serverOperation").Execute(buf, gO))

	ff, err = opt.LanguageOpts.FormatContent("operation.go", buf.Bytes())
	require.NoErrorf(t, err, buf.String())

	res = string(ff)
	assertNotInCode(t, "DataItems0", res)
	assertNotInCode(t, "models", res)

	buf = bytes.NewBuffer(nil)
	require.NoError(t, opt.templates.MustGet("serverResponses").Execute(buf, gO))

	ff, err = opt.LanguageOpts.FormatContent("operation.go", buf.Bytes())
	require.NoErrorf(t, err, buf.String())

	res = string(ff)
	assertInCode(t, "Payload *models.GetMediaSearchOKBody", res)
}

func methodPathOpBuilder(method, path, fname string) (codeGenOpBuilder, error) {
	defer discardOutput()()

	if fname == "" {
		fname = "../fixtures/codegen/todolist.simple.yml"
	}
	o := opts()
	o.Spec = fname
	specDoc, analyzed, err := o.analyzeSpec()
	if err != nil {
		return codeGenOpBuilder{}, err
	}
	op, ok := analyzed.OperationFor(method, path)
	if !ok {
		return codeGenOpBuilder{}, errors.New("No operation could be found for " + method + " " + path)
	}

	return codeGenOpBuilder{
		Name:          method + " " + path,
		Method:        method,
		Path:          path,
		APIPackage:    "restapi",
		ModelsPackage: "models",
		Principal:     "models.User",
		Target:        ".",
		Operation:     *op,
		Doc:           specDoc,
		Analyzed:      analyzed,
		Authed:        false,
		ExtraSchemas:  make(map[string]GenSchema),
		GenOpts:       o,
	}, nil
}

// methodPathOpBuilderWithFlatten prepares an operation build based on method and path, with spec full flattening
func methodPathOpBuilderWithFlatten(method, path, fname string) (codeGenOpBuilder, error) {
	defer discardOutput()()

	if fname == "" {
		fname = "../fixtures/codegen/todolist.simple.yml"
	}

	o := opBuildGetOpts(fname, true, false) // flatten: true, minimal: false
	o.Spec = fname
	specDoc, analyzed, err := o.analyzeSpec()
	if err != nil {
		return codeGenOpBuilder{}, err
	}
	op, ok := analyzed.OperationFor(method, path)
	if !ok {
		return codeGenOpBuilder{}, errors.New("No operation could be found for " + method + " " + path)
	}

	return codeGenOpBuilder{
		Name:          method + " " + path,
		Method:        method,
		Path:          path,
		APIPackage:    "restapi",
		ModelsPackage: "models",
		Principal:     "models.User",
		Target:        ".",
		Operation:     *op,
		Doc:           specDoc,
		Analyzed:      analyzed,
		Authed:        false,
		ExtraSchemas:  make(map[string]GenSchema),
		GenOpts:       opts(),
	}, nil
}

// opBuilderWithOpts prepares the making of an operation with spec flattening options
func opBuilderWithOpts(name, fname string, o *GenOpts) (codeGenOpBuilder, error) {
	defer discardOutput()()

	if fname == "" {
		// default fixture
		fname = "../fixtures/codegen/todolist.simple.yml"
	}

	o.Spec = fname
	specDoc, analyzed, err := o.analyzeSpec()
	if err != nil {
		return codeGenOpBuilder{}, err
	}

	method, path, op, ok := analyzed.OperationForName(name)
	if !ok {
		return codeGenOpBuilder{}, errors.New("No operation could be found for " + name)
	}

	return codeGenOpBuilder{
		Name:          name,
		Method:        method,
		Path:          path,
		BasePath:      specDoc.BasePath(),
		APIPackage:    "restapi",
		ModelsPackage: "models",
		Principal:     "models.User",
		Target:        ".",
		Operation:     *op,
		Doc:           specDoc,
		Analyzed:      analyzed,
		Authed:        false,
		ExtraSchemas:  make(map[string]GenSchema),
		GenOpts:       o,
	}, nil
}

func opBuildGetOpts(specName string, withFlatten bool, withMinimalFlatten bool) (opts *GenOpts) {
	opts = &GenOpts{}
	opts.ValidateSpec = true
	opts.FlattenOpts = &analysis.FlattenOpts{
		Expand:  !withFlatten,
		Minimal: withMinimalFlatten,
	}
	opts.Spec = specName
	if err := opts.EnsureDefaults(); err != nil {
		panic("Cannot initialize GenOpts")
	}
	return
}

// opBuilderWithFlatten prepares the making of an operation with spec full flattening prior to rendering
func opBuilderWithFlatten(name, fname string) (codeGenOpBuilder, error) {
	o := opBuildGetOpts(fname, true, false) // flatten: true, minimal: false
	return opBuilderWithOpts(name, fname, o)
}

/*
// opBuilderWithMinimalFlatten prepares the making of an operation with spec minimal flattening prior to rendering
func opBuilderWithMinimalFlatten(name, fname string) (codeGenOpBuilder, error) {
	o := opBuildGetOpts(fname, true, true) // flatten: true, minimal: true
	return opBuilderWithOpts(name, fname, o)
}
*/

// opBuilderWithExpand prepares the making of an operation with spec expansion prior to rendering
func opBuilderWithExpand(name, fname string) (codeGenOpBuilder, error) {
	o := opBuildGetOpts(fname, false, false) // flatten: false => expand
	return opBuilderWithOpts(name, fname, o)
}

// opBuilder prepares the making of an operation with spec minimal flattening (default for CLI)
func opBuilder(name, fname string) (codeGenOpBuilder, error) {
	o := opBuildGetOpts(fname, true, true) // flatten:true, minimal: true
	// some fixtures do not fully validate - skip this
	o.ValidateSpec = false
	return opBuilderWithOpts(name, fname, o)
}

func findResponseHeader(op *spec.Operation, code int, name string) *spec.Header {
	resp := op.Responses.Default
	if code > 0 {
		bb, ok := op.Responses.StatusCodeResponses[code]
		if ok {
			resp = &bb
		}
	}

	if resp == nil {
		return nil
	}

	hdr, ok := resp.Headers[name]
	if !ok {
		return nil
	}

	return &hdr
}

func TestDateFormat_Spec1(t *testing.T) {
	b, err := opBuilder("putTesting", "../fixtures/bugs/193/spec1.json")
	require.NoError(t, err)

	op, err := b.MakeOperation()
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	opts := opts()
	opts.defaultsEnsured = false
	opts.IsClient = true
	require.NoError(t, opts.EnsureDefaults())

	require.NoError(t, opts.templates.MustGet("clientParameter").Execute(buf, op))

	ff, err := opts.LanguageOpts.FormatContent("put_testing.go", buf.Bytes())
	require.NoErrorf(t, err, buf.String())

	assertInCode(t, "frTestingThis.String()", string(ff))
}

func TestDateFormat_Spec2(t *testing.T) {
	b, err := opBuilder("putTesting", "../fixtures/bugs/193/spec2.json")
	require.NoError(t, err)

	op, err := b.MakeOperation()
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	opts := opts()
	opts.defaultsEnsured = false
	opts.IsClient = true
	require.NoError(t, opts.EnsureDefaults())

	require.NoError(t, opts.templates.MustGet("clientParameter").Execute(buf, op))

	ff, err := opts.LanguageOpts.FormatContent("put_testing.go", buf.Bytes())
	require.NoErrorf(t, err, buf.String())

	res := string(ff)
	assertInCode(t, "o.TestingThis != nil {", res)
	assertInCode(t, "joinedTestingThis := o.bindParamTestingThis(reg)", res)
	assertInCode(t, `if err := r.SetFormParam("testingThis", joinedTestingThis...); err != nil {`, res)
	assertInCode(t, "func (o *PutTestingParams) bindParamTestingThis(formats strfmt.Registry) []string {", res)
	assertInCode(t, "testingThisIR := o.TestingThis", res)
	assertInCode(t, "var testingThisIC []string", res)
	assertInCode(t, "for _, testingThisIIR := range testingThisIR {", res)
	assertInCode(t, "testingThisIIV := testingThisIIR.String()", res)
	assertInCode(t, "testingThisIC = append(testingThisIC, testingThisIIV)", res)
	assertInCode(t, `testingThisIS := swag.JoinByFormat(testingThisIC, "")`, res)
	assertInCode(t, "return testingThisIS", res)
}

func TestBuilder_Issue1703(t *testing.T) {
	defer discardOutput()()

	dr := testCwd(t)

	opts := &GenOpts{
		Spec:              filepath.FromSlash("../fixtures/codegen/existing-model.yml"),
		IncludeModel:      true,
		IncludeHandler:    true,
		IncludeParameters: true,
		IncludeResponses:  true,
		IncludeMain:       true,
		APIPackage:        "restapi",
		ModelPackage:      "model",
		ServerPackage:     "server",
		ClientPackage:     "client",
		Target:            dr,
	}
	require.NoError(t, opts.EnsureDefaults())

	appGen, err := newAppGenerator("x-go-type-import-bug", nil, nil, opts)
	require.NoError(t, err)

	op, err := appGen.makeCodegenApp()
	require.NoError(t, err)

	for _, o := range op.Operations {
		buf := bytes.NewBuffer(nil)
		require.NoError(t, opts.templates.MustGet("serverResponses").Execute(buf, o))

		ff, err := appGen.GenOpts.LanguageOpts.FormatContent("response.go", buf.Bytes())
		require.NoErrorf(t, err, buf.String())

		assertInCode(t, "jwk \"github.com/user/package\"", string(ff))
	}
}

func TestBuilder_Issue287(t *testing.T) {
	defer discardOutput()()

	dr := testCwd(t)

	opts := &GenOpts{
		Spec:              filepath.FromSlash("../fixtures/bugs/287/swagger.yml"),
		IncludeModel:      true,
		IncludeHandler:    true,
		IncludeParameters: true,
		IncludeResponses:  true,
		IncludeMain:       true,
		APIPackage:        "restapi",
		ModelPackage:      "model",
		ServerPackage:     "server",
		ClientPackage:     "client",
		Target:            dr,
	}
	require.NoError(t, opts.EnsureDefaults())

	appGen, err := newAppGenerator("plainTexter", nil, nil, opts)
	require.NoError(t, err)

	op, err := appGen.makeCodegenApp()
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	require.NoError(t, opts.templates.MustGet("serverBuilder").Execute(buf, op))

	ff, err := appGen.GenOpts.LanguageOpts.FormatContent("put_testing.go", buf.Bytes())
	require.NoErrorf(t, err, buf.String())

	assertInCode(t, "case \"text/plain\":", string(ff))
}

func TestBuilder_Issue465(t *testing.T) {
	defer discardOutput()()

	dr := testCwd(t)

	opts := &GenOpts{
		Spec:              filepath.FromSlash("../fixtures/bugs/465/swagger.yml"),
		IncludeModel:      true,
		IncludeHandler:    true,
		IncludeParameters: true,
		IncludeResponses:  true,
		IncludeMain:       true,
		APIPackage:        "restapi",
		ModelPackage:      "model",
		ServerPackage:     "server",
		ClientPackage:     "client",
		Target:            dr,
		IsClient:          true,
	}
	require.NoError(t, opts.EnsureDefaults())

	appGen, err := newAppGenerator("plainTexter", nil, nil, opts)
	require.NoError(t, err)

	op, err := appGen.makeCodegenApp()
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	require.NoError(t, opts.templates.MustGet("clientFacade").Execute(buf, op))

	ff, err := appGen.GenOpts.LanguageOpts.FormatContent("put_testing.go", buf.Bytes())
	require.NoErrorf(t, err, buf.String())

	assertInCode(t, "/v1/fancyAPI", string(ff))
}

func TestBuilder_Issue500(t *testing.T) {
	defer discardOutput()()

	dr := testCwd(t)

	opts := &GenOpts{
		Spec:              filepath.FromSlash("../fixtures/bugs/500/swagger.yml"),
		IncludeModel:      true,
		IncludeHandler:    true,
		IncludeParameters: true,
		IncludeResponses:  true,
		IncludeMain:       true,
		APIPackage:        "restapi",
		ModelPackage:      "model",
		ServerPackage:     "server",
		ClientPackage:     "client",
		Target:            dr,
	}
	require.NoError(t, opts.EnsureDefaults())

	appGen, err := newAppGenerator("multiTags", nil, nil, opts)
	require.NoError(t, err)

	op, err := appGen.makeCodegenApp()
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	require.NoError(t, opts.templates.MustGet("serverBuilder").Execute(buf, op))

	ff, err := appGen.GenOpts.LanguageOpts.FormatContent("put_testing.go", buf.Bytes())
	require.NoErrorf(t, err, buf.String())

	res := string(ff)
	assertNotInCode(t, `o.handlers["GET"]["/payment/{invoice_id}/payments/{payment_id}"] = NewGetPaymentByID(o.context, o.GetPaymentByIDHandler)`, res)
	assertInCode(t, `o.handlers["GET"]["/payment/{invoice_id}/payments/{payment_id}"] = invoices.NewGetPaymentByID(o.context, o.InvoicesGetPaymentByIDHandler)`, res)
}

func TestGenClient_IllegalBOM(t *testing.T) {
	b, err := methodPathOpBuilder("get", "/v3/attachments/{attachmentId}", "../fixtures/bugs/727/swagger.json")
	require.NoError(t, err)

	op, err := b.MakeOperation()
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	opts := opts()
	opts.defaultsEnsured = false
	opts.IsClient = true
	require.NoError(t, opts.EnsureDefaults())

	require.NoError(t, opts.templates.MustGet("clientResponse").Execute(buf, op))
}

func TestGenClient_CustomFormatPath(t *testing.T) {
	b, err := methodPathOpBuilder("get", "/mosaic/experimental/series/{SeriesId}/mosaics", "../fixtures/bugs/789/swagger.yml")
	require.NoError(t, err)

	op, err := b.MakeOperation()
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	opts := opts()
	opts.defaultsEnsured = false
	opts.IsClient = true
	require.NoError(t, opts.EnsureDefaults())

	require.NoError(t, opts.templates.MustGet("clientParameter").Execute(buf, op))

	assertInCode(t, `if err := r.SetPathParam("SeriesId", o.SeriesID.String()); err != nil`, buf.String())
}

func TestGenClient_Issue733(t *testing.T) {
	b, err := opBuilder("get_characters_character_id_mail_mail_id", "../fixtures/bugs/733/swagger.json")
	require.NoError(t, err)

	op, err := b.MakeOperation()
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	opts := opts()
	opts.defaultsEnsured = false
	opts.IsClient = true
	require.NoError(t, opts.EnsureDefaults())

	require.NoError(t, opts.templates.MustGet("clientResponse").Execute(buf, op))

	assertInCode(t, "Labels []*int64 `json:\"labels\"`", buf.String())
}

func TestGenServerIssue890_ValidationTrueFlatteningTrue(t *testing.T) {
	defer discardOutput()()

	dr := testCwd(t)

	opts := &GenOpts{
		Spec:              filepath.FromSlash("../fixtures/bugs/890/swagger.yaml"),
		IncludeModel:      true,
		IncludeHandler:    true,
		IncludeParameters: true,
		IncludeResponses:  true,
		IncludeMain:       true,
		ValidateSpec:      true,
		APIPackage:        "restapi",
		ModelPackage:      "model",
		ServerPackage:     "server",
		ClientPackage:     "client",
		Target:            dr,
		IsClient:          true,
	}

	// Testing Server Generation
	require.NoError(t, opts.EnsureDefaults())

	// Full flattening
	opts.FlattenOpts.Expand = false
	opts.FlattenOpts.Minimal = false
	appGen, err := newAppGenerator("JsonRefOperation", nil, nil, opts)
	require.NoError(t, err)

	op, err := appGen.makeCodegenApp()
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	require.NoError(t, opts.templates.MustGet("serverOperation").Execute(buf, op.Operations[0]))

	filecontent, err := appGen.GenOpts.LanguageOpts.FormatContent("operation.go", buf.Bytes())
	require.NoErrorf(t, err, buf.String())

	res := string(filecontent)
	assertInCode(t, "GetHealthCheck", res)
}

func TestGenClientIssue890_ValidationTrueFlatteningTrue(t *testing.T) {
	defer discardOutput()()

	dr := testCwd(t)
	defer func() {
		_ = os.RemoveAll(filepath.Join(filepath.FromSlash(dr), "restapi"))
	}()

	opts := testGenOpts()
	opts.Spec = "../fixtures/bugs/890/swagger.yaml"
	opts.ValidateSpec = true
	opts.FlattenOpts.Minimal = false

	// Testing this is enough as there is only one operation which is specified as $ref.
	// If this doesn't get resolved then there will be an error definitely.
	require.NoError(t, GenerateClient("foo", nil, nil, opts))
}

func TestGenServerIssue890_ValidationFalseFlattenTrue(t *testing.T) {
	defer discardOutput()()

	dr := testCwd(t)

	opts := &GenOpts{
		Spec:              filepath.FromSlash("../fixtures/bugs/890/swagger.yaml"),
		IncludeModel:      true,
		IncludeHandler:    true,
		IncludeParameters: true,
		IncludeResponses:  true,
		IncludeMain:       true,
		APIPackage:        "restapi",
		ModelPackage:      "model",
		ServerPackage:     "server",
		ClientPackage:     "client",
		Target:            dr,
		IsClient:          true,
	}

	// Testing Server Generation
	require.NoError(t, opts.EnsureDefaults())

	// full flattening
	opts.FlattenOpts.Minimal = false
	appGen, err := newAppGenerator("JsonRefOperation", nil, nil, opts)
	require.NoError(t, err)

	op, err := appGen.makeCodegenApp()
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	err = opts.templates.MustGet("serverOperation").Execute(buf, op.Operations[0])
	require.NoError(t, err)

	filecontent, err := appGen.GenOpts.LanguageOpts.FormatContent("operation.go", buf.Bytes())
	require.NoErrorf(t, err, buf.String())

	res := string(filecontent)
	assertInCode(t, "GetHealthCheck", res)
}

func TestGenClientIssue890_ValidationFalseFlatteningTrue(t *testing.T) {
	defer discardOutput()()

	dr := testCwd(t)
	defer func() {
		_ = os.RemoveAll(filepath.Join(filepath.FromSlash(dr), "restapi"))
	}()

	opts := testGenOpts()
	opts.Spec = "../fixtures/bugs/890/swagger.yaml"
	opts.ValidateSpec = false
	// full flattening
	opts.FlattenOpts.Minimal = false
	// Testing this is enough as there is only one operation which is specified as $ref.
	// If this doesn't get resolved then there will be an error definitely.
	assert.NoError(t, GenerateClient("foo", nil, nil, opts))
}

func TestGenServerIssue890_ValidationFalseFlattenFalse(t *testing.T) {
	defer discardOutput()()

	dr := testCwd(t)

	opts := &GenOpts{
		Spec:              filepath.FromSlash("../fixtures/bugs/890/swagger.yaml"),
		IncludeModel:      true,
		IncludeHandler:    true,
		IncludeParameters: true,
		IncludeResponses:  true,
		IncludeMain:       true,
		ValidateSpec:      false,
		APIPackage:        "restapi",
		ModelPackage:      "model",
		ServerPackage:     "server",
		ClientPackage:     "client",
		Target:            dr,
		IsClient:          true,
	}

	// Testing Server Generation
	require.NoError(t, opts.EnsureDefaults())

	// minimal flattening
	opts.FlattenOpts.Minimal = true
	_, err := newAppGenerator("JsonRefOperation", nil, nil, opts)
	// if flatten is not set, expand takes over so this would resume normally
	assert.NoError(t, err)
}

func TestGenClientIssue890_ValidationFalseFlattenFalse(t *testing.T) {
	defer discardOutput()()

	dr := testCwd(t)
	defer func() {
		_ = os.RemoveAll(filepath.Join(filepath.FromSlash(dr), "restapi"))
	}()

	opts := testGenOpts()
	opts.Spec = "../fixtures/bugs/890/swagger.yaml"
	opts.ValidateSpec = false
	// minimal flattening
	opts.FlattenOpts.Minimal = true
	// Testing this is enough as there is only one operation which is specified as $ref.
	// If this doesn't get resolved then there will be an error definitely.
	// New: Now if flatten is false, expand takes over so server generation should resume normally
	assert.NoError(t, GenerateClient("foo", nil, nil, opts))
}

func TestGenServerIssue890_ValidationTrueFlattenFalse(t *testing.T) {
	defer discardOutput()()

	dr := testCwd(t)

	opts := &GenOpts{
		Spec:              filepath.FromSlash("../fixtures/bugs/890/swagger.yaml"),
		IncludeModel:      true,
		IncludeHandler:    true,
		IncludeParameters: true,
		IncludeResponses:  true,
		IncludeMain:       true,
		ValidateSpec:      true,
		APIPackage:        "restapi",
		ModelPackage:      "model",
		ServerPackage:     "server",
		ClientPackage:     "client",
		Target:            dr,
		IsClient:          true,
	}

	// Testing Server Generation
	require.NoError(t, opts.EnsureDefaults())

	// minimal flattening
	opts.FlattenOpts.Minimal = true

	_, err := newAppGenerator("JsonRefOperation", nil, nil, opts)
	// now if flatten is false, expand takes over so server generation should resume normally
	assert.NoError(t, err)
}

func TestGenServerWithTemplate(t *testing.T) {
	defer discardOutput()()

	dr := testCwd(t)

	tests := []struct {
		name      string
		opts      *GenOpts
		wantError bool
	}{
		{
			name: "None_existing_contributor_template",
			opts: &GenOpts{
				Spec:              filepath.FromSlash("../fixtures/bugs/890/swagger.yaml"),
				IncludeModel:      true,
				IncludeHandler:    true,
				IncludeParameters: true,
				IncludeResponses:  true,
				IncludeMain:       true,
				ValidateSpec:      true,
				APIPackage:        "restapi",
				ModelPackage:      "model",
				ServerPackage:     "server",
				ClientPackage:     "client",
				Target:            dr,
				IsClient:          true,
				Template:          "InvalidTemplate",
			},
			wantError: true,
		},
		{
			name: "Existing_contributor",
			opts: &GenOpts{
				Spec:              filepath.FromSlash("../fixtures/bugs/890/swagger.yaml"),
				IncludeModel:      true,
				IncludeHandler:    true,
				IncludeParameters: true,
				IncludeResponses:  true,
				IncludeMain:       true,
				ValidateSpec:      true,
				APIPackage:        "restapi",
				ModelPackage:      "model",
				ServerPackage:     "server",
				ClientPackage:     "client",
				Target:            dr,
				IsClient:          true,
				Template:          "stratoscale",
			},
			wantError: false,
		},
	}

	t.Run("codegen operations", func(t *testing.T) {
		for _, toPin := range tests {
			tt := toPin
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				// Testing Server Generation
				require.NoError(t, tt.opts.EnsureDefaults())

				// minimal flattening
				tt.opts.FlattenOpts.Minimal = true
				_, err := newAppGenerator("JsonRefOperation", nil, nil, tt.opts)
				if tt.wantError {
					require.Error(t, err)
				} else {
					require.NoError(t, err)
				}
			})
		}
	})
}

func TestGenClientIssue890_ValidationTrueFlattenFalse(t *testing.T) {
	defer discardOutput()()

	dr := testCwd(t)
	defer func() {
		_ = os.RemoveAll(filepath.Join(filepath.FromSlash(dr), "restapi"))
	}()

	opts := testGenOpts()
	opts.Spec = filepath.FromSlash("../fixtures/bugs/890/swagger.yaml")
	opts.ValidateSpec = true
	// Testing this is enough as there is only one operation which is specified as $ref.
	// If this doesn't get resolved then there will be an error definitely.
	// same here: now if flatten doesn't resume, expand takes over
	assert.NoError(t, GenerateClient("foo", nil, nil, opts))
}

// This tests that securityDefinitions generate stable code
func TestBuilder_Issue1214(t *testing.T) {
	defer discardOutput()()

	dr := testCwd(t)
	const any = `(.|\n)+`

	opts := &GenOpts{
		Spec:              filepath.FromSlash("../fixtures/bugs/1214/fixture-1214.yaml"),
		IncludeModel:      true,
		IncludeHandler:    true,
		IncludeParameters: true,
		IncludeResponses:  true,
		IncludeMain:       true,
		APIPackage:        "restapi",
		ModelPackage:      "model",
		ServerPackage:     "server",
		ClientPackage:     "client",
		Target:            dr,
		IsClient:          false,
	}
	require.NoError(t, opts.EnsureDefaults())

	appGen, e := newAppGenerator("fixture-1214", nil, nil, opts)
	require.NoError(t, e)

	op, e := appGen.makeCodegenApp()
	require.NoError(t, e)

	for i := 0; i < 5; i++ {
		buf := bytes.NewBuffer(nil)
		err := templates.MustGet("serverConfigureapi").Execute(buf, op)
		require.NoError(t, err)

		ff, err := appGen.GenOpts.LanguageOpts.FormatContent("fixture_1214_configure_api.go", buf.Bytes())
		require.NoErrorf(t, err, buf.String())

		res := string(ff)
		assertRegexpInCode(t, any+
			`api\.AAuth = func\(user string, pass string\)`+any+
			`api\.BAuth = func\(token string\)`+any+
			`api\.CAuth = func\(token string\)`+any+
			`api\.DAuth = func\(token string\)`+any+
			`api\.EAuth = func\(token string, scopes \[\]string\)`+any, res)

		buf = bytes.NewBuffer(nil)
		require.NoError(t, opts.templates.MustGet("serverBuilder").Execute(buf, op))

		ff, err = appGen.GenOpts.LanguageOpts.FormatContent("fixture_1214_server.go", buf.Bytes())
		require.NoErrorf(t, err, buf.String())

		res = string(ff)
		assertRegexpInCode(t, any+
			`AAuth: func\(user string, pass string\) \(interface{}, error\) {`+any+
			`BAuth: func\(token string\) \(interface{}, error\) {`+any+
			`CAuth: func\(token string\) \(interface{}, error\) {`+any+
			`DAuth: func\(token string\) \(interface{}, error\) {`+any+
			`EAuth: func\(token string, scopes \[\]string\) \(interface{}, error\) {`+any+

			`AAuth func\(string, string\) \(interface{}, error\)`+any+
			`BAuth func\(string\) \(interface{}, error\)`+any+
			`CAuth func\(string\) \(interface{}, error\)`+any+
			`DAuth func\(string\) \(interface{}, error\)`+any+
			`EAuth func\(string, \[\]string\) \(interface{}, error\)`+any+

			`if o\.AAuth == nil {`+any+
			`unregistered = append\(unregistered, "AAuth"\)`+any+
			`if o\.BAuth == nil {`+any+
			`unregistered = append\(unregistered, "K1Auth"\)`+any+
			`if o\.CAuth == nil {`+any+
			`unregistered = append\(unregistered, "K2Auth"\)`+any+
			`if o\.DAuth == nil {`+any+
			`unregistered = append\(unregistered, "K3Auth"\)`+any+
			`if o\.EAuth == nil {`+any+
			`unregistered = append\(unregistered, "EAuth"\)`+any+

			`case "A":`+any+
			`case "B":`+any+
			`case "C":`+any+
			`case "D":`+any+
			`case "E":`+any, res)
	}
}

func TestGenSecurityRequirements(t *testing.T) {
	for i := 0; i < 5; i++ {
		operation := "asecOp"
		b, err := opBuilder(operation, "../fixtures/bugs/1214/fixture-1214.yaml")
		require.NoError(t, err)

		b.Security = b.Analyzed.SecurityRequirementsFor(&b.Operation)
		genRequirements := b.makeSecurityRequirements("o")
		assert.Len(t, genRequirements, 2)
		assert.Equal(t, []GenSecurityRequirements{
			{
				GenSecurityRequirement{
					Name:   "A",
					Scopes: []string{},
				},
				GenSecurityRequirement{
					Name:   "B",
					Scopes: []string{},
				},
				GenSecurityRequirement{
					Name:   "E",
					Scopes: []string{"s0", "s1", "s2", "s3", "s4"},
				},
			},
			{
				GenSecurityRequirement{
					Name:   "C",
					Scopes: []string{},
				},
				GenSecurityRequirement{
					Name:   "D",
					Scopes: []string{},
				},
				GenSecurityRequirement{
					Name:   "E",
					Scopes: []string{"s5", "s6", "s7", "s8", "s9"},
				},
			},
		}, genRequirements)

		operation = "bsecOp"
		b, err = opBuilder(operation, "../fixtures/bugs/1214/fixture-1214.yaml")
		require.NoError(t, err)

		b.Security = b.Analyzed.SecurityRequirementsFor(&b.Operation)
		genRequirements = b.makeSecurityRequirements("o")
		assert.Len(t, genRequirements, 2)
		assert.Equal(t, []GenSecurityRequirements{
			{
				GenSecurityRequirement{
					Name:   "A",
					Scopes: []string{},
				},
				GenSecurityRequirement{
					Name:   "E",
					Scopes: []string{"s0", "s1", "s2", "s3", "s4"},
				},
			},
			{
				GenSecurityRequirement{
					Name:   "D",
					Scopes: []string{},
				},
				GenSecurityRequirement{
					Name:   "E",
					Scopes: []string{"s5", "s6", "s7", "s8", "s9"},
				},
			},
		}, genRequirements)
	}

	operation := "csecOp"
	b, err := opBuilder(operation, "../fixtures/bugs/1214/fixture-1214.yaml")
	require.NoError(t, err)

	b.Security = b.Analyzed.SecurityRequirementsFor(&b.Operation)
	genRequirements := b.makeSecurityRequirements("o")
	assert.NotNil(t, genRequirements)
	assert.Len(t, genRequirements, 0)

	operation = "nosecOp"
	b, err = opBuilder(operation, "../fixtures/bugs/1214/fixture-1214-2.yaml")
	require.NoError(t, err)

	b.Security = b.Analyzed.SecurityRequirementsFor(&b.Operation)
	genRequirements = b.makeSecurityRequirements("o")
	assert.Nil(t, genRequirements)
}

func TestGenerateServerOperation(t *testing.T) {
	defer discardOutput()()

	fname := "../fixtures/codegen/todolist.simple.yml"

	tgt, _ := ioutil.TempDir(filepath.Dir(fname), "generated")
	defer func() {
		_ = os.RemoveAll(tgt)
	}()
	o := &GenOpts{
		ValidateSpec:      false,
		IncludeModel:      true,
		IncludeHandler:    true,
		IncludeParameters: true,
		IncludeResponses:  true,
		ModelPackage:      "models",
		Spec:              fname,
		Target:            tgt,
	}
	require.NoError(t, o.EnsureDefaults())

	require.Error(t, GenerateServerOperation([]string{"createTask"}, nil))

	d := o.TemplateDir
	o.TemplateDir = "./nowhere"
	require.Error(t, GenerateServerOperation([]string{"notFound"}, o))

	o.TemplateDir = d
	d = o.Spec
	o.Spec = "nowhere.yaml"
	require.Error(t, GenerateServerOperation([]string{"notFound"}, o))

	o.Spec = d
	require.Error(t, GenerateServerOperation([]string{"notFound"}, o))

	require.NoError(t, GenerateServerOperation([]string{"createTask"}, o))

	// check expected files are generated and that's it
	_, err := os.Stat(filepath.Join(tgt, "tasks", "create_task.go"))
	assert.NoError(t, err)
	_, err = os.Stat(filepath.Join(tgt, "tasks", "create_task_parameters.go"))
	assert.NoError(t, err)
	_, err = os.Stat(filepath.Join(tgt, "tasks", "create_task_responses.go"))
	assert.NoError(t, err)

	origStdout := os.Stdout
	defer func() {
		os.Stdout = origStdout
	}()
	os.Stdout, _ = os.Create(filepath.Join(tgt, "stdout"))
	o.DumpData = true
	// just checks this does not fail
	err = GenerateServerOperation([]string{"createTask"}, o)
	assert.NoError(t, err)
	_, err = os.Stat(filepath.Join(tgt, "stdout"))
	assert.NoError(t, err)
}

// This tests that mimetypes generate stable code
func TestBuilder_Issue1646(t *testing.T) {
	defer discardOutput()()

	dr := testCwd(t)

	opts := &GenOpts{
		Spec:              filepath.FromSlash("../fixtures/bugs/1646/fixture-1646.yaml"),
		IncludeModel:      true,
		IncludeHandler:    true,
		IncludeParameters: true,
		IncludeResponses:  true,
		IncludeMain:       true,
		APIPackage:        "restapi",
		ModelPackage:      "model",
		ServerPackage:     "server",
		ClientPackage:     "client",
		Target:            dr,
		IsClient:          false,
	}
	err := opts.EnsureDefaults()
	require.NoError(t, err)
	appGen, err := newAppGenerator("fixture-1646", nil, nil, opts)
	require.NoError(t, err)

	preCons, preConj := appGen.makeConsumes()
	preProds, preProdj := appGen.makeProduces()
	assert.True(t, preConj)
	assert.True(t, preProdj)
	for i := 0; i < 5; i++ {
		cons, conj := appGen.makeConsumes()
		prods, prodj := appGen.makeProduces()
		assert.True(t, conj)
		assert.True(t, prodj)
		assert.Equal(t, preConj, conj)
		assert.Equal(t, preProdj, prodj)
		assert.Equal(t, preCons, cons)
		assert.Equal(t, preProds, prods)
	}
}

func TestGenServer_StrictAdditionalProperties(t *testing.T) {
	defer discardOutput()()

	dr := testCwd(t)

	opts := &GenOpts{
		Spec:              filepath.FromSlash("../fixtures/codegen/strict-additional-properties.yml"),
		IncludeModel:      true,
		IncludeHandler:    true,
		IncludeParameters: true,
		IncludeResponses:  true,
		IncludeMain:       true,
		APIPackage:        "restapi",
		ModelPackage:      "model",
		ServerPackage:     "server",
		ClientPackage:     "client",
		Target:            dr,
		IsClient:          false,
	}
	err := opts.EnsureDefaults()
	require.NoError(t, err)

	opts.StrictAdditionalProperties = true

	appGen, err := newAppGenerator("StrictAdditionalProperties", nil, nil, opts)
	require.NoError(t, err)

	op, err := appGen.makeCodegenApp()
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	err = templates.MustGet("serverOperation").Execute(buf, op.Operations[0])
	require.NoError(t, err)

	ff, err := appGen.GenOpts.LanguageOpts.FormatContent("strictAdditionalProperties.go", buf.Bytes())
	require.NoErrorf(t, err, buf.String())

	res := string(ff)
	for _, tt := range []struct {
		name      string
		assertion func(testing.TB, string, string) bool
	}{
		{"PostTestBody", assertInCode},
		{"PostTestParamsBodyExplicit", assertInCode},
		{"PostTestParamsBodyImplicit", assertInCode},
		{"PostTestParamsBodyDisabled", assertNotInCode},
	} {
		fn := funcBody(res, "*"+tt.name+") UnmarshalJSON(data []byte) error")
		require.NotEmpty(t, fn, "Method UnmarshalJSON should be defined for type *"+tt.name)
		tt.assertion(t, "dec.DisallowUnknownFields()", fn)
	}
}

func makeClientTimeoutNameTest() []struct {
	seenIds  map[string]interface{}
	name     string
	expected string
} {
	return []struct {
		seenIds  map[string]interface{}
		name     string
		expected string
	}{
		{
			seenIds:  nil,
			name:     "witness",
			expected: "witness",
		},
		{
			seenIds: map[string]interface{}{
				"id": true,
			},
			name:     "timeout",
			expected: "timeout",
		},
		{
			seenIds: map[string]interface{}{
				"timeout":        true,
				"requesttimeout": true,
			},
			name:     "timeout",
			expected: "httpRequestTimeout",
		},
		{
			seenIds: map[string]interface{}{
				"timeout":            true,
				"requesttimeout":     true,
				"httprequesttimeout": true,
				"swaggertimeout":     true,
				"operationtimeout":   true,
				"optimeout":          true,
			},
			name:     "timeout",
			expected: "operTimeout",
		},
		{
			seenIds: map[string]interface{}{
				"timeout":            true,
				"requesttimeout":     true,
				"httprequesttimeout": true,
				"swaggertimeout":     true,
				"operationtimeout":   true,
				"optimeout":          true,
				"opertimeout":        true,
				"opertimeout1":       true,
			},
			name:     "timeout",
			expected: "operTimeout11",
		},
	}
}

func TestRenameTimeout(t *testing.T) {
	for idx, toPin := range makeClientTimeoutNameTest() {
		i := idx
		testCase := toPin
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			assert.Equalf(t, testCase.expected, renameTimeout(testCase.seenIds, testCase.name), "unexpected deconflicting value [%d]", i)
		})
	}
}

func testInvalidParams() map[string]spec.Parameter {
	return map[string]spec.Parameter{
		"query#param1": *spec.QueryParam("param1"),
		"path#param1":  *spec.PathParam("param1"),
		"body#param1":  *spec.BodyParam("param1", &spec.Schema{}),
	}
}

func TestParamMappings(t *testing.T) {
	// Test deconfliction of duplicate param names across param locations
	mappings, _ := paramMappings(testInvalidParams())
	require.Contains(t, mappings, "query")
	require.Contains(t, mappings, "path")
	require.Contains(t, mappings, "body")
	q := mappings["query"]
	p := mappings["path"]
	b := mappings["body"]
	require.Len(t, q, 1)
	require.Len(t, p, 1)
	require.Len(t, b, 1)
	require.Containsf(t, q, "param1", "unexpected content of %#v", q)
	require.Containsf(t, p, "param1", "unexpected content of %#v", p)
	require.Containsf(t, b, "param1", "unexpected content of %#v", b)
	assert.Equalf(t, "QueryParam1", q["param1"], "unexpected content of %#v", q["param1"])
	assert.Equalf(t, "PathParam1", p["param1"], "unexpected content of %#v", p["param1"])
	assert.Equalf(t, "BodyParam1", b["param1"], "unexpected content of %#v", b["param1"])
}

func TestDeconflictTag(t *testing.T) {
	assert.Equal(t, "runtimeops", deconflictTag(nil, "runtime"))
	assert.Equal(t, "apiops", deconflictTag([]string{"tag1"}, "api"))
	assert.Equal(t, "apiops1", deconflictTag([]string{"tag1", "apiops"}, "api"))
	assert.Equal(t, "tlsops", deconflictTag([]string{"tag1"}, "tls"))
	assert.Equal(t, "mytag", deconflictTag([]string{"tag1", "apiops"}, "mytag"))

	assert.Equal(t, "operationsops", renameOperationPackage([]string{"tag1"}, "operations"))
	assert.Equal(t, "operationsops11", renameOperationPackage([]string{"tag1", "operationsops1", "operationsops"}, "operations"))
}

func TestGenServer_2161_panic(t *testing.T) {
	t.Parallel()
	defer discardOutput()()

	generated, err := ioutil.TempDir(testCwd(t), "generated_2161")
	require.NoError(t, err)

	defer func() {
		_ = os.RemoveAll(generated)
	}()

	opts := &GenOpts{
		Spec:                       filepath.FromSlash("../fixtures/bugs/2161/fixture-2161-panic.json"),
		IncludeModel:               true,
		IncludeHandler:             true,
		IncludeParameters:          true,
		IncludeResponses:           true,
		IncludeMain:                true,
		APIPackage:                 "restapi",
		ModelPackage:               "model",
		ServerPackage:              "server",
		ClientPackage:              "client",
		Target:                     generated,
		IsClient:                   false,
		StrictAdditionalProperties: true,
	}
	require.NoError(t, opts.EnsureDefaults())

	appGen, err := newAppGenerator("inlinedSubtype", nil, nil, opts)
	require.NoError(t, err)

	op, err := appGen.makeCodegenApp()
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	var selectedOp int
	for i := range op.Operations {
		if op.Operations[i].Name == "configuration_update_configuration_module" {
			selectedOp = i
		}
	}
	require.NotEmpty(t, selectedOp, "dev error: invalid test vs fixture")

	require.NoError(t, opts.templates.MustGet("serverOperation").Execute(buf, op.Operations[selectedOp]))

	_, err = appGen.GenOpts.LanguageOpts.FormatContent(op.Operations[selectedOp].Name+".go", buf.Bytes())
	require.NoErrorf(t, err, buf.String())
	// NOTE(fred): I know that the generated model is wrong from this spec at the moment.
	// The test with this fix simply asserts that there is no panic / internal error with building this.
}

func TestGenServer_1659_Principal(t *testing.T) {
	defer discardOutput()()

	dr := testCwd(t)

	for _, toPin := range []struct {
		Title       string
		Opts        *GenOpts
		Expected    map[string][]string
		NotExpected map[string][]string
	}{
		{
			Title: "default",
			Opts: &GenOpts{
				Spec:              filepath.FromSlash("../fixtures/enhancements/1659/fixture-1659.yaml"),
				IncludeHandler:    true,
				IncludeParameters: true,
				IncludeResponses:  true,
				IncludeMain:       false,
				APIPackage:        "restapi",
				ModelPackage:      "models",
				ServerPackage:     "server",
				ClientPackage:     "client",
				Target:            dr,
				IsClient:          false,
			},
			Expected: map[string][]string{
				"configure": {
					`if api.ApikeyAuth == nil {`,
					`api.ApikeyAuth = func(token string) (interface{}, error) {`,
					`if api.BasicAuth == nil {`,
					`api.BasicAuth = func(user string, pass string) (interface{}, error) {`,
					`if api.PetstoreAuthAuth == nil {`,
					`api.PetstoreAuthAuth = func(token string, scopes []string) (interface{}, error) {`,
					`api.GetHandler =restapi.GetHandlerFunc(func(params restapi.GetParams, principal interface{}) middleware.Responder {`,
					`api.PostHandler =restapi.PostHandlerFunc(func(params restapi.PostParams) middleware.Responder {`,
				},
				"get": {
					`type GetHandlerFunc func(GetParams, interface{})  middleware.Responder`,
					`func (fn GetHandlerFunc) Handle(params GetParams, principal interface{})  middleware.Responder  {`,
					`return fn(params, principal)`,
					`type GetHandler interface {`,
					`Handle(GetParams, interface{})  middleware.Responder`,
					`uprinc, aCtx, err := o.Context.Authorize(r, route)`,
					`if uprinc != nil {`,
					`principal = uprinc.(interface{})`,
					`res := o.Handler.Handle(Params, principal)`,
				},
				"post": {
					`type PostHandlerFunc func(PostParams)  middleware.Responder`,
					`return fn(params)`,
					`type PostHandler interface {`,
					`Handle(PostParams)  middleware.Responder`,
					`res := o.Handler.Handle(Params)`,
				},
			},
			NotExpected: map[string][]string{
				"post": {
					`uprinc, aCtx, err := o.Context.Authorize(r, route)`,
					`principal = uprinc.(interface{})`,
				},
			},
		},
		{
			Title: "principal is struct",
			Opts: &GenOpts{
				Spec:              filepath.FromSlash("../fixtures/enhancements/1659/fixture-1659.yaml"),
				IncludeHandler:    true,
				IncludeParameters: true,
				IncludeResponses:  true,
				IncludeMain:       false,
				APIPackage:        "restapi",
				ModelPackage:      "models",
				ServerPackage:     "server",
				ClientPackage:     "client",
				Target:            dr,
				Principal:         "github.com/example/security.Principal",
				IsClient:          false,
			},
			Expected: map[string][]string{
				"configure": {
					`auth "github.com/example/security"`,
					`if api.ApikeyAuth == nil {`,
					`api.ApikeyAuth = func(token string) (*auth.Principal, error) {`,
					`if api.BasicAuth == nil {`,
					`api.BasicAuth = func(user string, pass string) (*auth.Principal, error) {`,
					`if api.PetstoreAuthAuth == nil {`,
					`api.PetstoreAuthAuth = func(token string, scopes []string) (*auth.Principal, error) {`,
					`api.GetHandler =restapi.GetHandlerFunc(func(params restapi.GetParams, principal *auth.Principal) middleware.Responder {`,
					`api.PostHandler =restapi.PostHandlerFunc(func(params restapi.PostParams) middleware.Responder {`,
				},
				"get": {
					`type GetHandlerFunc func(GetParams, *auth.Principal)  middleware.Responder`,
					`func (fn GetHandlerFunc) Handle(params GetParams, principal *auth.Principal)  middleware.Responder  {`,
					`return fn(params, principal)`,
					`type GetHandler interface {`,
					`Handle(GetParams, *auth.Principal)  middleware.Responder`,
					`uprinc, aCtx, err := o.Context.Authorize(r, route)`,
					`if uprinc != nil {`,
					`principal = uprinc.(*auth.Principal)`,
					`res := o.Handler.Handle(Params, principal)`,
				},
				"post": {
					`type PostHandlerFunc func(PostParams)  middleware.Responder`,
					`return fn(params)`,
					`type PostHandler interface {`,
					`Handle(PostParams)  middleware.Responder`,
					`res := o.Handler.Handle(Params)`,
				},
			},
			NotExpected: map[string][]string{
				"post": {
					`uprinc, aCtx, err := o.Context.Authorize(r, route)`,
					`principal = uprinc.(`,
				},
			},
		},
		{
			Title: "principal is interface",
			Opts: &GenOpts{
				Spec:                 filepath.FromSlash("../fixtures/enhancements/1659/fixture-1659.yaml"),
				IncludeHandler:       true,
				IncludeParameters:    true,
				IncludeResponses:     true,
				IncludeMain:          false,
				APIPackage:           "restapi",
				ModelPackage:         "models",
				ServerPackage:        "server",
				ClientPackage:        "client",
				Target:               dr,
				Principal:            "github.com/example/security.PrincipalIface",
				IsClient:             false,
				PrincipalCustomIface: true,
			},
			Expected: map[string][]string{
				"configure": {
					`auth "github.com/example/security"`,
					`if api.ApikeyAuth == nil {`,
					`api.ApikeyAuth = func(token string) (auth.PrincipalIface, error) {`,
					`if api.BasicAuth == nil {`,
					`api.BasicAuth = func(user string, pass string) (auth.PrincipalIface, error) {`,
					`if api.PetstoreAuthAuth == nil {`,
					`api.PetstoreAuthAuth = func(token string, scopes []string) (auth.PrincipalIface, error) {`,
					`api.GetHandler =restapi.GetHandlerFunc(func(params restapi.GetParams, principal auth.PrincipalIface) middleware.Responder {`,
					`api.PostHandler =restapi.PostHandlerFunc(func(params restapi.PostParams) middleware.Responder {`,
				},
				"get": {
					`type GetHandlerFunc func(GetParams, auth.PrincipalIface)  middleware.Responder`,
					`func (fn GetHandlerFunc) Handle(params GetParams, principal auth.PrincipalIface)  middleware.Responder  {`,
					`return fn(params, principal)`,
					`type GetHandler interface {`,
					`Handle(GetParams, auth.PrincipalIface)  middleware.Responder`,
					`uprinc, aCtx, err := o.Context.Authorize(r, route)`,
					`if uprinc != nil {`,
					`principal = uprinc.(auth.PrincipalIface)`,
					`res := o.Handler.Handle(Params, principal)`,
				},
				"post": {
					`type PostHandlerFunc func(PostParams)  middleware.Responder`,
					`return fn(params)`,
					`type PostHandler interface {`,
					`Handle(PostParams)  middleware.Responder`,
					`res := o.Handler.Handle(Params)`,
				},
			},
			NotExpected: map[string][]string{
				"post": {
					`uprinc, aCtx, err := o.Context.Authorize(r, route)`,
					`principal = uprinc.(`,
				},
			},
		},
		{
			Title: "stratoscale: principal is struct",
			Opts: &GenOpts{
				Spec:              filepath.FromSlash("../fixtures/enhancements/1659/fixture-1659.yaml"),
				IncludeHandler:    true,
				IncludeParameters: true,
				IncludeResponses:  true,
				IncludeMain:       false,
				APIPackage:        "restapi",
				ModelPackage:      "models",
				ServerPackage:     "server",
				ClientPackage:     "client",
				Target:            dr,
				Principal:         "github.com/example/security.Principal",
				IsClient:          false,
				Template:          "stratoscale",
			},
			Expected: map[string][]string{
				"configure": {
					`auth "github.com/example/security"`,
					`AuthApikey func(token string) (*auth.Principal, error)`,
					`AuthBasic func(user string, pass string) (*auth.Principal, error)`,
					`AuthPetstoreAuth func(token string, scopes []string) (*auth.Principal, error)`,
					`api.ApikeyAuth = func(token string) (*auth.Principal, error) {`,
					`if c.AuthApikey == nil {`,
					`panic("you specified a custom principal type, but did not provide the authenticator to provide this")`,
					`return c.AuthApikey(token)`,
					`api.BasicAuth = func(user string, pass string) (*auth.Principal, error) {`,
					`if c.AuthBasic == nil {`,
					`panic("you specified a custom principal type, but did not provide the authenticator to provide this")`,
					`return c.AuthBasic(user, pass)`,
					`api.PetstoreAuthAuth = func(token string, scopes []string) (*auth.Principal, error) {`,
					`if c.AuthPetstoreAuth == nil {`,
					`panic("you specified a custom principal type, but did not provide the authenticator to provide this")`,
					`return c.AuthPetstoreAuth(token, scopes)`,
					`api.APIAuthorizer = authorizer(c.Authorizer)`,
					`api.GetHandler =restapi.GetHandlerFunc(func(params restapi.GetParams, principal *auth.Principal) middleware.Responder {`,
					`ctx = storeAuth(ctx, principal)`,
					`return c.RestapiAPI.Get(ctx, params)`,
					`api.PostHandler =restapi.PostHandlerFunc(func(params restapi.PostParams) middleware.Responder {`,
					`return c.RestapiAPI.Post(ctx, params)`,
					`func (a authorizer) Authorize(req *http.Request, principal interface{}) error {`,
					`ctx := storeAuth(req.Context(), principal)`,
					`func storeAuth(ctx context.Context, principal interface{})`,
				},
			},
		},
		{
			Title: "stratoscale: principal is interface",
			Opts: &GenOpts{
				Spec:                 filepath.FromSlash("../fixtures/enhancements/1659/fixture-1659.yaml"),
				IncludeHandler:       true,
				IncludeParameters:    true,
				IncludeResponses:     true,
				IncludeMain:          false,
				APIPackage:           "restapi",
				ModelPackage:         "models",
				ServerPackage:        "server",
				ClientPackage:        "client",
				Target:               dr,
				Principal:            "github.com/example/security.PrincipalIface",
				IsClient:             false,
				PrincipalCustomIface: true,
				Template:             "stratoscale",
			},
			Expected: map[string][]string{
				"configure": {
					`auth "github.com/example/security"`,
					`AuthApikey func(token string) (auth.PrincipalIface, error)`,
					`AuthBasic func(user string, pass string) (auth.PrincipalIface, error)`,
					`AuthPetstoreAuth func(token string, scopes []string) (auth.PrincipalIface, error)`,
					`api.ApikeyAuth = func(token string) (auth.PrincipalIface, error) {`,
					`if c.AuthApikey == nil {`,
					`panic("you specified a custom principal type, but did not provide the authenticator to provide this")`,
					`return c.AuthApikey(token)`,
					`api.BasicAuth = func(user string, pass string) (auth.PrincipalIface, error) {`,
					`if c.AuthBasic == nil {`,
					`panic("you specified a custom principal type, but did not provide the authenticator to provide this")`,
					`return c.AuthBasic(user, pass)`,
					`api.PetstoreAuthAuth = func(token string, scopes []string) (auth.PrincipalIface, error) {`,
					`if c.AuthPetstoreAuth == nil {`,
					`panic("you specified a custom principal type, but did not provide the authenticator to provide this")`,
					`return c.AuthPetstoreAuth(token, scopes)`,
					`api.APIAuthorizer = authorizer(c.Authorizer)`,
					`api.GetHandler =restapi.GetHandlerFunc(func(params restapi.GetParams, principal auth.PrincipalIface) middleware.Responder {`,
					`ctx = storeAuth(ctx, principal)`,
					`return c.RestapiAPI.Get(ctx, params)`,
					`api.PostHandler =restapi.PostHandlerFunc(func(params restapi.PostParams) middleware.Responder {`,
					`return c.RestapiAPI.Post(ctx, params)`,
					`func (a authorizer) Authorize(req *http.Request, principal interface{}) error {`,
					`ctx := storeAuth(req.Context(), principal)`,
					`func storeAuth(ctx context.Context, principal interface{})`,
				},
			},
		},
	} {
		fixture := toPin
		t.Run(fixture.Title, func(t *testing.T) {
			t.Parallel()

			opts := fixture.Opts
			require.NoError(t, opts.EnsureDefaults())
			require.NoError(t, opts.setTemplates())

			appGen, err := newAppGenerator(fixture.Title, nil, nil, opts)
			require.NoError(t, err)

			op, err := appGen.makeCodegenApp()
			require.NoError(t, err)

			bufC := bytes.NewBuffer(nil)
			require.NoError(t, opts.templates.MustGet("serverConfigureapi").Execute(bufC, op))

			_, err = appGen.GenOpts.LanguageOpts.FormatContent("configure_api.go", bufC.Bytes())
			require.NoErrorf(t, err, bufC.String())

			for _, line := range fixture.Expected["configure"] {
				assertInCode(t, line, bufC.String())
			}
			for _, line := range fixture.NotExpected["configure"] {
				assertNotInCode(t, line, bufC.String())
			}

			for i := range op.Operations {
				bufO := bytes.NewBuffer(nil)
				require.NoError(t, opts.templates.MustGet("serverOperation").Execute(bufO, op.Operations[i]))

				_, erf := appGen.GenOpts.LanguageOpts.FormatContent(op.Operations[i].Name+".go", bufO.Bytes())
				require.NoErrorf(t, erf, bufO.String())

				for _, line := range fixture.Expected[op.Operations[i].Name] {
					assertInCode(t, line, bufO.String())
				}
				for _, line := range fixture.NotExpected[op.Operations[i].Name] {
					assertNotInCode(t, line, bufO.String())
				}
			}
		})
	}
}
