---
title: Spec transformations
date: 2023-01-01T01:01:01-08:00
draft: true
---
# Transforming specifications

The `swagger` toolkit allows some transformation to be carried out with a specification.

Currently it is possible to apply the following transforms:
- expansion: this expands all `$ref`'s in a spec
- minimal flattening: carries on minimal transformation to a spec to be workable for the swagger codegen
- full flattening: performs minimal flattening and in addition, replaces all complex constructs in schemas by named definitions
- mixin: merges one or more specifications into a primary spec

In addition, it is possible to compare specs (diff) to inspect breaking changes in the API.

### Expansion

Expanding a spec may prove useful to validate a schema, produce documentation or test cases. The primary intent is not code generation.

> **NOTE**: Circular `$ref` are detected and remain as local `$ref`. Remote circular `$ref` remain remote 
> (favor `flatten` to regroup all remote `$ref` under one single root document).

Usage:

`swagger expand {spec}`
or

`swagger flatten --with-expand` {spec}`

Full list of available options [for expand](../../usage/expand.md) and [for flatten](../../usage/flatten.md).

or with codegen commands:

`swagger generate [model|server|client|operation|...] --spec={spec} --with-expand`

> **NOTE**: codegen may fail in some situations with spec expansion:
> - polymorphism: the original intent of the `$ref` pointing to a discriminated type (i.e. _extends type..._) is lost with expansion
> - duplicate names: expansion of `$ref` intended for type reuse may produce duplicate type identifiers
> - too many generated files: expansion of large specs may lead to many independent anonymous structures

### Minimal flattening

This transform makes complex JSON `$ref` amenable to analysis and code generation.

Minimal flattening attempts to minimally distort the original specification intent, while ensuring workable codegen.

Minimal flattening does:

- bundle all external `$ref`s into the local document
- resolve JSON pointers to anonymous places in the document (e.g. `"$ref": "#/definitions/thisModel/properties/codeName"`)
- resolve `$ref`s in swagger-specific sections: `parameters` and `responses`, so the only remaining `$ref`s are located in schemas

All `$ref`s in the resulting spec are thus only found in schema, with the _canonical_ form: `"$ref": "#/definitions/modelName"`.

Usage:

`swagger flatten {spec}`

or more explicitly:

`swagger flatten --with-flatten=minimal {spec}`

This is the default option for codegen commands:

`swagger generate [model|server|client|operation|...] --spec={spec}`

> **NOTE**: `$ref` bundling / pointer resolving may in some cases produce duplicate names.
> The flattener tries to resolve duplicates whenever possible (such as when a child identifier duplicates its parent).
> When such resolution is not possible, an "OAIGen" suffix is added to the definition name (a warning is issued).
> You may use the `x-go-name` extension here to set a more suitable name for the generated type.

### Full flattening

Full flattening is useful to factorize data model objects into simpler structures.

> Complex structures (i.e. objects with properties or schemas with an `allOf` composition) are moved to standalone definitions.
> Arrays and map constructs (e.g. AdditionalProperties) are not considered complex.

Usage:

`swagger flatten --with-flatten=full {spec}`

Or with codegen commands:

`swagger generate [model|server|client|operation|...] --spec={spec} --with-flatten=full`

> **NOTE**: this used to be the default for codegen commands with releases 0.13 and 0.14. 
> This behavior has been reverted with release 0.15.

You may not like the automatic names chosen for the new structures (e.g. `myModelAllOf3`, `myModelAdditionalPropertiesAdditionalProperties` ...).
Again, the `x-go-name` extension is the way to generate custom names that are easier to use in your API code.

Ex:
```yaml
definitions:
  complexArray:
    schema:
      type: array
      items:
        type: object        # <- inline schema for items
        properties:
          prop1:
            type: integer
```

Is factorized as:

```yaml
definitions:
  complexArray:
    schema:
      type: array
      items:
        $ref: '#/definitions/complexArrayItems'
  complexArrayItems:
    type: object
    properties:
      prop1:
        type: integer
```

### Other flattening options

The `--with-flatten` option supports the following additional arguments:

- `verbose`, `noverbose`: allows/mute warnings about the transformation
- `remove-unused`: removes unused definitions after expansion or flattening

> **NOTE**: you may specify multiple options like this:
>
> `swagger flatten {spec} --with-flatten=remove-unused --with-flatten=full`

### Mixin

Usage:

`swagger mixin {primary spec} [{spec to merge}...]`

Full list of available options [here](../../usage/mixin.md).

### Roadmap

This set of features is essentially provided by the `github.com/go-openapi/analysis` package.
Feel free to contribute new features to this repo.

Currently, here is a todo list of improvements planned to the spec preprocessing feature:

- in full flatten mode, more options to control how allOf are handled
- in full flatten mode propagating `x-go-name` (as prefix) to newly created definitions. 
As creating new definitions introduces some naming conventions.
If a x-go-name has been set in the original structure, new names should remain consistent with their container...
- struct name analysis and better prediction/resolution of duplicate identifiers
