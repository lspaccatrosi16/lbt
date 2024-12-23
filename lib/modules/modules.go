package modules

import (
	"github.com/lspaccatrosi16/lbt/lib/modules/cbuild"
	"github.com/lspaccatrosi16/lbt/lib/modules/cleanup"
	"github.com/lspaccatrosi16/lbt/lib/modules/compress"
	"github.com/lspaccatrosi16/lbt/lib/modules/gobuild"
	"github.com/lspaccatrosi16/lbt/lib/modules/javabuild"
	"github.com/lspaccatrosi16/lbt/lib/modules/odinbuild"
	"github.com/lspaccatrosi16/lbt/lib/modules/output"
	"github.com/lspaccatrosi16/lbt/lib/modules/setup"
	"github.com/lspaccatrosi16/lbt/lib/modules/static"
	"github.com/lspaccatrosi16/lbt/lib/modules/vbuild"
	"github.com/lspaccatrosi16/lbt/lib/modules/version"
	"github.com/lspaccatrosi16/lbt/lib/types"
)

var Pre = map[string]types.Module{
	"setup":   &setup.SetupModule{},
	"version": &version.VersionModule{},
}

var PreOrder = []string{
	"setup",
	"version",
}

var Main = map[string]types.Module{
	"gobuild":   &gobuild.GobuildModule{},
	"javabuild": &javabuild.JavabuildModule{},
	"cbuild":    &cbuild.CbuildModule{},
	"odinbuild": &odinbuild.OdinbuildModule{},
	"vbuild":    &vbuild.VbuildModule{},
	"output":    &output.OutputModule{},
	"static":    &static.StaticModule{},
	"compress":  &compress.CompressModule{},
}

var Post = map[string]types.Module{
	"cleanup": &cleanup.CleanupModule{},
}

var PostOrder = []string{
	"cleanup",
}
