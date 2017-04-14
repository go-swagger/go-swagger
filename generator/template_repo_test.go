package generator

import (
	"bytes"
	"testing"

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
