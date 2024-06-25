package build

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/lspaccatrosi16/lbt/lib/log"
	"github.com/lspaccatrosi16/lbt/lib/types"
)

type BuildModule struct {
	bc     *types.BuildConfig
	config *ModConfig
}

type Command struct {
	Name string `yaml:"name" validate:"required"`
	Path string `yaml:"path" validate:"required"`
}

type ModConfig struct {
	Commands    []Command      `yaml:"commands" validate:"required"`
	Ldflags     string         `yaml:"ldflags"`
	Targets     []types.Target `yaml:"targets" validate:"required"`
	UsesVersion bool           `yaml:"version"`
}

func (b *BuildModule) Configure(config *types.BuildConfig) error {
	b.bc = config
	cfg, err := types.GetModConfig[ModConfig](config, "build")
	if err != nil {
		return err
	}
	b.config = cfg
	return nil
}

func (b *BuildModule) RunModule(modLogger *log.Logger) error {
	ml := modLogger.ChildLogger("build")

	if len(b.config.Commands) == 0 {
		ml.Logf(log.Info, "No commands to build")
		return nil
	}

	if len(b.config.Targets) == 0 {
		ml.Logf(log.Info, "No targets to build")
		return nil
	}

	for _, cmd := range b.config.Commands {
		cmdPath := filepath.Join(b.bc.Cwd, cmd.Path)
		for _, target := range b.config.Targets {
			ml.Logf(log.Info, "Building %s %s", cmd, target.String())
			err := target.Validate()
			if err != nil {
				return err
			}

			outPath := filepath.Join(b.bc.Cwd, "tmp", "build", cmd.Name+"-"+target.String())
			if target.OS == types.Windows {
				outPath += ".exe"
			}
			args := []string{"build", "-o", outPath}
			if b.config.Ldflags != "" {
				args = append(args, "-ldflags", b.config.Ldflags)
			}
			args = append(args, cmdPath)
			eCmd := exec.Command("go", args...)
			eCmd.Env = os.Environ()

			eCmd.Env = append(eCmd.Env, "GOOS="+string(target.OS))
			eCmd.Env = append(eCmd.Env, "GOARCH="+string(target.Arch))

			var out, stderr bytes.Buffer
			eCmd.Stdout = &out
			eCmd.Stderr = &stderr

			err = eCmd.Run()
			if err != nil {
				return errors.New(stderr.String())
			}

			f, err := os.Open(outPath)
			if err != nil {
				return err
			}
			f.Chmod(0777)
			defer f.Close()

			ml.Logf(log.Info, "Built %s %s", cmd, target.String())
		}
	}

	return nil
}

func (b *BuildModule) Requires() []string {
	if b.config.UsesVersion {
		return []string{"version"}
	}
	return nil
}

func (b *BuildModule) Name() string {
	return "build"
}
