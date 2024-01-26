---
title: Install
date: 2023-01-01T01:01:01-08:00
draft: true
weight: 20
---
# Installing go-swagger

{{<hint warning>}}
**IMPORTANT NOTE**: `go-swagger` is a tool to mainly generate or analyze source code. In order to make it work after successful
installation, see [the prerequisites](https://goswagger.io/generate/requirements.html) on your development environment.
{{</hint>}}

## Installing from binary distributions

go-swagger releases are distributed as binaries that are built from signed tags. It is published [as github release](https://github.com/go-swagger/go-swagger/tags),
rpm, deb and docker image.

{{<tabs "installType">}}
{{<tab "MacOS" >}}
### Docker image [![Docker Repository on Quay](https://quay.io/repository/goswagger/swagger/status "Docker Repository on Quay")](https://quay.io/repository/goswagger/swagger)

First grab the image:

```sh
docker pull quay.io/goswagger/swagger
```

or 

```sh
docker pull ghcr.io/go-swagger/go-swagger
```

#### For Mac And Linux users

```sh
alias swagger='docker run --rm -it  --user $(id -u):$(id -g) -v $HOME:$HOME -w $PWD quay.io/goswagger/swagger'
swagger version
```

or 

```sh
alias swagger='docker run --rm -it  --user $(id -u):$(id -g) -v $HOME:$HOME -w $PWD ghcr.io/go-swagger/go-swagger'
swagger version
```

#### For windows users

```cmd
docker run --rm -it  -v %CD%:/app -w /app quay.io/goswagger/swagger
```

or

```cmd
docker run --rm -it  -v %CD%:/app -w /app ghcr.io/go-swagger/go-swagger
```

You can put the following in a file called **swagger.bat** and include it in your path environment variable to act as an alias.

```cmd
@echo off
echo.
docker run --rm -it -v %CD%:/app -w /app quay.io/goswagger/swagger %*
```

or

```cmd
@echo off
echo.
docker run --rm -it -v %CD%:/app -w /app ghcr.io/go-swagger/go-swagger %*
```
{{</tab>}}

{{<tab "Linux" >}}
### Homebrew/Linuxbrew

```sh
brew tap go-swagger/go-swagger
brew install go-swagger
```

### Debian packages [![Download](https://api-prd.cloudsmith.io/v1/badges/version/go-swagger/go-swagger/deb/swagger/latest/a=amd64;d=debian%252Fany-version;t=binary/?render=true&show_latest=true)](https://cloudsmith.io/~go-swagger/repos/go-swagger/packages/detail/deb/swagger/latest/a=amd64;d=debian%252Fany-version;t=binary/)

This repo will work for any debian, the only file it contains gets copied to `/usr/bin`

without sudo:

```sh
apt update
apt install -y apt-transport-https gnupg curl debian-keyring debian-archive-keyring
curl -1sLf 'https://dl.cloudsmith.io/public/go-swagger/go-swagger/gpg.2F8CB673971B5C9E.key' | gpg --dearmor -o /usr/share/keyrings/go-swagger-go-swagger-archive-keyring.gpg
curl -1sLf 'https://dl.cloudsmith.io/public/go-swagger/go-swagger/config.deb.txt?distro=debian&codename=any-version' > /etc/apt/sources.list.d/go-swagger-go-swagger.list
apt update 
apt install swagger
```

with sudo:

```sh
sudo apt update
sudo apt install -y apt-transport-https gnupg curl debian-keyring debian-archive-keyring
curl -1sLf 'https://dl.cloudsmith.io/public/go-swagger/go-swagger/gpg.2F8CB673971B5C9E.key' | sudo gpg --dearmor -o /usr/share/keyrings/go-swagger-go-swagger-archive-keyring.gpg
curl -1sLf 'https://dl.cloudsmith.io/public/go-swagger/go-swagger/config.deb.txt?distro=debian&codename=any-version' | sudo tee /etc/apt/sources.list.d/go-swagger-go-swagger.list
sudo apt update 
sudo apt install swagger
```

### RPM packages [![Download](https://api-prd.cloudsmith.io/v1/badges/version/go-swagger/go-swagger/rpm/swagger/latest/a=x86_64;d=fedora%252Fany-version;t=binary/?render=true&show_latest=true)](https://cloudsmith.io/~go-swagger/repos/go-swagger/packages/detail/rpm/swagger/latest/a=x86_64;d=fedora%252Fany-version;t=binary/)

This repo should work on any distro that wants rpm packages, the only file it contains gets copied to `/usr/bin`

```sh
dnf install -y yum-utils
rpm --import 'https://dl.cloudsmith.io/public/go-swagger/go-swagger/gpg.2F8CB673971B5C9E.key'
curl -1sLf 'https://dl.cloudsmith.io/public/go-swagger/go-swagger/config.rpm.txt?distro=fedora&codename=any-version' > /tmp/go-swagger-go-swagger.repo
dnf config-manager --add-repo '/tmp/go-swagger-go-swagger.repo'
dnf -q makecache -y --disablerepo='*' --enablerepo='go-swagger-go-swagger' --enablerepo='go-swagger-go-swagger-source'
dnf install -y swagger
```
{{</tab>}}

{{<tab "Windows" >}}
{{</tab>}}
{{</tabs>}}

## Installing from source

If you have `go` version `1.16` or greater installed the binary  can be installed by running:

```sh
go install github.com/go-swagger/go-swagger/cmd/swagger@latest
```

### Static binary

You can download a binary for your platform from github:
<https://github.com/go-swagger/go-swagger/releases/latest>

```sh
download_url=$(curl -s https://api.github.com/repos/go-swagger/go-swagger/releases/latest | \
  jq -r '.assets[] | select(.name | contains("'"$(uname | tr '[:upper:]' '[:lower:]')"'_amd64")) | .browser_download_url')
curl -o /usr/local/bin/swagger -L'#' "$download_url"
chmod +x /usr/local/bin/swagger
```

### Installing from source

Install or update from current source master:

```sh
dir=$(mktemp -d) 
git clone https://github.com/go-swagger/go-swagger "$dir" 
cd "$dir"
go install ./cmd/swagger
```

To install a specific version from source an appropriate tag needs to be checked out first (e.g. `v0.25.0`). Additional `-ldflags` are just to make `swagger version` command print the version and commit id instead of `dev`.

```sh
dir=$(mktemp -d)
git clone https://github.com/go-swagger/go-swagger "$dir" 
cd "$dir"
git checkout v0.25.0
go install -ldflags "-X github.com/go-swagger/go-swagger/cmd/swagger/commands.Version=$(git describe --tags) -X github.com/go-swagger/go-swagger/cmd/swagger/commands.Commit=$(git rev-parse HEAD)" ./cmd/swagger
```

You are welcome to clone this repo and start contributing:
```sh
git clone https://github.com/go-swagger/go-swagger
```

> **NOTE**: go-swagger works on *nix as well as Windows OS 
