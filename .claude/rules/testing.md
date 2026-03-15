---
paths:
  - "**/*_test.go"
---

# Testing conventions (go-openapi)

## Running tests

**Single module repos:**

```sh
go test ./...
```

**Mono-repos (with `go.work`):**

```sh
# All modules
go test work ./...

# Single module
go test ./conv/...
```

Note: in mono-repos, plain `go test ./...` only tests the root module.
The `work` pattern expands to all modules listed in `go.work`.

CI runs tests on `{ubuntu, macos, windows} x {stable, oldstable}` with `-race` via `gotestsum`.

## Fuzz tests

```sh
# List all fuzz targets
go test -list Fuzz ./...

# Run a specific target (go test -fuzz cannot span multiple packages)
go test -fuzz=Fuzz -run='FuzzTargetName$' -fuzztime=1m30s ./package
```

Fuzz corpus lives in `testdata/fuzz/` within each package. CI runs each fuzz target for 1m30s
with a 5m minimize timeout.

## Test framework

`github.com/go-openapi/testify/v2` — a zero-dep fork of `stretchr/testify`.
Because it's a fork, `testifylint` does not work.
