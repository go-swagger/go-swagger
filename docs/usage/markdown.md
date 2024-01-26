---
title: swagger generate markdown
date: 2023-01-01T01:01:01-08:00
draft: true
weight: 40
---
# Generate markdown documentation

This is a command to generate a markdown document from a swagger spec.

The generated doc is no substitute for advanced swagger documentation tools such as redoc:
it provides a simple documentation for your API.

The spec is canonicalized just like for code generation: the generated markdown represents
operations and models just like your generated code sees them.

The spec is flattened to be rendered as a self-contained document and all complex inlined models are
defined as standalone models (documented as "inlined schemas").

Known limitations:
* validations are not rendered, for the sake of brevity

### Usage

```
Usage:
  swagger [OPTIONS] generate markdown [markdown-OPTIONS]

generate a markdown representation from the swagger spec

Application Options:
  -q, --quiet                                                                     silence logs
      --log-output=LOG-FILE                                                       redirect logs to file

Help Options:
  -h, --help                                                                      Show this help message

[markdown command options]
          --output=                                                               the file to write the generated markdown. (default: markdown.md)

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
          --dump-data                                                             when present dumps the json for the template generator instead of generating
                                                                                  files
          --strict-responders                                                     Use strict type for the handler return value
          --with-expand                                                           expands all $ref's in spec prior to generation (shorthand to
                                                                                  --with-flatten=expand)
          --with-flatten=[minimal|full|expand|verbose|noverbose|remove-unused]    flattens all $ref's in spec prior to generation (default: minimal, verbose)

    Options for model generation:
      -m, --model-package=                                                        the package to save the models (default: models)
      -M, --model=                                                                specify a model to include in generation, repeat for multiple (defaults to
                                                                                  all)
          --existing-models=                                                      use pre-generated models e.g. github.com/foobar/model
          --strict-additional-properties                                          disallow extra properties when additionalProperties is set to false
          --keep-spec-order                                                       keep schema properties order identical to spec file
          --struct-tags=                                                          the struct tags to generate, repeat for multiple (defaults to json)

    Options for operation generation:
      -O, --operation=                                                            specify an operation to include, repeat for multiple (defaults to all)
          --tags=                                                                 the tags to include, if not specified defaults to all
      -a, --api-package=                                                          the package to save the operations (default: operations)
          --with-enum-ci                                                          allow case-insensitive enumerations
          --skip-tag-packages                                                     skips the generation of tag-based operation packages, resulting in a flat generation
```
