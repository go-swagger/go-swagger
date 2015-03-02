package generate

import (
	"errors"

	"github.com/casualjim/go-swagger/generator"
)

// Operation the generate operation files command
type Operation struct {
	shared
	Name      []string `long:"name" short:"n" required:"true" description:"the operations to generate, repeat for multiple"`
	Tags      []string `long:"tags" description:"the tags to include, if not specified defaults to all"`
	Principal string   `long:"principal" description:"the model to use for the security principal"`
	NoHandler bool     `long:"skip-handler" description:"when present will not generate an operation handler"`
	NoStruct  bool     `long:"skip-parameters" description:"when present will not generate the parameter model struct"`
	DumpData  bool     `long:"dump-data" description:"when present dumps the json for the template generator instead of generating files"`
}

// Execute generates a model file
func (o *Operation) Execute(args []string) error {
	if o.DumpData && len(o.Name) > 1 {
		return errors.New("only 1 operation at a time is supported for dumping data")
	}
	return generator.GenerateServerOperation(
		o.Name,
		o.Tags,
		!o.NoHandler,
		!o.NoStruct,
		generator.GenOpts{
			Spec:          string(o.Spec),
			Target:        string(o.Target),
			APIPackage:    o.APIPackage,
			ModelPackage:  o.ModelPackage,
			ServerPackage: o.ServerPackage,
			ClientPackage: o.ClientPackage,
			Principal:     o.Principal,
			DumpData:      o.DumpData,
		})
}
