package spec

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/casualjim/go-swagger/jsonpointer"
	"github.com/casualjim/go-swagger/util"
)

// ResolutionCache a cache for resolving urls
type ResolutionCache interface {
	Get(string) (interface{}, bool)
	Set(string, interface{})
}

type simpleCache struct {
	store map[string]interface{}
}

func defaultResolutionCache() ResolutionCache {
	return &simpleCache{store: map[string]interface{}{
		"http://swagger.io/v2/schema.json":       swaggerSchema,
		"http://json-schema.org/draft-04/schema": jsonSchema,
	}}
}

func (s *simpleCache) Get(uri string) (interface{}, bool) {
	v, ok := s.store[uri]
	return v, ok
}

func (s *simpleCache) Set(uri string, data interface{}) {
	s.store[uri] = data
}

type schemaLoader struct {
	loadingRef  *Ref
	startingRef *Ref
	currentRef  *Ref
	root        interface{}
	cache       ResolutionCache
	loadDoc     DocLoader
	schemaRef   *Ref
}

var idPtr, _ = jsonpointer.New("/id")
var schemaPtr, _ = jsonpointer.New("/$schema")
var refPtr, _ = jsonpointer.New("/$ref")

func defaultSchemaLoader(root interface{}, ref *Ref) (*schemaLoader, error) {
	startingRef := ref

	if startingRef != nil {
		idRef, err := idFromNode(root)
		if err != nil {
			return nil, err
		}
		if idRef != nil {
			startingRef, err = ref.Inherits(*idRef)
			if err != nil {
				return nil, err
			}
		}
	}
	var ptr *jsonpointer.Pointer
	if ref != nil {
		ptr = ref.GetPointer()
	}

	currentRef := nextRef(root, startingRef, ptr)
	return &schemaLoader{
		root:        root,
		loadingRef:  ref,
		startingRef: startingRef,
		cache:       defaultResolutionCache(),
		loadDoc:     util.JSONDoc,
		currentRef:  currentRef,
	}, nil
}

func idFromNode(node interface{}) (*Ref, error) {
	if idValue, _, err := idPtr.Get(node); err == nil {
		if refStr, ok := idValue.(string); ok && refStr != "" {
			idRef, err := NewRef(refStr)
			if err != nil {
				return nil, err
			}
			return &idRef, nil
		}
	}
	return nil, nil
}

func nextRef(startingNode interface{}, startingRef *Ref, ptr *jsonpointer.Pointer) *Ref {
	if startingRef == nil {
		return nil
	}
	if ptr == nil {
		return startingRef
	}
	ret := startingRef
	var idRef *Ref
	node := startingNode

	for _, tok := range ptr.DecodedTokens() {
		node, _, _ = jsonpointer.GetForToken(node, tok)
		if node == nil {
			break
		}

		idRef, _ = idFromNode(node)
		if idRef != nil {
			nw, err := ret.Inherits(*idRef)
			if err != nil {
				break
			}
			ret = nw
		}

		refRef, _, _ := refPtr.Get(node)
		if refRef != nil {
			rf, _ := NewRef(refRef.(string))
			nw, err := ret.Inherits(rf)
			if err != nil {
				break
			}
			ret = nw
		}

	}
	return ret
}

func (r *schemaLoader) resolveRef(currentRef, ref *Ref, node, target interface{}) error {
	tgt := reflect.ValueOf(target)
	if tgt.Kind() != reflect.Ptr {
		return fmt.Errorf("resolve ref: target needs to be a pointer")
	}

	fmt.Printf("[before] current ref: %s, ref: %s for target %T\n", currentRef, ref, target)
	// pretty.Println(node)

	oldRef := currentRef
	if currentRef != nil {
		var err error
		currentRef, err = currentRef.Inherits(*nextRef(node, ref, currentRef.GetPointer()))
		if err != nil {
			return err
		}
	}
	if currentRef == nil {
		currentRef = ref
	}

	refURL := currentRef.GetURL()
	if refURL == nil {
		fmt.Println("ref url was nil")
		return nil
	}
	if refURL.String() == "" {
		nv := reflect.ValueOf(node)
		reflect.Indirect(tgt).Set(reflect.Indirect(nv))
		return nil
	}

	if strings.HasPrefix(refURL.String(), "#") {
		fmt.Println("local ref so getting data and bailing")
		// pretty.Println(node)
		res, _, err := ref.GetPointer().Get(node)
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Printf("[result] %+v\n", res)
		reflect.Indirect(tgt).Set(reflect.Indirect(reflect.ValueOf(res)))
		return nil
	}

	if refURL.Scheme != "" && refURL.Host != "" {
		// most definitely take the red pill
		toFetch := *refURL
		toFetch.Fragment = ""
		fmt.Println("loading:", toFetch.String())

		data, fromCache := r.cache.Get(toFetch.String())
		if !fromCache {
			fmt.Println("fetching:", toFetch.String())
			b, err := r.loadDoc(toFetch.String())
			if err != nil {
				return err
			}

			if err := json.Unmarshal(b, &data); err != nil {
				return err
			}
			r.cache.Set(toFetch.String(), data)
		}

		if oldRef != currentRef && oldRef != ref {
			return r.resolveRef(currentRef, ref, data, target)
		}

		var err error
		var res interface{}
		if currentRef.String() != "" {
			res, _, err = currentRef.GetPointer().Get(data)
			if err != nil {
				return err
			}
		} else {
			res = data
		}

		// fmt.Printf("[result] %+v\n", res)

		bb, err := json.Marshal(res)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(bb, target); err != nil {
			return err
		}
		if !fromCache {
			r.cache.Set(toFetch.String(), target)
		}

	}
	return nil
}

func (r *schemaLoader) Resolve(ref *Ref, target interface{}) error {
	if err := r.resolveRef(r.currentRef, ref, r.root, target); err != nil {
		return err
	}

	return nil
}

type specExpander struct {
	spec     *Swagger
	resolver *schemaLoader
}

func expandSpec(spec *Swagger) error {

	resolver, err := defaultSchemaLoader(spec, nil)
	if err != nil {
		return err
	}
	for key, defintition := range spec.Definitions {
		if err := expandSchema(&defintition, resolver); err != nil {
			return err
		}
		spec.Definitions[key] = defintition
	}

	for key, parameter := range spec.Parameters {
		if err := expandParameter(&parameter, resolver); err != nil {
			return err
		}
		spec.Parameters[key] = parameter
	}

	for key, response := range spec.Responses {
		if err := expandResponse(&response, resolver); err != nil {
			return err
		}
		spec.Responses[key] = response
	}

	if spec.Paths != nil {
		for key, path := range spec.Paths.Paths {
			if err := expandPathItem(&path, resolver); err != nil {
				return err
			}
			spec.Paths.Paths[key] = path
		}
	}

	return nil
}

// ExpandSchema expands the refs in the schema object
func ExpandSchema(schema *Schema, root interface{}) error {
	if schema == nil {
		return nil
	}
	if root == nil {
		root = schema
	}

	resolver, err := defaultSchemaLoader(root, nil)
	if err != nil {
		return err
	}

	return expandSchema(schema, resolver)
}

func expandSchema(schema *Schema, resolver *schemaLoader) error {
	if schema == nil {
		return nil
	}

	// create a schema expander and run that
	if schema.Ref.String() != "" || schema.Ref.IsRoot() {
		var newSchema Schema
		if err := resolver.Resolve(&schema.Ref, &newSchema); err != nil {
			return err
		}
		*schema = newSchema
		return nil
	}

	if schema.Items != nil {
		if schema.Items.Schema != nil {
			if err := expandSchema(schema.Items.Schema, resolver); err != nil {
				return err
			}
		}
		for i := range schema.Items.Schemas {
			sch := &(schema.Items.Schemas[i])
			if err := expandSchema(sch, resolver); err != nil {
				return err
			}
		}
	}
	for i := range schema.AllOf {
		sch := &(schema.AllOf[i])
		if err := expandSchema(sch, resolver); err != nil {
			return err
		}
	}
	for i := range schema.AnyOf {
		sch := &(schema.AnyOf[i])
		if err := expandSchema(sch, resolver); err != nil {
			return err
		}
	}
	for i := range schema.OneOf {
		sch := &(schema.OneOf[i])
		if err := expandSchema(sch, resolver); err != nil {
			return err
		}
	}
	if schema.Not != nil {
		if err := expandSchema(schema.Not, resolver); err != nil {
			return err
		}
	}
	for k, v := range schema.Properties {
		if err := expandSchema(&v, resolver); err != nil {
			return err
		}
		schema.Properties[k] = v
	}
	if schema.AdditionalProperties != nil && schema.AdditionalProperties.Schema != nil {
		if err := expandSchema(schema.AdditionalProperties.Schema, resolver); err != nil {
			return err
		}
	}
	for k, v := range schema.PatternProperties {
		if err := expandSchema(&v, resolver); err != nil {
			return err
		}
		schema.PatternProperties[k] = v
	}
	for k, v := range schema.Dependencies {
		if v.Schema != nil {
			if err := expandSchema(v.Schema, resolver); err != nil {
				return err
			}
			schema.Dependencies[k] = v
		}
	}
	if schema.AdditionalItems != nil && schema.AdditionalItems.Schema != nil {
		if err := expandSchema(schema.AdditionalItems.Schema, resolver); err != nil {
			return err
		}
	}
	for k, v := range schema.Definitions {
		if err := expandSchema(&v, resolver); err != nil {
			return err
		}
		schema.Definitions[k] = v
	}
	return nil
}

func expandPathItem(pathItem *PathItem, resolver *schemaLoader) error {
	if pathItem == nil {
		return nil
	}
	if err := resolver.Resolve(&pathItem.Ref, &pathItem); err != nil {
		return err
	}

	if err := expandOperation(pathItem.Get, resolver); err != nil {
		return err
	}
	if err := expandOperation(pathItem.Head, resolver); err != nil {
		return err
	}
	if err := expandOperation(pathItem.Options, resolver); err != nil {
		return err
	}
	if err := expandOperation(pathItem.Put, resolver); err != nil {
		return err
	}
	if err := expandOperation(pathItem.Post, resolver); err != nil {
		return err
	}
	if err := expandOperation(pathItem.Patch, resolver); err != nil {
		return err
	}
	if err := expandOperation(pathItem.Delete, resolver); err != nil {
		return err
	}
	return nil
}

func expandOperation(op *Operation, resolver *schemaLoader) error {
	if op == nil {
		return nil
	}
	for i, param := range op.Parameters {
		if err := expandParameter(&param, resolver); err != nil {
			return err
		}
		op.Parameters[i] = param
	}

	if op.Responses != nil {
		responses := op.Responses
		if err := expandResponse(responses.Default, resolver); err != nil {
			return err
		}
		for code, response := range responses.StatusCodeResponses {
			if err := expandResponse(&response, resolver); err != nil {
				return err
			}
			responses.StatusCodeResponses[code] = response
		}
	}
	return nil
}

func expandResponse(response *Response, resolver *schemaLoader) error {
	if response == nil {
		return nil
	}
	if err := resolver.Resolve(&response.Ref, response); err != nil {
		return err
	}

	if response.Schema != nil {
		if err := expandSchema(response.Schema, resolver); err != nil {
			return err
		}
	}
	return nil
}

func expandParameter(parameter *Parameter, resolver *schemaLoader) error {
	if parameter == nil {
		return nil
	}
	if err := resolver.Resolve(&parameter.Ref, parameter); err != nil {
		return err
	}
	if parameter.Schema != nil {
		if err := expandSchema(parameter.Schema, resolver); err != nil {
			return err
		}

	}
	return nil
}
