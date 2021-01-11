package generator

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"testing"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"
	"github.com/go-openapi/swag"
	"github.com/stretchr/testify/require"
)

type externalTypeFixture struct {
	title     string
	schema    string
	expected  *externalTypeDefinition
	knownDefs struct{ tpe, pkg, alias string }
	resolved  resolvedType
}

func makeResolveExternalTypes() []externalTypeFixture {
	return []externalTypeFixture{
		{
			title: "hint as map",
			schema: `{
		"type": "object",
		"x-go-type": {
			"type": "Mytype",
			"import": {
				"package": "github.com/fredbi/mymodels",
				"alias": "external"
			},
			"hints": {
			  "kind": "map"
			},
			"embedded": false
		}
	}`,
			expected: &externalTypeDefinition{
				Type: "Mytype",
				Import: struct {
					Package string
					Alias   string
				}{
					Package: "github.com/fredbi/mymodels",
					Alias:   "external",
				},
				Hints: struct {
					Kind         string
					Nullable     *bool
					NoValidation *bool
				}{
					Kind: "map",
				},
				Embedded: false,
			},
			knownDefs: struct{ tpe, pkg, alias string }{
				tpe:   "external.Mytype",
				pkg:   "github.com/fredbi/mymodels",
				alias: "external",
			},
			resolved: resolvedType{
				GoType:         "external.Mytype",
				IsMap:          true,
				SwaggerType:    "object",
				IsEmptyOmitted: true,
				Pkg:            "github.com/fredbi/mymodels",
				PkgAlias:       "external",
			},
		},
		{
			title: "hint as map, embedded",
			schema: `{
		"type": "object",
		"x-go-type": {
			"type": "Mytype",
			"import": {
				"package": "github.com/fredbi/mymodels",
				"alias": "external"
			},
			"hints": {
			  "kind": "map"
			},
			"embedded": true
		}
	}`,
			expected: &externalTypeDefinition{
				Type: "Mytype",
				Import: struct {
					Package string
					Alias   string
				}{
					Package: "github.com/fredbi/mymodels",
					Alias:   "external",
				},
				Hints: struct {
					Kind         string
					Nullable     *bool
					NoValidation *bool
				}{
					Kind: "map",
				},
				Embedded: true,
			},
			knownDefs: struct{ tpe, pkg, alias string }{
				tpe:   "A",
				pkg:   "",
				alias: "",
			},
			resolved: resolvedType{
				GoType:         "A",
				IsMap:          true,
				SwaggerType:    "object",
				IsEmptyOmitted: true,
			},
		},
		{
			title: "hint as array, nullable",
			schema: `{
		"type": "object",
		"x-go-type": {
			"type": "Mytype",
			"import": {
				"package": "github.com/fredbi/mymodels"
			},
			"hints": {
			  "kind": "array",
				"nullable": true
			}
		}
	}`,
			expected: &externalTypeDefinition{
				Type: "Mytype",
				Import: struct {
					Package string
					Alias   string
				}{
					Package: "github.com/fredbi/mymodels",
					// Alias:   "mymodels",
				},
				Hints: struct {
					Kind         string
					Nullable     *bool
					NoValidation *bool
				}{
					Kind:     "array",
					Nullable: swag.Bool(true),
				},
				Embedded: false,
			},
			knownDefs: struct{ tpe, pkg, alias string }{tpe: "mymodels.Mytype", pkg: "github.com/fredbi/mymodels", alias: "mymodels"},
			resolved: resolvedType{
				GoType:         "mymodels.Mytype",
				IsArray:        true,
				SwaggerType:    "array",
				IsEmptyOmitted: false,
				Pkg:            "github.com/fredbi/mymodels",
				PkgAlias:       "mymodels",
				IsNullable:     true,
			},
		},
		{
			title: "hint as map, unaliased",
			schema: `{
		"type": "object",
		"x-go-type": {
			"type": "Mytype",
			"import": {
				"package": "github.com/fredbi/mymodels"
			},
			"hints": {
			  "kind": "map"
			}
		}
	}`,
			expected: &externalTypeDefinition{
				Type: "Mytype",
				Import: struct {
					Package string
					Alias   string
				}{
					Package: "github.com/fredbi/mymodels",
					// Alias:   "mymodels",
				},
				Hints: struct {
					Kind         string
					Nullable     *bool
					NoValidation *bool
				}{
					Kind: "map",
				},
			},
			knownDefs: struct{ tpe, pkg, alias string }{tpe: "mymodels.Mytype", pkg: "github.com/fredbi/mymodels", alias: "mymodels"},
			resolved: resolvedType{
				GoType:         "mymodels.Mytype",
				IsMap:          true,
				SwaggerType:    "object",
				IsEmptyOmitted: true,
				Pkg:            "github.com/fredbi/mymodels",
				PkgAlias:       "mymodels",
			},
		},
		{
			title: "hint as tuple, unaliased",
			schema: `{
		"type": "object",
		"x-go-type": {
			"type": "Mytype",
			"import": {
				"package": "github.com/fredbi/mymodels"
			},
			"hints": {
			  "kind": "tuple"
			}
		}
	}`,
			expected: &externalTypeDefinition{
				Type: "Mytype",
				Import: struct {
					Package string
					Alias   string
				}{
					Package: "github.com/fredbi/mymodels",
					// Alias:   "mymodels",
				},
				Hints: struct {
					Kind         string
					Nullable     *bool
					NoValidation *bool
				}{
					Kind: "tuple",
				},
			},
			knownDefs: struct{ tpe, pkg, alias string }{tpe: "mymodels.Mytype", pkg: "github.com/fredbi/mymodels", alias: "mymodels"},
			resolved: resolvedType{
				GoType:         "mymodels.Mytype",
				IsTuple:        true,
				SwaggerType:    "array",
				IsEmptyOmitted: true,
				Pkg:            "github.com/fredbi/mymodels",
				PkgAlias:       "mymodels",
			},
		},
		{
			title: "hint as primitive, unaliased",
			schema: `{
		"type": "number",
		"x-go-type": {
			"type": "Mytype",
			"import": {
				"package": "github.com/fredbi/mymodels"
			},
			"hints": {
			  "kind": "primitive"
			}
		}
	}`,
			expected: &externalTypeDefinition{
				Type: "Mytype",
				Import: struct {
					Package string
					Alias   string
				}{
					Package: "github.com/fredbi/mymodels",
					// Alias:   "mymodels",
				},
				Hints: struct {
					Kind         string
					Nullable     *bool
					NoValidation *bool
				}{
					Kind: "primitive",
				},
			},
			knownDefs: struct{ tpe, pkg, alias string }{tpe: "mymodels.Mytype", pkg: "github.com/fredbi/mymodels", alias: "mymodels"},
			resolved: resolvedType{
				GoType:         "mymodels.Mytype",
				IsPrimitive:    true,
				SwaggerType:    "",
				IsEmptyOmitted: true,
				Pkg:            "github.com/fredbi/mymodels",
				PkgAlias:       "mymodels",
			},
		},
		{
			title: "default model package",
			schema: `{
		"type": "number",
		"x-go-type": {
			"type": "Mytype",
			"hints": {
			  "kind": "primitive"
			}
		}
	}`,
			expected: &externalTypeDefinition{
				Type: "Mytype",
				Import: struct {
					Package string
					Alias   string
				}{
					// Package: "github.com/example/custom",
					// Alias:   "custom",
				},
				Hints: struct {
					Kind         string
					Nullable     *bool
					NoValidation *bool
				}{
					Kind: "primitive",
				},
			},
			knownDefs: struct{ tpe, pkg, alias string }{tpe: "Mytype", pkg: "", alias: ""},
			resolved: resolvedType{
				GoType:         "Mytype",
				IsPrimitive:    true,
				SwaggerType:    "",
				IsEmptyOmitted: true,
				Pkg:            "",
				PkgAlias:       "",
			},
		},
	}
}

func TestShortCircuitResolveExternal(t *testing.T) {
	defer discardOutput()()

	for i, toPin := range makeResolveExternalTypes() {
		fixture := toPin
		var title string
		if fixture.title == "" {
			title = strconv.Itoa(i)
		} else {
			title = fixture.title
		}
		t.Run(title, func(t *testing.T) {
			jazonDoc := fixture.schema
			doc, err := loads.Embedded([]byte(jazonDoc), []byte(jazonDoc))
			require.NoErrorf(t, err, "fixture %d", i)

			r := newTypeResolver("models", "github.com/example/custom", doc)
			var schema spec.Schema
			err = json.Unmarshal([]byte(jazonDoc), &schema)
			require.NoErrorf(t, err, "fixture %d", i)

			extType, ok := hasExternalType(schema.Extensions)
			require.Truef(t, ok, "fixture %d", i)
			require.NotNil(t, extType)

			tpe, pkg, alias := r.knownDefGoType("A", schema, r.goTypeName)
			require.EqualValuesf(t, fixture.knownDefs, struct{ tpe, pkg, alias string }{tpe, pkg, alias}, "fixture %d", i)

			resolved := r.shortCircuitResolveExternal(tpe, pkg, alias, extType, &schema, false)

			require.EqualValues(t, fixture.expected, extType)

			resolved.Extensions = nil // don't assert this
			require.EqualValuesf(t, fixture.resolved, resolved, "fixture %d", i)
		})
	}
}

type guardValidationsFixture struct {
	Title        string
	ResolvedType string
	Type         interface {
		Validations() spec.SchemaValidations
		SetValidations(spec.SchemaValidations)
	}
	Asserter func(testing.TB, spec.SchemaValidations)
}

func makeGuardValidationFixtures() []guardValidationsFixture {
	return []guardValidationsFixture{
		{
			Title:        "simple schema: guard array",
			ResolvedType: "array",
			Type: spec.NewItems().
				Typed("number", "int64").
				WithValidations(spec.CommonValidations{MinLength: swag.Int64(15), Maximum: swag.Float64(12.00)}).
				UniqueValues(),
			Asserter: func(t testing.TB, val spec.SchemaValidations) {
				require.False(t, val.HasNumberValidations(), "expected no number validations, got: %#v", val)
				require.False(t, val.HasStringValidations(), "expected no string validations, got: %#v", val)
				require.True(t, val.HasArrayValidations(), "expected array validations, got: %#v", val)
			},
		},
		{
			Title:        "simple schema: guard string",
			ResolvedType: "string",
			Type: spec.QueryParam("p1").
				Typed("string", "uuid").
				WithValidations(spec.CommonValidations{MinItems: swag.Int64(15), Maximum: swag.Float64(12.00)}).
				WithMinLength(12),
			Asserter: func(t testing.TB, val spec.SchemaValidations) {
				require.False(t, val.HasNumberValidations(), "expected no number validations, got: %#v", val)
				require.False(t, val.HasArrayValidations(), "expected no array validations, got: %#v", val)
				require.True(t, val.HasStringValidations(), "expected string validations, got: %#v", val)
			},
		},
		{
			Title:        "simple schema: guard file (1/3)",
			ResolvedType: "file",
			Type: spec.FileParam("p1").
				WithValidations(spec.CommonValidations{MinItems: swag.Int64(15), Maximum: swag.Float64(12.00)}).
				WithMinLength(12),
			Asserter: func(t testing.TB, val spec.SchemaValidations) {
				require.False(t, val.HasNumberValidations(), "expected no number validations, got: %#v", val)
				require.False(t, val.HasArrayValidations(), "expected no array validations, got: %#v", val)
				require.True(t, val.HasStringValidations(), "expected string validations, got: %#v", val)
			},
		},
		{
			Title:        "simple schema: guard file (2/3)",
			ResolvedType: "file",
			Type: spec.FileParam("p1").
				WithValidations(spec.CommonValidations{
					MinItems: swag.Int64(15),
					Maximum:  swag.Float64(12.00),
					Pattern:  "xyz",
					Enum:     []interface{}{"x", 34},
				}),
			Asserter: func(t testing.TB, val spec.SchemaValidations) {
				require.False(t, val.HasNumberValidations(), "expected no number validations, got: %#v", val)
				require.False(t, val.HasArrayValidations(), "expected no array validations, got: %#v", val)
				require.False(t, val.HasStringValidations(), "expected no string validations, got: %#v", val)
				require.False(t, val.HasEnum(), "expected no enum validations, got: %#v", val)
			},
		},
		{
			Title:        "schema: guard object",
			ResolvedType: "object",
			Type: spec.RefSchema("#/definitions/nowhere").
				WithValidations(spec.SchemaValidations{
					CommonValidations: spec.CommonValidations{
						MinItems: swag.Int64(15),
						Maximum:  swag.Float64(12.00),
					},
					MinProperties: swag.Int64(10),
				}).
				WithMinLength(12),
			Asserter: func(t testing.TB, val spec.SchemaValidations) {
				require.False(t, val.HasNumberValidations(), "expected no number validations, got: %#v", val)
				require.False(t, val.HasArrayValidations(), "expected no array validations, got: %#v", val)
				require.False(t, val.HasStringValidations(), "expected no string validations, got: %#v", val)
				require.True(t, val.HasObjectValidations(), "expected object validations, got: %#v", val)
			},
		},
		{
			Title:        "simple schema: guard number",
			ResolvedType: "number",
			Type: spec.QueryParam("p1").
				Typed("number", "double").
				WithValidations(spec.CommonValidations{MinItems: swag.Int64(15), MultipleOf: swag.Float64(12.00), Pattern: "xyz"}).
				WithMinLength(12),
			Asserter: func(t testing.TB, val spec.SchemaValidations) {
				require.False(t, val.HasArrayValidations(), "expected no array validations, got: %#v", val)
				require.False(t, val.HasStringValidations(), "expected no string validations, got: %#v", val)
				require.True(t, val.HasNumberValidations(), "expected number validations, got: %#v", val)
			},
		},
	}
}

func TestGuardValidations(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	defer func() {
		log.SetOutput(os.Stdout)
	}()

	for _, toPin := range makeGuardValidationFixtures() {
		testCase := toPin
		t.Run(testCase.Title, func(t *testing.T) {
			t.Parallel()
			input := testCase.Type
			guardValidations(testCase.ResolvedType, input)
			if testCase.Asserter != nil {
				testCase.Asserter(t, input.Validations())
			}
		})
	}
}

func makeGuardFormatFixtures() []guardValidationsFixture {
	return []guardValidationsFixture{
		{
			Title:        "schema: guard date format",
			ResolvedType: "date",
			Type: spec.StringProperty().
				WithValidations(spec.SchemaValidations{
					CommonValidations: spec.CommonValidations{
						MinLength: swag.Int64(15),
						Pattern:   "xyz",
						Enum:      []interface{}{"x", 34},
					}}),
			Asserter: func(t testing.TB, val spec.SchemaValidations) {
				require.True(t, val.HasStringValidations(), "expected string validations, got: %#v", val)
				require.True(t, val.HasEnum())
			},
		},
		{
			Title:        "simple schema: guard binary format",
			ResolvedType: "binary",
			Type: spec.StringProperty().
				WithValidations(spec.SchemaValidations{
					CommonValidations: spec.CommonValidations{
						MinLength: swag.Int64(15),
						Pattern:   "xyz",
						Enum:      []interface{}{"x", 34},
					}}),
			Asserter: func(t testing.TB, val spec.SchemaValidations) {
				require.False(t, val.HasStringValidations(), "expected no string validations, got: %#v", val)
				require.False(t, val.HasEnum())
			},
		},
	}
}

func TestGuardFormatConflicts(t *testing.T) {
	defer discardOutput()()

	for _, toPin := range makeGuardFormatFixtures() {
		testCase := toPin
		t.Run(testCase.Title, func(t *testing.T) {
			t.Parallel()
			input := testCase.Type
			guardFormatConflicts(testCase.ResolvedType, input)
			if testCase.Asserter != nil {
				testCase.Asserter(t, input.Validations())
			}
		})
	}
}
