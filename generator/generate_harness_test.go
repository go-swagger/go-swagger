// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package generator

import (
	"bytes"
	"math/rand/v2"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/go-swagger/go-swagger/generator/internal/gentest"

	"github.com/stretchr/testify/require"
)

// testHarnessBuildServers iterates over a map of [generateFixture] s to prepare the code generation layout,
// generate a server from spec and run a verify method provided by fixture entry.
//
// Optionally, a filter may be added to execute only the filtered keys in the map (e.g. debug mode).
func testHarnessBuildServers(root string, fixtures map[string]generateFixture, filtered ...string) func(*testing.T) {
	return func(t *testing.T) {
		for name, cas := range filterMap(fixtures, filtered...) {
			thisCas := cas
			thisName := name

			t.Run(thisName, func(t *testing.T) {
				t.Parallel()

				defer func() {
					thisCas.warnFailed(t)
				}()

				opts := testGenOpts() // default opts

				t.Run("should prepare directory layout",
					thisCas.prepareTarget(thisName, "server_test", root, opts),
				)

				// preparation before generation, e.g. some test cases need to rework their input depending on the target location
				if thisCas.prepare != nil {
					t.Run("should prepare testcase", thisCas.prepare(opts))
				}

				t.Run("generating test server from "+opts.Spec, func(t *testing.T) {
					err := GenerateServer("", nil, nil, opts)
					if thisCas.wantError {
						require.Errorf(t, err, "expected an error for server build fixture: %s", opts.Spec)
					} else {
						require.NoError(t, err, "unexpected error for server build fixture: %s", opts.Spec)
					}

					if thisCas.verify != nil {
						t.Run("should verify results", thisCas.verify(opts.Target))
					}
				})
			})
		}
	}
}

func testHarnessBuildClients(root string, fixtures map[string]generateFixture, filtered ...string) func(*testing.T) {
	return func(t *testing.T) {
		for name, cas := range filterMap(fixtures, filtered...) {
			thisCas := cas
			thisName := name

			t.Run(thisName, func(t *testing.T) {
				t.Parallel()

				defer thisCas.warnFailed(t)
				opts := testClientGenOpts() // default opts for client codegen

				t.Run("should prepare directory layout",
					thisCas.prepareTarget(thisName, "client_test", root, opts),
				)

				// preparation before generation
				if thisCas.prepare != nil {
					t.Run("should prepare testcase", thisCas.prepare(opts))
				}

				t.Run("generating test client from "+opts.Spec, func(t *testing.T) {
					err := GenerateClient(thisName, nil, nil, opts)
					if thisCas.wantError {
						require.Errorf(t, err, "expected an error for client build fixture: %s", opts.Spec)
					} else {
						require.NoError(t, err, "unexpected error for client build fixture: %s", opts.Spec)
					}

					if thisCas.verify != nil {
						t.Run("should verify results", thisCas.verify(opts.Target))
					}
				})
			})
		}
	}
}

// generateFixture defines a test case for a code generation target.
//
// It takes a swagger specification, optionally expects an error to occur
// when generating, an optional additional preparation step and an
// optional verify method to run extra assertions.
type generateFixture struct {
	name string
	spec string
	// target    string
	wantError bool
	prepare   func(opts *GenOpts) func(*testing.T)
	verify    func(target string) func(*testing.T)
}

func (f generateFixture) base(root string) string {
	return filepath.Join(root, randWithPattern("generated"))
}

// prepareTarget prepares the directory layout necessary to generate code.
//
// It initializes a go.mod in the target directory.
func (f generateFixture) prepareTarget(name, base, root string, opts *GenOpts) func(*testing.T) {
	return func(t *testing.T) {
		if name == "" {
			name = f.name
		}

		spec := filepath.FromSlash(f.spec)
		opts.Spec = spec

		generated := f.base(root)
		opts.Target = filepath.Join(generated, opts.LanguageOpts.ManglePackageName(name, base))

		require.NoErrorf(t, os.MkdirAll(opts.Target, 0o700), "error in test creating target dir")
		t.Run("init module", gentest.GoModInit(opts.Target))
	}
}

func (f generateFixture) warnFailed(t *testing.T) {
	t.Helper()

	if t.Failed() {
		t.Log("ERROR: generation failed")
	}
}

// filterMap produces a map with only the keys provided as a filter.
//
// It returns the input map if no filter is specified.
func filterMap[Map ~map[K]V, K comparable, V any](m Map, keys ...K) Map {
	if len(keys) == 0 {
		return m
	}

	index := make(map[K]struct{}, len(keys))
	for _, k := range keys {
		index[k] = struct{}{}
	}

	filtered := make(Map, len(m))
	for k, v := range m {
		_, ok := index[k]
		if !ok {
			continue
		}
		filtered[k] = v
	}

	return filtered
}

// randWithPattern returns a random folder name similar to t.TempDir but does not create it.
func randWithPattern(pattern string) string {
	const limit uint64 = 1_000_000

	return pattern + "-" + strconv.FormatUint(rand.N(limit), 10) //nolint:gosec // OK: it's okay for tests not to use crypto/rand
}

// removeBuildTags strips "build ignore" tags in go code provided with test data.
func removeBuildTags(content []byte) []byte {
	content = bytes.ReplaceAll(content, []byte("//go:build ignore"), []byte(""))

	return bytes.ReplaceAll(content, []byte("// +build ignore"), []byte(""))
}
