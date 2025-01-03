package types

import (
	"fmt"
	"path/filepath"

	"github.com/lspaccatrosi16/lbt/lib/log"
	"gopkg.in/yaml.v3"
)

type ModuleConfig struct {
	Name   string                 `yaml:"name"`
	Config map[string]interface{} `yaml:"config"`
}

type VerConfig struct {
	Path string `yaml:"path"`
	VtS  string `yaml:"type"`
}

func NewBuildConfig(loc string) *BuildConfig {
	return &BuildConfig{loc: loc}
}

type BuildConfig struct {
	Name        string         `yaml:"name"`
	Targets     []Target       `yaml:"targets"`
	Modules     []ModuleConfig `yaml:"modules"`
	IncludeDirs []string       `yaml:"includeDirs"`
	Version     VerConfig      `yaml:"version"`
	Produced    []string
	loc         string
}

func (b *BuildConfig) RelCfgPath(paths ...string) string {
	return filepath.Join(append([]string{b.loc}, paths...)...)
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
	RunModule(*log.Logger, Target) bool
	Configure(*BuildConfig) error
	Requires() []string
	TargetAgnostic() bool
	RunOnCached() bool
}
