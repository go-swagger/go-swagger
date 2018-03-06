# Generate a data model from swagger spec

The toolkit allows for generating go native structures from a swagger compliant model definition.
The generated objects follow the specified validation rules, including extended format directives for strings and numbers.

Generated models support most JSON-schema draft4 features, including dependencies, AllOf, AdditionalProperties and AdditionalItems.

NOTE: AnyOf and OneOf constructs are not supported at the moment (use swagger polymorphism instead).

Schema generation rules are described [here](../use/schemas.md)

#### Minimal use of go's reflection

A generated model uses _no reflection_ except for enum and required validations.

#### Usage

```
Usage:
  swagger [OPTIONS] generate model [model-OPTIONS]

generate one or more models from the swagger spec

Application Options:
  -q, --quiet                         silence logs
  -o, --output=LOG-FILE               redirect logs to file

Help Options:
  -h, --help                          Show this help message

[model command options]
      -f, --spec=                     the spec file to use (default swagger.{json,yml,yaml})
      -a, --api-package=              the package to save the operations (default: operations)
      -m, --model-package=            the package to save the models (default: models)
      -s, --server-package=           the package to save the server specific code (default: restapi)
      -c, --client-package=           the package to save the client specific code (default: client)
      -t, --target=                   the base directory for generating the files (default: ./)
      -T, --template-dir=             alternative template override directory
      -C, --config-file=              configuration file to use for overriding template options
      -r, --copyright-file=           copyright file used to add copyright header
          --existing-models=          use pre-generated models e.g. github.com/foobar/model
          --additional-initialism=    consecutive capitals that should be considered intialisms
      -n, --name=                     the model to generate
          --skip-struct               when present will not generate the model struct
          --dump-data                 when present dumps the json for the template generator instead of generating files
          --skip-validation           skips validation of spec prior to generation
```

