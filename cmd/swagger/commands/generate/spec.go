package generate

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/go-swagger/go-swagger/scan"
	"github.com/go-swagger/go-swagger/spec"
	"github.com/jessevdk/go-flags"
)

// SpecFile command to generate a swagger spec from a go application
type SpecFile struct {
	BasePath string         `long:"base-path" short:"b" description:"the base path to use" default:"."`
	Output   flags.Filename `long:"output" short:"o" description:"the file to write to"`
	Input    flags.Filename `long:"input" short:"i" description:"the file to use as input"`
}

// Execute runs this command
func (s *SpecFile) Execute(args []string) error {
	input, err := loadSpec(string(s.Input))
	if err != nil {
		return err
	}

	swspec, err := scan.Application(s.BasePath, input, nil, nil)
	if err != nil {
		return err
	}

	return writeToFile(swspec, string(s.Output))
}

var (
	newLine = []byte("\n")
)

func loadSpec(input string) (*spec.Swagger, error) {
	if fi, err := os.Stat(input); err == nil {
		if fi.IsDir() {
			return nil, fmt.Errorf("expected %q to be a file not a directory", input)
		}
		sp, err := spec.Load(input)
		if err != nil {
			return nil, err
		}
		return sp.Spec(), nil
	}
	return nil, nil
}

func writeToFile(swspec *spec.Swagger, output string) error {
	b, err := json.Marshal(swspec)
	if err != nil {
		return err
	}
	if output == "" {
		fmt.Println(string(b))
		return nil
	}
	return ioutil.WriteFile(output, b, 0644)
}
