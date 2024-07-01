package util

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

func Copy(dst string, src string) error {
	return filepath.Walk(filepath.Join(src), func(file string, fi fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
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
