package generate

import "github.com/jessevdk/go-flags"

type shared struct {
	Spec         flags.Filename `long:"spec" short:"f" default:"./swagger.json"`
	APIPackage   string         `long:"api-package" short:"a" default:"operations"`
	ModelPackage string         `long:"model-package" short:"m" default:"models"`
	Target       flags.Filename `long:"target" short:"t" default:"./"`
	// TemplateDir  flags.Filename `long:"template-dir"`

}

// All the command to generate an entire application
// both server and client will be generated
type All struct {
}
