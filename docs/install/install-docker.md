---
title: Install from Docker
date: 2023-01-01T01:01:01-08:00
draft: true
weight: 20
description: Installation instructions for go-swagger, from a docker image
---

## Running from Docker

{{<hint "info">}}
We release current master as docker images too.

Images are built for architectures: linux/amd64, linux/arm/v7, linux/arm64 and linux/ppc64le,linux/s390x.
{{</hint>}}

### Docker image

First grab the image:

1. From quay.io [![quay.io image](https://img.shields.io/badge/dynamic/json?url=https%3A%2F%2Fquay.io%2Fapi%2Fv1%2Frepository%2Fgoswagger%2Fswagger%2Ftag%2F%3Flimit%3D1%26onlyActiveTags%3Dtrue%26filter_tag_name%3Dlike%3Av&label=quay.io%20images&query=%24.tags[:1].name)](https://quay.io/repository/goswagger/swagger?tab=tags)

[![quay.io latest master](https://img.shields.io/badge/dynamic/json?url=https%3A%2F%2Fquay.io%2Fapi%2Fv1%2Frepository%2Fgoswagger%2Fswagger%2Ftag%2F%3Flimit%3D1%26onlyActiveTags%3Dtrue%26filter_tag_name%3Dlike%3Amaster&label=quay.io%20images&query=%24.tags[:1].name)](https://quay.io/repository/goswagger/swagger?tab=tags)
```sh
docker pull quay.io/goswagger/swagger
```

or 

2. From Github registry [![ghcr.io image](https://ghcr-badge.egpl.dev/go-swagger/go-swagger/latest_tag?trim=major&ignore=sha-*&label=ghcr.io%20image)](https://github.com/orgs/go-swagger/packages/container/go-swagger/versions)

[![ghcr.io latest master](https://ghcr-badge.egpl.dev/go-swagger/go-swagger/tags?trim=major&ignore=sha-%2A,v%2A,%5B0-9%5D%2A,publish%2A,latest&label=ghcr.io%20image)](https://github.com/orgs/go-swagger/packages/container/go-swagger/versions)

```sh
docker pull ghcr.io/go-swagger/go-swagger
```


{{<tabs "dockerInstallType">}}
{{<tab "MacOS/Linux" >}}
#### For Mac and Linux users

```sh
REPO="quay.io"    #<- or "ghcr.io"
alias swagger='docker run --rm -it  --user $(id -u):$(id -g) -v $HOME:$HOME -w $PWD $REPO/goswagger/swagger'
swagger version
```
{{</tab>}}
{{<tab "Windows" >}}
#### For windows users

```cmd
REM <- or "ghcr.io"
set REPO=quay.io
docker run --rm -it  -v %CD%:/app -w /app %REPO%/goswagger/swagger
```

You can put the following in a file called **swagger.bat** and include it in your path environment variable to act as an alias.

```cmd
@echo off
REM <- or "ghcr.io"
set REPO=quay.io
echo.
docker run --rm -it -v %CD%:/app -w /app %REPO%/goswagger/swagger %*
```
{{</tab>}}

{{</tabs>}}
