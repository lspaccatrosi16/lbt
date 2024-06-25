package config

import (
	"fmt"
	"os"
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

	var config types.BuildConfig
	err = yaml.NewDecoder(f).Decode(&config)
	if err != nil {
		return nil, err
	}

	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	config.Cwd = wd

	if config.Name == "" {
		return nil, fmt.Errorf("config file requires name field")
	} else if strings.Contains(config.Name, ".") {
		return nil, fmt.Errorf("name field cannot contain '.'")
	}

	os.Chdir(config.Cwd)
	return &config, nil
}
