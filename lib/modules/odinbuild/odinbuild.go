package odinbuild

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/lspaccatrosi16/lbt/lib/log"
	"github.com/lspaccatrosi16/lbt/lib/types"
)

type OdinbuildModule struct {
	bc     *types.BuildConfig
	config *ModConfig
}

type ModConfig struct {
	Src      string   `yaml:"src"`
	Optimise string   `yaml:"optimise"`
	Debug    bool     `yaml:"debug"`
	Flags    []string `yaml:"flags"`
}

func (b *OdinbuildModule) Configure(config *types.BuildConfig) error {
	b.bc = config
	cfg, err := types.GetModConfig[ModConfig](config, "odinbuild")
	if err != nil {
		return err
	}

	if cfg.Optimise == "" {
		cfg.Optimise = "minimal"
	}

	b.config = cfg
	return nil
}

func (b *OdinbuildModule) RunModule(modLogger *log.Logger, target types.Target) bool {
	ml := modLogger.ChildLogger("odinbuild")

	outPath := filepath.Join(target.TempDir(), "odinbuild", target.ExeName(b.bc.Name, true))

	// var err error
	var stdout, stderr bytes.Buffer
	cmd := exec.Command("mkdir", "-p", filepath.Dir(outPath))

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		ml.Logln(log.Error, stderr.String())
		ml.Logln(log.Error, stdout.String())
		return false
	}

	stdout.Reset()
	stderr.Reset()

	args := []string{"build", b.bc.RelCfgPath(b.config.Src)}
	args = append(args, fmt.Sprintf("-out:%s", outPath))
	args = append(args, fmt.Sprintf("-o:%s", b.config.Optimise))
	args = append(args, fmt.Sprintf("-target:%s", target.String()))
	args = append(args, b.config.Flags...)

	if b.config.Debug {
		args = append(args, "-debug")
	}

	cmd = exec.Command("odin", args...)
	ml.Logf(log.Info, "command: odin %s\n", strings.Join(args, " "))

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		ml.Logln(log.Error, stdout.String())
		ml.Logln(log.Error, stderr.String())
		return false
	}

	f, err := os.Open(outPath)
	if err != nil {
		ml.Logln(log.Error, stdout.String())
		return false
	}
	f.Chmod(0777)
	f.Close()

	ml.Logf(log.Info, "Built %s", target.ExeName(b.bc.Name, true))

	return true
}

func (b *OdinbuildModule) Requires() []string {
	return nil
}

func (b *OdinbuildModule) Name() string {
	return "odinbuild"
}

func (b *OdinbuildModule) OnFail() error {
	return nil
}

func (b *OdinbuildModule) TargetAgnostic() bool {
	return false
}

func (*OdinbuildModule) RunOnCached() bool {
	return false
}
