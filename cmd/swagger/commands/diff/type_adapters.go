package diff

import (
	"github.com/go-openapi/spec"
)

func forItems(items *spec.Items) *spec.Schema {
	if items == nil {
		return nil
	}
	valids := items.CommonValidations
	schema := spec.Schema{
		SchemaProps: spec.SchemaProps{
			Type:             []string{items.SimpleSchema.Type},
			Format:           items.SimpleSchema.Format,
			Maximum:          valids.Maximum,
			ExclusiveMaximum: valids.ExclusiveMaximum,
			Minimum:          valids.Minimum,
			ExclusiveMinimum: valids.ExclusiveMinimum,
			MaxLength:        valids.MaxLength,
			MinLength:        valids.MinLength,
			Pattern:          valids.Pattern,
			MaxItems:         valids.MaxItems,
			MinItems:         valids.MinItems,
			UniqueItems:      valids.UniqueItems,
			MultipleOf:       valids.MultipleOf,
			Enum:             valids.Enum,
		},
	}
	return &schema
}

func forParam(param spec.Parameter) spec.SchemaProps {
	return spec.SchemaProps{
		Type:             []string{param.Type},
		Format:           param.Format,
		Items:            &spec.SchemaOrArray{Schema: forItems(param.Items)},
		Maximum:          param.Maximum,
		ExclusiveMaximum: param.ExclusiveMaximum,
		Minimum:          param.Minimum,
		ExclusiveMinimum: param.ExclusiveMinimum,
		MaxLength:        param.MaxLength,
		MinLength:        param.MinLength,
		Pattern:          param.Pattern,
		MaxItems:         param.MaxItems,
		MinItems:         param.MinItems,
		UniqueItems:      param.UniqueItems,
		MultipleOf:       param.MultipleOf,
		Enum:             param.Enum,
	}
}

// OperationMap saves indexing operations in PathItems individually
type OperationMap map[string]*spec.Operation

func toMap(item *spec.PathItem) OperationMap {
	m := make(OperationMap)

	if item.Post != nil {
		m["post"] = item.Post
	}
	if item.Get != nil {
		m["get"] = item.Get
	}
	if item.Put != nil {
		m["put"] = item.Put
	}
	if item.Patch != nil {
		m["patch"] = item.Patch
	}
	if item.Head != nil {
		m["head"] = item.Head
	}
	if item.Options != nil {
		m["options"] = item.Options
	}
	if item.Delete != nil {
		m["delete"] = item.Delete
	}
	return m
}

func getURLMethodsFor(spec *spec.Swagger) URLMethods {
	returnURLMethods := URLMethods{}

	for url, eachPath := range spec.Paths.Paths {
		opsMap := toMap(&eachPath)
		for method, op := range opsMap {
			returnURLMethods[URLMethod{url, method}] = &PathItemOp{&eachPath, op}
		}
	}
	return returnURLMethods
}

func sliceToStrMap(elements []string) map[string]bool {
	elementMap := make(map[string]bool)
	for _, s := range elements {
		elementMap[s] = true
	}
	return elementMap
}
