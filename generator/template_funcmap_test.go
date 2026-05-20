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
	// singleTemplate        = `test`.
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

	opts := opts()
	modelTpl := testModelTpl()

	require.NoError(t, opts.templates.AddFile("modeltpl", modelTpl))

	templ, err := opts.templates.Get("modeltpl")
	require.NoError(t, err)

	genModel, err := getModelEnvironment("../fixtures/codegen/todolist.models.yml", opts)
	require.NoError(t, err)

	genModel.DependsOn = []string{"x", "z"}
	rendered := bytes.NewBuffer(nil)
	require.NoError(t, templ.Execute(rendered, genModel))

	assert.StringContainsT(t, rendered.String(), "ContainsString=true\n")
	assert.StringContainsT(t, rendered.String(), "DoesNotContainString=false\n")
	assert.StringContainsT(t, rendered.String(), `Json={`+
		`"conv":"github.com/go-openapi/swag/conv",`+
		`"errors":"github.com/go-openapi/errors",`+
		`"jsonutils":"github.com/go-openapi/swag/jsonutils",`+
		`"netutils":"github.com/go-openapi/swag/netutils",`+
		`"runtime":"github.com/go-openapi/runtime",`+
		`"strfmt":"github.com/go-openapi/strfmt",`+
		`"stringutils":"github.com/go-openapi/swag/stringutils",`+
		`"typeutils":"github.com/go-openapi/swag/typeutils",`+
		`"validate":"github.com/go-openapi/validate"`+
		`}`,
	)
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
