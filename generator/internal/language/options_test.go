// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package language

import (
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
)

func TestOptions_Init(t *testing.T) {
	opts := &Options{}
	assert.Empty(t, opts.BaseImport("x"))
	res, err := opts.FormatContent("x", []byte("y"))
	require.NoError(t, err)
	assert.Equal(t, []byte("y"), res)
	opts = GolangOpts()
	o := opts
	o.Init()
	assert.Equal(t, opts, o)
}

func TestOptions_MangleVarName(t *testing.T) {
	o := GolangOpts()

	// non-reserved word: returned as-is (after ToVarName)
	assert.EqualT(t, "myVar", o.MangleVarName("myVar"))

	// reserved word: gets "Var" suffix
	assert.EqualT(t, "breakVar", o.MangleVarName("break"))
	assert.EqualT(t, "selectVar", o.MangleVarName("select"))
}

func TestOptions_SetFormatFunc(t *testing.T) {
	o := &Options{}
	o.Init()

	// without formatFunc: returns content as-is
	res, err := o.FormatContent("test.go", []byte("hello"))
	require.NoError(t, err)
	assert.Equal(t, []byte("hello"), res)

	// with formatFunc: delegates
	o.SetFormatFunc(func(_ string, src []byte, _ ...FormatOption) ([]byte, error) {
		return []byte("formatted:" + string(src)), nil
	})
	res, err = o.FormatContent("test.go", []byte("hello"))
	require.NoError(t, err)
	assert.Equal(t, []byte("formatted:hello"), res)
}

func TestOptions_Imports_Nil(t *testing.T) {
	// nil ImportsFunc: returns ""
	o := &Options{}
	o.Init()
	assert.Empty(t, o.Imports(map[string]string{"fmt": "fmt"}))
}

func TestOptions_ArrayInitializer_Nil(t *testing.T) {
	// nil func: returns "", nil
	o := &Options{}
	o.Init()
	res, err := o.ArrayInitializer([]string{"a"})
	require.NoError(t, err)
	assert.Empty(t, res)
}

func TestOptions_BaseImport(t *testing.T) {
	// nil func: returns ""
	o := &Options{}
	o.Init()
	assert.Empty(t, o.BaseImport("anything"))

	// with custom func: delegates
	o.BaseImportFunc = func(s string) string { return "custom/" + s }
	assert.EqualT(t, "custom/target", o.BaseImport("target"))
}

func TestFormatOptions(t *testing.T) {
	// WithFormatLocalPrefixes
	opts := FormatOptsWithDefault([]FormatOption{
		WithFormatLocalPrefixes("github.com/myorg"),
	})
	assert.Equal(t, []string{"github.com/go-openapi", "github.com/myorg"}, opts.LocalPrefixes)

	// WithFormatOnly
	opts = FormatOptsWithDefault([]FormatOption{
		WithFormatOnly(true),
	})
	assert.TrueT(t, opts.FormatOnly)

	// defaults preserved when no options
	opts = FormatOptsWithDefault(nil)
	assert.EqualT(t, DefaultIndent, opts.TabWidth)
	assert.TrueT(t, opts.TabIndent)
	assert.TrueT(t, opts.Fragment)
	assert.TrueT(t, opts.Comments)
}
