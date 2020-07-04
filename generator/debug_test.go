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
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	logMutex = &sync.Mutex{}
)

func TestDebug(t *testing.T) {
	// test debugLog()
	tmpFile, _ := ioutil.TempFile("", "debug-test")
	tmpName := tmpFile.Name()
	logMutex.Lock()
	defer func() {
		Debug = false
		// mutex for -race
		logMutex.Unlock()
		_ = os.Remove(tmpName)
	}()

	// mutex for -race
	Debug = true
	debugOptions()
	defer func() {
		generatorLogger.SetOutput(os.Stdout)
	}()
	generatorLogger.SetOutput(tmpFile)

	debugLog("A debug")
	_ = tmpFile.Close()

	flushed, _ := os.Open(tmpName)
	buf := make([]byte, 500)
	_, _ = flushed.Read(buf)
	assert.Contains(t, string(buf), "A debug")
	_ = flushed.Close()

	// test debugLogAsJSON()
	tmpJSONFile, _ := ioutil.TempFile("", "debug-test")
	tmpJSONName := tmpJSONFile.Name()
	defer func() {
		_ = os.Remove(tmpJSONName)
	}()
	generatorLogger.SetOutput(tmpJSONFile)
	debugLogAsJSON("A short debug")

	sch := struct {
		FieldOne string `json:"fieldOne"`
	}{
		FieldOne: "content",
	}
	debugLogAsJSON("A long debug:%t", true, sch)
	_ = tmpJSONFile.Close()

	flushed, _ = os.Open(tmpJSONName)
	buf2 := make([]byte, 500)
	_, _ = flushed.Read(buf2)
	_ = flushed.Close()
	assert.Contains(t, string(buf2), "A short debug")
	assert.Contains(t, string(buf2), "A long debug:true")
	assert.Contains(t, string(buf2), `"fieldOne":`)
	assert.Contains(t, string(buf2), `"content"`)
}

func TestDebugAsJSON(t *testing.T) {
	t.SkipNow()

	var body struct {
		A string `json:"a"`
		B int    `json:"b"`
	}
	Debug = true
	body.A = "abc"
	body.B = 123
	debugLogAsJSON("No arg")
	debugLogAsJSON("With arg:%t", true, body)
}
