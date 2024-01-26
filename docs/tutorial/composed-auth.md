---
title: Composed security requirements
date: 2023-01-01T01:01:01-08:00
draft: true
---
# Composed Security Requirements

The full code of this example is [here][example_code].

This sample API demonstrates how to compose several authentication schemes 
and configure complex security requirements for your operations.

In this example, we mix security requirements with AND and OR constraints.

This API apes a very simple market place with customers and resellers of items.

Personas:

  - as a first time user, I want to see all items on sales
  - as a registered customer, I want to post orders for items and 
    consult my past orders
  - as a registered reseller, I want to see all pending orders on the items 
    I am selling on the market place
  - as a reseller managing my own inventories, I want to post replenishment orders for the items I provide
  - as a registered user, I want to consult my personal account infos

The playground situation we defined is as follows:

  - every known user is authenticated using a basic username:password pair
  - resellers are authenticated using API keys - we leave them the option to authenticate using a header or a query param
  - any registered user (customer or reseller) will add a signed JWT to access more API endpoints

Authentication with tokens allows us to inspect the signed claims in this token.

Obviously, there are several ways to achieve the same result. We just wanted to demonstrate here how
security requirements may be composed out of several schemes, and use API authorizers.

> Note that we used the "OAuth2" declaration here but don't actually follow an OAuth2 workflow:
> our intend here is to be able to extract scopes from the claims passed in a JWT token 
> (the only way to manipulate scoped authorizers with Swagger 2.0 is to declare them with type `oauth2`).


### Caveats

1. There should be at most one Authorization header: mixing "Authorization Basic" and "Authorization Bearer" won't work well
2. There should be at most one scoped authentication scheme: if we define several such authorizers they would all use the same bearer token
3. The "OAuth2" type supports other methods than the "Authorization: Bearer" header: the token may be passed
   using the `access_token` query param (or urlEncoded form value)
4. Unfortunately, Swagger 2.0 only supports "OAuth2" as a scoped method
5. There is one single principal and several methods to define it. Getting to these different intermediary principals requires some 
   interaction with the http request's context (e.g. using `middleware.SecurityPrincipalFrom(req)`). This is not demonstrated here for now.

### Prerequisites

`golang-jwt/jwt` ships with a nice JWT CLI utility. Although not required, you might want to install it and 
play with your own tokens:

- `go install github.com/golang-jwt/jwt/cmd/jwt`

### Swagger specification

We defined the following security schemes (in `swagger.yml` specification document):

```yaml
securityDefinitions:
  isRegistered:
    # This scheme uses the header: "Authorization: Basic {base64 encoded string defined by username:password}"
    # Scopes are not supported with this type of authorization.
    type: basic
  isReseller:
    # This scheme uses the header: "X-Custom-Key: {base64 encoded string}"
    # Scopes are not supported with this type of authorization.
    type: apiKey
    in: header
    name: X-Custom-Key
  isResellerQuery:
    # This scheme uses the query parameter "CustomKeyAsQuery"
    # Scopes are not supported with this type of authorization.
    type: apiKey
    in: query
    name: CustomKeyAsQuery
  hasRole:
    # This scheme uses the header: "Authorization: Bearer {base64 encoded string representing a JWT}"
    # Alternatively, the query param: "access_token" may be used.
    #
    # In our scenario, we must use the query param version in order to avoid 
    # passing several headers with key 'Authorization'
    type: oauth2
    # The flow and URLs in spec are for documentary purpose: go-swagger does not implement OAuth workflows
    flow: accessCode
    authorizationUrl: 'https://dummy.oauth.net/auth'
    tokenUrl: 'https://dumy.oauth.net/token'
    # Required scopes are passed by the runtime to the authorizer
    scopes:
      customer: scope of registered customers
      inventoryManager: scope of resellers acting as inventory managers
```

We specify the following security requirements:

- A default requirements for all endpoints: so by default, all endpoints use the Basic auth.

```yaml
security:
  - isRegistered: []
```

- Some endpoints are not restricted at all: this is made explicit by overriding the default security requirement with an empty array.

```yaml
paths:
  /items:
    get:
      summary: items on sale
      operationId: GetItems
      description: |
        Everybody should be able to access this operation
      security: []
...
```

- We created endpoints with various compositions of our 3 security schemes

Example: `isRegistered` **AND** `hasRole[ customer ]`
```yaml
  /order/{orderID}:
    get:
      summary: retrieves an order
      operationId: GetOrder
      description: |
        Only registered customers should be able to retrieve orders
      security: 
        - isRegistered: []
          hasRole: [ customer ]  
...
```

Example: (`isRegistered` **AND** `hasRole[ customer ]`) **OR** (`isReseller` **AND** `hasRole[ inventoryManager ]`) **OR** (`isResellerQuery` **AND** `hasRole[ inventoryManager ]`)

```yaml
  /order/add:
    post:
      summary: post a new order
      operationId: AddOrder
      description: |
        Registered customers should be able to add purchase orders.
        Registered inventory managers should be able to add replenishment orders.

      security:
        - isRegistered: []
          hasRole: [ customer ]  
        - isReseller: []
          hasRole: [ inventoryManager ]  
        - isResellerQuery: []
          hasRole: [ inventoryManager ]  
...
```

Example: isReseller **OR** isResellerQuery

This one allows to pass an API key either by header or by query param.

```yaml
  /orders/{itemID}:
    get:
      summary: retrieves all orders for an item
      operationId: GetOrdersForItem
      description: |
        Only registered resellers should be able to search orders for an item
      security:
        - isReseller: []
        - isResellerQuery: []
...
```
We need to specify a security principal in the model and generate the server with this. Operations will be passed this principal as 
parameter upon successful authentication.

When using the scoped authentication ("oauth2"), our custom authorizer with pass all claimed roles that match the security requirement in the principal.

```yaml
definitions:
  ...
  principal:
    type: object 
    properties: 
      name: 
        type: string
      roles:
        type: array 
        items: 
          type: string
```

### Generate the server 

```shell
swagger generate server -A multi-auth-example -P models.Principal -f ./swagger.yml
```

Files `restapi/configure_multi_auth_example.go` and `auth/authorizers.go` are not generated.

### Testing configuration

#### Test tokens and keys
In `./tokens`, we provided with some ready made tokens. If you have installed the `jwt` CLI, 
you can play around an build some different claims as JWT (see the `make-tokens.sh` script for usage).

> **NOTE:** tokens need a pair of public / private keys (for the signer and the verifier). We generated these keys 
> for testing purpose in the `keys` directory (RSA256 keys).

Our JWT defines "roles" as custom claim (in `auth/authorizers.go`): this means the signer of the token acknowledges the 
holder of the token to be enabled for these.

```go
// roleClaims describes the format of our JWT token's claims
type roleClaims struct {
	Roles []string `json:"roles"`
	jwt.StandardClaims
}
```

#### Configure the API with custom authorizers

In `configure_multi_auth_example.go` we have set up our custom authorizers:

```go
func configureAPI(api *operations.MultiAuthExampleAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

    ...
	logger := logging.MustGetLogger("api")

	api.Logger = logger.Infof

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	// Applies when the "Authorization: Basic" header is set with the Basic scheme
	api.IsRegisteredAuth = func(user string, pass string) (*models.Principal, error) {
		// The header: Authorization: Basic {base64 string} has already been decoded by the runtime as a username:password pair
		api.Logger("IsRegisteredAuth handler called")
		return auth.IsRegistered(user, pass)
	}

	// Applies when the "Authorization: Bearer" header or the "access_token" query is set
	api.HasRoleAuth = func(token string, scopes []string) (*models.Principal, error) {
		// The header: Authorization: Bearer {base64 string} (or ?access_token={base 64 string} param) has already
		// been decoded by the runtime as a token
		api.Logger("HasRoleAuth handler called")
		return auth.HasRole(token, scopes)
	}

	// Applies when the "CustomKeyAsQuery" query is set
	api.IsResellerQueryAuth = func(token string) (*models.Principal, error) {
		api.Logger("ResellerQueryAuth handler called")
		return auth.IsReseller(token)
	}

	// Applies when the "X-Custom-Key" header is set
	api.IsResellerAuth = func(token string) (*models.Principal, error) {
		api.Logger("IsResellerAuth handler called")
		return auth.IsReseller(token)
	}
    ...
```

These authorizers are implemented in `auth/authorizers.go`.

Here is the basic one:
```go 
// IsRegistered determines if the user is properly registered,
// i.e if a valid username:password pair has been provided
func IsRegistered(user, pass string) (*models.Principal, error) {
	logger.Debugf("Credentials: %q:%q", user, pass)
	if password, ok := userDb[user]; ok {
		if pass == password {
			return &models.Principal{
				Name: user,
			}, nil
		}
	}
	logger.Debug("Bad credentials")
	return nil, errors.New(401, "Unauthorized: not a registered user")
}
```

We did not set up actual operations: they are mere debug loggers, returning a "Not implemented" error.
We log on the serve console how the principal is passed to the operation.

```go 
	api.AddOrderHandler = operations.AddOrderHandlerFunc(func(params operations.AddOrderParams, principal *models.Principal) middleware.Responder {
		logger.Warningf("AddOrder called with params: %s, and principal: %s", spew.Sdump(params.Order), spew.Sdump(principal))
		return middleware.NotImplemented("operation .AddOrder has not yet been implemented")
	})
```

### Run the server

```shell
go run ./cmd/multi-auth-example-server/main.go --port 43016
```

### Exercise your authorizers

There is a little exercising utility script: `exerciser.sh`. 
This script pushes a sequence of curl requests. You may customize it to your liking and further exercise the API.

Authorizations actions and operations are logged on the server console.

Example:
```shell
curl \
  --verbose \
  --get \
  --header "X-Custom-Key: `cat tokens/token-apikey-reseller.jwt`" \
  "http://localhost:43016/api/orders/myItem"

*   Trying 127.0.0.1...
* Connected to localhost (127.0.0.1) port 43016 (#0)
> GET /api/orders/myItem HTTP/1.1
> Host: localhost:43016
> User-Agent: curl/7.47.0
> Accept: */*
> X-Custom-Key: eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJleGFtcGxlLmNvbSIsImp0aSI6ImZyZWQiLCJyb2xlcyI6WyJyZXNlbGxlciJdfQ.gvI8J3vNaXHOFCLF46Db-9tIf9Y_4xhN5ZKd0-z7AMRyrElVjG46epT_ld5p2YRyMQNXn4LPESGWkxdJVsnZPmXYkKHBUeDSb0hj523Eue-Ayf-pwMIN4DpvcAToU0XY8srlrlLIUWINn1tOPZGtprksxMfh7TkXcWHKkI8Q0P8-3JBTkoq4HBL1DzcAwYh4EGcFcgoXMUuR_TfE3SIOjUUE5Zs3c6UswPpvZv82jAGhFIs6uJI-73BvEZ084OmI0gCJNfHEms-79nDkqh5DHf6biQsABSdBfjDLNo24nkOhlOr7IOY0LSGws9xeaM8gY58lYN3Evpia642OUxwYI55fZzku4VGm7Ia2-uK_tD8AoNLquufmPP9ROAY63cZF0wnlw_6IM1gP4LQknVWb4gcdC0j7dk4SG01u4j9OhCXy2SLqx_SI9ZM5kfgAq6kGzQULRGmBbkSCFQfEzPn5v2WzAl_XmQ7uF5KJqgjDQlbamugXlz69w5eUECRpJGNjlGxb11Q-LBKgJ9An_nOSp0p3TfIIQOXTTz5W5CzC0DRsslN50l-6z0xTwtqiy47u8JhZk-073YkDWT_NS3MEAkgb48fFwLZIlnH5bAM5kZbZ4B7fql1j_G6UGY1tcmMXhfKP6ePE0PtMPSE1U7sF-nHPE7spwD5_56BjdBQf4pM
>
< HTTP/1.1 501 Not Implemented
< Content-Type: application/json
< Date: Tue, 17 Apr 2018 17:55:45 GMT
< Content-Length: 59
<
"operation .GetOrdersForItem has not yet been implemented"
```

Another example:
```shell
basic=`echo "ivan:terrible"|tr -d '\n'|base64 -i`
curl \
  --verbose \
  -i \
  -X POST \
  --data '{"orderID": "myorder", "orderLines": [{"quantity": 10, "purchasedItem": "myItem"}]}' \
  --header "Content-Type: application/json" \
  --header "Authorization: Basic ${basic}" \
  --header "X-Custom-Key: `cat tokens/token-apikey-reseller.jwt`" \
  "http://localhost:43016/api/order/add?access_token=`cat tokens/token-bearer-inventory-manager.jwt`"

*   Trying 127.0.0.1...
* Connected to localhost (127.0.0.1) port 43016 (#0)
> POST /api/order/add?access_token=eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJleGFtcGxlLmNvbSIsImp0aSI6ImZyZWQiLCJyb2xlcyI6WyJpbnZlbnRvcnlNYW5hZ2VyIl19.LaWEWhrDdrcwIatu7dVk-jhmzI4RlGgV0OFa1kLY2C6qMKQDibIActx1SVYuCxLLOafycbXlzCBGddoaHUHwjtuoOgftog2jHAR8-JJYyzSHCuz68cEngRtvY-MKgGApJqYInRhhdbV-DUiadPJjryxf9NNbyrdHjSMhSJOVDQp9Rj9VEGoK0zoufKOy_YrQEfcWl8OHHS17H7CI5_L44MsyC2Z6U3-HGo2eBdoIIVe5dUINA_PZ-U6netGOYuQ_T4GJ9IYjdkUOOpd_LJrFhCCE7vs4QnxVgnhSBzu5mL_ygJMyoA0yEvP01wSBxZIgFqJtiNXB5LYf-O_T2UuPLLfHiLRkmAJMDUxH_WYkat_qBILpzuRLFjGHrM3beouwVoHvdd7P3EYkr3vhElHs3GDKC9VVnqFL2a5hQOHdgJ4w6CVesNZGZkL-ZZKNzqeTzWch2WQcrYO26sG_XsS8Y1_mZ1UkMc8aJHGaQwnDt7SXUBTgvn57XuH_Ny6S2NfHJ6TX_rPanDlMQASOLR_yAYMJsqZjZlh9qXLR1bWjv7SQ4uuKb7YTf_gcYA3axEmYSXWZZA34gG3Q6eCDjPhEdN-RPj7C8zPiIa7VUG5ay4yp8A_hTtjsWKjvC7Kh3jZpyaF3M2QcBWkwQyFxljrMyxyHdAujkh7M-Y70O9U_YLU HTTP/1.1
> Host: localhost:43016
> User-Agent: curl/7.47.0
> Accept: */*
> Content-Type: application/json
> Authorization: Basic aXZhbjp0ZXJyaWJsZQ==
> X-Custom-Key: eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJleGFtcGxlLmNvbSIsImp0aSI6ImZyZWQiLCJyb2xlcyI6WyJyZXNlbGxlciJdfQ.gvI8J3vNaXHOFCLF46Db-9tIf9Y_4xhN5ZKd0-z7AMRyrElVjG46epT_ld5p2YRyMQNXn4LPESGWkxdJVsnZPmXYkKHBUeDSb0hj523Eue-Ayf-pwMIN4DpvcAToU0XY8srlrlLIUWINn1tOPZGtprksxMfh7TkXcWHKkI8Q0P8-3JBTkoq4HBL1DzcAwYh4EGcFcgoXMUuR_TfE3SIOjUUE5Zs3c6UswPpvZv82jAGhFIs6uJI-73BvEZ084OmI0gCJNfHEms-79nDkqh5DHf6biQsABSdBfjDLNo24nkOhlOr7IOY0LSGws9xeaM8gY58lYN3Evpia642OUxwYI55fZzku4VGm7Ia2-uK_tD8AoNLquufmPP9ROAY63cZF0wnlw_6IM1gP4LQknVWb4gcdC0j7dk4SG01u4j9OhCXy2SLqx_SI9ZM5kfgAq6kGzQULRGmBbkSCFQfEzPn5v2WzAl_XmQ7uF5KJqgjDQlbamugXlz69w5eUECRpJGNjlGxb11Q-LBKgJ9An_nOSp0p3TfIIQOXTTz5W5CzC0DRsslN50l-6z0xTwtqiy47u8JhZk-073YkDWT_NS3MEAkgb48fFwLZIlnH5bAM5kZbZ4B7fql1j_G6UGY1tcmMXhfKP6ePE0PtMPSE1U7sF-nHPE7spwD5_56BjdBQf4pM
> Content-Length: 83
>
* upload completely sent off: 83 out of 83 bytes
< HTTP/1.1 501 Not Implemented
HTTP/1.1 501 Not Implemented
< Content-Type: application/json
Content-Type: application/json
< Date: Tue, 17 Apr 2018 17:55:45 GMT
Date: Tue, 17 Apr 2018 17:55:45 GMT
< Content-Length: 51
Content-Length: 51
<
"operation .AddOrder has not yet been implemented"
```

[example_code]: https://github.com/go-swagger/go-swagger/blob/master/examples/composed-auth/
