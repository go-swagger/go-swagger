package parser

import (
	"fmt"
	"go/ast"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/tools/go/loader"

	"github.com/casualjim/go-swagger/spec"
	"github.com/kr/pretty"
)

var (
	rxStrFmt = regexp.MustCompile("\\+swagger:strfmt[^\\w]*((?:\\w+-?)+)")
)

type schemaSetter func(*spec.Schema, []string) error

func newSchemaTitle(setter schemaSetter) (t *sectionTagger) {
	t = newTitleTagger()
	t.rxStripComments = rxStripComments
	t.set = func(obj interface{}, lines []string) error { return setter(obj.(*spec.Schema), lines) }
	return
}
func newSchemaDescription(setter schemaSetter) (t *sectionTagger) {
	t = newDescriptionTagger()
	t.set = func(obj interface{}, lines []string) error { return setter(obj.(*spec.Schema), lines) }
	return
}
func newSchemaSection(name string, multiLine bool, setter schemaSetter) (t *sectionTagger) {
	t = newSectionTagger(name, multiLine)
	t.set = func(obj interface{}, lines []string) error { return setter(obj.(*spec.Schema), lines) }
	return
}

type structCommentParser struct {
	taggers []*sectionTagger
	header  struct {
		taggers   []*sectionTagger
		otherTags []string
	}
	program *loader.Program
}

func schemaParser(prog *loader.Program) *structCommentParser {
	scp := new(structCommentParser)
	scp.program = prog
	scp.header.taggers = []*sectionTagger{newSchemaTitle(setSchemaTitle), newSchemaDescription(setSchemaDescription)}
	scp.header.otherTags = []string{"+swagger"}
	return scp
}

func (scp *structCommentParser) Parse(gofile *ast.File, target interface{}) error {
	tgt := target.(map[string]spec.Schema)
	pretty.Println(gofile.Imports)
	for _, decl := range gofile.Decls {
		switch gd := decl.(type) {
		case *ast.GenDecl:
			if len(gd.Specs) > 0 {
				if ts, ok := gd.Specs[0].(*ast.TypeSpec); ok {
					// TODO: really infer type name
					// check if there is a +swagger:model tag that is followed by a word, this word is the type name for swagger
					// that word will also go in the X-GO- extensions map
					// once type name is found convert it to a schema, by looking up the schema in the
					// definitions dictionary that got passed into this Parse method
					schema := tgt[ts.Name.Name]

					// analyze doc comment for the model
					// first line of the doc comment is the title
					// all following lines are description
					// all other things are ignored and by definition added to the last matched tag unless
					// preceded by 2 new lines
					if err := parseDocComments(gd.Doc, &schema, scp.header.taggers, scp.header.otherTags); err != nil {
						return err
					}

					// analyze struct body for fields etc
					// each exported struct field:
					// * gets a type mapped to a go primitive
					// * perhaps gets a format
					// * has to document the validations that apply for the type and the field
					// * when the struct field points to a model it becomes a ref: #/definitions/ModelName
					// * the first line of the comment is the title
					// * the following lines are the description
					if tpe, ok := ts.Type.(*ast.StructType); ok {
						if err := scp.parseStructType(gofile, &schema, tpe); err != nil {
							return err
						}
					}

					tgt[ts.Name.Name] = schema
				}
			}
		default:
			//fmt.Println("unhandled decl:", gd)
		}
	}
	return nil
}

func (scp *structCommentParser) parseStructType(gofile *ast.File, schema *spec.Schema, tpe *ast.StructType) error {
	schema.Type = spec.StringOrArray([]string{"object"})
	if tpe.Fields != nil {
		// TODO: remove the properties that no longer exist
		for _, fld := range tpe.Fields.List {
			var nm, gnm string
			if len(fld.Names) > 0 && fld.Names[0] != nil && fld.Names[0].IsExported() {
				nm = fld.Names[0].Name
				gnm = nm
				if fld.Tag != nil && len(strings.TrimSpace(fld.Tag.Value)) > 0 {
					tv, err := strconv.Unquote(fld.Tag.Value)
					if err != nil {
						return err
					}

					if strings.TrimSpace(tv) != "" {
						st := reflect.StructTag(tv)
						if st.Get("json") != "" {
							nm = strings.Split(st.Get("json"), ",")[0]
						}
					}
				}

				ps := schema.Properties[nm]
				if err := scp.parseProperty(gofile, fld, &ps); err != nil {
					return err
				}
				if nm != gnm {
					if ps.Extensions == nil {
						ps.Extensions = make(map[string]interface{})
					}
					ps.Extensions.Add("x-go-name", gnm)
				}
				if schema.Properties == nil {
					schema.Properties = make(map[string]spec.Schema)
				}
				schema.Properties[nm] = ps

				// pretty.Println(fld)
			}
		}
	}

	return nil
}

func (scp *structCommentParser) parseProperty(gofile *ast.File, fld *ast.Field, prop *spec.Schema) error {
	switch ftpe := fld.Type.(type) {
	case *ast.Ident: // simple value
		return swaggerSchemaForType(ftpe.Name, prop)
	case *ast.StarExpr: // pointer to something, optional by default
		fmt.Println("star", ftpe.X.(*ast.Ident).Name)
	case *ast.ArrayType: // slice type
		fmt.Println("slice", ftpe.Elt)
	case *ast.StructType:
		fmt.Println("embedded", ftpe.Fields)
	case *ast.InterfaceType:
		fmt.Println("interface")
	case *ast.SelectorExpr:
		return swaggerSchemaForStringFormat(scp.program, gofile, ftpe, prop)
	default:
		fmt.Println("???", ftpe)
	}
	// add title and description
	// add validations
	return nil
}

func setSchemaTitle(schema *spec.Schema, lines []string) error {
	schema.Title = joinDropLast(lines)
	return nil
}

func setSchemaDescription(schema *spec.Schema, lines []string) error {
	schema.Description = joinDropLast(lines)
	return nil
}

func joinDropLast(lines []string) string {
	l := len(lines)
	lns := lines
	if l > 0 && len(strings.TrimSpace(lines[l-1])) == 0 {
		lns = lines[:l-1]
	}
	return strings.Join(lns, "\n")
}

func swaggerSchemaForStringFormat(prog *loader.Program, gofile *ast.File, expr *ast.SelectorExpr, prop *spec.Schema) error {
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

		pkg := prog.Package(selPath)
		if pkg == nil {
			return fmt.Errorf("no package found for %s", selPath)
		}

		// look at doc comments for +swagger:strfmt [name]
		// when found this is the format name, create a schema with that name
		for _, file := range pkg.Files {
			for _, decl := range file.Decls {
				if gd, ok := decl.(*ast.GenDecl); ok {
					for _, gs := range gd.Specs {
						if ts, ok := gs.(*ast.TypeSpec); ok {
							if ts.Name != nil && ts.Name.Name == expr.Sel.Name {
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
						}
					}
				}
			}
		}
		return fmt.Errorf("schema parser: no string format for %s.%s", pth.Name, expr.Sel.Name)
	}
	return fmt.Errorf("schema parser: no string format for %v", expr.Sel.Name)
}

func swaggerSchemaForType(typeName string, prop *spec.Schema) error {
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

//var strfmtMapping = map[string]func() spec.Schema{
//"strfmt.Base64": func() spec.Schema { return *spec.StrFmtProperty("base64") },
//"strfmt.Date",
//"strfmt.DateTime",
//"strfmt.URI",
//"strfmt.Email",
//"strfmt.Hostname",
//"strfmt.IPv4",
//"strfmt.IPv6",
//"strfmt.UUID",
//"strfmt.UUID3",
//"strfmt.UUID4",
//"strfmt.UUID5",
//"strfmt.ISBN",
//"strfmt.ISBN10",
//"strfmt.ISBN13",
//"strfmt.CreditCard",
//"strfmt.SSN",
//"strfmt.HexColor",
//"strfmt.RGBColor",
//"strfmt.Duration",
//}
