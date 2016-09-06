package generator

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// LanguageDefinition in the configuration file.
type LanguageDefinition struct {
	// Language       string            `mapstructure:"language"`
	// ReservedWords  []string          `mapstructure:"reserved_words"`
	// FormatScript   string            `mapstructure:"format_script"`
	// DefaultImports map[string]string `mapstructure:"default_imports"`
	Layout SectionOpts `mapstructure:"layout"`
}

// ConfigureOpts for generation
func (d *LanguageDefinition) ConfigureOpts(opts *GenOpts) error {
	// var lopts *LanguageOpts
	// if d.Language == "go" {
	// 	lopts = GoLangOpts()
	// }
	// if lopts == nil {
	// 	lopts = new(LanguageOpts)
	// }

	// lopts.ReservedWords = append(lopts.ReservedWords, d.ReservedWords...)

	// lopts.initialized = false
	// lopts.Init()

	opts.Sections = d.Layout
	opts.LanguageOpts = GoLangOpts()
	return nil
}

// LanguageConfig structure that is obtained from parsing a config file
type LanguageConfig map[string]LanguageDefinition

// ReadConfig at the specified path, when no path is specified it will look into
// the current directory and load a .swagger.{yml,json,hcl,toml,properties} file
// Returns a viper config or an error
func ReadConfig(fpath string) (*viper.Viper, error) {
	v := viper.New()
	if fpath != "" {
		if !fileExists(fpath, "") {
			return nil, fmt.Errorf("can't find file for %q", fpath)
		}
		file, err := os.Open(fpath)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		if err := v.ReadConfig(file); err != nil {
			return nil, err
		}
		return v, nil
	}

	v.SetConfigName(".swagger")
	v.AddConfigPath(".")
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.UnsupportedConfigError); !ok && v.ConfigFileUsed() != "" {
			return nil, err
		}
	}
	return v, nil
}
