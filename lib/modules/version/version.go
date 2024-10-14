package version

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"

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
	VerType VersionType
}

func (v *VersionModule) Configure(config *types.BuildConfig) error {
	v.bc = config

	if config.Version.Path == "" {
		return nil
	}

	if config.Version.VtS == "" {
		return nil
	}

	vt, err := ParseVersionType(config.Version.VtS)
	if err != nil {
		return err
	}

	cfg := &ModuleConfig{
		VerType: vt,
	}

	v.config = cfg
	return nil
}

func (v *VersionModule) RunModule(modLogger *log.Logger, _ types.Target) bool {
	if v.config == nil {
		return true 
	}

	ml := modLogger.ChildLogger("version")
	var newVersion string

	f, err := os.Open(filepath.Join(v.bc.Cwd, v.bc.Version.Path))
	if err == nil {
		by, err := io.ReadAll(f)
		if err != nil {
			ml.Logln(log.Error, err.Error())
			return false
		}
		v.prev = string(by)
		f.Close()
	}

	switch v.config.VerType {
	case VersionBuildStr:
		newVersion = strconv.FormatInt(rand.Int63(), 36)
	case VersionBuildInt:
		f, err := os.Open(v.bc.Version.Path)
		if err != nil {
			ml.Logln(log.Error, err.Error())
			return false
		}
		buf := bytes.NewBuffer(nil)
		io.Copy(buf, f)
		f.Close()
		curVer, err := strconv.Atoi(strings.Trim(buf.String(), " \r\n\t"))
		if err != nil {
			ml.Logln(log.Error, err.Error())
			return false
		}
		newVersion = strconv.Itoa(curVer + 1)
	case VersionSemVer:
		f, err := os.Open(v.bc.Version.Path)
		if err != nil {
			ml.Logln(log.Error, err.Error())
			return false
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

	f, err = os.Create(filepath.Join(v.bc.Cwd, v.bc.Version.Path))
	if err != nil {
		ml.Logln(log.Error, err.Error())
		return false
	}
	f.WriteString(newVersion)
	f.Close()
	return true
}

func (v *VersionModule) Name() string {
	return "version"
}

func (v *VersionModule) Requires() []string {
	return nil
}

func (v *VersionModule) OnFail() error {
	if v.prev != "" {
		f, err := os.Create(filepath.Join(v.bc.Cwd, v.bc.Version.Path))
		if err != nil {
			return err
		}
		f.WriteString(v.prev)
		f.Close()
	}
	return nil
}

func (v *VersionModule) TargetAgnostic() bool {
	return true
}

func (*VersionModule) RunOnCached() bool {
	return false
}
