// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package generate

import (
	"io"
	"log"
	"os"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"

	"github.com/go-openapi/analysis"
)

func TestMain(m *testing.M) {
	// initializations to run tests in this package
	log.SetOutput(io.Discard)
	os.Exit(m.Run())
}

func Test_Shared_SetFlattenOptions(t *testing.T) {
	// testing multiple options settings for flatten:
	// - verbose | noverbose
	// - remove-unused
	// - expand
	// - minimal

	var fixt *FlattenCmdOptions

	res := fixt.SetFlattenOptions(nil)
	assert.NotNil(t, res)

	defaultOpts := resetDefaultOpts()

	res = fixt.SetFlattenOptions(defaultOpts)
	require.NotNil(t, res)
	assert.Equal(t, *defaultOpts, *res)

	fixt = &FlattenCmdOptions{
		WithExpand:  false,
		WithFlatten: []string{"noverbose"},
	}
	res = fixt.SetFlattenOptions(defaultOpts)
	assert.Equal(t, analysis.FlattenOpts{
		Verbose:      false,
		Minimal:      true,
		Expand:       false,
		RemoveUnused: false,
	}, *res)

	fixt = &FlattenCmdOptions{
		WithExpand:  false,
		WithFlatten: []string{"noverbose", "full"},
	}
	res = fixt.SetFlattenOptions(defaultOpts)
	assert.Equal(t, analysis.FlattenOpts{
		Verbose:      false,
		Minimal:      false,
		Expand:       false,
		RemoveUnused: false,
	}, *res)

	fixt = &FlattenCmdOptions{
		WithExpand:  false,
		WithFlatten: []string{"verbose", "noverbose", "full"},
	}
	res = fixt.SetFlattenOptions(defaultOpts)
	assert.Equal(t, analysis.FlattenOpts{
		Verbose:      true,
		Minimal:      false,
		Expand:       false,
		RemoveUnused: false,
	}, *res)

	fixt = &FlattenCmdOptions{
		WithExpand:  false,
		WithFlatten: []string{"verbose", "noverbose", "full", "expand", "remove-unused"},
	}
	res = fixt.SetFlattenOptions(defaultOpts)
	assert.Equal(t, analysis.FlattenOpts{
		Verbose:      true,
		Minimal:      false,
		Expand:       true,
		RemoveUnused: true,
	}, *res)

	fixt = &FlattenCmdOptions{
		WithExpand:  false,
		WithFlatten: []string{"minimal", "verbose", "noverbose", "full"},
	}
	res = fixt.SetFlattenOptions(defaultOpts)
	assert.Equal(t, analysis.FlattenOpts{
		Verbose:      true,
		Minimal:      true,
		Expand:       false,
		RemoveUnused: false,
	}, *res)

	fixt = &FlattenCmdOptions{
		WithExpand:  true,
		WithFlatten: []string{"minimal", "noverbose", "full"},
	}
	res = fixt.SetFlattenOptions(defaultOpts)
	assert.Equal(t, analysis.FlattenOpts{
		Verbose:      false,
		Minimal:      true,
		Expand:       true,
		RemoveUnused: false,
	}, *res)
}

func Test_Shared_ReadConfig(t *testing.T) {
	tmpFile, errio := os.CreateTemp(t.TempDir(), "tmp-config*.yaml")
	require.NoError(t, errio)
	tmpConfig := tmpFile.Name()
	require.NoError(t,
		os.WriteFile(tmpConfig, []byte(`param: 123
other: abc
`), 0o600))
	_ = tmpFile.Close()

	for _, toPin := range []struct {
		Name        string
		Filename    string
		ExpectError bool
		Expected    map[string]any
	}{
		{
			Name:        "empty",
			Filename:    "",
			ExpectError: true,
		},
		{
			Name:        "no file",
			Filename:    "nowhere",
			ExpectError: true,
		},
		{
			Name:     "happy path",
			Filename: tmpConfig,
			Expected: map[string]any{
				"param": 123,
				"other": "abc",
			},
		},
	} {
		testCase := toPin
		t.Run(testCase.Name, func(t *testing.T) {
			v, err := readConfig(testCase.Filename)
			if testCase.ExpectError {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			if testCase.Expected == nil {
				require.Nil(t, v)
				return
			}
			require.NotNil(t, v)
			m := v.AllSettings()
			for k, expectedValue := range testCase.Expected {
				require.MapContainsT(t, m, k)
				assert.Equal(t, expectedValue, m[k])
			}
		})
	}
}

func resetDefaultOpts() *analysis.FlattenOpts {
	return &analysis.FlattenOpts{
		Verbose:      true,
		Minimal:      true,
		Expand:       false,
		RemoveUnused: false,
	}
}
