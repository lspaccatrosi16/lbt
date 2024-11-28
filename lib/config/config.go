package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/lspaccatrosi16/go-cli-tools/args"
	"github.com/lspaccatrosi16/lbt/lib/types"
	"gopkg.in/yaml.v3"
)

func ParseConfig() (*types.BuildConfig, error) {
	cf, err := args.GetFlagValue[string]("config")
	if err != nil {
		return nil, err
	}

	f, err := os.Open(cf)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	cfgAPath, err := filepath.Abs(cf)
	if err != nil {
		return nil, err
	}

	config := types.NewBuildConfig(filepath.Dir(cfgAPath))
	err = yaml.NewDecoder(f).Decode(config)
	if err != nil {
		return nil, err
	}

	if config.Name == "" {
		return nil, fmt.Errorf("config file requires name field")
	} else if strings.Contains(config.Name, ".") {
		return nil, fmt.Errorf("name field cannot contain '.'")
	}

	if len(config.Targets) == 0 {
		return nil, fmt.Errorf("no targets provided")
	}

	for _, t := range config.Targets {
		if err := t.Validate(); err != nil {
			return nil, err
		}
	}

	return config, nil
}
