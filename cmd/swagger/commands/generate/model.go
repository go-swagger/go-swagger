package generate

import "github.com/casualjim/go-swagger/generator2"

// Model the generate model file command
type Model struct {
	shared
	Name        string `long:"name" short:"n" required:"true" description:"the model to generate"`
	NoValidator bool   `long:"skip-validator" description:"when present will not generate a model validator"`
	NoStruct    bool   `long:"skip-struct" description:"when present will not generate the model struct"`
}

// Execute generates a model file
func (m *Model) Execute(args []string) error {
	return generator2.GenerateModel(
		m.Name,
		!m.NoStruct,
		!m.NoValidator,
		generator2.GenOpts{
			Spec:         string(m.Spec),
			Target:       string(m.Target),
			APIPackage:   m.APIPackage,
			ModelPackage: m.ModelPackage,
		})
}
