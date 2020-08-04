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
	"io/ioutil"
	"log"
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
	log.SetOutput(ioutil.Discard)

	// exercise safeguards
	err := GenerateClient("test", []string{"model1"}, []string{"op1", "op2"}, nil)
	assert.Error(t, err)

	opts := testClientGenOpts()
	opts.TemplateDir = "dir/nowhere"
	err = GenerateClient("test", []string{"model1"}, []string{"op1", "op2"}, opts)
	assert.Error(t, err)

	opts = testClientGenOpts()
	opts.TemplateDir = "http://nowhere.com"
	err = GenerateClient("test", []string{"model1"}, []string{"op1", "op2"}, opts)
	assert.Error(t, err)

	opts = testClientGenOpts()
	opts.Spec = "dir/nowhere.yaml"
	err = GenerateClient("test", []string{"model1"}, []string{"op1", "op2"}, opts)
	assert.Error(t, err)

	opts = testClientGenOpts()
	opts.Spec = basicFixture
	err = GenerateClient("test", []string{"model1"}, []string{}, opts)
	assert.Error(t, err)

	opts = testClientGenOpts()
	// bad content in spec (HTML...)
	opts.Spec = "https://github.com/OAI/OpenAPI-Specification/blob/master/examples/v2.0/json/petstore.json"
	err = GenerateClient("test", []string{}, []string{}, opts)
	assert.Error(t, err)

	opts = testClientGenOpts()
	// no operations selected
	opts.Spec = "https://raw.githubusercontent.com/OAI/OpenAPI-Specification/master/examples/v2.0/yaml/petstore.yaml"
	err = GenerateClient("test", []string{}, []string{"wrongOperationID"}, opts)
	assert.Error(t, err)

	opts = testClientGenOpts()
	// generate remote spec
	opts.Spec = "https://raw.githubusercontent.com/OAI/OpenAPI-Specification/master/examples/v2.0/yaml/petstore.yaml"
	cwd, _ := os.Getwd()
	tft, _ := ioutil.TempDir(cwd, "generated")
	defer func() {
		_ = os.RemoveAll(tft)
	}()
	opts.Target = tft
	opts.IsClient = true
	DefaultSectionOpts(opts)

	defer func() {
		_ = os.RemoveAll(opts.Target)
	}()
	err = GenerateClient("test", []string{}, []string{}, opts)
	assert.NoError(t, err)

	// just checks this does not fail
	origStdout := os.Stdout
	defer func() {
		os.Stdout = origStdout
	}()
	tgt, _ := ioutil.TempDir(cwd, "dumped")
	defer func() {
		_ = os.RemoveAll(tgt)
	}()
	os.Stdout, _ = os.Create(filepath.Join(tgt, "stdout"))
	opts.DumpData = true
	err = GenerateClient("test", []string{}, []string{}, opts)
	assert.NoError(t, err)
	_, err = os.Stat(filepath.Join(tgt, "stdout"))
	assert.NoError(t, err)
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
	log.SetOutput(ioutil.Discard)

	base := os.Getenv("GOPATH")
	if base == "" {
		base = "."
	} else {
		base = filepath.Join(base, "src")
		err := os.MkdirAll(base, 0755)
		require.NoError(t, err)
	}
	targetdir, err := ioutil.TempDir(base, "swagger_nogo")
	require.NoError(t, err, "Failed to create a test target directory: %v", err)

	defer func() {
		_ = os.RemoveAll(targetdir)
		log.SetOutput(os.Stdout)
	}()

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

				buf, err := ioutil.ReadFile(filepath.Join(target, "foo_client.go"))
				require.NoError(t, err)

				// assert client import, with deconfliction
				code := string(buf)
				baseImport := `swagger_nogo\d+/packages_mangling/client`
				assertImports(t, baseImport, code)

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

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := testClientGenOpts()
			opts.Spec = basicFixture
			opts.Target = filepath.Join(targetdir, opts.LanguageOpts.ManglePackageName(tt.name, "client_test"+strconv.Itoa(i)))
			err := os.MkdirAll(opts.Target, 0755)
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
}

func TestGenClient_1518(t *testing.T) {
	// test client response handling when unexpected success response kicks in
	log.SetOutput(ioutil.Discard)
	defer func() {
		log.SetOutput(os.Stdout)
	}()

	opts := testClientGenOpts()
	opts.Spec = filepath.Join("..", "fixtures", "bugs", "1518", "fixture-1518.yaml")

	cwd, _ := os.Getwd()
	tft, _ := ioutil.TempDir(cwd, "generated")
	opts.Target = tft

	defer func() {
		_ = os.RemoveAll(opts.Target)
	}()

	err := GenerateClient("client", []string{}, []string{}, opts)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

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
		code, err := ioutil.ReadFile(filepath.Join(opts.Target, filepath.FromSlash(fileToInspect)))
		require.NoError(t, err)
		for line, codeLine := range expectedCode {
			if !assertInCode(t, strings.TrimSpace(codeLine), string(code)) {
				t.Logf("Code expected did not match in codegenfile %s for expected line %d: %q", fileToInspect, line, expectedCode[line])
			}
		}
	}
}
