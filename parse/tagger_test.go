package parse

import (
	"fmt"
	"go/ast"
	"regexp"
	"strings"
	"testing"

	"github.com/casualjim/go-swagger/spec"
	"github.com/stretchr/testify/assert"
)

func TestSectionedParser_TitleDescription(t *testing.T) {
	text := `This has a title, separated by a whitespace line

In this example the punctuation for the title should not matter for swagger.
For go it will still make a difference though.
`
	text2 := `This has a title without whitespace.
The punctuation here does indeed matter. But it won't for go.
`

	st := &sectionedParser{}
	st.setTitle = func(lines []string) {}
	st.Parse(ascg(text))

	assert.EqualValues(t, []string{"This has a title, separated by a whitespace line"}, st.Title())
	assert.EqualValues(t, []string{"In this example the punctuation for the title should not matter for swagger.", "For go it will still make a difference though."}, st.Description())

	st = &sectionedParser{}
	st.setTitle = func(lines []string) {}
	st.Parse(ascg(text2))

	assert.EqualValues(t, []string{"This has a title without whitespace."}, st.Title())
	assert.EqualValues(t, []string{"The punctuation here does indeed matter. But it won't for go."}, st.Description())
}

func dummyBuilder() schemaValidations {
	return schemaValidations{new(spec.Schema)}
}

func TestSectionedParser_TagsDescription(t *testing.T) {
	block := `This has a title without whitespace.
The punctuation here does indeed matter. But it won't for go.
minimum: 10
maximum: 20
`
	block2 := `This has a title without whitespace.
The punctuation here does indeed matter. But it won't for go.

minimum: 10
maximum: 20
`

	st := &sectionedParser{}
	st.setTitle = func(lines []string) {}
	st.taggers = []tagParser{
		{"Maximum", false, nil, &setMaximum{dummyBuilder(), regexp.MustCompile(fmt.Sprintf(rxMaximumFmt, ""))}},
		{"Minimum", false, nil, &setMinimum{dummyBuilder(), regexp.MustCompile(fmt.Sprintf(rxMinimumFmt, ""))}},
		{"MultipleOf", false, nil, &setMultipleOf{dummyBuilder(), regexp.MustCompile(fmt.Sprintf(rxMultipleOfFmt, ""))}},
	}

	st.Parse(ascg(block))
	assert.EqualValues(t, []string{"This has a title without whitespace."}, st.Title())
	assert.EqualValues(t, []string{"The punctuation here does indeed matter. But it won't for go."}, st.Description())
	assert.Len(t, st.matched, 2)
	_, ok := st.matched["Maximum"]
	assert.True(t, ok)
	_, ok = st.matched["Minimum"]
	assert.True(t, ok)

	st = &sectionedParser{}
	st.setTitle = func(lines []string) {}
	st.taggers = []tagParser{
		{"Maximum", false, nil, &setMaximum{dummyBuilder(), regexp.MustCompile(fmt.Sprintf(rxMaximumFmt, ""))}},
		{"Minimum", false, nil, &setMinimum{dummyBuilder(), regexp.MustCompile(fmt.Sprintf(rxMinimumFmt, ""))}},
		{"MultipleOf", false, nil, &setMultipleOf{dummyBuilder(), regexp.MustCompile(fmt.Sprintf(rxMultipleOfFmt, ""))}},
	}

	st.Parse(ascg(block2))
	assert.EqualValues(t, []string{"This has a title without whitespace."}, st.Title())
	assert.EqualValues(t, []string{"The punctuation here does indeed matter. But it won't for go."}, st.Description())
	assert.Len(t, st.matched, 2)
	_, ok = st.matched["Maximum"]
	assert.True(t, ok)
	_, ok = st.matched["Minimum"]
	assert.True(t, ok)
}

func TestSectionedParser_SkipSectionAnnotation(t *testing.T) {
	block := `+swagger:model someModel

This has a title without whitespace.
The punctuation here does indeed matter. But it won't for go.

minimum: 10
maximum: 20
`
	st := &sectionedParser{}
	st.setTitle = func(lines []string) {}
	ap := newSchemaAnnotationParser("SomeModel")
	st.annotation = ap
	st.taggers = []tagParser{
		{"Maximum", false, nil, &setMaximum{dummyBuilder(), regexp.MustCompile(fmt.Sprintf(rxMaximumFmt, ""))}},
		{"Minimum", false, nil, &setMinimum{dummyBuilder(), regexp.MustCompile(fmt.Sprintf(rxMinimumFmt, ""))}},
		{"MultipleOf", false, nil, &setMultipleOf{dummyBuilder(), regexp.MustCompile(fmt.Sprintf(rxMultipleOfFmt, ""))}},
	}

	st.Parse(ascg(block))
	assert.EqualValues(t, []string{"This has a title without whitespace."}, st.Title())
	assert.EqualValues(t, []string{"The punctuation here does indeed matter. But it won't for go."}, st.Description())
	assert.Len(t, st.matched, 2)
	_, ok := st.matched["Maximum"]
	assert.True(t, ok)
	_, ok = st.matched["Minimum"]
	assert.True(t, ok)
	assert.Equal(t, "SomeModel", ap.GoName)
	assert.Equal(t, "someModel", ap.Name)
}

func TestSectionedParser_TerminateOnNewAnnotation(t *testing.T) {
	block := `+swagger:model someModel

This has a title without whitespace.
The punctuation here does indeed matter. But it won't for go.

minimum: 10
+swagger:meta
maximum: 20
`
	st := &sectionedParser{}
	st.setTitle = func(lines []string) {}
	ap := newSchemaAnnotationParser("SomeModel")
	st.annotation = ap
	st.taggers = []tagParser{
		{"Maximum", false, nil, &setMaximum{dummyBuilder(), regexp.MustCompile(fmt.Sprintf(rxMaximumFmt, ""))}},
		{"Minimum", false, nil, &setMinimum{dummyBuilder(), regexp.MustCompile(fmt.Sprintf(rxMinimumFmt, ""))}},
		{"MultipleOf", false, nil, &setMultipleOf{dummyBuilder(), regexp.MustCompile(fmt.Sprintf(rxMultipleOfFmt, ""))}},
	}

	st.Parse(ascg(block))
	assert.EqualValues(t, []string{"This has a title without whitespace."}, st.Title())
	assert.EqualValues(t, []string{"The punctuation here does indeed matter. But it won't for go."}, st.Description())
	assert.Len(t, st.matched, 1)
	_, ok := st.matched["Maximum"]
	assert.False(t, ok)
	_, ok = st.matched["Minimum"]
	assert.True(t, ok)
	assert.Equal(t, "SomeModel", ap.GoName)
	assert.Equal(t, "someModel", ap.Name)
}

func ascg(txt string) *ast.CommentGroup {
	var cg ast.CommentGroup
	for _, line := range strings.Split(txt, "\n") {
		var cmt ast.Comment
		cmt.Text = "// " + line
		cg.List = append(cg.List, &cmt)
	}
	return &cg
}
