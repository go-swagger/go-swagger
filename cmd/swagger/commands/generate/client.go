package generate

import "github.com/casualjim/go-swagger/generator"

// Client the command to generate a swagger client
type Client struct {
	shared
	Name           string   `long:"name" short:"A" description:"the name of the application, defaults to a mangled value of info.title"`
	Operations     []string `long:"operation" short:"O" description:"specify an operation to include, repeat for multiple"`
	Tags           []string `long:"tags" description:"the tags to include, if not specified defaults to all"`
	Principal      string   `long:"principal" short:"P" description:"the model to use for the security principal"`
	Models         []string `long:"model" short:"M" description:"specify a model to include, repeat for multiple"`
	SkipModels     bool     `long:"skip-models" description:"no models will be generated when this flag is specified"`
	SkipOperations bool     `long:"skip-operations" description:"no operations will be generated when this flag is specified"`
}

// Execute runs this command
func (c *Client) Execute(args []string) error {
	opts := generator.GenOpts{
		Spec:          string(c.Spec),
		Target:        string(c.Target),
		APIPackage:    c.APIPackage,
		ModelPackage:  c.ModelPackage,
		ServerPackage: c.ServerPackage,
		ClientPackage: c.ClientPackage,
		Principal:     c.Principal,
	}

	if !c.SkipModels && (len(c.Models) > 0 || len(c.Operations) == 0) {
		if err := generator.GenerateModel(c.Models, true, true, opts); err != nil {
			return err
		}
	}

	if !c.SkipOperations && (len(c.Operations) > 0 || len(c.Models) == 0) {
		if err := generator.GenerateClient(c.Name, c.Models, c.Operations, opts); err != nil {
			return err
		}
	}

	return nil
}
