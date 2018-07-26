// Copyright 2015 go-swagger maintainers
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package generator

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	targetdir, err := ioutil.TempDir(os.TempDir(), "swagger_nogo")
	if err != nil {
		t.Fatalf("Failed to create a test target directory: %v", err)
	}
	log.SetOutput(ioutil.Discard)
	defer func() {
		_ = os.RemoveAll(targetdir)
		log.SetOutput(os.Stdout)
	}()

	tests := []struct {
		name      string
		template  string
		wantError bool
		prepare   func(opts *GenOpts)
	}{
		{
			name:      "InvalidSpec",
			wantError: true,
			prepare: func(opts *GenOpts) {
				opts.Spec = invalidSpecExample
				opts.ValidateSpec = true
			},
		},
		{
			name:      "BaseImportDisabled",
			wantError: false,
		},
		{
			name:      "None_existing_contributor_tempalte",
			template:  "NonExistingContributorTemplate",
			wantError: true,
		},
		{
			name:      "Existing_contributor",
			template:  "stratoscale",
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := testGenOpts()
			opts.Target = targetdir
			opts.Spec = "../fixtures/petstores/petstore.json"
			opts.LanguageOpts.BaseImportFunc = nil
			opts.Template = tt.template

			if tt.prepare != nil {
				tt.prepare(&opts)
			}

			err := GenerateClient("foo", nil, nil, &opts)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
