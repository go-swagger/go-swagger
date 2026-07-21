# Swagger 2.0

<!-- Badges: status  -->
[![Tests][test-badge]][test-url] [![Coverage][cov-badge]][cov-url] [![CI vuln scan][vuln-scan-badge]][vuln-scan-url] [![CodeQL][codeql-badge]][codeql-url]
<!-- Badges: release & docker images  -->
[![Release][release-badge]][release-url] [![Container Registry on Quay.io][quay-badge]][quay-url] [![Container Registry on Github][ghcr-badge]][ghcr-url]
<!-- Badges: code quality  -->
[![Go Report Card][gocard-badge]][gocard-url] [![CodeFactor Grade][codefactor-badge]][codefactor-url]
<!-- Badges: license & compliance -->
[![License][license-badge]][license-url] [![Open SSF Scorecard][ossf-badge]][ossf-url] [![OpenSSF Best Practices][ossf-cci-badge]][ossf-cci-url] [![OpenSSF Baseline][ossf-baseline]][ossf-cci-url] [![OSS licences status][fossa-badge]][fossa-url]
<!-- Badges: documentation & support -->
<!-- Badges: others & stats -->
[![Documentation][doc-badge]][doc-url] [![GoDoc][godoc-badge]][godoc-url] [![Discord Channel][discord-badge]][discord-url] [![go version][goversion-badge]][goversion-url] ![Top language][top-badge] ![Commits since latest release][commits-badge]

---

This project contains a golang implementation of Swagger 2.0 (aka [OpenAPI 2.0](https://github.com/OAI/OpenAPI-Specification/blob/old-v3.2.0-dev/versions/2.0.md)).
It provide tools to work with swagger specifications.

[Swagger](https://swagger.io/) is a simple yet powerful representation of your RESTful API.<br>

## Announcements

You may join the discord community by clicking the invite link on the discord badge. [![Discord Channel][discord-badge]][discord-url].

* **2026-07-21** : v0.35.2 is out
  * security release: with this release, we have completed the hardening of go-openapi libraries and how go-swagger
    consumes those.
  * The class of threats that we've tried to mitigate is "adverserial specs", e.g. openapi documents that go-swagger
     would happily parse, resolve and generate code from, but that are crafted maliciously.
  * Attack vectors can be: `$ref` (to resolve against a malicious site, to expand local files and try to leak, or simply
    crash), to induce the code generator into producing malicious code.
  * For users of go-openapi libraries such as `loads` and `swag/loading`, the defaults remain unrestricted: secure options
    are available to developer to guard spec load against possibly unsafe input documents.
  * `go-openapi/validate` now warns about dubious `$ref`'s such as multiple remote hosts, or absolute local paths not
    beneath the base path.
  * `go-swagger` codegen commands may now be restricted to operate within a rooted workspace (no local $ref is resolved
    outside) or using a restricted client.

* **2026-07-20** : v0.36.0 will land in July (soon!)
  * **documentation**: restyle doc site like go-openapi doc sites, dedicated doc site for examples, cover a significant
    part of doc-related issues.
  * **spec generation**: codescan will publish its own lightweight CLI at a faster pace, as well as a TUI tool to
    instantly check how you annotated code looks like as a spec.
    (preview: <https://github.com/go-swagger/go-swagger/issues/3372#issuecomment-4733107554>). We hope we'll be able
    to land a Web playground on a similar principle (WASI build on top of codescan). go-swagger will still receive updates.
  * **code generation**: we'll try our best to land a few requested enhancements among the 50-60 reachable ones.
    (most issues in codegen now have hit an "architecture wall": work has started on a v2 to overcome these limitations).
  * our Slack channel is now closed and superseded by the [![Discord Channel][discord-badge]][discord-url].

* **2026-07-20** : v0.35.1 landed!
  * re-instated binary release for windows ARM64
  * **spec generation**: another round of fixes. Another small lots of features for more control over your rendered specs.
    Check out the [documentation site dedicated to spec generation][codescan-doc-url].

* **2026-06-22** : v0.35.0 landed! (end of June)
  * **code generation**: security fixes that prevent generated code to produce code injected from an erroneous
    or malicious spec. swagger validate now warns about possibly harmful $ref (e.g. from multiple origins).
  * **spec generation**: major bug-bashing action on go-openapi/codescan (which has eventually become fixable...).
    v0.35.0 closes ~200+ "generate spec" issues: bug fixes and requested enhancements.
  * Spec generation now produces a detailed diagnostic of how your code annotations may be misinterpreted.
    A complete [documentation site][codescan-doc-url] is now published.
  * Please check it out from master (or dev docker image). Your feedback is super important!

* **2026-05-28** : v0.34.0 ships!
  * **major refactoring actions**: the repo has been split in smaller chunks, easier to understand:
    * code examples have moved to `go-swagger/examples`, with a CI to automate code regeneration
    * `codescan` (the part that underpins `swagger generate spec`) has moved as a standalone
      library `go-openapi/codescan`. This library has been heavily refactored to prepare more
      significant improvements. The current version gets a few quirks already fixed.
    * `diff` (the implement of `swagger diff`) has joined `go-openapi/analysis`
    * package `generator` has been refactored to expose internal utilities (template repo, funcmaps, etc)
      as packages.
  * Generated code now requires `go-openapi/runtime` v0.32.x and benefits from many bug fixes and there too a
    rechunking of the code.
  * Generated code now requires `go-openapi/swag` v0.26+ and directly imports all sub-modules.
  * Many long-awaited improvements on the generated client.

## Documentation

<https://goswagger.io>

##  Features

`go-swagger` brings to the go community a complete suite of fully-featured, high-performance, API components to  work with a Swagger API: server, client and data model.

* Generates a server from a swagger specification
* Generates a client from a swagger specification
* Generates a CLI (command line tool) from a swagger specification (alpha stage)
* Supports most features offered by jsonschema and swagger, including polymorphism
* Generates a swagger specification from annotated go code
* Additional tools to work with a swagger spec
* Great customization features, with vendor extensions and customizable templates

Our focus with code generation is to produce idiomatic, fast go code, which plays nice with golint, go vet etc.

[More details](https://goswagger.io/go-swagger/features).

##  Project status

This project supports OpenAPI 2.0. **At this moment it does not support OpenAPI 3.x**.

`go-swagger` is now feature complete and has stabilized its API.

Most features and building blocks are now in a stable state, with a rich set of CI tests.

The go-openapi community actively continues bringing fixes and enhancements to this code base.

There is still much room for improvement: contributors and PR's are welcome.
You may also get in touch with maintainers on our [![Discord Channel][discord-badge]][discord-url].
## Installing

```sh
go install github.com/go-swagger/go-swagger/cmd/swagger@latest
```

`go-swagger` is also available as binary or docker releases as well as from source: [more details](https://goswagger.io/go-swagger/install).

## Try it

Try `go-swagger` in a free online workspace using Gitpod:

[![Open in Gitpod](https://gitpod.io/button/open-in-gitpod.svg)](https://gitpod.io#https://github.com/go-swagger/go-swagger)

## Security

`go-swagger` turns an OpenAPI 2.0 specification into source code. **Treat a specification like any other untrusted input: if you obtained it from a remote or untrusted location, review its contents before generating code from it.**

The generator never executes the spec, and the generated code runs only when *you* build and import it. We have hardened the generators against an adversarial spec that tries to inject unwanted Go into the artifacts it produces — identifiers, struct tags, doc comments and CLI string literals are sanitized or escaped — which substantially reduces the exposure. It is not, however, a substitute for reviewing what you generate. In particular:

- **Remote `$ref`s.** A spec may reference other documents, possibly over the network. Those references are resolved and folded into the generated code, so inspect any external reference you do not control.
- **The `x-go-type` extension.** By design, this extension lets the spec choose the Go type for a field — including an arbitrary imported package. That capability *cannot easily be safeguarded*: a spec using `x-go-type` can make your generated code import and depend on a package of its choosing. Always review specs that rely on it.

When in doubt, generate into a scratch directory, read the diff, and only then wire it into your build.

## Licensing

The toolkit itself is licensed under an Apache Software License 2.0:
[SPDX-License-Identifier: Apache-2.0](./LICENSE).

Just like swagger, this does not cover code generated by the toolkit. That code is entirely yours to license however you see fit.

### Licence scan on dependencies

[![FOSSA Status][fossa-badge-large]][fossa-url-large]

<!-- Badges: status  -->
[test-badge]: https://github.com/go-swagger/go-swagger/actions/workflows/test.yaml/badge.svg
[test-url]: https://github.com/go-swagger/go-swagger/actions/workflows/test.yaml
[cov-badge]: https://codecov.io/gh/go-swagger/go-swagger/branch/master/graph/badge.svg
[cov-url]: https://codecov.io/gh/go-swagger/go-swagger
[vuln-scan-badge]: https://github.com/go-swagger/go-swagger/actions/workflows/scanner.yaml/badge.svg
[vuln-scan-url]: https://github.com/go-swagger/go-swagger/actions/workflows/scanner.yaml
[codeql-badge]: https://github.com/go-swagger/go-swagger/actions/workflows/codeql.yaml/badge.svg
[codeql-url]: https://github.com/go-swagger/go-swagger/actions/workflows/codeql.yaml
<!-- Badges: release & docker images  -->
[release-badge]: https://badge.fury.io/gh/go-swagger%2Fgo-swagger.svg
[release-url]: https://badge.fury.io/gh/go-swagger%2Fgo-swagger
[quay-badge]: https://img.shields.io/badge/dynamic/json?url=https%3A%2F%2Fquay.io%2Fapi%2Fv1%2Frepository%2Fgoswagger%2Fswagger%2Ftag%2F%3Flimit%3D1%26onlyActiveTags%3Dtrue%26filter_tag_name%3Dlike%3Av&label=Container%20Registry%20on%20quay.io&query=%24.tags[:1].name&logo=redhatopenshift&logoColor=#EE0000?color=green
[quay-url]: https://quay.io/repository/goswagger/swagger?tab=tags
[ghcr-badge]: https://ghcr-badge-ipv2.onrender.com/go-swagger/go-swagger/latest_tag?ignore=sha-*,edge,master&label=Container%20Registry%20on%20Github
[ghcr-url]: https://github.com/orgs/go-swagger/packages/container/go-swagger/versions?filters[version_type]=tagged
<!-- Badges: code quality  -->
[gocard-badge]: https://goreportcard.com/badge/github.com/go-swagger/go-swagger
[gocard-url]: https://goreportcard.com/report/github.com/go-swagger/go-swagger
[codefactor-badge]: https://img.shields.io/codefactor/grade/github/go-swagger/go-swagger
[codefactor-url]: https://www.codefactor.io/repository/github/go-swagger/go-swagger
<!-- Badges: documentation & support -->
[doc-badge]: https://img.shields.io/badge/doc-site-blue?link=https%3A%2F%2Fgoswagger.io%2Fgo-swagger%2F
[doc-url]: https://goswagger.io/go-swagger
[godoc-badge]: https://godoc.org/github.com/go-swagger/go-swagger?status.svg
[godoc-url]: http://godoc.org/github.com/go-swagger/go-swagger
[discord-badge]: https://img.shields.io/discord/1446918742398341256?logo=discord&label=discord&color=blue
[discord-url]: https://discord.gg/FfnFYaC3k5
[codescan-doc-url]: https://go-openapi.github.io/codescan/
<!-- Badges: license & compliance -->
[license-badge]: http://img.shields.io/badge/license-Apache%20v2-orange.svg
[license-url]: https://github.com/go-swagger/go-swagger/?tab=Apache-2.0-1-ov-file#readme
[ossf-badge]: https://api.securityscorecards.dev/projects/github.com/go-swagger/go-swagger/badge
[ossf-url]: https://securityscorecards.dev/viewer/?uri=github.com/go-swagger/go-swagger
[ossf-cci-badge]: https://www.bestpractices.dev/projects/11359/badge
[ossf-cci-url]: https://www.bestpractices.dev/projects/11359
[ossf-baseline]: https://www.bestpractices.dev/projects/11359/baseline
[fossa-badge]: https://app.fossa.io/api/projects/git%2Bgithub.com%2Fgo-swagger%2Fgo-swagger.svg?type=shield
[fossa-url]: https://app.fossa.io/projects/git%2Bgithub.com%2Fgo-swagger%2Fgo-swagger?ref=badge_shield
[fossa-badge-large]: https://app.fossa.io/api/projects/git%2Bgithub.com%2Fgo-swagger%2Fgo-swagger.svg?type=large
[fossa-url-large]: https://app.fossa.io/projects/git%2Bgithub.com%2Fgo-swagger%2Fgo-swagger?ref=badge_large
[oai-url]: https://raw.githubusercontent.com/swagger-api/swagger-spec/master/LICENSE
<!-- Badges: others & stats -->
[goversion-badge]: https://img.shields.io/github/go-mod/go-version/go-swagger/go-swagger
[goversion-url]: https://github.com/go-swagger/go-swagger/blob/master/go.mod
[top-badge]: https://img.shields.io/github/languages/top/go-swagger/go-swagger
[commits-badge]: https://img.shields.io/github/commits-since/go-swagger/go-swagger/latest
