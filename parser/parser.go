package parser

import (
	"fmt"
	"go/ast"
	"log"
	"regexp"
	"strings"

	"github.com/casualjim/go-swagger/spec"
	"github.com/casualjim/go-swagger/util"
)

// Many thanks go to https://github.com/yvasiyarov/swagger
// this is loosely based on that implementation but for swagger 2.0

type setter func(interface{}, []string) error

// apiParser the global context for parsing a go application
// into a swagger specification
type apiParser struct {
	BasePath string
}

// newAPIParser creates a new api parser
func newAPIParser(bp string) *apiParser {
	return &apiParser{BasePath: bp}
}

// MustExpandPackagePath gets the real package path on disk
func (a *apiParser) MustExpandPackagePath(packagePath string) string {
	pkgRealpath := util.FindInGoSearchPath(packagePath)
	if pkgRealpath == "" {
		log.Fatalf("Can't find package %s \n", packagePath)
	}

	return pkgRealpath
}

type taggedSection struct {
	Name  string
	Lines []string
}

func tagNameMatcher(name string) *regexp.Regexp {
	return regexp.MustCompile(fmt.Sprintf("[^\\w]*(?:%s)[^:]*:[^\\w]*", name))
}

var stripCommentsExpr = regexp.MustCompile("^[^\\w]*")

func stripComment(line string) string {
	return stripCommentsExpr.ReplaceAllString(line, "")
}

func newSectionTagger(name string, multiLine bool) *sectionTagger {
	return &sectionTagger{
		taggedSection: taggedSection{
			Name: name,
		},
		matcher:   tagNameMatcher(name),
		Multiline: multiLine,
		isFirst:   true,
	}
}

// a section tagger analyzes comment lines
// and groups them together by tag
type sectionTagger struct {
	taggedSection
	matcher *regexp.Regexp
	// when true this section needs to read multiple lines
	// a section terminates when:
	//   * another tag is found
	//   * 2 consecutive new lines are found
	//   * the comment group ends
	// this applies to title, description and TOS for the info object for example
	Multiline bool

	// starts out being true
	isFirst  bool
	wasEmpty bool
	set      setter
}
type unmatchedSection struct{} // marker struct

type singleLineSection struct {
	taggedSection
}

func (sls singleLineSection) Line() string {
	if len(sls.Lines) == 0 {
		return ""
	}
	return sls.Lines[0]
}

// grabs lines to aggregate
type multiLineSectionPart struct {
	taggedSection
}

// stops and discards the last line
type multiLineSectionTerminator struct {
	taggedSection
}

// this means that whoever is aggregating the tagged section needs to move on to the next
// section tagger but probably not advance the lines
type newTagSectionTerminator struct {
	taggedSection
}

// not exported, ambiguous return type
func (st *sectionTagger) Tag(text string, terminatingTags []string) interface{} {
	if st.isFirst && st.matcher.MatchString(text) {
		st.isFirst = false
		txt := strings.TrimSpace(stripComment(st.matcher.ReplaceAllString(text, "")))

		if !st.Multiline {
			st.Lines = []string{txt}
			return singleLineSection{st.taggedSection}
		}
		if len(txt) > 0 {
			st.Lines = []string{txt}
		}
		return multiLineSectionPart{st.taggedSection}
	}

	if !st.isFirst && st.Multiline {
		txt := strings.TrimSpace(stripComment(text))
		isEmpty := txt == ""
		if isEmpty && st.wasEmpty {
			return multiLineSectionTerminator{st.taggedSection} // terminated
		}
		st.wasEmpty = isEmpty
		// check for tags that terminate
		if len(terminatingTags) > 0 {
			terminates := tagNameMatcher(strings.Join(terminatingTags, "|"))
			if terminates.MatchString(text) {
				return newTagSectionTerminator{st.taggedSection}
			}
		}
		if len(st.Lines) > 0 || !isEmpty {
			st.Lines = append(st.Lines, txt)
		}
		return multiLineSectionPart{st.taggedSection}
	}
	return unmatchedSection{}
}

func newDocCommentParser(otherTags []string, taggers ...*sectionTagger) *docCommentParser {
	return &docCommentParser{taggers: taggers, otherTags: otherTags}
}

type docCommentParser struct {
	taggers   []*sectionTagger
	otherTags []string
}

func (ai *docCommentParser) Parse(gofile *ast.File) (*spec.Info, error) {
	info := new(spec.Info)

	// var currentLines []string
	var selectedTagger *sectionTagger
	var otherTags []string
	taggers := ai.taggers
	for _, c := range gofile.Doc.List {
		text := c.Text
		lines := strings.Split(text, "\n")

	LINES:
		for _, line := range lines {
			// this is an aggregating tagger
			if selectedTagger != nil {
				switch res := selectedTagger.Tag(line, otherTags).(type) {
				case multiLineSectionPart:
					continue LINES
				case multiLineSectionTerminator:
					if err := selectedTagger.set(info, res.taggedSection.Lines); err != nil {
						return nil, err
					}
					selectedTagger = nil
					continue LINES
				case newTagSectionTerminator:
					if err := selectedTagger.set(info, res.taggedSection.Lines); err != nil {
						return nil, err
					}
				}
			}

			selectedTagger = nil
			for i, tagger := range taggers {
				switch res := tagger.Tag(line, nil).(type) {
				case singleLineSection:
					if err := tagger.set(info, res.taggedSection.Lines); err != nil {
						return nil, err
					}
					// once it has matched we don't care for probing for it again
					taggers = append(taggers[:i], taggers[i+1:]...)
					continue LINES

				case multiLineSectionPart:
					selectedTagger = tagger
					otherTags = ai.otherTags
					for _, t := range ai.taggers {
						if t.Name != tagger.Name {
							otherTags = append(otherTags, t.Name)
						}
					}
					// once it has matched we don't care for probing for it again
					taggers = append(taggers[:i], taggers[i+1:]...)
					continue LINES

				case unmatchedSection:
					// TODO: something slightly smarter than nothing???
				}
			}
		}
	}

	if selectedTagger != nil {
		if err := selectedTagger.set(info, selectedTagger.Lines); err != nil {
			return nil, err
		}
	}

	return info, nil
}
