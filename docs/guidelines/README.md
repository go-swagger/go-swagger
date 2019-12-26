# Guidelines to maintainers

A quick guide on how to contribute to `go-swagger` and the other `go-openapi` repos.

### Getting started

Repos follow standard go building and testing rules.

Cloning `go-swagger`:
```bash
mkdir -p $GOPATH/src/github.com/go-swagger
cd $GOPATH/src/github.com/go-swagger
git clone https://github.com/go-swagger/go-swagger
```

All dependencies are available in the checked out `vendor` directory.

Building and installing go-swagger from source on your system:
```
go install github.com/go-swagger/go-swagger/cmd/swagger
```

Running standard unit tests:
```bash
go test ./...
```

More advanced tests are run by CI. See [below](#continuous-integration).

### Generally accepted rules for pull requests

All PR's are welcome and generally accepted (so far, 95% were accepted...).
There are just a few common sense rules to be followed.

1. PRs which are not ready to merge should be prefixed with `WIP:`
2. Generally, contributors should squash their commits (use `git rebase -i master`)
3. Provide sufficient test coverage with changes
4. Do not bring in uncontrolled dependencies, including from fixtures or examples
5. Use the `fixes #xxx` github feature in PR to automate issue closing
6. Sign-off commits with `git commit -s`. PGP-signed commits with verified signatures are not mandatory (but much appreciated)
7. Use the "draft PR" github feature to draft some work and exercise it against our CI before review

### Go environment

We want to always support the **two most recent go versions**.

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

All repos should remain go-gettable (i.e. available with the `go get ./...` command)
and testable with `go test ./...`

### Linting

Check your work with golangci linter:
`go get -u github.com/golangci/golangci-lint/cmd/golangci-lint`
`golangci-lint run --new-from-rev HEAD`

### Continuous integration

All PR's require a review by a team member, whatever the CI engines tell.

##### go-swagger/go-swagger

Enabled CI engines and bots:
- CircleCI (linux)
- Appveyor (windows)
- GolangCI
- Codecov
- DCO (enforce signed-off commits)
- WIP (blocks PRs with title WIP/do not merge, etc...)

Codecov results are not blocking.

CI runs description/configuration:

| CI engine | Test type     | Configuration             | Comment |
|---        |---            |---                        |---      |
| CircleCI  | unit test     | .circleci/config.yml      |         |
|           | build test (1)| ./hack/codegen_nonreg_test.go  | Codegen and build test on many (~ 80) specs in `fixtures/codegen` and `fixtures/bugs``|
|           | build test (2)| ./hack/codegen_nonreg_test.go  | Codegen and build test on (large) specs in fixtures/canary`|
| Appveyor  | unit test     | appveyor.yml              | `go test -v ./...` |
| GolangCI  | linting       | .golangci.yml             | equ. `golangci-lint run` |
| Codecov   | test coverage | -                         | project test coverage and PR diff coverage|
| DCO       | commit signed | -                         | https://probot.github.io/apps/dco|
| WIP       | PR title      | -                         | https://github.com/apps/wip|

Deprecated engines:
- hound (`.hound.yml`): previous linting checker before we moved to golangCI

> **NOTE on Appveyor**:
> Appveyor runs our UT on Windows. This makes sure everything works fine
> on this platform as well.
> The peculiarity with this CI is that it does not tolerate output to `stderr`
> from test programs (this is actually a Powershell limitation).
> Therefore, please make sure your UT remain mute on stderr or capture the
> output if you need to assert something from the output.

The go test program `./hack/codegen_nonreg_test.go` runs on CI with various generation options.
You may run it manually to explore more generation options (expand spec, flatten, etc...).

##### Releases

Released are cut from CircleCI, with a separate workflow to build artifacts, and bake docker images.

##### go-openapi repos

Enabled CI engines:
- Travis (linux)
- GolangCI
- Codecov
- experimental: Appveyor

| CI engine | Test type     | Configuration             | Comment |
|---        |---            |---                        |---      |
| Travis    | unit test     | .travis.yml               | `go test -v ./...` |
| GolangCI  | linting       | .golangci.yml             | equ. `golangci-lint run` |

### Vendoring

The whole `go-openapi` and `go-swagger` has adopted go modules and no more use vendoring.

### Update templates

`go-swagger` is built with an in-memory image of templates.

Binary encoded assets are auto-generated from the `generator/templates` directory using `bindata`
(the result is `generator/bindata.go`).

While developing, you may work with dynamically updated templates (i.e. no need to rebuild)
using `bindata.go` generated in debug mode (use script: `./generator/gen-debug.sh`).

There is a `.githook` script configured as pre-commit: if you configure githooks locally for this repo,
every time you commit,`generator/bindata.go`
will be regenerated and added to the current commit (without debug mode).

For `bindata` please use the fork found at: `github.com/kevinburke/go-bindata`.

> **NOTE**: we are carrying out unit tests on codegen mostly by asserting lines in generated code.
> There is a bunch of test utility functions for this. See `generator/*_test.go`.
>
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
