---
title: Install from binary
date: 2023-01-01T01:01:01-08:00
draft: true
weight: 20
description: Installation instructions for go-swagger, from a binary distribution
---

## Installing from binary distributions
{{<tabs "binaryInstallType">}}
{{<tab "MacOS" >}}
### Homebrew

```sh
brew tap go-swagger/go-swagger
brew install go-swagger
```
{{</tab>}}
{{<tab "Linux (LinuxBrew)" >}}
### Linuxbrew

```sh
brew tap go-swagger/go-swagger
brew install go-swagger
```
{{</tab>}}
{{< tab "Linux (Ubuntu, Debian)" >}}
### Debian packages ![debian logo](../icons/debian.png)
[![Download](https://api-prd.cloudsmith.io/v1/badges/version/go-swagger/go-swagger/deb/swagger/latest/a=amd64;d=debian%252Fany-version;t=binary/?render=true&show_latest=true)](https://cloudsmith.io/~go-swagger/repos/go-swagger/packages/detail/deb/swagger/latest/a=amd64;d=debian%252Fany-version;t=binary/)

This package will work for any Debian. The only file it contains gets copied to `/usr/bin`

* Update and install prerequisite packages
```sh
sudo apt update
sudo apt install -y apt-transport-https gnupg curl debian-keyring debian-archive-keyring
```

* Register our GPG signing key
```sh
curl -1sLf 'https://dl.cloudsmith.io/public/go-swagger/go-swagger/gpg.2F8CB673971B5C9E.key' | sudo gpg --dearmor -o /usr/share/keyrings/go-swagger-go-swagger-archive-keyring.gpg
curl -1sLf 'https://dl.cloudsmith.io/public/go-swagger/go-swagger/config.deb.txt?distro=debian&codename=any-version' | sudo tee /etc/apt/sources.list.d/go-swagger-go-swagger.list
```

* Install
```sh
sudo apt update 
sudo apt install swagger
```
{{< /tab >}}

{{<tab "Linux (CentOS, Fedora)" >}}
### RPM packages ![fedora logo](../icons/fedora.png)
[![Download](https://api-prd.cloudsmith.io/v1/badges/version/go-swagger/go-swagger/rpm/swagger/latest/a=x86_64;d=fedora%252Fany-version;t=binary/?render=true&show_latest=true)](https://cloudsmith.io/~go-swagger/repos/go-swagger/packages/detail/rpm/swagger/latest/a=x86_64;d=fedora%252Fany-version;t=binary/)

This package should work on any distro that wants RPM packages. The only file it contains gets copied to `/usr/bin`

* Update and install prerequisite packages
```sh
dnf install -y yum-utils
```

* Register our GPG signing key
```sh
rpm --import 'https://dl.cloudsmith.io/public/go-swagger/go-swagger/gpg.2F8CB673971B5C9E.key'
```

* Install
```sh
curl -1sLf 'https://dl.cloudsmith.io/public/go-swagger/go-swagger/config.rpm.txt?distro=fedora&codename=any-version' > /tmp/go-swagger-go-swagger.repo
dnf config-manager --add-repo '/tmp/go-swagger-go-swagger.repo'
dnf -q makecache -y --disablerepo='*' --enablerepo='go-swagger-go-swagger' --enablerepo='go-swagger-go-swagger-source'
dnf install -y swagger
```
{{</tab>}}
{{</tabs>}}

{{<hint warning>}}
At this moment, we do not support pre-packaged binaries for Windows.
{{</hint>}}

<!-- TODO apk package for alpine -->
<!-- TODO msi package for windows -->

## Static binary

You can also download a binary for your platform from github:

<https://github.com/go-swagger/go-swagger/releases/latest>

{{<hint info>}}
We currently release binary builds for the following platforms:
* darwin AMD64 (MacOS)
* darwin ARM64
* linux AMD64
* linux ARM
* linux ARM64
* linux PPC64le
* linux s390x
* windows AMD64
* window ARM64
{{</hint>}}

```sh
dir=$(mktemp -d) 
download_url=$(curl -s https://api.github.com/repos/go-swagger/go-swagger/releases/latest | \
  jq -r '.assets[] | select(.name | contains("'"$(uname | tr '[:upper:]' '[:lower:]')"'_amd64")) | .browser_download_url')

curl -o $dir/swagger -L'#' "$download_url"
sudo install $dir/swagger /usr/local/bin
```
