package create

import _ "embed"

//go:embed build.yaml
var buildTemplate string

//go:embed output.yaml
var outputTemplate string

//go:embed static.yaml
var staticTemplate string

//go:embed version.yaml
var versionTemplate string

var templates = map[string]string{
	"build":   buildTemplate,
	"output":  outputTemplate,
	"static":  staticTemplate,
	"version": versionTemplate,
}
