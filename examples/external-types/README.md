# External types examples

The sample specification provided here demonstrates how to bind your generated API
with some externally defined types for schemas.

This [swagger specification](./example-external-types.yaml) illustrates the following use cases:

1. refer to an external type imported from some external package, as its own definition
2. refer to an external type imported from some external package, as part of another type (object, slice, map or tuple)
3. refer to an external type imported from the default location of models (i.e. besides generated models)
4. embed external type to add the Validatable interface.
5. use hints in annotations to solve nullable/struct vs interface issues 
6. use hints in annotations to skip the validation of external types

> NOTE: due to the addition of the "additionalItems" clause to illustrate tuples, the spec is not formally
> valid against the swagger 2.0 schema.

[Reference documentation](../../docs/use/models/schemas.md#external-types)
