package diff

import (
	"fmt"
	"strings"

	"github.com/go-openapi/spec"
)

func getRef(item interface{}) spec.Ref {
	switch s := item.(type) {
	case *spec.Refable:
		return s.Ref
	case *spec.Schema:
		return s.Ref
	case *spec.SchemaProps:
		return s.Ref
	default:
		return spec.Ref{}
	}
}

// CheckToFromArrayType check for changes to or from an Array type
func CheckToFromArrayType(diffs []TypeDiff, type1, type2 interface{}) []TypeDiff {
	// Single to Array or Array to Single
	typString1, isArray1 := getSchemaType(type1)
	typString2, isArray2 := getSchemaType(type2)

	if isArray1 != isArray2 {
		return addTypeDiff(diffs, TypeDiff{Change: ChangedType, FromType: formatTypeString(typString1, isArray1), ToType: formatTypeString(typString2, isArray2)})
	}

	return diffs
}

// CheckToFromPrimitiveType check for diff to or from a primitive
func CheckToFromPrimitiveType(diffs []TypeDiff, type1, type2 interface{}) []TypeDiff {

	type1IsPrimitive := isPrimitive(type1)
	type2IsPrimitive := isPrimitive(type2)

	// Primitive to Obj or Obj to Primitive
	if type1IsPrimitive != type2IsPrimitive {
		typeStr1, isarray1 := getSchemaType(type1)
		typeStr2, isarray2 := getSchemaType(type2)
		return addTypeDiff(diffs, TypeDiff{Change: ChangedType, FromType: formatTypeString(typeStr1, isarray1), ToType: formatTypeString(typeStr2, isarray2)})
	}

	return diffs
}

// CheckRefChange has the property ref changed
func CheckRefChange(diffs []TypeDiff, type1, type2 interface{}) (diffReturn []TypeDiff) {

	diffReturn = diffs
	if isRefType(type1) && isRefType(type2) {
		// both refs but to different objects (TODO detect renamed object)
		ref1 := definitionFromRef(getRef(type1))
		ref2 := definitionFromRef(getRef(type2))
		if ref1 != ref2 {
			diffReturn = addTypeDiff(diffReturn, TypeDiff{Change: RefTargetChanged, FromType: getSchemaTypeStr(type1), ToType: getSchemaTypeStr(type2)})
		}
	} else {
		if isRefType(type1) != isRefType(type2) {
			diffReturn = addTypeDiff(diffReturn, TypeDiff{Change: ChangedType, FromType: getSchemaTypeStr(type1), ToType: getSchemaTypeStr(type2)})
		}
	}
	return
}

func compareEnums(left, right []interface{}) []TypeDiff {
	diffs := []TypeDiff{}

	leftStrs := []string{}
	rightStrs := []string{}
	for _, eachLeft := range left {
		leftStrs = append(leftStrs, fmt.Sprintf("%v", eachLeft))
	}
	for _, eachRight := range right {
		rightStrs = append(rightStrs, fmt.Sprintf("%v", eachRight))
	}
	added, deleted, _ := fromStringArray(leftStrs).DiffsTo(rightStrs)
	if len(added) > 0 {
		typeChange := strings.Join(added, ",")
		diffs = append(diffs, TypeDiff{Change: AddedEnumValue, Description: typeChange})
	}
	if len(deleted) > 0 {
		typeChange := strings.Join(deleted, ",")
		diffs = append(diffs, TypeDiff{Change: DeletedEnumValue, Description: typeChange})
	}

	return diffs
}

// checkNumericTypeChanges checks for changes to or from a numeric type
func checkNumericTypeChanges(diffs []TypeDiff, type1, type2 *spec.SchemaProps) []TypeDiff {
	// Number
	_, type1IsNumeric := numberWideness[type1.Type[0]]
	_, type2IsNumeric := numberWideness[type2.Type[0]]

	if type1IsNumeric && type2IsNumeric {
		foundDiff := false
		if type1.ExclusiveMaximum && !type2.ExclusiveMaximum {
			diffs = addTypeDiff(diffs, TypeDiff{Change: WidenedType, Description: fmt.Sprintf("Exclusive Maximum Removed:%v->%v", type1.ExclusiveMaximum, type2.ExclusiveMaximum)})
			foundDiff = true
		}
		if !type1.ExclusiveMaximum && type2.ExclusiveMaximum {
			diffs = addTypeDiff(diffs, TypeDiff{Change: NarrowedType, Description: fmt.Sprintf("Exclusive Maximum Added:%v->%v", type1.ExclusiveMaximum, type2.ExclusiveMaximum)})
			foundDiff = true
		}
		if type1.ExclusiveMinimum && !type2.ExclusiveMinimum {
			diffs = addTypeDiff(diffs, TypeDiff{Change: WidenedType, Description: fmt.Sprintf("Exclusive Minimum Removed:%v->%v", type1.ExclusiveMaximum, type2.ExclusiveMaximum)})
			foundDiff = true
		}
		if !type1.ExclusiveMinimum && type2.ExclusiveMinimum {
			diffs = addTypeDiff(diffs, TypeDiff{Change: NarrowedType, Description: fmt.Sprintf("Exclusive Minimum Added:%v->%v", type1.ExclusiveMinimum, type2.ExclusiveMinimum)})
			foundDiff = true
		}
		if !foundDiff {
			diffs = addTypeDiff(diffs, compareFloatValues("Maximum", type1.Maximum, type2.Maximum, WidenedType, NarrowedType))
			diffs = addTypeDiff(diffs, compareFloatValues("Minimum", type1.Minimum, type2.Minimum, NarrowedType, WidenedType))
		}
	}
	return diffs
}

// CheckStringTypeChanges checks for changes to or from a string type
func CheckStringTypeChanges(diffs []TypeDiff, type1, type2 *spec.SchemaProps) []TypeDiff {
	// string changes
	if type1.Type[0] == StringType &&
		type2.Type[0] == StringType {
		diffs = addTypeDiff(diffs, compareIntValues("MinLength", type1.MinLength, type2.MinLength, NarrowedType, WidenedType))
		diffs = addTypeDiff(diffs, compareIntValues("MaxLength", type1.MinLength, type2.MinLength, WidenedType, NarrowedType))
		if type1.Pattern != type2.Pattern {
			diffs = addTypeDiff(diffs, TypeDiff{Change: ChangedType, Description: fmt.Sprintf("Pattern Changed:%s->%s", type1.Pattern, type2.Pattern)})
		}
		if type1.Type[0] == StringType {
			if len(type1.Enum) > 0 {
				enumDiffs := compareEnums(type1.Enum, type2.Enum)
				diffs = append(diffs, enumDiffs...)
			}
		}
	}
	return diffs
}

// CheckToFromRequired checks for changes to or from a required property
func CheckToFromRequired(required1, required2 bool) (diffs []TypeDiff) {
	if required1 != required2 {
		code := AddedRequiredProperty
		if required1 {
			code = ChangedRequiredToOptional
		}
		diffs = addTypeDiff(diffs, TypeDiff{Change: code})
	}
	return diffs
}

func compareProperties(location DifferenceLocation, schema1 *spec.Schema, schema2 *spec.Schema, getRefFn1 SchemaFromRefFn, getRefFn2 SchemaFromRefFn, cmp CompareSchemaFn) []SpecDifference {
	propDiffs := []SpecDifference{}

	requiredProps2 := sliceToStrMap(schema2.Required)
	requiredProps1 := sliceToStrMap(schema1.Required)
	schema1Props := propertiesFor(schema1, getRefFn1)
	schema2Props := propertiesFor(schema2, getRefFn2)
	// find deleted and changed properties

	for eachProp1Name, eachProp1 := range schema1Props {
		eachProp1 := eachProp1
		_, required1 := requiredProps1[eachProp1Name]
		_, required2 := requiredProps2[eachProp1Name]
		childLoc := addChildDiffNode(location, eachProp1Name, &eachProp1)

		if eachProp2, ok := schema2Props[eachProp1Name]; ok {
			diffs := CheckToFromRequired(required1, required2)
			if len(diffs) > 0 {
				for _, diff := range diffs {
					propDiffs = append(propDiffs, SpecDifference{DifferenceLocation: childLoc, Code: diff.Change})
				}
			}
			cmp(childLoc, &eachProp1, &eachProp2)
		} else {
			propDiffs = append(propDiffs, SpecDifference{DifferenceLocation: childLoc, Code: DeletedProperty})
		}
	}

	// find added properties
	for eachProp2Name, eachProp2 := range schema2.Properties {
		if _, ok := schema1.Properties[eachProp2Name]; !ok {
			childLoc := addChildDiffNode(location, eachProp2Name, &eachProp2)
			propDiffs = append(propDiffs, SpecDifference{DifferenceLocation: childLoc, Code: AddedProperty})
		}
	}
	return propDiffs

}

// SchemaFromRefFn define this to get a schema for a ref
type SchemaFromRefFn func(spec.Ref) (*spec.Schema, string)

func propertiesFor(schema *spec.Schema, getRefFn SchemaFromRefFn) PropertyMap {
	schemaFromRef, _ := getRefFn(schema.Ref)
	if schemaFromRef != nil {
		schema = schemaFromRef
	}
	props := PropertyMap{}

	if schema.Properties != nil {
		for name, prop := range schema.Properties {
			props[name] = prop
		}
	}
	for _, eachAllOf := range schema.AllOf {
		eachAllOf := eachAllOf
		eachAllOfActual, _ := getRefFn(eachAllOf.SchemaProps.Ref)
		if eachAllOfActual == nil {
			eachAllOfActual = &eachAllOf
		}
		for name, prop := range eachAllOfActual.Properties {
			props[name] = prop
		}
	}
	return props
}

func isRefType(item interface{}) bool {
	switch s := item.(type) {
	case spec.Refable:
		return s.Ref.String() != ""
	case *spec.Schema:
		return s.Ref.String() != ""
	case *spec.SchemaProps:
		return s.Ref.String() != ""
	case *spec.SimpleSchema:
		return false
	default:
		return false
	}
}
