---
paths:
  - "**/*.go"
---

# Linting conventions (go-openapi)

```sh
golangci-lint run
```

Config: `.golangci.yml` — posture is `default: all` with explicit disables.
See `docs/STYLE.md` for the rationale behind each disabled linter.

Key rules:
- Every `//nolint` directive **must** have an inline comment explaining why.
- Prefer disabling a linter over scattering `//nolint` across the codebase.
