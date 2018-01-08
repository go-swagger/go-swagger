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

func TestClient_InvalidSpec(t *testing.T) {
	opts := testGenOpts()
	opts.Spec = "../fixtures/bugs/825/swagger.yml"
	opts.ValidateSpec = true
	assert.Error(t, GenerateClient("foo", nil, nil, &opts))
}

func TestClient_BaseImportDisabled(t *testing.T) {
	targetdir, err := ioutil.TempDir(os.TempDir(), "swagger_nogo")
	if err != nil {
		t.Fatalf("Failed to create a test target directory: %v", err)
	}
	log.SetOutput(ioutil.Discard)
	defer func() {
		os.RemoveAll(targetdir)
		log.SetOutput(os.Stdout)
	}()
	opts := testGenOpts()
	opts.Target = targetdir
	opts.Spec = "../fixtures/petstores/petstore.json"
	opts.LanguageOpts.BaseImportFunc = nil
	assert.NoError(t, GenerateClient("foo", nil, nil, &opts))
}
