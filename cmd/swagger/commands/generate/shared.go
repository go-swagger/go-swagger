package generate

import (
	"log"
	"os"
	"path/filepath"

	"github.com/go-openapi/swag"
	"github.com/go-swagger/go-swagger/generator"
	flags "github.com/jessevdk/go-flags"
	"github.com/spf13/viper"
)

type shared struct {
	Spec                  flags.Filename `long:"spec" short:"f" description:"the spec file to use (default swagger.{json,yml,yaml})"`
	APIPackage            string         `long:"api-package" short:"a" description:"the package to save the operations" default:"operations"`
	ModelPackage          string         `long:"model-package" short:"m" description:"the package to save the models" default:"models"`
	ServerPackage         string         `long:"server-package" short:"s" description:"the package to save the server specific code" default:"restapi"`
	ClientPackage         string         `long:"client-package" short:"c" description:"the package to save the client specific code" default:"client"`
	Target                flags.Filename `long:"target" short:"t" default:"./" description:"the base directory for generating the files"`
	TemplateDir           flags.Filename `long:"template-dir" short:"T" description:"alternative template override directory"`
	ConfigFile            flags.Filename `long:"config-file" short:"C" description:"configuration file to use for overriding template options"`
	CopyrightFile         flags.Filename `long:"copyright-file" short:"r" description:"copyright file used to add copyright header"`
	ExistingModels        string         `long:"existing-models" description:"use pre-generated models e.g. github.com/foobar/model"`
	AdditionalInitialisms []string       `long:"additional-initialism" description:"consecutive capitals that should be considered intialisms"`
	genOpts               *generator.GenOpts
}

type sharedCommand interface {
	getOpts() (*generator.GenOpts, error)
	getConfigFile() flags.Filename
	getAdditionalInitialisms() []string
	generate(*generator.GenOpts) error
	log(string)
}

func (s *shared) getConfigFile() flags.Filename {
	return s.ConfigFile
}

func (s *shared) getAdditionalInitialisms() []string {
	return s.AdditionalInitialisms
}

func createSwagger(s sharedCommand) error {
	cfg, err := readConfig(string(s.getConfigFile()))
	if err != nil {
		return err
	}
	setDebug(cfg)

	opts, err := s.getOpts()
	if err != nil {
		return err
	}

	if err := opts.EnsureDefaults(); err != nil {
		return err
	}

	if err := configureOptsFromConfig(cfg, opts); err != nil {
		return err
	}

	swag.AddInitialisms(s.getAdditionalInitialisms()...)

	if err := s.generate(opts); err != nil {
		return err
	}

	basepath, err := filepath.Abs(".")
	if err != nil {
		return err
	}

	targetAbs, err := filepath.Abs(opts.Target)
	if err != nil {
		return err
	}
	rp, err := filepath.Rel(basepath, targetAbs)
	if err != nil {
		return err
	}

	s.log(rp)

	return nil
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
