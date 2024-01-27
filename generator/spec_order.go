package generator

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/go-openapi/spec"
	"github.com/go-openapi/swag"
	yamlv2 "gopkg.in/yaml.v2"
	yaml "gopkg.in/yaml.v3"
)

// withAutoXOrder reloads the spec to specify property order as they appear
// in the spec (supports yaml documents only).
//
// At this moment, this is exported as it is used independently by the mixin command.
func WithAutoXOrder(specPath string) (string, error) {
	loader := orderedLoader{specPath: specPath}

	return loader.withAutoXOrder()
}

type orderedLoader struct {
	specPath   string
	resultPath string
}

func (l orderedLoader) lookForMapSlice(ele any, key string) (yaml.Node, bool) {
	slice, ok := ele.(yamlv2.MapSlice)
	if !ok {
		return nil, false
	}

	for _, v := range slice {
		if v.Key != key {
			continue
		}

		if slice, ok := v.Value.(yamlv2.MapSlice); ok {
			return slice, ok
		}
	}

	return nil, false
}

func yamlNode(root *yaml.Node) (interface{}, error) {
	switch root.Kind {
	case yaml.DocumentNode:
		return yamlDocument(root)
	case yaml.SequenceNode:
		return yamlSequence(root)
	case yaml.MappingNode:
		return yamlMapping(root)
	case yaml.ScalarNode:
		return yamlScalar(root)
	case yaml.AliasNode:
		return yamlNode(root.Alias)
	default:
		return nil, fmt.Errorf("unsupported YAML node type: %v", root.Kind)
	}
}

func (l orderedLoader) isRemoteRef(key string, value any) (bool, string, error) {
	if key != "$ref" {
		return false, "", nil
	}
	asString, isString := value.(string)
	if !isString {
		return false, "", nil
	}

	ref, err := spec.NewRef(asString)
	if err != nil {
		return false, "", err
	}

	if ref.HasFragmentOnly {
		return false, "", nil
	}

	u := ref.GetURL()
	u.Fragment = ""

	return true, u.String(), nil
}

func (l orderedLoader) rewriteRemoteRef(element any) error {
	return nil // TODO
}

func (l orderedLoader) addXOrder(element any) {
	props, isProperties := l.lookForMapSlice(element, "properties")
	if !isProperties {
		return
	}

	for i, prop := range props {
		log.Printf("DEBUG(fred): property (%[1]T: %+[1]v", prop)
		pSlice, ok := prop.Value.(yamlv2.MapSlice)
		if !ok {
			continue
		}

		isObject := false
		xOrderIndex := -1 // find if x-order already exists

		for i, v := range pSlice {
			if v.Key == "type" && v.Value == object {
				isObject = true
			}
			if v.Key == xOrder {
				xOrderIndex = i
				break
			}
		}

		if xOrderIndex > -1 { // override existing x-order
			pSlice[xOrderIndex] = yamlv2.MapItem{Key: xOrder, Value: i}
		} else { // append new x-order
			pSlice = append(pSlice, yamlv2.MapItem{Key: xOrder, Value: i})
		}

		prop.Value = pSlice
		props[i] = prop

		if isObject {
			l.addXOrder(pSlice)
		}
	}
}

func (l orderedLoader) withAutoXOrder() (string, error) {
	data, err := swag.LoadFromFileOrHTTP(l.specPath)
	if err != nil {
		return "", err
	}

	//yamlDoc, err := bytesToYAMLv2Doc(data)
	yamlDoc, err := swag.BytesToYAMLDoc(data)
	if err != nil {
		return "", err
	}

	if defs, ok := l.lookForMapSlice(yamlDoc, "definitions"); ok {
		for _, def := range defs {
			// TODO: rewrite remote $ref
			log.Printf("DEBUG(fred): definition (%[1]T: %+[1]v", def)
			l.addXOrder(def.Value)
		}
	}

	if params, ok := l.lookForMapSlice(yamlDoc, "parameters"); ok {
		for _, parameter := range params {
			log.Printf("DEBUG(fred):  param (%[1]T: %+[1]v", parameter)
			if in, ok := l.lookForMapSlice(parameter, "in"); ok {
				if asString, isString := in.Value.(string); isString && asString == "body" {
					if schema, ok := l.lookForMapSlice(parameter, "schema"); ok {
					}
				}
				log.Printf("DEBUG(fred): in param (%[1]T: %+[1]v", in)
			}
		}
	}

	if responses, ok := l.lookForMapSlice(yamlDoc, "responses"); ok {
		for _, response := range responses {
			log.Printf("DEBUG(fred): response (%[1]T: %+[1]v", response)
		}
	}

	if pathsItem, ok := l.lookForMapSlice(yamlDoc, "paths"); ok {
		if operation, ok := l.lookForMapSlice(pathsItem, "delete"); ok {
			log.Printf("DEBUG(fred): delete op (%[1]T: %+[1]v", operation)
		}
		if operation, ok := l.lookForMapSlice(pathsItem, "get"); ok {
			log.Printf("DEBUG(fred): get op (%[1]T: %+[1]v", operation)
		}
		if operation, ok := l.lookForMapSlice(pathsItem, "patch"); ok {
			log.Printf("DEBUG(fred): patch op (%[1]T: %+[1]v", operation)
		}
		if operation, ok := l.lookForMapSlice(pathsItem, "post"); ok {
			log.Printf("DEBUG(fred): post op (%[1]T: %+[1]v", operation)
		}
		if operation, ok := l.lookForMapSlice(pathsItem, "put"); ok {
			log.Printf("DEBUG(fred): put op (%[1]T: %+[1]v", operation)
		}
	}

	l.addXOrder(yamlDoc)

	out, err := yamlv2.Marshal(yamlDoc)
	if err != nil {
		return "", err
	}

	tmpDir, err := os.MkdirTemp("", "go-swagger-")
	if err != nil {
		return "", err
	}

	tmpFile := filepath.Join(tmpDir, filepath.Base(l.specPath))
	if err := os.WriteFile(tmpFile, out, 0o600); err != nil {
		return "", err
	}
	return tmpFile, nil
}

// bytesToYAMLDoc converts a byte slice into a YAML document
func bytesToYAMLv2Doc(data []byte) (interface{}, error) {
	var canary map[interface{}]interface{} // validate this is an object and not a different type
	if err := yamlv2.Unmarshal(data, &canary); err != nil {
		return nil, err
	}

	var document yamlv2.MapSlice // preserve order that is present in the document
	if err := yamlv2.Unmarshal(data, &document); err != nil {
		return nil, err
	}
	return document, nil
}
