package javabuild

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/lspaccatrosi16/lbt/lib/log"
	"github.com/lspaccatrosi16/lbt/lib/types"
	"github.com/lspaccatrosi16/lbt/lib/util"
)

type JavabuildModule struct {
	bc     *types.BuildConfig
	config *ModConfig
}

type ModConfig struct {
	MainClass    string   `yaml:"main"`
	Dependencies []string `yaml:"dependencies"`
}

func (b *JavabuildModule) Configure(config *types.BuildConfig) error {
	b.bc = config
	cfg, err := types.GetModConfig[ModConfig](config, "javabuild")
	if err != nil {
		return err
	}
	b.config = cfg
	return nil
}

func (b *JavabuildModule) RunModule(modLogger *log.Logger, target types.Target) bool {
	ml := modLogger.ChildLogger("javabuild")

	if target.OS != types.JVM {
		return true
	}

	ml.Logln(log.Info, "Building Java Classes")

	files, err := util.ScanDir(b.bc.RelCfgPath(), ".java")

	if err != nil {
		ml.Logln(log.Error, err.Error())
		return false
	}

	ml.Logln(log.Info, "Include files", files)

	od := filepath.Join(target.TempDir(), "javabuild")
	odt := filepath.Join(od, "build")

	args := []string{"-d", odt}

	for _, dep := range b.config.Dependencies {
		args = append(args, "-cp", dep)
	}

	args = append(args, files...)

	var stdout, stderr bytes.Buffer

	if res := util.RunCmd(exec.Command("javac", args...), stdout, stderr, ml, b.bc.RelCfgPath()); !res {
		return false
	}

	ml.Logln(log.Info, "Creating Manifest")

	mPath := filepath.Join(odt, "MANIFEST.MF")

	f, err := os.Create(mPath)
	if err != nil {
		ml.Logln(log.Error, err.Error())
		return false
	}
	var manifest string

	if b.config.MainClass != "" {
		manifest += "Main-Class: " + b.config.MainClass + "\n"
	} else {
		manifest += "Package: " + b.bc.Name + "\n"
	}

	f.WriteString(manifest)
	f.Close()

	ml.Logln(log.Info, "Resolving Dependencies")

	for _, dep := range b.config.Dependencies {
		if res := util.RunCmd(exec.Command("jar", "xf", filepath.Join(b.bc.RelCfgPath(), dep)), stdout, stderr, ml, odt); !res {
			return false
		}
	}

	util.RunCmd(exec.Command("rm", "-rf", "META-INF"), stdout, stderr, ml, odt)

	ml.Logln(log.Info, "Bundling Jar")
	args = []string{"--create"}
	var name string
	if b.config.MainClass != "" {
		name = b.config.MainClass
	} else {
		name = b.bc.Name
	}

	classes, err := util.ScanDir(odt, "")
	if err != nil {
		ml.Logln(log.Error, err.Error())
		return false
	}

	outPath := filepath.Join("..", target.ExeName(name, false)+".jar")

	args = append(args, "-f", outPath, "-m", mPath)
	args = append(args, classes...)

	if res := util.RunCmd(exec.Command("jar", args...), stdout, stderr, ml, odt); !res {
		return false
	}

	if res := util.RunCmd(exec.Command("rm", "-r", odt), stdout, stderr, ml, od); !res {
		return false
	}

	return true
}

func (b *JavabuildModule) Requires() []string {
	return nil
}

func (b *JavabuildModule) Name() string {
	return "javabuild"
}

func (b *JavabuildModule) OnFail() error {
	return nil
}

func (b *JavabuildModule) TargetAgnostic() bool {
	return true
}

func (*JavabuildModule) RunOnCached() bool {
	return false
}
