- Phase 1 (Core MVP)
  - Add --oas 3.1 flag and Options.TargetOAS.
  - Introduce SpecWriter abstraction (V2Writer existing, new OAS31Writer).
  - Map swagger:route/operation to OAS 3.1 operations and paths.
  - Implement requestBody/content and single-media-type responses.
  - Schema generator (2020-12 core): objects, arrays, primitives, required, null-unions.
  - Components.schemas registry and $ref wiring.
  - Golden tests + 3.1 validation in CI.

- Phase 2 (I/O richness and security)
  - Multiple media types (request/response), headers, examples.
  - Parameter style/explode, cookie params.
  - components.securitySchemes and global/operation security.

- Phase 3 (Schema depth and graph features)
  - oneOf/anyOf/allOf, discriminator, const.
  - Links and Callbacks.

- Phase 4 (3.1-only features + polish)
  - webhooks, components.pathItems.
  - Strict mode (--strict-3.1), helpful diagnostics and migration warnings.
  - Expand docs and finalize annotation reference.

- Continuous
  - Back-compat guard: assert v2 output golden files unchanged.
  - Interop tests against popular viewers/generators.
