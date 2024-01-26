---
title: swagger flatten
date: 2023-01-01T01:01:01-08:00
draft: true
weight: 30
---
# Flatten a swagger spec

The toolkit has a command to flatten a swagger specification.

Flattening a specification bundles all remote `$ref` in the main spec document.

Depending on the selected flattening options, additional preprocessing may take place:

- full flattening: replacing all inline complex constructs by a named entry in #/definitions
- expand: replace all $ref's in the document by their expanded content
- remove-unused: remove unused definitions from the spec
- keep-names: does not attempt to jsonify names in definitions and keep them as-is. Generally not suited to use with "full".

The default behavior of flatten is to bundles remote refs into definitions and
normalize JSON pointers to definitions.

### Usage

To flatten a specification:

```
Usage:
  swagger [OPTIONS] flatten [flatten-OPTIONS]

expand the remote references in a spec and move inline schemas to definitions, after flattening there are no complex inlined anymore

Application Options:
  -q, --quiet                                                                                silence logs
      --log-output=LOG-FILE                                                                  redirect logs to file

Help Options:
  -h, --help                                                                                 Show this help message

[flatten command options]
          --compact                                                                          applies to JSON formatted specs. When present, doesn't prettify the json
      -o, --output=                                                                          the file to write to
          --format=[yaml|json]                                                               the format for the spec document (default: json)
          --with-expand                                                                      expands all $ref's in spec prior to generation (shorthand to --with-flatten=expand)
          --with-flatten=[minimal|full|expand|verbose|noverbose|remove-unused|keep-names]    flattens all $ref's in spec prior to generation (default: minimal, verbose)
```
