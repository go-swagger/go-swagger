package codescan

import (
	"fmt"
	"go/ast"
	"regexp"
	"strings"
	"testing"

	"github.com/go-openapi/spec"
	"github.com/stretchr/testify/assert"
)

// only used within this group of tests but never used within actual code base.
func newSchemaAnnotationParser(goName string) *schemaAnnotationParser {
	return &schemaAnnotationParser{GoName: goName, rx: rxModelOverride}
}

type schemaAnnotationParser struct {
	GoName string
	Name   string
	rx     *regexp.Regexp
}

func (sap *schemaAnnotationParser) Matches(line string) bool {
	return sap.rx.MatchString(line)
}

func (sap *schemaAnnotationParser) Parse(lines []string) error {
	if sap.Name != "" {
		return nil
	}

	if len(lines) > 0 {
		for _, line := range lines {
			matches := sap.rx.FindStringSubmatch(line)
			if len(matches) > 1 && len(matches[1]) > 0 {
				sap.Name = matches[1]
				return nil
			}
		}
	}
	return nil
}

func TestSectionedParser_TitleDescription(t *testing.T) {
	text := `This has a title, separated by a whitespace line

In this example the punctuation for the title should not matter for swagger.
For go it will still make a difference though.
`
	text2 := `This has a title without whitespace.
The punctuation here does indeed matter. But it won't for go.
`

	text3 := `This has a title, and markdown in the description

See how markdown works now, we can have lists:

+ first item
+ second item
+ third item

[Links works too](http://localhost)
`

	text4 := `This has whitespace sensitive markdown in the description

|+ first item
|    + nested item
|    + also nested item

Sample code block:

|    fmt.Println("Hello World!")

`

	var err error

	st := &sectionedParser{}
	st.setTitle = func(lines []string) {}
	err = st.Parse(ascg(text))
	assert.NoError(t, err)

	assert.EqualValues(t, []string{"This has a title, separated by a whitespace line"}, st.Title())
	assert.EqualValues(t, []string{"In this example the punctuation for the title should not matter for swagger.", "For go it will still make a difference though."}, st.Description())

	st = &sectionedParser{}
	st.setTitle = func(lines []string) {}
	err = st.Parse(ascg(text2))
	assert.NoError(t, err)

	assert.EqualValues(t, []string{"This has a title without whitespace."}, st.Title())
	assert.EqualValues(t, []string{"The punctuation here does indeed matter. But it won't for go."}, st.Description())

	st = &sectionedParser{}
	st.setTitle = func(lines []string) {}
	err = st.Parse(ascg(text3))
	assert.NoError(t, err)

	assert.EqualValues(t, []string{"This has a title, and markdown in the description"}, st.Title())
	assert.EqualValues(t, []string{"See how markdown works now, we can have lists:", "", "+ first item", "+ second item", "+ third item", "", "[Links works too](http://localhost)"}, st.Description())

	st = &sectionedParser{}
	st.setTitle = func(lines []string) {}
	err = st.Parse(ascg(text4))
	assert.NoError(t, err)

	assert.EqualValues(t, []string{"This has whitespace sensitive markdown in the description"}, st.Title())
	assert.EqualValues(t, []string{"+ first item", "    + nested item", "    + also nested item", "", "Sample code block:", "", "    fmt.Println(\"Hello World!\")"}, st.Description())
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

	var err error

	st := &sectionedParser{}
	st.setTitle = func(lines []string) {}
	st.taggers = []tagParser{
		{"Maximum", false, false, nil, &setMaximum{dummyBuilder(), regexp.MustCompile(fmt.Sprintf(rxMaximumFmt, ""))}},
		{"Minimum", false, false, nil, &setMinimum{dummyBuilder(), regexp.MustCompile(fmt.Sprintf(rxMinimumFmt, ""))}},
		{"MultipleOf", false, false, nil, &setMultipleOf{dummyBuilder(), regexp.MustCompile(fmt.Sprintf(rxMultipleOfFmt, ""))}},
	}

	err = st.Parse(ascg(block))
	assert.NoError(t, err)
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
		{"Maximum", false, false, nil, &setMaximum{dummyBuilder(), regexp.MustCompile(fmt.Sprintf(rxMaximumFmt, ""))}},
		{"Minimum", false, false, nil, &setMinimum{dummyBuilder(), regexp.MustCompile(fmt.Sprintf(rxMinimumFmt, ""))}},
		{"MultipleOf", false, false, nil, &setMultipleOf{dummyBuilder(), regexp.MustCompile(fmt.Sprintf(rxMultipleOfFmt, ""))}},
	}

	err = st.Parse(ascg(block2))
	assert.NoError(t, err)
	assert.EqualValues(t, []string{"This has a title without whitespace."}, st.Title())
	assert.EqualValues(t, []string{"The punctuation here does indeed matter. But it won't for go."}, st.Description())
	assert.Len(t, st.matched, 2)
	_, ok = st.matched["Maximum"]
	assert.True(t, ok)
	_, ok = st.matched["Minimum"]
	assert.True(t, ok)
}

func TestSectionedParser_Empty(t *testing.T) {
	block := `swagger:response someResponse`

	var err error

	st := &sectionedParser{}
	st.setTitle = func(lines []string) {}
	ap := newSchemaAnnotationParser("SomeResponse")
	ap.rx = rxResponseOverride
	st.annotation = ap

	err = st.Parse(ascg(block))
	assert.NoError(t, err)
	assert.Empty(t, st.Title())
	assert.Empty(t, st.Description())
	assert.Empty(t, st.taggers)
	assert.Equal(t, "SomeResponse", ap.GoName)
	assert.Equal(t, "someResponse", ap.Name)
}

func TestSectionedParser_SkipSectionAnnotation(t *testing.T) {
	block := `swagger:model someModel

This has a title without whitespace.
The punctuation here does indeed matter. But it won't for go.

minimum: 10
maximum: 20
`
	var err error

	st := &sectionedParser{}
	st.setTitle = func(lines []string) {}
	ap := newSchemaAnnotationParser("SomeModel")
	st.annotation = ap
	st.taggers = []tagParser{
		{"Maximum", false, false, nil, &setMaximum{dummyBuilder(), regexp.MustCompile(fmt.Sprintf(rxMaximumFmt, ""))}},
		{"Minimum", false, false, nil, &setMinimum{dummyBuilder(), regexp.MustCompile(fmt.Sprintf(rxMinimumFmt, ""))}},
		{"MultipleOf", false, false, nil, &setMultipleOf{dummyBuilder(), regexp.MustCompile(fmt.Sprintf(rxMultipleOfFmt, ""))}},
	}

	err = st.Parse(ascg(block))
	assert.NoError(t, err)
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
	block := `swagger:model someModel

This has a title without whitespace.
The punctuation here does indeed matter. But it won't for go.

minimum: 10
swagger:meta
maximum: 20
`
	var err error

	st := &sectionedParser{}
	st.setTitle = func(lines []string) {}
	ap := newSchemaAnnotationParser("SomeModel")
	st.annotation = ap
	st.taggers = []tagParser{
		{"Maximum", false, false, nil, &setMaximum{dummyBuilder(), regexp.MustCompile(fmt.Sprintf(rxMaximumFmt, ""))}},
		{"Minimum", false, false, nil, &setMinimum{dummyBuilder(), regexp.MustCompile(fmt.Sprintf(rxMinimumFmt, ""))}},
		{"MultipleOf", false, false, nil, &setMultipleOf{dummyBuilder(), regexp.MustCompile(fmt.Sprintf(rxMultipleOfFmt, ""))}},
	}

	err = st.Parse(ascg(block))
	assert.NoError(t, err)
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

func TestShouldAcceptTag(t *testing.T) {
	var tagTests = []struct {
		tags        []string
		includeTags map[string]bool
		excludeTags map[string]bool
		expected    bool
	}{
		{nil, nil, nil, true},
		{[]string{"app"}, map[string]bool{"app": true}, nil, true},
		{[]string{"app"}, nil, map[string]bool{"app": true}, false},
	}
	for _, tt := range tagTests {
		actual := shouldAcceptTag(tt.tags, tt.includeTags, tt.excludeTags)
		assert.Equal(t, tt.expected, actual)
	}
}

func TestShouldAcceptPkg(t *testing.T) {
	var pkgTests = []struct {
		path        string
		includePkgs []string
		excludePkgs []string
		expected    bool
	}{
		{"", nil, nil, true},
		{"", nil, []string{"app"}, true},
		{"", []string{"app"}, nil, false},
		{"app", []string{"app"}, nil, true},
		{"app", nil, []string{"app"}, false},
		{"vendor/app", []string{"app"}, nil, true},
		{"vendor/app", nil, []string{"app"}, false},
	}
	for _, tt := range pkgTests {
		actual := shouldAcceptPkg(tt.path, tt.includePkgs, tt.excludePkgs)
		assert.Equal(t, tt.expected, actual)
	}
}
