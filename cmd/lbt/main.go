package main

import (
	_ "embed"

	"github.com/lspaccatrosi16/go-cli-tools/args"
	"github.com/lspaccatrosi16/lbt/lib/commands/build"
	"github.com/lspaccatrosi16/lbt/lib/commands/create"
	"github.com/lspaccatrosi16/lbt/lib/log"
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

	var cmd string
	a := args.GetArgs()

	if len(a) >= 1 {
		cmd = a[0]
	} else {
		cmd = "build"
	}

	switch cmd {
	case "build":
		err = build.Run()
	case "create":
		err = create.Run()
	default:
		log.Fatalf("Unknown command: %s", cmd)
	}

	if err != nil {
		log.Fatalln(err)
	}
}
