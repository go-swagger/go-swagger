---
title: Guidelines
date: 2023-01-01T01:01:01-08:00
draft: true
description: Contributing guide - Code style
weight: 20
---
# Linting

Our CI run `golangci-lint` to enforce a few linting rules defined there: https://github.com/go-swagger/go-swagger/blob/master/.golangci.yml

Before you push, check your work with the golangci meta linter:
`go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest`
`golangci-lint run --new-from-rev HEAD`

We try to remain consistent across repositories regarding linting rules.

# Vendoring

All `go-openapi` and `go-swagger` repositories have adopted go modules and are no longer using vendoring.
