package generator

import (
	"testing"

	"github.com/go-openapi/testify/v2/assert"
)

func TestRenameTimeout(t *testing.T) {
	for idx, toPin := range makeClientTimeoutNameTest() {
		i := idx
		testCase := toPin
		renameTimeout := rename(timeoutVarNamePreferences)
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			assert.EqualTf(t, testCase.expected, renameTimeout(testCase.seenIDs, testCase.name, 0), "unexpected deconflicting value [%d]", i)
		})
	}
}

func TestDeconflictTag(t *testing.T) {
	assert.EqualT(t, "runtimeops", deconflictTag(nil, "runtime"))
	assert.EqualT(t, "apiops", deconflictTag([]string{"tag1"}, apiPkg))
	assert.EqualT(t, "apiops1", deconflictTag([]string{"tag1", "apiops"}, apiPkg))
	assert.EqualT(t, "tlsops", deconflictTag([]string{"tag1"}, "tls"))
	assert.EqualT(t, "mytag", deconflictTag([]string{"tag1", "apiops"}, "mytag"))

	assert.EqualT(t, "operationsops", renameOperationPackage([]string{"tag1"}, "operations"))
	assert.EqualT(t, "operationsops11", renameOperationPackage([]string{"tag1", "operationsops1", "operationsops"}, "operations"))
}

func makeClientTimeoutNameTest() []struct {
	seenIDs  map[string]any
	name     string
	expected string
} {
	return []struct {
		seenIDs  map[string]any
		name     string
		expected string
	}{
		{
			seenIDs:  nil,
			name:     "witness",
			expected: "witness",
		},
		{
			seenIDs: map[string]any{
				"id": true,
			},
			name:     "timeout",
			expected: "timeout",
		},
		{
			seenIDs: map[string]any{
				"timeout":        true,
				"requesttimeout": true,
			},
			name:     "timeout",
			expected: "httpRequestTimeout",
		},
		{
			seenIDs: map[string]any{
				"timeout":            true,
				"requesttimeout":     true,
				"httprequesttimeout": true,
				"swaggertimeout":     true,
				"operationtimeout":   true,
				"optimeout":          true,
			},
			name:     "timeout",
			expected: "operTimeout",
		},
		{
			seenIDs: map[string]any{
				"timeout":            true,
				"requesttimeout":     true,
				"httprequesttimeout": true,
				"swaggertimeout":     true,
				"operationtimeout":   true,
				"optimeout":          true,
				"opertimeout":        true,
				"opertimeout1":       true,
			},
			name:     "timeout",
			expected: "operTimeout11",
		},
	}
}
