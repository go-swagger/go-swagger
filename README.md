# Swagger 2.0

Contains an implementation of Swagger 2.0.
It knows how to serialize and deserialize swagger specifications.

At present it's got the entire object model defined, and I'm writing tests to make that work completely.

Planned:
* Generate validations based on the swagger spec
* Later it will also know how to generate those specifications from your source code.
* Generate a stub api based on a swagger spec
* Generate a client from a swagger spec
* Build a full swagger spec by inspecting your source code and embedding it in a go file.
