---
title: About documentating your API
date: 2023-01-01T01:01:01-08:00
draft: true
---
<!-- Questions about the serve UI use-case -->
### Serving swagger-ui with the API Server
Update: You can visit $API_BASE/docs and it should show the UI, as mentioned in [#comment](https://github.com/go-swagger/go-swagger/issues/2401#issuecomment-688962519) 

_Use-Case_: I was trying to serve swagger-ui from the generated API Server and
I didn't find a straightforward enough way in the docs,
so I've created my own swagger-ui middleware:

```golang
func setupGlobalMiddleware(handler http.Handler) http.Handler {
    return uiMiddleware(handler)
}

func uiMiddleware(handler http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Shortcut helpers for swagger-ui
        if r.URL.Path == "/swagger-ui" || r.URL.Path == "/api/help" {
            http.Redirect(w, r, "/swagger-ui/", http.StatusFound)
            return
        }
        // Serving ./swagger-ui/
        if strings.Index(r.URL.Path, "/swagger-ui/") == 0 {
            http.StripPrefix("/swagger-ui/", http.FileServer(http.Dir("swagger-ui"))).ServeHTTP(w, r)
            return
        }
        handler.ServeHTTP(w, r)
    })
}
```

*Can this be achieved any easier?*

**Answer**: I think that's the way to do it.

At some point I included an embedded version of the swagger-ui in this toolkit but it became annoying to keep up to date
and severely bloated the size of my binary.

*What do you say if we add swagger-ui as git submodule, include this middleware in your binary and update the docs?*

I'm reluctant to do so at this point in time because a git submodule break go-gettability.

>I've had it included at one point but it's so much of a moving target that it would always be outdated.
>On top of it is a lot of javascript and html files and people haven't been over the moon when go-swagger gets
>vendored and they see all of that.

Originally from issue [#370](https://github.com/go-swagger/go-swagger/issues/370).

See also: How to serve Swagger UI from a preexisting web app? [#1029](https://github.com/go-swagger/go-swagger/issues/1029).

### How to serve Swagger UI from a preexisting web app?
_Use-Case_: Does go-swagger provide an `http.HandlerFunc` or other easy method for serving Swagger UI from a preexisting web app? 
I want my web app to expose `/swagger-ui`, without using code generation, and without hosting a separate server.

**Answer**: there are a few ways you can serve a UI.

Use the middleware provided in the go-openapi/runtime package: https://github.com/go-openapi/runtime/blob/master/middleware/redoc.go

Originally from issues [#1029](https://github.com/go-swagger/go-swagger/issues/1029) and [#976](https://github.com/go-swagger/go-swagger/issues/976)

### How to use swagger-ui cors?

**Answer**: you can add a cors middleware.

Like: https://github.com/rs/cors

[Documentation on how to customize middleware](reference/middleware.md)

Working example (in `configure_name.go`):

```golang
import "github.com/rs/cors"

func setupGlobalMiddleware(handler http.Handler) http.Handler {
    handleCORS := cors.Default().Handler
    return handleCORS(handler)
}
```

Originally from issue [#481](https://github.com/go-swagger/go-swagger/issues/481).

### How to serve my UI files?
_Use-Case_: I generated server code using go-swagger with my swagger.yaml file like below.
```bash
$ swagger generate server --exclude-main -A myapp -t gen -f ./swagger.yaml
```
And I want to add new handler to serve my own UI files. 
In this case, is middleware only solution to serve UI files? Or can I add new handler to serve files without middleware?
```go
// Handler example which I want to add to swagger server
func pastaWorkspacePage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/workspace.html")
}
```

I solved the problem using middleware.
```go
func FileServerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api") {
			next.ServeHTTP(w, r)
		} else {
			http.FileServer(http.Dir("./static")).ServeHTTP(w, r)
		}
	})
}
```
But I'm not sure it is the best solution.

**Hint:** more info on using middlewares is found here: https://goswagger.io/use/middleware.html
That page also contains a link to a good explanation on how to create net/http middlewares.

> An implementation example is provided by the go-swagger serve UI command. It constructs a server with a redoc middleware:
> https://github.com/go-swagger/go-swagger/blob/f552963ac0dfdec0450f6749aeeeeb2d31cd5544/cmd/swagger/commands/serve.go#L35.

Besides, every swagger generated server comes with the redoc UI baked in at `/{basepath}/docs`

Originally from issue [#1375](https://github.com/go-swagger/go-swagger/issues/1375).

