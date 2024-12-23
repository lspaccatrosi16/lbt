package vbuild

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/lspaccatrosi16/lbt/lib/log"
	"github.com/lspaccatrosi16/lbt/lib/types"
)

type VbuildModule struct {
	bc     *types.BuildConfig
	config *ModConfig
}

type ModConfig struct {
	Src     string   `yaml:"src"`
	Backend string   `yaml:"backend"`
	Debug   bool     `yaml:"debug"`
	Flags   []string `yaml:"flags"`
}

func (b *VbuildModule) Configure(config *types.BuildConfig) error {
	b.bc = config
	cfg, err := types.GetModConfig[ModConfig](config, "vbuild")
	if err != nil {
		return err
	}

	if cfg.Backend == "" {
		cfg.Backend = "c"
	}

	b.config = cfg
	return nil
}

func (b *VbuildModule) RunModule(modLogger *log.Logger, target types.Target) bool {
	ml := modLogger.ChildLogger("vbuild")

	outPath := filepath.Join(target.TempDir(), "vbuild", target.ExeName(b.bc.Name, true))

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

	args := []string{b.bc.RelCfgPath(b.config.Src)}
	args = append(args, "-o", outPath)
	args = append(args, "-os", string(target.OS))
	args = append(args, "-arch", string(target.Arch))
	args = append(args, "-backend", b.config.Backend)

	if b.config.Debug {
		args = append(args, "-g")
	}

	args = append(args, b.config.Flags...)

	cmd = exec.Command("v", args...)
	ml.Logf(log.Info, "command: v %s\n", strings.Join(args, " "))

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
		ml.Logln(log.Error, stderr.String())
		ml.Logln(log.Error, err.Error())
		return false
	}
	f.Chmod(0777)
	f.Close()

	ml.Logf(log.Info, "Built %s", target.ExeName(b.bc.Name, true))

	return true
}

func (b *VbuildModule) Requires() []string {
	return nil
}

func (b *VbuildModule) Name() string {
	return "vbuild"
}

func (b *VbuildModule) OnFail() error {
	return nil
}

func (b *VbuildModule) TargetAgnostic() bool {
	return false
}

func (*VbuildModule) RunOnCached() bool {
	return false
}
