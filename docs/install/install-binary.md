---
title: Install from static binary
date: 2023-01-01T01:01:01-08:00
draft: true
weight: 20
description: Installation instructions for go-swagger, from a static binary
---

## Static binary

You can also download a binary for your platform from github:

<https://github.com/go-swagger/go-swagger/releases/latest>

{{<hint info>}}
We currently release binary builds for the following platforms:

* MacOS (darwin)
  * AMD64 (x86_64)
  * ARM64
* Linux
  * AMD64 (x86_64)
  * ARM (v6)
  * ARM64
  * PPC64le
  * s390x
* Windows
  * AMD64 (x86_64)

{{</hint>}}

Binaries are compressed with UPX (when supported by the platform) and are available as tarballs:
```
swagger_${version}_${os}_${arch}.tar.{gz|zip}
```

Where version is `{major}.{minor}.{patch}`, os is one of `Linux`, `Darwin`, `Windows`, and arch is one of
`x86_64` (for `x86_64/amd64` architectures), `arm64`, `armv6`, `ppc64le`, `os390s`.

Compression is `zip` on Windows.

For Linux/Darwin:

```sh
dir=$(mktemp -d) 
download_url=$(\
  curl -s https://api.github.com/repos/go-swagger/go-swagger/releases/latest | \
  jq -r --arg arch $(arch) --arg os $(uname -o |cut -d'/' -f 2) '.assets[] | select(.name | test($os+"_"+$arch+"\\.tar\\.gz$"))|.browser_download_url' \
)
curl -o "$dir"/swagger.tar.gz -L'#' "$download_url"
(cd "$dir" && tar xvf swagger.tar.gz swagger)
sudo install "$dir"/swagger /usr/local/bin
/usr/local/bin/swagger version
```
