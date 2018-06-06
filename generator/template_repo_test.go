package generator

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/go-openapi/loads"
	"github.com/stretchr/testify/assert"
)

var (
	singleTemplate        = `test`
	multipleDefinitions   = `{{ define "T1" }}T1{{end}}{{ define "T2" }}T2{{end}}`
	dependantTemplate     = `{{ template "T1" }}D1`
	cirularDeps1          = `{{ define "T1" }}{{ .Name }}: {{ range .Children }}{{ template "T2" . }}{{end}}{{end}}{{template "T1" . }}`
	cirularDeps2          = `{{ define "T2" }}{{if .Recurse }}{{ template "T1" . }}{{ else }}Children{{end}}{{end}}`
	customHeader          = `custom header`
	customMultiple        = `{{define "bindprimitiveparam" }}custom primitive{{end}}`
	customNewTemplate     = `new template`
	customExistingUsesNew = `{{define "bindprimitiveparam" }}{{ template "newtemplate" }}{{end}}`
	// Test template environment
	copyright        = `{{ .Copyright }}`
	targetImportPath = `{{ .TargetImportPath }}`
	funcTpl          = `
Pascalize={{ pascalize "WeArePonies_Of_the_round table" }}
Snakize={{ snakize "WeArePonies_Of_the_round table" }}
Humanize={{ humanize "WeArePonies_Of_the_round table" }}
PluralizeFirstWord={{ pluralizeFirstWord "pony of the round table" }}
PluralizeFirstOfOneWord={{ pluralizeFirstWord "dwarf" }}
PluralizeFirstOfNoWord={{ pluralizeFirstWord "" }}
StripPackage={{ stripPackage "prefix.suffix" "xyz"}}
StripNoPackage={{ stripPackage "suffix" "xyz"}}
StripEmptyPackage={{ stripPackage "" "xyz" }}
DropPackage={{ dropPackage "prefix.suffix" }}
DropNoPackage={{ dropPackage "suffix" }}
DropEmptyPackage={{ dropPackage "" }}
ImportRuntime={{ contains .DefaultImports "github.com/go-openapi/runtime"}}
DoNotImport={{ contains .DefaultImports "github.com/go-openapi/xruntime"}}
PadSurround1={{ padSurround "padme" "-" 3 12}}
PadSurround2={{ padSurround "padme" "-" 0 12}}
Json={{ json .DefaultImports }}
PrettyJson={{ prettyjson . }}
`
)

func TestTemplates_CustomTemplates(t *testing.T) {

	var buf bytes.Buffer
	headerTempl, err := templates.Get("bindprimitiveparam")

	assert.Nil(t, err)

	err = headerTempl.Execute(&buf, nil)

	assert.Nil(t, err)
	assert.Equal(t, "\n", buf.String())

	buf.Reset()
	err = templates.AddFile("bindprimitiveparam", customHeader)

	assert.Nil(t, err)
	headerTempl, err = templates.Get("bindprimitiveparam")

	assert.Nil(t, err)

	err = headerTempl.Execute(&buf, nil)

	assert.Nil(t, err)
	assert.Equal(t, "custom header", buf.String())

}

func TestTemplates_CustomTemplatesMultiple(t *testing.T) {
	var buf bytes.Buffer

	err := templates.AddFile("differentFileName", customMultiple)

	assert.Nil(t, err)
	headerTempl, err := templates.Get("bindprimitiveparam")

	assert.Nil(t, err)

	err = headerTempl.Execute(&buf, nil)

	assert.Nil(t, err)
	assert.Equal(t, "custom primitive", buf.String())
}

func TestTemplates_CustomNewTemplates(t *testing.T) {
	var buf bytes.Buffer

	err := templates.AddFile("newtemplate", customNewTemplate)
	assert.Nil(t, err)

	err = templates.AddFile("existingUsesNew", customExistingUsesNew)
	assert.Nil(t, err)

	headerTempl, err := templates.Get("bindprimitiveparam")
	assert.Nil(t, err)

	err = headerTempl.Execute(&buf, nil)
	assert.Nil(t, err)

	assert.Equal(t, "new template", buf.String())
}

func TestTemplates_RepoLoadingTemplates(t *testing.T) {

	repo := NewRepository(nil)

	err := repo.AddFile("simple", singleTemplate)
	assert.NoError(t, err)

	templ, err := repo.Get("simple")

	assert.Nil(t, err)

	var b bytes.Buffer

	err = templ.Execute(&b, nil)

	assert.Nil(t, err)

	assert.Equal(t, "test", b.String())
}

func TestTemplates_RepoLoadsAllTemplatesDefined(t *testing.T) {

	var b bytes.Buffer
	repo := NewRepository(nil)

	err := repo.AddFile("multiple", multipleDefinitions)
	assert.NoError(t, err)

	templ, err := repo.Get("multiple")
	assert.Nil(t, err)
	err = templ.Execute(&b, nil)
	assert.Nil(t, err)

	assert.Equal(t, "", b.String())

	templ, err = repo.Get("T1")
	assert.Nil(t, err)
	err = templ.Execute(&b, nil)
	assert.Nil(t, err)

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
	assert.NoError(t, err)
	err = repo.AddFile("dependant", dependantTemplate)
	assert.NoError(t, err)

	templ, err := repo.Get("dependant")
	assert.Nil(t, err)

	err = templ.Execute(&b, nil)

	assert.Nil(t, err)

	assert.Equal(t, "T1D1", b.String())

}

func TestTemplates_RepoRecursiveTemplates(t *testing.T) {

	var b bytes.Buffer
	repo := NewRepository(nil)

	err := repo.AddFile("c1", cirularDeps1)
	assert.NoError(t, err)
	err = repo.AddFile("c2", cirularDeps2)
	assert.NoError(t, err)

	templ, err := repo.Get("c1")
	assert.Nil(t, err)
	data := testData{
		Name: "Root",
		Children: []testData{
			{Recurse: false},
		},
	}
	expected := `Root: Children`
	err = templ.Execute(&b, data)

	assert.Nil(t, err)

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

	assert.Nil(t, err)

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

	assert.Nil(t, err)

	assert.Equal(t, expected, b.String())
}

// Test that definitions are available to templates
// TODO: should test also with the codeGenApp context

// Test copyright definition
func TestTemplates_DefinitionCopyright(t *testing.T) {
	log.SetOutput(os.Stdout)

	repo := NewRepository(nil)

	err := repo.AddFile("copyright", copyright)
	assert.NoError(t, err)

	templ, err := repo.Get("copyright")
	assert.Nil(t, err)

	opts := opts()
	opts.Copyright = "My copyright clause"
	expected := opts.Copyright

	// executes template against model definitions
	genModel, err := getModelEnvironment("../fixtures/codegen/todolist.models.yml", opts)
	assert.Nil(t, err)

	rendered := bytes.NewBuffer(nil)
	err = templ.Execute(rendered, genModel)
	assert.Nil(t, err)

	assert.Equal(t, expected, rendered.String())

	// executes template against operations definitions
	genOperation, err := getOperationEnvironment("get", "/media/search", "../fixtures/codegen/instagram.yml", opts)
	assert.Nil(t, err)

	rendered.Reset()

	err = templ.Execute(rendered, genOperation)
	assert.Nil(t, err)

	assert.Equal(t, expected, rendered.String())

}

// Test TargetImportPath definition
func TestTemplates_DefinitionTargetImportPath(t *testing.T) {
	log.SetOutput(os.Stdout)

	repo := NewRepository(nil)

	err := repo.AddFile("targetimportpath", targetImportPath)
	assert.NoError(t, err)

	templ, err := repo.Get("targetimportpath")
	assert.Nil(t, err)

	opts := opts()
	// Non existing target would panic: to be tested too, but in another module
	opts.Target = "../fixtures"
	var expected = "github.com/go-swagger/go-swagger/fixtures"

	// executes template against model definitions
	genModel, err := getModelEnvironment("../fixtures/codegen/todolist.models.yml", opts)
	assert.Nil(t, err)

	rendered := bytes.NewBuffer(nil)
	err = templ.Execute(rendered, genModel)
	assert.Nil(t, err)

	assert.Equal(t, expected, rendered.String())

	// executes template against operations definitions
	genOperation, err := getOperationEnvironment("get", "/media/search", "../fixtures/codegen/instagram.yml", opts)
	assert.Nil(t, err)

	rendered.Reset()

	err = templ.Execute(rendered, genOperation)
	assert.Nil(t, err)

	assert.Equal(t, expected, rendered.String())

}

// Simulates a definition environment for model templates
func getModelEnvironment(spec string, opts *GenOpts) (*GenDefinition, error) {
	// Don't want stderr output to pollute CI
	log.SetOutput(ioutil.Discard)
	defer log.SetOutput(os.Stdout)

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
	// Don't want stderr output to pollute CI
	log.SetOutput(ioutil.Discard)
	defer log.SetOutput(os.Stdout)

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
// NOTE: We note that functions StripPackage() and DropPackage() behave the same way... and StripPackage()
// function is not sensitive to its second arg... Probably not what was intended in the first place but not
// blocking anyone for now.
func TestTemplates_FuncMap(t *testing.T) {
	log.SetOutput(os.Stdout)

	err := templates.AddFile("functpl", funcTpl)
	if assert.NoError(t, err) {
		templ, err := templates.Get("functpl")
		if assert.Nil(t, err) {
			opts := opts()
			// executes template against model definitions
			genModel, err := getModelEnvironment("../fixtures/codegen/todolist.models.yml", opts)
			if assert.Nil(t, err) {
				rendered := bytes.NewBuffer(nil)
				err = templ.Execute(rendered, genModel)
				if assert.Nil(t, err) {
					assert.Contains(t, rendered.String(), "Pascalize=WeArePoniesOfTheRoundTable\n")
					assert.Contains(t, rendered.String(), "Snakize=we_are_ponies_of_the_round_table\n")
					assert.Contains(t, rendered.String(), "Humanize=we are ponies of the round table\n")
					assert.Contains(t, rendered.String(), "PluralizeFirstWord=ponies of the round table\n")
					assert.Contains(t, rendered.String(), "PluralizeFirstOfOneWord=dwarves\n")
					assert.Contains(t, rendered.String(), "PluralizeFirstOfNoWord=\n")
					assert.Contains(t, rendered.String(), "StripPackage=suffix\n")
					assert.Contains(t, rendered.String(), "StripNoPackage=suffix\n")
					assert.Contains(t, rendered.String(), "StripEmptyPackage=\n")
					assert.Contains(t, rendered.String(), "DropPackage=suffix\n")
					assert.Contains(t, rendered.String(), "DropNoPackage=suffix\n")
					assert.Contains(t, rendered.String(), "DropEmptyPackage=\n")
					assert.Contains(t, rendered.String(), "DropEmptyPackage=\n")
					assert.Contains(t, rendered.String(), "ImportRuntime=true\n")
					assert.Contains(t, rendered.String(), "DoNotImport=false\n")
					assert.Contains(t, rendered.String(), "PadSurround1=-,-,-,padme,-,-,-,-,-,-,-,-\n")
					assert.Contains(t, rendered.String(), "PadSurround2=padme,-,-,-,-,-,-,-,-,-,-,-\n")
					assert.Contains(t, rendered.String(), "Json=[\"github.com/go-openapi/errors\",\"github.com/go-openapi/runtime\",\"github.com/go-openapi/swag\",\"github.com/go-openapi/validate\"]")
					assert.Contains(t, rendered.String(), "\"TargetImportPath\": \"github.com/go-swagger/go-swagger/generator\"")
					//fmt.Println(rendered.String())
				}
			}
		}
	}
}

// AddFile() global package function (protected vs unprotected)
// Mostly unused in tests, since the Repository.AddFile()
// is generally prefered.
func TestTemplates_AddFile(t *testing.T) {
	log.SetOutput(os.Stdout)

	// unprotected
	err := AddFile("functpl", funcTpl)
	if assert.NoError(t, err) {
		_, err := templates.Get("functpl")
		assert.Nil(t, err)
	}
	// protected
	err = AddFile("schemabody", funcTpl)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Cannot overwrite protected template")
}

// Test LoadDir
func TestTemplates_LoadDir(t *testing.T) {
	log.SetOutput(os.Stdout)

	// Fails
	err := templates.LoadDir("")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Could not complete")

	// Fails again (from any dir?)
	err = templates.LoadDir("templates")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Cannot overwrite protected template")

	// TODO: success case
	// To force a success, we need to empty the global list of protected
	// templates...
	origProtectedTemplates := protectedTemplates

	defer func() {
		// Restore variable initialized with package
		protectedTemplates = origProtectedTemplates
	}()

	protectedTemplates = make(map[string]bool)
	repo := NewRepository(FuncMap)
	err = repo.LoadDir("templates")
	assert.NoError(t, err)
}

// TODO: test error case in LoadDefaults()
// test DumpTemplates()
func TestTemplates_DumpTemplates(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	log.SetOutput(buf)
	defer func() {
		log.SetOutput(os.Stdout)
	}()

	templates.DumpTemplates()
	assert.NotEmpty(t, buf)
	// Sample output
	assert.Contains(t, buf.String(), "## tupleSerializer")
	assert.Contains(t, buf.String(), "Defined in `tupleserializer.gotmpl`")
	assert.Contains(t, buf.String(), "####requires \n - schemaType")
	//fmt.Println(buf)
}

// Go literal initializer func
func TestTemplates_GoSliceInitializer(t *testing.T) {
	a0 := []interface{}{"a", "b"}
	res, err := goSliceInitializer(a0)
	assert.NoError(t, err)
	assert.Equal(t, `{"a","b",}`, res)

	a1 := []interface{}{[]interface{}{"a", "b"}, []interface{}{"c", "d"}}
	res, err = goSliceInitializer(a1)
	assert.NoError(t, err)
	assert.Equal(t, `{{"a","b",},{"c","d",},}`, res)

	a2 := map[string]interface{}{"a": "y", "b": "z"}
	res, err = goSliceInitializer(a2)
	assert.NoError(t, err)
	assert.Equal(t, `{"a":"y","b":"z",}`, res)
}
