package gobuild

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
	Commands   []Command `yaml:"commands" validate:"required"`
	Ldflags    string    `yaml:"ldflags"`
	DisableCgo bool      `yaml:"cgoOff"`
	Root       string    `yaml:"root"`
}

func (b *BuildModule) Configure(config *types.BuildConfig) error {
	b.bc = config
	cfg, err := types.GetModConfig[ModConfig](config, "gobuild")
	if err != nil {
		return err
	}
	b.config = cfg
	return nil
}

func (b *BuildModule) RunModule(modLogger *log.Logger, target types.Target) bool {
	ml := modLogger.ChildLogger("gobuild")

	if len(b.config.Commands) == 0 {
		ml.Logf(log.Info, "No commands to build")
		return true
	}

	for _, cmd := range b.config.Commands {
		cmdPath := filepath.Join(b.bc.Cwd, cmd.Path)
		err := b.buildCommandTarget(ml, cmd, target, cmdPath)
		if err != nil {
			ml.Logln(log.Error, err.Error())
			return false
		}
	}

	return true
}

func (b *BuildModule) buildCommandTarget(ml *log.Logger, cmd Command, target types.Target, cmdPath string) error {
	ml.Logf(log.Info, "Building %s", target.ExeName(cmd.Name, true))
	err := target.Validate()
	if err != nil {
		return err
	}

	outPath := filepath.Join(target.TempDir(b.bc.Cwd), "gobuild", target.ExeName(cmd.Name, true))
	args := []string{"build", "-o", outPath}
	if b.config.Ldflags != "" {
		args = append(args, "-ldflags", b.config.Ldflags)
	}
	args = append(args, cmdPath)
	eCmd := exec.Command("go", args...)
	eCmd.Env = os.Environ()

	eCmd.Env = append(eCmd.Env, "GOOS="+string(target.OS))
	eCmd.Env = append(eCmd.Env, "GOARCH="+string(target.Arch))
	if b.config.DisableCgo {
		eCmd.Env = append(eCmd.Env, "CGO_ENABLED=0")
	}

	if b.config.Root != "" {
		eCmd.Dir = b.config.Root
	}

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

	ml.Logf(log.Info, "Built %s", target.ExeName(cmd.Name, true))
	return nil
}

func (b *BuildModule) Requires() []string {
	return nil
}

func (b *BuildModule) Name() string {
	return "gobuild"
}

func (b *BuildModule) OnFail() error {
	return nil
}

func (b *BuildModule) TargetAgnostic() bool {
	return false
}

func (*BuildModule) RunOnCached() bool {
	return false
}
