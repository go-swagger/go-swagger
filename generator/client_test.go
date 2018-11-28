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
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	defaultAPIPackage    = "operations"
	defaultClientPackage = "client"
	defaultModelPackage  = "models"
	defaultServerPackage = "restapi"

	basicFixture = "../fixtures/petstores/petstore.json"
)

func testClientGenOpts() (g GenOpts) {
	g.Target = "."
	g.APIPackage = defaultAPIPackage
	g.ModelPackage = defaultModelPackage
	g.ServerPackage = defaultServerPackage
	g.ClientPackage = defaultClientPackage
	g.Principal = ""
	g.DefaultScheme = "http"
	g.IncludeModel = true
	g.IncludeValidator = true
	g.IncludeHandler = true
	g.IncludeParameters = true
	g.IncludeResponses = true
	g.IncludeSupport = true
	g.TemplateDir = ""
	g.DumpData = false
	g.IsClient = true
	_ = g.EnsureDefaults()
	return
}

func Test_GenerateClient(t *testing.T) {
	log.SetOutput(ioutil.Discard)

	// exercise safeguards
	err := GenerateClient("test", []string{"model1"}, []string{"op1", "op2"}, nil)
	assert.Error(t, err)

	opts := testClientGenOpts()
	opts.TemplateDir = "dir/nowhere"
	err = GenerateClient("test", []string{"model1"}, []string{"op1", "op2"}, &opts)
	assert.Error(t, err)

	opts = testClientGenOpts()
	opts.TemplateDir = "http://nowhere.com"
	err = GenerateClient("test", []string{"model1"}, []string{"op1", "op2"}, &opts)
	assert.Error(t, err)

	opts = testClientGenOpts()
	opts.Spec = "dir/nowhere.yaml"
	err = GenerateClient("test", []string{"model1"}, []string{"op1", "op2"}, &opts)
	assert.Error(t, err)

	opts = testClientGenOpts()
	opts.Spec = basicFixture
	err = GenerateClient("test", []string{"model1"}, []string{}, &opts)
	assert.Error(t, err)

	opts = testClientGenOpts()
	// bad content in spec (HTML...)
	opts.Spec = "https://github.com/OAI/OpenAPI-Specification/blob/master/examples/v2.0/json/petstore.json"
	err = GenerateClient("test", []string{}, []string{}, &opts)
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
	DefaultSectionOpts(&opts)

	defer func() {
		_ = os.RemoveAll(opts.Target)
	}()
	err = GenerateClient("test", []string{}, []string{}, &opts)
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
	err = GenerateClient("test", []string{}, []string{}, &opts)
	assert.NoError(t, err)
	_, err = os.Stat(filepath.Join(tgt, "stdout"))
	assert.NoError(t, err)
}

func TestClient(t *testing.T) {
	targetdir, err := ioutil.TempDir(os.TempDir(), "swagger_nogo")
	if err != nil {
		t.Fatalf("Failed to create a test target directory: %v", err)
	}
	log.SetOutput(ioutil.Discard)
	defer func() {
		_ = os.RemoveAll(targetdir)
		log.SetOutput(os.Stdout)
	}()

	tests := []struct {
		name      string
		template  string
		wantError bool
		prepare   func(opts *GenOpts)
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
			name:      "BaseImportDisabled",
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := testClientGenOpts()
			opts.Target = targetdir
			opts.Spec = basicFixture
			opts.LanguageOpts.BaseImportFunc = nil
			opts.Template = tt.template

			if tt.prepare != nil {
				tt.prepare(&opts)
			}

			err := GenerateClient("foo", nil, nil, &opts)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
