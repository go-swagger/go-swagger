---
title: swagger expand
date: 2023-01-01T01:01:01-08:00
draft: true
weight: 20
---
# Expand a swagger spec

The toolkit has a command to expand a swagger specification.

Expanding a specification resolve all `$ref` (remote or local) and replace them by their expanded
content in the main spec document.

### Usage

To expand a specification:

```
Usage:
  swagger [OPTIONS] expand [expand-OPTIONS]

expands the $refs in a swagger document to inline schemas

Application Options:
  -q, --quiet                     silence logs
      --log-output=LOG-FILE       redirect logs to file

Help Options:
  -h, --help                      Show this help message

[expand command options]
          --compact               applies to JSON formatted specs. When present, doesn't prettify the json
      -o, --output=               the file to write to
          --format=[yaml|json]    the format for the spec document (default: json)
```
