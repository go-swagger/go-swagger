// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package generator

import (
	"bytes"
	"errors"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"

	"github.com/go-swagger/go-swagger/generator/internal/gentest"
	templatesrepo "github.com/go-swagger/go-swagger/generator/internal/templates-repo"
)

func TestTemplates_CustomTemplates(t *testing.T) {
	var buf bytes.Buffer
	opts := opts()
	headerTempl, err := opts.templates.Get("bindprimitiveparam")
	require.NoError(t, err)
	require.NoError(t, headerTempl.Execute(&buf, nil))
	assert.EqualT(t, "\n", buf.String())

	buf.Reset()
	require.NoError(t, opts.templates.AddFile("bindprimitiveparam", customHeader))

	headerTempl, err = opts.templates.Get("bindprimitiveparam")
	require.NoError(t, err)
	assert.NotNil(t, headerTempl)
	require.NoError(t, headerTempl.Execute(&buf, nil))
	assert.EqualT(t, "custom header", buf.String())
}

func TestTemplates_CustomTemplatesMultiple(t *testing.T) {
	var buf bytes.Buffer
	opts := opts()
	require.NoError(t, opts.templates.AddFile("differentFileName", customMultiple))
	headerTempl, err := opts.templates.Get("bindprimitiveparam")
	require.NoError(t, err)
	require.NoError(t, headerTempl.Execute(&buf, nil))
	assert.EqualT(t, "custom primitive", buf.String())
}

func TestTemplates_CustomNewTemplates(t *testing.T) {
	var buf bytes.Buffer
	opts := opts()
	require.NoError(t, opts.templates.AddFile("newtemplate", customNewTemplate))
	require.NoError(t, opts.templates.AddFile("existingUsesNew", customExistingUsesNew))
	headerTempl, err := opts.templates.Get("bindprimitiveparam")
	require.NoError(t, err)
	require.NoError(t, headerTempl.Execute(&buf, nil))
	assert.EqualT(t, "new template", buf.String())
}

// Test that definitions are available to templates
// TODO: should test also with the codeGenApp context

// Test copyright definition.
func TestTemplates_DefinitionCopyright(t *testing.T) {
	defer discardOutput()()

	const copyright = `{{ .Copyright }}`

	repo := templatesrepo.NewRepository(nil)
	require.NoError(t, repo.AddFile("copyright", copyright))
	templ, err := repo.Get("copyright")
	require.NoError(t, err)
	require.NotNil(t, templ)

	opts := opts()
	const expected = "My copyright clause"
	opts.Copyright = expected

	// executes template against model definitions
	genModel, err := getModelEnvironment("../fixtures/codegen/todolist.models.yml", opts)
	require.NoError(t, err)
	require.NotNil(t, genModel)
	rendered := bytes.NewBuffer(nil)
	require.NoError(t, templ.Execute(rendered, genModel))
	assert.EqualT(t, expected, rendered.String())

	// executes template against operations definitions
	genOperation, err := getOperationEnvironment("get", "/media/search", "../fixtures/codegen/instagram.yml", opts)
	require.NoError(t, err)
	require.NotNil(t, genOperation)
	rendered.Reset()
	require.NoError(t, templ.Execute(rendered, genOperation))
	assert.EqualT(t, expected, rendered.String())
}

// Test TargetImportPath definition.
func TestTemplates_DefinitionTargetImportPath(t *testing.T) {
	const targetImportPath = `{{ .TargetImportPath }}`
	defer discardOutput()()

	repo := templatesrepo.NewRepository(nil)
	require.NoError(t, repo.AddFile("targetimportpath", targetImportPath))
	templ, err := repo.Get("targetimportpath")
	require.NoError(t, err)
	require.NotNil(t, templ)

	// Non existing target would panic: to be tested too, but in another module
	opts := opts()
	opts.Target = "../fixtures"
	expected := "github.com/go-swagger/go-swagger/fixtures"

	// executes template against model definitions
	genModel, err := getModelEnvironment("../fixtures/codegen/todolist.models.yml", opts)
	require.NoError(t, err)
	require.NotNil(t, genModel)

	rendered := bytes.NewBuffer(nil)
	require.NoError(t, templ.Execute(rendered, genModel))
	assert.EqualT(t, expected, rendered.String())

	// executes template against operations definitions
	genOperation, err := getOperationEnvironment("get", "/media/search", "../fixtures/codegen/instagram.yml", opts)
	require.NoError(t, err)
	require.NotNil(t, genOperation)

	rendered.Reset()
	require.NoError(t, templ.Execute(rendered, genOperation))
	assert.EqualT(t, expected, rendered.String())
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

	genModel, err := makeGenDefinition(name, "models", schema, specDoc, opts)
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

// AddFile on the global repository (protected vs unprotected).
func TestTemplates_AddFile(t *testing.T) {
	defer discardOutput()()

	const funcTpl = `{{ pascalize "hello world" }}`
	opts := opts()

	t.Run("should load as override an unprotected template", func(t *testing.T) {
		require.NoError(t, opts.templates.AddFile("functpl", funcTpl))
		_, err := opts.templates.Get("functpl")
		require.NoError(t, err)
	})

	t.Run("should not load a protected template", func(t *testing.T) {
		err := opts.templates.AddFile("schemabody", funcTpl)
		require.Error(t, err)
		assert.ErrorContains(t, err, "cannot overwrite protected template")
	})
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
	opts := opts()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := opts.templates.LoadContrib(tt.template, embeddedAssets{})
			if tt.wantError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// test DumpTemplates() with actual templates.
func TestTemplates_DumpTemplates(t *testing.T) {
	var buf bytes.Buffer
	defer gentest.CaptureOutput(&buf)()

	opts := opts()
	opts.templates.DumpTemplates()
	assert.NotEmpty(t, buf)

	// Sample output
	str := buf.String()
	assert.StringContainsT(t, str, "## tupleSerializer")
	assert.StringContainsT(t, str, "Defined in `tupleserializer.gotmpl`")
	assert.StringContainsT(t, str, "####requires \n - schemaType")
}
