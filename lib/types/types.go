package types

import (
	"fmt"

	"github.com/lspaccatrosi16/lbt/lib/log"
	"gopkg.in/yaml.v3"
)

type ModuleConfig struct {
	Name   string                 `yaml:"name"`
	Config map[string]interface{} `yaml:"config"`
}

type BuildConfig struct {
	Name    string `yaml:"name"`
	Cwd     string
	Modules []ModuleConfig `yaml:"modules"`
}

func GetModConfig[T any](b *BuildConfig, name string) (*T, error) {
	cfg, err := b.modConfig(name)
	if err != nil {
		return nil, err
	}
	by, err := yaml.Marshal(cfg)
	if err != nil {
		return nil, err
	}

	var out T
	err = yaml.Unmarshal(by, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (b *BuildConfig) modConfig(name string) (map[string]interface{}, error) {
	for _, mod := range b.Modules {
		if mod.Name == name {
			return mod.Config, nil
		}
	}
	return nil, fmt.Errorf("module %s has not been configured", name)
}

type Module interface {
	Name() string
	RunModule(*log.Logger) error
	Configure(*BuildConfig) error
	Requires() []string
	OnFail() error
}
