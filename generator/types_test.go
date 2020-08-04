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
	"github.com/stretchr/testify/require"
)

type externalTypeFixture struct {
	schema    string
	expected  *externalTypeDefinition
	knownDefs struct{ tpe, pkg, alias string }
	resolved  resolvedType
}

func makeResolveExternalTypes() []externalTypeFixture {
	return []externalTypeFixture{
		{
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
					Kind     string
					Nullable bool
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
					Kind     string
					Nullable bool
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
					Kind     string
					Nullable bool
				}{
					Kind:     "array",
					Nullable: true,
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
			},
		},
		{
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
					Kind     string
					Nullable bool
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
					Kind     string
					Nullable bool
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
					Kind     string
					Nullable bool
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
	}
}

func TestShortCircuitResolveExternal(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	defer func() {
		log.SetOutput(os.Stdout)
	}()

	for i, toPin := range makeResolveExternalTypes() {
		fixture := toPin
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			jazonDoc := fixture.schema
			doc, err := loads.Embedded([]byte(jazonDoc), []byte(jazonDoc))
			require.NoErrorf(t, err, "fixture %d", i)

			r := newTypeResolver("models", doc)
			var schema spec.Schema
			err = json.Unmarshal([]byte(jazonDoc), &schema)
			require.NoErrorf(t, err, "fixture %d", i)

			extType, ok := hasExternalType(schema.Extensions)
			require.Truef(t, ok, "fixture %d", i)

			require.EqualValuesf(t, fixture.expected, extType, "fixture %d", i)

			tpe, pkg, alias := knownDefGoType("A", schema, r.goTypeName)
			require.EqualValuesf(t, fixture.knownDefs, struct{ tpe, pkg, alias string }{tpe, pkg, alias}, "fixture %d", i)

			resolved := r.shortCircuitResolveExternal(tpe, pkg, alias, extType, &schema)

			resolved.Extensions = nil // don't assert this
			require.EqualValuesf(t, fixture.resolved, resolved, "fixture %d", i)
		})
	}
}
