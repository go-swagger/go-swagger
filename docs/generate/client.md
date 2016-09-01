# Generate an API client

The toolkit has a command that will let you generate a client.

<!--more-->

## Usage

There is an example client in https://github.com/go-swagger/go-swagger/tree/master/examples/todo-list/client

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
