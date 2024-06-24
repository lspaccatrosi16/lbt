package static

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/lspaccatrosi16/lbt/lib/log"
	"github.com/lspaccatrosi16/lbt/lib/types"
)

type StaticModule struct {
	bc     *types.BuildConfig
	config *ModConfig
}
type ModConfig struct {
	Structure string `yaml:"structure" validate:"required"`
	ExePath   string `yaml:"exePath" validate:"required"`
}

func (s *StaticModule) Configure(config *types.BuildConfig) error {
	s.bc = config
	cfg, err := types.GetModConfig[ModConfig](config, "static")
	if err != nil {
		return err
	}

	if cfg.Structure == "" {
		return fmt.Errorf("static module requires structure field")
	}

	if cfg.ExePath == "" {
		return fmt.Errorf("static module requires exePath field")
	}

	s.config = cfg
	return nil
}

func (s *StaticModule) RunModule(modLogger *log.Logger) error {
	ml := modLogger.ChildLogger("static")

	iPath := filepath.Join(s.bc.Cwd, s.config.Structure)
	exePath := filepath.Join(s.bc.Cwd, "tmp", "build")
	dE, err := os.ReadDir(exePath)
	if err != nil {
		return err
	}

	for _, entry := range dE {
		if !entry.IsDir() {
			name := strings.Split(entry.Name(), ".")[0]

			ml.Logf(log.Info, "Copying structure to %s", name)

			target, err := types.ParseTarget(strings.Split(name, "-")[1])
			if err != nil {
				return err
			}

			oPath := filepath.Join(s.bc.Cwd, "tmp", "static", name)

			err = os.MkdirAll(oPath, 0755)
			if err != nil {
				return err
			}
			err = copy(oPath, iPath)
			if err != nil {
				return err
			}

			exe, err := os.Open(filepath.Join(exePath, entry.Name()))
			if err != nil {
				return err
			}

			newExe := filepath.Join(oPath, s.config.ExePath, s.bc.Name)
			if target.OS == types.Windows {
				newExe += ".exe"
			}

			out, err := os.Create(newExe)
			if err != nil {
				return err
			}
			io.Copy(out, exe)
			exe.Close()
			out.Close()

			compressed := bytes.NewBuffer(nil)
			err = compress(oPath, compressed)
			if err != nil {
				return err
			}

			f, err := os.Create(filepath.Join(s.bc.Cwd, "tmp", "static", name+".tar.gz"))
			if err != nil {
				return err
			}
			io.Copy(f, compressed)
			f.Close()
			ml.Logf(log.Info, "Created %s.tar.gz", name)
			os.RemoveAll(oPath)
		}
	}
	return nil
}

func (s *StaticModule) Name() string {
	return "static"
}

func (s *StaticModule) Requires() []string {
	return []string{"build"}
}

func copy(dst string, src string) error {
	return filepath.Walk(filepath.Join(src), func(file string, fi fs.FileInfo, err error) error {
		if fi.IsDir() {
			relPath, err := filepath.Rel(src, file)
			if err != nil {
				return err
			}

			dstPath := filepath.Join(dst, relPath)

			os.MkdirAll(dstPath, 0755)
		} else {
			srcF, err := os.Open(file)
			if err != nil {
				return err
			}
			relPath, err := filepath.Rel(src, file)
			if err != nil {
				return err
			}

			dst, err := os.Create(filepath.Join(dst, relPath))
			if err != nil {
				return err
			}
			dst.Chmod(0755)
			io.Copy(dst, srcF)
			srcF.Close()
			dst.Close()
		}
		return nil
	})
}

func compress(src string, buf io.Writer) error {
	zr := gzip.NewWriter(buf)
	tw := tar.NewWriter(zr)

	err := filepath.Walk(src, func(file string, fi fs.FileInfo, _ error) error {
		header, err := tar.FileInfoHeader(fi, "")
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, file)
		if err != nil {
			return err
		}

		if relPath == "." {
			return nil
		}

		header.Name = filepath.ToSlash(relPath)

		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		if !fi.IsDir() {
			data, err := os.Open(file)
			if err != nil {
				return err
			}
			io.Copy(tw, data)
			data.Close()
		}
		return nil
	})
	if err != nil {
		return err
	}

	if err := tw.Close(); err != nil {
		return err
	}

	if err := zr.Close(); err != nil {
		return err
	}

	return nil
}
