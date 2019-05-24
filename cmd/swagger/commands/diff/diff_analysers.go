package diff

import (
	"fmt"
	"strings"

	"github.com/go-openapi/spec"
)

// URLMethodResponse encapsulates these three elements to act as a map key
type URLMethodResponse struct {
	Path     string `json:"path"`
	Method   string `json:"method"`
	Response string `json:"response"`
}

// MarshalText - for serializing as a map key
func (p URLMethod) MarshalText() (text []byte, err error) {
	return []byte(fmt.Sprintf("%s %s", p.Path, p.Method)), nil
}

// URLMethods allows iteration of endpoints based on url and method
type URLMethods map[URLMethod]*PathItemOp

// SpecDiff contains all the differences for a Spec
type SpecDiff struct {
	BreakingChangeCount int
	Diffs               []SpecDifference
	urlMethods1         URLMethods
	urlMethods2         URLMethods
	Definitions1        spec.Definitions
	Definitions2        spec.Definitions
}

// NewSpecDiffs returns an empty SpecDiffs
func NewSpecDiffs() *SpecDiff {
	return &SpecDiff{
		Diffs: []SpecDifference{},
	}
}

// Analyse the differences in two specs
func (sd *SpecDiff) Analyse(spec1, spec2 *spec.Swagger) error {

	sd.Definitions1 = spec1.Definitions
	sd.Definitions2 = spec2.Definitions
	sd.urlMethods1 = getURLMethodsFor(spec1)
	sd.urlMethods2 = getURLMethodsFor(spec2)

	sd.analyseSpecMetadata(spec1, spec2)
	sd.analyseEndpoints()
	sd.analyseParams()
	sd.analyseResponseParams()

	return nil
}

func (sd *SpecDiff) analyseSpecMetadata(spec1, spec2 *spec.Swagger) {
	// breaking if it no longer consumes any formats
	added, deleted, _ := FromStringArray(spec1.Consumes).DiffsTo(spec1.Consumes)

	for _, eachAdded := range added {
		sd.addDiff(SpecDifference{DifferenceLocation: DifferenceLocation{URL: "consumes"}, Code: AddedConsumesFormat, Compatability: NonBreaking, DiffInfo: eachAdded})
	}
	for _, eachDeleted := range deleted {
		sd.addDiff(SpecDifference{DifferenceLocation: DifferenceLocation{URL: "consumes"}, Code: DeletedConsumesFormat, Compatability: Breaking, DiffInfo: eachDeleted})
	}

	// // breaking if it no longer produces any formats
	added, deleted, _ = FromStringArray(spec1.Produces).DiffsTo(spec1.Produces)

	for _, eachAdded := range added {
		sd.addDiff(SpecDifference{DifferenceLocation: DifferenceLocation{URL: "produces"}, Code: AddedProducesFormat, Compatability: NonBreaking, DiffInfo: eachAdded})
	}
	for _, eachDeleted := range deleted {
		sd.addDiff(SpecDifference{DifferenceLocation: DifferenceLocation{URL: "produces"}, Code: DeletedProducesFormat, Compatability: Breaking, DiffInfo: eachDeleted})
	}

	// // breaking if it no longer supports a scheme
	added, deleted, _ = FromStringArray(spec1.Schemes).DiffsTo(spec1.Schemes)

	for _, eachAdded := range added {
		sd.addDiff(SpecDifference{DifferenceLocation: DifferenceLocation{URL: "schemes"}, Code: AddedSchemes, Compatability: NonBreaking, DiffInfo: eachAdded})
	}
	for _, eachDeleted := range deleted {
		sd.addDiff(SpecDifference{DifferenceLocation: DifferenceLocation{URL: "schemes"}, Code: DeletedSchemes, Compatability: Breaking, DiffInfo: eachDeleted})
	}

	// // host should be able to change without any issues?
	sd.diffStringMetaData(spec1.Host, spec2.Host, ChangedHostURL, Breaking)
	// sd.Host = compareStrings(spec1.Host, spec2.Host)

	// // Base Path change will break non generated clients
	sd.diffStringMetaData(spec1.BasePath, spec2.BasePath, ChangedBasePath, Breaking)

	// TODO: what to do about security?
	// Missing security scheme will break a client
	// Security            []map[string][]string  `json:"security,omitempty"`
	// Tags                []Tag                  `json:"tags,omitempty"`
	// ExternalDocs        *ExternalDocumentation `json:"externalDocs,omitempty"`
}

func (sd *SpecDiff) analyseParams() {
	locations := []string{"query", "path", "body", "header"}

	for _, paramLocation := range locations {
		rootNode := getNameOnlyDiffNode(strings.Title(paramLocation))
		for URLMethod, op2 := range sd.urlMethods2 {
			if op1, ok := sd.urlMethods1[URLMethod]; ok {

				params1 := getParams(op1.ParentPathItem.Parameters, op1.Operation.Parameters, paramLocation)
				params2 := getParams(op2.ParentPathItem.Parameters, op2.Operation.Parameters, paramLocation)

				location := DifferenceLocation{URL: URLMethod.Path, Method: URLMethod.Method, Node: rootNode}
				// detect deleted params
				for paramName1, param1 := range params1 {
					if _, ok := params2[paramName1]; !ok {
						childLocation := location.AddNode(getSchemaDiffNode(paramName1, param1.Schema))
						code := DeletedOptionalParam
						if param1.Required{
							code = DeletedRequiredParam
						}
						sd.addDiff(SpecDifference{DifferenceLocation: childLocation, Code: code})
					}
				}

				// detect added changed params
				for paramName2, param2 := range params2 {
					//changed?
					if param1, ok := params1[paramName2]; ok {
						sd.compareParams(URLMethod, paramLocation, paramName2, param1, param2)
					} else {
						// Added
						childLocation := location.AddNode(getSchemaDiffNode(paramName2, param2.Schema))
						code := AddedOptionalParam
						if param2.Required {
							code = AddedRequiredParam
						}
						sd.addDiff(SpecDifference{DifferenceLocation: childLocation, Code: code})
					}
				}
			}
		}
	}
}

func (sd *SpecDiff) analyseEndpoints() {
	sd.findDeletedEndpoints()
	sd.findAddedEndpoints()
}

func (sd *SpecDiff) analyseResponseParams() {

	for URLMethod2, eachOp2 := range sd.urlMethods2 {
		if op1, ok := sd.urlMethods1[URLMethod2]; ok {
			op1Responses := op1.Operation.Responses.StatusCodeResponses
			op2Responses := eachOp2.Operation.Responses.StatusCodeResponses

			// deleted responses
			for code1 := range op1Responses {
				if _, ok := op2Responses[code1]; !ok {
					location := DifferenceLocation{URL: URLMethod2.Path, Method: URLMethod2.Method, Response: code1}
					sd.addDiff(SpecDifference{DifferenceLocation: location, Code: DeletedResponse})
				}
			}
			// Added updated Response Codes
			for code2, response2 := range op2Responses {

				if op1Response, ok := op1Responses[code2]; ok {
					op1Headers := op1Response.ResponseProps.Headers
					for op2HeaderName, op2Header := range response2.ResponseProps.Headers {

						node := getNameOnlyDiffNode("Headers")
						location := DifferenceLocation{URL: URLMethod2.Path, Method: URLMethod2.Method, Response: code2, Node: node}

						if op1Header, ok := op1Headers[op2HeaderName]; ok {
							childLocation := location.AddNode(getNameOnlyDiffNode(op2HeaderName))
							sd.compareSimpleSchema(childLocation, &op1Header.SimpleSchema, &op2Header.SimpleSchema, false, false)
						} else {
							childLocation := location.AddNode(getNameOnlyDiffNode(op2HeaderName))
							sd.addDiff(SpecDifference{DifferenceLocation: childLocation, Code: AddedResponseHeader})
						}
					}
					for op1HeaderName := range op1Response.ResponseProps.Headers {
						if _, ok := response2.ResponseProps.Headers[op1HeaderName]; !ok {
							node := getNameOnlyDiffNode("Headers")
							location := DifferenceLocation{URL: URLMethod2.Path, Method: URLMethod2.Method, Response: code2, Node: node}
							childLocation := location.AddNode(getNameOnlyDiffNode(op1HeaderName))
							sd.addDiff(SpecDifference{DifferenceLocation: childLocation, Code: DeletedResponseHeader})
						}
					}

					if op1Response.Schema != nil {
						diffLocation := DifferenceLocation{URL: URLMethod2.Path, Method: URLMethod2.Method, Response: code2}
						sd.compareSchema(diffLocation, op1Response.Schema, response2.Schema, true, true)
					}
				} else {
					location := DifferenceLocation{URL: URLMethod2.Path, Method: URLMethod2.Method, Response: code2}
					sd.addDiff(SpecDifference{DifferenceLocation: location, Code: AddedResponse})
				}
			}
		}
	}
}

func addTypeDiff(diffs []TypeDiff, diff TypeDiff) []TypeDiff {
	if diff.Change != NoChangeDetected {
		diffs = append(diffs, diff)
	}
	return diffs
}

// CompareTypes computes type specific property diffs
func (sd *SpecDiff) CompareTypes(type1, type2 spec.SchemaProps) []TypeDiff {

	diffs := []TypeDiff{}

	type1IsPrimitive := len(type1.Type) > 0
	type2IsPrimitive := len(type2.Type) > 0

	// Primitive to Obj or Obj to Primitive
	if type1IsPrimitive && !type2IsPrimitive {
		return addTypeDiff(diffs, TypeDiff{Change: ChangedType, FromType: type1.Type[0], ToType: "obj"})
	}

	if !type1IsPrimitive && type2IsPrimitive {
		return addTypeDiff(diffs, TypeDiff{Change: ChangedType, FromType: type2.Type[0], ToType: "obj"})
	}

	// Single to Array or Array to Single
	type1Array := type1.Type[0] == "array"
	type2Array := type2.Type[0] == "array"

	if type1Array && !type2Array {
		return addTypeDiff(diffs, TypeDiff{Change: ChangedType, FromType: "obj", ToType: type2.Type[0]})
	}

	if !type1Array && type2Array {
		return addTypeDiff(diffs, TypeDiff{Change: ChangedType, FromType: type1.Type[0], ToType: "array"})
	}

	// check type hierarchy change eg string -> integer = NarrowedChange
	//Type
	//Format
	if type1.Type[0] != type2.Type[0] ||
		type1.Format != type2.Format {
		diff := getTypeHierarchyChange(primitiveTypeString(type1.Type[0], type1.Format), primitiveTypeString(type2.Type[0], type2.Format))
		diffs = addTypeDiff(diffs, diff)
	}

	// string changes
	if type1.Type[0] == "string" &&
		type2.Type[0] == "string" {
		//Pattern
		diffs = addTypeDiff(diffs, compareIntValues("MinLength", type1.MinLength, type2.MinLength, NarrowedType, WidenedType))
		diffs = addTypeDiff(diffs, compareIntValues("MaxLength", type1.MinLength, type2.MinLength, WidenedType, NarrowedType))
		if type1.Pattern != type2.Pattern {
			diffs = addTypeDiff(diffs, TypeDiff{Change: ChangedType, Description: fmt.Sprintf("Pattern Changed:%s->%s", type1.Pattern, type2.Pattern)})
		}
		if type1.Type[0] == "string" {
			if len(type1.Enum) > 0 {
				enumDiffs := sd.compareEnums(type1.Enum, type2.Enum)
				for _, eachDiff := range enumDiffs {
					diffs = append(diffs, eachDiff)
				}
			}
		}
	}

	if type1.Type[0] == "array" &&
		type2.Type[0] == "array" {
		// array
		// TODO: Items??
		diffs = addTypeDiff(diffs, compareIntValues("MaxItems", type1.MaxItems, type2.MaxItems, WidenedType, NarrowedType))
		diffs = addTypeDiff(diffs, compareIntValues("MinItems", type1.MinItems, type2.MinItems, NarrowedType, WidenedType))

	}
	// Number
	_, type1IsNumeric := numberWideness[type1.Type[0]]
	_, type2IsNumeric := numberWideness[type2.Type[0]]

	if type1IsNumeric && type2IsNumeric {
		diffs = addTypeDiff(diffs, compareFloatValues("Maximum", type1.Maximum, type2.Maximum, WidenedType, NarrowedType))
		diffs = addTypeDiff(diffs, compareFloatValues("Minimum", type1.Minimum, type2.Minimum, NarrowedType, WidenedType))
		if type1.ExclusiveMaximum && !type2.ExclusiveMaximum {
			diffs = addTypeDiff(diffs, TypeDiff{Change: WidenedType, Description: fmt.Sprintf("Exclusive Maximum Removed:%v->%v", type1.ExclusiveMaximum, type2.ExclusiveMaximum)})
		}
		if !type1.ExclusiveMaximum && type2.ExclusiveMaximum {
			diffs = addTypeDiff(diffs, TypeDiff{Change: NarrowedType, Description: fmt.Sprintf("Exclusive Maximum Added:%v->%v", type1.ExclusiveMaximum, type2.ExclusiveMaximum)})
		}
		if type1.ExclusiveMinimum && !type2.ExclusiveMinimum {
			diffs = addTypeDiff(diffs, TypeDiff{Change: WidenedType, Description: fmt.Sprintf("Exclusive Minimum Removed:%v->%v", type1.ExclusiveMaximum, type2.ExclusiveMaximum)})
		}
		if !type1.ExclusiveMinimum && type2.ExclusiveMinimum {
			diffs = addTypeDiff(diffs, TypeDiff{Change: NarrowedType, Description: fmt.Sprintf("Exclusive Minimum Added:%v->%v", type1.ExclusiveMaximum, type2.ExclusiveMaximum)})
		}
	}
	return diffs
}

func compareFloatValues(fieldName string, val1 *float64, val2 *float64, ifGreaterCode SpecChangeCode, ifLessCode SpecChangeCode) TypeDiff {
	if val1 != nil && val2 != nil {
		if *val1 > *val2 {
			return TypeDiff{Change: ifGreaterCode, Description: ""}
		}
		if *val1 < *val2 {
			return TypeDiff{Change: ifLessCode, Description: ""}
		}
	}
	return TypeDiff{Change: NoChangeDetected, Description: ""}
}

func compareIntValues(fieldName string, val1 *int64, val2 *int64, ifGreaterCode SpecChangeCode, ifLessCode SpecChangeCode) TypeDiff {
	if val1 != nil && val2 != nil {
		if *val1 > *val2 {
			return TypeDiff{Change: ifGreaterCode, Description: ""}
		}
		if *val1 < *val2 {
			return TypeDiff{Change: ifLessCode, Description: ""}
		}

	}
	return TypeDiff{Change: NoChangeDetected, Description: ""}
}

func (sd *SpecDiff) compareParams(urlMethod URLMethod, location string, name string, param1, param2 spec.Parameter) {
	diffLocation := DifferenceLocation{URL: urlMethod.Path, Method: urlMethod.Method}

	if param1.Schema != nil && param2.Schema != nil {
		childLocation := diffLocation.AddNode(getNameOnlyDiffNode(strings.Title(location)))
		childLocation = childLocation.AddNode(getSchemaDiffNode(name, param2.Schema))
		sd.compareSchema(childLocation, param1.Schema, param2.Schema, param1.Required, param2.Required)
	}
	diffs := sd.CompareTypes(forParam(param1), forParam(param2))

	childLocation := diffLocation.AddNode(getNameOnlyDiffNode(strings.Title(location)))
	childLocation = childLocation.AddNode(getSchemaDiffNode(name, param2.Schema))
	for _, eachDiff := range diffs {
		sd.addDiff(SpecDifference{
			DifferenceLocation: childLocation,
			Code:               eachDiff.Change,
			DiffInfo:           eachDiff.Description})
	}
	if param1.Required != param2.Required {
		code := ChangedRequiredToOptionalParam
		if param2.Required {
			code = ChangedOptionalToRequiredParam
		}
		sd.addDiff(SpecDifference{DifferenceLocation: childLocation, Code: code})
	}
}

// TODO - not much here yet
func (sd *SpecDiff) compareHeaders(op1Header, op2Header spec.Header) bool {
	return false
}

func (sd *SpecDiff) compareSimpleSchema(location DifferenceLocation, schema1, schema2 *spec.SimpleSchema, required1, required2 bool) {
	if schema1 == nil || schema2 == nil {
		return
	}

	if schema1.Type == "array" {
		refSchema1 := schema1.Items.SimpleSchema
		refSchema2 := schema2.Items.SimpleSchema

		childLocation := location.AddNode(getSimpleSchemaDiffNode("", schema1))
		sd.compareSimpleSchema(childLocation, &refSchema1, &refSchema2, required1, required2)
		return
	}
	if required1 != required2 {
		code := AddedRequiredProperty
		if required1 {
			code = ChangedRequiredToOptional

		}
		sd.addDiff(SpecDifference{DifferenceLocation: location, Code: code})
	}

}

func (sd *SpecDiff) compareSchema(location DifferenceLocation, schema1, schema2 *spec.Schema, required1, required2 bool) {

	if schema1 == nil || schema2 == nil {
		return
	}

	if len(schema1.Type) == 0 {
		refSchema1, definition1 := sd.schemaFromRef(schema1, &sd.Definitions1)
		refSchema2, definition2 := sd.schemaFromRef(schema2, &sd.Definitions2)
		info := fmt.Sprintf("[%s -> %s]", definition1, definition2)

		if definition1 != definition2 {
			sd.addDiff(SpecDifference{DifferenceLocation: location,
				Code:     ChangedType,
				DiffInfo: info,
			})
		}
		sd.compareSchema(location, refSchema1, refSchema2, required1, required2)
		return
	}
	diffs := sd.CompareTypes(schema1.SchemaProps, schema2.SchemaProps)

	for _, eachTypeDiff := range diffs {
		if eachTypeDiff.Change != NoChangeDetected {
			sd.addDiff(SpecDifference{DifferenceLocation: location, Code: eachTypeDiff.Change, DiffInfo: eachTypeDiff.Description})
		}
	}
	if schema1.Type[0] == "array" {
		refSchema1, _ := sd.schemaFromRef(schema1.Items.Schema, &sd.Definitions1)
		refSchema2, _ := sd.schemaFromRef(schema2.Items.Schema, &sd.Definitions2)

		childLocation := location.AddNode(getSchemaDiffNode("", schema1))
		sd.compareSchema(childLocation, refSchema1, refSchema2, required1, required2)
		return
	}
	if required1 != required2 {
		code := AddedRequiredProperty
		if required1 {
			code = ChangedRequiredToOptional

		}
		sd.addDiff(SpecDifference{DifferenceLocation: location, Code: code})
	}
	requiredProps2 := sliceToStrMap(schema2.Required)
	requiredProps1 := sliceToStrMap(schema1.Required)
	schema1Props := sd.propertiesFor(schema1, &sd.Definitions1)
	schema2Props := sd.propertiesFor(schema2, &sd.Definitions2)
	// find deleted and changed properties
	for eachProp1Name, eachProp1 := range schema1Props {
		_, required1 := requiredProps1[eachProp1Name]
		_, required2 := requiredProps2[eachProp1Name]
		childLoc := sd.addChildDiffNode(location, eachProp1Name, &eachProp1)

		//		thisContext := schemaContext(currentPath, eachProp1Name, &eachProp1)
		if eachProp2, ok := schema2Props[eachProp1Name]; ok {
			sd.compareSchema(childLoc, &eachProp1, &eachProp2, required1, required2)
		} else {
			sd.addDiff(SpecDifference{DifferenceLocation: childLoc, Code: DeletedOptionalParam})
		}
	}

	// find added properties
	for eachProp2Name, eachProp2 := range schema2.Properties {
		if _, ok := schema1.Properties[eachProp2Name]; !ok {
			childLoc := sd.addChildDiffNode(location, eachProp2Name, &eachProp2)
			_, required2 := requiredProps2[eachProp2Name]
			code := AddedProperty
			if required2 {
				code = AddedRequiredProperty
			}
			sd.addDiff(SpecDifference{DifferenceLocation: childLoc, Code: code})
		}
	}
}
func (sd *SpecDiff) addChildDiffNode(location DifferenceLocation, propName string, propSchema *spec.Schema) DifferenceLocation {
	newLoc := location
	if newLoc.Node != nil {
		newLoc.Node = newLoc.Node.Copy()
	}

	childNode := sd.fromSchemaProps(propName, &propSchema.SchemaProps)
	if newLoc.Node != nil {
		newLoc.Node.AddLeafNode(&childNode)
	} else {
		newLoc.Node = &childNode
	}
	return newLoc
}

func (sd *SpecDiff) fromSchemaProps(fieldName string, props *spec.SchemaProps) Node {
	node := Node{}
	node.IsArray = props.Type[0] == "array"
	if !node.IsArray {
		node.TypeName = props.Type[0]
	} //else {
	// if len(props.Items.Schema.Type) == 0 { // reference
	// 	// refSchema := sd.schemaFromRef(props.Items.Schema)
	// 	node.TypeName = "REF_NOT_RESOLVED:" + props.Items.Schema.Ref.GetURL().String()
	// } else {
	// 	node.TypeName = props.Items.Schema.Type[0]
	// }
	//}
	node.Field = fieldName
	return node
}

func (sd *SpecDiff) compareEnums(left, right []interface{}) []TypeDiff {
	diffs := []TypeDiff{}

	leftStrs := []string{}
	rightStrs := []string{}
	for _, eachLeft := range left {
		leftStrs = append(leftStrs, fmt.Sprintf("%v", eachLeft))
	}
	for _, eachRight := range right {
		rightStrs = append(rightStrs, fmt.Sprintf("%v", eachRight))
	}
	added, deleted, _ := FromStringArray(leftStrs).DiffsTo(rightStrs)
	if len(added) > 0 {
		typeChange := "<" + strings.Join(added, ",") + ">"
		diffs = append(diffs, TypeDiff{Change: AddedEnumValue, Description: typeChange})
	}
	if len(deleted) > 0 {
		typeChange := "<" + strings.Join(deleted, ",") + ">"
		diffs = append(diffs, TypeDiff{Change: DeletedEnumValue, Description: typeChange})
	}
	return diffs
}

func isNumericType(typeName string) (wideness int, isNumeric bool) {
	wideness, isNumeric = numberWideness[typeName]
	return
}

func isStringType(typeName string) bool {
	return typeName == "string" || typeName == "password"
}

func getTypeHierarchyChange(type1, type2 string) TypeDiff {
	if type1 == type2 {
		return TypeDiff{Change: NoChangeDetected, Description: ""}
	}
	diffDescription := fmt.Sprintf("%s -> %s", type1, type2)
	if isStringType(type1) && !isStringType(type2) {
		return TypeDiff{Change: NarrowedType, Description: diffDescription}
	}
	if !isStringType(type1) && isStringType(type2) {
		return TypeDiff{Change: WidenedType, Description: diffDescription}
	}
	type1Wideness, type1IsNumeric := numberWideness[type1]
	type2Wideness, type2IsNumeric := numberWideness[type2]
	if type1IsNumeric && type2IsNumeric {
		if type1Wideness == type2Wideness {
			return TypeDiff{Change: ChangedToCompatibleType, Description: diffDescription}
		}
		if type1Wideness > type2Wideness {
			return TypeDiff{Change: NarrowedType, Description: diffDescription}
		}
		if type1Wideness < type2Wideness {
			return TypeDiff{Change: WidenedType, Description: diffDescription}
		}
	}
	return TypeDiff{Change: ChangedType, Description: diffDescription}
}

func (sd *SpecDiff) findAddedEndpoints() {
	for URLMethod := range sd.urlMethods2 {
		if _, ok := sd.urlMethods1[URLMethod]; !ok {
			sd.addDiff(SpecDifference{DifferenceLocation: DifferenceLocation{URL: URLMethod.Path, Method: URLMethod.Method}, Code: AddedEndpoint})
		}
	}
}

func (sd *SpecDiff) findDeletedEndpoints() {
	for URLMethod := range sd.urlMethods1 {
		if _, ok := sd.urlMethods2[URLMethod]; !ok {
			sd.addDiff(SpecDifference{DifferenceLocation: DifferenceLocation{URL: URLMethod.Path, Method: URLMethod.Method}, Code: DeletedEndpoint})
		}
	}
}

func (sd *SpecDiff) diffStringMetaData(item1, item2 string, codeIfDiff SpecChangeCode, compatIfDiff Compatability) {
	if item1 != item2 {
		diffSpec := fmt.Sprintf("%s -> %s", item1, item2)
		sd.addDiff(SpecDifference{DifferenceLocation: DifferenceLocation{URL: ""}, Code: codeIfDiff, Compatability: compatIfDiff, DiffInfo: diffSpec})
	}
}

func (sd *SpecDiff) addDiff(diff SpecDifference) {
	context := Request
	if diff.DifferenceLocation.Response > 0 {
		context = Response
	}
	diff.Compatability = getCompatabilityForChange(diff.Code, context)

	if diff.Compatability == Breaking {
		sd.BreakingChangeCount++
	}
	sd.Diffs = append(sd.Diffs, diff)
}
