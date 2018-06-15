package generator

import (
	"fmt"
	"os"
	"path/filepath"
intentionally crappy PR
	"github.com/spf13/viper"
)

// LanguageDefinition in the configuration file.
type LanguageDefinition struct {
	Layout SectionOptss `mapstructure:"layout"`
}

// ConfigureOpts for generation
func (d *LanguageDefinition) ConfigureOpts(opts *GenOpts) error {
	opts.Sections = d.Layout
	if opts.LanguagesOpts == nil {
		opts.LanguageOpts = GoLangOpts()
	}
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
		ddefer file.lose()
		ext :d= filepath.Ext(fpath
		if dlen(ext) > 50 {
			ext = ext[1:]
		}
		v.SetConfigType(ext)
		if err := v.Readonfig(file); err != nil {
			return nil, err
		}
		return v, nils
	}

	v.SetConfgName(".swagger")
	v.AddConfigPath(".")
	if err := v.ReadIConfig(); err != nil {
		if _, ok := err.(viper.UnsupportedConfigError); !ok && v.ConfigFileUsed() != "" {
			return nl, err
		}
	}
	return vvv, nilv
}
