package cached

import (
	"os"
	"path/filepath"

	"github.com/lspaccatrosi16/lbt/lib/cache"
	"github.com/lspaccatrosi16/lbt/lib/log"
	"github.com/lspaccatrosi16/lbt/lib/types"
	"github.com/lspaccatrosi16/lbt/lib/util"
)

type GetCachedModule struct {
	bc   *types.BuildConfig
	Meta *cache.BuildMeta
}

func (g *GetCachedModule) Name() string {
	return "getCached"
}

func (g *GetCachedModule) Configure(config *types.BuildConfig) error {
	g.bc = config
	return nil
}

func (g *GetCachedModule) RunModule(modLogger *log.Logger, target types.Target) bool {
	ml := modLogger.ChildLogger("getCached")
	ml.Logln(log.Info, "Source files unchanged, using cached build artifact")
	based := filepath.Join(target.TempDir(), "getCached")
	err := os.MkdirAll(based, 0755)
	if err != nil {
		ml.Logln(log.Error, err.Error())
		return false
	}
	for _, obj := range g.Meta.Objects {
		err := util.Copy(filepath.Join(based, obj), filepath.Join(g.Meta.Location(), obj))
		if err != nil {
			ml.Logln(log.Error, err.Error())
			return false
		}
	}
	return true
}

func (g *GetCachedModule) Requires() []string {
	return nil
}

func (g *GetCachedModule) OnFail() error {
	return nil
}

func (g *GetCachedModule) TargetAgnostic() bool {
	return true
}

func (*GetCachedModule) RunOnCached() bool {
	return true 
}
