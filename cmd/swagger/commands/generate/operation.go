package generate

import "github.com/casualjim/go-swagger/generator"

// Operation the generate operation files command
type Operation struct {
	shared
	Name      string   `long:"name" short:"n" required:"true" description:"the operation to generate"`
	Tags      []string `long:"tags" description:"the tags to include, if not specified defaults to all"`
	Principal string   `long:"principal" description:"the model to use for the security principal"`
	NoHandler bool     `long:"skip-handler" description:"when present will not generate an operation handler"`
	NoStruct  bool     `long:"skip-parameters" description:"when present will not generate the parameter model struct"`
}

// Execute generates a model file
func (o *Operation) Execute(args []string) error {
	return generator.GenerateServerOperation(
		o.Name,
		o.Tags,
		!o.NoHandler,
		!o.NoStruct,
		generator.GenOpts{
			Spec:         string(o.Spec),
			Target:       string(o.Target),
			APIPackage:   o.APIPackage,
			ModelPackage: o.ModelPackage,
			Principal:    o.Principal,
		})
}
