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
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	basicFixture = "../fixtures/petstores/petstore.json"
)

func testClientGenOpts() *GenOpts {
	g := &GenOpts{}
	g.Target = "."
	g.APIPackage = defaultAPIPackage
	g.ModelPackage = defaultModelPackage
	g.ServerPackage = defaultServerPackage
	g.ClientPackage = defaultClientPackage
	g.Principal = ""
	g.IncludeModel = true
	g.IncludeHandler = true
	g.IncludeParameters = true
	g.IncludeResponses = true
	g.IncludeSupport = true
	g.TemplateDir = ""
	g.DumpData = false
	g.IsClient = true
	if err := g.EnsureDefaults(); err != nil {
		panic(err)
	}
	return g
}

func Test_GenerateClient(t *testing.T) {
	t.Parallel()
	defer discardOutput()()

	const clientName = "test"

	t.Run("exercise codegen safeguards", func(t *testing.T) {
		t.Run("should fail on nil options", func(t *testing.T) {
			require.Error(t,
				GenerateClient(clientName, []string{"model1"}, []string{"op1", "op2"}, nil),
			)
		})

		t.Run("should fail on invalid templates location (1)", func(t *testing.T) {
			opts := testClientGenOpts()
			opts.TemplateDir = "dir/nowhere"
			require.Error(t,
				GenerateClient(clientName, []string{"model1"}, []string{"op1", "op2"}, opts),
			)
		})

		t.Run("should fail on invalid templates location (2)", func(t *testing.T) {
			opts := testClientGenOpts()
			opts.TemplateDir = "http://nowhere.com"
			require.Error(t,
				GenerateClient(clientName, []string{"model1"}, []string{"op1", "op2"}, opts),
			)
		})

		t.Run("should fail on invalid spec location", func(t *testing.T) {
			opts := testClientGenOpts()
			opts.Spec = "dir/nowhere.yaml"
			require.Error(t,
				GenerateClient(clientName, []string{"model1"}, []string{"op1", "op2"}, opts),
			)
		})

		t.Run("should fail on invalid model name", func(t *testing.T) {
			opts := testClientGenOpts()
			opts.Spec = basicFixture
			require.Error(t,
				GenerateClient(clientName, []string{"model1"}, []string{}, opts),
			)
		})

		t.Run("should fail on bad content in spec (HTML, not json)", func(t *testing.T) {
			opts := testClientGenOpts()
			opts.Spec = "https://github.com/OAI/OpenAPI-Specification/blob/master/examples/v2.0/json/petstore.json"
			require.Error(t,
				GenerateClient(clientName, []string{}, []string{}, opts),
			)
		})

		t.Run("should fail when no valid operation is selected", func(t *testing.T) {
			opts := testClientGenOpts()
			opts.Spec = "https://raw.githubusercontent.com/OAI/OpenAPI-Specification/master/examples/v2.0/yaml/petstore.yaml"
			require.Error(t,
				GenerateClient(clientName, []string{}, []string{"wrongOperationID"}, opts),
			)
		})

		t.Run("should refuse to generate from garbled parameters", func(t *testing.T) {
			opts := testClientGenOpts()
			opts.Spec = filepath.Join("..", "fixtures", "bugs", "2527", "swagger.yml")
			opts.ValidateSpec = false
			err := GenerateClient(clientName, []string{}, []string{"GetDeposits"}, opts)
			require.Error(t, err)
			require.ErrorContains(t, err, `GET /deposits, "" has an invalid parameter definition`)
		})
	})

	t.Run("should generate client", func(t *testing.T) {
		cwd, err := os.Getwd()
		require.NoError(t, err)

		t.Run("from remote spec", func(t *testing.T) {
			opts := testClientGenOpts()
			opts.Spec = "https://raw.githubusercontent.com/OAI/OpenAPI-Specification/master/examples/v2.0/yaml/petstore.yaml"

			tft, err := os.MkdirTemp(cwd, "generated")
			require.NoError(t, err)
			t.Cleanup(func() {
				_ = os.RemoveAll(tft)
			})

			opts.Target = tft
			opts.IsClient = true
			DefaultSectionOpts(opts)

			t.Cleanup(func() {
				_ = os.RemoveAll(opts.Target)
			})
			require.NoError(t,
				GenerateClient(clientName, []string{}, []string{}, opts),
			)
		})

		t.Run("from fixed spec (issue #2527)", func(t *testing.T) {
			opts := testClientGenOpts()
			opts.Spec = filepath.Join("..", "fixtures", "bugs", "2527", "swagger-fixed.yml")

			tft, err := os.MkdirTemp(cwd, "generated")
			require.NoError(t, err)
			t.Cleanup(func() {
				_ = os.RemoveAll(tft)
			})

			opts.Target = tft
			opts.IsClient = true
			DefaultSectionOpts(opts)

			t.Cleanup(func() {
				_ = os.RemoveAll(opts.Target)
			})
			require.NoError(t,
				GenerateClient(clientName, []string{}, []string{}, opts),
			)
		})

		t.Run("should dump template data", func(t *testing.T) {
			opts := testClientGenOpts()
			opts.Spec = "https://raw.githubusercontent.com/OAI/OpenAPI-Specification/master/examples/v2.0/yaml/petstore.yaml"

			origStdout := os.Stdout
			defer func() {
				os.Stdout = origStdout
			}()
			tgt, err := os.MkdirTemp(cwd, "dumped")
			require.NoError(t, err)
			t.Cleanup(func() {
				_ = os.RemoveAll(tgt)
			})
			os.Stdout, err = os.Create(filepath.Join(tgt, "stdout"))
			require.NoError(t, err)

			opts.DumpData = true
			require.NoError(t,
				GenerateClient(clientName, []string{}, []string{}, opts),
			)
			t.Run("make sure this did not fail and we have some output", func(t *testing.T) {
				stat, err := os.Stat(filepath.Join(tgt, "stdout"))
				require.NoError(t, err)
				require.Positive(t, stat.Size())
			})
		})
	})
}

func assertImports(t testing.TB, baseImport, code string) {
	assertRegexpInCode(t, baseImport, code)
	assertRegexpInCode(t, `"`+baseImport+`/abc_linux"`, code)
	assertRegexpInCode(t, `"`+baseImport+`/abc_linux"`, code)
	assertRegexpInCode(t, `"`+baseImport+`/abc_test"`, code)
	assertRegexpInCode(t, `apiops\s+"`+baseImport+`/api"`, code)
	assertRegexpInCode(t, `"`+baseImport+`/custom"`, code)
	assertRegexpInCode(t, `"`+baseImport+`/hash_tag_donuts"`, code)
	assertRegexpInCode(t, `"`+baseImport+`/nr123abc"`, code)
	assertRegexpInCode(t, `"`+baseImport+`/nr_at_donuts"`, code)
	assertRegexpInCode(t, `"`+baseImport+`/plus_donuts`, code)
	assertRegexpInCode(t, `strfmtops "`+baseImport+`/strfmt`, code)
	assertRegexpInCode(t, `"`+baseImport+`/forced`, code)
	assertRegexpInCode(t, `"`+baseImport+`/nr12nasty`, code)
	assertRegexpInCode(t, `"`+baseImport+`/override`, code)
	assertRegexpInCode(t, `"`+baseImport+`/gtl`, code)
	assertRegexpInCode(t, `"`+baseImport+`/operationsops`, code)
}

func TestClient(t *testing.T) {
	t.Parallel()
	defer discardOutput()()

	base := os.Getenv("GOPATH")
	var importBase string
	if base == "" {
		base = "."
		importBase = "github.com/go-swagger/go-swagger/generator/"
	} else {
		base = filepath.Join(base, "src")
		err := os.MkdirAll(base, 0o755)
		require.NoError(t, err)
	}
	targetdir, err := os.MkdirTemp(base, "swagger_nogo")
	require.NoError(t, err, "Failed to create a test target directory: %v", err)

	t.Cleanup(func() {
		_ = os.RemoveAll(targetdir)
	})

	tests := []struct {
		name      string
		spec      string
		template  string
		wantError bool
		prepare   func(opts *GenOpts)
		verify    func(testing.TB, string)
	}{
		{
			name:      "InvalidSpec",
			wantError: true,
			prepare: func(opts *GenOpts) {
				opts.Spec = invalidSpecExample
				opts.ValidateSpec = true
			},
		},
		{
			name: "BaseImportDisabled",
			prepare: func(opts *GenOpts) {
				opts.LanguageOpts.BaseImportFunc = nil
			},
			wantError: false,
		},
		{
			name:      "Non_existing_contributor_template",
			template:  "NonExistingContributorTemplate",
			wantError: true,
		},
		{
			name:      "Existing_contributor",
			template:  "stratoscale",
			wantError: false,
		},
		{
			name:      "packages mangling",
			wantError: false,
			spec:      filepath.Join("..", "fixtures", "bugs", "2111", "fixture-2111.yaml"),
			verify: func(t testing.TB, target string) {
				require.True(t, fileExists(target, "client"))

				// assert package generation based on mangled tags
				target = filepath.Join(target, "client")
				assert.True(t, fileExists(target, "abc_linux"))
				assert.True(t, fileExists(target, "abc_test"))
				assert.True(t, fileExists(target, "api"))
				assert.True(t, fileExists(target, "custom"))
				assert.True(t, fileExists(target, "hash_tag_donuts"))
				assert.True(t, fileExists(target, "nr123abc"))
				assert.True(t, fileExists(target, "nr_at_donuts"))
				assert.True(t, fileExists(target, "operations"))
				assert.True(t, fileExists(target, "plus_donuts"))
				assert.True(t, fileExists(target, "strfmt"))
				assert.True(t, fileExists(target, "forced"))
				assert.True(t, fileExists(target, "gtl"))
				assert.True(t, fileExists(target, "nr12nasty"))
				assert.True(t, fileExists(target, "override"))
				assert.True(t, fileExists(target, "operationsops"))

				buf, err := os.ReadFile(filepath.Join(target, "foo_client.go"))
				require.NoError(t, err)

				// assert client import, with deconfliction
				code := string(buf)
				importRegexp := importBase + `swagger_nogo\d+/packages_mangling/client`
				assertImports(t, importRegexp, code)

				assertInCode(t, `cli.Strfmt = strfmtops.New(transport, formats)`, code)
				assertInCode(t, `cli.API = apiops.New(transport, formats)`, code)
				assertInCode(t, `cli.Operations = operations.New(transport, formats)`, code)
			},
		},
		{
			name:      "packages flattening",
			wantError: false,
			spec:      filepath.Join("..", "fixtures", "bugs", "2111", "fixture-2111.yaml"),
			prepare: func(opts *GenOpts) {
				opts.SkipTagPackages = true
			},
			verify: func(t testing.TB, target string) {
				require.True(t, fileExists(target, "client"))

				// packages are not created here
				target = filepath.Join(target, "client")
				assert.False(t, fileExists(target, "abc_linux"))
				assert.False(t, fileExists(target, "abc_test"))
				assert.False(t, fileExists(target, "api"))
				assert.False(t, fileExists(target, "custom"))
				assert.False(t, fileExists(target, "hash_tag_donuts"))
				assert.False(t, fileExists(target, "nr123abc"))
				assert.False(t, fileExists(target, "nr_at_donuts"))
				assert.False(t, fileExists(target, "plus_donuts"))
				assert.False(t, fileExists(target, "strfmt"))
				assert.False(t, fileExists(target, "forced"))
				assert.False(t, fileExists(target, "gtl"))
				assert.False(t, fileExists(target, "nr12nasty"))
				assert.False(t, fileExists(target, "override"))
				assert.False(t, fileExists(target, "operationsops"))

				assert.True(t, fileExists(target, "operations"))
			},
		},
		{
			name:      "name with trailing API",
			spec:      filepath.Join("..", "fixtures", "bugs", "2278", "fixture-2278.yaml"),
			wantError: false,
		},
	}

	t.Run("generate client", func(t *testing.T) {
		for idx, toPin := range tests {
			tt := toPin
			i := idx
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				opts := testClientGenOpts()
				opts.Spec = basicFixture
				opts.Target = filepath.Join(targetdir, opts.LanguageOpts.ManglePackageName(tt.name, "client_test"+strconv.Itoa(i)))
				err := os.MkdirAll(opts.Target, 0o755)
				require.NoError(t, err)

				if tt.spec == "" {
					opts.Spec = basicFixture
				} else {
					opts.Spec = tt.spec
				}
				opts.Template = tt.template

				if tt.prepare != nil {
					tt.prepare(opts)
				}

				err = GenerateClient("foo", nil, nil, opts)
				if tt.wantError {
					require.Errorf(t, err, "expected an error for client build fixture: %s", opts.Spec)
				} else {
					require.NoError(t, err, "unexpected error for client build fixture: %s", opts.Spec)
				}

				if tt.verify != nil {
					tt.verify(t, opts.Target)
				}
			})
		}
	})
}

func TestGenClient_1518(t *testing.T) {
	t.Parallel()
	defer discardOutput()()

	// test client response handling when unexpected success response kicks in

	opts := testClientGenOpts()
	opts.Spec = filepath.Join("..", "fixtures", "bugs", "1518", "fixture-1518.yaml")

	cwd, _ := os.Getwd()
	tft, _ := os.MkdirTemp(cwd, "generated")
	opts.Target = tft

	defer func() {
		_ = os.RemoveAll(tft)
	}()

	err := GenerateClient("client", []string{}, []string{}, opts)
	require.NoError(t, err)

	fixtureConfig := map[string][]string{
		"client/operations/operations_client.go": { // generated file
			// expected code lines
			`success, ok := result.(*GetRecords1OK)`,
			`if ok {`,
			`return success, nil`,
			`msg := fmt.Sprintf(`,
			`panic(msg)`,
			// expected code lines
			`success, ok := result.(*GetRecords2OK)`,
			`if ok {`,
			`return success, nil`,
			`unexpectedSuccess := result.(*GetRecords2Default)`,
			`return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())`,
			// expected code lines
			`switch value := result.(type) {`,
			`case *GetRecords3OK:`,
			`return value, nil, nil`,
			`case *GetRecords3Created:`,
			`return nil, value, nil`,
			`msg := fmt.Sprintf(`,
			`panic(msg)`,
			// expected code lines
			`switch value := result.(type) {`,
			`case *GetRecords4OK:`,
			`return value, nil, nil`,
			`case *GetRecords4Created:`,
			`return nil, value, nil`,
			`unexpectedSuccess := result.(*GetRecords4Default)`,
			`return nil, nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())`,
		},
	}

	for fileToInspect, expectedCode := range fixtureConfig {
		code, err := os.ReadFile(filepath.Join(opts.Target, filepath.FromSlash(fileToInspect)))
		require.NoError(t, err)

		for line, codeLine := range expectedCode {
			if !assertInCode(t, strings.TrimSpace(codeLine), string(code)) {
				t.Logf("Code expected did not match in codegenfile %s for expected line %d: %q", fileToInspect, line, expectedCode[line])
			}
		}
	}
}

func TestGenClient_2945(t *testing.T) {
	t.Parallel()
	defer discardOutput()()

	opts := testClientGenOpts()
	opts.Spec = filepath.Join("..", "fixtures", "bugs", "2945", "fixture-2945.yaml")

	cwd, _ := os.Getwd()
	tft, _ := os.MkdirTemp(cwd, "generated")
	opts.Target = tft

	defer func() {
		_ = os.RemoveAll(tft)
	}()

	err := GenerateClient("client", []string{}, []string{}, opts)
	require.NoError(t, err)

	fixtureConfig := map[string][]string{
		"client/operations/get_version_responses.go": { // generated file
			// expected code lines
			`return nil, runtime.NewAPIError("[GET /version] getVersion", response, response.Code())`,
		},
	}

	for fileToInspect, expectedCode := range fixtureConfig {
		code, err := os.ReadFile(filepath.Join(opts.Target, filepath.FromSlash(fileToInspect)))
		require.NoError(t, err)

		for line, codeLine := range expectedCode {
			if !assertInCode(t, strings.TrimSpace(codeLine), string(code)) {
				t.Logf("Code expected did not match in codegenfile %s for expected line %d: %q", fileToInspect, line, expectedCode[line])
			}
		}
	}
}

func TestGenClient_2471(t *testing.T) {
	t.Parallel()
	defer discardOutput()()

	opts := testClientGenOpts()
	opts.Spec = filepath.Join("..", "fixtures", "bugs", "2471", "fixture-2471.yaml")

	cwd, _ := os.Getwd()
	tft, _ := os.MkdirTemp(cwd, "generated")
	opts.Target = tft

	defer func() {
		_ = os.RemoveAll(tft)
	}()

	err := GenerateClient("client", []string{}, []string{}, opts)
	require.NoError(t, err)

	fixtureConfig := map[string][]string{
		"client/operations/example_post_parameters.go": { // generated file
			`func (o *ExamplePostParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {`,
			`	if err := r.SetTimeout(o.timeout); err != nil {`,
			`	joinedFoo := o.bindParamFoo(reg)`,
			`	if len(joinedFoo) > 0 {`,
			`		if err := r.SetHeaderParam("Foo", joinedFoo[0]); err != nil {`,
			` joinedFooPath := o.bindParamFooPath(reg)`,
			`	if len(joinedFooPath) > 0 {`,
			`		if err := r.SetPathParam("FooPath", joinedFooPath[0]); err != nil {`,
			` joinedFooQuery := o.bindParamFooQuery(reg)`,
			`	if err := r.SetQueryParam("FooQuery", joinedFooQuery...); err != nil {`,
			`func (o *ExamplePostParams) bindParamFoo(formats strfmt.Registry) []string {`,
			`		fooIR := o.Foo`,
			`   var fooIC []string`,
			` 	for _, fooIIR := range fooIR {`,
			` 	  fooIIV := fooIIR`,
			` 	  fooIC = append(fooIC, fooIIV)`,
			` 	  fooIS := swag.JoinByFormat(fooIC, "")`,
			` 	  return fooIS`,
			`func (o *ExamplePostParams) bindParamFooPath(formats strfmt.Registry) []string {`,
			` 		fooPathIR := o.FooPath`,
			` 	 	var fooPathIC []string`,
			` 	 	for _, fooPathIIR := range fooPathIR {`,
			` 	 		fooPathIIV := fooPathIIR`,
			` 	    fooPathIC = append(fooPathIC, fooPathIIV)`,
			` 	    fooPathIS := swag.JoinByFormat(fooPathIC, "")`,
			`  return fooPathIS`,
			`func (o *ExamplePostParams) bindParamFooQuery(formats strfmt.Registry) []string {`,
			` 	  fooQueryIR := o.FooQuery`,
			` 	  var fooQueryIC []string`,
			` 	  for _, fooQueryIIR := range fooQueryIR {`,
			` 	    fooQueryIIV := fooQueryIIR`,
			` 	    fooQueryIC = append(fooQueryIC, fooQueryIIV)`,
			` 	    fooQueryIS := swag.JoinByFormat(fooQueryIC, "")`,
			`  return fooQueryIS`,
		},
	}

	for fileToInspect, expectedCode := range fixtureConfig {
		code, err := os.ReadFile(filepath.Join(opts.Target, filepath.FromSlash(fileToInspect)))
		require.NoError(t, err)

		for line, codeLine := range expectedCode {
			if !assertInCode(t, strings.TrimSpace(codeLine), string(code)) {
				t.Logf("Code expected did not match in codegenfile %s for expected line %d: %q", fileToInspect, line, expectedCode[line])
			}
		}
	}
}

func TestGenClient_2096(t *testing.T) {
	t.Parallel()
	defer discardOutput()()

	opts := testClientGenOpts()
	opts.Spec = filepath.Join("..", "fixtures", "bugs", "2096", "fixture-2096.yaml")

	cwd, _ := os.Getwd()
	tft, _ := os.MkdirTemp(cwd, "generated")
	opts.Target = tft

	defer func() {
		_ = os.RemoveAll(tft)
	}()

	err := GenerateClient("client", []string{}, []string{}, opts)
	require.NoError(t, err)

	fixtureConfig := map[string][]string{
		"client/operations/list_resources_parameters.go": { // generated file
			`type ListResourcesParams struct {`,
			`	Fields []string`,
			`func (o *ListResourcesParams) SetDefaults() {`,
			`	var (`,
			`		fieldsDefault = []string{"first", "second", "third"}`,
			`	val := ListResourcesParams{`,
			`		Fields: fieldsDefault,`,
			`	val.timeout = o.timeout`,
			`	val.Context = o.Context`,
			`	val.HTTPClient = o.HTTPClient`,
			`	*o = val`,
			`	joinedFields := o.bindParamFields(reg)`,
			`	if err := r.SetQueryParam("fields", joinedFields...); err != nil {`,
		},
	}

	for fileToInspect, expectedCode := range fixtureConfig {
		code, err := os.ReadFile(filepath.Join(opts.Target, filepath.FromSlash(fileToInspect)))
		require.NoError(t, err)

		for line, codeLine := range expectedCode {
			if !assertInCode(t, strings.TrimSpace(codeLine), string(code)) {
				t.Logf("Code expected did not match in codegenfile %s for expected line %d: %q", fileToInspect, line, expectedCode[line])
			}
		}
	}
}

func TestGenClient_909_3(t *testing.T) {
	t.Parallel()
	defer discardOutput()()

	opts := testClientGenOpts()
	opts.Spec = filepath.Join("..", "fixtures", "bugs", "909", "fixture-909-3.yaml")

	cwd, _ := os.Getwd()
	tft, _ := os.MkdirTemp(cwd, "generated")
	opts.Target = tft

	defer func() {
		_ = os.RemoveAll(tft)
	}()

	err := GenerateClient("client", []string{}, []string{}, opts)
	require.NoError(t, err)

	fixtureConfig := map[string][]string{
		"client/operations/get_optional_parameters.go": { // generated file
			`func (o *GetOptionalParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {`,
			`	if err := r.SetTimeout(o.timeout); err != nil {`,
			`	if o.IsAnOption2 != nil {`,
			`		joinedIsAnOption2 := o.bindParamIsAnOption2(reg)`,
			`		if err := r.SetQueryParam("isAnOption2", joinedIsAnOption2...); err != nil {`,
			`	if o.IsAnOption4 != nil {`,
			`		joinedIsAnOption4 := o.bindParamIsAnOption4(reg)`,
			`		if err := r.SetQueryParam("isAnOption4", joinedIsAnOption4...); err != nil {`,
			`	if o.IsAnOptionalHeader != nil {`,
			`		joinedIsAnOptionalHeader := o.bindParamIsAnOptionalHeader(reg)`,
			`		if len(joinedIsAnOptionalHeader) > 0 {`,
			`			if err := r.SetHeaderParam("isAnOptionalHeader", joinedIsAnOptionalHeader[0]); err != nil {`,
			`	if o.NotAnOption1 != nil {`,
			`		joinedNotAnOption1 := o.bindParamNotAnOption1(reg)`,
			`		if err := r.SetQueryParam("notAnOption1", joinedNotAnOption1...); err != nil {`,
			`	if o.NotAnOption3 != nil {`,
			`		if err := r.SetBodyParam(o.NotAnOption3); err != nil {`,
			`func (o *GetOptionalParams) bindParamIsAnOption2(formats strfmt.Registry) []string {`,
			`	isAnOption2IR := o.IsAnOption2`,
			`	var isAnOption2IC []string`,
			`	for _, isAnOption2IIR := range isAnOption2IR { // explode [][]strfmt.UUID`,
			`		var isAnOption2IIC []string`,
			`		for _, isAnOption2IIIR := range isAnOption2IIR { // explode []strfmt.UUID`,
			`			isAnOption2IIIV := isAnOption2IIIR.String() // strfmt.UUID as string`,
			`			isAnOption2IIC = append(isAnOption2IIC, isAnOption2IIIV)`,
			`		isAnOption2IIS := swag.JoinByFormat(isAnOption2IIC, "")`,
			`		isAnOption2IIV := isAnOption2IIS[0]`,
			`		isAnOption2IC = append(isAnOption2IC, isAnOption2IIV)`,
			`	isAnOption2IS := swag.JoinByFormat(isAnOption2IC, "pipes")`,
			`	return isAnOption2IS`,
			`func (o *GetOptionalParams) bindParamIsAnOption4(formats strfmt.Registry) []string {`,
			`	isAnOption4IR := o.IsAnOption4`,
			`	var isAnOption4IC []string`,
			`	for _, isAnOption4IIR := range isAnOption4IR { // explode [][][]strfmt.UUID`,
			`		var isAnOption4IIC []string`,
			`		for _, isAnOption4IIIR := range isAnOption4IIR { // explode [][]strfmt.UUID`,
			`			var isAnOption4IIIC []string`,
			`			for _, isAnOption4IIIIR := range isAnOption4IIIR { // explode []strfmt.UUID`,
			`				isAnOption4IIIIV := isAnOption4IIIIR.String() // strfmt.UUID as string`,
			`				isAnOption4IIIC = append(isAnOption4IIIC, isAnOption4IIIIV)`,
			`			isAnOption4IIIS := swag.JoinByFormat(isAnOption4IIIC, "pipes")`,
			`			isAnOption4IIIV := isAnOption4IIIS[0]`,
			`			isAnOption4IIC = append(isAnOption4IIC, isAnOption4IIIV)`,
			`		}`,
			`		isAnOption4IIS := swag.JoinByFormat(isAnOption4IIC, "tsv")`,
			`		isAnOption4IIV := isAnOption4IIS[0]`,
			`		isAnOption4IC = append(isAnOption4IC, isAnOption4IIV)`,
			`	isAnOption4IS := swag.JoinByFormat(isAnOption4IC, "csv")`,
			`	return isAnOption4IS`,
			`func (o *GetOptionalParams) bindParamIsAnOptionalHeader(formats strfmt.Registry) []string {`,
			`	isAnOptionalHeaderIR := o.IsAnOptionalHeader`,
			`	var isAnOptionalHeaderIC []string`,
			`	for _, isAnOptionalHeaderIIR := range isAnOptionalHeaderIR { // explode [][]strfmt.UUID`,
			`		var isAnOptionalHeaderIIC []string`,
			`		for _, isAnOptionalHeaderIIIR := range isAnOptionalHeaderIIR { // explode []strfmt.UUID`,
			`			isAnOptionalHeaderIIIV := isAnOptionalHeaderIIIR.String() // strfmt.UUID as string`,
			`			isAnOptionalHeaderIIC = append(isAnOptionalHeaderIIC, isAnOptionalHeaderIIIV)`,
			`		isAnOptionalHeaderIIS := swag.JoinByFormat(isAnOptionalHeaderIIC, "")`,
			`		isAnOptionalHeaderIIV := isAnOptionalHeaderIIS[0]`,
			`		isAnOptionalHeaderIC = append(isAnOptionalHeaderIC, isAnOptionalHeaderIIV)`,
			`	isAnOptionalHeaderIS := swag.JoinByFormat(isAnOptionalHeaderIC, "pipes")`,
			`	return isAnOptionalHeaderIS`,
			`func (o *GetOptionalParams) bindParamNotAnOption1(formats strfmt.Registry) []string {`,
			`	notAnOption1IR := o.NotAnOption1`,
			`	var notAnOption1IC []string`,
			`	for _, notAnOption1IIR := range notAnOption1IR { // explode [][]strfmt.DateTime`,
			`		var notAnOption1IIC []string`,
			`		for _, notAnOption1IIIR := range notAnOption1IIR { // explode []strfmt.DateTime`,
			`			notAnOption1IIIV := notAnOption1IIIR.String() // strfmt.DateTime as string`,
			`			notAnOption1IIC = append(notAnOption1IIC, notAnOption1IIIV)`,
			`		notAnOption1IIS := swag.JoinByFormat(notAnOption1IIC, "pipes")`,
			`		notAnOption1IIV := notAnOption1IIS[0]`,
			`		notAnOption1IC = append(notAnOption1IC, notAnOption1IIV)`,
			`	notAnOption1IS := swag.JoinByFormat(notAnOption1IC, "csv")`,
			`	return notAnOption1IS`,
		},
	}

	for fileToInspect, expectedCode := range fixtureConfig {
		code, err := os.ReadFile(filepath.Join(opts.Target, filepath.FromSlash(fileToInspect)))
		require.NoError(t, err)

		for line, codeLine := range expectedCode {
			if !assertInCode(t, strings.TrimSpace(codeLine), string(code)) {
				t.Logf("Code expected did not match in codegenfile %s for expected line %d: %q", fileToInspect, line, expectedCode[line])
			}
		}
	}
}

func TestGenClient_909_5(t *testing.T) {
	t.Parallel()
	defer discardOutput()()

	opts := testClientGenOpts()
	opts.Spec = filepath.Join("..", "fixtures", "bugs", "909", "fixture-909-5.yaml")

	cwd, _ := os.Getwd()
	tft, _ := os.MkdirTemp(cwd, "generated")
	opts.Target = tft

	defer func() {
		_ = os.RemoveAll(tft)
	}()

	err := GenerateClient("client", []string{}, []string{}, opts)
	require.NoError(t, err)

	fixtureConfig := map[string][]string{
		"client/operations/get_optional_responses.go": { // generated file
			`func NewGetOptionalOK() *GetOptionalOK {`,
			`	var (`,
			`		xIsAnOptionalHeader0Default = strfmt.DateTime{}`,
			`		xIsAnOptionalHeader0PrimitiveDefault = float32(345.55)`,
			`		xIsAnOptionalHeader0StringerDefault = strfmt.UUID("524fc6d5-66c6-46f6-90bc-34e0d0139e43")`,
			`		xIsAnOptionalHeader1Default = make([]strfmt.DateTime, 0, 50)`,
			`		xIsAnOptionalHeader2Default = make([][]int32, 0, 50)`,
			`		xIsAnOptionalHeader2NoFormatDefault = make([][]int64, 0, 50)`,
			`		xIsAnOptionalHeader3Default = make([][][]strfmt.UUID, 0, 50)`,
			`	)`,
			`	if err := xIsAnOptionalHeader0Default.UnmarshalText([]byte("2018-01-28T23:54:00.000Z")); err != nil {`,
			`		msg := fmt.Sprintf("invalid default value for xIsAnOptionalHeader0: %v", err)`,
			"	if err := json.Unmarshal([]byte(`[\"2018-01-28T23:54:00.000Z\",\"2018-02-28T23:54:00.000Z\",\"2018-03-28T23:54:00.000Z\",\"2018-04-28T23:54:00.000Z\"]`), &xIsAnOptionalHeader1Default); err != nil {",
			`		msg := fmt.Sprintf("invalid default value for xIsAnOptionalHeader1: %v", err)`,
			"	if err := json.Unmarshal([]byte(`[[21,22,23],[31,32,33]]`), &xIsAnOptionalHeader2Default); err != nil {",
			`		msg := fmt.Sprintf("invalid default value for xIsAnOptionalHeader2: %v", err)`,
			"	if err := json.Unmarshal([]byte(`[[21,22,23],[31,32,33]]`), &xIsAnOptionalHeader2NoFormatDefault); err != nil {",
			`		msg := fmt.Sprintf("invalid default value for xIsAnOptionalHeader2NoFormat: %v", err)`,
			"	if err := json.Unmarshal([]byte(`[[[\"524fc6d5-66c6-46f6-90bc-34e0d0139e43\",\"c8199a5f-f7ce-4fb1-b8af-082256125e89\"],[\"c8199a5f-f7ce-4fb1-b8af-082256125e89\",\"524fc6d5-66c6-46f6-90bc-34e0d0139e43\"]],[[\"c8199a5f-f7ce-4fb1-b8af-082256125e89\",\"524fc6d5-66c6-46f6-90bc-34e0d0139e43\"],[\"524fc6d5-66c6-46f6-90bc-34e0d0139e43\",\"c8199a5f-f7ce-4fb1-b8af-082256125e89\"]]]`), &xIsAnOptionalHeader3Default); err != nil {",
			`		msg := fmt.Sprintf("invalid default value for xIsAnOptionalHeader3: %v", err)`,
			`	return &GetOptionalOK{`,
			`		XIsAnOptionalHeader0:          xIsAnOptionalHeader0Default,`,
			`		XIsAnOptionalHeader0Primitive: xIsAnOptionalHeader0PrimitiveDefault,`,
			`		XIsAnOptionalHeader0Stringer:  xIsAnOptionalHeader0StringerDefault,`,
			`		XIsAnOptionalHeader1:          xIsAnOptionalHeader1Default,`,
			`		XIsAnOptionalHeader2:          xIsAnOptionalHeader2Default,`,
			`		XIsAnOptionalHeader2NoFormat:  xIsAnOptionalHeader2NoFormatDefault,`,
			`		XIsAnOptionalHeader3:          xIsAnOptionalHeader3Default,`,
			`func (o *GetOptionalOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {`,
			`	hdrXIsAnOptionalHeader0 := response.GetHeader("x-isAnOptionalHeader0")`,
			`	if hdrXIsAnOptionalHeader0 != "" {`,
			`		valxIsAnOptionalHeader0, err := formats.Parse("date-time", hdrXIsAnOptionalHeader0)`,
			`		if err != nil {`,
			`			return errors.InvalidType("x-isAnOptionalHeader0", "header", "strfmt.DateTime", hdrXIsAnOptionalHeader0)`,
			`		o.XIsAnOptionalHeader0 = *(valxIsAnOptionalHeader0.(*strfmt.DateTime))`,
			`	hdrXIsAnOptionalHeader0DirtSimple := response.GetHeader("x-isAnOptionalHeader0DirtSimple")`,
			`	if hdrXIsAnOptionalHeader0DirtSimple != "" {`,
			`		o.XIsAnOptionalHeader0DirtSimple = hdrXIsAnOptionalHeader0DirtSimple`,
			`	hdrXIsAnOptionalHeader0DirtSimpleArray := response.GetHeader("x-isAnOptionalHeader0DirtSimpleArray")`,
			`	if hdrXIsAnOptionalHeader0DirtSimpleArray != "" {`,
			`		valXIsAnOptionalHeader0DirtSimpleArray, err := o.bindHeaderXIsAnOptionalHeader0DirtSimpleArray(hdrXIsAnOptionalHeader0DirtSimpleArray, formats)`,
			`		o.XIsAnOptionalHeader0DirtSimpleArray = valXIsAnOptionalHeader0DirtSimpleArray`,
			`	hdrXIsAnOptionalHeader0DirtSimpleInteger := response.GetHeader("x-isAnOptionalHeader0DirtSimpleInteger")`,
			`	if hdrXIsAnOptionalHeader0DirtSimpleInteger != "" {`,
			`		valxIsAnOptionalHeader0DirtSimpleInteger, err := swag.ConvertInt64(hdrXIsAnOptionalHeader0DirtSimpleInteger)`,
			`			return errors.InvalidType("x-isAnOptionalHeader0DirtSimpleInteger", "header", "int64", hdrXIsAnOptionalHeader0DirtSimpleInteger)`,
			`		o.XIsAnOptionalHeader0DirtSimpleInteger = valxIsAnOptionalHeader0DirtSimpleInteger`,
			`	hdrXIsAnOptionalHeader0Primitive := response.GetHeader("x-isAnOptionalHeader0Primitive")`,
			`	if hdrXIsAnOptionalHeader0Primitive != "" {`,
			`		valxIsAnOptionalHeader0Primitive, err := swag.ConvertFloat32(hdrXIsAnOptionalHeader0Primitive)`,
			`			return errors.InvalidType("x-isAnOptionalHeader0Primitive", "header", "float32", hdrXIsAnOptionalHeader0Primitive)`,
			`		o.XIsAnOptionalHeader0Primitive = valxIsAnOptionalHeader0Primitive`,
			`	hdrXIsAnOptionalHeader0Stringer := response.GetHeader("x-isAnOptionalHeader0Stringer")`,
			`	if hdrXIsAnOptionalHeader0Stringer != "" {`,
			`		valxIsAnOptionalHeader0Stringer, err := formats.Parse("uuid", hdrXIsAnOptionalHeader0Stringer)`,
			`			return errors.InvalidType("x-isAnOptionalHeader0Stringer", "header", "strfmt.UUID", hdrXIsAnOptionalHeader0Stringer)`,
			`		o.XIsAnOptionalHeader0Stringer = *(valxIsAnOptionalHeader0Stringer.(*strfmt.UUID))`,
			`	hdrXIsAnOptionalHeader1 := response.GetHeader("x-isAnOptionalHeader1")`,
			`	if hdrXIsAnOptionalHeader1 != "" {`,
			`		valXIsAnOptionalHeader1, err := o.bindHeaderXIsAnOptionalHeader1(hdrXIsAnOptionalHeader1, formats)`,
			`		o.XIsAnOptionalHeader1 = valXIsAnOptionalHeader1`,
			`	hdrXIsAnOptionalHeader2 := response.GetHeader("x-isAnOptionalHeader2")`,
			`	if hdrXIsAnOptionalHeader2 != "" {`,
			`		valXIsAnOptionalHeader2, err := o.bindHeaderXIsAnOptionalHeader2(hdrXIsAnOptionalHeader2, formats)`,
			`		o.XIsAnOptionalHeader2 = valXIsAnOptionalHeader2`,
			`	hdrXIsAnOptionalHeader2NoFormat := response.GetHeader("x-isAnOptionalHeader2NoFormat")`,
			`	if hdrXIsAnOptionalHeader2NoFormat != "" {`,
			`		valXIsAnOptionalHeader2NoFormat, err := o.bindHeaderXIsAnOptionalHeader2NoFormat(hdrXIsAnOptionalHeader2NoFormat, formats)`,
			`		o.XIsAnOptionalHeader2NoFormat = valXIsAnOptionalHeader2NoFormat`,
			`	hdrXIsAnOptionalHeader3 := response.GetHeader("x-isAnOptionalHeader3")`,
			`	if hdrXIsAnOptionalHeader3 != "" {`,
			`		valXIsAnOptionalHeader3, err := o.bindHeaderXIsAnOptionalHeader3(hdrXIsAnOptionalHeader3, formats)`,
			`		o.XIsAnOptionalHeader3 = valXIsAnOptionalHeader3`,
			`	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {`,
			`func (o *GetOptionalOK) bindHeaderXIsAnOptionalHeader0DirtSimpleArray(hdr string, formats strfmt.Registry) ([]int64, error) {`,
			`	xIsAnOptionalHeader0DirtSimpleArrayIV := hdr`,
			`	var (`,
			`		xIsAnOptionalHeader0DirtSimpleArrayIC []int64`,
			`	xIsAnOptionalHeader0DirtSimpleArrayIR := swag.SplitByFormat(xIsAnOptionalHeader0DirtSimpleArrayIV, "")`,
			`	for i, xIsAnOptionalHeader0DirtSimpleArrayIIV := range xIsAnOptionalHeader0DirtSimpleArrayIR {`,
			`		val, err := swag.ConvertInt64(xIsAnOptionalHeader0DirtSimpleArrayIIV)`,
			`			return nil, errors.InvalidType(fmt.Sprintf("%s.%v", "header", i), "headeritems.", "int64", xIsAnOptionalHeader0DirtSimpleArrayIIV)`,
			`		xIsAnOptionalHeader0DirtSimpleArrayIIC := val`,
			`		xIsAnOptionalHeader0DirtSimpleArrayIC = append(xIsAnOptionalHeader0DirtSimpleArrayIC, xIsAnOptionalHeader0DirtSimpleArrayIIC) // roll-up int64 into []int64`,
			`	return xIsAnOptionalHeader0DirtSimpleArrayIC, nil`,
			`func (o *GetOptionalOK) bindHeaderXIsAnOptionalHeader1(hdr string, formats strfmt.Registry) ([]strfmt.DateTime, error) {`,
			`	xIsAnOptionalHeader1IV := hdr`,
			`		xIsAnOptionalHeader1IC []strfmt.DateTime`,
			`	xIsAnOptionalHeader1IR := swag.SplitByFormat(xIsAnOptionalHeader1IV, "tsv")`,
			`	for i, xIsAnOptionalHeader1IIV := range xIsAnOptionalHeader1IR {`,
			`		val, err := formats.Parse("date-time", xIsAnOptionalHeader1IIV)`,
			`			return nil, errors.InvalidType(fmt.Sprintf("%s.%v", "header", i), "headeritems.", "strfmt.DateTime", xIsAnOptionalHeader1IIV)`,
			`		xIsAnOptionalHeader1IIC := val.(strfmt.DateTime)`,
			`		xIsAnOptionalHeader1IC = append(xIsAnOptionalHeader1IC, xIsAnOptionalHeader1IIC) // roll-up strfmt.DateTime into []strfmt.DateTime`,
			`	return xIsAnOptionalHeader1IC, nil`,
			`func (o *GetOptionalOK) bindHeaderXIsAnOptionalHeader2(hdr string, formats strfmt.Registry) ([][]int32, error) {`,
			`	xIsAnOptionalHeader2IV := hdr`,
			`		xIsAnOptionalHeader2IC [][]int32`,
			`	xIsAnOptionalHeader2IR := swag.SplitByFormat(xIsAnOptionalHeader2IV, "")`,
			`	for _, xIsAnOptionalHeader2IIV := range xIsAnOptionalHeader2IR {`,
			`			xIsAnOptionalHeader2IIC []int32`,
			`		xIsAnOptionalHeader2IIR := swag.SplitByFormat(xIsAnOptionalHeader2IIV, "pipes")`,
			`		for ii, xIsAnOptionalHeader2IIIV := range xIsAnOptionalHeader2IIR {`,
			`			val, err := swag.ConvertInt32(xIsAnOptionalHeader2IIIV)`,
			`				return nil, errors.InvalidType(fmt.Sprintf("%s.%v", "header", ii), "headeritems.items.", "int32", xIsAnOptionalHeader2IIIV)`,
			`			xIsAnOptionalHeader2IIIC := val`,
			`			xIsAnOptionalHeader2IIC = append(xIsAnOptionalHeader2IIC, xIsAnOptionalHeader2IIIC) // roll-up int32 into []int32`,
			`		xIsAnOptionalHeader2IC = append(xIsAnOptionalHeader2IC, xIsAnOptionalHeader2IIC) // roll-up []int32 into [][]int32`,
			`	return xIsAnOptionalHeader2IC, nil`,
			`func (o *GetOptionalOK) bindHeaderXIsAnOptionalHeader2NoFormat(hdr string, formats strfmt.Registry) ([][]int64, error) {`,
			`	xIsAnOptionalHeader2NoFormatIV := hdr`,
			`		xIsAnOptionalHeader2NoFormatIC [][]int64`,
			`	xIsAnOptionalHeader2NoFormatIR := swag.SplitByFormat(xIsAnOptionalHeader2NoFormatIV, "pipes")`,
			`	for _, xIsAnOptionalHeader2NoFormatIIV := range xIsAnOptionalHeader2NoFormatIR {`,
			`			xIsAnOptionalHeader2NoFormatIIC []int64`,
			`		xIsAnOptionalHeader2NoFormatIIR := swag.SplitByFormat(xIsAnOptionalHeader2NoFormatIIV, "tsv")`,
			`		for ii, xIsAnOptionalHeader2NoFormatIIIV := range xIsAnOptionalHeader2NoFormatIIR {`,
			`			val, err := swag.ConvertInt64(xIsAnOptionalHeader2NoFormatIIIV)`,
			`				return nil, errors.InvalidType(fmt.Sprintf("%s.%v", "header", ii), "headeritems.items.", "int64", xIsAnOptionalHeader2NoFormatIIIV)`,
			`			xIsAnOptionalHeader2NoFormatIIIC := val`,
			`			xIsAnOptionalHeader2NoFormatIIC = append(xIsAnOptionalHeader2NoFormatIIC, xIsAnOptionalHeader2NoFormatIIIC) // roll-up int64 into []int64`,
			`		xIsAnOptionalHeader2NoFormatIC = append(xIsAnOptionalHeader2NoFormatIC, xIsAnOptionalHeader2NoFormatIIC) // roll-up []int64 into [][]int64`,
			`	return xIsAnOptionalHeader2NoFormatIC, nil`,
			`func (o *GetOptionalOK) bindHeaderXIsAnOptionalHeader3(hdr string, formats strfmt.Registry) ([][][]strfmt.UUID, error) {`,
			`	xIsAnOptionalHeader3IV := hdr`,
			`		xIsAnOptionalHeader3IC [][][]strfmt.UUID`,
			`	xIsAnOptionalHeader3IR := swag.SplitByFormat(xIsAnOptionalHeader3IV, "")`,
			`	for _, xIsAnOptionalHeader3IIV := range xIsAnOptionalHeader3IR {`,
			`			xIsAnOptionalHeader3IIC [][]strfmt.UUID`,
			`		xIsAnOptionalHeader3IIR := swag.SplitByFormat(xIsAnOptionalHeader3IIV, "pipes")`,
			`		for _, xIsAnOptionalHeader3IIIV := range xIsAnOptionalHeader3IIR {`,
			`				xIsAnOptionalHeader3IIIC []strfmt.UUID`,
			`			xIsAnOptionalHeader3IIIR := swag.SplitByFormat(xIsAnOptionalHeader3IIIV, "")`,
			`			for iii, xIsAnOptionalHeader3IIIIV := range xIsAnOptionalHeader3IIIR {`,
			`				val, err := formats.Parse("uuid", xIsAnOptionalHeader3IIIIV)`,
			`					return nil, errors.InvalidType(fmt.Sprintf("%s.%v", "header", iii), "headeritems.items.", "strfmt.UUID", xIsAnOptionalHeader3IIIIV)`,
			`				xIsAnOptionalHeader3IIIIC := val.(strfmt.UUID)`,
			`				xIsAnOptionalHeader3IIIC = append(xIsAnOptionalHeader3IIIC, xIsAnOptionalHeader3IIIIC) // roll-up strfmt.UUID into []strfmt.UUID`,
			`			xIsAnOptionalHeader3IIC = append(xIsAnOptionalHeader3IIC, xIsAnOptionalHeader3IIIC) // roll-up []strfmt.UUID into [][]strfmt.UUID`,
			`		xIsAnOptionalHeader3IC = append(xIsAnOptionalHeader3IC, xIsAnOptionalHeader3IIC) // roll-up [][]strfmt.UUID into [][][]strfmt.UUID`,
			`	return xIsAnOptionalHeader3IC, nil`,
		},
	}

	for fileToInspect, expectedCode := range fixtureConfig {
		code, err := os.ReadFile(filepath.Join(opts.Target, filepath.FromSlash(fileToInspect)))
		require.NoError(t, err)

		for line, codeLine := range expectedCode {
			if !assertInCode(t, strings.TrimSpace(codeLine), string(code)) {
				t.Logf("Code expected did not match in codegenfile %s for expected line %d: %q", fileToInspect, line, expectedCode[line])
			}
		}
	}
}

func TestGenClient_909_6(t *testing.T) {
	t.Parallel()
	defer discardOutput()()

	opts := testClientGenOpts()
	opts.Spec = filepath.Join("..", "fixtures", "bugs", "909", "fixture-909-6.yaml")

	cwd, _ := os.Getwd()
	tft, _ := os.MkdirTemp(cwd, "generated")
	opts.Target = tft

	defer func() {
		_ = os.RemoveAll(tft)
	}()

	err := GenerateClient("client", []string{}, []string{}, opts)
	require.NoError(t, err)

	fixtureConfig := map[string][]string{
		"client/operations/get_optional_responses.go": { // generated file
			`func NewGetOptionalOK() *GetOptionalOK {`,
			`	var (`,
			`		xaBoolDefault = bool(true)`,
			`		xaBsonObjectIDDefault = strfmt.ObjectId{}`,
			`		xaByteDefault = strfmt.Base64([]byte(nil))`,
			`		xaCreditCardDefault = strfmt.CreditCard("4111-1111-1111-1111")`,
			`		xaDateDefault = strfmt.Date{}`,
			`		xaDateTimeDefault = strfmt.DateTime{}`,
			`		xaDoubleDefault = float64(99.99)`,
			`		xaDurationDefault = strfmt.Duration(0)`,
			`		xaFloatDefault = float32(99.99)`,
			`		xaHexColorDefault = strfmt.HexColor("#FFFFFF")`,
			`		xaHostnameDefault = strfmt.Hostname("www.example.com")`,
			`		xaInt32Default = int32(-99)`,
			`		xaInt64Default = int64(-99)`,
			`		xaMacDefault = strfmt.MAC("01:02:03:04:05:06")`,
			`		xaPasswordDefault = strfmt.Password("secret")`,
			`		xaRGBColorDefault = strfmt.RGBColor("rgb(255,255,255)")`,
			`		xaSsnDefault = strfmt.SSN("111-11-1111")`,
			`		xaUUIDDefault = strfmt.UUID("a8098c1a-f86e-11da-bd1a-00112444be1e")`,
			`		xaUUID3Default = strfmt.UUID3("bcd02e22-68f0-3046-a512-327cca9def8f")`,
			`		xaUUID4Default = strfmt.UUID4("025b0d74-00a2-4048-bf57-227c5111bb34")`,
			`		xaUUID5Default = strfmt.UUID5("886313e1-3b8a-5372-9b90-0c9aee199e5d")`,
			`		xaUint32Default = uint32(99)`,
			`		xaUint64Default = uint64(99)`,
			`		xaURIDefault = strfmt.URI("http://foo.bar/?baz=qux#quux")`,
			`		xAnEmailDefault = strfmt.Email("fredbi@github.com")`,
			`		xAnISBNDefault = strfmt.ISBN("0321751043")`,
			`		xAnISBN10Default = strfmt.ISBN10("0321751043")`,
			`		xAnISBN13Default = strfmt.ISBN13("978 3401013190")`,
			`		xAnIPV4Default = strfmt.IPv4("192.168.224.1")`,
			`		xAnIPV6Default = strfmt.IPv6("::1")`,
			`	if err := xaBsonObjectIDDefault.UnmarshalText([]byte("507f1f77bcf86cd799439011")); err != nil {`,
			`		msg := fmt.Sprintf("invalid default value for xaBsonObjectID: %v", err)`,
			`	if err := xaByteDefault.UnmarshalText([]byte("ZWxpemFiZXRocG9zZXk=")); err != nil {`,
			`		msg := fmt.Sprintf("invalid default value for xaByte: %v", err)`,
			`	if err := xaDateDefault.UnmarshalText([]byte("1970-01-01")); err != nil {`,
			`		msg := fmt.Sprintf("invalid default value for xaDate: %v", err)`,
			`	if err := xaDateTimeDefault.UnmarshalText([]byte("1970-01-01T11:01:05.283185Z")); err != nil {`,
			`		msg := fmt.Sprintf("invalid default value for xaDateTime: %v", err)`,
			`	if err := xaDurationDefault.UnmarshalText([]byte("1 ms")); err != nil {`,
			`		msg := fmt.Sprintf("invalid default value for xaDuration: %v", err)`,
			`	return &GetOptionalOK{`,
			`		XaBool:         xaBoolDefault,`,
			`		XaBsonObjectID: xaBsonObjectIDDefault,`,
			`		XaByte:         xaByteDefault,`,
			`		XaCreditCard:   xaCreditCardDefault,`,
			`		XaDate:         xaDateDefault,`,
			`		XaDateTime:     xaDateTimeDefault,`,
			`		XaDouble:       xaDoubleDefault,`,
			`		XaDuration:     xaDurationDefault,`,
			`		XaFloat:        xaFloatDefault,`,
			`		XaHexColor:     xaHexColorDefault,`,
			`		XaHostname:     xaHostnameDefault,`,
			`		XaInt32:        xaInt32Default,`,
			`		XaInt64:        xaInt64Default,`,
			`		XaMac:          xaMacDefault,`,
			`		XaPassword:     xaPasswordDefault,`,
			`		XaRGBColor:     xaRGBColorDefault,`,
			`		XaSsn:          xaSsnDefault,`,
			`		XaUUID:         xaUUIDDefault,`,
			`		XaUUID3:        xaUUID3Default,`,
			`		XaUUID4:        xaUUID4Default,`,
			`		XaUUID5:        xaUUID5Default,`,
			`		XaUint32:       xaUint32Default,`,
			`		XaUint64:       xaUint64Default,`,
			`		XaURI:          xaURIDefault,`,
			`		XAnEmail:       xAnEmailDefault,`,
			`		XAnISBN:        xAnISBNDefault,`,
			`		XAnISBN10:      xAnISBN10Default,`,
			`		XAnISBN13:      xAnISBN13Default,`,
			`		XAnIPV4:        xAnIPV4Default,`,
			`		XAnIPV6:        xAnIPV6Default,`,
			`func (o *GetOptionalOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {`,
			`	hdrXaBool := response.GetHeader("X-aBool")`,
			`	if hdrXaBool != "" {`,
			`		valxABool, err := swag.ConvertBool(hdrXaBool)`,
			`		if err != nil {`,
			`			return errors.InvalidType("X-aBool", "header", "bool", hdrXaBool)`,
			`		o.XaBool = valxABool`,
			`	hdrXaBsonObjectID := response.GetHeader("X-aBsonObjectId")`,
			`	if hdrXaBsonObjectID != "" {`,
			`		valxABsonObjectId, err := formats.Parse("bsonobjectid", hdrXaBsonObjectID)`,
			`		if err != nil {`,
			`			return errors.InvalidType("X-aBsonObjectId", "header", "strfmt.ObjectId", hdrXaBsonObjectID)`,
			`		o.XaBsonObjectID = *(valxABsonObjectId.(*strfmt.ObjectId))`,
			`	hdrXaByte := response.GetHeader("X-aByte")`,
			`	if hdrXaByte != "" {`,
			`		valxAByte, err := formats.Parse("byte", hdrXaByte)`,
			`		if err != nil {`,
			`			return errors.InvalidType("X-aByte", "header", "strfmt.Base64", hdrXaByte)`,
			`		o.XaByte = *(valxAByte.(*strfmt.Base64))`,
			`	hdrXaCreditCard := response.GetHeader("X-aCreditCard")`,
			`	if hdrXaCreditCard != "" {`,
			`		valxACreditCard, err := formats.Parse("creditcard", hdrXaCreditCard)`,
			`		if err != nil {`,
			`			return errors.InvalidType("X-aCreditCard", "header", "strfmt.CreditCard", hdrXaCreditCard)`,
			`		o.XaCreditCard = *(valxACreditCard.(*strfmt.CreditCard))`,
			`	hdrXaDate := response.GetHeader("X-aDate")`,
			`	if hdrXaDate != "" {`,
			`		valxADate, err := formats.Parse("date", hdrXaDate)`,
			`		if err != nil {`,
			`			return errors.InvalidType("X-aDate", "header", "strfmt.Date", hdrXaDate)`,
			`		o.XaDate = *(valxADate.(*strfmt.Date))`,
			`	hdrXaDateTime := response.GetHeader("X-aDateTime")`,
			`	if hdrXaDateTime != "" {`,
			`		valxADateTime, err := formats.Parse("date-time", hdrXaDateTime)`,
			`		if err != nil {`,
			`			return errors.InvalidType("X-aDateTime", "header", "strfmt.DateTime", hdrXaDateTime)`,
			`		o.XaDateTime = *(valxADateTime.(*strfmt.DateTime))`,
			`	hdrXaDouble := response.GetHeader("X-aDouble")`,
			`	if hdrXaDouble != "" {`,
			`		valxADouble, err := swag.ConvertFloat64(hdrXaDouble)`,
			`		if err != nil {`,
			`			return errors.InvalidType("X-aDouble", "header", "float64", hdrXaDouble)`,
			`		o.XaDouble = valxADouble`,
			`	hdrXaDuration := response.GetHeader("X-aDuration")`,
			`	if hdrXaDuration != "" {`,
			`		valxADuration, err := formats.Parse("duration", hdrXaDuration)`,
			`		if err != nil {`,
			`			return errors.InvalidType("X-aDuration", "header", "strfmt.Duration", hdrXaDuration)`,
			`		o.XaDuration = *(valxADuration.(*strfmt.Duration))`,
			`	hdrXaFloat := response.GetHeader("X-aFloat")`,
			`	if hdrXaFloat != "" {`,
			`		valxAFloat, err := swag.ConvertFloat32(hdrXaFloat)`,
			`		if err != nil {`,
			`			return errors.InvalidType("X-aFloat", "header", "float32", hdrXaFloat)`,
			`		o.XaFloat = valxAFloat`,
			`	hdrXaHexColor := response.GetHeader("X-aHexColor")`,
			`	if hdrXaHexColor != "" {`,
			`		valxAHexColor, err := formats.Parse("hexcolor", hdrXaHexColor)`,
			`		if err != nil {`,
			`			return errors.InvalidType("X-aHexColor", "header", "strfmt.HexColor", hdrXaHexColor)`,
			`		o.XaHexColor = *(valxAHexColor.(*strfmt.HexColor))`,
			`	hdrXaHostname := response.GetHeader("X-aHostname")`,
			`	if hdrXaHostname != "" {`,
			`		valxAHostname, err := formats.Parse("hostname", hdrXaHostname)`,
			`		if err != nil {`,
			`			return errors.InvalidType("X-aHostname", "header", "strfmt.Hostname", hdrXaHostname)`,
			`		o.XaHostname = *(valxAHostname.(*strfmt.Hostname))`,
			`	hdrXaInt32 := response.GetHeader("X-aInt32")`,
			`	if hdrXaInt32 != "" {`,
			`		valxAInt32, err := swag.ConvertInt32(hdrXaInt32)`,
			`		if err != nil {`,
			`			return errors.InvalidType("X-aInt32", "header", "int32", hdrXaInt32)`,
			`		o.XaInt32 = valxAInt32`,
			`	hdrXaInt64 := response.GetHeader("X-aInt64")`,
			`	if hdrXaInt64 != "" {`,
			`		valxAInt64, err := swag.ConvertInt64(hdrXaInt64)`,
			`		if err != nil {`,
			`			return errors.InvalidType("X-aInt64", "header", "int64", hdrXaInt64)`,
			`		o.XaInt64 = valxAInt64`,
			`	hdrXaMac := response.GetHeader("X-aMac")`,
			`	if hdrXaMac != "" {`,
			`		valxAMac, err := formats.Parse("mac", hdrXaMac)`,
			`		if err != nil {`,
			`			return errors.InvalidType("X-aMac", "header", "strfmt.MAC", hdrXaMac)`,
			`		o.XaMac = *(valxAMac.(*strfmt.MAC))`,
			`	hdrXaPassword := response.GetHeader("X-aPassword")`,
			`	if hdrXaPassword != "" {`,
			`		valxAPassword, err := formats.Parse("password", hdrXaPassword)`,
			`		if err != nil {`,
			`			return errors.InvalidType("X-aPassword", "header", "strfmt.Password", hdrXaPassword)`,
			`		o.XaPassword = *(valxAPassword.(*strfmt.Password))`,
			`	hdrXaRGBColor := response.GetHeader("X-aRGBColor")`,
			`	if hdrXaRGBColor != "" {`,
			`		valxARGBColor, err := formats.Parse("rgbcolor", hdrXaRGBColor)`,
			`		if err != nil {`,
			`			return errors.InvalidType("X-aRGBColor", "header", "strfmt.RGBColor", hdrXaRGBColor)`,
			`		o.XaRGBColor = *(valxARGBColor.(*strfmt.RGBColor))`,
			`	hdrXaSsn := response.GetHeader("X-aSsn")`,
			`	if hdrXaSsn != "" {`,
			`		valxASsn, err := formats.Parse("ssn", hdrXaSsn)`,
			`		if err != nil {`,
			`			return errors.InvalidType("X-aSsn", "header", "strfmt.SSN", hdrXaSsn)`,
			`		o.XaSsn = *(valxASsn.(*strfmt.SSN))`,
			`	hdrXaUUID := response.GetHeader("X-aUUID")`,
			`	if hdrXaUUID != "" {`,
			`		valxAUuid, err := formats.Parse("uuid", hdrXaUUID)`,
			`		if err != nil {`,
			`			return errors.InvalidType("X-aUUID", "header", "strfmt.UUID", hdrXaUUID)`,
			`		o.XaUUID = *(valxAUuid.(*strfmt.UUID))`,
			`	hdrXaUUID3 := response.GetHeader("X-aUUID3")`,
			`	if hdrXaUUID3 != "" {`,
			`		valxAUuid3, err := formats.Parse("uuid3", hdrXaUUID3)`,
			`		if err != nil {`,
			`			return errors.InvalidType("X-aUUID3", "header", "strfmt.UUID3", hdrXaUUID3)`,
			`		o.XaUUID3 = *(valxAUuid3.(*strfmt.UUID3))`,
			`	hdrXaUUID4 := response.GetHeader("X-aUUID4")`,
			`	if hdrXaUUID4 != "" {`,
			`		valxAUuid4, err := formats.Parse("uuid4", hdrXaUUID4)`,
			`		if err != nil {`,
			`			return errors.InvalidType("X-aUUID4", "header", "strfmt.UUID4", hdrXaUUID4)`,
			`		o.XaUUID4 = *(valxAUuid4.(*strfmt.UUID4))`,
			`	hdrXaUUID5 := response.GetHeader("X-aUUID5")`,
			`	if hdrXaUUID5 != "" {`,
			`		valxAUuid5, err := formats.Parse("uuid5", hdrXaUUID5)`,
			`		if err != nil {`,
			`			return errors.InvalidType("X-aUUID5", "header", "strfmt.UUID5", hdrXaUUID5)`,
			`		o.XaUUID5 = *(valxAUuid5.(*strfmt.UUID5))`,
			`	hdrXaUint32 := response.GetHeader("X-aUint32")`,
			`	if hdrXaUint32 != "" {`,
			`		valxAUint32, err := swag.ConvertUint32(hdrXaUint32)`,
			`		if err != nil {`,
			`			return errors.InvalidType("X-aUint32", "header", "uint32", hdrXaUint32)`,
			`		o.XaUint32 = valxAUint32`,
			`	hdrXaUint64 := response.GetHeader("X-aUint64")`,
			`	if hdrXaUint64 != "" {`,
			`		valxAUint64, err := swag.ConvertUint64(hdrXaUint64)`,
			`		if err != nil {`,
			`			return errors.InvalidType("X-aUint64", "header", "uint64", hdrXaUint64)`,
			`		o.XaUint64 = valxAUint64`,
			`	hdrXaURI := response.GetHeader("X-aUri")`,
			`	if hdrXaURI != "" {`,
			`		valxAUri, err := formats.Parse("uri", hdrXaURI)`,
			`		if err != nil {`,
			`			return errors.InvalidType("X-aUri", "header", "strfmt.URI", hdrXaURI)`,
			`		o.XaURI = *(valxAUri.(*strfmt.URI))`,
			`	hdrXAnEmail := response.GetHeader("X-anEmail")`,
			`	if hdrXAnEmail != "" {`,
			`		valxAnEmail, err := formats.Parse("email", hdrXAnEmail)`,
			`		if err != nil {`,
			`			return errors.InvalidType("X-anEmail", "header", "strfmt.Email", hdrXAnEmail)`,
			`		o.XAnEmail = *(valxAnEmail.(*strfmt.Email))`,
			`	hdrXAnISBN := response.GetHeader("X-anISBN")`,
			`	if hdrXAnISBN != "" {`,
			`		valxAnISBN, err := formats.Parse("isbn", hdrXAnISBN)`,
			`		if err != nil {`,
			`			return errors.InvalidType("X-anISBN", "header", "strfmt.ISBN", hdrXAnISBN)`,
			`		o.XAnISBN = *(valxAnISBN.(*strfmt.ISBN))`,
			`	hdrXAnISBN10 := response.GetHeader("X-anISBN10")`,
			`	if hdrXAnISBN10 != "" {`,
			`		valxAnISBN10, err := formats.Parse("isbn10", hdrXAnISBN10)`,
			`		if err != nil {`,
			`			return errors.InvalidType("X-anISBN10", "header", "strfmt.ISBN10", hdrXAnISBN10)`,
			`		o.XAnISBN10 = *(valxAnISBN10.(*strfmt.ISBN10))`,
			`	hdrXAnISBN13 := response.GetHeader("X-anISBN13")`,
			`	if hdrXAnISBN13 != "" {`,
			`		valxAnISBN13, err := formats.Parse("isbn13", hdrXAnISBN13)`,
			`		if err != nil {`,
			`			return errors.InvalidType("X-anISBN13", "header", "strfmt.ISBN13", hdrXAnISBN13)`,
			`		o.XAnISBN13 = *(valxAnISBN13.(*strfmt.ISBN13))`,
			`	hdrXAnIPV4 := response.GetHeader("X-anIpv4")`,
			`	if hdrXAnIPV4 != "" {`,
			`		valxAnIpv4, err := formats.Parse("ipv4", hdrXAnIPV4)`,
			`		if err != nil {`,
			`			return errors.InvalidType("X-anIpv4", "header", "strfmt.IPv4", hdrXAnIPV4)`,
			`		o.XAnIPV4 = *(valxAnIpv4.(*strfmt.IPv4))`,
			`	hdrXAnIPV6 := response.GetHeader("X-anIpv6")`,
			`	if hdrXAnIPV6 != "" {`,
			`		valxAnIpv6, err := formats.Parse("ipv6", hdrXAnIPV6)`,
			`		if err != nil {`,
			`			return errors.InvalidType("X-anIpv6", "header", "strfmt.IPv6", hdrXAnIPV6)`,
			`		o.XAnIPV6 = *(valxAnIpv6.(*strfmt.IPv6))`,
		},
	}

	for fileToInspect, expectedCode := range fixtureConfig {
		code, err := os.ReadFile(filepath.Join(opts.Target, filepath.FromSlash(fileToInspect)))
		require.NoError(t, err)

		for line, codeLine := range expectedCode {
			if !assertInCode(t, strings.TrimSpace(codeLine), string(code)) {
				t.Logf("Code expected did not match in codegenfile %s for expected line %d: %q", fileToInspect, line, expectedCode[line])
			}
		}
	}
}

func TestGenClient_2590(t *testing.T) {
	t.Parallel()
	defer discardOutput()()

	opts := testClientGenOpts()
	opts.Spec = filepath.Join("..", "fixtures", "bugs", "2590", "2590.yaml")

	cwd, err := os.Getwd()
	require.NoError(t, err)
	tft, err := os.MkdirTemp(cwd, "generated")
	require.NoError(t, err)
	opts.Target = tft

	t.Cleanup(func() {
		_ = os.RemoveAll(tft)
	})

	require.NoError(t,
		GenerateClient("client", []string{}, []string{}, opts),
	)

	fixtureConfig := map[string][]string{
		"client/abc/create_responses.go": { // generated file
			// expected code lines
			`payload, _ := json.Marshal(o.Payload)`,
			`return fmt.Sprintf("[POST /abc][%d] createAccepted %s", 202, payload)`,
			`return fmt.Sprintf("[POST /abc][%d] createInternalServerError %s", 500, payload)`,
		},
	}

	for fileToInspect, expectedCode := range fixtureConfig {
		code, err := os.ReadFile(filepath.Join(opts.Target, filepath.FromSlash(fileToInspect)))
		require.NoError(t, err)

		for line, codeLine := range expectedCode {
			if !assertInCode(t, strings.TrimSpace(codeLine), string(code)) {
				t.Logf("Code expected did not match in codegenfile %s for expected line %d: %q", fileToInspect, line, expectedCode[line])
			}
		}
	}
}

func TestGenClient_2773(t *testing.T) {
	t.Parallel()
	defer discardOutput()()

	opts := testClientGenOpts()
	opts.Spec = filepath.Join("..", "fixtures", "bugs", "2773", "2773.yaml")

	cwd, err := os.Getwd()
	require.NoError(t, err)
	tft, err := os.MkdirTemp(cwd, "generated")
	require.NoError(t, err)
	opts.Target = tft

	t.Cleanup(func() {
		_ = os.RemoveAll(tft)
	})

	require.NoError(t,
		GenerateClient("client", []string{}, []string{}, opts),
	)

	t.Run("generated operation should keep content type in the specified order", func(t *testing.T) {
		fixtureConfig := map[string][]string{
			"client/uploads/uploads_client.go": { // generated file
				// expected code lines
				`ProducesMediaTypes: []string{"application/octet-stream", "application/json"},`,
				`ConsumesMediaTypes: []string{"multipart/form-data", "application/x-www-form-urlencoded"},`,
			},
		}

		for fileToInspect, expectedCode := range fixtureConfig {
			code, err := os.ReadFile(filepath.Join(opts.Target, filepath.FromSlash(fileToInspect)))
			require.NoError(t, err)

			for line, codeLine := range expectedCode {
				if !assertInCode(t, strings.TrimSpace(codeLine), string(code)) {
					t.Logf("Code expected did not match in codegenfile %s for expected line %d: %q", fileToInspect, line, expectedCode[line])
				}
			}
		}
	})

	t.Run("generated operation should have options to set media type", func(t *testing.T) {
		fixtureConfig := map[string][]string{
			"client/uploads/uploads_client.go": { // generated file
				// free mime consumes option
				`func WithContentType(mime string) ClientOption {`,
				`	return func(r *runtime.ClientOperation) {`,
				`	r.ConsumesMediaTypes = []string{mime}`,
				// shorthand options
				`func WithContentTypeApplicationJSON(r *runtime.ClientOperation) {`,
				`	r.ConsumesMediaTypes = []string{"application/json"}`,
				`func WithContentTypeApplicationxWwwFormUrlencoded(r *runtime.ClientOperation) {`,
				`	r.ConsumesMediaTypes = []string{"application/x-www-form-urlencoded"}`,
				`func WithContentTypeMultipartFormData(r *runtime.ClientOperation) {`,
				`	r.ConsumesMediaTypes = []string{"multipart/form-data"}`,
				// free mime produces option
				`func WithAccept(mime string) ClientOption {`,
				`	return func(r *runtime.ClientOperation) {`,
				`		r.ProducesMediaTypes = []string{mime}`,
				// shorthand options
				`func WithAcceptApplicationJSON(r *runtime.ClientOperation) {`,
				`	r.ProducesMediaTypes = []string{"application/json"}`,
				`func WithAcceptApplicationOctetStream(r *runtime.ClientOperation) {`,
				`	r.ProducesMediaTypes = []string{"application/octet-stream"}`,
			},
		}

		for fileToInspect, expectedCode := range fixtureConfig {
			code, err := os.ReadFile(filepath.Join(opts.Target, filepath.FromSlash(fileToInspect)))
			require.NoError(t, err)

			for line, codeLine := range expectedCode {
				if !assertInCode(t, strings.TrimSpace(codeLine), string(code)) {
					t.Logf("Code expected did not match in codegenfile %s for expected line %d: %q", fileToInspect, line, expectedCode[line])
				}
			}
		}
	})
}
