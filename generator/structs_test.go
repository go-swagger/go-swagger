package generator

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPrintTags(t *testing.T) {
	type tagFixture struct {
		Title        string
		Schema       GenSchema
		ExpectedTags string
	}

	mustJSON := func(in interface{}) string {
		b, _ := asJSON(in)
		return b
	}

	fixtures := []tagFixture{
		{
			Title: "no extra: default json",
			Schema: GenSchema{
				OriginalName: "field",
			},
			ExpectedTags: "`json:\"field\"`",
		},
		{
			Title: "no extra: default json, omitempty",
			Schema: GenSchema{
				OriginalName: "field",
				resolvedType: resolvedType{
					IsEmptyOmitted: true,
				},
			},
			ExpectedTags: "`json:\"field,omitempty\"`",
		},
		{
			Title: "no extra: default json, required, omitempty",
			Schema: GenSchema{
				OriginalName: "field",
				resolvedType: resolvedType{
					IsEmptyOmitted: true,
				},
				sharedValidations: sharedValidations{
					Required: true,
				},
			},
			ExpectedTags: "`json:\"field\"`",
		},
		{
			Title: "no extra: default json, omitempty, as JSON string",
			Schema: GenSchema{
				OriginalName: "field",
				resolvedType: resolvedType{
					IsEmptyOmitted: true,
					IsJSONString:   true,
				},
			},
			ExpectedTags: "`json:\"field,omitempty,string\"`",
		},
		{
			Title: "with xml name",
			Schema: GenSchema{
				OriginalName: "field",
				resolvedType: resolvedType{
					IsEmptyOmitted: true,
				},
				sharedValidations: sharedValidations{
					Required: true,
				},
				XMLName: "xmlfield",
			},
			ExpectedTags: "`json:\"field\" xml:\"xmlfield\"`",
		},
		{
			Title: "with example (1/3)",
			Schema: GenSchema{
				OriginalName: "field",
				resolvedType: resolvedType{
					IsEmptyOmitted: true,
				},
				sharedValidations: sharedValidations{
					Required: true,
				},
				Example:    mustJSON("xyz"),
				StructTags: []string{"example"},
			},
			ExpectedTags: "`json:\"field\" example:\"\\\"xyz\\\"\"`",
		},
		{
			Title: "with example (2/3)",
			Schema: GenSchema{
				OriginalName: "field",
				resolvedType: resolvedType{
					IsEmptyOmitted: true,
				},
				sharedValidations: sharedValidations{
					Required: true,
				},
				Example:    mustJSON(15),
				StructTags: []string{"example"},
			},
			ExpectedTags: "`json:\"field\" example:\"15\"`",
		},
		{
			Title: "with example (3/3)",
			Schema: GenSchema{
				OriginalName: "field",
				resolvedType: resolvedType{
					IsEmptyOmitted: true,
				},
				sharedValidations: sharedValidations{
					Required: true,
				},
				Example: mustJSON(struct {
					A string `json:"a"`
					B int64  `json:"b"`
				}{A: "xyz", B: 12}),
				StructTags: []string{"example"},
			},
			ExpectedTags: "`json:\"field\" example:\"{\\\"a\\\":\\\"xyz\\\",\\\"b\\\":12}\"`",
		},
		{
			Title: "with example, xml, omitempty, custom tag",
			Schema: GenSchema{
				OriginalName: "field",
				resolvedType: resolvedType{
					IsEmptyOmitted: true,
				},
				Example:    mustJSON(15),
				StructTags: []string{"example"},
				XMLName:    "xmlfield,attr",
				CustomTag:  `metric:"on"`,
			},
			ExpectedTags: "`json:\"field,omitempty\" xml:\"xmlfield,attr,omitempty\" example:\"15\" metric:\"on\"`",
		},
		{
			Title: "with backticks",
			Schema: GenSchema{
				OriginalName: "field",
				Example: mustJSON(struct {
					A string `json:"a"`
					B int64  `json:"b"`
				}{A: "`xyz`", B: 12}),
				StructTags: []string{"example"},
				CustomTag:  "metric:\"`on`\"",
			},
			ExpectedTags: "\"json:\\\"field\\\" example:\\\"{\\\\\\\"a\\\\\\\":\\\\\\\"`xyz`\\\\\\\",\\\\\\\"b\\\\\\\":12}\\\" metric:\\\"`on`\\\"\"",
		},
		{
			Title: "with description",
			Schema: GenSchema{
				OriginalName: "field",
				Description: "some description",
				StructTags: []string{"description"},
			},
			ExpectedTags: "`json:\"field\" description:\"some description\"`",
		},
		{
			Title: "with multiline description",
			Schema: GenSchema{
				OriginalName: "field",
				Description: "a\ndescription\nspanning\nmultiple\nlines",
				StructTags: []string{"description"},
			},
			ExpectedTags: "\"json:\\\"field\\\" description:\\\"a\\\\ndescription\\\\nspanning\\\\nmultiple\\\\nlines\\\"\"",
		},
	}

	for _, toPin := range fixtures {
		fixture := toPin
		t.Run(fixture.Title, func(t *testing.T) {
			require.Equal(t, fixture.ExpectedTags, fixture.Schema.PrintTags())
		})
	}
}
