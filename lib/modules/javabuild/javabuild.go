package javabuild

import (
	"bytes"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/lspaccatrosi16/lbt/lib/log"
	"github.com/lspaccatrosi16/lbt/lib/types"
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

	files, err := scanDir(b.bc.Cwd, ".java")

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

	if res := runCmd(exec.Command("javac", args...), stdout, stderr, ml, b.bc.Cwd); !res {
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
		if res := runCmd(exec.Command("jar", "xf", filepath.Join(b.bc.Cwd, dep)), stdout, stderr, ml, odt); !res {
			return false
		}
	}

	runCmd(exec.Command("rm", "-rf", "META-INF"), stdout, stderr, ml, odt)

	ml.Logln(log.Info, "Bundling Jar")
	args = []string{"--create"}
	var name string
	if b.config.MainClass != "" {
		name = b.config.MainClass
	} else {
		name = b.bc.Name
	}

	classes, err := scanDir(odt, "")
	if err != nil {
		ml.Logln(log.Error, err.Error())
		return false
	}

	args = append(args, "-f", filepath.Join("..", target.ExeName(name, false)+".jar"), "-m", mPath)
	args = append(args, classes...)

	ml.Logln(log.Info, args)

	if res := runCmd(exec.Command("jar", args...), stdout, stderr, ml, odt); !res {
		return false
	}

	if res := runCmd(exec.Command("rm", "-r", odt), stdout, stderr, ml, od); !res {
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

func runCmd(eCmd *exec.Cmd, stdout, stderr bytes.Buffer, ml *log.Logger, dir string) bool {
	eCmd.Dir = dir
	eCmd.Stdout = &stdout
	eCmd.Stderr = &stderr

	err := eCmd.Run()
	if err != nil {
		ml.Logln(log.Error, err.Error())
		ml.Logln(log.Error, stdout.String())
		ml.Logln(log.Error, stderr.String())
		return false
	}

	stdout.Reset()
	stderr.Reset()

	return true
}

func scanDir(startDir, suffix string) ([]string, error) {
	files := []string{}
	err := filepath.WalkDir(startDir, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		if strings.HasSuffix(d.Name(), suffix) {
			rp, err := filepath.Rel(startDir, path)
			if err != nil {
				return err
			}
			files = append(files, rp)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return files, nil
}
