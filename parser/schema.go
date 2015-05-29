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
)

var (
	rxStrFmt        = regexp.MustCompile("\\+swagger:strfmt\\p{Zs}*(\\p{L}[\\p{L}\\p{N}-]+)$")
	rxModelOverride = regexp.MustCompile("\\+swagger:model\\p{Zs}*(\\p{L}[\\p{L}\\p{N}-]+)?$")

	rxMaximumFmt    = "%s[Mm]ax(?:imum)?\\p{Zs}*:\\p{Zs}*([\\<=])?\\p{Zs}*([\\+-]?(?:\\p{N}+\\.)?\\p{N}+)$"
	rxMinimumFmt    = "%s[Mm]in(?:imum)?\\p{Zs}*:\\p{Zs}*([\\>=])?\\p{Zs}*([\\+-]?(?:\\p{N}+\\.)?\\p{N}+)$"
	rxMultipleOfFmt = "%s[Mm]ultiple\\p{Zs}*[Oo]f\\p{Zs}*:\\p{Zs}*([\\+-]?(?:\\p{N}+\\.)?\\p{N}+)$"

	rxMaxLengthFmt = "%s[Mm]ax(?:imum)?(?:\\p{Zs}*-?[Ll]en(?:gth)?)\\p{Zs}*:\\p{Zs}*(\\p{N}+)$"
	rxMinLengthFmt = "%s[Mm]in(?:imum)?(?:\\p{Zs}*-?[Ll]en(?:gth)?)\\p{Zs}*:\\p{Zs}*(\\p{N}+)$"
	rxPatternFmt   = "%s(?:P|p)attern\\p{Zs}*:\\p{Zs}*(.*)$"

	rxMaxItemsFmt = "%s[Mm]ax(?:imum)?(?:\\p{Zs}*|-)?[Ii]tems\\p{Zs}*:\\p{Zs}*(\\p{N}+)$"
	rxMinItemsFmt = "%s[Mm]in(?:imum)?(?:\\p{Zs}*|-)?[Ii]tems\\p{Zs}*:\\p{Zs}*(\\p{N}+)$"
	rxUniqueFmt   = "%s[Uu]nique\\p{Zs}*:\\p{Zs}*(true|false)$"

	rxRequired = regexp.MustCompile("[Rr]equired\\p{Zs}*:\\p{Zs}*(true|false)$")
	rxReadOnly = regexp.MustCompile("[Rr]ead(?:\\p{Zs}*|-)?[Oo]nly\\p{Zs}*:\\p{Zs}*(true|false)$")

	rxItemsPrefix = "(?:[Ii]tems[\\.\\p{Zs}]?)+"
)

type schemaSetter func(*spec.Schema, []string) error
type matchingSchemaSetter func(*regexp.Regexp) schemaSetter

func rxf(rxp, ar string) *regexp.Regexp {
	return regexp.MustCompile(fmt.Sprintf(rxp, ar))
}

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
	t.stripsTag = false

	t.set = func(obj interface{}, lines []string) error { return setter(obj.(*spec.Schema), lines) }
	return
}

func newFieldSection(name string, matcher *regexp.Regexp, ms matchingSchemaSetter) (t *sectionTagger) {
	t = newSectionTagger(name, false)
	t.stripsTag = false
	t.matcher = matcher
	setter := ms(matcher)
	t.set = func(obj interface{}, lines []string) error { return setter(obj.(*spec.Schema), lines) }
	return
}

func newSchemaFieldSection(name string, matcher *regexp.Regexp, setter schemaSetter) (t *sectionTagger) {
	t = newSectionTagger(name, false)
	t.stripsTag = false
	t.matcher = matcher
	t.set = func(obj interface{}, lines []string) error { return setter(obj.(*spec.Schema), lines) }
	return
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

func setSchemaMaximum(rx *regexp.Regexp) schemaSetter {
	return func(schema *spec.Schema, lines []string) error {
		if len(lines) == 0 || (len(lines) == 1 && len(lines[0]) == 0) {
			return nil
		}
		matches := rx.FindStringSubmatch(lines[0])
		if len(matches) > 2 && len(matches[2]) > 0 {
			max, err := strconv.ParseFloat(matches[2], 64)
			if err != nil {
				return err
			}
			schema.Maximum = &max
			schema.ExclusiveMaximum = matches[1] == "<"
		}
		return nil
	}
}

func setSchemaMinimum(rx *regexp.Regexp) schemaSetter {
	return func(schema *spec.Schema, lines []string) error {
		if len(lines) == 0 || (len(lines) == 1 && len(lines[0]) == 0) {
			return nil
		}
		matches := rx.FindStringSubmatch(lines[0])
		if len(matches) > 2 && len(matches[2]) > 0 {
			min, err := strconv.ParseFloat(matches[2], 64)
			if err != nil {
				return err
			}
			schema.Minimum = &min
			schema.ExclusiveMinimum = matches[1] == ">"
		}
		return nil
	}
}

func setSchemaMultipleOf(rx *regexp.Regexp) schemaSetter {
	return func(schema *spec.Schema, lines []string) error {
		if len(lines) == 0 || (len(lines) == 1 && len(lines[0]) == 0) {
			return nil
		}
		matches := rx.FindStringSubmatch(lines[0])
		if len(matches) > 1 && len(matches[1]) > 0 {
			multipleOf, err := strconv.ParseFloat(matches[1], 64)
			if err != nil {
				return err
			}
			schema.MultipleOf = &multipleOf
		}
		return nil
	}
}

func setSchemaMaxItems(rx *regexp.Regexp) schemaSetter {
	return func(schema *spec.Schema, lines []string) error {
		if len(lines) == 0 || (len(lines) == 1 && len(lines[0]) == 0) {
			return nil
		}
		matches := rx.FindStringSubmatch(lines[0])
		if len(matches) > 1 && len(matches[1]) > 0 {
			maxItems, err := strconv.ParseInt(matches[1], 10, 64)
			if err != nil {
				return err
			}
			schema.MaxItems = &maxItems
		}
		return nil
	}
}

func setSchemaMinItems(rx *regexp.Regexp) schemaSetter {
	return func(schema *spec.Schema, lines []string) error {
		if len(lines) == 0 || (len(lines) == 1 && len(lines[0]) == 0) {
			return nil
		}
		matches := rx.FindStringSubmatch(lines[0])
		if len(matches) > 1 && len(matches[1]) > 0 {
			minItems, err := strconv.ParseInt(matches[1], 10, 64)
			if err != nil {
				return err
			}
			schema.MinItems = &minItems
		}
		return nil
	}
}

func setSchemaMaxLength(rx *regexp.Regexp) schemaSetter {
	return func(schema *spec.Schema, lines []string) error {
		if len(lines) == 0 || (len(lines) == 1 && len(lines[0]) == 0) {
			return nil
		}
		matches := rx.FindStringSubmatch(lines[0])
		if len(matches) > 1 && len(matches[1]) > 0 {
			maxLength, err := strconv.ParseInt(matches[1], 10, 64)
			if err != nil {
				return err
			}
			schema.MaxLength = &maxLength
		}
		return nil
	}
}

func setSchemaMinLength(rx *regexp.Regexp) schemaSetter {
	return func(schema *spec.Schema, lines []string) error {
		if len(lines) == 0 || (len(lines) == 1 && len(lines[0]) == 0) {
			return nil
		}
		matches := rx.FindStringSubmatch(lines[0])
		if len(matches) > 1 && len(matches[1]) > 0 {
			minLength, err := strconv.ParseInt(matches[1], 10, 64)
			if err != nil {
				return err
			}
			schema.MinLength = &minLength
		}
		return nil
	}
}

func setSchemaPattern(rx *regexp.Regexp) schemaSetter {
	return func(schema *spec.Schema, lines []string) error {
		if len(lines) == 0 || (len(lines) == 1 && len(lines[0]) == 0) {
			return nil
		}
		matches := rx.FindStringSubmatch(lines[0])
		if len(matches) > 1 && len(matches[1]) > 0 {
			schema.Pattern = matches[1]
		}
		return nil
	}
}

func setSchemaUnique(rx *regexp.Regexp) schemaSetter {
	return func(schema *spec.Schema, lines []string) error {
		if len(lines) == 0 || (len(lines) == 1 && len(lines[0]) == 0) {
			return nil
		}
		matches := rx.FindStringSubmatch(lines[0])
		if len(matches) > 1 && len(matches[1]) > 0 {
			req, err := strconv.ParseBool(matches[1])
			if err != nil {
				return err
			}
			schema.UniqueItems = req
		}
		return nil
	}
}

func setSchemaReadOnly(schema *spec.Schema, lines []string) error {
	if len(lines) == 0 || (len(lines) == 1 && len(lines[0]) == 0) {
		return nil
	}
	matches := rxReadOnly.FindStringSubmatch(lines[0])
	if len(matches) > 1 && len(matches[1]) > 0 {
		req, err := strconv.ParseBool(matches[1])
		if err != nil {
			return err
		}
		schema.ReadOnly = req
	}
	return nil
}

func setSchemaRequired(parent *spec.Schema, value string) func(*spec.Schema, []string) error {
	return func(schema *spec.Schema, lines []string) error {
		if len(lines) == 0 || (len(lines) == 1 && len(lines[0]) == 0) {
			return nil
		}
		matches := rxRequired.FindStringSubmatch(lines[0])
		if len(matches) > 1 && len(matches[1]) > 0 {
			req, err := strconv.ParseBool(matches[1])
			if err != nil {
				return err
			}
			midx := -1
			for i, nm := range parent.Required {
				if nm == value {
					midx = i
					break
				}
			}
			if req {
				if midx < 0 {
					parent.Required = append(parent.Required, value)
				}
			} else if midx >= 0 {
				parent.Required = append(parent.Required[:midx], parent.Required[midx+1:]...)
			}
		}
		return nil
	}
}

type structDecl struct {
	File     *ast.File
	Decl     *ast.GenDecl
	TypeSpec *ast.TypeSpec
	GoName   string
	Name     string
}

func (sd *structDecl) inferNames() (goName string, name string) {
	if sd.GoName != "" {
		goName, name = sd.GoName, sd.Name
		return
	}
	goName = sd.TypeSpec.Name.Name
	name = goName
	if sd.Decl.Doc != nil {
		for _, cmt := range sd.Decl.Doc.List {
			for _, ln := range strings.Split(cmt.Text, "\n") {
				matches := rxModelOverride.FindStringSubmatch(ln)
				if len(matches) > 1 && len(matches[1]) > 0 {
					name = matches[1]
				}
			}
		}
	}
	sd.GoName = goName
	sd.Name = name
	return
}

type structCommentParser struct {
	taggers []*sectionTagger
	header  struct {
		taggers   []*sectionTagger
		otherTags []string
	}
	program   *loader.Program
	postDecls []structDecl
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
	for _, decl := range gofile.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		for _, spc := range gd.Specs {
			if ts, ok := spc.(*ast.TypeSpec); ok {
				sd := structDecl{gofile, gd, ts, "", ""}
				sd.inferNames()
				if err := scp.parseDecl(tgt, sd); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (scp *structCommentParser) parseDecl(definitions map[string]spec.Schema, decl structDecl) error {
	// check if there is a +swagger:model tag that is followed by a word,
	// this word is the type name for swagger
	// the package and type are recorded in the extensions
	// once type name is found convert it to a schema, by looking up the schema in the
	// definitions dictionary that got passed into this parse method
	decl.inferNames()
	schema := definitions[decl.Name]
	schPtr := &schema

	// analyze doc comment for the model
	// first line of the doc comment is the title
	// all following lines are description
	// all other things are ignored and by definition added to the last matched tag unless
	// preceded by 2 new lines
	if err := parseDocComments(decl.Decl.Doc, schPtr, scp.header.taggers, scp.header.otherTags); err != nil {
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
	if tpe, ok := decl.TypeSpec.Type.(*ast.StructType); ok {
		if err := scp.parseStructType(decl.File, schPtr, tpe); err != nil {
			return err
		}
	}
	if decl.Name != decl.GoName {
		schPtr.AddExtension("x-go-name", decl.GoName)
	}
	// TODO: perhaps move this to the classifier
	// and build a map from file pos to package
	for _, pkgInfo := range scp.program.AllPackages {
		if pkgInfo.Importable {
			for _, fil := range pkgInfo.Files {
				if fil.Pos() == decl.File.Pos() {
					schPtr.AddExtension("x-go-package", pkgInfo.Pkg.Path())
				}
			}
		}
	}
	definitions[decl.Name] = schema
	return nil
}

func (scp *structCommentParser) parseStructType(gofile *ast.File, schema *spec.Schema, tpe *ast.StructType) error {
	schema.Type = spec.StringOrArray([]string{"object"})
	if tpe.Fields != nil {
		seenProperties := make(map[string]struct{})
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
				if err := scp.parseProperty(gofile, fld.Type, &ps); err != nil {
					return err
				}

				// check if this is a primitive, if so parse the validations from the
				// doc comments of the slice declaration.
				if ftpe, ok := fld.Type.(*ast.ArrayType); ok {
					if iftpe, ok := ftpe.Elt.(*ast.Ident); ok && iftpe.Obj == nil {
						if ps.Items.Schema != nil {
							if err := scp.parseItemsDocComments(gofile, fld, ps.Items.Schema); err != nil {
								return err
							}
						} else {
							for _, sch := range ps.Items.Schemas {
								if err := scp.parseItemsDocComments(gofile, fld, &sch); err != nil {
									return err
								}
							}
						}
					}
				}

				if ps.Ref.GetURL() == nil {
					// add title and description for property
					// add validations for property
					taggers := []*sectionTagger{
						newSchemaDescription(setSchemaDescription),
						newFieldSection("maximum", rxf(rxMaximumFmt, ""), setSchemaMaximum),
						newFieldSection("minimum", rxf(rxMinimumFmt, ""), setSchemaMinimum),
						newFieldSection("multipleOf", rxf(rxMultipleOfFmt, ""), setSchemaMultipleOf),
						newFieldSection("minLength", rxf(rxMinLengthFmt, ""), setSchemaMinLength),
						newFieldSection("maxLength", rxf(rxMaxLengthFmt, ""), setSchemaMaxLength),
						newFieldSection("pattern", rxf(rxPatternFmt, ""), setSchemaPattern),
						newFieldSection("minItems", rxf(rxMinItemsFmt, ""), setSchemaMinItems),
						newFieldSection("maxItems", rxf(rxMaxItemsFmt, ""), setSchemaMaxItems),
						newFieldSection("unique", rxf(rxUniqueFmt, ""), setSchemaUnique),
						newSchemaFieldSection("readOnly", rxReadOnly, setSchemaReadOnly),
						newSchemaFieldSection("required", rxRequired, setSchemaRequired(schema, nm)),
					}
					parseDocComments(fld.Doc, &ps, taggers, nil)
				}

				if nm != gnm {
					ps.AddExtension("x-go-name", gnm)
				}
				if schema.Properties == nil {
					schema.Properties = make(map[string]spec.Schema)
				}
				seenProperties[nm] = struct{}{}
				schema.Properties[nm] = ps
			}
		}

		for k := range schema.Properties {
			if _, ok := seenProperties[k]; !ok {
				delete(schema.Properties, k)
			}
		}
	}

	return nil
}

func (scp *structCommentParser) parseItemsDocComments(gofile *ast.File, fld *ast.Field, prop *spec.Schema) error {
	// add title and description for property
	// add validations for property
	taggers := []*sectionTagger{
		newFieldSection("maximum", rxf(rxMaximumFmt, rxItemsPrefix), setSchemaMaximum),
		newFieldSection("minimum", rxf(rxMinimumFmt, rxItemsPrefix), setSchemaMinimum),
		newFieldSection("multipleOf", rxf(rxMultipleOfFmt, rxItemsPrefix), setSchemaMultipleOf),
		newFieldSection("minLength", rxf(rxMinLengthFmt, rxItemsPrefix), setSchemaMinLength),
		newFieldSection("maxLength", rxf(rxMaxLengthFmt, rxItemsPrefix), setSchemaMaxLength),
		newFieldSection("pattern", rxf(rxPatternFmt, rxItemsPrefix), setSchemaPattern),
		newFieldSection("minItems", rxf(rxMinItemsFmt, rxItemsPrefix), setSchemaMinItems),
		newFieldSection("maxItems", rxf(rxMaxItemsFmt, rxItemsPrefix), setSchemaMaxItems),
		newFieldSection("unique", rxf(rxUniqueFmt, rxItemsPrefix), setSchemaUnique),
	}
	return parseDocComments(fld.Doc, prop, taggers, nil)
}

func (scp *structCommentParser) parseProperty(gofile *ast.File, fld ast.Expr, prop *spec.Schema) error {
	switch ftpe := fld.(type) {
	case *ast.Ident: // simple value
		if ftpe.Obj == nil {
			return swaggerSchemaForType(ftpe.Name, prop)
		}
		// we're probably looking at a struct here
		// make sure it is one. Try to find it in the package
		// when found make sure the struct gets added as a schema too
		// and turn this property into a ref
		if ftpe.Obj.Kind == ast.Typ {
			if ts, ok := ftpe.Obj.Decl.(*ast.TypeSpec); ok {
				if _, ok := ts.Type.(*ast.StructType); ok {
					ref, err := spec.NewRef("#/definitions/" + ts.Name.Name)
					if err != nil {
						return err
					}
					prop.Ref = ref
				DECLS:
					for _, d := range gofile.Decls {
						if gd, ok := d.(*ast.GenDecl); ok {
							for _, tss := range gd.Specs {
								if tss.Pos() == ts.Pos() {
									sd := structDecl{gofile, gd, ts, "", ""}
									sd.inferNames()
									scp.postDecls = append(scp.postDecls, sd)
									break DECLS
								}
							}
						}
					}
				}
			}
		}

	case *ast.StarExpr: // pointer to something, optional by default
		scp.parseProperty(gofile, ftpe.X, prop)

	case *ast.ArrayType: // slice type
		var items *spec.Schema
		if prop.Items != nil && prop.Items.Schema != nil {
			items = prop.Items.Schema
		}
		if items == nil {
			items = new(spec.Schema)
		}
		if err := scp.parseProperty(gofile, ftpe.Elt, items); err != nil {
			return err
		}
		prop.Typed("array", "")
		if prop.Items == nil {
			prop.Items = new(spec.SchemaOrArray)
		}
		prop.Items.Schema = items
		prop.Items.Schemas = nil

	case *ast.StructType:
		return scp.parseStructType(gofile, prop, ftpe)

	case *ast.SelectorExpr:
		return scp.swaggerSchemaForSelector(gofile, ftpe, prop)

	case *ast.MapType:
		// check if key is a string type, if not print a message
		// and skip the map property. Only maps with string keys can go into additional properties
		if keyIdent, ok := ftpe.Key.(*ast.Ident); ok {
			if keyIdent.Name == "string" {
				if prop.AdditionalProperties == nil {
					prop.AdditionalProperties = new(spec.SchemaOrBool)
				}
				prop.AdditionalProperties.Allows = false
				if prop.AdditionalProperties.Schema == nil {
					prop.AdditionalProperties.Schema = new(spec.Schema)
				}
				scp.parseProperty(gofile, ftpe.Value, prop.AdditionalProperties.Schema)
				prop.Typed("object", "")
			}
		}

	case *ast.InterfaceType:
		// FIXME:
		// what to do with an interface? support it?
		// ignoring it for now
		// I guess something can be done with a discriminator field
		// but is it worth the trouble?
	default:
		fmt.Println("???", ftpe)
	}
	return nil
}

func (scp *structCommentParser) swaggerSchemaForSelector(gofile *ast.File, expr *ast.SelectorExpr, prop *spec.Schema) error {
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

		pkg := scp.program.Package(selPath)
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
									prop.Ref = ref
									sd := structDecl{file, gd, ts, "", ""}
									sd.inferNames()
									scp.postDecls = append(scp.postDecls, sd)
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
