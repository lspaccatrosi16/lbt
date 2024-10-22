package output

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/lspaccatrosi16/lbt/lib/log"
	"github.com/lspaccatrosi16/lbt/lib/types"
	"github.com/lspaccatrosi16/lbt/lib/util"
)

type OutputModule struct {
	bc     *types.BuildConfig
	config *ModuleConfig
}

func (o *OutputModule) Name() string {
	return "output"
}

type ModuleConfig struct {
	Module string `yaml:"module" validate:"required"`
	OutDir string `yaml:"outDir" validate:"required"`
}

func (o *OutputModule) Configure(config *types.BuildConfig) error {
	o.bc = config
	cfg, err := types.GetModConfig[ModuleConfig](config, "output")
	if err != nil {
		return err
	}

	if cfg.Module == "" {
		return fmt.Errorf("output module requires module field")
	}

	if cfg.OutDir == "" {
		return fmt.Errorf("output module requires outDir field")
	}

	o.config = cfg
	return nil
}

func (o *OutputModule) RunModule(modLogger *log.Logger, target types.Target) bool {
	ml := modLogger.ChildLogger("output")

	oPath := o.bc.RelCfgPath(o.config.OutDir)

	s, err := os.Stat(oPath)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(oPath, 0755)
			if err != nil {
				ml.Logln(log.Error, err.Error())
				return false
			}
		} else {
			ml.Logln(log.Error, err.Error())
			return false
		}
	} else {
		if !s.IsDir() {
			ml.Logln(log.Error, fmt.Errorf("output directory %s is not a directory", o.config.OutDir))
			return false
		}
	}

	objDir := filepath.Join(target.TempDir(), o.config.Module)

	dE, err := os.ReadDir(objDir)
	if err != nil {
				ml.Logln(log.Error, err.Error())
				return false
	}

	for _, e := range dE {
		err = util.Copy(filepath.Join(oPath, e.Name()), filepath.Join(objDir, e.Name()))
		if err != nil {
				ml.Logln(log.Error, err.Error())
				return false
		}
		ml.Logf(log.Info, "Copied %s to %s", e.Name(), o.config.OutDir)
		o.bc.Produced = append(o.bc.Produced, filepath.Join(oPath, e.Name()))
	}

	return true 
}

func (o *OutputModule) Requires() []string {
	return []string{o.config.Module}
}

func (o *OutputModule) OnFail() error {
	return nil
}

func (o *OutputModule) TargetAgnostic() bool {
	return false
}

func (*OutputModule) RunOnCached() bool {
	return true 
}
