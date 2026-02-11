// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package generator

import (
	"bytes"
	"io"
	"os"
	"sync"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
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
				assert.StringContainsT(t, str, "A debug")
			})
			t.Run("log should contain string argument", func(t *testing.T) {
				assert.StringContainsT(t, str, "with arg")
			})
			t.Run("log should sanitize stuff like token, password...", func(t *testing.T) {
				assert.StringContainsT(t, str, "***REDACTED**")
				assert.StringNotContainsT(t, str, "123")
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
				assert.StringContainsT(t, str, "A short debug")
			})
			t.Run("log should contain long message with args", func(t *testing.T) {
				assert.StringContainsT(t, str, "A long debug:true")
			})
			t.Run("log should contain struct fields as JSON", func(t *testing.T) {
				assert.StringContainsT(t, str, `"fieldOne":`)
				assert.StringContainsT(t, str, `"content"`)
			})
		})
	})
}
