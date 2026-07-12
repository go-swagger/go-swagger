---
title: swagger generate
date: 2023-01-01T01:01:01-08:00
draft: true
weight: 35 
---
# Generating

The toolkit has a command to generate different targets.

### Usage

Available subcommands.

```
Usage:
  swagger [OPTIONS] generate <command>

generate go code for the swagger spec file

Application Options:
  -q, --quiet                  silence logs
      --log-output=LOG-FILE    redirect logs to file

Help Options:
  -h, --help                   Show this help message

Available commands:
  cli        generate a command line client tool from the swagger spec
  client     generate all the files for a client library
  markdown   generate a markdown representation from the swagger spec
  model      generate one or more models from the swagger spec
  operation  generate one or more server operations from the swagger spec
  server     generate all the files for a server application
  spec       generate a swagger spec document from a go application
  support    generate supporting files like the main function and the api builder
```

For code generation targets (`cli`, `client`, `model`, `operation`, `server`, `support`), read more [here](../generate/).

For spec generation targets (`spec`), read more [there](../generate-spec/).

For markdown generation target (`markdown`), read [this](markdown.md).

## Troubleshooting generate and validate

Use this checklist when `swagger validate` succeeds locally but generation fails, or when a spec path works in one shell and not another.

| Symptom | What to check | Suggested command or fix |
| --- | --- | --- |
| `swagger validate` reports that the spec cannot be loaded | Confirm the spec path is relative to the directory where you run `swagger` | Run `pwd` and then use an explicit path, for example `swagger validate ./api/swagger.yml` |
| `$ref` targets cannot be resolved | Check whether each local or remote `$ref` is a valid absolute reference, or is relative to the containing spec file | Validate the same file you pass to generation: `swagger validate ./api/swagger.yml` |
| Generation writes files in an unexpected location | Check `--target`; generated packages are created under that directory, and the target directory must exist before generation | Create the target first: `mkdir -p gen`, then run `swagger generate server -f ./api/swagger.yml --target ./gen` |
| Generated code does not compile after a successful run | Ensure the generated target is a Go module that can resolve the generated dependencies | In a new target, run `go mod init` and `go mod tidy` after generation |
| `goimports` or formatting-related failures appear | Verify Go is installed and available on `PATH` in the same shell running `swagger` | Run `go version` before retrying generation |
| Operations or models are missing | Check that the spec validates and that flags such as `--operation`, `--model`, `--include-tag`, or `--exclude-tag` match the spec exactly | Start without selection flags, then add them back one at a time |
| The shell prints only `zsh: killed` or `Killed` after generation starts | The OS terminated the process; two common causes are (1) CGO is enabled in a source build but not in the release binary, or (2) the process was killed by the OOM killer or a resource limit | See the section below |

### `zsh: killed` or `Killed` during generation

The official release binaries are built with `CGO_ENABLED=0`.
When you build from source the Go toolchain enables CGO by default, even though go-swagger does not use any CGO code.
The resulting dynamically-linked binary can fail to run on some macOS versions and configurations where the statically-linked release binary works correctly.
Rebuild with CGO disabled to match the release behaviour:

```sh
CGO_ENABLED=0 go build -o swagger ./cmd/swagger
```

If the `SIGKILL` occurs on Linux or inside a container it is more likely a kernel OOM event or a configured resource limit (`ulimit`, container memory cap).
Check `dmesg` for OOM evidence and compare `ulimit -a` with the memory available during the run.

When reporting this problem, include the exact `swagger generate` command, `swagger version`, operating system, whether the binary was installed from a release or built from source, and any relevant evidence from the system logs.

A minimal validation-first flow is:

```sh
swagger validate ./api/swagger.yml
mkdir -p gen
swagger generate server -f ./api/swagger.yml --target ./gen
cd ./gen
go mod init example.com/generated-api
go mod tidy
go test ./...
```

When reporting an issue, include the exact `swagger` command, the path to the spec as passed with `-f`, the current working directory, and the validation output. Those details usually determine whether the failure is a spec problem, a path problem, or a generated-code setup problem.
