// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package templates

import (
	"bytes"
	"testing"
	"text/template"

	"github.com/go-openapi/swag/mangling"
	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"

	golangfuncs "github.com/go-swagger/go-swagger/generator/internal/funcmaps/golang"
	templatesrepo "github.com/go-swagger/go-swagger/generator/internal/templates-repo"
)

func TestIntegrationTemplates(t *testing.T) {
	t.Parallel()

	m := mangling.NewNameMangler(mangling.WithGoNamePrefixFunc(golangfuncs.PrefixForName))
	fm := golangfuncs.FuncMap(m)

	t.Run("pascalize and camelize should handle $ in names", func(t *testing.T) {
		t.Parallel()

		// issue #2821
		const tpl = `
Pascalize={{ pascalize . }}
Camelize={{ camelize . }}
`
		rendered := renderTemplateInRepo(fm, tpl, "get$ref")(t)
		assert.StringContainsT(t, rendered, "Pascalize=GetDollarRef\n")
		assert.StringContainsT(t, rendered, "Camelize=getDollarRef\n")
	})

	t.Run("should execute templates with funcmap-provided functions", func(t *testing.T) {
		t.Parallel()

		// Exercises the funcmap functions that use only
		// literal arguments — no generator types or LanguageOpts required.
		const tpl = `
Pascalize={{ pascalize "WeArePonies_Of_the_round table" }}
Humanize={{ humanize "WeArePonies_Of_the_round table" }}
PluralizeFirstWord={{ pluralizeFirstWord "pony of the round table" }}
PluralizeFirstOfOneWord={{ pluralizeFirstWord "dwarf" }}
PluralizeFirstOfNoWord={{ pluralizeFirstWord "" }}
DropPackage={{ dropPackage "prefix.suffix" }}
DropNoPackage={{ dropPackage "suffix" }}
DropEmptyPackage={{ dropPackage "" }}
PadSurround1={{ padSurround "padme" "-" 3 12}}
PadSurround2={{ padSurround "padme" "-" 0 12}}
PascalizeSpecialChar1={{ pascalize "+1" }}
PascalizeSpecialChar2={{ pascalize "-1" }}
PascalizeSpecialChar3={{ pascalize "1" }}
PascalizeSpecialChar4={{ pascalize "-" }}
PascalizeSpecialChar5={{ pascalize "+" }}
PascalizeCleanupEnumVariant1={{ pascalize (cleanupEnumVariant "2.4Ghz") }}
Dict={{ template "dictTemplate" dict "Animal" "Pony" "Shape" "round" "Furniture" "table" }}
{{ define "dictTemplate" }}{{ .Animal }} of the {{ .Shape }} {{ .Furniture }}{{ end }}
`

		rendered := renderTemplateInRepo(fm, tpl, nil)(t)

		assert.StringContainsT(t, rendered, "Pascalize=WeArePoniesOfTheRoundTable\n")
		assert.StringContainsT(t, rendered, "Humanize=we are ponies of the round table\n")
		assert.StringContainsT(t, rendered, "PluralizeFirstWord=ponies of the round table\n")
		assert.StringContainsT(t, rendered, "PluralizeFirstOfOneWord=dwarves\n")
		assert.StringContainsT(t, rendered, "PluralizeFirstOfNoWord=\n")
		assert.StringContainsT(t, rendered, "DropPackage=suffix\n")
		assert.StringContainsT(t, rendered, "DropNoPackage=suffix\n")
		assert.StringContainsT(t, rendered, "DropEmptyPackage=\n")
		assert.StringContainsT(t, rendered, "PadSurround1=-,-,-,padme,-,-,-,-,-,-,-,-\n")
		assert.StringContainsT(t, rendered, "PadSurround2=padme,-,-,-,-,-,-,-,-,-,-,-\n")
		assert.StringContainsT(t, rendered, "PascalizeSpecialChar1=Plus1\n")
		assert.StringContainsT(t, rendered, "PascalizeSpecialChar2=Minus1\n")
		assert.StringContainsT(t, rendered, "PascalizeSpecialChar3=Nr1\n")
		assert.StringContainsT(t, rendered, "PascalizeSpecialChar4=Minus\n")
		assert.StringContainsT(t, rendered, "PascalizeSpecialChar5=Plus\n")
		assert.StringContainsT(t, rendered, "PascalizeCleanupEnumVariant1=Nr2Dot4Ghz")
		assert.StringContainsT(t, rendered, "Dict=Pony of the round table\n")
	})

	t.Run("should execute with inline closure functions from funcmaps", func(t *testing.T) {
		t.Parallel()

		// Exercises the inline closure functions defined
		// in [golangfuncs.FuncMap] via template execution.
		const tpl = `
HasInsecureHTTP={{ hasInsecure .Schemes }}
HasInsecureWS={{ hasInsecure .WSSchemes }}
HasInsecureHTTPS={{ hasInsecure .SecureSchemes }}
HasSecureHTTPS={{ hasSecure .SecureSchemes }}
HasSecureWSS={{ hasSecure .WSSSchemes }}
HasSecureHTTP={{ hasSecure .Schemes }}
EscapeNone={{ escapeBackticks "no ticks" }}
FlagName={{ flagNameVar "myField" }}
FlagValue={{ flagValueVar "myField" }}
FlagDefault={{ flagDefaultVar "myField" }}
FlagModel={{ flagModelVar "myField" }}
FlagDescription={{ flagDescriptionVar "myField" }}
PrintGoLiteral={{ printGoLiteral "hello" }}
`
		data := map[string]any{
			"Schemes":       []string{"http"},
			"WSSchemes":     []string{"ws"},
			"SecureSchemes": []string{"https"},
			"WSSSchemes":    []string{"wss"},
		}

		rendered := renderTemplateInRepo(fm, tpl, data)(t)

		assert.StringContainsT(t, rendered, "HasInsecureHTTP=true\n")
		assert.StringContainsT(t, rendered, "HasInsecureWS=true\n")
		assert.StringContainsT(t, rendered, "HasInsecureHTTPS=false\n")
		assert.StringContainsT(t, rendered, "HasSecureHTTPS=true\n")
		assert.StringContainsT(t, rendered, "HasSecureWSS=true\n")
		assert.StringContainsT(t, rendered, "HasSecureHTTP=false\n")
		assert.StringContainsT(t, rendered, "EscapeNone=no ticks\n")
		assert.StringContainsT(t, rendered, "FlagName=flagMyFieldName\n")
		assert.StringContainsT(t, rendered, "FlagValue=flagMyFieldValue\n")
		assert.StringContainsT(t, rendered, "FlagDefault=flagMyFieldDefault\n")
		assert.StringContainsT(t, rendered, "FlagModel=flagMyFieldModel\n")
		assert.StringContainsT(t, rendered, "FlagDescription=flagMyFieldDescription\n")
		assert.StringContainsT(t, rendered, `PrintGoLiteral="hello"`)
	})
}

func renderTemplateInRepo(fm template.FuncMap, tpl string, data any) func(*testing.T) string {
	return func(t *testing.T) (rendered string) {
		t.Helper()

		repo := templatesrepo.NewRepository(fm)
		require.NoError(t, repo.AddFile("test", tpl))

		templ, err := repo.Get("test")
		require.NoError(t, err)

		var buf bytes.Buffer
		require.NoError(t, templ.Execute(&buf, data))

		return buf.String()
	}
}
