// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package generator

import (
	"github.com/go-swagger/go-swagger/generator/internal/language"
)

// Type aliases for backward compatibility within the generator package.
type (
	LanguageOpts  = language.Options
	FormatterFunc = language.FormatterFunc
	FormatOption  = language.FormatOption
)

var (

	// WithFormatLocalPrefixes adds local prefixes to group imports.
	WithFormatLocalPrefixes = language.WithFormatLocalPrefixes

	// WithFormatOnly tells the formatter to skip imports processing.
	WithFormatOnly = language.WithFormatOnly
)

// DefaultLanguageFunc defines the default generation language.
func DefaultLanguageFunc() *LanguageOpts {
	return language.GolangOpts()
}

// GolangOpts for rendering items as golang code.
func GolangOpts() *LanguageOpts {
	return language.GolangOpts()
}
