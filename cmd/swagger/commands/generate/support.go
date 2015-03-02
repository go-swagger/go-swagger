package generate

import "github.com/casualjim/go-swagger/generator"

// Support generates the supporting files
type Support struct {
	shared
	Name       string   `long:"name" short:"A" description:"the name of the application, defaults to a mangled value of info.title"`
	Operations []string `long:"operation" short:"O" description:"specify an operation to include, repeat for multiple"`
	// Tags       []string `long:"tags" description:"the tags to include, if not specified defaults to all"`
	Principal string   `long:"principal" description:"the model to use for the security principal"`
	Models    []string `long:"model" short:"M" description:"specify a model to include, repeat for multiple"`
	DumpData  bool     `long:"dump-data" description:"when present dumps the json for the template generator instead of generating files"`
	IncludeUI bool     `long:"with-ui" description:"when generating a main package it uses a middleware that also serves a swagger-ui for the swagger json"`
}

// Execute generates the supporting files file
func (s *Support) Execute(args []string) error {
	return generator.GenerateSupport(
		s.Name,
		nil,
		nil,
		s.IncludeUI,
		generator.GenOpts{
			Spec:         string(s.Spec),
			Target:       string(s.Target),
			APIPackage:   s.APIPackage,
			ModelPackage: s.ModelPackage,
			Principal:    s.Principal,
			DumpData:     s.DumpData,
		})
}
