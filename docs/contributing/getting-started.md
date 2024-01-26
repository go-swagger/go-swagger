---
title: Getting started
date: 2023-01-01T01:01:01-08:00
draft: true
description: Contributing guide - Getting started
weight: 10
---

# Getting started

## Your development environment

You only need a go compiler >= {{< param goswagger.goVersion >}}.
You may develop on Linux, MacOS or Windows.

## Cloning `go-swagger`

```sh
mkdir -p $GOPATH/src/github.com/go-swagger
cd $GOPATH/src/github.com/go-swagger
git clone https://github.com/go-swagger/go-swagger
```

Building and installing go-swagger from source on your system:
```sh
go install github.com/go-swagger/go-swagger/cmd/swagger@latest

swagger version
dev
```

Building and installing go-swagger from your local clone:
```sh
cd $GOPATH/src/github.com/go-swagger
go install ./cmd/swagger

swagger version
dev
```

## Sending us a Pull Request

All PR's are welcome and generally accepted (so far, 95% were accepted...).

All PR's require a review by a team member.

These are just a few common sense rules to be followed.

### Before you push

1. Please open an issue on github to describe your proposal, feature or bug fix
   We encourage PR authors to engage other maintainers in issues.
2. Please refer to the issue as `* fixes #xxx` in your commit body to automate issue closing

### Your PR

1. PRs that are not ready to merge should be prefixed with `WIP:`
2. Use the "draft PR" github feature to signal that your work is in progress and exercise it against our CI before review.
3. Squash your commits (use `git rebase -i master`). Please make sure the resulting commit remains readable and
   meaningful.
4. Provide sufficient test coverage with changes
5. Do not bring in uncontrolled dependencies, including from fixtures or examples. If your idea really requires it, engage a conversation with
   the other maintainers to collect feedback.
6. Sign-off commits with `git commit -s`. PGP-signed commits with verified signatures are not mandatory (but much appreciated)

## Working with our repositories

`go-swagger` exposes a command line interface (CLI) to functionality largely built on top of
the `go-openapi` packages.

All these repos are go-gettable (i.e. available with the `go get ./...` command) and follow the standard go building and testing procedures.

Running standard unit tests:
```sh
go test ./...
```

Specifically for `go-swagger`, we run additional integration tests in CI. See [Continuous Integration](ci.md).

---

The diagram below displays a family picture of the `go-openapi` eco-system.

{{< mermaid class="optional" >}}
---
title: "Direct package dependencies"
---
flowchart TD
    A((go-swagger))

    B[go-openapi/runtime]
    C[go-openapi/loads]
    D[go-openapi/analysis]
    E[go-openapi/validate]
    F[go-openapi/spec]
    G[go-openapi/jsonreference]
    H[go-openapi/jsonpointer]
    I[go-openapi/strfmt]
    J[go-openapi/errors]
    K[go-openapi/swag]
    L[go-openapi/inflect]

    A--> D
    A--> J 
    A--> L 
    A--> B 
    A--> C 
    A--> F 
    A--> I 
    A--> J
    A--> K 
    A--> E 
    
    B--> D
    B--> J
    B--> C
    B--> F
    B--> I
    B--> K
    B--> E

    C--> D
    C--> F
    C--> K
    
    D--> H
    D--> F 
    D--> I
    D--> J
    D--> K

    E--> D 
    E--> J
    E--> H 
    E--> C 
    E--> F
    E--> I
    E--> K

    F--> H
    F--> G 
    F--> K

    G--> H

    H--> K

    I--> J
{{< /mermaid >}}

## Supported go versions

We want to always support the **two most recent minor versions of the go compiler**.

However, we try to avoid introducing breaking changes, especially on the more
stable `go-openapi` repos.

When a deprecation or a change comes from the go language or the standard library, we handle this with build tags.

Notice the very important blank line after your build tag comment line.

Example (from `go-openapi/swag`):
```go
//go:build !go1.8

package swag

import "net/url"

func pathUnescape(path string) (string, error) {
	return url.QueryUnescape(path)
}
```
