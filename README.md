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

## Base Config

A `lbt.yaml` defines the build process.
| Name | Type | Description |
| ---- | ---- | ----------- |
| `name` | string | The name of the program. |
| `targets` | []{`os`: string; `arch`: string} | A list of build targets. |
| `modules` | {name: string, config: moduleConfig} | A list of all modules used, and their respective configurations. |
| `includeDirs` | []string | A list of directories to watch for file changes. |

> The currently supported `os` are `linux`, `darwin`, `windows`  
> The currently supported `arch` are `amd64`, `i386`, `arm64`, `arm`

## Modules

### Build

The core build module of `lbt`. Configure commands (e.g. each main.go that you want to be produced into a binary).

#### Build Module Config

| Name | Type | Description |
| ---- | ---- | ----------- |
| `commands` | {`name`: string, `path`: string} | A list of commands, which are objects that contain a `name` and a `path`, which points to the `main.go` file. |
| `version` | boolean | Whether the version module is being used, to ensure that the `version` module is run before this.
| `ldflags` | string | Any flags that should be passed to the `go build` command in the `-ldflags` argument. |
| `cgOff` | boolean | Disables CGO for the build (useful if you want to make sure your program is statically linked). |

### Output
Writes built objects to a given output directory.

#### Output Module Config

| Name | Type | Description |
| ---- | ---- | ----------- |
| `module` | string | The module's output that will be outputed. |
| `outDir` | string | The path relative to the config file that the outputted files will be placed in.  |

### Static
Generates file structures based on a static template, inserting built assets into the structure. 

#### Static Module Config

| Name | Type | Description |
| ---- | ---- | ----------- |
| `structures` | structure | A list of `structure` objects. |

#### Structure Config

| Name | Type | Description |
| ---- | ---- | ----------- |
| `name` | string | The name of the structure. |
| `path` | string | The path of the structure's static template. |
| `executables` | []{`command`: string, `path`: string} | A list of executables, where `command` is the name of the built command, and `path` is the path relative to the root of the structure where it should be inserted. |

### Compress

Compresses built objects into archives.

#### Compress Module Config

| Name | Type | Description |
| ---- | ---- | ----------- |
| `module` | string | The module's output that will be compressed. |
| `format` | string | The compression format to use. |

> The currently supported `format` are `tar.gz` and `zip`

### Version
Updates a plaintext file with a version string, which can be included into the executable with a `//go:embed` tag. To ensure that this is executed before the executable is built, the build module's `version` setting should be enabled.

#### Version Module Config

| Name | Type | Description |
| ---- | ---- | ----------- |
| `path` | string | The path of the version file relative to the project config file |
| `type` | string | The versioning type to use |

#### Versioning Types

| Type | Description |
| ---- | ----------- | 
| `buildint` | Increments an integer in the file on each build |
| `buildstr` | Generates a unique string that can be used to identify the build. |
| `semver` | Increments a 4th component in a version field, that corresponds to the build number e.g. `x.x.x.40` would be the 40th build for version `x.x.x` |


---

## Licence

See [Licence](./LICENCE)