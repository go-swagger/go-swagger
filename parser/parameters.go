package parser

import (
	"fmt"
	"go/ast"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/casualjim/go-swagger/spec"
	"golang.org/x/tools/go/loader"
)

type paramSetter func(*spec.Parameter, []string) error
type itemsSetter func(*spec.Items, []string) error
type matchingParamSetter func(*regexp.Regexp) paramSetter
type matchingItemsSetter func(*regexp.Regexp) itemsSetter

type operationTypable interface {
	swaggerTypable
	Items() *spec.Items
	CollectionOf(*spec.Items, string)
}

type operationValidationBuilder interface {
	validationBuilder
	SetCollectionFormat(string)
}

type paramTypable struct {
	param *spec.Parameter
}

func (pt paramTypable) Typed(tpe, format string) {
	pt.param.Typed(tpe, format)
}

func (pt paramTypable) SetRef(ref spec.Ref) {
	pt.param.Ref = ref
}

func (pt paramTypable) Items() *spec.Items {
	return pt.param.Items
}

func (pt paramTypable) Schema() *spec.Schema {
	return pt.param.Schema
}

func (pt paramTypable) SetSchema(schema *spec.Schema) {
	pt.param.Schema = schema
}

func (pt paramTypable) CollectionOf(items *spec.Items, format string) {
	pt.param.CollectionOf(items, format)
}

type itemsTypable struct {
	items *spec.Items
}

func (pt itemsTypable) Typed(tpe, format string) {
	pt.items.Typed(tpe, format)
}

func (pt itemsTypable) SetRef(ref spec.Ref) {
	pt.items.Ref = ref
}

func (pt itemsTypable) Items() *spec.Items {
	return pt.items.Items
}

func (pt itemsTypable) CollectionOf(items *spec.Items, format string) {
	pt.items.CollectionOf(items, format)
}

type paramValidations struct {
	current *spec.Parameter
}

func (sv paramValidations) SetMaximum(val float64, exclusive bool) {
	sv.current.Maximum = &val
	sv.current.ExclusiveMaximum = exclusive
}
func (sv paramValidations) SetMinimum(val float64, exclusive bool) {
	sv.current.Minimum = &val
	sv.current.ExclusiveMinimum = exclusive
}
func (sv paramValidations) SetMultipleOf(val float64)      { sv.current.MultipleOf = &val }
func (sv paramValidations) SetMinItems(val int64)          { sv.current.MinItems = &val }
func (sv paramValidations) SetMaxItems(val int64)          { sv.current.MaxItems = &val }
func (sv paramValidations) SetMinLength(val int64)         { sv.current.MinLength = &val }
func (sv paramValidations) SetMaxLength(val int64)         { sv.current.MaxLength = &val }
func (sv paramValidations) SetPattern(val string)          { sv.current.Pattern = val }
func (sv paramValidations) SetUnique(val bool)             { sv.current.UniqueItems = val }
func (sv paramValidations) SetCollectionFormat(val string) { sv.current.CollectionFormat = val }

type itemsValidations struct {
	current *spec.Items
}

func (sv itemsValidations) SetMaximum(val float64, exclusive bool) {
	sv.current.Maximum = &val
	sv.current.ExclusiveMaximum = exclusive
}
func (sv itemsValidations) SetMinimum(val float64, exclusive bool) {
	sv.current.Minimum = &val
	sv.current.ExclusiveMinimum = exclusive
}
func (sv itemsValidations) SetMultipleOf(val float64)      { sv.current.MultipleOf = &val }
func (sv itemsValidations) SetMinItems(val int64)          { sv.current.MinItems = &val }
func (sv itemsValidations) SetMaxItems(val int64)          { sv.current.MaxItems = &val }
func (sv itemsValidations) SetMinLength(val int64)         { sv.current.MinLength = &val }
func (sv itemsValidations) SetMaxLength(val int64)         { sv.current.MaxLength = &val }
func (sv itemsValidations) SetPattern(val string)          { sv.current.Pattern = val }
func (sv itemsValidations) SetUnique(val bool)             { sv.current.UniqueItems = val }
func (sv itemsValidations) SetCollectionFormat(val string) { sv.current.CollectionFormat = val }

type paramDecl struct {
	File         *ast.File
	Decl         *ast.GenDecl
	TypeSpec     *ast.TypeSpec
	OperationIDs []string
}

func (sd paramDecl) inferOperationIDs() (opids []string) {
	if len(sd.OperationIDs) > 0 {
		opids = sd.OperationIDs
		return
	}

	if sd.Decl.Doc != nil {
	DECLS:
		for _, cmt := range sd.Decl.Doc.List {
			for _, ln := range strings.Split(cmt.Text, "\n") {
				matches := rxParametersOverride.FindStringSubmatch(ln)
				if len(matches) > 1 && len(matches[1]) > 0 {
					for _, pt := range strings.Split(matches[1], " ") {
						tr := strings.TrimSpace(pt)
						if len(tr) > 0 {
							opids = append(opids, tr)
						}
					}
					break DECLS
				}
			}
		}
	}
	sd.OperationIDs = opids
	return
}

func newParamDescription(setter paramSetter) (t *sectionTagger) {
	t = newDescriptionTagger()
	t.set = func(obj interface{}, lines []string) error { return setter(obj.(*spec.Parameter), lines) }
	return
}

func newParamSection(name string, multiLine bool, setter paramSetter) (t *sectionTagger) {
	t = newSectionTagger(name, multiLine)
	t.stripsTag = false
	t.set = func(obj interface{}, lines []string) error { return setter(obj.(*spec.Parameter), lines) }
	return
}

func newParameterFieldSection(name string, matcher *regexp.Regexp, ms matchingParamSetter) (t *sectionTagger) {
	t = newSectionTagger(name, false)
	t.stripsTag = false
	t.matcher = matcher
	setter := ms(matcher)
	t.set = func(obj interface{}, lines []string) error { return setter(obj.(*spec.Parameter), lines) }
	return
}

func newParamValidatorSection(name string, matcher *regexp.Regexp, setter paramSetter) (t *sectionTagger) {
	t = newSectionTagger(name, false)
	t.stripsTag = false
	t.matcher = matcher
	t.set = func(obj interface{}, lines []string) error { return setter(obj.(*spec.Parameter), lines) }
	return
}

func newItemsFieldSection(name string, matcher *regexp.Regexp, ms matchingItemsSetter) (t *sectionTagger) {
	t = newSectionTagger(name, false)
	t.stripsTag = false
	t.matcher = matcher
	setter := ms(matcher)
	t.set = func(obj interface{}, lines []string) error { return setter(obj.(*spec.Items), lines) }
	return
}

func setParamDescription(param *spec.Parameter, lines []string) error {
	param.Description = joinDropLast(lines)
	return nil
}

func setParamMaximum(rx *regexp.Regexp) paramSetter {
	return func(schema *spec.Parameter, lines []string) error {
		bldr := setMaximum{paramValidations{schema}, rx}
		return bldr.Parse(lines)
	}
}

func setParamMinimum(rx *regexp.Regexp) paramSetter {
	return func(schema *spec.Parameter, lines []string) error {
		bldr := setMinimum{paramValidations{schema}, rx}
		return bldr.Parse(lines)
	}
}

func setParamMultipleOf(rx *regexp.Regexp) paramSetter {
	return func(schema *spec.Parameter, lines []string) error {
		bldr := setMultipleOf{paramValidations{schema}, rx}
		return bldr.Parse(lines)
	}
}

func setParamMaxItems(rx *regexp.Regexp) paramSetter {
	return func(schema *spec.Parameter, lines []string) error {
		bldr := setMaxItems{paramValidations{schema}, rx}
		return bldr.Parse(lines)
	}
}

func setParamMinItems(rx *regexp.Regexp) paramSetter {
	return func(schema *spec.Parameter, lines []string) error {
		bldr := setMinItems{paramValidations{schema}, rx}
		return bldr.Parse(lines)
	}
}

func setParamMaxLength(rx *regexp.Regexp) paramSetter {
	return func(schema *spec.Parameter, lines []string) error {
		bldr := setMaxLength{paramValidations{schema}, rx}
		return bldr.Parse(lines)
	}
}

func setParamMinLength(rx *regexp.Regexp) paramSetter {
	return func(schema *spec.Parameter, lines []string) error {
		bldr := setMinLength{paramValidations{schema}, rx}
		return bldr.Parse(lines)

	}
}

func setParamPattern(rx *regexp.Regexp) paramSetter {
	return func(schema *spec.Parameter, lines []string) error {
		bldr := setPattern{paramValidations{schema}, rx}
		return bldr.Parse(lines)
	}
}

func setParamUnique(rx *regexp.Regexp) paramSetter {
	return func(schema *spec.Parameter, lines []string) error {
		bldr := setUnique{paramValidations{schema}, rx}
		return bldr.Parse(lines)
	}
}

func setParamRequired(rx *regexp.Regexp) paramSetter {
	return func(param *spec.Parameter, lines []string) error {
		if len(lines) == 0 || (len(lines) == 1 && len(lines[0]) == 0) {
			return nil
		}
		matches := rx.FindStringSubmatch(lines[0])
		if len(matches) > 1 && len(matches[1]) > 0 {
			req, err := strconv.ParseBool(matches[1])
			if err != nil {
				return err
			}
			param.Required = req
		}
		return nil
	}
}

func setParamIn(rx *regexp.Regexp) paramSetter {
	return func(param *spec.Parameter, lines []string) error {
		if len(lines) == 0 || (len(lines) == 1 && len(lines[0]) == 0) {
			return nil
		}
		matches := rx.FindStringSubmatch(lines[0])
		if len(matches) > 1 && len(matches[1]) > 0 {
			param.In = strings.TrimSpace(matches[1])
		}
		return nil
	}
}

func setParamCollectionFormat(rx *regexp.Regexp) paramSetter {
	return func(schema *spec.Parameter, lines []string) error {
		bldr := setCollectionFormat{paramValidations{schema}, rx}
		return bldr.Parse(lines)
	}
}

func setItemsMaximum(rx *regexp.Regexp) itemsSetter {
	return func(schema *spec.Items, lines []string) error {
		bldr := setMaximum{itemsValidations{schema}, rx}
		return bldr.Parse(lines)
	}
}

func setItemsMinimum(rx *regexp.Regexp) itemsSetter {
	return func(schema *spec.Items, lines []string) error {
		bldr := setMinimum{itemsValidations{schema}, rx}
		return bldr.Parse(lines)
	}
}

func setItemsMultipleOf(rx *regexp.Regexp) itemsSetter {
	return func(schema *spec.Items, lines []string) error {
		bldr := setMultipleOf{itemsValidations{schema}, rx}
		return bldr.Parse(lines)
	}
}

func setItemsMaxItems(rx *regexp.Regexp) itemsSetter {
	return func(schema *spec.Items, lines []string) error {
		bldr := setMaxItems{itemsValidations{schema}, rx}
		return bldr.Parse(lines)
	}
}

func setItemsMinItems(rx *regexp.Regexp) itemsSetter {
	return func(schema *spec.Items, lines []string) error {
		bldr := setMinItems{itemsValidations{schema}, rx}
		return bldr.Parse(lines)
	}
}

func setItemsMaxLength(rx *regexp.Regexp) itemsSetter {
	return func(schema *spec.Items, lines []string) error {
		bldr := setMaxLength{itemsValidations{schema}, rx}
		return bldr.Parse(lines)
	}
}

func setItemsMinLength(rx *regexp.Regexp) itemsSetter {
	return func(schema *spec.Items, lines []string) error {
		bldr := setMinLength{itemsValidations{schema}, rx}
		return bldr.Parse(lines)

	}
}

func setItemsPattern(rx *regexp.Regexp) itemsSetter {
	return func(schema *spec.Items, lines []string) error {
		bldr := setPattern{itemsValidations{schema}, rx}
		return bldr.Parse(lines)
	}
}

func setItemsUnique(rx *regexp.Regexp) itemsSetter {
	return func(schema *spec.Items, lines []string) error {
		bldr := setUnique{itemsValidations{schema}, rx}
		return bldr.Parse(lines)
	}
}

func setItemsCollectionFormat(rx *regexp.Regexp) itemsSetter {
	return func(schema *spec.Items, lines []string) error {
		bldr := setCollectionFormat{itemsValidations{schema}, rx}
		return bldr.Parse(lines)
	}
}

func newParameterParser(prog *loader.Program) *paramStructParser {
	scp := new(paramStructParser)
	scp.program = prog
	scp.header.taggers = []*sectionTagger{newParamDescription(setParamDescription)}
	scp.header.otherTags = []string{"+swagger"}
	return scp
}

type paramStructParser struct {
	taggers []*sectionTagger
	header  struct {
		taggers   []*sectionTagger
		otherTags []string
	}
	program   *loader.Program
	postDecls []schemaDecl
}

func (pp *paramStructParser) Parse(gofile *ast.File, target interface{}) error {
	tgt := target.(map[string]spec.Operation)
	for _, decl := range gofile.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		for _, spc := range gd.Specs {
			if ts, ok := spc.(*ast.TypeSpec); ok {
				sd := paramDecl{gofile, gd, ts, nil}
				sd.inferOperationIDs()
				if err := pp.parseDecl(tgt, sd); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (pp *paramStructParser) parseDecl(operations map[string]spec.Operation, decl paramDecl) error {
	// check if there is a +swagger:parameters tag that is followed by one or more words,
	// these words are the ids of the operations this parameter struct applies to
	// once type name is found convert it to a schema, by looking up the schema in the
	// parameters dictionary that got passed into this parse method
	for _, opid := range decl.inferOperationIDs() {
		operation, ok := operations[opid]
		if !ok {
			operation.ID = opid
		}
		opPtr := &operation

		// analyze struct body for fields etc
		// each exported struct field:
		// * gets a type mapped to a go primitive
		// * perhaps gets a format
		// * has to document the validations that apply for the type and the field
		// * when the struct field points to a model it becomes a ref: #/definitions/ModelName
		// * comments that aren't tags is used as the description
		if tpe, ok := decl.TypeSpec.Type.(*ast.StructType); ok {
			if err := pp.parseStructType(decl.File, opPtr, tpe); err != nil {
				return err
			}
		}

		operations[opid] = operation
	}
	return nil
}

func (pp *paramStructParser) parseStructType(gofile *ast.File, operation *spec.Operation, tpe *ast.StructType) error {
	if tpe.Fields != nil {
		pt := make(map[string]spec.Parameter)

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

				in := "query"
				// scan for param location first, this changes some behavior down the line
				for _, cmt := range fld.Doc.List {
					for _, line := range strings.Split(cmt.Text, "\n") {
						matches := rxIn.FindStringSubmatch(line)
						if len(matches) > 0 && len(strings.TrimSpace(matches[1])) > 0 {
							in = strings.TrimSpace(matches[1])
						}
					}
				}

				ps := pt[nm]
				if err := pp.parseProperty(gofile, fld.Type, paramTypable{&ps}, in); err != nil {
					return err
				}
				ps.In = in

				// check if this is a primitive, if so parse the validations from the
				// doc comments of the slice declaration.
				if ftpe, ok := fld.Type.(*ast.ArrayType); ok {
					if iftpe, ok := ftpe.Elt.(*ast.Ident); ok && iftpe.Obj == nil {
						if ps.Items != nil {
							if err := pp.parseItemsDocComments(gofile, fld, ps.Items); err != nil {
								return err
							}
						}
					}
				}

				var taggers []*sectionTagger
				if ps.Ref.GetURL() == nil {
					// add title and description for property
					// add validations for property
					taggers = []*sectionTagger{
						newParamDescription(setParamDescription),
						newParameterFieldSection("maximum", rxf(rxMaximumFmt, ""), setParamMaximum),
						newParameterFieldSection("minimum", rxf(rxMinimumFmt, ""), setParamMinimum),
						newParameterFieldSection("multipleOf", rxf(rxMultipleOfFmt, ""), setParamMultipleOf),
						newParameterFieldSection("minLength", rxf(rxMinLengthFmt, ""), setParamMinLength),
						newParameterFieldSection("maxLength", rxf(rxMaxLengthFmt, ""), setParamMaxLength),
						newParameterFieldSection("pattern", rxf(rxPatternFmt, ""), setParamPattern),
						newParameterFieldSection("collectionFormat", rxf(rxCollectionFormatFmt, ""), setParamCollectionFormat),
						newParameterFieldSection("minItems", rxf(rxMinItemsFmt, ""), setParamMinItems),
						newParameterFieldSection("maxItems", rxf(rxMaxItemsFmt, ""), setParamMaxItems),
						newParameterFieldSection("unique", rxf(rxUniqueFmt, ""), setParamUnique),
						newParameterFieldSection("required", rxRequired, setParamRequired),
					}
				} else {
					taggers = []*sectionTagger{
						newParamDescription(setParamDescription),
					}
				}
				if err := parseDocComments(fld.Doc, &ps, taggers, nil); err != nil {
					return err
				}

				if ps.Name == "" {
					ps.Name = nm
				}

				if nm != gnm {
					ps.AddExtension("x-go-name", gnm)
				}
				seenProperties[nm] = struct{}{}
				pt[nm] = ps
			}
		}

		for k := range seenProperties {
			if p, ok := pt[k]; ok {
				for i, v := range operation.Parameters {
					if v.Name == k {
						operation.Parameters = append(operation.Parameters[:i], operation.Parameters[i+1:]...)
						break
					}
				}
				operation.Parameters = append(operation.Parameters, p)
			}
		}
	}

	return nil
}

func (pp *paramStructParser) parseItemsDocComments(gofile *ast.File, fld *ast.Field, prop *spec.Items) error {
	// add title and description for property
	// add validations for property
	taggers := []*sectionTagger{
		newItemsFieldSection("maximum", rxf(rxMaximumFmt, rxItemsPrefix), setItemsMaximum),
		newItemsFieldSection("minimum", rxf(rxMinimumFmt, rxItemsPrefix), setItemsMinimum),
		newItemsFieldSection("multipleOf", rxf(rxMultipleOfFmt, rxItemsPrefix), setItemsMultipleOf),
		newItemsFieldSection("minLength", rxf(rxMinLengthFmt, rxItemsPrefix), setItemsMinLength),
		newItemsFieldSection("maxLength", rxf(rxMaxLengthFmt, rxItemsPrefix), setItemsMaxLength),
		newItemsFieldSection("pattern", rxf(rxPatternFmt, rxItemsPrefix), setItemsPattern),
		newItemsFieldSection("minItems", rxf(rxMinItemsFmt, rxItemsPrefix), setItemsMinItems),
		newItemsFieldSection("maxItems", rxf(rxMaxItemsFmt, rxItemsPrefix), setItemsMaxItems),
		newItemsFieldSection("unique", rxf(rxUniqueFmt, rxItemsPrefix), setItemsUnique),
		newItemsFieldSection("collectionFormat", rxf(rxCollectionFormatFmt, rxItemsPrefix), setItemsCollectionFormat),
	}
	return parseDocComments(fld.Doc, prop, taggers, nil)
}

func (pp *paramStructParser) parseProperty(gofile *ast.File, fld ast.Expr, prop operationTypable, in string) error {
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
					prop.SetRef(ref)

					for _, d := range gofile.Decls {
						if gd, ok := d.(*ast.GenDecl); ok {
							for _, tss := range gd.Specs {
								if tss.Pos() == ts.Pos() {
									sd := schemaDecl{gofile, gd, ts, "", ""}
									sd.inferNames()
									pp.postDecls = append(pp.postDecls, sd)
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
		pp.parseProperty(gofile, ftpe.X, prop, in)

	case *ast.ArrayType: // slice type
		if in == "body" {
			var items spec.Schema
			scp := newSchemaParser(pp.program)
			if err := scp.parseProperty(gofile, ftpe.Elt, &items); err != nil {
				return err
			}
			prop.(paramTypable).SetSchema(new(spec.Schema).CollectionOf(items))
			return nil
		}
		var items *spec.Items
		if prop.Items() != nil {
			items = prop.Items()
		}
		if items == nil {
			items = new(spec.Items)
		}
		if err := pp.parseProperty(gofile, ftpe.Elt, itemsTypable{items}, in); err != nil {
			return err
		}
		prop.CollectionOf(items, "")

	case *ast.StructType:
		// this is an embedded struct, we want to parse this to a schema
		ptb, ok := prop.(paramTypable)
		if !ok && in != "body" {
			return fmt.Errorf("items doesn't support embedded structs")
		}
		schema := ptb.Schema()
		if schema == nil {
			schema = new(spec.Schema)
		}
		scp := newSchemaParser(pp.program)
		if err := scp.parseStructType(gofile, schema, ftpe); err != nil {
			return err
		}
		ptb.SetSchema(schema)
		pp.postDecls = append(pp.postDecls, scp.postDecls...)

	case *ast.SelectorExpr:
		sp := selectorParser{
			program:     pp.program,
			AddPostDecl: func(sd schemaDecl) { pp.postDecls = append(pp.postDecls) },
		}
		return sp.TypeForSelector(gofile, ftpe, prop)

	default:
		return fmt.Errorf("%s is unsupported as parameter", ftpe)
	}
	return nil
}
