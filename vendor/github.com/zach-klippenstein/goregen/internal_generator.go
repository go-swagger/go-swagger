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
	"math"
	"regexp/syntax"
)

// generatorFactory is a function that creates a random string generator from a regular expression AST.
type generatorFactory func(regexp *syntax.Regexp, args *GeneratorArgs) (*internalGenerator, error)

// Must be initialized in init() to avoid "initialization loop" compile error.
var generatorFactories map[syntax.Op]generatorFactory

const noBound = -1

func init() {
	generatorFactories = map[syntax.Op]generatorFactory{
		syntax.OpEmptyMatch:     opEmptyMatch,
		syntax.OpLiteral:        opLiteral,
		syntax.OpAnyCharNotNL:   opAnyCharNotNl,
		syntax.OpAnyChar:        opAnyChar,
		syntax.OpQuest:          opQuest,
		syntax.OpStar:           opStar,
		syntax.OpPlus:           opPlus,
		syntax.OpRepeat:         opRepeat,
		syntax.OpCharClass:      opCharClass,
		syntax.OpConcat:         opConcat,
		syntax.OpAlternate:      opAlternate,
		syntax.OpCapture:        opCapture,
		syntax.OpBeginLine:      noop,
		syntax.OpEndLine:        noop,
		syntax.OpBeginText:      noop,
		syntax.OpEndText:        noop,
		syntax.OpWordBoundary:   noop,
		syntax.OpNoWordBoundary: noop,
	}
}

type internalGenerator struct {
	Name         string
	GenerateFunc func() string
}

func (gen *internalGenerator) Generate() string {
	return gen.GenerateFunc()
}

func (gen *internalGenerator) String() string {
	return gen.Name
}

// Create a new generator for each expression in regexps.
func newGenerators(regexps []*syntax.Regexp, args *GeneratorArgs) ([]*internalGenerator, error) {
	generators := make([]*internalGenerator, len(regexps), len(regexps))
	var err error

	// create a generator for each alternate pattern
	for i, subR := range regexps {
		generators[i], err = newGenerator(subR, args)
		if err != nil {
			return nil, err
		}
	}

	return generators, nil
}

// Create a new generator for r.
func newGenerator(regexp *syntax.Regexp, args *GeneratorArgs) (generator *internalGenerator, err error) {
	simplified := regexp.Simplify()

	factory, ok := generatorFactories[simplified.Op]
	if ok {
		return factory(simplified, args)
	}

	return nil, fmt.Errorf("invalid generator pattern: /%s/ as /%s/\n%s",
		regexp, simplified, inspectRegexpToString(simplified))
}

// Generator that does nothing.
func noop(regexp *syntax.Regexp, args *GeneratorArgs) (*internalGenerator, error) {
	return &internalGenerator{regexp.String(), func() string {
		return ""
	}}, nil
}

func opEmptyMatch(regexp *syntax.Regexp, args *GeneratorArgs) (*internalGenerator, error) {
	enforceOp(regexp, syntax.OpEmptyMatch)
	return &internalGenerator{regexp.String(), func() string {
		return ""
	}}, nil
}

func opLiteral(regexp *syntax.Regexp, args *GeneratorArgs) (*internalGenerator, error) {
	enforceOp(regexp, syntax.OpLiteral)
	return &internalGenerator{regexp.String(), func() string {
		return runesToString(regexp.Rune...)
	}}, nil
}

func opAnyChar(regexp *syntax.Regexp, args *GeneratorArgs) (*internalGenerator, error) {
	enforceOp(regexp, syntax.OpAnyChar)
	return &internalGenerator{regexp.String(), func() string {
		return runesToString(rune(args.rng.Int31()))
	}}, nil
}

func opAnyCharNotNl(regexp *syntax.Regexp, args *GeneratorArgs) (*internalGenerator, error) {
	enforceOp(regexp, syntax.OpAnyCharNotNL)
	charClass := newCharClass(1, rune(math.MaxInt32))
	return createCharClassGenerator(regexp.String(), charClass, args)
}

func opQuest(regexp *syntax.Regexp, args *GeneratorArgs) (*internalGenerator, error) {
	enforceOp(regexp, syntax.OpQuest)
	return createRepeatingGenerator(regexp, args, 0, 1)
}

func opStar(regexp *syntax.Regexp, args *GeneratorArgs) (*internalGenerator, error) {
	enforceOp(regexp, syntax.OpStar)
	return createRepeatingGenerator(regexp, args, noBound, noBound)
}

func opPlus(regexp *syntax.Regexp, args *GeneratorArgs) (*internalGenerator, error) {
	enforceOp(regexp, syntax.OpPlus)
	return createRepeatingGenerator(regexp, args, 1, noBound)
}

func opRepeat(regexp *syntax.Regexp, args *GeneratorArgs) (*internalGenerator, error) {
	enforceOp(regexp, syntax.OpRepeat)
	return createRepeatingGenerator(regexp, args, regexp.Min, regexp.Max)
}

// Handles syntax.ClassNL because the parser uses that flag to generate character
// classes that respect it.
func opCharClass(regexp *syntax.Regexp, args *GeneratorArgs) (*internalGenerator, error) {
	enforceOp(regexp, syntax.OpCharClass)
	charClass := parseCharClass(regexp.Rune)
	return createCharClassGenerator(regexp.String(), charClass, args)
}

func opConcat(regexp *syntax.Regexp, genArgs *GeneratorArgs) (*internalGenerator, error) {
	enforceOp(regexp, syntax.OpConcat)

	generators, err := newGenerators(regexp.Sub, genArgs)
	if err != nil {
		return nil, generatorError(err, "error creating generators for concat pattern /%s/", regexp)
	}

	return &internalGenerator{regexp.String(), func() string {
		var result bytes.Buffer
		for _, generator := range generators {
			result.WriteString(generator.Generate())
		}
		return result.String()
	}}, nil
}

func opAlternate(regexp *syntax.Regexp, genArgs *GeneratorArgs) (*internalGenerator, error) {
	enforceOp(regexp, syntax.OpAlternate)

	generators, err := newGenerators(regexp.Sub, genArgs)
	if err != nil {
		return nil, generatorError(err, "error creating generators for alternate pattern /%s/", regexp)
	}

	numGens := len(generators)

	return &internalGenerator{regexp.String(), func() string {
		i := genArgs.rng.Intn(numGens)
		generator := generators[i]
		return generator.Generate()
	}}, nil
}

func opCapture(regexp *syntax.Regexp, args *GeneratorArgs) (*internalGenerator, error) {
	enforceOp(regexp, syntax.OpCapture)

	if err := enforceSingleSub(regexp); err != nil {
		return nil, err
	}

	groupRegexp := regexp.Sub[0]
	generator, err := newGenerator(groupRegexp, args)
	if err != nil {
		return nil, err
	}

	// Group indices are 0-based, but index 0 is the whole expression.
	index := regexp.Cap - 1

	return &internalGenerator{regexp.String(), func() string {
		return args.CaptureGroupHandler(index, regexp.Name, groupRegexp, generator, args)
	}}, nil
}

func defaultCaptureGroupHandler(index int, name string, group *syntax.Regexp, generator Generator, args *GeneratorArgs) string {
	return generator.Generate()
}

// Panic if r.Op != op.
func enforceOp(r *syntax.Regexp, op syntax.Op) {
	if r.Op != op {
		panic(fmt.Sprintf("invalid Op: expected %s, was %s", opToString(op), opToString(r.Op)))
	}
}

// Return an error if r has 0 or more than 1 sub-expression.
func enforceSingleSub(regexp *syntax.Regexp) error {
	if len(regexp.Sub) != 1 {
		return generatorError(nil,
			"%s expected 1 sub-expression, but got %d: %s", opToString(regexp.Op), len(regexp.Sub), regexp)
	}
	return nil
}

func createCharClassGenerator(name string, charClass *tCharClass, args *GeneratorArgs) (*internalGenerator, error) {
	return &internalGenerator{name, func() string {
		i := args.rng.Int31n(charClass.TotalSize)
		r := charClass.GetRuneAt(i)
		return runesToString(r)
	}}, nil
}

// Returns a generator that will run the generator for r's sub-expression [min, max] times.
func createRepeatingGenerator(regexp *syntax.Regexp, genArgs *GeneratorArgs, min, max int) (*internalGenerator, error) {
	if err := enforceSingleSub(regexp); err != nil {
		return nil, err
	}

	generator, err := newGenerator(regexp.Sub[0], genArgs)
	if err != nil {
		return nil, generatorError(err, "failed to create generator for subexpression: /%s/", regexp)
	}

	if min == noBound {
		min = int(genArgs.MinUnboundedRepeatCount)
	}
	if max == noBound {
		max = int(genArgs.MaxUnboundedRepeatCount)
	}

	return &internalGenerator{regexp.String(), func() string {
		n := min + genArgs.rng.Intn(max-min+1)

		var result bytes.Buffer
		for i := 0; i < n; i++ {
			result.WriteString(generator.Generate())
		}
		return result.String()
	}}, nil
}
