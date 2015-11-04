+++
categories = ["generate"]
date = "2015-10-23T22:11:54-07:00"
tags = ["generate"]
title = "Generate an API client"

+++

The toolkit has a command that will let you generate a client.

## Usage

To generate a client:

```
swagger generate client -f [http-url|filepath] -A [application-name] [--principal [principal-name]]
```

To then use the client with a HTTP transport:

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
  // load the swagger spec from URL or local file
  doc, err := spec.Load("https://raw.githubusercontent.com/go-swagger/go-swagger/master/examples/todo-list/swagger.yml")
  if err != nil {
    log.Fatal(err)
  }

  // create the API client
  client := apiclient.New(httptransport.New(doc), strfmt.Default)

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

* Basic Auth
* API key auth in header or query
* Bearer token header for oauth2

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
  // load the swagger spec from URL or local file
  doc, err := spec.Load("https://raw.githubusercontent.com/go-swagger/go-swagger/master/examples/todo-list/swagger.yml")
  if err != nil {
    log.Fatal(err)
  }

  // create the API client
  client := apiclient.New(httptransport.New(doc), strfmt.Default)

  // make the authenticated request to get all items
  bearerTokenAuth := httptransport.BearerToken(os.Getenv("API_ACCESS_TOKEN"))
  resp, err := client.Operations.All(operations.AllParams{}, bearerTokenAuth)
  if err != nil {
    log.Fatal(err)
  }
  fmt.Printf("%#v\n", resp.Payload)
}
```
