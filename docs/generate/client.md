---
title: Generating a client
date: 2023-01-01T01:01:01-08:00
draft: true
weight: 20
---
# Generate an API client from a swagger spec

The toolkit has a command that will let you generate a client.

### Client usage

```
Usage:
  swagger [OPTIONS] generate client [client-OPTIONS]

generate all the files for a client library

Application Options:
  -q, --quiet                                                                     silence logs
      --log-output=LOG-FILE                                                       redirect logs to file

Help Options:
  -h, --help                                                                      Show this help message

[client command options]
      -c, --client-package=                                                       the package to save the client specific code (default: client)
      -P, --principal=                                                            the model to use for the security principal
          --default-scheme=                                                       the default scheme for this API (default: http)
          --default-produces=                                                     the default mime type that API operations produce (default: application/json)
          --default-consumes=                                                     the default mime type that API operations consume (default: application/json)
          --skip-models                                                           no models will be generated when this flag is specified
          --skip-operations                                                       no operations will be generated when this flag is specified
      -A, --name=                                                                 the name of the application, defaults to a mangled value of info.title

    Options common to all code generation commands:
      -f, --spec=                                                                 the spec file to use (default swagger.{json,yml,yaml})
      -t, --target=                                                               the base directory for generating the files (default: ./)
          --template=[stratoscale]                                                load contributed templates
      -T, --template-dir=                                                         alternative template override directory
      -C, --config-file=                                                          configuration file to use for overriding template options
      -r, --copyright-file=                                                       copyright file used to add copyright header
          --additional-initialism=                                                consecutive capitals that should be considered intialisms
          --allow-template-override                                               allows overriding protected templates
          --skip-validation                                                       skips validation of spec prior to generation
          --dump-data                                                             when present dumps the json for the template generator instead of generating files
          --with-expand                                                           expands all $ref's in spec prior to generation (shorthand to --with-flatten=expand)
          --with-flatten=[minimal|full|expand|verbose|noverbose|remove-unused]    flattens all $ref's in spec prior to generation (default: minimal, verbose)

    Options for model generation:
      -m, --model-package=                                                        the package to save the models (default: models)
      -M, --model=                                                                specify a model to include in generation, repeat for multiple (defaults to all)
          --existing-models=                                                      use pre-generated models e.g. github.com/foobar/model
          --strict-additional-properties                                          disallow extra properties when additionalProperties is set to false
          --keep-spec-order                                                       keep schema properties order identical to spec file

    Options for operation generation:
      -O, --operation=                                                            specify an operation to include, repeat for multiple (defaults to all)
          --tags=                                                                 the tags to include, if not specified defaults to all
      -a, --api-package=                                                          the package to save the operations (default: operations)
          --with-enum-ci                                                          set all enumerations case-insensitive by default
          --skip-tag-packages                                                     skips the generation of tag-based operation packages, resulting in a flat generation
```

### Build a client

There is an example client provided at: https://github.com/go-swagger/go-swagger/tree/master/examples/todo-list/client

To generate a client:

```
swagger generate client -f [http-url|filepath] -A [application-name] [--principal [principal-name]]
```

If you want to debug what the client is sending and receiving you can set the environment value DEBUG to a non-empty
value.


Use a default client, which has an HTTP transport:

```go
import (
  "log"

  "github.com/myproject/client/operations"
  "github.com/go-openapi/strfmt"
  "github.com/go-openapi/spec"

  apiclient "github.com/myproject/client"
  httptransport "github.com/go-openapi/runtime/client"
)

func main() {

  // make the request to get all items
  resp, err := apiclient.Default.Operations.All(operations.AllParams{})
  if err != nil {
    log.Fatal(err)
  }
  fmt.Printf("%#v\n", resp.Payload)
}
```

The client runtime allows for a number of [configuration
options](https://godoc.org/github.com/go-openapi/runtime/client#Runtime) to be set.
To then use the client, and override the host, with a HTTP transport:

```go
import (
  "os"
  "log"

  "github.com/myproject/client/operations"
  "github.com/go-openapi/strfmt"
  "github.com/go-openapi/spec"

  apiclient "github.com/myproject/client"
  httptransport "github.com/go-openapi/runtime/client"
)

func main() {

  // create the transport
  transport := httptransport.New(os.Getenv("TODOLIST_HOST"), "", nil)

  // create the API client, with the transport
  client := apiclient.New(transport, strfmt.Default)

  // to override the host for the default client
  // apiclient.Default.SetTransport(transport)

  // make the request to get all items
  resp, err := client.Operations.All(operations.AllParams{})
  if err != nil {
    log.Fatal(err)
  }
  fmt.Printf("%#v\n", resp.Payload)
}
```

### Authentication

The client supports 3 authentication schemes:

* [Basic Auth](https://godoc.org/github.com/go-openapi/runtime/client#BasicAuth)
* [API key auth in header or query](https://godoc.org/github.com/go-openapi/runtime/client#APIKeyAuth)
* [Bearer token header for oauth2](https://godoc.org/github.com/go-openapi/runtime/client#BearerToken)

```go
import (
  "os"
  "log"

  "github.com/myproject/client/operations"
  "github.com/go-openapi/strfmt"
  "github.com/go-openapi/spec"

  apiclient "github.com/myproject/client"
  httptransport "github.com/go-openapi/runtime/client"
)

func main() {

  // create the API client
  client := apiclient.New(httptransport.New("", "", nil), strfmt.Default)

  // make the authenticated request to get all items
  bearerTokenAuth := httptransport.BearerToken(os.Getenv("API_ACCESS_TOKEN"))
  // basicAuth := httptransport.BasicAuth(os.Getenv("API_USER"), os.Getenv("API_PASSWORD"))
  // apiKeyQueryAuth := httptransport.APIKeyAuth("apiKey", "query", os.Getenv("API_KEY"))
  // apiKeyHeaderAuth := httptransport.APIKeyAuth("X-API-TOKEN", "header", os.Getenv("API_KEY"))
  resp, err := client.Operations.All(operations.AllParams{}, bearerTokenAuth)
  // resp, err := client.Operations.All(operations.AllParams{}, basicAuth)
  // resp, err := client.Operations.All(operations.AllParams{}, apiKeyQueryAuth)
  // resp, err := client.Operations.All(operations.AllParams{}, apiKeyHeaderAuth)
  if err != nil {
    log.Fatal(err)
  }
  fmt.Printf("%#v\n", resp.Payload)
}
```

### Consuming an XML API

In order to enable XML support, you need to set the command options `--default-consumes` or `--default-produces` to an XML mime type like `application/xml` when generating the client:

```
swagger generate client -f [http-url|filepath] -A [application-name] --default-consumes application/xml
```

This is necessary regardless of whether your swagger specification already specifies XML in the consumes and produces properties.

An example using the generated client with default Bearer authentication:

```go
import (
  "os"

  "github.com/go-openapi/strfmt"
  "github.com/go-openapi/runtime"

  apiclient "github.com/myproject/client"
  httptransport "github.com/go-openapi/runtime/client"
)

func main() {
  r := httptransport.New(apiclient.DefaultHost, apiclient.DefaultBasePath, apiclient.DefaultSchemes)
  r.DefaultAuthentication = httptransport.BearerToken(os.Getenv("API_ACCESS_TOKEN"))
  /*
  r.DefaultMediaType = runtime.XMLMime
  r.Consumers = map[string]runtime.Consumer{
    runtime.XMLMime: runtime.XMLConsumer(),
  }
  r.Producers = map[string]runtime.Producer{
    "application/xhtml+xml": runtime.XMLProducer(),
  }
  */
  client := apiclient.New(r, strfmt.Default)

  resp, err := client.Operations.MyGreatEndpoint()
  if err != nil {
    log.Fatal(err)
  }
  fmt.Printf("%#v\n", resp.Payload)
}
```

It can under certain circumstances be necessary to manually set the DefaultMediaType and the Consumers and Producers similar to the commented-out code above, particularly if you're using special mime types like `application/xhtml+xml`.
