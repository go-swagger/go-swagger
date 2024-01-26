---
title: Documentation
date: 2023-01-01T01:01:01-08:00
draft: true
weight: 50
---
# Writing documentation

The `go-swagger` documentation site (`goswagger.io`) is built with `HUGO`.
Configuration is in `hack/hugo/hugo.yaml`. The documents root is in `./docs`

Previous releases used to be documented using `gitbooks`.

Assets (images, css, other hugo resources) are located in `hack/doc-site/hugo/themes`.

We systematically copy the repository main `README.md` to `docs/README.md`.
Please make sure links work both from github and gitbook.

There is also a minimal godoc for goswagger, available on pkg.go.dev.

Please make sure new CLI options remain well documented in `./docs/usage`.

## go-openapi repos

Documentation is limited to the repo's README.md and godoc, published on pkg.go.dev.
