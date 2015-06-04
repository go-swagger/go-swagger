package parse

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

func newSchemaAnnotationParser(goName string) *schemaAnnotationParser {
	return &schemaAnnotationParser{GoName: goName, rx: rxModelOverride}
}

type schemaAnnotationParser struct {
	GoName string
	Name   string
	rx     *regexp.Regexp
}

func (sap *schemaAnnotationParser) Matches(line string) bool {
	return sap.rx.MatchString(line)
}

func (sap *schemaAnnotationParser) Parse(lines []string) error {
	if sap.Name != "" {
		return nil
	}

	if len(lines) > 0 {
		for _, line := range lines {
			matches := sap.rx.FindStringSubmatch(line)
			if len(matches) > 1 && len(matches[1]) > 0 {
				sap.Name = matches[1]
				return nil
			}
		}
	}
	return nil
}

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

type schemaParser struct {
	program   *loader.Program
	postDecls []schemaDecl
}

func newSchemaParser(prog *loader.Program) *schemaParser {
	scp := new(schemaParser)
	scp.program = prog
	return scp
}

func (scp *schemaParser) Parse(gofile *ast.File, target interface{}) error {
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

func (scp *schemaParser) parseDecl(definitions map[string]spec.Schema, decl schemaDecl) error {
	// check if there is a +swagger:model tag that is followed by a word,
	// this word is the type name for swagger
	// the package and type are recorded in the extensions
	// once type name is found convert it to a schema, by looking up the schema in the
	// definitions dictionary that got passed into this parse method
	decl.inferNames()
	schema := definitions[decl.Name]
	schPtr := &schema

	// analyze doc comment for the model
	sp := new(sectionedParser)
	sp.setTitle = func(lines []string) { schema.Title = joinDropLast(lines) }
	sp.setDescription = func(lines []string) { schema.Description = joinDropLast(lines) }
	if err := sp.Parse(decl.Decl.Doc); err != nil {
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

func (scp *schemaParser) parseStructType(gofile *ast.File, schema *spec.Schema, tpe *ast.StructType) error {
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

				sp := new(sectionedParser)
				sp.setDescription = func(lines []string) { ps.Description = joinDropLast(lines) }
				if ps.Ref.GetURL() == nil {
					sp.taggers = []tagParser{
						newSingleLineTagParser("maximum", &setMaximum{schemaValidations{&ps}, rxf(rxMaximumFmt, "")}),
						newSingleLineTagParser("minimum", &setMinimum{schemaValidations{&ps}, rxf(rxMinimumFmt, "")}),
						newSingleLineTagParser("multipleOf", &setMultipleOf{schemaValidations{&ps}, rxf(rxMultipleOfFmt, "")}),
						newSingleLineTagParser("minLength", &setMinLength{schemaValidations{&ps}, rxf(rxMinLengthFmt, "")}),
						newSingleLineTagParser("maxLength", &setMaxLength{schemaValidations{&ps}, rxf(rxMaxLengthFmt, "")}),
						newSingleLineTagParser("pattern", &setPattern{schemaValidations{&ps}, rxf(rxPatternFmt, "")}),
						newSingleLineTagParser("minItems", &setMinItems{schemaValidations{&ps}, rxf(rxMinItemsFmt, "")}),
						newSingleLineTagParser("maxItems", &setMaxItems{schemaValidations{&ps}, rxf(rxMaxItemsFmt, "")}),
						newSingleLineTagParser("unique", &setUnique{schemaValidations{&ps}, rxf(rxUniqueFmt, "")}),
						newSingleLineTagParser("required", &setRequiredSchema{schema, nm}),
						newSingleLineTagParser("readOnly", &setReadOnlySchema{&ps}),
					}

					// check if this is a primitive, if so parse the validations from the
					// doc comments of the slice declaration.
					if ftpe, ok := fld.Type.(*ast.ArrayType); ok {
						if iftpe, ok := ftpe.Elt.(*ast.Ident); ok && iftpe.Obj == nil {
							if ps.Items != nil && ps.Items.Schema != nil {
								itemsTaggers := []tagParser{
									newSingleLineTagParser("itemsMaximum", &setMaximum{schemaValidations{ps.Items.Schema}, rxf(rxMaximumFmt, rxItemsPrefix)}),
									newSingleLineTagParser("itemsMinimum", &setMinimum{schemaValidations{ps.Items.Schema}, rxf(rxMinimumFmt, rxItemsPrefix)}),
									newSingleLineTagParser("itemsMultipleOf", &setMultipleOf{schemaValidations{ps.Items.Schema}, rxf(rxMultipleOfFmt, rxItemsPrefix)}),
									newSingleLineTagParser("itemsMinLength", &setMinLength{schemaValidations{ps.Items.Schema}, rxf(rxMinLengthFmt, rxItemsPrefix)}),
									newSingleLineTagParser("itemsMaxLength", &setMaxLength{schemaValidations{ps.Items.Schema}, rxf(rxMaxLengthFmt, rxItemsPrefix)}),
									newSingleLineTagParser("itemsPattern", &setPattern{schemaValidations{ps.Items.Schema}, rxf(rxPatternFmt, rxItemsPrefix)}),
									newSingleLineTagParser("itemsMinItems", &setMinItems{schemaValidations{ps.Items.Schema}, rxf(rxMinItemsFmt, rxItemsPrefix)}),
									newSingleLineTagParser("itemsMaxItems", &setMaxItems{schemaValidations{ps.Items.Schema}, rxf(rxMaxItemsFmt, rxItemsPrefix)}),
									newSingleLineTagParser("itemsUnique", &setUnique{schemaValidations{ps.Items.Schema}, rxf(rxUniqueFmt, rxItemsPrefix)}),
								}

								// items matchers should go before the default matchers so they match first
								sp.taggers = append(itemsTaggers, sp.taggers...)
							}
						}
					}
				} else {
					sp.taggers = []tagParser{
						newSingleLineTagParser("required", &setRequiredSchema{schema, nm}),
					}
				}
				if err := sp.Parse(fld.Doc); err != nil {
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

func (scp *schemaParser) parseProperty(gofile *ast.File, fld ast.Expr, prop *spec.Schema) error {
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

					for _, d := range gofile.Decls {
						if gd, ok := d.(*ast.GenDecl); ok {
							for _, tss := range gd.Specs {
								if tss.Pos() == ts.Pos() {
									sd := schemaDecl{gofile, gd, ts, "", ""}
									sd.inferNames()
									ref, err := spec.NewRef("#/definitions/" + sd.Name)
									if err != nil {
										return err
									}
									prop.Ref = ref
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
