package setup

import (
	"os"
	"path/filepath"

	"github.com/lspaccatrosi16/lbt/lib/log"
	"github.com/lspaccatrosi16/lbt/lib/types"
)

type SetupModule struct {
	bc *types.BuildConfig
}

func (i *SetupModule) Name() string {
	return "init"
}

func (i *SetupModule) Configure(config *types.BuildConfig) error {
	i.bc = config
	return nil
}

func (i *SetupModule) RunModule(*log.Logger) error {
	return os.Mkdir(filepath.Join(i.bc.Cwd, "tmp"), 0755)
}

func (i *SetupModule) Requires() []string {
	return nil
}

func (i *SetupModule) OnFail() error {
	return nil
}
