// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package generator

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/go-swagger/go-swagger/generator/internal/gentest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type clientGenerateFixture struct {
	name      string
	spec      string
	template  string
	wantError bool
	prepare   func(opts *GenOpts)
	verify    func(*testing.T, string)
}

func clientFixtures() []clientGenerateFixture {
	return []clientGenerateFixture{
		{
			name:      "InvalidSpec",
			wantError: true,
			prepare: func(opts *GenOpts) {
				opts.Spec = invalidSpecExample
				opts.ValidateSpec = true
			},
		},
		{
			name: "BaseImportDisabled",
			prepare: func(opts *GenOpts) {
				opts.LanguageOpts.BaseImportFunc = nil
			},
			wantError: false,
		},
		{
			name:      "Non_existing_contributor_template",
			template:  "NonExistingContributorTemplate",
			wantError: true,
		},
		{
			name:      "Existing_contributor",
			template:  "stratoscale",
			wantError: false,
		},
		{
			name:      "packages mangling",
			wantError: false,
			spec:      filepath.Join("..", "fixtures", "bugs", "2111", "fixture-2111.yaml"),
			verify: func(t *testing.T, target string) {
				require.True(t, fileExists(target, "client"))

				// assert package generation based on mangled tags
				target = filepath.Join(target, "client")
				assert.True(t, fileExists(target, "abc_linux"))
				assert.True(t, fileExists(target, "abc_test"))
				assert.True(t, fileExists(target, apiPkg))
				assert.True(t, fileExists(target, "custom"))
				assert.True(t, fileExists(target, "hash_tag_donuts"))
				assert.True(t, fileExists(target, "nr123abc"))
				assert.True(t, fileExists(target, "nr_at_donuts"))
				assert.True(t, fileExists(target, "operations"))
				assert.True(t, fileExists(target, "plus_donuts"))
				assert.True(t, fileExists(target, "strfmt"))
				assert.True(t, fileExists(target, "forced"))
				assert.True(t, fileExists(target, "gtl"))
				assert.True(t, fileExists(target, "nr12nasty"))
				assert.True(t, fileExists(target, "override"))
				assert.True(t, fileExists(target, "operationsops"))

				buf, err := os.ReadFile(filepath.Join(target, "foo_client.go"))
				require.NoError(t, err)

				// assert client import, with deconfliction
				code := string(buf)
				importBase := gentest.SanitizeGoModPath(filepath.Dir(filepath.Dir(target)))
				importRegexp := importBase + `/packages_mangling/client`
				assertImports(t, importRegexp, code)

				assertInCode(t, `cli.Strfmt = strfmtops.New(transport, formats)`, code)
				assertInCode(t, `cli.API = apiops.New(transport, formats)`, code)
				assertInCode(t, `cli.Operations = operations.New(transport, formats)`, code)
			},
		},
		{
			name:      "packages flattening",
			wantError: false,
			spec:      filepath.Join("..", "fixtures", "bugs", "2111", "fixture-2111.yaml"),
			prepare: func(opts *GenOpts) {
				opts.SkipTagPackages = true
			},
			verify: func(t *testing.T, target string) {
				require.True(t, fileExists(target, "client"))

				// packages are not created here
				target = filepath.Join(target, "client")
				assert.False(t, fileExists(target, "abc_linux"))
				assert.False(t, fileExists(target, "abc_test"))
				assert.False(t, fileExists(target, apiPkg))
				assert.False(t, fileExists(target, "custom"))
				assert.False(t, fileExists(target, "hash_tag_donuts"))
				assert.False(t, fileExists(target, "nr123abc"))
				assert.False(t, fileExists(target, "nr_at_donuts"))
				assert.False(t, fileExists(target, "plus_donuts"))
				assert.False(t, fileExists(target, "strfmt"))
				assert.False(t, fileExists(target, "forced"))
				assert.False(t, fileExists(target, "gtl"))
				assert.False(t, fileExists(target, "nr12nasty"))
				assert.False(t, fileExists(target, "override"))
				assert.False(t, fileExists(target, "operationsops"))

				assert.True(t, fileExists(target, "operations"))
			},
		},
		{
			name:      "name with trailing API",
			spec:      filepath.Join("..", "fixtures", "bugs", "2278", "fixture-2278.yaml"),
			wantError: false,
		},
	}
}
