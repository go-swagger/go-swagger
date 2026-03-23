// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package generator

import (
	"bytes"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
)

const (
	// Test template environment.
	singleTemplate        = `test`
	multipleDefinitions   = `{{ define "T1" }}T1{{end}}{{ define "T2" }}T2{{end}}`
	dependantTemplate     = `{{ template "T1" }}D1`
	cirularDeps1          = `{{ define "T1" }}{{ .Name }}: {{ range .Children }}{{ template "T2" . }}{{end}}{{end}}{{template "T1" . }}`
	cirularDeps2          = `{{ define "T2" }}{{if .Recurse }}{{ template "T1" . }}{{ else }}Children{{end}}{{end}}`
	customHeader          = `custom header`
	customMultiple        = `{{define "bindprimitiveparam" }}custom primitive{{end}}`
	customNewTemplate     = `new template`
	customExistingUsesNew = `{{define "bindprimitiveparam" }}{{ template "newtemplate" }}{{end}}`
)

// testModelTpl exercises the full funcmap against a GenDefinition context,
// including LanguageOpts-dependent and model-dependent functions.
func testModelTpl() string {
	return `
ContainsString={{ contains .DependsOn "x"}}
DoesNotContainString={{ contains .DependsOn "y"}}
Json={{ json .DefaultImports }}
PrettyJson={{ prettyjson . }}
Snakize={{ snakize "WeArePonies_Of_the_round table" }}
Snakize1={{ snakize "endingInOsNameLinux" }}
Snakize2={{ snakize "endingInArchNameLinuxAmd64" }}
Snakize3={{ snakize "endingInTest" }}
toPackage1={{ toPackage "a/b-c/d-e" }}
toPackage2={{ toPackage "a.a/b_c/d_e" }}
toPackage3={{ toPackage "d_e" }}
toPackage4={{ toPackage "d-e" }}
toPackageName={{ toPackageName "d-e/f-g" }}
`
}

// TestTemplates_FuncMap_Model exercises the funcmap against a real GenDefinition.
// This is an integration test between the template repository, the funcmap, and
// the generator types (LanguageOpts, GenDefinition).
func TestTemplates_FuncMap_Model(t *testing.T) {
	defer discardOutput()()

	modelTpl := testModelTpl()

	err := templates.AddFile("modeltpl", modelTpl)
	require.NoError(t, err)

	templ, err := templates.Get("modeltpl")
	require.NoError(t, err)

	opts := opts()
	genModel, err := getModelEnvironment("../fixtures/codegen/todolist.models.yml", opts)
	require.NoError(t, err)

	genModel.DependsOn = []string{"x", "z"}
	rendered := bytes.NewBuffer(nil)
	err = templ.Execute(rendered, genModel)
	require.NoError(t, err)

	assert.StringContainsT(t, rendered.String(), "ContainsString=true\n")
	assert.StringContainsT(t, rendered.String(), "DoesNotContainString=false\n")
	assert.StringContainsT(t, rendered.String(), `Json={"errors":"github.com/go-openapi/errors","runtime":"github.com/go-openapi/runtime","strfmt":"github.com/go-openapi/strfmt","swag":"github.com/go-openapi/swag","validate":"github.com/go-openapi/validate"}`)
	assert.StringContainsT(t, rendered.String(), "\"TargetImportPath\": \"github.com/go-swagger/go-swagger/generator\"")

	// LanguageOpts-dependent assertions
	assert.StringContainsT(t, rendered.String(), "Snakize=we_are_ponies_of_the_round_table\n")
	assert.StringContainsT(t, rendered.String(), "Snakize1=ending_in_os_name_linux_swagger\n")
	assert.StringContainsT(t, rendered.String(), "Snakize2=ending_in_arch_name_linux_amd64_swagger\n")
	assert.StringContainsT(t, rendered.String(), "Snakize3=ending_in_test_swagger\n")
	assert.StringContainsT(t, rendered.String(), "toPackage1=a/b-c/d_e\n")
	assert.StringContainsT(t, rendered.String(), "toPackage2=a.a/b_c/d_e\n")
	assert.StringContainsT(t, rendered.String(), "toPackage3=d_e\n")
	assert.StringContainsT(t, rendered.String(), "toPackage4=d_e\n")
	assert.StringContainsT(t, rendered.String(), "toPackageName=f_g\n")
}
