package swag

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type translationSample struct {
	str, out string
}

func titleize(s string) string { return strings.ToTitle(s[:1]) + lower(s[1:]) }

func TestToGoName(t *testing.T) {
	samples := []translationSample{
		{"sample text", "SampleText"},
		{"sample-text", "SampleText"},
		{"sample_text", "SampleText"},
		{"sampleText", "SampleText"},
		{"sample 2 Text", "Sample2Text"},
		{"findThingById", "FindThingByID"},
	}

	for k := range commonInitialisms {
		samples = append(samples,
			translationSample{"sample " + lower(k) + " text", "Sample" + k + "Text"},
			translationSample{"sample-" + lower(k) + "-text", "Sample" + k + "Text"},
			translationSample{"sample_" + lower(k) + "_text", "Sample" + k + "Text"},
			translationSample{"sample" + titleize(k) + "Text", "Sample" + k + "Text"},
			translationSample{"sample " + lower(k), "Sample" + k},
			translationSample{"sample-" + lower(k), "Sample" + k},
			translationSample{"sample_" + lower(k), "Sample" + k},
			translationSample{"sample" + titleize(k), "Sample" + k},
			translationSample{"sample " + titleize(k) + " text", "Sample" + k + "Text"},
			translationSample{"sample-" + titleize(k) + "-text", "Sample" + k + "Text"},
			translationSample{"sample_" + titleize(k) + "_text", "Sample" + k + "Text"},
		)
	}

	for _, sample := range samples {
		assert.Equal(t, sample.out, ToGoName(sample.str))
	}
}

func TestContainsStringsCI(t *testing.T) {
	list := []string{"hello", "world", "and", "such"}

	assert.True(t, ContainsStringsCI(list, "hELLo"))
	assert.True(t, ContainsStringsCI(list, "world"))
	assert.True(t, ContainsStringsCI(list, "AND"))
	assert.False(t, ContainsStringsCI(list, "nuts"))
}

func TestSplitByFormat(t *testing.T) {
	expected := []string{"one", "two", "three"}
	for _, fmt := range []string{"csv", "pipes", "tsv", "ssv", "multi"} {

		var actual []string
		switch fmt {
		case "multi":
			assert.Nil(t, SplitByFormat("", fmt))
			assert.Nil(t, SplitByFormat("blah", fmt))
		case "ssv":
			actual = SplitByFormat(strings.Join(expected, " "), fmt)
			assert.EqualValues(t, expected, actual)
		case "pipes":
			actual = SplitByFormat(strings.Join(expected, "|"), fmt)
			assert.EqualValues(t, expected, actual)
		case "tsv":
			actual = SplitByFormat(strings.Join(expected, "\t"), fmt)
			assert.EqualValues(t, expected, actual)
		default:
			actual = SplitByFormat(strings.Join(expected, ","), fmt)
			assert.EqualValues(t, expected, actual)
		}
	}
}

func TestJoinByFormat(t *testing.T) {
	for _, fmt := range []string{"csv", "pipes", "tsv", "ssv", "multi"} {

		lval := []string{"one", "two", "three"}
		var expected []string
		switch fmt {
		case "multi":
			expected = lval
		case "ssv":
			expected = []string{strings.Join(lval, " ")}
		case "pipes":
			expected = []string{strings.Join(lval, "|")}
		case "tsv":
			expected = []string{strings.Join(lval, "\t")}
		default:
			expected = []string{strings.Join(lval, ",")}
		}
		assert.Nil(t, JoinByFormat(nil, fmt))
		assert.EqualValues(t, expected, JoinByFormat(lval, fmt))
	}
}

func TestToFileName(t *testing.T) {
	samples := []translationSample{
		{"SampleText", "sample_text"},
		{"FindThingByID", "find_thing_by_id"},
	}

	for k := range commonInitialisms {
		samples = append(samples,
			translationSample{"Sample" + k + "Text", "sample_" + lower(k) + "_text"},
		)
	}

	for _, sample := range samples {
		assert.Equal(t, sample.out, ToFileName(sample.str))
	}
}

func TestToCommandName(t *testing.T) {
	samples := []translationSample{
		{"SampleText", "sample-text"},
		{"FindThingByID", "find-thing-by-id"},
	}

	for k := range commonInitialisms {
		samples = append(samples,
			translationSample{"Sample" + k + "Text", "sample-" + lower(k) + "-text"},
		)
	}

	for _, sample := range samples {
		assert.Equal(t, sample.out, ToCommandName(sample.str))
	}
}

func TestToHumanName(t *testing.T) {
	samples := []translationSample{
		{"SampleText", "sample text"},
		{"FindThingByID", "find thing by ID"},
	}

	for k := range commonInitialisms {
		samples = append(samples,
			translationSample{"Sample" + k + "Text", "sample " + k + " text"},
		)
	}

	for _, sample := range samples {
		assert.Equal(t, sample.out, ToHumanNameLower(sample.str))
	}
}

func TestToJSONName(t *testing.T) {
	samples := []translationSample{
		{"SampleText", "sampleText"},
		{"FindThingByID", "findThingById"},
	}

	for k := range commonInitialisms {
		samples = append(samples,
			translationSample{"Sample" + k + "Text", "sample" + titleize(k) + "Text"},
		)
	}

	for _, sample := range samples {
		assert.Equal(t, sample.out, ToJSONName(sample.str))
	}
}
