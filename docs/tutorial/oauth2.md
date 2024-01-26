---
title: Authentication with Oauth2
date: 2023-01-01T01:01:01-08:00
draft: true
---
# Oauth2 Authentication sample: AccessCode workflow

The full code of this example is [here][example_code].

This example illustrates a complete OAuth2 handshake.

We want to implement a simple access control based on a user's Google account (i.e. OpenID).

Personas:

- the user logs in on its Google account, which returns an access token that we will use
  with our API. This mechanism follows the 'accessCode' OAuth2 workflow.


### Swagger specification

Given the following security definitions (in `swagger.yml` specification document):

```yaml
securityDefinitions:
  OauthSecurity:
    type: oauth2
    flow: accessCode
    authorizationUrl: 'https://accounts.google.com/o/oauth2/v2/auth'
    tokenUrl: 'https://www.googleapis.com/oauth2/v4/token'
    scopes:
      admin: Admin scope
      user: User scope
```

We specify the following security requirements:

- A default requirements for all endpoints: users need to be authenticated within the "user" scope by providing 
a OAuth token (e.g. Authentication: Bearer header or `access_token` query parameter).

```yaml
security:
  - OauthSecurity:
    - user
```

- Login and callback endpoints are not restricted: this is made explicit by overriding the default security requirement with an empty array.

```yaml
paths:
  /login:
    get:
      summary: login through oauth2 server
      security: []

...

  /auth/callback:
    get:
      summary: return access_token
      security: []
```

We need to specify a security principal in the model, to generate the server. Operations will be passed this principal as 
parameter upon successful authentication:

```yaml
definitions:
  ...
  principal:
    type: string
```

In this example, the principal (descriptor of an identity for our API) 
is just a string (i.e. the token itself).

### Generate the server 

```shell
swagger generate server -A oauthSample -P models.Principal -f ./swagger.yml
```

### Prepare the configuration

In `restapi/implementation.go` (this is not a generated file), we defined an
implementation for our workflow.

First, we need some extra packages to work with OAuth2, OpenID and HTTP redirections:

```go
import (
	oidc "github.com/coreos/go-oidc"            // Google OpenID client
	"context"
	"golang.org/x/oauth2"                       // OAuth2 client
)
```

```go
var (
    // state carries an internal token during the oauth2 workflow
    // we just need a non empty initial value
	state = "foobar" // Don't make this a global in production.

    // the credentials for this API (adapt values when registering API)
	clientID     = "" // <= enter registered API client ID here
	clientSecret = "" // <= enter registered API client secret here

    //  unused in this example: the signer of the delivered token
	issuer       = "https://accounts.google.com"

    // the Google login URL
	authURL      = "https://accounts.google.com/o/oauth2/v2/auth"

    // the Google OAuth2 resource provider which delivers access tokens
	tokenURL     = "https://www.googleapis.com/oauth2/v4/token"
	userInfoURL  = "https://www.googleapis.com/oauth2/v3/userinfo"

    // our endpoint to be called back by the redirected client
	callbackURL  = "http://127.0.0.1:12345/api/auth/callback"

    // the description of the OAuth2 flow
	endpoint = oauth2.Endpoint{
		AuthURL:  authURL,
		TokenURL: tokenURL,
	}

	config = oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     endpoint,
		RedirectURL:  callbackURL,
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}
)
```

### Configure the API in `restapi/configure_auth_sample.go`

```go
func configureAPI(api *operations.OauthSampleAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	api.OauthSecurityAuth = func(token string, scopes []string) (*models.Principal, error) {
        // This handler is called by the runtime whenever a route needs authentication 
        // against the 'OAuthSecurity' scheme.
        // It is passed a token extracted from the Authentication Bearer header, and 
        // the list of scopes mentioned by the spec for this route.

        // NOTE: in this simple implementation, we do not check scopes against  
        // the signed claims in the JWT token.
        // So whatever the required scope (passed a parameter by the runtime), 
        // this will succeed provided we get a valid token.

        // authenticated validates a JWT token at userInfoURL
		ok, err := authenticated(token)
		if err != nil {
			return nil, errors.New(401, "error authenticate")
		}
		if !ok {
			return nil, errors.New(401, "invalid token")
		}

        // returns the authenticated principal (here just filled in with its token)
		prin := models.Principal(token)
		return &prin, nil
	}

	api.GetAuthCallbackHandler = operations.GetAuthCallbackHandlerFunc(func(params operations.GetAuthCallbackParams) middleware.Responder {
        // implements the callback operation
		token, err := callback(params.HTTPRequest)
		if err != nil {
			return middleware.NotImplemented("operation .GetAuthCallback error")
		}
		log.Println("Token", token)
		return operations.NewGetAuthCallbackDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(token)})
	})

	api.GetLoginHandler = operations.GetLoginHandlerFunc(func(params operations.GetLoginParams) middleware.Responder {
        // implements the login operation
		login(params.HTTPRequest)
		return middleware.NotImplemented("operation .GetLogin has not yet been implemented")
	})

	api.CustomersCreateHandler = customers.CreateHandlerFunc(func(params customers.CreateParams, principal *models.Principal) middleware.Responder {
        // other API endpoint ...
		log.Println("hit customer API")
		return middleware.NotImplemented("operation customers.Create has not yet been implemented")
	})

	api.CustomersGetIDHandler = customers.GetIDHandlerFunc(func(params customers.GetIDParams, principal *models.Principal) middleware.Responder {
        // other API endpoint ...
		log.Println("hit customer API")
		return middleware.NotImplemented("operation customers.GetID has not yet been implemented")
	})

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}
```

We set the following implementation for authentication in `restapi/implementation.go` (**this is not generated code** and may be fully customized):

- Redirecting to the login page

```go
func login(r *http.Request) string {
	// implements the login with a redirection and an access token
	var accessToken string
	wG := r.Context().Value(ctxResponseWriter).(http.ResponseWriter)
	http.Redirect(wG, r, config.AuthCodeURL(state), http.StatusFound)
	return accessToken
}
```
- Retrieving the access token

```go
func callback(r *http.Request) (string, error) {
    // we expect the redirected client to call us back 
    // with 2 query params: state and code.
    // We use directly the Request params here, since we did not 
    // bother to document these parameters in the spec.

	if r.URL.Query().Get("state") != state {
		log.Println("state did not match")
		return "", fmt.Errorf("state did not match")
	}

	myClient := &http.Client{}

	parentContext := context.Background()
	ctx := oidc.ClientContext(parentContext, myClient)

	authCode := r.URL.Query().Get("code")
	log.Printf("Authorization code: %v\n", authCode)

    // Exchange converts an authorization code into a token.
    // Under the hood, the oauth2 client POST a request to do so
    // at tokenURL, then redirects...
	oauth2Token, err := config.Exchange(ctx, authCode)
	if err != nil {
		log.Println("failed to exchange token", err.Error())
		return "", fmt.Errorf("failed to exchange token")
	}

    // the authorization server's returned token
	log.Println("Raw token data:", oauth2Token)
	return oauth2Token.AccessToken, nil
}
```


- Validating the token

```go
func authenticated(token string) (bool, error) {
	// validates the token by sending a request at userInfoURL
	bearToken := "Bearer " + token
	req, err := http.NewRequest("GET", userInfoURL, nil)
	if err != nil {
		return false, fmt.Errorf("http request: %v", err)
	}

	req.Header.Add("Authorization", bearToken)
    
	cli := &http.Client{}
	resp, err := cli.Do(req)
	if err != nil {
		return false, fmt.Errorf("http request: %v", err)
	}
	defer resp.Body.Close()

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("fail to get response: %v", err)
	}
	if resp.StatusCode != 200 {
		return false, nil
	}
	return true, nil
}
```

### Register the callback URL 
Register your API at [google oauth2 server][google_credential], with 
an OAuth ID.  
Make sure that the callback URL is the same as set in the above code (``./restapi/configure_auth_sample.go``), e.g.:

```
http://127.0.0.1:12345/api/auth/callback
```

![Google api screenshot](https://github.com/go-swagger/go-swagger/blob/master/examples/oauth2/img/google-api.png)

>**NOTE:** you may specify a client ID for your API during the registration process.
>A password (the API client's secret) is then delivered.
>Those are the credentials of the API itself, not the end user.
>Put these values (client ID and client's secret) in the initial 
>var declarations in `implementation.go`.

### Run the server

```shell
go run ./cmd/oauth-sample-server/main.go --port 12345
```

### Login to get the access token

Get the access token through Google's oauth2 server. 

Open the browser and access the API login url on: 
http://127.0.0.1:12345/api/login,  which will direct you to the Google 
login page. 

Once you login with your google ID (e.g., your gmail account), the oauth2  
``access_token`` is returned and displayed on the browser.

### Exercise your authorizer

``TOKEN`` is obtained from the previous step.

Now we may use this token to access the other endpoints published by our API.

Let's try this with curl. Copy the received token and reuse it as shown below:
```shellsession
± ivan@avalon:~  
 »  curl -i  -H 'Authorization: Bearer TOKEN' http://127.0.0.1:12345/api/customers
```
```http
HTTP/1.1 501 Not Implemented
Content-Type: application/keyauth.api.v1+json
Date: Fri, 25 Nov 2016 19:14:14 GMT
Content-Length: 57

"operation customers.GetID has not yet been implemented"
```

Use an random string as the token:

```shellsession
± ivan@avalon:~  
 » curl -i  -H 'Authorization: Bearer RAMDOM_TOKEN' http://127.0.0.1:12345/api/customers
```
```http
HTTP/1.1 401 Unauthorized
Content-Type: application/keyauth.api.v1+json
Date: Fri, 25 Nov 2016 19:16:49 GMT
Content-Length: 47

{"code":401,"message":"unauthenticated for invalid credentials"}       
```

[google_credential]: https://console.cloud.google.com/apis/credentials/
[example_code]: https://github.com/go-swagger/go-swagger/blob/master/examples/oauth2/
