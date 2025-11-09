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
	"bytes"
	"io"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mutex for -race because this test alters a global.
var logMutex = &sync.Mutex{}

func TestDebugLog(t *testing.T) {
	// mutex for -race: we don't want several instances of this test to collide with each other
	logMutex.Lock()
	original := Debug

	defer func() {
		Debug = original
		logMutex.Unlock()
	}()

	tmp := t.TempDir()
	Debug = true

	t.Run("with debug log capture", func(t *testing.T) {
		tmpFile, _ := os.CreateTemp(tmp, "debug-test")
		tmpName := tmpFile.Name()

		debugOptions()
		defer func() {
			generatorLogger.SetOutput(os.Stdout)
		}()
		generatorLogger.SetOutput(tmpFile)

		t.Run("debug output with formatted args", func(t *testing.T) {
			debugLogf("A debug %v", map[string]any{"with arg": "arg", "token": "123"})
			err := tmpFile.Close()
			require.NoErrorf(t, err, "should flush the captured log, but got: %v", err)

			flushed, err := os.Open(tmpName)
			require.NoError(t, err)
			t.Cleanup(func() {
				_ = flushed.Close()
			})

			var buf bytes.Buffer
			_, err = io.Copy(&buf, flushed)
			require.NoError(t, err)
			str := buf.String()

			t.Run("log should contain format message", func(t *testing.T) {
				assert.Contains(t, str, "A debug")
			})
			t.Run("log should contain string argument", func(t *testing.T) {
				assert.Contains(t, str, "with arg")
			})
			t.Run("log should sanitize stuff like token, password...", func(t *testing.T) {
				assert.Contains(t, str, "***REDACTED**")
				assert.NotContains(t, str, "123")
			})
		})

		t.Run("debugAsJSON output with struct args", func(t *testing.T) {
			tmpJSONFile, _ := os.CreateTemp(tmp, "debug-as-json-test")
			tmpJSONName := tmpJSONFile.Name()
			generatorLogger.SetOutput(tmpJSONFile)
			debugLogAsJSONf("A short debug")

			sch := struct {
				FieldOne string `json:"fieldOne"`
			}{
				FieldOne: "content",
			}
			debugLogAsJSONf("A long debug:%t", true, sch)
			err := tmpJSONFile.Close()
			require.NoErrorf(t, err, "should flush the captured log, but got: %v", err)

			flushed, err := os.Open(tmpJSONName)
			require.NoError(t, err)
			t.Cleanup(func() {
				_ = flushed.Close()
			})

			var buf bytes.Buffer
			_, err = io.Copy(&buf, flushed)
			require.NoError(t, err)
			str := buf.String()

			t.Run("log should contain short message", func(t *testing.T) {
				assert.Contains(t, str, "A short debug")
			})
			t.Run("log should contain long message with args", func(t *testing.T) {
				assert.Contains(t, str, "A long debug:true")
			})
			t.Run("log should contain struct fields as JSON", func(t *testing.T) {
				assert.Contains(t, str, `"fieldOne":`)
				assert.Contains(t, str, `"content"`)
			})
		})
	})
}
