/*
Copyright 2014 Zachary Klippenstein

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package regen

import (
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"regexp/syntax"
	"testing"

	"github.com/google/gxui/math"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
)

const (
	// Each expression is generated and validated this many times.
	SampleSize = 999

	// Arbitrary limit in the standard package.
	// See https://golang.org/src/regexp/syntax/parse.go?s=18885:18935#L796
	MaxSupportedRepeatCount = 1000
)

func ExampleGenerate() {
	pattern := "[ab]{5}"
	str, _ := Generate(pattern)

	if matched, _ := regexp.MatchString(pattern, str); matched {
		fmt.Println("Matches!")
	}

	// Output:
	// Matches!
}

func ExampleNewGenerator() {
	pattern := "[ab]{5}"

	generator, _ := NewGenerator(pattern, &GeneratorArgs{
		RngSource: rand.NewSource(0),
	})

	str := generator.Generate()

	if matched, _ := regexp.MatchString(pattern, str); matched {
		fmt.Println("Matches!")
	}

	// Output:
	// Matches!
}

func ExampleNewGenerator_perl() {
	pattern := `\d{5}`

	generator, _ := NewGenerator(pattern, &GeneratorArgs{
		Flags: syntax.Perl,
	})

	str := generator.Generate()

	if matched, _ := regexp.MatchString("[[:digit:]]{5}", str); matched {
		fmt.Println("Matches!")
	}
	// Output:
	// Matches!
}

func ExampleCaptureGroupHandler() {
	pattern := `Hello, (?P<firstname>[A-Z][a-z]{2,10}) (?P<lastname>[A-Z][a-z]{2,10})`

	generator, _ := NewGenerator(pattern, &GeneratorArgs{
		Flags: syntax.Perl,
		CaptureGroupHandler: func(index int, name string, group *syntax.Regexp, generator Generator, args *GeneratorArgs) string {
			if name == "firstname" {
				return fmt.Sprintf("FirstName (e.g. %s)", generator.Generate())
			}
			return fmt.Sprintf("LastName (e.g. %s)", generator.Generate())
		},
	})

	// Print to stderr since we're generating random output and can't assert equality.
	fmt.Fprintln(os.Stderr, generator.Generate())

	// Needed for "go test" to run this example. (Must be a blank line before.)

	// Output:
}

func TestGeneratorArgs(t *testing.T) {
	t.Parallel()

	Convey("initialize", t, func() {
		Convey("Handles empty struct", func() {
			args := GeneratorArgs{}

			var err error
			So(func() { err = args.initialize() }, ShouldNotPanic)
			So(err, ShouldBeNil)
		})

		Convey("Unicode groups not supported", func() {
			args := &GeneratorArgs{
				Flags: syntax.UnicodeGroups,
			}

			err := args.initialize()
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "UnicodeGroups not supported")
		})

		Convey("Panics if repeat bounds are invalid", func() {
			args := &GeneratorArgs{
				MinUnboundedRepeatCount: 2,
				MaxUnboundedRepeatCount: 1,
			}

			So(func() { args.initialize() },
				ShouldPanicWith,
				"MinUnboundedRepeatCount(2) > MaxUnboundedRepeatCount(1)")
		})

		Convey("Allows equal repeat bounds", func() {
			args := &GeneratorArgs{
				MinUnboundedRepeatCount: 1,
				MaxUnboundedRepeatCount: 1,
			}

			var err error
			So(func() { err = args.initialize() }, ShouldNotPanic)
			So(err, ShouldBeNil)
		})
	})

	Convey("Rng", t, func() {
		Convey("Panics if called before initialization", func() {
			args := GeneratorArgs{}
			So(func() { args.Rng() }, ShouldPanic)
		})

		Convey("Non-nil after initialization", func() {
			args := GeneratorArgs{}
			err := args.initialize()
			So(err, ShouldBeNil)
			So(args.Rng(), ShouldNotBeNil)
		})
	})
}

func TestNewGenerator(t *testing.T) {
	t.Parallel()

	Convey("NewGenerator", t, func() {

		Convey("Handles nil GeneratorArgs", func() {
			generator, err := NewGenerator("", nil)
			So(generator, ShouldNotBeNil)
			So(err, ShouldBeNil)
		})

		Convey("Handles empty GeneratorArgs", func() {
			generator, err := NewGenerator("", &GeneratorArgs{})
			So(generator, ShouldNotBeNil)
			So(err, ShouldBeNil)
		})

		Convey("Forwards errors from args initialization", func() {
			args := &GeneratorArgs{
				Flags: syntax.UnicodeGroups,
			}

			_, err := NewGenerator("", args)
			So(err, ShouldNotBeNil)
		})
	})
}

func TestGenEmpty(t *testing.T) {
	t.Parallel()

	Convey("Empty", t, func() {
		args := &GeneratorArgs{
			RngSource: rand.NewSource(0),
		}
		ConveyGeneratesStringMatching(args, "", "^$")
	})
}

func TestGenLiterals(t *testing.T) {
	t.Parallel()

	Convey("Literals", t, func() {
		ConveyGeneratesStringMatchingItself(nil,
			"a",
			"abc",
		)
	})
}

func TestGenDotNotNl(t *testing.T) {
	t.Parallel()

	Convey("DotNotNl", t, func() {
		ConveyGeneratesStringMatchingItself(nil, ".")

		Convey("No newlines are generated", func() {
			generator, _ := NewGenerator(".", nil)

			// Not a very strong assertion, but not sure how to do better. Exploring the entire
			// generation space (2^32) takes far too long for a unit test.
			for i := 0; i < SampleSize; i++ {
				So(generator.Generate(), ShouldNotContainSubstring, "\n")
			}
		})
	})
}

func TestGenStringStartEnd(t *testing.T) {
	t.Parallel()

	Convey("String start/end", t, func() {
		args := &GeneratorArgs{
			RngSource: rand.NewSource(0),
			Flags:     0,
		}

		ConveyGeneratesStringMatching(args, `^abc$`, `^abc$`)
		ConveyGeneratesStringMatching(args, `$abc^`, `^abc$`)
		ConveyGeneratesStringMatching(args, `a^b$c`, `^abc$`)
	})
}

func TestGenQuestionMark(t *testing.T) {
	t.Parallel()

	Convey("QuestionMark", t, func() {
		ConveyGeneratesStringMatchingItself(nil,
			"a?",
			"(abc)?",
			"[ab]?",
			".?")
	})
}

func TestGenPlus(t *testing.T) {
	t.Parallel()

	Convey("Plus", t, func() {
		ConveyGeneratesStringMatchingItself(nil, "a+")
	})
}

func TestGenStar(t *testing.T) {
	t.Parallel()

	Convey("Star", t, func() {
		ConveyGeneratesStringMatchingItself(nil, "a*")

		Convey("HitsDefaultMin", func() {
			regexp := "a*"
			args := &GeneratorArgs{
				RngSource: rand.NewSource(0),
			}
			counts := generateLenHistogram(regexp, DefaultMaxUnboundedRepeatCount, args)

			So(counts[0], ShouldBeGreaterThan, 0)
		})

		Convey("HitsCustomMin", func() {
			regexp := "a*"
			args := &GeneratorArgs{
				RngSource:               rand.NewSource(0),
				MinUnboundedRepeatCount: 200,
			}
			counts := generateLenHistogram(regexp, DefaultMaxUnboundedRepeatCount, args)

			So(counts[200], ShouldBeGreaterThan, 0)
			for i := 0; i < 200; i++ {
				So(counts[i], ShouldEqual, 0)
			}
		})

		Convey("HitsDefaultMax", func() {
			regexp := "a*"
			args := &GeneratorArgs{
				RngSource: rand.NewSource(0),
			}
			counts := generateLenHistogram(regexp, DefaultMaxUnboundedRepeatCount, args)

			So(len(counts), ShouldEqual, DefaultMaxUnboundedRepeatCount+1)
			So(counts[DefaultMaxUnboundedRepeatCount], ShouldBeGreaterThan, 0)
		})

		Convey("HitsCustomMax", func() {
			regexp := "a*"
			args := &GeneratorArgs{
				RngSource:               rand.NewSource(0),
				MaxUnboundedRepeatCount: 200,
			}
			counts := generateLenHistogram(regexp, 200, args)

			So(len(counts), ShouldEqual, 200+1)
			So(counts[200], ShouldBeGreaterThan, 0)
		})
	})
}

func TestGenCharClassNotNl(t *testing.T) {
	t.Parallel()

	Convey("CharClassNotNl", t, func() {
		ConveyGeneratesStringMatchingItself(nil,
			"[a]",
			"[abc]",
			"[a-d]",
			"[ac]",
			"[0-9]",
			"[a-z0-9]",
		)

		Convey("No newlines are generated", func() {
			// Try to narrow down the generation space. Still not a very strong assertion.
			generator, _ := NewGenerator("[^a-zA-Z0-9]", nil)
			for i := 0; i < SampleSize; i++ {
				assert.NotEqual(t, "\n", generator.Generate())
			}
		})
	})
}

func TestGenNegativeCharClass(t *testing.T) {
	t.Parallel()

	Convey("NegativeCharClass", t, func() {
		ConveyGeneratesStringMatchingItself(nil, "[^a-zA-Z0-9]")
	})
}

func TestGenAlternate(t *testing.T) {
	t.Parallel()

	Convey("Alternate", t, func() {
		ConveyGeneratesStringMatchingItself(nil,
			"a|b",
			"abc|def|ghi",
			"[ab]|[cd]",
			"foo|bar|baz", // rewrites to foo|ba[rz]
		)
	})
}

func TestGenCapture(t *testing.T) {
	t.Parallel()

	Convey("Capture", t, func() {
		ConveyGeneratesStringMatching(nil, "(abc)", "^abc$")
		ConveyGeneratesStringMatching(nil, "()", "^$")
	})
}

func TestGenConcat(t *testing.T) {
	t.Parallel()

	Convey("Concat", t, func() {
		ConveyGeneratesStringMatchingItself(nil, "[ab][cd]")
	})
}

func TestGenRepeat(t *testing.T) {
	t.Parallel()

	Convey("Repeat", t, func() {

		Convey("Unbounded", func() {
			ConveyGeneratesStringMatchingItself(nil, `a{1,}`)

			Convey("HitsDefaultMax", func() {
				regexp := "a{0,}"
				args := &GeneratorArgs{
					RngSource: rand.NewSource(0),
				}
				counts := generateLenHistogram(regexp, DefaultMaxUnboundedRepeatCount, args)

				So(len(counts), ShouldEqual, DefaultMaxUnboundedRepeatCount+1)
				So(counts[DefaultMaxUnboundedRepeatCount], ShouldBeGreaterThan, 0)
			})

			Convey("HitsCustomMax", func() {
				regexp := "a{0,}"
				args := &GeneratorArgs{
					RngSource:               rand.NewSource(0),
					MaxUnboundedRepeatCount: 200,
				}
				counts := generateLenHistogram(regexp, 200, args)

				So(len(counts), ShouldEqual, 200+1)
				So(counts[200], ShouldBeGreaterThan, 0)
			})
		})

		Convey("HitsMin", func() {
			regexp := "a{0,3}"
			args := &GeneratorArgs{
				RngSource: rand.NewSource(0),
			}
			counts := generateLenHistogram(regexp, 3, args)

			So(len(counts), ShouldEqual, 3+1)
			So(counts[0], ShouldBeGreaterThan, 0)
		})

		Convey("HitsMax", func() {
			regexp := "a{0,3}"
			args := &GeneratorArgs{
				RngSource: rand.NewSource(0),
			}
			counts := generateLenHistogram(regexp, 3, args)

			So(len(counts), ShouldEqual, 3+1)
			So(counts[3], ShouldBeGreaterThan, 0)
		})

		Convey("IsWithinBounds", func() {
			regexp := "a{5,10}"
			args := &GeneratorArgs{
				RngSource: rand.NewSource(0),
			}
			counts := generateLenHistogram(regexp, 10, args)

			So(len(counts), ShouldEqual, 11)

			for i := 0; i < 11; i++ {
				if i < 5 {
					So(counts[i], ShouldEqual, 0)
				} else if i < 11 {
					So(counts[i], ShouldBeGreaterThan, 0)
				}
			}
		})
	})
}

func TestGenCharClasses(t *testing.T) {
	t.Parallel()

	Convey("CharClasses", t, func() {

		Convey("Ascii", func() {
			ConveyGeneratesStringMatchingItself(nil,
				"[[:alnum:]]",
				"[[:alpha:]]",
				"[[:ascii:]]",
				"[[:blank:]]",
				"[[:cntrl:]]",
				"[[:digit:]]",
				"[[:graph:]]",
				"[[:lower:]]",
				"[[:print:]]",
				"[[:punct:]]",
				"[[:space:]]",
				"[[:upper:]]",
				"[[:word:]]",
				"[[:xdigit:]]",
				"[[:^alnum:]]",
				"[[:^alpha:]]",
				"[[:^ascii:]]",
				"[[:^blank:]]",
				"[[:^cntrl:]]",
				"[[:^digit:]]",
				"[[:^graph:]]",
				"[[:^lower:]]",
				"[[:^print:]]",
				"[[:^punct:]]",
				"[[:^space:]]",
				"[[:^upper:]]",
				"[[:^word:]]",
				"[[:^xdigit:]]",
			)
		})

		Convey("Perl", func() {
			args := &GeneratorArgs{
				Flags: syntax.Perl,
			}

			ConveyGeneratesStringMatchingItself(args,
				`\d`,
				`\s`,
				`\w`,
				`\D`,
				`\S`,
				`\W`,
			)
		})
	})
}

func TestCaptureGroupHandler(t *testing.T) {
	t.Parallel()

	Convey("CaptureGroupHandler", t, func() {
		callCount := 0

		gen, err := NewGenerator(`(?:foo) (bar) (?P<name>baz)`, &GeneratorArgs{
			Flags: syntax.PerlX,
			CaptureGroupHandler: func(index int, name string, group *syntax.Regexp, generator Generator, args *GeneratorArgs) string {
				callCount++

				So(index, ShouldBeLessThan, 2)

				if index == 0 {
					So(name, ShouldEqual, "")
					So(group.String(), ShouldEqual, "bar")
					So(generator.Generate(), ShouldEqual, "bar")
					return "one"
				}

				// Index 1
				So(name, ShouldEqual, "name")
				So(group.String(), ShouldEqual, "baz")
				So(generator.Generate(), ShouldEqual, "baz")
				return "two"
			},
		})
		So(err, ShouldBeNil)

		So(gen.Generate(), ShouldEqual, "foo one two")
		So(callCount, ShouldEqual, 2)
	})
}

func ConveyGeneratesStringMatchingItself(args *GeneratorArgs, patterns ...string) {
	for _, pattern := range patterns {
		Convey(fmt.Sprintf("String generated from /%s/ matches itself", pattern), func() {
			So(pattern, ShouldGenerateStringMatching, pattern, args)
		})
	}
}

func ConveyGeneratesStringMatching(args *GeneratorArgs, pattern string, expectedPattern string) {
	Convey(fmt.Sprintf("String generated from /%s/ matches /%s/", pattern, expectedPattern), func() {
		So(pattern, ShouldGenerateStringMatching, expectedPattern, args)
	})
}

func ShouldGenerateStringMatching(actual interface{}, expected ...interface{}) string {
	return ShouldGenerateStringMatchingTimes(actual, expected[0], expected[1], SampleSize)
}

func ShouldGenerateStringMatchingTimes(actual interface{}, expected ...interface{}) string {
	pattern := actual.(string)
	expectedPattern := expected[0].(string)
	args := expected[1].(*GeneratorArgs)
	times := expected[2].(int)

	generator, err := NewGenerator(pattern, args)
	if err != nil {
		panic(err)
	}

	for i := 0; i < times; i++ {
		result := generator.Generate()
		matched, err := regexp.MatchString(expectedPattern, result)
		if err != nil {
			panic(err)
		}
		if !matched {
			return fmt.Sprintf("string “%s” generated from /%s/ did not match /%s/.",
				result, pattern, expectedPattern)
		}
	}

	return ""
}

func generateLenHistogram(regexp string, maxLen int, args *GeneratorArgs) (counts []int) {
	generator, err := NewGenerator(regexp, args)
	if err != nil {
		panic(err)
	}

	iterations := math.Max(maxLen*4, SampleSize)

	for i := 0; i < iterations; i++ {
		str := generator.Generate()

		// Grow the slice if necessary.
		if len(str) >= len(counts) {
			newCounts := make([]int, len(str)+1)
			copy(newCounts, counts)
			counts = newCounts
		}

		counts[len(str)]++
	}

	return
}
