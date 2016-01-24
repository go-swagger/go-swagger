package generator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegistryOnlyCompilesOnce(t *testing.T) {
	registry := NewTemplateRegistry()

	asset := []byte(`{{ define "test" }}{{ end }}`)

	template := TemplateDefinition{
		Files: []string{"test.gotmpl"},
	}
	registry.AddFile("test.gotmpl", asset)

	registry.AddTemplate("test", template)

	compiled := registry.MustGet("test")
	compiled2 := registry.MustGet("test")

	assert.Equal(t, compiled, compiled2)
}

func TestRegistryRecompilesIfAssetsChange(t *testing.T) {
	registry := NewTemplateRegistry()

	asset := []byte(`{{ define "test" }}{{ end }}`)
	asset2 := []byte(`{{ define "test2" }}{{ end }}`)
	template := TemplateDefinition{
		Files: []string{"test.gotmpl"},
	}
	registry.AddFile("test.gotmpl", asset)

	registry.AddTemplate("test", template)

	compiled := registry.MustGet("test")

	registry.AddFile("test.gotmpl", asset2)

	compiled2 := registry.MustGet("test")

	assert.NotEqual(t, compiled, compiled2)

}
