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
					Alias:   "mymodels",
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
					Alias:   "mymodels",
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
					Alias:   "mymodels",
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
					Alias:   "mymodels",
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
					Package: "github.com/example/custom",
					Alias:   "custom",
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
	log.SetOutput(ioutil.Discard)
	defer func() {
		log.SetOutput(os.Stdout)
	}()

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

			resolved.Extensions = nil // don't assert this
			require.EqualValuesf(t, fixture.resolved, resolved, "fixture %d", i)
		})
	}
}
