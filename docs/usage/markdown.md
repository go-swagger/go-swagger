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

Since this command generates no go code, it only takes the options that bear on the document it
produces. The options that shape go source, such as `--struct-tags` or `--strict-responders`, are
not available here, and neither are the ones naming the packages the code is written to
(`--model-package`, `--api-package`): the "Go type" column always reports the conventional
`models` and `operations` layout.

### Usage

```
Usage:
  swagger [OPTIONS] generate markdown [markdown-OPTIONS] [spec]

generate a markdown representation from the swagger spec

Application Options:
  -q, --quiet                                                                                silence logs
      --log-output=LOG-FILE                                                                  redirect logs to file

Help Options:
  -h, --help                                                                                 Show this help message

[markdown command options]
          --output=                                                                          the file to write the generated markdown. (default:
                                                                                             markdown.md)

    Options for reading the spec and writing the documentation:
          --with-expand                                                                      expands all $ref's in the spec (shorthand to
                                                                                             --with-flatten=expand)
          --with-flatten=[minimal|full|expand|verbose|noverbose|remove-unused|keep-names]    flattens all $ref's in the spec (default: minimal,
                                                                                             verbose)
      -f, --spec=                                                                            the spec file to use (default swagger.{json,yml,yaml})
          --skip-validation                                                                  skips validation of spec prior to generation
          --restricted                                                                       Use restricted http client for remote $ref
          --rooted=                                                                          Local $ref resolution contained relative to root FS
      -t, --target=                                                                          the base directory for generating the files (default: ./)
      -T, --template-dir=                                                                    alternative template override directory
      -C, --config-file=                                                                     configuration file to use for overriding template options
          --additional-initialism=                                                           consecutive capitals that should be considered intialisms
          --allow-template-override                                                          allows overriding protected templates
          --dump-data                                                                        when present dumps the json for the template generator
                                                                                             instead of generating files
          --ensure-target                                                                    Create the target directory if it does not already exist
      -p, --template-plugin=                                                                 the template plugin to use

    Options for selecting the documented models:
      -M, --model=                                                                           specify a model to include in generation, repeat for
                                                                                             multiple (defaults to all)
          --keep-spec-order                                                                  keep schema properties order identical to spec file

    Options for selecting the documented operations:
      -O, --operation=                                                                       specify an operation to include, repeat for multiple
                                                                                             (defaults to all)
          --tags=                                                                            the tags to include, if not specified defaults to all
          --skip-tag-packages                                                                skips the generation of tag-based operation packages,
                                                                                             resulting in a flat generation
```

`--dump-data` dumps the data handed to the template instead of rendering it, which is how you find
out what a custom markdown template passed with `--template-dir` can use.
