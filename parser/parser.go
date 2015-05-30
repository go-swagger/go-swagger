package parser

import (
	"fmt"
	"go/ast"
	goparser "go/parser"
	"log"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/tools/go/loader"

	"github.com/casualjim/go-swagger/spec"
	"github.com/casualjim/go-swagger/util"
)

var (
	rxSwaggerAnnotation  = regexp.MustCompile("[^+]*\\+\\p{Zs}*swagger:([\\p{L}\\p{N}-_]+)")
	rxStrFmt             = regexp.MustCompile("\\+swagger:strfmt\\p{Zs}*(\\p{L}[\\p{L}\\p{N}-]+)$")
	rxModelOverride      = regexp.MustCompile("\\+swagger:model\\p{Zs}*(\\p{L}[\\p{L}\\p{N}-]+)?$")
	rxParametersOverride = regexp.MustCompile("\\+swagger:parameters\\p{Zs}*(\\p{L}[\\p{L}\\p{N}-\\p{Zs}]+)$")
	rxRoute              = regexp.MustCompile("\\+swagger:route\\p{Zs}*(\\p{L}+)\\p{Zs}*((?:/[\\p{L}\\p{N}-_{}]*)+/?)\\p{Zs}+((?:\\p{L}[\\p{L}\\p{N}-]+)+)\\p{Zs}*(\\p{L}[\\p{L}\\p{N}-\\p{Zs}]+)$")

	rxMaximumFmt    = "%s[Mm]ax(?:imum)?\\p{Zs}*:\\p{Zs}*([\\<=])?\\p{Zs}*([\\+-]?(?:\\p{N}+\\.)?\\p{N}+)$"
	rxMinimumFmt    = "%s[Mm]in(?:imum)?\\p{Zs}*:\\p{Zs}*([\\>=])?\\p{Zs}*([\\+-]?(?:\\p{N}+\\.)?\\p{N}+)$"
	rxMultipleOfFmt = "%s[Mm]ultiple\\p{Zs}*[Oo]f\\p{Zs}*:\\p{Zs}*([\\+-]?(?:\\p{N}+\\.)?\\p{N}+)$"

	rxMaxLengthFmt        = "%s[Mm]ax(?:imum)?(?:\\p{Zs}*-?[Ll]en(?:gth)?)\\p{Zs}*:\\p{Zs}*(\\p{N}+)$"
	rxMinLengthFmt        = "%s[Mm]in(?:imum)?(?:\\p{Zs}*-?[Ll]en(?:gth)?)\\p{Zs}*:\\p{Zs}*(\\p{N}+)$"
	rxPatternFmt          = "%s[Pp]attern\\p{Zs}*:\\p{Zs}*(.*)$"
	rxCollectionFormatFmt = "%s[Cc]ollection(?:\\p{Zs}*-?[Ff]ormat)\\p{Zs}*:\\p{Zs}*(.*)$"

	rxMaxItemsFmt = "%s[Mm]ax(?:imum)?(?:\\p{Zs}*|-)?[Ii]tems\\p{Zs}*:\\p{Zs}*(\\p{N}+)$"
	rxMinItemsFmt = "%s[Mm]in(?:imum)?(?:\\p{Zs}*|-)?[Ii]tems\\p{Zs}*:\\p{Zs}*(\\p{N}+)$"
	rxUniqueFmt   = "%s[Uu]nique\\p{Zs}*:\\p{Zs}*(true|false)$"

	rxIn       = regexp.MustCompile("(?:[Ii]n|[Ss]ource)\\p{Zs}*:\\p{Zs}*(query|path|header|body)$")
	rxRequired = regexp.MustCompile("[Rr]equired\\p{Zs}*:\\p{Zs}*(true|false)$")
	rxReadOnly = regexp.MustCompile("[Rr]ead(?:\\p{Zs}*|-)?[Oo]nly\\p{Zs}*:\\p{Zs}*(true|false)$")

	rxItemsPrefix = "(?:[Ii]tems[\\.\\p{Zs}]?)+"
)

// Many thanks go to https://github.com/yvasiyarov/swagger
// this is loosely based on that implementation but for swagger 2.0

type setter func(interface{}, []string) error

func rxf(rxp, ar string) *regexp.Regexp {
	return regexp.MustCompile(fmt.Sprintf(rxp, ar))
}

// apiParser the global context for parsing a go application
// into a swagger specification
type apiParser struct {
	loader     *loader.Config
	prog       *loader.Program
	classifier *programClassifier
	discovered []schemaDecl

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
		if err := a.parseSchema(modFile, definitions); err != nil {
			return nil, err
		}
	}
	// loop over discovered until all the items are in definitions
	keepGoing := len(a.discovered) > 0
	for keepGoing {
		var queue []schemaDecl
		for _, d := range a.discovered {
			if _, ok := definitions[d.Name]; !ok {
				queue = append(queue, d)
			}
		}
		a.discovered = nil
		for _, sd := range queue {
			if err := a.parseSchema(sd.File, definitions); err != nil {
				return nil, err
			}
		}
		keepGoing = len(a.discovered) > 0
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

func (a *apiParser) parseSchema(file *ast.File, definitions map[string]spec.Schema) error {
	sp := newSchemaParser(a.prog)
	if err := sp.Parse(file, definitions); err != nil {
		return err
	}
	a.discovered = append(a.discovered, sp.postDecls...)
	return nil
}

func (a *apiParser) parseRoutes(file *ast.File, paths *spec.Paths) error {
	return nil
}

func (a *apiParser) parseMeta(file *ast.File, swspec *spec.Swagger) error {
	return newMetaParser().Parse(file, swspec)
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

func parseDocComments(doc *ast.CommentGroup, target interface{}, tgrs []*sectionTagger, ot []string) error {
	if doc == nil {
		return nil
	}
	var selectedTagger *sectionTagger
	var otherTags []string
	taggers := tgrs
	for _, c := range doc.List {
		text := c.Text
		lines := strings.Split(text, "\n")

	LINES:
		for _, line := range lines {
			// this is an aggregating tagger
			if selectedTagger != nil {
				switch res := selectedTagger.Tag(line, otherTags).(type) {
				case multiLineSectionPart:
					continue
				case multiLineSectionTerminator:
					if err := selectedTagger.set(target, res.taggedSection.Lines); err != nil {
						return err
					}
					selectedTagger = nil
					continue
				case newTagSectionTerminator:
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
					if err := tagger.set(target, res.taggedSection.Lines); err != nil {
						return err
					}
					// once it has matched we don't care for probing for it again
					taggers = append(taggers[:i], taggers[i+1:]...)
					continue LINES

				case multiLineSectionPart:
					selectedTagger = tagger
					otherTags = ot
					for _, t := range tgrs {
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

type docCommentParser struct {
	taggers   []*sectionTagger
	otherTags []string
}

func (ai *docCommentParser) Parse(gofile *ast.File, target interface{}) error {
	return parseDocComments(gofile.Doc, target, ai.taggers, ai.otherTags)
}

type swaggerTypable interface {
	Typed(string, string)
	SetRef(spec.Ref)
}

type selectorParser struct {
	program     *loader.Program
	AddPostDecl func(schemaDecl)
}

func (sp *selectorParser) TypeForSelector(gofile *ast.File, expr *ast.SelectorExpr, prop swaggerTypable) error {
	if pth, ok := expr.X.(*ast.Ident); ok {
		// lookup import
		var selPath string
		for _, imp := range gofile.Imports {
			pv, err := strconv.Unquote(imp.Path.Value)
			if err != nil {
				pv = imp.Path.Value
			}
			if imp.Name != nil {
				if imp.Name.Name == pth.Name {
					selPath = pv
					break
				}
			} else {
				parts := strings.Split(pv, "/")
				if len(parts) > 0 && parts[len(parts)-1] == pth.Name {
					selPath = pv
					break
				}
			}
		}
		// find actual struct
		if selPath == "" {
			return fmt.Errorf("no import found for %s", pth.Name)
		}

		pkg := sp.program.Package(selPath)
		if pkg == nil {
			return fmt.Errorf("no package found for %s", selPath)
		}

		// find the file this selector points to
		for _, file := range pkg.Files {
			for _, decl := range file.Decls {
				if gd, ok := decl.(*ast.GenDecl); ok {
					for _, gs := range gd.Specs {
						if ts, ok := gs.(*ast.TypeSpec); ok {
							if ts.Name != nil && ts.Name.Name == expr.Sel.Name {
								// look at doc comments for +swagger:strfmt [name]
								// when found this is the format name, create a schema with that name
								if gd.Doc != nil {
									for _, cmt := range gd.Doc.List {
										for _, ln := range strings.Split(cmt.Text, "\n") {
											matches := rxStrFmt.FindStringSubmatch(ln)
											if len(matches) > 1 && len(matches[1]) > 0 {
												prop.Typed("string", matches[1])
												return nil
											}
										}
									}
								}
								// ok so not a string format, perhaps a model?
								if _, ok := ts.Type.(*ast.StructType); ok {
									ref, err := spec.NewRef("#/definitions/" + ts.Name.Name)
									if err != nil {
										return err
									}
									prop.SetRef(ref)
									sd := schemaDecl{file, gd, ts, "", ""}
									sd.inferNames()
									sp.AddPostDecl(sd)
									return nil
								}
							}
						}
					}
				}
			}
		}

		return fmt.Errorf("schema parser: no string format for %s.%s", pth.Name, expr.Sel.Name)
	}
	return fmt.Errorf("schema parser: no string format for %v", expr.Sel.Name)
}

func swaggerSchemaForType(typeName string, prop swaggerTypable) error {
	switch typeName {
	case "bool":
		prop.Typed("boolean", "")
	case "rune", "string":
		prop.Typed("string", "")
	case "int8":
		prop.Typed("number", "int8")
	case "int16":
		prop.Typed("number", "int16")
	case "int32":
		prop.Typed("number", "int32")
	case "int", "int64":
		prop.Typed("number", "int64")
	case "uint8":
		prop.Typed("number", "uint8")
	case "uint16":
		prop.Typed("number", "uint16")
	case "uint32":
		prop.Typed("number", "uint32")
	case "uint", "uint64":
		prop.Typed("number", "uint64")
	case "float32":
		prop.Typed("number", "float")
	case "float64":
		prop.Typed("number", "double")
	}
	return nil
}
