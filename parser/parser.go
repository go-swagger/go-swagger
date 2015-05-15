package parser

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/casualjim/go-swagger/util"
)

// Many thanks go to https://github.com/yvasiyarov/swagger
// this is loosely based on that implementation but for swagger 2.0

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
