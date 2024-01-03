package generator

import (
	"os"
	"path/filepath"

	"github.com/go-openapi/swag"
	yamlv2 "gopkg.in/yaml.v2"
)

// withAutoXOrder reloads the spec to specify property order as they appear
// in the spec (supports yaml documents only).
//
// At this moment, this is exported as it is used independently by the mixin command.
func WithAutoXOrder(specPath string) (string, error) {
	return withAutoXOrder(specPath)
}

func withAutoXOrder(specPath string) (string, error) {
	lookFor := func(ele interface{}, key string) (yamlv2.MapSlice, bool) {
		if slice, ok := ele.(yamlv2.MapSlice); ok {
			for _, v := range slice {
				if v.Key == key {
					if slice, ok := v.Value.(yamlv2.MapSlice); ok {
						return slice, ok
					}
				}
			}
		}
		return nil, false
	}

	var addXOrder func(interface{})
	addXOrder = func(element interface{}) {
		if props, ok := lookFor(element, "properties"); ok {
			for i, prop := range props {
				if pSlice, ok := prop.Value.(yamlv2.MapSlice); ok {
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
						addXOrder(pSlice)
					}
				}
			}
		}
	}

	data, err := swag.LoadFromFileOrHTTP(specPath)
	if err != nil {
		return "", err
	}

	yamlDoc, err := bytesToYAMLv2Doc(data)
	if err != nil {
		return "", err
	}

	if defs, ok := lookFor(yamlDoc, "definitions"); ok {
		for _, def := range defs {
			addXOrder(def.Value)
		}
	}

	addXOrder(yamlDoc)

	out, err := yamlv2.Marshal(yamlDoc)
	if err != nil {
		return "", err
	}

	tmpDir, err := os.MkdirTemp("", "go-swagger-")
	if err != nil {
		return "", err
	}

	tmpFile := filepath.Join(tmpDir, filepath.Base(specPath))
	if err := os.WriteFile(tmpFile, out, 0o600); err != nil {
		panic(err)
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
