package cleanup

import (
	"os"
	"path/filepath"

	"github.com/lspaccatrosi16/lbt/lib/log"
	"github.com/lspaccatrosi16/lbt/lib/types"
)

type CleanupModule struct {
	bc *types.BuildConfig
}

func (c *CleanupModule) Name() string {
	return "cleanup"
}

func (c *CleanupModule) Configure(config *types.BuildConfig) error {
	c.bc = config
	return nil
}

func (c *CleanupModule) RunModule(modLogger *log.Logger, _ types.Target) bool {
	ml := modLogger.ChildLogger("cleanup")
	err := os.RemoveAll(filepath.Join(c.bc.Cwd, "tmp"))
	if err != nil {
		ml.Logln(log.Error, err.Error())
		return false
	}
	return true
}

func (c *CleanupModule) Requires() []string {
	return nil
}

func (c *CleanupModule) OnFail() error {
	return nil
}

func (c *CleanupModule) TargetAgnostic() bool {
	return true
}

func (*CleanupModule) RunOnCached() bool {
	return true 
}
