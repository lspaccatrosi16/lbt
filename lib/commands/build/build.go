package build

import (
	"github.com/lspaccatrosi16/lbt/lib/config"
	"github.com/lspaccatrosi16/lbt/lib/modules"
	"github.com/lspaccatrosi16/lbt/lib/runner"
)

func Run() error {
	config, err := config.ParseConfig()
	if err != nil {
		return err
	}

	err = runner.RunModules(config, modules.List)
	if err != nil {
		return err
	}

	return nil
}
