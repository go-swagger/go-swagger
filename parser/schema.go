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

type schemaSetter func(*spec.Schema, []string) error
type matchingSchemaSetter func(*regexp.Regexp) schemaSetter

type schemaTypable struct {
	schema *spec.Schema
}

func (st *schemaTypable) Typed(tpe, format string) {
	st.schema.Typed(tpe, format)
}

func (st *schemaTypable) SetRef(ref spec.Ref) {
	st.schema.Ref = ref
}

type schemaValidations struct {
	current *spec.Schema
}

func (sv schemaValidations) SetMaximum(val float64, exclusive bool) {
	sv.current.Maximum = &val
	sv.current.ExclusiveMaximum = exclusive
}
func (sv schemaValidations) SetMinimum(val float64, exclusive bool) {
	sv.current.Minimum = &val
	sv.current.ExclusiveMinimum = exclusive
}
func (sv schemaValidations) SetMultipleOf(val float64) { sv.current.MultipleOf = &val }
func (sv schemaValidations) SetMinItems(val int64)     { sv.current.MinItems = &val }
func (sv schemaValidations) SetMaxItems(val int64)     { sv.current.MaxItems = &val }
func (sv schemaValidations) SetMinLength(val int64)    { sv.current.MinLength = &val }
func (sv schemaValidations) SetMaxLength(val int64)    { sv.current.MaxLength = &val }
func (sv schemaValidations) SetPattern(val string)     { sv.current.Pattern = val }
func (sv schemaValidations) SetUnique(val bool)        { sv.current.UniqueItems = val }

type schemaDecl struct {
	File     *ast.File
	Decl     *ast.GenDecl
	TypeSpec *ast.TypeSpec
	GoName   string
	Name     string
}

func (sd *schemaDecl) inferNames() (goName string, name string) {
	if sd.GoName != "" {
		goName, name = sd.GoName, sd.Name
		return
	}
	goName = sd.TypeSpec.Name.Name
	name = goName
	if sd.Decl.Doc != nil {
	DECLS:
		for _, cmt := range sd.Decl.Doc.List {
			for _, ln := range strings.Split(cmt.Text, "\n") {
				matches := rxModelOverride.FindStringSubmatch(ln)
				if len(matches) > 1 && len(matches[1]) > 0 {
					name = matches[1]
					break DECLS
				}
			}
		}
	}
	sd.GoName = goName
	sd.Name = name
	return
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
		bldr := setMaximum{schemaValidations{schema}, rx}
		return bldr.Parse(lines)
	}
}

func setSchemaMinimum(rx *regexp.Regexp) schemaSetter {
	return func(schema *spec.Schema, lines []string) error {
		bldr := setMinimum{schemaValidations{schema}, rx}
		return bldr.Parse(lines)
	}
}

func setSchemaMultipleOf(rx *regexp.Regexp) schemaSetter {
	return func(schema *spec.Schema, lines []string) error {
		bldr := setMultipleOf{schemaValidations{schema}, rx}
		return bldr.Parse(lines)
	}
}

func setSchemaMaxItems(rx *regexp.Regexp) schemaSetter {
	return func(schema *spec.Schema, lines []string) error {
		bldr := setMaxItems{schemaValidations{schema}, rx}
		return bldr.Parse(lines)
	}
}

func setSchemaMinItems(rx *regexp.Regexp) schemaSetter {
	return func(schema *spec.Schema, lines []string) error {
		bldr := setMinItems{schemaValidations{schema}, rx}
		return bldr.Parse(lines)
	}
}

func setSchemaMaxLength(rx *regexp.Regexp) schemaSetter {
	return func(schema *spec.Schema, lines []string) error {
		bldr := setMaxLength{schemaValidations{schema}, rx}
		return bldr.Parse(lines)
	}
}

func setSchemaMinLength(rx *regexp.Regexp) schemaSetter {
	return func(schema *spec.Schema, lines []string) error {
		bldr := setMinLength{schemaValidations{schema}, rx}
		return bldr.Parse(lines)

	}
}

func setSchemaPattern(rx *regexp.Regexp) schemaSetter {
	return func(schema *spec.Schema, lines []string) error {
		bldr := setPattern{schemaValidations{schema}, rx}
		return bldr.Parse(lines)
	}
}

func setSchemaUnique(rx *regexp.Regexp) schemaSetter {
	return func(schema *spec.Schema, lines []string) error {
		bldr := setUnique{schemaValidations{schema}, rx}
		return bldr.Parse(lines)
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

type structCommentParser struct {
	taggers []*sectionTagger
	header  struct {
		taggers   []*sectionTagger
		otherTags []string
	}
	program   *loader.Program
	postDecls []schemaDecl
}

func newSchemaParser(prog *loader.Program) *structCommentParser {
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
				sd := schemaDecl{gofile, gd, ts, "", ""}
				sd.inferNames()
				if err := scp.parseDecl(tgt, sd); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (scp *structCommentParser) parseDecl(definitions map[string]spec.Schema, decl schemaDecl) error {
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
	schema.Typed("object", "")
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

				var taggers []*sectionTagger
				if ps.Ref.GetURL() == nil {
					// add title and description for property
					// add validations for property
					taggers = []*sectionTagger{
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
				} else {
					// add title and description for property
					// add validations for property
					taggers = []*sectionTagger{
						newSchemaDescription(setSchemaDescription),
						newSchemaFieldSection("required", rxRequired, setSchemaRequired(schema, nm)),
					}
				}
				if err := parseDocComments(fld.Doc, &ps, taggers, nil); err != nil {
					return err
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
	sct := &schemaTypable{prop}
	switch ftpe := fld.(type) {
	case *ast.Ident: // simple value
		if ftpe.Obj == nil {
			return swaggerSchemaForType(ftpe.Name, sct)
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

					for _, d := range gofile.Decls {
						if gd, ok := d.(*ast.GenDecl); ok {
							for _, tss := range gd.Specs {
								if tss.Pos() == ts.Pos() {
									sd := schemaDecl{gofile, gd, ts, "", ""}
									sd.inferNames()
									scp.postDecls = append(scp.postDecls, sd)
									return nil
								}
							}
						}
					}
				}
			}
			return nil
		}

		return fmt.Errorf("couldn't infer type for %v", ftpe.Name)

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
		sp := selectorParser{
			program:     scp.program,
			AddPostDecl: func(sd schemaDecl) { scp.postDecls = append(scp.postDecls) },
		}
		return sp.TypeForSelector(gofile, ftpe, sct)

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
		return fmt.Errorf("%s is unsupported for a schema", ftpe)
	}
	return nil
}
