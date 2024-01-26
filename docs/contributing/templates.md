---
title: Contributing templates
date: 2023-01-01T01:01:01-08:00
draft: true
description: Contributing guide - Templates
weight: 30
---
# Maintaining templates

For its code generation features, `go-swagger` uses a bunch of go text/templates.

You'll find them all there: https://github.com/go-swagger/go-swagger/blob/master/generator/templates

The `go-swagger` executable is built with an in-memory image of templates.
Binary encoded assets are auto-generated from the `generator/templates` directory using `go:embed`.

Most templates can be overriden at run time with a config setup.

> **NOTE**: we are carrying out unit tests on codegen mostly by asserting lines in generated code.
> There is a bunch of test utility functions for this. See `generator/*_test.go`.
>
> If you want to bring in more advanced testing go programs with your fixtures, please tag
> those so they don't affect the `go ./...` command (e.g. with `// +build +integration`).

