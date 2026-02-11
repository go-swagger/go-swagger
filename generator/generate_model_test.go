// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package generator

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
)

func TestGenerateModels(t *testing.T) {
	t.Parallel()
	defer discardOutput()()

	root := t.TempDir()

	t.Run("generate models", func(t *testing.T) {
		for k, cas := range generateModelFixtures() {
			name := k
			thisCas := cas

			t.Run(name, func(t *testing.T) {
				t.Parallel()

				defer func() {
					thisCas.warnFailed(t)
				}()

				opts := testGenOpts()
				t.Run("should prepare directory layout",
					thisCas.prepareTarget(name, "model_test", root, opts),
				)

				if thisCas.prepare != nil {
					t.Run("should prepare testcase", thisCas.prepare(opts))
				}

				t.Run("generating test models from: "+opts.Spec, func(t *testing.T) {
					err := GenerateModels([]string{"", ""}, opts) // NOTE: generate all models, ignore ""
					if thisCas.wantError {
						require.Errorf(t, err, "expected an error for models build fixture: %s", opts.Spec)
					} else {
						require.NoError(t, err, "unexpected error for models build fixture: %s", opts.Spec)
					}

					if thisCas.verify != nil {
						t.Run("should verify results", thisCas.verify(opts.Target))
					}
				})
			})
		}
	})
}

func generateModelFixtures() map[string]generateFixture {
	return map[string]generateFixture{
		"allDefinitions": {
			spec: "../fixtures/bugs/1042/fixture-1042.yaml",
			// target: "../fixtures/bugs/1042",
			verify: func(target string) func(*testing.T) {
				return func(t *testing.T) {
					target = filepath.Join(target, defaultModelsTarget)
					require.TrueT(t, fileExists(target, ""))
					assert.TrueT(t, fileExists(target, "a.go"))
					assert.TrueT(t, fileExists(target, "b.go"))
				}
			},
		},
		"acceptDefinitions": {
			spec: "../fixtures/enhancements/2333/fixture-definitions.yaml",
			// target: "../fixtures/enhancements/2333",
			prepare: func(opts *GenOpts) func(*testing.T) {
				return func(_ *testing.T) {
					opts.AcceptDefinitionsOnly = true
				}
			},
			verify: func(target string) func(*testing.T) {
				return func(t *testing.T) {
					target = filepath.Join(target, defaultModelsTarget)
					require.TrueT(t, fileExists(target, ""))
					assert.TrueT(t, fileExists(target, "model_interface.go"))
					assert.TrueT(t, fileExists(target, "records_model.go"))
					assert.TrueT(t, fileExists(target, "records_model_with_max.go"))
					assert.FalseT(t, fileExists(target, "restapi"))
				}
			},
		},
		"mangleNames": {
			spec: "../fixtures/bugs/2821/ServiceManagementBody.json",
			// target: "../fixtures/bugs/2821",
			verify: func(target string) func(*testing.T) {
				return func(t *testing.T) {
					target = filepath.Join(target, defaultModelsTarget)
					require.TrueT(t, fileExists(target, "schema.go"))
					content, err := os.ReadFile(filepath.Join(target, "schema.go"))
					require.NoError(t, err)
					assert.StringContainsT(t, string(content), "getDollarRefField string")
				}
			},
		},
	}
}
