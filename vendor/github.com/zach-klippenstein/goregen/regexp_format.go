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
	"bytes"
	"fmt"
	"io"
	"regexp/syntax"
)

// inspectRegexpToString returns a string describing a regular expression.
func inspectRegexpToString(r *syntax.Regexp) string {
	var buffer bytes.Buffer
	inspectRegexpToWriter(&buffer, r)
	return buffer.String()
}

// inspectPatternsToString returns a string describing one or more regular expressions.
func inspectPatternsToString(simplify bool, patterns ...string) string {
	var buffer bytes.Buffer
	for _, pattern := range patterns {
		inspectPatternsToWriter(simplify, &buffer, pattern)
	}
	return buffer.String()
}
func inspectPatternsToWriter(simplify bool, w io.Writer, patterns ...string) {
	for _, pattern := range patterns {
		inspectRegexpToWriter(w, parseOrPanic(simplify, pattern))
	}
}

func inspectRegexpToWriter(w io.Writer, r ...*syntax.Regexp) {
	for _, regexp := range r {
		inspectWithIndent(regexp, "", w)
	}
}

func inspectWithIndent(r *syntax.Regexp, indent string, w io.Writer) {
	fmt.Fprintf(w, "%s{\n", indent)
	fmt.Fprintf(w, "%s  Op: %s\n", indent, opToString(r.Op))
	fmt.Fprintf(w, "%s  Flags: %x\n", indent, r.Flags)
	if len(r.Sub) > 0 {
		fmt.Fprintf(w, "%s  Sub: [\n", indent)
		for _, subR := range r.Sub {
			inspectWithIndent(subR, indent+"    ", w)
		}
		fmt.Fprintf(w, "%s  ]\n", indent)
	} else {
		fmt.Fprintf(w, "%s  Sub: []\n", indent)
	}
	fmt.Fprintf(w, "%s  Rune: %s (%s)\n", indent, runesToString(r.Rune...), runesToDecimalString(r.Rune))
	fmt.Fprintf(w, "%s  [Min, Max]: [%d, %d]\n", indent, r.Min, r.Max)
	fmt.Fprintf(w, "%s  Cap: %d\n", indent, r.Cap)
	fmt.Fprintf(w, "%s  Name: %s\n", indent, r.Name)
}

// ParseOrPanic parses a regular expression into an AST.
// Panics on error.
func parseOrPanic(simplify bool, pattern string) *syntax.Regexp {
	regexp, err := syntax.Parse(pattern, 0)
	if err != nil {
		panic(err)
	}
	if simplify {
		regexp = regexp.Simplify()
	}
	return regexp
}

// runesToString converts a slice of runes to the string they represent.
func runesToString(runes ...rune) string {
	defer func() {
		if err := recover(); err != nil {
			panic(fmt.Errorf("RunesToString panicked"))
		}
	}()
	var buffer bytes.Buffer
	for _, r := range runes {
		buffer.WriteRune(r)
	}
	return buffer.String()
}

// RunesToDecimalString converts a slice of runes to their comma-separated decimal values.
func runesToDecimalString(runes []rune) string {
	var buffer bytes.Buffer
	for _, r := range runes {
		buffer.WriteString(fmt.Sprintf("%d, ", r))
	}
	return buffer.String()
}

// opToString gets the string name of a regular expression operation.
func opToString(op syntax.Op) string {
	switch op {
	case syntax.OpNoMatch:
		return "OpNoMatch"
	case syntax.OpEmptyMatch:
		return "OpEmptyMatch"
	case syntax.OpLiteral:
		return "OpLiteral"
	case syntax.OpCharClass:
		return "OpCharClass"
	case syntax.OpAnyCharNotNL:
		return "OpAnyCharNotNL"
	case syntax.OpAnyChar:
		return "OpAnyChar"
	case syntax.OpBeginLine:
		return "OpBeginLine"
	case syntax.OpEndLine:
		return "OpEndLine"
	case syntax.OpBeginText:
		return "OpBeginText"
	case syntax.OpEndText:
		return "OpEndText"
	case syntax.OpWordBoundary:
		return "OpWordBoundary"
	case syntax.OpNoWordBoundary:
		return "OpNoWordBoundary"
	case syntax.OpCapture:
		return "OpCapture"
	case syntax.OpStar:
		return "OpStar"
	case syntax.OpPlus:
		return "OpPlus"
	case syntax.OpQuest:
		return "OpQuest"
	case syntax.OpRepeat:
		return "OpRepeat"
	case syntax.OpConcat:
		return "OpConcat"
	case syntax.OpAlternate:
		return "OpAlternate"
	}

	panic(fmt.Sprintf("invalid op: %d", op))
}
