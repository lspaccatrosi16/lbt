package types

import (
	"fmt"
	"path/filepath"
	"strings"
)

type OS string
type Arch string

const (
	Windows OS = "windows"
	Linux   OS = "linux"
	MacOS   OS = "darwin"
)

const (
	AMD64 Arch = "amd64"
	ARM64 Arch = "arm64"
	ARM   Arch = "arm"
	i386  Arch = "386"
)

func ParseOS(s string) (OS, error) {
	switch s {
	case "windows":
		return Windows, nil
	case "linux":
		return Linux, nil
	case "darwin":
		return MacOS, nil
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
	case "386":
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

func (t Target) TempDir(cwd string) string {
	return filepath.Join(cwd, "tmp", t.String())
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
