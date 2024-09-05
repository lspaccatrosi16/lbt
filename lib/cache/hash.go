package cache

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/lspaccatrosi16/lbt/lib/types"
)

func HashDirectories(bc *types.BuildConfig, dirs []string) (string, error) {
	tHashes := ""

	var vfile string
	if bc.Version.Path != "" {
		vfile = filepath.Join(bc.Cwd, bc.Version.Path)
	}

	for _, dir := range dirs {
		dHash := ""
		err := filepath.WalkDir(filepath.Join(bc.Cwd, dir), func(path string, d fs.DirEntry, err error) error {
			if !d.IsDir() && (vfile == "" || vfile != path){
				pHash := hash([]byte(path))
				di, err := d.Info()
				if err != nil {
					return err
				}
				dsc := fmt.Sprintf("%d", di.Size())
				iHash := hash([]byte(dsc))

				dHash += pHash + iHash
			}
			return nil
		})
		if err != nil {
			return "", err
		}
		tHashes += hash([]byte(dHash))
	}
	return hash([]byte(tHashes)), nil
}

func hash(data []byte) string {
	hasher := sha1.New()
	hasher.Write(data)
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}
