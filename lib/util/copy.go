package util

import (
	"io"
	"os"
	"path/filepath"
)

func Copy(dst string, src string) error {
	s, err := os.Stat(src)
	if err != nil {
		return err
	}
	if s.IsDir() {
		return cpy_dir(dst, src)
	} else {
		return cpy_file(dst, src)

	}
}

func cpy_dir(dst, src string) error {
	de, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, d := range de {
		srcPath := filepath.Join(src, d.Name())
		dstPath := filepath.Join(dst, d.Name())
		if d.IsDir() {
			err = os.MkdirAll(dstPath, 0755)
			if err != nil {
				return err
			}
			return cpy_dir(dstPath, srcPath)
		} else {
			return cpy_file(dstPath, srcPath)
		}
	}
	return nil
}

func cpy_file(dstPath, srcPath string) error {
	srcF, err := os.Open(srcPath)
	if err != nil {
		return err
	}

	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	dst.Chmod(0755)
	io.Copy(dst, srcF)
	srcF.Close()
	dst.Close()
	return nil

}
