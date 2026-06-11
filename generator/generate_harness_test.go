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
	"github.com/go-swagger/go-swagger/generator/internal/language"

	"github.com/go-openapi/testify/v2/require"
)

// optionsBuilder builds a complete set of generation options for a harness
// fixture, given the resolved spec and the prepared target directory.
//
// It returns options ready to hand to a Generate* function (which performs the
// single Prepare). A builder may run its own setup sub-steps and assertions
// (hence the *testing.T) — e.g. relocating a test program into the target or
// generating an external models module.
//
// This single, one-shot construction replaces the former two-step dance
// (prepareTarget mutating the target, then a prepare hook mutating the opts):
// each fixture either relies on the harness default builder or provides its own.
type optionsBuilder func(t *testing.T, spec, target string) *GenOpts

// defaultServerOpts builds the standard server generation options for the harness.
func defaultServerOpts(t *testing.T, spec, target string) *GenOpts {
	t.Helper()

	g := NewGenOpts(ForServer(), WithSpec(spec), WithTarget(target))
	g.ExcludeSpec = true

	return g
}

// defaultClientOpts builds the standard client generation options for the harness.
func defaultClientOpts(t *testing.T, spec, target string) *GenOpts {
	t.Helper()

	return NewGenOpts(ForClient(), WithSpec(spec), WithTarget(target))
}

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

				spec := filepath.FromSlash(thisCas.spec)
				target := harnessTarget(t, thisName, "server_test", root)

				build := thisCas.prepare
				if build == nil {
					build = defaultServerOpts
				}
				opts := build(t, spec, target)

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

				spec := filepath.FromSlash(thisCas.spec)
				target := harnessTarget(t, thisName, "client_test", root)

				build := thisCas.prepare
				if build == nil {
					build = defaultClientOpts
				}
				opts := build(t, spec, target)

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
// It takes a swagger specification, optionally expects an error to occur when
// generating, an optional builder to construct the options (defaulting to the
// harness default builder) and an optional verify method to run extra assertions.
type generateFixture struct {
	spec      string
	wantError bool
	prepare   optionsBuilder
	verify    func(target string) func(*testing.T)
}

// harnessTarget computes and prepares (mkdir + go.mod init) a unique target
// directory for a fixture.
//
// It mangles the directory name with a standalone language mangler, so no
// generation options need to exist yet.
func harnessTarget(t *testing.T, name, base, root string) string {
	t.Helper()

	generated := filepath.Join(root, randWithPattern("generated"))
	target := filepath.Join(generated, language.GolangOpts().ManglePackageName(name, base))

	require.NoErrorf(t, os.MkdirAll(target, 0o700), "error in test creating target dir")
	t.Run("init module", gentest.GoModInit(target))

	return target
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
