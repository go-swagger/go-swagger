---
title: Authentication with API key
date: 2023-01-01T01:01:01-08:00
draft: true
---
# Authentication sample

The full code of this example is [here][example_code].

Define the following security scheme (in `swagger.yml` specification document):

```yaml
securityDefinitions:
  key:
    type: apiKey
    in: header
    name: x-token
```

Specify the following security requirements for all endpoints: so by default,
all endpoints use the API key auth.

```yaml
security:
  - key: []
```

Add security principal model definition:

```yaml
definitions:

...

  principal:
    type: string
```

Generate the code with a security principal:

```shell
swagger generate server -A AuthSample -P models.Principal -f ./swagger.yml
```

Edit the ./restapi/configure_auth_sample.go file

```go
func configureAPI(api *operations.AuthSampleAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	// Applies when the "x-token" header is set
	api.KeyAuth = func(token string) (*models.Principal, error) {
		if token == "abcdefuvwxyz" {
			prin := models.Principal(token)
			return &prin, nil
		}
		api.Logger("Access attempt with incorrect api key auth: %s", token)
		return nil, errors.New(401, "incorrect api key auth")
	}

	api.CustomersCreateHandler = customers.CreateHandlerFunc(func(params customers.CreateParams, principal *models.Principal) middleware.Responder {
		return middleware.NotImplemented("operation customers.Create has not yet been implemented")
	})
	api.CustomersGetIDHandler = customers.GetIDHandlerFunc(func(params customers.GetIDParams, principal *models.Principal) middleware.Responder {
		return middleware.NotImplemented("operation customers.GetID has not yet been implemented")
	})

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}
```

Run the server:

```shell
go run ./cmd/auth-sample-server/main.go --port 35307
```

Exercise auth:

```shell
± ivan@avalon:~  
 » curl -i -H 'Content-Type: application/keyauth.api.v1+json' -H 'X-Token: abcdefuvwxyz' http://127.0.0.1:35307/api/customers
```
```http
HTTP/1.1 501 Not Implemented
Content-Type: application/keyauth.api.v1+json
Date: Fri, 25 Nov 2016 19:14:14 GMT
Content-Length: 57

"operation customers.GetID has not yet been implemented"
```
```shell
± ivan@avalon:~  
 » curl -i -H 'Content-Type: application/keyauth.api.v1+json' -H 'X-Token: abcdefu' http://127.0.0.1:35307/api/customers
```
```http
HTTP/1.1 401 Unauthorized
Content-Type: application/keyauth.api.v1+json
Date: Fri, 25 Nov 2016 19:16:49 GMT
Content-Length: 47

{"code":401,"message":"incorrect api key auth"}
```

[example_code]: https://github.com/go-swagger/go-swagger/tree/master/examples/authentication
