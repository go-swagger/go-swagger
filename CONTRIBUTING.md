# Contributing to go-swagger

Thank you for your interest in contributing to `go-swagger`!
This document explains how to set up a development environment, propose changes,
and submit a pull request that is easy for maintainers to review.

`go-swagger` is a community-driven project. By participating you agree to follow
our [Code of Conduct](CODE_OF_CONDUCT.md).

---

## Table of contents

1. [Project layout](#project-layout)
2. [Getting started](#getting-started)
3. [Development workflow](#development-workflow)
4. [Coding conventions](#coding-conventions)
5. [Testing](#testing)
6. [Commit & pull request guidelines](#commit--pull-request-guidelines)
7. [Issue triage](#issue-triage)
8. [Release process](#release-process)
9. [Getting help](#getting-help)

---

## Project layout

```
.
├── cmd/                # CLI entry points (swagger, swagger generate, ...)
├── generator/          # Code generation engine (server, client, model, ...)
├── codescan/           # Annotation scanner used by the generator
├── docs/               # User-facing documentation that ships to goswagger.io
├── hack/               # Internal codegen test programs and tooling
├── fixtures/           # Codegen test data (specs, expected output, ...)
├── .github/workflows/  # CI: tests, builds, releases, security scans
├── CODE_OF_CONDUCT.md
├── SECURITY.md
└── README.md
```

The `generator` package is the heart of the project: most behavioural changes
either live there or in the `cmd/swagger/commands` tree that drives it.

## Getting started

### Prerequisites

* **Go 1.25 or newer** — see `go.mod` for the exact minimum version.
  The repository uses Go toolchain directives, so any 1.25+ toolchain will
  fetch the right version automatically.
* **git** (recent enough to handle shallow clones and partial checkouts).
* **make** (optional but recommended — see the targets in `hack/`).
* A POSIX shell on Linux or macOS. Windows users typically work inside WSL.

### Fork & clone

```bash
# 1. Fork the repo on GitHub, then:
git clone https://github.com/<your-username>/go-swagger.git
cd go-swagger

# 2. Add the upstream remote so you can sync with master:
git remote add upstream https://github.com/go-swagger/go-swagger.git
git fetch upstream
```

### Build

```bash
# Build the CLI to ./bin/swagger
make swagger
# or, without make:
go build -o bin/swagger ./cmd/swagger
```

You should now be able to run `./bin/swagger version` and see a version string.

## Development workflow

1. **Sync with upstream** before starting work:
   ```bash
   git checkout master
   git pull upstream master
   git push origin master
   ```
2. **Create a topic branch** off `master`:
   ```bash
   git checkout -b type/short-description
   ```
   Use one of the conventional prefixes: `feat/`, `fix/`, `docs/`, `refactor/`,
   `test/`, `chore/`, `ci/`, `build/`, `perf/`.
3. **Make focused commits**. Each commit should compile and, where reasonable,
   leave the tree in a green test state. Avoid mixing refactors with behaviour
   changes.
4. **Run the test suite locally** before pushing (see [Testing](#testing)).
5. **Push your branch** and open a pull request against `go-swagger:master`.

## Coding conventions

* **Formatting**: `gofmt -s -w .` (or `go fmt ./...`). Imports are grouped with
  `goimports`.
* **Linting**: the project uses `golangci-lint` (see `.golangci.yml`). Run it
  locally with:
  ```bash
  golangci-lint run ./...
  ```
* **Static checks**: keep `go vet ./...` clean.
* **Naming**: follow effective Go. Exported identifiers carry a doc comment
  that starts with the identifier name.
* **Errors**: wrap with `fmt.Errorf("...: %w", err)` and provide enough context
  for end users to act on.
* **Public API**: changes to `generator/` or the CLI command surface are
  considered API changes and require a clear changelog entry (see below).
* **Backwards compatibility**: avoid breaking exported Go APIs and CLI flags
  unless the change is coordinated in an issue first.

## Testing

`go-swagger` relies heavily on golden-file testing against the fixtures in
`fixtures/`. Please make sure new behaviour is covered.

```bash
# Fast unit tests
go test ./...

# Tests with the race detector (recommended for generator changes)
go test -race ./generator/...

# Regenerate codegen fixtures after intentional output changes
go generate ./...
# or
make codegen
```

The CI pipeline (see `.github/workflows/test.yaml`) runs a matrix of Go
versions, the race detector, `golangci-lint`, CodeQL, and a codegen harness.
Reproduce a CI failure locally by running the same `go test` invocation with
the relevant `-run` filter.

If you add a new fixture, document it briefly in `fixtures/README.md` (if it
exists) or in the PR description so reviewers can navigate it.

## Commit & pull request guidelines

### Commit messages

We use [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<optional scope>): <short summary>

<body explaining motivation and approach>

<footer with references, e.g. Closes #1234 or BREAKING CHANGE: ...>
```

Allowed types: `feat`, `fix`, `docs`, `refactor`, `test`, `chore`, `ci`,
`build`, `perf`, `revert`.

### Pull request checklist

Before requesting review, make sure:

- [ ] The branch is up to date with `go-swagger:master`.
- [ ] `gofmt`, `go vet`, and `golangci-lint` are clean.
- [ ] `go test ./...` passes locally; for generator changes, `go test -race`
      also passes.
- [ ] New behaviour is covered by tests (and fixtures, where applicable).
- [ ] Public API/CLI changes are documented in the PR description and have a
      changelog entry (see below).
- [ ] Commits are logically grouped; consider `git rebase -i upstream/master`
      before pushing.
- [ ] The PR description links the issue it addresses (`Closes #NNNN`) and
      explains *why*, not just *what*.

A typical PR description looks like:

```markdown
## Summary
- One bullet describing the change.

## Motivation
- Why is this needed? Link the issue. What alternatives were considered?

## Test Plan
- [ ] `go test ./...`
- [ ] `go test -race ./generator/...`
- [ ] Manual run of `./bin/swagger generate spec -w fixtures/...`
```

### Review expectations

* Expect at least one maintainer review before merge.
* Address review comments with follow-up commits (do not force-push during
  review unless asked — it makes review harder).
* CI must be green before merge. Maintainers may push minor fixes to your
  branch; please grant write access to the fork if you want to opt out.

## Issue triage

* **Bug reports**: include the `swagger version` output, the input spec (or a
  minimal reproducer), the expected vs. actual behaviour, and any error logs.
  Use the "Bug report" issue template.
* **Feature requests**: explain the use case, not just the proposed solution.
  Large generator features usually warrant a design note first.
* **Questions**: prefer the [Discord channel][discord] or the project
  discussions; only file an issue if you are reporting a defect or proposing a
  tracked change.

If you want to help triage, look for issues labelled
[`good first issue`][good-first-issue] or [`help wanted`][help-wanted].

## Release process

Releases are cut by maintainers from `master` and follow [semver](https://semver.org/).
The release pipeline (`.github/workflows/release.yaml`) signs tags and
artifacts. You usually do **not** need to do anything special to get a change
shipped — once merged, it lands in the next tagged release.

## Getting help

* Documentation: <https://goswagger.io>
* Discord: see the badge in [README.md](README.md)
* Slack: see the badge in [README.md](README.md) (open until 2026-06-30)
* Security issues: see [SECURITY.md](SECURITY.md) — **do not** file
  vulnerabilities as public issues.

[discord]: https://github.com/go-swagger/go-swagger#readme
[good-first-issue]: https://github.com/go-swagger/go-swagger/issues?q=is%3Aissue+is%3Aopen+label%3A%22good+first+issue%22
[help-wanted]: https://github.com/go-swagger/go-swagger/issues?q=is%3Aissue+is%3Aopen+label%3A%22help+wanted%22

---

Thanks again for contributing — `go-swagger` is better because of people like
you. ❤️