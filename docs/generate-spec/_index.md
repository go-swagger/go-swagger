---
title: Generate spec
date: 2023-01-01T01:01:01-08:00
draft: true
bookCollapseSection: true
weight: 70
---
# Generate a spec from source code

```
Usage:
  swagger [OPTIONS] generate spec [spec-OPTIONS]

generate a swagger spec document from a go application

Application Options:
  -q, --quiet                            silence logs
      --log-output=LOG-FILE              redirect logs to file

Help Options:
  -h, --help                             Show this help message

[spec command options]
      -w, --work-dir=                    the base path to use (default: .)
      -t, --tags=                        build tags
      -m, --scan-models                  includes models that were annotated with 'swagger:model'
          --compact                      when present, doesn't prettify the json
      -o, --output=                      the file to write to
      -i, --input=                       an input swagger file with which to merge
      -c, --include=                     include packages matching pattern
      -x, --exclude=                     exclude packages matching pattern
          --include-tag=                 include routes having specified tags (can be specified many times)
          --exclude-tag=                 exclude routes having specified tags (can be specified many times)
          --exclude-deps                 exclude all dependencies of project
      -n, --nullable-pointers            set x-nullable extension to true automatically for fields of pointer types without 'omitempty'
      -r, --ref-aliases                  transform aliased types into $ref rather than expanding their definition
          --transparent-aliases          treat type aliases as completely transparent, never creating definitions for them
          --skip-extensions              skip generation of x-go-* go-swagger extensions
          --skip-enum-desc               controls whether descriptions of enum values in field are preserved in the main description
          --allow-desc-with-ref          allow descriptions to flow alongside $ref
          --format=[yaml|json]           the format for the spec document (default: json)
          --emit-x-go-type               controls whether special extension x-go-type is emitted
          --emit-hierarchical-defs       controls how name conflicts are handled - this enables the last resort, failsafe method using nested definitions
          --single-line-comment-desc     controls how single line comments are handled. Default (false): as title. When true, title is skipped and only description is hydrated
          --enable-allof-compounding     controls compounded validations & descriptions with $ref. Default is to drop. When enabled, construct a allOf compound that preserves all siblings
          --default-allof-embeds         render plain (untagged) struct embeds as allOf composition instead of inlining their properties
          --name-from-tag=               ordered list of struct tag types consulted to derive property names, e.g. 'form' then 'json' (can be specified many times); defaults to 'json'
          --skip-jsonify-methods         emit interface method names verbatim, skipping the auto-jsonify (ToJSONName) mangler
          --name-concat-budget=          readability cutoff in [0,1] for concatenating package segments when deconflicting colliding definition names; 0 selects the built-in default (0.65)
          --after-decl-comments          allow swagger annotations inside a declaration body (leading comment of a struct body) or as a trailing inline comment
          --clean-godoc                  rewrite godoc-specific syntax (doc-link brackets, reference-style link definitions) when carried from a Go doc comment into the spec
          --prune                        with --scan-models, drop discovered definitions not transitively referenced from a path, response, parameter or input spec
          --colorized                    enable colorized diagnostics on stderr
      -q, --quiet                        mute diagnostics on stderr
```

See code annotation rules [here](../reference/annotations)

See also the [complete documentation of the `codescan` package](https://go-openapi.github.io/codescan/index.html) for additional guidance, tutorials and examples.
