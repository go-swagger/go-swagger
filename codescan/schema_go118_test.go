//go:build go1.18

package codescan

import (
	"testing"

	"github.com/go-openapi/spec"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	go118ClassificationCtx *scanCtx
)

func loadGo118ClassificationPkgsCtx(t testing.TB, extra ...string) *scanCtx {
	if go118ClassificationCtx != nil {
		return go118ClassificationCtx
	}
	sctx, err := newScanCtx(&Options{
		Packages: append([]string{
			"github.com/go-swagger/go-swagger/fixtures/goparsing/go118",
		}, extra...),
	})
	require.NoError(t, err)
	go118ClassificationCtx = sctx
	return go118ClassificationCtx
}

func getGo118ClassificationModel(sctx *scanCtx, nm string) *entityDecl {
	decl, ok := sctx.FindDecl("github.com/go-swagger/go-swagger/fixtures/goparsing/go118", nm)
	if !ok {
		return nil
	}
	return decl
}

func TestGo118SwaggerTypeNamed(t *testing.T) {
	sctx := loadGo118ClassificationPkgsCtx(t)
	decl := getGo118ClassificationModel(sctx, "NamedWithType")
	require.NotNil(t, decl)
	prs := &schemaBuilder{
		ctx:  sctx,
		decl: decl,
	}
	models := make(map[string]spec.Schema)
	require.NoError(t, prs.Build(models))
	schema := models["namedWithType"]

	assertProperty(t, &schema, "object", "some_map", "", "SomeMap")
}

func TestGo118AliasedModels(t *testing.T) {
	sctx := loadGo118ClassificationPkgsCtx(t)

	names := []string{
		"SomeObject",
	}

	defs := make(map[string]spec.Schema)
	for _, nm := range names {
		decl := getGo118ClassificationModel(sctx, nm)
		require.NotNil(t, decl)

		prs := &schemaBuilder{
			decl: decl,
			ctx:  sctx,
		}
		require.NoError(t, prs.Build(defs))
	}

	for k := range defs {
		for i, b := range names {
			if b == k {
				// remove the entry from the collection
				names = append(names[:i], names[i+1:]...)
			}
		}
	}
	if assert.Empty(t, names) {
		// map types
		assertMapDefinition(t, defs, "SomeObject", "object", "", "")
	}
}

func TestGo118InterfaceField(t *testing.T) {
	sctx := loadGo118ClassificationPkgsCtx(t)
	decl := getGo118ClassificationModel(sctx, "Interfaced")
	require.NotNil(t, decl)
	prs := &schemaBuilder{
		ctx:  sctx,
		decl: decl,
	}
	models := make(map[string]spec.Schema)
	require.NoError(t, prs.Build(models))

	schema := models["Interfaced"]
	assertProperty(t, &schema, "", "custom_data", "", "CustomData")
}

func TestGo118ParameterParser_Issue2011(t *testing.T) {
	sctx := loadGo118ClassificationPkgsCtx(t)
	operations := make(map[string]*spec.Operation)
	td := getParameter(sctx, "NumPlates")
	prs := &parameterBuilder{
		ctx:  sctx,
		decl: td,
	}
	require.NoError(t, prs.Build(operations))

	op := operations["putNumPlate"]
	require.NotNil(t, op)
	require.Len(t, op.Parameters, 1)
	sch := op.Parameters[0].Schema
	require.NotNil(t, sch)
}

func TestGo118ParseResponses_Issue2011(t *testing.T) {
	sctx := loadGo118ClassificationPkgsCtx(t)
	responses := make(map[string]spec.Response)
	td := getResponse(sctx, "NumPlatesResp")
	prs := &responseBuilder{
		ctx:  sctx,
		decl: td,
	}
	require.NoError(t, prs.Build(responses))

	resp := responses["NumPlatesResp"]
	require.Len(t, resp.Headers, 0)
	require.NotNil(t, resp.Schema)
}

func TestGo118_Issue2809(t *testing.T) {
	sctx := loadGo118ClassificationPkgsCtx(t)
	decl := getGo118ClassificationModel(sctx, "transportErr")
	require.NotNil(t, decl)
	prs := &schemaBuilder{
		ctx:  sctx,
		decl: decl,
	}
	models := make(map[string]spec.Schema)
	require.NoError(t, prs.Build(models))

	schema := models["transportErr"]
	assertProperty(t, &schema, "", "data", "", "Data")
}
