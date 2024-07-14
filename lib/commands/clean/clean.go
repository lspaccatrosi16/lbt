package clean

import (
	"os"
	"path/filepath"

	"github.com/lspaccatrosi16/lbt/lib/cache"
)

func Run() error {
	cd, err := cache.GetArtifactCacheDir("foo")
	if err != nil {
		return err
	}

	err = os.RemoveAll(filepath.Dir(cd))
	return err
}
