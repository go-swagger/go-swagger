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

package gentest

import (
	"context"
	"io"
	"log"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var lockLogger sync.Mutex

// DiscardOutput discards the standard logger and returns
// a rollback function.
//
// Typical usage:
//
//	defer gentest.DiscardOutput()()
func DiscardOutput() func() {
	return setOutput(io.Discard)
}

// CaptureOutput captures the standard logger to the passed writer
// and returns a rollback function.
// Typical usage:
//
//	var buf bytes.Buffer
//	defer gentest.CaptureOutput(&buf)()
func CaptureOutput(w io.Writer) func() {
	return setOutput(w)
}

func setOutput(w io.Writer) func() {
	lockLogger.Lock()
	defer lockLogger.Unlock()

	original := log.Writer()
	// discards log output then sends a function to set it back to its original value
	log.SetOutput(w)

	return func() {
		lockLogger.Lock()
		log.SetOutput(original)
		lockLogger.Unlock()
	}
}

// GoExecInDir executes a go commands from a target current directory.
//
// It returns a test runner func(*testing.T).
//
// Typical usage:
//
//	t.Run("should execute mycommand", gentest.GoExecInDir(folder, args))
func GoExecInDir(target string, args ...string) func(*testing.T) {
	return func(t *testing.T) {
		t.Helper()

		const minute = 60 * time.Second
		ctx, cancel := context.WithTimeout(t.Context(), minute)
		defer cancel()

		cmd := exec.CommandContext(ctx, "go", args...)
		cmd.Dir = target
		p, err := cmd.CombinedOutput()
		require.NoErrorf(t, err, "unexpected error: %s: %v\n%s", cmd.String(), err, string(p))
	}
}

var sanitizer = strings.NewReplacer("(", "-", ")", "-", ".", "-", "_", "-", "\\", "/", ":", "-", " ", "x")

func SanitizeGoModPath(pth string) string {
	return path.Clean(sanitizer.Replace(filepath.Base(pth)))
}
