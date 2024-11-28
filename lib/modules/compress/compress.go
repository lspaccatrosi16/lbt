package compress

import (
	"archive/tar"
	"archive/zip"
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

type compressionFormat string

const (
	cfTarGz = "tar.gz"
	cfZip   = "zip"
)

func parseCompressionFormat(format string) (compressionFormat, error) {
	switch format {
	case "tar.gz":
		return cfTarGz, nil
	case "zip":
		return cfZip, nil
	default:
		return "", fmt.Errorf("unknown compression format: %s", format)
	}
}

type CompressModule struct {
	bc     *types.BuildConfig
	config *ModConfig
}
type ModConfig struct {
	Module  string `yaml:"module" validate:"required"`
	Fts     string `yaml:"format" validate:"required"`
	Sformat compressionFormat
}

func (s *CompressModule) Configure(config *types.BuildConfig) error {
	s.bc = config
	cfg, err := types.GetModConfig[ModConfig](config, "compress")
	if err != nil {
		return err
	}

	if cfg.Module == "" {
		return fmt.Errorf("compress module requires input module")
	}

	if cfg.Fts == "" {
		return fmt.Errorf("compress module requires format")
	}

	ft, err := parseCompressionFormat(cfg.Fts)
	if err != nil {
		return err
	}

	cfg.Sformat = ft
	s.config = cfg
	return nil
}

func (s *CompressModule) RunModule(modLogger *log.Logger, target types.Target) bool {
	ml := modLogger.ChildLogger("compress")

	objDir := filepath.Join(target.TempDir(), s.config.Module)
	dE, err := os.ReadDir(objDir)
	if err != nil {
		log.Logln(log.Error, err.Error())
		return false
	}
	outDir := filepath.Join(target.TempDir(), "compress")
	err = os.MkdirAll(outDir, 0755)
	if err != nil {
		log.Logln(log.Error, err.Error())
		return false
	}

	for _, entry := range dE {
		err = s.compressTarget(ml, entry, objDir, outDir)
		if err != nil {
			log.Logln(log.Error, err.Error())
			return false
		}
	}

	return true
}

func (s *CompressModule) compressTarget(ml *log.Logger, entry fs.DirEntry, objDir, oDir string) error {
	name := strings.Split(entry.Name(), ".")[0]
	ml.Logf(log.Info, "Compressing %s", name)
	compressed := bytes.NewBuffer(nil)
	iPath := filepath.Join(objDir, entry.Name())

	var err error

	switch s.config.Sformat {
	case cfTarGz:
		err = tarGzCompress(iPath, compressed)
	case cfZip:
		err = zipCompress(iPath, compressed)
	}

	if err != nil {
		return err
	}

	oPath := filepath.Join(oDir, name+"."+string(s.config.Sformat))
	f, err := os.Create(oPath)
	if err != nil {
		return err
	}

	_, err = io.Copy(f, compressed)
	if err != nil {
		return err
	}
	f.Close()
	ml.Logf(log.Info, "Compressed %s", name)
	return nil
}

func (s *CompressModule) Name() string {
	return "compress"
}

func (s *CompressModule) Requires() []string {
	return []string{s.config.Module}
}

func (c *CompressModule) OnFail() error {
	return nil
}

func tarGzCompress(src string, buf io.Writer) error {
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

func zipCompress(src string, buf io.Writer) error {
	w := zip.NewWriter(buf)

	err := filepath.Walk(src, func(file string, fi fs.FileInfo, _ error) error {
		header, err := zip.FileInfoHeader(fi)
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

		f, err := w.CreateHeader(header)
		if err != nil {
			return err
		}

		if !fi.IsDir() {
			data, err := os.Open(file)
			if err != nil {
				return err
			}
			io.Copy(f, data)
			data.Close()
		}
		return nil
	})
	if err != nil {
		return err
	}

	if err := w.Close(); err != nil {
		return err
	}
	return nil
}

func (c *CompressModule) TargetAgnostic() bool {
	return false
}

func (*CompressModule) RunOnCached() bool {
	return false 
}
