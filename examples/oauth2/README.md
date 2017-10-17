# Oauth2 Authentication sample

Generate the code with a security principal:

```shell
swagger generate server -A oauthSample -P models.Principal -f ./swagger.yml
```

### Edit the ./restapi/configure_auth_sample.go file

```go
var (
	state = "foobar" // Don't do this in production.

	clientID     = "YOUR CLIENT ID"
	clientSecret = "YOUR CLIENT SECRET"
	issuer       = "https://accounts.google.com"
	authURL      = "https://accounts.google.com/o/oauth2/v2/auth"
	tokenURL     = "https://www.googleapis.com/oauth2/v4/token"
	userInfoURL  = "https://www.googleapis.com/oauth2/v3/userinfo"
	callback_url = "http://127.0.0.1:12345/api/auth/callback" // must be registered with google API credential
)

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
		ok, err := authenticated(token)
		if err != nil {
			return nil, errors.New(401, "error authenticate")
		}
		if !ok {
			return nil, errors.New(401, "invalid token")
		}
		prin := models.Principal(token)
		return &prin, nil
	}
}
```

### Register the callback URL in [google oauth2 server][google_credential], e.g.,:

```
http://127.0.0.1:12345/api/auth/callback
```
Make sure that the callback URL is the same as set in the above code (``./restapi/configure_auth_sample.go``)

### Run the server:

```shell
go run ./cmd/auth-sample-server/main.go --port 12345
```

### Login to get the access token

Get the access token through google's oauth2 server. Open the browser and access the url of http://127.0.0.1:12345/api/login, which will direct you to the google login page. Once you login with your google ID (e.g., your gmail account), the oauth2 ``access_token`` is returned and displayed on the browser.

### Exercise auth:

``TOKEN`` is obtained from the previous step.

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

{"code":401,"message":"invalid token"}       
```

[google_credential]: https://console.cloud.google.com/apis/credentials/
