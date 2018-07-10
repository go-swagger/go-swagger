# Guidelines to maintainers

A quick guide on how to contribute to `go-swagger` and the other `go-openapi` repos.

### Getting started

Repos follow standard go building and testing rules.

Building and installing go-swagger from source on your system:
```
go install github.com/go-swagger/go-swagger/cmd/swagger
```

Running standard unit tests:
```
 go test ./...
``` 

More advanced tests are run by CI. See [below](#continuous-integration).

### Generally accepted rules for pull requests

All PR's are welcome and generally accepted (so far, 95% were accepted...).
There are just a few common sense rules to be followed.

1. PRs which are not ready to merge should be prefixed with "WIP:"
2. Generally, contributors should squash their commits (use `git rebase -i ...`)
3. Provide sufficient test coverage with changes
4. Do not bring in uncontrolled dependencies, including from fixtures or examples.
Adding dependencies is possible with a vendor update (`dep ensure -update`).
5. Use the "fixes #xxx" feature in PR to automate issue closing

### Go environment

We want to always support the two most recent go versions.

However, we try to avoid introducing breaking changes, especially on the more
stable `go-openapi` repos. We manage this with build tags. Notice the very
important blank line after your build tag comment line.

Example (from `go-openapi/swag`):
```go 
// +build !go1.8

package swag

import "net/url"

func pathUnescape(path string) (string, error) {
	return url.QueryUnescape(path)
}
```

All repos should remain go-gettable and testable with `go test ./...`

### Continuous integration

All PR's require a review by a team member, whatever the CI engines tell.

##### go-swagger/go-swagger

Enabled CI engines:
- CircleCI (linux)
- Appveyor (windows)
- GolangCI
- Codecov

Codecov results are not blocking.

CI runs description/configuration:

| CI engine | Test type     | Configuration             | Comment |
|---        |---            |---                        |---      |
| CircleCI  | unit test     | .circleci/config.yml      |         |
|           | build test (1)| ./hack/codegen-nonreg.sh  | Codegen and build test on many (~ 80) specs in `fixtures/codegen` and `fixtures/bugs``|
|           | build test (2)| ./hack/run-canary.sh      | Codegen and build test on (large) specs in fixtures/canary`|
| Appveyor  | unit test     | appveyor.yml              | `go test -v ./...` |
| GolangCI  | linting       | .golangci.yml             | equ. `golangci-lint run` |
| Codecov   | test coverage | -                         | project test coverage and PR diff coverage|

Deprecated engines:
- hound (`.hound.yml`): previous linting checker before we moved to golangCI

> **NOTE on Appveyor**:
> Appveyor runs our UT on Windows. This makes sure everything works fine
> on this platform as well.
> The peculiarity with this CI is that it does not tolerate output to `stderr`
> from test programs (this is actually a Powershell limitation).
> Therefore, please make sure your UT remain mute on stderr or capture the
> output if you need to assert something from the output.

The `./hackcodegen-nonreg.sh` runs on CI with a single generation option.
You may run it manually to explore more generation options (expand spec, flatten, etc...).

CircleCI has a separate CI workflow to build releases, baking and checking newly released docker 
images.

##### go-openapi repos

Enabled CI engines:
- Travis (linux)
- GolangCI
- Codecov

> **NOTE**: setting up Appveyor on go-openapi/spec and validate is on the todo list.

| CI engine | Test type     | Configuration             | Comment |
|---        |---            |---                        |---      |
| Travis    | unit test     | .travis.yml               | `go test -v ./...` |
| GolangCI  | linting       | .golangci.yml             | equ. `golangci-lint run` |

### Vendoring

`go-swagger/go-swagger` repo comes with a vendor directory. This is because
we release binary distributions (docker, debian...).

The `go-openapi` packages are **not** vendored.

Vendoring is managed using the current _official_ `dep` tool.
Configuration is in `Gopkg.toml`.

Run `dep ensure -update` to update dependencies. Please do not cherry-pick
updates manually.

> **NOTE**: since there are many dependencies, running an update may update
> many things.
> We prefer to get vendor updates as separate commits with changes to vendor only.

##### Testing PRs requiring integration of a dependency (e.g. another pending PR on `go-openapi`)

This happens for instance whenever you want to test the full integration in `go-swagger` of an unmerged PR
in one of the `go-openapi` repos.

With `go-swagger` (vendored):
- prepare a "WIP" PR with a temporary vendor update commit
- this vendor update temporarily alters the `branch` and `source` in `Gopkg.toml`
to get the proper version of the required dependency from the unmerged branch (e.g. from your fork)

With `go-openapi` (non vendored):
- for your "WIP" PR, temporarily alter the CI config script (e.g. `.travis.yml`) and 
replace the `go get` requirements to build your CI with the adequate `git clone`
pointing to the required branches

### Update templates

`go-swagger` is built with an in-memory image of templates.
Binary encoded assets are auto-generated from the `generator/templates` directory using `bindata`
(the result is `generator/bindata.go`).

While developing, you may work with dynamically updated templates (i.e. no need to rebuild)
using `bindata.go` generated in debug mode (use script: `./generator/gen-debug.sh`).

There is a `.githook` script configured as pre-commit: every time you commit to the repo, `generator/bindata.go`
is regenerated and added to the current commit (without debug mode).

> **NOTE**: we are carrying out unit tests on codegen mostly by asserting lines in generated code.
> There is a bunch of test utility functions for this. See `generator/*_test.go`.
> If you want to bring in more advanced testing go programs with your fixtures, please tag 
> those so they don't affect the `go ./...` command (e.g. with `// +build +integration`).

### Updating examples

Whenever code generation rules change, we feel it is important to maintain
consistency with the generated code provided as examples.  

The script `./hack/regen-samples.sh` does just that.

Do not forget to update this script when you add a new example.

### Writing documentation

##### go-swagger/go-swagger

The `go-swagger` documentation site (`goswagger.io`) is built with GitBooks.
Configuration is in `book.json`. The documents root is in `./docs`

We systematically copy the repository main `README.md` to `docs/README.md`.
Please make sure links work both from github and gitbook.

There is also a minimal godoc for goswagger.

Please make sure new CLI options remain well documented in `./docs/usage` and `./docs/generate`.

##### go-openapi repos

Documentation is limited to the repo's README.md and godoc, published on godoc.org.
