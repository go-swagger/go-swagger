# Installing

**IMPORTANT NOTE**: `go-swagger` is a tool to mainly generate or analyze source code. In order to make it work after successful
installation, see [the prerequisites](https://goswagger.io/generate/requirements.html) on your development environment.

## Installing from binary distributions

go-swagger releases are distributed as binaries that are built from signed tags. It is published [as github release](https://github.com/go-swagger/go-swagger/tags),
rpm, deb and docker image.

### Docker image [![Docker Repository on Quay](https://quay.io/repository/goswagger/swagger/status "Docker Repository on Quay")](https://quay.io/repository/goswagger/swagger)

First grab the image:

```
docker pull quay.io/goswagger/swagger
```

#### For Mac And Linux users:

```bash
alias swagger='docker run --rm -it  --user $(id -u):$(id -g) -e GOPATH=$(go env GOPATH):/go -v $HOME:$HOME -w $(pwd) quay.io/goswagger/swagger'
swagger version
```

#### For windows users:

```cmd
docker run --rm -it --env GOPATH=/go -v %CD%:/go/src -w /go/src quay.io/goswagger/swagger
```

You can put the following in a file called **swagger.bat** and include it in your path environment variable to act as an alias.

```batch
@echo off
echo.
docker run --rm -it --env GOPATH=/go -v %CD%:/go/src -w /go/src quay.io/goswagger/swagger %*
```

### Homebrew/Linuxbrew

```
brew tap go-swagger/go-swagger
brew install go-swagger
```

### Static binary

You can download a binary for your platform from github:
<https://github.com/go-swagger/go-swagger/releases/latest>

```
download_url=$(curl -s https://api.github.com/repos/go-swagger/go-swagger/releases/latest | \
  jq -r '.assets[] | select(.name | contains("'"$(uname | tr '[:upper:]' '[:lower:]')"'_amd64")) | .browser_download_url')
curl -o /usr/local/bin/swagger -L'#' "$download_url"
chmod +x /usr/local/bin/swagger
```

### Debian packages [ ![Download](https://api.bintray.com/packages/go-swagger/goswagger-debian/swagger/images/download.svg) ](https://bintray.com/go-swagger/goswagger-debian/swagger/_latestVersion)

This repo will work for any debian, the only file it contains gets copied to `/usr/bin`

```
sudo apt install gnupg ca-certificates
sudo apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys 379CE192D401AB61
echo "deb https://dl.bintray.com/go-swagger/goswagger-debian ubuntu main" | sudo tee /etc/apt/sources.list.d/goswagger.list
sudo apt update 
sudo apt install swagger
```

### RPM packages [ ![Download](https://api.bintray.com/packages/go-swagger/goswagger-rpm/swagger/images/download.svg) ](https://bintray.com/go-swagger/goswagger-rpm/swagger/_latestVersion)

This repo should work on any distro that wants rpm packages, the only file it contains gets copied to `/usr/bin`

```
wget https://bintray.com/go-swagger/goswagger-rpm/rpm -O bintray-go-swagger-goswagger-rpm.repo
```


### Installing from source

Install or update from current source master:

```
dir=$(mktemp -d) 
git clone https://github.com/go-swagger/go-swagger "$dir" 
cd "$dir"
go install ./cmd/swagger
```

To install a specific version from source an appropriate tag needs to be checked out first (e.g. `v0.25.0`). Additional `-ldflags` are just to make `swagger version` command print the version and commit id instead of `dev`.

```
dir=$(mktemp -d)
git clone https://github.com/go-swagger/go-swagger "$dir" 
cd "$dir"
git checkout v0.25.0
go install -ldflags "-X github.com/go-swagger/go-swagger/cmd/swagger/commands.Version=$(git describe --tags) -X github.com/go-swagger/go-swagger/cmd/swagger/commands.Commit=$(git rev-parse HEAD)" ./cmd/swagger
```

You are welcome to clone this repo and start contributing:
```
git clone https://github.com/go-swagger/go-swagger
```

> **NOTE**: go-swagger works on *nix as well as Windows OS 
