# LBT
> A powerful build tool for projects written in golang, java, and c.

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

> The currently supported `os` are `linux`, `darwin`, `windows`, `jvm`, `android`
> The currently supported `arch` are `amd64`, `i386`, `arm64`, `arm`

## Modules

### GoBuild

The build module of `lbt` for golang. Configure commands (e.g. each main.go that you want to be produced into a binary).

#### GoBuild Module Config

| Name | Type | Description |
| ---- | ---- | ----------- |
| `commands` | []{`name`: string, `path`: string} | A list of commands, which are objects that contain a `name` and a `path`, which points to the `main.go` file. |
| `ldflags` | string | Any flags that should be passed to the `go build` command in the `-ldflags` argument. |
| `cgOff` | boolean | Disables CGO for the build (useful if you want to make sure your program is statically linked). |
| `root` | string | The root directory that the go build commands will be run from. |

### JavaBuild

The build module of `lbt` for java.

#### JavaBuild Module Config

| Name | Type | Description |
| ---- | ---- | ----------- |
| `main` | string | The name of the main class. |
| `dependencies` | []string | A list of paths to dependencies. Can be directories or jar files. |

### CBuild

The build module of `lbt` for C.

#### CBuild Module Config

| Name | Type | Description |
| ---- | ---- | ----------- |
| `name` | string | The name of the executable. |
| `source` | string | The path to the `src` directory, which is the main code directory of the project. |
| `include` | []string | A list of paths to the `include` directory for header files. |
| `compiler` | string | The name of the compiler to run e.g. `gcc` or `clang`, |
| `flags` | []string | A list of extra flags to pass the compiler. |
| `main` | string | The path to the `main.c` file. |
| `librarymode` | `shared` \| `static` | The type of library to build. |
| `cc` | boolean | Generate a `compile_commands.json` file (useful for clangd lsp). |
| `libs` | []string | List of libraries to include and compile against. |
| `libdirs` | []string | List of directories to search for 3rd party libraries. |

### OdinBuild

The build module of `lbt` for odin.

#### OdinBuild Module Comfig

| Name | Type | Description |
| ---- | ---- | ----------- |
| `src` | string | The path to the project's src folder. |
| `optimise` | string | The optimisation preset to use in odin compilation. |
| `debug` | boolean | Enables debug compilation. |
| `flags` | []string | Additional flags to be passed to the compiler. |


### VBuild

The build module of `lbt` for v.

#### OdinBuild Module Comfig

| Name | Type | Description |
| ---- | ---- | ----------- |
| `src` | string | The path to the project's src folder. |
| `backend` | string | The backend  to use in v compilation. |
| `debug` | boolean | Enables debug compilation. |
| `flags` | []string | Additional flags to be passed to the compiler. |
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
Updates a plaintext file with a version string, which can be included into the executable with a `//go:embed` tag. 

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
