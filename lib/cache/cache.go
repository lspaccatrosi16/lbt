package cache

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type BuildMeta struct {
	BuildTime int64    `json:"build_time"`
	BuildName string   `json:"build_name"`
	Hash      string   `json:"hash"`
	Objects   []string `json:"objects"`
	location  string
}

func (b *BuildMeta) Location() string {
	return b.location
}

func getCacheDir() (string, error) {
	cd, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	d := filepath.Join(cd, "lbt_cache")
	err = os.MkdirAll(d, 0755)
	if err != nil {
		return "", err
	}
	return d, nil
}

func GetArtifactCacheDir(name string) (string, error) {
	cd, err := getCacheDir()
	if err != nil {
		return "", err
	}
	d := filepath.Join(cd, name)
	err = os.MkdirAll(d, 0755)
	if err != nil {
		return "", err
	}

	return d, nil
}

func GetLatestBuildArtifact(name string) (*BuildMeta, error) {
	cd, err := getCacheDir()
	if err != nil {
		return nil, err
	}
	mf := filepath.Join(cd, name, "meta.json")

	fd, err := os.ReadFile(mf)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	meta := BuildMeta{}
	err = json.Unmarshal(fd, &meta)
	if err != nil {
		return nil, err
	}
	meta.location = filepath.Dir(mf)
	return &meta, nil
}

func WriteBuildMeta(meta BuildMeta) error {
	cd, err := getCacheDir()
	if err != nil {
		return err
	}
	err = os.MkdirAll(filepath.Join(cd, meta.BuildName), 0755)
	if err != nil {
		return err
	}
	mf := filepath.Join(cd, meta.BuildName, "meta.json")

	fd, err := json.Marshal(meta)
	if err != nil {
		return err
	}

	err = os.WriteFile(mf, fd, 0644)
	if err != nil {
		return err
	}
	return nil
}
