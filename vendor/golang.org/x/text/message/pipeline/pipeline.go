// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package pipeline provides tools for creating translation pipelines.
//
// NOTE: UNDER DEVELOPMENT. API MAY CHANGE.
package pipeline

import (
	"fmt"
	"go/build"
	"go/parser"
	"log"

	"golang.org/x/text/language"
	"golang.org/x/tools/go/loader"
)

const (
	extractFile  = "extracted.gotext.json"
	outFile      = "out.gotext.json"
	gotextSuffix = ".gotext.json"
)

// Config contains configuration for the translation pipeline.
type Config struct {
	// Supported indicates the languages for which data should be generated.
	// The default is to support all locales for which there are matching
	// translation files.
	Supported []language.Tag

	// --- Extraction

	SourceLanguage language.Tag

	Packages []string

	// --- File structure

	// Dir is the root dir for all operations.
	Dir string

	// TranslationsPattern is a regular expression for input translation files
	// that match anywhere in the directory structure rooted at Dir.
	TranslationsPattern string

	// OutPattern defines the location for translation files for a certain
	// language. The default is "{{.Dir}}/{{.Language}}/out.{{.Ext}}"
	OutPattern string

	// Format defines the file format for generated translation files.
	// The default is XMB. Alternatives are GetText, XLIFF, L20n, GoText.
	Format string

	Ext string

	// TODO:
	// Actions are additional actions to be performed after the initial extract
	// and merge.
	// Actions []struct {
	// 	Name    string
	// 	Options map[string]string
	// }

	// --- Generation

	// CatalogFile may be in a different package. It is not defined, it will
	// be written to stdout.
	CatalogFile string

	// DeclareVar defines a variable to which to assing the generated Catalog.
	DeclareVar string

	// SetDefault determines whether to assign the generated Catalog to
	// message.DefaultCatalog. The default for this is true if DeclareVar is
	// not defined, false otherwise.
	SetDefault bool

	// TODO:
	// - Printf-style configuration
	// - Template-style configuration
	// - Extraction options
	// - Rewrite options
	// - Generation options
}

// Operations:
// - extract:       get the strings
// - disambiguate:  find messages with the same key, but possible different meaning.
// - create out:    create a list of messages that need translations
// - load trans:    load the list of current translations
// - merge:         assign list of translations as done
// - (action)expand:    analyze features and create example sentences for each version.
// - (action)googletrans:   pre-populate messages with automatic translations.
// - (action)export:    send out messages somewhere non-standard
// - (action)import:    load messages from somewhere non-standard
// - vet program:   don't pass "foo" + var + "bar" strings. Not using funcs for translated strings.
// - vet trans:     coverage: all translations/ all features.
// - generate:      generate Go code

// State holds all accumulated information on translations during processing.
type State struct {
	Config Config

	Package string
	program *loader.Program

	Extracted Messages `json:"messages"`

	// Messages includes all messages for which there need to be translations.
	// Duplicates may be eliminated. Generation will be done from these messages
	// (usually after merging).
	Messages []Messages

	// Translations are incoming translations for the application messages.
	Translations []Messages
}

// A full-cycle pipeline example:
//
//  func updateAll(c *Config) error {
//  	s := Extract(c)
//  	s.Import()
//  	s.Merge()
//  	for range s.Config.Actions {
//  		if s.Err != nil {
//  			return s.Err
//  		}
//  		//  TODO: do the actions.
//  	}
//  	if err := s.Export(); err != nil {
//  		return err
//  	}
//  	if err := s.Generate(); err != nil {
//  		return err
//  	}
//  	return nil
//  }

// Import loads existing translation files.
func (s *State) Import() error {
	panic("unimplemented")
	return nil
}

// Merge merges the extracted messages with the existing translations.
func (s *State) Merge() error {
	panic("unimplemented")
	return nil

}

// Export writes out the messages to translation out files.
func (s *State) Export() error {
	panic("unimplemented")
	return nil
}

// NOTE: The command line tool already prefixes with "gotext:".
var (
	wrap = func(err error, msg string) error {
		return fmt.Errorf("%s: %v", msg, err)
	}
	wrapf = func(err error, msg string, args ...interface{}) error {
		return wrap(err, fmt.Sprintf(msg, args...))
	}
	errorf = fmt.Errorf
)

func warnf(format string, args ...interface{}) {
	// TODO: don't log.
	log.Printf(format, args...)
}

func loadPackages(conf *loader.Config, args []string) (*loader.Program, error) {
	if len(args) == 0 {
		args = []string{"."}
	}

	conf.Build = &build.Default
	conf.ParserMode = parser.ParseComments

	// Use the initial packages from the command line.
	args, err := conf.FromArgs(args, false)
	if err != nil {
		return nil, wrap(err, "loading packages failed")
	}

	// Load, parse and type-check the whole program.
	return conf.Load()
}
