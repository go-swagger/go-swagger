// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

// TODO: these definitions should be moved to a package so that the can be used
// by other tools.

// The file contains the structures used to define translations of a certain
// messages.
//
// A translation may have multiple translations strings, or messages, depending
// on the feature values of the various arguments. For instance, consider
// a hypothetical translation from English to English, where the source defines
// the format string "%d file(s) remaining".
// See the examples directory for examples of extracted messages.

// A Message describes a message to be translated.
type Message struct {
	// Key contains a list of identifiers for the message. If this list is empty
	// the message itself is used as the key.
	Key         []string `json:"key,omitempty"`
	Meaning     string   `json:"meaning,omitempty"`
	Message     Text     `json:"message"`
	Translation *Text    `json:"translation,omitempty"`

	Comment           string `json:"comment,omitempty"`
	TranslatorComment string `json:"translatorComment,omitempty"`

	Placeholders []Placeholder `json:"placeholders,omitempty"`

	// TODO: default placeholder syntax is {foo}. Allow alternative escaping
	// like `foo`.

	// Extraction information.
	Position string `json:"position,omitempty"` // filePosition:line
}

// A Placeholder is a part of the message that should not be changed by a
// translator. It can be used to hide or prettify format strings (e.g. %d or
// {{.Count}}), hide HTML, or mark common names that should not be translated.
type Placeholder struct {
	// ID is the placeholder identifier without the curly braces.
	ID string `json:"id"`

	// String is the string with which to replace the placeholder. This may be a
	// formatting string (for instance "%d" or "{{.Count}}") or a literal string
	// (<div>).
	String string `json:"string"`

	Type           string `json:"type"`
	UnderlyingType string `json:"underlyingType"`
	// ArgNum and Expr are set if the placeholder is a substitution of an
	// argument.
	ArgNum int    `json:"argNum,omitempty"`
	Expr   string `json:"expr,omitempty"`

	Comment string `json:"comment,omitempty"`
	Example string `json:"example,omitempty"`

	// Features contains the features that are available for the implementation
	// of this argument.
	Features []Feature `json:"features,omitempty"`
}

// An argument contains information about the arguments passed to a message.
type argument struct {
	// ArgNum corresponds to the number that should be used for explicit argument indexes (e.g.
	// "%[1]d").
	ArgNum int `json:"argNum,omitempty"`

	used           bool   // Used by Placeholder
	Type           string `json:"type"`
	UnderlyingType string `json:"underlyingType"`
	Expr           string `json:"expr"`
	Value          string `json:"value,omitempty"`
	Comment        string `json:"comment,omitempty"`
	Position       string `json:"position,omitempty"`
}

// Feature holds information about a feature that can be implemented by
// an Argument.
type Feature struct {
	Type string `json:"type"` // Right now this is only gender and plural.

	// TODO: possible values and examples for the language under consideration.

}

// Text defines a message to be displayed.
type Text struct {
	// Msg and Select contains the message to be displayed. Msg may be used as
	// a fallback value if none of the select cases match.
	Msg    string  `json:"msg,omitempty"`
	Select *Select `json:"select,omitempty"`

	// Var defines a map of variables that may be substituted in the selected
	// message.
	Var map[string]Text `json:"var,omitempty"`

	// Example contains an example message formatted with default values.
	Example string `json:"example,omitempty"`
}

// Select selects a Text based on the feature value associated with a feature of
// a certain argument.
type Select struct {
	Feature string          `json:"feature"` // Name of Feature type (e.g plural)
	Arg     string          `json:"arg"`     // The placeholder ID
	Cases   map[string]Text `json:"cases"`
}

// TODO: order matters, but can we derive the ordering from the case keys?
// type Case struct {
// 	Key   string `json:"key"`
// 	Value Text   `json:"value"`
// }
