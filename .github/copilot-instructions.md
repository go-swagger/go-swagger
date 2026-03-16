# Copilot Instructions

## Project Overview

Go implementation of [OpenAPI 2.0](https://swagger.io/specification/v2/) (Swagger).
Generates server, client, CLI, and model code from specs — or specs from annotated Go source.
Entry point: `cmd/swagger/swagger.go` (CLI built on `jessevdk/go-flags`).

### Key packages

| Package | Purpose |
|---------|---------|
| `cmd/swagger/commands/` | CLI commands: validate, generate, serve, expand, flatten, mixin, diff |
| `codescan/` | Scans Go source for swagger annotations, builds spec |
| `generator/` | Template-based code generation engine |
| `generator/templates/` | Go templates by target (client, server, cli, markdown, contrib) |
| `docs/` | Hugo documentation site (goswagger.io) |

### Key dependencies

- `go-openapi/spec` — OpenAPI 2.0 document model
- `go-openapi/analysis` — spec analysis, flattening, mixin
- `go-openapi/loads` — spec loading and unmarshaling
- `go-openapi/validate` — spec and data validation
- `go-openapi/runtime` — runtime support for generated code
- `go-openapi/strfmt` — string format types
- `go-openapi/swag` — JSON/YAML utilities, name mangling

## Conventions

Coding conventions are found beneath `.github/copilot/`

### Summary

- All `.go` files must have SPDX license headers (Apache-2.0).
- Commits require DCO sign-off (`git commit -s`).
- Linting: `golangci-lint run` — config in `.golangci.yml` (posture: `default: all` with explicit disables).
- Every `//nolint` directive **must** have an inline comment explaining why.
- Tests: `go test ./...` (single module). CI runs on `{ubuntu, macos, windows} x {stable, oldstable}` with `-race`.
- Test framework: `github.com/go-openapi/testify/v2` (not `stretchr/testify`; `testifylint` does not work).

See `.github/copilot/` (symlinked to `.claude/rules/`) for detailed rules on Go conventions, linting, testing, and contributions.
