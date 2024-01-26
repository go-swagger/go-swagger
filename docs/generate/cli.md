---
title: Generate a CLI client
date: 2023-01-01T01:01:01-08:00
draft: true
weight: 25
---
# Generate a CLI (command line tool) from swagger spec
This toolkit can generate a CLI to interact with your server

### Features of generated app
* auto-completion for bash, zsh, fish and powershell.
* use config file to specify common flags.
* each param and each field in body has a cli flag. etc.

### CLI usage
```
Usage:
  swagger [OPTIONS] generate cli [cli-OPTIONS]

generate a command line client tool from the swagger spec

Application Options:
  -q, --quiet                                                                     silence logs
      --log-output=LOG-FILE                                                       redirect logs to file

Help Options:
  -h, --help                                                                      Show this help message

[cli command options]
      -c, --client-package=                                                       the package to save the client specific code (default: client)
      -P, --principal=                                                            the model to use for the security principal
          --default-scheme=                                                       the default scheme for this API (default: http)
          --principal-is-interface                                                the security principal provided is an interface, not a struct
          --default-produces=                                                     the default mime type that API operations produce (default:
                                                                                  application/json)
          --default-consumes=                                                     the default mime type that API operations consume (default:
                                                                                  application/json)
          --skip-models                                                           no models will be generated when this flag is specified
          --skip-operations                                                       no operations will be generated when this flag is specified
      -A, --name=                                                                 the name of the application, defaults to a mangled value of info.title
          --cli-app-name=                                                         the app name for the cli executable. useful for go install. (default:
                                                                                  cli)

    Options common to all code generation commands:
      -f, --spec=                                                                 the spec file to use (default swagger.{json,yml,yaml})
      -t, --target=                                                               the base directory for generating the files (default: ./)
          --template=[stratoscale]                                                load contributed templates
      -T, --template-dir=                                                         alternative template override directory
      -C, --config-file=                                                          configuration file to use for overriding template options
      -r, --copyright-file=                                                       copyright file used to add copyright header
          --additional-initialism=                                                consecutive capitals that should be considered intialisms
          --allow-template-override                                               allows overriding protected templates
          --skip-validation                                                       skips validation of spec prior to generation
          --dump-data                                                             when present dumps the json for the template generator instead of
                                                                                  generating files
          --strict-responders                                                     Use strict type for the handler return value
          --with-expand                                                           expands all $ref's in spec prior to generation (shorthand to
                                                                                  --with-flatten=expand)
          --with-flatten=[minimal|full|expand|verbose|noverbose|remove-unused]    flattens all $ref's in spec prior to generation (default: minimal,
                                                                                  verbose)

    Options for model generation:
      -m, --model-package=                                                        the package to save the models (default: models)
      -M, --model=                                                                specify a model to include in generation, repeat for multiple
                                                                                  (defaults to all)
          --existing-models=                                                      use pre-generated models e.g. github.com/foobar/model
          --strict-additional-properties                                          disallow extra properties when additionalProperties is set to false
          --keep-spec-order                                                       keep schema properties order identical to spec file
          --struct-tags=                                                          the struct tags to generate, repeat for multiple (defaults to json)

    Options for operation generation:
      -O, --operation=                                                            specify an operation to include, repeat for multiple (defaults to all)
          --tags=                                                                 the tags to include, if not specified defaults to all
      -a, --api-package=                                                          the package to save the operations (default: operations)
          --with-enum-ci                                                          allow case-insensitive enumerations
          --skip-tag-packages                                                     skips the generation of tag-based operation packages, resulting in a
                                                                                  flat generation
```

### Build a CLI
There is an example cli and tutorial provided at: https://github.com/go-swagger/go-swagger/tree/master/examples/cli

To generate a CLI:
```
swagger generate cli -f [http-url|filepath] --cli-app-name [app-name]
```
Cli is a wrapper of generated client code (see [client](./client.md) for details), so all client generation options are honored.

To build the generated CLI code:
```
go build cmd/<app-name>/main.go 
```
Or install in your go/bin
```
go install cmd/<app-name>/main.go
```

See details of the generated app help message for usage
```
<app name> help
```

A more detailed/complicated example is generated CLI for docker engine: https://github.com/go-swagger/dockerctl
