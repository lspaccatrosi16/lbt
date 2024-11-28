package cbuild

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	"github.com/lspaccatrosi16/lbt/lib/log"
	"github.com/lspaccatrosi16/lbt/lib/types"
	"github.com/lspaccatrosi16/lbt/lib/util"
)

type CbuildModule struct {
	bc     *types.BuildConfig
	config *ModConfig
}

type ModConfig struct {
	Name        string   `yaml:"name"`
	SrcDir      string   `yaml:"source"`
	IncDir      []string `yaml:"include"`
	Compiler    string   `yaml:"compiler"`
	Flags       []string `yaml:"flags"`
	Main        string   `yaml:"main"`
	GenCC       bool     `yaml:"cc"`
	LibraryMode string   `yaml:"librarymode"`
	Libs        []string `yaml:"libs"`
	LibDirs     []string `yaml:"libdirs"`
}

type Command struct {
	Arguments []string `json:"arguments"`
	Directory string   `json:"directory"`
	File      string   `json:"file"`
	Output    string   `json:"output"`
}

type Commands []Command

func (c *Commands) Produce() ([]byte, error) {
	return json.Marshal(c)
}

func (b *CbuildModule) Configure(config *types.BuildConfig) error {
	b.bc = config
	cfg, err := types.GetModConfig[ModConfig](config, "cbuild")
	if err != nil {
		return err
	}

	if cfg.Name == "" {
		return fmt.Errorf("cbuild requires \"name\" to be set")
	} else if cfg.SrcDir == "" {
		return fmt.Errorf("cbuild requires \"source\" to be set")
	} else if cfg.Compiler == "" {
		return fmt.Errorf("cbuild requires \"compiler\" to be set")
	} else if cfg.Main == "" && cfg.LibraryMode == "" {
		return fmt.Errorf("cbuild requires either \"main\" or \"librarymode\" to be set")
	} else if cfg.Main != "" && cfg.LibraryMode != "" {
		return fmt.Errorf("cbuild requires only 1 of \"main\" and \"librarymode\" to be set")
	}

	if cfg.LibraryMode != "" {
		switch cfg.LibraryMode {
		case "shared", "static":
			break
		default:
			return fmt.Errorf("value of \"%s\" is not valid for field \"library\"", cfg.LibraryMode)
		}
	}

	b.config = cfg
	return nil
}

func (b *CbuildModule) RunModule(modLogger *log.Logger, target types.Target) bool {
	ml := modLogger.ChildLogger("cbuild")
	if !target.CmpRuntimeOS() {
		ml.Logln(log.Error, "cbuild does not support building for alternate OS")
		return false
	} else if !target.CmpRuntimeArch() {
		ml.Logln(log.Error, "cbuild does not support building for alternate arch")
		return false
	}

	var stdout, stderr bytes.Buffer
	var cmds = Commands{}

	buildDir := filepath.Join(target.TempDir(), "cbuild")
	if ok := util.RunCmd(exec.Command("mkdir", "-p", buildDir), stdout, stderr, ml, ""); !ok {
		return false
	}

	srcFiles, err := util.ScanDir(b.bc.RelCfgPath(b.config.SrcDir), ".c")
	if err != nil {
		ml.Logln(log.Error, err.Error())
		return false
	}

	ml.Logln(log.Info, "files", srcFiles)

	var incStrs []string
	for _, i := range b.config.IncDir {
		incStrs = append(incStrs, "-I"+b.bc.RelCfgPath(i))
	}

	var libStrs []string
	for _, l := range b.config.Libs {
		libStrs = append(libStrs, "-l"+l)
	}

	var libDirs []string
	for _, l := range b.config.LibDirs {
		libDirs = append(libDirs, "-L", l)
	}

	for _, f := range srcFiles {
		args := []string{"-o", filepath.Join(buildDir, repCO(f))}
		args = append(args, incStrs...)
		args = append(args, libStrs...)
		args = append(args, libDirs...)
		args = append(args, b.config.Flags...)
		args = append(args, "-c", filepath.Join(b.bc.RelCfgPath(b.config.SrcDir), f))

		if ok := util.RunCmd(exec.Command(b.config.Compiler, args...), stdout, stderr, ml, b.bc.RelCfgPath()); !ok {
			return false
		}

		absF, _ := filepath.Abs(f)
		absO, _ := filepath.Abs(filepath.Join(buildDir, repCO(f)))

		exe, _ := exec.LookPath(b.config.Compiler)
		cmds = append(cmds, Command{
			Arguments: append([]string{exe}, args...),
			Directory: b.bc.RelCfgPath(),
			File:      absF,
			Output:    absO,
		})
	}

	objFiles, err := util.ScanDir(buildDir, ".o")

	if err != nil {
		ml.Logln(log.Error, err.Error())
		return false
	}

	if b.config.Main != "" {
		args := []string{"-o", filepath.Join(buildDir, b.config.Name)}
		args = append(args, incStrs...)
		args = append(args, libStrs...)
		args = append(args, libDirs...)
		args = append(args, b.config.Flags...)
		args = append(args, objFiles...)

		if ok := util.RunCmd(exec.Command(b.config.Compiler, args...), stdout, stderr, ml, buildDir); !ok {
			return false
		}
	} else if b.config.LibraryMode == "static" {
		args := []string{"rcs", filepath.Join(buildDir, b.config.Name) + ".a"}
		args = append(args, objFiles...)

		if ok := util.RunCmd(exec.Command("ar", args...), stdout, stderr, ml, buildDir); !ok {
			return false
		}
	} else if b.config.LibraryMode == "shared" {
		args := []string{"-o", filepath.Join(buildDir, b.config.Name) + ".so", "-shared"}
		args = append(args, incStrs...)
		args = append(args, libStrs...)
		args = append(args, libDirs...)
		args = append(args, b.config.Flags...)
		args = append(args, objFiles...)
		if ok := util.RunCmd(exec.Command(b.config.Compiler, args...), stdout, stderr, ml, buildDir); !ok {
			return false
		}
	}

	if ok := util.RunCmd(exec.Command("rm", objFiles...), stdout, stderr, ml, buildDir); !ok {
		return false
	}

	if b.config.GenCC {
		f, err := os.Create(b.bc.RelCfgPath("compile_commands.json"))
		if err != nil {
			ml.Logln(log.Error, err.Error())
			return false
		}
		by, err := cmds.Produce()
		if err != nil {
			ml.Logln(log.Error, err.Error())
			return false
		}
		f.Write(by)
		f.Close()
	}

	return true
}

func repCO(s string) string {
	return regexp.MustCompilePOSIX(`\.c$`).ReplaceAllString(s, ".o")
}

func (b *CbuildModule) Requires() []string {
	return nil
}

func (b *CbuildModule) Name() string {
	return "cbuild"
}

func (b *CbuildModule) OnFail() error {
	return nil
}

func (b *CbuildModule) TargetAgnostic() bool {
	return false
}

func (*CbuildModule) RunOnCached() bool {
	return false
}
