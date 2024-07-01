package static

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/lspaccatrosi16/lbt/lib/log"
	"github.com/lspaccatrosi16/lbt/lib/types"
	"github.com/lspaccatrosi16/lbt/lib/util"
)

type StaticModule struct {
	bc     *types.BuildConfig
	config *ModConfig
}

type StaticExec struct {
	Command string `yaml:"command" validate:"required"`
	Path    string `yaml:"path" validate:"required"`
}

type Structure struct {
	Name        string       `yaml:"name" validate:"required"`
	Path        string       `yaml:"path" validate:"required"`
	Executables []StaticExec `yaml:"executables" validate:"required"`
}

type ModConfig struct {
	Structures []Structure `yaml:"structures" validate:"required"`
}

func (s *StaticModule) Configure(config *types.BuildConfig) error {
	s.bc = config
	cfg, err := types.GetModConfig[ModConfig](config, "static")
	if err != nil {
		return err
	}

	for i, str := range cfg.Structures {
		if str.Name == "" {
			return fmt.Errorf("static module requires name field in structure %d", i+1)
		}

		if str.Path == "" {
			return fmt.Errorf("static module requires path field in structure %d", i+1)
		}

		for j, exe := range str.Executables {
			if exe.Command == "" {
				return fmt.Errorf("static module requires command field in structure %d executable %d", i+1, j+1)
			}
			if exe.Path == "" {
				return fmt.Errorf("static module requires path field in structure %d executable %d", i+1, j+1)
			}
		}
	}
	s.config = cfg
	return nil
}

func (s *StaticModule) RunModule(modLogger *log.Logger) error {
	ml := modLogger.ChildLogger("static")

	exeDir := filepath.Join(s.bc.Cwd, "tmp", "build")
	oPath := filepath.Join(s.bc.Cwd, "tmp", "static")
	for _, str := range s.config.Structures {
		iPath := filepath.Join(s.bc.Cwd, str.Path)
		for _, target := range s.bc.Targets {
			sName := target.ExeName(str.Name, false)
			ml.Logf(log.Info, "Creating structure %s", sName)
			sOut := filepath.Join(oPath, sName)
			err := os.MkdirAll(sOut, 0755)
			if err != nil {
				return err
			}

			err = util.Copy(sOut, iPath)
			if err != nil {
				return err
			}

			for _, exe := range str.Executables {
				exeName := target.ExeName(exe.Command, true)
				newName := target.CleanName(exe.Command, true)
				err = util.Copy(filepath.Join(sOut, exe.Path, newName), filepath.Join(exeDir, exeName))
				if err != nil {
					return err
				}
			}
			ml.Logf(log.Info, "Created structure %s", sName)
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
