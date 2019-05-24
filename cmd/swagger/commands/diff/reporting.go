package diff

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"sort"
	"strings"

	"github.com/go-openapi/spec"
)

// Compare returns the result of analysing breaking and non breaking changes
// between to Swagger specs
func Compare(spec1, spec2 *spec.Swagger) *SpecAnalyser {
	specDiffs := NewSpecDiffs()
	specDiffs.Analyse(spec1, spec2)
	return specDiffs
}

// PathItemOp - combines path and operation into a single keyed entity
type PathItemOp struct {
	ParentPathItem *spec.PathItem  `json:"pathitem"`
	Operation      *spec.Operation `json:"operation"`
}

// URLMethod - combines url and method into a single keyed entity
type URLMethod struct {
	Path   string `json:"path"`
	Method string `json:"method"`
}



// DataDirection indicates the direction of change Request vs Response
type DataDirection int

const (
	// Request Used for messages/param diffs in a request
	Request DataDirection = iota
	// Response Used for messages/param diffs in a response
	Response
)

// SpecDifference encapsulates the details of an individual diff in part of a spec
type SpecDifference struct {
	DifferenceLocation DifferenceLocation `json:"location"`
	Code               SpecChangeCode     `json:"code"`
	Compatability      Compatability      `json:"compatibility"`
	DiffInfo           string             `json:"info,omitempty"`
}


func (sd SpecDifference) String() string {
	optionalMethod := ""
	direction := "Request Param:"
	if len(sd.DifferenceLocation.Method) > 0 {
		optionalMethod = fmt.Sprintf(":%s", sd.DifferenceLocation.Method)
	}
	optionalResponse := ""
	if sd.DifferenceLocation.Response > 0 {
		direction = "Response Body:"
		optionalResponse = fmt.Sprintf("->%d", sd.DifferenceLocation.Response)
	}

	paramOrPropertyLocation := ""
	if sd.DifferenceLocation.Node != nil {
		paramOrPropertyLocation = " - " + sd.DifferenceLocation.Node.String()
	} else {
		direction = ""
	}
	return fmt.Sprintf("%s%s%s - %s %s %s %s", sd.DifferenceLocation.URL, optionalMethod, optionalResponse, direction, sd.Code.Description(), sd.DiffInfo, paramOrPropertyLocation)
}

func findParam(name string, params []spec.Parameter) (spec.Parameter, bool) {
	for _, eachCandidate := range params {
		if eachCandidate.Name == name {
			return eachCandidate, true
		}
	}
	return spec.Parameter{}, false
}

func getParams(pathParams, opParams []spec.Parameter, location string) map[string]spec.Parameter {
	params := map[string]spec.Parameter{}
	// add shared path params
	for _, eachParam := range pathParams {
		if eachParam.In == location {
			params[eachParam.Name] = eachParam
		}
	}
	// add any overridden params
	for _, eachParam := range opParams {
		if eachParam.In == location {
			params[eachParam.Name] = eachParam
		}
	}
	return params
}



func getNameOnlyDiffNode(forLocation string) *Node {
	node := Node{
		Field: forLocation,
	}
	return &node
}

func schemaContext(location string, propertyName string, schema *spec.Schema) string {
	schemaType := ""
	if len(schema.Type) > 0 {
		schemaType = schema.Type[0]
	}
	arraySuffix := ""
	if schemaType == "array" {
		arraySuffix = "[]"
	}
	elName := "obj"
	if len(propertyName) > 0 {
		elName = propertyName
	}
	if len(propertyName) == 0 {
		return fmt.Sprintf("%s.%s%s type:[%s]", strings.Title(location), elName, arraySuffix, schemaType)
	}
	return fmt.Sprintf("%s.%s%s  type:[%s]", strings.Title(location), elName, arraySuffix, schemaType)
}

func getParamDiffNode(paramName string, param spec.Parameter, includeType bool) *Node {

	node := getSchemaDiffNode(paramName, param.Schema)

	return node
}

func getSimpleSchemaDiffNode(name string, schema *spec.SimpleSchema) *Node {
	node := Node{
		Field: name,
	}
	if schema != nil {
		node.TypeName, node.IsArray = getSimpleSchemaType(schema)
	}

	return &node
}

func getSchemaDiffNode(name string, schema *spec.Schema) *Node {
	node := Node{
		Field: name,
	}
	if schema != nil {
		node.TypeName, node.IsArray = getSchemaType(&schema.SchemaProps)
	}

	return &node
}

func definitonFromURL(url *url.URL) string {
	if url == nil {
		return ""
	}
	fragmentParts := strings.Split(url.Fragment, "/")
	numParts := len(fragmentParts)
	if numParts == 0 {
		return ""
	}

	return fragmentParts[numParts-1]

}

func getSimpleSchemaType(schema *spec.SimpleSchema) (typeName string, isArray bool) {
	typeName = schema.Type
	if typeName == "array" {
		typeName, _ = getSimpleSchemaType(&schema.Items.SimpleSchema)
		return typeName, true
	}
	return typeName, false
}

func getSchemaType(schema *spec.SchemaProps) (typeName string, isArray bool) {
	refStr := definitonFromURL(schema.Ref.GetURL())
	if len(refStr) > 0 {
		return refStr, false
	}
	typeName = schema.Type[0]
	if typeName == "array" {
		typeName, _ = getSchemaType(&schema.Items.Schema.SchemaProps)
		return typeName, true
	}
	return typeName, false
}

func paramContext(location string, paramName string, param spec.Parameter, includeType bool) string {
	arraySuffix := ""

	paramType := ""
	if includeType {
		paramType = "type: " + param.Type + " "
		if paramType == "array" {
			arraySuffix = "[]"
			refStr := param.Ref.String()
			paramType = refStr + "::" + param.Items.ItemsTypeName()
		}
	}
	elName := "obj"
	if len(paramName) > 0 {
		elName = paramName
	}
	return fmt.Sprintf("%s.%s%s %s", strings.Title(location), elName, arraySuffix, paramType)
}

func primitiveTypeString(typeName, typeFormat string) string {
	if typeFormat != "" {
		return fmt.Sprintf("%s.%s", typeName, typeFormat)
	}
	return typeName
}

// TypeDiff - describes a primitive type change
type TypeDiff struct {
	Change      SpecChangeCode `json:"change-type,omitempty"`
	Description string         `json:"description,omitempty"`
	FromType    string         `json:"from-type,omitempty"`
	ToType      string         `json:"to-type,omitempty"`
}

// didn't use 'width' so as not to confuse with bit width
var numberWideness = map[string]int{
	"number":        3,
	"number.double": 3,
	"double":        3,
	"number.float":  2,
	"float":         2,
	"long":          1,
	"integer.int64": 1,
	"integer":       0,
	"integer.int32": 0,
}

func typeAndFormat(property *spec.Schema) string {
	if len(property.Type) == 0 {
		return "obj"
	}
	if property.Format != "" {
		return fmt.Sprintf("%s:%s", property.Type[0], property.Format)
	}
	return fmt.Sprintf("%s", property.Type[0])
}

func onlyOneNil(left, right interface{}) (onlyOne bool, isNil bool) {
	leftIsNil := left == nil
	rightIsNil := right == nil

	if (leftIsNil && !rightIsNil) || (!leftIsNil && rightIsNil) {
		return true, false
	}
	return false, leftIsNil // could use either
}

// ReportCompatability lists and spec
func (sd *SpecAnalyser) ReportCompatability() error {
	if sd.BreakingChangeCount > 0 {
		fmt.Printf("\nBREAKING CHANGES:\n=================\n")
		sd.reportChanges(Breaking)
		return fmt.Errorf("Compatability Test FAILED: %d Breaking changes detected", sd.BreakingChangeCount)
	}
	log.Printf("Compatability test OK. No breaking changes identified.")
	return nil
}

func (sd *SpecAnalyser) reportChanges(compat Compatability) {
	toReportList := []string{}

	for _, diff := range sd.Diffs {
		if diff.Compatability == compat {
			toReportList = append(toReportList, diff.String())
		}
	}

	sort.Slice(toReportList, func(i, j int) bool {
		return toReportList[i] < toReportList[j]
	})

	for _, eachDiff := range toReportList {
		fmt.Println(eachDiff)
	}
}

func prettyprint(b []byte) ([]byte, error) {
	var out bytes.Buffer
	err := json.Indent(&out, b, "", "  ")
	return out.Bytes(), err
}

// JSONMarshal allows the item to be correctly rendered to json
func JSONMarshal(t interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(t)
	return buffer.Bytes(), err
}

// ReportAllDiffs lists all the diffs between two specs
func (sd *SpecAnalyser) ReportAllDiffs(fmtJSON bool) error {
	if fmtJSON {

		b, err := JSONMarshal(sd.Diffs)
		if err != nil {
			log.Fatalf("Couldn't print results: %v", err)
		}
		pretty, err := prettyprint(b)
		if err != nil {
			log.Fatalf("Couldn't print results: %v", err)
		}
		fmt.Println(string(pretty))
		return nil
	}
	numDiffs := len(sd.Diffs)
	if numDiffs == 0 {
		fmt.Println("No changes identified")
		return nil
	}

	if numDiffs != sd.BreakingChangeCount {
		fmt.Println("NON-BREAKING CHANGES:\n=====================")
		sd.reportChanges(NonBreaking)
	}

	return sd.ReportCompatability()
}
