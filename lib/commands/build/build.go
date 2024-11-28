package build

import (
	"path/filepath"
	"time"

	"github.com/lspaccatrosi16/go-cli-tools/args"
	"github.com/lspaccatrosi16/lbt/lib/cache"
	"github.com/lspaccatrosi16/lbt/lib/config"
	"github.com/lspaccatrosi16/lbt/lib/modules"
	"github.com/lspaccatrosi16/lbt/lib/modules/cached"
	"github.com/lspaccatrosi16/lbt/lib/modules/output"
	"github.com/lspaccatrosi16/lbt/lib/runner"
	"github.com/lspaccatrosi16/lbt/lib/types"
	"github.com/lspaccatrosi16/lbt/lib/util"
)

func Run() error {
	config, err := config.ParseConfig()
	if err != nil {
		return err
	}

	buildMeta := cache.BuildMeta{
		BuildName: config.Name,
		BuildTime: time.Now().Unix(),
	}

	modList := modules.Main
	force, err := args.GetFlagValue[bool]("force")
	if err != nil {
		return err
	}

	var usesCache bool

	if len(config.IncludeDirs) > 0 {
		buildMeta.Hash, err = cache.HashDirectories(config, config.IncludeDirs)
		if err != nil {
			return err
		}

		prevMeta, err := cache.GetLatestBuildArtifact(config.Name)
		if err != nil {
			return err
		}

		oCfg, oErr := types.GetModConfig[output.ModuleConfig](config, "output")
		if prevMeta != nil && prevMeta.Hash == buildMeta.Hash && oErr == nil && !force {
			modList = map[string]types.Module{
				"getCached": &cached.GetCachedModule{Meta: prevMeta},
				"output": &output.OutputModule{},
			}
			config.Modules = []types.ModuleConfig{
				{Name: "getCached", Config: map[string]interface{}{}},
				{Name: "output", Config: map[string]interface{}{"module": "getCached", "outDir": oCfg.OutDir}},
			}
			buildMeta = *prevMeta
			usesCache = true
		}
	}

	err = runner.RunModules(config, modList, usesCache)
	if err != nil {
		return err
	}

	cd, err := cache.GetArtifactCacheDir(buildMeta.BuildName)
	if err != nil {
		return err
	}

	pName := []string{}
	for _, p := range config.Produced {
		name := filepath.Base(p)
		pName = append(pName, name)
		err = util.Copy(filepath.Join(cd, name), p)
		if err != nil {
			return err
		}
	}

	buildMeta.Objects = pName

	err = cache.WriteBuildMeta(buildMeta)
	return err
}
