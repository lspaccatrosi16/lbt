package modules

import (
	"github.com/lspaccatrosi16/lbt/lib/modules/build"
	"github.com/lspaccatrosi16/lbt/lib/modules/compress"
	"github.com/lspaccatrosi16/lbt/lib/modules/output"
	"github.com/lspaccatrosi16/lbt/lib/modules/static"
	"github.com/lspaccatrosi16/lbt/lib/modules/version"
	"github.com/lspaccatrosi16/lbt/lib/types"
)

var List = []types.Module{
	&build.BuildModule{},
	&output.OutputModule{},
	&version.VersionModule{},
	&static.StaticModule{},
	&compress.CompressModule{},
}
