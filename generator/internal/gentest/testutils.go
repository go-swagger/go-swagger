// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package gentest

import (
	"context"
	"fmt"
	"io"
	"log"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/go-openapi/testify/v2/require"
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

const minute = 60 * time.Second

// GoExecInDir executes a go commands from a target current directory.
//
// It returns a test runner func(*testing.T).
//
// Typical usage:
//
//	t.Run("should execute mycommand", gentest.GoExecInDir(folder, args))
func GoExecInDir(target string, args ...string) func(*testing.T) {
	return ExecInDir(target, "go", args...)
}

func ExecInDir(target string, command string, args ...string) func(*testing.T) {
	return func(t *testing.T) {
		t.Helper()

		ctx, cancel := context.WithTimeout(t.Context(), minute)
		defer cancel()

		cmd := exec.CommandContext(ctx, command, args...)
		cmd.Dir = target
		p, err := cmd.CombinedOutput()
		require.NoErrorf(t, err, "unexpected error: %s: %v\n%s", cmd.String(), err, string(p))
	}
}

var sanitizer = strings.NewReplacer(
	"(", "-",
	")", "-",
	".", "-",
	"_", "-",
	"\\", "/",
	":", "-",
	" ", "-",
)

func SanitizeGoModPath(pth string) string {
	return path.Clean(sanitizer.Replace(filepath.Base(pth)))
}

type GoModOption func(o *goModOptions)

type goModOptions struct {
	moduleName string
}

func WithGoModuleName(name string) GoModOption {
	return func(o *goModOptions) {
		o.moduleName = name
	}
}

func GoModInit(pth string, opts ...GoModOption) func(*testing.T) {
	var o goModOptions
	for _, apply := range opts {
		apply(&o)
	}

	if o.moduleName == "" {
		o.moduleName = SanitizeGoModPath(pth)
	}

	return func(t *testing.T) {
		t.Helper()

		t.Run(fmt.Sprintf("should initialize go.mod for %q", o.moduleName), func(t *testing.T) {
			ctx, cancel := context.WithTimeout(t.Context(), minute)
			defer cancel()

			mod := exec.CommandContext(ctx, "go", "mod", "init", o.moduleName) //nolint:gosec // "tainted" args exec is actually okay
			mod.Dir = pth
			output, err := mod.CombinedOutput()
			require.NoErrorf(t, err, "go mod init returned: %s", string(output))
		})
	}
}

func GoModTidy(pth string) func(*testing.T) {
	return func(t *testing.T) {
		t.Helper()

		t.Run("should tidy go.mod", func(t *testing.T) {
			ctx, cancel := context.WithTimeout(t.Context(), minute)
			defer cancel()

			vet := exec.CommandContext(ctx, "go", "mod", "tidy")
			vet.Dir = pth
			output, err := vet.CombinedOutput()
			require.NoError(t, err, string(output))
		})
	}
}

func GoBuild(pth string) func(*testing.T) {
	return func(t *testing.T) {
		t.Helper()

		t.Run("should build go", func(t *testing.T) {
			ctx, cancel := context.WithTimeout(t.Context(), minute)
			defer cancel()

			mod := exec.CommandContext(ctx, "go", "build")
			mod.Dir = pth
			output, err := mod.CombinedOutput()
			require.NoErrorf(t, err, "go build returned: %s", string(output))
		})
	}
}
