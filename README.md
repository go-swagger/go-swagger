# OpenAPI 3.0 Migration Fork [![Run CI](https://github.com/go-swagger/go-swagger/actions/workflows/test.yaml/badge.svg)](https://github.com/go-swagger/go-swagger/actions/workflows/test.yaml) [![codecov](https://codecov.io/gh/go-swagger/go-swagger/branch/master/graph/badge.svg)](https://codecov.io/gh/go-swagger/go-swagger)[![Go Report Card](https://goreportcard.com/badge/github.com/go-swagger/go-swagger)](https://goreportcard.com/report/github.com/go-swagger/go-swagger)

[![GitHub version](https://badge.fury.io/gh/go-swagger%2Fgo-swagger.svg)](https://badge.fury.io/gh/go-swagger%2Fgo-swagger) [![Docker Repository on Quay](https://img.shields.io/badge/dynamic/json?url=https%3A%2F%2Fquay.io%2Fapi%2Fv1%2Frepository%2Fgoswagger%2Fswagger%2Ftag%2F%3Flimit%3D1%26onlyActiveTags%3Dtrue%26filter_tag_name%3Dlike%3Av&label=Docker%20Repository%20on%20Quay&query=%24.tags[:1].name)](https://quay.io/repository/goswagger/swagger?tab=tags) [![Docker Repository on Github](https://ghcr-badge.egpl.dev/go-swagger/go-swagger/latest_tag?trim=major&ignore=sha-*&label=Docker%20Repository%20on%20Github)](https://github.com/orgs/go-swagger/packages/container/go-swagger/versions) ![GitHub commits since latest release](https://img.shields.io/github/commits-since/go-swagger/go-swagger/latest)

[![Slack Status](https://slackin.goswagger.io/badge.svg)](https://slackin.goswagger.io)
[![license](http://img.shields.io/badge/license-Apache%20v2-orange.svg)](https://raw.githubusercontent.com/swagger-api/swagger-spec/master/LICENSE)
[![GoDoc](https://godoc.org/github.com/go-swagger/go-swagger?status.svg)](http://godoc.org/github.com/go-swagger/go-swagger)

[![Open SSF Scorecard](https://api.securityscorecards.dev/projects/github.com/go-swagger/go-swagger/badge)](https://securityscorecards.dev/viewer/?uri=github.com/go-swagger/go-swagger)
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fgo-swagger%2Fgo-swagger.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fgo-swagger%2Fgo-swagger?ref=badge_shield)

---

## About This Fork

This is a fork of the original [go-swagger](https://github.com/go-swagger/go-swagger) project with the specific goal of **migrating from Swagger 2.0/OpenAPI 2.0 to OpenAPI 3.0 specification support**.

### Migration Goals

- Update the codebase to support [OpenAPI 3.0](https://spec.openapis.org/oas/v3.0.3) specification
- Maintain backward compatibility where possible
- Preserve the existing feature set while adding OpenAPI 3.0 capabilities
- Update code generation to work with OpenAPI 3.0 schemas

### Original Project

This project was originally a golang implementation of Swagger 2.0 (aka [OpenAPI 2.0](https://github.com/OAI/OpenAPI-Specification/blob/old-v3.2.0-dev/versions/2.0.md)) that provided tools to work with swagger specifications.

[Swagger](https://swagger.io/) is a simple yet powerful representation of your RESTful API.<br>

## Documentation

Original documentation: <https://goswagger.io>

##  Features

`go-swagger` brings to the go community a complete suite of fully-featured, high-performance, API components to work with a Swagger API: server, client and data model.

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

⚠️ **Migration in Progress**: This fork is actively being updated to support OpenAPI 3.0 specification.

The original project supported OpenAPI 2.0 only. **This fork aims to add OpenAPI 3.x support** while maintaining the existing feature set.

The migration is a work in progress. The go-openapi community actively continues bringing fixes and enhancements to the original code base, and this fork will incorporate relevant changes while focusing on OpenAPI 3.0 support.

Contributors and PR's are welcome for the migration effort. You may also get in touch with maintainers on [our slack channel](https://slackin.goswagger.io).

## Installing
