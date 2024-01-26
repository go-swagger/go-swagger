---
title: Continuous Integration
date: 2023-01-01T01:01:01-08:00
draft: true
weight: 60
description: Contributing guide - CI setup
---
# Continous integration


## `go-swagger/go-swagger`

Enabled CI engines and bots and apps:
- GitHub Actions (Linux, MacOS, Windows)
- Codecov
- DCO (enforce signed-off commits)
- WIP (blocks PRs with title WIP/do not merge, etc...)

### Build
The CI pipeline builds the `swagger` binary and runs a few spec validation commands: this is a smoke test.

### Test
Codecov results are not blocking.

We run unit tests on the two most recent go versions of the 3 platforms above.
Tests run with race detection: `go test -race`.

Integration tests run a real swagger CLI command to generate servers, clients and models.
The CI pipeline uses this test tool to iterate over the swagger spec fixtures: `hack/codegen_nonreg_test.go`.

Integration tests are divided in 2 groups:
* "canary" specs: a bunch of rather larger real life specs (e.g. kubernetes, docker, quay.io...)
* fixtures: many trickier specs intended to exercise the code generation

The go test program `codegen_nonreg_test.go` runs on CI with various generation options.
You may alos run it manually on your local environment to explore more generation options (expand spec, flatten, etc...).

### Releases
Releases are cut with a separate workflow to build artifacts and bake docker images.

### Documentation update
A github action takes care of generating this web site. It builds a github pages artifact that is deployed.

## `go-openapi/...`

Enabled CI engines and apps:
- GitHub Actions (Linux, MacOS, Windows)
- Codecov
- DCO

{{< hint "warning" >}}
We are trying to keep the CI configuration in all go-openapi repos aligned.
{{< /hint >}}
