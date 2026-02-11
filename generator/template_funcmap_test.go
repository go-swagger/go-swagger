// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package generator

import (
	"bytes"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"

	"github.com/go-openapi/swag"
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

func testFuncTpl() string {
	return `
Pascalize={{ pascalize "WeArePonies_Of_the_round table" }}
Snakize={{ snakize "WeArePonies_Of_the_round table" }}
Humanize={{ humanize "WeArePonies_Of_the_round table" }}
PluralizeFirstWord={{ pluralizeFirstWord "pony of the round table" }}
PluralizeFirstOfOneWord={{ pluralizeFirstWord "dwarf" }}
PluralizeFirstOfNoWord={{ pluralizeFirstWord "" }}
DropPackage={{ dropPackage "prefix.suffix" }}
DropNoPackage={{ dropPackage "suffix" }}
DropEmptyPackage={{ dropPackage "" }}
ContainsString={{ contains .DependsOn "x"}}
DoesNotContainString={{ contains .DependsOn "y"}}
PadSurround1={{ padSurround "padme" "-" 3 12}}
PadSurround2={{ padSurround "padme" "-" 0 12}}
Json={{ json .DefaultImports }}
PrettyJson={{ prettyjson . }}
Snakize1={{ snakize "endingInOsNameLinux" }}
Snakize2={{ snakize "endingInArchNameLinuxAmd64" }}
Snakize3={{ snakize "endingInTest" }}
toPackage1={{ toPackage "a/b-c/d-e" }}
toPackage2={{ toPackage "a.a/b_c/d_e" }}
toPackage3={{ toPackage "d_e" }}
toPackage4={{ toPackage "d-e" }}
toPackageName={{ toPackageName "d-e/f-g" }}
PascalizeSpecialChar1={{ pascalize "+1" }}
PascalizeSpecialChar2={{ pascalize "-1" }}
PascalizeSpecialChar3={{ pascalize "1" }}
PascalizeSpecialChar4={{ pascalize "-" }}
PascalizeSpecialChar5={{ pascalize "+" }}
PascalizeCleanupEnumVariant1={{ pascalize (cleanupEnumVariant "2.4Ghz") }}
Dict={{ template "dictTemplate" dict "Animal" "Pony" "Shape" "round" "Furniture" "table" }}
{{ define "dictTemplate" }}{{ .Animal }} of the {{ .Shape }} {{ .Furniture }}{{ end }}
`
}

// Exercises FuncMap
// Just running basic tests to make sure the function map works and all functions are available as expected.
// More complete unit tests are provided by go-openapi/swag.
func TestTemplates_FuncMap(t *testing.T) {
	defer discardOutput()()

	funcTpl := testFuncTpl()

	err := templates.AddFile("functpl", funcTpl)
	require.NoError(t, err)

	templ, err := templates.Get("functpl")
	require.NoError(t, err)

	opts := opts()
	// executes template against model definitions
	genModel, err := getModelEnvironment("../fixtures/codegen/todolist.models.yml", opts)
	require.NoError(t, err)

	genModel.DependsOn = []string{"x", "z"}
	rendered := bytes.NewBuffer(nil)
	err = templ.Execute(rendered, genModel)
	require.NoError(t, err)

	assert.StringContainsT(t, rendered.String(), "Pascalize=WeArePoniesOfTheRoundTable\n")
	assert.StringContainsT(t, rendered.String(), "Snakize=we_are_ponies_of_the_round_table\n")
	assert.StringContainsT(t, rendered.String(), "Humanize=we are ponies of the round table\n")
	assert.StringContainsT(t, rendered.String(), "PluralizeFirstWord=ponies of the round table\n")
	assert.StringContainsT(t, rendered.String(), "PluralizeFirstOfOneWord=dwarves\n")
	assert.StringContainsT(t, rendered.String(), "PluralizeFirstOfNoWord=\n")
	assert.StringContainsT(t, rendered.String(), "DropPackage=suffix\n")
	assert.StringContainsT(t, rendered.String(), "DropNoPackage=suffix\n")
	assert.StringContainsT(t, rendered.String(), "DropEmptyPackage=\n")
	assert.StringContainsT(t, rendered.String(), "DropEmptyPackage=\n")
	assert.StringContainsT(t, rendered.String(), "ContainsString=true\n")
	assert.StringContainsT(t, rendered.String(), "DoesNotContainString=false\n")
	assert.StringContainsT(t, rendered.String(), "PadSurround1=-,-,-,padme,-,-,-,-,-,-,-,-\n")
	assert.StringContainsT(t, rendered.String(), "PadSurround2=padme,-,-,-,-,-,-,-,-,-,-,-\n")
	assert.StringContainsT(t, rendered.String(), `Json={"errors":"github.com/go-openapi/errors","runtime":"github.com/go-openapi/runtime","strfmt":"github.com/go-openapi/strfmt","swag":"github.com/go-openapi/swag","validate":"github.com/go-openapi/validate"}`)
	assert.StringContainsT(t, rendered.String(), "\"TargetImportPath\": \"github.com/go-swagger/go-swagger/generator\"")
	assert.StringContainsT(t, rendered.String(), "Snakize1=ending_in_os_name_linux_swagger\n")
	assert.StringContainsT(t, rendered.String(), "Snakize2=ending_in_arch_name_linux_amd64_swagger\n")
	assert.StringContainsT(t, rendered.String(), "Snakize3=ending_in_test_swagger\n")
	assert.StringContainsT(t, rendered.String(), "toPackage1=a/b-c/d_e\n")
	assert.StringContainsT(t, rendered.String(), "toPackage2=a.a/b_c/d_e\n")
	assert.StringContainsT(t, rendered.String(), "toPackage3=d_e\n")
	assert.StringContainsT(t, rendered.String(), "toPackage4=d_e\n")
	assert.StringContainsT(t, rendered.String(), "toPackageName=f_g\n")
	assert.StringContainsT(t, rendered.String(), "PascalizeSpecialChar1=Plus1\n")
	assert.StringContainsT(t, rendered.String(), "PascalizeSpecialChar2=Minus1\n")
	assert.StringContainsT(t, rendered.String(), "PascalizeSpecialChar3=Nr1\n")
	assert.StringContainsT(t, rendered.String(), "PascalizeSpecialChar4=Minus\n")
	assert.StringContainsT(t, rendered.String(), "PascalizeSpecialChar5=Plus\n")
	assert.StringContainsT(t, rendered.String(), "PascalizeCleanupEnumVariant1=Nr2Dot4Ghz")
	assert.StringContainsT(t, rendered.String(), "Dict=Pony of the round table\n")
}

func TestFuncMap_DropPackage(t *testing.T) {
	assert.EqualT(t, "trail", dropPackage("base.trail"))
	assert.EqualT(t, "trail", dropPackage("base.another.trail"))
	assert.EqualT(t, "trail", dropPackage("trail"))
}

func TestFuncMap_Pascalize(t *testing.T) {
	assert.EqualT(t, "Plus1", pascalize("+1"))
	assert.EqualT(t, "Plus", pascalize("+"))
	assert.EqualT(t, "Minus1", pascalize("-1"))
	assert.EqualT(t, "Minus", pascalize("-"))
	assert.EqualT(t, "Nr8", pascalize("8"))
	assert.EqualT(t, "Asterisk", pascalize("*"))
	assert.EqualT(t, "ForwardSlash", pascalize("/"))
	assert.EqualT(t, "EqualSign", pascalize("="))

	assert.EqualT(t, "Hello", pascalize("+hello"))

	// other values from swag rules
	assert.EqualT(t, "At8", pascalize("@8"))
	assert.EqualT(t, "AtHello", pascalize("@hello"))
	assert.EqualT(t, "Bang8", pascalize("!8"))
	assert.EqualT(t, "At", pascalize("@"))

	// # values
	assert.EqualT(t, "Hello", pascalize("#hello"))
	assert.EqualT(t, "BangHello", pascalize("#!hello"))
	assert.EqualT(t, "HashTag8", pascalize("#8"))
	assert.EqualT(t, "HashTag", pascalize("#"))

	// single '_'
	assert.EqualT(t, "Nr", pascalize("_"))
	assert.EqualT(t, "Hello", pascalize("_hello"))

	// remove spaces
	assert.EqualT(t, "HelloWorld", pascalize("# hello world"))
	assert.EqualT(t, "HashTag8HelloWorld", pascalize("# 8 hello world"))

	assert.EqualT(t, "Empty", pascalize(""))
}

func TestFuncMap_AsJSON(t *testing.T) {
	for _, jsonFunc := range []func(any) (string, error){
		asJSON,
		asPrettyJSON,
	} {
		res, err := jsonFunc(struct {
			A string `json:"a"`
			B int
		}{A: "good", B: 3})
		require.NoError(t, err)
		assert.JSONEqT(t, `{"a":"good","B":3}`, res)

		_, err = jsonFunc(struct {
			A string `json:"a"`
			B func() string
		}{A: "good", B: func() string { return "" }})
		require.Error(t, err)
	}
}

func TestFuncMap_Dict(t *testing.T) {
	d, err := dict("a", "b", "c", "d")
	require.NoError(t, err)
	assert.Equal(t, map[string]any{"a": "b", "c": "d"}, d)

	// odd number of arguments
	_, err = dict("a", "b", "c")
	require.Error(t, err)

	// none-string key
	_, err = dict("a", "b", 3, "d")
	require.Error(t, err)
}

func TestIsInteger(t *testing.T) {
	var (
		nilString *string
		nilInt    *int
		nilFloat  *float32
	)

	for _, anInteger := range []any{
		int8(4),
		int16(4),
		int32(4),
		int64(4),
		int(4),
		swag.Int(4),
		swag.Int32(4),
		swag.Int64(4),
		swag.Uint(4),
		swag.Uint32(4),
		swag.Uint64(4),
		float32(12),
		float64(12),
		swag.Float32(12),
		swag.Float64(12),
		"12",
		swag.String("12"),
	} {
		val := anInteger
		require.Truef(t, isInteger(val), "expected %#v to be detected an integer value", val)
	}

	for _, notAnInteger := range []any{
		float32(12.5),
		float64(12.5),
		swag.Float32(12.5),
		swag.Float64(12.5),
		[]string{"a"},
		struct{}{},
		nil,
		map[string]int{"a": 1},
		"abc",
		"2.34",
		swag.String("2.34"),
		nilString,
		nilInt,
		nilFloat,
	} {
		val := notAnInteger
		require.Falsef(t, isInteger(val), "did not expect %#v to be detected an integer value", val)
	}
}

func TestGt0(t *testing.T) {
	require.TrueT(t, gt0(swag.Int64(1)))  //nolint:staticcheck // have to migrate to the new swag API
	require.FalseT(t, gt0(swag.Int64(0))) //nolint:staticcheck // have to migrate to the new swag API
	require.FalseT(t, gt0(nil))
}

func TestIssue2821(t *testing.T) {
	tpl := `
Pascalize={{ pascalize . }}
Camelize={{ camelize . }}
`

	require.NoError(t,
		templates.AddFile("functpl", tpl),
	)

	compiled, err := templates.Get("functpl")
	require.NoError(t, err)

	rendered := bytes.NewBuffer(nil)
	require.NoError(t,
		compiled.Execute(rendered, "get$ref"),
	)

	assert.StringContainsT(t, rendered.String(), "Pascalize=GetDollarRef\n")
	assert.StringContainsT(t, rendered.String(), "Camelize=getDollarRef\n")
}
