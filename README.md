# LBT
> A powerful build tool for golang projects

---

1. Generate a new project with

```shell
lbt create
```

2. Select the modules you want
3. Build

```shell
lbt -c <your config path>.yaml
```

---

## Modules

### Build

The core build module of `lbt`. Configure commands (e.g. each main.go that you want to be produced into a binary).
Add build targets doubles in the `targets` field. If you want the version module to be used, the `version` field should be enabled.
Add optional compiler flags in the `ldflags` field.

> The currently supported os are `linux`, `darwin`, `windows`  
> The currently supported arch are `amd64`, `i386`, `arm64`, `arm`

### Output
Writes the output of the specified `module` to the `outDir`. If this is not used, all generated files will be deleted after the build is finished.

### Static
Generates an archive including static files. Specify the directory of the static files with `structure`. The `exePath` is a path relative to structure in which the generated executable should be placed.

### Version
Updates a plaintext file with a version string, which can be included into the executable with a `//go:embed` tag. To ensure that this is executed before the executable is built, the build module's `version` tag should be enabled.
Specify the version file relative to the config file of the project with `path`.
The `type` field can be one of:
- `buildint` - Increments an integer in the file on each build
- `buildstr` - Generates a unique string that can be used to identify the build.
- `semver` - Increments a 4th component in a version field, that corresponds to the build number e.g. `x.x.x.40` would be the 40th build for version `x.x.x`

---

## Licence

See [Licence](./LICENCE)