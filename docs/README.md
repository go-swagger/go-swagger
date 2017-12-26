# Swagger 2.0 [![Build Status](https://circleci.com/gh/go-swagger/go-swagger.svg?style=shield)](https://circleci.com/gh/go-swagger/go-swagger) [![Build status](https://ci.appveyor.com/api/projects/status/x377t5o9ennm847o/branch/master?svg=true)](https://ci.appveyor.com/project/casualjim/go-swagger/branch/master) [![codecov](https://codecov.io/gh/go-swagger/go-swagger/branch/master/graph/badge.svg)](https://codecov.io/gh/go-swagger/go-swagger) [![Slack Status](https://slackin.goswagger.io/badge.svg)](https://slackin.goswagger.io)

[![license](http://img.shields.io/badge/license-Apache%20v2-orange.svg)](https://raw.githubusercontent.com/swagger-api/swagger-spec/master/LICENSE) [![GoDoc](https://godoc.org/github.com/go-swagger/go-swagger?status.svg)](http://godoc.org/github.com/go-swagger/go-swagger) [![GitHub version](https://badge.fury.io/gh/go-swagger%2Fgo-swagger.svg)](https://badge.fury.io/gh/go-swagger%2Fgo-swagger) [![Docker Repository on Quay](https://quay.io/repository/goswagger/swagger/status "Docker Repository on Quay")](https://quay.io/repository/goswagger/swagger)

Development of this toolkit is sponsored by VMware<br>[![VMWare](https://avatars2.githubusercontent.com/u/473334?v=3&s=200)](https://vmware.github.io)  

Contains an implementation of Swagger 2.0 (aka [OpenAPI 2.0](https://github.com/OAI/OpenAPI-Specification/blob/master/versions/2.0.md)). It knows how to serialize and deserialize swagger specifications.

[Swagger](https://swagger.io/) is a simple yet powerful representation of your RESTful API.<br>With the largest ecosystem of API tooling on the planet, thousands of developers are supporting Swagger in almost every modern programming language and deployment environment.

With a Swagger-enabled API, you get interactive documentation, client SDK generation and discoverability. We created Swagger to help fulfill the promise of APIs.

Swagger helps companies like Apigee, Getty Images, Intuit, LivingSocial, McKesson, Microsoft, Morningstar, and PayPal build the best possible services with RESTful APIs. Now in version 2.0, Swagger is more enabling than ever. And it's 100% open source software.

## How is this different from go generator in swagger-codegen?

**tl;dr** The main difference at this moment is that this one will actually work.

The swagger-codegen project only generates a client and even there it will only support flat models.

* This project supports most features offered by jsonschema including polymorphism.
* It allows for generating a swagger specification from go code.
* It allows for generating a server from a swagger definition and to generate an equivalent spec back from that codebase.
* It allows for generating a client from a swagger definition.
* It has support for several common swagger vendor extensions.

Why is this not done in the swagger-codegen project? Because:

* I don't really know java very well and so I'd be learning both java and the object model of the codegen which was in heavy flux as opposed to doing go and I really wanted to go experience of designing a large codebase with it.
* Go's super limited type system makes it so that it doesn't fit well in the model of swagger-codegen
* Go's idea of polymorphism doesn't reconcile very well with a solution designed for languages that actually have inheritance and so forth.
* For supporting types like `[][][]map[string][][]int64` I don't think it's possible with mustache
* I gravely underestimated the amount of work that would be involved in making something useful out of it.
* My personal mission: I want the jvm to go away, it was great way back when now it's just silly (vm in container on vm in vm in container)

## What's inside?

Here is an outline of available features (see the full list [here](https://goswagger.io/features.html)):

- [x] An object model that serializes to swagger-compliant yaml or json
- [x] A tool to work with swagger

  - [x] Serve swagger UI for any swagger spec file
  - [x] Flexible code generation, with customizable templates

    - [x] Generate API based on swagger spec
    - [x] Generate go client from a swagger spec

  - [x] Validate a swagger spec document, with extra rules outlined [here](https://github.com/apigee-127/sway/blob/master/docs/README.md#semantic-validation)
  - [x] Validate against jsonschema (Draft 4), with full $ref support
  - [x] Generate spec document based on annotated code

- [x] Middlewares

  - [x] serve spec
  - [x] routing
  - [x] validation
  - [x] additional validation through an interface
  - [x] authorization
  - [x] swagger docs UI

- [x] Typed JSON Schema implementation
- [x] Extended string and numeric formats
- [x] Project documentation site
- [x] Play nice with golint, go vet etc.

## Documentation

<https://goswagger.io>

## Installing

### Installing from binary distributions

go-swagger releases are distributed as binaries that are built from signed tags. It is published [as github release](https://github.com/go-swagger/go-swagger/tags),
rpm, deb and docker image.

#### Docker image

```
docker pull quay.io/goswagger/swagger

alias swagger="docker run --rm -it -v $HOME:$HOME -w $(pwd) quay.io/goswagger/swagger"
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

Install or update from source:

```
go get -u github.com/go-swagger/go-swagger/cmd/swagger
```

## Use-cases
The main package of the toolkit, go-swagger/go-swagger, provides a number of command line tools to help working with swagger.

The toolkit is highly customizable and allows endless possibilities to work with OpenAPI2.0 specifications.

Beside the go-swagger CLI tool and generator, the [go-openapi packages](https://github.com/go-openapi) provide modular functionality to build custom solutions
on top of OpenAPI.

The CLI supports shell autocompletion utilities: see [here](https://goswagger.io/cli_helpers.html).

### Serve spec UI
Most basic use-case: serve a UI for your spec:

```
swagger serve https://raw.githubusercontent.com/swagger-api/swagger-spec/master/examples/v2.0/json/petstore-expanded.json
```

### Validate an OpenAPI 2.0 spec
To [validate](https://goswagger.io/usage/validate.html) a Swagger specification:

```
swagger validate https://raw.githubusercontent.com/swagger-api/swagger-spec/master/examples/v2.0/json/petstore-expanded.json
```

### Generate an API server
To generate a [server for a swagger spec](https://goswagger.io/generate/server.html) document:

```
swagger generate server [-f ./swagger.json] -A [application-name [--principal [principal-name]]
```

### Generate an API client
To generate a [client for a swagger spec](https://goswagger.io/generate/client.html) document:

```
swagger generate client [-f ./swagger.json] -A [application-name [--principal [principal-name]]
```

### Generate a spec from source
To generate a [swagger spec document for a go application](https://goswagger.io/generate/spec.html):

```
swagger generate spec -o ./swagger.json
```

### Generate a data model
To generate model structures and validators exposed by the API:

```
swagger generate model
```

### Transform specs

Resolve and expand $ref's in your spec as inline definitions:
```
swagger expand {spec}
```

Flatten you spec: all external $ref's are imported into the main document and inline schemas reorganized as definitions.
```
swagger flatten {spec}
```

Merge specifications (composition):
```
swagger mixin {spec1} {spec2}
```

## Note to users migrating from older releases

### Using 0.5.0

Because 0.5.0 and master have diverged significantly, you should checkout the tag 0.5.0 for go-swagger when you use the currently released version.

### Migrating from 0.5.0 to 0.6.0

You will have to rename some imports:

```
github.com/go-swagger/go-swagger/httpkit/validate to github.com/go-openapi/validate
github.com/go-swagger/go-swagger/httpkit to github.com/go-openapi/runtime
github.com/naoina/denco to github.com/go-openapi/runtime/middleware/denco
github.com/go-swagger/go-swagger to github.com/go-openapi
```

## Licensing

The toolkit itself is licensed as Apache Software License 2.0. Just like swagger, this does not cover code generated by the toolkit. That code is entirely yours to license however you see fit.

## FAQ

<https://goswagger.io/faq/>
