---
title: swagger mixin
date: 2023-01-01T01:01:01-08:00
draft: true
weight: 50
---
# Mixin several swagger specs

The toolkit has a command to mix several swagger specifications into one.

Mixin merges several specs into the first (primary) spec given, and issues warnings when conflicts
are detected.

### Usage

To mixin several specifications:

```
Usage:
  swagger [OPTIONS] mixin [mixin-OPTIONS]

merge additional specs into first/primary spec by copying their paths and definitions

Application Options:
  -q, --quiet                     silence logs
      --log-output=LOG-FILE       redirect logs to file

Help Options:
  -h, --help                      Show this help message

[mixin command options]
      -c=                         expected # of rejected mixin paths, defs, etc due to existing key.
                                  Non-zero exit if does not match actual.
          --compact               applies to JSON formatted specs. When present, doesn't prettify
                                  the json
      -o, --output=               the file to write to
          --keep-spec-order       keep schema properties order identical to spec file
                                  (applies to schemas only, not to paths)
          --format=[yaml|json]    the format for the spec document (default: json)
          --ignore-conflicts      ignore conflict
```

### How merging works

The first argument is the primary spec. Subsequent arguments are mixins, in decreasing priority
order.

* On any collision, the **primary always wins**; among mixins, the one given **earliest** wins.
* Top-level scalar fields (`Info`, `BasePath`, `Host`, `ExternalDocs`) on the primary are
  **filled from the first mixin that provides a value**, but only when the primary's field is
  empty.
* `paths`, `definitions`, `parameters`, `responses`, `securityDefinitions`, `tags`, `security`
  and extensions are merged entry by entry. Duplicate keys (or equal security requirements, or
  equal tag names) are skipped with a warning.
* `schemes`, `consumes` and `produces` are merged as the union of distinct values; duplicates are
  silently dropped.
* Operation-id collisions are auto-resolved by appending `Mixin<N>` to the mixin operation id, so
  the merged spec keeps unique operation ids.

Example: if `primary.yaml` has `host: a.example.com` and `mixin.yaml` has `host: b.example.com`,
the output keeps `host: a.example.com`. If `primary.yaml` has no `host`, the output uses
`host: b.example.com` from the mixin.

### Limitations

#### YAML anchors are not preserved

YAML anchors (`&name` and `*name`) are resolved by the YAML parser into a fully expanded object
tree **before** mixin sees the document. The merged output has no anchor information, and
cross-file anchors are not legal YAML in the first place.

For type re-use across files, use [JSON Reference](https://datatracker.ietf.org/doc/html/draft-pbryan-zyp-json-ref-03)
(`$ref`) rather than YAML anchors:

```yaml
# common.yaml
definitions:
  Severity:
    type: string
    enum: [low, medium, high]
```

```yaml
# main.yaml
definitions:
  Issue:
    type: object
    properties:
      severity:
        $ref: "./common.yaml#/definitions/Severity"
```

For a literal concatenation of YAML files (no semantic merge, no conflict handling), `cat` is the
simpler tool.

#### Path and operation order in the output is alphabetical

The merged spec stores paths and definitions in Go maps, which serialize with alphabetically
sorted keys. Source-file order is not preserved, and there is no option to reorder by tag, by
source priority, or otherwise. This is an architectural constraint inherited from the underlying
spec model (`spec.Paths.Paths` is a map keyed by path string).

Note that the `--keep-spec-order` flag applies to **schema properties** order only, not to paths
or definitions.

If you need a specific non-alphabetical order, the only workaround today is a post-processing step
on the YAML or JSON output.

#### Empty response descriptions

Unmarshalling responses from JSON can produce empty `description` fields where the source had one.
Mixin internally calls a fixer to repopulate them where possible. If you observe missing response
descriptions in the output, double-check they are present and non-empty in your source files.
