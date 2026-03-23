// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package templatesrepo

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
)

func TestLoadDir_EmptyPath(t *testing.T) {
	repo := NewRepository(nil)

	err := repo.LoadDir("")
	require.Error(t, err)
	assert.StringContainsT(t, err.Error(), "could not complete")
}

func TestLoadDir_ProtectedTemplateBlocks(t *testing.T) {
	repo := NewRepository(nil)
	repo.SetProtectedTemplates(map[string]bool{
		"myProtected": true,
	})

	// Create a temp dir with a .gotmpl that defines a protected template
	dir := t.TempDir()
	err := os.WriteFile(
		filepath.Join(dir, "test.gotmpl"),
		[]byte(`{{ define "myProtected" }}hello{{ end }}`),
		0o600,
	)
	require.NoError(t, err)

	err = repo.LoadDir(dir)
	require.Error(t, err)
	assert.StringContainsT(t, err.Error(), "cannot overwrite protected template")
}

func TestLoadDir_Success(t *testing.T) {
	repo := NewRepository(nil)

	dir := t.TempDir()
	err := os.WriteFile(
		filepath.Join(dir, "greeting.gotmpl"),
		[]byte(`hello world`),
		0o600,
	)
	require.NoError(t, err)

	err = repo.LoadDir(dir)
	require.NoError(t, err)

	tmpl, err := repo.Get("greeting")
	require.NoError(t, err)
	require.NotNil(t, tmpl)
}

func TestSetAllowOverride(t *testing.T) {
	repo := NewRepository(nil)
	repo.SetProtectedTemplates(map[string]bool{
		"secret": true,
	})

	// Seed the repo with the protected template
	err := repo.addFile("secret.gotmpl", "original", true)
	require.NoError(t, err)

	// Without allowOverride, adding a file that redefines "secret" fails
	repo.SetAllowOverride(false)
	err = repo.AddFile("other.gotmpl", `{{ define "secret" }}replaced{{ end }}`)
	require.Error(t, err)
	assert.StringContainsT(t, err.Error(), "cannot overwrite protected template secret")

	// With allowOverride, it succeeds
	repo.SetAllowOverride(true)
	err = repo.AddFile("other.gotmpl", `{{ define "secret" }}replaced{{ end }}`)
	require.NoError(t, err)
}

func TestShallowClone(t *testing.T) {
	repo := NewRepository(nil)
	err := repo.AddFile("hello", "world")
	require.NoError(t, err)
	repo.SetProtectedTemplates(map[string]bool{"hello": true})
	repo.SetAllowOverride(true)

	clone := repo.ShallowClone()

	// clone has the same template
	tmpl, err := clone.Get("hello")
	require.NoError(t, err)
	require.NotNil(t, tmpl)

	// adding to clone doesn't affect original
	err = clone.AddFile("extra", "data")
	require.NoError(t, err)
	_, err = repo.Get("extra")
	require.Error(t, err)
}

func TestLoadDefaults(t *testing.T) {
	repo := NewRepository(nil)

	assets := map[string][]byte{
		"greeting.gotmpl": []byte("hello {{ . }}"),
		"farewell.gotmpl": []byte("goodbye {{ . }}"),
	}
	err := repo.LoadDefaults(assets)
	require.NoError(t, err)

	tmpl, err := repo.Get("greeting")
	require.NoError(t, err)
	require.NotNil(t, tmpl)

	tmpl, err = repo.Get("farewell")
	require.NoError(t, err)
	require.NotNil(t, tmpl)
}

func TestLoadDefaults_ParseError(t *testing.T) {
	repo := NewRepository(nil)

	assets := map[string][]byte{
		"bad.gotmpl": []byte("{{ .Broken"),
	}
	err := repo.LoadDefaults(assets)
	require.Error(t, err)
}

type mockAssetProvider struct {
	assets map[string][]byte
}

func (m mockAssetProvider) AssetNames() []string {
	names := make([]string, 0, len(m.assets))
	for k := range m.assets {
		names = append(names, k)
	}
	return names
}

func (m mockAssetProvider) MustAsset(name string) []byte {
	return m.assets[name]
}

func TestLoadContrib(t *testing.T) {
	repo := NewRepository(nil)
	provider := mockAssetProvider{
		assets: map[string][]byte{
			"templates/contrib/mycontrib/model.gotmpl":  []byte("model template"),
			"templates/contrib/mycontrib/server.gotmpl": []byte("server template"),
			"templates/contrib/other/skip.gotmpl":       []byte("should be skipped"),
			"templates/contrib/mycontrib/readme.md":     []byte("not a template"),
		},
	}

	err := repo.LoadContrib("mycontrib", provider)
	require.NoError(t, err)

	_, err = repo.Get("model")
	require.NoError(t, err)
	_, err = repo.Get("server")
	require.NoError(t, err)

	// "skip" from another contrib should not be loaded
	_, err = repo.Get("skip")
	require.Error(t, err)
}

func TestLoadContrib_NoFiles(t *testing.T) {
	repo := NewRepository(nil)
	provider := mockAssetProvider{
		assets: map[string][]byte{
			"templates/contrib/other/model.gotmpl": []byte("wrong contrib"),
		},
	}

	err := repo.LoadContrib("nonexistent", provider)
	require.Error(t, err)
	assert.StringContainsT(t, err.Error(), "no files added")
}

func TestMustGet(t *testing.T) {
	repo := NewRepository(nil)
	err := repo.AddFile("present", "content")
	require.NoError(t, err)

	// success
	tmpl := repo.MustGet("present")
	require.NotNil(t, tmpl)

	// panic on missing
	assert.Panics(t, func() {
		repo.MustGet("missing")
	})
}

func TestGet_NotFound(t *testing.T) {
	repo := NewRepository(nil)
	_, err := repo.Get("nonexistent")
	require.Error(t, err)
	assert.StringContainsT(t, err.Error(), "template doesn't exist")
}

func TestAddFile_ParseError(t *testing.T) {
	repo := NewRepository(nil)
	err := repo.AddFile("bad", "{{ .Broken")
	require.Error(t, err)
	assert.StringContainsT(t, err.Error(), "failed to load template")
}

func TestDumpTemplates(t *testing.T) {
	repo := NewRepository(nil)
	err := repo.AddFile("base", `base content {{ template "sub" }}`)
	require.NoError(t, err)
	err = repo.AddFile("sub", `sub content`)
	require.NoError(t, err)

	// should not panic
	repo.DumpTemplates()
}

func TestFuncs(t *testing.T) {
	fm := make(map[string]any)
	fm["myfunc"] = func() string { return "hi" }
	repo := NewRepository(fm)

	funcs := repo.Funcs()
	require.NotNil(t, funcs["myfunc"])
}

func TestFuncs_NilInit(t *testing.T) {
	repo := NewRepository(nil)
	funcs := repo.Funcs()
	require.NotNil(t, funcs)
}

func TestDependencies_TemplateNode(t *testing.T) {
	repo := NewRepository(nil)
	err := repo.AddFile("child", "child content")
	require.NoError(t, err)
	err = repo.AddFile("parent", `parent calls {{ template "child" }}`)
	require.NoError(t, err)

	tmpl, err := repo.Get("parent")
	require.NoError(t, err)

	// executing should work since dependency is resolved
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, nil)
	require.NoError(t, err)
	assert.StringContainsT(t, buf.String(), "child content")
}

func TestDependencies_MissingDep(t *testing.T) {
	repo := NewRepository(nil)
	err := repo.AddFile("orphan", `calls {{ template "missing" }}`)
	require.NoError(t, err)

	_, err = repo.Get("orphan")
	require.Error(t, err)
	assert.StringContainsT(t, err.Error(), "could not find template missing")
}

func TestDependencies_IfNode(t *testing.T) {
	repo := NewRepository(nil)
	err := repo.AddFile("ifchild", "if-child")
	require.NoError(t, err)
	err = repo.AddFile("elsechild", "else-child")
	require.NoError(t, err)
	err = repo.AddFile("iftpl", `{{ if . }}{{ template "ifchild" }}{{ else }}{{ template "elsechild" }}{{ end }}`)
	require.NoError(t, err)

	tmpl, err := repo.Get("iftpl")
	require.NoError(t, err)

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, true)
	require.NoError(t, err)
	assert.StringContainsT(t, buf.String(), "if-child")

	buf.Reset()
	err = tmpl.Execute(&buf, false)
	require.NoError(t, err)
	assert.StringContainsT(t, buf.String(), "else-child")
}

func TestDependencies_RangeNode(t *testing.T) {
	repo := NewRepository(nil)
	err := repo.AddFile("item", "item:{{ . }}")
	require.NoError(t, err)
	err = repo.AddFile("rangetpl", `{{ range . }}{{ template "item" . }}{{ end }}`)
	require.NoError(t, err)

	tmpl, err := repo.Get("rangetpl")
	require.NoError(t, err)

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, []string{"a", "b"})
	require.NoError(t, err)
	assert.StringContainsT(t, buf.String(), "item:a")
	assert.StringContainsT(t, buf.String(), "item:b")
}

func TestDependencies_WithNode(t *testing.T) {
	repo := NewRepository(nil)
	err := repo.AddFile("withchild", "with-child:{{ . }}")
	require.NoError(t, err)
	err = repo.AddFile("withtpl", `{{ with . }}{{ template "withchild" . }}{{ end }}`)
	require.NoError(t, err)

	tmpl, err := repo.Get("withtpl")
	require.NoError(t, err)

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, "data")
	require.NoError(t, err)
	assert.StringContainsT(t, buf.String(), "with-child:data")
}

func TestDependencies_Transitive(t *testing.T) {
	repo := NewRepository(nil)
	err := repo.AddFile("leaf", "leaf")
	require.NoError(t, err)
	err = repo.AddFile("mid", `mid->{{ template "leaf" }}`)
	require.NoError(t, err)
	err = repo.AddFile("root", `root->{{ template "mid" }}`)
	require.NoError(t, err)

	tmpl, err := repo.Get("root")
	require.NoError(t, err)

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, nil)
	require.NoError(t, err)
	assert.StringContainsT(t, buf.String(), "root->mid->leaf")
}

func TestFindDependencies_Nil(t *testing.T) {
	deps := findDependencies(nil)
	assert.Nil(t, deps)
}

func TestLoadContrib_ParseError(t *testing.T) {
	repo := NewRepository(nil)
	provider := mockAssetProvider{
		assets: map[string][]byte{
			"templates/contrib/bad/broken.gotmpl": []byte("{{ .Broken"),
		},
	}

	err := repo.LoadContrib("bad", provider)
	require.Error(t, err)
	assert.StringContainsT(t, err.Error(), "failed to load template")
}

func TestDependencies_RangeWithElse(t *testing.T) {
	repo := NewRepository(nil)
	err := repo.AddFile("rangeitem", "item")
	require.NoError(t, err)
	err = repo.AddFile("rangeempty", "empty")
	require.NoError(t, err)
	err = repo.AddFile("rangeelse", `{{ range . }}{{ template "rangeitem" }}{{ else }}{{ template "rangeempty" }}{{ end }}`)
	require.NoError(t, err)

	tmpl, err := repo.Get("rangeelse")
	require.NoError(t, err)

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, []string{})
	require.NoError(t, err)
	assert.StringContainsT(t, buf.String(), "empty")
}

func TestDependencies_WithElse(t *testing.T) {
	repo := NewRepository(nil)
	err := repo.AddFile("withpresent", "present")
	require.NoError(t, err)
	err = repo.AddFile("withnil", "nil")
	require.NoError(t, err)
	err = repo.AddFile("withelse", `{{ with . }}{{ template "withpresent" }}{{ else }}{{ template "withnil" }}{{ end }}`)
	require.NoError(t, err)

	tmpl, err := repo.Get("withelse")
	require.NoError(t, err)

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, nil)
	require.NoError(t, err)
	assert.StringContainsT(t, buf.String(), "nil")
}
