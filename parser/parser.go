package parser

import (
	"fmt"
	"go/ast"
	goparser "go/parser"
	"log"
	"strings"

	"golang.org/x/tools/go/loader"

	"github.com/casualjim/go-swagger/spec"
	"github.com/casualjim/go-swagger/util"
)

// Many thanks go to https://github.com/yvasiyarov/swagger
// this is loosely based on that implementation but for swagger 2.0

type setter func(interface{}, []string) error

// apiParser the global context for parsing a go application
// into a swagger specification
type apiParser struct {
	loader     *loader.Config
	prog       *loader.Program
	classifier *programClassifier

	// MainPackage the path to find the main class in
	MainPackage string
}

// newAPIParser creates a new api parser
func newAPIParser(bp string, includes, excludes packageFilters) (*apiParser, error) {
	var ldr loader.Config
	ldr.ParserMode = goparser.ParseComments
	ldr.Import(bp)
	prog, err := ldr.Load()
	if err != nil {
		return nil, err
	}
	return &apiParser{MainPackage: bp, prog: prog, loader: &ldr, classifier: &programClassifier{
		Includes: includes,
		Excludes: excludes,
	}}, nil
}

// Parse produces a swagger object for an application
func (a *apiParser) Parse() (*spec.Swagger, error) {
	// classification still includes files that are completely commented out
	cp, err := a.classifier.Classify(a.prog)
	if err != nil {
		return nil, err
	}

	// build definitions dictionary
	var definitions = make(map[string]spec.Schema)
	for _, modFile := range cp.Models {
		if err := a.parseSchema(modFile, definitions, nil); err != nil {
			return nil, err
		}
	}

	// build paths dictionary
	var paths spec.Paths
	for _, routeFile := range cp.Operations {
		if err := a.parseRoutes(routeFile, &paths); err != nil {
			return nil, err
		}
	}

	// build swagger object
	result := new(spec.Swagger)
	for _, metaFile := range cp.Meta {
		if err := a.parseMeta(metaFile, result); err != nil {
			return nil, err
		}
	}
	result.Paths = &paths
	result.Definitions = definitions
	return result, nil
}

func (a *apiParser) parseSchema(file *ast.File, definitions map[string]spec.Schema, parent *spec.Schema) error {
	return nil
}

func (a *apiParser) parseRoutes(file *ast.File, paths *spec.Paths) error {
	return nil
}

func (a *apiParser) parseMeta(file *ast.File, swspec *spec.Swagger) error {
	newMetaParser().Parse(file, swspec)
	return nil
}

// MustExpandPackagePath gets the real package path on disk
func (a *apiParser) MustExpandPackagePath(packagePath string) string {
	pkgRealpath := util.FindInGoSearchPath(packagePath)
	if pkgRealpath == "" {
		log.Fatalf("Can't find package %s \n", packagePath)
	}

	return pkgRealpath
}

func newDocCommentParser(otherTags []string, taggers ...*sectionTagger) *docCommentParser {
	return &docCommentParser{taggers: taggers, otherTags: otherTags}
}

type docCommentParser struct {
	taggers   []*sectionTagger
	otherTags []string
	header    []string
	body      []string
}

func (ai *docCommentParser) Parse(gofile *ast.File, target interface{}) error {
	var selectedTagger *sectionTagger
	var otherTags []string
	taggers := ai.taggers
	for _, c := range gofile.Doc.List {
		text := c.Text
		lines := strings.Split(text, "\n")

	LINES:
		for _, line := range lines {
			fmt.Printf("processing: %q\n", line)
			// this is an aggregating tagger
			if selectedTagger != nil {
				switch res := selectedTagger.Tag(line, otherTags).(type) {
				case multiLineSectionPart:
					fmt.Println("this is a multi line section part for", selectedTagger.Name)
					continue
				case multiLineSectionTerminator:
					fmt.Println("this is a multi line section terminator for", selectedTagger.Name)
					if err := selectedTagger.set(target, res.taggedSection.Lines); err != nil {
						return err
					}
					selectedTagger = nil
					continue
				case newTagSectionTerminator:
					fmt.Println("this is a multi line section tag terminator for", selectedTagger.Name)
					if err := selectedTagger.set(target, res.taggedSection.Lines); err != nil {
						return err
					}
					selectedTagger = nil
				}
			}
			if len(taggers) == 0 {
				break
			}
			selectedTagger = nil
			for i, tagger := range taggers {
				switch res := tagger.Tag(line, nil).(type) {
				case singleLineSection:
					fmt.Println("this is a single line section for", tagger.Name)
					if err := tagger.set(target, res.taggedSection.Lines); err != nil {
						return err
					}
					// once it has matched we don't care for probing for it again
					taggers = append(taggers[:i], taggers[i+1:]...)
					continue LINES

				case multiLineSectionPart:
					fmt.Println("this is a multi line section for", tagger.Name)
					selectedTagger = tagger
					otherTags = ai.otherTags
					for _, t := range ai.taggers {
						if t.Name != tagger.Name {
							otherTags = append(otherTags, t.Name)
						}
					}
					// once it has matched we don't care for probing for it again
					taggers = append(taggers[:i], taggers[i+1:]...)
					continue LINES

				case unmatchedSection:
				}
			}
		}
	}

	if selectedTagger != nil {
		if err := selectedTagger.set(target, selectedTagger.Lines); err != nil {
			return err
		}
	}

	return nil
}
