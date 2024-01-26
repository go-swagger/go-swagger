---
title: Generating models
date: 2023-01-01T01:01:01-08:00
draft: true
weight: 40
---
# Generate models from a swagger spec


```
Usage:
  swagger [OPTIONS] generate model [model-OPTIONS]

generate one or more models from the swagger spec

Application Options:
  -q, --quiet                                                                     silence logs
      --log-output=LOG-FILE                                                       redirect logs to file

Help Options:
  -h, --help                                                                      Show this help message

[model command options]
      -n, --name=                                                                 the model to generate, repeat for multiple (defaults to all). Same as --models
          --accept-definitions-only                                               accepts a partial swagger spec wih only the definitions key

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
          --dump-data                                                             when present dumps the json for the template generator instead of generating files
          --strict-responders                                                     Use strict type for the handler return value
          --with-expand                                                           expands all $ref's in spec prior to generation (shorthand to --with-flatten=expand)
          --with-flatten=[minimal|full|expand|verbose|noverbose|remove-unused]    flattens all $ref's in spec prior to generation (default: minimal, verbose)

    Options for model generation:
      -m, --model-package=                                                        the package to save the models (default: models)
      -M, --model=                                                                specify a model to include in generation, repeat for multiple (defaults to all)
          --existing-models=                                                      use pre-generated models e.g. github.com/foobar/model
          --strict-additional-properties                                          disallow extra properties when additionalProperties is set to false
          --keep-spec-order                                                       keep schema properties order identical to spec file
          --struct-tags=                                                          the struct tags to generate, repeat for multiple (defaults to json)
          --rooted-error-path                                                     extends validation errors with the type name instead of an empty path, in the case of arrays and maps
```

Schema generation rules are detailed [here](../reference/models/schemas.md).
