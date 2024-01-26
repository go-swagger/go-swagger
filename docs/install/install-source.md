---
title: Install from source
date: 2023-01-01T01:01:01-08:00
draft: true
weight: 20
description: Installation instructions for go-swagger, from source
---
## Installing from source

### `go install`
If you have `go` version `{{< param goswagger.goVersion >}}` or greater installed, `go-swagger` can be installed by running:

```sh
go install github.com/go-swagger/go-swagger/cmd/swagger@latest
```

## Alternative methods

#### Using a released source tarball

[![GitHub Downloads (all assets, latest release)](https://img.shields.io/badge/dynamic/json?url=https%3A%2F%2Fapi.github.com%2Frepos%2Fgo-swagger%2Fgo-swagger%2Freleases%2Flatest&label=Latest%20tarball&query=%24.tarball_url)](https://github.com/go-swagger/go-swagger/releases/latest)

```sh
dir=$(mktemp -d) 
download_url=$(curl -s https://api.github.com/repos/go-swagger/go-swagger/releases/latest | \
  jq -r '.tarball_url')

curl -o $dir/swagger -L'#' "$download_url"
cd "$dir"
tar xf swagger
cd go-swagger*
go install ./cmd/swagger
```

#### `git clone`

Install or update from current source master:

```sh
dir=$(mktemp -d) 
git clone https://github.com/go-swagger/go-swagger "$dir" 
cd "$dir"
go install ./cmd/swagger
```

To install a specific version from source an appropriate tag needs to be checked out first (e.g. `{{< param goswagger.latestRelease >}}`). Additional `-ldflags` are just to make `swagger version` command print the version and commit id instead of `dev`.

```sh
dir=$(mktemp -d)
git clone https://github.com/go-swagger/go-swagger "$dir" 
cd "$dir"
git checkout {{< param goswagger.latestRelease >}}
go install -ldflags "-X github.com/go-swagger/go-swagger/cmd/swagger/commands.Version=$(git describe --tags) -X github.com/go-swagger/go-swagger/cmd/swagger/commands.Commit=$(git rev-parse HEAD)" ./cmd/swagger
```
