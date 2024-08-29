package build

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/lspaccatrosi16/lbt/lib/log"
	"github.com/lspaccatrosi16/lbt/lib/types"
)

type BuildModule struct {
	bc     *types.BuildConfig
	config *ModConfig
	wg     sync.WaitGroup
	errors chan error
}

type Command struct {
	Name string `yaml:"name" validate:"required"`
	Path string `yaml:"path" validate:"required"`
}

type ModConfig struct {
	Commands    []Command `yaml:"commands" validate:"required"`
	Ldflags     string    `yaml:"ldflags"`
	UsesVersion bool      `yaml:"version"`
	DisableCgo  bool      `yaml:"cgoOff"`
}

func (b *BuildModule) Configure(config *types.BuildConfig) error {
	b.bc = config
	cfg, err := types.GetModConfig[ModConfig](config, "build")
	if err != nil {
		return err
	}
	b.config = cfg
	b.errors = make(chan error)
	return nil
}

func (b *BuildModule) RunModule(modLogger *log.Logger) error {
	ml := modLogger.ChildLogger("build")

	if len(b.config.Commands) == 0 {
		ml.Logf(log.Info, "No commands to build")
		return nil
	}

	if len(b.bc.Targets) == 0 {
		ml.Logf(log.Info, "No targets to build")
		return nil
	}

	for _, cmd := range b.config.Commands {
		cmdPath := filepath.Join(b.bc.Cwd, cmd.Path)
		for _, target := range b.bc.Targets {
			b.wg.Add(1)
			go b.buildCommandTarget(ml, cmd, target, cmdPath)
		}
	}

	b.wg.Wait()
	close(b.errors)

	errs := []error{}
	for e := range b.errors {
		errs = append(errs, e)
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

func (b *BuildModule) buildCommandTarget(ml *log.Logger, cmd Command, target types.Target, cmdPath string) {
	defer b.wg.Done()
	ml.Logf(log.Info, "Building %s", target.ExeName(cmd.Name, true))
	err := target.Validate()
	if err != nil {
		b.errors <- err
		return
	}

	outPath := filepath.Join(b.bc.Cwd, "tmp", "build", target.ExeName(cmd.Name, true))
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

	var out, stderr bytes.Buffer
	eCmd.Stdout = &out
	eCmd.Stderr = &stderr

	err = eCmd.Run()
	if err != nil {
		b.errors <- errors.New(stderr.String())
		return
	}

	f, err := os.Open(outPath)
	if err != nil {
		b.errors <- err
		return
	}
	f.Chmod(0777)
	defer f.Close()

	ml.Logf(log.Info, "Built %s", target.ExeName(cmd.Name, true))
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

func (b *BuildModule) OnFail() error {
	return nil
}
