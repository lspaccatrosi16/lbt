package util

import (
	"io/fs"
	"path/filepath"
	"strings"
)


func ScanDir(startDir, suffix string) ([]string, error) {
	files := []string{}
	err := filepath.WalkDir(startDir, func(path string, d fs.DirEntry, err error) error {
		if d == nil {
			panic("dir entry is nil. Path: " + path)
		}
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
