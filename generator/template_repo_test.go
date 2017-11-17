package generator

import (
	"bytes"
	"github.com/go-openapi/loads"
	"github.com/stretchr/testify/assert"
	"testing"
	//"fmt"
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
)

func TestCustomTemplates(t *testing.T) {

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

func TestCustomTemplatesMultiple(t *testing.T) {
	var buf bytes.Buffer

	err := templates.AddFile("differentFileName", customMultiple)

	assert.Nil(t, err)
	headerTempl, err := templates.Get("bindprimitiveparam")

	assert.Nil(t, err)

	err = headerTempl.Execute(&buf, nil)

	assert.Nil(t, err)
	assert.Equal(t, "custom primitive", buf.String())
}

func TestCustomNewTemplates(t *testing.T) {
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

func TestRepoLoadingTemplates(t *testing.T) {

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

func TestRepoLoadsAllTemplatesDefined(t *testing.T) {

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

func TestRepoLoadsAllDependantTemplates(t *testing.T) {

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

func TestRepoRecursiveTemplates(t *testing.T) {

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
// ..todo:: should test also with the codeGenApp context

// Test copyright definition
func TestDefinitionCopyright(t *testing.T) {

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
func TestDefinitionTargetImportPath(t *testing.T) {

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
