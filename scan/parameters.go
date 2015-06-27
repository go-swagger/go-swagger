package scan

import (
	"fmt"
	"go/ast"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-swagger/go-swagger/spec"
	"golang.org/x/tools/go/loader"
)

type paramSetter func(*spec.Parameter, []string) error
type itemsSetter func(*spec.Items, []string) error
type matchingParamSetter func(*regexp.Regexp) paramSetter
type matchingItemsSetter func(*regexp.Regexp) itemsSetter

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

func (pt paramTypable) Items() swaggerTypable {
	if pt.param.In == "body" {
		// get the schema for items on the schema property
		if pt.param.Schema == nil {
			pt.param.Schema = new(spec.Schema)
		}
		if pt.param.Schema.Items == nil {
			pt.param.Schema.Items = new(spec.SchemaOrArray)
		}
		if pt.param.Schema.Items.Schema == nil {
			pt.param.Schema.Items.Schema = new(spec.Schema)
		}
		pt.param.Schema.Typed("array", "")
		return schemaTypable{pt.param.Schema.Items.Schema}
	}

	if pt.param.Items == nil {
		pt.param.Items = new(spec.Items)
	}
	pt.param.Type = "array"
	return itemsTypable{pt.param.Items}
}

func (pt paramTypable) Schema() *spec.Schema {
	if pt.param.In != "body" {
		return nil
	}
	if pt.param.Schema == nil {
		pt.param.Schema = new(spec.Schema)
	}
	return pt.param.Schema
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

func (pt itemsTypable) Schema() *spec.Schema {
	return nil
}

func (pt itemsTypable) Items() swaggerTypable {
	if pt.items.Items == nil {
		pt.items.Items = new(spec.Items)
	}
	pt.items.Type = "array"
	return itemsTypable{pt.items.Items}
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

func newParameterParser(prog *loader.Program) *paramStructParser {
	scp := new(paramStructParser)
	scp.program = prog
	scp.scp = newSchemaParser(prog)
	return scp
}

type paramStructParser struct {
	program   *loader.Program
	postDecls []schemaDecl
	scp       *schemaParser
}

func (pp *paramStructParser) Parse(gofile *ast.File, target interface{}) error {
	tgt := target.(map[string]*spec.Operation)
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

func (pp *paramStructParser) parseDecl(operations map[string]*spec.Operation, decl paramDecl) error {
	// check if there is a swagger:parameters tag that is followed by one or more words,
	// these words are the ids of the operations this parameter struct applies to
	// once type name is found convert it to a schema, by looking up the schema in the
	// parameters dictionary that got passed into this parse method
	for _, opid := range decl.inferOperationIDs() {
		operation, ok := operations[opid]
		if !ok {
			operation = new(spec.Operation)
			operations[opid] = operation
			operation.ID = opid
		}

		// analyze struct body for fields etc
		// each exported struct field:
		// * gets a type mapped to a go primitive
		// * perhaps gets a format
		// * has to document the validations that apply for the type and the field
		// * when the struct field points to a model it becomes a ref: #/definitions/ModelName
		// * comments that aren't tags is used as the description
		if tpe, ok := decl.TypeSpec.Type.(*ast.StructType); ok {
			if err := pp.parseStructType(decl.File, operation, tpe, make(map[string]spec.Parameter)); err != nil {
				return err
			}
		}

		//operations[opid] = operation
	}
	return nil
}

func (pp *paramStructParser) parseEmbeddedStruct(gofile *ast.File, operation *spec.Operation, expr ast.Expr, seenPreviously map[string]spec.Parameter) error {
	switch tpe := expr.(type) {
	case *ast.Ident:
		// do lookup of type
		// take primitives into account, they should result in an error for swagger
		pkg, err := pp.scp.packageForFile(gofile)
		if err != nil {
			return fmt.Errorf("embedded struct: %v", err)
		}
		file, _, ts, err := findSourceFile(pkg, tpe.Name)
		if err != nil {
			return fmt.Errorf("embedded struct: %v", err)
		}
		if st, ok := ts.Type.(*ast.StructType); ok {
			return pp.parseStructType(file, operation, st, seenPreviously)
		}
	case *ast.SelectorExpr:
		// look up package, file and then type
		pkg, err := pp.scp.packageForSelector(gofile, tpe.X)
		if err != nil {
			return fmt.Errorf("embedded struct: %v", err)
		}
		file, _, ts, err := findSourceFile(pkg, tpe.Sel.Name)
		if err != nil {
			return fmt.Errorf("embedded struct: %v", err)
		}
		if st, ok := ts.Type.(*ast.StructType); ok {
			return pp.parseStructType(file, operation, st, seenPreviously)
		}
	}
	return fmt.Errorf("unable to resolve embedded struct for: %v\n", expr)
}

func (pp *paramStructParser) parseStructType(gofile *ast.File, operation *spec.Operation, tpe *ast.StructType, seenPreviously map[string]spec.Parameter) error {
	if tpe.Fields != nil {
		pt := seenPreviously

		for _, fld := range tpe.Fields.List {
			if len(fld.Names) == 0 {
				// when the embedded struct is annotated with swagger:allOf it will be used as allOf property
				// otherwise the fields will just be included as normal properties
				if err := pp.parseEmbeddedStruct(gofile, operation, fld.Type, pt); err != nil {
					return err
				}
			}
		}

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

				ps := pt[nm]
				ps.In = in
				var pty swaggerTypable = paramTypable{&ps}
				if in == "body" {
					pty = schemaTypable{pty.Schema()}
				}
				if err := parseProperty(pp.scp, gofile, fld.Type, pty); err != nil {
					return err
				}

				sp := new(sectionedParser)
				sp.setDescription = func(lines []string) { ps.Description = joinDropLast(lines) }
				if ps.Ref.GetURL() == nil {
					sp.taggers = []tagParser{
						newSingleLineTagParser("maximum", &setMaximum{paramValidations{&ps}, rxf(rxMaximumFmt, "")}),
						newSingleLineTagParser("minimum", &setMinimum{paramValidations{&ps}, rxf(rxMinimumFmt, "")}),
						newSingleLineTagParser("multipleOf", &setMultipleOf{paramValidations{&ps}, rxf(rxMultipleOfFmt, "")}),
						newSingleLineTagParser("minLength", &setMinLength{paramValidations{&ps}, rxf(rxMinLengthFmt, "")}),
						newSingleLineTagParser("maxLength", &setMaxLength{paramValidations{&ps}, rxf(rxMaxLengthFmt, "")}),
						newSingleLineTagParser("pattern", &setPattern{paramValidations{&ps}, rxf(rxPatternFmt, "")}),
						newSingleLineTagParser("collectionFormat", &setCollectionFormat{paramValidations{&ps}, rxf(rxCollectionFormatFmt, "")}),
						newSingleLineTagParser("minItems", &setMinItems{paramValidations{&ps}, rxf(rxMinItemsFmt, "")}),
						newSingleLineTagParser("maxItems", &setMaxItems{paramValidations{&ps}, rxf(rxMaxItemsFmt, "")}),
						newSingleLineTagParser("unique", &setUnique{paramValidations{&ps}, rxf(rxUniqueFmt, "")}),
						newSingleLineTagParser("required", &setRequiredParam{&ps}),
						newSingleLineTagParser("in", &matchOnlyParam{&ps, rxIn}),
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

				} else {

					sp.taggers = []tagParser{
						newSingleLineTagParser("required", &matchOnlyParam{&ps, rxRequired}),
						newSingleLineTagParser("in", &matchOnlyParam{&ps, rxIn}),
					}
				}
				if err := sp.Parse(fld.Doc); err != nil {
					return err
				}

				if ps.Name == "" {
					ps.Name = nm
				}

				if nm != gnm {
					ps.AddExtension("x-go-name", gnm)
				}
				pt[nm] = ps
			}
		}

		for k, p := range pt {
			for i, v := range operation.Parameters {
				if v.Name == k {
					operation.Parameters = append(operation.Parameters[:i], operation.Parameters[i+1:]...)
					break
				}
			}
			operation.Parameters = append(operation.Parameters, p)
		}
	}

	return nil
}
