# CLAUDE.md — go-swagger

## Project Overview

Go implementation of [OpenAPI 2.0](https://swagger.io/specification/v2/) (Swagger).
Generates server, client, CLI, and model code from specs — or specs from annotated Go source.

**Entry point:** `cmd/swagger/swagger.go` — CLI built on `jessevdk/go-flags`.

### Key packages

| Package | Purpose |
|---------|---------|
| `cmd/swagger/commands/` | CLI commands: validate, generate (server/client/model/spec/cli/markdown), serve, expand, flatten, mixin, diff |
| `codescan/` | Scans Go source for swagger annotations → builds spec |
| `generator/` | Template-based code generation engine (server, client, model, CLI, markdown) |
| `generator/templates/` | Go templates organized by target (client/, server/, cli/, markdown/, contrib/) |
| `hack/` | Build scripts, fixture configs, helper tools |
| `docs/` | Hugo documentation site (published to goswagger.io via GitHub Pages) |
| `examples/` | ~20 example projects with specs and generated code |
| `fixtures/` | Test specs and expected outputs for regression testing |

### Key dependencies

| Dependency | Role |
|------------|------|
| `go-openapi/spec` | OpenAPI 2.0 document model |
| `go-openapi/analysis` | Spec analysis, flattening, mixin |
| `go-openapi/loads` | Spec loading and unmarshaling |
| `go-openapi/validate` | Spec and data validation |
| `go-openapi/runtime` | Runtime support for generated code |
| `go-openapi/strfmt` | String format types |
| `go-openapi/swag` | JSON/YAML utilities, name mangling |
| `jessevdk/go-flags` | CLI flag parsing |
| `Masterminds/sprig/v3` | Template functions |

## Conventions

Coding conventions are found beneath `.claude/rules/` (also symlinked as `.github/copilot/`).

### Summary

- All `.go` files must have SPDX license headers (Apache-2.0).
- Commits require DCO sign-off (`git commit -s`).
- Go version policy: support the 2 latest stable Go minor versions.
- Linting: `golangci-lint run` — config in `.golangci.yml` (posture: `default: all` with explicit disables).
- Every `//nolint` directive **must** have an inline comment explaining why.
- Tests: `go test ./...` (single module). CI runs on `{ubuntu, macos, windows} x {stable, oldstable}` with `-race`.
- Test framework: `github.com/go-openapi/testify/v2` (not `stretchr/testify`; `testifylint` does not work).

See `.claude/rules/` for detailed rules on Go conventions, linting, testing, GitHub Actions, and contributions.
See `.github/STYLE.md` for the linting posture rationale.

## CI / Release pipeline

| Workflow | Trigger | Purpose |
|----------|---------|---------|
| `test.yaml` | PR, push to master | Full test matrix, linting, coverage |
| `release.yaml` | Tag `v*` | goreleaser + git-cliff + GPG signing + Discord announce |
| `build-docker.yaml` | Called by release/master | Multi-arch Docker images → ghcr.io, quay.io |
| `update-doc.yaml` | Called by release | Publish docs to GitHub Pages |
| `auto-merge.yaml` | Dependabot PRs | Auto-merge dependency updates |
| `codeql.yaml` | Schedule, PR | Security analysis |
| `scorecard.yaml` | Schedule, push | OpenSSF scorecard |

### Release artifacts

A tagged release produces:
- Signed binaries (linux/darwin/windows × amd64/arm64/arm)
- Signed deb/rpm packages
- Signed source archive + checksums
- Docker multi-arch images (ghcr.io, quay.io)
- Release notes (git-cliff with tag message)
- Discord announcement
