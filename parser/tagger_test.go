package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSingleLineTag(t *testing.T) {

	lines := []string{
		"// Version: 0.0.1",
		"//Version: 0.0.1",
		"Version: 0.0.1",
		"  Version: 0.0.1",
		"// Version : 0.0.1",
		"//Version : 0.0.1",
		"Version : 0.0.1",
		"  Version : 0.0.1",
	}

	for _, line := range lines {
		tagger := newSectionTagger("Version", false)
		ts := tagger.Tag(line, nil)
		assert.IsType(t, singleLineSection{}, ts)
		tss := ts.(singleLineSection)
		assert.Equal(t, "Version", tss.Name)
		assert.Equal(t, "0.0.1", tss.Line())
	}

	invalid := []string{
		"Version",
		"Versoin: 0.0.1",
	}
	for _, line := range invalid {
		tagger := newSectionTagger("Version", false)
		ts := tagger.Tag(line, nil)
		assert.IsType(t, unmatchedSection{}, ts)
	}
}

func TestMultilineTag_DoubleSpaceTerminator(t *testing.T) {
	lines := []string{
		"// Description: some content here",
		"// and also on another line",
		"// ",
		"// there is a linebreak in this one",
		"// ",
		"// ",
		"// this line should not be included",
	}
	parsed := []string{
		"some content here",
		"and also on another line",
		"",
		"there is a linebreak in this one",
		"",
		"Version: 0.0.1",
	}

	var collected multiLineSectionPart
	var res interface{}
	tagger := newSectionTagger("Description", true)
	res = tagger.Tag(lines[0], nil)
	assert.IsType(t, multiLineSectionPart{}, res)
	collected = res.(multiLineSectionPart)
	assert.EqualValues(t, parsed[:1], collected.Lines)

	res = tagger.Tag(lines[1], nil)
	assert.IsType(t, multiLineSectionPart{}, res)
	collected = res.(multiLineSectionPart)
	assert.EqualValues(t, parsed[:2], collected.Lines)

	res = tagger.Tag(lines[2], nil)
	assert.IsType(t, multiLineSectionPart{}, res)
	collected = res.(multiLineSectionPart)
	assert.EqualValues(t, parsed[:3], collected.Lines)

	res = tagger.Tag(lines[3], nil)
	assert.IsType(t, multiLineSectionPart{}, res)
	collected = res.(multiLineSectionPart)
	assert.EqualValues(t, parsed[:4], collected.Lines)

	res = tagger.Tag(lines[4], nil)
	assert.IsType(t, multiLineSectionPart{}, res)
	collected = res.(multiLineSectionPart)
	assert.EqualValues(t, parsed[:5], collected.Lines)

	res = tagger.Tag(lines[5], nil)
	assert.IsType(t, multiLineSectionTerminator{}, res)
}

func TestMultilineTag_NextTagTerminator(t *testing.T) {
	lines := []string{
		"// Description: some content here",
		"// and also on another line",
		"// ",
		"// there is a linebreak in this one",
		"// ",
		"// Version: 0.0.1",
	}
	parsed := []string{
		"some content here",
		"and also on another line",
		"",
		"there is a linebreak in this one",
		"",
		"Version: 0.0.1",
	}

	var collected multiLineSectionPart
	var res interface{}
	tagger := newSectionTagger("Description", true)
	res = tagger.Tag(lines[0], []string{"Version"})
	assert.IsType(t, multiLineSectionPart{}, res)
	collected = res.(multiLineSectionPart)
	assert.EqualValues(t, parsed[:1], collected.Lines)

	res = tagger.Tag(lines[1], []string{"Version"})
	assert.IsType(t, multiLineSectionPart{}, res)
	collected = res.(multiLineSectionPart)
	assert.EqualValues(t, parsed[:2], collected.Lines)

	res = tagger.Tag(lines[2], []string{"Version"})
	assert.IsType(t, multiLineSectionPart{}, res)
	collected = res.(multiLineSectionPart)
	assert.EqualValues(t, parsed[:3], collected.Lines)

	res = tagger.Tag(lines[3], []string{"Version"})
	assert.IsType(t, multiLineSectionPart{}, res)
	collected = res.(multiLineSectionPart)
	assert.EqualValues(t, parsed[:4], collected.Lines)

	res = tagger.Tag(lines[4], []string{"Version"})
	assert.IsType(t, multiLineSectionPart{}, res)
	collected = res.(multiLineSectionPart)
	assert.EqualValues(t, parsed[:5], collected.Lines)

	res = tagger.Tag(lines[5], []string{"Version"})
	assert.IsType(t, newTagSectionTerminator{}, res)
}
