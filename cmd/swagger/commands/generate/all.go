package generate

import (
	"time"

	"github.com/casualjim/go-swagger/generator"
	"github.com/jessevdk/go-flags"
)

type shared struct {
	Spec          flags.Filename `long:"spec" short:"f" description:"the spec file to use" default:"./swagger.json"`
	APIPackage    string         `long:"api-package" short:"a" description:"the package to save the operations" default:"operations"`
	ModelPackage  string         `long:"model-package" short:"m" description:"the package to save the models" default:"models"`
	ServerPackage string         `long:"server-package" short:"s" description:"the package to save the server specific code" default:"server"`
	ClientPackage string         `long:"client-package" short:"c" description:"the package to save the client specific code" default:"client"`
	Target        flags.Filename `long:"target" short:"t" default:"./" description:"the base directory for generating the files"`
	// TemplateDir  flags.Filename `long:"template-dir"`

}

// All the command to generate an entire server application
type All struct {
	shared
	Name           string   `long:"name" short:"A" description:"the name of the application, defaults to a mangled value of info.title"`
	Operations     []string `long:"operation" short:"O" description:"specify an operation to include, repeat for multiple"`
	Tags           []string `long:"tags" description:"the tags to include, if not specified defaults to all"`
	Principal      string   `long:"principal" short:"P" description:"the model to use for the security principal"`
	Models         []string `long:"model" short:"M" description:"specify a model to include, repeat for multiple"`
	SkipModels     bool     `long:"skip-models" description:"no models will be generated when this flag is specified"`
	SkipOperations bool     `long:"skip-operations" description:"no operations will be generated when this flag is specified"`
	SkipSupport    bool     `long:"skip-support" description:"no supporting files will be generated when this flag is specified"`
	IncludeUI      bool     `long:"with-ui" description:"when generating a main package it uses a middleware that also serves a swagger-ui for the swagger json"`
}

// Execute runs this command
func (a *All) Execute(args []string) error {
	opts := generator.GenOpts{
		Spec:         string(a.Spec),
		Target:       string(a.Target),
		APIPackage:   a.APIPackage,
		ModelPackage: a.ModelPackage,
		Principal:    a.Principal,
	}

	if !a.SkipModels && (len(a.Models) > 0 || len(a.Operations) == 0) {
		if err := generator.GenerateModel(a.Models, true, true, opts); err != nil {
			return err
		}
	}

	if !a.SkipOperations && (len(a.Operations) > 0 || len(a.Models) == 0) {
		if err := generator.GenerateServerOperation(a.Operations, a.Tags, true, true, opts); err != nil {
			return err
		}
	}
	<-time.After(1 * time.Second)
	if !a.SkipSupport {
		if err := generator.GenerateSupport(a.Name, a.Models, a.Operations, a.IncludeUI, opts); err != nil {
			return err
		}
	}

	return nil
}
