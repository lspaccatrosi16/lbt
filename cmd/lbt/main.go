package main

import (
	_ "embed"
	"os"
	"path/filepath"

	"github.com/lspaccatrosi16/go-cli-tools/args"
	"github.com/lspaccatrosi16/lbt/lib/config"
	"github.com/lspaccatrosi16/lbt/lib/log"
	"github.com/lspaccatrosi16/lbt/lib/modules"
	"github.com/lspaccatrosi16/lbt/lib/runner"
)

//go:embed version
var version string

func setup() error {
	args.RegisterEntry(args.NewStringEntry("config", "c", "config file", "lbt.yaml"))
	args.RegisterEntry(args.NewStringEntry("logLevel", "l", "log level", "info"))
	args.SetVersion(version)

	return args.ParseOpts()
}

func main() {
	err := setup()
	if err != nil {
		log.Fatalln(err)
	}

	ll, err := args.GetFlagValue[string]("logLevel")
	if err != nil {
		log.Fatalln(err)
	}
	logLev, err := log.ParseLogLevel(ll)
	if err != nil {
		log.Fatalln(err)
	}

	log.SetLogLevel(logLev)

	config, err := config.ParseConfig()
	if err != nil {
		log.Fatalln(err)
	}

	err = runner.RunModules(config, modules.List)
	if err != nil {
		log.Fatalln(err)
	}

	err = os.RemoveAll(filepath.Join(config.Cwd, "tmp"))
	if err != nil {
		log.Fatalln(err)
	}
}
