package runner

import (
	"bytes"
	"fmt"
	"slices"
	"strings"

	"github.com/lspaccatrosi16/go-cli-tools/args"
	"github.com/lspaccatrosi16/lbt/lib/log"
	"github.com/lspaccatrosi16/lbt/lib/modules"
	"github.com/lspaccatrosi16/lbt/lib/progress"
	"github.com/lspaccatrosi16/lbt/lib/types"
)

func RunModules(config *types.BuildConfig, mainMods map[string]types.Module, cached bool) error {
	ml := log.Default.ChildLogger("build")

	job := progress.NewJob("lbt")
	preHooksJob := job.NewChild("pre-build")
	for _, modName := range modules.PreOrder {
		mod := modules.Pre[modName]
		if !cached || mod.RunOnCached() {
			preHooksJob.NewChild(mod.Name()).WithFunc(mod.RunModule).WithConfigure(WrapConfig(mod.Configure, config)).WithLog(ml)
		}
	}
	tls := map[types.Target]*bytes.Buffer{}

	if len(config.Modules) > 0 {
		targFilter, err := args.GetFlagValue[string]("targFilter")
		if err != nil {
			return err
		}

		filters := strings.Split(targFilter, ",")
		for i := range filters {
			filters[i] = strings.TrimSpace(filters[i])
		}
		mainJob := job.NewChild("build").WithParallel()

		order := []string{}
		for _, m := range config.Modules {
			order, err = orderModules(config, m.Name, order, mainMods)
			if err != nil {
				return err
			}
		}

		for _, targ := range config.Targets {
			if targFilter == "" || slices.Contains(filters, targ.String()) {
				buf := bytes.NewBuffer(nil)
				tl := ml.ChildLogger(targ.String()).OverrideWriter(buf)
				tls[targ] = buf
				tg := mainJob.NewChild(targ.String()).WithLog(tl)
				for _, modName := range order {
					mod := mainMods[modName]
					if !cached || mod.RunOnCached() {
						tg.NewChild(modName).WithFunc(mod.RunModule).WithTarget(targ).WithLog(tl)
					}
				}
			}
		}
	}

	cleanupJob := progress.NewJob("post-build")
	nc, err := args.GetFlagValue[bool]("nc")
	if err != nil {
		return err
	}

	if nc {
		fmt.Println(types.NoTarget.TempDir())
	}

	for _, modName := range modules.PostOrder {
		mod := modules.Post[modName]
		if (!cached || mod.RunOnCached()) && !nc {
			cleanupJob.NewChild(mod.Name()).WithFunc(mod.RunModule).WithConfigure(WrapConfig(mod.Configure, config)).WithLog(ml)
		}
	}

	progress := progress.
		NewProgress(job, cleanupJob)
	res := progress.Render(fmt.Sprintf("build %s", config.Name))

	for t, b := range tls {
		if b.String() == "" {
			continue
		}
		fmt.Printf("[%s]\n", t.String())
		fmt.Println(b.String())
	}

	if !res {
		return fmt.Errorf("tasks encountered errors")
	}

	return nil
}

var modsSeen = map[string]bool{}

func orderModules(config *types.BuildConfig, modName string, order []string, mainMods map[string]types.Module) ([]string, error) {
	if ok := modsSeen[modName]; ok {
		return nil, fmt.Errorf("requirement cycle detected around module %s", modName)
	}

	modsSeen[modName] = true

	mod, ok := mainMods[modName]
	if !ok {
		return nil, fmt.Errorf("module %s was specified, but could not be found", modName)
	}

	err := mod.Configure(config)
	if err != nil {
		return nil, err
	}

	requirements := mod.Requires()
	if len(requirements) == 0 {
		return append([]string{modName}, order...), nil
	} else {
		for _, req := range requirements {
			inOrder := slices.Contains(order, req)
			if !inOrder {
				withReq, err := orderModules(config, req, order, mainMods)
				if err != nil {
					return nil, err
				}
				order = withReq
			}
		}
		order = append(order, modName)
	}
	return order, nil
}

func WrapConfig(configurer func(*types.BuildConfig) error, config *types.BuildConfig) func() error {
	return func() error {
		return configurer(config)
	}
}
