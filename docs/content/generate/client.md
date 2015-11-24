+++
categories = ["generate", "client"]
date = "2015-10-23T22:11:54-07:00"
title = "Generate an API client"
series = ["home"]
weight = 2
+++

The toolkit has a command that will let you generate a client.

<!--more-->

## Usage

There is an example client in 

To generate a client:

```
swagger generate client -f [http-url|filepath] -A [application-name] [--principal [principal-name]]
```

Use a default client, which has an HTTP transport:

```go
import (
  "log"

  "github.com/myproject/client/operations"
  "github.com/go-swagger/go-swagger/strfmt"
  "github.com/go-swagger/go-swagger/spec"

  apiclient "github.com/myproject/client"
  httptransport "github.com/go-swagger/go-swagger/httpkit/client"
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
options](https://godoc.org/github.com/go-swagger/go-swagger/httpkit/client#Runtime) to be set.  
To then use the client, and override the host, with a HTTP transport:

```go
import (
  "os"
  "log"

  "github.com/myproject/client/operations"
  "github.com/go-swagger/go-swagger/strfmt"
  "github.com/go-swagger/go-swagger/spec"

  apiclient "github.com/myproject/client"
  httptransport "github.com/go-swagger/go-swagger/httpkit/client"
)

func main() {
  // load the swagger spec from URL or local file
  doc, err := spec.Load("https://raw.githubusercontent.com/go-swagger/go-swagger/master/examples/todo-list/swagger.yml")
  if err != nil {
    log.Fatal(err)
  }

  // create the transport
  transport := httptransport.New(doc)
  // configure the host
  if os.Getenv("TODOLIST_HOST") != "" {
    transport.Host = os.Getenv("TODOLIST_HOST")
  }

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

* [Basic Auth](https://godoc.org/github.com/go-swagger/go-swagger/httpkit/client#BasicAuth)
* [API key auth in header or query](https://godoc.org/github.com/go-swagger/go-swagger/httpkit/client#APIKeyAuth)
* [Bearer token header for oauth2](https://godoc.org/github.com/go-swagger/go-swagger/httpkit/client#BearerToken)

```go
import (
  "os"
  "log"

  "github.com/myproject/client/operations"
  "github.com/go-swagger/go-swagger/strfmt"
  "github.com/go-swagger/go-swagger/spec"

  apiclient "github.com/myproject/client"
  httptransport "github.com/go-swagger/go-swagger/httpkit/client"
)

func main() {
  // load the swagger spec from URL or local file
  doc, err := spec.Load("./swagger.yml")
  if err != nil {
    log.Fatal(err)
  }

  // create the API client
  client := apiclient.New(httptransport.New(doc), strfmt.Default)

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
