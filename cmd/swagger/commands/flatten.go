package commands

import (
	"errors"
	"log"

	"github.com/go-openapi/analysis"
	"github.com/go-openapi/loads"
	flags "github.com/jessevdk/go-flags"
)

// FlattenSpec is a command that flattens a swagger document
// Which will expand the remote references in a spec and move inline schemas to definitions
// after flattening there are no complex inlined anymore
type FlattenSpec struct {
	Compact bool           `long:"compact" description:"when present, doesn't prettify the json"`
	Output  flags.Filename `long:"output" short:"o" description:"the file to write to"`
}

// Execute expands the spec
func (c *FlattenSpec) Execute(args []string) error {
	if len(args) == 0 {
		return errors.New("The validate command requires the swagger document url to be specified")
	}

	swaggerDoc := args[0]
	specDoc, err := loads.Spec(swaggerDoc)
	if err != nil {
		log.Fatalln(err)
	}

	if er := analysis.Flatten(analysis.FlattenOpts{
		BasePath: specDoc.SpecFilePath(),
		Spec:     analysis.New(specDoc.Spec()),
	}); er != nil {
		log.Fatalln(er)
	}

	return writeToFile(specDoc.Spec(), !c.Compact, string(c.Output))
}
