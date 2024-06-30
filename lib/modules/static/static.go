package static

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/lspaccatrosi16/lbt/lib/log"
	"github.com/lspaccatrosi16/lbt/lib/types"
	"github.com/lspaccatrosi16/lbt/lib/util"
)

type StaticModule struct {
	bc     *types.BuildConfig
	config *ModConfig
}
type ModConfig struct {
	Structure string `yaml:"structure" validate:"required"`
	ExePath   string `yaml:"exePath" validate:"required"`
}

func (s *StaticModule) Configure(config *types.BuildConfig) error {
	s.bc = config
	cfg, err := types.GetModConfig[ModConfig](config, "static")
	if err != nil {
		return err
	}

	if cfg.Structure == "" {
		return fmt.Errorf("static module requires structure field")
	}

	if cfg.ExePath == "" {
		return fmt.Errorf("static module requires exePath field")
	}

	s.config = cfg
	return nil
}

func (s *StaticModule) RunModule(modLogger *log.Logger) error {
	ml := modLogger.ChildLogger("static")

	iPath := filepath.Join(s.bc.Cwd, s.config.Structure)
	exePath := filepath.Join(s.bc.Cwd, "tmp", "build")
	dE, err := os.ReadDir(exePath)
	if err != nil {
		return err
	}

	for _, entry := range dE {
		if !entry.IsDir() {
			name := strings.Split(entry.Name(), ".")[0]

			ml.Logf(log.Info, "Copying structure to %s", name)

			target, err := types.ParseTarget(strings.Split(name, "-")[1])
			if err != nil {
				return err
			}

			oPath := filepath.Join(s.bc.Cwd, "tmp", "static", name)

			err = os.MkdirAll(oPath, 0755)
			if err != nil {
				return err
			}
			err = util.Copy(oPath, iPath)
			if err != nil {
				return err
			}

			exe, err := os.Open(filepath.Join(exePath, entry.Name()))
			if err != nil {
				return err
			}

			newExe := filepath.Join(oPath, s.config.ExePath, s.bc.Name)
			if target.OS == types.Windows {
				newExe += ".exe"
			}

			out, err := os.Create(newExe)
			if err != nil {
				return err
			}
			io.Copy(out, exe)
			exe.Close()
			out.Close()
			ml.Logf(log.Info, "Created static object %s", name)
		}
	}
	return nil
}

func (s *StaticModule) Name() string {
	return "static"
}

func (s *StaticModule) Requires() []string {
	return []string{"build"}
}

func (s *StaticModule) OnFail() error {
	return nil
}
