package version

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"

	"github.com/lspaccatrosi16/lbt/lib/log"
	"github.com/lspaccatrosi16/lbt/lib/types"
)

type VersionModule struct {
	bc     *types.BuildConfig
	config *ModuleConfig
	prev   string
}

type VersionType int

const (
	VersionBuildStr VersionType = iota
	VersionBuildInt
	VersionSemVer
)

func ParseVersionType(version string) (VersionType, error) {
	switch version {
	case "buildstr":
		return VersionBuildStr, nil
	case "buildint":
		return VersionBuildInt, nil
	case "semver":
		return VersionSemVer, nil
	default:
		return VersionBuildInt, fmt.Errorf("unknown version type: %s", version)
	}
}

type ModuleConfig struct {
	Path    string `yaml:"path" validate:"required"`
	VtS     string `yaml:"type" validate:"required"`
	VerType VersionType
}

func (v *VersionModule) Configure(config *types.BuildConfig) error {
	v.bc = config
	cfg, err := types.GetModConfig[ModuleConfig](config, "version")

	if err != nil {
		return err
	}

	if cfg.Path == "" {
		return fmt.Errorf("version module requires path field")
	}

	if cfg.VtS == "" {
		return fmt.Errorf("version module requires type field")
	}

	vt, err := ParseVersionType(cfg.VtS)
	if err != nil {
		return err
	}

	cfg.VerType = vt
	v.config = cfg

	return nil
}

func (v *VersionModule) RunModule(modLogger *log.Logger) error {
	ml := modLogger.ChildLogger("version")
	var newVersion string

	f, err := os.Open(filepath.Join(v.bc.Cwd, v.config.Path))
	if err == nil {
		by, err := io.ReadAll(f)
		if err != nil {
			return err
		}
		v.prev = string(by)
		f.Close()
	}

	switch v.config.VerType {
	case VersionBuildStr:
		newVersion = strconv.FormatInt(rand.Int63(), 36)
	case VersionBuildInt:
		f, err := os.Open(v.config.Path)
		if err != nil {
			return err
		}
		buf := bytes.NewBuffer(nil)
		io.Copy(buf, f)
		f.Close()
		curVer, err := strconv.Atoi(buf.String())
		if err != nil {
			return err
		}
		newVersion = strconv.Itoa(curVer + 1)
	case VersionSemVer:
		f, err := os.Open(v.config.Path)
		if err != nil {
			return err
		}
		buf := bytes.NewBuffer(nil)
		io.Copy(buf, f)
		f.Close()
		curVer := buf.String()
		verParts := []int{0, 0, 0, 0}
		fmt.Sscanf(curVer, "%d.%d.%d.%d", &verParts[0], &verParts[1], &verParts[2], &verParts[3])
		verParts[3]++
		newVersion = fmt.Sprintf("%d.%d.%d.%d", verParts[0], verParts[1], verParts[2], verParts[3])
	}

	ml.Logf(log.Info, "new version: %s", newVersion)

	f, err = os.Create(filepath.Join(v.bc.Cwd, v.config.Path))
	if err != nil {
		return err
	}
	f.WriteString(newVersion)
	f.Close()
	return nil
}

func (v *VersionModule) Name() string {
	return "version"
}

func (v *VersionModule) Requires() []string {
	return nil
}

func (v *VersionModule) OnFail() error {
	if v.prev != "" {
		f, err := os.Create(filepath.Join(v.bc.Cwd, v.config.Path))
		if err != nil {
			return err
		}
		f.WriteString(v.prev)
		f.Close()
	}
	return nil
}
