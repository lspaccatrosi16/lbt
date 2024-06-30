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

func (c *CleanupModule) RunModule(modLogger *log.Logger) error {
	return os.RemoveAll(filepath.Join(c.bc.Cwd, "tmp"))
}

func (c *CleanupModule) Requires() []string {
	return nil
}

func (c *CleanupModule) OnFail() error {
	return nil
}
