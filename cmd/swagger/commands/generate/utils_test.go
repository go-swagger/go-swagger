// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package generate_test

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var enableGoVet bool

func init() {
	enableGoVet = os.Getenv("GOVET_TEST") != "" // enable go vet to run on generated CLI code - can't do this with args as the flags parser doesn't like it
}

func testBase() string {
	return filepath.FromSlash("../../../../")
}

const minute = 60 * time.Second

func gomodinit(pth string) func(*testing.T) {
	return func(t *testing.T) {
		t.Run("should initialize go.mod", func(t *testing.T) {
			ctx, cancel := context.WithTimeout(t.Context(), minute)
			defer cancel()

			mod := exec.CommandContext(ctx, "go", "mod", "init", filepath.FromSlash(filepath.Base(pth))) //nolint:gosec // "tainted" args exec is actually okay
			mod.Dir = pth
			output, err := mod.CombinedOutput()
			require.NoErrorf(t, err, "go mod init returned: %s", string(output))
		})
	}
}

func govet(pth string, expectError bool) func(*testing.T) {
	return func(t *testing.T) {
		t.Run("should pass go vet", func(t *testing.T) {
			ctx, cancel := context.WithTimeout(t.Context(), minute)
			defer cancel()

			vet := exec.CommandContext(ctx, "go", "vet", "./...")
			vet.Dir = pth
			output, err := vet.CombinedOutput()
			if expectError {
				require.Errorf(t, err, string(output))

				return
			}
			require.NoError(t, err, string(output))
		})
	}
}

func gomodtidy(pth string) func(*testing.T) {
	return func(t *testing.T) {
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

func gomoddownload(pth string) func(*testing.T) {
	return func(t *testing.T) {
		t.Run("should download dependencies", func(t *testing.T) {
			ctx, cancel := context.WithTimeout(t.Context(), minute)
			defer cancel()

			vet := exec.CommandContext(ctx, "go", "mod", "tidy")
			vet.Dir = pth
			output, err := vet.CombinedOutput()
			require.NoError(t, err, string(output))
		})
	}
}
