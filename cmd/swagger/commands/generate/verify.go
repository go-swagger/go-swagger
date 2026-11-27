// SPDX-FileCopyrightText: Copyright 2026 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package generate

import (
	"bytes"
	"errors"
	"fmt"
	"os"

	"github.com/go-swagger/go-swagger/generator"
)

// Verify checks the project's generated server against a spec.
type Verify struct {
	Server
}

// Usage shows that the spec can be supplied as a positional argument.
func (Verify) Usage() string {
	return usageWithSpec("verify")
}

// Execute compares the generated server files with their on-disk contents.
func (v *Verify) Execute(args []string) error {
	if len(args) == 0 && v.Shared.Spec == "" {
		return errors.New("verify requires exactly one swagger document as a positional argument or with --spec")
	}

	if v.Shared.DumpData {
		return errors.New("verify cannot be combined with --dump-data")
	}
	return createSwagger(v, args)
}

func (v *Verify) generate(opts *generator.GenOpts) error {
	var differences []error
	opts.VerifyFile = func(path string, generated []byte) {
		current, err := os.ReadFile(path)
		if err != nil {
			differences = append(differences, err)
		} else if !bytes.Equal(current, generated) {
			differences = append(differences, fmt.Errorf("generated file %q is out of date", path))
		}
	}

	err := v.Server.generate(opts)
	return errors.Join(append(differences, err)...)
}

// Verification has no post-generation instructions.
func (*Verify) log(string) {}
