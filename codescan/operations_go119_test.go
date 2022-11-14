package codescan

import (
	"testing"

	"github.com/go-openapi/spec"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIndentedYAMLBlock(t *testing.T) {
	sctx, err := newScanCtx(&Options{
		Packages: []string{
			"github.com/go-swagger/go-swagger/fixtures/goparsing/go119",
		},
	})
	require.NoError(t, err)
	var ops spec.Paths
	for _, apiPath := range sctx.app.Operations {
		prs := &operationsBuilder{
			ctx:        sctx,
			path:       apiPath,
			operations: make(map[string]*spec.Operation),
		}
		require.NoError(t, prs.Build(&ops))
	}

	assert.Len(t, ops.Paths, 2)

	po, ok := ops.Paths["/api/v1/somefunc"]
	assert.True(t, ok)
	assert.NotNil(t, po.Post)
	op := po.Post
	assert.Empty(t, op.Summary)
	assert.Equal(t, "Do something", op.Description)
	assert.Equal(t, "someFunc", op.ID)

	assert.Contains(t, op.Extensions, "x-codeSamples")

	samples := op.Extensions["x-codeSamples"].([]interface{})
	assert.Len(t, samples, 1)
	sample := samples[0].(map[string]interface{})
	assert.Contains(t, sample, "lang")
	assert.Equal(t, "curl", sample["lang"])

	assert.Contains(t, sample, "source")
	const expectedSource = `curl -u "${LOGIN}:${PASSWORD}" -d '{"key": "value"}' -X POST   "https://{host}/api/v1/somefunc"
curl -u "${LOGIN}:${PASSWORD}" -d '{"key2": "value2"}' -X POST   "https://{host}/api/v1/somefunc"
`
	assert.Equal(t, expectedSource, sample["source"])

	po2, ok := ops.Paths["/api/v1/somefuncTabs"]
	assert.True(t, ok)
	assert.NotNil(t, po2.Post)
	op2 := po2.Post
	assert.Empty(t, op2.Summary)
	assert.Equal(t, "Do something", op2.Description)
	assert.Equal(t, "someFuncTabs", op2.ID)

	assert.Contains(t, op2.Extensions, "x-codeSamples")

	samples2 := op2.Extensions["x-codeSamples"].([]interface{})
	assert.Len(t, samples2, 1)
	sample2 := samples2[0].(map[string]interface{})
	assert.Contains(t, sample2, "lang")
	assert.Equal(t, "curl", sample2["lang"])

	assert.Contains(t, sample2, "source")
	assert.Equal(t, expectedSource, sample2["source"])
}
