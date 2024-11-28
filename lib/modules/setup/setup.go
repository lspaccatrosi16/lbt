package setup

import (
	"os"

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

func (i *SetupModule) RunModule(modLogger *log.Logger, _ types.Target) bool {
	ml := modLogger.ChildLogger("setup")
	err := os.MkdirAll(types.NoTarget.TempDir(), 0755)
	if err != nil {
		ml.Logln(log.Error, err.Error())
		return false
	}
	return true
}

func (i *SetupModule) Requires() []string {
	return nil
}

func (i *SetupModule) OnFail() error {
	return nil
}

func (i *SetupModule) TargetAgnostic() bool {
	return true
}

func (*SetupModule) RunOnCached() bool {
	return true
}
