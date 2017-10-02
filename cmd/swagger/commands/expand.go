package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"
	flags "github.com/jessevdk/go-flags"
)

// ExpandSpec is a command that expands the $refs in a swagger document
type ExpandSpec struct {
	Compact bool           `long:"compact" description:"when present, doesn't prettify the json"`
	Output  flags.Filename `long:"output" short:"o" description:"the file to write to"`
}

// Execute expands the spec
func (c *ExpandSpec) Execute(args []string) error {
	if len(args) == 0 {
		return errors.New("The validate command requires the swagger document url to be specified")
	}

	swaggerDoc := args[0]
	specDoc, err := loads.Spec(swaggerDoc)
	if err != nil {
		return err
	}

	exp, err := specDoc.Expanded()
	if err != nil {
		return err
	}

	return writeToFile(exp.Spec(), !c.Compact, string(c.Output))
}

func writeToFile(swspec *spec.Swagger, pretty bool, output string) error {
	var b []byte
	var err error
	if pretty {
		b, err = json.MarshalIndent(swspec, "", "  ")
	} else {
		b, err = json.Marshal(swspec)
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
