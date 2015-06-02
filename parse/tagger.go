package parse

import (
	"go/ast"
	"regexp"
	"strings"
)

func newMultiLineTagParser(name string, parser valueParser) tagParser {
	return tagParser{
		Name:      name,
		MultiLine: true,
		Parser:    parser,
	}
}

func newSingleLineTagParser(name string, parser valueParser) tagParser {
	return tagParser{
		Name:      name,
		MultiLine: false,
		Parser:    parser,
	}
}

type tagParser struct {
	Name      string
	MultiLine bool
	Lines     []string
	Parser    valueParser
}

func (st *tagParser) Matches(line string) bool {
	return st.Parser.Matches(line)
}

func (st *tagParser) Parse(lines []string) error {
	return st.Parser.Parse(lines)
}

// aggregates lines in header until it sees a tag.
type sectionedParser struct {
	header     []string
	matched    map[string]tagParser
	annotation valueParser

	seenTag        bool
	skipHeader     bool
	setTitle       func([]string)
	setDescription func([]string)
	workedOutTitle bool
	taggers        []tagParser
	currentTagger  *tagParser
	title          []string
	description    []string
}

func (st *sectionedParser) cleanup(lines []string) []string {
	seenLine := -1
	var lastContent int
	var uncommented []string
	for i, v := range lines {
		str := regexp.MustCompile("^[^\\p{L}\\p{N}\\+]*").ReplaceAllString(v, "")
		uncommented = append(uncommented, str)
		if str != "" {
			if seenLine < 0 {
				seenLine = i
			}
			lastContent = i
		}
	}
	return uncommented[seenLine : lastContent+1]
}

func (st *sectionedParser) collectTitleDescription() {
	if st.workedOutTitle {
		return
	}
	if st.setTitle == nil {
		st.header = st.cleanup(st.header)
		return
	}
	hdrs := st.cleanup(st.header)

	st.workedOutTitle = true
	idx := -1
	for i, line := range hdrs {
		if strings.TrimSpace(line) == "" {
			idx = i
			break
		}
	}

	if idx > -1 {

		st.title = hdrs[:idx]
		if len(hdrs) > idx+1 {
			st.header = hdrs[idx+1:]
		} else {
			st.header = nil
		}
		return
	}

	if len(hdrs) > 0 {
		line := hdrs[0]
		if rxPunctuationEnd.MatchString(line) {
			st.title = []string{line}
			st.header = hdrs[1:]
		} else {
			st.header = hdrs
		}
	}
}

func (st *sectionedParser) Title() []string {
	st.collectTitleDescription()
	return st.title
}

func (st *sectionedParser) Description() []string {
	st.collectTitleDescription()
	return st.header
}

func (st *sectionedParser) Parse(doc *ast.CommentGroup) error {
	if doc == nil {
		return nil
	}
COMMENTS:
	for _, c := range doc.List {
		for _, line := range strings.Split(c.Text, "\n") {
			if rxSwaggerAnnotation.MatchString(line) {
				if st.annotation == nil || !st.annotation.Matches(line) {
					break COMMENTS // a new +swagger: annotation terminates this parser
				}

				st.annotation.Parse([]string{line})
				if len(st.header) > 0 {
					st.seenTag = true
				}
				continue
			}

			var matched bool
			for _, tagger := range st.taggers {
				if tagger.Matches(line) {
					st.seenTag = true
					st.currentTagger = &tagger
					matched = true
					break
				}
			}

			if st.currentTagger == nil {
				if !st.skipHeader && !st.seenTag {
					st.header = append(st.header, line)
				}
				// didn't match a tag, moving on
				continue
			}

			if st.currentTagger.MultiLine && matched {
				// the first line of a multiline tagger doesn't count
				continue
			}

			ts, ok := st.matched[st.currentTagger.Name]
			if !ok {
				ts = *st.currentTagger
			}
			ts.Lines = append(ts.Lines, line)
			if st.matched == nil {
				st.matched = make(map[string]tagParser)
			}
			st.matched[st.currentTagger.Name] = ts

			if !st.currentTagger.MultiLine {
				st.currentTagger = nil
			}
		}
	}
	if st.setTitle != nil {
		st.setTitle(st.Title())
	}
	if st.setDescription != nil {
		st.setDescription(st.Description())
	}
	for _, mt := range st.matched {
		if err := mt.Parse(st.cleanup(mt.Lines)); err != nil {
			return err
		}
	}
	return nil
}
