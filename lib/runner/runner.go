package runner

import (
	"fmt"

	"github.com/lspaccatrosi16/lbt/lib/log"
	"github.com/lspaccatrosi16/lbt/lib/types"
)

var modsRun = map[string]bool{}
var mods map[string]types.Module

func RunModules(config *types.BuildConfig, modList []types.Module) error {
	if mods == nil {
		setup(modList)
	}

	for _, mod := range config.Modules {
		rmod, ok := mods[mod.Name]
		if !ok {
			return fmt.Errorf("module %s not found", mod.Name)
		}
		err := RunModule(rmod, config)
		if err != nil {
			return err
		}
	}

	return nil
}

func RunModule(mod types.Module, config *types.BuildConfig) error {
	if modsRun[mod.Name()] {
		return nil
	}
	modsRun[mod.Name()] = true

	ml := log.Default.ChildLogger("module")

	err := mod.Configure(config)
	if err != nil {
		return err
	}

	reqs := mod.Requires()

	for _, req := range reqs {
		if rmod, ok := mods[req]; ok {
			err = RunModule(rmod, config)
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("module %s requires %s, but it is not available", mod.Name(), req)
		}
	}

	return mod.RunModule(ml)
}

func setup(list []types.Module) error {
	mods = map[string]types.Module{}

	for _, mod := range list {
		mods[mod.Name()] = mod
	}

	return nil
}
