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
)

// CharClass represents a regular expression character class as a list of ranges.
// The runes contained in the class can be accessed by index.
type tCharClass struct {
	Ranges    []tCharClassRange
	TotalSize int32
}

// CharClassRange represents a single range of characters in a character class.
type tCharClassRange struct {
	Start rune
	Size  int32
}

// NewCharClass creates a character class with a single range.
func newCharClass(start rune, end rune) *tCharClass {
	charRange := newCharClassRange(start, end)
	return &tCharClass{
		Ranges:    []tCharClassRange{charRange},
		TotalSize: charRange.Size,
	}
}

/*
ParseCharClass parses a character class as represented by syntax.Parse into a slice of CharClassRange structs.

Char classes are encoded as pairs of runes representing ranges:
[0-9] = 09, [a0] = aa00 (2 1-len ranges).

e.g.

"[a0-9]" -> "aa09" -> a, 0-9

"[^a-z]" -> "…" -> 0-(a-1), (z+1)-(max rune)
*/
func parseCharClass(runes []rune) *tCharClass {
	var totalSize int32
	numRanges := len(runes) / 2
	ranges := make([]tCharClassRange, numRanges, numRanges)

	for i := 0; i < numRanges; i++ {
		start := runes[i*2]
		end := runes[i*2+1]

		// indicates a negative class
		if start == 0 {
			// doesn't make sense to generate null bytes, so all ranges must start at
			// no less than 1.
			start = 1
		}

		r := newCharClassRange(start, end)

		ranges[i] = r
		totalSize += r.Size
	}

	return &tCharClass{ranges, totalSize}
}

// GetRuneAt gets a rune from CharClass as a contiguous array of runes.
func (class *tCharClass) GetRuneAt(i int32) rune {
	for _, r := range class.Ranges {
		if i < r.Size {
			return r.Start + rune(i)
		}
		i -= r.Size
	}
	panic("index out of bounds")
}

func (class *tCharClass) String() string {
	return fmt.Sprintf("%s", class.Ranges)
}

func newCharClassRange(start rune, end rune) tCharClassRange {
	if start < 1 {
		panic("char class range cannot contain runes less than 1")
	}

	size := end - start + 1

	if size < 1 {
		panic("char class range size must be at least 1")
	}

	return tCharClassRange{
		Start: start,
		Size:  size,
	}
}

func (r tCharClassRange) String() string {
	if r.Size == 1 {
		return fmt.Sprintf("%s:1", runesToString(r.Start))
	}
	return fmt.Sprintf("%s-%s:%d", runesToString(r.Start), runesToString(r.Start+rune(r.Size-1)), r.Size)

}
