// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package codescan

import (
	"encoding/json"
	"go/ast"
	"testing"

	"github.com/go-openapi/testify/v2/require"
)

func TestYamlParser(t *testing.T) {
	t.Parallel()

	setter := func(out *string, called *int) func(json.RawMessage) error {
		return func(in json.RawMessage) error {
			*called++
			*out = string(in)

			return nil
		}
	}

	t.Run("with happy path", func(t *testing.T) {
		t.Run("should parse security definitions object as YAML", func(t *testing.T) {
			setterCalled := 0
			var actualJSON string
			parser := newYamlParser(rxSecurity, setter(&actualJSON, &setterCalled))

			lines := []string{
				"SecurityDefinitions:",
				"  api_key:",
				"    type: apiKey",
				"    name: X-API-KEY",
				"  petstore_auth:",
				"    type: oauth2",
				"    scopes:",
				"      'write:pets': modify pets in your account",
				"      'read:pets': read your pets",
			}

			require.TrueT(t, parser.Matches(lines[0]))
			require.NoError(t, parser.Parse(lines))
			require.EqualT(t, 1, setterCalled)

			const expectedJSON = `{"SecurityDefinitions":{"api_key":{"name":"X-API-KEY","type":"apiKey"},"petstore_auth":{"scopes":{"read:pets":"read your pets","write:pets":"modify pets in your account"},"type":"oauth2"}}}`

			require.JSONEqT(t, expectedJSON, actualJSON)
		})
	})

	t.Run("with edge cases", func(t *testing.T) {
		t.Run("should handle empty input", func(t *testing.T) {
			setterCalled := 0
			var actualJSON string
			parser := newYamlParser(rxSecurity, setter(&actualJSON, &setterCalled))

			require.FalseT(t, parser.Matches(""))
			require.NoError(t, parser.Parse([]string{}))
			require.Zero(t, setterCalled)
		})

		t.Run("should handle nil input", func(t *testing.T) {
			setterCalled := 0
			var actualJSON string
			parser := newYamlParser(rxSecurity, setter(&actualJSON, &setterCalled))

			require.NoError(t, parser.Parse(nil))
			require.Zero(t, setterCalled)
		})

		t.Run("should handle bad indentation", func(t *testing.T) {
			setterCalled := 0
			var actualJSON string
			parser := newYamlParser(rxSecurity, setter(&actualJSON, &setterCalled))
			lines := []string{
				"SecurityDefinitions:",
				"\t\tapi_key:",
				"  type: apiKey",
			}

			require.TrueT(t, parser.Matches(lines[0]))
			err := parser.Parse(lines)
			require.Error(t, err)
			require.StringContainsT(t, err.Error(), "yaml: line 2:")
			require.Zero(t, setterCalled)
		})

		t.Run("should catch YAML errors", func(t *testing.T) {
			setterCalled := 0
			var actualJSON string
			parser := newYamlParser(rxSecurity, setter(&actualJSON, &setterCalled))
			lines := []string{
				"SecurityDefinitions:",
				"  api_key",
				"    type: apiKey",
			}

			require.TrueT(t, parser.Matches(lines[0]))
			err := parser.Parse(lines)
			require.Error(t, err)
			require.StringContainsT(t, err.Error(), "yaml: line 3: mapping value")
			require.Zero(t, setterCalled)
		})
	})
}

func TestYamlSpecScanner(t *testing.T) {
	t.Parallel()

	t.Run("with happy path", func(t *testing.T) {
		t.Run("should parse operation definition object as YAML", func(t *testing.T) {
			parser := new(yamlSpecScanner)
			var title, description []string
			parser.setTitle = func(lines []string) { title = lines }
			parser.setDescription = func(lines []string) { description = lines }

			lines := []string{
				// from issue #3225, reindented
				// `swagger:operation POST /v1/example-endpoint addExampleConfig`,
				`title for this operation`,
				``, // blank line elided
				`description of this operation`,
				``, // blank line preserved
				`continuation of the description`,
				`---`, // YAML block
				`summary: Adds a new configuration entry`,
				`description: |-`,
				`  Creates and validates a new configuration request.`,
				``,
				`security:`,
				`- AuthToken: []`,
				`consumes:`,
				`- application/json`,
				`tags:`,
				`- Example|Configuration`,
				`responses:`,
				`  201:`,
				`    $ref: "#/responses/createdResponse"`,
				`  400:`,
				`    $ref: "#/responses/badRequestResponse"`,
				`  412:`,
				`    $ref: "#/responses/preconditionFailedResponse"`,
				`  500:`,
				`    $ref: "#/responses/internalServerErrorResponse"`,
			}

			doc := buildRawTestComments(lines)
			require.NoError(t, parser.Parse(doc))
			require.Equal(t, title, parser.Title())
			require.Equal(t, []string{"title for this operation"}, parser.Title())
			require.Equal(t, description, parser.Description())
			require.Equal(t, []string{"description of this operation", "", "continuation of the description"}, parser.Description())

			var receivedJSON string
			yamlReceiver := func(b []byte) error {
				receivedJSON = string(b)
				return nil
			}

			require.NoError(t, parser.UnmarshalSpec(yamlReceiver))

			const expectedJSON = `{
				"summary":"Adds a new configuration entry",
				"description":"Creates and validates a new configuration request.",
				"security":[
					{"AuthToken":[]}
				],
				"consumes":["application/json"],
				"tags":["Example|Configuration"],
				"responses":{
					"201":{"$ref":"#/responses/createdResponse"},
					"400":{"$ref":"#/responses/badRequestResponse"},
					"412":{"$ref":"#/responses/preconditionFailedResponse"},
					"500":{"$ref":"#/responses/internalServerErrorResponse"}
				}
			}`

			require.JSONEqT(t, expectedJSON, receivedJSON)
		})

		t.Run("should stop yaml operation block when new tag is found", func(t *testing.T) {
			parser := new(yamlSpecScanner)
			var title, description []string
			parser.setTitle = func(lines []string) { title = lines }
			parser.setDescription = func(lines []string) { description = lines }

			lines := []string{
				`title for this operation`,
				``, // blank line elided
				`description of this operation`,
				`---`, // YAML block
				`summary: Adds a new configuration entry`,
				``,
				`swagger:enum`, // yaml block ended at this tag. Rest is ignored
				`security:`,
				`- AuthToken: []`,
			}

			doc := buildRawTestComments(lines)
			require.NoError(t, parser.Parse(doc))
			require.Equal(t, title, parser.Title())
			require.Equal(t, []string{"title for this operation"}, parser.Title())
			require.Equal(t, description, parser.Description())
			require.Equal(t, []string{"description of this operation"}, parser.Description())

			var receivedJSON string
			yamlReceiver := func(b []byte) error {
				receivedJSON = string(b)
				return nil
			}

			require.NoError(t, parser.UnmarshalSpec(yamlReceiver))

			const expectedJSON = `{
				"summary":"Adds a new configuration entry"
			}`

			require.JSONEqT(t, expectedJSON, receivedJSON)
		})

		t.Run("should stop yaml operation block when new yaml document separator is found", func(t *testing.T) {
			parser := new(yamlSpecScanner)
			var title, description []string
			parser.setTitle = func(lines []string) { title = lines }
			parser.setDescription = func(lines []string) { description = lines }

			lines := []string{
				`title for this operation`,
				``, // blank line elided
				`description of this operation`,
				`---`, // YAML block
				`summary: Adds a new configuration entry`,
				``,
				`---`, // yaml block ended at mark. Rest is ignored
				`security:`,
				`- AuthToken: []`,
			}

			doc := buildRawTestComments(lines)
			require.NoError(t, parser.Parse(doc))
			require.Equal(t, title, parser.Title())
			require.Equal(t, []string{"title for this operation"}, parser.Title())
			require.Equal(t, description, parser.Description())
			require.Equal(t, []string{"description of this operation"}, parser.Description())

			var receivedJSON string
			yamlReceiver := func(b []byte) error {
				receivedJSON = string(b)
				return nil
			}

			require.NoError(t, parser.UnmarshalSpec(yamlReceiver))

			const expectedJSON = `{
				"summary":"Adds a new configuration entry"
			}`

			require.JSONEqT(t, expectedJSON, receivedJSON)
		})
	})

	t.Run("with edge cases", func(t *testing.T) {
		t.Run("with empty comment block", func(t *testing.T) {
			parser := new(yamlSpecScanner)
			var title, description []string
			parser.setTitle = func(lines []string) { title = lines }
			parser.setDescription = func(lines []string) { description = lines }
			doc := buildRawTestComments(nil)
			require.NoError(t, parser.Parse(doc))
			require.Empty(t, title)
			require.Empty(t, description)
		})

		t.Run("with nil comment block", func(t *testing.T) {
			parser := new(yamlSpecScanner)
			var title, description []string
			parser.setTitle = func(lines []string) { title = lines }
			parser.setDescription = func(lines []string) { description = lines }
			require.NoError(t, parser.Parse(nil))
			require.Empty(t, title)
			require.Empty(t, description)
		})

		t.Run("without setTitle", func(t *testing.T) {
			parser := new(yamlSpecScanner)
			var description []string
			parser.setDescription = func(lines []string) { description = lines }

			lines := []string{
				`title for this operation`,
				``, // blank line preserved
				`description of this operation`,
				`---`, // YAML block
			}

			doc := buildRawTestComments(lines)
			require.NoError(t, parser.Parse(doc))
			require.Nil(t, parser.Title())
			require.Equal(t, description, parser.Description())
			require.Equal(t, []string{"title for this operation", "", "description of this operation"}, parser.Description())

			var receivedJSON string
			yamlReceiver := func(b []byte) error {
				receivedJSON = string(b)
				return nil
			}
			require.NoError(t, parser.UnmarshalSpec(yamlReceiver))
			require.JSONEqT(t, `{}`, receivedJSON)
		})
	})
}

func TestRemoveIndent(t *testing.T) {
	t.Parallel()

	t.Run("with removeIndent", func(t *testing.T) {
		t.Run("should tolerate empty input", func(t *testing.T) {
			res := removeIndent([]string{})
			require.Empty(t, res)
			require.NotNil(t, res)
		})

		t.Run("should tolerate nil input", func(t *testing.T) {
			res := removeIndent(nil)
			require.Empty(t, res)
			require.Nil(t, res)
		})

		t.Run("should support headline without indentation", func(t *testing.T) {
			lines := []string{
				"xyz",
				"  abc",
			}
			res := removeIndent(lines)
			require.Equal(t, lines, res)
		})

		t.Run("should tolerate lines with only indents", func(t *testing.T) {
			lines := []string{
				"  xyz",
				"",
				"    ",
				"    ",
			}
			res := removeIndent(lines)

			expected := []string{
				"xyz",
				"",   // empty line preserved
				"  ", // blank lines unindented
				"  ",
			}
			require.Equal(t, expected, res)
		})

		t.Run("should replace tabs with spaces in indentation", func(t *testing.T) {
			lines := []string{
				"\t\txyz",
				"",
				"    ",
				"\t  \t",
			}
			res := removeIndent(lines)

			expected := []string{
				"xyz",
				"",   // empty line preserved
				"  ", // blank lines unindented
				" \t",
			}
			require.Equal(t, expected, res)
		})
	})

	t.Run("with removeYamlIndent", func(t *testing.T) {
		t.Run("should tolerate empty input", func(t *testing.T) {
			res := removeYamlIndent([]string{})
			require.Empty(t, res)
			require.NotNil(t, res)
		})

		t.Run("should tolerate nil input", func(t *testing.T) {
			res := removeYamlIndent(nil)
			require.Empty(t, res)
			require.Nil(t, res)
		})

		t.Run("should support headline without indentation", func(t *testing.T) {
			lines := []string{
				"xyz",
				"  abc",
			}
			res := removeYamlIndent(lines)
			require.Equal(t, lines, res)
		})

		t.Run("should support headline without indentation", func(t *testing.T) {
			lines := []string{
				"xyz",
				"  abc",
			}
			res := removeYamlIndent(lines)
			require.Equal(t, lines, res)
		})
	})
}

func buildRawTestComments(lines []string) *ast.CommentGroup {
	// build raw doc comments like ast provides
	doc := &ast.CommentGroup{
		List: make([]*ast.Comment, 0, len(lines)),
	}
	for _, line := range lines {
		doc.List = append(doc.List, &ast.Comment{Text: "// " + line})
	}

	return doc
}
