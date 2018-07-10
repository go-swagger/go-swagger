# Generate a spec from source code

```
Usage:
  swagger [OPTIONS] generate spec [spec-OPTIONS]

generate a swagger spec document from a go application

Application Options:
  -q, --quiet               silence logs
  -o, --output=LOG-FILE     redirect logs to file

Help Options:
  -h, --help                Show this help message

[spec command options]
      -b, --base-path=      the base path to use (default: .)
      -t, --tags=           build tags
      -m, --scan-models     includes models that were annotated with 'swagger:model'
          --compact         when present, doesn't prettify the json
      -o, --output=         the file to write to
      -i, --input=          the file to use as input
```

See code annotation rules [here](../use/spec.md)
