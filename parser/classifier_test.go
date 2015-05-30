package parser

import (
	goparser "go/parser"
	"log"
	"path/filepath"
	"sort"
	"testing"

	"golang.org/x/tools/go/loader"

	"github.com/stretchr/testify/assert"
)

func TestAnnotationMatcher(t *testing.T) {
	variations := []string{
		"// +swagger",
		" +swagger",
		"+swagger",
		" * +swagger",
	}
	known := []string{
		"meta",
		"route",
		"model",
		"parameters",
		"strfmt",
	}

	for _, variation := range variations {
		for _, tpe := range known {
			assert.True(t, rxSwaggerAnnotation.MatchString(variation+":"+tpe))
		}
	}
}

func classifierProgram() *loader.Program {
	var ldr loader.Config
	ldr.ParserMode = goparser.ParseComments
	ldr.Import("../fixtures/goparsing/classification")
	ldr.Import("../fixtures/goparsing/classification/models")
	ldr.Import("../fixtures/goparsing/classification/operations")
	prog, err := ldr.Load()
	if err != nil {
		log.Fatal(err)
	}
	return prog
}

func TestClassifier(t *testing.T) {

	prog := classificationProg
	classifier := &programClassifier{}
	classified, err := classifier.Classify(prog)
	assert.NoError(t, err)

	// ensure all the dependencies are there
	assert.Len(t, classified.Meta, 1)
	assert.Len(t, classified.Operations, 1)
	assert.Len(t, classified.Models, 3)

	var fNames []string
	for _, file := range classified.Models {
		fNames = append(
			fNames,
			filepath.Base(prog.Fset.File(file.Pos()).Name()))
	}

	sort.Sort(sort.StringSlice(fNames))
	assert.EqualValues(t, []string{"order.go", "pet.go", "user.go"}, fNames)
}

func TestClassifierInclude(t *testing.T) {

	prog := classificationProg
	classifier := &programClassifier{
		Includes: packageFilters{
			Models: []packageFilter{
				packageFilter{"github.com/casualjim/go-swagger/fixtures/goparsing/classification"},
				packageFilter{"github.com/casualjim/go-swagger/fixtures/goparsing/classification/transitive/mods"},
				packageFilter{"github.com/casualjim/go-swagger/fixtures/goparsing/classification/operations"},
			},
		},
	}
	classified, err := classifier.Classify(prog)
	assert.NoError(t, err)

	// ensure all the dependencies are there
	assert.Len(t, classified.Meta, 1)
	assert.Len(t, classified.Operations, 1)
	assert.Len(t, classified.Models, 1)

	var fNames []string
	for _, file := range classified.Models {
		fNames = append(
			fNames,
			filepath.Base(prog.Fset.File(file.Pos()).Name()))
	}

	sort.Sort(sort.StringSlice(fNames))
	assert.EqualValues(t, []string{"pet.go"}, fNames)
}

func TestClassifierExclude(t *testing.T) {

	prog := classificationProg
	classifier := &programClassifier{
		Excludes: packageFilters{
			Models: []packageFilter{
				packageFilter{"github.com/casualjim/go-swagger/fixtures/goparsing/classification/transitive/mods"},
			},
		},
	}
	classified, err := classifier.Classify(prog)
	assert.NoError(t, err)

	// ensure all the dependencies are there
	assert.Len(t, classified.Meta, 1)
	assert.Len(t, classified.Operations, 1)
	assert.Len(t, classified.Models, 2)

	var fNames []string
	for _, file := range classified.Models {
		fNames = append(
			fNames,
			filepath.Base(prog.Fset.File(file.Pos()).Name()))
	}

	sort.Sort(sort.StringSlice(fNames))
	assert.EqualValues(t, []string{"order.go", "user.go"}, fNames)
}
