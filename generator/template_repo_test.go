package generator

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/swag"
)

const (
	// Test template environment
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

// Test copyright definition
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

// Test TargetImportPath definition
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

// Simulates a definition environment for model templates
func getModelEnvironment(_ string, opts *GenOpts) (*GenDefinition, error) {
	defer discardOutput()()

	specDoc, err := loads.Spec("../fixtures/codegen/todolist.models.yml")
	if err != nil {
		return nil, err
	}
	definitions := specDoc.Spec().Definitions

	for k, schema := range definitions {
		genModel, err := makeGenDefinition(k, "models", schema, specDoc, opts)
		if err != nil {
			return nil, err
		}
		// One is enough
		return genModel, nil
	}
	return nil, nil
}

// Simulates a definition environment for operation templates
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

	assert.Contains(t, rendered.String(), "Pascalize=WeArePoniesOfTheRoundTable\n")
	assert.Contains(t, rendered.String(), "Snakize=we_are_ponies_of_the_round_table\n")
	assert.Contains(t, rendered.String(), "Humanize=we are ponies of the round table\n")
	assert.Contains(t, rendered.String(), "PluralizeFirstWord=ponies of the round table\n")
	assert.Contains(t, rendered.String(), "PluralizeFirstOfOneWord=dwarves\n")
	assert.Contains(t, rendered.String(), "PluralizeFirstOfNoWord=\n")
	assert.Contains(t, rendered.String(), "DropPackage=suffix\n")
	assert.Contains(t, rendered.String(), "DropNoPackage=suffix\n")
	assert.Contains(t, rendered.String(), "DropEmptyPackage=\n")
	assert.Contains(t, rendered.String(), "DropEmptyPackage=\n")
	assert.Contains(t, rendered.String(), "ContainsString=true\n")
	assert.Contains(t, rendered.String(), "DoesNotContainString=false\n")
	assert.Contains(t, rendered.String(), "PadSurround1=-,-,-,padme,-,-,-,-,-,-,-,-\n")
	assert.Contains(t, rendered.String(), "PadSurround2=padme,-,-,-,-,-,-,-,-,-,-,-\n")
	assert.Contains(t, rendered.String(), `Json={"errors":"github.com/go-openapi/errors","runtime":"github.com/go-openapi/runtime","strfmt":"github.com/go-openapi/strfmt","swag":"github.com/go-openapi/swag","validate":"github.com/go-openapi/validate"}`)
	assert.Contains(t, rendered.String(), "\"TargetImportPath\": \"github.com/go-swagger/go-swagger/generator\"")
	assert.Contains(t, rendered.String(), "Snakize1=ending_in_os_name_linux_swagger\n")
	assert.Contains(t, rendered.String(), "Snakize2=ending_in_arch_name_linux_amd64_swagger\n")
	assert.Contains(t, rendered.String(), "Snakize3=ending_in_test_swagger\n")
	assert.Contains(t, rendered.String(), "toPackage1=a/b-c/d_e\n")
	assert.Contains(t, rendered.String(), "toPackage2=a.a/b_c/d_e\n")
	assert.Contains(t, rendered.String(), "toPackage3=d_e\n")
	assert.Contains(t, rendered.String(), "toPackage4=d_e\n")
	assert.Contains(t, rendered.String(), "toPackageName=f_g\n")
	assert.Contains(t, rendered.String(), "PascalizeSpecialChar1=Plus1\n")
	assert.Contains(t, rendered.String(), "PascalizeSpecialChar2=Minus1\n")
	assert.Contains(t, rendered.String(), "PascalizeSpecialChar3=Nr1\n")
	assert.Contains(t, rendered.String(), "PascalizeSpecialChar4=Minus\n")
	assert.Contains(t, rendered.String(), "PascalizeSpecialChar5=Plus\n")
	assert.Contains(t, rendered.String(), "PascalizeCleanupEnumVariant1=Nr2Dot4Ghz")
	assert.Contains(t, rendered.String(), "Dict=Pony of the round table\n")
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

// Test LoadDir
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

// Test LoadDir
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

// Test LoadContrib
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
// test DumpTemplates()
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

func TestFuncMap_Pascalize(t *testing.T) {
	assert.Equal(t, "Plus1", pascalize("+1"))
	assert.Equal(t, "Plus", pascalize("+"))
	assert.Equal(t, "Minus1", pascalize("-1"))
	assert.Equal(t, "Minus", pascalize("-"))
	assert.Equal(t, "Nr8", pascalize("8"))
	assert.Equal(t, "Asterisk", pascalize("*"))
	assert.Equal(t, "ForwardSlash", pascalize("/"))
	assert.Equal(t, "EqualSign", pascalize("="))

	assert.Equal(t, "Hello", pascalize("+hello"))

	// other values from swag rules
	assert.Equal(t, "At8", pascalize("@8"))
	assert.Equal(t, "AtHello", pascalize("@hello"))
	assert.Equal(t, "Bang8", pascalize("!8"))
	assert.Equal(t, "At", pascalize("@"))

	// # values
	assert.Equal(t, "Hello", pascalize("#hello"))
	assert.Equal(t, "BangHello", pascalize("#!hello"))
	assert.Equal(t, "HashTag8", pascalize("#8"))
	assert.Equal(t, "HashTag", pascalize("#"))

	// single '_'
	assert.Equal(t, "Nr", pascalize("_"))
	assert.Equal(t, "Hello", pascalize("_hello"))

	// remove spaces
	assert.Equal(t, "HelloWorld", pascalize("# hello world"))
	assert.Equal(t, "HashTag8HelloWorld", pascalize("# 8 hello world"))

	assert.Equal(t, "Empty", pascalize(""))
}

func TestFuncMap_DropPackage(t *testing.T) {
	assert.Equal(t, "trail", dropPackage("base.trail"))
	assert.Equal(t, "trail", dropPackage("base.another.trail"))
	assert.Equal(t, "trail", dropPackage("trail"))
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
		assert.JSONEq(t, `{"a":"good","B":3}`, res)

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
	require.True(t, gt0(swag.Int64(1)))
	require.False(t, gt0(swag.Int64(0)))
	require.False(t, gt0(nil))
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

	assert.Contains(t, rendered.String(), "Pascalize=GetDollarRef\n")
	assert.Contains(t, rendered.String(), "Camelize=getDollarRef\n")
}
