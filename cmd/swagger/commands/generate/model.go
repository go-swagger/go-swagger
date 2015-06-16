package generate

import (
	"errors"

	"github.com/go-swagger/go-swagger/generator"
)

// Model the generate model file command
type Model struct {
	shared
	Name        []string `long:"name" short:"n" required:"true" description:"the model to generate"`
	NoValidator bool     `long:"skip-validator" description:"when present will not generate a model validator"`
	NoStruct    bool     `long:"skip-struct" description:"when present will not generate the model struct"`
	DumpData    bool     `long:"dump-data" description:"when present dumps the json for the template generator instead of generating files"`
}

// Execute generates a model file
func (m *Model) Execute(args []string) error {
	if m.DumpData && len(m.Name) > 1 {
		return errors.New("only 1 model at a time is supported for dumping data")
	}
	return generator.GenerateDefinition(
		m.Name,
		!m.NoStruct,
		!m.NoValidator,
		generator.GenOpts{
			Spec:          string(m.Spec),
			Target:        string(m.Target),
			APIPackage:    m.APIPackage,
			ModelPackage:  m.ModelPackage,
			ServerPackage: m.ServerPackage,
			ClientPackage: m.ClientPackage,
			DumpData:      m.DumpData,
		})
}
