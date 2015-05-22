package parser

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	rxStripComments      = regexp.MustCompile("^[^\\w]*")
	rxStripTitleComments = regexp.MustCompile("^[^\\w]*(:?P|p)ackage\\s+\\w+[^\\w]*")
)

func tagNamesMatcher(tags []string) *regexp.Regexp {
	var escaped []string
	for _, t := range tags {
		escaped = append(escaped, regexp.QuoteMeta(t))
	}
	return unsafeTagNameMatcher(strings.Join(escaped, "|"))
}

func tagNameMatcher(name string) *regexp.Regexp {
	return unsafeTagNameMatcher(regexp.QuoteMeta(name))
}

func unsafeTagNameMatcher(name string) *regexp.Regexp {
	return regexp.MustCompile(fmt.Sprintf("[^\\w]*(?:%s)[^:]*:\\s*", name))
}

func swaggerClassifier(name string) *regexp.Regexp {
	return regexp.MustCompile(fmt.Sprintf("[^+]*\\+\\s*swagger:%s", regexp.QuoteMeta(name)))
}

type taggedSection struct {
	Name  string
	Lines []string
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

func newSectionTagger(name string, multiLine bool) *sectionTagger {
	return &sectionTagger{
		taggedSection: taggedSection{
			Name: name,
		},
		matcher:         tagNameMatcher(name),
		Multiline:       multiLine,
		isFirst:         true,
		stripsTag:       true,
		rxStripComments: rxStripComments,
	}
}

// a title matcher is a special section tagger that collects until the first open line, but doesn't really
// have a tag name. It expects that the first line of the comment flagged with swagger:meta is the title
func newTitleTagger() *sectionTagger {
	return &sectionTagger{
		taggedSection: taggedSection{
			Name: "Title",
		},
		matcher:         regexp.MustCompile(".*"),
		Multiline:       true,
		isFirst:         true,
		wasEmpty:        true,
		rxStripComments: rxStripTitleComments,
	}
}

// a description matcher is a special section tagger that collects after the first paragraph and until
// the first open line, but doesn't really have a tag name. It expects everything but the first paragraph
// of the comment before any tags appear is the description
func newDescriptionTagger() *sectionTagger {
	return &sectionTagger{
		taggedSection: taggedSection{
			Name: "Description",
		},
		matcher:         regexp.MustCompile(".*"),
		Multiline:       true,
		isFirst:         true,
		rxStripComments: rxStripComments,
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
	isFirst   bool
	wasEmpty  bool
	stripsTag bool
	set       setter

	rxStripComments *regexp.Regexp
}

// not exported, ambiguous return type
func (st *sectionTagger) Tag(text string, terminatingTags []string) interface{} {
	if st.isFirst && st.matcher.MatchString(text) {
		st.isFirst = false

		txt := strings.TrimSpace(st.stripComment(text))
		if st.stripsTag {
			txt = strings.TrimSpace(st.stripComment(st.matcher.ReplaceAllString(text, "")))
		}

		if !st.Multiline {
			st.Lines = []string{txt}
			return singleLineSection{st.taggedSection}
		}
		if len(txt) > 0 {
			st.Lines = []string{txt}
		}
		// so far only title wants to have a different pattern
		// and this is only for the first line, after that we replace
		// the strip comments with the default pattern
		st.rxStripComments = rxStripComments
		return multiLineSectionPart{st.taggedSection}
	}

	if !st.isFirst && st.Multiline {
		txt := strings.TrimSpace(st.stripComment(text))
		isEmpty := txt == ""

		if isEmpty && st.wasEmpty {
			return multiLineSectionTerminator{st.taggedSection} // terminated
		}
		st.wasEmpty = isEmpty

		// check for tags that terminate
		if len(terminatingTags) > 0 {
			if tagNamesMatcher(terminatingTags).MatchString(text) {
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
func (st *sectionTagger) stripComment(line string) string {
	return st.rxStripComments.ReplaceAllString(line, "")
}
