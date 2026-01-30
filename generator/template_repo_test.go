// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package generator

import (
	"bytes"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/go-openapi/analysis"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"
)

func TestTemplates_CustomTemplates(t *testing.T) {
	var buf bytes.Buffer
	headerTempl, err := templates.Get("bindprimitiveparam")
	require.NoError(t, err)

	err = headerTempl.Execute(&buf, nil)
	require.NoError(t, err)
	require.NotNil(t, buf)
	assert.Equal(t, "\n", buf.String())

	buf.Reset()
	err = templates.AddFile("bindprimitiveparam", customHeader)
	require.NoError(t, err)

	headerTempl, err = templates.Get("bindprimitiveparam")
	require.NoError(t, err)
	assert.NotNil(t, headerTempl)

	err = headerTempl.Execute(&buf, nil)
	require.NoError(t, err)
	assert.Equal(t, "custom header", buf.String())
}

func TestTemplates_CustomTemplatesMultiple(t *testing.T) {
	var buf bytes.Buffer

	err := templates.AddFile("differentFileName", customMultiple)
	require.NoError(t, err)

	headerTempl, err := templates.Get("bindprimitiveparam")
	require.NoError(t, err)

	err = headerTempl.Execute(&buf, nil)
	require.NoError(t, err)

	assert.Equal(t, "custom primitive", buf.String())
}

func TestTemplates_CustomNewTemplates(t *testing.T) {
	var buf bytes.Buffer

	err := templates.AddFile("newtemplate", customNewTemplate)
	require.NoError(t, err)

	err = templates.AddFile("existingUsesNew", customExistingUsesNew)
	require.NoError(t, err)

	headerTempl, err := templates.Get("bindprimitiveparam")
	require.NoError(t, err)

	err = headerTempl.Execute(&buf, nil)
	require.NoError(t, err)

	assert.Equal(t, "new template", buf.String())
}

func TestTemplates_RepoLoadingTemplates(t *testing.T) {
	repo := NewRepository(nil)

	err := repo.AddFile("simple", singleTemplate)
	require.NoError(t, err)

	templ, err := repo.Get("simple")
	require.NoError(t, err)

	var b bytes.Buffer
	err = templ.Execute(&b, nil)
	require.NoError(t, err)

	assert.Equal(t, "test", b.String())
}

func TestTemplates_RepoLoadsAllTemplatesDefined(t *testing.T) {
	var b bytes.Buffer
	repo := NewRepository(nil)

	err := repo.AddFile("multiple", multipleDefinitions)
	require.NoError(t, err)

	templ, err := repo.Get("multiple")
	require.NoError(t, err)

	err = templ.Execute(&b, nil)
	require.NoError(t, err)

	assert.Empty(t, b.String())

	templ, err = repo.Get("T1")
	require.NoError(t, err)
	require.NotNil(t, templ)

	err = templ.Execute(&b, nil)
	require.NoError(t, err)

	assert.Equal(t, "T1", b.String())
}

type testData struct {
	Children []testData
	Name     string
	Recurse  bool
}

func TestTemplates_RepoLoadsAllDependantTemplates(t *testing.T) {
	var b bytes.Buffer
	repo := NewRepository(nil)

	err := repo.AddFile("multiple", multipleDefinitions)
	require.NoError(t, err)

	err = repo.AddFile("dependant", dependantTemplate)
	require.NoError(t, err)

	templ, err := repo.Get("dependant")
	require.NoError(t, err)
	require.NotNil(t, templ)

	err = templ.Execute(&b, nil)
	require.NoError(t, err)

	assert.Equal(t, "T1D1", b.String())
}

func TestTemplates_RepoRecursiveTemplates(t *testing.T) {
	var b bytes.Buffer
	repo := NewRepository(nil)

	err := repo.AddFile("c1", cirularDeps1)
	require.NoError(t, err)

	err = repo.AddFile("c2", cirularDeps2)
	require.NoError(t, err)

	templ, err := repo.Get("c1")
	require.NoError(t, err)
	require.NotNil(t, templ)

	data := testData{
		Name: "Root",
		Children: []testData{
			{Recurse: false},
		},
	}
	expected := `Root: Children`
	err = templ.Execute(&b, data)
	require.NoError(t, err)
	assert.Equal(t, expected, b.String())

	data = testData{
		Name: "Root",
		Children: []testData{
			{Name: "Child1", Recurse: true, Children: []testData{{Name: "Child2"}}},
		},
	}

	b.Reset()

	expected = `Root: Child1: Children`

	err = templ.Execute(&b, data)
	require.NoError(t, err)

	assert.Equal(t, expected, b.String())

	data = testData{
		Name: "Root",
		Children: []testData{
			{Name: "Child1", Recurse: false, Children: []testData{{Name: "Child2"}}},
		},
	}

	b.Reset()

	expected = `Root: Children`

	err = templ.Execute(&b, data)
	require.NoError(t, err)

	assert.Equal(t, expected, b.String())
}

// Test that definitions are available to templates
// TODO: should test also with the codeGenApp context

// Test copyright definition.
func TestTemplates_DefinitionCopyright(t *testing.T) {
	defer discardOutput()()

	const copyright = `{{ .Copyright }}`

	repo := NewRepository(nil)

	err := repo.AddFile("copyright", copyright)
	require.NoError(t, err)

	templ, err := repo.Get("copyright")
	require.NoError(t, err)
	require.NotNil(t, templ)

	opts := opts()
	opts.Copyright = "My copyright clause"
	expected := opts.Copyright

	// executes template against model definitions
	genModel, err := getModelEnvironment("../fixtures/codegen/todolist.models.yml", opts)
	require.NoError(t, err)
	require.NotNil(t, genModel)

	rendered := bytes.NewBuffer(nil)
	err = templ.Execute(rendered, genModel)
	require.NoError(t, err)
	assert.Equal(t, expected, rendered.String())

	// executes template against operations definitions
	genOperation, err := getOperationEnvironment("get", "/media/search", "../fixtures/codegen/instagram.yml", opts)
	require.NoError(t, err)
	require.NotNil(t, genOperation)

	rendered.Reset()

	err = templ.Execute(rendered, genOperation)
	require.NoError(t, err)

	assert.Equal(t, expected, rendered.String())
}

// Test TargetImportPath definition.
func TestTemplates_DefinitionTargetImportPath(t *testing.T) {
	const targetImportPath = `{{ .TargetImportPath }}`
	defer discardOutput()()

	repo := NewRepository(nil)

	err := repo.AddFile("targetimportpath", targetImportPath)
	require.NoError(t, err)

	templ, err := repo.Get("targetimportpath")
	require.NoError(t, err)
	require.NotNil(t, templ)

	opts := opts()
	// Non existing target would panic: to be tested too, but in another module
	opts.Target = "../fixtures"
	expected := "github.com/go-swagger/go-swagger/fixtures"

	// executes template against model definitions
	genModel, err := getModelEnvironment("../fixtures/codegen/todolist.models.yml", opts)
	require.NoError(t, err)
	require.NotNil(t, genModel)

	rendered := bytes.NewBuffer(nil)
	err = templ.Execute(rendered, genModel)
	require.NoError(t, err)

	assert.Equal(t, expected, rendered.String())

	// executes template against operations definitions
	genOperation, err := getOperationEnvironment("get", "/media/search", "../fixtures/codegen/instagram.yml", opts)
	require.NoError(t, err)
	require.NotNil(t, genOperation)

	rendered.Reset()

	err = templ.Execute(rendered, genOperation)
	require.NoError(t, err)

	assert.Equal(t, expected, rendered.String())
}

// Simulates a definition environment for model templates.
func getModelEnvironment(_ string, opts *GenOpts) (*GenDefinition, error) {
	defer discardOutput()()

	specDoc, err := loads.Spec("../fixtures/codegen/todolist.models.yml")
	if err != nil {
		return nil, err
	}

	definitions := specDoc.Spec().Definitions
	if len(definitions) == 0 {
		return nil, errors.New("todolist.models.yml did not return any definition")
	}

	var (
		name   string
		schema spec.Schema
	)
	for k, sch := range definitions {
		name = k
		schema = sch
		// one is enough
		break
	}

	di := discriminatorInfo(analysis.New(specDoc.Spec()))
	genModel, err := makeGenDefinition(name, "models", schema, specDoc, opts, di)
	if err != nil {
		return nil, err
	}

	return genModel, nil
}

// Simulates a definition environment for operation templates.
func getOperationEnvironment(operation string, path string, spec string, opts *GenOpts) (*GenOperation, error) {
	defer discardOutput()()

	b, err := methodPathOpBuilder(operation, path, spec)
	if err != nil {
		return nil, err
	}
	b.GenOpts = opts
	g, err := b.MakeOperation()
	if err != nil {
		return nil, err
	}
	return &g, nil
}

// AddFile() global package function (protected vs unprotected)
// Mostly unused in tests, since the Repository.AddFile()
// is generally preferred.
func TestTemplates_AddFile(t *testing.T) {
	defer discardOutput()()

	funcTpl := testFuncTpl()

	// unprotected
	err := AddFile("functpl", funcTpl)
	require.NoError(t, err)

	_, err = templates.Get("functpl")
	require.NoError(t, err)

	// protected
	err = AddFile("schemabody", funcTpl)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cannot overwrite protected template")
}

// Test LoadDir.
func TestTemplates_LoadDir(t *testing.T) {
	defer discardOutput()()

	// Fails
	err := templates.LoadDir("")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "could not complete")

	// Fails again (from any dir?)
	err = templates.LoadDir("templates")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cannot overwrite protected template")

	// TODO: success case
	// To force a success, we need to empty the global list of protected
	// templates...
	origProtectedTemplates := protectedTemplates

	defer func() {
		// Restore variable initialized with package
		protectedTemplates = origProtectedTemplates
	}()

	protectedTemplates = make(map[string]bool)
	repo := NewRepository(FuncMapFunc(DefaultLanguageFunc()))
	err = repo.LoadDir("templates")
	require.NoError(t, err)
}

// Test LoadDir.
func TestTemplates_SetAllowOverride(t *testing.T) {
	defer discardOutput()()

	// adding protected file with allowOverride set to false fails
	templates.SetAllowOverride(false)
	err := templates.AddFile("schemabody", "some data")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cannot overwrite protected template schemabody")

	// adding protected file with allowOverride set to true should not fail
	templates.SetAllowOverride(true)
	err = templates.AddFile("schemabody", "some data")
	require.NoError(t, err)
}

// Test LoadContrib.
func TestTemplates_LoadContrib(t *testing.T) {
	tests := []struct {
		name      string
		template  string
		wantError bool
	}{
		{
			name:      "None_existing_contributor_template",
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
			err := templates.LoadContrib(tt.template)
			if tt.wantError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// TODO: test error case in LoadDefaults()
// test DumpTemplates().
func TestTemplates_DumpTemplates(t *testing.T) {
	var buf bytes.Buffer
	defer captureOutput(&buf)()

	templates.DumpTemplates()
	assert.NotEmpty(t, buf)
	// Sample output
	assert.Contains(t, buf.String(), "## tupleSerializer")
	assert.Contains(t, buf.String(), "Defined in `tupleserializer.gotmpl`")
	assert.Contains(t, buf.String(), "####requires \n - schemaType")
}
