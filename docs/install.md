## Installing

### Installing from binary distributions

go-swagger releases are distributed as binaries that are built from signed tags. It is published [as github release](https://github.com/go-swagger/go-swagger/tags),
rpm, deb and docker image.

#### Docker image [![Docker Repository on Quay](https://quay.io/repository/goswagger/swagger/status "Docker Repository on Quay")](https://quay.io/repository/goswagger/swagger)

```
docker pull quay.io/goswagger/swagger

alias swagger="docker run --rm -it -e GOPATH=$HOME/go:/go -v $HOME:$HOME -w $(pwd) quay.io/goswagger/swagger"
swagger version
```

#### Homebrew/Linuxbrew

```
brew tap go-swagger/go-swagger
brew install go-swagger
```

#### Static binary

You can download a binary for your platform from github:
<https://github.com/go-swagger/go-swagger/releases/latest>

```
latestv=$(curl -s https://api.github.com/repos/go-swagger/go-swagger/releases/latest | jq -r .tag_name)
curl -o /usr/local/bin/swagger -L'#' https://github.com/go-swagger/go-swagger/releases/download/$latestv/swagger_$(echo `uname`|tr '[:upper:]' '[:lower:]')_amd64
chmod +x /usr/local/bin/swagger
```

#### Debian packages [ ![Download](https://api.bintray.com/packages/go-swagger/goswagger-debian/swagger/images/download.svg) ](https://bintray.com/go-swagger/goswagger-debian/swagger/_latestVersion)

This repo will work for any debian, the only file it contains gets copied to `/usr/bin`

```
echo "deb https://dl.bintray.com/go-swagger/goswagger-debian ubuntu main" | sudo tee -a /etc/apt/sources.list
```

#### RPM packages [ ![Download](https://api.bintray.com/packages/go-swagger/goswagger-rpm/swagger/images/download.svg) ](https://bintray.com/go-swagger/goswagger-rpm/swagger/_latestVersion)

This repo should work on any distro that wants rpm packages, the only file it contains gets copied to `/usr/bin`

```
wget https://bintray.com/go-swagger/goswagger-rpm/rpm -O bintray-go-swagger-goswagger-rpm.repo
```


### Installing from source

Install or update from current source master:

```
go get -u github.com/go-swagger/go-swagger/cmd/swagger
```

You are welcome to clone this repo and start contributing:
```
cd $GOPATH/src
mkdir -p github.com/go-swagger
cd github.com/go-swagger
git clone https://github.com/go-swagger/go-swagger
```

> **NOTE**: go-swagger works on *nix as well as Windows OS 
