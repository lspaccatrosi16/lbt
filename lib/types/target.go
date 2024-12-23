package types

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type OS string
type Arch string

var timestamp = fmt.Sprint(time.Now().Unix())

const (
	Windows OS = "windows"
	Linux   OS = "linux"
	MacOS   OS = "darwin"
	JVM     OS = "jvm"
	Android OS = "android"
)

const (
	AMD64 Arch = "amd64"
	ARM64 Arch = "arm64"
	ARM   Arch = "arm"
	i386  Arch = "i386"
)

func ParseOS(s string) (OS, error) {
	switch s {
	case "windows":
		return Windows, nil
	case "linux":
		return Linux, nil
	case "darwin":
		return MacOS, nil
	case "jvm":
		return JVM, nil
	case "android":
		return Android, nil
	default:
		return "", fmt.Errorf("unknown OS: %s", s)
	}
}

func ParseArch(s string) (Arch, error) {
	switch s {
	case "amd64":
		return AMD64, nil
	case "arm64":
		return ARM64, nil
	case "arm":
		return ARM, nil
	case "i386":
		return i386, nil
	default:
		return "", fmt.Errorf("unknown arch: %s", s)
	}
}

type Target struct {
	OS   OS   `yaml:"os"`
	Arch Arch `yaml:"arch"`
}

func (t *Target) Validate() error {
	var err error
	t.OS, err = ParseOS(string(t.OS))
	if err != nil {
		return err
	}

	t.Arch, err = ParseArch(string(t.Arch))
	if err != nil {
		return err
	}
	return nil
}

func (t Target) String() string {
	return fmt.Sprintf("%s_%s", t.OS, t.Arch)
}

func (t Target) ExeName(n string, addExe bool) string {
	str := fmt.Sprintf("%s-%s", n, t.String())
	if t.OS == Windows && addExe {
		return str + ".exe"
	}
	return str
}

func (t Target) CleanName(n string, addExe bool) string {
	if t.OS == Windows && addExe {
		return n + ".exe"
	}
	return n
}

func (t Target) TempDir() string {
	tmpDir := os.TempDir()
	tS := t.String()
	if t.Arch == "" && t.OS == "" {
		tS = ""
	}
	return filepath.Join(tmpDir, "lbt", timestamp, tS)
}

func (t Target) CmpRuntimeOS() bool {
	return string(t.OS) == runtime.GOOS
}

func (t Target) CmpRuntimeArch() bool {
	return string(t.Arch) == runtime.GOARCH
}

func ParseTarget(s string) (Target, error) {
	comps := strings.Split(s, "_")
	if len(comps) != 2 {
		return Target{}, fmt.Errorf("invalid target string: %s", s)
	}

	os, err := ParseOS(comps[0])
	if err != nil {
		return Target{}, err
	}

	arch, err := ParseArch(comps[1])
	if err != nil {
		return Target{}, err
	}

	return Target{OS: os, Arch: arch}, nil
}

var NoTarget = Target{}
