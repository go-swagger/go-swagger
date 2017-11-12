package generate

import (
	"log"
	"os"
	"path/filepath"

	"github.com/go-swagger/go-swagger/generator"
	flags "github.com/jessevdk/go-flags"
	"github.com/spf13/viper"
)

type shared struct {
	Spec           flags.Filename `long:"spec" short:"f" description:"the spec file to use (default swagger.{json,yml,yaml})"`
	APIPackage     string         `long:"api-package" short:"a" description:"the package to save the operations" default:"operations"`
	ModelPackage   string         `long:"model-package" short:"m" description:"the package to save the models" default:"models"`
	ServerPackage  string         `long:"server-package" short:"s" description:"the package to save the server specific code" default:"restapi"`
	ClientPackage  string         `long:"client-package" short:"c" description:"the package to save the client specific code" default:"client"`
	Target         flags.Filename `long:"target" short:"t" default:"./" description:"the base directory for generating the files"`
	TemplateDir    flags.Filename `long:"template-dir" short:"T" description:"alternative template override directory"`
	ConfigFile     flags.Filename `long:"config-file" short:"C" description:"configuration file to use for overriding template options"`
	CopyrightFile  flags.Filename `long:"copyright-file" short:"r" description:"copyright file used to add copyright header"`
	ExistingModels string         `long:"existing-models" description:"use pre-generated models e.g. github.com/foobar/model"`
}

func readConfig(filename string) (*viper.Viper, error) {
	if filename == "" {
		return nil, nil
	}

	abspath, err := filepath.Abs(filename)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("trying to read config from", abspath)
	return generator.ReadConfig(abspath)
}

func configureOptsFromConfig(cfg *viper.Viper, opts *generator.GenOpts) error {
	if cfg == nil {
		return nil
	}

	var def generator.LanguageDefinition
	if err := cfg.Unmarshal(&def); err != nil {
		return err
	}
	return def.ConfigureOpts(opts)
}

func setDebug(cfg *viper.Viper) {
	if os.Getenv("DEBUG") != "" || os.Getenv("SWAGGER_DEBUG") != "" {
		if cfg != nil {
			cfg.Debug()
		} else {
			log.Println("NO config read")
		}
	}
}
