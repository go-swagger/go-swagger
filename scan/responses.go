package scan

import (
	"fmt"
	"go/ast"
	"reflect"
	"strconv"
	"strings"

	"golang.org/x/tools/go/loader"

	"github.com/casualjim/go-swagger/spec"
)

type responseTypable struct {
	header   *spec.Header
	response *spec.Response
}

func (ht responseTypable) Typed(tpe, format string) {
	ht.header.Typed(tpe, format)
}

func (ht responseTypable) Items() *spec.Items {
	return ht.header.Items
}

func (ht responseTypable) SetRef(ref spec.Ref) {
	// having trouble seeing the usefulness of this one here
}

func (ht responseTypable) Schema() *spec.Schema {
	return ht.response.Schema
}

func (ht responseTypable) SetSchema(schema *spec.Schema) {
	ht.response.Schema = schema
}
func (ht responseTypable) CollectionOf(items *spec.Items, format string) {
	ht.header.CollectionOf(items, format)
}

type headerValidations struct {
	current *spec.Header
}

func (sv headerValidations) SetMaximum(val float64, exclusive bool) {
	sv.current.Maximum = &val
	sv.current.ExclusiveMaximum = exclusive
}
func (sv headerValidations) SetMinimum(val float64, exclusive bool) {
	sv.current.Minimum = &val
	sv.current.ExclusiveMinimum = exclusive
}
func (sv headerValidations) SetMultipleOf(val float64)      { sv.current.MultipleOf = &val }
func (sv headerValidations) SetMinItems(val int64)          { sv.current.MinItems = &val }
func (sv headerValidations) SetMaxItems(val int64)          { sv.current.MaxItems = &val }
func (sv headerValidations) SetMinLength(val int64)         { sv.current.MinLength = &val }
func (sv headerValidations) SetMaxLength(val int64)         { sv.current.MaxLength = &val }
func (sv headerValidations) SetPattern(val string)          { sv.current.Pattern = val }
func (sv headerValidations) SetUnique(val bool)             { sv.current.UniqueItems = val }
func (sv headerValidations) SetCollectionFormat(val string) { sv.current.CollectionFormat = val }

type responseDecl struct {
	File     *ast.File
	Decl     *ast.GenDecl
	TypeSpec *ast.TypeSpec
	GoName   string
	Name     string
}

func (sd *responseDecl) inferNames() (goName string, name string) {
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
				matches := rxResponseOverride.FindStringSubmatch(ln)
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

func newResponseParser(prog *loader.Program) *responseParser {
	return &responseParser{prog, nil}
}

type responseParser struct {
	program   *loader.Program
	postDecls []schemaDecl
}

func (rp *responseParser) Parse(gofile *ast.File, target interface{}) error {
	tgt := target.(map[string]spec.Response)
	for _, decl := range gofile.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		for _, spc := range gd.Specs {
			if ts, ok := spc.(*ast.TypeSpec); ok {
				sd := responseDecl{gofile, gd, ts, "", ""}
				sd.inferNames()
				if err := rp.parseDecl(tgt, sd); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (rp *responseParser) parseDecl(responses map[string]spec.Response, decl responseDecl) error {
	// check if there is a +swagger:parameters tag that is followed by one or more words,
	// these words are the ids of the operations this parameter struct applies to
	// once type name is found convert it to a schema, by looking up the schema in the
	// parameters dictionary that got passed into this parse method
	response := responses[decl.Name]
	resPtr := &response

	// analyze struct body for fields etc
	// each exported struct field:
	// * gets a type mapped to a go primitive
	// * perhaps gets a format
	// * has to document the validations that apply for the type and the field
	// * when the struct field points to a model it becomes a ref: #/definitions/ModelName
	// * comments that aren't tags is used as the description
	if tpe, ok := decl.TypeSpec.Type.(*ast.StructType); ok {
		if err := rp.parseStructType(decl.File, resPtr, tpe); err != nil {
			return err
		}
	}

	responses[decl.Name] = response
	return nil
}

func (rp *responseParser) parseStructType(gofile *ast.File, response *spec.Response, tpe *ast.StructType) error {
	if tpe.Fields != nil {

		seenProperties := make(map[string]struct{})
		for _, fld := range tpe.Fields.List {
			var nm string
			if len(fld.Names) > 0 && fld.Names[0] != nil && fld.Names[0].IsExported() {
				nm = fld.Names[0].Name
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

				var in string
				// scan for param location first, this changes some behavior down the line
				if fld.Doc != nil {
					for _, cmt := range fld.Doc.List {
						for _, line := range strings.Split(cmt.Text, "\n") {
							matches := rxIn.FindStringSubmatch(line)
							if len(matches) > 0 && len(strings.TrimSpace(matches[1])) > 0 {
								in = strings.TrimSpace(matches[1])
							}
						}
					}
				}

				ps := response.Headers[nm]
				if err := rp.parseProperty(gofile, fld.Type, responseTypable{&ps, response}, in); err != nil {
					return err
				}

				sp := new(sectionedParser)
				sp.setDescription = func(lines []string) { ps.Description = joinDropLast(lines) }
				sp.taggers = []tagParser{
					newSingleLineTagParser("maximum", &setMaximum{headerValidations{&ps}, rxf(rxMaximumFmt, "")}),
					newSingleLineTagParser("minimum", &setMinimum{headerValidations{&ps}, rxf(rxMinimumFmt, "")}),
					newSingleLineTagParser("multipleOf", &setMultipleOf{headerValidations{&ps}, rxf(rxMultipleOfFmt, "")}),
					newSingleLineTagParser("minLength", &setMinLength{headerValidations{&ps}, rxf(rxMinLengthFmt, "")}),
					newSingleLineTagParser("maxLength", &setMaxLength{headerValidations{&ps}, rxf(rxMaxLengthFmt, "")}),
					newSingleLineTagParser("pattern", &setPattern{headerValidations{&ps}, rxf(rxPatternFmt, "")}),
					newSingleLineTagParser("collectionFormat", &setCollectionFormat{headerValidations{&ps}, rxf(rxCollectionFormatFmt, "")}),
					newSingleLineTagParser("minItems", &setMinItems{headerValidations{&ps}, rxf(rxMinItemsFmt, "")}),
					newSingleLineTagParser("maxItems", &setMaxItems{headerValidations{&ps}, rxf(rxMaxItemsFmt, "")}),
					newSingleLineTagParser("unique", &setUnique{headerValidations{&ps}, rxf(rxUniqueFmt, "")}),
				}
				itemsTaggers := func() []tagParser {
					return []tagParser{
						newSingleLineTagParser("itemsMaximum", &setMaximum{itemsValidations{ps.Items}, rxf(rxMaximumFmt, rxItemsPrefix)}),
						newSingleLineTagParser("itemsMinimum", &setMinimum{itemsValidations{ps.Items}, rxf(rxMinimumFmt, rxItemsPrefix)}),
						newSingleLineTagParser("itemsMultipleOf", &setMultipleOf{itemsValidations{ps.Items}, rxf(rxMultipleOfFmt, rxItemsPrefix)}),
						newSingleLineTagParser("itemsMinLength", &setMinLength{itemsValidations{ps.Items}, rxf(rxMinLengthFmt, rxItemsPrefix)}),
						newSingleLineTagParser("itemsMaxLength", &setMaxLength{itemsValidations{ps.Items}, rxf(rxMaxLengthFmt, rxItemsPrefix)}),
						newSingleLineTagParser("itemsPattern", &setPattern{itemsValidations{ps.Items}, rxf(rxPatternFmt, rxItemsPrefix)}),
						newSingleLineTagParser("itemsCollectionFormat", &setCollectionFormat{itemsValidations{ps.Items}, rxf(rxCollectionFormatFmt, rxItemsPrefix)}),
						newSingleLineTagParser("itemsMinItems", &setMinItems{itemsValidations{ps.Items}, rxf(rxMinItemsFmt, rxItemsPrefix)}),
						newSingleLineTagParser("itemsMaxItems", &setMaxItems{itemsValidations{ps.Items}, rxf(rxMaxItemsFmt, rxItemsPrefix)}),
						newSingleLineTagParser("itemsUnique", &setUnique{itemsValidations{ps.Items}, rxf(rxUniqueFmt, rxItemsPrefix)}),
					}
				}

				// check if this is a primitive, if so parse the validations from the
				// doc comments of the slice declaration.
				if ftpe, ok := fld.Type.(*ast.ArrayType); ok {
					if iftpe, ok := ftpe.Elt.(*ast.Ident); ok && iftpe.Obj == nil {
						if ps.Items != nil {
							// items matchers should go before the default matchers so they match first
							sp.taggers = append(itemsTaggers(), sp.taggers...)
						}
					}
				}

				if err := sp.Parse(fld.Doc); err != nil {
					return err
				}

				if in != "body" {
					seenProperties[nm] = struct{}{}
					if response.Headers == nil {
						response.Headers = make(map[string]spec.Header)
					}
					response.Headers[nm] = ps
				}
			}
		}

		for k := range response.Headers {
			if _, ok := seenProperties[k]; !ok {
				delete(response.Headers, k)
			}
		}
	}

	return nil
}

func (rp *responseParser) parseProperty(gofile *ast.File, fld ast.Expr, prop operationTypable, in string) error {
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
									prop.SetRef(ref)
									rp.postDecls = append(rp.postDecls, sd)
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

	case *ast.StarExpr: // pointer to something
		rp.parseProperty(gofile, ftpe.X, prop, in)

	case *ast.ArrayType: // slice type
		if in == "body" {
			var items spec.Schema
			scp := newSchemaParser(rp.program)
			if err := scp.parseProperty(gofile, ftpe.Elt, &items); err != nil {
				return err
			}
			prop.(responseTypable).SetSchema(new(spec.Schema).CollectionOf(items))
			return nil
		}
		var items *spec.Items
		if prop.Items() != nil {
			items = prop.Items()
		}
		if items == nil {
			items = new(spec.Items)
		}
		if err := rp.parseProperty(gofile, ftpe.Elt, itemsTypable{items}, in); err != nil {
			return err
		}
		prop.CollectionOf(items, "")

	case *ast.StructType:
		// this is an embedded struct, we want to parse this to a schema
		ptb, ok := prop.(responseTypable)
		if !ok && in != "body" {
			return fmt.Errorf("items doesn't support embedded structs")
		}
		schema := ptb.Schema()
		if schema == nil {
			schema = new(spec.Schema)
		}
		scp := newSchemaParser(rp.program)
		if err := scp.parseStructType(gofile, schema, ftpe); err != nil {
			return err
		}
		ptb.SetSchema(schema)
		rp.postDecls = append(rp.postDecls, scp.postDecls...)

	case *ast.SelectorExpr:
		sp := selectorParser{
			program:     rp.program,
			AddPostDecl: func(sd schemaDecl) { rp.postDecls = append(rp.postDecls, sd) },
		}
		return sp.TypeForSelector(gofile, ftpe, prop)

	default:
		return fmt.Errorf("%s is unsupported as parameter", ftpe)
	}
	return nil
}
