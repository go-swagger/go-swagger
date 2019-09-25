package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/go-openapi/swag"
	"github.com/iancoleman/orderedmap"
	"github.com/jessevdk/go-flags"
	"gopkg.in/yaml.v2"
)

const xOrder = "x-order" // sort order for properties (or any schema)

// OrderSpec is a command that flattens a swagger document
// which will expand the remote references in a spec and move inline schemas to definitions
// after flattening there are no complex inlined anymore
type OrderSpec struct {
	Compact bool           `long:"compact" description:"applies to JSON formatted specs. When present, doesn't prettify the json"`
	Output  flags.Filename `long:"output" short:"o" description:"the file to write to"`
	Format  string         `long:"format" description:"the format for the spec document" default:"json" choice:"yaml" choice:"json"`
}

// Execute flattens the spec
func (c *OrderSpec) Execute(args []string) error {
	if len(args) != 1 {
		return errors.New("flatten command requires the single swagger document url to be specified")
	}

	swaggerDoc := args[0]

	doc, err := OrderByXOrder(swaggerDoc)
	if err != nil {
		return err
	}

	return writeOrderedSpecToFile(doc, !c.Compact, c.Format, string(c.Output))
}

func OrderByXOrder(specPath string) (*orderedmap.OrderedMap, error) {
	var convertToOrderedOutput func(ele interface{}) *orderedmap.OrderedMap
	convertToOrderedOutput = func(ele interface{}) *orderedmap.OrderedMap {
		o := orderedmap.New()
		if slice, ok := ele.(yaml.MapSlice); ok {
			for _, v := range slice {
				if slice, ok := v.Value.(yaml.MapSlice); ok {
					o.Set(v.Key.(string), convertToOrderedOutput(slice))
				} else if items, ok := v.Value.([]interface{}); ok {
					newItems := []interface{}{}
					for _, item := range items {
						if slice, ok := item.(yaml.MapSlice); ok {
							newItems = append(newItems, convertToOrderedOutput(slice))
						} else {
							newItems = append(newItems, item)
						}
					}
					o.Set(v.Key.(string), newItems)
				} else {
					o.Set(v.Key.(string), v.Value)
				}
			}
		}
		o.Sort(func(a *orderedmap.Pair, b *orderedmap.Pair) bool {
			return getXOrder(a.Value()) < getXOrder(b.Value())
		})
		return o
	}

	yamlDoc, err := swag.YAMLData(specPath)
	if err != nil {
		panic(err)
	}

	return convertToOrderedOutput(yamlDoc), nil
}

func getXOrder(val interface{}) int {
	if prop, ok := val.(*orderedmap.OrderedMap); ok {
		if pSlice, ok := prop.Get(xOrder); ok {
			return pSlice.(int)
		}
	}
	return 0
}

func writeOrderedSpecToFile(swspec *orderedmap.OrderedMap, pretty bool, format string, output string) error {
	var b []byte
	var err error
	asJSON := format == "json"

	if pretty && asJSON {
		b, err = json.MarshalIndent(swspec, "", "  ")
	} else if asJSON {
		b, err = json.Marshal(swspec)
	} else {
		// marshals as YAML
		b, err = json.Marshal(swspec)
		if err == nil {
			d, ery := swag.BytesToYAMLDoc(b)
			if ery != nil {
				return ery
			}
			b, err = yaml.Marshal(d)
		}
	}
	if err != nil {
		return err
	}
	if output == "" {
		fmt.Println(string(b))
		return nil
	}
	return ioutil.WriteFile(output, b, 0644)
}
