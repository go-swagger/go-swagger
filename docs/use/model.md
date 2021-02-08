# Generate a data model from swagger spec

The toolkit allows for generating go native structures to serialize and validate a swagger compliant model definition.

The generated objects follow the specified validation rules, including extended format directives for strings and numbers.

Generated models support most Swagger 2.0 features, including polymorphism.

Generated models support most JSON-schema draft4 features, including `AllOf`, `AdditionalProperties` and `AdditionalItems`.

### Usage

`generate model -f {spec}`

See the full list of available options [here](../generate/model.md).

### Model building rules

* [Schema generation rules](./models/schemas.md#schema-generation-rules)
  * [About schemas](./models/schemas.md#about-schema)
  * [Interfaces](./models/schemas.md#interfaces)
  * [Mapping patterns](./models/schemas.md#mapping-patterns)
    * [Minimal use of go's reflection](./models/schemas.md#minimal-use-of-go-s-reflection)
    * [Doc strings](./models/schemas.md#doc-strings)
    * [Reusability](./models/schemas.md#reusability)
  * [Swagger vs JSONSchema](./models/schemas.md#swagger-vs-jsonschema)
  * [Go-swagger vs Swagger](./models/schemas.md#go-swagger-vs-swagger)
  * [Known limitations with go-swagger models](./models/schemas.md#known-limitations-with-go-swagger-models)
  * [Custom extensions](./models/schemas.md#custom-extensions)
  * [Primitive types](./models/schemas.md#primitive-types)
  * [Formatted types](./models/schemas.md#formatted-types)
  * [Nullability](./models/schemas.md#nullability)
  * [Validation](./models/schemas.md#validation)
  * [Type aliasing](./models/schemas.md#type-aliasing)
  * [Extensible types](./models/schemas.md#extensible-types)
    * [Additional properties](./models/schemas.md#additional-properties)
    * [Tuples and additional items](./models/schemas.md#tuples-and-additional-items)
  * [Polymorphic types](./models/schemas.md#polymorphic-types)
  * [Serialization interfaces](./models/schemas.md#serialization-interfaces)
  * [External types](./models/schemas.md#external-types)
  * [Customizing struct tags](./models/schemas.md#customizing-struct-tags)
