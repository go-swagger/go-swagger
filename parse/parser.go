package parse

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

const (
	rxMethod = "(\\p{L}+)"
	rxPath   = "((?:/[\\p{L}\\p{N}\\p{Pd}\\p{Pc}{}]*)+/?)"
	rxOpTags = "(\\p{L}[\\p{L}\\p{N}\\p{Pd}\\p{Pc}\\p{Zs}]+)"
	rxOpID   = "((?:\\p{L}[\\p{L}\\p{N}\\p{Pd}\\p{Pc}]+)+)"

	rxMaximumFmt    = "%s[Mm]ax(?:imum)?\\p{Zs}*:\\p{Zs}*([\\<=])?\\p{Zs}*([\\+-]?(?:\\p{N}+\\.)?\\p{N}+)$"
	rxMinimumFmt    = "%s[Mm]in(?:imum)?\\p{Zs}*:\\p{Zs}*([\\>=])?\\p{Zs}*([\\+-]?(?:\\p{N}+\\.)?\\p{N}+)$"
	rxMultipleOfFmt = "%s[Mm]ultiple\\p{Zs}*[Oo]f\\p{Zs}*:\\p{Zs}*([\\+-]?(?:\\p{N}+\\.)?\\p{N}+)$"

	rxMaxLengthFmt        = "%s[Mm]ax(?:imum)?(?:\\p{Zs}*[\\p{Pd}\\p{Pc}]?[Ll]en(?:gth)?)\\p{Zs}*:\\p{Zs}*(\\p{N}+)$"
	rxMinLengthFmt        = "%s[Mm]in(?:imum)?(?:\\p{Zs}*[\\p{Pd}\\p{Pc}]?[Ll]en(?:gth)?)\\p{Zs}*:\\p{Zs}*(\\p{N}+)$"
	rxPatternFmt          = "%s[Pp]attern\\p{Zs}*:\\p{Zs}*(.*)$"
	rxCollectionFormatFmt = "%s[Cc]ollection(?:\\p{Zs}*[\\p{Pd}\\p{Pc}]?[Ff]ormat)\\p{Zs}*:\\p{Zs}*(.*)$"

	rxMaxItemsFmt = "%s[Mm]ax(?:imum)?(?:\\p{Zs}*|[\\p{Pd}\\p{Pc}]|\\.)?[Ii]tems\\p{Zs}*:\\p{Zs}*(\\p{N}+)$"
	rxMinItemsFmt = "%s[Mm]in(?:imum)?(?:\\p{Zs}*|[\\p{Pd}\\p{Pc}]|\\.)?[Ii]tems\\p{Zs}*:\\p{Zs}*(\\p{N}+)$"
	rxUniqueFmt   = "%s[Uu]nique\\p{Zs}*:\\p{Zs}*(true|false)$"

	rxItemsPrefix = "(?:[Ii]tems[\\.\\p{Zs}]?)+"
)

var (
	rxSwaggerAnnotation  = regexp.MustCompile("[^+]*\\+\\p{Zs}*swagger:([\\p{L}\\p{N}\\p{Pd}\\p{Pc}]+)")
	rxMeta               = regexp.MustCompile("\\+swagger:meta")
	rxStrFmt             = regexp.MustCompile("\\+swagger:strfmt\\p{Zs}*(\\p{L}[\\p{L}\\p{N}\\p{Pd}\\p{Pc}]+)$")
	rxModelOverride      = regexp.MustCompile("\\+swagger:model\\p{Zs}*(\\p{L}[\\p{L}\\p{N}\\p{Pd}\\p{Pc}]+)?$")
	rxResponseOverride   = regexp.MustCompile("\\+swagger:response\\p{Zs}*(\\p{L}[\\p{L}\\p{N}\\p{Pd}\\p{Pc}]+)?$")
	rxParametersOverride = regexp.MustCompile("\\+swagger:parameters\\p{Zs}*(\\p{L}[\\p{L}\\p{N}\\p{Pd}\\p{Pc}\\p{Zs}]+)$")
	rxRoute              = regexp.MustCompile(
		"\\+swagger:route\\p{Zs}*" +
			rxMethod +
			"\\p{Zs}*" +
			rxPath +
			"\\p{Zs}+" +
			rxOpTags +
			"\\p{Zs}+" +
			rxOpID + "$")

	rxIn                 = regexp.MustCompile("(?:[Ii]n|[Ss]ource)\\p{Zs}*:\\p{Zs}*(query|path|header|body)$")
	rxRequired           = regexp.MustCompile("[Rr]equired\\p{Zs}*:\\p{Zs}*(true|false)$")
	rxReadOnly           = regexp.MustCompile("[Rr]ead(?:\\p{Zs}*|[\\p{Pd}\\p{Pc}])?[Oo]nly\\p{Zs}*:\\p{Zs}*(true|false)$")
	rxSpace              = regexp.MustCompile("\\p{Zs}+")
	rxNotAlNumSpaceComma = regexp.MustCompile("[^\\p{L}\\p{N}\\p{Zs},]")
	rxPunctuationEnd     = regexp.MustCompile("\\p{Po}$")
	rxStripComments      = regexp.MustCompile("^[^\\w\\+]*")
	rxStripTitleComments = regexp.MustCompile("^[^\\p{L}]*(:?P|p)ackage\\p{Zs}+[^\\p{Zs}]+\\p{Zs}*")

	rxConsumes  = regexp.MustCompile("[Cc]onsumes\\p{Zs}*:")
	rxProduces  = regexp.MustCompile("[Pp]roduces\\p{Zs}*:")
	rxSecurity  = regexp.MustCompile("[Ss]ecurity\\p{Zs}*:")
	rxResponses = regexp.MustCompile("[Rr]esponses\\p{Zs}*:")
	rxSchemes   = regexp.MustCompile("[Ss]chemes\\p{Zs}*:\\p{Zs}*((?:(?:https?|HTTPS?|wss?|WSS?)[\\p{Zs},]*)+)$")
	rxVersion   = regexp.MustCompile("[Vv]ersion\\p{Zs}*:\\p{Zs}*(.+)$")
	rxHost      = regexp.MustCompile("[Hh]ost\\p{Zs}*:\\p{Zs}*(.+)$")
	rxBasePath  = regexp.MustCompile("[Bb]ase\\p{Zs}*-*[Pp]ath\\p{Zs}*:\\p{Zs}*" + rxPath + "$")
	rxLicense   = regexp.MustCompile("[Ll]icense\\p{Zs}*:\\p{Zs}*(.+)$")
	rxContact   = regexp.MustCompile("[Cc]ontact\\p{Zs}*-?(?:[Ii]info\\p{Zs}*)?:\\p{Zs}*(.+)$")
	rxTOS       = regexp.MustCompile("[Tt](:?erms)?\\p{Zs}*-?[Oo]f?\\p{Zs}*-?[Ss](?:ervice)?\\p{Zs}*:")
)

// Many thanks go to https://github.com/yvasiyarov/swagger
// this is loosely based on that implementation but for swagger 2.0

type setter func(interface{}, []string) error

func joinDropLast(lines []string) string {
	l := len(lines)
	lns := lines
	if l > 0 && len(strings.TrimSpace(lines[l-1])) == 0 {
		lns = lines[:l-1]
	}
	return strings.Join(lns, "\n")
}

func removeEmptyLines(lines []string) (notEmpty []string) {
	for _, l := range lines {
		if len(strings.TrimSpace(l)) > 0 {
			notEmpty = append(notEmpty, l)
		}
	}
	return
}

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

	// find the operations
	// build the responses
	// build the paramters
	// build schemas for discovered models
	// add meta data

	// build definitions dictionary
	var definitions = make(map[string]spec.Schema)
	for _, modFile := range cp.Models {
		if err := a.parseSchema(modFile, definitions); err != nil {
			return nil, err
		}
	}

	// build parameters dictionary
	var parameters = make(map[string]spec.Operation)
	for _, paramsFile := range cp.Parameters {
		if err := a.parseParameters(paramsFile, parameters); err != nil {
			return nil, err
		}
	}

	// build responses dictionary
	var responses = make(map[string]spec.Response)
	for _, responseFile := range cp.Responses {
		if err := a.parseResponses(responseFile, responses); err != nil {
			return nil, err
		}
	}

	if err := a.processDiscovered(definitions); err != nil {
		return nil, err
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
	// assemble swagger object, only including things that are in actual use
	result.Paths = &paths
	result.Definitions = definitions
	return result, nil
}

func (a *apiParser) processDiscovered(definitions map[string]spec.Schema) error {
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
				return err
			}
		}
		keepGoing = len(a.discovered) > 0
	}

	return nil
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
	rp := newRoutesParser(a.prog)
	if err := rp.Parse(file, paths); err != nil {
		return err
	}
	return nil
}

func (a *apiParser) parseParameters(file *ast.File, operations map[string]spec.Operation) error {
	rp := newParameterParser(a.prog)
	if err := rp.Parse(file, operations); err != nil {
		return err
	}
	a.discovered = append(a.discovered, rp.postDecls...)
	return nil
}

func (a *apiParser) parseResponses(file *ast.File, responses map[string]spec.Response) error {
	rp := newResponseParser(a.prog)
	if err := rp.Parse(file, responses); err != nil {
		return err
	}
	a.discovered = append(a.discovered, rp.postDecls...)
	return nil
}

func (a *apiParser) parseMeta(file *ast.File, swspec *spec.Swagger) error {
	return newMetaParser(swspec).Parse(file.Doc)
}

// MustExpandPackagePath gets the real package path on disk
func (a *apiParser) MustExpandPackagePath(packagePath string) string {
	pkgRealpath := util.FindInGoSearchPath(packagePath)
	if pkgRealpath == "" {
		log.Fatalf("Can't find package %s \n", packagePath)
	}

	return pkgRealpath
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
