package output

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/lspaccatrosi16/lbt/lib/log"
	"github.com/lspaccatrosi16/lbt/lib/types"
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

func (o *OutputModule) RunModule(modLogger *log.Logger) error {
	ml := modLogger.ChildLogger("output")

	oPath := filepath.Join(o.bc.Cwd, o.config.OutDir)

	s, err := os.Stat(oPath)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(oPath, 0755)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		if !s.IsDir() {
			return fmt.Errorf("output directory %s is not a directory", o.config.OutDir)
		}
	}

	iPath := filepath.Join(o.bc.Cwd, "tmp", o.config.Module)

	dE, err := os.ReadDir(iPath)
	if err != nil {
		return err
	}

	for _, e := range dE {
		if !e.IsDir() {
			ml.Logf(log.Info, "Copying %s to %s", e.Name(), oPath)
			src, err := os.Open(filepath.Join(iPath, e.Name()))
			if err != nil {
				return err
			}
			dst, err := os.Create(filepath.Join(oPath, e.Name()))
			if err != nil {
				return err
			}
			io.Copy(dst, src)
			src.Close()
			dst.Close()
		}
	}

	return nil

}

func (o *OutputModule) Requires() []string {
	return []string{o.config.Module}
}
